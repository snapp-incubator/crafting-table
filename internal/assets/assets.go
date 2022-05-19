package assets

import (
	"github.com/snapp-incubator/crafting-table/internal/assets/sqlx"
)

var A Assets

type Assets struct {
	Sqlx sqlx.Sqlx
}

func init() {
	A.Sqlx = sqlx.NewSqlx()
}
