package models

import (
	"fmt"
	"strconv"
	"time"
)

type DataType string

const (
	Income  DataType = "Income"
	Expense DataType = "Expense"

	RowLength = 7
)

const (
	Year = iota
	Month
	Day
	Description
	Cost
	Category
	Code
	Type
	Date
)

type Record struct {
	Year        int
	Month       int
	Day         int
	Description string
	Cost        int
	Category    string
	Code        int
	Type        DataType
	Date        time.Time
}

func NewRecord(year, month, day int, description string, cost int, category string, code int) *Record {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	recordType := Expense
	if category == string(Income) || code == 0 {
		recordType = Income
	}
	return &Record{
		Year:        year,
		Month:       month,
		Day:         day,
		Description: description,
		Cost:        cost,
		Category:    category,
		Code:        code,
		Type:        recordType,
		Date:        date,
	}
}

func (r *Record) FormattedDate() string {
	return fmt.Sprintf("%02d-%02d-%d", r.Month, r.Day, r.Year)
}

func (r *Record) FormattedCost() string {
	return fmt.Sprintf("$%s", strconv.FormatInt(int64(r.Cost), 10))
}

func (r *Record) IsIncome() bool {
	return r.Type == Income
}

func (r *Record) IsExpense() bool {
	return r.Type == Expense
}

func (r *Record) ToCSVRow() []string {
	return []string{
		fmt.Sprintf("%d", r.Year),
		fmt.Sprintf("%d", r.Month),
		fmt.Sprintf("%d", r.Day),
		r.Description,
		fmt.Sprintf("%d", r.Cost),
		r.Category,
		fmt.Sprintf("%d", r.Code),
	}
}

func RecordFromCSVRow(row []string) (*Record, error) {
	if len(row) < RowLength {
		return nil, fmt.Errorf("invalid CSV row: expected at least 7 columns, got %d", len(row))
	}

	year, err := parseIntWithDefault(row[Year], 0)
	if err != nil {
		return nil, fmt.Errorf("invalid year: %v", err)
	}
	month, err := parseIntWithDefault(row[Month], 0)
	if err != nil {
		return nil, fmt.Errorf("invalid month: %v", err)
	}
	day, err := parseIntWithDefault(row[Day], 0)
	if err != nil {
		return nil, fmt.Errorf("invalid day: %v", err)
	}
	cost, err := parseIntWithDefault(row[Cost], 0)
	if err != nil {
		return nil, fmt.Errorf("invalid cost: %v", err)
	}
	code, err := parseIntWithDefault(row[Code], 0)
	if err != nil {
		return nil, fmt.Errorf("invalid code: %v", err)
	}

	return NewRecord(
		year,
		month,
		day,
		row[Description],
		cost,
		row[Category],
		code,
	), nil
}

func parseIntWithDefault(s string, defaultValue int) (int, error) {
	if s == "" {
		return defaultValue, nil
	}

	var result int
	_, err := fmt.Sscanf(s, "%d", &result)
	if err != nil {
		return defaultValue, err
	}
	return result, nil
}
