package build

import (
	"strings"

	"github.com/doug-martin/goqu/v9"
)

type OrderType string

const (
	OrderTypeAsc  OrderType = "asc"
	OrderTypeDesc OrderType = "desc"
)

type OperatorType string

const (
	OperatorTypeEqual     OperatorType = "equal"
	OperatorTypeNotEqual  OperatorType = "not_equal"
	OperatorTypeIn        OperatorType = "in"
	OperatorTypeNotIn     OperatorType = "not_in"
	OperatorTypeGt        OperatorType = "gt"
	OperatorTypeGte       OperatorType = "gte"
	OperatorTypeLt        OperatorType = "lt"
	OperatorTypeLte       OperatorType = "lte"
	OperatorTypeIsNull    OperatorType = "is_null"
	OperatorTypeIsNotNull OperatorType = "is_not_null"
)

type JoinType string

const (
	JoinTypeJoin         JoinType = "Join"
	JoinTypeInner        JoinType = "inner"
	JoinTypeFullOuter    JoinType = "fullOuter"
	JoinTypeRightOuter   JoinType = "rightOuter"
	JoinTypeLeftOuter    JoinType = "leftOuter"
	JoinTypeFull         JoinType = "full"
	JoinTypeLeft         JoinType = "left"
	JoinTypeRight        JoinType = "right"
	JoinTypeNatural      JoinType = "natural"
	JoinTypeNaturalLeft  JoinType = "naturalLeft"
	JoinTypeNaturalRight JoinType = "naturalRight"
	JoinTypeNaturalFull  JoinType = "naturalFull"
	JoinTypeCross        JoinType = "cross"
)

// AggregateField is a struct for aggregate field
type AggregateField struct {
	Function string
	On       string
	As       string
}

// WhereCondition is a condition for where clause
type WhereCondition struct {
	Column   string
	Operator OperatorType
}

// JoinField is a struct for join field
type JoinField struct {
	Table    string
	As       string
	OnSource string
	OnJoin   string
	Function JoinType
}

// TODO: ADD DB NAME TO INPUT

// BuildSelectQuery builds a select query
func BuildSelectQuery(
	table string,
	fields []interface{},
	where []WhereCondition,
	aggregate []AggregateField,
	orderBy *string,
	orderType *OrderType,
	limit *uint,
	groupBy []interface{},
	join []JoinField,
) string {
	ds := goqu.From(table)

	// Aggregate: e.g. COUNT, SUM, MIN, MAX, AVG, FIRST, LAST
	// TODO: ADD AGGREGATE FUNCTION TO MAP
	aggregateExpressions := make([]interface{}, 0)
	if len(aggregate) > 0 {
		for _, agg := range aggregate {
			switch agg.Function {
			case "COUNT":
				aggregateExpressions = append(aggregateExpressions, goqu.COUNT(agg.On).As(agg.As))
			case "SUM":
				aggregateExpressions = append(aggregateExpressions, goqu.SUM(agg.On).As(agg.As))
			case "AVG":
				aggregateExpressions = append(aggregateExpressions, goqu.AVG(agg.On).As(agg.As))
			case "MAX":
				aggregateExpressions = append(aggregateExpressions, goqu.MAX(agg.On).As(agg.As))
			case "MIN":
				aggregateExpressions = append(aggregateExpressions, goqu.MIN(agg.On).As(agg.As))
			case "FIRST":
				aggregateExpressions = append(aggregateExpressions, goqu.FIRST(agg.On).As(agg.As))
			case "LAST":
				aggregateExpressions = append(aggregateExpressions, goqu.LAST(agg.On).As(agg.As))
			}
		}
	}
	if len(aggregateExpressions) > 0 {
		ds = ds.Select(aggregateExpressions...)
	}

	// Fields
	if len(fields) > 0 {
		ds = ds.Select(fields...)
	}

	// Where
	whereConditions := make([]interface{}, 0)
	if len(where) > 0 {
		for _, cond := range where {
			switch cond.Operator {
			case OperatorTypeEqual:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: 9999999999999999})
			case OperatorTypeNotEqual:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"neq": 9999999999999999}})
			case OperatorTypeIn:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"in": 9999999999999999}})
			case OperatorTypeNotIn:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"not_in": 9999999999999999}})
			case OperatorTypeGt:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"gt": 9999999999999999}})
			case OperatorTypeGte:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"gte": 9999999999999999}})
			case OperatorTypeLt:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"lt": 9999999999999999}})
			case OperatorTypeLte:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"lte": 9999999999999999}})
			case OperatorTypeIsNull:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"is_null": true}})
			case OperatorTypeIsNotNull:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"is_null": false}})
			}
		}
	}

	// Order By
	if orderBy != nil {
		if *orderType == OrderTypeAsc {
			ds = ds.Order(goqu.I("a").Asc())
		} else if *orderType == OrderTypeDesc {
			ds = ds.Order(goqu.I("a").Desc())
		}
	}

	// Limit
	if *limit > 0 {
		ds = ds.Limit(*limit)
	}

	// Group By
	if len(groupBy) > 0 {
		ds = ds.GroupBy(groupBy...)
	}

	// Join
	for _, j := range join {
		switch j.Function {
		case JoinTypeJoin:
			ds = ds.Join(
				goqu.T(j.Table).As(j.As),
				goqu.On(goqu.Ex{table + j.OnSource: j.Table + j.OnJoin}),
			)
		case JoinTypeFullOuter:
			ds = ds.FullOuterJoin(
				goqu.T(j.Table).As(j.As),
				goqu.On(goqu.Ex{table + j.OnSource: j.Table + j.OnJoin}),
			)
		case JoinTypeLeft:
			ds = ds.LeftJoin(
				goqu.T(j.Table).As(j.As),
				goqu.On(goqu.Ex{table + j.OnSource: j.Table + j.OnJoin}),
			)
		case JoinTypeRight:
			ds = ds.RightJoin(
				goqu.T(j.Table).As(j.As),
				goqu.On(goqu.Ex{table + j.OnSource: j.Table + j.OnJoin}),
			)
		case JoinTypeInner:
			ds = ds.InnerJoin(
				goqu.T(j.Table).As(j.As),
				goqu.On(goqu.Ex{table + j.OnSource: j.Table + j.OnJoin}),
			)
		case JoinTypeRightOuter:
			ds = ds.RightOuterJoin(
				goqu.T(j.Table).As(j.As),
				goqu.On(goqu.Ex{table + j.OnSource: j.Table + j.OnJoin}),
			)
		case JoinTypeLeftOuter:
			ds = ds.LeftOuterJoin(
				goqu.T(j.Table).As(j.As),
				goqu.On(goqu.Ex{table + j.OnSource: j.Table + j.OnJoin}),
			)
		case JoinTypeFull:
			ds = ds.FullJoin(
				goqu.T(j.Table).As(j.As),
				goqu.On(goqu.Ex{table + j.OnSource: j.Table + j.OnJoin}),
			)
		case JoinTypeNatural:
			ds = ds.NaturalJoin(
				goqu.T(j.Table).As(j.As),
			)
		case JoinTypeNaturalLeft:
			ds = ds.NaturalLeftJoin(
				goqu.T(j.Table).As(j.As),
			)
		case JoinTypeNaturalRight:
			ds = ds.NaturalRightJoin(
				goqu.T(j.Table).As(j.As),
			)
		case JoinTypeNaturalFull:
			ds = ds.NaturalFullJoin(
				goqu.T(j.Table).As(j.As),
			)
		case JoinTypeCross:
			ds = ds.CrossJoin(
				goqu.T(j.Table).As(j.As),
			)
		}
	}

	// Build
	query, _, _ := ds.ToSQL()

	// Replace 9999999999999999 with "?"
	strings.ReplaceAll(query, "'9999999999999999'", "?")
	strings.ReplaceAll(query, "9999999999999999", "?")

	return query
}

// BuildUpdateQuery Building a query to update a table.
func BuildUpdateQuery(
	table string,
	fields []interface{},
	where []WhereCondition,
) string {
	ds := goqu.Update(table)

	// Set
	setRecords := make(goqu.Record, 0)
	for _, f := range fields {
		setRecords[f.(string)] = 9999999999999999
	}
	ds = ds.Set(setRecords)

	// Where
	whereConditions := make([]interface{}, 0)
	if len(where) > 0 {
		for _, cond := range where {
			switch cond.Operator {
			case OperatorTypeEqual:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: 9999999999999999})
			case OperatorTypeNotEqual:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"neq": 9999999999999999}})
			case OperatorTypeIn:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"in": 9999999999999999}})
			case OperatorTypeNotIn:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"not_in": 9999999999999999}})
			case OperatorTypeGt:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"gt": 9999999999999999}})
			case OperatorTypeGte:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"gte": 9999999999999999}})
			case OperatorTypeLt:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"lt": 9999999999999999}})
			case OperatorTypeLte:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"lte": 9999999999999999}})
			case OperatorTypeIsNull:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"is_null": true}})
			case OperatorTypeIsNotNull:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"is_null": false}})
			}
		}
	}

	// Build
	query, _, _ := ds.ToSQL()

	// Replace 9999999999999999 with "?"
	strings.ReplaceAll(query, "'9999999999999999'", "?")
	strings.ReplaceAll(query, "9999999999999999", "?")

	return query
}

// BuildInsertQuery build insert query
func BuildInsertQuery(
	table string,
	fields []interface{},
	where []WhereCondition,
) string {
	ds := goqu.Insert(table)

	// Set
	setRecords := make(goqu.Record, 0)
	for _, f := range fields {
		setRecords[f.(string)] = ":" + f.(string)
	}
	ds = ds.Rows(setRecords)

	// Where
	whereConditions := make([]interface{}, 0)
	if len(where) > 0 {
		for _, cond := range where {
			switch cond.Operator {
			case OperatorTypeEqual:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: 9999999999999999})
			case OperatorTypeNotEqual:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"neq": 9999999999999999}})
			case OperatorTypeIn:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"in": 9999999999999999}})
			case OperatorTypeNotIn:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"not_in": 9999999999999999}})
			case OperatorTypeGt:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"gt": 9999999999999999}})
			case OperatorTypeGte:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"gte": 9999999999999999}})
			case OperatorTypeLt:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"lt": 9999999999999999}})
			case OperatorTypeLte:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"lte": 9999999999999999}})
			case OperatorTypeIsNull:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"is_null": true}})
			case OperatorTypeIsNotNull:
				whereConditions = append(whereConditions, goqu.Ex{cond.Column: goqu.Op{"is_null": false}})
			}
		}
	}

	// Build
	query, _, _ := ds.ToSQL()

	// Replace 9999999999999999 with "?"
	strings.ReplaceAll(query, "'9999999999999999'", "?")
	strings.ReplaceAll(query, "9999999999999999", "?")

	return query
}
