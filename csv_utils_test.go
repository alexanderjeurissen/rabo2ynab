package main

import (
	"encoding/csv"
	"os"
	"reflect"
	"testing"
)

func TestReadCSV(t *testing.T) {
	// Create a temporary CSV file for testing
	file, err := os.CreateTemp("", "test.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// Write test data to the temporary CSV file
	writer := csv.NewWriter(file)
	testData := [][]string{
		{"Header1", "Header2"},
		{"Value1", "Value2"},
	}
	if err := writer.WriteAll(testData); err != nil {
		t.Fatalf("Failed to write test data to temp file: %v", err)
	}
	writer.Flush()
	file.Close()

	// Test the readCSV function
	records, err := readCSV(file.Name())
	if err != nil {
		t.Fatalf("readCSV returned an error: %v", err)
	}

	// Verify the records
	if !reflect.DeepEqual(records, testData) {
		t.Errorf("Expected %v, but got %v", testData, records)
	}
}

func TestGetColumnIndices(t *testing.T) {
	headers := []string{"Datum", "Naam tegenpartij", "Omschrijving-1", "Bedrag"}
	expectedIndices := map[string]int{
		"Date":   0,
		"Payee":  1,
		"Memo":   2,
		"Amount": 3,
	}

	indices, err := getColumnIndices(headers)
	if err != nil {
		t.Fatalf("getColumnIndices returned an error: %v", err)
	}

	if !reflect.DeepEqual(indices, expectedIndices) {
		t.Errorf("Expected %v, but got %v", expectedIndices, indices)
	}

	// Test missing columns
	headers = []string{"Datum", "Naam tegenpartij", "Omschrijving-1"}
	_, err = getColumnIndices(headers)
	if err == nil {
		t.Fatalf("Expected an error for missing columns, but got nil")
	}
}

func TestPrepareOutputRecords(t *testing.T) {
	indices := map[string]int{
		"Date":   0,
		"Payee":  1,
		"Memo":   2,
		"Amount": 3,
	}
	data := [][]string{
		{"2021-01-01", "Payee1", "Memo1", "100"},
		{"2021-01-02", "Payee2", "Memo2", "200"},
	}
	expectedOutput := [][]string{
		{"Date", "Payee", "Memo", "Amount"},
		{"2021-01-01", "Payee1", "Memo1", "100"},
		{"2021-01-02", "Payee2", "Memo2", "200"},
	}

	outputRecords := prepareOutputRecords(indices, data)

	if !reflect.DeepEqual(outputRecords, expectedOutput) {
		t.Errorf("Expected %v, but got %v", expectedOutput, outputRecords)
	}
}

func TestWriteCSV(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "output.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(file.Name())

	// Test data to write
	records := [][]string{
		{"Header1", "Header2"},
		{"Value1", "Value2"},
	}

	// Test the writeCSV function
	if err := writeCSV(file.Name(), records); err != nil {
		t.Fatalf("writeCSV returned an error: %v", err)
	}

	// Read the written file and verify the contents
	file, err = os.Open(file.Name())
	if err != nil {
		t.Fatalf("Failed to open temp file: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	writtenRecords, err := reader.ReadAll()
	if err != nil {
		t.Fatalf("Failed to read written file: %v", err)
	}

	if !reflect.DeepEqual(writtenRecords, records) {
		t.Errorf("Expected %v, but got %v", records, writtenRecords)
	}
}
