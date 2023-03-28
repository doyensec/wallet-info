package dapp_file

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/doyensec/safeurl"
	"github.com/doyensec/wallet-info/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/status-im/keycard-go/hexutils"
)

type DappFile struct {
	Name            string    `json:"name"`
	Domain          string    `json:"domain"`
	ContractAddress string    `json:"contract"`
	Timestamp       time.Time `json:"timestamp"`
	Signature       string    `json:"signature,omitempty"`
}

func (dapp *DappFile) ValidateHostOnly(host, deployer string) (bool, error) {
	signatureValid := dapp.isSignatureValid(deployer)

	err := dapp.validateHost(host)
	if err != nil {
		return signatureValid, err
	}

	return signatureValid, nil
}

func (dapp *DappFile) Validate(host, contract, deployer string) (bool, error) {
	signatureValid := dapp.isSignatureValid(deployer)

	err := dapp.validateHost(host)
	if err != nil {
		return signatureValid, err
	}

	err = dapp.validateContractAddress(contract)
	if err != nil {
		return signatureValid, err
	}

	return signatureValid, nil
}

func (dapp *DappFile) validateHost(requestedHost string) error {
	if dapp.Domain == "" {
		utils.Logger.Errorw("empty dapp domain")
		return errors.New("invalid domain")
	}

	if requestedHost != dapp.Domain {
		utils.Logger.Errorw("dapp file domain mismatch", "requested", requestedHost, "dapp_file", dapp.Domain)
		return errors.New("domain mismatch")
	}

	return nil
}

func (dapp *DappFile) validateContractAddress(contract string) error {
	if dapp.ContractAddress == "" {
		utils.Logger.Errorw("empty dapp contract address")
		return errors.New("invalid contract address")
	}

	if !strings.EqualFold(contract, dapp.ContractAddress) {
		utils.Logger.Errorw("dapp contract address mismatch", "requested", contract, "dapp_file", dapp.ContractAddress)
		return errors.New("contract address mismatch")
	}

	return nil
}

func (dapp *DappFile) toSigningBytes() []byte {
	copy := &DappFile{
		Name:            dapp.Name,
		Domain:          dapp.Domain,
		ContractAddress: dapp.ContractAddress,
		Timestamp:       dapp.Timestamp,
		// omit signature
	}

	bytes, _ := json.Marshal(copy)
	return bytes
}

func (dapp *DappFile) isSignatureValid(deployer string) bool {
	msg := dapp.toSigningBytes()
	msgHash := crypto.Keccak256Hash(msg)

	signaturePublicKey, err := crypto.Ecrecover(msgHash.Bytes(), hexutils.HexToBytes(dapp.Signature[2:]))
	if err != nil {
		log.Fatal(err)
	}

	signatureAddress := publicKeyBytesToAddress(signaturePublicKey).String()

	return deployer == signatureAddress
}

func Get(domain string) (*DappFile, error) {
	dappFileUrl := fmt.Sprintf("%s/.well-known/dapp_file", domain)

	utils.Logger.Infow("reading dapp file", "url", dappFileUrl)

	client := getClient()

	resp, err := client.Get(dappFileUrl)
	if err != nil {
		utils.Logger.Errorw("failed to get dapp file", err)
		return nil, errors.New("failed to process `dapp_file`")
	}

	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.Errorw("failed to read response body", err)
		return nil, errors.New("failed to process `dapp_file`")
	}

	body := &DappFile{}
	err = json.Unmarshal(bytes, &body)
	if err != nil {
		utils.Logger.Errorw("failed to unmarshal dapp file", err)
		return nil, errors.New("failed to process `dapp_file`")
	}

	return body, nil
}

func getClient() *safeurl.WrappedClient {
	config := safeurl.GetConfigBuilder().
		SetAllowedSchemes("https").
		Build()
	return safeurl.Client(config)
}

func publicKeyBytesToAddress(publicKey []byte) common.Address {
	buf := crypto.Keccak256Hash(publicKey[1:]) // remove ec prefix 04
	address := buf[12:]
	return common.HexToAddress(hex.EncodeToString(address))
}
