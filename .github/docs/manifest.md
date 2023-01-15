# Manifest File
Manifest file is a yaml file that contains the information about the functions that you want to create.
We describe the manifest file in the following sections.

## Manifest File Structure
The manifest file has the following structure:
```yaml
tag: string
source: string
destination: string
package_name: string
struct_name: string
table_name: string
db_library: string
test: bool
select:
  - type : string
    fields : ArrayOfString
    aggregate_fields :
      - function: string
        on: string
        as: string
    where_conditions :
      - column: string
        operator: string
    join_fields :
      - table: string
        as: string
        on_source: string
        on_join: string
        function: string
    order_by : string
    order_type : string
    limit : int
    group_by : ArrayOfString
```

### Tag
Tag is a string that is used to identify the version of the manifest file.

### Source
Source is a string that is used to identify the path of the source file. Source file is a file that contains the struct
that you want to create repository for it.

### Destination
Destination is a string that is used to identify the path of the destination file. Destination file is a file that 
contains the functions that you want to create.

### Package Name
Package name is a string that is used to identify the name of the package that you want to create functions in it.
As default, the package name is `repository`.

### Struct Name
Struct name is a string that is used to identify the name of the struct that you want to create functions for it.
As default, crafting table uses the first struct in the source file.

### Table Name
Table name is a string that is used to identify the name of the table that you want to create functions for it.
As default, crafting table uses the name of the snake case of the struct name.

### DB Library
DB Library is a string that is used to identify the name of the database library that you want to use in the functions.
Crafting table just support `sqlx` as the database library.

### Test
Test is a boolean that is used to identify if you want to create test file for the functions or not.

### Select
Select is an array of objects that is used to identify the information about the select functions 
that you want to create. Crafting table supports the following fields for select functions:

#### Type
Type is a string that is used to identify the type of the select function. Crafting table supports the following types:
- `select`
    - Get more than one row from the database.
- `get`
  - Get one row from the database. (In this case, if database returns more than one row, query will return an error.)

#### Fields
Fields is an array of strings that is used to identify the fields that you want to select from the database.
If you want to select all fields, you can use empty array as the value of the fields.

#### Aggregate Fields
Aggregate Fields is an array of objects that is used to identify the aggregate fields 
that you want to select from the database.
Crafting table supports the following fields for aggregate fields:
- `function`
    - The function that you want to use for aggregate fields. Crafting table supports the following functions:
        - COUNT
        - SUM
        - AVG
        - MAX
        - MIN
        - FIRST
        - LAST
- `on`
    - The field that you want to use for aggregate function.
- `as`
    - The name of the field that you want to use for the result of the aggregate function.

#### Where Conditions
Where Conditions is an array of objects that is used to identify the where conditions.
Crafting table supports the following fields for where conditions:
- `column`
    - The column that you want to use for where condition.
- `operator`
    - The operator that you want to use for where condition. Crafting table supports the following operators:
        - equal
        - not_equal
        - in
        - not_in
        - gt
        - gte
        - lt
        - lte 
        - is_null
        - is_not_null


#### Join Fields
Join Fields is an array of objects that is used to identify the join fields.
Crafting table supports the following fields for join fields:
- `table`
    - The table that you want to join with.
- `as`
    - The alias of the table that you want to join with.
- `on_source`
    - The column of the source table that you want to use for join condition.
- `on_join`
    - The column of the join table that you want to use for join condition.
- `function`
    - The function that you want to use for join fields. Crafting table supports the following functions:
        - Join
        - inner
        - fullOuter
        - rightOuter
        - leftOuter
        - full
        - left
        - right
        - natural
        - naturalLeft
        - naturalRight
        - naturalFull
        - cross

#### Order By
Order By is a string that is used to identify the field that you want to use for order by.

#### Order Type
Order Type is a string that is used to identify the type of the order by.
Crafting table supports the following fields for order by:
- `asc`
    - Order by ascending.
- `desc`
    - Order by descending.

#### Limit
Limit is an integer that is used to identify the limit of the select function.

#### Group By
Group By is an array of strings that is used to identify the fields that you want to use for group by.
