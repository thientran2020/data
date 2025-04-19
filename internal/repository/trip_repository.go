package repository

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/models"
)

type JsonTripRepository struct {
	cfg      *config.Config
	filePath string
}

func NewTripRepository(cfg *config.Config) (*JsonTripRepository, error) {
	repo := &JsonTripRepository{
		cfg:      cfg,
		filePath: cfg.TripFile,
	}

	if err := ensureJSONFile(repo.filePath, []models.Trip{}); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *JsonTripRepository) GetAll() ([]*models.Trip, error) {
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var trips []*models.Trip
	if err := json.Unmarshal(file, &trips); err != nil {
		return nil, err
	}

	return trips, nil
}

func (r *JsonTripRepository) GetByID(id int) (*models.Trip, error) {
	trips, err := r.GetAll()
	if err != nil {
		return nil, err
	}

	for _, t := range trips {
		if t.ID == id {
			return t, nil
		}
	}

	return nil, fmt.Errorf("trip with ID %d not found", id)
}

func (r *JsonTripRepository) Add(t *models.Trip) error {
	trips, err := r.GetAll()
	if err != nil {
		return err
	}

	for _, existingTrip := range trips {
		if existingTrip.ID == t.ID {
			return fmt.Errorf("trip with ID %d already exists", t.ID)
		}
	}

	trips = append(trips, t)
	return r.writeTrips(trips)
}

func (r *JsonTripRepository) AddRecord(tripID int, rec *models.Record, shared bool) error {
	trips, err := r.GetAll()
	if err != nil {
		return err
	}

	var found bool
	for i, t := range trips {
		if t.ID == tripID {
			t.AddRecord(rec, shared)
			trips[i] = t
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("trip with ID %d not found", tripID)
	}

	return r.writeTrips(trips)
}

func (r *JsonTripRepository) writeTrips(trips []*models.Trip) error {
	data, err := json.MarshalIndent(trips, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}
