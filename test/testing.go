package main

import (
	"fmt"
	"github.com/heraldSim/csv4go"
)

func main() {
	myCsv := csv4go.CSV{}
	err := myCsv.LoadCSV("./mlb_players.csv", ",")
	if err != nil {
		fmt.Print("Load error")
	}

	fmt.Println(myCsv.RowNum)

	row, ok := myCsv.NextRow()
	for ok {
		fmt.Println(row)
		row, ok = myCsv.NextRow()
	}
}
