package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/thientran2020/financial-cli/models"
)

func AppInit() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addSubscription := addCmd.Bool("s", false, "Specify when adding subscription/membership data")

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	showCurrent := showCmd.Bool("c", false, "Specify when you want to retrieve current month data")
	showMonth := showCmd.Int("m", -1, "Specify month you want to retrive financial data")
	showYear := showCmd.Int("y", -1, "Specify year you want to retrive financial data")
	showIncome := showCmd.Bool("i", false, "Specify when you want to retrieve income data")
	showExpense := showCmd.Bool("e", false, "Specify when you want to retrieve income data")
	showKeyword := showCmd.String("k", "", "Specify keyword for filtering")

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	ctgCmd := flag.NewFlagSet("category", flag.ExitOnError)

	searchCmd := flag.NewFlagSet("search", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected at least 1 subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		HandleAdd(addCmd, addSubscription)
	case "show":
		HandleShow(showCmd, showCurrent, showMonth, showYear, showIncome, showExpense, showKeyword)
	case "help":
		HandleHelp(helpCmd)
	case "category":
		HandleCategory(ctgCmd)
	case "search":
		HandleSearch(searchCmd)
	default:
		fmt.Print(models.INSTRUCTION)
	}
}
