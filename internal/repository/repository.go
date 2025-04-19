package repository

import (
	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/models"
)

type Repositories struct {
	Record       RecordRepository
	Subscription SubscriptionRepository
	Trip         TripRepository
}

type RecordRepository interface {
	GetAll() ([]*models.Record, error)
	GetByYear(year int) ([]*models.Record, error)
	Add(record *models.Record) error
	Search(keyword string) ([]*models.Record, error)
	Filter(month int, recordType string, keyword string) ([]*models.Record, error)
}

type SubscriptionRepository interface {
	GetAll() (*models.SubscriptionList, error)
	Add(sub *models.Subscription) error
	Update() error
}

type TripRepository interface {
	GetAll() ([]*models.Trip, error)
	GetByID(id int) (*models.Trip, error)
	Add(trip *models.Trip) error
	AddRecord(tripID int, record *models.Record, shared bool) error
}

func NewRepositories(cfg *config.Config) (*Repositories, error) {
	recordRepo, err := NewRecordRepository(cfg)
	if err != nil {
		return nil, err
	}

	subscriptionRepo, err := NewSubscriptionRepository(cfg, recordRepo)
	if err != nil {
		return nil, err
	}

	tripRepo, err := NewTripRepository(cfg)
	if err != nil {
		return nil, err
	}

	return &Repositories{
		Record:       recordRepo,
		Subscription: subscriptionRepo,
		Trip:         tripRepo,
	}, nil
}
