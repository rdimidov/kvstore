package wal

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/rdimidov/kvstore/internal/domain"

	"github.com/stretchr/testify/assert"
)

type testConfig struct{}

func (testConfig) WALBatchSize() int                   { return 2 }
func (testConfig) WALBatchFlushTimeout() time.Duration { return 20 * time.Millisecond }
func (testConfig) WALDirName() string                  { return "./test_wal" }
func (testConfig) WALMaxSegmentSize() int              { return 1 } // MB

func cleanupTestDir(t *testing.T, path string) {
	t.Helper()
	err := os.RemoveAll(path)
	assert.NoError(t, err)
}

func TestWriteSetAndFlushOnBatchLimit(t *testing.T) {
	cfg := testConfig{}
	defer cleanupTestDir(t, cfg.WALDirName())
	w, err := New(context.Background(), cfg, newMockinterpreter(t))
	assert.NoError(t, err)

	key1, _ := domain.NewKey("foo")
	val1, _ := domain.NewValue("bar")
	err = w.WriteSet(key1, val1)

	assert.NoError(t, err)

	err = w.WriteSet("key", "val")
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond)

	reader := NewReader(cfg.WALDirName())
	lines, err := reader.Read()
	assert.NoError(t, err)
	assert.Contains(t, lines, "SET foo bar")
	assert.Contains(t, lines, "SET key val")
}

func TestWriteDelAndFlushOnTimeout(t *testing.T) {
	cfg := testConfig{}
	defer cleanupTestDir(t, cfg.WALDirName())
	w, err := New(context.Background(), cfg, newMockinterpreter(t))
	assert.NoError(t, err)

	err = w.WriteDel("somekey")
	assert.NoError(t, err)

	time.Sleep(50 * time.Millisecond) // flush on timeout

	reader := NewReader("./test_wal")
	lines, err := reader.Read()
	assert.NoError(t, err)
	assert.Contains(t, lines, "DEL somekey")
}

func TestRecoverExecutesCommands(t *testing.T) {
	cfg := testConfig{}
	defer cleanupTestDir(t, cfg.WALDirName())

	_ = os.MkdirAll(cfg.WALDirName(), 0o755)
	f, err := os.Create(filepath.Join(cfg.WALDirName(), "manual.wal"))
	assert.NoError(t, err)
	_, _ = f.WriteString("SET foo bar\nDEL foo\n")
	f.Close()

	ctx := context.Background()

	mockInterpreter := newMockinterpreter(t)
	mockInterpreter.On("Execute", ctx, "SET foo bar").Return(&domain.Entry{}, nil).Once()
	mockInterpreter.On("Execute", ctx, "DEL foo").Return(&domain.Entry{}, nil).Once()

	w := WAL{
		reader:      NewReader(cfg.WALDirName()),
		interpreter: mockInterpreter,
	}
	err = w.Recover(ctx)
	assert.NoError(t, err)
	mockInterpreter.AssertExpectations(t)
}

func TestRecoverFailsIfInterpreterFails(t *testing.T) {
	cfg := testConfig{}
	defer cleanupTestDir(t, cfg.WALDirName())

	_ = os.MkdirAll(cfg.WALDirName(), 0o755)
	f, err := os.Create(filepath.Join(cfg.WALDirName(), "bad.wal"))
	assert.NoError(t, err)
	_, _ = f.WriteString("SET foo bar\n")
	f.Close()

	ctx := context.Background()

	mockInterpreter := newMockinterpreter(t)
	mockInterpreter.On("Execute", ctx, "SET foo bar").Return(nil, errors.New("fail")).Once()

	w := WAL{
		reader:      NewReader(cfg.WALDirName()),
		interpreter: mockInterpreter,
	}
	err = w.Recover(ctx)
	assert.EqualError(t, err, "fail")
}

func TestNewFailsOnNilInterpreter(t *testing.T) {
	ctx := context.Background()
	_, err := New(ctx, testConfig{}, nil)
	assert.Error(t, err)
}
