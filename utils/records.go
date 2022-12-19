package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	m "github.com/thientran2020/financial-cli/models"
)

func PrintSingleRecord(record m.Record, color string) {
	date := fmt.Sprintf("    %s/%s/%s   ", GetStringDateFromNumber(record.Month), GetStringDateFromNumber(record.Day), strconv.Itoa(record.Year))
	description := fmt.Sprintf(" %-38s", record.Description)
	costString := fmt.Sprintf(" $%-6s", strconv.Itoa(record.Cost))
	category := fmt.Sprintf(" %-18s", record.Category)

	message := fmt.Sprintf("\n%s\n|%s|%s|%s|%s|\n%s\n", m.DASH, date, description, costString, category, m.DASH)
	PrintCustomizedMessage(message, color, true)
}

func PrintSingleTripRecord(record m.TripRecord, color string) {
	date := fmt.Sprintf("    %s/%s/%s   ", GetStringDateFromNumber(record.Month), GetStringDateFromNumber(record.Day), strconv.Itoa(record.Year))
	description := fmt.Sprintf(" %-35s", record.Description)
	costString := fmt.Sprintf(" $%-6s", strconv.Itoa(record.Cost))
	category := fmt.Sprintf(" %-18s", record.Category)

	message := fmt.Sprintf("\n%s\n|%s|%s|%s|%s|\n%s\n", m.DASH, date, description, costString, category, m.DASH)
	PrintCustomizedMessage(message, color, true)
}

func PrintTable(data [][]interface{}, headers []string, typeFlag string, style *simpletable.Style) {
	table := simpletable.New()

	// Generate Table Header
	for _, header := range headers {
		headerCell := simpletable.Cell{Align: simpletable.AlignCenter, Text: Colorize(header, Green)}
		table.Header.Cells = append(table.Header.Cells, &headerCell)
	}

	// Generate Table Body
	totalIncome := 0
	totalExpense := 0
	totalShared := 0
	n := len(headers)
	for _, rowData := range data {
		row := []*simpletable.Cell{}
		isIncome := string(rowData[n-2].(string)) == "Income"
		isShared := strings.Contains(rowData[n-2].(string), "shared")
		for idx, rowCellData := range rowData {
			var rowCell *simpletable.Cell
			if _, ok := rowCellData.(int); ok {
				rowCell = &simpletable.Cell{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", rowCellData)}
			} else {
				rowCell = &simpletable.Cell{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", rowCellData)}
			}
			if idx == n-1 {
				if isIncome {
					totalIncome += int(rowCellData.(int))
					rowCell.Text = Colorize(rowCell.Text, Yellow)
				} else {
					totalExpense += int(rowCellData.(int))
					if isShared {
						totalShared += int(rowCellData.(int))
					}
				}
			}
			row = append(row, rowCell)
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}

	// Generate Table Footer
	var footer *simpletable.Footer
	switch typeFlag {
	case "income":
		footer = &simpletable.Footer{
			Cells: []*simpletable.Cell{
				{
					Align: simpletable.AlignRight,
					Span:  4,
					Text:  fmt.Sprintf("%s", Colorize("TOTAL INCOME", Red)),
				},
				{
					Align: simpletable.AlignRight,
					Text:  fmt.Sprintf("%s", Colorize(fmt.Sprintf("%d", totalIncome), Yellow)),
				},
			},
		}
	case "expense":
		footer = &simpletable.Footer{
			Cells: []*simpletable.Cell{
				{
					Align: simpletable.AlignRight,
					Span:  4,
					Text:  fmt.Sprintf("%s", Colorize("TOTAL EXPENSE", Red)),
				},
				{
					Align: simpletable.AlignRight,
					Text:  fmt.Sprintf("%s", Colorize(fmt.Sprintf("%d", totalExpense), White)),
				},
			},
		}
	case "all":
		footer = &simpletable.Footer{
			Cells: []*simpletable.Cell{
				{
					Align: simpletable.AlignRight,
					Span:  4,
					Text: fmt.Sprintf("%s\n%s",
						Colorize("T-INCOME", Yellow),
						Colorize("T-EXPENSE", Red)),
				},
				{
					Align: simpletable.AlignRight,
					Text: fmt.Sprintf("%s\n%s",
						Colorize(fmt.Sprintf("%d", totalIncome), Yellow),
						Colorize(fmt.Sprintf("%d", totalExpense), White)),
				},
			},
		}
	case "trip":
		footer = &simpletable.Footer{
			Cells: []*simpletable.Cell{
				{
					Align: simpletable.AlignRight,
					Span:  4,
					Text: fmt.Sprintf("%s\n%s",
						Colorize("T-SHARED", Yellow),
						Colorize("T-EXPENSE", Red)),
				},
				{
					Align: simpletable.AlignRight,
					Text: fmt.Sprintf("%s\n%s",
						Colorize(fmt.Sprintf("%d", totalShared), Yellow),
						Colorize(fmt.Sprintf("%d", totalExpense), White)),
				},
			},
		}
	default:
		footer = &simpletable.Footer{
			Cells: []*simpletable.Cell{
				{
					Align: simpletable.AlignCenter,
					Span:  4,
					Text:  fmt.Sprintf("%s", Colorize("__________ ^o^ __________", Red)),
				},
			},
		}
	}

	table.Footer = footer
	table.SetStyle(style)
	table.Println()
}

func AddRecordToFile(filepath string, record m.Record, color string) {
	success := CsvWrite(filepath, record)
	if success && color != "" {
		PrintCustomizedMessage("Record has been successfully added at "+filepath, color, true)
	}
}

func AddRecord(record m.Record) {
	sharedFile := GetSharedFile()
	currentYearFile := GetSpecificYearFile(time.Now().Year())
	AddRecordToFile(sharedFile, record, Yellow)
	AddRecordToFile(currentYearFile, record, Red)
}
