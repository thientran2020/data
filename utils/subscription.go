package utils

import (
	"fmt"
	"strings"
	"time"

	m "github.com/thientran2020/financial-cli/models"
)

func AddSubscription() {
	subscriptionList := ReadJson(m.BASE_FILEPATH_SUBCRIPTION)

	// Prompt user to enter neccessary information
	startDate := strings.Split(time.Now().String(), " ")[0]
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
		WriteJson(m.BASE_FILEPATH_SUBCRIPTION, subscriptionList)
		PrintCustomizedMessage("Successfully added at "+m.BASE_FILEPATH_SUBCRIPTION, Yellow, true)
	} else {
		PrintCustomizedMessage("Subscription ignored "+CheckMark, Red, true)
	}
}

func PrintSubcriptionList(billingCycle string, subcriptions []m.Subscription) {
	if subcriptions != nil {
		title := fmt.Sprintf("\n***** %s SUBSCRIPTION & MEMBERSHIP *****\n", strings.ToUpper(billingCycle))
		fmt.Printf(Colorize(title, Yellow))
		for _, sub := range subcriptions {
			fmt.Printf("%-14s [%-7s] $%-7d || Start date: %s\n", sub.Name, sub.Type, sub.Cost, sub.StartDate)
		}
		fmt.Println()
	}
}

func GetSubscription() {
	data := ReadJson(m.BASE_FILEPATH_SUBCRIPTION)
	PrintSubcriptionList("monthly", data.Monthly)
	PrintSubcriptionList("yearly", data.Yearly)
}

// TODOs
func UpdateSubscription() {

}