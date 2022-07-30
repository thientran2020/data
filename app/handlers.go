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
		AddSubscription()
		return
	}

	// Get current date
	year, month, day := time.Now().Date()

	// Create filepath for all-in-one file and current year file
	// Check if file exists - if not create a new one
	filepathCommon := strings.Replace(models.BASE_FILEPATH, "?????", "", -1)
	if !fileExists(filepathCommon) {
		createFile(filepathCommon)
	}

	filepathCurrentYear := strings.Replace(models.BASE_FILEPATH, "?????", fmt.Sprintf("_%d", year), -1)
	if !fileExists(filepathCurrentYear) {
		createFile(filepathCurrentYear)
	}

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
	utils.PrintSingleRecord(record, utils.BWhite)

	confirmed := utils.ConfirmYesNoPromt("Do you confirm to enter above record")
	if confirmed {
		if csvWrite(filepathCurrentYear, record) {
			utils.PrintCustomizedMessage(
				"Record has been successfully added at "+filepathCurrentYear,
				utils.BRed,
				true,
			)
		}
		if csvWrite(filepathCommon, record) {
			utils.PrintCustomizedMessage(
				"Record has been successfully added at "+filepathCommon,
				utils.BYellow,
				true,
			)
		}
	}
}

func HandleShow(showCmd *flag.FlagSet, month *int, year *int, income *bool, expense *bool) {
	showCmd.Parse(os.Args[2:])

	if showCmd.NArg() != 0 {
		showCmd.PrintDefaults()
		os.Exit(1)
	}

	if *year == -1 {
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

	fmt.Printf("%s %d/%d: \n\n",
		utils.Colorize("\nSHOW FINANCIAL DATA FOR", utils.UWhite),
		*month,
		*year,
	)

	data := csvRead(*year, *month, flag)
	headers := []string{"#", "DATE", "DESCRIPTION", "CATEGORY", "COST"}
	utils.PrintTable(data, headers, flag, simpletable.StyleDefault)
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

func AddSubscription() {
	subscriptionList := readJson(models.BASE_FILEPATH_SUBCRIPTION)

	// Prompt user to enter neccessary information
	startDate := strings.Split(time.Now().String(), " ")[0]
	name, _ := utils.PromptEnter("What is your new subscription/membership", false)
	ftype, _ := utils.InteractiveSelect(
		"What type of your subscription",
		[]string{"income", "expense"},
	)
	billingCycle, _ := utils.InteractiveSelect(
		"Choose your billing cycle",
		[]string{"monthly", "yearly"},
	)
	cost, _ := utils.NumberEnter("How much per billing period")

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
	utils.PrintCustomizedMessage(message, utils.BGreen, true)
	confirmed := utils.ConfirmYesNoPromt("Do you confirm to enter above subscription")
	if confirmed {
		writeJson(models.BASE_FILEPATH_SUBCRIPTION, subscriptionList)
		utils.PrintCustomizedMessage("Successfully added at "+models.BASE_FILEPATH_SUBCRIPTION, utils.BYellow, true)
	} else {
		utils.PrintCustomizedMessage("Subscription ignored "+utils.CheckMark, utils.BRed, true)
	}
}
