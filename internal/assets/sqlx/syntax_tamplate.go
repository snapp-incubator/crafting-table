package sqlx

import (
	"errors"
	"fmt"
	"log"
	"reflect"
	"strings"
	"text/template"

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
	insertFuncSignatureTemplate    *template.Template
	insertFuncBodyTemplate         *template.Template
	selectAllFuncSignatureTemplate *template.Template
	selectAllFuncBodyTemplate      *template.Template
	selectFuncSignatureTemplate    *template.Template
	selectFuncBodyTemplate         *template.Template
	updateAllFuncSignatureTemplate *template.Template
	updateAllFuncBodyTemplate      *template.Template
	updateFuncSignatureTemplate    *template.Template
	updateFuncBodyTemplate         *template.Template
	joinFuncBodyTemplate           *template.Template
	joinFuncSignatureTemplate      *template.Template
	aggregateFuncBodyTemplate      *template.Template
	aggregateFuncSignatureTemplate *template.Template
}

func NewSqlx() Sqlx {
	s := sqlx{}

	s.insertFuncBodyTemplate = template.Must(template.New("insert").Parse(`
		func (r *mysql{{.Name}}) Insert(ctx context.Context, {{.NameLowerCamel}} *{{.PackageName}}.{{.Name}}) error {
			_, err := r.db.NamedExecContext(ctx, "INSERT INTO {{.NameSnake}} (" +
			{{.Fields}} +
			") VALUES (" +
			{{.FieldsWithColon}})", 
			{{.NameLowerCamel}})
			
			if err != nil {
				return err
			}
		
			return nil
		}
		`))
	s.insertFuncSignatureTemplate = template.Must(template.New("insert-signiture").Parse(`Insert(ctx context.Context, {{.NameLowerCamel}} *{{.PackageName}}.{{.Name}}) error`))

	s.selectAllFuncBodyTemplate = template.Must(template.New("select-all").Parse(`
func (r *mysql{{.Name}}) Get{{.Name}}s(ctx context.Context) (*[]{{.PackageName}}.{{.Name}}, error) {
	var {{.TableNameLowerCamel}} []{{.PackageName}}.{{.Name}}
	err := r.db.SelectContext(ctx, &{{.NameLowerCamel}}, "SELECT * from {{.TableName}}")
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, Err{{.Name}}NotFound
		}

		return nil, err
	}

	return &{{.NameLowerCamel}}, nil
}
`))

	s.selectAllFuncSignatureTemplate = template.Must(template.New("select-all-signiture").Parse(`Get{{.Name}}s(ctx context.Context) (*[]{{.PackageName}}.{{.Name}}, error)`))

	s.selectFuncBodyTemplate = template.Must(template.New("select").Parse(`
func (r *mysql{{.Name}}) GetBy{{.FuncName}}(ctx context.Context, {{.Inputs}}) (*{{.PackageName}}.{{.Name}}, error) {
	var {{.NameLowerCamel}} {{.PackageName}}.{{.Name}}

	err := r.db.GetContext(ctx, &{{.NameLowerCamel}}, "SELECT * FROM {{.TableName}} " +
		"{{.Conditions}}",
		{{.ContextVariables}},
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, Err{{.Name}}NotFound
		}

		return nil, err
	}

	return &{{.NameLowerCamel}}, nil
}
`))

	s.selectFuncSignatureTemplate = template.Must(template.New("select-signiture").Parse(`GetBy{{.FuncName}}(ctx context.Context, {{.Inputs}}) (*{{.PackageName}}.{{.Name}}, error)`))

	s.updateAllFuncBodyTemplate = template.Must(template.New("update-all").Parse(`
func (r *mysql{{.Name}}) Update(ctx context.Context, {{.FieldNameLowerCamel}} {{.FieldType}}, {{.NameLowerCamel}} {{.PackageName}}.{{.Name}}) (int64, error) {
	{{.NameLowerCamel}}.{{.FieldName}} = {{.FieldNameLowerCamel}}

	result, err := r.db.NamedExecContext(ctx, "UPDATE {{.TableName}} "+
		"SET "+
		{{.FieldsContextKeys}} +
		"{{.Conditions}}",
		{{.NameLowerCamel}},
	)

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
`))

	s.updateAllFuncSignatureTemplate = template.Must(template.New("update-all").Parse(`Update(ctx context.Context, {{.FieldNameLowerCamel}} {{.FieldType}}, {{.NameLowerCamel}} {{.PackageName}}.{{.Name}}) (int64, error)`))

	s.updateFuncBodyTemplate = template.Must(template.New("update").Parse(`
func (r *mysql{{.Name}}) Update{{.FuncName}}(ctx context.Context, {{.InputConditions}}, {{.InputFields}}) (int64, error) {
	query := "UPDATE {{.TableName}} SET " +
			{{.FieldsContextKeys}} +
			"{{.Conditions}};" 

	result, err := r.db.ExecContext(ctx, query, {{.ExecVars}})

	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}
`))

	s.updateFuncSignatureTemplate = template.Must(template.New("update-signiture").Parse(`Update{{.FuncName}}(ctx context.Context, {{.InputConditions}}, {{.InputFields}}) (int64, error)`))

	s.joinFuncBodyTemplate = template.Must(template.New("join").Parse(`
func (r *mysql{{.Name}}) GetJoined{{.Name}}(ctx context.Context, limit uint) ([]{{.PackageName}}.{{.Name}}, error) {
	query := "SELECT " +
		{{.JoinFields}}
		"FROM {{.TableName}} AS {{.TableNameShort}} " +
		{{.Joins}} +
		"LIMIT ?"

	var {{.NameLowerCamel}} []{{.PackageName}}.{{.Name}}
	err := r.db.SelectContext(ctx, &{{.NameLowerCamel}}, query, limit)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, {{.NotFoundErr}}
		}

		return nil, err
	}

	return {{.NameLowerCamel}}, nil
}
`))

	s.joinFuncSignatureTemplate = template.Must(template.New("join-signiture").Parse(`GetJoined{{.Name}}(ctx context.Context, limit uint) ([]{{.PackageName}}.{{.Name}}, error)`))

	s.aggregateFuncBodyTemplate = template.Must(template.New("aggregate").Parse(`
func (r *mysql{{.Name}}) GetAggregateBy{{.FuncName}}(ctx context.Context, {{.InputConditions}}) (*int, error) {
	var res struct{` +
		"result int `db:\"{{.AggregateFieldAs}}\"`" + `
	}

	err := r.db.SelectContext(ctx, &res, "SELECT {{.AggregateSyntax}} FROM {{.TableName}} " +
		"{{.Conditions}}{{.GroupBy}}",
		 {{.ContextVariables}},
	)

	if err != nil {
		return nil, err
	}

	return &res.result, nil
}
`))

	s.aggregateFuncSignatureTemplate = template.Must(template.New("aggregate-signiture").Parse(`GetAggregateBy{{.FuncName}}(ctx context.Context, {{.Inputs}}) (*int, error)`))

	return s
}

func (s sqlx) Insert(structure *structure.Structure) (string, string) {
	fields := structure.GetDBFields(":")
	data := struct {
		Name            string
		NameLowerCamel  string
		PackageName     string
		NameSnake       string
		Fields          string
		FieldsWithColon string
	}{
		structure.Name,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		strcase.ToSnake(structure.Name),
		structure.GetDBFields(""),
		fields[:len(fields)-1],
	}
	var syntax strings.Builder
	if err := s.insertFuncSignatureTemplate.Execute(&syntax, data); err != nil {
		panic(err)
	}

	var signiture strings.Builder
	if err := s.insertFuncSignatureTemplate.Execute(&signiture, data); err != nil {
		panic(err)
	}

	return syntax.String(), signiture.String()
}

func (s sqlx) UpdateAll(structure *structure.Structure) (string, string) {
	data := struct {
		Name                string
		FieldNameLowerCamel string
		FieldType           string
		NameLowerCamel      string
		PackageName         string
		FieldName           string
		TableName           string
		FieldsContextKeys   string
		Conditions          string
	}{
		structure.Name,
		strcase.ToLowerCamel(structure.Fields[0].Name),
		structure.Fields[0].Type,
		strcase.ToLowerCamel(structure.Name),
		structure.PackageName,
		structure.Fields[0].Name,
		structure.TableName,
		contextKeys(structure.Fields),
		conditions([]string{
			structure.FieldMapNameToDBFlag[structure.Fields[0].Name],
		}, structure, false),
	}
	var syntax strings.Builder
	if err := s.updateAllFuncBodyTemplate.Execute(&syntax, data); err != nil {
		panic(err)
	}

	var signiture strings.Builder
	if err := s.updateAllFuncSignatureTemplate.Execute(&signiture, data); err != nil {
		panic(err)
	}

	return syntax.String(), signiture.String()
}

func (s sqlx) UpdateBy(structure *structure.Structure, vars *[]structure.UpdateVariables) (syntax string, signitures []string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, name := range v.Fields {
			functionNameList = append(functionNameList, structure.FieldMapDBFlagToName[name])
		}
		functionName := strings.Join(functionNameList, "And")

		data := struct {
			Name              string
			FuncName          string
			InputConditions   string
			InputFields       string
			TableName         string
			FieldsContextKeys string
			Conditions        string
			ExecVars          string
		}{
			structure.Name,
			functionName,
			inputFunctionVariables(v.Conditions, structure),
			inputFunctionVariables(v.Fields, structure),
			structure.TableName,
			contextKeys(v.Fields),
			conditions(v.Conditions, structure, true),
			execContextVariables(v, structure, false),
		}

		var syntaxBuilder strings.Builder
		if err := s.updateFuncBodyTemplate.Execute(&syntaxBuilder, data); err != nil {
			panic(err)
		}

		syntax += syntaxBuilder.String()

		var signaturesBuilder strings.Builder
		if err := s.updateFuncSignatureTemplate.Execute(&signaturesBuilder, data); err != nil {
			panic(err)
		}

		signitures = append(signitures, signaturesBuilder.String())
	}

	return syntax, signitures
}

func (s sqlx) SelectAll(structure *structure.Structure) (string, string) {

	data := struct {
		Name                string
		PackageName         string
		TableNameLowerCamel string
		NameLowerCamel      string
		TableName           string
	}{
		structure.Name,
		structure.PackageName,
		strcase.ToLowerCamel(structure.TableName),
		strcase.ToLowerCamel(structure.Name),
		structure.TableName,
	}

	var syntax strings.Builder
	if err := s.selectAllFuncBodyTemplate.Execute(&syntax, data); err != nil {
		panic(err)
	}

	var signiture strings.Builder
	if err := s.selectAllFuncSignatureTemplate.Execute(&signiture, data); err != nil {
		panic(err)
	}

	return syntax.String(), signiture.String()
}

func (s sqlx) SelectBy(structure *structure.Structure, vars *[]structure.GetVariable) (syntax string, signatures []string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, condition := range v.Conditions {
			functionNameList = append(functionNameList, structure.FieldMapDBFlagToName[condition])
		}
		functionName := strings.Join(functionNameList, "And")

		data := struct {
			Name             string
			FuncName         string
			Inputs           string
			PackageName      string
			NameLowerCamel   string
			TableName        string
			Conditions       string
			ContextVariables string
		}{
			structure.Name,
			functionName,
			inputFunctionVariables(v.Conditions, structure),
			structure.PackageName,
			strcase.ToLowerCamel(structure.Name),
			structure.TableName,
			conditions(v.Conditions, structure, true),
			contextVariables(v.Conditions, structure),
		}

		var syntaxBuilder strings.Builder
		if err := s.selectFuncBodyTemplate.Execute(&syntaxBuilder, data); err != nil {
			panic(err)
		}

		syntax += syntaxBuilder.String()

		var signitureBuilder strings.Builder
		if err := s.selectFuncSignatureTemplate.Execute(&signitureBuilder, data); err != nil {
			panic(err)
		}
		signatures = append(signatures, signitureBuilder.String())
	}

	return syntax, signatures
}

func (s sqlx) Join(structure *structure.Structure, joinVariables *structure.JoinVariables) (string, string) {

	data := struct {
		Name           string
		PackageName    string
		JoinFields     string
		TableName      string
		TableNameShort string
		Joins          string
		NameLowerCamel string
		NotFoundErr    string
	}{
		structure.Name,
		structure.PackageName,
		joinField(structure, joinVariables),
		structure.TableName,
		string(structure.TableName[0]),
		joins(structure, joinVariables),
		strcase.ToLowerCamel(structure.Name),
		"Err" + structure.Name + "NotFound",
	}

	var syntax strings.Builder
	if err := s.joinFuncBodyTemplate.Execute(&syntax, data); err != nil {
		panic(err)
	}

	var header strings.Builder
	if err := s.joinFuncSignatureTemplate.Execute(&header, data); err != nil {
		panic(err)
	}

	return syntax.String(), header.String()
}

func (s sqlx) Aggregate(structure *structure.Structure, vars *[]structure.AggregateField) (syntax string, signatures []string) {

	for _, v := range *vars {
		functionNameList := make([]string, 0)
		for _, condition := range v.Conditions {
			functionNameList = append(functionNameList, structure.FieldMapDBFlagToName[condition])
		}
		functionName := strings.Join(functionNameList, "And")

		data := struct {
			Name             string
			FuncName         string
			InputConditions  string
			AggregateFieldAs string
			AggregateSyntax  string
			TableName        string
			Conditions       string
			GroupBy          string
			ContextVariables string
		}{
			structure.Name,
			functionName,
			inputFunctionVariables(v.Conditions, structure),
			v.As,
			aggregateSyntax(v),
			structure.TableName,
			conditions(v.Conditions, structure, true),
			groupBy(v),
			contextVariables(v.Conditions, structure),
		}

		var syntaxBuilder strings.Builder
		if err := s.aggregateFuncBodyTemplate.Execute(&syntaxBuilder, data); err != nil {
			panic(err)
		}

		syntax += syntaxBuilder.String()

		var signitureBuilder strings.Builder
		if err := s.aggregateFuncSignatureTemplate.Execute(&signitureBuilder, data); err != nil {
			panic(err)
		}

		signatures = append(signatures, signitureBuilder.String())
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
