package wal

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
)

type reader struct {
	dir string
}

func NewReader(dir string) reader {
	return reader{dir: dir}
}

func (r *reader) Read() ([]string, error) {
	var lines []string

	entries, err := os.ReadDir(r.dir)
	if err != nil {
		return nil, err
	}

	// Sort entries by filename (assumes filenames are time-prefixed and sortable)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		path := filepath.Join(r.dir, entry.Name())

		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}

		file.Close()

		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return lines, nil
}
