package querybuilder

import (
	"strings"
)

const ModelAnnotation = "ct: model"

func Generate(pkg string, typeName string, fields map[string]string, tags []string, args map[string]string, dialect string) string {
	var buff strings.Builder
	td := templateData{
		ModelName: typeName,
		Fields:    fields,
		Pkg:       pkg,
	}
	err := baseOutputFileTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}

	err = queryBuilderTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}

	err = selectsTemplate.Execute(&buff, td)

	if err != nil {
		panic(err)
	}

	err = selectQueryBuilderTemplate.Execute(&buff, td)

	if err != nil {
		panic(err)
	}

	err = updateQueryBuilderTemplate.Execute(&buff, td)

	if err != nil {
		panic(err)
	}

	err = deleteQueryBuilderTemplate.Execute(&buff, td)

	if err != nil {
		panic(err)
	}

	err = eqWhereTemplate.Execute(&buff, td)

	if err != nil {
		panic(err)
	}

	err = scalarWhereTemplate.Execute(&buff, td)

	if err != nil {
		panic(err)
	}

	err = setsTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}

	err = fromRowsTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}
	return buff.String()
}
