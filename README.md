# jzero-contrib

jzero contrib

## swaggerv2

在线展示 swagger ui 文档

![](https://oss.jaronnie.com/image-20240627175804999.png)

### Usage

将 swagger.json 放在 docs 文件夹下

```go
package main

import (
	"github.com/jzero-io/jzero-contrib/swaggerv2"
	"github.com/zeromicro/go-zero/rest"
)

func main() {
	server := rest.MustNewServer(rest.RestConf{
		Port: 8001,
	})
	swaggerv2.RegisterRoutes(server, swaggerv2.WithSwaggerPath("docs"))

	server.Start()
}
```

访问 localhost:8001/swagger

## condition

查询/更新/删除 条件构建器

```go
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
```