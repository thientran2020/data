package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/thientran2020/financial-cli/models"
)

func PrintSingleRecord(record models.Record, color string) {
	date := fmt.Sprintf("    %s/%s/%s   ", GetStringDateFromNumber(record.Month), GetStringDateFromNumber(record.Day), strconv.Itoa(record.Year))
	description := fmt.Sprintf(" %-35s", record.Description)
	costString := fmt.Sprintf(" $%-6s", strconv.Itoa(record.Cost))
	category := fmt.Sprintf(" %-18s", record.Category)

	message := fmt.Sprintf("\n%s\n|%s|%s|%s|%s|\n%s\n", models.DASH, date, description, costString, category, models.DASH)
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
	n := len(headers)
	for _, rowData := range data {
		row := []*simpletable.Cell{}
		isIncome := string(rowData[n-2].(string)) == "Income"
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
	default:
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
	}

	table.Footer = footer
	table.SetStyle(style)
	fmt.Println(Colorize(models.HEADER_LINE, Green))
	table.Println()
}

func AddRecord(filepath string, record models.Record, color string) {
	if CsvWrite(filepath, record) {
		PrintCustomizedMessage("Record has been successfully added at "+filepath, color, true)
	}
}

func AddSubscription() {
	subscriptionList := ReadJson(models.BASE_FILEPATH_SUBCRIPTION)

	// Prompt user to enter neccessary information
	startDate := strings.Split(time.Now().String(), " ")[0]
	name, _ := PromptEnter("What is your new subscription/membership", false)
	ftype, _ := InteractiveSelect(
		"What type of your subscription",
		[]string{"income", "expense"},
	)
	billingCycle, _ := InteractiveSelect(
		"Choose your billing cycle",
		[]string{"monthly", "yearly"},
	)
	cost, _ := NumberEnter("How much per billing period")

	// Create new subscription and add to existing list
	subscription := models.Subscription{
		Name:         name,
		Type:         ftype,
		Cost:         int(cost),
		BillingCycle: billingCycle,
		StartDate:    startDate,
	}

	switch billingCycle {
	case "monthly":
		subscriptionList.Monthly = append(subscriptionList.Monthly, subscription)
	case "yearly":
		subscriptionList.Yearly = append(subscriptionList.Yearly, subscription)
	}

	// Print new subscription and ask for confirmation before adding
	message := fmt.Sprintf("%s: $%d/%s", name, cost, strings.ToLower(billingCycle[:len(billingCycle)-2]))
	PrintCustomizedMessage(message, Green, true)
	confirmed := ConfirmYesNoPromt("Do you confirm to enter above subscription")
	if confirmed {
		WriteJson(models.BASE_FILEPATH_SUBCRIPTION, subscriptionList)
		PrintCustomizedMessage("Successfully added at "+models.BASE_FILEPATH_SUBCRIPTION, Yellow, true)
	} else {
		PrintCustomizedMessage("Subscription ignored "+CheckMark, Red, true)
	}
}

// TODOs
func UpdateSubscription() {

}
