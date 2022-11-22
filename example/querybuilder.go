//go:generate crafting-table qb --dialect=mysql -f $GOFILE
package cmd

// ct: model
type User struct {
	ID   int
	Name string
	Age  int
}
