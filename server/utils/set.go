package utils

var _exist = struct{}{}

type Set[T comparable] struct {
	m map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		m: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(v T) {
	s.m[v] = _exist
}

func (s *Set[T]) Remove(v T) {
	delete(s.m, v)
}

func (s *Set[T]) Contain(v T) (exist bool) {
	_, exist = s.m[v]
	return
}

func (s *Set[T]) GetAndRemove() (v T) {
	for _, v = range s.m {
		break
	}
	delete(s.m, v)
	return
}

func (s *Set[T]) Size() int {
	return len(s.m)
}
