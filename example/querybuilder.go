//go:generate crafting-table qb $GOFILE
package cmd

type Int int

// ct: model
type User struct {
	ID   int
	Name string
	Age  Int
}
