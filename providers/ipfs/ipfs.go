package ipfs

import (
	//"github.com/nomasters/killcord"
	"github.com/nomasters/killcord/storage"
)

const (
	defaultRPCURL = "https://ipfs.infura.io:5001"
)

var (
	payloadEncryptedPath string
	payloadTempPath      string
	payloadDecryptPath   string
	RPCURL               string = defaultRPCURL
)

func init() {
	payloadEncryptedPath = filepath.Join(strings.Split(encryptPrefix, "/")...)
	payloadTempPath = filepath.Join(strings.Split(tempPrefix, "/")...)
	payloadDecryptPath = filepath.Join(strings.Split(decryptPrefix, "/")...)
}

type Options struct {
	RPCURL string
}

type Storage struct {
	Options Options
	Config  storage.Config
}

func NewStorage(opts Options, config storage.Config) Storage {
	s := new(Storage)
	s.setOptions(opts)
	return s
}

func (s *Storage) setOptions(opts Options) {
}

func (s *Storage) setRPCURL() {
	if s.Options.RPCURL != "" {
		RPCURL = s.Options.RPCURL
		return
	}
	if x, ok := s.Config.Settings.RPCURL.(string); ok {
		RPCURL = x
	}
	RPCURL = defaultRPCURL
}
