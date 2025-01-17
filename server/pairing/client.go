package pairing

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/btcsuite/btcutil/base58"

	"github.com/status-im/status-go/api"
	"github.com/status-im/status-go/logutils"
	"github.com/status-im/status-go/signal"
)

type Client struct {
	*http.Client
	PayloadManager
	rawMessagePayloadManager *RawMessagePayloadManager

	baseAddress     *url.URL
	certPEM         []byte
	serverPK        *ecdsa.PublicKey
	serverMode      Mode
	serverCert      *x509.Certificate
	serverChallenge []byte
}

func NewPairingClient(backend *api.GethStatusBackend, c *ConnectionParams, config *AccountPayloadManagerConfig) (*Client, error) {
	u, err := c.URL()
	if err != nil {
		return nil, err
	}

	serverCert, err := getServerCert(u)
	if err != nil {
		return nil, err
	}
	err = verifyCert(serverCert, c.publicKey)
	if err != nil {
		return nil, err
	}
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverCert.Raw})

	rootCAs, err := x509.SystemCertPool()
	if err != nil {
		return nil, err
	}
	if ok := rootCAs.AppendCertsFromPEM(certPem); !ok {
		return nil, fmt.Errorf("failed to append certPem to rootCAs")
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion:         tls.VersionTLS12,
			InsecureSkipVerify: false, // MUST BE FALSE
			RootCAs:            rootCAs,
		},
	}

	cj, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	pm, err := NewAccountPayloadManager(c.aesKey, config, logutils.ZapLogger().Named("Client"))
	if err != nil {
		return nil, err
	}

	rmpm, err := NewRawMessagePayloadManager(logutils.ZapLogger().Named("Client"), pm.accountPayload, c.aesKey, backend, config.KeystorePath)
	if err != nil {
		return nil, err
	}
	return &Client{
		Client:                   &http.Client{Transport: tr, Jar: cj},
		baseAddress:              u,
		certPEM:                  certPem,
		serverCert:               serverCert,
		serverPK:                 c.publicKey,
		serverMode:               c.serverMode,
		PayloadManager:           pm,
		rawMessagePayloadManager: rmpm,
	}, nil
}

func (c *Client) PairAccount() error {
	switch c.serverMode {
	case Receiving:
		return c.sendAccountData()
	case Sending:
		err := c.getChallenge()
		if err != nil {
			return err
		}
		return c.receiveAccountData()
	default:
		return fmt.Errorf("unrecognised server mode '%d'", c.serverMode)
	}
}

func (c *Client) PairSyncDevice() error {
	switch c.serverMode {
	case Receiving:
		return c.sendSyncDeviceData()
	case Sending:
		return c.receiveSyncDeviceData()
	default:
		return fmt.Errorf("unrecognised server mode '%d'", c.serverMode)
	}
}

func (c *Client) sendSyncDeviceData() error {
	err := c.rawMessagePayloadManager.Mount()
	if err != nil {
		return err
	}

	c.baseAddress.Path = pairingSyncDeviceReceive
	resp, err := c.Post(c.baseAddress.String(), "application/octet-stream", bytes.NewBuffer(c.rawMessagePayloadManager.ToSend()))
	if err != nil {
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionSyncDevice})
		return err
	}

	if resp.StatusCode != http.StatusOK {
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionSyncDevice})
		return fmt.Errorf("status not ok, received '%s'", resp.Status)
	}

	signal.SendLocalPairingEvent(Event{Type: EventTransferSuccess, Action: ActionSyncDevice})
	return nil
}

func (c *Client) receiveSyncDeviceData() error {
	c.baseAddress.Path = pairingSyncDeviceSend
	req, err := http.NewRequest(http.MethodGet, c.baseAddress.String(), nil)
	if err != nil {
		return err
	}

	if c.serverChallenge != nil {
		ec, err := c.PayloadManager.EncryptPlain(c.serverChallenge)
		if err != nil {
			return err
		}

		req.Header.Set(sessionChallenge, base58.Encode(ec))
	}

	resp, err := c.Do(req)
	if err != nil {
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionSyncDevice})
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status not ok, received '%s'", resp.Status)
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionSyncDevice})
		return err
	}

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionSyncDevice})
		return err
	}
	signal.SendLocalPairingEvent(Event{Type: EventTransferSuccess, Action: ActionSyncDevice})

	err = c.rawMessagePayloadManager.Receive(payload)
	if err != nil {
		signal.SendLocalPairingEvent(Event{Type: EventProcessError, Error: err.Error(), Action: ActionSyncDevice})
		return err
	}
	signal.SendLocalPairingEvent(Event{Type: EventProcessSuccess, Action: ActionSyncDevice})
	return nil
}

func (c *Client) sendAccountData() error {
	err := c.Mount()
	if err != nil {
		return err
	}

	c.baseAddress.Path = pairingReceiveAccount
	resp, err := c.Post(c.baseAddress.String(), "application/octet-stream", bytes.NewBuffer(c.PayloadManager.ToSend()))
	if err != nil {
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionPairingAccount})
		return err
	}

	if resp.StatusCode != http.StatusOK {
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionPairingAccount})
		return fmt.Errorf("status not ok, received '%s'", resp.Status)
	}

	signal.SendLocalPairingEvent(Event{Type: EventTransferSuccess, Action: ActionPairingAccount})

	c.PayloadManager.LockPayload()
	return nil
}

func (c *Client) receiveAccountData() error {
	c.baseAddress.Path = pairingSendAccount
	req, err := http.NewRequest(http.MethodGet, c.baseAddress.String(), nil)
	if err != nil {
		return err
	}

	if c.serverChallenge != nil {
		ec, err := c.PayloadManager.EncryptPlain(c.serverChallenge)
		if err != nil {
			return err
		}

		req.Header.Set(sessionChallenge, base58.Encode(ec))
	}

	resp, err := c.Do(req)
	if err != nil {
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionPairingAccount})
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status not ok, received '%s'", resp.Status)
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionPairingAccount})
		return err
	}

	payload, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		signal.SendLocalPairingEvent(Event{Type: EventTransferError, Error: err.Error(), Action: ActionPairingAccount})
		return err
	}
	signal.SendLocalPairingEvent(Event{Type: EventTransferSuccess, Action: ActionPairingAccount})

	err = c.PayloadManager.Receive(payload)
	if err != nil {
		signal.SendLocalPairingEvent(Event{Type: EventProcessError, Error: err.Error(), Action: ActionPairingAccount})
		return err
	}
	signal.SendLocalPairingEvent(Event{Type: EventProcessSuccess, Action: ActionPairingAccount})
	return nil
}

func (c *Client) getChallenge() error {
	c.baseAddress.Path = pairingChallenge
	resp, err := c.Get(c.baseAddress.String())
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("status not ok, received '%s'", resp.Status)
	}

	c.serverChallenge, err = ioutil.ReadAll(resp.Body)
	return err
}

func StartUpPairingClient(backend *api.GethStatusBackend, cs, configJSON string) error {
	c, err := setupClient(backend, cs, configJSON)
	if err != nil {
		return err
	}
	err = c.PairAccount()
	if err != nil {
		return err
	}
	return c.PairSyncDevice()
}

func setupClient(backend *api.GethStatusBackend, cs string, configJSON string) (*Client, error) {
	var conf PayloadSourceConfig
	err := json.Unmarshal([]byte(configJSON), &conf)
	if err != nil {
		return nil, err
	}

	ccp := new(ConnectionParams)
	err = ccp.FromString(cs)
	if err != nil {
		return nil, err
	}

	c, err := NewPairingClient(backend, ccp, &AccountPayloadManagerConfig{backend.GetMultiaccountDB(), &conf})
	if err != nil {
		return nil, err
	}
	return c, nil
}
