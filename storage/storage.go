package storage

import (
	"io"
)

type Storage interface {
	DeployPayload(io.Reader) error
	GetPayload() (io.ReadCloser, error)
	GetConfig() Config
	SetConfig(Config) error
}

type Config struct {
	Status   string                 `toml:"status"`
	Provider string                 `toml:"provider"`
	Settings map[string]interface{} `toml:"settings"`
}
