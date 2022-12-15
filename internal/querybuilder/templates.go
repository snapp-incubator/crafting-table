package querybuilder

import (
	"text/template"

	"github.com/iancoleman/strcase"
)

var funcMap = template.FuncMap{
	"toSnakeCase": func(name string) string {
		return strcase.ToSnake(name)
	},
	"ToLowerCamelCase": func(name string) string {
		return strcase.ToLowerCamel(name)
	},
}

var (
	baseOutputFileTemplate        = template.Must(template.New("ct-base").Funcs(funcMap).Parse(baseOutputFile))
	setsTemplate                  = template.Must(template.New("ct-sets").Funcs(funcMap).Parse(sets))
	eqWhereTemplate               = template.Must(template.New("ct-eq-where").Funcs(funcMap).Parse(eqWhere))
	scalarWhereTemplate           = template.Must(template.New("ct-scalar-where").Funcs(funcMap).Parse(scalarWhere))
	limitOffsetTemplate           = template.Must(template.New("ct-limit-offset").Funcs(funcMap).Parse(limitOffset))
	queryBuilderTemplate          = template.Must(template.New("ct-query-builder").Funcs(funcMap).Parse(queryBuilder))
	selectQueryBuilderTemplate    = template.Must(template.New("ct-select-builder").Funcs(funcMap).Parse(selectQueryBuilder))
	updateQueryBuilderTemplate    = template.Must(template.New("ct-update-builder").Funcs(funcMap).Parse(updateQueryBuilder))
	deleteQueryBuilderTemplate    = template.Must(template.New("ct-delete-builder").Funcs(funcMap).Parse(deleteQueryBuilder))
	fromRowsTemplate              = template.Must(template.New("ct-from-rows").Funcs(funcMap).Parse(fromRows))
	toRowsTemplate                = template.Must(template.New("ct-to-rows").Funcs(funcMap).Parse(toRows))
	finishersTemplate             = template.Must(template.New("ct-finishers").Funcs(funcMap).Parse(finishers))
	placeholderGeneratorTemplate  = template.Must(template.New("ct-placeholder").Funcs(funcMap).Parse(placeholderGenerator))
	schemaTemplate                = template.Must(template.New("ct-schema").Funcs(funcMap).Parse(schema))
	queryBuilderInterfaceTemplate = template.Must(template.New("ct-interface").Funcs(funcMap).Parse(queryBuilderInterface))
	orderByTemplate               = template.Must(template.New("ct-orderby").Funcs(funcMap).Parse(orderBy))
)

var queryBuilderTemplates = []*template.Template{
	baseOutputFileTemplate,
	queryBuilderInterfaceTemplate,
	schemaTemplate,
	orderByTemplate,
	queryBuilderTemplate,
	selectQueryBuilderTemplate,
	limitOffsetTemplate,
	updateQueryBuilderTemplate,
	deleteQueryBuilderTemplate,
	eqWhereTemplate,
	scalarWhereTemplate,
	setsTemplate,
	fromRowsTemplate,
	toRowsTemplate,
	placeholderGeneratorTemplate,
	finishersTemplate,
}

type templateData struct {
	Pkg       string
	ModelName string
	TableName string
	Fields    []structField
	Dialect   string
}

const queryBuilderInterface = `
type {{.ModelName}}SQLQueryBuilder interface{
	{{ range .Fields }}Set{{.Name}}({{.Type}}) {{$.ModelName}}SQLQueryBuilder
	Where{{.Name}}Eq({{.Type}}) {{$.ModelName}}SQLQueryBuilder
	{{ if .IsComparable  }}Where{{.Name}}GT({{ .Type }}) {{$.ModelName}}SQLQueryBuilder
	Where{{.Name}}GE({{ .Type }}) {{$.ModelName}}SQLQueryBuilder
	Where{{.Name}}LT({{ .Type }}) {{$.ModelName}}SQLQueryBuilder
	Where{{.Name}}LE({{ .Type }}) {{$.ModelName}}SQLQueryBuilder {{ end }}
	{{ end }}
	OrderByAsc(column {{$.ModelName}}Column) {{$.ModelName}}SQLQueryBuilder
	OrderByDesc(column {{$.ModelName}}Column) {{$.ModelName}}SQLQueryBuilder

	Limit(int) {{$.ModelName}}SQLQueryBuilder
	Offset(int) {{$.ModelName}}SQLQueryBuilder

    getPlaceholder() string
	
	First(db *sql.DB) ({{ .ModelName }}, error)
	Last(db *sql.DB) ({{ .ModelName }}, error)
	Update(db *sql.DB) (sql.Result, error)
	Delete(db *sql.DB) (sql.Result, error)
	Fetch(db *sql.DB) ([]{{ .ModelName }}, error)
	SQL() (string, error)
}
`

const schema = `
type {{.ModelName}}Column string 
var {{.ModelName}}Columns = struct {
	{{ range .Fields }} {{.Name}} {{$.ModelName}}Column 
	{{ end }}
}{
	{{ range .Fields }} {{.Name}}: {{$.ModelName}}Column("{{ toSnakeCase .Name }}"), 
	{{ end }}
}
`

const placeholderGenerator = `
{{ if eq .Dialect "mysql" }}func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) getPlaceholder() string {
	return "?"
}
{{ end }}
{{ if eq .Dialect "sqlite" }}func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) getPlaceholder() string {
	return "?"
}
{{ end }}
{{ if eq .Dialect "postgres" }}func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) getPlaceholder() string {
	return fmt.Sprintf("$", len(q.whereArgs) + len(q.setArgs) + 1) 
}
{{ end }}
`

const limitOffset = `
func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) Limit(l int) {{ .ModelName }}SQLQueryBuilder {
	q.mode = "select"
	q.limit = l	
	return q
}

func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) Offset(l int) {{ .ModelName }}SQLQueryBuilder {
	q.mode = "select"
	q.offset = l
	return q
}
`

const toRows = `
func (q {{ .ModelName }}) Values() []interface{} {
    var values []interface{}
	{{ range .Fields }}values = append(values, &q.{{ .Name }})
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
    var q {{ .ModelName }}
    err := row.Scan(
        {{ range .Fields }}&q.{{ .Name }},
        {{ end }}
    )
    if err != nil {
        return {{.ModelName}}{}, err
    }

    return q, nil
}
`

const finishers = `
func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) Update(db *sql.DB) (sql.Result, error) {
	q.mode = "update"
	args := append(q.setArgs, q.whereArgs...)
	query, err := q.SQL()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, args...)
}

func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) Delete(db *sql.DB) (sql.Result, error) {
	q.mode = "delete"
	query, err := q.SQL()
	if err != nil {
		return nil, err
	}
	return db.Exec(query, q.whereArgs...)
}

func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) Fetch(db *sql.DB) ([]{{ .ModelName }}, error) {
	q.mode = "select"
	query, err := q.SQL()
	if err != nil {
		return nil, err
	}
	rows, err := db.Query(query, q.whereArgs...)
	if err != nil {
		return nil, err
	}
	return {{.ModelName}}sFromRows(rows)
}

func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) First(db *sql.DB) ({{ .ModelName }}, error) {
	q.mode = "select"
	q.orderBy = []string{"ORDER BY id ASC"}
	q.Limit(1)
	query, err := q.SQL()
	if err != nil {
		return {{ .ModelName }}{}, err
	}
	row := db.QueryRow(query, q.whereArgs...)
	if row.Err() != nil {
		return {{ .ModelName }}{}, row.Err()
	}
	return UserFromRow(row)
}


func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) Last(db *sql.DB) ({{ .ModelName }}, error) {
	q.mode = "select"
	q.orderBy = []string{"ORDER BY id DESC"}
	q.Limit(1)
	query, err := q.SQL()
	if err != nil {
		return {{ .ModelName }}{}, err
	}
	row := db.QueryRow(query, q.whereArgs...)
	if row.Err() != nil {
		return {{ .ModelName}}{}, row.Err()
	}
	return UserFromRow(row)
}
`

const queryBuilder = `
type {{ToLowerCamelCase .ModelName }}SQLQueryBuilder struct {
	mode string

    where {{ToLowerCamelCase .ModelName }}Where

	set {{ToLowerCamelCase .ModelName }}Set

	orderBy []string
	groupBy string

	projected []string

	limit int
	offset int

	whereArgs []interface{}
    setArgs []interface{}
}

func New{{.ModelName}}QueryBuilder() {{ .ModelName }}SQLQueryBuilder {
	return &{{ToLowerCamelCase .ModelName }}SQLQueryBuilder{}
}

func (q *{{ToLowerCamelCase .ModelName }}SQLQueryBuilder) SQL() (string, error) {
	if q.mode == "select" {
		return q.sqlSelect()
	} else if q.mode == "update" {
		return q.sqlUpdate()
	} else if q.mode == "delete" {
		return q.sqlDelete()
	} else {
		return "", fmt.Errorf("unsupported query mode '%s'", q.mode)
	}
}

`

const newselect = `
func (q *{{ToLowerCamelCase .ModelName}}SQLQueryBuilder) Select(column {{.ModelName}}Column) {{ $.ModelName }}SQLQueryBuilder {
	q.projected = append(q.projected, string(column))
	return q
}
`

const orderBy = `
func (q *{{ToLowerCamelCase $.ModelName}}SQLQueryBuilder) OrderByAsc(column {{.ModelName}}Column) {{ $.ModelName }}SQLQueryBuilder {
    q.mode = "select"
	q.orderBy = append(q.orderBy, fmt.Sprintf("%s ASC", string(column)))
	return q
}

func (q *{{ToLowerCamelCase $.ModelName}}SQLQueryBuilder) OrderByDesc(column {{.ModelName}}Column) {{ $.ModelName }}SQLQueryBuilder {
    q.mode = "select"
	q.orderBy = append(q.orderBy, fmt.Sprintf("%s DESC", string(column)))
	return q
}
`
const selectQueryBuilder = `
func (q *{{ToLowerCamelCase .ModelName}}SQLQueryBuilder) sqlSelect() (string, error) {
	if q.projected == nil {
		q.projected = append(q.projected, "*")
	}
	base := fmt.Sprintf("SELECT %s FROM {{ .TableName }}", strings.Join(q.projected, ", "))

	var wheres []string 
	{{ range .Fields }}
	if q.where.{{.Name}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "{{ toSnakeCase .Name }}", q.where.{{ .Name }}.operator, fmt.Sprint(q.where.{{ .Name }}.argument)))
	}
	{{ end }}
	if len(wheres) > 0 {
		base += " WHERE " + strings.Join(wheres, " AND ")
	}

	if len(q.orderBy) > 0 {
		base += fmt.Sprintf(" ORDER BY %s", strings.Join(q.orderBy, ", "))
	}

	if q.limit != 0 {
		base += " LIMIT " + fmt.Sprint(q.limit)
	}
	if q.offset != 0 {
		base += " OFFSET " + fmt.Sprint(q.offset)
	}
	return base, nil
}
`

var updateQueryBuilder = `
func (q *{{ToLowerCamelCase .ModelName}}SQLQueryBuilder) sqlUpdate() (string, error) {
	base := fmt.Sprintf("UPDATE {{.TableName}} ")

	var wheres []string 
    var sets []string 

    {{ range .Fields }}
	if q.where.{{.Name}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "{{ toSnakeCase .Name }}", q.where.{{ .Name }}.operator, fmt.Sprint(q.where.{{ .Name }}.argument)))
	}
	if q.set.{{ .Name }} != "" {
		sets = append(sets, fmt.Sprintf("%s = %s", "{{ toSnakeCase .Name }}", fmt.Sprint(q.set.{{ .Name }})))
	}
    {{ end }}

	if len(sets) > 0 {
		base += "SET " + strings.Join(sets, " , ")
	}

	if len(wheres) > 0 {
		base += " WHERE " + strings.Join(wheres, " AND ")
	}

	

	return base, nil
}
`

const deleteQueryBuilder = `
func (q *{{ToLowerCamelCase .ModelName}}SQLQueryBuilder) sqlDelete() (string, error) {
    base := fmt.Sprintf("DELETE FROM {{ .TableName }}")

	var wheres []string 
	{{ range .Fields }}
	if q.where.{{.Name}}.operator != "" {
		wheres = append(wheres, fmt.Sprintf("%s %s %s", "{{ toSnakeCase .Name  }}", q.where.{{.Name }}.operator, fmt.Sprint(q.where.{{.Name }}.argument)))
	}
	{{ end }}
	if len(wheres) > 0 {
		base += " WHERE " + strings.Join(wheres, " AND ")
	}

	return base, nil
}
`

const scalarWhere = `
{{ range .Fields }}
{{ if .IsComparable  }}
func (q *{{ToLowerCamelCase $.ModelName}}SQLQueryBuilder) Where{{.Name}}GE({{.Name }} {{.Type}}) {{$.ModelName}}SQLQueryBuilder {
    q.whereArgs = append(q.whereArgs, {{.Name }})
    q.where.{{.Name }}.argument = q.getPlaceholder()
    q.where.{{.Name }}.operator = ">="
	return q
}

func (q *{{ToLowerCamelCase $.ModelName}}SQLQueryBuilder) Where{{.Name }}GT({{.Name }} {{.Type}}) {{$.ModelName}}SQLQueryBuilder {
    q.whereArgs = append(q.whereArgs, {{.Name }})
    q.where.{{.Name }}.argument = q.getPlaceholder()
    q.where.{{.Name }}.operator = ">"
	return q
}

func (q *{{ToLowerCamelCase $.ModelName}}SQLQueryBuilder) Where{{.Name }}LE({{.Name }} {{.Type}}) {{$.ModelName}}SQLQueryBuilder {
    q.whereArgs = append(q.whereArgs, {{.Name }})
    q.where.{{.Name }}.argument = q.getPlaceholder()
    q.where.{{.Name }}.operator = "<="
	return q
}

func (q *{{ToLowerCamelCase $.ModelName}}SQLQueryBuilder) Where{{.Name }}LT({{.Name }} {{.Type}}) {{$.ModelName}}SQLQueryBuilder {
    q.whereArgs = append(q.whereArgs, {{.Name }})
    q.where.{{.Name }}.argument = q.getPlaceholder()
    q.where.{{.Name }}.operator = "<"
	return q
}
{{ end }}{{ end }}
`

const eqWhere = `
type {{ToLowerCamelCase .ModelName }}Where struct {
	{{ range .Fields }}
	{{.Name}} struct {
        argument interface{}
        operator string
    }
	{{ end }}
}
{{ range .Fields }}
func (q *{{ToLowerCamelCase $.ModelName }}SQLQueryBuilder) Where{{.Name}}Eq({{ .Name }} {{ .Type }}) {{ $.ModelName }}SQLQueryBuilder {
    q.whereArgs = append(q.whereArgs, {{.Name}})
    q.where.{{.Name}}.argument = q.getPlaceholder()
    q.where.{{.Name}}.operator = "="
	return q
}
{{ end }}
`

const sets = `
type {{ToLowerCamelCase .ModelName }}Set struct {
	{{ range .Fields }}
	{{.Name }} string
    {{ end }}
}
{{ range .Fields }}
func (q *{{ToLowerCamelCase $.ModelName }}SQLQueryBuilder) Set{{ .Name }}({{ .Name }} {{ .Type }}) {{ $.ModelName }}SQLQueryBuilder {
	q.mode = "update"
    q.setArgs = append(q.setArgs, {{ .Name }})
	q.set.{{.Name}} = q.getPlaceholder()
	return q
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
