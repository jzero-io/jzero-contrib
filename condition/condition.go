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
	OrderBy          Operator = "ORDER BY"
)

type Condition struct {
	Skip     bool
	SkipFunc func() bool

	// or condition
	Or           bool
	OrOperators  []Operator
	OrFields     []string
	OrValues     []any
	OrValuesFunc func() []any

	// and condition
	Field     string
	Operator  Operator
	Value     any
	ValueFunc func() any
}

func New(conditions ...Condition) []Condition {
	return conditions
}

func Apply(sb *sqlbuilder.SelectBuilder, conditions ...Condition) {
	ApplySelect(sb, conditions...)
}

func ApplySelect(sb *sqlbuilder.SelectBuilder, conditions ...Condition) {
	for _, cond := range conditions {
		if cond.SkipFunc != nil {
			cond.Skip = cond.SkipFunc()
		}
		if cond.Skip {
			continue
		}
		if cond.Or {
			if cond.OrValuesFunc != nil {
				cond.OrValues = cond.OrValuesFunc()
			}
			var expr []string
			for i, field := range cond.OrFields {
				switch Operator(strings.ToUpper(string(cond.OrOperators[i]))) {
				case Equal:
					expr = append(expr, sb.Equal(field, cond.OrValues[i]))
				case NotEqual:
					expr = append(expr, sb.NotEqual(field, cond.OrValues[i]))
				case GreaterThan:
					expr = append(expr, sb.GreaterThan(field, cond.OrValues[i]))
				case LessThan:
					expr = append(expr, sb.LessThan(field, cond.OrValues[i]))
				case GreaterEqualThan:
					expr = append(expr, sb.GreaterEqualThan(field, cond.OrValues[i]))
				case LessEqualThan:
					expr = append(expr, sb.LessEqualThan(field, cond.OrValues[i]))
				case In:
					if len(castx.ToSlice(cond.OrValues[i])) > 0 {
						expr = append(expr, sb.In(field, castx.ToSlice(cond.OrValues[i])...))
					}
				case NotIn:
					if len(castx.ToSlice(cond.OrValues[i])) > 0 {
						expr = append(expr, sb.NotIn(field, castx.ToSlice(cond.OrValues[i])...))
					}
				case Like:
					expr = append(expr, sb.Like(field, cond.OrValues[i]))
				case NotLike:
					expr = append(expr, sb.NotLike(field, cond.OrValues[i]))
				case Between:
					value := castx.ToSlice(cond.OrValues[i])
					if len(value) == 2 {
						expr = append(expr, sb.Between(field, value[0], value[1]))
					}
				}
			}
			sb.Where(sb.Or(expr...))
		} else {
			if cond.ValueFunc != nil {
				cond.Value = cond.ValueFunc()
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
}

func ApplyUpdate(sb *sqlbuilder.UpdateBuilder, conditions ...Condition) {
	for _, cond := range conditions {
		if cond.SkipFunc != nil {
			cond.Skip = cond.SkipFunc()
		}
		if cond.Skip {
			continue
		}
		if cond.Or {
			if cond.OrValuesFunc != nil {
				cond.OrValues = cond.OrValuesFunc()
			}
			var expr []string
			for i, field := range cond.OrFields {
				switch Operator(strings.ToUpper(string(cond.OrOperators[i]))) {
				case Equal:
					expr = append(expr, sb.Equal(field, cond.OrValues[i]))
				case NotEqual:
					expr = append(expr, sb.NotEqual(field, cond.OrValues[i]))
				case GreaterThan:
					expr = append(expr, sb.GreaterThan(field, cond.OrValues[i]))
				case LessThan:
					expr = append(expr, sb.LessThan(field, cond.OrValues[i]))
				case GreaterEqualThan:
					expr = append(expr, sb.GreaterEqualThan(field, cond.OrValues[i]))
				case LessEqualThan:
					expr = append(expr, sb.LessEqualThan(field, cond.OrValues[i]))
				case In:
					if len(castx.ToSlice(cond.OrValues[i])) > 0 {
						expr = append(expr, sb.In(field, castx.ToSlice(cond.OrValues[i])...))
					}
				case NotIn:
					if len(castx.ToSlice(cond.OrValues[i])) > 0 {
						expr = append(expr, sb.NotIn(field, castx.ToSlice(cond.OrValues[i])...))
					}
				case Like:
					expr = append(expr, sb.Like(field, cond.OrValues[i]))
				case NotLike:
					expr = append(expr, sb.NotLike(field, cond.OrValues[i]))
				case Between:
					value := castx.ToSlice(cond.OrValues[i])
					if len(value) == 2 {
						expr = append(expr, sb.Between(field, value[0], value[1]))
					}
				}
			}
			sb.Where(sb.Or(expr...))
		} else {
			if cond.ValueFunc != nil {
				cond.Value = cond.ValueFunc()
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
}

func ApplyDelete(sb *sqlbuilder.DeleteBuilder, conditions ...Condition) {
	for _, cond := range conditions {
		if cond.SkipFunc != nil {
			cond.Skip = cond.SkipFunc()
		}
		if cond.Skip {
			continue
		}
		if cond.Or {
			if cond.OrValuesFunc != nil {
				cond.OrValues = cond.OrValuesFunc()
			}
			var expr []string
			for i, field := range cond.OrFields {
				switch Operator(strings.ToUpper(string(cond.OrOperators[i]))) {
				case Equal:
					expr = append(expr, sb.Equal(field, cond.OrValues[i]))
				case NotEqual:
					expr = append(expr, sb.NotEqual(field, cond.OrValues[i]))
				case GreaterThan:
					expr = append(expr, sb.GreaterThan(field, cond.OrValues[i]))
				case LessThan:
					expr = append(expr, sb.LessThan(field, cond.OrValues[i]))
				case GreaterEqualThan:
					expr = append(expr, sb.GreaterEqualThan(field, cond.OrValues[i]))
				case LessEqualThan:
					expr = append(expr, sb.LessEqualThan(field, cond.OrValues[i]))
				case In:
					if len(castx.ToSlice(cond.OrValues[i])) > 0 {
						expr = append(expr, sb.In(field, castx.ToSlice(cond.OrValues[i])...))
					}
				case NotIn:
					if len(castx.ToSlice(cond.OrValues[i])) > 0 {
						expr = append(expr, sb.NotIn(field, castx.ToSlice(cond.OrValues[i])...))
					}
				case Like:
					expr = append(expr, sb.Like(field, cond.OrValues[i]))
				case NotLike:
					expr = append(expr, sb.NotLike(field, cond.OrValues[i]))
				case Between:
					value := castx.ToSlice(cond.OrValues[i])
					if len(value) == 2 {
						expr = append(expr, sb.Between(field, value[0], value[1]))
					}
				}
			}
			sb.Where(sb.Or(expr...))
		} else {
			if cond.ValueFunc != nil {
				cond.Value = cond.ValueFunc()
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
}
