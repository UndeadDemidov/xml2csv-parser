package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copySrc(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file: %w", src, err)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err = source.Close(); err != nil {
			fmt.Printf("can't close file: %s", err)
		}
	}()

	destination, err := os.Create(filepath.Join(dst, source.Name()))
	if err != nil {
		return 0, err
	}
	defer func() {
		if err = destination.Close(); err != nil {
			fmt.Printf("can't close file: %s", err)
		}
	}()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
