package internal

import (
	"encoding/csv"
	"os"
	"sync"
)

// CsvWriter writes csv file with mutex lock.
type CsvWriter struct {
	mu sync.Mutex
	w  *csv.Writer
	f  *os.File
}

// NewCsvWriter creates new CsvWriter instance with given filename.
func NewCsvWriter(filename string) (*CsvWriter, error) {
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &CsvWriter{w: csv.NewWriter(file), f: file}, nil
}

// WriteToFile writes given line to csv file.
func (cw *CsvWriter) WriteToFile(line []string) {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	err := cw.w.Write(line)
	if err != nil {
		panic(err)
	}
}

// Flush flushes early wrote data to disk.
func (cw *CsvWriter) Flush() {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	cw.w.Flush()
}

// Close closes underlying files.
func (cw *CsvWriter) Close() {
	cw.mu.Lock()
	defer cw.mu.Unlock()

	err := cw.f.Close()
	if err != nil {
		panic(err)
	}
}
