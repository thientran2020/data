package app

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	m "github.com/thientran2020/financial-cli/models"
	u "github.com/thientran2020/financial-cli/utils"
)

func HandleAdd(addCmd *flag.FlagSet, sub *bool) {
	addCmd.Parse(os.Args[2:])
	if addCmd.NArg() != 0 {
		addCmd.PrintDefaults()
		os.Exit(1)
	}

	if *sub == true {
		u.AddSubscription()
		return
	}

	// Get current date
	year, month, day := time.Now().Date()

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
	u.PrintSingleRecord(record, u.Green)

	// Confirm record and enter to files
	if u.ConfirmYesNoPromt("Do you confirm to enter above record") {
		sharedFile := u.GetSharedFile()
		currentYearFile := u.GetSpecificYearFile(time.Now().Year())
		u.AddRecord(sharedFile, record, u.Yellow)
		u.AddRecord(currentYearFile, record, u.Red)
	} else {
		u.PrintCustomizedMessage("Record ignored "+u.CheckMark, u.Red, true)
	}
}

func HandleShow(showCmd *flag.FlagSet, month, year *int, current, income, expense *bool, keyword *string) {
	// Update subscriptions
	u.UpdateSubscriptionRecord()

	// Handle Show command
	showCmd.Parse(os.Args[2:])
	if showCmd.NArg() != 0 {
		showCmd.PrintDefaults()
		os.Exit(1)
	}

	if *year != -1 && (*year < m.START_YEAR || *year > time.Now().Year()) {
		fmt.Println(u.Colorize("No data found for the requested year...!", u.Red))
		return
	}

	if *current == true {
		*month = int(time.Now().Month())
		*year = time.Now().Year()
	}

	var flag string
	if *income == true && *expense == false {
		flag = "income"
	} else if *income == false && *expense == true {
		flag = "expense"
	} else {
		flag = "all"
	}

	// Choose file for retrieving financial data
	filepath := u.GetSharedFile()
	if *year >= m.START_YEAR && *year <= time.Now().Year() {
		filepath = u.GetSpecificYearFile(*year)
	}

	// Retrieve, filter and display data
	data, _ := u.CsvRead(filepath)
	filteredData := u.FilterData(data, *month, flag, *keyword)
	u.PrintTable(filteredData, m.HEADERS, flag, simpletable.StyleDefault)
}

func HandleHelp(helpCmd *flag.FlagSet) {
	helpCmd.Parse(os.Args[2:])

	if helpCmd.NFlag() > 0 || helpCmd.NArg() > 0 {
		fmt.Println("Please don't specific any argument/flag.")
		fmt.Println("Correct usage: 'data help'")
		return
	}
	fmt.Print(m.INSTRUCTION)
}

func HandleGet(getCmd *flag.FlagSet, getCategory, getSubscription *bool) {
	getCmd.Parse(os.Args[2:])
	if getCmd.NArg() > 0 || getCmd.NFlag() != 1 {
		getCmd.PrintDefaults()
		os.Exit(1)
	}

	if *getCategory == true {
		u.PrintCustomizedMessage(m.CATEGORY_TABLE, u.White, true)
		return
	}

	if *getSubscription == true {
		data := u.GetSubscription()
		u.PrintSubcriptionList("monthly", data.Monthly)
		u.PrintSubcriptionList("yearly", data.Yearly)
		fmt.Println()
		return
	}
}

func HandleSearch(searchCmd *flag.FlagSet) {
	searchCmd.Parse(os.Args[2:])

	if searchCmd.NArg() < 1 {
		fmt.Println("Please specific keyword. Correct usage: 'data search keyword'")
		return
	}

	keyword := strings.Join(os.Args[2:], " ")
	data, _ := u.CsvRead(u.GetSharedFile())
	filteredData := u.FilterData(data, -1, "all", keyword)
	u.PrintTable(filteredData, m.HEADERS, "all", simpletable.StyleDefault)
}
