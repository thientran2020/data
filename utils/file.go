package utils

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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
func CsvWrite(filepath string, record m.Record) bool {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Cannot open file to write ", err)
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

func CsvRead(filepath string) Data {
	data := Data{}

	file, err := os.Open(filepath)
	if err != nil {
		fmt.Printf("Error opening file %s", err)
	}
	reader := csv.NewReader(file)

	// discard the header
	_, err = reader.Read()
	if err != nil {
		fmt.Printf("Error reading file %s", err)
	}

	count := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading file %s", err)
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
	}
	return data
}

func FilterData(data Data, month int, typeFlag, keyword string) Data {
	filteredData := Data{}
	for _, row := range data {
		dateArray := strings.Split(row[1].(string), "-")
		row_month, _ := strconv.Atoi(dateArray[0])
		skipRowByMonth := month != -1 && row_month != month

		skipRowByTypeFlag :=
			(typeFlag == "income" && row[3] != "Income") ||
				(typeFlag == "expense" && row[3] == "Income")

		skipRowByKeyWord := !ContainString(row[2].(string), keyword) && !ContainString(row[3].(string), keyword)

		skip := skipRowByMonth || skipRowByTypeFlag || skipRowByKeyWord
		if !skip {
			formatted_date := fmt.Sprintf("%s-%s-%s",
				Colorize(dateArray[0], Yellow),
				dateArray[1],
				Colorize(dateArray[2], UGreen),
			)
			row[1] = formatted_date
			filteredData = append(filteredData, row)
		}
	}
	return filteredData
}

// json file processing
func ReadJson(filepath string) m.MySubscriptionList {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		CreateFile(filepath)
		fmt.Printf("Subcription list was created at %v\n", filepath)
		return ReadJson(filepath)
	}
	result := m.MySubscriptionList{}

	_ = json.Unmarshal([]byte(file), &result)
	return result
}

func WriteJson(filepath string, subscriptionList m.MySubscriptionList) {
	file, _ := json.MarshalIndent(subscriptionList, "", " ")
	err := ioutil.WriteFile(filepath, file, 0644)
	if err != nil {
		fmt.Printf("Error writing subscription data %v\n", err)
	}
}

// Return filepaths for all-in-one file "finance.csv" and current year file "finance_<year>.csv"
// Check if files exist - if not create new ones
func GetSharedFile() string {
	sharedFilePath := strings.Replace(m.BASE_FILEPATH, "<YEAR>", "", -1)
	if !FileExists(sharedFilePath) {
		CreateFile(sharedFilePath)
	}
	return sharedFilePath
}

func GetSpecificYearFile(year int) string {
	currentYearFilePath := strings.Replace(m.BASE_FILEPATH, "<YEAR>", fmt.Sprintf("_%d", year), -1)
	if !FileExists(currentYearFilePath) {
		CreateFile(currentYearFilePath)
	}
	return currentYearFilePath
}
