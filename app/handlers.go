package app

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	"github.com/thientran2020/financial-cli/models"
	"github.com/thientran2020/financial-cli/utils"
)

func HandleAdd(addCmd *flag.FlagSet, sub *bool) {
	addCmd.Parse(os.Args[2:])

	if addCmd.NArg() != 0 {
		addCmd.PrintDefaults()
		os.Exit(1)
	}

	if *sub == true {
		utils.AddSubscription()
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
	ftype, _ := utils.InteractiveSelect(
		"What type of financial data are you entering",
		[]string{"Income", "Expense"},
	)
	description, _ := utils.PromptEnter(models.LABELS[ftype]["Description"], false)
	cost, _ := utils.NumberEnter(models.LABELS[ftype]["Cost"])
	category, _ := utils.InteractiveSelect(
		"Pick the category that describe best your entered data",
		models.CATEGORY,
	)
	var code int
	for index := range models.CATEGORY {
		if models.CATEGORY[index] == category {
			code = index
		}
	}

	// Create record and ask for confirmation before adding
	record := models.Record{
		Year:        year,
		Month:       int(month),
		Day:         day,
		Description: description,
		Cost:        int(cost),
		Category:    category,
		Code:        code,
	}
	utils.PrintSingleRecord(record, utils.Green)

	// Confirm record and enter to files
	if utils.ConfirmYesNoPromt("Do you confirm to enter above record") {
		sharedFile := utils.GetSharedFile()
		currentYearFile := utils.GetCurrentYearFile()
		utils.AddRecord(sharedFile, record, utils.Yellow)
		utils.AddRecord(currentYearFile, record, utils.Red)
	} else {
		utils.PrintCustomizedMessage("Record ignored "+utils.CheckMark, utils.Red, true)
	}
}

func HandleShow(showCmd *flag.FlagSet, current *bool, month *int, year *int, income *bool, expense *bool, keyword *string) {
	showCmd.Parse(os.Args[2:])

	if showCmd.NArg() != 0 {
		showCmd.PrintDefaults()
		os.Exit(1)
	}

	if *year != -1 && (*year < models.START_YEAR || *year > time.Now().Year()) {
		fmt.Println(utils.Colorize("No data found for the requested year...!", utils.Red))
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

	data := utils.CsvRead(*year, *month, flag, *keyword)
	utils.PrintTable(data, models.HEADERS, flag, simpletable.StyleDefault)
}

func HandleHelp(helpCmd *flag.FlagSet) {
	helpCmd.Parse(os.Args[2:])

	if helpCmd.NFlag() > 0 || helpCmd.NArg() > 0 {
		fmt.Println("Please don't specific any argument/flag.")
		fmt.Println("Correct usage: 'data help'")
		return
	}
	fmt.Print(models.INSTRUCTION)
}

func HandleCategory(ctgCmd *flag.FlagSet) {
	ctgCmd.Parse(os.Args[2:])

	if ctgCmd.NFlag() > 0 || ctgCmd.NArg() > 0 {
		fmt.Println("Please don't specific any argument/flag.")
		fmt.Println("Correct usage: 'data help'")
		return
	}
	fmt.Print(models.CATEGORY_TABLE)
}

func HandleSearch(searchCmd *flag.FlagSet) {
	searchCmd.Parse(os.Args[2:])

	if searchCmd.NArg() < 1 {
		fmt.Println("Please specific keyword. Correct usage: 'data search keyword'")
		return
	}

	keyword := strings.Join(os.Args[2:], " ")
	data := utils.CsvRead(-1, -1, "all", keyword)
	utils.PrintTable(data, models.HEADERS, "all", simpletable.StyleDefault)
}
