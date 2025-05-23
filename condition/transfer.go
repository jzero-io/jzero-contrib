package condition

import (
	"fmt"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

func SelectByWhereRawSql(sb *sqlbuilder.SelectBuilder, originalField string, args ...any) {
	originalFields := strings.Split(originalField, " and ")
	for i, v := range originalFields {
		field := strings.Split(v, " = ")[0]
		if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL {
			field = util.Unquote(field)
			field = fmt.Sprintf(`"%s"`, field)
		}
		sb.Where(sb.EQ(field, args[i]))
	}
}
