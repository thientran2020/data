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
