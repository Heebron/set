// set provides set operations over comparable types.

package set

import (
	"sync"
)

type void struct{}

var member void

type Set[T comparable] struct {
	members map[T]void
	mutex   *sync.RWMutex // if nil, don't use mutex
}

// NewConcurrent returns a set that is thread safe.
func NewConcurrent[T comparable]() Set[T] {
	return Set[T]{members: make(map[T]void), mutex: new(sync.RWMutex)}
}

// NewConcurrentWithInitializer returns a set that is thread safe and contains the provided initial set of members.
func NewConcurrentWithInitializer[T comparable](members ...T) Set[T] {
	s := Set[T]{members: make(map[T]void), mutex: new(sync.RWMutex)}
	for _, v := range members {
		s.members[v] = member
	}
	return s
}

// New returns a set that is not thread safe.
func New[T comparable]() Set[T] {
	return Set[T]{members: make(map[T]void)}
}

// NewWithInitializer returns a set that is not thread safe and contains the provided initial set of members.
func NewWithInitializer[T comparable](members ...T) Set[T] {
	s := Set[T]{members: make(map[T]void)}
	for _, v := range members {
		s.members[v] = member
	}
	return s
}

// Add member.
func (s Set[T]) Add(e T) bool {
	if s.mutex != nil {
		s.mutex.Lock()
		defer s.mutex.Unlock()
	}
	_, exists := s.members[e]
	if exists {
		return false
	}
	s.members[e] = member
	return true
}

// Remove member.
func (s Set[T]) Remove(e T) bool {
	if s.mutex != nil {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
	}
	_, exists := s.members[e]
	if exists {
		delete(s.members, e)
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
	_, exists := s.members[e]
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

	for k := range rhs.members {
		_, exists := s.members[k]
		if exists {
			newSet.members[k] = member
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

	for k := range s.members {
		newSet.members[k] = member
	}

	for k := range rhs.members {
		newSet.members[k] = member
	}

	return newSet
}

// Size returns the number of members in this.
func (s Set[T]) Size() int {
	return len(s.members)
}

// Members returns a slice containing all the members of the set. This is a shallow copy.
func (s Set[T]) Members() []T {
	if s.mutex != nil {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
	}

	result := make([]T, 0, len(s.members))

	for k := range s.members {
		result = append(result, k)
	}
	return result
}
