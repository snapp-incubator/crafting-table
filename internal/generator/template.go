package generator

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"
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
	syntax = fmt.Sprintf(
		`
func (r *mysql%s) Create(ctx context.Context, %s *%s.%s) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO %s (" +
`,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToSnake(structure.Name),
	)

	fields := structure.GetDBFields(":")

	syntax += "\t" + structure.GetDBFields("") + "+\n\t\") VALUES (\"+\n\t" +
		fields[:len(fields)-1] + ")\",\n\t" +
		fmt.Sprintf(
			` %s)

	if err != nil {
		return err
	}

	return nil
}
`, strcase.ToLowerCamel(structure.Name))

	return syntax, fmt.Sprintf(
		"Create(ctx context.Context, %s *%s.%s) error",
		strcase.ToLowerCamel(structure.Name), structure.PackageName, structure.Name), nil
}

func getFunctionCreator(structure *Structure, vars *[]Variables) (syntax string, functions []string, err error) {

	syntax += fmt.Sprintf(
		`
func (r *mysql%s) Get%ss(ctx context.Context) (*[]%s.%s, error) {
	var %s []%s.%s
	err := r.db.SelectContext(ctx, &%s, "SELECT * from %s")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, Err%sNotFound
		}

		return nil, err
	}

	return &%s, nil
}
`,
		structure.Name,
		structure.Name,
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		strcase.ToSnake(structure.Name),
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
	)

	functions = append(functions,
		fmt.Sprintf("Get%ss(ctx context.Context) (*[]%s.%s, error)",
			structure.Name, structure.PackageName, structure.Name))

	for _, v := range *vars {
		syntax += fmt.Sprintf(
			`
func (r *mysql%s) GetBy%s(ctx context.Context, `,
			structure.Name,
			structure.FieldDBNameToName[v.Name[0]],
		) +
			getFunctionVars(v.Name, structure) +
			fmt.Sprintf(`) (*%s.%s, error) {
	var %s %s.%s
	err := r.db.GetContext(ctx, &%s, "SELECT * FROM %s " +
`,
				structure.PackageName,
				structure.Name,
				strcase.ToLowerCamel(structure.Name),
				structure.PackageName,
				structure.Name,
				strcase.ToLowerCamel(structure.Name),
				strcase.ToSnake(structure.Name),
			) +
			"\t\"" +
			getConditions(v.Name, structure) +
			"" +
			"\", " +
			getUpdateVariables(v.Name, structure) +
			")" + fmt.Sprintf(`

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, Err%sNotFound
		}

		return nil, err
	}

	return &%s, nil
}
`,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
		)

		functions = append(functions,
			fmt.Sprintf("GetBy%s(ctx context.Context, ",
				structure.FieldDBNameToName[v.Name[0]],
			)+
				getFunctionVars(v.Name, structure)+
				fmt.Sprintf(") (*%s.%s, error)",
					structure.PackageName,
					structure.Name),
		)

	}
	return syntax, functions, nil
}

func updateFunctionCreator(structure *Structure, updateVars *[]UpdateVariables) (syntax string, functions []string, err error) {
	syntax += fmt.Sprintf(`
func (r *mysql%s) Update(ctx context.Context, %s %s, %s %s.%s) (int64, error) {
	%s.%s = %s

	result, err := r.db.NamedExecContext(ctx, "UPDATE %s "+
		"SET "`,
		structure.Name,
		strcase.ToLowerCamel(structure.Fields[0].Name),
		structure.Fields[0].Type,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		structure.Fields[0].Name,
		strcase.ToLowerCamel(structure.Fields[0].Name),
		strcase.ToSnake(structure.Name),
	) + "\n" + setKeys(structure.Fields) + "\n" +
		getConditions([]string{
			structure.FieldNameToDBName[structure.Fields[0].Name],
		}, structure) + `
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}`

	functions = append(functions,
		fmt.Sprintf("Update(ctx context.Context, %s %s, %s %s.%s) (int64, error)",
			strcase.ToLowerCamel(structure.Fields[0].Name),
			structure.Fields[0].Type,
			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name),
	)

	for _, vars := range *updateVars {

		syntax += fmt.Sprintf(`
func (r *mysql%s) Update%s`,
			structure.Name,
			structure.FieldDBNameToName[vars.Fields[0]],
		) + "(ctx context.Context, " +
			getFunctionVars(vars.By, structure) +
			getFunctionVars(vars.Fields, structure) + ") (int64, error) {\n" +
			`
	query := "UPDATE cancellation_events SET ` + setKeysWithQuestion(vars.Fields) +
			getConditions(vars.By, structure) + `;"
	result, err := r.db.ExecContext(ctx, query, ` + execContextVariables(vars, structure) + `)
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}`

		functions = append(functions,
			fmt.Sprintf("Update%s",
				structure.FieldDBNameToName[vars.Fields[0]],
			)+"(ctx context.Context, "+
				getFunctionVars(vars.By, structure)+
				getFunctionVars(vars.Fields, structure)+") (int64, error)",
		)
	}

	return syntax, functions, err
}
