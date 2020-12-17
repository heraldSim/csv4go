package csv4go

import (
	"bufio"
	"os"
	"strings"
)

// Cell
type Cell map[string]string

// Record
type Record []Cell

type RecordMapper func (value Cell) (Cell)
type RecordFilter func (value Cell) (bool)
type RecordReducer func (a, b interface{}) (interface{})

// CSV
type CSV struct {
	HeaderNum int
	RowNum    int
	Header    []string
	Records   Record
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
			for i :=0; i < csv.HeaderNum; i++ {
				cell := make(map[string]string)
				cell[csv.Header[i]] = tempRow[i]

				csv.Records = append(csv.Records, cell)
				idx++
			}
		}
	}
	csv.RowNum = idx - 1
	return nil
}


func (csv *CSV) Map(mapper RecordMapper) (*CSV) {
	newRecord := make(Record, 0, len(csv.Header))
	for _, r := range csv.Records {
		newRecord = append(newRecord, mapper(r))
	}

	return csv
}

func (csv *CSV) Filter(filter RecordFilter) (*CSV) {
	newRecord := make(Record, 0, len(csv.Header))
	for _, v := range csv.Records {
		if filter(v) {
			newRecord = append(newRecord, v)
		}
	}
	return csv
}

func (csv *CSV) Reduce(identity interface{}, reducer RecordReducer) (interface{}) {
	res := identity
	for _, v := range csv.Records {
		res = reducer(res, v)
	}

	return res
}