package condition

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero-contrib/castx"
	"github.com/spf13/cast"
	"strings"
)

type Operators string

const (
	Equal            Operators = "="
	NotEqual         Operators = "!="
	GreaterThan      Operators = ">"
	LessThan         Operators = "<"
	GreaterEqualThan Operators = ">="
	LessEqualThan    Operators = "<="
	In               Operators = "IN"
	NotIn            Operators = "NOT IN"
	Like             Operators = "LIKE"
	NotLike          Operators = "NOT LIKE"
	Limit            Operators = "LIMIT"
	Offset           Operators = "OFFSET"
	Between          Operators = "BETWEEN"
	OrderBy          Operators = "ORDER BY"
)

type Condition struct {
	Skip bool

	Field    string
	Operator string
	Value    any
}

func New(conditions ...Condition) []Condition {
	return conditions
}

func Apply(sb *sqlbuilder.SelectBuilder, conditions ...Condition) {
	for _, cond := range conditions {
		if cond.Skip {
			continue
		}
		switch Operators(strings.ToUpper(cond.Operator)) {
		case Equal:
			sb.Where(sb.Equal(cond.Field, cond.Value))
		case NotEqual:
			sb.Where(sb.NotEqual(cond.Field, cond.Value))
		case GreaterThan:
			sb.Where(sb.GreaterThan(cond.Field, cond.Value))
		case LessThan:
			sb.Where(sb.LessThan(cond.Field, cond.Value))
		case GreaterEqualThan:
			sb.Where(sb.GreaterEqualThan(cond.Field, cond.Value))
		case LessEqualThan:
			sb.Where(sb.LessEqualThan(cond.Field, cond.Value))
		case In:
			if len(castx.ToSlice(cond.Value)) > 0 {
				sb.Where(sb.In(cond.Field, castx.ToSlice(cond.Value)...))
			}
		case NotIn:
			if len(castx.ToSlice(cond.Value)) > 0 {
				sb.Where(sb.NotIn(cond.Field, castx.ToSlice(cond.Value)...))
			}
		case Like:
			sb.Where(sb.Like(cond.Field, cond.Value))
		case NotLike:
			sb.Where(sb.NotLike(cond.Field, cond.Value))
		case Limit:
			sb.Limit(cast.ToInt(cond.Value))
		case Offset:
			sb.Offset(cast.ToInt(cond.Value))
		case Between:
			value := castx.ToSlice(cond.Value)
			if len(value) == 2 {
				sb.Where(sb.Between(cond.Field, value[0], value[1]))
			}
		case OrderBy:
			if len(castx.ToSlice(cond.Value)) > 0 {
				sb.OrderBy(cast.ToStringSlice(castx.ToSlice(cond.Value))...)
			}
		}
	}
}
