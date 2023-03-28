package ethereum

import (
	"github.com/doyensec/safeurl"
	"github.com/doyensec/wallet-info/config"
)

type EthereumClient struct {
	etherscanEndpoint string
	etherscanApiKey   string
}

func BuildClient(config *config.Config) *EthereumClient {
	return &EthereumClient{
		etherscanEndpoint: config.EtherscanEndpoint,
		etherscanApiKey:   config.EtherscanApiKey,
	}
}

func buildHttpClient() *safeurl.WrappedClient {
	config := safeurl.GetConfigBuilder().
		Build()
	return safeurl.Client(config)
}
