package condition

import (
	"fmt"
	"github.com/huandu/go-assert"
	"testing"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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
	builder := Select(*sb, cds...)

	sql, args := builder.Build()
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
	builder := Update(*sb, cds...)
	builder.Set(sb.Equal("name", "gocloudcoder"))

	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)

	builder = Update(*sb, Condition{
		Field:    "age",
		Operator: Equal,
		Value:    30,
	})
	sql, args = builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
}

func TestDeleteWithCondition(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

	var values []any
	values = append(values, []int{24, 48}, []int{170, 175})

	cds := New(Condition{
		SkipFunc: func() bool {
			return true
		},
		Field:    "name",
		Operator: Equal,
		Value:    "jaronnie",
		ValueFunc: func() any {
			return "jaronnie2"
		},
	}, Condition{
		Or:          true,
		OrFields:    []string{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
		OrValuesFunc: func() []any {
			return []any{[]int{24, 49}, []int{170, 176}}
		},
	})

	sb := sqlbuilder.NewDeleteBuilder().DeleteFrom("user")
	builder := Delete(*sb, cds...)

	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
}

func TestSqlBuilder(t *testing.T) {
	builder := sqlbuilder.NewSelectBuilder().Select("id", "name").From("user")
	builder.Where(builder.Or(builder.Equal("id", 1), builder.Equal("id", "2")))
	builder.Where(builder.And(builder.Equal("name", "jaronnie")))
	fmt.Println(builder.Build())
}

func TestWhereClause(t *testing.T) {
	var values []any
	values = append(values, []int{24, 48}, []int{170, 175})
	cds := New(Condition{
		SkipFunc: func() bool {
			return true
		},
		Field:    "name",
		Operator: Equal,
		Value:    "jaronnie",
		ValueFunc: func() any {
			return "jaronnie2"
		},
	}, Condition{
		Or:          true,
		OrFields:    []string{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
		OrValuesFunc: func() []any {
			return []any{[]int{24, 49}, []int{170, 176}}
		},
	})
	clause := whereClause(cds...)
	statement, args := clause.Build()
	fmt.Println(statement)
	fmt.Println(args)
}

func TestGroupBySelect(t *testing.T) {
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

	var values []any
	values = append(values, []int{24, 48}, []int{170, 175})

	cds := New(Condition{
		Field:    "name",
		Operator: Equal,
		Value:    "jaronnie",
	}, Condition{
		Field:    "money",
		Operator: GreaterEqualThan,
		Value:    100000,
	}, Condition{
		Or:          true,
		OrFields:    []string{"age", "height"},
		OrOperators: []Operator{Between, Between},
		OrValues:    values,
	}, Condition{
		Operator: GroupBy,
		Value:    []string{"class", "subject"},
	}, Condition{
		Operator: Having,
		NestedCondition: []Condition{
			{
				Field:    "classNum",
				Operator: GreaterThan,
				Value:    1,
			},
			{
				Or:          true,
				OrFields:    []string{"subjectNum", "subjectNum"},
				OrOperators: []Operator{LessThan, GreaterThan},
				OrValues:    []any{10, 20},
			},
		},
	})

	sb := sqlbuilder.NewSelectBuilder().Select("name", "age", "height", "COUNT(class) as classNum", "COUNT(subject) as subjectNum").From("user")
	builder := Select(*sb, cds...)

	sql, args := builder.Build()
	fmt.Println(sql)
	fmt.Println(args)
	assert.Equal(t, len(args), 9)
}

func TestGroupBySelectInAdmin(t *testing.T) {
	conn := sqlx.NewMysql("root:123456@tcp(localhost:3306)/jzeroadmin?charset=utf8mb4&parseTime=True&loc=Local")
	sqlbuilder.DefaultFlavor = sqlbuilder.MySQL

	cds := New(Condition{
		Field:    "menu_id",
		Operator: GreaterThan,
		Value:    5,
	}, Condition{
		Operator: GroupBy,
		Value:    []string{"role_id"},
	}, Condition{
		Operator: Having,
		NestedCondition: []Condition{
			{
				Or:          true,
				OrFields:    []string{"menuNum", "menuNum"},
				OrOperators: []Operator{LessThan, GreaterThan},
				OrValues:    []any{10, 20},
			},
		},
	})

	sb := sqlbuilder.NewSelectBuilder().Select("role_id", "COUNT(id) as menuNum").From("manage_role_menu")
	builder := Select(*sb, cds...)

	sql, args := builder.Build()
	assert.Equal(t, len(args), 3)
	fmt.Println(sql)
	fmt.Println(args)

	type Res struct {
		A int `db:"role_id"`
		B int `db:"menuNum"`
	}

	var res []Res
	err := conn.QueryRowsPartial(&res, sql, args...)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}
