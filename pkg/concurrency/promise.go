package concurrency

import "sync"

type PromiseError = Promise[error]

type Promise[T any] struct {
	result chan T

	mu       sync.Mutex
	promised bool
}

func NewPromise[T any]() *Promise[T] {
	return &Promise[T]{
		result: make(chan T, 1),
	}
}

func (p *Promise[T]) Set(value T) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.promised {
		return
	}

	p.promised = true

	p.result <- value
	close(p.result)
}

func (p *Promise[T]) GetFuture() Future[T] {
	return NewFuture(p.result)
}
