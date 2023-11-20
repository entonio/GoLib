package lang

import "reflect"

func FieldValue[T any, A any](name string, a A) T {
	if len(name) == 0 {
		var i any
		i = a
		return i.(T)
	}
	r := reflect.ValueOf(a)
	for {
		isPointer := r.Kind() == reflect.Pointer
		if !isPointer {
			break
		}
		r = r.Elem()
	}
	if r.IsZero() {
		var none T
		return none
	}
	v := r.FieldByName(name)
	if v == (reflect.Value{}) {
		v = r.MethodByName(name)
		if v == (reflect.Value{}) {
			var none T
			return none
		}
		v = v.Call(nil)[0]
	}
	return v.Interface().(T)
}
