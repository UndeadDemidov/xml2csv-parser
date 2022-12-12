package internal

import (
	"encoding/csv"
	"os"
	"sync"
)

type CsvWriter struct {
	mu sync.Mutex
	w  *csv.Writer
	f  *os.File
}

func NewCsvWriter(filename string) (*CsvWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &CsvWriter{w: csv.NewWriter(file), f: file}, nil
}

func (cw *CsvWriter) WriteToFile(line []string) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	cw.w.Write(line)
}

func (cw *CsvWriter) Flush() {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	cw.w.Flush()
}

func (cw *CsvWriter) Close() {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	cw.f.Close()
}
