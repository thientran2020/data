package models

type Record struct {
	Year         int16
	Month        int16
	Day          int16
	Description  string
	Code         int16
	Cost         int32
	Category     string
	Subscription *Subscription
}

type Subscription struct {
	Name         string `default:"none"`
	BillingCycle string `default:"none"`
	Details      string `default:""`
}
