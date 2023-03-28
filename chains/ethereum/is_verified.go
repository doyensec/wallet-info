package ethereum

import (
	"encoding/json"
	"io"
	"net/url"

	werr "github.com/doyensec/wallet-info/errors"
	"github.com/doyensec/wallet-info/utils"
)

type IsVerifiedResult struct {
	ContractName    string   `json:"ContractName"`
	ABI             []string `json:"ABI"`
	CompilerVersion string   `json:"CompilerVersion"`
	SourceCode      string   `json:"SourceCode"`
}

type IsVerifiedResp struct {
	Message string             `json:"message"`
	Status  string             `json:"status"`
	Result  []IsVerifiedResult `json:"result"`
}

func (ec *EthereumClient) IsVerified(address string) (*IsVerifiedResp, error) {
	url, _ := url.Parse(ec.etherscanEndpoint)
	q := url.Query()
	q.Set("module", "contract")
	q.Set("action", "getsourcecode")
	q.Set("address", address)
	q.Set("apikey", ec.etherscanApiKey)
	url.RawQuery = q.Encode()

	u := url.String()
	utils.Logger.Infof("url => %s", u)

	client := buildHttpClient()
	resp, err := client.Get(u)
	if err != nil {
		utils.Logger.Errorw("failed to get verification info", err)
		return nil, &werr.BlockchainError{}
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.Errorw("failed to parse response", err)
		return nil, &werr.BlockchainError{}
	}

	body := &IsVerifiedResp{}
	json.Unmarshal(bytes, body)

	return body, nil
}
