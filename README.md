## DATA - FINANCIAL CLI TOOL

A great tool helps you manage your financial data with beautiful color-coded output ^^

### INSTALLATION
There are 2 ways: either use `brew tap` (recommended) or clone the project with `git`. </br>
1. Use `brew tap`. Open your terminal, run the following commands:
```
 ➜ brew tap thientran2020/homebrew-thientran2020
```

```
 ➜ brew install thientran2020/thientran2020/data
```
Verify if it was already installed successfully:
```
 ➜ data help

---------------------------------------------------------------------------------
				FINANCIAL CLI
   		 A great tool helps you manage your financial data ^^
---------------------------------------------------------------------------------
.....
```

2. Clone project to your local machine
```
git clone git@github.com:thientran2020/data.git
cd data
cd cmd/data && go build -o data main.go && export PATH=$(pwd):$PATH
```

### FEATURES
- Beautiful color-coded output for better readability
  - Income amounts in green
  - Expense amounts in red
  - Headers in bold
  - Net amounts color-coded based on value
- Consistent date formatting (MM-DD-YYYY)
- Properly aligned financial data display
- Subscription management
- Trip expense tracking
- Category-based organization

### COMMANDS
+ `data add` (used to add a single record: either income or expense) <br/>
*-s* &emsp;adding subscription <br/>
*-t* &emsp;adding trip <br/>
*-i* &emsp;adding income <br/>
*-c* &emsp;specify category code <br/>
*-d* &emsp;add description <br/>
*-n* &emsp;set date to today <br/>
*-y* &emsp;set date to yesterday

```
➜ data add
Enter date: : 12-19-2022
Use the arrow keys to navigate: ↓ ↑ → ←
? What type of financial data are you entering:
  ▸ Income
    Expense
? What did you work for today: ABC
Awesome. How much did you earn: 100
Use the arrow keys to navigate: ↓ ↑ → ←
? Pick the category that describe best your entered data:
  ▸ Income
    Mortgage
    Utilities
    Insurance
    Vehicle Services
    Fuel - Car Wash
    Subcription
    Restaurants
    Amazon Shopping
    Merchandise
    Travel
    Personal
    Trip

✔ Income
+-----+-----+-----+------------------------------------+--------+-------------------+
|    12/19/2022   | ABC                                | $100   | Income            |
+-----+-----+-----+------------------------------------+--------+-------------------+
? Do you confirm to enter above record? [y/N]
```

+ `data show` (used to retrieve your financial data) <br/>
*-i* &emsp;&nbsp;&nbsp;show only income data <br/>
*-e* &emsp;&nbsp;show only expense data <br/>
*-m* &emsp;show only data filtered by specified month <br/>
*-y* &emsp;&nbsp;show only data filtered by specified year <br/>
*-k* &emsp;&nbsp;show only data filtered by specified keyword <br/>
*-c* &emsp;&nbsp;show current month's records

<ins>ESPECIALLY:</ins> This command will update subscriptions automatically. <br/> 
Please run once if you see your subscriptions/memberships are not renewed yet in the financial table.

Example output with color formatting:
```
➜ data show -m=7 -y=2024
SHOW FINANCIAL DATA FOR 7/2024:

+---+------------+-------------+----------+------+
| # |    DATE    | DESCRIPTION | CATEGORY | COST |
+---+------------+-------------+----------+------+
| 1 | 07-09-2024 | Salary      | Income   | $5000|
| 2 | 07-15-2024 | Groceries   | Food     | $150 |
+---+------------+-------------+----------+------+

TOTAL INCOME:      $5000
TOTAL EXPENSE:     $150
NET:               $4850
```

+ `data get` (used to retrieve category mapping or subscription details) <br/>
*-c* &emsp;show category mapping table <br/>
*-s* &emsp;show subscription details <br/>
*-t* &emsp;show trip details

+ `data search <keyword>` (used to search records by keyword) <br/>
Searches through descriptions and categories for the specified keyword.

### DATA STORAGE
All financial data is stored in JSON format under the `./finance` directory:
- Yearly records: `./finance/YYYY.json`
- Subscriptions: `./finance/subscription.json`
- Trips: `./finance/trip.json`

### CATEGORIES
The tool comes with predefined categories for better organization:
- Income
- Mortgage
- Utilities
- Insurance
- Vehicle Services
- Fuel - Car Wash
- Subscription
- Restaurants
- Amazon Shopping
- Merchandise
- Travel
- Personal
- Trip

You can view the complete category mapping using `data get -c`.

### COPYRIGHT
```
Copyright (c) 2022 Thien Tran
```
