package testing

import (
	"crypto/ecdsa"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func GenSignature(data string) {
	privateKey, err := crypto.HexToECDSA("fad9c8855b740a0b7ed4c221dbad0f33a83a49cad6b3fe8d5817ac83d38b6a19")
	if err != nil {
		log.Fatal(err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	fmt.Printf("pub => %s\n", hexutil.Encode(publicKeyBytes))

	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	fmt.Printf("address => %s\n", address)

	msg := []byte(data)
	msgHash := crypto.Keccak256Hash(msg)
	fmt.Printf("msg hash => %s\n", msgHash.Hex())

	signature, err := crypto.Sign(msgHash.Bytes(), privateKey)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("signature => %s\n", hexutil.Encode(signature))

	sigPublicKey, err := crypto.Ecrecover(msgHash.Bytes(), signature)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("sig pub => %s\n", hexutil.Encode(sigPublicKey))
}
