package wal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRotatingWalWriter_WriteAndRotate(t *testing.T) {
	dir := t.TempDir()
	maxBytes := 5 // small threshold to trigger rotation

	writer, err := newRotatingWalWriter(dir, maxBytes)
	if err != nil {
		t.Fatalf("failed to create writer: %v", err)
	}

	batch := []entry{
		newEntry(strings.Repeat("a", 30)),
		newEntry(strings.Repeat("b", 30)),
	}

	writer.Write(batch)
	time.Sleep(time.Second)
	writer.Write(batch)

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("failed to read WAL directory: %v", err)
	}

	if len(files) < 2 {
		t.Errorf("expected at least 2 WAL files due to rotation, got %d", len(files))
	}

	var totalLines int
	for _, file := range files {
		data, err := os.ReadFile(filepath.Join(dir, file.Name()))
		if err != nil {
			t.Errorf("failed to read file %s: %v", file.Name(), err)
			continue
		}
		lines := strings.Split(strings.TrimSpace(string(data)), "\n")
		totalLines += len(lines)
	}

	allEntries := len(batch) + len(batch)
	if totalLines != allEntries {
		t.Errorf("expected %d log lines, found %d", allEntries, totalLines)
	}
}
