package condition

import (
	"github.com/eddieowens/opts"
)

type Chain struct {
	conditions []Condition
}

type ChainOperatorOpts struct {
	Skip     bool
	SkipFunc func() bool

	ValueFunc func() any

	OrValuesFunc func() []any
}

func (opts ChainOperatorOpts) DefaultOptions() ChainOperatorOpts {
	return ChainOperatorOpts{}
}

func WithSkip(skip bool) opts.Opt[ChainOperatorOpts] {
	return func(c *ChainOperatorOpts) {
		c.Skip = skip
	}
}

func WithSkipFunc(skipFunc func() bool) opts.Opt[ChainOperatorOpts] {
	return func(c *ChainOperatorOpts) {
		c.SkipFunc = skipFunc
	}
}

func WithValue(value any) opts.Opt[ChainOperatorOpts] {
	return func(c *ChainOperatorOpts) {
		c.ValueFunc = func() any {
			return value
		}
	}
}

func WithOrValues(orValues []any) opts.Opt[ChainOperatorOpts] {
	return func(c *ChainOperatorOpts) {
		c.OrValuesFunc = func() []any {
			return orValues
		}
	}
}

func NewChain() Chain {
	return Chain{}
}

func NewChainWithConditions(conditions ...Condition) Chain {
	return Chain{conditions: conditions}
}

func (c Chain) Equal(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  Equal,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) NotEqual(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  NotEqual,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) GreaterThan(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  GreaterThan,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) LessThan(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  LessThan,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) GreaterEqualThan(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  GreaterEqualThan,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) LessEqualThan(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  LessEqualThan,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) Like(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  Like,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) NotLike(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  NotLike,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) In(field string, values any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  In,
		Value:     values,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) NotIn(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  NotIn,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) Between(field string, value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Field:     field,
		Operator:  Between,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) Or(fields []string, values []any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Or:           true,
		OrFields:     fields,
		OrValues:     values,
		Skip:         o.Skip,
		SkipFunc:     o.SkipFunc,
		OrValuesFunc: o.OrValuesFunc,
	})
	return c
}

func (c Chain) OrderBy(value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Operator:  OrderBy,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) Limit(value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Operator:  Limit,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) Offset(value any, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Operator:  Offset,
		Value:     value,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) Page(page, pageSize int, op ...opts.Opt[ChainOperatorOpts]) Chain {
	o := opts.DefaultApply(op...)
	c.conditions = append(c.conditions, Condition{
		Operator:  Offset,
		Value:     (page - 1) * pageSize,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	c.conditions = append(c.conditions, Condition{
		Operator:  Limit,
		Value:     pageSize,
		Skip:      o.Skip,
		SkipFunc:  o.SkipFunc,
		ValueFunc: o.ValueFunc,
	})
	return c
}

func (c Chain) ToCondition() []Condition {
	return c.conditions
}
