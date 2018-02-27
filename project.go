//go:generate ./scripts/contract-gen.sh

package killcord

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	Version                 = "0.0.1-alpha"
	defaultContractProvider = "ethereum"
	defaultPayloadProvider  = "ipfs"
)

var ProjectPath string

type Session struct {
	Config  ProjectConfig
	Options ProjectOptions
}

type ProjectOptions struct {
	DevMode   bool
	Type      string
	Audience  string
	Version   string
	Status    StatusOptions
	Payload   PayloadOptions
	Contract  ContractOptions
	Publisher PublisherOptions
}

type StatusOptions struct {
	ViewAll bool
}

type PayloadOptions struct {
	SourcePath      string
	DestinationPath string
	RPCURL          string
	Provider        string
	Secret          string
	ID              string
}

type ContractOptions struct {
	Kill      bool
	Provider  string
	ID        string
	RPCURL    string
	Owner     AccountConfig
	Publisher AccountConfig
}

type PublisherOptions struct {
	WarningThreshold int64
	PublishThreshold int64
	Address          string
	Password         string
	KeyStore         string
}

type ProjectConfig struct {
	Type      string          `toml:"type"`
	Version   string          `toml:"version"`
	Audience  string          `toml:"audience"`
	Created   time.Time       `toml:"created"`
	Updated   time.Time       `toml:"updated"`
	DevMode   bool            `toml:"devMode"`
	Payload   PayloadConfig   `toml:"payload"`
	Contract  ContractConfig  `toml:"contract"`
	Publisher PublisherConfig `toml:"publisher"`
	Status    string          `toml:"status"`
}

type PayloadConfig struct {
	Status   string `toml:"status"`
	Provider string `toml:"provider"`
	ID       string `toml:"id"`
	Secret   string `toml:"secret"`
	Mode     string `toml:"mode"`
	RPCURL   string `toml:"rpcUrl"`
}

type ContractConfig struct {
	Status    string        `toml:"status"`
	Provider  string        `toml:"provider"`
	ID        string        `toml:"id"`
	Owner     AccountConfig `toml:"owner"`
	Publisher AccountConfig `toml:"publisher"`
	Mode      string        `toml:"mode"`
	RPCURL    string        `toml:"rpcUrl"`
}

type AccountConfig struct {
	Address  string `toml:"address"`
	Password string `toml:"password"`
	KeyStore string `toml:"keystore"`
}

type PublisherConfig struct {
	Status           string `toml:"status"`
	WarningThreshold int64  `toml:"warningThreshold"`
	PublishThreshold int64  `toml:"publishThreshold"`
}

func init() {
	setAbsPath()
}

// Returns a new Killcord Session with defaults
func New() *Session {
	return &Session{}
}

// initializes ProjectConfig and ProjectOptions session variables
// for use by functions that use Payload and Contract RPC settings
func (s *Session) Init() {
	s.setEthereumRPCPath()
	s.setPayloadRPCPath()
}

func (s *Session) NewProject() error {
	if err := s.initProject(); err != nil {
		return err
	}
	// bootstrap owner and watcher projects
	switch s.Options.Type {
	case "owner":
		if err := s.initOwner(); err != nil {
			return err
		}
	case "watcher":
		if err := s.initWatcher(); err != nil {
			return err
		}
	default:
		return errors.New("new projects must be created for either owner or watcher")
	}
	return nil
}

func (s *Session) initWatcher() error {
	if err := s.initWatcherContract(); err != nil {
		return err
	}
	if s.Config.Payload.ID != "" {
		if err := s.GetPayload(); err != nil {
			return err
		}
	}
	s.Config.Payload.Status = "active"
	return nil
}

func (s *Session) initWatcherContract() error {
	s.Config.Contract.Status = "initialized"
	id, err := sanitizeContractId(s.Options.Contract.ID)
	if err != nil {
		return err
	}
	s.Config.Contract.ID = id
	if s.Config.DevMode == true {
		s.Config.Contract.Mode = "testnet"
	} else {
		s.Config.Contract.Mode = "mainnet"
	}
	s.setEthereumRPCPath()
	s.setPayloadRPCPath()
	endpoint, err := GetPayloadEndpoint(id)
	if err != nil {
		return err
	}
	if endpoint == "" {
		fmt.Println("Payload endpoint not yet registered, skipping")
		return nil
	} else {
		fmt.Println("adding payload endpoint to project configuration")
		s.Config.Payload.ID = endpoint
	}
	key, err := GetKey(s.Config.Contract.ID)
	if err != nil {
		return err
	}
	if key != "" {
		fmt.Println("adding decryption key to project configuration")
		s.Config.Payload.Secret = key
		fmt.Println("decrypt the payload by running: killcord decrypt")
	}
	return nil
}

func (s *Session) initOwner() error {
	// configure payload based on Provider
	if err := s.initPayload(); err != nil {
		return err
	}
	// configure payload based on Provider
	if err := s.initContract(); err != nil {
		return err
	}
	return nil
}

func (s *Session) initProject() error {
	// check that directory is empty
	if err := ensureEmptyDir(); err != nil {
		return err
	}
	if err := createDirs(); err != nil {
		return err
	}
	if err := validateProjectType(s.Options.Type); err != nil {
		return err
	}
	if err := s.setProviders(); err != nil {
		return err
	}
	// set project level values
	s.Config.Type = s.Options.Type
	s.Config.Version = Version
	s.Config.Audience = s.Options.Audience // TODO: not sure I remember what this is
	s.Config.DevMode = s.Options.DevMode
	s.Config.Created = time.Now()
	s.Config.Updated = time.Now()

	return nil
}

func (s *Session) setProviders() error {
	// attempt to set Contract Provider from Options, or set default provider
	if s.Options.Contract.Provider != "" {
		if err := validateContractProvider(s.Options.Contract.Provider); err != nil {
			return err
		}
		s.Config.Contract.Provider = s.Options.Contract.Provider
	} else {
		s.Config.Contract.Provider = defaultContractProvider
	}

	// attempt to set Payload Provider from Options, or set default provider
	if s.Options.Payload.Provider != "" {
		if err := validatePayloadProvider(s.Options.Payload.Provider); err != nil {
			return err
		}
		s.Config.Payload.Provider = s.Options.Payload.Provider
	} else {
		s.Config.Payload.Provider = defaultPayloadProvider
	}
	return nil
}

func sanitizeContractId(id string) (string, error) {
	// check and remove "0x" if hex prefix exists
	if id[:2] == "0x" {
		id = id[2:]
	}
	// check that hex encoding is valid
	if _, err := hex.DecodeString(id); err != nil {
		return id, err
	}
	return id, nil
}

func validateContractProvider(p string) error {
	switch p {
	case "ethereum":
	default:
		return errors.New("invalide contract provider: " + p)
	}
	return nil
}

func validatePayloadProvider(p string) error {
	switch p {
	case "ipfs":
	default:
		return errors.New("invalide payload provider: " + p)
	}
	return nil
}

func validateProjectType(t string) error {
	// filter ProjectType for one of proper commands
	switch t {
	case "owner":
	case "publisher":
	case "watcher":
	default:
		return errors.New("invalide project type in options")
	}
	return nil
}

func (s *Session) initPayload() error {
	switch s.Config.Payload.Provider {
	case "ipfs":
		s.Config.Payload.Status = "initialized"
	default:
		return fmt.Errorf("error: unrecognized provider: %v\n", s.Config.Payload.Provider)
	}
	return nil
}

func (s *Session) initContract() error {
	switch s.Config.Contract.Provider {
	case "ethereum":
		if s.Config.DevMode == true {
			s.Config.Contract.Mode = "testnet"
		} else {
			s.Config.Contract.Mode = "mainnet"
		}
		// TODO: this needs to be reimagined thorugh interfaces
		if err := s.ConfigEthereum(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("error: unrecognized provider: %v\n", s.Config.Contract.Provider)
	}
	return nil
}

func ensureEmptyDir() error {
	var files []string
	allFiles, err := filepath.Glob("*")
	if err != nil {
		return err
	}
	for _, f := range allFiles {
		// omit dotFiles
		if strings.HasPrefix(f, ".") {
			// skip dotFiles
		} else {
			files = append(files, f)
		}
	}
	if len(files) > 0 {
		return fmt.Errorf("error: init failed. Only init an empty directory.")
	}
	return nil
}

func setAbsPath() {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Printf("error: cannot set absolute path for project: %v\n", err)
		os.Exit(1)
	}
	ProjectPath = dir
}

func generateKey() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}

func createDirs() error {
	dirs := []string{
		"payload/source",
		"payload/encrypted",
		"payload/decrypted",
	}

	for _, dir := range dirs {
		// ensure path is contructed properly for either platform
		pathSplit := strings.Split(dir, "/")
		path := filepath.Join(pathSplit...)
		if err := os.MkdirAll(path, 0777); err != nil {
			return err
		}
	}
	return nil
}
