package cfg

import (
	"fmt"
	"os"

	"xml2csv-parser/internal"

	"gopkg.in/xmlpath.v2"
	"gopkg.in/yaml.v2"
)

type Column struct {
	Name     string `yaml:"name"`
	XPath    string `yaml:"xpath"`
	Optional bool   `yaml:"optional,omitempty"`
}

type Line struct {
	MessageType string   `yaml:"messageType"`
	Columns     []Column `yaml:"columns"`
}

type XMLParser struct {
	IncludeFilename bool   `yaml:"includeFilename,omitempty"`
	Set             []Line `yaml:"set"`
}

func (xp *XMLParser) Load(filename string) (err error) {
	var yamlBytes []byte
	if filename == "" {
		yamlBytes = []byte(defaultYaml)
	} else {
		yamlBytes, err = os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("can't read file %s: %w", filename, err)
		}
	}
	err = yaml.Unmarshal(yamlBytes, xp)
	if err != nil {
		return fmt.Errorf("can't parse file %s: %w", filename, err)
	}
	return nil
}

func (xp *XMLParser) CreateCompiled() *internal.XMLParser {
	cmpLines := make([]internal.Line, 0, len(xp.Set))
	for _, line := range xp.Set {
		cmpCols := make([]internal.Column, 0, len(line.Columns))
		for _, col := range line.Columns {
			cmpCols = append(cmpCols, internal.Column{
				Name:  col.Name,
				XPath: xmlpath.MustCompile(col.XPath),
			})
		}
		cmpLines = append(cmpLines, internal.Line{
			MessageType: line.MessageType,
			Columns:     cmpCols,
		})
	}
	return &internal.XMLParser{
		IncludeFilename: xp.IncludeFilename,
		Set:             cmpLines,
	}
}
