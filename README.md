![alt text](https://user-images.githubusercontent.com/49960770/161902750-b853f8ad-5ab1-4676-9868-9be63ed3f8c3.jpeg)
![new logo](https://user-images.githubusercontent.com/49960770/212420359-aea418eb-e5e5-49b4-8564-8257d0cbae57.png)  
#### logo by: Midjourney


# Crafting Table
A crafting-table is third-party for creating database functions of repository layer.
It uses `sqlx` package to connect to database and querying to database.


## Installation
```bash
go install github.com/snapp-incubator/crafting-table@latest
```

# How to use Crafting Table?
The command for creating functions is as below:

```bash
crafting-table manifest apply <manifest-file-path>
```

## Manifest file
The manifest file is a yaml file that contains the information about the functions that you want to create.  
The manifest file is as below:

```yaml
tag: "Example"
source: "./example/src/example.go"
destination: "./example/dst/example.go"
package_name: "repository"
struct_name: "Example"
table_name: ""
db_library: "sqlx"
test: true
select:
#   get query
  - type : "get"
    fields : []
    aggregate_fields : []
    where_conditions : []
    join_fields : []
    order_by : "var1"
    order_type : "asc"
    limit : 0
    group_by : []

#   select query by conditions and group by and limit
  - type: "select"
    fields: []
    aggregate_fields: []
    where_conditions:
      - column: var1
        operator: gt
      - column: var2
        operator: equal
    join_fields: [ ]
    order_by: ""
    order_type: ""
    limit: 10
    group_by: ["var1", "var4"]

#   select query by aggregate functions
  - type: "select"
    fields: [ ]
    aggregate_fields:
      - function: COUNT
        on: var1
        as: count_var1
      - function: SUM
        on: var2
        as: sum_var2
      - function: AVG
        on: var1
        as: avg_var1
      - function: MAX
        on: var2
        as: max_var2
      - function: MIN
        on: var2
        as: min_var2
      - function: FIRST
        on: var1
        as: first_var1
      - function: LAST
        on: var2
        as: last_var2
    where_conditions: [ ]
    join_fields: [ ]
    order_by: ""
    order_type: ""
    limit: 0
    group_by: []
```

# Help Us
You can contribute to improving this tool by sending pull requests or issues on GitHub.  
Please send us your feedback. Thanks!
