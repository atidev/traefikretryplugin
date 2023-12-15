package set

var e = struct{}{}

type set[T comparable] struct {
	ul map[T]struct{}
}

type Set[T comparable] interface {
	Add(value T) bool
	Delete(value T) bool
	Includes(value T) bool
	Len() int
	Values() []T
}

func (s *set[T]) Add(value T) bool {
	if _, ok := s.ul[value]; ok {
		return false
	}

	s.ul[value] = e

	return true
}

func (s *set[T]) Delete(value T) bool {
	if _, ok := s.ul[value]; ok {
		return false
	}

	delete(s.ul, value)

	return true
}

func (s *set[T]) Includes(value T) bool {
	_, ok := s.ul[value]

	return ok
}

func (s *set[T]) Len() int {
	return len(s.ul)
}

func (s *set[T]) Values() []T {
	values := make([]T, 0, len(s.ul))

	for v := range s.ul {
		values = append(values, v)
	}

	return values
}

func NewSet[T comparable]() Set[T] {
	return &set[T]{
		ul: make(map[T]struct{}),
	}
}
