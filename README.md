![alt text](https://raw.githubusercontent.com/n25a/crafting-table/logo/logo.jpeg?token=GHSAT0AAAAAABQEANQSFSJIZ2NPCILY2TLMYSMGALA)

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

# Help Us
You can contribute to improving this tool by sending pull requests or issues on GitHub.
Please send us your feedback. Thanks!
