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
	"time"

	"github.com/thientran2020/financial-cli/models"
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
func CsvWrite(filepath string, record models.Record) bool {
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

func CsvRead(requestedYear int, requestedMonth int, typeFlag string, keyword string) [][]interface{} {
	filepath := strings.Replace(models.BASE_FILEPATH, "<YEAR>", "", -1)
	if requestedYear >= models.START_YEAR && requestedYear <= time.Now().Year() {
		filepath = strings.Replace(models.BASE_FILEPATH, "<YEAR>", fmt.Sprintf("_%d", requestedYear), -1)
	}

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

	data := [][]interface{}{}
	count := 0
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error reading file %s", err)
		}

		month, _ := strconv.Atoi(row[month])
		skipRowByRequestedMonth := requestedMonth != -1 && month != requestedMonth
		skipRowByTypeFlag :=
			(typeFlag == "income" && strings.Trim(row[category], " ") != "Income") ||
				(typeFlag == "expense" && strings.Trim(row[category], " ") == "Income")
		skipRowByKeyWord := !ContainString(row[content], keyword) && !ContainString(row[category], keyword)

		if skipRowByRequestedMonth || skipRowByTypeFlag || skipRowByKeyWord {
			continue
		}

		count++
		cost, _ := strconv.Atoi(row[cost])

		rowData := []interface{}{
			count,
			fmt.Sprintf("%s-%s-%s",
				Colorize(GetStringDateFromNumber(month), Yellow),
				GetStringDateFromString(row[day]),
				Colorize(row[year], UGreen),
			),
			strings.Trim(row[content], " "),
			strings.Trim(row[category], " "),
			cost,
		}
		data = append(data, rowData)
	}

	return data
}

// json file processing
func ReadJson(filepath string) models.MySubscriptionList {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		CreateFile(filepath)
		fmt.Printf("Subcription list was created at %v\n", filepath)
		return ReadJson(filepath)
	}
	result := models.MySubscriptionList{}

	_ = json.Unmarshal([]byte(file), &result)
	return result
}

func WriteJson(filepath string, subscriptionList models.MySubscriptionList) {
	file, _ := json.MarshalIndent(subscriptionList, "", " ")
	err := ioutil.WriteFile(filepath, file, 0644)
	if err != nil {
		fmt.Printf("Error writing subscription data %v\n", err)
	}
}

// Return filepaths for all-in-one file "finance.csv" and current year file "finance_<year>.csv"
// Check if files exist - if not create new ones
func GetSharedFile() string {
	sharedFilePath := strings.Replace(models.BASE_FILEPATH, "<YEAR>", "", -1)
	if !FileExists(sharedFilePath) {
		CreateFile(sharedFilePath)
	}
	return sharedFilePath
}

func GetCurrentYearFile() string {
	currentYearFilePath := strings.Replace(models.BASE_FILEPATH, "<YEAR>", fmt.Sprintf("_%d", time.Now().Year()), -1)
	if !FileExists(currentYearFilePath) {
		CreateFile(currentYearFilePath)
	}
	return currentYearFilePath
}
