package condition

import (
	"fmt"
	"testing"

	"github.com/huandu/go-sqlbuilder"
)

func TestChain(t *testing.T) {
	sb := sqlbuilder.NewSelectBuilder().Select("name", "age").From("user")

	chain := NewChain()
	conds := chain.
		Equal("field1", "value1", WithSkip(true)).
		Equal("field2", "value2").
		OrderBy("create_time desc").
		OrderBy("sort desc").
		ToCondition()
	ApplySelect(sb, conds...)

	sql, args := sb.Build()
	fmt.Println(sql)
	fmt.Println(args)
}
