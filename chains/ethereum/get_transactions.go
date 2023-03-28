package ethereum

import (
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"

	"github.com/doyensec/wallet-info/chains"
	werr "github.com/doyensec/wallet-info/errors"
	"github.com/doyensec/wallet-info/utils"
)

type GetTransactionsResp struct {
	Message int32              `json:"message"`
	Status  string             `json:"status"`
	Result  []TransactionsResp `json:"result"`
}

type TransactionsResp struct {
	BlockNumber uint64 `json:"blockNumber"`
	Timestamp   uint64 `json:"timeStamp"`
	Hash        string `json:"hash"`
	From        string `json:"from"`
	To          string `json:"to"`
}

func (ec *EthereumClient) GetTransactions(address string) (*chains.Transactions, error) {
	url, _ := url.Parse(ec.etherscanEndpoint)
	q := url.Query()
	q.Set("module", "account")
	q.Set("action", "txlist")
	q.Set("address", address)
	q.Set("startblock", "0")
	q.Set("endblock", strconv.Itoa(math.MaxInt64))
	q.Set("page", "0")
	q.Set("offset", "1")
	q.Set("apikey", ec.etherscanApiKey)
	url.RawQuery = q.Encode()

	u := url.String()
	utils.Logger.Infof("url => %s", u)

	client := http.Client{}
	resp, err := client.Get(u)
	if err != nil {
		utils.Logger.Errorw("failed to get transaction info", err)
		return nil, &werr.BlockchainError{}
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.Errorw("failed to parse response", err)
		return nil, &werr.BlockchainError{}
	}

	body := &GetTransactionsResp{}
	json.Unmarshal(bytes, body)

	unique := getUniqueInteractions(body.Result)

	result := &chains.Transactions{
		Count:  len(body.Result),
		Unique: unique,
	}

	return result, nil
}

func getUniqueInteractions(txs []TransactionsResp) int {
	m := make(map[string]int)

	for _, tx := range txs {
		if val, ok := m[tx.From]; ok {
			m[tx.From] = val + 1
		} else {
			m[tx.From] = 1
		}
	}

	return len(m)
}
