package killcord

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/nomasters/killcord/contract"
)

const (
	defaultETHRPCPDev  = "https://ropsten.infura.io/tDr4BM10GNjusw3XnkrT"
	defaultETHRPCDProd = "https://mainnet.infura.io/tDr4BM10GNjusw3XnkrT"
)

var (
	contractDir          string
	fullKeyStorePath     string
	relativeKeyStorePath string
	ethereumRPCPath      string = defaultETHRPCPDev
)

func init() {
	contractDir = filepath.Join(ProjectPath, "contract")
	fullKeyStorePath = filepath.Join(contractDir, "data", "keystore")
	relativeKeyStorePath = filepath.Join("contract", "data", "keystore")
}

// A simple func to convert wei to ETH
func weiToETH(i *big.Int) float64 {
	f := float64(i.Int64())
	return f / 1000000000000000000
}

// Gracefully resolve Ethereum RPC path by waterfalling through
// Options > Project > Defaults

func (s *Session) setEthereumRPCPath() {
	if s.Options.Contract.RPCURL != "" {
		ethereumRPCPath = s.Options.Contract.RPCURL
		return
	}
	if s.Config.Contract.RPCURL != "" {
		ethereumRPCPath = s.Config.Contract.RPCURL
		return
	}
	if s.Config.Contract.Mode == "mainnet" {
		ethereumRPCPath = defaultETHRPCDProd
		return
	}
	s.Config.Contract.Mode = "testnet"
	ethereumRPCPath = defaultETHRPCPDev
}

func (s *Session) ConfigEthereum() error {
	ks := newKeyStore()
	if err := s.Config.Contract.Owner.New(ks); err != nil {
		return err
	}
	if err := s.Config.Contract.Publisher.New(ks); err != nil {
		return err
	}
	fmt.Println("ethereum: initializing")
	s.Config.Contract.Provider = "ethereum"
	s.Config.Contract.Status = "initialized"
	fmt.Println("ethereum: configured")

	fmt.Printf(`

Congrats! You've successfully initialized your ethereum owner and publisher accounts. 
Next, you'll need to add a little bit of ETH to both accounts to move forward.

You should add a minimum of:

- 0.01  ETH to your owner account
- 0.005 ETH to your publisher account

your owner account address is:     0x%v
your publisher account address is: 0x%v

You can check your ethereum account balances with the "killcord status" command.
If this is your first time using Ethereum, metamask (https://metamask.io/) is the
easiest way to get started.

`, s.Config.Contract.Owner.Address, s.Config.Contract.Publisher.Address)

	return nil
}

func newKeyStore() *keystore.KeyStore {
	return keystore.NewKeyStore(fullKeyStorePath, keystore.StandardScryptN, keystore.StandardScryptP)
}

func (a *AccountConfig) New(ks *keystore.KeyStore) error {
	pw := generateKey()
	newAcc, err := ks.NewAccount(pw)
	if err != nil {
		return err
	}
	a.Address = hex.EncodeToString(newAcc.Address[:])
	a.Password = pw
	a.KeyStore, err = getKeyStore(a)
	if err != nil {
		return err
	}
	os.RemoveAll(contractDir)
	return nil
}

func newContractAuthSession(account AccountConfig, contractID string) (*contract.KillCordSession, error) {
	conn, err := ethclient.Dial(ethereumRPCPath)
	if err != nil {
		return &contract.KillCordSession{}, fmt.Errorf("Failed to connect to the Ethereum client: %v", err)
	}
	auth, err := bind.NewTransactor(strings.NewReader(account.KeyStore), account.Password)
	if err != nil {
		return &contract.KillCordSession{}, fmt.Errorf("Failed to create authorized transactor: %v", err)
	}
	killcord, err := contract.NewKillCord(common.HexToAddress("0x"+contractID), conn)
	if err != nil {
		return &contract.KillCordSession{}, fmt.Errorf("Failed to instantiate a killcord contract: %v", err)
	}
	return &contract.KillCordSession{
		Contract: killcord,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
		TransactOpts: bind.TransactOpts{
			From:   auth.From,
			Signer: auth.Signer,
		},
	}, nil
}

func newContractCallerSession(contractID string) (*contract.KillCordSession, error) {
	conn, err := ethclient.Dial(ethereumRPCPath)
	if err != nil {
		return &contract.KillCordSession{}, fmt.Errorf("Failed to connect to the Ethereum client: %v", err)
	}
	killcord, err := contract.NewKillCord(common.HexToAddress("0x"+contractID), conn)
	if err != nil {
		return &contract.KillCordSession{}, fmt.Errorf("Failed to instantiate a killcord contract: %v", err)
	}
	return &contract.KillCordSession{
		Contract: killcord,
		CallOpts: bind.CallOpts{
			Pending: true,
		},
	}, nil
}

func GetLastCheckIn(contractID string) (time.Time, error) {
	session, err := newContractCallerSession(contractID)
	if err != nil {
		return time.Now(), err
	}
	timeStamp, err := session.GetLastCheckIn()
	if err != nil {
		return time.Now(), fmt.Errorf("Failed to get last checkin: %v", err)
	}
	return time.Unix(timeStamp.Int64(), 0), nil
}

func GetKey(contractID string) (string, error) {
	session, err := newContractCallerSession(contractID)
	if err != nil {
		return "", err
	}
	key, err := session.GetKey()
	if err != nil {
		return "", fmt.Errorf("Failed to get last checkin: %v", err)
	}
	return key, nil
}

func GetOwner(contractID string) (string, error) {
	session, err := newContractCallerSession(contractID)
	if err != nil {
		return "", err
	}
	address, err := session.GetOwner()
	if err != nil {
		return "", fmt.Errorf("Failed to get last checkin: %v", err)
	}
	return address.String(), nil
}

func GetPublisher(contractID string) (string, error) {
	session, err := newContractCallerSession(contractID)
	if err != nil {
		return "", err
	}
	address, err := session.GetPublisher()
	if err != nil {
		return "", fmt.Errorf("Failed to get last checkin: %v", err)
	}
	return address.String(), nil
}

func GetPayloadEndpoint(contractID string) (string, error) {
	session, err := newContractCallerSession(contractID)
	if err != nil {
		return "", err
	}
	endpoint, err := session.GetPayloadEndpoint()
	if err != nil {
		return "", fmt.Errorf("Failed to get last checkin: %v", err)
	}
	return endpoint, nil
}

// Runs a simple checkin to the contract with the owner account.
// TODO: support options for confirming checkin, not just submitting it
func CheckIn(account AccountConfig, contractID string) error {
	session, err := newContractAuthSession(account, contractID)
	if err != nil {
		return err
	}
	if _, err = session.CheckIn(); err != nil {
		return fmt.Errorf("Failed to set Endpoint: %v", err)
	}
	fmt.Println("checkin successfully submitted")
	return nil
}

func (s *Session) CheckIn() error {
	if err := CheckIn(s.Config.Contract.Owner, s.Config.Contract.ID); err != nil {
		return err
	}
	return nil
}

func SetKey(account AccountConfig, contractID string, secret string) error {
	session, err := newContractAuthSession(account, contractID)
	if err != nil {
		return err
	}
	if _, err := session.SetKey(secret); err != nil {
		return fmt.Errorf("Failed to set Publishable Key: %v", err)
	}
	fmt.Printf("key publication submitted with %v\n", account.Address)
	return nil
}

func KillContract(account AccountConfig, contractID string) error {
	session, err := newContractAuthSession(account, contractID)
	if err != nil {
		return err
	}
	if _, err := session.Kill(); err != nil {
		return fmt.Errorf("Failed to set Endpoint: %v", err)
	}
	fmt.Println("contract kill submitted")
	return nil
}

func (s *Session) KillContract() error {
	if err := KillContract(s.Config.Contract.Owner, s.Config.Contract.ID); err != nil {
		return err
	}
	return nil
}

func SetPayloadEndpoint(account AccountConfig, contractID string, payloadID string) error {
	session, err := newContractAuthSession(account, contractID)
	if err != nil {
		return err
	}
	if _, err := session.SetPayloadEndpoint(payloadID); err != nil {
		return fmt.Errorf("Failed to set Endpoint: %v", err)
	}
	fmt.Println("payload endpoint successfully submitted to contract")
	return nil
}

func (s *Session) DeployContract() error {
	if s.Config.Contract.ID != "" {
		return fmt.Errorf("contract 0x%v already deployed, skipping", s.Config.Contract.ID)
	}
	conn, err := ethclient.Dial(ethereumRPCPath)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return err
	}
	auth, err := bind.NewTransactor(strings.NewReader(s.Config.Contract.Owner.KeyStore), s.Config.Contract.Owner.Password)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
		return err
	}
	// TODO: this was set arbitrarily, should dive into this more
	auth.GasLimit = big.NewInt(2200000)
	publisher := common.HexToAddress("0x" + s.Config.Contract.Publisher.Address)
	address, tx, _, err := contract.DeployKillCord(auth, conn, publisher)
	if err != nil {
		log.Fatalf("Failed to deploy new killcord contract: %v", err)
		return err
	}
	fmt.Printf("Contract pending deploy: 0x%x\n", address)
	fmt.Printf("Transaction waiting to be mined: 0x%x\n\n", tx.Hash())
	s.Config.Contract.ID = hex.EncodeToString(address[:])
	return nil
}

func getBalance(account string) float64 {
	conn, err := ethclient.Dial(ethereumRPCPath)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	a := common.HexToAddress("0x" + account)
	balance, err := conn.BalanceAt(context.TODO(), a, nil)
	if err != nil {
		log.Fatalf("balance check failed %v\n", err)
	}
	b := balance
	return weiToETH(b)
}

func getKeyStore(account *AccountConfig) (string, error) {
	var file string
	files, err := filepath.Glob(relativeKeyStorePath + "/*")
	if err != nil {
		return "", err
	}
	for _, f := range files {
		if strings.Contains(f, account.Address) {
			file = f
			break
		}
	}
	if file == "" {
		return "", errors.New("No Contract Account Found")
	}

	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
