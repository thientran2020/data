package models

var LABELS = map[string]map[string]string{
	"Expense": {
		"Description": "What did you spend money for",
		"Cost":        "Nice. How much did you spend",
	},
	"Income": {
		"Description": "What did you work for today",
		"Cost":        "Awesome. How much did you earn",
	},
}

var CATEGORY = []string{
	"Income",
	"Rent",
	"Utilities",
	"Insurance",
	"Vehicle Services",
	"Fuel - Car Wash",
	"Subcription",
	"Restaurants",
	"Amazon Shopping",
	"Merchandise",
	"Travel",
	"Personal",
}

const (
	BASE_FILEPATH             = "./finance/finance?????.csv"
	BASE_FILEPATH_SUBCRIPTION = "./finance/subscription.json"
	INSTRUCTION               = `
---------------------------------------------------------------------------------
				FINANCIAL CLI
   		 A great tool helps you manage your financial data ^^
---------------------------------------------------------------------------------

COMMANDS
	- data add
	USAGE: add financial data (expense or income)
	FLAGS:
		-s: Specify this is a subscription data which will be handled automatically by billing cycle

	- data show
	USAGE: display financial data in table format - current date by default 
	FLAGS:
		-m: Specify month you want to retrive financial data 
			Ex: 1, 2, 10,....
		-y: Specify year you want to retrive financial data
			Ex: 2021, 2022,...
		-h: Print default instruction

	- data help
	USAGE: display app's instruction
	
	- data category
	USAGE: display category map for specific details 

	- data search <keyword>
	USAGE: search specific financial key

---------------------------------------------------------------------------------
`

	CATEGORY_TABLE = `
			CATEGORY TABLE
	|------|-----------------------------------------|
	|  0   |  Income                                 |
	|------|-----------------------------------------|
	|  1   |  Rent                                   |
	|------|-----------------------------------------|
	|  2   |  Ultilites (AT&T, PG&E, Water)          |
	|------|-----------------------------------------|
	|  3   |  Insurance (Allstate, Medical, Dental)  |
	|------|-----------------------------------------|
	|  4   |  Vehicle Services                       |
	|------|-----------------------------------------|
	|  5   |  Fuel & Car Wash                        |
	|------|-----------------------------------------|
	|  6   |  Subscription & Membership              |
	|------|-----------------------------------------|
	|  7   |  Restaurants                            |
	|------|-----------------------------------------|
	|  8   |  Amazon Shopping                        |
	|------|-----------------------------------------|
	|  9   |  Merchandise                            |
	|------|-----------------------------------------|
	|  10  |  Travel                                 |
	|------|-----------------------------------------|
	|  11  |  Personal                               |
	|------|-----------------------------------------|

`
)
