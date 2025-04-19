package ui

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
)

type Prompt struct {
	// Any configuration or dependencies can be added here
}

func NewPrompt() *Prompt {
	return &Prompt{}
}

func (p *Prompt) AskForText(message string, defaultValue string) (string, error) {
	prompt := &survey.Input{
		Message: message,
		Default: defaultValue,
	}

	var result string
	err := survey.AskOne(prompt, &result)

	return result, err
}

func (p *Prompt) AskForNumber(message string) (int, error) {
	var result string
	prompt := &survey.Input{
		Message: message,
	}

	err := survey.AskOne(prompt, &result, survey.WithValidator(survey.Required))
	if err != nil {
		return 0, err
	}

	number, err := strconv.Atoi(result)
	if err != nil {
		return 0, fmt.Errorf("please enter a valid number: %w", err)
	}

	return number, nil
}

func (p *Prompt) AskForDate(message string) (string, error) {
	prompt := &survey.Input{
		Message: message,
	}

	var result string
	err := survey.AskOne(prompt, &result)
	if err != nil {
		return "", err
	}

	// If input is empty, return today's date
	if result == "" {
		now := time.Now()
		return fmt.Sprintf("%02d-%02d-%d", int(now.Month()), now.Day(), now.Year()), nil
	}

	// Validate the date if provided
	if err := validateDate(result); err != nil {
		return "", err
	}

	return result, nil
}

func (p *Prompt) AskForConfirmation(message string) (bool, error) {
	prompt := &survey.Confirm{
		Message: message,
		Default: false,
	}

	var result bool
	err := survey.AskOne(prompt, &result)

	return result, err
}

func (p *Prompt) AskForSelection(message string, options []string) (string, error) {
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}

	var result string
	err := survey.AskOne(prompt, &result)

	return result, err
}

func (p *Prompt) AskForSelectionWithIndex(message string, options []string) (string, int, error) {
	prompt := &survey.Select{
		Message: message,
		Options: options,
	}

	var result string
	err := survey.AskOne(prompt, &result)
	if err != nil {
		return "", -1, err
	}

	for i, option := range options {
		if option == result {
			return result, i, nil
		}
	}

	return result, -1, nil
}

func (p *Prompt) DisplayMessage(message string) {
	fmt.Println(message)
}

func (p *Prompt) DisplayError(message string) {
	fmt.Printf("Error: %s\n", message)
}

func validateDate(val interface{}) error {
	str, ok := val.(string)
	if !ok {
		return fmt.Errorf("date must be a string")
	}

	parts := strings.Split(str, "-")
	if len(parts) != 3 {
		return fmt.Errorf("date must be in MM-DD-YYYY format")
	}

	if len(parts[0]) != 2 || len(parts[1]) != 2 || len(parts[2]) != 4 {
		return fmt.Errorf("date must be in MM-DD-YYYY format")
	}

	month, err := strconv.Atoi(parts[0])
	if err != nil || month < 1 || month > 12 {
		return fmt.Errorf("month must be between 01 and 12")
	}

	day, err := strconv.Atoi(parts[1])
	if err != nil || day < 1 || day > 31 {
		return fmt.Errorf("day must be between 01 and 31")
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil || year < 1900 || year > 2100 {
		return fmt.Errorf("year must be between 1900 and 2100")
	}

	if (month == 4 || month == 6 || month == 9 || month == 11) && day > 30 {
		return fmt.Errorf("day must be between 01 and 30 for this month")
	}

	if month == 2 {
		isLeapYear := year%4 == 0 && (year%100 != 0 || year%400 == 0)
		if (isLeapYear && day > 29) || (!isLeapYear && day > 28) {
			return fmt.Errorf("invalid day for February")
		}
	}

	return nil
}
