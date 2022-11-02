//go:generate sqlgen $GOFILE
package main

// sqlgen: table=users
type User struct {
	Id   int
	Name string
	Age  int
}
