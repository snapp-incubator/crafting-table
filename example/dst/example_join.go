// Code generated by Crafting-Table.
// Source code: https://github.com/snapp-incubator/crafting-table

package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/snapp-incubator/crafting-table/example/src"
)

type JoinExample interface {
	GetJoinedJoinExample(ctx context.Context, limit uint) ([]src.JoinExample, error)
}

var ErrJoinExampleNotFound = errors.New("join example not found")

type mysqlJoinExample struct {
	db *sqlx.DB
}

func NewMySQLJoinExample(db *sqlx.DB) JoinExample {
	return &mysqlJoinExample{db: db}
}

func (r *mysqlJoinExample) GetJoinedJoinExample(ctx context.Context, limit uint) ([]src.JoinExample, error) {
	query := "SELECT " +
		"j.var1 AS var1, " +
		"j.var9 AS var9, " +
		"j.var10 AS var10, " +
		"j.var11 AS var11, " +
		"j.var12 AS var12, " +
		"j.var13 AS var13, " +
		"e.var1 AS \"var13.var1\", " +
		"e.var2 AS \"var13.var2\", " +
		"e.var3 AS \"var13.var3\", " +
		"e.var4 AS \"var13.var4\" " +
		"FROM join_example AS j " +
		"LEFT OUTER JOIN example AS e ON j.var1 = e.var1 " +
		"LIMIT ?"

	var joinExample []src.JoinExample
	err := r.db.SelectContext(ctx, &joinExample, query, limit)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrJoinExampleNotFound
		}

		return nil, err
	}

	return joinExample, nil
}
