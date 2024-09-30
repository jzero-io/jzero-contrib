package condition

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestSelectWithCondition(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

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

func TestUpdateWithCondition(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

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

	sb := sqlbuilder.NewUpdateBuilder().Update("user")
	ApplyUpdate(sb, cds...)
	sb.Set(sb.Equal("name", "gocloudcoder"))

	sql, args := sb.Build()
	fmt.Println(sql)
	fmt.Println(args)
}

func TestDeleteWithCondition(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

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

	sb := sqlbuilder.NewDeleteBuilder().DeleteFrom("user")
	ApplyDelete(sb, cds...)

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
