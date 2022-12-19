package models

type Record struct {
	Year        int
	Month       int
	Day         int
	Description string
	Cost        int
	Category    string
	Code        int
}

type Subscription struct {
	Name         string `json:"name"`
	Type         string `json:"type"`
	Cost         int    `json:"cost"`
	BillingCycle string `json:"billingCycle"`
	StartDate    string `json:"startDate"`
}

type MySubscriptionList struct {
	Monthly []Subscription `json:"monthly"`
	Yearly  []Subscription `json:"yearly"`
}

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
	Record
	Shared bool `json:"shared"`
}
