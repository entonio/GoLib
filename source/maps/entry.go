package maps

type Entry[K comparable, V any] struct {
	key   K
	value V
}

func (self Entry[K, V]) Key() K {
	return self.key
}

func (self Entry[K, V]) Value() V {
	return self.value
}
