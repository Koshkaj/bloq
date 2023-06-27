package node

import (
	"testing"

	"github.com/koshkaj/bloq/crypto"
	"github.com/koshkaj/bloq/proto"
	"github.com/koshkaj/bloq/types"
	"github.com/koshkaj/bloq/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randomBlock(t *testing.T, chain *Chain) *proto.Block {
	privKey := crypto.GeneratePrivateKey()
	b := util.RandomBlock()
	prevBlock, err := chain.GetBlockByHeight(chain.Height())
	require.Nil(t, err)
	b.Header.PrevHash = types.HashBlock(prevBlock)
	types.SignBlock(privKey, b)
	return b
}

func TestNewChain(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	require.Equal(t, 0, chain.Height())
	_, err := chain.GetBlockByHeight(0)
	require.Nil(t, err)
}

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	for i := 0; i < 100; i++ {
		block := randomBlock(t, chain)
		blockHash := types.HashBlock(block)

		require.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i + 1)
		require.Nil(t, err)
		require.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
	for i := 0; i < 100; i++ {
		b := randomBlock(t, chain)
		require.Nil(t, chain.AddBlock(b))
		require.Equal(t, chain.Height(), i+1)
	}
}

func TestAddBlockWithTx(t *testing.T) {
	var (
		privKey   = crypto.NewPrivateKeyFromSeedStr(seed)
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
		block     = randomBlock(t, chain)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)
	ftt, err := chain.txStore.Get("3621faa661c206054d614eaa103687095e6dc8722bb31a03143bd3deeb0613b0") // Genisis Transaction
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(ftt),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  100,
			Address: recipient,
		},
		{
			Amount:  8788,
			Address: privKey.Public().Address().Bytes(), // nashti racaa gavigzavnot ukan chven addressze
		},
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}
	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)
	require.Nil(t, chain.AddBlock(block))
}

func TestAddBlockWithTxLowFunds(t *testing.T) {
	var (
		privKey   = crypto.NewPrivateKeyFromSeedStr(seed)
		chain     = NewChain(NewMemoryBlockStore(), NewMemoryTXStore())
		block     = randomBlock(t, chain)
		recipient = crypto.GeneratePrivateKey().Public().Address().Bytes()
	)
	ftt, err := chain.txStore.Get("3621faa661c206054d614eaa103687095e6dc8722bb31a03143bd3deeb0613b0") // Genisis Transaction
	assert.Nil(t, err)

	inputs := []*proto.TxInput{
		{
			PrevTxHash:   types.HashTransaction(ftt),
			PrevOutIndex: 0,
			PublicKey:    privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  9999,
			Address: recipient,
		},
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}
	sig := types.SignTransaction(privKey, tx)
	tx.Inputs[0].Signature = sig.Bytes()

	block.Transactions = append(block.Transactions, tx)
	require.NotNil(t, chain.AddBlock(block))
}
