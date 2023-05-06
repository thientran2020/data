package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
)

const (
	ColorOff = "\033[0m"

	// Bold Colors
	Red    = "\033[1;31m"
	Green  = "\033[1;32m"
	Yellow = "\033[1;33m"
	White  = "\033[1;37m"

	// Underline Colors
	URed    = "\033[4;31m"
	UGreen  = "\033[4;32m"
	UYellow = "\033[4;33m"
	UWhite  = "\033[4;37m"

	UnderlineGreenCommandColor = UGreen + "%s" + ColorOff
	CheckMark                  = "\u2713"
)

// Different Input Types
func ConfirmYesNoPromt(label string) bool {
	confirmed := false
	prompt := &survey.Confirm{
		Message: label,
	}
	err := survey.AskOne(prompt, &confirmed)
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	return confirmed
}

func InteractiveSelect(label string, items []string) string {
	selected := ""
	prompt := &survey.Select{
		Message:  label,
		Options:  items,
		PageSize: 14,
	}
	err := survey.AskOne(prompt, &selected)
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	return selected
}

func NumberEnter(label string) int64 {
	q := &survey.Question{
		Prompt: &survey.Input{Message: label},
		Validate: func(val interface{}) error {
			strNumber, ok := val.(string)
			number, err := strconv.ParseInt(strNumber, 10, 64)
			if !ok || len(strNumber) == 0 || err != nil {
				return errors.New("Please enter a valid number!")
			}
			if number <= 0 {
				return errors.New("Please enter a positive number!")
			}
			return nil
		},
	}
	answer := ""
	survey.Ask([]*survey.Question{q}, &answer)
	numAnswer, _ := strconv.ParseInt(answer, 10, 64)
	return numAnswer
}

func PromptEnter(label string) string {
	answer := ""
	prompt := &survey.Input{Message: label}
	err := survey.AskOne(prompt, &answer, survey.WithValidator(survey.Required))
	if err == terminal.InterruptErr {
		os.Exit(0)
	}
	return answer
}

// Return today's date by default if no input entered
func DateEnter(label string) string {
	q := &survey.Question{
		Prompt: &survey.Input{Message: label},
		Validate: func(val interface{}) error {
			input, ok := val.(string)
			if !ok {
				return errors.New("Valid input required!")
			}
			if len(input) != 0 && !IsValidDate(input) {
				return errors.New("Invalid date. Please enter with format mm-dd-yyyy...!")
			}
			return nil
		},
	}
	date := ""
	survey.Ask([]*survey.Question{q}, &date)
	if date == "" {
		date = time.Now().Format("01-02-2006")
	}
	return date
}

// Print customized messages with color
func PrintCustomizedMessage(message string, color string, newline bool) {
	message = strings.ReplaceAll(message, ColorOff, "")
	if newline {
		fmt.Println(Colorize(message, color))
	} else {
		fmt.Print(Colorize(message, color))
	}
}

// Helper functions
func GetStringDateFromNumber(number int) string {
	if number < 10 {
		return "0" + strconv.Itoa(number)
	}
	return strconv.Itoa(number)
}

func GetStringDateFromString(number string) string {
	if len(number) < 2 {
		return "0" + number
	}
	return number
}

func ContainString(s, ss string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(ss))
}

func Colorize(text string, color string) string {
	return fmt.Sprintf("%s%s%s", color, text, ColorOff)
}

// date format: mm-dd-yyyy
func IsValidDate(dateString string) bool {
	_, err := time.Parse("01-02-2006", dateString)
	return err == nil
}

func GenerateDateFromStartDate(startDate, billingCycle string) []string {
	month, day, year := GetDateNumber(startDate)
	generatedDate := []string{}

	if billingCycle == "monthly" {
		for y := year; y <= time.Now().Year(); y++ {
			lower, upper := 1, 12
			if y == year {
				lower = month
			}
			if y == time.Now().Year() {
				upper = int(time.Now().Month())
			}
			for m := lower; m <= upper; m++ {
				newDate := fmt.Sprintf("%s-%s-%s",
					GetStringDateFromNumber(m),
					GetStringDateFromNumber(day),
					GetStringDateFromNumber(y),
				)
				generatedDate = append(generatedDate, newDate)
			}
		}
	} else if billingCycle == "yearly" {
		for y := year; y <= time.Now().Year(); y++ {
			newDate := fmt.Sprintf("%s-%s-%s",
				GetStringDateFromNumber(month),
				GetStringDateFromNumber(day),
				GetStringDateFromNumber(y),
			)
			generatedDate = append(generatedDate, newDate)
		}
	}
	return generatedDate
}

// Return date string with format "mm-dd-yyyy" to month, day, year numbers
func GetDateNumber(dateString string) (int, int, int) {
	date := strings.Split(dateString, "-")
	m, _ := strconv.Atoi(date[0])
	d, _ := strconv.Atoi(date[1])
	y, _ := strconv.Atoi(date[2])
	return m, d, y
}

// Filter []Record data based on month, type (income/expense) & keyword
func FilterData(data Data, month int, typeFlag, keyword string) Data {
	filteredData := Data{}
	count := 0
	for _, row := range data {
		dateArray := strings.Split(row[1].(string), "-")
		row_month, _ := strconv.Atoi(dateArray[0])
		skipRowByMonth := month != 0 && row_month != month

		skipRowByTypeFlag :=
			(typeFlag == "income" && row[3] != "Income") ||
				(typeFlag == "expense" && row[3] == "Income")

		skipRowByKeyWord := !ContainString(row[2].(string), keyword) && !ContainString(row[3].(string), keyword)

		skip := skipRowByMonth || skipRowByTypeFlag || skipRowByKeyWord
		if !skip {
			count++
			formatted_date := fmt.Sprintf("%s-%s-%s",
				Colorize(dateArray[0], Yellow),
				dateArray[1],
				Colorize(dateArray[2], UGreen),
			)
			row[0] = count
			row[1] = formatted_date
			filteredData = append(filteredData, row)
		}
	}
	return filteredData
}

func FilterSubscriptionByName(data Data, subscription string) map[string]bool {
	dateMap := map[string]bool{}
	for _, row := range data {
		if row[2].(string) == subscription &&
			(row[3].(string) == "Subscription" || row[3].(string) == "Rent" || row[3].(string) == "Income") {
			dateMap[row[1].(string)] = true
		}
	}
	return dateMap
}

func CenterString(s string, width int) string {
	return fmt.Sprintf("%[1]*s", -width, fmt.Sprintf("%[1]*s", (width+len(s))/2, s))
}
