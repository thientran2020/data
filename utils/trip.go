package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/alexeyco/simpletable"
	m "github.com/thientran2020/financial-cli/models"
)

func AddNewTrip() {
	// interactive shell for information
	var tripName, startDate, endDate string
	tripName, _ = PromptEnter("Where do you visit", false)
	numOfParticipants, _ := NumberEnter("How many participants are there")
	if startDate, _ = DateEnter("When does your trip start (default is today's date)"); startDate == "" {
		startDate = time.Now().Format("01-02-2006")
	}
	if endDate, _ = DateEnter("When does your trip end (default is today's date)"); endDate == "" {
		endDate = time.Now().Format("01-02-2006")
	}

	// Create new trip, append to existing list and overwrite new file "trip.json" if existed
	trips := ReadTripJson(GetUserHomeDirectory() + m.BASE_FILEPATH_TRIP)
	trip := m.Trip{
		ID:            len(trips),
		Name:          tripName,
		NParticipants: int(numOfParticipants),
		Costs:         m.Costs{Shared: 0, Total: 0},
		Records:       []m.TripRecord{},
		StartDate:     startDate,
		EndDate:       endDate,
	}

	// Confirm record and enter to files
	if ConfirmYesNoPromt("Do you confirm to enter above record") {
		trips = append(trips, trip)
		WriteTripJson(GetUserHomeDirectory()+m.BASE_FILEPATH_TRIP, trips)
	} else {
		PrintCustomizedMessage("Trip ignored "+CheckMark, Red, true)
	}
}

func AddTripRecord(record m.Record) {
	trips := ReadTripJson(GetUserHomeDirectory() + m.BASE_FILEPATH_TRIP)
	tripsList := []string{}
	for _, trip := range trips {
		tripsList = append(tripsList, trip.Name)
	}
	selected, _ := InteractiveSelect("What trip you would like to add financial data", tripsList)
	shared := ConfirmYesNoPromt("Is it financially shared")

	for i := range trips {
		if trips[i].Name != selected {
			continue
		}
		record.Description = fmt.Sprintf("[%s] %s", selected, record.Description)
		tripRecord := m.TripRecord{
			Record: record,
			Shared: shared,
		}
		if tripRecord.Shared {
			trips[i].Costs.Shared += tripRecord.Cost
		}
		trips[i].Costs.Total += tripRecord.Cost
		trips[i].Records = append(trips[i].Records, tripRecord)
	}

	PrintSingleRecord(record, Green)
	if ConfirmYesNoPromt("Do you confirm to enter above record") {
		AddRecord(record)
		WriteTripJson(GetUserHomeDirectory()+m.BASE_FILEPATH_TRIP, trips)
	} else {
		PrintCustomizedMessage("Record ignored "+CheckMark, Red, true)
	}
}

func ReadTripJson(filepath string) []m.Trip {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		CreateFile(filepath)
		return ReadTripJson(filepath)
	}
	result := []m.Trip{}

	_ = json.Unmarshal([]byte(file), &result)
	return result
}

func WriteTripJson(filepath string, trips []m.Trip) {
	file, _ := json.MarshalIndent(trips, "", " ")
	err := ioutil.WriteFile(filepath, file, 0644)
	if err != nil {
		fmt.Printf("Error writing trip data %v\n", err)
	}
}

func PritnTrip() {
	trips := ReadTripJson(GetUserHomeDirectory() + m.BASE_FILEPATH_TRIP)
	tripsList := []string{}
	for _, trip := range trips {
		tripsList = append(tripsList, trip.Name)
	}
	selected, _ := InteractiveSelect("What trip you would like to add financial data", tripsList)
	data := Data{}
	for i := range trips {
		if trips[i].Name != selected {
			continue
		}
		for i, record := range trips[i].Records {
			formatted_date := fmt.Sprintf("%s-%s-%s",
				Colorize(GetStringDateFromNumber(record.Month), Yellow),
				GetStringDateFromNumber(record.Day),
				Colorize(GetStringDateFromNumber(record.Year), UGreen),
			)
			category := record.Category
			if record.Shared {
				category += " (shared)"
			}
			row := []interface{}{
				i,
				formatted_date,
				record.Description,
				category,
				record.Cost,
			}
			data = append(data, row)
		}
	}

	// Display trip as table
	intro := fmt.Sprintf("__________[%s] TRIP_________\n", strings.ToUpper(selected))
	headers := []string{"#", "DATE", "DESCRIPTION", "CATEGORY", "COST"}
	PrintCustomizedMessage(intro, Red, true)
	PrintTable(data, headers, "trip", simpletable.StyleDefault)
}
