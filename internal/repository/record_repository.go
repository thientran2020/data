package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/models"
)

var CSVHeader = []string{"Year", "Month", "Day", "Content", "Cost", "Category", "Code"}

type CsvRecordRepository struct {
	cfg         *config.Config
	sharedFile  string
	yearlyFiles map[int]string
}

func NewRecordRepository(cfg *config.Config) (*CsvRecordRepository, error) {
	repo := &CsvRecordRepository{
		cfg:         cfg,
		sharedFile:  cfg.SharedFinanceFile,
		yearlyFiles: make(map[int]string),
	}

	if err := ensureCSVFile(repo.sharedFile); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *CsvRecordRepository) GetAll() ([]*models.Record, error) {
	return r.readCSVFile(r.sharedFile)
}

func (r *CsvRecordRepository) GetByYear(year int) ([]*models.Record, error) {
	if year < r.cfg.StartYear {
		return nil, fmt.Errorf("year %d is below the start year %d", year, r.cfg.StartYear)
	}

	yearlyFile, err := r.getYearlyFilePath(year)
	if err != nil {
		return nil, err
	}

	return r.readCSVFile(yearlyFile)
}

func (r *CsvRecordRepository) Add(rec *models.Record) error {
	if err := r.appendToCSV(r.sharedFile, rec); err != nil {
		return err
	}

	yearlyFile, err := r.getYearlyFilePath(rec.Year)
	if err != nil {
		return err
	}

	return r.appendToCSV(yearlyFile, rec)
}

func (r *CsvRecordRepository) Search(keyword string) ([]*models.Record, error) {
	records, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	keyword = strings.ToLower(keyword)
	var results []*models.Record

	for _, rec := range records {
		if strings.Contains(strings.ToLower(rec.Description), keyword) ||
			strings.Contains(strings.ToLower(rec.Category), keyword) {
			results = append(results, rec)
		}
	}

	return results, nil
}

// Filter filters records by month, type, and/or keyword
func (r *CsvRecordRepository) Filter(month int, recordType string, keyword string) ([]*models.Record, error) {
	records, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	var filtered []*models.Record
	for _, rec := range records {
		if month > 0 && rec.Month != month {
			continue
		}

		if recordType == "income" && !rec.IsIncome() {
			continue
		} else if recordType == "expense" && !rec.IsExpense() {
			continue
		}

		if keyword != "" {
			keywordLower := strings.ToLower(keyword)
			if !strings.Contains(strings.ToLower(rec.Description), keywordLower) &&
				!strings.Contains(strings.ToLower(rec.Category), keywordLower) {
				continue
			}
		}

		filtered = append(filtered, rec)
	}

	return filtered, nil
}

func (r *CsvRecordRepository) getYearlyFilePath(year int) (string, error) {
	// Check if we already have the path cached
	if path, ok := r.yearlyFiles[year]; ok {
		return path, nil
	}

	path := r.cfg.GetYearlyFinanceFile(year)
	if err := ensureCSVFile(path); err != nil {
		return "", err
	}

	r.yearlyFiles[year] = path
	return path, nil
}

func (r *CsvRecordRepository) readCSVFile(filepath string) ([]*models.Record, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var records []*models.Record

	for i, row := range rows {
		if i == 0 {
			continue
		}

		rec, err := models.RecordFromCSVRow(row)
		if err != nil {
			return nil, fmt.Errorf("error parsing row %d: %w", i, err)
		}

		records = append(records, rec)
	}

	return records, nil
}

func (r *CsvRecordRepository) appendToCSV(filepath string, rec *models.Record) error {
	file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(rec.ToCSVRow())
}

func ensureCSVFile(filepath string) error {
	_, err := os.Stat(filepath)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return writer.Write(CSVHeader)
}
