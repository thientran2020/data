package models

import (
	"fmt"
	"time"
)

type BillingCycle string

const (
	Monthly BillingCycle = "monthly"
	Yearly  BillingCycle = "yearly"
)

type Subscription struct {
	Name         string       `json:"name"`
	Type         DataType     `json:"type"`
	Cost         int          `json:"cost"`
	BillingCycle BillingCycle `json:"billingCycle"`
	StartDate    string       `json:"startDate"`
}

type SubscriptionList struct {
	Monthly []Subscription `json:"monthly"`
	Yearly  []Subscription `json:"yearly"`
}

func NewSubscription(name string, subType DataType, cost int, cycle BillingCycle, startDate string) *Subscription {
	return &Subscription{
		Name:         name,
		Type:         subType,
		Cost:         cost,
		BillingCycle: cycle,
		StartDate:    startDate,
	}
}

func (s *Subscription) FormattedCost() string {
	return fmt.Sprintf("$%d", s.Cost)
}

func (s *Subscription) IsExpense() bool {
	return s.Type == Expense
}

func (s *Subscription) IsIncome() bool {
	return s.Type == Income
}

func (s *Subscription) NextBillingDate(currentDate time.Time) (time.Time, error) {
	startDate, err := time.Parse("01-02-2006", s.StartDate)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid start date: %v", err)
	}

	if currentDate.Before(startDate) {
		return startDate, nil
	}

	switch s.BillingCycle {
	case Monthly:
		monthsPassed := (currentDate.Year()-startDate.Year())*12 + int(currentDate.Month()-startDate.Month())
		nextDate := startDate.AddDate(0, monthsPassed, 0)
		if !nextDate.After(currentDate) {
			nextDate = nextDate.AddDate(0, 1, 0)
		}
		return nextDate, nil
	case Yearly:
		yearsPassed := currentDate.Year() - startDate.Year()
		nextDate := startDate.AddDate(yearsPassed, 0, 0)
		if !nextDate.After(currentDate) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}
		return nextDate, nil
	default:
		return time.Time{}, fmt.Errorf("unknown billing cycle: %s", s.BillingCycle)
	}
}

func (s *Subscription) ToRecord(date time.Time) *Record {
	recordType := Expense
	code := Code

	if s.Type == Income {
		recordType = Income
		code = 0
	}

	return &Record{
		Year:        date.Year(),
		Month:       int(date.Month()),
		Day:         date.Day(),
		Description: s.Name,
		Cost:        s.Cost,
		Category:    "Subscription",
		Code:        code,
		Type:        recordType,
		Date:        date,
	}
}

func NewSubscriptionList() *SubscriptionList {
	return &SubscriptionList{
		Monthly: []Subscription{},
		Yearly:  []Subscription{},
	}
}

func (sl *SubscriptionList) AddSubscription(sub Subscription) {
	switch sub.BillingCycle {
	case Monthly:
		sl.Monthly = append(sl.Monthly, sub)
	case Yearly:
		sl.Yearly = append(sl.Yearly, sub)
	}
}
