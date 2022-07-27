package app

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

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
		fmt.Println("Adding subscription....")
		AddSubscription()
		return
	}

	// Get current date

	year, month, day := time.Now().Date()

	// Create filepath for all-in-one file and current year file
	// Check if file exists - if not create a new one
	filepathCommon := "./finance/finance.csv"
	if !fileExists(filepathCommon) {
		createFile(filepathCommon)
	}

	filepathCurrentYear := "./finance/finance-" + strconv.Itoa(year) + ".csv"
	if !fileExists(filepathCurrentYear) {
		createFile(filepathCurrentYear)
	}

	// Prompt to input data
	// 1. Check data entered is for expense or income
	ftype, _ := utils.InteractiveSelect(
		"What type of financial data are you entering",
		[]string{"Income", "Expense"},
	)
	// 2. Prompt enter description
	description, _ := utils.PromptEnter(models.LABELS[ftype]["Description"], false)
	// 3. Get $$$ spent
	cost, _ := utils.NumberEnter(models.LABELS[ftype]["Cost"])
	// 4. Choose category
	category, _ := utils.InteractiveSelect(
		"Pick the category that describe best your entered data",
		models.CATEGORY,
	)
	// 5. Convert category to code
	var code int
	for index := range models.CATEGORY {
		if models.CATEGORY[index] == category {
			code = index
		}
	}

	record := models.Record{
		Year:        year,
		Month:       int(month),
		Day:         day,
		Description: description,
		Cost:        int(cost),
		Category:    category,
		Code:        code,
	}
	utils.PrintSingleRecord(record, utils.BGreen)

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

func HandleShow(showCmd *flag.FlagSet, month *int, year *int) {
	showCmd.Parse(os.Args[2:])

	if showCmd.NArg() != 0 {
		showCmd.PrintDefaults()
		os.Exit(1)
	}

	if *month == -1 {
		*month = int(time.Now().Month())
	}
	if *year == -1 {
		*year = time.Now().Year()
	}
	fmt.Printf("Show financial data for %d, %d\n", *month, *year)
	// TODO
	// Retrieve financial data for specific date
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
	filepath := "./finance/subscription.json"
	// if !fileExists(filepath) {
	// 	createFile(filepath)
	// }
	subscriptionList := readJson(filepath)

	// Prompt user to enter neccessary information
	startDate := strings.Split(time.Now().String(), " ")[0]
	name, _ := utils.PromptEnter("What is your new subscription/membership", false)
	details, _ := utils.PromptEnter("Any detail you want to provide", true)
	billingCycle, _ := utils.InteractiveSelect(
		"Choose your billing cycle",
		[]string{"Monthly", "Yearly"},
	)
	cost, _ := utils.NumberEnter("How much per billing period")

	// Create new subscription and add to existing list
	subscription := models.Subscription{
		Name:         name,
		Details:      details,
		Cost:         int(cost),
		BillingCycle: billingCycle,
		StartDate:    startDate,
	}

	switch billingCycle {
	case "Monthly":
		subscriptionList.Monthly = append(subscriptionList.Monthly, subscription)
	case "Yearly":
		subscriptionList.Yearly = append(subscriptionList.Yearly, subscription)
	}

	// Print new subscription and ask for confirmation before adding
	message := name + ": $" + strconv.Itoa(int(cost)) + "/" + strings.ToLower(billingCycle[:len(billingCycle)-2])
	utils.PrintCustomizedMessage(message, utils.BGreen, true)
	confirmed := utils.ConfirmYesNoPromt("Do you confirm to enter above subscription")
	if confirmed {
		writeJson(filepath, subscriptionList)
		utils.PrintCustomizedMessage("Successfully added at "+filepath, utils.BYellow, true)
	} else {
		utils.PrintCustomizedMessage("Subscription ignored "+utils.CheckMark, utils.BRed, true)
	}
}
