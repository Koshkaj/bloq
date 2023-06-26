package node

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/koshkaj/bloq/proto"
	"github.com/koshkaj/bloq/types"
)

type UTXOStorer interface {
	Put(*UTXO) error
	Get(string) (*UTXO, error)
}

type MemoryUTXOStore struct {
	mu   sync.RWMutex
	data map[string]*UTXO
}

func (s *MemoryUTXOStore) Get(hash string) (*UTXO, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	utxo, ok := s.data[hash]
	if !ok {
		return nil, fmt.Errorf("could not find utxo with hash %s", hash)
	}
	return utxo, nil
}

func (s *MemoryUTXOStore) Put(utxo *UTXO) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	key := fmt.Sprintf("%s_%d", utxo.Hash, utxo.OutIndex)
	s.data[key] = utxo
	return nil
}
func NewMemoryUTXOStore() *MemoryUTXOStore {
	return &MemoryUTXOStore{
		data: make(map[string]*UTXO),
	}
}

type TXStorer interface {
	Put(*proto.Transaction) error
	Get(string) (*proto.Transaction, error)
}

type MemoryTXStore struct {
	mu  sync.RWMutex
	txx map[string]*proto.Transaction
}

func NewMemoryTXStore() *MemoryTXStore {
	return &MemoryTXStore{
		txx: make(map[string]*proto.Transaction),
	}
}

func (s *MemoryTXStore) Put(tx *proto.Transaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hash := hex.EncodeToString(types.HashTransaction(tx))
	s.txx[hash] = tx
	return nil
}

func (s *MemoryTXStore) Get(hash string) (*proto.Transaction, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tx, ok := s.txx[hash]
	if !ok {
		return nil, fmt.Errorf("could not find tx with hash %s", hash)
	}
	return tx, nil
}

type BlockStorer interface {
	Put(*proto.Block) error
	Get(string) (*proto.Block, error)
}

type MemoryBlockStore struct {
	mu     sync.RWMutex
	blocks map[string]*proto.Block
}

func NewMemoryBlockStore() *MemoryBlockStore {
	return &MemoryBlockStore{
		blocks: make(map[string]*proto.Block),
	}
}

func (s *MemoryBlockStore) Put(b *proto.Block) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	hash := hex.EncodeToString(types.HashBlock(b))
	s.blocks[hash] = b
	return nil
}

func (s *MemoryBlockStore) Get(hash string) (*proto.Block, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	block, ok := s.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash [%s] does not exist", hash)
	}
	return block, nil
}
