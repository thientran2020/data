package app

import (
	"flag"
	"fmt"
	"os"
	"time"
)

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
	fmt.Println("HELP COMMAND HERE....")
	// TODO
	// Add general instruction here again...
}
