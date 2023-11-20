package lang

func Call[A any](f func() A) A {
	if f != nil {
		return f()
	}
	var zero A
	return zero
}
