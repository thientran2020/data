package app

import (
	"fmt"

	"github.com/alecthomas/kong"
	m "github.com/thientran2020/financial-cli/models"
)

type CLI struct {
	Add struct {
		Subscription bool `optional:"" short:"s" help:"Add subscription/membership data"`
	} `cmd:"" help:"Add financial data (expense or income)"`
	Get struct {
		Category     bool `optional:"" short:"c" help:"Display category mapping table"`
		Subscription bool `optional:"" short:"s" help:"Display current subscriptions' details"`
	} `cmd:"" help:"Get category mapping table or subscriptions' details"`
	Show struct {
		Current bool   `optional:"" short:"c" help:"Retrieve current month data"`
		Month   int    `optional:"" short:"m" help:"Retrive financial data for specific month"`
		Year    int    `optional:"" short:"y" help:"Retrive financial data for specific year"`
		Income  bool   `optional:"" short:"i" help:"Retrive income data"`
		Expense bool   `optional:"" short:"e" help:"Retrive expense data"`
		Keyword string `optional:"" short:"k" help:"Retrive financial data by filtering specific keyword"`
	} `cmd:"" help:"Display financial data in table format - current date by default "`
	Search struct {
		Keyword string `arg:"" required:"" help:"Keyword to search. For example: restaurant, travel,..."`
	} `cmd:"" help:"Display financial data by specific keyword"`
}

func (cmd *CLI) Run(ctx *kong.Context) error {
	switch ctx.Command() {
	case "add":
		HandleAdd(cmd)
	case "get":
		if len(ctx.Args) == 1 {
			fmt.Print(m.INSTRUCTION)
		}
		HandleGet(cmd)
	case "show":
		HandleShow(cmd)
	case "search <keyword>":
		HandleSearch(ctx.Args[1])
	default:
		panic(ctx.Command())
	}
	return nil
}

func AppInit() {
	var cli CLI
	ctx := kong.Parse(&cli)
	err := ctx.Run()
	ctx.FatalIfErrorf(err)
}
