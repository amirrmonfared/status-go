package wallet

import (
	"database/sql"

	"github.com/ethereum/go-ethereum/event"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/p2p"
	gethrpc "github.com/ethereum/go-ethereum/rpc"

	"github.com/status-im/status-go/account"
	"github.com/status-im/status-go/multiaccounts/accounts"
	"github.com/status-im/status-go/params"
	"github.com/status-im/status-go/rpc"
	"github.com/status-im/status-go/services/ens"
	"github.com/status-im/status-go/services/stickers"
	"github.com/status-im/status-go/services/wallet/token"
	"github.com/status-im/status-go/services/wallet/transfer"
	"github.com/status-im/status-go/services/wallet/walletevent"
	"github.com/status-im/status-go/transactions"
)

// NewService initializes service instance.
func NewService(
	db *sql.DB,
	accountsDB *accounts.Database,
	rpcClient *rpc.Client,
	accountFeed *event.Feed,
	openseaAPIKey string,
	gethManager *account.GethManager,
	transactor *transactions.Transactor,
	config *params.NodeConfig,
	ens *ens.Service,
	stickers *stickers.Service,
) *Service {
	cryptoOnRampManager := NewCryptoOnRampManager(&CryptoOnRampOptions{
		dataSourceType: DataSourceStatic,
	})
	walletFeed := &event.Feed{}
	signals := &walletevent.SignalsTransmitter{
		Publisher: walletFeed,
	}
	tokenManager := token.NewTokenManager(db, rpcClient, rpcClient.NetworkManager)
	savedAddressesManager := &SavedAddressesManager{db: db}
	transactionManager := &TransactionManager{db: db, transactor: transactor, gethManager: gethManager, config: config, accountsDB: accountsDB}
	transferController := transfer.NewTransferController(db, rpcClient, accountFeed, walletFeed)
	reader := NewReader(rpcClient, tokenManager, accountsDB, walletFeed)
	return &Service{
		db:                    db,
		accountsDB:            accountsDB,
		rpcClient:             rpcClient,
		tokenManager:          tokenManager,
		savedAddressesManager: savedAddressesManager,
		transactionManager:    transactionManager,
		transferController:    transferController,
		cryptoOnRampManager:   cryptoOnRampManager,
		openseaAPIKey:         openseaAPIKey,
		feesManager:           &FeeManager{rpcClient},
		gethManager:           gethManager,
		transactor:            transactor,
		ens:                   ens,
		stickers:              stickers,
		feed:                  accountFeed,
		signals:               signals,
		reader:                reader,
	}
}

// Service is a wallet service.
type Service struct {
	db                    *sql.DB
	accountsDB            *accounts.Database
	rpcClient             *rpc.Client
	savedAddressesManager *SavedAddressesManager
	tokenManager          *token.Manager
	transactionManager    *TransactionManager
	cryptoOnRampManager   *CryptoOnRampManager
	transferController    *transfer.Controller
	feesManager           *FeeManager
	started               bool
	openseaAPIKey         string
	gethManager           *account.GethManager
	transactor            *transactions.Transactor
	ens                   *ens.Service
	stickers              *stickers.Service
	feed                  *event.Feed
	signals               *walletevent.SignalsTransmitter
	reader                *Reader
}

// Start signals transmitter.
func (s *Service) Start() error {
	s.transferController.Start()
	err := s.signals.Start()
	s.started = true
	return err
}

// GetFeed returns signals feed.
func (s *Service) GetFeed() *event.Feed {
	return s.transferController.TransferFeed
}

// Stop reactor and close db.
func (s *Service) Stop() error {
	log.Info("wallet will be stopped")
	s.signals.Stop()
	s.transferController.Stop()
	s.reader.Stop()
	s.started = false
	log.Info("wallet stopped")
	return nil
}

// APIs returns list of available RPC APIs.
func (s *Service) APIs() []gethrpc.API {
	return []gethrpc.API{
		{
			Namespace: "wallet",
			Version:   "0.1.0",
			Service:   NewAPI(s),
			Public:    true,
		},
	}
}

// Protocols returns list of p2p protocols.
func (s *Service) Protocols() []p2p.Protocol {
	return nil
}

func (s *Service) IsStarted() bool {
	return s.started
}
