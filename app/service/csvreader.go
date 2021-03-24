package service

import (
	"encoding/csv"
	"errors"
	"github.com/golang/glog"
	"os"
)

var (
	errLoadCSVFile  = errors.New("load csv file failed")
	errWriteCSVFile = errors.New("write csv file failed")
)

type CSVOperator struct {
	filepath string
}

func NewCSVOperator(filename string) *CSVOperator {
	return &CSVOperator{
		filepath: filename,
	}
}

func (c *CSVOperator) Load() ([][]string, error) {
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

func (c *CSVOperator) Write(content [][]string) error {
	f, err := os.Create(c.filepath)
	if err != nil {
		glog.Errorf("failed to open csv file: %s, error: %v", c.filepath, err)
		return errWriteCSVFile
	}
	defer f.Close()

	w := csv.NewWriter(f)
	err = w.WriteAll(content)
	if err != nil {
		glog.Errorf("failed to open csv file: %s, error: %v", c.filepath, err)
		return errWriteCSVFile
	}
	w.Flush()

	return nil
}
