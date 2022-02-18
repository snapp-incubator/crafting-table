package generator

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func createTemplate(filename, packageName, interfaceSyntax string, structure Structure) string {
	return fmt.Sprintf(
		`package %s
import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

`+interfaceSyntax+
			`

var %sNotFoundErr = errors.New("%s event not found")

type mysql%s struct {
	db *sqlx.DB
}

func NewMySQL%s(db *sqlx.DB) %s {
	return &mysqlCancellationEvent{db: db}
}

`,
		structure.PackageName,
		structure.Name,
		strings.Replace(toSnakeCase(structure.Name), "_", " ", -1),
		structure.Name,
		structure.Name,
		structure.Name,
	)
}

func writeFile(content string, dst string) error {
	f, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)

	if err != nil {
		return err
	}

	return nil
}
