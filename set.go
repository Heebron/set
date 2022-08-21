// set provides set operations over comparable types.

package set

import (
	"sync"
)

type void struct{}

var member void

type Set[T comparable] struct {
	elements map[T]void
	mutex    *sync.RWMutex // if nil, don't use mutex
}

// NewConcurrent returns a set that is concurrent safe.
func NewConcurrent[T comparable]() Set[T] {
	return Set[T]{elements: make(map[T]void), mutex: new(sync.RWMutex)}
}

// New returns a set that is not concurrent safe.
func New[T comparable]() Set[T] {
	return Set[T]{elements: make(map[T]void)}
}

// Add element.
func (s Set[T]) Add(e T) bool {
	if s.mutex != nil {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}
	_, exists := s.elements[e]
	if exists {
		return false
	}
	s.elements[e] = member
	return true
}

// Remove element.
func (s Set[T]) Remove(e T) bool {
	if s.mutex != nil {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
	}
	_, exists := s.elements[e]
	if exists {
		delete(s.elements, e)
		return true
	}
	return false
}

// Contains returns true if this contains e, else false.
func (s Set[T]) Contains(e T) bool {
	if s.mutex != nil {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}
	_, exists := s.elements[e]
	if exists {
		return true
	}
	return false
}

// Intersect returns the set intersection of this and rhs.
func (s Set[T]) Intersect(rhs Set[T]) Set[T] {
	if s.mutex != nil {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
	}

	var newSet Set[T]
	if s.mutex != nil {
		newSet = NewConcurrent[T]()
	} else {
		newSet = New[T]()
	}

	for k := range rhs.elements {
		_, exists := s.elements[k]
		if exists {
			newSet.elements[k] = member
		}
	}
	return newSet
}

// Union returns the set union of this and rhs.
func (s Set[T]) Union(rhs Set[T]) Set[T] {
	if s.mutex != nil {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
	}

	var newSet Set[T]
	if s.mutex != nil {
		newSet = NewConcurrent[T]()
	} else {
		newSet = New[T]()
	}

	for k := range s.elements {
		newSet.elements[k] = member
	}

	for k := range rhs.elements {
		newSet.elements[k] = member
	}

	return newSet
}

// Size returns the number of elements in this.
func (s Set[T]) Size() int {
	return len(s.elements)
}
