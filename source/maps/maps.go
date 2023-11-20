package maps

/*
func CreateSM(values ...struct{}) map[struct{}]struct{} {
	m := make(map[struct{}]struct{})
	for i := 0; i < len(values); i += 2 {
		m[values[i]] = values[i+1]
	}
	return m
}

func CreateIM(values ...any) map[any]any {
	m := make(map[any]any)
	for i := 0; i < len(values); i += 2 {
		m[values[i]] = values[i+1]
	}
	return m
}

func Create(keyType reflect.Type, valueType reflect.Type, values ...any) (mapKV any) {
	m := reflect.MakeMap(reflect.MapOf(
		keyType,
		valueType,
	))
	for i := 0; i < len(values); i += 2 {
		m.SetMapIndex(reflect.ValueOf(values[i]), reflect.ValueOf(values[i+1]))
	}
	return m.Interface()
}
*/

func Keys[K comparable, V any](dictionary map[K]V) []K {
	results := make([]K, len(dictionary))
	i := -1
	for k, _ := range dictionary {
		i += 1
		results[i] = k
	}
	return results
	/*
		dv := reflect.ValueOf(dictionary)
		kt := dv.Type().Key()
		mk := dv.MapKeys()
		slice := reflect.MakeSlice(reflect.SliceOf(kt), 0, len(mk))
		for i := 0; i < len(mk); i += 1 {
			slice = reflect.Append(slice, mk[i])
		}
		return slice.Interface().([]K)
	*/
}

/*
func AnyValue[K comparable, V any](dictionary map[K]V) V {
	for _, v := range dictionary {
		return v
	}
	var none V
	return none
	/*
		dv := reflect.ValueOf(dictionary)
		mk := dv.MapKeys()
		v := dv.MapIndex(mk[0])
		return v.Interface()
	* /
}
*/

func Values[K comparable, V any](dictionary map[K]V) []V {
	results := make([]V, len(dictionary))
	i := -1
	for _, v := range dictionary {
		i += 1
		results[i] = v
	}
	return results
	/*
		dv := reflect.ValueOf(dictionary)
		vt := dv.Type().Elem()
		mk := dv.MapKeys()
		slice := reflect.MakeSlice(reflect.SliceOf(vt), 0, len(mk))
		for i := 0; i < len(mk); i += 1 {
			v := dv.MapIndex(mk[i])
			slice = reflect.Append(slice, v)
		}
		return slice.Interface()
	*/
}

func ContainsKey[K comparable, V comparable](dictionary map[K]V, key K) bool {
	for k, _ := range dictionary {
		if k == key {
			return true
		}
	}
	return false
}

func ContainsValue[K comparable, V comparable](dictionary map[K]V, value V) bool {
	for _, v := range dictionary {
		if v == value {
			return true
		}
	}
	return false
	/*
		m := reflect.ValueOf(dictionary)
		mk := m.MapKeys()
		for i := 0; i < len(mk); i += 1 {
			if value == m.MapIndex(mk[i]).Interface() {
				return true
			}
		}
		return false
	*/
}

func Merge[K comparable, V any](parts ...map[K]V) map[K]V {
	return PutInto(nil, parts...)
}

func PutInto[K comparable, V any](into map[K]V, others ...map[K]V) map[K]V {
	if into == nil {
		into = make(map[K]V)
	}
	for _, m := range others {
		for k, v := range m {
			into[k] = v
		}
	}
	return into
}

func ToArray[K comparable, V any](m map[K]V) []Entry[K, V] {
	a := make([]Entry[K, V], len(m))
	i := -1
	for k, v := range m {
		i += 1
		a[i] = Entry[K, V]{key: k, value: v}
	}
	return a
}
