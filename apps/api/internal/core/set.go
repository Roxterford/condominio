package core

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values ...T) Set[T] {
	set := make(Set[T])
	for _, v := range values {
		set.Add(v)
	}
	return set
}

// Has checks if an element is in the set
func (s Set[T]) Has(v T) bool {
	_, ok := s[v]
	return ok
}

// Add adds an element to the set
func (s Set[T]) Add(v T) Set[T] {
	s[v] = struct{}{}
	return s
}

// Remove removes an element from the set
func (s Set[T]) Remove(v T) Set[T] {
	delete(s, v)
	return s
}

// Clear removes all elements from the set
func (s Set[T]) Clear() Set[T] {
	for k := range s {
		delete(s, k)
	}
	return s
}

// IsEmpty checks if the set is empty
func (s Set[T]) IsEmpty() bool {
	return len(s) == 0
}
