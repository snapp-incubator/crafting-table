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
// Code generated by Crafting-Table. DO NOT EDIT
// Source code: https://github.com/snapp-incubator/crafting-table
package cmd

import (
	"fmt"
	"strings"
	"database/sql"
)

type UserSQLQueryBuilder interface{
	SetID(int) UserSQLQueryBuilder
	WhereIDEq(int) UserSQLQueryBuilder
	WhereIDGT(int) UserSQLQueryBuilder
	WhereIDGE(int) UserSQLQueryBuilder
	WhereIDLT(int) UserSQLQueryBuilder
	WhereIDLE(int) UserSQLQueryBuilder
	SetName(string) UserSQLQueryBuilder
	WhereNameEq(string) UserSQLQueryBuilder

	SetAge(int) UserSQLQueryBuilder
	WhereAgeEq(int) UserSQLQueryBuilder
	WhereAgeGT(int) UserSQLQueryBuilder
	WhereAgeGE(int) UserSQLQueryBuilder
	WhereAgeLT(int) UserSQLQueryBuilder
	WhereAgeLE(int) UserSQLQueryBuilder

	OrderByAsc(column UserColumn) UserSQLQueryBuilder
	OrderByDesc(column UserColumn) UserSQLQueryBuilder

	Limit(int) UserSQLQueryBuilder
	Offset(int) UserSQLQueryBuilder

	getPlaceholder() string

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

func (q *userSQLQueryBuilder) OrderByAsc(column UserColumn) UserSQLQueryBuilder {
	q.mode = "select"
	q.orderBy = append(q.orderBy, fmt.Sprintf("%s ASC", string(column)))
	return q
}

func (q *userSQLQueryBuilder) OrderByDesc(column UserColumn) UserSQLQueryBuilder {
	q.mode = "select"
	q.orderBy = append(q.orderBy, fmt.Sprintf("%s DESC", string(column)))
	return q
}

type userSQLQueryBuilder struct {
	mode string

	where userWhere

	set userSet

	orderBy []string
	groupBy string

	projected []string

	limit int
	offset int

	whereArgs []interface{}
	setArgs []interface{}
}

func NewUserQueryBuilder() UserSQLQueryBuilder {
	return &userSQLQueryBuilder{}
}

func (q *userSQLQueryBuilder) SQL() (string, error) {
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


func (q *userSQLQueryBuilder) sqlSelect() (string, error) {
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

	if len(q.orderBy) > 0 {
		base += fmt.Sprintf(" ORDER BY %s", strings.Join(q.orderBy, ", "))
	}

	if q.limit != 0 {
		base += " LIMIT " + fmt.Sprint(q.limit)
	}
	if q.offset != 0 {
		base += " OFFSET " + fmt.Sprint(q.offset)
	}
	return base, nil
}

func (q *userSQLQueryBuilder) Limit(l int) UserSQLQueryBuilder {
	q.mode = "select"
	q.limit = l
	return q
}

func (q *userSQLQueryBuilder) Offset(l int) UserSQLQueryBuilder {
	q.mode = "select"
	q.offset = l
	return q
}

func (q *userSQLQueryBuilder) sqlUpdate() (string, error) {
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

func (q *userSQLQueryBuilder) sqlDelete() (string, error) {
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

type userWhere struct {

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

func (q *userSQLQueryBuilder) WhereIDEq(ID int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, ID)
	q.where.ID.argument = q.getPlaceholder()
	q.where.ID.operator = "="
	return q
}

func (q *userSQLQueryBuilder) WhereNameEq(Name string) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, Name)
	q.where.Name.argument = q.getPlaceholder()
	q.where.Name.operator = "="
	return q
}

func (q *userSQLQueryBuilder) WhereAgeEq(Age int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, Age)
	q.where.Age.argument = q.getPlaceholder()
	q.where.Age.operator = "="
	return q
}




func (q *userSQLQueryBuilder) WhereIDGE(ID int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, ID)
	q.where.ID.argument = q.getPlaceholder()
	q.where.ID.operator = ">="
	return q
}

func (q *userSQLQueryBuilder) WhereIDGT(ID int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, ID)
	q.where.ID.argument = q.getPlaceholder()
	q.where.ID.operator = ">"
	return q
}

func (q *userSQLQueryBuilder) WhereIDLE(ID int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, ID)
	q.where.ID.argument = q.getPlaceholder()
	q.where.ID.operator = "<="
	return q
}

func (q *userSQLQueryBuilder) WhereIDLT(ID int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, ID)
	q.where.ID.argument = q.getPlaceholder()
	q.where.ID.operator = "<"
	return q
}



func (q *userSQLQueryBuilder) WhereAgeGE(Age int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, Age)
	q.where.Age.argument = q.getPlaceholder()
	q.where.Age.operator = ">="
	return q
}

func (q *userSQLQueryBuilder) WhereAgeGT(Age int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, Age)
	q.where.Age.argument = q.getPlaceholder()
	q.where.Age.operator = ">"
	return q
}

func (q *userSQLQueryBuilder) WhereAgeLE(Age int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, Age)
	q.where.Age.argument = q.getPlaceholder()
	q.where.Age.operator = "<="
	return q
}

func (q *userSQLQueryBuilder) WhereAgeLT(Age int) UserSQLQueryBuilder {
	q.whereArgs = append(q.whereArgs, Age)
	q.where.Age.argument = q.getPlaceholder()
	q.where.Age.operator = "<"
	return q
}

type userSet struct {
	ID string
	Name string
	Age string
}

func (q *userSQLQueryBuilder) SetID(ID int) UserSQLQueryBuilder {
	q.mode = "update"
	q.setArgs = append(q.setArgs, ID)
	q.set.ID = q.getPlaceholder()
	return q
}

func (q *userSQLQueryBuilder) SetName(Name string) UserSQLQueryBuilder {
	q.mode = "update"
	q.setArgs = append(q.setArgs, Name)
	q.set.Name = q.getPlaceholder()
	return q
}

func (q *userSQLQueryBuilder) SetAge(Age int) UserSQLQueryBuilder {
	q.mode = "update"
	q.setArgs = append(q.setArgs, Age)
	q.set.Age = q.getPlaceholder()
	return q
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
	var q User
	err := row.Scan(
		&q.ID,
		&q.Name,
		&q.Age,

	)
	if err != nil {
		return User{}, err
	}

	return q, nil
}

func (q User) Values() []interface{} {
	var values []interface{}
	values = append(values, &q.ID)
	values = append(values, &q.Name)
	values = append(values, &q.Age)

	return values
}

func (q *userSQLQueryBuilder) getPlaceholder() string {
	return "?"
}

func (q *userSQLQueryBuilder) Update(db *sql.DB) (sql.Result, error) {
	q.mode = "update"
	args := append(q.setArgs, q.whereArgs...)
	query, err := q.SQL()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}

func (q *userSQLQueryBuilder) Delete(db *sql.DB) (sql.Result, error) {
	q.mode = "delete"
	query, err := q.SQL()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, q.whereArgs...)
}

func (q *userSQLQueryBuilder) Fetch(db *sql.DB) ([]User, error) {
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

func (q *userSQLQueryBuilder) First(db *sql.DB) (User, error) {
	q.mode = "select"
	q.orderBy = []string{"ORDER BY id ASC"}
	q.Limit(1)
	query, err := q.SQL()
	if err != nil {
		return User{}, err
	}
	row := db.QueryRow(query, q.whereArgs...)
	if row.Err() != nil {
		return User{}, row.Err()
	}
	return UserFromRow(row)
}


func (q *userSQLQueryBuilder) Last(db *sql.DB) (User, error) {
	q.mode = "select"
	q.orderBy = []string{"ORDER BY id DESC"}
	q.Limit(1)
	query, err := q.SQL()
	if err != nil {
		return User{}, err
	}
	row := db.QueryRow(query, q.whereArgs...)
	if row.Err() != nil {
		return User{}, row.Err()
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
