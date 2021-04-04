package env

import (
	"github.com/lbryio/lbry.go/v2/extras/errors"

	e "github.com/caarlos0/env"
)

// Config holds the environment configuration used by lighthouse.
type Config struct {
	APIServerPort      int    `env:"API_SERVER_PORT" envDefault:"7000"`
	IsDebug            bool   `env:"IS_DEBUG"`
	ChainQueryDSN      string `env:"CHAINQUERY_DSN"`
	InternalAPIsToken  string `env:"INTERNAL_APIS_TOKEN"`
	SlackHookURL       string `env:"SLACKHOOKURL"`
	SlackChannel       string `env:"SLACKCHANNEL"`
	BucketPath         string `env:"BUCKET_PATH"`
	BucketPrefixLength int    `env:"BUCKET_PREFIX_LENGTH"`
}

// NewWithEnvVars creates an Config from environment variables
func NewWithEnvVars() (*Config, error) {
	cfg := &Config{}
	err := e.Parse(cfg)
	if err != nil {
		return nil, errors.Err(err)
	}

	return cfg, nil
}
