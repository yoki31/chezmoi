package chezmoi

type set[T comparable] map[T]struct{}

func newSet[T comparable](elements ...T) set[T] {
	s := set[T](make(map[T]struct{}))
	s.add(elements...)
	return s
}

func (s set[T]) add(elements ...T) {
	for _, element := range elements {
		s[element] = struct{}{}
	}
}

func (s set[T]) contains(element T) bool {
	_, ok := s[element]
	return ok
}

func (s set[T]) element() T {
	for element := range s {
		return element
	}
	var zero T
	return zero
}

func (s set[T]) elements() []T {
	elements := make([]T, 0, len(s))
	for element := range s {
		elements = append(elements, element)
	}
	return elements
}

func (s set[T]) remove(element T) {
	delete(s, element)
}
