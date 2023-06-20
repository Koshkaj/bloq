package node

import (
	"testing"

	"github.com/koshkaj/bloq/types"
	"github.com/koshkaj/bloq/util"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	for i := 0; i < 100; i++ {
		var (
			block       = util.RandomBlock()
			hashedBlock = types.HashBlock(block)
		)
		assert.Nil(t, chain.AddBlock(block))

		fetchedBlock, err := chain.GetBlockByHash(hashedBlock)
		assert.Nil(t, err)
		assert.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i)
		assert.Nil(t, err)
		assert.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(NewMemoryBlockStore())
	for i := 0; i < 100; i++ {
		b := util.RandomBlock()
		assert.Nil(t, chain.AddBlock(b))
		assert.Equal(t, chain.Height(), i)
	}
}
