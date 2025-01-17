package pairing

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/status-im/status-go/api"

	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/status-im/status-go/account/generator"
	"github.com/status-im/status-go/eth-node/keystore"
	"github.com/status-im/status-go/multiaccounts"
	"github.com/status-im/status-go/protocol/common"
	"github.com/status-im/status-go/protocol/protobuf"
)

// PayloadManager is the interface for PayloadManagers and wraps the basic functions for fulfilling payload management
type PayloadManager interface {
	// Mount Loads the payload into the PayloadManager's state
	Mount() error

	// Receive stores data from an inbound source into the PayloadManager's state
	Receive(data []byte) error

	// ToSend returns an outbound safe (encrypted) payload
	ToSend() []byte

	// Received returns a decrypted and parsed payload from an inbound source
	Received() []byte

	// ResetPayload resets all payloads the PayloadManager has in its state
	ResetPayload()

	// EncryptPlain encrypts the given plaintext using internal key(s)
	EncryptPlain(plaintext []byte) ([]byte, error)

	// LockPayload prevents future excess to outbound safe and received data
	LockPayload()
}

// PayloadSourceConfig represents location and access data of the pairing payload
// ONLY available from the application client
type PayloadSourceConfig struct {
	// required
	KeystorePath string `json:"keystorePath"`
	// following 2 fields r optional.
	// optional cases:
	// 1. server mode is Receiving and server side doesn't contain this info
	// 2. server mode is Sending and client side doesn't contain this info
	// they are required in other cases
	KeyUID   string `json:"keyUID"`
	Password string `json:"password"`
}

// AccountPayloadManagerConfig represents the initialisation parameters required for a AccountPayloadManager
type AccountPayloadManagerConfig struct {
	DB *multiaccounts.Database
	*PayloadSourceConfig
}

// AccountPayloadManager is responsible for the whole lifecycle of a AccountPayload
type AccountPayloadManager struct {
	logger         *zap.Logger
	accountPayload *AccountPayload
	*PayloadEncryptionManager
	accountPayloadMarshaller *AccountPayloadMarshaller
	payloadRepository        PayloadRepository
}

// NewAccountPayloadManager generates a new and initialised AccountPayloadManager
func NewAccountPayloadManager(aesKey []byte, config *AccountPayloadManagerConfig, logger *zap.Logger) (*AccountPayloadManager, error) {
	l := logger.Named("AccountPayloadManager")
	l.Debug("fired", zap.Binary("aesKey", aesKey), zap.Any("config", config))

	pem, err := NewPayloadEncryptionManager(aesKey, l)
	if err != nil {
		return nil, err
	}

	// A new SHARED AccountPayload
	p := new(AccountPayload)

	return &AccountPayloadManager{
		logger:                   l,
		accountPayload:           p,
		PayloadEncryptionManager: pem,
		accountPayloadMarshaller: NewPairingPayloadMarshaller(p, l),
		payloadRepository:        NewAccountPayloadRepository(p, config),
	}, nil
}

// Mount loads and prepares the payload to be stored in the AccountPayloadManager's state ready for later access
func (apm *AccountPayloadManager) Mount() error {
	l := apm.logger.Named("Mount()")
	l.Debug("fired")

	err := apm.payloadRepository.LoadFromSource()
	if err != nil {
		return err
	}
	l.Debug("after LoadFromSource")

	pb, err := apm.accountPayloadMarshaller.MarshalToProtobuf()
	if err != nil {
		return err
	}
	l.Debug(
		"after MarshalToProtobuf",
		zap.Any("accountPayloadMarshaller.accountPayloadMarshaller.keys", apm.accountPayloadMarshaller.keys),
		zap.Any("accountPayloadMarshaller.accountPayloadMarshaller.multiaccount", apm.accountPayloadMarshaller.multiaccount),
		zap.String("accountPayloadMarshaller.accountPayloadMarshaller.password", apm.accountPayloadMarshaller.password),
		zap.Binary("pb", pb),
	)

	return apm.Encrypt(pb)
}

// Receive takes a []byte representing raw data, parses and stores the data
func (apm *AccountPayloadManager) Receive(data []byte) error {
	l := apm.logger.Named("Receive()")
	l.Debug("fired")

	err := apm.Decrypt(data)
	if err != nil {
		return err
	}
	l.Debug("after Decrypt")

	err = apm.accountPayloadMarshaller.UnmarshalProtobuf(apm.Received())
	if err != nil {
		return err
	}
	l.Debug(
		"after UnmarshalProtobuf",
		zap.Any("accountPayloadMarshaller.accountPayloadMarshaller.keys", apm.accountPayloadMarshaller.keys),
		zap.Any("accountPayloadMarshaller.accountPayloadMarshaller.multiaccount", apm.accountPayloadMarshaller.multiaccount),
		zap.String("accountPayloadMarshaller.accountPayloadMarshaller.password", apm.accountPayloadMarshaller.password),
		zap.Binary("accountPayloadMarshaller.Received()", apm.Received()),
	)

	return apm.payloadRepository.StoreToSource()
}

// ResetPayload resets all payload state managed by the AccountPayloadManager
func (apm *AccountPayloadManager) ResetPayload() {
	apm.accountPayload.ResetPayload()
	apm.PayloadEncryptionManager.ResetPayload()
}

// EncryptionPayload represents the plain text and encrypted text of payload data
type EncryptionPayload struct {
	plain     []byte
	encrypted []byte
	locked    bool
}

func (ep *EncryptionPayload) lock() {
	ep.locked = true
}

// PayloadEncryptionManager is responsible for encrypting and decrypting payload data
type PayloadEncryptionManager struct {
	logger   *zap.Logger
	aesKey   []byte
	toSend   *EncryptionPayload
	received *EncryptionPayload
}

func NewPayloadEncryptionManager(aesKey []byte, logger *zap.Logger) (*PayloadEncryptionManager, error) {
	return &PayloadEncryptionManager{logger.Named("PayloadEncryptionManager"), aesKey, new(EncryptionPayload), new(EncryptionPayload)}, nil
}

// EncryptPlain encrypts any given plain text using the internal AES key and returns the encrypted value
// This function is different to Encrypt as the internal EncryptionPayload.encrypted value is not set
func (pem *PayloadEncryptionManager) EncryptPlain(plaintext []byte) ([]byte, error) {
	l := pem.logger.Named("EncryptPlain()")
	l.Debug("fired")

	return common.Encrypt(plaintext, pem.aesKey, rand.Reader)
}

func (pem *PayloadEncryptionManager) Encrypt(data []byte) error {
	l := pem.logger.Named("Encrypt()")
	l.Debug("fired")

	ep, err := common.Encrypt(data, pem.aesKey, rand.Reader)
	if err != nil {
		return err
	}

	pem.toSend.plain = data
	pem.toSend.encrypted = ep

	l.Debug(
		"after common.Encrypt",
		zap.Binary("data", data),
		zap.Binary("pem.aesKey", pem.aesKey),
		zap.Binary("ep", ep),
	)

	return nil
}

func (pem *PayloadEncryptionManager) Decrypt(data []byte) error {
	l := pem.logger.Named("Decrypt()")
	l.Debug("fired")

	pd, err := common.Decrypt(data, pem.aesKey)
	l.Debug(
		"after common.Decrypt(data, pem.aesKey)",
		zap.Binary("data", data),
		zap.Binary("pem.aesKey", pem.aesKey),
		zap.Binary("pd", pd),
		zap.Error(err),
	)
	if err != nil {
		return err
	}

	pem.received.encrypted = data
	pem.received.plain = pd
	return nil
}

func (pem *PayloadEncryptionManager) ToSend() []byte {
	if pem.toSend.locked {
		return nil
	}
	return pem.toSend.encrypted
}

func (pem *PayloadEncryptionManager) Received() []byte {
	if pem.toSend.locked {
		return nil
	}
	return pem.received.plain
}

func (pem *PayloadEncryptionManager) ResetPayload() {
	pem.toSend = new(EncryptionPayload)
	pem.received = new(EncryptionPayload)
}

func (pem *PayloadEncryptionManager) LockPayload() {
	l := pem.logger.Named("LockPayload")
	l.Debug("fired")

	pem.toSend.lock()
	pem.received.lock()
}

// AccountPayload represents the payload structure a Server handles
type AccountPayload struct {
	keys         map[string][]byte
	multiaccount *multiaccounts.Account
	password     string
}

func (ap *AccountPayload) ResetPayload() {
	*ap = AccountPayload{}
}

// AccountPayloadMarshaller is responsible for marshalling and unmarshalling Server payload data
type AccountPayloadMarshaller struct {
	logger *zap.Logger
	*AccountPayload
}

func NewPairingPayloadMarshaller(ap *AccountPayload, logger *zap.Logger) *AccountPayloadMarshaller {
	return &AccountPayloadMarshaller{logger: logger, AccountPayload: ap}
}

func (ppm *AccountPayloadMarshaller) MarshalToProtobuf() ([]byte, error) {
	return proto.Marshal(&protobuf.LocalPairingPayload{
		Keys:         ppm.accountKeysToProtobuf(),
		Multiaccount: ppm.multiaccount.ToProtobuf(),
		Password:     ppm.password,
	})
}

func (ppm *AccountPayloadMarshaller) accountKeysToProtobuf() []*protobuf.LocalPairingPayload_Key {
	var keys []*protobuf.LocalPairingPayload_Key
	for name, data := range ppm.keys {
		keys = append(keys, &protobuf.LocalPairingPayload_Key{Name: name, Data: data})
	}
	return keys
}

func (ppm *AccountPayloadMarshaller) UnmarshalProtobuf(data []byte) error {
	l := ppm.logger.Named("UnmarshalProtobuf()")
	l.Debug("fired")

	pb := new(protobuf.LocalPairingPayload)
	err := proto.Unmarshal(data, pb)
	l.Debug(
		"after protobuf.LocalPairingPayload",
		zap.Any("pb", pb),
		zap.Any("pb.Multiaccount", pb.Multiaccount),
		zap.Any("pb.Keys", pb.Keys),
	)
	if err != nil {
		return err
	}

	ppm.accountKeysFromProtobuf(pb.Keys)
	ppm.multiaccountFromProtobuf(pb.Multiaccount)
	ppm.password = pb.Password
	return nil
}

func (ppm *AccountPayloadMarshaller) accountKeysFromProtobuf(pbKeys []*protobuf.LocalPairingPayload_Key) {
	l := ppm.logger.Named("accountKeysFromProtobuf()")
	l.Debug("fired")

	if ppm.keys == nil {
		ppm.keys = make(map[string][]byte)
	}

	for _, key := range pbKeys {
		ppm.keys[key.Name] = key.Data
	}
	l.Debug(
		"after for _, key := range pbKeys",
		zap.Any("pbKeys", pbKeys),
		zap.Any("accountPayloadMarshaller.keys", ppm.keys),
	)
}

func (ppm *AccountPayloadMarshaller) multiaccountFromProtobuf(pbMultiAccount *protobuf.MultiAccount) {
	ppm.multiaccount = new(multiaccounts.Account)
	ppm.multiaccount.FromProtobuf(pbMultiAccount)
}

type PayloadRepository interface {
	LoadFromSource() error
	StoreToSource() error
}

// AccountPayloadRepository is responsible for loading, parsing, validating and storing Server payload data
type AccountPayloadRepository struct {
	*AccountPayload

	multiaccountsDB *multiaccounts.Database

	keystorePath, keyUID string
}

func NewAccountPayloadRepository(p *AccountPayload, config *AccountPayloadManagerConfig) *AccountPayloadRepository {
	ppr := &AccountPayloadRepository{
		AccountPayload: p,
	}

	if config == nil {
		return ppr
	}

	ppr.multiaccountsDB = config.DB
	ppr.keystorePath = config.KeystorePath
	ppr.keyUID = config.KeyUID
	ppr.password = config.Password
	return ppr
}

func (apr *AccountPayloadRepository) LoadFromSource() error {
	err := apr.loadKeys(apr.keystorePath)
	if err != nil {
		return err
	}

	err = apr.validateKeys(apr.password)
	if err != nil {
		return err
	}

	apr.multiaccount, err = apr.multiaccountsDB.GetAccount(apr.keyUID)
	if err != nil {
		return err
	}

	return nil
}

func (apr *AccountPayloadRepository) loadKeys(keyStorePath string) error {
	apr.keys = make(map[string][]byte)

	fileWalker := func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if fileInfo.IsDir() || filepath.Dir(path) != keyStorePath {
			return nil
		}

		rawKeyFile, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("invalid account key file: %v", err)
		}

		accountKey := new(keystore.EncryptedKeyJSONV3)
		if err := json.Unmarshal(rawKeyFile, &accountKey); err != nil {
			return fmt.Errorf("failed to read key file: %s", err)
		}

		if len(accountKey.Address) != 40 {
			return fmt.Errorf("account key address has invalid length '%s'", accountKey.Address)
		}

		apr.keys[fileInfo.Name()] = rawKeyFile

		return nil
	}

	err := filepath.Walk(keyStorePath, fileWalker)
	if err != nil {
		return fmt.Errorf("cannot traverse key store folder: %v", err)
	}

	return nil
}

func (apr *AccountPayloadRepository) StoreToSource() error {
	err := apr.validateKeys(apr.password)
	if err != nil {
		return err
	}

	err = apr.storeKeys(apr.keystorePath)
	if err != nil {
		return err
	}

	err = apr.storeMultiAccount()
	if err != nil {
		return err
	}

	// TODO install PublicKey into settings, probably do this outside of StoreToSource
	return nil
}

func (apr *AccountPayloadRepository) validateKeys(password string) error {
	for _, key := range apr.keys {
		k, err := keystore.DecryptKey(key, password)
		if err != nil {
			return err
		}

		err = generator.ValidateKeystoreExtendedKey(k)
		if err != nil {
			return err
		}
	}

	return nil
}

func (apr *AccountPayloadRepository) storeKeys(keyStorePath string) error {
	if keyStorePath == "" {
		return fmt.Errorf("keyStorePath can not be empty")
	}

	_, lastDir := filepath.Split(keyStorePath)

	// If lastDir == "keystore" we presume we need to create the rest of the keystore path
	// else we presume the provided keystore is valid
	if lastDir == "keystore" {
		if apr.multiaccount == nil || apr.multiaccount.KeyUID == "" {
			return fmt.Errorf("no known Key UID")
		}
		keyStorePath = filepath.Join(keyStorePath, apr.multiaccount.KeyUID)

		err := os.MkdirAll(keyStorePath, 0777)
		if err != nil {
			return err
		}
	}

	for name, data := range apr.keys {
		accountKey := new(keystore.EncryptedKeyJSONV3)
		if err := json.Unmarshal(data, &accountKey); err != nil {
			return fmt.Errorf("failed to read key file: %s", err)
		}

		if len(accountKey.Address) != 40 {
			return fmt.Errorf("account key address has invalid length '%s'", accountKey.Address)
		}

		err := ioutil.WriteFile(filepath.Join(keyStorePath, name), data, 0600)
		if err != nil {
			return err
		}
	}
	return nil
}

func (apr *AccountPayloadRepository) storeMultiAccount() error {
	return apr.multiaccountsDB.SaveAccount(*apr.multiaccount)
}

type RawMessagePayloadManager struct {
	logger *zap.Logger
	// reference from AccountPayloadManager#accountPayload
	accountPayload *AccountPayload
	*PayloadEncryptionManager
	payloadRepository *RawMessageRepository
}

func NewRawMessagePayloadManager(logger *zap.Logger, accountPayload *AccountPayload, aesKey []byte, backend *api.GethStatusBackend, keystorePath string) (*RawMessagePayloadManager, error) {
	l := logger.Named("RawMessagePayloadManager")
	pem, err := NewPayloadEncryptionManager(aesKey, l)
	if err != nil {
		return nil, err
	}
	return &RawMessagePayloadManager{
		logger:                   l,
		accountPayload:           accountPayload,
		PayloadEncryptionManager: pem,
		payloadRepository:        NewRawMessageRepository(backend, keystorePath, accountPayload),
	}, nil
}

func (r *RawMessagePayloadManager) Mount() error {
	err := r.payloadRepository.LoadFromSource()
	if err != nil {
		return err
	}
	return r.Encrypt(r.payloadRepository.payload)
}

func (r *RawMessagePayloadManager) Receive(data []byte) error {
	err := r.Decrypt(data)
	if err != nil {
		return err
	}
	r.payloadRepository.payload = r.Received()
	return r.payloadRepository.StoreToSource()
}

func (r *RawMessagePayloadManager) ResetPayload() {
	r.payloadRepository.payload = make([]byte, 0)
	r.PayloadEncryptionManager.ResetPayload()
}

type RawMessageRepository struct {
	payload               []byte
	syncRawMessageHandler *SyncRawMessageHandler
	keystorePath          string
	accountPayload        *AccountPayload
}

func NewRawMessageRepository(backend *api.GethStatusBackend, keystorePath string, accountPayload *AccountPayload) *RawMessageRepository {
	return &RawMessageRepository{
		syncRawMessageHandler: NewSyncRawMessageHandler(backend),
		keystorePath:          keystorePath,
		payload:               make([]byte, 0),
		accountPayload:        accountPayload,
	}
}

func (r *RawMessageRepository) LoadFromSource() error {
	account := r.accountPayload.multiaccount
	if account == nil || account.KeyUID == "" {
		return fmt.Errorf("no known KeyUID when loading raw messages")
	}
	payload, err := r.syncRawMessageHandler.PrepareRawMessage(account.KeyUID)
	if err != nil {
		return err
	}
	r.payload = payload
	return nil
}

func (r *RawMessageRepository) StoreToSource() error {
	accountPayload := r.accountPayload
	if accountPayload == nil || accountPayload.multiaccount == nil {
		return fmt.Errorf("no known multiaccount when storing raw messages")
	}
	return r.syncRawMessageHandler.HandleRawMessage(accountPayload.multiaccount, accountPayload.password, r.keystorePath, r.payload)
}
