package command

import (
	"fmt"
	"time"

	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/models"
	"github.com/thientran2020/financial-cli/internal/repository"
	"github.com/thientran2020/financial-cli/internal/ui"
)

type AddHandler struct {
	cfg   *config.Config
	repos *repository.Repositories
	ui    *ui.UI
}

func NewAddHandler(cfg *config.Config, repos *repository.Repositories, ui *ui.UI) *AddHandler {
	return &AddHandler{
		cfg:   cfg,
		repos: repos,
		ui:    ui,
	}
}

func (h *AddHandler) Handle(cmd *struct {
	Income       bool   `optional:"" short:"i" help:"Add income"`
	Subscription bool   `optional:"" short:"s" help:"Add subscription/membership data"`
	Trip         bool   `optional:"" short:"t" help:"Add new trip"`
	Today        bool   `optional:"" short:"n" help:"Set date is today"`
	Yes          bool   `optional:"" help:"Add data without confirming yes/no"`
	Yesterday    bool   `optional:"" short:"y" help:"Set date is yesterday"`
	Category     int    `optional:"" short:"c" help:"Set category by code. To get list of categories, please run 'data get -c'"`
	Cost         int    `optional:"" help:"Expense cost"`
	Description  string `optional:"" short:"d" help:"Set a short description for the expense/income"`
}) error {
	if cmd.Subscription {
		return h.addSubscription()
	}

	if cmd.Trip {
		return h.addTrip()
	}

	recordType := models.Expense
	if cmd.Income {
		recordType = models.Income
	}

	date := time.Now()
	if cmd.Yesterday {
		date = date.AddDate(0, 0, -1)
	} else if !cmd.Today {
		dateStr, err := h.ui.Prompt.AskForDate("Enter date (MM-DD-YYYY): ")
		if err != nil {
			return err
		}

		parsedDate, err := time.Parse("01-02-2006", dateStr)
		if err != nil {
			return fmt.Errorf("invalid date: %w", err)
		}
		date = parsedDate
	}

	description := cmd.Description
	if description == "" {
		var err error
		var promptMsg string

		if recordType == models.Income {
			promptMsg = "What did you work for today"
		} else {
			promptMsg = "What did you spend money for"
		}

		description, err = h.ui.Prompt.AskForText(promptMsg+": ", "")
		if err != nil {
			return err
		}
	}

	cost := cmd.Cost
	if cost == 0 {
		var err error
		var promptMsg string

		if recordType == models.Income {
			promptMsg = "Awesome. How much did you earn"
		} else {
			promptMsg = "Nice. How much did you spend"
		}

		cost, err = h.ui.Prompt.AskForNumber(promptMsg + ": ")
		if err != nil {
			return err
		}
	}

	code := cmd.Category
	category := h.cfg.CategoryByCode(code)

	if (code == 0) && (recordType == models.Expense) {
		categoryOptions := h.cfg.Categories

		var err error
		category, code, err = h.ui.Prompt.AskForSelectionWithIndex(
			"Pick the category that describe best your entered data:",
			categoryOptions,
		)
		if err != nil {
			return err
		}
	}

	rec := models.NewRecord(
		date.Year(),
		int(date.Month()),
		date.Day(),
		description,
		cost,
		category,
		code,
	)

	if category == "Trip" {
		return h.addTripRecord(rec)
	}

	h.ui.Table.DisplayRecords([]*models.Record{rec}, "New Record")

	if cmd.Yes {
		return h.repos.Record.Add(rec)
	}

	confirm, err := h.ui.Prompt.AskForConfirmation("Do you confirm to enter the above record?")
	if err != nil {
		return err
	}

	if confirm {
		return h.repos.Record.Add(rec)
	}

	h.ui.Prompt.DisplayMessage("Record ignored.")
	return nil
}

func (h *AddHandler) addSubscription() error {
	name, err := h.ui.Prompt.AskForText("What is your new subscription/membership: ", "")
	if err != nil {
		return err
	}

	typeOptions := []string{string(models.Income), string(models.Expense)}
	subType, _, err := h.ui.Prompt.AskForSelectionWithIndex("Select type:", typeOptions)
	if err != nil {
		return err
	}

	cycleOptions := []string{string(models.Monthly), string(models.Yearly)}
	cycle, _, err := h.ui.Prompt.AskForSelectionWithIndex("Select billing cycle:", cycleOptions)
	if err != nil {
		return err
	}

	cost, err := h.ui.Prompt.AskForNumber("How much per billing period: ")
	if err != nil {
		return err
	}

	startDate, err := h.ui.Prompt.AskForDate("What was the start date (MM-DD-YYYY): ")
	if err != nil {
		return err
	}

	sub := models.NewSubscription(
		name,
		models.DataType(subType),
		cost,
		models.BillingCycle(cycle),
		startDate,
	)

	h.ui.Prompt.DisplayMessage(fmt.Sprintf("\nNew %s Subscription: %s, Cost: %s, Start Date: %s\n",
		cycle, name, sub.FormattedCost(), startDate))

	confirm, err := h.ui.Prompt.AskForConfirmation("Do you confirm to add this subscription?")
	if err != nil {
		return err
	}

	if confirm {
		return h.repos.Subscription.Add(sub)
	}

	h.ui.Prompt.DisplayMessage("Subscription ignored.")
	return nil
}

func (h *AddHandler) addTrip() error {
	name, err := h.ui.Prompt.AskForText("Enter trip name: ", "")
	if err != nil {
		return err
	}

	participants, err := h.ui.Prompt.AskForNumber("How many participants: ")
	if err != nil {
		return err
	}

	startDate, err := h.ui.Prompt.AskForDate("Start date (MM-DD-YYYY): ")
	if err != nil {
		return err
	}

	endDate, err := h.ui.Prompt.AskForDate("End date (MM-DD-YYYY): ")
	if err != nil {
		return err
	}

	trips, err := h.repos.Trip.GetAll()
	if err != nil {
		return err
	}

	nextID := 1
	for _, t := range trips {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}

	newTrip := models.NewTrip(nextID, name, participants, startDate, endDate)

	h.ui.Prompt.DisplayMessage(fmt.Sprintf("\nNew Trip: %s, Participants: %d, Date: %s\n",
		name, participants, newTrip.DateRange()))

	confirm, err := h.ui.Prompt.AskForConfirmation("Do you confirm to add this trip?")
	if err != nil {
		return err
	}

	if confirm {
		return h.repos.Trip.Add(newTrip)
	}

	h.ui.Prompt.DisplayMessage("Trip ignored.")
	return nil
}

func (h *AddHandler) addTripRecord(rec *models.Record) error {
	trips, err := h.repos.Trip.GetAll()
	if err != nil {
		return err
	}

	if len(trips) == 0 {
		return fmt.Errorf("no trips found. Please add a trip first with 'data add -t'")
	}

	tripOptions := make([]string, len(trips))
	tripMap := make(map[string]int)

	for i, t := range trips {
		tripOptions[i] = fmt.Sprintf("%s (%s)", t.Name, t.DateRange())
		tripMap[tripOptions[i]] = t.ID
	}

	selectedTrip, err := h.ui.Prompt.AskForSelection("Select a trip for this expense:", tripOptions)
	if err != nil {
		return err
	}

	tripID := tripMap[selectedTrip]

	shared, err := h.ui.Prompt.AskForConfirmation("Is this a shared expense?")
	if err != nil {
		return err
	}

	if err := h.repos.Trip.AddRecord(tripID, rec, shared); err != nil {
		return err
	}

	return h.repos.Record.Add(rec)
}
