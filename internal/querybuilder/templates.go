package querybuilder

import (
	"text/template"

	"github.com/iancoleman/strcase"
)

var funcMap template.FuncMap = template.FuncMap{
	"toSnakeCase": func(name string) string {
		return strcase.ToSnake(name)
	},
}

var (
	baseOutputFileTemplate        = template.Must(template.New("ct-base").Funcs(funcMap).Parse(baseOutputFile))
	setsTemplate                  = template.Must(template.New("ct-sets").Funcs(funcMap).Parse(sets))
	eqWhereTemplate               = template.Must(template.New("ct-eq-where").Funcs(funcMap).Parse(eqWhere))
	scalarWhereTemplate           = template.Must(template.New("ct-scalar-where").Funcs(funcMap).Parse(scalarWhere))
	limitOffsetTemplate           = template.Must(template.New("ct-limit-offset").Funcs(funcMap).Parse(limitOffset))
	queryBuilderTemplate          = template.Must(template.New("ct-query-builder").Funcs(funcMap).Parse(queryBuilder))
	selectsTemplate               = template.Must(template.New("ct-selects").Funcs(funcMap).Parse(selects))
	selectQueryBuilderTemplate    = template.Must(template.New("ct-select-builder").Funcs(funcMap).Parse(selectQueryBuilder))
	updateQueryBuilderTemplate    = template.Must(template.New("ct-update-builder").Funcs(funcMap).Parse(updateQueryBuilder))
	deleteQueryBuilderTemplate    = template.Must(template.New("ct-delete-builder").Funcs(funcMap).Parse(deleteQueryBuilder))
	fromRowsTemplate              = template.Must(template.New("ct-from-rows").Funcs(funcMap).Parse(fromRows))
	toRowsTemplate                = template.Must(template.New("ct-to-rows").Funcs(funcMap).Parse(toRows))
	finishersTemplate             = template.Must(template.New("ct-finishers").Funcs(funcMap).Parse(finishers))
	queryBuilderInterfaceTemplate = template.Must(template.New("ct-interface").Funcs(funcMap).Parse(queryBuilderInterface))
)

type templateData struct {
	Pkg       string
	ModelName string
	Fields    map[string]string
}

const queryBuilderInterface = `
type {{.ModelName}}SQLQueryBuilder interface{
	{{ range $field,$type := .Fields }}
	Where{{$field}}Eq({{$type}}) {{$.ModelName}}SQLQueryBuilder
	{{ if eq $type "int" "int8" "int16" "int32" "int64" "uint8" "uint16" "uint32" "uint64" "uint" "float32" "float64"  }}
	Where{{$field}}GT({{ $type }}) {{$.ModelName}}SQLQueryBuilder
	Where{{$field}}GE({{ $type }}) {{$.ModelName}}SQLQueryBuilder
	Where{{$field}}LT({{ $type }}) {{$.ModelName}}SQLQueryBuilder
	Where{{$field}}LE({{ $type }}) {{$.ModelName}}SQLQueryBuilder
	{{ end }}
	Select{{$field}}() {{$.ModelName}}SQLQueryBuilder
	Set{{$field}}({{$type}}) {{$.ModelName}}SQLQueryBuilder
	{{ end }}
	Limit(int) {{$.ModelName}}SQLQueryBuilder
	Offset(int) {{$.ModelName}}SQLQueryBuilder
	SelectAll() {{$.ModelName}}SQLQueryBuilder
	
	// finishers
	First(db *sql.DB) ({{ .ModelName }}, error)
	Last(db *sql.DB) ({{ .ModelName }}, error)
	Update(db *sql.DB) (sql.Result, error)
	Delete(db *sql.DB) (sql.Result, error)
	Fetch(db *sql.DB) ([]{{ .ModelName }}, error)
}
`

const limitOffset = `
func (q *__{{ .ModelName }}SQLQueryBuilder) Limit(l int) {{ .ModelName }}SQLQueryBuilder {
	q.mode = "select"
	q.limit = l	
	return q
}
func (q *__{{ .ModelName }}SQLQueryBuilder) Offset(l int) {{ .ModelName }}SQLQueryBuilder {
	q.mode = "select"
	q.offset = l
	return q
}
`

const toRows = `
func (m {{ .ModelName }}) Values() []interface{} {
    var values []interface{}
	{{ range $field, $type := .Fields }}
	values = append(values, &m.{{ $field }})
	{{ end }}
    return values
}
`

const fromRows = `
func {{ .ModelName }}sFromRows(rows *sql.Rows) ([]{{.ModelName}}, error) {
    var {{.ModelName}}s []{{.ModelName}}
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
        {{.ModelName}}s = append({{.ModelName}}s, m)
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

const finishers = `
func (q *__{{ .ModelName }}SQLQueryBuilder) Update(db *sql.DB) (sql.Result, error) {
	q.mode = "update"
	return db.Exec(q.SQL(), q.args...)
}

func (q *__{{ .ModelName }}SQLQueryBuilder) Delete(db *sql.DB) (sql.Result, error) {
	q.mode = "delete"
	return db.Exec(q.SQL(), q.args...)
}

func (q *__{{ .ModelName }}SQLQueryBuilder) Fetch(db *sql.DB) ([]{{ .ModelName }}, error) {
	q.mode = "select"
	rows, err := db.Query(q.SQL(), q.args...)
	if err != nil {
		return nil, err
	}
	return {{.ModelName}}sFromRows(rows)
}

func (q *__{{ .ModelName }}SQLQueryBuilder) First(db *sql.DB) ({{ .ModelName }}, error) {
	q.mode = "select"
	q.orderby = "ORDER BY ID ASC"
	q.Limit(1)
	row := db.QueryRow(q.SQL(), q.args...)
	if row.Err() != nil {
		return {{ .ModelName }} {}, row.Err()
	}
	return UserFromRow(row)
}


func (q *__{{ .ModelName }}SQLQueryBuilder) Last(db *sql.DB) ({{ .ModelName }}, error) {
	q.mode = "select"
	q.orderby = "ORDER BY ID DESC"
	q.Limit(1)
	row := db.QueryRow(q.SQL(), q.args...)
	if row.Err() != nil {
		return {{ .ModelName}} {}, row.Err()
	}
	return UserFromRow(row)
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

	limit int
	offset int

	args []interface{}
}

func {{.ModelName}}QueryBuilder() {{ .ModelName }}SQLQueryBuilder {
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
func (q *__{{ $.ModelName}}SQLQueryBuilder) Select{{$field}}() {{ $.ModelName }}SQLQueryBuilder {
	q.projected = append(q.projected, "{{ toSnakeCase $field }}")
	return q
}
{{ end }}

func (q *__{{ $.ModelName}}SQLQueryBuilder) SelectAll() {{ $.ModelName }}SQLQueryBuilder {
	q.projected = append(q.projected, "*")
	return q
}
`

const selectQueryBuilder = `
func (q *__{{ .ModelName}}SQLQueryBuilder) sqlSelect() string {
	if q.projected == nil {
		q.projected = append(q.projected, "*")
	}
	base := fmt.Sprintf("SELECT %s FROM %s", strings.Join(q.projected, ", "), q.table)

	var wheres []string 
	{{ range $field, $type := .Fields }}
	if q.where.{{$field}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "{{ toSnakeCase $field }}", q.where.{{$field}}.operator, fmt.Sprint(q.where.{{$field}}.argument)))
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
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "{{ toSnakeCase $field }}", q.where.{{$field}}.operator, fmt.Sprint(q.where.{{$field}}.argument)))
	}
	if q.set.{{$field}} != nil {
		sets = append(sets, fmt.Sprintf("%s = %s", "{{ toSnakeCase $field }}"), fmt.Sprint(q.set.{{$field}}))
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
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "{{ toSnakeCase $field }}"), q.where.{{$field}}.operator, fmt.Sprint(q.where.{{$field}}.argument))
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
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{$field}}GE({{$field}} {{$type}}) {{$.ModelName}}SQLQueryBuilder {
	m.mode = "select"
	m.where.{{$field}}.argument = {{$field}}
    m.where.{{$field}}.operator = ">="
	return m
}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{$field}}GT({{$field}} {{$type}}) {{$.ModelName}}SQLQueryBuilder {
	m.mode = "select"
    m.where.{{$field}}.argument = {{$field}}
    m.where.{{$field}}.operator = ">="
	return m
}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{$field}}LE({{$field}} {{$type}}) {{$.ModelName}}SQLQueryBuilder {
	m.mode = "select"
    m.where.{{$field}}.argument = {{$field}}
    m.where.{{$field}}.operator = "<="
	return m
}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{$field}}LT({{$field}} {{$type}}) {{$.ModelName}}SQLQueryBuilder {
	m.mode = "select"
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
func (m *__{{ $.ModelName }}SQLQueryBuilder) Where{{$field}}Eq({{ $field }} {{ $type }}) {{ $.ModelName }}SQLQueryBuilder {
	m.mode = "select"
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
func (m *__{{ $.ModelName }}SQLQueryBuilder) Set{{ $field }}({{ $field }} {{ $type }}) {{ $.ModelName }}SQLQueryBuilder {
	m.mode = "update"
	m.set.{{$field}} = &{{ $field }}
	return m
}
{{ end }}
`

const baseOutputFile = `// Code generated by Crafting-Table. DO NOT EDIT
// Source code: https://github.com/snapp-incubator/crafting-table
package {{ .Pkg }}

import (
    "fmt"
    "strings"
    "database/sql"
)
`
