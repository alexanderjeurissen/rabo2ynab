package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alexanderjeurissen/rabodebit2ynab/csv_utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: csv_parser <input_file.csv>")
		return
	}

	inputFile := os.Args[1]
	outputFile := fmt.Sprintf("ynab_%d.csv", time.Now().Unix())

	records, err := csv_utils.readCSV(inputFile)
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return
	}

	headers, data := records[0], records[1:]
	columnIndices, err := csv_utils.getColumnIndices(headers)
	if err != nil {
		fmt.Printf("Error getting column indices: %v\n", err)
		return
	}

	outputRecords := csv_utils.prepareOutputRecords(columnIndices, data)

	if err := csv_utils.writeCSV(outputFile, outputRecords); err != nil {
		fmt.Printf("Error writing to CSV file: %v\n", err)
		return
	}

	fmt.Printf("Processed CSV written to %s\n", filepath.Join(".", outputFile))
}
