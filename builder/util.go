package builder

import "reflect"

func toSlice(arg any) []any {
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
