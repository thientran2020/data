package app

import (
	"fmt"
	"time"

	"github.com/alexeyco/simpletable"
	m "github.com/thientran2020/financial-cli/models"
	u "github.com/thientran2020/financial-cli/utils"
)

func HandleAdd(cmd *CLI) {
	if cmd.Add.Subscription {
		u.AddSubscription()
		return
	}

	if cmd.Add.Trip {
		u.AddNewTrip()
		return
	}

	// By default, "add" for expense
	ftype := "Expense"
	if cmd.Add.Income {
		ftype = "Income"
	}

	// Get date - default is current date
	var year, month, day int
	date := time.Now().Format("01-02-2006")
	if cmd.Add.Yesterday {
		date = time.Now().AddDate(0, 0, -1).Format("01-02-2006")
	} else if !cmd.Add.Today {
		date = u.DateEnter("Enter date: ")
	}
	month, day, year = u.GetDateNumber(date)

	// Prompt to input data if no specific flags were passed
	// 1. Prompt enter description
	// 2. Get $$$ spent
	// 3. Choose category
	// 4. Convert category to code

	description := cmd.Add.Description
	if description == "" {
		description = u.PromptEnter(m.LABELS[ftype]["Description"])
	}

	cost := cmd.Add.Cost
	if cost == 0 {
		cost = u.NumberEnter(m.LABELS[ftype]["Cost"])
	}

	code := cmd.Add.Category
	category := m.CATEGORY[code]
	if (code == 0) && !(cmd.Add.Income) {
		category = u.InteractiveSelect(
			"Pick the category that describe best your entered data!",
			m.CATEGORY,
		)
		for index := range m.CATEGORY {
			if m.CATEGORY[index] == category {
				code = index
			}
		}
	}

	// Create record and ask for confirmation before adding
	record := m.Record{
		Year:        year,
		Month:       int(month),
		Day:         day,
		Description: description,
		Cost:        int(cost),
		Category:    category,
		Code:        code,
	}

	if category == "Trip" {
		u.AddTripRecord(record)
	} else {
		u.PrintSingleRecord(record, u.Green)
		if cmd.Add.Yes {
			u.AddRecord(record)
		} else if u.ConfirmYesNoPromt("Do you confirm to enter above record?") {
			u.AddRecord(record)
		} else {
			u.PrintCustomizedMessage("Record ignored "+u.CheckMark, u.Red, true)
		}
	}
}

func HandleShow(cmd *CLI) {
	// Update subscriptions
	u.UpdateSubscriptionRecord()

	year := cmd.Show.Year
	month := cmd.Show.Month
	current := cmd.Show.Current
	income := cmd.Show.Income
	expense := cmd.Show.Expense
	keyword := cmd.Show.Keyword

	if year != 0 && (year < m.START_YEAR || year > time.Now().Year()) {
		fmt.Println(u.Colorize("No data found for the requested year...!", u.Red))
		return
	}

	if current {
		month = int(time.Now().Month())
		year = time.Now().Year()
	}

	var flag string
	if income && !expense {
		flag = "income"
	} else if !income && expense {
		flag = "expense"
	} else {
		flag = "all"
	}

	// Choose file for retrieving financial data
	filepath := u.GetSharedFile()
	if year >= m.START_YEAR && year <= time.Now().Year() {
		filepath = u.GetSpecificYearFile(year)
	}

	// Retrieve, filter and display data
	data, _ := u.CsvRead(filepath)
	filteredData := u.FilterData(data, month, flag, keyword)
	u.PrintTable(filteredData, m.HEADERS, flag, simpletable.StyleDefault)
}

func HandleGet(cmd *CLI) {
	if cmd.Get.Trip {
		onlyShared := u.ConfirmYesNoPromt("Would you like to display shared (Y) or all expenses (N)?")
		u.PrintTrip(onlyShared)
	} else if !cmd.Get.Subscription && !cmd.Get.Category {
		fmt.Print(m.INSTRUCTION)
	} else if cmd.Get.Category {
		u.PrintCustomizedMessage(m.CATEGORY_TABLE, u.ColorOff, true)
	} else if cmd.Get.Subscription {
		data := u.GetSubscription()
		u.PrintSubcriptionList("monthly", data.Monthly)
		u.PrintSubcriptionList("yearly", data.Yearly)
		fmt.Println()
	}
}

func HandleSearch(keyword string) {
	data, _ := u.CsvRead(u.GetSharedFile())
	filteredData := u.FilterData(data, 0, "all", keyword)
	u.PrintTable(filteredData, m.HEADERS, "all", simpletable.StyleDefault)
}
