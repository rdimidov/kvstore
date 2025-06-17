package wal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReader_Read_MultipleFiles(t *testing.T) {
	dir := t.TempDir()

	// create files with predictable timestamps in filenames
	files := []struct {
		name    string
		content string
	}{
		{"20240101T000000.wal", "first1\nfirst2\n"},
		{"20240101T000001.wal", "second1\nsecond2\n"},
		{"20240101T000002.wal", "third1\nthird2\n"},
	}

	for _, f := range files {
		err := os.WriteFile(filepath.Join(dir, f.name), []byte(f.content), 0o644)
		if err != nil {
			t.Fatalf("failed to write %s: %v", f.name, err)
		}
	}

	r := NewReader(dir)
	lines, err := r.Read()
	if err != nil {
		t.Fatalf("Read failed: %v", err)
	}

	expected := []string{
		"first1", "first2",
		"second1", "second2",
		"third1", "third2",
	}
	if len(lines) != len(expected) {
		t.Fatalf("expected %d lines, got %d", len(expected), len(lines))
	}
	for i := range expected {
		if lines[i] != expected[i] {
			t.Errorf("at index %d: expected %q, got %q", i, expected[i], lines[i])
		}
	}
}
