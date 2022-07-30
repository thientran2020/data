## FINANCIAL CLI

A great tool helps you manage your financial data ^^

### INSTALLATION

```
git clone git@github.com:thientran2020/financial-cli.git
cd financial-cli
go build -o data main.go && export PATH=$(pwd):$PATH
```

### COMMANDS
+ `data add` (used to add a single record: either income or expense) <br/>
*-s* (optional): &emsp;adding subscription

```
➜ data add
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
```

+ `data show` (used to retrieve your financial data) <br/>
*-m*: &emsp;specify month <br/>
*-y*: &emsp;&nbsp;specify year <br/>
*-i*: &emsp;&nbsp;&nbsp;show only income data <br/>
*-e*: &emsp;&nbsp;show only expense data <br/>

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

+ `data category` (used to display a table of app categories) <br/>

```
➜ data category

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
```

+ `data search` <br/>
*-k*:&emsp; keyword to look up 

+ `data help` (used to show general instruction) <br/>

### WHERE ARE DATA SAVED?
There are 2 types of csv files: <br/>
- `./finance/finance.csv`: this holds all financial data. <br/>
- `./finance/finance-<year>.csv`: this holds specific-year financial data. <br/>
- `./finance/subcriptions.json`: this holds all subcriptionb/membership data 
where all of these will be automatically updated based on their billing cycle.<br/>
 

```
➜ ls -l ./finance
-rw-rw-r--  1 ---------  staff   7850 Jul 29 19:38 finance.csv
-rw-rw-r--  1 ---------  staff   7850 Jul 29 19:38 finance_2022.csv
```

### COPYRIGHT
```
Copyright (c) 2022 Thien Tran
```
