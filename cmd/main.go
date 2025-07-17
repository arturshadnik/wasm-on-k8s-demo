package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// Get input and output paths from environment variables
	inputPath := os.Getenv("INPUT_PATH")
	if inputPath == "" {
		fmt.Fprintf(os.Stderr, "INPUT_PATH environment variable is required\n")
		os.Exit(1)
	}

	outputPath := os.Getenv("OUTPUT_PATH")
	if outputPath == "" {
		fmt.Fprintf(os.Stderr, "OUTPUT_PATH environment variable is required\n")
		os.Exit(1)
	}

	// Read and process the CSV file
	sums, err := processCSV(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error processing CSV: %v\n", err)
		os.Exit(1)
	}

	// Sort sums in descending order
	sort.Sort(sort.Reverse(sort.Float64Slice(sums)))

	// Write results to output file
	err = writeResults(outputPath, sums)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
		os.Exit(1)
	}

	// Write results to stdout
	for _, sum := range sums {
		fmt.Printf("%.2f\n", sum)
	}
}

func processCSV(inputPath string) ([]float64, error) {
	// Read the input file
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read input file: %w", err)
	}

	// Parse CSV
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to parse CSV: %w", err)
	}

	var sums []float64

	// Process each row (skip header if present)
	for _, record := range records {
		if len(record) < 3 {
			continue // Skip rows with fewer than 3 columns
		}

		// Parse the three numeric columns
		var rowSum float64
		for j := 0; j < 3; j++ {
			value, err := strconv.ParseFloat(strings.TrimSpace(record[j]), 64)
			if err != nil {
				// Skip rows with invalid numbers
				continue
			}
			rowSum += value
		}

		// Only add to results if we successfully parsed all 3 numbers
		if len(record) >= 3 {
			sums = append(sums, rowSum)
		}
	}

	return sums, nil
}

func writeResults(outputPath string, sums []float64) error {
	// Create output file
	dir := filepath.Dir(outputPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
			return fmt.Errorf("failed to create output directory: %w", mkErr)
		}
	}
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Write CSV header
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write each sum as a row
	for _, sum := range sums {
		err := writer.Write([]string{fmt.Sprintf("%.2f", sum)})
		if err != nil {
			return fmt.Errorf("failed to write to CSV: %w", err)
		}
	}

	return nil
}
