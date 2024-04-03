package builder

import (
	"reflect"
	"strings"
)

func IsExpression(obj any) (b bool) {
	rfv := reflect.Indirect(reflect.ValueOf(obj))
	if rfv.Kind() == reflect.String && strings.Contains(rfv.String(), "?") {
		b = true
	}
	return
}

func ToSlice(arg any) []any {
	ref := reflect.Indirect(reflect.ValueOf(arg))
	var res []any
	switch ref.Kind() {
	case reflect.Slice:
		l := ref.Len()
		v := ref.Slice(0, l)
		for i := 0; i < l; i++ {
			res = append(res, v.Index(i).Interface())
		}
	default:
		res = append(res, ref.Interface())
	}
	return res
}
