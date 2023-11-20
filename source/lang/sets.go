package lang

type Set[K comparable] map[K]bool

func NewSet[K comparable](keys ...K) Set[K] {
	s := make(map[K]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}

//func (self Set[K]) Delete(keys ...K) {
//	for _, k := range keys {
//		delete(self, k)
//	}
//}
