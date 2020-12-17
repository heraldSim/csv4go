package csv4go

import (
	"bufio"
	"os"
	"strings"
)

const INVALID_INT_VAL = -1

// Cell
type Cell map[string]string

// Record
type Record []Cell

type RecordMapper func(value Record) Record
type RecordFilter func(value Record) bool
type RecordReducer func(a, b interface{}) interface{}

// CSV
type CSV struct {
	index     int
	HeaderNum int
	RowNum    int
	Header    []string
	Records   []Record
}

// LoadCSV
func (csv *CSV) LoadCSV(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	idx := 0
	for scanner.Scan() {
		if idx == 0 {
			csv.Header = ParseLine(strings.Split(scanner.Text(),
				","),
				",")
			csv.HeaderNum = len(csv.Header)
			idx++
		} else {
			tempRow := ParseLine(strings.Split(scanner.Text(),
				","),
				",")
			var tempRecord Record
			for i := 0; i < csv.HeaderNum; i++ {
				cell := make(map[string]string)
				cell[csv.Header[i]] = tempRow[i]

				tempRecord = append(tempRecord, cell)
				idx++
			}
			csv.Records = append(csv.Records, tempRecord)
		}
	}
	csv.RowNum = idx - 1
	csv.index = INVALID_INT_VAL
	
	return nil
}

func (csv *CSV) NextRow() (record Record, ok bool) {
	csv.index++
	if csv.index >= csv.RowNum {
		return Record{}, false
	}
	return csv.Records[csv.index], true
}

func (csv *CSV) Map(mapper RecordMapper) *CSV {
	newRecords := make([]Record, 0, csv.RowNum)
	for _, r := range csv.Records {
		newRecords = append(newRecords, mapper(r))
	}

	return &CSV{
		index: INVALID_INT_VAL,
		HeaderNum: csv.HeaderNum,
		RowNum:    csv.RowNum,
		Header:    csv.Header,
		Records:   newRecords}
}

func (csv *CSV) Filter(filter RecordFilter) *CSV {
	newRecords := make([]Record, 0, csv.HeaderNum)
	for _, v := range csv.Records {
		if filter(v) {
			newRecords = append(newRecords, v)
		}
	}
	return &CSV{
		index: INVALID_INT_VAL,
		HeaderNum: csv.HeaderNum,
		RowNum:    csv.RowNum,
		Header:    csv.Header,
		Records:   newRecords}
}

func (csv *CSV) Reduce(identity interface{}, reducer RecordReducer) interface{} {
	res := identity
	for _, v := range csv.Records {
		res = reducer(res, v)
	}

	return res
}
