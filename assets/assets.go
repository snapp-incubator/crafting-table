package assets

var A Assets

type Assets struct {
	Sqlx Sqlx
}

func init() {
	A.Sqlx = NewSqlx()
}
