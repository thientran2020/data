package ui

type UI struct {
	Prompt *Prompt
	Table  *Table
}

func NewUI() *UI {
	return &UI{
		Prompt: NewPrompt(),
		Table:  NewTable(),
	}
}

func (ui *UI) DisplayHelp() {
	help := `
---------------------------------------------------------------------------------
				FINANCIAL CLI
   		 A great tool helps you manage your financial data ^^
---------------------------------------------------------------------------------

COMMANDS
	- data add
	USAGE: add financial data (expense or income)
	FLAGS:
		-i: add income
		-s: add subscription/membership data 
		    this will be handled automatically based on its billing cycle
		-t: add new trip
		-n: use today's date
		-y: use yesterday's date
		-c: set category by code (see 'data get -c' for codes)
		--cost: set the cost
		-d: set a description
		--yes: add without confirmation
	NOTE: 
		This command will update subscriptions automatically.
		Please run once if you see your subscriptions/memberships are not renewed yet in the financial table.

	- data show
	USAGE: display financial data in table format - current date by default 
	FLAGS:
	  -i: show only income data
		-e: show only expense data
		-m: show only data filtered by specified month
			Ex: 1, 2, 10,....
		-y: show only data filtered by specified year
			Ex: 2021, 2022,...
		-k: show only data filtered by specified keyword 
		-c: show current month's data
	
	- data get
	USAGE: get category mapping table or subscriptions' details
	FLAGS:
		-c: display category map for specific details
		-s: display subscriptions' details
		-t: display all trips
	
	- data search <keyword>
	USAGE: search specific financial key

	- data help
	USAGE: display app's instruction

---------------------------------------------------------------------------------
`
	ui.Prompt.DisplayMessage(help)
}
