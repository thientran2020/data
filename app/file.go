package app

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/thientran2020/financial-cli/models"
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

// json file processing
func readJson(filepath string) models.MySubscriptionList {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Printf("Error reading subscription file %v\n", err)
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
