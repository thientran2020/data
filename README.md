## DATA - FINANCIAL CLI TOOL

A great tool helps you manage your financial data ^^

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
go build -o data main.go && export PATH=$(pwd):$PATH
```

### COMMANDS
+ `data add` (used to add a single record: either income or expense) <br/>
*-s* &emsp;adding subscription

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
    Rent
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
|    08/07/2022   | ABC                                | $100   | Income            |
+-----+-----+-----+------------------------------------+--------+-------------------+
? Do you confirm to enter above record? [y/N]

 ➜ data add -s
Adding subscription....
Subcription list was created at ./finance/subscription.json
What is your new subscription/membership: Costco
✔ Expense
✔ Monthly
How much per billing period: 70
What was the start date (mm-dd-yyyy):

✔ Expense
+-----+-----+-----+------------------------------------+--------+-------------------+
|    08/07/2022   | Costco                             | $70    | Expense           |
+-----+-----+-----+------------------------------------+--------+-------------------+
? Do you confirm to enter above subscription? [y/N]
```

+ `data show` (used to retrieve your financial data) <br/>
*-i* &emsp;&nbsp;&nbsp;show only income data <br/>
*-e* &emsp;&nbsp;show only expense data <br/>
*-m* &emsp;show only data filtered by specified month <br/>
*-y* &emsp;&nbsp;show only data filtered by specified year <br/>
*-k* &emsp;&nbsp;show only data filtered by specified keyword <br/>

<ins>ESPECIALLY:</ins> This command will update subscriptions automatically. <br/> 
Please run once if you see your subscriptions/memberships are not renewed yet in the financial table.

```
➜ data show -m=7 -i
SHOW FINANCIAL DATA FOR 7/2022:

+---+------------+-------------+----------+------+
| # |    DATE    | DESCRIPTION | CATEGORY | COST |
+---+------------+-------------+----------+------+
| 1 | 07-29-2022 | ABC         | Income   |  100 |
+---+------------+-------------+----------+------+
|                            TOTAL INCOME |  100 |
+---+------------+-------------+----------+------+
```

+ `data get` (used to display more details) <br/>
*-c* &emsp;&nbsp;display category mapping table <br/>
*-s* &emsp;&nbsp;display current subscriptions' details <br/>
*-t* &emsp;&nbsp;display a list of trips to select for details <br/>

```
➜ data get -c

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
	|  12  |  Trip                                   |
	|------|-----------------------------------------|
```

+ `data search *agrs` (used to look up data for specific keyword) <br/>

+ `data help` (used to show general instruction) <br/>

### WHERE ARE DATA SAVED?
There are 2 types of csv files at `$HOME/finance` <br/>
- `~/finance/finance.csv`: this holds all financial data. <br/>
- `~/finance/finance-<year>.csv`: this holds specific-year financial data. <br/>
There are 2 json files at `$HOME/finance` <br/>
- `~/finance/subcriptions.json`: this holds all subscription/membership data (ex: monthly rent, Costco membership,...) where all of these will be automatically updated based on their billing cycle.<br/>
- `~/finance/trips.json`: this holds all trips data. <br/>
 
```
➜ ls -l ~/finance
-rw-rw-r--  1 ---------  staff   7850 Jul 29 19:38 finance.csv
-rw-rw-r--  1 ---------  staff   7850 Jul 29 19:38 finance_2022.csv
-rw-rw-r--  1 ---------  staff   7850 Dec 24 19:38 subscription.json
-rw-rw-r--  1 ---------  staff   7850 Dec 26 19:38 trip.json
```

### COPYRIGHT
```
Copyright (c) 2022 Thien Tran
```
