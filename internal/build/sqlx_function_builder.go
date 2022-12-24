package build

import (
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
func BuildUpdateFunction(
	structure *structure.Structure,
	dialect DialectType,
	table string,
	fields []string,
	where []WhereCondition,
	customFunctionName string,
) (function string, signature string) {
	// converting a []string to a []interface{}
	fieldsInterface := make([]interface{}, len(fields))
	for i, v := range fields {
		fieldsInterface[i] = v
	}

	// create query
	q := BuildUpdateQuery(
		dialect,
		table,
		fieldsInterface,
		where,
	)

	// fields: prepare functionName using field instead of where columns since updated fields are the main focus
	fieldColumns := make([]string, len(fields))
	for i, v := range fields {
		fieldColumns[i] = strcase.ToCamel(v)
	}

	var functionName string
	if customFunctionName == "" {
		functionName = "Update"
		if len(fieldColumns) > 0 {
			functionName += "" + strings.Join(fieldColumns, "And") // UpdateColumn1AndColumn2AndColumn3
		}
	} else {
		functionName = customFunctionName
	}
	// fields: prepare where inputs
	// "current" is appended to field in signature to make selecting and updating the same field possible
	// find row (where) using var1 (var1Current) and further down update row (set) using new value for var1 (var1New)
	whereInputWithTypeList := make([]string, len(where))
	whereInputList := make([]string, len(where))
	for i, v := range where {
		whereInputList[i] = strcase.ToLowerCamel(v.Column)
		whereInputWithTypeList[i] = strcase.ToLowerCamel(v.Column) + "Current" + " " +
			structure.FieldMapNameToType[structure.FieldMapDBFlagToName[v.Column]]
	}
	whereInputsWithType := strings.Join(whereInputWithTypeList, ", ")
	whereInputs := strings.Join(whereInputList, ", ")

	// fields: prepare field inputs
	// "new" is appended to field in signature to make selecting and updating the same field possible
	// further up find row (where) using var1 (var1Current) and update row (set) using new value for var1 (var1New)
	fieldInputWithTypeList := make([]string, len(fields))
	fieldInputList := make([]string, len(fields))
	for i, v := range fields {
		fieldInputList[i] = strcase.ToLowerCamel(v)
		fieldInputWithTypeList[i] = strcase.ToLowerCamel(v) + "New" + " " +
			structure.FieldMapNameToType[structure.FieldMapDBFlagToName[v]]
	}
	fieldInputsWithType := strings.Join(fieldInputWithTypeList, ", ")
	fieldInputs := strings.Join(fieldInputList, ", ")

	// fields: Combine field and where inputs for function signature
	completeInputListWithType := make([]string, 2)
	completeInputListWithType[0] = whereInputsWithType
	completeInputListWithType[1] = fieldInputsWithType
	completeInputsWithType := strings.Join(completeInputListWithType, ", ")

	// fields: prepare outputs (update only returns updated row count or error)
	var outputList []string

	outputList = append(outputList, "int64")

	outputList = append(outputList, "error")
	outputs := strings.Join(outputList, ", ")

	// fields: prepare real outputs without error (result is added in updateContextTemplate)
	var realOutputList []string
	realOutputList = append(realOutputList, "result.RowsAffected()")

	// create signature
	signatureData := struct {
		FuncName string
		Inputs   string
		Outputs  string
	}{
		FuncName: functionName,
		Inputs:   completeInputsWithType,
		Outputs:  outputs,
	}
	var signatureBuilder strings.Builder
	if err := signatureTemplate.Execute(&signatureBuilder, signatureData); err != nil {
		panic(err)
	}
	signature = signatureBuilder.String()

	// create exec query
	// error output should return 0 affected rows along with error
	outputsWithError := make([]string, len(realOutputList)+1)
	for i, _ := range realOutputList {
		outputsWithError[i] = "0"
	}
	outputsWithError[len(realOutputList)] = "err"

	specialQuery := false
	if dialect == MySQL || dialect == SQLite3 {
		specialQuery = true
	}

	// seperated where from set (field) inputs since they should come first
	execQueryData := struct {
		Query          string
		SpecialQuery   bool
		OutputsWithErr string
		WhereVars      string
		FieldVars      string
	}{
		Query:          q,
		SpecialQuery:   specialQuery,
		OutputsWithErr: strings.Join(outputsWithError, ", "),
		WhereVars:      whereInputs,
		FieldVars:      fieldInputs,
	}
	var updateContextBuilder strings.Builder
	if err := updateContextTemplate.Execute(&updateContextBuilder, execQueryData); err != nil {
		panic(err)
	}
	updateContextQuery := updateContextBuilder.String()

	// create function
	functionData := struct {
		ModelName         string
		Signature         string
		ExecQueryTemplate string
		Outputs           string
	}{
		ModelName:         structure.Name,
		Signature:         signature,
		ExecQueryTemplate: updateContextQuery,
		Outputs:           strings.Join(append(realOutputList, "nil"), ", "),
	}
	var updateFunctionBuilder strings.Builder
	if err := updateFunctionTemplate.Execute(&updateFunctionBuilder, functionData); err != nil {
		panic(err)
	}
	function = updateFunctionBuilder.String()

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

// Reused the execContext template and put set and where variables in separate inputs
// Renamed main output to result instead of ignoring it in order to get affected rows count further down
var updateContextTemplate *template.Template = template.Must(
	template.New("updateContext").Parse("{{ if .SpecialQuery }}query := \"{{.Query}}\"" +
		"{{ else }}query := `{{.Query}}`{{ end }} \n" +
		`result, err := d.db.ExecContext(ctx, query, {{.FieldVars}}, {{.WhereVars}})
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
	return {{.OutputsWithErr}}
}
`))

var execContextTemplate *template.Template = template.Must(
	template.New("execContext").Parse("{{ if .SpecialQuery }}query := \"{{.Query}}\"" +
		"{{ else }}query := `{{.Query}}`{{ end }} \n" +
		`_, err := d.db.ExecContext(ctx, query, {{.ExecVars}})
if err != nil {
	return {{.OutputsWithErr}}
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

// Simplified the function template for use in update queries since there is no output except the row count
var updateFunctionTemplate *template.Template = template.Must(template.New("updateFunction").Parse(`
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
