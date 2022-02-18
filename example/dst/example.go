package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/n25a/repogen/example/src"
)

type Example interface {
	Create(ctx context.Context, example *src.Example) error
	GetExamples(ctx context.Context) (*[]src.Example, error)
	GetByVar2(ctx context.Context, var2 string, var3 bool) (*src.Example, error)
	GetByVar1(ctx context.Context, var1 int) (*src.Example, error)
	GetByVar3(ctx context.Context, var3 bool, var1 int) (*src.Example, error)
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

func (r *mysqlExample) GetByVar2(ctx context.Context, var2 string, var3 bool) (*src.Example, error) {
	var example src.Example
	err := r.db.GetContext(ctx, &example, "SELECT * FROM example "+
		"WHERE var2 = ? AND var3 = ?", var2, var3)

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
		"WHERE var1 = ?", var1)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrExampleNotFound
		}

		return nil, err
	}

	return &example, nil
}

func (r *mysqlExample) GetByVar3(ctx context.Context, var3 bool, var1 int) (*src.Example, error) {
	var example src.Example
	err := r.db.GetContext(ctx, &example, "SELECT * FROM example "+
		"WHERE var3 = ? AND var1 = ?", var3, var1)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrExampleNotFound
		}

		return nil, err
	}

	return &example, nil
}
