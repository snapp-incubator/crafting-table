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
) string {
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
	signature := builder.String()

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
		Signature:         signature,
		DesStructTemplate: desStructTemplate,
		Model:             model,
		ExecQueryTemplate: "",
		Outputs:           "",
	}

	if err := function.Execute(&builder, functionData); err != nil {
		panic(err)
	}
	signature := builder.String()

	return ""
}
