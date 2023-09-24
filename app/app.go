package app

import (
	"github.com/alecthomas/kong"
)

type CLI struct {
	Add struct {
		Income       bool   `optional:"" short:"i" help:"Add income"`
		Subscription bool   `optional:"" short:"s" help:"Add subscription/membership data"`
		Trip         bool   `optional:"" short:"t" help:"Add new trip"`
		Today        bool   `optional:"" short:"n" help:"Set date is today"`
		Yes          bool   `optional:"" help:"Add data without confirming yes/no"`
		Yesterday    bool   `optional:"" short:"y" help:"Set date is yesterday"`
		Category     int    `optional:"" short:"c" help:"Set category by code. To get list of categories, please run 'data get -c'"`
		Cost         int64  `optional:"" help:"Expense cost"`
		Description  string `optional:"" short:"d" help:"Set a short description for the expense/income"`
	} `cmd:"" help:"Add financial data (expense or income)"`
	Get struct {
		Category     bool `optional:"" short:"c" help:"Display category mapping table"`
		Subscription bool `optional:"" short:"s" help:"Display current subscriptions' details"`
		Trip         bool `optional:"" short:"t" help:"List all trips in interactive shell for details"`
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
