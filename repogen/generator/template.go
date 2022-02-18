package generator

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

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
		strings.Replace(strcase.ToSnake(structure.Name), "_", " ", -1),
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

func createFunctionRepository(structure Structure) (syntax string, funcDeclare string, err error) {
	syntax = fmt.Sprintf(
		`
func (r *mysql%s) Create(ctx context.Context, %s *%s.%s) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO %s (" +
`+
			structure.GetDBFields("")+
			") VALUES ("+
			structure.GetDBFields(":")+
			"), %s)"+
			`
	if err != nil {
		return err
	}

	return nil
}
`, structure.Name,
		strcase.ToCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToSnake(structure.Name),
		strcase.ToCamel(structure.Name),
	)

	return syntax, fmt.Sprintf(
		"Create(ctx context.Context, %s *%s.%s) error {",
		strcase.ToCamel(structure.Name), structure.PackageName, structure.Name), nil
}
