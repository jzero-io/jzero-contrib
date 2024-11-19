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
)

type Condition struct {
	Skip     bool
	SkipFunc func() bool

	// or condition
	Or bool

	OrOperators  []Operator
	OrFields     []string
	OrValues     []any
	OrValuesFunc func() []any

	// and condition
	Field string

	Operator  Operator
	Value     any
	ValueFunc func() any
}

func New(conditions ...Condition) []Condition {
	return conditions
}

func buildWhereClause(conditions ...Condition) *sqlbuilder.WhereClause {
	clause := sqlbuilder.NewWhereClause()
	cond := sqlbuilder.NewCond()

	for _, c := range conditions {
		if c.SkipFunc != nil {
			c.Skip = c.SkipFunc()
		}
		if c.Skip {
			continue
		}
		if c.Or {
			if c.OrValuesFunc != nil {
				c.OrValues = c.OrValuesFunc()
			}
			var expr []string
			for i, field := range c.OrFields {
				switch Operator(strings.ToUpper(string(c.OrOperators[i]))) {
				case Equal:
					expr = append(expr, cond.Equal(field, c.OrValues[i]))
				case NotEqual:
					expr = append(expr, cond.NotEqual(field, c.OrValues[i]))
				case GreaterThan:
					expr = append(expr, cond.GreaterThan(field, c.OrValues[i]))
				case LessThan:
					expr = append(expr, cond.LessThan(field, c.OrValues[i]))
				case GreaterEqualThan:
					expr = append(expr, cond.GreaterEqualThan(field, c.OrValues[i]))
				case LessEqualThan:
					expr = append(expr, cond.LessEqualThan(field, c.OrValues[i]))
				case In:
					if len(castx.ToSlice(c.OrValues[i])) > 0 {
						expr = append(expr, cond.In(field, castx.ToSlice(c.OrValues[i])...))
					}
				case NotIn:
					if len(castx.ToSlice(c.OrValues[i])) > 0 {
						expr = append(expr, cond.NotIn(field, castx.ToSlice(c.OrValues[i])...))
					}
				case Like:
					expr = append(expr, cond.Like(field, c.OrValues[i]))
				case NotLike:
					expr = append(expr, cond.NotLike(field, c.OrValues[i]))
				case Between:
					value := castx.ToSlice(c.OrValues[i])
					if len(value) == 2 {
						expr = append(expr, cond.Between(field, value[0], value[1]))
					}
				case NotBetween:
					value := castx.ToSlice(c.OrValues[i])
					if len(value) == 2 {
						expr = append(expr, cond.NotBetween(field, value[0], value[1]))
					}
				}
			}
			if len(expr) > 0 {
				clause.AddWhereExpr(cond.Args, cond.Or(expr...))
			}
		} else {
			if c.ValueFunc != nil {
				c.Value = c.ValueFunc()
			}
			switch Operator(strings.ToUpper(string(c.Operator))) {
			case Equal:
				clause.AddWhereExpr(cond.Args, cond.Equal(c.Field, c.Value))
			case NotEqual:
				clause.AddWhereExpr(cond.Args, cond.NotEqual(c.Field, c.Value))
			case GreaterThan:
				clause.AddWhereExpr(cond.Args, cond.GreaterThan(c.Field, c.Value))
			case LessThan:
				clause.AddWhereExpr(cond.Args, cond.LessThan(c.Field, c.Value))
			case GreaterEqualThan:
				clause.AddWhereExpr(cond.Args, cond.GreaterThan(c.Field, c.Value))
			case LessEqualThan:
				clause.AddWhereExpr(cond.Args, cond.LessThan(c.Field, c.Value))
			case In:
				if len(castx.ToSlice(c.Value)) > 0 {
					clause.AddWhereExpr(cond.Args, cond.In(c.Field, castx.ToSlice(c.Value)...))
				}
			case NotIn:
				if len(castx.ToSlice(c.Value)) > 0 {
					clause.AddWhereExpr(cond.Args, cond.NotIn(c.Field, castx.ToSlice(c.Value)...))
				}
			case Like:
				clause.AddWhereExpr(cond.Args, cond.Like(c.Field, c.Value))
			case NotLike:
				clause.AddWhereExpr(cond.Args, cond.NotLike(c.Field, c.Value))
			case Between:
				value := castx.ToSlice(c.Value)
				if len(value) == 2 {
					clause.AddWhereExpr(cond.Args, cond.Between(c.Field, value[0], value[1]))
				}
			case NotBetween:
				value := castx.ToSlice(c.Value)
				if len(value) == 2 {
					clause.AddWhereExpr(cond.Args, cond.NotBetween(c.Field, value[0], value[1]))
				}
			}
		}
	}
	return clause
}

func ApplySelect(sb *sqlbuilder.SelectBuilder, conditions ...Condition) {
	clause := buildWhereClause(conditions...)
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
			if len(castx.ToSlice(c.Value)) > 0 {
				sb.OrderBy(cast.ToStringSlice(castx.ToSlice(c.Value))...)
			}
		}
	}
	if clause != nil {
		sb = sb.AddWhereClause(clause)
	}
}

func ApplyUpdate(sb *sqlbuilder.UpdateBuilder, conditions ...Condition) {
	clause := buildWhereClause(conditions...)
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
		case OrderBy:
			if len(castx.ToSlice(c.Value)) > 0 {
				sb.OrderBy(cast.ToStringSlice(castx.ToSlice(c.Value))...)
			}
		}
	}
	if clause != nil {
		sb = sb.AddWhereClause(clause)
	}
}

func ApplyDelete(sb *sqlbuilder.DeleteBuilder, conditions ...Condition) {
	clause := buildWhereClause(conditions...)
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
		case OrderBy:
			if len(castx.ToSlice(c.Value)) > 0 {
				sb.OrderBy(cast.ToStringSlice(castx.ToSlice(c.Value))...)
			}
		}
	}
	if clause != nil {
		sb = sb.AddWhereClause(clause)
	}
}
