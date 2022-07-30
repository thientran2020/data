package utils

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"
	"github.com/manifoldco/promptui"
	"github.com/thientran2020/financial-cli/models"
)

const (
	ColorOff = "\033[0m"

	// Regular Colors
	Black  = "\033[0;30m"
	Red    = "\033[0;31m"
	Green  = "\033[0;32m"
	Yellow = "\033[0;33m"
	Blue   = "\033[0;34m"
	Purple = "\033[0;35m"
	Cyan   = "\033[0;36m"
	White  = "\033[0;37m"

	// Bold Colors
	BBlack  = "\033[1;30m"
	BRed    = "\033[1;31m"
	BGreen  = "\033[1;32m"
	BYellow = "\033[1;33m"
	BBlue   = "\033[1;34m"
	BPurple = "\033[1;35m"
	BCyan   = "\033[1;36m"
	BWhite  = "\033[1;37m"

	// Underline Colors
	UBlack  = "\033[4;30m"
	URed    = "\033[4;31m"
	UGreen  = "\033[4;32m"
	UYellow = "\033[4;33m"
	UBlue   = "\033[4;34m"
	UPurple = "\033[4;35m"
	UCyan   = "\033[4;36m"
	UWhite  = "\033[4;37m"

	UnderlineCommandColor = UGreen + "%s" + ColorOff
	BoldCommandColor      = BGreen + "%s" + ColorOff
	RedCommandColor       = Red + "%s" + ColorOff
	CheckMark             = "\u2713"
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

func PrintSingleRecord(record models.Record, color string) {
	date := fmt.Sprintf("    %s/%s/%s   ", GetStringDateFromNumber(record.Month), GetStringDateFromNumber(record.Day), strconv.Itoa(record.Year))
	description := fmt.Sprintf(" %-35s", record.Description)
	costString := fmt.Sprintf(" $%-6s", strconv.Itoa(record.Cost))
	category := fmt.Sprintf(" %-18s", record.Category)

	message := fmt.Sprintf("\n%s\n|%s|%s|%s|%s|\n%s\n", models.DASH, date, description, costString, category, models.DASH)
	PrintCustomizedMessage(message, color, true)
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

// Print customized table with given styles
func PrintTable(data [][]interface{}, headers []string, typeFlag string, style *simpletable.Style) {
	table := simpletable.New()

	// Generate Table Header
	for _, header := range headers {
		headerCell := simpletable.Cell{Align: simpletable.AlignCenter, Text: Colorize(header, BGreen)}
		table.Header.Cells = append(table.Header.Cells, &headerCell)
	}

	// Generate Table Body
	totalIncome := 0
	totalExpense := 0
	n := len(headers)
	for _, rowData := range data {
		row := []*simpletable.Cell{}
		isIncome := string(rowData[n-2].(string)) == "Income"
		for idx, rowCellData := range rowData {
			var rowCell *simpletable.Cell
			if _, ok := rowCellData.(int); ok {
				rowCell = &simpletable.Cell{Align: simpletable.AlignRight, Text: fmt.Sprintf("%d", rowCellData)}
			} else {
				rowCell = &simpletable.Cell{Align: simpletable.AlignLeft, Text: fmt.Sprintf("%s", rowCellData)}
			}
			if idx == n-1 {
				if isIncome {
					totalIncome += int(rowCellData.(int))
					rowCell.Text = Colorize(rowCell.Text, BYellow)
				} else {
					totalExpense += int(rowCellData.(int))
				}
			}
			row = append(row, rowCell)
		}
		table.Body.Cells = append(table.Body.Cells, row)
	}

	// Generate Table Footer
	var footer *simpletable.Footer
	switch typeFlag {
	case "income":
		footer = &simpletable.Footer{
			Cells: []*simpletable.Cell{
				{
					Align: simpletable.AlignRight,
					Span:  4,
					Text:  fmt.Sprintf("%s", Colorize("TOTAL INCOME", BRed)),
				},
				{
					Align: simpletable.AlignRight,
					Text:  fmt.Sprintf("%s", Colorize(fmt.Sprintf("%d", totalIncome), BYellow)),
				},
			},
		}
	case "expense":
		footer = &simpletable.Footer{
			Cells: []*simpletable.Cell{
				{
					Align: simpletable.AlignRight,
					Span:  4,
					Text:  fmt.Sprintf("%s", Colorize("TOTAL EXPENSE", BRed)),
				},
				{
					Align: simpletable.AlignRight,
					Text:  fmt.Sprintf("%s", Colorize(fmt.Sprintf("%d", totalExpense), BWhite)),
				},
			},
		}
	default:
		footer = &simpletable.Footer{
			Cells: []*simpletable.Cell{
				{
					Align: simpletable.AlignRight,
					Span:  4,
					Text: fmt.Sprintf("%s\n%s",
						Colorize("T-INCOME", BYellow),
						Colorize("T-EXPENSE", BRed)),
				},
				{
					Align: simpletable.AlignRight,
					Text: fmt.Sprintf("%s\n%s",
						Colorize(fmt.Sprintf("%d", totalIncome), BYellow),
						Colorize(fmt.Sprintf("%d", totalExpense), BWhite)),
				},
			},
		}
	}

	table.Footer = footer
	table.SetStyle(style)
	fmt.Println(Colorize(models.HEADER_LINE, BGreen))
	table.Println()
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
