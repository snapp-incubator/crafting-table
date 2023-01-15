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



