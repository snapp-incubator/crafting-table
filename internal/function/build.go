package function

import (
	"strings"

	"github.com/iancoleman/strcase"

	"github.com/snapp-incubator/crafting-table/internal/structure"

	"github.com/snapp-incubator/crafting-table/internal/query"
)

func BuildGetFunction(
	structure *structure.Structure,
	table string,
	fields []string,
	where []query.WhereCondition,
	aggregate []query.AggregateField,
	orderBy *string,
	orderType *query.OrderType,
	limit *uint,
	groupBy []string,
	join []query.JoinField,
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
	q := query.BuildSelectQuery(
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

	functionName := "Get"
	if len(whereColumns) > 0 {
		functionName += "By" + strings.Join(whereColumns, "And") // GetByColumn1AndColumn2AndColumn3
	}

	// fields: prepare inputs
	inputList := make([]string, len(where))
	for i, v := range where {
		inputList[i] = v.Column + " " + structure.FieldMapNameToType[v.Column]
	}
	inputs := strings.Join(inputList, ", ")

	// fields: prepare DesStructTemplate
	desStructTemplate := ""
	if len(aggregate) > 0 {
		desStructTemplate += "type structDes struct {\n"
		for _, v := range aggregate {
			desStructTemplate += strcase.ToCamel(v.As) + " int  " + "`db:\"" + v.As + "\"`" + "\n"
		}
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
		outputList = append(outputList, "*"+structure.PackageName+structure.Name)
	}
	outputList = append(outputList, "error")
	outputs := strings.Join(outputList, ", ")

	// fields: prepare real outputs without error
	var realOutputList []string
	if len(aggregate) > 0 {
		for _, v := range aggregate {
			realOutputList = append(realOutputList, "structDes"+strcase.ToCamel(v.As))
		}
	} else {
		realOutputList = append(realOutputList, strcase.ToLowerCamel(structure.Name))
	}

	// fields: prepare Dest
	dest := strcase.ToLowerCamel(structure.Name)
	if len(aggregate) > 0 {
		dest = "structDes"
	}

	// fields: prepare builder
	var builder strings.Builder

	// create signature
	signatureData := struct {
		FuncName string
		Inputs   string
		Outputs  string
	}{
		FuncName: functionName,
		Inputs:   inputs,
		Outputs:  outputs,
	}
	if err := signature.Execute(&builder, signatureData); err != nil {
		panic(err)
	}
	signatureTemplate = builder.String()

	// create exec query
	outputsWithNotFoundError := make([]string, len(realOutputList)+1)
	for i, _ := range realOutputList {
		outputsWithNotFoundError[i] = "nil"
	}
	outputsWithNotFoundError[len(realOutputList)] = "Err" + structure.Name + "NotFound"

	execQueryData := struct {
		Query                    string
		Dest                     string
		OutputsWithNotFoundError string
	}{
		Query:                    q,
		Dest:                     dest,
		OutputsWithNotFoundError: strings.Join(outputsWithNotFoundError, ", "),
	}
	if err := getContext.Execute(&builder, execQueryData); err != nil {
		panic(err)
	}
	getContextQuery := builder.String()

	// create function
	functionData := struct {
		ModelName         string
		Signature         string
		DesStructTemplate string
		Model             string
		ExecQueryTemplate string
		Outputs           string
	}{
		ModelName:         structure.Name,
		Signature:         signatureTemplate,
		DesStructTemplate: desStructTemplate,
		Model:             model,
		ExecQueryTemplate: getContextQuery,
		Outputs:           strings.Join(append(realOutputList, "nil"), ", "),
	}

	if err := function.Execute(&builder, functionData); err != nil {
		panic(err)
	}
	functionTemplate = builder.String()

	return functionTemplate, signatureTemplate
}

func BuildSelectFunction(
	structure *structure.Structure,
	table string,
	fields []string,
	where []query.WhereCondition,
	aggregate []query.AggregateField,
	orderBy *string,
	orderType *query.OrderType,
	limit *uint,
	groupBy []string,
	join []query.JoinField,
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
	q := query.BuildSelectQuery(
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

	functionName := "Get"
	if len(whereColumns) > 0 {
		functionName += "By" + strings.Join(whereColumns, "And") // GetByColumn1AndColumn2AndColumn3
	}

	// fields: prepare inputs
	inputList := make([]string, len(where))
	for i, v := range where {
		inputList[i] = v.Column + " " + structure.FieldMapNameToType[v.Column]
	}
	inputs := strings.Join(inputList, ", ")

	// fields: prepare DesStructTemplate
	desStructTemplate := ""
	if len(aggregate) > 0 {
		desStructTemplate += "type structDes struct {\n"
		for _, v := range aggregate {
			desStructTemplate += strcase.ToCamel(v.As) + " int  " + "`db:\"" + v.As + "\"`" + "\n"
		}
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
		outputList = append(outputList, "*"+structure.PackageName+structure.Name)
	}
	outputList = append(outputList, "error")
	outputs := strings.Join(outputList, ", ")

	// fields: prepare real outputs without error
	var realOutputList []string
	if len(aggregate) > 0 {
		for _, v := range aggregate {
			realOutputList = append(realOutputList, "structDes"+strcase.ToCamel(v.As))
		}
	} else {
		realOutputList = append(realOutputList, strcase.ToLowerCamel(structure.Name))
	}

	// fields: prepare Dest
	dest := strcase.ToLowerCamel(structure.Name)
	if len(aggregate) > 0 {
		dest = "structDes"
	}

	// fields: prepare builder
	var builder strings.Builder

	// create signature
	signatureData := struct {
		FuncName string
		Inputs   string
		Outputs  string
	}{
		FuncName: functionName,
		Inputs:   inputs,
		Outputs:  outputs,
	}
	if err := signature.Execute(&builder, signatureData); err != nil {
		panic(err)
	}
	signatureTemplate = builder.String()

	// create exec query
	outputsWithError := make([]string, len(realOutputList)+1)
	for i, _ := range realOutputList {
		outputsWithError[i] = "nil"
	}
	outputsWithError[len(realOutputList)] = "err"

	execQueryData := struct {
		Query                    string
		Dest                     string
		OutputsWithNotFoundError string
	}{
		Query:                    q,
		Dest:                     dest,
		OutputsWithNotFoundError: strings.Join(outputsWithError, ", "),
	}
	if err := selectContext.Execute(&builder, execQueryData); err != nil {
		panic(err)
	}
	selectContextQuery := builder.String()

	// create function
	functionData := struct {
		ModelName         string
		Signature         string
		DesStructTemplate string
		Model             string
		ExecQueryTemplate string
		Outputs           string
	}{
		ModelName:         structure.Name,
		Signature:         signatureTemplate,
		DesStructTemplate: desStructTemplate,
		Model:             model,
		ExecQueryTemplate: selectContextQuery,
		Outputs:           strings.Join(append(realOutputList, "nil"), ", "),
	}

	if err := function.Execute(&builder, functionData); err != nil {
		panic(err)
	}
	functionTemplate = builder.String()

	return functionTemplate, signatureTemplate
}

// TODO: add more functions like: update, insert.

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
