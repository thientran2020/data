package utils

import (
	"fmt"
	"strings"

	m "github.com/thientran2020/financial-cli/models"
)

func AddSubscription() {
	subscriptionList := ReadSubscriptionJson(GetUserHomeDirectory() + m.BASE_FILEPATH_SUBCRIPTION)

	// Prompt user to enter neccessary information
	name, _ := PromptEnter("What is your new subscription/membership", false)
	ftype, _ := InteractiveSelect(
		"What type of your subscription",
		[]string{"income", "expense"},
	)
	billingCycle, _ := InteractiveSelect(
		"Choose your billing cycle",
		[]string{"monthly", "yearly"},
	)
	cost, _ := NumberEnter("How much per billing period")
	startDate, _ := PromptEnter("What was the start date (mm-dd-yyyy)", false)
	if !IsValidDate(startDate) {
		fmt.Println("Not a valid date. Please rerun and enter required format mm-dd-yyyy!")
		return
	}

	// Create new subscription and add to existing list
	subscription := m.Subscription{
		Name:         name,
		Type:         ftype,
		Cost:         int(cost),
		BillingCycle: billingCycle,
		StartDate:    startDate,
	}

	switch billingCycle {
	case "monthly":
		subscriptionList.Monthly = append(subscriptionList.Monthly, subscription)
	case "yearly":
		subscriptionList.Yearly = append(subscriptionList.Yearly, subscription)
	}

	// Print new subscription and ask for confirmation before adding
	message := fmt.Sprintf("%s: $%d/%s", name, cost, strings.ToLower(billingCycle[:len(billingCycle)-2]))
	PrintCustomizedMessage(message, Green, true)
	confirmed := ConfirmYesNoPromt("Do you confirm to enter above subscription")
	if confirmed {
		WriteSubscriptionJson(GetUserHomeDirectory()+m.BASE_FILEPATH_SUBCRIPTION, subscriptionList)
		PrintCustomizedMessage("Successfully added at "+GetUserHomeDirectory()+m.BASE_FILEPATH_SUBCRIPTION, Yellow, true)
	} else {
		PrintCustomizedMessage("Subscription ignored "+CheckMark, Red, true)
	}
}

func PrintSubcriptionList(billingCycle string, subcriptions []m.Subscription) {
	if subcriptions != nil {
		title := fmt.Sprintf("\n%s SUBSCRIPTION & MEMBERSHIP\n", strings.ToUpper(billingCycle))
		fmt.Printf(Colorize(title, Yellow))
		for _, sub := range subcriptions {
			fmt.Printf("%-28s [%-7s] $%-7d || Start date: %s\n", sub.Name, sub.Type, sub.Cost, sub.StartDate)
		}
	}
}

func GetSubscription() m.MySubscriptionList {
	return ReadSubscriptionJson(GetUserHomeDirectory() + m.BASE_FILEPATH_SUBCRIPTION)
}

func UpdateSubscriptionRecord() {
	data, _ := CsvRead(GetSharedFile())
	subscriptions := GetSubscription()
	UpdateSubRecordByBCycle(data, subscriptions.Monthly, "monthly")
	UpdateSubRecordByBCycle(data, subscriptions.Yearly, "yearly")
}

func UpdateSubRecordByBCycle(data Data, subscriptions []m.Subscription, billingCycle string) {
	yearMap := map[int]bool{}
	for _, s := range subscriptions {
		dateMap := FilterSubscriptionByName(data, s.Name)
		generatedDateFromStartDate := GenerateDateFromStartDate(s.StartDate, billingCycle)
		for _, date := range generatedDateFromStartDate {
			if _, ok := dateMap[date]; !ok {
				month, day, year := GetDateNumber(date)
				var category string
				switch s.Name {
				case "Block":
					category = "Income"
				case "Monthly rent":
					category = "Rent"
				default:
					category = "Subscription"
				}
				record := m.Record{
					Year:        year,
					Month:       month,
					Day:         day,
					Description: s.Name,
					Cost:        s.Cost,
					Category:    category,
					Code:        6,
				}
				yearMap[year] = true
				AddRecordToFile(GetSpecificYearFile(year), record, "")
				AddRecordToFile(GetSharedFile(), record, "")
			}
		}
	}

	if len(yearMap) == 0 {
		return
	}

	for year := range yearMap {
		CsvUpdate(GetSpecificYearFile(year))
		CsvUpdate(GetSharedFile())
	}
}
