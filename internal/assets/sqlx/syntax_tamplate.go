package sqlx

import (
	"errors"
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
	SelectBy(structure *structure.Structure, vars *[]structure.GetVariable) (syntax string, header []string)
	Join(structure *structure.Structure, joinVariables *structure.JoinVariables) (syntax string, header string)
	Aggregate(structure *structure.Structure, vars *[]structure.AggregateField) (syntax string, signatures []string)
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
	joinFuncBody           string
	joinFuncSignature      string
	aggregateFuncBody      string
	aggregateFuncSignature string
}

func NewSqlx() Sqlx {
	s := sqlx{}

	s.insertFuncBody = `
func (r *mysql%s) Insert(ctx context.Context, %s *%s.%s) error {
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

	s.insertFuncSignature = `Insert(ctx context.Context, %s *%s.%s) error`

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
		"SET "+
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

	s.joinFuncBody = `
func (r *mysql%s) GetJoined%s(ctx context.Context, limit uint) ([]%s.%s, error) {
	query := "SELECT " +
		%s
		"FROM %s AS %s " +
		%s +
		"LIMIT ?"

	var %s []%s.%s
	err := r.db.SelectContext(ctx, &%s, query, limit)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, %s
		}

		return nil, err
	}

	return %s, nil
}
`

	s.joinFuncSignature = `GetJoined%s(ctx context.Context, limit uint) ([]%s.%s, error)`

	s.aggregateFuncBody = `
func (r *mysql%s) GetAggregateBy%s(ctx context.Context, %s) (*int, error) {
	var res struct{` +
		"result int `db:\"%s\"`" + `
	}

	err := r.db.SelectContext(ctx, &res, "SELECT %s FROM %s " +
		"%s%s",
		 %s,
	)

	if err != nil {
		return nil, err
	}

	return &res.result, nil
}
`

	s.aggregateFuncSignature = `GetAggregateBy%s(ctx context.Context, %s) (*int, error)`

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
		structure.TableName,
		contextKeys(structure.Fields),
		conditions([]string{
			structure.FieldMapNameToDBFlag[structure.Fields[0].Name],
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
			functionNameList = append(functionNameList, structure.FieldMapDBFlagToName[name])
		}
		functionName := strings.Join(functionNameList, "And")

		syntax += fmt.Sprintf(
			s.updateFuncBody,
			structure.Name,
			functionName,
			inputFunctionVariables(v.Conditions, structure),
			inputFunctionVariables(v.Fields, structure),
			structure.TableName,
			contextKeys(v.Fields),
			conditions(v.Conditions, structure, true),
			execContextVariables(v, structure, false),
		)

		signatures = append(
			signatures,
			fmt.Sprintf(
				s.updateFuncSignature,
				functionName,
				inputFunctionVariables(v.Conditions, structure),
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
		strcase.ToLowerCamel(structure.TableName),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		structure.TableName,
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

func (s sqlx) SelectBy(structure *structure.Structure, vars *[]structure.GetVariable) (syntax string, signatures []string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, condition := range v.Conditions {
			functionNameList = append(functionNameList, structure.FieldMapDBFlagToName[condition])
		}
		functionName := strings.Join(functionNameList, "And")

		syntax += fmt.Sprintf(
			s.selectFuncBody,
			structure.Name,
			functionName,
			inputFunctionVariables(v.Conditions, structure),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
			structure.TableName,
			conditions(v.Conditions, structure, true),
			contextVariables(v.Conditions, structure),
			structure.Name,
			strcase.ToLowerCamel(structure.Name),
		)

		signatures = append(
			signatures,
			fmt.Sprintf(
				s.selectFuncSignature,
				functionName,
				inputFunctionVariables(v.Conditions, structure),
				structure.PackageName,
				structure.Name,
			),
		)
	}

	return syntax, signatures
}

func (s sqlx) Join(structure *structure.Structure, joinVariables *structure.JoinVariables) (syntax string, header string) {

	syntax = fmt.Sprintf(
		s.joinFuncBody,
		structure.Name,
		structure.Name,
		structure.PackageName,
		structure.Name,

		joinField(structure, joinVariables),

		structure.TableName,
		string(structure.TableName[0]),

		joins(structure, joinVariables),

		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Name,
		strcase.ToLowerCamel(structure.Name),

		"Err"+structure.Name+"NotFound",

		strcase.ToLowerCamel(structure.Name),
	)

	header = fmt.Sprintf(
		s.joinFuncSignature,
		structure.Name,
		structure.PackageName,
		structure.Name,
	)

	return syntax, header
}

func (s sqlx) Aggregate(structure *structure.Structure, vars *[]structure.AggregateField) (syntax string, signatures []string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, condition := range v.Conditions {
			functionNameList = append(functionNameList, structure.FieldMapDBFlagToName[condition])
		}
		functionName := strings.Join(functionNameList, "And")

		syntax += fmt.Sprintf(
			s.aggregateFuncBody,
			structure.Name,
			functionName,
			inputFunctionVariables(v.Conditions, structure),
			v.As,
			aggregateSyntax(v),
			structure.TableName,
			conditions(v.Conditions, structure, true),
			groupBy(v),
			contextVariables(v.Conditions, structure),
		)

		signatures = append(
			signatures,
			fmt.Sprintf(
				s.aggregateFuncSignature,
				functionName,
				inputFunctionVariables(v.Conditions, structure),
			),
		)
	}

	return syntax, signatures
}

func groupBy(v structure.AggregateField) string {
	if len(v.GroupBy) == 0 {
		return ""
	}

	return " GROUP BY " + strings.Join(v.GroupBy, ", ")
}

func aggregateSyntax(v structure.AggregateField) string {
	function, ok := structure.AggregateMap[strings.ToLower(v.Function)]
	if !ok {
		panic("aggregate function not found")
	}

	return function + "(" + v.On + ")" + " AS " + v.As
}

func joinField(s *structure.Structure, joinVariables *structure.JoinVariables) string {
	fields := ""

	// add first struct fields
	firstCharOfStruct := string(s.TableName[0])
	for dbName, _ := range s.FieldMapDBFlagToName {
		fields += fmt.Sprintf("\"%s.%s AS %s, \" + \n\t\t", firstCharOfStruct, dbName, dbName)
	}

	// add join fields
	for _, joinVariable := range joinVariables.Fields {
		source := strings.Replace(joinVariable.JoinStructPath, " ", "", -1)

		ss, err := structure.BindStruct(source, joinVariable.JoinStructName)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error in bindStruct: %s", err.Error()))
			panic(err)
		}

		firstCharOfStruct := string(ss.TableName[0])
		for dbName, _ := range ss.FieldMapDBFlagToName {
			fields += fmt.Sprintf("\"%s.%s AS \\\"%s.%s\\\", \" + \n\t\t ",
				firstCharOfStruct, dbName, joinVariable.JoinFieldAs, dbName)
		}

		fields = strings.TrimSuffix(fields, ", \" + \n\t\t ") + " \" +"
	}

	return fields
}

func joins(structure *structure.Structure, vars *structure.JoinVariables) string {
	var syntax string
	for _, v := range vars.Fields {
		syntax += fmt.Sprintf(
			"\"%s JOIN %s AS %s ON %s = %s \"",
			v.JoinType,
			strcase.ToSnake(v.JoinStructName),
			string(strcase.ToSnake(v.JoinStructName)[0]),
			string(structure.TableName[0])+"."+v.JoinOn,
			string(strcase.ToSnake(v.JoinStructName)[0])+"."+v.JoinOn,
		)
	}
	return syntax
}

func contextKeys[T FieldType](fields []T) string {
	result := "\""
	tmp := ""
	for _, field := range fields {
		if len(tmp) > 80 {
			result += tmp[:len(tmp)-2] + ", \"+\n\t\t\""
			tmp = ""
		}

		switch reflect.TypeOf(field).String() {
		case "structure.Field":
			tmp += interface{}(field).(structure.Field).DBFlag +
				" = :" +
				interface{}(field).(structure.Field).DBFlag + ", "
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
		res := structure.FieldMapDBFlagToName[value]
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
		res := structure.FieldMapDBFlagToName[value]
		if res == "" {
			log.Fatalf("%s is not valid db_field", value)
		}
	}

	res := ""
	for _, value := range v {
		res += fmt.Sprintf("%s %s, ",
			strcase.ToLowerCamel(structure.FieldMapDBFlagToName[value]),
			structure.FieldMapNameToType[structure.FieldMapDBFlagToName[value]])
	}

	return res[:len(res)-2]
}

func contextVariables(v []string, structure *structure.Structure) string {
	for _, value := range v {
		res := structure.FieldMapDBFlagToName[value]
		if res == "" {
			log.Fatalf("%s is not valid db_field", value)
		}
	}

	res := ""
	for _, value := range v {
		res += fmt.Sprintf("%s, ",
			strcase.ToLowerCamel(structure.FieldMapDBFlagToName[value]))
	}

	return res[:len(res)-2]
}

func execContextVariables(vars structure.UpdateVariables, structure *structure.Structure, reverse bool) string {
	if reverse {
		return contextVariables(vars.Conditions, structure) + ", " + contextVariables(vars.Fields, structure)
	}
	return contextVariables(vars.Fields, structure) + ", " + contextVariables(vars.Conditions, structure)
}
