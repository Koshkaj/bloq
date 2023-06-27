package node

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/koshkaj/bloq/crypto"
	"github.com/koshkaj/bloq/proto"
	"github.com/koshkaj/bloq/types"
)

const seed = "13be19fc5de106d87f9deaec7de204bd5b3a36bb50d66f81fdd5d1482dbeab6e"

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	return &HeaderList{
		headers: []*proto.Header{},
	}
}

func (list *HeaderList) Height() int {
	return list.Len() - 1
}

func (list *HeaderList) Get(index int) *proto.Header {
	if index > list.Height() {
		log.Fatal("index too high")
	}
	return list.headers[index]
}

func (list *HeaderList) Add(h *proto.Header) {
	list.headers = append(list.headers, h)
}

func (list *HeaderList) Len() int {
	return len(list.headers)
}

type UTXO struct {
	Hash     string
	OutIndex int
	Amount   int64
	Spent    bool
}

type Chain struct {
	txStore    TXStorer
	blockStore BlockStorer
	utxoStore  UTXOStorer
	headers    *HeaderList
}

func NewChain(bs BlockStorer, txStore TXStorer) *Chain {
	chain := &Chain{
		blockStore: bs,
		headers:    NewHeaderList(),
		utxoStore:  NewMemoryUTXOStore(),
		txStore:    txStore,
	}
	chain.addBlock(createGenesisBlock())
	return chain
}

func (c *Chain) Height() int {
	return c.headers.Height()

}

func (c *Chain) AddBlock(b *proto.Block) error {
	if err := c.ValidateBlock(b); err != nil {
		return err
	}
	return c.addBlock(b)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	if c.Height() < height {
		return nil, fmt.Errorf("given height [%d] is invalid, current height [%d]", height, c.Height())
	}
	header := c.headers.Get(height)
	hash := types.HashHeader(header)
	return c.GetBlockByHash(hash)
}

func (c *Chain) ValidateBlock(b *proto.Block) error {
	if !types.VerifyBlock(b) {
		return fmt.Errorf("invalid block signature")
	}

	currentBlock, err := c.GetBlockByHeight(c.Height())
	if err != nil {
		return err
	}
	currentBlockHash := types.HashBlock(currentBlock)
	if !bytes.Equal(currentBlockHash, b.Header.PrevHash) {
		return fmt.Errorf("invalid previous block hash")
	}
	for _, tx := range b.Transactions {
		if err := c.ValidateTransaction(tx); err != nil {
			return err
		}
	}
	return nil
}

func (c *Chain) ValidateTransaction(tx *proto.Transaction) error {
	if !types.VerifyTransaction(tx) {
		return fmt.Errorf("invalid tx signature")
	}
	var (
		nIutputs   = len(tx.Inputs)
		hash       = types.HashTransaction(tx)
		sumInputs  = 0
		sumOutputs = 0
	)

	for i := 0; i < nIutputs; i++ {
		prevHash := hex.EncodeToString(tx.Inputs[i].PrevTxHash)
		key := fmt.Sprintf("%s_%d", prevHash, i)
		utxo, err := c.utxoStore.Get(key)
		sumInputs += int(utxo.Amount)
		if err != nil {
			return err
		}
		if utxo.Spent {
			return fmt.Errorf("input %d of tx %s is already spent", i, hash)
		}
	}
	for _, output := range tx.Outputs {
		sumOutputs += int(output.Amount)
	}
	if sumInputs < sumOutputs {
		return fmt.Errorf("insufficient balance")
	}

	return nil
}

func (c *Chain) addBlock(b *proto.Block) error {
	c.headers.Add(b.Header)
	for _, tx := range b.Transactions {
		if err := c.txStore.Put(tx); err != nil {
			return err
		}
		for idx, output := range tx.Outputs {
			hash := hex.EncodeToString(types.HashTransaction(tx))
			utxo := &UTXO{
				Hash:     hash,
				Amount:   output.Amount,
				OutIndex: idx,
				Spent:    false,
			}
			if err := c.utxoStore.Put(utxo); err != nil {
				return err
			}
		}
		for _, input := range tx.Inputs {
			key := fmt.Sprintf("%s_%d", hex.EncodeToString(input.PrevTxHash), input.PrevOutIndex)
			utxo, err := c.utxoStore.Get(key)
			if err != nil {
				return err
			}
			utxo.Spent = true
			if err := c.utxoStore.Put(utxo); err != nil {
				return err
			}
		}
	}

	return c.blockStore.Put(b)
}

func createGenesisBlock() *proto.Block {
	privKey := crypto.NewPrivateKeyFromSeedStr(seed)
	block := &proto.Block{
		Header: &proto.Header{
			Version: 1,
		},
	}
	tx := &proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{},
		Outputs: []*proto.TxOutput{
			{
				Amount:  8888,
				Address: privKey.Public().Address().Bytes(),
			},
		},
	}
	block.Transactions = append(block.Transactions, tx)
	types.SignBlock(privKey, block)

	return block
}
