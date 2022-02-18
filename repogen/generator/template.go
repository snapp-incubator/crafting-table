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

func createFunctionRepository(structure *Structure) (syntax, funcDeclare string, err error) {
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

func getConditions(v []string, structure *Structure) string {
	var conditions []string

	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res != "" {
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
		if res != "" {
			panic(fmt.Sprintf("%s is not valid db_field", value))
		}
	}

	res := ""
	for _, value := range v {
		res += fmt.Sprintf("%s %s, ",
			strcase.ToCamel(structure.FieldDBNameToName[value]),
			structure.FieldNameToType[structure.FieldDBNameToName[value]])
	}
	return res[:len(res)-2]
}

func getUpdateVariables(v []string, structure *Structure) string {
	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res != "" {
			panic(fmt.Sprintf("%s is not valid db_field", value))
		}
	}

	res := ""
	for _, value := range v {
		res += fmt.Sprintf("%s, ",
			strcase.ToCamel(structure.FieldDBNameToName[value]))
	}

	return res[:len(res)-2]
}

func getFunctionCreator(structure *Structure, vars *[]Variables) (syntax string, fucntions []string, err error) {

	syntax += fmt.Sprintf(
		`
func (r *mysql%s) Get%ss(ctx context.Context) (*[]%s.%s, error) {
	var %s []%s.%s
	err := r.db.SelectContext(ctx, %s, "SELECT * from %s")
	if err != nil {
		return nil, err
	}

	return &%s, nil
}
`,
		structure.Name,
		structure.Name,
		structure.PackageName,
		structure.Name,
		strcase.ToCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToCamel(structure.Name),
		strcase.ToSnake(structure.Name),
		strcase.ToCamel(structure.Name),
	)

	fucntions = append(fucntions,
		fmt.Sprintf("Get%ss(ctx context.Context) (*[]%s.%s, error)",
			structure.Name, structure.PackageName, structure.Name))

	for _, v := range *vars {
		syntax += fmt.Sprintf(
			`
func (r *mysql%s) GetBy%s(ctx context.Context,`+getFunctionVars(v.Name, structure)+`) (*[]%s.%s, error) {
	var %s %s.%s
	err := r.db.GetContext(ctx, %s, "SELECT * FROM %s "
`+getConditions(v.Name, structure)+`, `+getUpdateVariables(v.Name, structure)+`) 
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
			structure.FieldDBNameToName[v.Name[0]],
			structure.PackageName,
			structure.Name,
			strcase.ToCamel(structure.Name),
			structure.PackageName,
			structure.Name,
			strcase.ToCamel(structure.Name),
			strcase.ToSnake(structure.Name),
			structure.Name,
			strcase.ToCamel(structure.Name),
		)

		fucntions = append(fucntions,
			fmt.Sprintf("GetBy%s(ctx context.Context, "+getFunctionVars(v.Name, structure)+") (*[]%s.%s, error)",
				structure.FieldDBNameToName[v.Name[0]],
				getFunctionVars(v.Name, structure),
				structure.PackageName,
				structure.Name))

	}
	return syntax, functions, nil
}
