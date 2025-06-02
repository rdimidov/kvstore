package wal

import (
	"github.com/rdimidov/kvstore/pkg/concurrency"
)

type entry struct {
	data    string
	promise *concurrency.PromiseError
}

func newEntry(s string) entry {
	return entry{
		data:    s,
		promise: concurrency.NewPromise[error](),
	}
}

func (e *entry) FutureResponse() concurrency.FutureError {
	return e.promise.GetFuture()
}

func (e *entry) SetResponse(err error) {
	e.promise.Set(err)
}
