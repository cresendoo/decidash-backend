package xaptos

import (
	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/aptos-labs/aptos-go-sdk/crypto"
)

type SimpleTransaction struct {
	RawTxn   *aptos.RawTransaction
	FeePayer *aptos.AccountAddress
}

func (txn *SimpleTransaction) MarshalBCS(ser *bcs.Serializer) {
	ser.Struct(txn.RawTxn)
	if txn.FeePayer == nil {
		ser.Bool(false)
	} else {
		ser.Bool(true)
		ser.Struct(txn.FeePayer)
	}
}

func (txn *SimpleTransaction) UnmarshalBCS(des *bcs.Deserializer) {
	txn.RawTxn = &aptos.RawTransaction{}
	des.Struct(txn.RawTxn)
	feePayerPresent := des.Bool()
	if feePayerPresent {
		txn.FeePayer = &aptos.AccountAddress{}
		des.Struct(txn.FeePayer)
	}
}

func (txn *SimpleTransaction) SigningMessage() (message []byte, err error) {
	txnBytes, err := bcs.Serialize(&aptos.RawTransactionWithData{
		Variant: aptos.MultiAgentWithFeePayerRawTransactionWithDataVariant,
		Inner: &aptos.MultiAgentWithFeePayerRawTransactionWithData{
			RawTxn:           txn.RawTxn,
			FeePayer:         txn.FeePayer,
			SecondarySigners: []aptos.AccountAddress{},
		},
	})
	if err != nil {
		return
	}
	prehash := aptos.RawTransactionWithDataPrehash()
	message = make([]byte, len(prehash)+len(txnBytes))
	copy(message, prehash)
	copy(message[len(prehash):], txnBytes)
	return message, nil
}

func (txn *SimpleTransaction) Sign(signer crypto.Signer) (*crypto.AccountAuthenticator, error) {
	message, err := txn.SigningMessage()
	if err != nil {
		return nil, err
	}
	return signer.Sign(message)
}
