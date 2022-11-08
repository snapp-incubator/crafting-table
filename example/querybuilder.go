//go:generate crafting-table qb $GOFILE
package cmd

type Int int

// ct: model
type User struct {
	ID   int
	Name string
	Age  Int
}

func some() {
	UserQueryBuilder().
		OrderByAsc(UserColumns.ID).
		Select(UserColumns.Name).
		First(db)
}
