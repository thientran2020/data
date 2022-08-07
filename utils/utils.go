package utils

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
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
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}
	_, err := prompt.Run()
	return err == nil
}

func InteractiveSelect(label string, items []string) (string, error) {
	prompt := promptui.Select{
		Label:  label,
		Items:  items,
		Size:   len(items),
		Stdout: &BellSkipper{},
	}
	_, result, err := prompt.Run()

	if err != nil {
		return "", err
	}
	return result, nil
}

func NumberEnter(label string) (int64, error) {
	validate := func(input string) error {
		number, err := strconv.ParseInt(input, 10, 64)
		if err != nil {
			return errors.New("Invalid number")
		}
		if number < 0 {
			return errors.New("Negative number")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	stringNum, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	result, _ := strconv.ParseInt(stringNum, 10, 64)
	return result, nil
}

func PromptEnter(label string, empty bool) (string, error) {
	validate := func(input string) error {
		if len(input) == 0 && !empty {
			return errors.New("Invalid input")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    label,
		Validate: validate,
	}

	result, err := prompt.Run()
	return result, err
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

func IsValidDate(dateString string) bool {
	re := regexp.MustCompile("(0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])-((19|20)\\d\\d)")
	return re.MatchString(dateString)
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
		skipRowByMonth := month != -1 && row_month != month

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

// Resolve terminal's bell ring issue when moving between interactive select
// The following implementation followed from: https://github.com/manifoldco/promptui/issues/49
type BellSkipper struct{}

func (bs *BellSkipper) Write(b []byte) (int, error) {
	const charBell = 7
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

func (bs *BellSkipper) Close() error {
	return os.Stderr.Close()
}
