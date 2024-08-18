package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	m "github.com/thientran2020/financial-cli/models"
)

const (
	year int = iota
	month
	day
	content
	cost
	category
	code
)

type Data [][]interface{}
type String2D [][]string

// Get user's home directory
func GetUserHomeDirectory() string {
	userHomeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return userHomeDir
}

// Create folder if not exist
func CreateFolderIfNotExist() bool {
	path := GetUserHomeDirectory() + "/finance"
	var err error
	if !FileExists(path) {
		err = os.Mkdir(path, os.ModePerm)
	}
	return path != "" || err != nil
}

// file processing with os
func FileExists(filepath string) bool {
	file, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !file.IsDir()
}

func CreateFile(filepath string) bool {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Cannot create new file at %s\n", filepath)
		return false
	}
	defer file.Close()
	return true
}

// csv file processing
func CsvWriteHeader(filepath string) bool {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return false
	}
	writer := csv.NewWriter(file)
	writer.Write(m.ROW_HEADER)
	writer.Flush()
	return true
}

func CsvWrite(filepath string, record m.Record) bool {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return false
	}

	recordString := []string{
		strconv.Itoa(record.Year),
		strconv.Itoa(record.Month),
		strconv.Itoa(record.Day),
		record.Description,
		strconv.Itoa(record.Cost),
		record.Category,
		strconv.Itoa(record.Code),
	}

	writer := csv.NewWriter(file)
	writer.Write(recordString)
	writer.Flush()
	return true
}

func CsvRead(filepath string) (Data, String2D) {
	original_data := String2D{}
	data := Data{}

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error opening file %s\n", err)
	}
	reader := csv.NewReader(file)

	// discard the header
	reader.Read()

	count := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		count++
		month, _ := strconv.Atoi(row[month])
		cost, _ := strconv.Atoi(row[cost])
		rowData := []interface{}{
			count,
			fmt.Sprintf("%s-%s-%s",
				GetStringDateFromNumber(month),
				GetStringDateFromString(row[day]),
				row[year],
			),
			strings.Trim(row[content], " "),
			strings.Trim(row[category], " "),
			cost,
		}

		data = append(data, rowData)
		original_data = append(original_data, row)
	}
	return data, original_data
}

// JSON file processing
func ReadSubscriptionJson(filepath string) m.MySubscriptionList {
	file, err := os.ReadFile(filepath)
	if err != nil {
		CreateFile(filepath)
		return ReadSubscriptionJson(filepath)
	}
	result := m.MySubscriptionList{}

	_ = json.Unmarshal([]byte(file), &result)
	return result
}

func WriteSubscriptionJson(filepath string, subscriptionList m.MySubscriptionList) {
	file, _ := json.MarshalIndent(subscriptionList, "", " ")
	err := os.WriteFile(filepath, file, 0644)
	if err != nil {
		fmt.Printf("Error writing subscription data %v\n", err)
	}
}

// Return filepaths for all-in-one file "finance.csv" and current year file "finance_<year>.csv"
// Check if files exist - if not create new ones
func GetSharedFile() string {
	sharedFilePath := GetUserHomeDirectory() + strings.Replace(m.BASE_FILEPATH, "<YEAR>", "", -1)
	if !FileExists(sharedFilePath) {
		CreateFile(sharedFilePath)
		CsvWriteHeader(sharedFilePath)
	}
	return sharedFilePath
}

func GetSpecificYearFile(year int) string {
	currentYearFilePath := GetUserHomeDirectory() + strings.Replace(m.BASE_FILEPATH, "<YEAR>", fmt.Sprintf("_%d", year), -1)
	if !FileExists(currentYearFilePath) {
		CreateFile(currentYearFilePath)
		CsvWriteHeader(currentYearFilePath)
	}
	return currentYearFilePath
}

// Read + Sort by date + Overwrite CSV file with sorted data
func CsvUpdate(filepath string) {
	_, data := CsvRead(filepath)
	sort.Sort(data)

	// Overwrite existing file
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Cannot open file at %s: %s", filepath, err)
	}
	defer file.Close()

	// Add row headers for csv file
	data = append([][]string{m.ROW_HEADER}, data...)

	writer := csv.NewWriter(file)
	writer.WriteAll(data)
	if err := writer.Error(); err != nil {
		fmt.Printf("Error overwriting at %s: %s\n", filepath, err)
	}
}

// Sort 2D array
// Call sort.Sort(Data(data)) to sort
func (data String2D) Less(i, j int) bool {
	y1, _ := strconv.Atoi(data[i][0])
	m1, _ := strconv.Atoi(data[i][1])
	d1, _ := strconv.Atoi(data[i][2])

	y2, _ := strconv.Atoi(data[j][0])
	m2, _ := strconv.Atoi(data[j][1])
	d2, _ := strconv.Atoi(data[j][2])

	I_before_J := y1 < y2 || (y1 == y2 && (m1 < m2 || (m1 == m2 && d1 <= d2)))
	return I_before_J
}

func (data String2D) Len() int {
	return len(data)
}

func (data String2D) Swap(i, j int) {
	data[i], data[j] = data[j], data[i]
}
