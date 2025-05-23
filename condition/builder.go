package condition

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/huandu/go-sqlbuilder"
	"github.com/zeromicro/go-zero/tools/goctl/util"
)

const dbTag = "db"

// RawFieldNames converts golang struct field into slice string.
func RawFieldNames(in any) []string {
	out := make([]string, 0)
	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		tagv := fi.Tag.Get(dbTag)
		switch tagv {
		case "-":
			continue
		case "":
			if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL {
				out = append(out, fmt.Sprintf(`"%s"`, tagv))
			} else {
				out = append(out, fmt.Sprintf("`%s`", fi.Name))
			}
		default:
			if strings.Contains(tagv, ",") {
				tagv = strings.TrimSpace(strings.Split(tagv, ",")[0])
			}
			if tagv == "-" {
				continue
			}
			if len(tagv) == 0 {
				tagv = fi.Name
			}
			if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL {
				out = append(out, fmt.Sprintf(`"%s"`, tagv))
			} else {
				out = append(out, fmt.Sprintf("`%s`", fi.Name))
			}
		}
	}

	return out
}

func RemoveIgnoreColumns(strings []string, strs ...string) []string {
	out := append([]string(nil), strings...)

	for _, str := range strs {
		if sqlbuilder.DefaultFlavor == sqlbuilder.PostgreSQL {
			str = fmt.Sprintf(`"%s"`, util.Unquote(str))
		}

		var n int
		for _, v := range out {
			if v != str {
				out[n] = v
				n++
			}
		}
		out = out[:n]
	}

	return out
}

func Table(table string) string {
	return format(table)
}

func Field(field string) string {
	return format(field)
}

func format(str string) string {
	switch sqlbuilder.DefaultFlavor {
	case sqlbuilder.PostgreSQL:
		str = util.Unquote(str)
		return fmt.Sprintf(`"%s"`, str)
	default:
		return str
	}
}
