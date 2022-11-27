package build

import (
	"fmt"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"

	"github.com/snapp-incubator/crafting-table/internal/structure"
)

type signatureParameters struct {
	FuncName string
	Inputs   string
	Outputs  string
}

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
) (functionTemplate string, signatureTemplate string) {
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
		whereColumns[i] = v.Column
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
		inputWithTypeList[i] = strcase.ToLowerCamel(v.Column) + " " + structure.FieldMapNameToType[v.Column]
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
	if err := signature.Execute(&signatureBuilder, signatureData); err != nil {
		panic(err)
	}
	signatureTemplate = signatureBuilder.String()

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
	if err := getContext.Execute(&getContextBuilder, getQueryData); err != nil {
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
		Signature:         signatureTemplate,
		DesStructTemplate: desStructTemplate,
		DstModel:          model,
		ExecQueryTemplate: getContextQuery,
		Outputs:           strings.Join(append(realOutputList, "nil"), ", "),
	}

	var functionBuilder strings.Builder
	if err := function.Execute(&functionBuilder, functionData); err != nil {
		panic(err)
	}
	functionTemplate = functionBuilder.String() // This is the result

	return functionTemplate, signatureTemplate
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
) (functionTemplate string, signatureTemplate string) {
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
		whereColumns[i] = v.Column
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
		inputWithTypeList[i] = strcase.ToLowerCamel(v.Column) + " " + structure.FieldMapNameToType[v.Column]
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
	// this struct is duplicated among all methods ,is this correct ?
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
	if err := signature.Execute(&signatureBuilder, signatureData); err != nil {
		panic(err)
	}
	signatureTemplate = signatureBuilder.String()

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
	if err := selectContext.Execute(&selectContextBuilder, execQueryData); err != nil {
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
		Signature:         signatureTemplate,
		DesStructTemplate: desStructTemplate,
		DstModel:          model,
		ExecQueryTemplate: selectContextQuery,
		Outputs:           strings.Join(append(realOutputList, "nil"), ", "),
	}

	var functionBuilder strings.Builder
	if err := function.Execute(&functionBuilder, functionData); err != nil {
		panic(err)
	}
	functionTemplate = functionBuilder.String()

	return functionTemplate, signatureTemplate
}

// TODO: add more functions like: update, insert.

func BuildInsertFunction(
	structure *structure.Structure,
	dialect DialectType,
	table string,
	fields []string,
	where []WhereCondition,
	objectName string,
	customFunctionName string,
) (functionTemplate string, signatureTemplate string) {
	// bring an example of wanted result
	/*
		_, err := r.db.NamedExecContext(ctx, "INSERT INTO cancellation_schedule_ride "+
			"(ride_id, match_time, driver_delay, driver_pickup_waiting_time, created_at) "+
			"values ( :ride_id,  :match_time,  :driver_delay,  :driver_pickup_waiting_time,  now())", ride)
		if err != nil {
			return err
		}
	*/
	var functionName string
	if customFunctionName == "" {
		functionName = "Create"
	} else {
		functionName = customFunctionName
	}

	inputs := fmt.Sprintf("%s *%s.%s", objectName, structure.PackageName, structure.Name)

	// make functions signature
	signatureData := signatureParameters{
		FuncName: functionName,
		Inputs:   inputs,
		Outputs:  "error",
	}
	var signatureBuilder strings.Builder
	if err := signature.Execute(&signatureBuilder, signatureData); err != nil {
		panic(err)
	}

	signatureTemplate = signatureBuilder.String()
	// is this template a good name ? because this is not a template actually .

	// make functions body
	insertQuery := BuildInsertQuery(
		dialect,
		table,
		fields,
		where,
	)
	functionData := struct {
		Query      string
		ObjectName string
	}{
		Query:      insertQuery,
		ObjectName: objectName,
	}
	var builder strings.Builder
	if err := insertContext.Execute(&builder, functionData); err != nil {
		panic(err)
	}

	functionTemplate = builder.String()

	return functionTemplate, signatureTemplate

	// return the result
}

func BuildRepository(
	signatureTemplateList []string,
	functionTemplateList []string,
	packageName string,
	tableName string,
	modelName string,
) (repositoryTemplate string) {
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
	if err := repository.Execute(&builder, repositoryData); err != nil {
		panic(err)
	}
	repositoryTemplate = builder.String()

	return repositoryTemplate
}

// Query to database
var selectContext *template.Template = template.Must(
	template.New("selectContext").Parse("{{ if .SpecialQuery }}query := \"{{.Query}}\"" +
		"{{ else }}query := `{{.Query}}`{{ end }} \n" +
		`err := d.db.SelectContext(ctx, &{{.Dest}}, query, {{.Inputs}})
if err != nil {
	return {{.OutputsWithErr}}
}
`))

var getContext *template.Template = template.Must(
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

// Complete it
var insertContext *template.Template = template.Must(
	template.New("insertContext").Parse("query := '{{.Query}}' \n" +
		`_, err := r.db.NamedExecContext(ctx, query , {{.ObjectName}})
if err != nil {
	return err
}
`))

var namedExecContext *template.Template = template.Must(
	template.New("namedExecContext").Parse("query := `{{.Query}}`\n" +
		`_, err := d.db.NamedExecContext(ctx, query, {{.Dest}})
if err != nil {
	return {{.OutputsWithErr}}
}
`))

var execContext *template.Template = template.Must(
	template.New("execContext").Parse("query := `{{.Query}}`\n" +
		`_, err := d.db.ExecContext(ctx, query, {{.ExecVars}})
if err != nil {
	return {{.OutputsWithErr}}
}
`))

// signature is function's signature
// I think it would be better if this name became "signatureTemplate"
var signature *template.Template = template.Must(
	template.New("signature").Parse(`{{.FuncName}}(ctx context.Context, {{.Inputs}}) ({{.Outputs}})`))

// function is function's body
var function *template.Template = template.Must(template.New("function").Parse(`
func (d *database{{.ModelName}}) {{.Signature}} {
	{{.DesStructTemplate}}

	var dst {{.DstModel}}

	{{.ExecQueryTemplate}}

	return {{.Outputs}}
}
`))

// repository is file's body
var repository *template.Template = template.Must(template.New("repository").Parse(`
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
