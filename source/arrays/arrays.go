package arrays

import (
	"math/rand"
	"reflect"
	"strings"

	"golib/lang"
)

func New[A any](elements ...A) (array []A) {
	return append(array, elements...)
}

func NotNil[A any](array ...*A) *A {
	for _, a := range array {
		if a != nil {
			return a
		}
	}
	return nil
}

func NotEmpty(array ...string) string {
	for _, s := range array {
		if len(s) > 0 {
			return s
		}
	}
	return ""
}

func First[A any](array []A) A {
	return array[0]
}

func Last[A any](array []A) A {
	return array[len(array)-1]
}

func At[A any](array []A, index int) A {
	if array == nil {
		var none A
		return none
	}
	if index < 0 {
		index = len(array) + index
	}
	if index < 0 || index >= len(array) {
		var none A
		return none
	}
	return array[index]
	/*
		if array == nil {
			return nil
		}
		a := reflect.ValueOf(array)
		len := a.Len()
		if index < 0 {
			index = len + index
		}
		if index < 0 || index >= len {
			return nil
		}
		return a.Index(index).Interface()
	*/
}

func Contains[A comparable](array []A, value A) bool {
	return IndexOf(array, value) >= 0
}

func IndexOf[A comparable](array []A, value A) int {
	for i, a := range array {
		if a == value {
			return i
		}
	}
	/*
		if array != nil {
			a := reflect.ValueOf(array)
			len := a.Len()
			for i := 0; i < len; i++ {
				if value == a.Index(i).Interface() {
					return i
				}
			}
		}
	*/
	return -1
}

func WithoutEmptyString(array ...string) []string {
	return Without(array, "")
}

func Without[A comparable](array []A, excluded ...A) (arrayWithoutExcluded []A) {
	return Filter(array, func(candidate A) bool {
		for _, ex := range excluded {
			if candidate == ex {
				return false
			}
		}
		return true
	})
}

func Replace[A comparable](array []A, replaced []A, replacement A) []A {
	results := make([]A, len(array))
	for i, a := range array {
		for _, r := range replaced {
			if a == r {
				a = replacement
				break
			}
		}
		results[i] = a
	}
	return results
	/*
		a := reflect.ValueOf(array)
		la := a.Len()
		ri := reflect.ValueOf(replacement)
		r := reflect.ValueOf(replaced)
		lr := r.Len()
		slice := reflect.MakeSlice(a.Type(), 0, la)
		for i := 0; i < la; i += 1 {
			ai := a.Index(i)
			aii := ai.Interface()
			for j := 0; j < lr; j++ {
				if aii == r.Index(j).Interface() {
					ai = ri
					break
				}
			}
			slice = reflect.Append(slice, ai)
		}
		return slice.Interface().([]A)
	*/
}

func WithoutDuplicates[A comparable](array []A) []A {
	var results []A
	seen := make(map[any]bool)
	for _, a := range array {
		if !seen[a] {
			results = append(results, a)
			seen[a] = true
		}
	}
	return results
	/*
		a := reflect.ValueOf(array)
		len := a.Len()
		slice := reflect.MakeSlice(a.Type(), 0, len)
		present := make(map[any]bool)
		for i := 0; i < len; i += 1 {
			ai := a.Index(i)
			vi := ai.Interface()
			if !present[vi] {
				slice = reflect.Append(slice, ai)
				present[vi] = true
			}
		}
		return slice.Interface().([]A)
	*/
}

func Sum[A any](array []A, valueFunction func(A) float64) float64 {
	return Reduce(array, 0.0, func(current float64, a A) float64 {
		return current + valueFunction(a)
	})
	/*
		av := reflect.ValueOf(valueFunction)
		return Reduce(array, 0.0, func(current float64, element any) float64 {
			ai := reflect.ValueOf(element)
			return current + av.Call([]reflect.Value{ai})[0].Float()
		}).(float64)
	*/
}

func Reduce[A any, V any](array []A, startValue V, combine func(V, A) V) V {
	current := startValue
	for _, a := range array {
		current = combine(current, a)
	}
	return current
	/*
		a := reflect.ValueOf(array)
		len := a.Len()
		ac := reflect.ValueOf(combineFunction)
		current := reflect.ValueOf(startValue)
		for i := 0; i < len; i += 1 {
			ai := a.Index(i)
			current = ac.Call([]reflect.Value{current, ai})[0]
			/*
				if len > 1 {
					fmt.Printf("CR: %d %#v\n", i+1, current)
				}
			* /
		}
		return current.Interface()
	*/
}

func Detect[A comparable](array []A, accept func(A) bool) (firstMatch A, index int) {
	for i, a := range array {
		if accept(a) {
			return a, i
		}
	}
	var none A
	return none, -1
	/*
		a := reflect.ValueOf(array)
		len := a.Len()
		af := reflect.ValueOf(acceptFunction)
		for i := 0; i < len; i += 1 {
			ai := a.Index(i)
			accepted := af.Call([]reflect.Value{ai})[0]
			if accepted.Bool() {
				return ai.Interface(), i
			}
		}
		return nil, -1
	*/
}

/*
func AllMatch(array any, acceptFunction any) bool {
	return _match(array, acceptFunction, true)
}
*/

func AnyMatch[A comparable](array []A, accept func(A) bool) bool {
	for _, a := range array {
		if accept(a) {
			return true
		}
	}
	return false
}

/*
func _match(array any, acceptFunction any, matchAll bool) bool {
	a := reflect.ValueOf(array)
	len := a.Len()
	af := reflect.ValueOf(acceptFunction)
	in := af.Type().NumIn()
	for i := 0; i < len; i += 1 {
		ai := a.Index(i)
		parameters := []reflect.Value{ai}
		if in == 2 {
			parameters = append(parameters, reflect.ValueOf(i))
		}
		accepted := af.Call(parameters)[0]
		if accepted.Bool() != matchAll {
			return !matchAll
		}
	}
	return matchAll
}
*/

func Shuffle(array any) {
	a := reflect.ValueOf(array)
	len := a.Len()
	for i := len - 1; i > 0; i -= 1 {
		j := rand.Intn(i + 1) // 0 <= j <= i
		vi := a.Index(i)
		vj := a.Index(j)
		vvi := vi.Interface()
		vvj := vj.Interface()
		vi.Set(reflect.ValueOf(vvj))
		vj.Set(reflect.ValueOf(vvi))
	}
}

func Shuffled(array any) (shuffledArray any) {
	a := reflect.ValueOf(array)
	len := a.Len()
	slice := reflect.MakeSlice(a.Type(), len, len)
	for i := len - 1; i > 0; i -= 1 {
		j := rand.Intn(i + 1) // 0 <= j <= i
		slice.Index(i).Set(a.Index(j))
		slice.Index(j).Set(a.Index(i))
	}
	return slice.Interface()
}

func Partition[A any](array []A, accept func(A) bool) (accepted []A, notAccepted []A) {
	for _, a := range array {
		if accept(a) {
			accepted = append(accepted, a)
		} else {
			notAccepted = append(notAccepted, a)
		}
	}
	return accepted, notAccepted
}

func Filter[A any](array []A, accept func(A) bool) []A {
	var results []A
	for _, a := range array {
		if accept(a) {
			results = append(results, a)
		}
	}
	return results
	/*
		a := reflect.ValueOf(array)
		len := a.Len()
		af := reflect.ValueOf(acceptFunction)
		slice := reflect.MakeSlice(a.Type(), 0, len)
		in := af.Type().NumIn()
		for i := 0; i < len; i += 1 {
			ai := a.Index(i)
			parameters := []reflect.Value{ai}
			if in == 2 {
				parameters = append(parameters, reflect.ValueOf(i))
			}
			accepted := af.Call(parameters)[0]
			if accepted.Bool() {
				slice = reflect.Append(slice, ai)
			}
		}
		return slice.Interface().([]A)
	*/
}

func FilterMap[A any, T any](array []A, acceptTransform func(A) (bool, T)) (filteredAndTransfomedArray []T) {
	var results []T
	for _, a := range array {
		accepted, t := acceptTransform(a)
		if accepted {
			results = append(results, t)
		}
	}
	return results
}

func FilterMapIndex[A any, T any](array []A, acceptTransform func(A, int) (bool, T)) []T {
	var results []T
	for i, a := range array {
		accepted, t := acceptTransform(a, i)
		if accepted {
			results = append(results, t)
		}
	}
	return results
	/*
		a := reflect.ValueOf(array)
		len := a.Len()
		atf := reflect.ValueOf(acceptTransformFunction)
		tt := reflect.TypeOf(acceptTransformFunction).Out(1)
		slice := reflect.MakeSlice(reflect.SliceOf(tt), 0, len)
		for i := 0; i < len; i += 1 {
			ai := a.Index(i)
			parameters := []reflect.Value{ai, reflect.ValueOf(i)}
			result := atf.Call(parameters)
			accepted := result[0]
			if accepted.Bool() {
				ti := result[1]
				slice = reflect.Append(slice, ti)
			}
		}
		return slice.Interface().([]T)
	*/
}

/*
func FilterThenMap(array any, acceptFunction any, transformFunction any) (filteredAndTransfomedArray any) {
	a := reflect.ValueOf(array)
	len := a.Len()
	af := reflect.ValueOf(acceptFunction)
	tf := reflect.ValueOf(transformFunction)
	tt := reflect.TypeOf(transformFunction).Out(0)
	slice := reflect.MakeSlice(reflect.SliceOf(tt), 0, len)
	for i := 0; i < len; i += 1 {
		ai := a.Index(i)
		accepted := af.Call([]reflect.Value{ai})[0]
		if accepted.Bool() {
			ti := tf.Call([]reflect.Value{ai})[0]
			slice = reflect.Append(slice, ti)
		}
	}
	return slice.Interface()
}
*/

func Map[A any, T any](array []A, transform func(A) T) (transfomed []T) {
	results := make([]T, len(array))
	for i, a := range array {
		results[i] = transform(a)
	}
	return results
}

func MapUnique[A comparable, T comparable](array []A, transform func(A) T) []T {
	var results []T
	seen := make(map[T]bool)
	for _, a := range array {
		t := transform(a)
		if !seen[t] {
			results = append(results, t)
			seen[t] = true
		}
	}
	return results
}

func CollectF[A any, K comparable, V any](array []A, acceptKeyValue func(A) (bool, K, V)) (filteredAndCollectedMap map[K]V) {
	m := make(map[K]V)
	for _, a := range array {
		accepted, k, v := acceptKeyValue(a)
		if accepted {
			m[k] = v
		}
	}
	return m
	/*
		a := reflect.ValueOf(array)
		len := a.Len()
		akvf := reflect.ValueOf(acceptKeyValue)
		kt := reflect.TypeOf(acceptKeyValue).Out(1)
		vt := reflect.TypeOf(acceptKeyValue).Out(2)
		m := reflect.MakeMap(reflect.MapOf(kt, vt))
		for i := 0; i < len; i += 1 {
			ai := a.Index(i)
			result := akvf.Call([]reflect.Value{ai})
			accepted := result[0]
			if accepted.Bool() {
				ki := result[1]
				vi := result[2]
				m.SetMapIndex(ki, vi)
			}
		}
		return m.Interface().(map[K]V)
	*/
}

func GroupF[A any, K comparable, V any](array []A, acceptKeyValueFunction func(A) (bool, K, V)) (filteredAndCollectedMap map[K][]V) {
	m := make(map[K][]V)
	for _, a := range array {
		accepted, k, v := acceptKeyValueFunction(a)
		if accepted {
			m[k] = append(m[k], v)
		}
	}
	return m
	/*
		a := reflect.ValueOf(array)
		len := a.Len()
		akvf := reflect.ValueOf(acceptKeyValueFunction)
		kt := reflect.TypeOf(acceptKeyValueFunction).Out(1)
		st := reflect.SliceOf(reflect.TypeOf(acceptKeyValueFunction).Out(2))
		m := reflect.MakeMap(reflect.MapOf(kt, st))
		for i := 0; i < len; i += 1 {
			ai := a.Index(i)
			result := akvf.Call([]reflect.Value{ai})
			accepted := result[0]
			if accepted.Bool() {
				ki := result[1]
				vi := result[2]
				si := m.MapIndex(ki)
				if !si.IsValid() {
					si = reflect.MakeSlice(st, 0, 0)
				}
				m.SetMapIndex(ki, reflect.Append(si, vi))
			}
		}
		return m.Interface().(map[K][]V)
	*/
}

/*
func CollectT[A any, K comparable, V any](array []A, m map[K]V) map[K]V {
	mt := reflect.TypeOf(m)
	kt := mt.Key()
	vt := mt.Elem()
	var fields []reflect.StructField
	for i := 0; i < len(array); i += 1 {
		ai := reflect.ValueOf(array[i])
		if fields == nil {
			fields = reflect.VisibleFields(ai.Type())
		}
		var k K
		var v V
		kf := false
		vf := false
		for _, field := range fields {
			if !kf && field.Type == kt {
				k = ai.FieldByName(field.Name).Interface().(K)
				if vf {
					break
				}
				kf = true
			} else if !vf && field.Type == vt {
				v = ai.FieldByName(field.Name).Interface().(V)
				if kf {
					break
				}
				vf = true
			}
		}
		m[k] = v
	}
	return m
}

func GroupT[A any, K comparable, V any](array []A, m map[K][]V) map[K][]V {
	mt := reflect.TypeOf(m)
	kt := mt.Key()
	vt := mt.Elem()
	var fields []reflect.StructField
	for i := 0; i < len(array); i += 1 {
		ai := reflect.ValueOf(array[i])
		if fields == nil {
			fields = reflect.VisibleFields(ai.Type())
		}
		var k K
		var v V
		kf := false
		vf := false
		for _, field := range fields {
			if !kf && field.Type == kt {
				k = ai.FieldByName(field.Name).Interface().(K)
				if vf {
					break
				}
				kf = true
			} else if !vf && field.Type == vt {
				v = ai.FieldByName(field.Name).Interface().(V)
				if kf {
					break
				}
				vf = true
			}
		}
		m[k] = append(m[k], v)
	}
	return m
}

func CollectN[A any, K comparable, V any](array []A, m map[K]V, keyName string, valueName string) map[K]V {
	for i := 0; i < len(array); i += 1 {
		ai := reflect.ValueOf(array[i])
		k := ai.FieldByName(keyName).Interface().(K)
		v := ai.FieldByName(valueName).Interface().(V)
		m[k] = v
	}
	return m
}

func GroupN[A any, K comparable, V any](array []A, m map[K][]V, keyName string, valueName string) map[K][]V {
	for i := 0; i < len(array); i += 1 {
		ai := reflect.ValueOf(array[i])
		k := ai.FieldByName(keyName).Interface().(K)
		v := ai.FieldByName(valueName).Interface().(V)
		m[k] = append(m[k], v)
	}
	return m
}
*/

/*
func ReverseArray[A any](array []A) (reversed []A) {
	a := reflect.ValueOf(array)
	len := a.Len()
	slice := reflect.MakeSlice(a.Type(), 0, len)
	for i := len - 1; i >= 0; i -= 1 {
		ai := a.Index(i)
		slice = reflect.Append(slice, ai)
	}
	return slice.Interface().([]A)
}
*/

func AsMap[A comparable, V any](array []A, value V) map[A]V {
	m := make(map[A]V, len(array))
	for _, a := range array {
		m[a] = value
	}
	return m
	/*
		a := reflect.ValueOf(array)
		v := reflect.ValueOf(value)
		m := reflect.MakeMap(reflect.MapOf(
			a.Type().Elem(),
			v.Type(),
		))
		len := a.Len()
		for i := 0; i < len; i += 1 {
			m.SetMapIndex(a.Index(i), v)
		}
		return m.Interface().(map[A]V)
	*/
}

func InterfaceArray(original []string, from int) (array []any) {
	array = make([]any, lang.MaxInt(0, len(original)-from))
	for i := 0; i < len(array); i += 1 {
		array[i] = original[from+i]
	}
	return
}

func AsInts(s string) []int {
	return Map(strings.Split(s, ","), func(s string) int {
		return lang.AsInt(strings.TrimSpace(s))
	})
}

type UniqueTransform[A comparable, T any] interface {
	Append(A)
	Contents() []T
}

func NewUniqueTransform[A comparable, T any](transform func(A) T) UniqueTransform[A, T] {
	return &_UniqueTransform[A, T]{
		transform: transform,
		appended:  make(map[A]bool),
	}
}

// this has to be hidden because it needs to be created using NewUniqueTransform, because:
// - it has to be a pointer so that mutating ops aren't lost
// - it has to initialise `appended` to make things easier
type _UniqueTransform[A comparable, T any] struct {
	appended  map[A]bool
	transform func(A) T
	contents  []T
}

func (self *_UniqueTransform[A, T]) Append(element A) {
	if !self.appended[element] {
		t := self.transform(element)
		self.contents = append(self.contents, t)
		self.appended[element] = true
		//fmt.Printf("contents: %d\n", len(self.contents))
	}
}

func (self *_UniqueTransform[A, T]) Contents() []T {
	return self.contents
}

func MapField[T any, A any](name string, array []A) []T {
	t := make([]T, len(array))
	for i, a := range array {
		t[i] = lang.FieldValue[T](name, a)
	}
	return t
}

func CollectField[K comparable, V any, A any](keyName string, valueName string, array []A) map[K]V {
	m := make(map[K]V)
	for _, a := range array {
		key := lang.FieldValue[K](keyName, a)
		m[key] = lang.FieldValue[V](valueName, a)
	}
	return m
}

func GroupField[K comparable, V any, A any](keyName string, valueName string, array []A) map[K][]V {
	m := make(map[K][]V)
	for _, a := range array {
		key := lang.FieldValue[K](keyName, a)
		m[key] = append(m[key], lang.FieldValue[V](valueName, a))
	}
	return m
}

func Merge[A any](parts ...[]A) []A {
	size := 0
	for _, part := range parts {
		size += len(part)
	}
	result := make([]A, 0, size)
	for _, part := range parts {
		result = append(result, part...)
	}
	return result
}

func Print1d[A any](array1d []A, separator1d string) string {
	var result string
	for index1d, value := range array1d {
		if index1d > 0 {
			result += separator1d
		}
		result += lang.Print(value)
	}
	return result
}

func Print2d[A any](array2d [][]A, separator1d string, separator2d string, nil1d string) string {
	var result string
	for index2d, array1d := range array2d {
		if index2d > 0 {
			result += separator2d
		}
		if array1d == nil {
			result += nil1d
			continue
		}
		result += Print1d(array1d, separator1d)
	}
	return result
}
