package command

import (
	"github.com/alecthomas/kong"
	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/repository"
	"github.com/thientran2020/financial-cli/internal/ui"
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
		Cost         int    `optional:"" help:"Expense cost"`
		Description  string `optional:"" short:"d" help:"Set a short description for the expense/income"`
	} `cmd:"" help:"Add financial data (expense or income)"`

	Get struct {
		Category     bool `optional:"" short:"c" help:"Get category table"`
		Subscription bool `optional:"" short:"s" help:"Get subscription details"`
		Trip         bool `optional:"" short:"t" help:"Get trip details"`
	} `cmd:"" help:"Get category mapping table or subscriptions' details"`

	Show struct {
		Current bool   `optional:"" short:"c" help:"Show current month's records"`
		Month   int    `optional:"" short:"m" help:"Month to show records for"`
		Year    int    `optional:"" short:"y" help:"Year to show records for"`
		Income  bool   `optional:"" short:"i" help:"Show income records"`
		Expense bool   `optional:"" short:"e" help:"Show expense records"`
		Keyword string `optional:"" short:"k" help:"Filter records by keyword"`
	} `cmd:"" help:"Display financial data in table format - current date by default "`

	Search struct {
		Keyword string `arg:"" required:"" help:"Keyword to search. For example: restaurant, travel,..."`
	} `cmd:"" help:"Display financial data by specific keyword"`
}

type Handler struct {
	cfg           *config.Config
	repos         *repository.Repositories
	ui            *ui.UI
	addHandler    *AddHandler
	getHandler    *GetHandler
	showHandler   *ShowHandler
	searchHandler *SearchHandler
}

func NewHandler(cfg *config.Config, repos *repository.Repositories) *Handler {
	uiManager := ui.NewUI()

	return &Handler{
		cfg:           cfg,
		repos:         repos,
		ui:            uiManager,
		addHandler:    NewAddHandler(cfg, repos, uiManager),
		getHandler:    NewGetHandler(cfg, repos, uiManager),
		showHandler:   NewShowHandler(cfg, repos, uiManager),
		searchHandler: NewSearchHandler(cfg, repos, uiManager),
	}
}

func (h *Handler) Execute() error {
	var cli CLI
	ctx := kong.Parse(&cli)

	switch ctx.Command() {
	case "add":
		return h.addHandler.Handle(&cli.Add)
	case "get":
		return h.getHandler.Handle(&cli.Get)
	case "show":
		return h.showHandler.Handle(&cli.Show)
	case "search <keyword>":
		return h.searchHandler.Handle(cli.Search.Keyword)
	default:
		return ctx.PrintUsage(true)
	}
}
