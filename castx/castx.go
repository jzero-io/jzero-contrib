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
	default:
		return []interface{}{}
	}
}
