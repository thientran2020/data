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

	// Get date - default is current date
	var year, month, day int
	var date string
	if date, _ = u.DateEnter("Enter date: "); date == "" {
		date = time.Now().Format("01-02-2006")
	}
	month, day, year = u.GetDateNumber(date)

	// Prompt to input data
	// 1. Check data entered is for expense or income
	// 2. Prompt enter description
	// 3. Get $$$ spent
	// 4. Choose category
	// 5. Convert category to code
	ftype, _ := u.InteractiveSelect(
		"What type of financial data are you entering",
		[]string{"Income", "Expense"},
	)
	description, _ := u.PromptEnter(m.LABELS[ftype]["Description"], false)
	cost, _ := u.NumberEnter(m.LABELS[ftype]["Cost"])
	category, _ := u.InteractiveSelect(
		"Pick the category that describe best your entered data",
		m.CATEGORY,
	)
	var code int
	for index := range m.CATEGORY {
		if m.CATEGORY[index] == category {
			code = index
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

	if record.Category == "Trip" {
		u.AddTripRecord(record)
	} else {
		u.PrintSingleRecord(record, u.Green)
		if u.ConfirmYesNoPromt("Do you confirm to enter above record") {
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
		u.PritnTrip()
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
