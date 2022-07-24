package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addIncome := addCmd.Bool("income", false, "Flag to add income data - False by default")

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	showMonth := showCmd.String("m", "", "Flag | Specify month you want to retrive financial data | Ex: Jan, Feb, Mar,....")
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

func HandleAdd(addCmd *flag.FlagSet, income *bool) {
	addCmd.Parse(os.Args[2:])

	if addCmd.NArg() != 0 {
		addCmd.PrintDefaults()
		os.Exit(1)
	}

	if *income == true {
		fmt.Println("Adding income data...")
		// TODO
		// Handle income input
		return
	}
	// TODO
	// Handle expense input
	fmt.Println("Adding expense data...")
}

func HandleShow(showCmd *flag.FlagSet, month *string, year *int) {
	showCmd.Parse(os.Args[2:])

	if showCmd.NArg() != 0 {
		showCmd.PrintDefaults()
		os.Exit(1)
	}

	if *month == "" {
		*month = time.Now().Month().String()
	}
	if *year == -1 {
		*year = time.Now().Year()
	}
	fmt.Printf("Show financial data for %s, %d\n", *month, *year)
	// TODO
	// Retrieve financial data for specific date
}

func HandleHelp(helpCmd *flag.FlagSet) {
	helpCmd.Parse(os.Args[2:])
	fmt.Println("HELP COMMAND HERE....")
	// TODO
	// Add general instruction here again...
}
