package node

import (
	"encoding/hex"
	"sync"

	"github.com/koshkaj/bloq/proto"
	"github.com/koshkaj/bloq/types"
)

type Mempool struct {
	// Binary tree would be more performant to search
	mu  sync.RWMutex
	txx map[string]*proto.Transaction
}

func (pool *Mempool) Clear() []*proto.Transaction {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	txx := make([]*proto.Transaction, len(pool.txx))
	iterate := 0
	for k, v := range pool.txx {
		delete(pool.txx, k)
		txx[iterate] = v
		iterate++
	}
	return txx
}

func (pool *Mempool) Has(tx *proto.Transaction) bool {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	hash := hex.EncodeToString(types.HashTransaction(tx))
	_, ok := pool.txx[hash]
	return ok
}

func (pool *Mempool) Add(tx *proto.Transaction) {
	pool.mu.Lock()
	defer pool.mu.Unlock()
	hash := hex.EncodeToString(types.HashTransaction(tx))
	pool.txx[hash] = tx
}

func (pool *Mempool) Len() int {
	pool.mu.RLock()
	defer pool.mu.RUnlock()
	return len(pool.txx)
}

func NewMempool() *Mempool {
	return &Mempool{
		txx: make(map[string]*proto.Transaction),
	}
}
