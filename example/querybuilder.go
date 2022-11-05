//go:generate crafting-table gogen $GOFILE
package cmd

// ct: model
type User struct {
	ID   int
	Name string
	Age  int
}

func some() {
}
