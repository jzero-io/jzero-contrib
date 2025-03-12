package condition

import (
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/spf13/cast"

	"github.com/jzero-io/jzero-contrib/castx"
)

type Operator string

func (o Operator) String() string {
	return string(o)
}

const (
	Equal            Operator = "="
	NotEqual         Operator = "!="
	GreaterThan      Operator = ">"
	LessThan         Operator = "<"
	GreaterEqualThan Operator = ">="
	LessEqualThan    Operator = "<="
	In               Operator = "IN"
	NotIn            Operator = "NOT IN"
	Like             Operator = "LIKE"
	NotLike          Operator = "NOT LIKE"
	Limit            Operator = "LIMIT"
	Offset           Operator = "OFFSET"
	Between          Operator = "BETWEEN"
	NotBetween       Operator = "NOT BETWEEN"
	OrderBy          Operator = "ORDER BY"
	GroupBy          Operator = "GROUP BY"
	Join             Operator = "JOIN"
)

type Condition struct {
	// Skip indicates whether the condition is effective.
	Skip bool

	// SkipFunc The priority is higher than Skip.
	SkipFunc func() bool

	// Or indicates an or condition
	Or bool

	OrOperators  []Operator
	OrFields     []string
	OrValues     []any
	OrValuesFunc func() []any

	// Field for default and condition
	Field string

	Operator Operator
	Value    any

	// ValueFunc The priority is higher than Value.
	ValueFunc func() any

	// JoinCondition
	JoinCondition

	WhereClause *sqlbuilder.WhereClause
}

type JoinCondition struct {
	Option sqlbuilder.JoinOption
	Table  string
	OnExpr []string
}

func New(conditions ...Condition) []Condition {
	return conditions
}

func buildExpr(cond *sqlbuilder.Cond, field string, operator Operator, value any) string {
	switch operator {
	case Equal:
		return cond.Equal(field, value)
	case NotEqual:
		return cond.NotEqual(field, value)
	case GreaterThan:
		return cond.GreaterThan(field, value)
	case LessThan:
		return cond.LessThan(field, value)
	case GreaterEqualThan:
		return cond.GreaterEqualThan(field, value)
	case LessEqualThan:
		return cond.LessEqualThan(field, value)
	case In:
		if len(castx.ToSlice(value)) == 0 {
			// if value is empty, force placeholder nil to avoid sql error
			return cond.In(field, nil)
		}
		return cond.In(field, castx.ToSlice(value)...)
	case NotIn:
		if len(castx.ToSlice(value)) == 0 {
			// if value is empty, force placeholder nil to avoid sql error
			return cond.NotIn(field, nil)
		}
		return cond.NotIn(field, castx.ToSlice(value)...)
	case Like:
		return cond.Like(field, value)
	case NotLike:
		return cond.NotLike(field, value)
	case Between:
		v := castx.ToSlice(value)
		return cond.Between(field, v[0], v[1])
	case NotBetween:
		v := castx.ToSlice(value)
		return cond.NotBetween(field, v[0], v[1])
	}
	return ""
}

func whereClause(conditions ...Condition) *sqlbuilder.WhereClause {
	clause := sqlbuilder.NewWhereClause()
	cond := sqlbuilder.NewCond()

	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.WhereClause != nil {
			clause.AddWhereClause(c.WhereClause)
			continue
		}
		if c.Or {
			if c.OrValuesFunc != nil {
				c.OrValues = c.OrValuesFunc()
			}
			var expr []string
			for i, field := range c.OrFields {
				if or := buildExpr(cond, field, c.OrOperators[i], c.OrValues[i]); or != "" {
					expr = append(expr, or)
				}
			}
			if len(expr) > 0 {
				clause.AddWhereExpr(cond.Args, cond.Or(expr...))
			}
		} else {
			if c.ValueFunc != nil {
				c.Value = c.ValueFunc()
			}
			if and := buildExpr(cond, c.Field, c.Operator, c.Value); and != "" {
				clause.AddWhereExpr(cond.Args, and)
			}
		}
	}
	return clause
}

func Select(sb sqlbuilder.SelectBuilder, conditions ...Condition) sqlbuilder.SelectBuilder {
	clause := whereClause(conditions...)
	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.ValueFunc != nil {
			c.Value = c.ValueFunc()
		}
		switch Operator(strings.ToUpper(string(c.Operator))) {
		case Limit:
			sb.Limit(cast.ToInt(c.Value))
		case Offset:
			sb.Offset(cast.ToInt(c.Value))
		case OrderBy:
			sb.OrderBy(cast.ToStringSlice(castx.ToSlice(c.Value))...)
		case GroupBy:
			sb.GroupBy(cast.ToStringSlice(castx.ToSlice(c.Value))...)
		case Join:
			sb.JoinWithOption(c.JoinCondition.Option, c.JoinCondition.Table, cast.ToStringSlice(castx.ToSlice(c.JoinCondition.OnExpr))...)
		}
	}
	if clause != nil {
		sb = *sb.AddWhereClause(clause)
	}
	return sb
}

func Update(builder sqlbuilder.UpdateBuilder, conditions ...Condition) sqlbuilder.UpdateBuilder {
	clause := whereClause(conditions...)
	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.ValueFunc != nil {
			c.Value = c.ValueFunc()
		}
		switch Operator(strings.ToUpper(string(c.Operator))) {
		case Limit:
			builder.Limit(cast.ToInt(c.Value))
		case OrderBy:
			if len(castx.ToSlice(c.Value)) > 0 {
				builder.OrderBy(cast.ToStringSlice(castx.ToSlice(c.Value))...)
			}
		}
	}
	if clause != nil {
		builder = *builder.AddWhereClause(clause)
	}
	return builder
}

func Delete(builder sqlbuilder.DeleteBuilder, conditions ...Condition) sqlbuilder.DeleteBuilder {
	clause := whereClause(conditions...)
	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.ValueFunc != nil {
			c.Value = c.ValueFunc()
		}
		switch Operator(strings.ToUpper(string(c.Operator))) {
		case Limit:
			builder.Limit(cast.ToInt(c.Value))
		case OrderBy:
			builder.OrderBy(cast.ToStringSlice(castx.ToSlice(c.Value))...)
		}
	}
	if clause != nil {
		builder = *builder.AddWhereClause(clause)
	}
	return builder
}
