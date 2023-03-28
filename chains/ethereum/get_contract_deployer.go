package ethereum

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	werr "github.com/doyensec/wallet-info/errors"
	"github.com/doyensec/wallet-info/utils"
)

type contractResp struct {
	Address string `json:"contractAddress"`
	Creator string `json:"contractCreator"`
	TxHash  string `json:"txHash"`
}

type contractDeployerResp struct {
	Message string         `json:"message"`
	Status  string         `json:"status"`
	Result  []contractResp `json:"result"`
}

func (ec *EthereumClient) GetContractDeployer(address string) (string, error) {
	url := fmt.Sprintf("%s?module=contract&action=getcontractcreation&contractaddresses=%s&apikey=%s",
		ec.etherscanEndpoint, address, ec.etherscanApiKey)

	client := http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		utils.Logger.Errorw("failed to get contract deployer", err)
		return "", &werr.BlockchainError{}
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.Errorw("failed to parse response", err)
		return "", &werr.BlockchainError{}
	}

	body := &contractDeployerResp{}
	json.Unmarshal(bytes, body)

	if body.Status != "1" {
		utils.Logger.Errorf("eoa supplied. address has no deployer.")
		return "", nil
	}

	return body.Result[0].Creator, nil
}
