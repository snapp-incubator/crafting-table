package querybuilder

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/gertd/go-pluralize"
	"github.com/iancoleman/strcase"
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
			if t.Name == "int" || t.Name == "int8" || t.Name == "int16" || t.Name == "int32" || t.Name == "int64" ||
				t.Name == "uint" || t.Name == "uint8" || t.Name == "uint16" || t.Name == "uint32" || t.Name == "uint64" ||
				t.Name == "float32" || t.Name == "float64" {
				return true
			}
			return false
		}
	}
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
				IsNullable:   false, // TODO: fix this
			}
			if field.Tag != nil {
				sf.Tag = field.Tag.Value
			}
			fields = append(fields, sf)
		}
	}
	return fields
}

func Generate(dialect string, pkg string, structDecl *ast.GenDecl) string {
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

	for _, t := range queryBuilderTemplates {
		if err := t.Execute(&buff, td); err != nil {
			panic(err)
		}
	}

	return buff.String()
}
