### FINANCIAL CLI

## INSTALLATION

```
go get -u "github.com/thientran2020/financial-cli/models"
```

## COMMANDS
+ `data add` (used to add a single record: either income or expense) <br/>
**-s** (optional):&emsp; adding subscription

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
**-m**:&emsp; specify month <br/>
**-y**:&emsp; specify year <br/>
**-i**:&emsp; show only income data <br/>
**-e**:&emsp; show only expense data <br/>

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
**-k**:&emsp; keyword to look up 

+ `data help` (used to show general instruction) <br/>

## WHERE ARE DATA SAVED?

```
Copyright (c) 2022 Thien Tran
```

## COPYRIGHT
