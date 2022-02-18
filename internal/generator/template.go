package generator

import (
	"fmt"
	"os"
	"os/exec"
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

func writeFile(content, dst string) error {
	f, err := os.Create(dst)

	if err != nil {
		return err
	}

	defer func(f *os.File) {
		_ = f.Close()
	}(f)

	_, err = f.WriteString(content)

	if err != nil {
		return err
	}

	return nil
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

func getConditions(v []string, structure *Structure) string {
	var conditions []string

	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res == "" {
			panic(fmt.Sprintf("%s is not valid db_field", value))
		}
	}

	for _, value := range v {
		conditions = append(conditions, fmt.Sprintf("%s = ?", value))
	}
	return "WHERE " + strings.Join(conditions, " AND ")
}

func getFunctionVars(v []string, structure *Structure) string {
	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res == "" {
			panic(fmt.Sprintf("%s is not valid db_field", value))
		}
	}

	res := ""
	for _, value := range v {
		res += fmt.Sprintf("%s %s, ",
			strcase.ToLowerCamel(structure.FieldDBNameToName[value]),
			structure.FieldNameToType[structure.FieldDBNameToName[value]])
	}

	return res[:len(res)-2]
}

func getUpdateVariables(v []string, structure *Structure) string {
	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res == "" {
			panic(fmt.Sprintf("%s is not valid db_field", value))
		}
	}

	res := ""
	for _, value := range v {
		res += fmt.Sprintf("%s, ",
			strcase.ToLowerCamel(structure.FieldDBNameToName[value]))
	}

	return res[:len(res)-2]
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

func interfaceSyntaxCreator(structure *Structure, functions []string) string {
	syntax := fmt.Sprintf(
		"type %s interface {",
		structure.Name,
	)

	for _, function := range functions {
		syntax += fmt.Sprintf("\n\t%s", function)
	}
	syntax += "\n}"

	return syntax
}

func linter(dst string) error {
	cmd := exec.Command("gofmt", "-s", "-w", dst)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("goimports", dst)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
