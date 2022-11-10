//go:generate crafting-table qb --dialect=mysql $GOFILE
package cmd

// ct: model
type User struct {
	ID   int
	Name string
	Age  int
}
