package storage

import (
	"context"
	"sync"

	"github.com/rdimidov/kvstore/internal/domain"
)

type Memory struct {
	mu sync.RWMutex
	hm map[string]domain.Entry
}

func NewMemory() *Memory {
	return &Memory{
		hm: make(map[string]domain.Entry),
	}
}

func (m *Memory) Set(_ context.Context, key domain.Key, value domain.Value) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.hm[key.String()] = domain.NewEntryFromKV(key, value)
	return nil
}

func (m *Memory) Get(_ context.Context, key domain.Key) (*domain.Entry, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if entry, ok := m.hm[key.String()]; ok {
		return &entry, nil
	}
	return nil, domain.ErrKeyNotFound
}

func (m *Memory) Delete(_ context.Context, key domain.Key) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.hm, key.String())
	return nil
}
