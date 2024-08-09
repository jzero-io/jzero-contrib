package condition

import (
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"testing"
)

func TestCondition(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL

	cds := New(Condition{
		Field:    "field1",
		Operator: Equal,
		Value:    "value1",
	}, Condition{
		Field:    "field2",
		Operator: "=",
		Value:    "value2",
	}, Condition{
		Field:    "field3",
		Operator: "IN",
		Value:    nil,
	})

	sb := sqlbuilder.NewSelectBuilder()
	Apply(sb, cds...)

	sql, args := sb.Build()
	fmt.Println(sql)
	fmt.Println(args)
}
