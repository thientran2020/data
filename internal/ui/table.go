package ui

import (
	"fmt"
	"sort"

	"github.com/alexeyco/simpletable"
	"github.com/thientran2020/financial-cli/internal/models"
)

const (
	Reset      = "\033[0m"
	Bold       = "\033[1m"
	Red        = "\033[31m"
	Green      = "\033[32m"
	Yellow     = "\033[33m"
	Blue       = "\033[34m"
	Magenta    = "\033[35m"
	Cyan       = "\033[36m"
	White      = "\033[37m"
	BoldRed    = "\033[1;31m"
	BoldGreen  = "\033[1;32m"
	BoldYellow = "\033[1;33m"
	BoldBlue   = "\033[1;34m"
	BoldWhite  = "\033[1;37m"
)

type Table struct {
	// Any configuration or dependencies can be added here
}

func NewTable() *Table {
	return &Table{}
}

// DisplayRecords displays a table of records
func (t *Table) DisplayRecords(records []*models.Record, title string) {
	// Sort records by date in ascending order
	sort.Slice(records, func(i, j int) bool {
		return records[i].Date.Before(records[j].Date)
	})

	table := simpletable.New()

	// Set headers
	table.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "#" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "DATE" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "DESCRIPTION" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "CATEGORY" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "COST" + Reset},
		},
	}

	var incomeTotal, expenseTotal int
	for i, rec := range records {
		var cost string
		if rec.IsIncome() {
			cost = Green + fmt.Sprintf("%d", rec.Cost) + Reset
			incomeTotal += rec.Cost
		} else {
			cost = White + fmt.Sprintf("%d", rec.Cost) + Reset
			expenseTotal += rec.Cost
		}

		table.Body.Cells = append(table.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", i+1)},
			{Align: simpletable.AlignCenter, Text: rec.FormattedDate()},
			{Align: simpletable.AlignLeft, Text: rec.Description},
			{Align: simpletable.AlignLeft, Text: rec.Category},
			{Align: simpletable.AlignRight, Text: cost},
		})
	}

	table.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Span: 4, Text: Bold + "TOTAL INCOME:" + Reset},
			{Align: simpletable.AlignRight, Text: BoldGreen + fmt.Sprintf("$%d", incomeTotal) + Reset},
		},
	}

	if title != "" {
		fmt.Printf("\n%s%s%s:\n\n", Bold, title, Reset)
	}

	table.SetStyle(simpletable.StyleDefault)
	fmt.Println(table.String())

	// Format totals with commas and align them
	net := incomeTotal - expenseTotal
	netColor := BoldGreen
	if net < 0 {
		netColor = BoldRed
	}

	fmt.Printf("\n%-25s %s$%d%s\n", Bold+"TOTAL INCOME:"+Reset, BoldGreen, incomeTotal, Reset)
	fmt.Printf("%-25s %s$%d%s\n", Bold+"TOTAL EXPENSE:"+Reset, BoldRed, expenseTotal, Reset)
	fmt.Printf("%-25s %s$%d%s\n\n", Bold+"NET:"+Reset, netColor, net, Reset)
}

func (t *Table) DisplayRecordsByType(records []*models.Record, recordType string, title string) {
	var filtered []*models.Record

	// Filter records by type
	for _, rec := range records {
		if (recordType == "income" && rec.IsIncome()) ||
			(recordType == "expense" && rec.IsExpense()) ||
			recordType == "all" {
			filtered = append(filtered, rec)
		}
	}

	t.DisplayRecords(filtered, title)
}

func (t *Table) DisplaySubscriptions(subscriptions *models.SubscriptionList) {
	monthlyTable := simpletable.New()
	monthlyTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: Bold + "#" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + "NAME" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + "TYPE" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + "COST" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + "START DATE" + Reset},
		},
	}

	var monthlyTotal int
	for i, sub := range subscriptions.Monthly {
		cost := fmt.Sprintf("$%d", sub.Cost)
		if sub.Type == models.Income {
			cost = Green + cost + Reset
		} else {
			cost = Red + cost + Reset
		}
		monthlyTotal += sub.Cost

		monthlyTable.Body.Cells = append(monthlyTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", i+1)},
			{Align: simpletable.AlignLeft, Text: sub.Name},
			{Align: simpletable.AlignCenter, Text: string(sub.Type)},
			{Align: simpletable.AlignRight, Text: cost},
			{Align: simpletable.AlignCenter, Text: sub.StartDate},
		})
	}

	monthlyTable.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Span: 3, Text: Bold + "TOTAL MONTHLY:" + Reset},
			{Align: simpletable.AlignRight, Text: BoldRed + fmt.Sprintf("$%d", monthlyTotal) + Reset},
			{Align: simpletable.AlignCenter, Text: ""},
		},
	}

	yearlyTable := simpletable.New()
	yearlyTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: Bold + "#" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + "NAME" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + "TYPE" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + "COST" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + "START DATE" + Reset},
		},
	}

	var yearlyTotal int
	for i, sub := range subscriptions.Yearly {
		cost := fmt.Sprintf("$%d", sub.Cost)
		if sub.Type == models.Income {
			cost = Green + cost + Reset
		} else {
			cost = Red + cost + Reset
		}
		yearlyTotal += sub.Cost

		yearlyTable.Body.Cells = append(yearlyTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", i+1)},
			{Align: simpletable.AlignLeft, Text: sub.Name},
			{Align: simpletable.AlignCenter, Text: string(sub.Type)},
			{Align: simpletable.AlignRight, Text: cost},
			{Align: simpletable.AlignCenter, Text: sub.StartDate},
		})
	}

	yearlyTable.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Span: 3, Text: Bold + "TOTAL YEARLY:" + Reset},
			{Align: simpletable.AlignRight, Text: BoldRed + fmt.Sprintf("$%d", yearlyTotal) + Reset},
			{Align: simpletable.AlignCenter, Text: ""},
		},
	}

	fmt.Printf("\n%sMONTHLY SUBSCRIPTIONS:%s\n", Bold, Reset)
	monthlyTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(monthlyTable.String())

	fmt.Printf("\n%sYEARLY SUBSCRIPTIONS:%s\n", Bold, Reset)
	yearlyTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(yearlyTable.String())

	yearlyMonthly := yearlyTotal / 12
	fmt.Printf("\n%-25s %s$%d%s\n", Bold+"YEARLY TOTAL:"+Reset, BoldRed, yearlyTotal, Reset)
	fmt.Printf("%-25s %s$%d%s\n", Bold+"YEARLY AVERAGE PER MONTH:"+Reset, BoldRed, yearlyMonthly, Reset)
	fmt.Printf("%-25s %s$%d%s\n\n", Bold+"TOTAL MONTHLY COST:"+Reset, BoldRed, monthlyTotal+yearlyMonthly, Reset)
}

func (t *Table) DisplayTrip(trip *models.Trip) {
	fmt.Printf("\n%-25s %s%s%s\n", Bold+"TRIP:"+Reset, BoldYellow, trip.Name, Reset)
	fmt.Printf("%-25s %s%s%s\n", Bold+"DATE:"+Reset, Bold, trip.DateRange(), Reset)
	fmt.Printf("%-25s %s%d%s\n", Bold+"PARTICIPANTS:"+Reset, Bold, trip.NParticipants, Reset)
	fmt.Printf("%-25s %s$%d%s\n", Bold+"TOTAL COST:"+Reset, BoldGreen, trip.Costs.Total, Reset)
	fmt.Printf("%-25s %s$%d%s\n", Bold+"SHARED COST:"+Reset, BoldGreen, trip.Costs.Shared, Reset)

	perPersonCost := trip.CalculatePerPersonCost()
	fmt.Printf("%-25s %s$%d%s\n\n", Bold+"COST PER PERSON:"+Reset, BoldRed, perPersonCost, Reset)

	recordTable := simpletable.New()
	recordTable.Header = &simpletable.Header{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "#" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "DATE" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "DESCRIPTION" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "COST" + Reset},
			{Align: simpletable.AlignCenter, Text: Bold + Yellow + "SHARED" + Reset},
		},
	}

	for i, rec := range trip.Records {
		var record *models.Record
		if rec.Record != nil {
			record = rec.Record
		} else {
			record = &models.Record{
				Year:        rec.Year,
				Month:       rec.Month,
				Day:         rec.Day,
				Description: rec.Description,
				Cost:        rec.Cost,
				Category:    rec.Category,
				Code:        rec.Code,
			}
		}

		shared := "No"
		if rec.Shared {
			shared = BoldGreen + "Yes" + Reset
		}

		recordTable.Body.Cells = append(recordTable.Body.Cells, []*simpletable.Cell{
			{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", i+1)},
			{Align: simpletable.AlignCenter, Text: record.FormattedDate()},
			{Align: simpletable.AlignLeft, Text: record.Description},
			{Align: simpletable.AlignRight, Text: White + fmt.Sprintf("%d", rec.Cost) + Reset},
			{Align: simpletable.AlignCenter, Text: shared},
		})
	}

	recordTable.Footer = &simpletable.Footer{
		Cells: []*simpletable.Cell{
			{Align: simpletable.AlignRight, Span: 3, Text: Bold + "TOTAL:" + Reset},
			{Align: simpletable.AlignRight, Text: BoldGreen + fmt.Sprintf("$%d", trip.Costs.Total) + Reset},
			{Align: simpletable.AlignCenter, Text: ""},
		},
	}

	fmt.Println(Bold + "TRIP RECORDS:" + Reset)
	recordTable.SetStyle(simpletable.StyleDefault)
	fmt.Println(recordTable.String())
}

func (t *Table) DisplayCategoryTable(categories []string) {
	fmt.Println("\nCATEGORY TABLE:")
	fmt.Println("|------|------------------------------------------|")

	for i, category := range categories {
		fmt.Printf("|  %-4d|  %-40s|\n", i, category)
		fmt.Println("|------|------------------------------------------|")
	}
}
