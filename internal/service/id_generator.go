package service

import (
	"sync"
	"sync/atomic"
)

var once sync.Once

type IDGenerator struct {
	counter atomic.Int64
}

func NewIDGenerator() *IDGenerator {
	return &IDGenerator{}
}

func (g *IDGenerator) SetInitValue(inintValue int64) {
	once.Do(func() {
		g.counter.Store(inintValue)
	})
}

func (g *IDGenerator) Generate() int64 {
	old := g.counter.Load()
	for !g.counter.CompareAndSwap(old, old+1) {
		old = g.counter.Load()
	}

	return old + 1
}
