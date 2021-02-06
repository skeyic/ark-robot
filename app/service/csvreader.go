package service

import (
	"encoding/csv"
	"errors"
	"github.com/golang/glog"
	"os"
)

var (
	errLoadCSVFile = errors.New("load csv file failed")
)

type CSVReader struct {
	filepath string
}

func NewCSVReader(filename string) *CSVReader {
	return &CSVReader{
		filepath: filename,
	}
}

func (c *CSVReader) Load() ([][]string, error) {
	f, err := os.Open(c.filepath)
	if err != nil {
		glog.Errorf("failed to open csv file: %s, error: %v", c.filepath, err)
		return nil, errLoadCSVFile
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		glog.Errorf("failed to parse as csv file: %s, error: %v", c.filepath, err)
		return nil, errLoadCSVFile
	}

	return records, nil
}
