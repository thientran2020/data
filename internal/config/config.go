package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	SharedFinanceFile = "finance.csv"
	YearlyFinanceFile = "finance_%d.csv"
	SubscriptionFile  = "subscription.json"
	TripFile          = "trip.json"
	StartYear         = 2017
)

var (
	ROW_HEADER = []string{"Year", "Month", "Day", "Content", "Cost", "Category", "Code"}
	HEADERS    = []string{"#", "DATE", "DESCRIPTION", "CATEGORY", "COST"}
	LABELS     = map[string]map[string]string{
		"Expense": {
			"Description": "What did you spend money for:",
			"Cost":        "Nice. How much did you spend:",
		},
		"Income": {
			"Description": "What did you work for today:",
			"Cost":        "Awesome. How much did you earn:",
		},
	}

	CATEGORIES = []string{
		"Income",
		"Mortgage",
		"Utilities",
		"Insurance",
		"Vehicle Services",
		"Fuel - Car Wash",
		"Subcription",
		"Restaurants",
		"Amazon Shopping",
		"Merchandise",
		"Travel",
		"Personal",
		"Trip",
	}
)

type Config struct {
	Appname string
	Version string

	DataDir           string
	SharedFinanceFile string
	YearlyFinanceFile string
	SubscriptionFile  string
	TripFile          string

	StartYear  int
	Categories []string
	DataFormat string
}

func Default() *Config {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	dataDir := filepath.Join(homeDir, "finance")

	return &Config{
		Appname: "data",
		Version: "2.0.0",

		DataDir:           dataDir,
		SharedFinanceFile: filepath.Join(dataDir, SharedFinanceFile),
		YearlyFinanceFile: filepath.Join(dataDir, YearlyFinanceFile),
		SubscriptionFile:  filepath.Join(dataDir, SubscriptionFile),
		TripFile:          filepath.Join(dataDir, TripFile),
		StartYear:         StartYear,
		DataFormat:        "01-02-2006", // MM-DD-YYYY
		Categories:        CATEGORIES,
	}
}

func Load() (*Config, error) {
	cfg := Default()

	if err := ensureDataDir(cfg.DataDir); err != nil {
		return nil, fmt.Errorf("failed to ensure data directory: %v", err)
	}
	return cfg, nil
}

func ensureDataDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, 0755)
	}
	return nil
}

func (c *Config) GetYearlyFinanceFile(year int) string {
	return fmt.Sprintf(c.YearlyFinanceFile, year)
}

func (c *Config) CategoryCode(category string) int {
	for i, cat := range c.Categories {
		if cat == category {
			return i
		}
	}
	return -1
}

func (c *Config) CategoryByCode(code int) string {
	if code >= 0 && code < len(c.Categories) {
		return c.Categories[code]
	}
	return ""
}
