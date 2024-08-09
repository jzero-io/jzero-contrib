package condition

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero-contrib/castx"
	"github.com/spf13/cast"
	"strings"
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
		switch strings.ToUpper(cond.Operator) {
		case "=":
			sb.Where(sb.Equal(cond.Field, cond.Value))
		case "!=":
			sb.Where(sb.NotEqual(cond.Field, cond.Value))
		case ">":
			sb.Where(sb.GreaterThan(cond.Field, cond.Value))
		case "<":
			sb.Where(sb.LessThan(cond.Field, cond.Value))
		case ">=":
			sb.Where(sb.GreaterEqualThan(cond.Field, cond.Value))
		case "<=":
			sb.Where(sb.LessEqualThan(cond.Field, cond.Value))
		case "IN":
			if len(castx.ToSlice(cond.Value)) > 0 {
				sb.Where(sb.In(cond.Field, castx.ToSlice(cond.Value)...))
			}
		case "NOT IN":
			if len(castx.ToSlice(cond.Value)) > 0 {
				sb.Where(sb.NotIn(cond.Field, castx.ToSlice(cond.Value)...))
			}
		case "LIKE":
			sb.Where(sb.Like(cond.Field, cond.Value))
		case "NOT LIKE":
			sb.Where(sb.NotLike(cond.Field, cond.Value))
		case "LIMIT":
			sb.Limit(cast.ToInt(cond.Value))
		case "OFFSET":
			sb.Offset(cast.ToInt(cond.Value))
		case "BETWEEN":
			value := castx.ToSlice(cond.Value)
			if len(value) == 2 {
				sb.Where(sb.Between(cond.Field, value[0], value[1]))
			}
		case "ORDER BY":
			if len(castx.ToSlice(cond.Value)) > 0 {
				sb.OrderBy(cast.ToStringSlice(castx.ToSlice(cond.Value))...)
			}
		}
	}
}
