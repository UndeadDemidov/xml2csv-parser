package internal

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

// Producer producing work for parser from a list of xml files in given path.
type Producer struct {
	root, ext string
	filenames *chan string
	quit      chan interface{}
}

// NewProducer creates of Producer instance.
func NewProducer(path, fileExt string, l *chan string, q chan interface{}) *Producer {
	return &Producer{root: path, ext: fileExt, filenames: l, quit: q}
}

// Produce sends filenames to consumer to parse files concurently.
func (p *Producer) Produce() {
	fileList := p.find(p.root, p.ext)
	cnt := 0

loop:
	for _, filename := range fileList {
		select {
		case *p.filenames <- filename:
			cnt++
		case <-p.quit:
			break loop
		}
	}
	close(*p.filenames)
	fmt.Printf("Processed %d files out of %d\n", cnt, len(fileList))
}

func (p *Producer) find(root, ext string) []string {
	var a []string
	err := filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ext {
			a = append(a, s)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return a
}
