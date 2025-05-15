package storage

import (
	"context"
	"strconv"
	"sync"
	"testing"

	"github.com/rdimidov/kvstore/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestMemory_SetGet(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mem := NewMemory()
	key := domain.Key("foo")
	val := domain.Value("bar")

	// Set
	err := mem.Set(ctx, key, val)
	assert.NoError(t, err)

	// Get existing
	entry, err := mem.Get(ctx, key)
	assert.NoError(t, err)
	assert.Equal(t, key, entry.Key)
	assert.Equal(t, val, entry.Value)
}

func TestMemory_GetNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mem := NewMemory()
	key := domain.Key("missing")

	// Get missing key
	entry, err := mem.Get(ctx, key)
	assert.Nil(t, entry)
	assert.ErrorIs(t, err, domain.ErrKeyNotFound)
}

func TestMemory_Delete(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mem := NewMemory()
	key := domain.Key("to_delete")
	val := domain.Value("value")

	// Set then delete
	assert.NoError(t, mem.Set(ctx, key, val))
	assert.NoError(t, mem.Delete(ctx, key))

	// Ensure it's gone
	entry, err := mem.Get(ctx, key)
	assert.Nil(t, entry)
	assert.ErrorIs(t, err, domain.ErrKeyNotFound)
}

func TestMemory_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mem := NewMemory()

	const n = 10000
	var wg sync.WaitGroup
	wg.Add(n)

	// Concurrently set keys
	for i := range n {
		go func(i int) {
			defer wg.Done()
			key := domain.Key("key-" + strconv.Itoa(i))
			val := domain.Value("val-" + strconv.Itoa(i))
			assert.NoError(t, mem.Set(ctx, key, val))
		}(i)
	}
	wg.Wait()

	// Concurrently read and delete keys
	wg.Add(n)
	for i := range n {
		go func(i int) {
			defer wg.Done()
			key := domain.Key("key-" + strconv.Itoa(i))

			// Read
			entry, err := mem.Get(ctx, key)
			assert.NoError(t, err)
			assert.Equal(t, key, entry.Key)

			// Delete
			assert.NoError(t, mem.Delete(ctx, key))

			// Then Get should error
			_, err = mem.Get(ctx, key)
			assert.ErrorIs(t, err, domain.ErrKeyNotFound)
		}(i)
	}
	wg.Wait()
}
