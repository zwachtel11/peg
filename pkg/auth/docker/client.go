package docker

import (
	"github.com/zwachtel11/peg/pkg/auth"

	"github.com/docker/cli/cli/config"
	"github.com/docker/cli/cli/config/configfile"
	"github.com/docker/cli/cli/config/credentials"
)

// Client provides authentication operations for docker registries.
type Client struct {
	configs []*configfile.ConfigFile
}

// NewClient
func NewClient() (auth.Client, error) {
	cfg, err := config.Load(config.Dir())
	if err != nil {
		return nil, err
	}
	if !cfg.ContainsAuth() {
		cfg.CredentialsStore = credentials.DetectDefaultStore(cfg.CredentialsStore)
	}

	return &Client{
		configs: []*configfile.ConfigFile{cfg},
	}, nil

}

func (c *Client) primaryCredentialsStore(hostname string) credentials.Store {
	return c.configs[0].GetCredentialsStore(hostname)
}
