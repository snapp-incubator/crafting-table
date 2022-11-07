package querybuilder

import (
	"fmt"
	"go/ast"
	"strings"
)

const ModelAnnotation = "ct: model"

type structField struct {
	Name         string
	Type         string
	IsComparable bool
	IsNullable   bool
	Tag          string
}

func (s structField) String() string {
	return s.Name
}

func resolveTypes(structDecl *ast.GenDecl) []structField {
	var fields []structField
	for _, field := range structDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
		for _, name := range field.Names {
			sf := structField{
				Name:         name.Name,
				Type:         fmt.Sprint(field.Type),
				IsComparable: false,
				IsNullable:   false,
			}
			if field.Tag != nil {
				sf.Tag = field.Tag.Value
			}
			fields = append(fields, sf)
		}
	}
	return fields
}

func Generate(pkg string, structDecl *ast.GenDecl, args map[string]string, dialect string) string {
	fields := resolveTypes(structDecl)
	typeName := structDecl.Specs[0].(*ast.TypeSpec).Name.String()
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

	err = queryBuilderInterfaceTemplate.Execute(&buff, td)
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

	err = limitOffsetTemplate.Execute(&buff, td)

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

	err = toRowsTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}
	err = orderByTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}
	err = placeholderGeneratorTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}

	err = finishersTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}

	return buff.String()
}
