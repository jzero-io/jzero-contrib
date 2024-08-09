package castx

import (
	"reflect"
)

func ToSlice(i interface{}) []interface{} {
	if i == nil {
		return []interface{}{}
	}

	switch v := i.(type) {
	case []interface{}:
		return v
	}

	kind := reflect.TypeOf(i).Kind()
	switch kind {
	case reflect.Slice, reflect.Array:
		s := reflect.ValueOf(i)
		a := make([]interface{}, s.Len())
		for i := range a {
			a[i] = s.Index(i).Interface()
		}
		return a
	case
		reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.String:
		return []interface{}{i}
	default:
		return []interface{}{}
	}
}
