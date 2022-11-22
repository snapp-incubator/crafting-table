package querybuilder

import (
	"fmt"
	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
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

func isComparable(typeExpr ast.Expr) bool {
	switch t := typeExpr.(type) {
	case *ast.Ident:
		if t.Obj == nil {
			// it's a primitive go type
			if t.Name == "int" ||
				t.Name == "int8" ||
				t.Name == "int16" ||
				t.Name == "int32" ||
				t.Name == "int64" ||
				t.Name == "uint" ||
				t.Name == "uint8" ||
				t.Name == "uint16" ||
				t.Name == "uint32" ||
				t.Name == "uint64" ||
				t.Name == "float32" ||
				t.Name == "float64" {
				return true
			}
			return false
		}
	}
	return false
}

func isNullable(typeExpr ast.Expr) bool {
	return false
}

func resolveTypes(structDecl *ast.GenDecl) []structField {
	var fields []structField
	for _, field := range structDecl.Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
		for _, name := range field.Names {
			sf := structField{
				Name:         name.Name,
				Type:         fmt.Sprint(field.Type),
				IsComparable: isComparable(field.Type),
				IsNullable:   isNullable(field.Type),
			}
			if field.Tag != nil {
				sf.Tag = field.Tag.Value
			}
			fields = append(fields, sf)
		}
	}
	return fields
}

func getSampleValues(fields []structField) []any {
    var values []any
    for _, field := range fields {
        switch (field.Type) {
            case "string":
               values = append(values, "some string") 
           case "int":
           case "int8":
           case "int16":
           case "int32":
           case "int64":
           case "uint":
           case "uint8":
           case "uint16":
           case "uint32":
           case "uint64":
               values = append(values, 1)
           case "float32":
           case "float64":
               values = append(values, 2.2)
           default:
               values = append(values, nil)
        }
    }

    return values
}

func Generate(dialect string, pkg string, structDecl *ast.GenDecl, args map[string]string) string {
	fields := resolveTypes(structDecl)
	typeName := structDecl.Specs[0].(*ast.TypeSpec).Name.String()
	var buff strings.Builder
	td := templateData{
		ModelName: typeName,
		Fields:    fields,
		Pkg:       pkg,
		Dialect:   dialect,
		TableName: strcase.ToSnake(pluralize.NewClient().Plural(typeName)),
	}
	err := baseOutputFileTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}

	err = queryBuilderInterfaceTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}

	err = schemaTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}

	err = orderByTemplate.Execute(&buff, td)
	if err != nil {
		panic(err)
	}

	err = queryBuilderTemplate.Execute(&buff, td)
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
