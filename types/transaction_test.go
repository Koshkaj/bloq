package types

import (
	"testing"

	"github.com/koshkaj/bloq/crypto"
	"github.com/koshkaj/bloq/proto"
	"github.com/koshkaj/bloq/util"
	"github.com/stretchr/testify/assert"
)

// balance : 100 coins
// send 5 coins
// balance becomes 95, recipient receives 5 coins

func TestNewTransaction(t *testing.T) {
	fromPrivKey := crypto.GeneratePrivateKey()
	fromAddress := fromPrivKey.Public().Address().Bytes()

	toPrivKey := crypto.GeneratePrivateKey()
	toAddress := toPrivKey.Public().Address().Bytes()

	input := &proto.TxInput{
		PrevTxHash:   util.RandomHash(),
		PrevOutIndex: 0,
		PublicKey:    fromPrivKey.Public().Bytes(),
	}
	output_1 := &proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}
	output_2 := &proto.TxOutput{
		Amount:  95,
		Address: fromAddress,
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{input},
		Outputs: []*proto.TxOutput{output_1, output_2},
	}
	sig := SignTransaction(fromPrivKey, tx)
	input.Signature = sig.Bytes()
	assert.True(t, VerifyTransaction(tx))

}
