package app

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
	"github.com/thientran2020/financial-cli/utils"
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
func fileExists(filepath string) bool {
	file, err := os.Stat(filepath)
	if os.IsNotExist(err) {
		return false
	}
	return !file.IsDir()
}

func createFile(filepath string) bool {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Printf("Cannot create new file at %s\n", filepath)
		return false
	}
	defer file.Close()
	return true
}

// csv file processing
func csvWrite(filepath string, record models.Record) bool {
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

func csvRead(requestedYear int, requestedMonth int, typeFlag string, keyword string) [][]interface{} {
	filepath := strings.Replace(models.BASE_FILEPATH, "?????", "", -1)
	if requestedYear >= 2017 && requestedYear <= time.Now().Year() {
		filepath = strings.Replace(models.BASE_FILEPATH, "?????", fmt.Sprintf("_%d", requestedYear), -1)
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
		skipRowByKeyWord := !utils.ContainString(row[content], keyword) && !utils.ContainString(row[category], keyword)

		if skipRowByRequestedMonth || skipRowByTypeFlag || skipRowByKeyWord {
			continue
		}

		count++
		cost, _ := strconv.Atoi(row[cost])

		rowData := []interface{}{
			count,
			fmt.Sprintf("%s-%s-%s",
				utils.Colorize(utils.GetStringDateFromNumber(month), utils.Yellow),
				utils.GetStringDateFromString(row[day]),
				utils.Colorize(row[year], utils.UGreen),
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
func readJson(filepath string) models.MySubscriptionList {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		createFile(filepath)
		fmt.Printf("Subcription list was created at %v\n", filepath)
		return readJson(filepath)
	}
	result := models.MySubscriptionList{}

	_ = json.Unmarshal([]byte(file), &result)
	return result
}

func writeJson(filepath string, subscriptionList models.MySubscriptionList) {
	file, _ := json.MarshalIndent(subscriptionList, "", " ")
	err := ioutil.WriteFile(filepath, file, 0644)
	if err != nil {
		fmt.Printf("Error writing subscription data %v\n", err)
	}
}
