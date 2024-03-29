package build

import (
	"fmt"
	"log"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

func BuildGetFunction(
	structure *structure.Structure,
	dialect DialectType,
	table string,
	fields []string,
	where []WhereCondition,
	aggregate []AggregateField,
	orderBy *string,
	orderType *OrderType,
	limit *uint,
	groupBy []string,
	join []JoinField,
	customFunctionName string,
) (function string, signature string) {
	// converting a []string to a []interface{}
	fieldsInterface := make([]interface{}, len(fields))
	groupByInterface := make([]interface{}, len(groupBy))
	for i, v := range fields {
		fieldsInterface[i] = v
	}
	for i, v := range groupBy {
		groupByInterface[i] = v
	}

	// create query
	q := BuildSelectQuery(
		dialect,
		table,
		fieldsInterface,
		where,
		aggregate,
		orderBy,
		orderType,
		limit,
		groupByInterface,
		join,
	)

	// fields: prepare functionName
	whereColumns := make([]string, len(where))
	for i, v := range where {
		whereColumns[i] = strcase.ToCamel(v.Column)
	}

	var functionName string
	if customFunctionName == "" {
		functionName = "Get"
		if len(whereColumns) > 0 {
			functionName += "By" + strings.Join(whereColumns, "And") // GetByColumn1AndColumn2AndColumn3
		}
	} else {
		functionName = customFunctionName
	}

	// fields: prepare inputs
	inputWithTypeList := make([]string, len(where))
	inputList := make([]string, len(where))

	for i, v := range where {
		inputList[i] = strcase.ToLowerCamel(v.Column)
		inputWithTypeList[i] = strcase.ToLowerCamel(v.Column) + " " +
			structure.FieldMapNameToType[structure.FieldMapDBFlagToName[v.Column]]
	}
	inputsWithType := strings.Join(inputWithTypeList, ", ")
	inputs := strings.Join(inputList, ", ")

	// fields: prepare DesStructTemplate
	// TODO: add fields to DesStructTemplate
	desStructTemplate := ""
	if len(aggregate) > 0 {
		desStructTemplate += "type structDes struct {\n"
		for _, v := range aggregate {
			desStructTemplate += strcase.ToCamel(v.As) + " int  " + "`db:\"" + v.As + "\"`" + "\n"
		}
		desStructTemplate += "}"
	}

	// fields: prepare model
	model := structure.PackageName + "." + structure.Name
	if desStructTemplate != "" {
		model += "structDes\n"
	}

	// fields: prepare outputs
	var outputList []string
	if len(aggregate) > 0 {
		for i := 1; i <= len(aggregate); i++ {
			outputList = append(outputList, "*int")
		}
	} else {
		outputList = append(outputList, "*"+structure.PackageName+"."+structure.Name)
	}
	outputList = append(outputList, "error")
	outputs := strings.Join(outputList, ", ")

	// fields: prepare real outputs without error
	var realOutputList []string
	if len(aggregate) > 0 {
		for _, v := range aggregate {
			realOutputList = append(realOutputList, "&dst"+strcase.ToCamel(v.As))
		}
	} else {
		realOutputList = append(realOutputList, "&dst")
	}

	// create signature
	signatureData := struct {
		FuncName string
		Inputs   string
		Outputs  string
	}{
		FuncName: functionName,
		Inputs:   inputsWithType,
		Outputs:  outputs,
	}

	var signatureBuilder strings.Builder
	if err := signatureTemplate.Execute(&signatureBuilder, signatureData); err != nil {
		panic(err)
	}
	signature = signatureBuilder.String()

	// create exec query
	var outputsWithError []string
	for i := 0; i < len(realOutputList); i++ {
		outputsWithError = append(outputsWithError, "nil")
	}

	specialQuery := false
	if dialect == MySQL || dialect == SQLite3 {
		specialQuery = true
	}

	getQueryData := struct {
		Query                  string
		SpecialQuery           bool
		Dest                   string
		OutputsWithNotFoundErr string
		OutputsWithErr         string
		Inputs                 string
	}{
		Query:        q,
		SpecialQuery: specialQuery,
		Dest:         "dst",
		OutputsWithNotFoundErr: strings.Join(
			append(outputsWithError, "Err"+structure.Name+"NotFound"), ", "),
		OutputsWithErr: strings.Join(append(outputsWithError, "err"), ", "),
		Inputs:         inputs,
	}
	var getContextBuilder strings.Builder
	if err := getContextTemplate.Execute(&getContextBuilder, getQueryData); err != nil {
		panic(err)
	}
	getContextQuery := getContextBuilder.String()

	// create function
	functionData := struct {
		ModelName         string
		Signature         string
		DesStructTemplate string
		DstModel          string
		ExecQueryTemplate string
		Outputs           string
	}{
		ModelName:         structure.Name,
		Signature:         signature,
		DesStructTemplate: desStructTemplate,
		DstModel:          model,
		ExecQueryTemplate: getContextQuery,
		Outputs:           strings.Join(append(realOutputList, "nil"), ", "),
	}

	var functionBuilder strings.Builder
	if err := functionTemplate.Execute(&functionBuilder, functionData); err != nil {
		panic(err)
	}
	function = functionBuilder.String()

	return function, signature
}

func BuildSelectFunction(
	structure *structure.Structure,
	dialect DialectType,
	table string,
	fields []string,
	where []WhereCondition,
	aggregate []AggregateField,
	orderBy *string,
	orderType *OrderType,
	limit *uint,
	groupBy []string,
	join []JoinField,
	customFunctionName string,
) (function string, signature string) {
	// converting a []string to a []interface{}
	fieldsInterface := make([]interface{}, len(fields))
	groupByInterface := make([]interface{}, len(groupBy))
	for i, v := range fields {
		fieldsInterface[i] = v
	}
	for i, v := range groupBy {
		groupByInterface[i] = v
	}

	// create query
	q := BuildSelectQuery(
		dialect,
		table,
		fieldsInterface,
		where,
		aggregate,
		orderBy,
		orderType,
		limit,
		groupByInterface,
		join,
	)

	// fields: prepare functionName
	whereColumns := make([]string, len(where))
	for i, v := range where {
		whereColumns[i] = strcase.ToCamel(v.Column)
	}

	var functionName string
	if customFunctionName == "" {
		functionName = "Select"
		if len(whereColumns) > 0 {
			functionName += "By" + strings.Join(whereColumns, "And") // GetByColumn1AndColumn2AndColumn3
		}
	} else {
		functionName = customFunctionName
	}

	// fields: prepare inputs
	inputWithTypeList := make([]string, len(where))
	inputList := make([]string, len(where))
	for i, v := range where {
		inputList[i] = strcase.ToLowerCamel(v.Column)
		inputWithTypeList[i] = strcase.ToLowerCamel(v.Column) + " " +
			structure.FieldMapNameToType[structure.FieldMapDBFlagToName[v.Column]]
	}
	inputsWithType := strings.Join(inputWithTypeList, ", ")
	inputs := strings.Join(inputList, ", ")

	// fields: prepare DesStructTemplate
	desStructTemplate := ""
	if len(aggregate) > 0 {
		desStructTemplate += "type structDes struct {\n"
		for _, v := range aggregate {
			desStructTemplate += strcase.ToCamel(v.As) + " int  " + "`db:\"" + v.As + "\"`" + "\n"
		}
		desStructTemplate += "}"
	}

	// fields: prepare model
	model := structure.PackageName + "." + structure.Name
	if desStructTemplate != "" {
		model = "structDes\n"
	}

	// fields: prepare outputs
	var outputList []string
	if len(aggregate) > 0 {
		for i := 1; i <= len(aggregate); i++ {
			outputList = append(outputList, "*int")
		}
	} else {
		outputList = append(outputList, "*"+structure.PackageName+"."+structure.Name)
	}
	outputList = append(outputList, "error")
	outputs := strings.Join(outputList, ", ")

	// fields: prepare real outputs without error
	var realOutputList []string
	if len(aggregate) > 0 {
		for _, v := range aggregate {
			realOutputList = append(realOutputList, "&dst."+strcase.ToCamel(v.As))
		}
	} else {
		realOutputList = append(realOutputList, "&dst")
	}

	// create signature
	signatureData := struct {
		FuncName string
		Inputs   string
		Outputs  string
	}{
		FuncName: functionName,
		Inputs:   inputsWithType,
		Outputs:  outputs,
	}
	var signatureBuilder strings.Builder
	if err := signatureTemplate.Execute(&signatureBuilder, signatureData); err != nil {
		panic(err)
	}
	signature = signatureBuilder.String()

	// create exec query
	outputsWithError := make([]string, len(realOutputList)+1)
	for i, _ := range realOutputList {
		outputsWithError[i] = "nil"
	}
	outputsWithError[len(realOutputList)] = "err"

	specialQuery := false
	if dialect == MySQL || dialect == SQLite3 {
		specialQuery = true
	}

	execQueryData := struct {
		Query          string
		SpecialQuery   bool
		Dest           string
		OutputsWithErr string
		Inputs         string
	}{
		Query:          q,
		SpecialQuery:   specialQuery,
		Dest:           "dst",
		OutputsWithErr: strings.Join(outputsWithError, ", "),
		Inputs:         inputs,
	}
	var selectContextBuilder strings.Builder
	if err := selectContextTemplate.Execute(&selectContextBuilder, execQueryData); err != nil {
		panic(err)
	}
	selectContextQuery := selectContextBuilder.String()

	// create function
	functionData := struct {
		ModelName         string
		Signature         string
		DesStructTemplate string
		DstModel          string
		ExecQueryTemplate string
		Outputs           string
	}{
		ModelName:         structure.Name,
		Signature:         signature,
		DesStructTemplate: desStructTemplate,
		DstModel:          model,
		ExecQueryTemplate: selectContextQuery,
		Outputs:           strings.Join(append(realOutputList, "nil"), ", "),
	}

	var functionBuilder strings.Builder
	if err := functionTemplate.Execute(&functionBuilder, functionData); err != nil {
		panic(err)
	}
	function = functionBuilder.String()

	return function, signature
}

// TODO: add more functions like: update, insert.

func BuildInsertFunction(
	structure *structure.Structure,
	dialect DialectType,
	table string,
	fields []string,
	withObject bool,
	customFunctionName string,
) (function string, signature string) {
	var functionName string
	if customFunctionName == "" {
		functionName = "Create"
	} else {
		functionName = customFunctionName
	}

	for _, f := range fields {
		_, ok := structure.FieldMapDBFlagToName[f]
		if !ok {
			log.Fatalf("field %s not found in structure", f)
		}
	}

	var inputs string
	var execVars string
	if withObject {
		inputs = fmt.Sprintf(
			"%s *%s.%s",
			strcase.ToLowerCamel(structure.Name),
			structure.PackageName,
			structure.Name)
	} else {
		if len(fields) == 0 {
			panic("fields is empty")
		}

		for _, f := range fields {
			name := structure.FieldMapDBFlagToName[f]
			fieldType := structure.FieldMapNameToType[name]
			inputs += fmt.Sprintf("%s %s, ", strcase.ToLowerCamel(name), fieldType)
			execVars += fmt.Sprintf("%s, ", strcase.ToLowerCamel(name))
		}
	}

	// make functions signature
	signatureData := struct {
		FuncName string
		Inputs   string
		Outputs  string
	}{
		FuncName: functionName,
		Inputs:   inputs,
		Outputs:  "error",
	}
	var signatureBuilder strings.Builder
	if err := signatureTemplate.Execute(&signatureBuilder, signatureData); err != nil {
		panic(err)
	}

	signature = signatureBuilder.String()

	// make functions body
	insertQuery := BuildInsertQuery(
		dialect,
		table,
		fields,
		withObject,
	)

	specialQuery := false
	if dialect == MySQL || dialect == SQLite3 {
		specialQuery = true
	}

	var execQueryBuilder strings.Builder
	if withObject {
		execQueryData := struct {
			SpecialQuery bool
			Query        string
			Dest         string
		}{
			SpecialQuery: specialQuery,
			Query:        insertQuery,
			Dest:         strcase.ToLowerCamel(structure.Name),
		}
		if err := namedExecContextTemplate.Execute(&execQueryBuilder, execQueryData); err != nil {
			panic(err)
		}
	} else {
		execQueryData := struct {
			SpecialQuery bool
			Query        string
			ExecVars     string
		}{
			SpecialQuery: specialQuery,
			Query:        insertQuery,
			ExecVars:     execVars,
		}
		if err := execContextTemplate.Execute(&execQueryBuilder, execQueryData); err != nil {
			panic(err)
		}
	}

	insertContextQuery := execQueryBuilder.String()

	functionData := struct {
		ModelName         string
		Signature         string
		ExecQueryTemplate string
		Outputs           string
	}{
		ModelName:         structure.Name,
		Signature:         signature,
		ExecQueryTemplate: insertContextQuery,
		Outputs:           "nil",
	}

	var functionBuilder strings.Builder
	if err := insertFunctionTemplate.Execute(&functionBuilder, functionData); err != nil {
		panic(err)
	}
	function = functionBuilder.String()

	return function, signature
}

func BuildRepository(
	signatureTemplateList []string,
	functionTemplateList []string,
	packageName string,
	tableName string,
	modelName string,
) (repository string) {
	// fields: prepare builder
	var builder strings.Builder

	// create repository
	repositoryData := struct {
		PackageName string
		ModelName   string
		Signatures  string
		TableName   string
		Functions   string
	}{
		PackageName: packageName,
		ModelName:   modelName,
		Signatures:  strings.Join(signatureTemplateList, "\n"),
		TableName:   tableName,
		Functions:   strings.Join(functionTemplateList, "\n"),
	}
	if err := repositoryTemplate.Execute(&builder, repositoryData); err != nil {
		panic(err)
	}
	repository = builder.String()

	return repository
}

// Query to database
var selectContextTemplate *template.Template = template.Must(
	template.New("selectContext").Parse("{{ if .SpecialQuery }}query := \"{{.Query}}\"" +
		"{{ else }}query := `{{.Query}}`{{ end }} \n" +
		`err := d.db.SelectContext(ctx, &{{.Dest}}, query, {{.Inputs}})
if err != nil {
	return {{.OutputsWithErr}}
}
`))

var getContextTemplate *template.Template = template.Must(
	template.New("getContext").Parse("{{ if .SpecialQuery }}query := \"{{.Query}}\"" +
		"{{ else }}query := `{{.Query}}`{{ end }} \n" +
		`err := d.db.GetContext(ctx, &{{.Dest}}, query, {{.Inputs}})
if err != nil {
	if err == sql.ErrNoRows {
		return {{.OutputsWithNotFoundErr}}
	}

	return {{.OutputsWithErr}}
}
`))

var namedExecContextTemplate *template.Template = template.Must(
	template.New("namedExecContext").Parse("{{ if .SpecialQuery }}query := \"{{.Query}}\"" +
		"{{ else }}query := `{{.Query}}`{{ end }} \n" +
		`_, err := d.db.NamedExecContext(ctx, query, {{.Dest}})
if err != nil {
	return err
}
`))

var execContextTemplate *template.Template = template.Must(
	template.New("execContext").Parse("{{ if .SpecialQuery }}query := \"{{.Query}}\"" +
		"{{ else }}query := `{{.Query}}`{{ end }} \n" +
		`_, err := d.db.ExecContext(ctx, query, {{.ExecVars}})
if err != nil {
	return err
}
`))

// signature is function's signature
var signatureTemplate *template.Template = template.Must(
	template.New("signature").Parse(`{{.FuncName}}(ctx context.Context, {{.Inputs}}) ({{.Outputs}})`))

// function is function's body
var functionTemplate *template.Template = template.Must(template.New("function").Parse(`
func (d *database{{.ModelName}}) {{.Signature}} {
	{{.DesStructTemplate}}

	var dst {{.DstModel}}

	{{.ExecQueryTemplate}}

	return {{.Outputs}}
}
`))

// insertFunctionTemplate is function's body for insert methods
var insertFunctionTemplate *template.Template = template.Must(template.New("function").Parse(`
func (d *database{{.ModelName}}) {{.Signature}} {
	{{.ExecQueryTemplate}}
	return {{.Outputs}}
}
`))

// repository is file's body
var repositoryTemplate *template.Template = template.Must(template.New("repository").Parse(`
// Code generated by Crafting-Table.
// Source code: https://github.com/snapp-incubator/crafting-table

package {{.PackageName}}

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
)

type {{.ModelName}} interface {
	{{.Signatures}}
}

var Err{{.ModelName}}NotFound = errors.New("{{.TableName}} not found")

type database{{.ModelName}} struct {
	db *sqlx.DB
}

func New{{.ModelName}}(db *sqlx.DB) {{.ModelName}} {
	return &database{{.ModelName}}{db: db}
}

{{.Functions}}
`))
