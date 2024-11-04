package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	reportFile := flag.String("report", "", "Path to the report CSV file")
	rulesFile := flag.String("rules", "", "Path to the rules CSV file")
	jsonOutput := flag.Bool("json", false, "Also output as JSON format")
	flag.Parse()

	if *reportFile == "" || *rulesFile == "" {
		fmt.Println("Both -report and -rules options must be specified")
		return
	}

	reportData, err := readCSV(*reportFile)
	if err != nil {
		fmt.Println("Error reading report file:", err)
		return
	}

	rulesData, err := readCSV(*rulesFile)
	if err != nil {
		fmt.Println("Error reading rules file:", err)
		return
	}

	reportRuleHREFIndex, err := findColumnIndex(reportData[0], "Rule HREF")
	if err != nil {
		fmt.Println(err)
		return
	}

	rulesRuleHrefIndex, err := findColumnIndex(rulesData[0], "rule_href")
	if err != nil {
		fmt.Println(err)
		return
	}

	newFileName := generateNewFileName(*reportFile)
	newFile, err := os.Create(newFileName)
	if err != nil {
		fmt.Println("Error creating new file:", err)
		return
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	defer writer.Flush()

	// Write the headers from both report and rules to the new CSV file
	combinedHeader := append(reportData[0], filterRulesHeaders(rulesData[0], rulesRuleHrefIndex)...)
	if err := writer.Write(combinedHeader); err != nil {
		fmt.Println("Error writing headers to new file:", err)
		return
	}

	// Iterate over report CSV data rows
	for _, reportRow := range reportData[1:] {
		ruleHref := reportRow[reportRuleHREFIndex]
		matched := false

		// Iterate over rules CSV data rows
		for _, rulesRow := range rulesData[1:] {
			if compareHrefs(ruleHref, rulesRow[rulesRuleHrefIndex]) {
				matched = true

				// Append the rules row excluding certain columns
				filteredRow := filterRulesRow(rulesRow, rulesRuleHrefIndex, rulesData[0])
				combinedRow := append(reportRow, filteredRow...)

				// Write the combined row to the new CSV file
				if err := writer.Write(combinedRow); err != nil {
					fmt.Println("Error writing combined row to new file:", err)
					return
				}
				break
			}
		}

		if !matched {
			// If no match found, write the original row with empty placeholders
			emptyPlaceholders := make([]string, len(combinedHeader)-len(reportRow))
			combinedRow := append(reportRow, emptyPlaceholders...)
			if err := writer.Write(combinedRow); err != nil {
				fmt.Println("Error writing report row to new file:", err)
				return
			}
		}
	}
	fmt.Printf("New file created successfully: %s\n", newFileName)

	// Add JSON output if requested
	if *jsonOutput {
		jsonFileName := strings.TrimSuffix(*reportFile, ".csv") + "_plus.json"
		if err := convertCSVToJSON(newFileName, jsonFileName); err != nil {
			fmt.Println("Error creating JSON file:", err)
			return
		}
		fmt.Printf("JSON file created successfully: %s\n", jsonFileName)
	}
}

func readCSV(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	return reader.ReadAll()
}

func findColumnIndex(header []string, columnName string) (int, error) {
	for i, column := range header {
		if column == columnName {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Column '%s' not found", columnName)
}

func generateNewFileName(filePath string) string {
	return filePath[:len(filePath)-4] + "_plus.csv"
}

// compareHrefs compares the rule hrefs ignoring the "active" and "draft" distinction
func compareHrefs(href1, href2 string) bool {
	part1 := strings.Split(href1, "/")
	part2 := strings.Split(href2, "/")

	if len(part1) != len(part2) {
		return false
	}

	// Only compare relevant parts, ignoring "active" and "draft"
	for i := range part1 {
		if i == 4 { // Index 4 corresponds to "active" or "draft"
			continue
		}
		if part1[i] != part2[i] {
			return false
		}
	}
	return true
}

// filterRulesHeaders returns the rules headers excluding certain columns
func filterRulesHeaders(headers []string, rulesRuleHrefIndex int) []string {
	var filteredHeaders []string
	for i, header := range headers {
		if i != rulesRuleHrefIndex && header != "ruleset_name" && header != "ruleset_href" {
			filteredHeaders = append(filteredHeaders, header)
		}
	}
	return filteredHeaders
}

// filterRulesRow returns the rules row excluding certain columns
func filterRulesRow(row []string, rulesRuleHrefIndex int, headers []string) []string {
	var filteredRow []string
	for i, value := range row {
		if i != rulesRuleHrefIndex && headers[i] != "ruleset_name" && headers[i] != "ruleset_href" {
			filteredRow = append(filteredRow, value)
		}
	}
	return filteredRow
}

func convertCSVToJSON(csvPath, jsonPath string) error {
	// Open the CSV file
	file, err := os.Open(csvPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)
	// Set FieldsPerRecord to -1 to allow different field counts
	reader.FieldsPerRecord = -1

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(records) < 2 {
		return fmt.Errorf("CSV file is empty or has no data rows")
	}

	// Get the headers
	headers := records[0]
	var jsonData []map[string]string

	// Process each row of data
	for _, row := range records[1:] {
		item := make(map[string]string)
		// Ensure we don't go out of bounds
		for i := 0; i < len(headers) && i < len(row); i++ {
			item[headers[i]] = row[i]
		}
		// If the row has fewer fields than the headers, fill with empty strings
		for i := len(row); i < len(headers); i++ {
			item[headers[i]] = ""
		}
		jsonData = append(jsonData, item)
	}

	// Create and write to JSON file
	jsonFile, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ")
	return encoder.Encode(jsonData)
}
