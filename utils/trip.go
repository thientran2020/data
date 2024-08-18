package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/alexeyco/simpletable"
	m "github.com/thientran2020/financial-cli/models"
)

func AddNewTrip() {
	// interactive shell for information
	tripName := PromptEnter("Where do you visit?")
	numOfParticipants := NumberEnter("How many participants are there?")
	startDate := DateEnter("When does your trip start (default is today's date)?")
	endDate := DateEnter("When does your trip end (default is today's date)?")

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
	if ConfirmYesNoPromt("Do you confirm to enter above record?") {
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
	selected := InteractiveSelect("What trip you would like to add financial data?", tripsList)
	shared := ConfirmYesNoPromt("Is it financially shared?")

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
	if ConfirmYesNoPromt("Do you confirm to enter above record?") {
		AddRecord(record)
		WriteTripJson(GetUserHomeDirectory()+m.BASE_FILEPATH_TRIP, trips)
	} else {
		PrintCustomizedMessage("Record ignored "+CheckMark, Red, true)
	}
}

func ReadTripJson(filepath string) []m.Trip {
	file, err := os.ReadFile(filepath)
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
	err := os.WriteFile(filepath, file, 0644)
	if err != nil {
		fmt.Printf("Error writing trip data %v\n", err)
	}
}

func PrintTrip(onlyShared bool) {
	trips := ReadTripJson(GetUserHomeDirectory() + m.BASE_FILEPATH_TRIP)
	tripsList := []string{}
	for _, trip := range trips {
		tripsList = append(tripsList, trip.Name)
	}

	var trip m.Trip
	selected := InteractiveSelect("What trip you would like to display financial data?", tripsList)
	data := String2D{}
	for i := range trips {
		if trips[i].Name != selected {
			continue
		}
		trip = trips[i]
		for _, record := range trips[i].Records {
			category := record.Category
			if onlyShared && !record.Shared {
				continue
			}
			if record.Shared {
				category += " (shared)"
			}
			row := []string{
				strconv.Itoa(record.Year),
				strconv.Itoa(record.Month),
				strconv.Itoa(record.Day),
				record.Description,
				category,
				strconv.Itoa(record.Cost),
			}
			data = append(data, row)
		}
	}
	sort.Sort(data)

	tableData := Data{}
	for i, row := range data {
		formatted_date := fmt.Sprintf("%s-%s-%s",
			Colorize(GetStringDateFromString(row[month]), Yellow),
			GetStringDateFromString(row[day]),
			Colorize(GetStringDateFromString(row[year]), UGreen),
		)
		cost, _ := strconv.Atoi(row[5])
		tableRow := []interface{}{
			i,
			formatted_date,
			row[3],
			row[4],
			cost,
		}
		tableData = append(tableData, tableRow)
	}

	// Display trip as table
	intro := fmt.Sprintf("__________TRIP TO %s_________\n\n   - From %s to %s", strings.ToUpper(selected), trip.StartDate, trip.EndDate)
	summary := fmt.Sprintf("   - Total shared cost is $%d, shared by %d --> Each pays $%d\n", trip.Costs.Shared, trip.NParticipants, trip.Costs.Shared/trip.NParticipants)
	if !onlyShared {
		summary += fmt.Sprintf("   - Total trip cost is $%d\n", trip.Costs.Total)
	}
	headers := []string{"#", "DATE", "DESCRIPTION", "CATEGORY", "COST"}
	PrintCustomizedMessage("\n"+CenterString(intro, len(m.DASH)), Red, true)
	PrintCustomizedMessage(summary, Yellow, true)
	PrintTable(tableData, headers, "trip", simpletable.StyleDefault)
}
