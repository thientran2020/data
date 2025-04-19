package command

import (
	"fmt"
	"strings"
	"time"

	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/models"
	"github.com/thientran2020/financial-cli/internal/repository"
	"github.com/thientran2020/financial-cli/internal/ui"
)

type ShowHandler struct {
	cfg   *config.Config
	repos *repository.Repositories
	ui    *ui.UI
}

func NewShowHandler(cfg *config.Config, repos *repository.Repositories, ui *ui.UI) *ShowHandler {
	return &ShowHandler{
		cfg:   cfg,
		repos: repos,
		ui:    ui,
	}
}

func (h *ShowHandler) Handle(cmd *struct {
	Current bool   `optional:"" short:"c" help:"Show current month's records"`
	Month   int    `optional:"" short:"m" help:"Month to show records for"`
	Year    int    `optional:"" short:"y" help:"Year to show records for"`
	Income  bool   `optional:"" short:"i" help:"Show income records"`
	Expense bool   `optional:"" short:"e" help:"Show expense records"`
	Keyword string `optional:"" short:"k" help:"Filter records by keyword"`
}) error {
	if err := h.repos.Subscription.Update(); err != nil {
		return fmt.Errorf("failed to update subscriptions: %w", err)
	}

	year := cmd.Year
	month := cmd.Month

	if cmd.Current {
		now := time.Now()
		year = now.Year()
		month = int(now.Month())
	}

	recordType := "all"
	if cmd.Income && !cmd.Expense {
		recordType = "income"
	} else if cmd.Expense && !cmd.Income {
		recordType = "expense"
	}

	yearRecords, err := h.getRecords(year)
	if err != nil {
		return err
	}

	var filteredRecords []*models.Record
	for _, rec := range yearRecords {
		if month > 0 && rec.Month != month {
			continue
		}

		if recordType == "income" && !rec.IsIncome() {
			continue
		} else if recordType == "expense" && !rec.IsExpense() {
			continue
		}

		if cmd.Keyword != "" {
			keywordLower := strings.ToLower(cmd.Keyword)
			if !strings.Contains(strings.ToLower(rec.Description), keywordLower) &&
				!strings.Contains(strings.ToLower(rec.Category), keywordLower) {
				continue
			}
		}

		filteredRecords = append(filteredRecords, rec)
	}

	title := "FINANCIAL DATA"
	if month > 0 {
		title += fmt.Sprintf(" FOR %d/%d", month, year)
	} else if year > 0 {
		title += fmt.Sprintf(" FOR %d", year)
	}

	if cmd.Keyword != "" {
		title += fmt.Sprintf(" MATCHING '%s'", cmd.Keyword)
	}

	h.ui.Table.DisplayRecordsByType(filteredRecords, recordType, title)
	return nil
}

func (h *ShowHandler) getRecords(year int) ([]*models.Record, error) {
	if year == 0 {
		return h.repos.Record.GetAll()
	}

	if year < h.cfg.StartYear || year > time.Now().Year() {
		return nil, fmt.Errorf("year must be between %d and %d", h.cfg.StartYear, time.Now().Year())
	}

	return h.repos.Record.GetByYear(year)
}
