package cmd

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var goGenCmd = &cobra.Command{
	Use: "gogen",
	Run: func(cmd *cobra.Command, args []string) {
		run(args)
	},
}

const modelAnnotation = "ct: model"
const repoAnnotation = "ct: repo"

func generateQueryBuilder(pkg string, typeName string, fields map[string]string, tags []string, args map[string]string) string {
	// base file
	// query builder
	// eq where
	// scalar where for each field
	// sets
	var buff strings.Builder
	err := baseOutputFileTemplate.Execute(&buff, baseOutputFileTemplateData{
		Pkg: pkg,
	})
	if err != nil {
		panic(err)
	}

	err = queryBuilderTemplate.Execute(&buff, queryBuilderTemplateData{
		ModelName: typeName,
		Fields:    fields,
	})
	if err != nil {
		panic(err)
	}

	err = eqWhereTemplate.Execute(&buff, eqWhereTemplateData{
		ModelName: typeName,
		Fields:    fields,
	})

	if err != nil {
		panic(err)
	}

	err = scalarWhereTemplate.Execute(&buff, scalarWhereTemplateData{
		ModelName: typeName,
		Fields:    fields,
	})

	if err != nil {
		panic(err)
	}

	err = setsTemplate.Execute(&buff, setsTemplateData{
		ModelName: typeName,
		Fields:    fields,
	})
	if err != nil {
		panic(err)
	}
	return buff.String()
}

func run(args []string) {
	if len(args) < 1 {
		log.Fatalln("needs a filename")
	}
	filename := args[0]

	inputFilePath, err := filepath.Abs(filename)
	if err != nil {
		panic(err)
	}
	pathList := filepath.SplitList(inputFilePath)
	pathList = pathList[:len(pathList)-1]
	dir := filepath.Join(pathList...)
	fset := token.NewFileSet()
	fast, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)

	if err != nil {
		panic(err)
	}
	actualName := strings.TrimSuffix(filename, filepath.Ext(filename))
	outputFilePath := filepath.Join(dir, fmt.Sprintf("%s_sqlgen_gen.go", actualName))
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()
	for _, decl := range fast.Decls {
		if _, ok := decl.(*ast.GenDecl); ok {
			declComment := decl.(*ast.GenDecl).Doc.Text()
			if len(declComment) > 0 && declComment[:len(modelAnnotation)] == modelAnnotation {
				name := decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Name.String()
				// arguments := strings.Split(strings.Trim(declComment[len(annotation)+1:], " \n\t\r"), " ")
				fields := make(map[string]string)
				for _, field := range decl.(*ast.GenDecl).Specs[0].(*ast.TypeSpec).Type.(*ast.StructType).Fields.List {
					for _, name := range field.Names {
						fields[name.String()] = fmt.Sprint(field.Type)
					}
				}
				args := make(map[string]string)
				// for _, argkv := range arguments {
				// 	splitted := strings.Split(argkv, "=")
				// 	args[splitted[0]] = splitted[1]
				// }
				output := generateQueryBuilder(fast.Name.String(), name, fields, nil, args)
				fmt.Fprint(outputFile, output)
			}
		}

	}

}

var (
	baseOutputFileTemplate = template.Must(template.New("sqlgen-base").Parse(baseOutputFile))
	setsTemplate           = template.Must(template.New("sqlgen-sets").Parse(sets))
	eqWhereTemplate        = template.Must(template.New("sqlgen-eq-where").Parse(eqWhere))
	scalarWhereTemplate    = template.Must(template.New("sqlgen-scalar-where").Parse(scalarWhere))
	queryBuilderTemplate   = template.Must(template.New("sqlgen-query-builder").Parse(queryBuilder))
)

type queryBuilderTemplateData struct {
	ModelName string
	Fields    map[string]string
}

const queryBuilder = `
type __{{ .ModelName }}SQLQueryBuilder struct {
	mode string
    where __{{ .ModelName }}Where
	set __{{ .ModelName }}Set
	orderby string
	groupby string
	table string
	projected []string
}

func {{.ModelName}}QueryBuilder() *__{{ .ModelName }}SQLQueryBuilder {
	return &__{{ .ModelName }}SQLQueryBuilder{}
}


{{ range $field, $type := .Fields }}
func (q *__{{ $.ModelName}}SQLQueryBuilder) Select{{$field}}() *__{{ $.ModelName }}SQLQueryBuilder {
	q.projected = append(q.projected, strcase.ToSnake("{{ $field }}"))
	return q
}
{{ end }}

func (q *__{{ $.ModelName}}SQLQueryBuilder) SelectAll() *__{{ $.ModelName }}SQLQueryBuilder {
	q.projected = append(q.projected, "*")
	return q
}

func (q *__{{ .ModelName}}SQLQueryBuilder) sqlSelect() string {
	base := fmt.Sprintf("SELECT %s FROM %s", strings.Join(q.projected, ", "), q.table)

	var wheres []string 
	{{ range $field, $type := .Fields }}
	if q.where.{{$field}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", strcase.ToSnake("{{ $field }}"), q.where.{{$field}}.operator, fmt.Sprint(q.where.{{$field}}.argument)))
	}
	{{ end }}
	if len(wheres) > 0 {
		base += "WHERE " + strings.Join(wheres, " AND ")
	}
	return base
}
func (q *__{{ .ModelName}}SQLQueryBuilder) sqlUpdate() string {
	base := fmt.Sprintf("UPDATE %s", q.table)

	var wheres []string 
    var sets []string 

    {{ range $field, $type := .Fields }}
	if q.where.{{$field}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", strcase.ToSnake("{{ $field }}"), q.where.{{$field}}.operator, fmt.Sprint(q.where.{{$field}}.argument)))
	}
	if q.set.{{$field}} != nil {
		sets = append(sets, fmt.Sprintf("%s = %s", strcase.ToSnake("{{ $field }}"), fmt.Sprint(q.set.{{$field}})))
	}
    {{ end }}


	if len(wheres) > 0 {
		base += " WHERE " + strings.Join(wheres, " AND ")
	}

	if len(sets) > 0 {
		base += " SET " + strings.Join(sets, " , ")
	}

	return base
}
func (q *__{{ .ModelName}}SQLQueryBuilder) sqlDelete() string {
    base := fmt.Sprintf("DELETE FROM %s", q.table)

	var wheres []string 
	{{ range $field, $type := .Fields }}
	if q.where.{{$field}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", strcase.ToSnake("{{ $field }}"), q.where.{{$field}}.operator, fmt.Sprint(q.where.{{$field}}.argument)))
	}
	{{ end }}
	if len(wheres) > 0 {
		base += " WHERE " + strings.Join(wheres, " AND ")
	}

	return base

}

func (q *__{{ .ModelName }}SQLQueryBuilder) SQL() string {
	if q.mode == "select" {
		return q.sqlSelect()
	} else if q.mode == "update" {
		return q.sqlUpdate()
	} else if q.mode == "delete" {
		return q.sqlDelete()
	} else {
		panic("unsupported query mode")
	}
}

`

type scalarWhereTemplateData struct {
	ModelName string
	Fields    map[string]string
}

const scalarWhere = `
{{ range $field, $type := .Fields }}
{{ if eq $type "int" "int8" "int16" "int32" "int64" "uint8" "uint16" "uint32" "uint64" "uint" "float32" "float64"  }}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{$field}}GE({{$field}} {{$type}}) *__{{$.ModelName}}SQLQueryBuilder {
	m.where.{{$field}}.argument = {{$field}}
    m.where.{{$field}}.operator = ">="
	return m
}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{$field}}GT({{$field}} {{$type}}) *__{{$.ModelName}}SQLQueryBuilder {
    m.where.{{$field}}.argument = {{$field}}
    m.where.{{$field}}.operator = ">="
	return m
}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{$field}}LE({{$field}} {{$type}}) *__{{$.ModelName}}SQLQueryBuilder {
    m.where.{{$field}}.argument = {{$field}}
    m.where.{{$field}}.operator = "<="
	return m
}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{$field}}LT({{$field}} {{$type}}) *__{{$.ModelName}}SQLQueryBuilder {
    m.where.{{$field}}.argument = {{$field}}
    m.where.{{$field}}.operator = "<="
	return m
}
{{ end }}
{{ end }}
`

type eqWhereTemplateData struct {
	ModelName string
	Fields    map[string]string
}

const eqWhere = `
type __{{ .ModelName }}Where struct {
	{{ range $field, $type := .Fields }}
	{{$field}} struct {
        argument {{$type}}
        operator string
    }
	{{ end }}
}
{{ range $field, $type := .Fields }}
func (m *__{{ $.ModelName }}SQLQueryBuilder) Where{{$field}}Eq({{ $field }} {{ $type }}) *__{{ $.ModelName }}SQLQueryBuilder {
	m.where.{{$field}}.argument = {{$field}}
    m.where.{{$field}}.operator = "="
	return m
}
{{ end }}
`

type setsTemplateData struct {
	ModelName string
	Fields    map[string]string
}

const sets = `
type __{{ .ModelName }}Set struct {
	{{ range $field, $type := .Fields }}
	{{$field}} *{{$type}}
	{{ end }}
}
{{ range $field, $type := .Fields }}
func (m *__{{ $.ModelName }}SQLQueryBuilder) Set{{ $field }}({{ $field }} {{ $type }}) *__{{ $.ModelName }}SQLQueryBuilder {
	m.set.{{$field}} = &{{ $field }}
	return m
}
{{ end }}
`

type baseOutputFileTemplateData struct {
	Pkg string
}

const baseOutputFile = `
// Code generated by Crafting-Table.
// Source code: https://github.com/snapp-incubator/crafting-table
package {{ .Pkg }}

import (
    "fmt"
    "strings"

    "github.com/iancoleman/strcase"
)
`
