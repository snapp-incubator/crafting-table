package generator

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/n25a/repogen/assets"
)

func createTemplate(structure *Structure, packageName, interfaceSyntax,
	createSyntax, updateSyntax, getSyntax string) string {
	syntax := fmt.Sprintf(`
package %s

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
)

`, packageName) + interfaceSyntax + fmt.Sprintf(`

var Err%sNotFound = errors.New("%s not found")

type mysql%s struct {
	db *sqlx.DB
}

func NewMySQL%s(db *sqlx.DB) %s {
	return &mysql%s{db: db}
}

`,
		structure.Name,
		strings.Replace(strcase.ToSnake(structure.Name), "_", " ", -1),
		structure.Name,
		structure.Name,
		structure.Name,
		structure.Name,
	)

	syntax += createSyntax + "\n" + updateSyntax + "\n" + getSyntax + "\n"

	return syntax
}

func createFunctionRepository(structure *Structure) (syntax, funcDeclare string, err error) {
	syntax, funcDeclare = assets.A.Sqlx.UpdateAll(structure)
	return syntax, funcDeclare, nil
}

func getFunctionCreator(structure *Structure, vars *[]Variables) (syntax string, functions []string, err error) {

	body, header := assets.A.Sqlx.SelectAll(structure)

	syntax += body
	functions = append(functions, header)

	body, headers := assets.A.Sqlx.SelectBy(structure, vars)

	syntax += body
	for _, header := range headers {
		functions = append(functions, header)
	}

	return syntax, functions, nil
}

func updateFunctionCreator(structure *Structure, updateVars *[]UpdateVariables) (syntax string, functions []string, err error) {

	body, header := assets.A.Sqlx.UpdateAll(structure)

	syntax += body
	functions = append(functions, header)

	body, headers := assets.A.Sqlx.UpdateBy(structure, updateVars)

	syntax += body
	for _, header := range headers {
		functions = append(functions, header)
	}

	return syntax, functions, nil
}
