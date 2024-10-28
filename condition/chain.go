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

func (c Chain) ToCondition() []Condition {
	return c.conditions
}
