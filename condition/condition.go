package condition

import (
	"github.com/huandu/go-sqlbuilder"
	"github.com/jzero-io/jzero-contrib/castx"
	"github.com/spf13/cast"
	"strings"
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
	OrderBy          Operator = "ORDER BY"
)

type Condition struct {
	Skip bool

	Field    string
	Operator Operator
	Value    any
}

func New(conditions ...Condition) []Condition {
	return conditions
}

func ApplySelect(sb *sqlbuilder.SelectBuilder, conditions ...Condition) {
	for _, cond := range conditions {
		if cond.Skip {
			continue
		}
		switch Operator(strings.ToUpper(string(cond.Operator))) {
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

func ApplyUpdate(sb *sqlbuilder.UpdateBuilder, conditions ...Condition) {
	for _, cond := range conditions {
		if cond.Skip {
			continue
		}
		switch Operator(strings.ToUpper(string(cond.Operator))) {
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

func ApplyDelete(sb *sqlbuilder.DeleteBuilder, conditions ...Condition) {
	for _, cond := range conditions {
		if cond.Skip {
			continue
		}
		switch Operator(strings.ToUpper(string(cond.Operator))) {
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
