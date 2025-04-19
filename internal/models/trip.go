package models

import (
	"time"
)

type Trip struct {
	ID            int          `json:"id"`
	Name          string       `json:"name"`
	NParticipants int          `json:"nparticipants"`
	Costs         Costs        `json:"costs"`
	Records       []TripRecord `json:"records"`
	StartDate     string       `json:"startDate"`
	EndDate       string       `json:"endDate"`
}

type Costs struct {
	Shared int `json:"shared"`
	Total  int `json:"total"`
}

type TripRecord struct {
	Record *Record `json:"record"`
	Shared bool    `json:"shared"`
	// Direct record fields for backward compatibility
	Year        int    `json:"Year,omitempty"`
	Month       int    `json:"Month,omitempty"`
	Day         int    `json:"Day,omitempty"`
	Description string `json:"Description,omitempty"`
	Cost        int    `json:"Cost,omitempty"`
	Category    string `json:"Category,omitempty"`
	Code        int    `json:"Code,omitempty"`
}

func NewTrip(id int, name string, participants int, startDate, endDate string) *Trip {
	return &Trip{
		ID:            id,
		Name:          name,
		NParticipants: participants,
		Costs: Costs{
			Shared: 0,
			Total:  0,
		},
		Records:   []TripRecord{},
		StartDate: startDate,
		EndDate:   endDate,
	}
}

func (t *Trip) AddRecord(rec *Record, shared bool) {
	tripRecord := TripRecord{
		Record: rec,
		Shared: shared,
	}

	t.Records = append(t.Records, tripRecord)

	t.Costs.Total += rec.Cost
	if shared {
		t.Costs.Shared += rec.Cost
	}
}

func (t *Trip) DateRange() string {
	return t.StartDate + " to " + t.EndDate
}

func (t *Trip) CalculatePerPersonCost() int {
	if t.NParticipants <= 0 {
		return 0
	}

	return t.Costs.Shared / t.NParticipants
}

func (t *Trip) IsActive() (bool, error) {
	startDate, err := time.Parse("01-02-2006", t.StartDate)
	if err != nil {
		return false, err
	}

	endDate, err := time.Parse("01-02-2006", t.EndDate)
	if err != nil {
		return false, err
	}

	now := time.Now()
	return (now.After(startDate) || now.Equal(startDate)) && (now.Before(endDate) || now.Equal(endDate)), nil
}

func (t *Trip) FilterSharedRecords() []TripRecord {
	sharedRecords := []TripRecord{}

	for _, rec := range t.Records {
		if rec.Shared {
			sharedRecords = append(sharedRecords, rec)
		}
	}

	return sharedRecords
}

func (t *Trip) Duration() (int, error) {
	startDate, err := time.Parse("01-02-2006", t.StartDate)
	if err != nil {
		return 0, err
	}

	endDate, err := time.Parse("01-02-2006", t.EndDate)
	if err != nil {
		return 0, err
	}

	duration := endDate.Sub(startDate)
	return int(duration.Hours()/24) + 1, nil
}
