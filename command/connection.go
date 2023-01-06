package command

import (
	"code.cloudfoundry.org/cli/plugin"
	"github.com/cloudfoundry-community/go-cfclient/v3/client"
	"github.com/cloudfoundry-community/go-cfclient/v3/config"
	"strings"
)

func createCFClient(cliConnection plugin.CliConnection) (*client.Client, error) {
	u, err := cliConnection.ApiEndpoint()
	if err != nil {
		return nil, err
	}

	t, err := cliConnection.AccessToken()
	if err != nil {
		return nil, err
	}
	t = strings.TrimPrefix(t, "bearer ")

	cfg, err := config.NewToken(u, t)
	if err != nil {
		return nil, err
	}

	skipSSLValidation, err := cliConnection.IsSSLDisabled()
	if err != nil {
		return nil, err
	}
	cfg.WithSkipTLSValidation(skipSSLValidation)

	return client.New(cfg)
}
