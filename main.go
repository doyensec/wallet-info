package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/doyensec/wallet-info/chains/ethereum"
	"github.com/doyensec/wallet-info/config"
	"github.com/doyensec/wallet-info/dapp_file"
	"github.com/doyensec/wallet-info/domain"
	"github.com/doyensec/wallet-info/tls"
	"github.com/doyensec/wallet-info/utils"
	"github.com/gin-gonic/gin"
)

type HostResp struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`

	Domain              string    `json:"domain"`
	IsTlsValid          bool      `json:"tls"`
	TlsIssuedOn         time.Time `json:"tls_issued_on"`
	TlsExpiresOn        time.Time `json:"tls_expires_on"`
	DomainRecordCreated time.Time `json:"domain_created_on"`
	DomainRecordUpdated time.Time `json:"domain_updated_on"`
	DomainRecordExpires time.Time `json:"domain_expired_on"`

	IsSignatureValid bool `json:"valid_signature"`
}

type ContractResp struct {
	IsContract                   bool      `json:"is_contract"`
	ContractAddress              string    `json:"contract_address"`
	ContractDeployer             string    `json:"contract_deployer"`
	ContractDeployedOn           time.Time `json:"contract_deployed_on"`
	ContractTxCount              int       `json:"contract_tx_count"`
	ContractUniqueInteractions   int       `json:"contract_unique_tx"`
	IsContractSourceCodeVerified bool      `json:"verified_source"`

	IsSignatureValid bool `json:"valid_signature"`
}

func handleError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"msg": err.Error(),
	})
}

func getHostInfo(url string) (*HostResp, error) {
	host := utils.GetHost(url)

	dappFile, err := dapp_file.Get(host)
	if err != nil {
		utils.Logger.Infow("failed to get dapp file", "err", err)
		return nil, err
	}

	ethClient := ethereum.BuildClient(config.ActiveConfig)
	deployer, err := ethClient.GetContractDeployer(dappFile.ContractAddress)
	if err != nil {
		utils.Logger.Infow("failed to get contract deployer", "err", err)
		return nil, err
	}

	signatureValid, err := dappFile.ValidateHostOnly(host, deployer)
	if err != nil {
		utils.Logger.Infow("failed to validate dapp file", "err", err)
		return nil, err
	}

	tlsInfo, err := tls.GetInfo(dappFile.Domain)
	if err != nil {
		utils.Logger.Infow("failed to get tls info", "err", err)
		return nil, err
	}

	domainInfo, err := domain.GetDomainRecordInfo(host)
	if err != nil {
		utils.Logger.Infow("failed to get domain record info", "err", err)
		return nil, err
	}

	resp := &HostResp{
		Domain:           dappFile.Domain,
		IsSignatureValid: signatureValid,
		Name:             dappFile.Name,
		Timestamp:        dappFile.Timestamp,

		IsTlsValid:          tlsInfo.IsValid(),
		TlsIssuedOn:         tlsInfo.IssuedOn,
		TlsExpiresOn:        tlsInfo.ExpiresOn,
		DomainRecordCreated: domainInfo.Created,
		DomainRecordUpdated: domainInfo.Updated,
		DomainRecordExpires: domainInfo.Expires,
	}

	return resp, nil
}

func getContractInfo(url, contract string) (*ContractResp, error) {
	host := utils.GetHost(url)

	dappFile, err := dapp_file.Get(host)
	if err != nil {
		utils.Logger.Infow("failed to get dapp file", "err", err)
		return nil, err
	}

	ethClient := ethereum.BuildClient(config.ActiveConfig)
	deployer, err := ethClient.GetContractDeployer(contract)
	if err != nil {
		utils.Logger.Infow("failed to get contract deployer", "err", err)
		return nil, err
	}

	isContract := deployer != ""

	var signatureValid = false
	if isContract {
		signatureValid, err = dappFile.Validate(host, contract, deployer)
		if err != nil {
			utils.Logger.Infow("failed to validate dapp file", "err", err)
			return nil, err
		}
	}

	txs, err := ethClient.GetTransactions(contract)
	if err != nil {
		utils.Logger.Infow("failed to get contract transactions", "err", err)
		return nil, err
	}

	verified, err := ethClient.IsVerified(contract)
	if err != nil {
		utils.Logger.Infow("failed to get contract verified status", "err", err)
		return nil, err
	}

	isVerified := len(verified.Result[0].SourceCode) > 0

	resp := &ContractResp{
		IsContract:                   isContract,
		ContractAddress:              contract,
		ContractDeployer:             deployer,
		ContractDeployedOn:           time.Now(),
		ContractTxCount:              txs.Count,
		ContractUniqueInteractions:   txs.Unique,
		IsContractSourceCodeVerified: isVerified,
		IsSignatureValid:             signatureValid,
	}

	return resp, nil
}

func runServer() {
	router := gin.Default()

	router.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hi.")
	})

	router.GET("/host", func(ctx *gin.Context) {
		url, found := ctx.GetQuery("url")
		if !found {
			ctx.String(http.StatusBadRequest, "missing 'url' parameter")
			return
		}

		resp, err := getHostInfo(url)
		if err != nil {
			handleError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	})

	router.GET("/contract", func(ctx *gin.Context) {
		url, found := ctx.GetQuery("url")
		if !found {
			ctx.String(http.StatusBadRequest, "missing 'url' parameter")
			return
		}

		address, found := ctx.GetQuery("address")
		if !found {
			ctx.String(http.StatusBadRequest, "missing 'address' parameter")
			return
		}

		resp, err := getContractInfo(url, address)
		if err != nil {
			handleError(ctx, err)
			return
		}

		ctx.JSON(http.StatusOK, resp)
	})

	port := 8000
	utils.Logger.Infow("listener started...", "port", port)
	router.Run(fmt.Sprintf("127.0.0.1:%d", port))
}

func main() {
	utils.InitializeLogger()

	err := config.LoadConfig()
	if err != nil {
		utils.Logger.Errorw("error loading config", "err", err)
	}

	runServer()
}
