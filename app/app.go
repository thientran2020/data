package app

import (
	"flag"
	"fmt"
	"os"

	m "github.com/thientran2020/financial-cli/models"
)

func AppInit() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addSubscription := addCmd.Bool("s", false, m.AddSubscriptionMessage)

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	showCurrent := showCmd.Bool("c", false, m.ShowCurrentMessage)
	showMonth := showCmd.Int("m", -1, m.ShowMonthMessage)
	showYear := showCmd.Int("y", -1, m.ShowYearMessage)
	showIncome := showCmd.Bool("i", false, m.ShowIncomeMessage)
	showExpense := showCmd.Bool("e", false, m.ShowExpenseMessage)
	showKeyword := showCmd.String("k", "", m.ShowKeywordMessage)

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	getCmd := flag.NewFlagSet("get", flag.ExitOnError)
	getCategory := getCmd.Bool("c", false, m.GetCategoryMessage)
	getSubscription := getCmd.Bool("s", false, m.GetSubscriptionMessage)

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected at least 1 subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		HandleAdd(addCmd, addSubscription)
	case "show":
		HandleShow(showCmd, showMonth, showYear, showCurrent, showIncome, showExpense, showKeyword)
	case "help":
		HandleHelp(helpCmd)
	case "get":
		HandleGet(getCmd, getCategory, getSubscription)
	case "search":
		HandleSearch(searchCmd)
	default:
		fmt.Print(m.INSTRUCTION)
	}
}
