//go:generate crafting-table qb $GOFILE
package cmd

// ct: model
type User struct {
	ID   int
	Name string
	Age  int
}
