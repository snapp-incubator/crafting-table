package querybuilder

import "text/template"

var (
	baseOutputFileTemplate     = template.Must(template.New("ct-base").Parse(baseOutputFile))
	setsTemplate               = template.Must(template.New("ct-sets").Parse(sets))
	eqWhereTemplate            = template.Must(template.New("ct-eq-where").Parse(eqWhere))
	scalarWhereTemplate        = template.Must(template.New("ct-scalar-where").Parse(scalarWhere))
	queryBuilderTemplate       = template.Must(template.New("ct-query-builder").Parse(queryBuilder))
	selectsTemplate            = template.Must(template.New("ct-selects").Parse(selects))
	selectQueryBuilderTemplate = template.Must(template.New("ct-select-builder").Parse(selectQueryBuilder))
	updateQueryBuilderTemplate = template.Must(template.New("ct-update-builder").Parse(updateQueryBuilder))
	deleteQueryBuilderTemplate = template.Must(template.New("ct-delete-builder").Parse(deleteQueryBuilder))
	fromRowsTemplate           = template.Must(template.New("ct-from-rows").Parse(fromRows))
)

type templateData struct {
	Pkg       string
	ModelName string
	Fields    map[string]string
}

const fromRows = `
func {{ .ModelName }}sFromRows(rows *sql.Rows) ([]*{{.ModelName}}, error) {
    var {{.ModelName}}s []*{{.ModelName}}
    for rows.Next() {
        var m {{ .ModelName }}
        err := rows.Scan(
            {{ range $field, $type := .Fields }}
            &m.{{ $field }},
            {{ end }}
        )
        if err != nil {
            return nil, err
        }
        {{.ModelName}}s = append({{.ModelName}}s, &m)
    }
    return {{.ModelName}}s, nil
}

func {{ .ModelName }}FromRow(row *sql.Row) ({{.ModelName}}, error) {
    if row.Err() != nil {
        return {{.ModelName}}{}, row.Err()
    }
    var m {{ .ModelName }}
    err := row.Scan(
        {{ range $field, $type := .Fields }}
        &m.{{ $field }},
        {{ end }}
    )
    if err != nil {
        return {{.ModelName}}{}, err
    }
    return m, nil
}
`

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

const selects = `
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
`

const selectQueryBuilder = `
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
`

var updateQueryBuilder = `
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
`

const deleteQueryBuilder = `
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
`

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

const baseOutputFile = `
// Code generated by Crafting-Table. DO NOT EDIT
// Source code: https://github.com/snapp-incubator/crafting-table
package {{ .Pkg }}

import (
    "fmt"
    "strings"
    "database/sql"

    "github.com/iancoleman/strcase"
)
`
