package assets

import "github.com/n25a/crafting-table/assets/sqlx"

var A Assets

type Assets struct {
	Sqlx sqlx.Sqlx
}

func init() {
	A.Sqlx = sqlx.NewSqlx()
}