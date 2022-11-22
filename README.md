![alt text](https://user-images.githubusercontent.com/49960770/161902750-b853f8ad-5ab1-4676-9868-9be63ed3f8c3.jpeg)

# Crafting Table
A function creation tool for querying to database from golang struct.
It uses `sqlx` package to connect to database.


# How to use Crafting Table?

The Command for creating functions is as below.

```bash
crafting-table generate \
      -s ./SOURCE/PATH/STRUCTFILE.go \
      -d ./DESTINATIOIN/PATH/REPOSITORY.go \
      --get "[ GetByVar1, (GetByVar1, GetByVar3) ]" \
      --update "[[(UpdateByVar3),(UpdateVar2, UpdateVar1)], [(UpdateByVar2, UpdateByVar3), (UpdateVar1)]]" \
      --create true
```

## Flags
|    Flag    | Description                                                                                                                                                                                                                                                                                                                                                                   |
|:----------:|:------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
|    `-s`    | Source file path.                                                                                                                                                                                                                                                                                                                                                             |
|    `-d`    | Destination file path.                                                                                                                                                                                                                                                                                                                                                        |
|  `--get`   | Get function name and arguments. You can use multiple arguments separated with comma `,` to define multiple `Get` functions. Also, You can use parentheses to get records by more than one parameter. If you don't define this flag, `Get` function will not be created.                                                                                                      |
| `--update` | Update function name and arguments. You can use multiple arguments separated with comma `,` and brackets `[ ]` to define multiple `Update` functions. You should define two pairs of parentheses in each bracket, the former for updating by parameters and the latter for updating fields parameters. If you don't define this flag, `Update` function will not be created.  |
| `--create` | Create function name and arguments. If you don't define this flag, `Create` function will not be created.                                                                                                                                                                                                                                                                     |

# Note
In each flag, you should use the database name as parameter.

## Example

Suppose that the following struct is a table in the database.

```go
package src

type Example struct {
	Var1 int    `db:"var1"`
	Var2 string `db:"var2"`
	Var3 bool   `db:"var3"`
}
```

You can create a functions to query to database using `crafting-table` as below.

Command:
```bash
crafting-table generate \
    -s ./example/src/example.go \
    -d ./example/dst/example.go \
    --get "[ var1, (var1, var3) ]" \
    --update "[[(var3), (var2, var1)], [(var2, var3), (var1)]]" \
    --create true
```

Result:
```go
package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/n25a/crafting-table/example/src"

	"github.com/jmoiron/sqlx"
)

type Example interface {
	Create(ctx context.Context, example *src.Example) error
	GetExamples(ctx context.Context) (*[]src.Example, error)
	GetByVar1(ctx context.Context, var1 int) (*src.Example, error)
	GetByVar1AndVar3(ctx context.Context, var1 int, var3 bool) (*src.Example, error)
	Update(ctx context.Context, var1 int, example src.Example) (int64, error)
	UpdateVar2AndVar1(ctx context.Context, var3 bool, var2 string, var1 int) (int64, error)
	UpdateVar1(ctx context.Context, var2 string, var3 bool, var1 int) (int64, error)
}

var ErrExampleNotFound = errors.New("example not found")

type mysqlExample struct {
	db *sqlx.DB
}

func NewMySQLExample(db *sqlx.DB) Example {
	return &mysqlExample{db: db}
}

func (r *mysqlExample) Create(ctx context.Context, example *src.Example) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO example ("+
		"var1, var2, var3"+
		") VALUES ("+
		":var1, :var2, :var3)",
		example)

	if err != nil {
		return err
	}

	return nil
}

func (r *mysqlExample) Update(ctx context.Context, var1 int, example src.Example) (int64, error) {
	example.Var1 = var1

	result, err := r.db.NamedExecContext(ctx, "UPDATE example "+
		"SET"+
		"var1 = :var1, var2 = :var2, var3 = :var3 "+
		"WHERE var1 = :var1",
		example,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *mysqlExample) UpdateVar2AndVar1(ctx context.Context, var3 bool, var2 string, var1 int) (int64, error) {
	query := "UPDATE example SET " +
		"var2 = ?, var1 = ? " +
		"WHERE var3 = ?;"

	result, err := r.db.ExecContext(ctx, query, var2, var1, var3)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *mysqlExample) UpdateVar1(ctx context.Context, var2 string, var3 bool, var1 int) (int64, error) {
	query := "UPDATE example SET " +
		"var1 = ? " +
		"WHERE var2 = ? AND var3 = ?;"

	result, err := r.db.ExecContext(ctx, query, var1, var2, var3)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

func (r *mysqlExample) GetExamples(ctx context.Context) (*[]src.Example, error) {
	var example []src.Example
	err := r.db.SelectContext(ctx, &example, "SELECT * from example")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrExampleNotFound
		}

		return nil, err
	}

	return &example, nil
}

func (r *mysqlExample) GetByVar1(ctx context.Context, var1 int) (*src.Example, error) {
	var example src.Example

	err := r.db.GetContext(ctx, &example, "SELECT * FROM example "+
		"WHERE var1 = ?",
		var1,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrExampleNotFound
		}

		return nil, err
	}

	return &example, nil
}

func (r *mysqlExample) GetByVar1AndVar3(ctx context.Context, var1 int, var3 bool) (*src.Example, error) {
	var example src.Example

	err := r.db.GetContext(ctx, &example, "SELECT * FROM example "+
		"WHERE var1 = ? AND var3 = ?",
		var1, var3,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrExampleNotFound
		}

		return nil, err
	}

	return &example, nil
}

```
# Query builder generator based on structs
There is another feature in `crafting-table` which lets you generate a type-safe and fast query builder
based on your model structs.
To use this feature you should annotate your structs with `ct: model` doc comment, and then run
```bash
crafting-table query-builder -f <your file name> -d <your sql dialect, defaults to mysql>
```
then it will generate a query builder for you.

## Example
imagine you have following go file,
```go
// user.go
package cmd

// ct: model
type User struct {
	ID   int
	Name string
	Age  int
}
```
and now you run `crafting-table query-builder -f user.go -d mysql`,
this command will generate a new file for you named `user_ct_gen.go`.
```go
// user_ct_gen.go
// Code generated by Crafting-Table. DO NOT EDIT
// Source code: https://github.com/snapp-incubator/crafting-table
package cmd

import (
    "fmt"
    "strings"
    "database/sql"
)

type UserSQLQueryBuilder interface{
	
	WhereIDEq(int) UserSQLQueryBuilder
	
	WhereIDGT(int) UserSQLQueryBuilder
	WhereIDGE(int) UserSQLQueryBuilder
	WhereIDLT(int) UserSQLQueryBuilder
	WhereIDLE(int) UserSQLQueryBuilder
	
	SetID(int) UserSQLQueryBuilder
	
	WhereNameEq(string) UserSQLQueryBuilder
	
	SetName(string) UserSQLQueryBuilder
	
	WhereAgeEq(int) UserSQLQueryBuilder
	
	WhereAgeGT(int) UserSQLQueryBuilder
	WhereAgeGE(int) UserSQLQueryBuilder
	WhereAgeLT(int) UserSQLQueryBuilder
	WhereAgeLE(int) UserSQLQueryBuilder
	
	SetAge(int) UserSQLQueryBuilder
	

	OrderByAsc(column UserColumn) UserSQLQueryBuilder
	OrderByDesc(column UserColumn) UserSQLQueryBuilder

	Limit(int) UserSQLQueryBuilder
	Offset(int) UserSQLQueryBuilder

    getPlaceholder() string
	
	// finishers
	First(db *sql.DB) (User, error)
	Last(db *sql.DB) (User, error)
	Update(db *sql.DB) (sql.Result, error)
	Delete(db *sql.DB) (sql.Result, error)
	Fetch(db *sql.DB) ([]User, error)
	SQL() (string, error)
}

type UserColumn string 
var UserColumns = struct {
	
	ID UserColumn
	
	Name UserColumn
	
	Age UserColumn
	
}{
	
	ID: UserColumn("id"),
	
	Name: UserColumn("name"),
	
	Age: UserColumn("age"),
	
}

func (q *__UserSQLQueryBuilder) OrderByAsc(column UserColumn) UserSQLQueryBuilder {
    q.mode = "select"
	q.orderby = append(q.orderby, fmt.Sprintf("%s ASC", string(column)))
	return q
}
func (q *__UserSQLQueryBuilder) OrderByDesc(column UserColumn) UserSQLQueryBuilder {
    q.mode = "select"
	q.orderby = append(q.orderby, fmt.Sprintf("%s DESC", string(column)))
	return q
}

type __UserSQLQueryBuilder struct {
	mode string

    where __UserWhere

	set __UserSet

	orderby []string
	groupby string

	projected []string

	limit int
	offset int

	whereArgs []interface{}
    setArgs []interface{}
}

func UserQueryBuilder() UserSQLQueryBuilder {
	return &__UserSQLQueryBuilder{}
}



func (q *__UserSQLQueryBuilder) SQL() (string, error) {
	if q.mode == "select" {
		return q.sqlSelect()
	} else if q.mode == "update" {
		return q.sqlUpdate()
	} else if q.mode == "delete" {
		return q.sqlDelete()
	} else {
		return "", fmt.Errorf("unsupported query mode '%s'", q.mode)
	}
}


func (q *__UserSQLQueryBuilder) sqlSelect() (string, error) {
	if q.projected == nil {
		q.projected = append(q.projected, "*")
	}
	base := fmt.Sprintf("SELECT %s FROM users", strings.Join(q.projected, ", "))

	var wheres []string 
	
	if q.where.ID.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "id", q.where.ID.operator, fmt.Sprint(q.where.ID.argument)))
	}
	
	if q.where.Name.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "name", q.where.Name.operator, fmt.Sprint(q.where.Name.argument)))
	}
	
	if q.where.Age.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "age", q.where.Age.operator, fmt.Sprint(q.where.Age.argument)))
	}
	
	if len(wheres) > 0 {
		base += " WHERE " + strings.Join(wheres, " AND ")
	}

	if len(q.orderby) > 0 {
		base += fmt.Sprintf(" ORDER BY %s", strings.Join(q.orderby, ", "))
	}

	if q.limit != 0 {
		base += " LIMIT " + fmt.Sprint(q.limit)
	}
	if q.offset != 0 {
		base += " OFFSET " + fmt.Sprint(q.offset)
	}
	return base, nil
}

func (q *__UserSQLQueryBuilder) Limit(l int) UserSQLQueryBuilder {
	q.mode = "select"
	q.limit = l	
	return q
}
func (q *__UserSQLQueryBuilder) Offset(l int) UserSQLQueryBuilder {
	q.mode = "select"
	q.offset = l
	return q
}

func (q *__UserSQLQueryBuilder) sqlUpdate() (string, error) {
	base := fmt.Sprintf("UPDATE users ")

	var wheres []string 
    var sets []string 

    
	if q.where.ID.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "id", q.where.ID.operator, fmt.Sprint(q.where.ID.argument)))
	}
	if q.set.ID != "" {
		sets = append(sets, fmt.Sprintf("%s = %s", "id", fmt.Sprint(q.set.ID)))
	}
    
	if q.where.Name.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "name", q.where.Name.operator, fmt.Sprint(q.where.Name.argument)))
	}
	if q.set.Name != "" {
		sets = append(sets, fmt.Sprintf("%s = %s", "name", fmt.Sprint(q.set.Name)))
	}
    
	if q.where.Age.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "age", q.where.Age.operator, fmt.Sprint(q.where.Age.argument)))
	}
	if q.set.Age != "" {
		sets = append(sets, fmt.Sprintf("%s = %s", "age", fmt.Sprint(q.set.Age)))
	}
    

	if len(sets) > 0 {
		base += "SET " + strings.Join(sets, " , ")
	}

	if len(wheres) > 0 {
		base += " WHERE " + strings.Join(wheres, " AND ")
	}

	

	return base, nil
}

func (q *__UserSQLQueryBuilder) sqlDelete() (string, error) {
    base := fmt.Sprintf("DELETE FROM users")

	var wheres []string 
	
	if q.where.ID.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "id", q.where.ID.operator, fmt.Sprint(q.where.ID.argument)))
	}
	
	if q.where.Name.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "name", q.where.Name.operator, fmt.Sprint(q.where.Name.argument)))
	}
	
	if q.where.Age.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "age", q.where.Age.operator, fmt.Sprint(q.where.Age.argument)))
	}
	
	if len(wheres) > 0 {
		base += " WHERE " + strings.Join(wheres, " AND ")
	}

	return base, nil
}

type __UserWhere struct {
	
	ID struct {
        argument interface{}
        operator string
    }
	
	Name struct {
        argument interface{}
        operator string
    }
	
	Age struct {
        argument interface{}
        operator string
    }
	
}

func (m *__UserSQLQueryBuilder) WhereIDEq(ID int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, ID)
    m.where.ID.argument = m.getPlaceholder()
    m.where.ID.operator = "="
	return m
}

func (m *__UserSQLQueryBuilder) WhereNameEq(Name string) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, Name)
    m.where.Name.argument = m.getPlaceholder()
    m.where.Name.operator = "="
	return m
}

func (m *__UserSQLQueryBuilder) WhereAgeEq(Age int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, Age)
    m.where.Age.argument = m.getPlaceholder()
    m.where.Age.operator = "="
	return m
}




func (m *__UserSQLQueryBuilder) WhereIDGE(ID int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, ID)
    m.where.ID.argument = m.getPlaceholder()
    m.where.ID.operator = ">="
	return m
}
func (m *__UserSQLQueryBuilder) WhereIDGT(ID int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, ID)
    m.where.ID.argument = m.getPlaceholder()
    m.where.ID.operator = ">"
	return m
}
func (m *__UserSQLQueryBuilder) WhereIDLE(ID int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, ID)
    m.where.ID.argument = m.getPlaceholder()
    m.where.ID.operator = "<="
	return m
}
func (m *__UserSQLQueryBuilder) WhereIDLT(ID int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, ID)
    m.where.ID.argument = m.getPlaceholder()
    m.where.ID.operator = "<"
	return m
}





func (m *__UserSQLQueryBuilder) WhereAgeGE(Age int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, Age)
    m.where.Age.argument = m.getPlaceholder()
    m.where.Age.operator = ">="
	return m
}
func (m *__UserSQLQueryBuilder) WhereAgeGT(Age int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, Age)
    m.where.Age.argument = m.getPlaceholder()
    m.where.Age.operator = ">"
	return m
}
func (m *__UserSQLQueryBuilder) WhereAgeLE(Age int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, Age)
    m.where.Age.argument = m.getPlaceholder()
    m.where.Age.operator = "<="
	return m
}
func (m *__UserSQLQueryBuilder) WhereAgeLT(Age int) UserSQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, Age)
    m.where.Age.argument = m.getPlaceholder()
    m.where.Age.operator = "<"
	return m
}



type __UserSet struct {
	
	ID string
    
	Name string
    
	Age string
    
}

func (m *__UserSQLQueryBuilder) SetID(ID int) UserSQLQueryBuilder {
	m.mode = "update"
    m.setArgs = append(m.setArgs, ID)
	m.set.ID = m.getPlaceholder()
	return m
}

func (m *__UserSQLQueryBuilder) SetName(Name string) UserSQLQueryBuilder {
	m.mode = "update"
    m.setArgs = append(m.setArgs, Name)
	m.set.Name = m.getPlaceholder()
	return m
}

func (m *__UserSQLQueryBuilder) SetAge(Age int) UserSQLQueryBuilder {
	m.mode = "update"
    m.setArgs = append(m.setArgs, Age)
	m.set.Age = m.getPlaceholder()
	return m
}


func UsersFromRows(rows *sql.Rows) ([]User, error) {
    var Users []User
    for rows.Next() {
        var m User
        err := rows.Scan(
            
            &m.ID,
            
            &m.Name,
            
            &m.Age,
            
        )
        if err != nil {
            return nil, err
        }
        Users = append(Users, m)
    }
    return Users, nil
}

func UserFromRow(row *sql.Row) (User, error) {
    if row.Err() != nil {
        return User{}, row.Err()
    }
    var m User
    err := row.Scan(
        
        &m.ID,
        
        &m.Name,
        
        &m.Age,
        
    )
    if err != nil {
        return User{}, err
    }
    return m, nil
}

func (m User) Values() []interface{} {
    var values []interface{}
	
	values = append(values, &m.ID)
	
	values = append(values, &m.Name)
	
	values = append(values, &m.Age)
	
    return values
}


func (q *__UserSQLQueryBuilder) getPlaceholder() string {
	return "?"
}




func (q *__UserSQLQueryBuilder) Update(db *sql.DB) (sql.Result, error) {
	q.mode = "update"
	args := append(q.setArgs, q.whereArgs...)
	query, err := q.SQL()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}

func (q *__UserSQLQueryBuilder) Delete(db *sql.DB) (sql.Result, error) {
	q.mode = "delete"
	query, err := q.SQL()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, q.whereArgs...)
}

func (q *__UserSQLQueryBuilder) Fetch(db *sql.DB) ([]User, error) {
	q.mode = "select"
	query, err := q.SQL()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query, q.whereArgs...)
	if err != nil {
		return nil, err
	}
	return UsersFromRows(rows)
}

func (q *__UserSQLQueryBuilder) First(db *sql.DB) (User, error) {
	q.mode = "select"
	q.orderby = []string{"ORDER BY id ASC"}
	q.Limit(1)
	query, err := q.SQL()
	if err != nil {
		return User{}, err
	}
	row := db.QueryRow(query, q.whereArgs...)
	if row.Err() != nil {
		return User {}, row.Err()
	}
	return UserFromRow(row)
}


func (q *__UserSQLQueryBuilder) Last(db *sql.DB) (User, error) {
	q.mode = "select"
	q.orderby = []string{"ORDER BY id DESC"}
	q.Limit(1)
	query, err := q.SQL()
	if err != nil {
		return User{}, err
	}
	row := db.QueryRow(query, q.whereArgs...)
	if row.Err() != nil {
		return User {}, row.Err()
	}
	return UserFromRow(row)
}
```
and example usage of this query builder something like
```go
func some(db *sql.DB) {
	q := UserQueryBuilder()
	users, err := q.WhereAgeEq(10).Limit(100).Offset(1000).OrderByAsc(UserColumns.ID).Fetch(db)
}
```

# Help Us
You can contribute to improving this tool by sending pull requests or issues on GitHub.
Please send us your feedback. Thanks!
