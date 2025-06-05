package wal

import (
	"github.com/rdimidov/kvstore/pkg/concurrency"
)

type FutureError = concurrency.Future[error]
type PromiseError = concurrency.Promise[error]

type entry struct {
	data    string
	promise *PromiseError
}

func newEntry(s string) entry {
	return entry{
		data:    s,
		promise: concurrency.NewPromise[error](),
	}
}

func (e *entry) FutureResponse() FutureError {
	return e.promise.GetFuture()
}

func (e *entry) SetResponse(err error) {
	e.promise.Set(err)
}
