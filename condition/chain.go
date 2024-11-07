package condition

type Chain struct {
	conditions []Condition
}

func NewChain() Chain {
	return Chain{}
}

func NewChainWithConditions(conditions ...Condition) Chain {
	return Chain{conditions: conditions}
}

func (c Chain) Equal(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: Equal,
		Value:    value,
	})
	return c
}

func (c Chain) NotEqual(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: NotEqual,
		Value:    value,
	})
	return c
}

func (c Chain) GreaterThan(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: GreaterThan,
		Value:    value,
	})
	return c
}

func (c Chain) LessThan(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: LessThan,
		Value:    value,
	})
	return c
}

func (c Chain) GreaterEqualThan(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: GreaterEqualThan,
		Value:    value,
	})
	return c
}

func (c Chain) LessEqualThan(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: LessEqualThan,
		Value:    value,
	})
	return c
}

func (c Chain) Like(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: Like,
		Value:    value,
	})
	return c
}

func (c Chain) NotLike(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: NotLike,
		Value:    value,
	})
	return c
}

func (c Chain) In(field string, values any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: In,
		Value:    values,
	})
	return c
}

func (c Chain) NotIn(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: NotIn,
		Value:    value,
	})
	return c
}

func (c Chain) Between(field string, value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Field:    field,
		Operator: Between,
		Value:    value,
	})
	return c
}

func (c Chain) Or(fields []string, values []any) Chain {
	c.conditions = append(c.conditions, Condition{
		Or:       true,
		OrFields: fields,
		OrValues: values,
	})
	return c
}

func (c Chain) OrderBy(value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Operator: OrderBy,
		Value:    value,
	})
	return c
}

func (c Chain) Limit(value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Operator: Limit,
		Value:    value,
	})
	return c
}

func (c Chain) Offset(value any) Chain {
	c.conditions = append(c.conditions, Condition{
		Operator: Offset,
		Value:    value,
	})
	return c
}

func (c Chain) Page(page, pageSize int) Chain {
	c.conditions = append(c.conditions, Condition{
		Operator: Offset,
		Value:    (page - 1) * pageSize,
	})
	c.conditions = append(c.conditions, Condition{
		Operator: Limit,
		Value:    pageSize,
	})
	return c
}

func (c Chain) ToCondition() []Condition {
	return c.conditions
}
