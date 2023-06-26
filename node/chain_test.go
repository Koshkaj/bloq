package node

import (
	"encoding/hex"
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

	inputs := []*proto.TxInput{
		{
			PublicKey: privKey.Public().Bytes(),
		},
	}
	outputs := []*proto.TxOutput{
		{
			Amount:  100,
			Address: recipient,
		},
	}

	tx := &proto.Transaction{
		Version: 1,
		Inputs:  inputs,
		Outputs: outputs,
	}
	block.Transactions = append(block.Transactions, tx)
	require.Nil(t, chain.AddBlock(block))
	txHash := hex.EncodeToString(types.HashTransaction(tx))
	fetchedTx, err := chain.txStore.Get(txHash)
	assert.Nil(t, err)
	assert.Equal(t, tx, fetchedTx)
}
