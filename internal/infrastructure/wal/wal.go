package wal

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/rdimidov/kvstore/internal/domain"
)

const (
	defaultWriteDir     = "./wal"
	defaultSegentSizeMB = 10
	defaultBatchSize    = 10
	defaultFlushTimeout = 10 * time.Millisecond
)

type writer interface {
	Write([]entry)
}

type config interface {
	WALBatchSize() int
	WALBatchFlushTimeout() time.Duration
	WALDirName() string
	WALMaxSegmentSize() int
}

type interpreter interface {
	Execute(ctx context.Context, raw string) (*domain.Entry, error)
}

type WAL struct {
	writer      writer
	reader      reader
	interpreter interpreter

	batchLimit int
	timeout    time.Duration

	readyCh chan []entry
	mu      sync.Mutex
	batch   []entry
}

func New(ctx context.Context, config config, interpreter interpreter) (*WAL, error) {
	if interpreter == nil {
		return nil, errors.New("interpreter is nil")
	}

	dirname := config.WALDirName()
	if dirname == "" {
		dirname = defaultWriteDir
	}

	mssMB := config.WALMaxSegmentSize()
	if mssMB == 0 {
		mssMB = defaultSegentSizeMB
	}

	writer, err := newRotatingWalWriter(dirname, mssMB*1024*1024)
	if err != nil {
		return nil, err
	}

	reader := NewReader(dirname)

	batch := config.WALBatchSize()
	if batch == 0 {
		batch = defaultBatchSize
	}

	timeout := config.WALBatchFlushTimeout()
	if timeout == 0 {
		timeout = defaultFlushTimeout
	}

	wal := &WAL{
		batchLimit:  batch,
		timeout:     defaultFlushTimeout,
		readyCh:     make(chan []entry, 1),
		writer:      writer,
		reader:      reader,
		interpreter: interpreter,
	}
	wal.start(ctx)
	return wal, nil
}

func (w *WAL) start(ctx context.Context) {
	go func() {
		ticker := time.NewTicker(w.timeout)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.dumpBatch()
				return

			case <-ticker.C:
				w.dumpBatch()

			case batch := <-w.readyCh:
				go w.writer.Write(batch)
				ticker.Reset(w.timeout)

			}
		}
	}()
}

func (w *WAL) WriteSet(key domain.Key, value domain.Value) error {
	fut := w.processInput(fmt.Sprintf("SET %s %s", key, value))
	return fut.Get()
}

func (w *WAL) WriteDel(key domain.Key) error {
	fut := w.processInput(fmt.Sprintf("DEL %s", key))
	return fut.Get()
}

func (w *WAL) processInput(input string) FutureError {
	entry := newEntry(input)

	w.mu.Lock()
	w.batch = append(w.batch, entry)

	if len(w.batch) == w.batchLimit {
		w.readyCh <- w.batch
		w.batch = nil
	}
	w.mu.Unlock()

	return entry.FutureResponse()
}

func (w *WAL) dumpBatch() {
	var batch []entry

	w.mu.Lock()
	batch = w.batch
	w.batch = nil
	w.mu.Unlock()

	if len(batch) != 0 {
		w.writer.Write(batch)
	}
}

func (w *WAL) Recover(ctx context.Context) error {
	commands, err := w.reader.Read()
	if err != nil {
		return err
	}
	for _, l := range commands {
		_, err := w.interpreter.Execute(ctx, l)
		if err != nil {
			return err
		}
	}
	return nil
}

type Noop struct{}

func (w *Noop) WriteSet(domain.Key, domain.Value) error { return nil }
func (w *Noop) WriteDel(domain.Key) error               { return nil }
func (w *Noop) Recover(context.Context) error           { return nil }
