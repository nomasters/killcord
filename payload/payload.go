package payload

import (
	"io"
)

type Payload interface {
	Encrypt() error
	Decrypt() error
	Reader() io.ReadCloser
	Writer() io.WriteCloser
	GetConfig() Config
	SetConfig(Config) error
}

type Config struct {
	Status   string                 `toml:"status"`
	Provider string                 `toml:"provider"`
	Key      string                 `toml:"key"`
	Settings map[string]interface{} `toml:"settings"`
}
