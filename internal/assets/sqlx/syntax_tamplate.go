package sqlx

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

type Sqlx interface {
	Insert(structure *structure.Structure) (syntax string, header string)
	UpdateAll(structure *structure.Structure) (syntax string, header string)
	UpdateBy(structure *structure.Structure, vars *[]structure.UpdateVariables) (syntax string, header []string)
	SelectAll(structure *structure.Structure) (syntax string, header string)
	SelectBy(structure *structure.Structure, vars *[]structure.Variables) (syntax string, header []string)
}

type FieldType interface{ structure.Field | string }

type sqlx struct {
	insertFuncSignature    string
	insertFuncBody         string
	selectAllFuncSignature string
	selectAllFuncBody      string
	selectFuncSignature    string
	selectFuncBody         string
	updateAllFuncSignature string
	updateAllFuncBody      string
	updateFuncSignature    string
	updateFuncBody         string
}

func NewSqlx() Sqlx {
	s := sqlx{}

	s.insertFuncBody = `
func (r *mysql%s) Create(ctx context.Context, %s *%s.%s) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO %s (" +
	%s +
	") VALUES (" +
	%s)", 
	%s)
	
	if err != nil {
		return err
	}

	return nil
}
`

	s.insertFuncSignature = `Create(ctx context.Context, %s *%s.%s) error`

	s.selectAllFuncBody = `
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
`

	s.selectAllFuncSignature = `Get%ss(ctx context.Context) (*[]%s.%s, error)`

	s.selectFuncBody = `
func (r *mysql%s) GetBy%s(ctx context.Context, %s) (*%s.%s, error) {
	var %s %s.%s

	err := r.db.GetContext(ctx, &%s, "SELECT * FROM %s " +
		"%s",
		%s,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, Err%sNotFound
		}

		return nil, err
	}

	return &%s, nil
}
`

	s.selectFuncSignature = `GetBy%s(ctx context.Context, %s) (*%s.%s, error)`

	s.updateAllFuncBody = `
func (r *mysql%s) Update(ctx context.Context, %s %s, %s %s.%s) (int64, error) {
	%s.%s = %s

	result, err := r.db.NamedExecContext(ctx, "UPDATE %s "+
		"SET"+
		%s +
		"%s",
		%s,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
`

	s.updateAllFuncSignature = `Update(ctx context.Context, %s %s, %s %s.%s) (int64, error)`

	s.updateFuncBody = `
func (r *mysql%s) Update%s(ctx context.Context, %s, %s) (int64, error) {
	query := "UPDATE %s SET " +
			%s +
			"%s;" 

	result, err := r.db.ExecContext(ctx, query, %s)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
`

	s.updateFuncSignature = `Update%s(ctx context.Context, %s, %s) (int64, error)`

	return s
}

func (s sqlx) Insert(structure *structure.Structure) (syntax string, signature string) {
	fields := structure.GetDBFields(":")

	syntax = fmt.Sprintf(
		s.insertFuncBody,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToSnake(structure.Name),
		structure.GetDBFields(""),
		fields[:len(fields)-1],
		strcase.ToLowerCamel(structure.Name),
	)

	signature = fmt.Sprintf(
		s.insertFuncSignature,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
	)

	return syntax, signature
}

func (s sqlx) UpdateAll(structure *structure.Structure) (syntax string, signature string) {
	syntax = fmt.Sprintf(
		s.updateAllFuncBody,
		structure.Name,
		strcase.ToLowerCamel(structure.Fields[0].Name),
		structure.Fields[0].Type,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		structure.Fields[0].Name,
		strcase.ToLowerCamel(structure.Fields[0].Name),
		structure.DBName,
		contextKeys(structure.Fields),
		conditions([]string{
			structure.FieldNameToDBName[structure.Fields[0].Name],
		}, structure, false),
		strcase.ToLowerCamel(structure.Name),
	)

	signature = fmt.Sprintf(
		s.updateAllFuncSignature,
		strcase.ToLowerCamel(structure.Fields[0].Name),
		structure.Fields[0].Type,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
	)

	return syntax, signature
}

func (s sqlx) UpdateBy(structure *structure.Structure, vars *[]structure.UpdateVariables) (syntax string, signatures []string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, name := range v.Fields {
			functionNameList = append(functionNameList, structure.FieldDBNameToName[name])
		}
		functionName := strings.Join(functionNameList, "And")

		syntax += fmt.Sprintf(
			s.updateFuncBody,
			structure.Name,
			functionName,
			inputFunctionVariables(v.By, structure),
			inputFunctionVariables(v.Fields, structure),
			structure.DBName,
			contextKeys(v.Fields),
			conditions(v.By, structure, true),
			execContextVariables(v, structure),
		)

		signatures = append(
			signatures,
			fmt.Sprintf(
				s.updateFuncSignature,
				functionName,
				inputFunctionVariables(v.By, structure),
				inputFunctionVariables(v.Fields, structure),
			),
		)
	}

	return syntax, signatures
}

func (s sqlx) SelectAll(structure *structure.Structure) (syntax string, signature string) {

	syntax = fmt.Sprintf(
		s.selectAllFuncBody,
		structure.Name,
		structure.Name,
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.DBName),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		structure.DBName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
	)

	signature = fmt.Sprintf(
		s.selectAllFuncSignature,
		structure.Name,
		structure.PackageName,
		structure.Name,
	)

	return syntax, signature
}

func (s sqlx) SelectBy(structure *structure.Structure, vars *[]structure.Variables) (syntax string, signatures []string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, name := range v.Name {
			functionNameList = append(functionNameList, structure.FieldDBNameToName[name])
		}
		functionName := strings.Join(functionNameList, "And")

		syntax += fmt.Sprintf(
			s.selectFuncBody,
			structure.Name,
			functionName,
			inputFunctionVariables(v.Name, structure),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
			structure.DBName,
			conditions(v.Name, structure, true),
			contextVariables(v.Name, structure),
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
		)

		signatures = append(
			signatures,
			fmt.Sprintf(
				s.selectFuncSignature,
				functionName,
				inputFunctionVariables(v.Name, structure),
				structure.PackageName,
				structure.Name,
			),
		)
	}

	return syntax, signatures
}

func contextKeys[T FieldType](fields []T) string {
	result := "\""
	tmp := ""
	for _, field := range fields {
		if len(tmp) > 80 {
			result += tmp[:len(tmp)-2] + "\"+\n\t\t\""
			tmp = ""
		}

		switch reflect.TypeOf(field).String() {
		case "structure.Field":
			tmp += interface{}(field).(structure.Field).DBName +
				" = :" +
				interface{}(field).(structure.Field).DBName + ", "
		case "string":
			tmp += interface{}(field).(string) + " = ?, "
		}
	}

	if tmp != "" {
		result += tmp[:len(tmp)-2] + " \""
	}

	return result
}

func conditions(v []string, structure *structure.Structure, withQuestionMark bool) string {
	var conditions []string

	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res == "" {
			log.Fatalf("%s is not valid db_field", value)
		}
	}

	for _, value := range v {
		if withQuestionMark {
			conditions = append(conditions, fmt.Sprintf("%s = ?", value))
		} else {
			conditions = append(conditions, fmt.Sprintf("%s = :%s", value, value))
		}
	}
	return "WHERE " + strings.Join(conditions, " AND ")
}

func inputFunctionVariables(v []string, structure *structure.Structure) string {
	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res == "" {
			log.Fatalf("%s is not valid db_field", value)
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

func contextVariables(v []string, structure *structure.Structure) string {
	for _, value := range v {
		res := structure.FieldDBNameToName[value]
		if res == "" {
			log.Fatalf("%s is not valid db_field", value)
		}
	}

	res := ""
	for _, value := range v {
		res += fmt.Sprintf("%s, ",
			strcase.ToLowerCamel(structure.FieldDBNameToName[value]))
	}

	return res[:len(res)-2]
}

func execContextVariables(vars structure.UpdateVariables, structure *structure.Structure) string {
	return contextVariables(vars.Fields, structure) + ", " + contextVariables(vars.By, structure)
}
