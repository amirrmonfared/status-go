package bridge

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/status-im/status-go/account"
	"github.com/status-im/status-go/contracts/celer"
	"github.com/status-im/status-go/eth-node/types"
	"github.com/status-im/status-go/rpc"

	"github.com/status-im/status-go/params"
	"github.com/status-im/status-go/services/wallet/bridge/cbridge"
	"github.com/status-im/status-go/services/wallet/token"
	"github.com/status-im/status-go/transactions"
)

const baseURL = "https://cbridge-prod2.celer.app"
const testBaseURL = "https://cbridge-v2-test.celer.network"

type CBridgeTxArgs struct {
	transactions.SendTxArgs
	ChainID   uint64         `json:"chainId"`
	Symbol    string         `json:"symbol"`
	Recipient common.Address `json:"recipient"`
	Amount    *hexutil.Big   `json:"amount"`
}

type CBridge struct {
	rpcClient          *rpc.Client
	transactor         *transactions.Transactor
	tokenManager       *token.Manager
	prodTransferConfig *cbridge.GetTransferConfigsResponse
	testTransferConfig *cbridge.GetTransferConfigsResponse
}

func NewCbridge(rpcClient *rpc.Client, transactor *transactions.Transactor, tokenManager *token.Manager) *CBridge {
	return &CBridge{
		rpcClient:    rpcClient,
		transactor:   transactor,
		tokenManager: tokenManager,
	}
}

func (s *CBridge) Name() string {
	return "CBridge"
}

func (s *CBridge) estimateAmt(from, to *params.Network, amountIn *big.Int, symbol string) (*cbridge.EstimateAmtResponse, error) {
	client := &http.Client{
		Timeout: time.Second * 5,
	}
	base := baseURL
	if from.IsTest {
		base = testBaseURL
	}
	url := fmt.Sprintf(
		"%s/v2/estimateAmt?src_chain_id=%d&dst_chain_id=%d&token_symbol=%s&amt=%s&usr_addr=0xaa47c83316edc05cf9ff7136296b026c5de7eccd&slippage_tolerance=500",
		base,
		from.ChainID,
		to.ChainID,
		symbol,
		amountIn.String(),
	)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("failed to close cbridge request body", "err", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res cbridge.EstimateAmtResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *CBridge) getTransferConfig(isTest bool) (*cbridge.GetTransferConfigsResponse, error) {
	if !isTest && s.prodTransferConfig != nil {
		return s.prodTransferConfig, nil
	}

	if isTest && s.testTransferConfig != nil {
		return s.testTransferConfig, nil
	}

	client := &http.Client{
		Timeout: time.Second * 5,
	}
	base := baseURL
	if isTest {
		base = testBaseURL
	}
	url := fmt.Sprintf("%s/v2/getTransferConfigs", base)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			fmt.Println("failed to close cbridge request body", "err", err)
		}
	}()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res cbridge.GetTransferConfigsResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	if isTest {
		s.testTransferConfig = &res
	} else {
		s.prodTransferConfig = &res
	}
	return &res, nil
}

func (s *CBridge) Can(from, to *params.Network, token *token.Token, balance *big.Int) (bool, error) {
	if from.ChainID == to.ChainID {
		return false, nil
	}

	transferConfig, err := s.getTransferConfig(from.IsTest)
	if err != nil {
		return false, err
	}
	if transferConfig.Err != nil {
		return false, errors.New(transferConfig.Err.Msg)
	}

	var fromAvailable *cbridge.Chain
	var toAvailable *cbridge.Chain
	for _, chain := range transferConfig.Chains {
		if uint64(chain.GetId()) == from.ChainID {
			fromAvailable = chain
		}

		if uint64(chain.GetId()) == to.ChainID {
			toAvailable = chain
		}
	}

	if fromAvailable == nil || toAvailable == nil {
		return false, nil
	}

	found := false
	for _, tokenInfo := range transferConfig.ChainToken[fromAvailable.GetId()].Token {
		if tokenInfo.Token.Symbol == token.Symbol {
			found = true
			break
		}
	}
	if !found {
		return false, nil
	}

	found = false
	for _, tokenInfo := range transferConfig.ChainToken[toAvailable.GetId()].Token {
		if tokenInfo.Token.Symbol == token.Symbol {
			found = true
			break
		}
	}
	if !found {
		return false, nil
	}
	return true, nil
}

func (s *CBridge) CalculateFees(from, to *params.Network, token *token.Token, amountIn *big.Int, nativeTokenPrice, tokenPrice float64, gasPrice *big.Float) (*big.Int, *big.Int, error) {
	amt, err := s.estimateAmt(from, to, amountIn, token.Symbol)
	if err != nil {
		return nil, nil, err
	}
	baseFee, _ := new(big.Int).SetString(amt.BaseFee, 10)
	percFee, _ := new(big.Int).SetString(amt.PercFee, 10)

	return big.NewInt(0), new(big.Int).Add(baseFee, percFee), nil
}

func (s *CBridge) EstimateGas(from, to *params.Network, token *token.Token, amountIn *big.Int) (uint64, error) {
	// TODO: replace by estimate function
	if token.IsNative() {
		return 22000, nil // default gas limit for eth transaction
	}

	return 200000, nil //default gas limit for erc20 transaction
}

func (s *CBridge) GetContractAddress(network *params.Network, token *token.Token) *common.Address {
	transferConfig, err := s.getTransferConfig(network.IsTest)
	if err != nil {
		return nil
	}
	if transferConfig.Err != nil {
		return nil
	}

	for _, chain := range transferConfig.Chains {
		if uint64(chain.Id) == network.ChainID {
			addr := common.HexToAddress(chain.ContractAddr)
			return &addr
		}
	}

	return nil
}

func (s *CBridge) Send(sendArgs *TransactionBridge, verifiedAccount *account.SelectedExtKey) (types.Hash, error) {
	fromNetwork := s.rpcClient.NetworkManager.Find(sendArgs.ChainID)
	if fromNetwork == nil {
		return types.HexToHash(""), errors.New("network not found")
	}
	tk := s.tokenManager.FindToken(fromNetwork, sendArgs.CbridgeTx.Symbol)
	if tk == nil {
		return types.HexToHash(""), errors.New("token not found")
	}
	addrs := s.GetContractAddress(fromNetwork, nil)
	if addrs == nil {
		return types.HexToHash(""), errors.New("contract not found")
	}

	backend, err := s.rpcClient.EthClient(sendArgs.ChainID)
	if err != nil {
		return types.HexToHash(""), err
	}
	contract, err := celer.NewCeler(*addrs, backend)
	if err != nil {
		return types.HexToHash(""), err
	}

	txOpts := sendArgs.CbridgeTx.ToTransactOpts(getSigner(sendArgs.ChainID, sendArgs.CbridgeTx.From, verifiedAccount))
	var tx *ethTypes.Transaction
	if tk.IsNative() {
		tx, err = contract.SendNative(
			txOpts,
			sendArgs.CbridgeTx.Recipient,
			(*big.Int)(sendArgs.CbridgeTx.Amount),
			sendArgs.CbridgeTx.ChainID,
			uint64(time.Now().UnixMilli()),
			500,
		)
	} else {
		tx, err = contract.Send(
			txOpts,
			sendArgs.CbridgeTx.Recipient,
			tk.Address,
			(*big.Int)(sendArgs.CbridgeTx.Amount),
			sendArgs.CbridgeTx.ChainID,
			uint64(time.Now().UnixMilli()),
			500,
		)
	}
	if err != nil {
		return types.HexToHash(""), err
	}

	return types.Hash(tx.Hash()), nil
}

func (s *CBridge) CalculateAmountOut(from, to *params.Network, amountIn *big.Int, symbol string) (*big.Int, error) {
	amt, err := s.estimateAmt(from, to, amountIn, symbol)
	if err != nil {
		return nil, err
	}
	if amt.Err != nil {
		return nil, err
	}
	amountOut, _ := new(big.Int).SetString(amt.EqValueTokenAmt, 10)
	return amountOut, nil
}