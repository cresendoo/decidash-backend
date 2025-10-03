package xaptos

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
)

func MustAccountFromEd25519PrivateKey(strPK string) *aptos.Account {
	privateKey := &crypto.Ed25519PrivateKey{}
	pk, err := crypto.FormatPrivateKey(strPK, crypto.PrivateKeyVariantEd25519)
	if err != nil {
		panic("Failed to format private key:" + err.Error())
	}
	if err := privateKey.FromHex(pk); err != nil {
		panic("Failed to parse private key:" + err.Error())
	}
	account, err := aptos.NewAccountFromSigner(privateKey)
	if err != nil {
		panic("Failed to create account:" + err.Error())
	}
	return account
}

func AccountFromEd25519PrivateKey(strPK string) (*aptos.Account, error) {
	privateKey := &crypto.Ed25519PrivateKey{}
	pk, err := crypto.FormatPrivateKey(strPK, crypto.PrivateKeyVariantEd25519)
	if err != nil {
		return nil, err
	}
	if err := privateKey.FromHex(pk); err != nil {
		return nil, err
	}
	account, err := aptos.NewAccountFromSigner(privateKey)
	if err != nil {
		return nil, err
	}
	return account, nil
}
