package storage

import (
	"context"

	"github.com/rdimidov/kvstore/internal/domain"
)

type Memory struct {
	hm map[string]domain.Entry
}

func NewMemory() *Memory {
	return &Memory{
		hm: make(map[string]domain.Entry),
	}
}

func (m *Memory) Set(_ context.Context, key domain.Key, value domain.Value) error {
	m.hm[key.String()] = domain.NewEntryFromKV(key, value)
	return nil
}

func (m *Memory) Get(_ context.Context, key domain.Key) (*domain.Entry, error) {
	if entry, ok := m.hm[key.String()]; ok {
		return &entry, nil
	}
	return nil, domain.ErrKeyNotFound
}

func (m *Memory) Delete(ctx context.Context, key domain.Key) error {
	delete(m.hm, key.String())
	return nil
}
