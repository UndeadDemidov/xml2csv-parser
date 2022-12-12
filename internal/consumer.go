package internal

import (
	"context"
	"log"
	"sync"
	"time"
)

type Consumer struct {
	xmlParser *XMLParser
	csvWriter *CsvWriter
	filenames *chan string
	jobs      chan string
}

func NewConsumer(xp *XMLParser, cw *CsvWriter, l *chan string, j chan string) *Consumer {
	c := &Consumer{xmlParser: xp, csvWriter: cw, filenames: l, jobs: j}
	c.csvWriter.WriteToFile(xp.GetHeader())
	return c
}

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
		}
	}
}

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
