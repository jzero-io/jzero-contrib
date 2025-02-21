package condition

import (
	"strconv"
	"strings"

	"github.com/huandu/go-sqlbuilder"
)

func SelectByWhereRawSql(sb *sqlbuilder.SelectBuilder, originalField string, args ...any) {
	originalFields := strings.Split(originalField, " and ")
	for i, v := range originalFields {
		field := strings.Split(v, " = ")
		unquote, _ := strconv.Unquote(field[0])
		sb.Where(sb.EQ(unquote, args[i]))
	}
}
