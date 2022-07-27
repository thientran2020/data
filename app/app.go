package app

import (
	"flag"
	"fmt"
	"os"

	"github.com/thientran2020/financial-cli/models"
)

func AppInit() {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	addSubscription := addCmd.Bool("sub", false, "Specify when adding subscription/membership data")

	showCmd := flag.NewFlagSet("show", flag.ExitOnError)
	showMonth := showCmd.Int("m", -1, "Flag | Specify month you want to retrive financial data")
	showYear := showCmd.Int("y", -1, "Flag | Specify year you want to retrive financial data")

	helpCmd := flag.NewFlagSet("help", flag.ExitOnError)

	ctgCmd := flag.NewFlagSet("category", flag.ExitOnError)

	if len(os.Args) < 2 {
		fmt.Println("Expected at least 1 subcommand")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		HandleAdd(addCmd, addSubscription)
	case "show":
		HandleShow(showCmd, showMonth, showYear)
	case "help":
		HandleHelp(helpCmd)
	case "category":
		HandleCategory(ctgCmd)
	default:
		fmt.Print(models.INSTRUCTION)
	}
}
