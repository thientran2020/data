package command

import (
	"fmt"

	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/repository"
	"github.com/thientran2020/financial-cli/internal/ui"
)

type SearchHandler struct {
	cfg   *config.Config
	repos *repository.Repositories
	ui    *ui.UI
}

func NewSearchHandler(cfg *config.Config, repos *repository.Repositories, ui *ui.UI) *SearchHandler {
	return &SearchHandler{
		cfg:   cfg,
		repos: repos,
		ui:    ui,
	}
}

func (h *SearchHandler) Handle(keyword string) error {
	if keyword == "" {
		return fmt.Errorf("search keyword cannot be empty")
	}

	records, err := h.repos.Record.Search(keyword)
	if err != nil {
		return err
	}

	title := fmt.Sprintf("SEARCH RESULTS FOR '%s'", keyword)
	h.ui.Table.DisplayRecords(records, title)

	if len(records) == 0 {
		h.ui.Prompt.DisplayMessage("No records found for the given keyword.")
	} else {
		h.ui.Prompt.DisplayMessage(fmt.Sprintf("Found %d records matching the keyword '%s'.",
			len(records), keyword))
	}
	return nil
}
