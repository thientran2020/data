package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/thientran2020/financial-cli/internal/config"
	"github.com/thientran2020/financial-cli/internal/models"
)

type JsonSubscriptionRepository struct {
	cfg        *config.Config
	filePath   string
	recordRepo RecordRepository
}

func NewSubscriptionRepository(cfg *config.Config, recordRepo RecordRepository) (*JsonSubscriptionRepository, error) {
	repo := &JsonSubscriptionRepository{
		cfg:        cfg,
		filePath:   cfg.SubscriptionFile,
		recordRepo: recordRepo,
	}

	if err := ensureJSONFile(repo.filePath, models.NewSubscriptionList()); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *JsonSubscriptionRepository) GetAll() (*models.SubscriptionList, error) {
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, err
	}

	var subscriptions models.SubscriptionList
	if err := json.Unmarshal(file, &subscriptions); err != nil {
		return nil, err
	}

	return &subscriptions, nil
}

func (r *JsonSubscriptionRepository) Add(sub *models.Subscription) error {
	subscriptions, err := r.GetAll()
	if err != nil {
		return err
	}

	sub.BillingCycle = models.BillingCycle(strings.ToLower(string(sub.BillingCycle)))
	switch sub.BillingCycle {
	case models.Monthly:
		subscriptions.Monthly = append(subscriptions.Monthly, *sub)
	case models.Yearly:
		subscriptions.Yearly = append(subscriptions.Yearly, *sub)
	default:
		return fmt.Errorf("unknown billing cycle: %s", sub.BillingCycle)
	}

	return r.writeSubscriptions(subscriptions)
}

func (r *JsonSubscriptionRepository) Update() error {
	subscriptions, err := r.GetAll()
	if err != nil {
		return err
	}

	now := time.Now()
	for _, sub := range subscriptions.Monthly {
		if err := r.processSubscription(&sub, now); err != nil {
			return err
		}
	}

	for _, sub := range subscriptions.Yearly {
		if err := r.processSubscription(&sub, now); err != nil {
			return err
		}
	}

	return nil
}

func (r *JsonSubscriptionRepository) processSubscription(sub *models.Subscription, now time.Time) error {
	nextBillingDate, err := sub.NextBillingDate(now)
	if err != nil {
		return err
	}

	if !nextBillingDate.After(now) {
		record := sub.ToRecord(nextBillingDate)
		if err := r.recordRepo.Add(record); err != nil {
			return err
		}

		sub.StartDate = nextBillingDate.Format("01-02-2006")
		subscriptions, err := r.GetAll()
		if err != nil {
			return err
		}

		updated := false
		switch sub.BillingCycle {
		case models.Monthly:
			for i, s := range subscriptions.Monthly {
				if s.Name == sub.Name {
					subscriptions.Monthly[i] = *sub
					updated = true
					break
				}
			}
		case models.Yearly:
			for i, s := range subscriptions.Yearly {
				if s.Name == sub.Name {
					subscriptions.Yearly[i] = *sub
					updated = true
					break
				}
			}
		}

		if !updated {
			return fmt.Errorf("subscription not found: %s", sub.Name)
		}

		return r.writeSubscriptions(subscriptions)
	}

	return nil
}

func (r *JsonSubscriptionRepository) writeSubscriptions(subscriptions *models.SubscriptionList) error {
	data, err := json.MarshalIndent(subscriptions, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, data, 0644)
}

func ensureJSONFile(filepath string, defaultContent interface{}) error {
	_, err := os.Stat(filepath)
	if err == nil {
		return nil
	}

	if !os.IsNotExist(err) {
		return err
	}

	dir := filepath[:len(filepath)-len(fmt.Sprintf("/%s", getFileName(filepath)))]
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(defaultContent, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath, data, 0644)
}

func getFileName(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' {
			return path[i+1:]
		}
	}
	return path
}
