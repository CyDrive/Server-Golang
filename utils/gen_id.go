package utils

import (
	"math"
	"sync"
	"sync/atomic"
)

type IdGenerator struct {
	minId int64
	maxId int64

	currentId int64

	idRefMap sync.Map
}

func NewIdGenerator() *IdGenerator {
	return &IdGenerator{
		minId:     math.MinInt64,
		maxId:     math.MaxInt64,
		currentId: math.MinInt64,
	}
}

func (idGen *IdGenerator) SetMinId(minId int64) {
	atomic.StoreInt64(&idGen.minId, minId)

	for {
		oldId := atomic.LoadInt64(&idGen.currentId)
		if oldId < minId &&
				atomic.CompareAndSwapInt64(&idGen.currentId, oldId, minId) {
			break
		}
	}
}

func (idGen *IdGenerator) SetMaxId(maxId int64) {
	atomic.StoreInt64(&idGen.maxId, maxId)

	for {
		oldId := atomic.LoadInt64(&idGen.currentId)
		if oldId > maxId &&
				atomic.CompareAndSwapInt64(&idGen.currentId, oldId, idGen.minId) {
			break
		}
	}
}

func (idGen *IdGenerator) Next() int64 {
	return atomic.AddInt64(&idGen.currentId, 1) - 1
}

func (idGen *IdGenerator) NextAndRef() int64 {
	id := atomic.AddInt64(&idGen.currentId, 1) - 1
	idGen.Ref(id)
	return id
}

func (idGen *IdGenerator) Ref(id int64) {
	idGen.idRefMap.Store(id, true)
}

func (idGen *IdGenerator) UnRef(id int64) {
	idGen.idRefMap.Delete(id)
}
