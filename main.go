package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: csv_parser <input_file.csv>")
		return
	}

	inputFile := os.Args[1]
	outputFile := fmt.Sprintf("ynab_%d.csv", time.Now().Unix())

	records, err := readCSV(inputFile)
	if err != nil {
		fmt.Printf("Error reading CSV file: %v\n", err)
		return
	}

	headers, data := records[0], records[1:]
	columnIndices := getColumnIndices(headers)

	outputRecords := prepareOutputRecords(columnIndices, data)

	if err := writeCSV(outputFile, outputRecords); err != nil {
		fmt.Printf("Error writing to CSV file: %v\n", err)
		return
	}

	fmt.Printf("Processed CSV written to %s\n", filepath.Join(".", outputFile))
}

// readCSV reads a CSV file and returns the records.
// It takes the filename as an argument and returns a 2D slice of strings
// and an error if any occurs during the file opening or reading process.
func readCSV(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("reading CSV file: %w", err)
	}

	return records, nil
}

// getColumnIndices maps the required column names to their indices.
// It takes a slice of strings representing the headers and returns a map
// with the new column names as keys and their respective indices as values.
func getColumnIndices(headers []string) map[string]int {
	columnIndices := make(map[string]int)
	for i, header := range headers {
		switch header {
		case "Datum":
			columnIndices["Date"] = i
		case "Naam tegenpartij":
			columnIndices["Payee"] = i
		case "Omschrijving-1":
			columnIndices["Memo"] = i
		case "Bedrag":
			columnIndices["Amount"] = i
		}
	}
	return columnIndices
}

// prepareOutputRecords prepares the output records with renamed columns.
// It takes a map of column indices and a 2D slice of strings representing the data,
// and returns a new 2D slice of strings with the columns renamed and reordered.
func prepareOutputRecords(indices map[string]int, data [][]string) [][]string {
	outputRecords := [][]string{
		{"Date", "Payee", "Memo", "Amount"},
	}

	for _, record := range data {
		outputRecords = append(outputRecords, []string{
			record[indices["Date"]],
			record[indices["Payee"]],
			record[indices["Memo"]],
			record[indices["Amount"]],
		})
	}

	return outputRecords
}

// writeCSV writes the records to a CSV file.
// It takes the filename and a 2D slice of strings representing the records,
// and returns an error if any occurs during the file creation or writing process.
func writeCSV(filename string, records [][]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("creating output file: %w", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("writing to CSV file: %w", err)
	}

	return nil
}
