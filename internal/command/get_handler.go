package command

import (
	"fmt"
	"sort"
	"time"

	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/repository"
	"github.com/thientran2020/financial-cli/internal/ui"
)

type GetHandler struct {
	cfg   *config.Config
	repos *repository.Repositories
	ui    *ui.UI
}

func NewGetHandler(cfg *config.Config, repos *repository.Repositories, ui *ui.UI) *GetHandler {
	return &GetHandler{
		cfg:   cfg,
		repos: repos,
		ui:    ui,
	}
}

func (h *GetHandler) Handle(cmd *struct {
	Category     bool `optional:"" short:"c" help:"Get category table"`
	Subscription bool `optional:"" short:"s" help:"Get subscription details"`
	Trip         bool `optional:"" short:"t" help:"Get trip details"`
}) error {
	if cmd.Trip {
		return h.getTripDetails()
	} else if cmd.Subscription {
		return h.getSubscriptionDetails()
	} else if cmd.Category {
		return h.getCategoryTable()
	} else {
		h.ui.DisplayHelp()
		return nil
	}
}

func (h *GetHandler) getCategoryTable() error {
	h.ui.Table.DisplayCategoryTable(h.cfg.Categories)
	return nil
}

func (h *GetHandler) getSubscriptionDetails() error {
	subscriptions, err := h.repos.Subscription.GetAll()
	if err != nil {
		return err
	}

	h.ui.Table.DisplaySubscriptions(subscriptions)
	return nil
}

func (h *GetHandler) getTripDetails() error {
	trips, err := h.repos.Trip.GetAll()
	if err != nil {
		return err
	}

	if len(trips) == 0 {
		h.ui.Prompt.DisplayMessage("No trips found. Add a trip first with 'data add -t'.")
		return nil
	}

	// Sort trips by date in descending order (most recent first)
	sort.Slice(trips, func(i, j int) bool {
		dateI, errI := time.Parse("01-02-2006", trips[i].StartDate)
		dateJ, errJ := time.Parse("01-02-2006", trips[j].StartDate)
		if errI != nil || errJ != nil {
			return false
		}
		return dateI.After(dateJ)
	})

	// Take only the 10 most recent trips
	if len(trips) > 10 {
		trips = trips[:10]
	}

	tripOptions := make([]string, len(trips))
	tripMap := make(map[string]int)

	for i, t := range trips {
		tripOptions[i] = fmt.Sprintf("%-30s %s", t.Name, t.DateRange())
		tripMap[tripOptions[i]] = t.ID
	}

	selectedOption, err := h.ui.Prompt.AskForSelection("Select a trip to view details:", tripOptions)
	if err != nil {
		return err
	}

	tripID := tripMap[selectedOption]

	trip, err := h.repos.Trip.GetByID(tripID)
	if err != nil {
		return err
	}

	h.ui.Table.DisplayTrip(trip)
	return nil
}
