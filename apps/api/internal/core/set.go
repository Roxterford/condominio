package core

import (
	"iter"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](values ...T) Set[T] {
	set := make(Set[T])
	for _, v := range values {
		set.Add(v)
	}
	return set
}

// NewSetFrom creates a new set from an iterable and a mapper function
func NewSetFrom[T any, R comparable](iterable iter.Seq[T], mapper func(it T) R) Set[R] {
	set := make(Set[R])
	for v := range iterable {
		set.Add(mapper(v))
	}
	return set
}

func NewSetFromSlice[T any, R comparable](iterable []T, mapper func(it T) R) Set[R] {
	// Al pasar len(iterable), el mapa reserva espacio y evitas
	// que tenga que crecer (re-hash) mientras lo llenas.
	set := make(Set[R], len(iterable))
	for _, v := range iterable {
		set.Add(mapper(v))
	}
	return set
}

func (s Set[T]) ForEach(iterator func(it T, i int)) {
	i := 0
	for k := range s {
		iterator(k, i)
		i++
	}
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
