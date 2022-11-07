package repository

import (
	"fmt"
	"strings"

	"github.com/snapp-incubator/crafting-table/internal/assets"

	"github.com/iancoleman/strcase"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

func createTemplate(structure *structure.Structure, packageName, interfaceSyntax,
	createFunc, updateFunc, getFunc, joinFunc string) string {
	syntax := fmt.Sprintf(`
// Code generated by Crafting-Table.
// Source code: https://github.com/snapp-incubator/crafting-table

package %s

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

%s

var Err%sNotFound = errors.New("%s not found")

type mysql%s struct {
	db *sqlx.DB
}

func NewMySQL%s(db *sqlx.DB) %s {
	return &mysql%s{db: db}
}

%s

%s

%s

%s
`,
		packageName,
		interfaceSyntax,
		structure.Name,
		strings.Replace(strcase.ToSnake(structure.Name), "_", " ", -1),
		structure.Name,
		structure.Name,
		structure.Name,
		structure.Name,
		createFunc,
		updateFunc,
		getFunc,
		joinFunc,
	)

	return syntax
}

func createTestTemplate(structure *structure.Structure, packageName,
	createTestSyntax, updateTestSyntax, getTestSyntax, joinTestSyntax, aggregateSyntax string) string {

	if createTestSyntax != "" {
		createTestSyntax = fmt.Sprintf(`
//-----------------------------------------------------------------
//							INSERT
//-----------------------------------------------------------------

%s
`,
			createTestSyntax,
		)
	}

	if updateTestSyntax != "" {
		updateTestSyntax = fmt.Sprintf(`
//-----------------------------------------------------------------
//							UPDATE
//-----------------------------------------------------------------

%s
`,
			updateTestSyntax,
		)
	}

	if getTestSyntax != "" {
		getTestSyntax = fmt.Sprintf(`
//-----------------------------------------------------------------
//							GET
//-----------------------------------------------------------------

%s
`,
			getTestSyntax,
		)
	}

	if joinTestSyntax != "" {
		joinTestSyntax = fmt.Sprintf(`
//-----------------------------------------------------------------
//							JOIN
//-----------------------------------------------------------------

%s
`,
			joinTestSyntax,
		)
	}

	if aggregateSyntax != "" {
		aggregateSyntax = fmt.Sprintf(`
//-----------------------------------------------------------------
//							AGGREGATE
//-----------------------------------------------------------------

%s
`,
			aggregateSyntax,
		)
	}

	return fmt.Sprintf(`
// Code generated by Crafting-Table.
// Source code: https://github.com/snapp-incubator/crafting-table

package %s

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/suite"
)

type %sRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	mock sqlmock.Sqlmock
	repo *mysql%s
}

func (suite *%sRepositoryTestSuite) SetupTest() {
	require := suite.Require()
	var err error

	suite.db, suite.mock, err = sqlmock.New()
	require.NoError(err)

	suite.repo = &mysql%s{
		db: sqlx.NewDb(suite.db, "mysql"),
	}
}

%s

%s

%s

%s

%s

//-----------------------------------------------------------------
//						 RUN ALL TESTS
//-----------------------------------------------------------------

func Test%sRepository(t *testing.T) {
	suite.Run(t, new(%sRepositoryTestSuite))
}
		`,
		packageName,
		structure.Name,
		structure.Name,
		structure.Name,
		structure.Name,
		createTestSyntax,
		updateTestSyntax,
		getTestSyntax,
		joinTestSyntax,
		aggregateSyntax,
		structure.Name,
		structure.Name,
	)
}

func createFunction(structure *structure.Structure) (syntax, funcDeclare string, err error) {
	syntax, funcDeclare = assets.A.Sqlx.Insert(structure)
	return syntax, funcDeclare, nil
}

func createTestFunction(structure *structure.Structure) (syntax string) {
	return assets.A.SqlxTest.Insert(structure)
}

func getFunction(structure *structure.Structure, vars *[]structure.GetVariable) (syntax string, functions []string, err error) {

	body, signature := assets.A.Sqlx.SelectAll(structure)

	syntax += body
	functions = append(functions, signature)

	body, signatures := assets.A.Sqlx.SelectBy(structure, vars)

	syntax += body
	for _, signature := range signatures {
		functions = append(functions, signature)
	}

	return syntax, functions, nil
}

func getTestFunction(structure *structure.Structure, vars *[]structure.GetVariable) (syntax string) {
	//syntax = assets.A.SqlxTest.SelectAll(structure)
	body := assets.A.SqlxTest.SelectBy(structure, vars)

	syntax += "\n" + body

	return syntax
}

func updateFunction(structure *structure.Structure, updateVars *[]structure.UpdateVariables) (syntax string, functions []string, err error) {

	body, signature := assets.A.Sqlx.UpdateAll(structure)

	syntax += body
	functions = append(functions, signature)

	body, signatures := assets.A.Sqlx.UpdateBy(structure, updateVars)

	syntax += body
	for _, signature := range signatures {
		functions = append(functions, signature)
	}

	return syntax, functions, nil
}

func updateTestFunction(structure *structure.Structure, updateVars *[]structure.UpdateVariables) (syntax string) {
	syntax = assets.A.SqlxTest.UpdateAll(structure)
	body := assets.A.SqlxTest.UpdateBy(structure, updateVars)

	syntax += "\n" + body

	return syntax
}

func joinFunction(structure *structure.Structure,
	vars *[]structure.JoinVariables) (syntax string, signatures []string, err error) {

	for _, joinVariables := range *vars {
		body, signature := assets.A.Sqlx.Join(structure, &joinVariables)
		syntax += body + "\n"
		signatures = append(signatures, signature)
	}

	return syntax, signatures, nil
}

func joinTestFunction(structure *structure.Structure,
	vars *[]structure.JoinVariables) (syntax string) {

	for _, joinVariables := range *vars {
		body := assets.A.SqlxTest.Join(structure, &joinVariables)

		syntax += "\n" + body
	}

	return syntax
}

func aggregateFunction(structure *structure.Structure,
	vars *[]structure.AggregateField) (syntax string, functions []string, err error) {

	body, signatures := assets.A.Sqlx.Aggregate(structure, vars)
	syntax += body

	for _, signature := range signatures {
		functions = append(functions, signature)
	}

	return syntax, functions, nil
}
