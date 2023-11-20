package lang

func T2[A any, B any](v1 A, v2 B) (A, B) {
	return v1, v2
}

func V1of2[A any](v1 A, v2 any) A {
	return v1
}
