package lang

func Zero[T any]() (zero T) {
	return
}

func IsZero[A comparable](a A) bool {
	return a == Zero[A]()
}
