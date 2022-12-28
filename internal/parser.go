package internal

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/xmlpath.v2"
)

// Column represents output column with XPath expression that will be processed against source xml file.
type Column struct {
	Name     string
	XPath    *xmlpath.Path
	Optional bool
}

// Line describes type of line (type of file or message).
type Line struct {
	MessageType string
	Columns     []Column
}

// XMLParser intended to search and extract data from xml file to csv.
type XMLParser struct {
	IncludeFilename bool
	Set             []Line
}

// GetHeader returns header for csv file.
func (xp *XMLParser) GetHeader() []string {
	s := make([]string, 0, 4)
	for _, cName := range xp.Set[0].Columns {
		s = append(s, cName.Name)
	}
	s = append(s, "message_type")
	if xp.IncludeFilename {
		s = append(s, "filename")
	}
	return s
}

// Parse parses xml file and returns line of fields for csv.
func (xp *XMLParser) Parse(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatalf("can't open file %s: %s", filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	root, err := xmlpath.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("can't parse file %s: %w", filename, err)
	}

	for _, cf := range xp.Set {
		allFound := true
		line := make([]string, 0, 4)
		for _, col := range cf.Columns {
			val, ok := col.XPath.String(root)
			if !ok && !col.Optional {
				allFound = false
				break
			}
			line = append(line, val)
		}
		if allFound {
			line = append(line, cf.MessageType)
			if xp.IncludeFilename {
				line = append(line, filename)
			}
			return line, nil
		}
	}

	log.Printf("File %s is dropped, useful info was not found.\n", filename)
	return nil, nil
}
