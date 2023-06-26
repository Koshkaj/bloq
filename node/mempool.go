package node

import (
	"container/list"
	"encoding/hex"

	"github.com/koshkaj/bloq/proto"
	"github.com/koshkaj/bloq/types"
)

type Mempool struct {
	txx *list.List
}

func (pool *Mempool) Has(tx *proto.Transaction) bool {
	for e := pool.txx.Front(); e != nil; e = e.Next() {
		hash := e.Value.(string)
		if hash == hex.EncodeToString(types.HashTransaction(tx)) {
			return true
		}
	}
	return false
}

func (pool *Mempool) Add(tx *proto.Transaction) {
	pool.txx.PushBack(hex.EncodeToString(types.HashTransaction(tx)))
}

func (pool *Mempool) Length() int {
	return pool.txx.Len()
}

func NewMempool() *Mempool {
	return &Mempool{
		txx: list.New(),
	}
}
