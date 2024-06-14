package tools

import "sync/atomic"

type Promise struct {
	result   chan interface{}
	promised int32
}

func NewPromise() *Promise {
	return &Promise{
		result: make(chan interface{}),
	}
}

func (p *Promise) GetFuture() *Future {
	return NewFuture(p.result)
}

func (p *Promise) Set(value interface{}) {
	if ok := atomic.CompareAndSwapInt32(&p.promised, 0, 1); ok {
		p.result <- value
		close(p.result)
	}
}
