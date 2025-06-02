package wal

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

const baseFileName = "wal"

// rotatingWalWriter implements walWriter and can "fold" logs into segments:
// as soon as one file grows to maxBytes, it is closed and a new one is started.
type rotatingWalWriter struct {
	mu       sync.Mutex
	dir      string // directory where to put segments
	baseName string // segment file  prefix
	maxBytes int    // max segment size
	curFile  *os.File
	curSize  int // curr segment size
}

func newRotatingWalWriter(dir string, maxBytes int) (*rotatingWalWriter, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	w := &rotatingWalWriter{
		dir:      dir,
		baseName: baseFileName,
		maxBytes: maxBytes,
	}
	w.openLastCreatedFile()

	if w.curFile == nil {
		if err := w.rotate(); err != nil {
			return nil, err
		}
	}

	return w, nil
}

// rotate closes the current file (if it is open) and creates a new segment.
func (w *rotatingWalWriter) rotate() error {
	if w.curFile != nil {
		if err := w.curFile.Close(); err != nil {
			return err
		}
		w.curFile = nil
	}

	timestamp := time.Now().Format("20060102T150405")
	filename := fmt.Sprintf("%s.%s", timestamp, w.baseName)
	fullpath := filepath.Join(w.dir, filename)

	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	w.curFile = f
	w.curSize = 0
	return nil
}

// Write writes a batch of entries to the current WAL segment.
func (w *rotatingWalWriter) Write(batch []entry) {
	w.mu.Lock()
	defer w.mu.Unlock()

	var buf bytes.Buffer
	for _, e := range batch {
		buf.WriteString(e.data)
		buf.WriteByte('\n')
	}

	if w.curSize+buf.Len() > w.maxBytes {
		if err := w.rotate(); err != nil {
			for _, e := range batch {
				e.SetResponse(err)
			}
			return
		}
	}

	n, err := w.curFile.Write(buf.Bytes())
	if err == nil && n < buf.Len() {
		err = errors.New("short write to WAL")
	}
	if err == nil {
		w.curSize += n
		err = w.curFile.Sync()
	}

	for _, e := range batch {
		e.SetResponse(err)
	}
}

func (w *rotatingWalWriter) openLastCreatedFile() {
	files, err := os.ReadDir(w.dir)
	if err != nil || len(files) == 0 {
		return
	}

	type namedTime struct {
		name string
		t    int64
	}

	var candidates []namedTime

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		info, err := file.Info()
		if err != nil {
			continue
		}
		candidates = append(candidates, namedTime{file.Name(), info.ModTime().UnixNano()})
	}
	if len(candidates) == 0 {
		return
	}

	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].t > candidates[j].t
	})

	f, err := os.OpenFile(filepath.Join(w.dir, candidates[0].name), os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}

	info, err := f.Stat()
	if err != nil {
		f.Close()
		return
	}

	w.curFile = f
	w.curSize = int(info.Size())
}
