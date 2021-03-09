package service

import (
	"flag"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/golang/glog"
	"strings"
	"testing"
)

func TestCSVRead(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()

	r := NewCSVReader("C:\\Users\\15902\\go\\src\\github.com\\skeyic\\ark-robot\\data\\ARK\\20210304\\a20210304ARKF.csv")
	records, err := r.Load()
	if err != nil {
		glog.Errorf("failed to load, error: %v", err)
		return
	}
	glog.V(4).Infof("RECORDS: %+v", records)

	for idx, k := range strings.TrimSpace(strings.Join(records[0], ",")) {
		glog.V(4).Infof("IDX: %d, value: %v", idx, k)
	}

	for idx, k := range strings.TrimSpace(ARKCSVTitle) {
		glog.V(4).Infof("IDX: %d, value: %v", idx, k)
	}

	glog.V(4).Infof("CSV Validate: %v", ValidateARKCSV(records))

	for idx, record := range records[1:] {
		glog.V(4).Infof("RECORD IDX: %d, stockHolding: %+v, value: %v", idx, NewStockHoldingFromRecord(record), record)
	}
}

func TestCSVReadWrite(t *testing.T) {
	flag.Set("logtostderr", "true")
	flag.Set("v", "10")
	flag.Parse()

	r := NewCSVReader("C:\\Users\\15902\\go\\src\\github.com\\skeyic\\ark-robot\\data\\ARK\\20210304\\arkqw.csv")
	records, err := r.Load()
	if err != nil {
		glog.Errorf("failed to load, error: %v", err)
		return
	}
	glog.V(4).Infof("RECORDS: %+v", records)

	var (
		shardsMap = make(map[string][]string)
	)

	for _, record := range records {
		shardsMap[record[0]] = []string{record[1], record[3]}
		glog.V(4).Infof("record: %+v", record[1])
	}

	f, err := excelize.OpenFile("C:\\Users\\15902\\go\\src\\github.com\\skeyic\\ark-robot\\data\\ARK\\20210304\\20210303ARKF.xlsx")
	if err != nil {
		glog.Errorf("failed to open excel, err: %v", err)
		return
	}

	sheet := "20210304ARKF"
	for i := 1; ; i++ {
		var (
			stockF  = fmt.Sprintf("D%d", i)
			tickerF = fmt.Sprintf("E%d", i)
			shardF  = fmt.Sprintf("F%d", i)
			weightF = fmt.Sprintf("H%d", i)
		)

		s, err := f.GetCellValue(sheet, stockF)
		if err != nil {
			glog.Errorf("failed to get value, err: %v", err)
			return
		}

		t, err := f.GetCellValue(sheet, tickerF)
		if err != nil {
			glog.Errorf("failed to get value, err: %v", err)
			return
		}

		if t == "" {
			break
		}

		r := shardsMap[s]
		if len(r) == 2 {
			f.SetCellValue(sheet, shardF, r[0])
			glog.V(4).Infof("SET %s shard to %s", s, r[0])
			f.SetCellValue(sheet, weightF, r[1])
			glog.V(4).Infof("SET %s weight to %s", s, r[1])
		} else {
			glog.V(4).Infof("No such record %s", s)
		}

	}

	err = f.Save()
	if err != nil {
		glog.Errorf("failed to save excel, err: %v", err)
		return
	}

}
