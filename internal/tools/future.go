package tools

type Future struct {
	result <-chan interface{}
}

func NewFuture(result <-chan interface{}) *Future {
	future := &Future{
		result: result,
	}

	return future
}

func (f *Future) Get() interface{} {
	return <-f.result
}
