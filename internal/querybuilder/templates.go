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
	placeholderGeneratorTemplate  = template.Must(template.New("ct-placeholder").Funcs(funcMap).Parse(placeholderGenerator))
	orderByTemplate               = template.Must(template.New("ct-orderby").Funcs(funcMap).Parse(orderby))
	queryBuilderInterfaceTemplate = template.Must(template.New("ct-interface").Funcs(funcMap).Parse(queryBuilderInterface))
)

type templateData struct {
	Pkg       string
	ModelName string
	Fields    []structField
}

const queryBuilderInterface = `
type {{.ModelName}}SQLQueryBuilder interface{
	{{ range .Fields }}
	Where{{.Name}}Eq({{.Type}}) {{$.ModelName}}SQLQueryBuilder
	{{ if eq .Type "int" "int8" "int16" "int32" "int64" "uint8" "uint16" "uint32" "uint64" "uint" "float32" "float64"  }}
	Where{{.Name}}GT({{ .Type }}) {{$.ModelName}}SQLQueryBuilder
	Where{{.Name}}GE({{ .Type }}) {{$.ModelName}}SQLQueryBuilder
	Where{{.Name}}LT({{ .Type }}) {{$.ModelName}}SQLQueryBuilder
	Where{{.Name}}LE({{ .Type }}) {{$.ModelName}}SQLQueryBuilder
	{{ end }}
	OrderBy{{.Name}}Asc() {{$.ModelName}}SQLQueryBuilder
	OrderBy{{.Name}}Desc() {{$.ModelName}}SQLQueryBuilder
	Select{{.Name}}() {{$.ModelName}}SQLQueryBuilder
	Set{{.Name}}({{.Type}}) {{$.ModelName}}SQLQueryBuilder
	{{ end }}
	Limit(int) {{$.ModelName}}SQLQueryBuilder
	Offset(int) {{$.ModelName}}SQLQueryBuilder
	SelectAll() {{$.ModelName}}SQLQueryBuilder

    getPlaceholder() string
	
	// finishers
	First(db *sql.DB) ({{ .ModelName }}, error)
	Last(db *sql.DB) ({{ .ModelName }}, error)
	Update(db *sql.DB) (sql.Result, error)
	Delete(db *sql.DB) (sql.Result, error)
	Fetch(db *sql.DB) ([]{{ .ModelName }}, error)


}
`

const placeholderGenerator = `
func (q *__{{ .ModelName }}SQLQueryBuilder) getPlaceholder() string {
     if q.dialect == "mysql" { return "?" }
     else if q.dialect == "postgres" { return fmt.Sprintf("$", len(q.args) + 1) }
     else { log.Fatalf("dialect %s not supported\n", q.dialect) }
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
	{{ range .Fields }}
	values = append(values, &m.{{ .Name }})
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
            {{ range .Fields }}
            &m.{{ .Name }},
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
        {{ range .Fields }}
        &m.{{ .Name }},
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
	q.orderby = []string{"ORDER BY ID ASC"}
	q.Limit(1)
	row := db.QueryRow(q.SQL(), q.args...)
	if row.Err() != nil {
		return {{ .ModelName }} {}, row.Err()
	}
	return UserFromRow(row)
}


func (q *__{{ .ModelName }}SQLQueryBuilder) Last(db *sql.DB) ({{ .ModelName }}, error) {
	q.mode = "select"
	q.orderby = []string{"ORDER BY ID DESC"}
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

    dialect string

    where __{{ .ModelName }}Where

	set __{{ .ModelName }}Set

	orderby []string
	groupby string

	table string

	projected []string

	limit int
	offset int

	whereArgs []interface{}
    setArgs []interface{}
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
{{ range .Fields }}
func (q *__{{ $.ModelName}}SQLQueryBuilder) Select{{.Name}}() {{ $.ModelName }}SQLQueryBuilder {
	q.projected = append(q.projected, "{{ toSnakeCase .Name }}")
	return q
}
{{ end }}

func (q *__{{ $.ModelName}}SQLQueryBuilder) SelectAll() {{ $.ModelName }}SQLQueryBuilder {
	q.projected = append(q.projected, "*")
	return q
}
`

const orderby = `
{{ range .Fields }}
func (q *__{{ $.ModelName}}SQLQueryBuilder) OrderBy{{.Name}}Asc() {{ $.ModelName }}SQLQueryBuilder {
	q.orderby = append(q.orderby, "{{ toSnakeCase .Name }} ASC")
	return q
}
func (q *__{{ $.ModelName}}SQLQueryBuilder) OrderBy{{.Name}}Desc() {{ $.ModelName }}SQLQueryBuilder {
	q.orderby = append(q.orderby, "{{ toSnakeCase .Name }} DESC")
	return q
}
{{ end }}
`

const selectQueryBuilder = `
func (q *__{{ .ModelName}}SQLQueryBuilder) sqlSelect() string {
	if q.projected == nil {
		q.projected = append(q.projected, "*")
	}
	base := fmt.Sprintf("SELECT %s FROM %s", strings.Join(q.projected, ", "), q.table)

	var wheres []string 
	{{ range .Fields }}
	if q.where.{{.Name}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "{{ toSnakeCase .Name }}", q.where.{{ .Name }}.operator, fmt.Sprint(q.where.{{ .Name }}.argument)))
	}
	{{ end }}
	if len(wheres) > 0 {
		base += "WHERE " + strings.Join(wheres, " AND ")
	}

	if len(q.orderby) > 0 {
		base += fmt.Sprintf(" ORDER BY %s", strings.Join(q.orderby, ", "))
	}
	return base
}
`

var updateQueryBuilder = `
func (q *__{{ .ModelName}}SQLQueryBuilder) sqlUpdate() string {
	base := fmt.Sprintf("UPDATE %s", q.table)

	var wheres []string 
    var sets []string 

    {{ range .Fields }}
	if q.where.{{.Name}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "{{ toSnakeCase .Name }}", q.where.{{ .Name }}.operator, fmt.Sprint(q.where.{{ .Name }}.argument)))
	}
	if q.set.{{ .Name }} != nil {
		sets = append(sets, fmt.Sprintf("%s = %s", "{{ toSnakeCase .Name }}", fmt.Sprint(q.set.{{ .Name }})))
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
	{{ range .Fields }}
	if q.where.{{.Name}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "{{ toSnakeCase .Name  }}", q.where.{{.Name }}.operator, fmt.Sprint(q.where.{{.Name }}.argument)))
	}
	{{ end }}
	if len(wheres) > 0 {
		base += " WHERE " + strings.Join(wheres, " AND ")
	}

	return base

}
`

const scalarWhere = `
{{ range .Fields }}
{{ if eq .Type "int" "int8" "int16" "int32" "int64" "uint8" "uint16" "uint32" "uint64" "uint" "float32" "float64"  }}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{.Name}}GE({{.Name }} {{.Type}}) {{$.ModelName}}SQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, {{.Name }})
    m.where.{{.Name }}.argument = m.getPlaceholder()
    m.where.{{.Name }}.operator = ">="
	return m
}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{.Name }}GT({{.Name }} {{.Type}}) {{$.ModelName}}SQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, {{.Name }})
    m.where.{{.Name }}.argument = m.getPlaceholder()
    m.where.{{.Name }}.operator = ">="
	return m
}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{.Name }}LE({{.Name }} {{.Type}}) {{$.ModelName}}SQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, {{.Name }})
    m.where.{{.Name }}.argument = m.getPlaceholder()
    m.where.{{.Name }}.operator = "<="
	return m
}
func (m *__{{$.ModelName}}SQLQueryBuilder) Where{{.Name }}LT({{.Name }} {{.Type}}) {{$.ModelName}}SQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, {{.Name }})
    m.where.{{.Name }}.argument = m.getPlaceholder()
    m.where.{{.Name }}.operator = "<="
	return m
}
{{ end }}
{{ end }}
`

const eqWhere = `
type __{{ .ModelName }}Where struct {
	{{ range .Fields }}
	{{.Name}} struct {
        argument {{.Type}}
        operator string
    }
	{{ end }}
}
{{ range .Fields }}
func (m *__{{ $.ModelName }}SQLQueryBuilder) Where{{.Name}}Eq({{ .Name }} {{ .Type }}) {{ $.ModelName }}SQLQueryBuilder {
    m.whereArgs = append(m.whereArgs, {{.Name}})
    m.where.{{.Name}}.argument = m.getPlaceholder()
    m.where.{{.Name}}.operator = "="
	return m
}
{{ end }}
`

const sets = `
type __{{ .ModelName }}Set struct {
	{{ range .Fields }}
	{{.Name }} *{{ .Type }}
	{{ end }}
}
{{ range .Fields }}
func (m *__{{ $.ModelName }}SQLQueryBuilder) Set{{ .Name }}({{ .Name }} {{ .Type }}) {{ $.ModelName }}SQLQueryBuilder {
	m.mode = "update"
	m.set.{{.Name}} = &{{ .Name }}
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
    "log"
)
`
