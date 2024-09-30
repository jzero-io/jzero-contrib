package condition

import (
	"fmt"
	"github.com/huandu/go-sqlbuilder"
	"testing"
)

func TestCondition(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.PostgreSQL

	var values []any
	values = append(values, []int{24, 48}, []int{170, 175})

	cds := New(Condition{
		Field:    "name",
		Operator: Equal,
		Value:    "jaronnie",
	}, Condition{
		Or:          true,
		OrFields:    []string{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
	})

	sb := sqlbuilder.NewSelectBuilder().Select("name", "age", "height").From("user")
	Apply(sb, cds...)

	sql, args := sb.Build()
	fmt.Println(sql)
	fmt.Println(args)
}

func TestSqlBuilder(t *testing.T) {
	builder := sqlbuilder.NewSelectBuilder().Select("id", "name").From("user")
	builder.Where(builder.Or(builder.Equal("id", 1), builder.Equal("id", "2")))
	builder.Where(builder.And(builder.Equal("name", "jaronnie")))
	fmt.Println(builder.Build())
}
