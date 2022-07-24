package app

import (
	"flag"
	"fmt"
	"os"
)

func AppInit() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addIncome := addCmd.Bool("income", false, "Flag to add income data - False by default")

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	showMonth := showCmd.Int("m", -1, "Flag | Specify month you want to retrive financial data | Ex: Jan, Feb, Mar,....")
	showYear := showCmd.Int("y", -1, "Flag | Specify year you want to retrive financial data")

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected at least 1 subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		HandleAdd(addCmd, addIncome)
	case "show":
		HandleShow(showCmd, showMonth, showYear)
	case "help":
		HandleHelp(helpCmd)
	default:
		fmt.Println("Invalid command...!")
		// TODO
		// Print general instruction for financial cli tool
	}
}
