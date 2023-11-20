package lang

func init() {
	SetStackStyle(1)
}

type SV = map[string]any

/*
var mainPackagePrefix *regexp.Regexp

func SetMainPackageName(name string) {
	if len(name) > 0 {
		mainPackagePrefix = regexp.MustCompile(`/[a-zA-Z0-9/_-]+/` + name + `/`)
	} else {
		mainPackagePrefix = nil
	}
}

func AsString(array []any, each func(any) string, separator string, conjunction string) string {
}

func Interface(array []*Entidade) []any {
	copy := make([]any, len(array))
	for i, v := range array {
		copy[i] = v
	}
	return copy
}

func InterfaceFn(fn func(Entidade) string) func(any) string {
	return func(o any) string {
		x, _ := o.(Entidade)
		return fn(x)
	}
}
*/
