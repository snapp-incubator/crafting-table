package assets

import (
	"fmt"
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/n25a/repogen/internal/generator"
)

type Sqlx interface {
	Insert(structure *generator.Structure) (syntax string, header string)
	UpdateAll(structure *generator.Structure) (syntax string, header string)
	UpdateBy(structure *generator.Structure, vars *[]generator.UpdateVariables) (syntax string, header []string)
	SelectAll(structure *generator.Structure) (syntax string, header string)
	SelectBy(structure *generator.Structure, vars *[]generator.Variables) (syntax string, header []string)
}

type sqlx struct {
	insertFuncHeader    string
	insertFuncBody      string
	selectAllFuncHeader string
	selectAllFuncBody   string
	selectFuncHeader    string
	selectFuncBody      string
	updateAllFuncHeader string
	updateAllFuncBody   string
	updateFuncHeader    string
	updateFuncBody      string
}

func NewSqlx() Sqlx {
	s := sqlx{}

	s.insertFuncBody = `
func (r *mysql%s) Create(ctx context.Context, %s *%s.%s) error {
	_, err := r.db.NamedExecContext(ctx, "INSERT INTO %s (" +
	%s 
	") VALUES (" +
	%s
	)",
	
	if err != nil {
		return err
	}

	return nil
}
`

	s.insertFuncHeader = `Create(ctx context.Context, %s *%s.%s) error`

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

	s.selectAllFuncHeader = `Get%ss(ctx context.Context) (*[]%s.%s, error)`

	s.selectAllFuncBody = `
func (r *mysql%s) GetBy%s(ctx context.Context, %s) (*%s.%s, error) {
	var %s %s.%s
	err := r.db.GetContext(ctx, &%s, "SELECT * FROM %s " +
	"%s",
	%s,
	")",

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, Err%sNotFound
		}

		return nil, err
	}

	return &%s, nil
}
`

	s.selectAllFuncHeader = `GetBy%s(ctx context.Context, %s) (*%s.%s, error)`

	s.updateAllFuncBody = `
func (r *mysql%s) Update(ctx context.Context, %s %s, %s %s.%s) (int64, error) {
	%s.%s = %s

	result, err := r.db.NamedExecContext(ctx, "UPDATE %s "+
		"SET"+
		"%s" +
		"WHERE %s = :%s",
		%s,
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
`

	s.updateAllFuncHeader = `Update(ctx context.Context, %s %s, %s %s.%s) (int64, error)`

	s.updateFuncBody = `
func (r *mysql%s) Update%s(ctx context.Context, %s, %s) (int64, error) {
	query := "UPDATE %s SET %s WHERE %s;" 

	result, err := r.db.ExecContext(ctx, query, %s)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
`

	s.updateFuncHeader = `Update%s(ctx context.Context, %s, %s) (int64, error)`

	return s
}

func (s sqlx) Insert(structure *generator.Structure) (syntax string, header string) {
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

	header = fmt.Sprintf(
		s.insertFuncHeader,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
	)

	return syntax, header
}

func (s sqlx) UpdateAll(structure *generator.Structure) (syntax string, header string) {
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
		strcase.ToSnake(structure.Name),
		SetKeys(structure.Fields),
		GetConditions([]string{
			structure.FieldNameToDBName[structure.Fields[0].Name],
		}, structure),
		strcase.ToLowerCamel(structure.Name),
	)

	header = fmt.Sprintf(
		s.updateAllFuncHeader,
		strcase.ToLowerCamel(structure.Fields[0].Name),
		structure.Fields[0].Type,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
	)

	return syntax, header
}

func (s sqlx) UpdateBy(structure *generator.Structure, vars *[]generator.UpdateVariables) (syntax string, headers []string) {

	for _, v := range *vars {
		syntax += fmt.Sprintf(
			s.updateFuncBody,
			structure.Name,
			structure.FieldDBNameToName[v.Fields[0]],
			GetFunctionVars(v.By, structure),
			GetFunctionVars(v.Fields, structure),
			SetKeysWithQuestion(v.Fields),
			GetConditions(v.By, structure),
			ExecContextVariables(v, structure),
		)

		headers = append(
			headers,
			fmt.Sprintf(
				s.updateFuncHeader,
				structure.FieldDBNameToName[v.Fields[0]],
				GetFunctionVars(v.By, structure),
				GetFunctionVars(v.Fields, structure),
			),
		)
	}

	return syntax, headers
}

func (s sqlx) SelectAll(structure *generator.Structure) (syntax string, header string) {

	syntax = fmt.Sprintf(
		s.selectAllFuncBody,
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

	header = fmt.Sprintf(
		s.selectAllFuncHeader,
		structure.Name,
		structure.PackageName,
		structure.Name,
	)

	return syntax, header
}

func (s sqlx) SelectBy(structure *generator.Structure, vars *[]generator.Variables) (syntax string, headers []string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, name := range v.Name {
			functionNameList = append(functionNameList, structure.FieldDBNameToName[name])
		}
		functionName := strings.Join(functionNameList, "And")

		syntax += fmt.Sprintf(
			s.selectFuncBody,
			functionName,
			structure.FieldDBNameToName[v.Name[0]],
			GetFunctionVars(v.Name, structure),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
			strcase.ToSnake(structure.Name),
			GetConditions(v.Name, structure),
			GetUpdateVariables(v.Name, structure),
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
		)

		headers = append(
			headers,
			fmt.Sprintf(
				s.selectFuncHeader,
				functionName,
				GetFunctionVars(v.Name, structure),
				structure.PackageName,
				structure.Name,
			),
		)
	}

	return syntax, headers
}
