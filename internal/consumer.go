package internal

import (
	"context"
	"log"
	"sync"
	"time"
)

// Consumer consumes a batch xml files to process then concurently.
type Consumer struct {
	xmlParser *XMLParser
	csvWriter *CsvWriter
	filenames *chan string
	jobs      chan string
	copyPath  string
}

// NewConsumer creates instance of Consumer with given parameters.
func NewConsumer(xp *XMLParser, cw *CsvWriter, l *chan string, j chan string, copyPath string) *Consumer {
	c := &Consumer{xmlParser: xp, csvWriter: cw, filenames: l, jobs: j, copyPath: copyPath}
	c.csvWriter.WriteToFile(xp.GetHeader())
	return c
}

// Work does work.
func (c *Consumer) Work(wg *sync.WaitGroup) {
	defer wg.Done()
	for job := range c.jobs {
		if job == "" {
			return
		}
		line, err := c.xmlParser.Parse(job)
		if err != nil {
			log.Println(err)
			continue
		}
		if line != nil {
			c.csvWriter.WriteToFile(line)
			if c.copyPath != "" {
				// ToDo make copy in goroutines.
				_, err = copySrc(job, c.copyPath)
				log.Println(err)
				return
			}
		}
	}
}

// Consume consumes workload and pushes for workers.
func (c *Consumer) Consume(ctx context.Context) {
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			c.csvWriter.Flush()
		case job := <-*c.filenames:
			c.jobs <- job
		case <-ctx.Done():
			close(c.jobs)
			return
		}
	}
}
