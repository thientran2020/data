package models

import "time"

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
	Name         string
	StartDate    time.Time
	BillingCycle string
	Note         string `default:"---"`
}

type MySubscriptionList struct {
	NSub int
	LSub []Subscription
}
