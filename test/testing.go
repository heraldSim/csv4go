package main

import (
	"fmt"
	"github.com/heraldSim/csv4go"
	"strconv"
	"strings"
)

func convertCap(record csv4go.Record) csv4go.Record {
	newRecod := csv4go.Record{}
	for _, cell := range record {
		if val, ok := cell[`"Name"`]; ok {
			capVal := strings.ToUpper(val)
			newCell := csv4go.Cell{"Name": capVal}

			newRecod = append(newRecod, newCell)
		} else {
			newRecod = append(newRecod, cell)
		}
	}
	return newRecod
}

func filteringAge(record csv4go.Record) bool {
	for _, cell := range record {
		if val, ok := cell[`"Age"`]; ok {
			age, _ := strconv.ParseFloat(strings.TrimSpace(val), 64)
			if age >= 21 {
				return true
			} else {
				return false
			}
		}
	}
	// No 'Age' Column
	return false
}

func totalAge(a, b interface{}) interface{} {
	var age float64
	for _, cell := range b.(csv4go.Record) {
		if val, ok := cell[`"Age"`]; ok {
			age, _ = strconv.ParseFloat(strings.TrimSpace(val), 64)
		}
	}
	return a.(float64) + age
}

func main() {
	fmt.Println("##### Example of Using Collection functions #####")

	myCsv := csv4go.CSV{}
	err := myCsv.LoadCSV("./mlb_players.csv", ", ")
	if err != nil {
		fmt.Print("Load error")
	}

	fmt.Println("Convert all Name as uppercase")
	// Apply Map
	mapCsv := myCsv.Map(convertCap)
	// Print all rows
	row, ok := mapCsv.NextRow()
	for ok {
		fmt.Println(row)
		row, ok = mapCsv.NextRow()
	}

	// Apply Filter and Reduce
	filtered := mapCsv.Filter(filteringAge)
	count := filtered.RowNum
	avg := filtered.Reduce(float64(0), totalAge)

	fmt.Printf("Get average age over 21: %f \n", avg.(float64) / float64(count))
}

