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
	s := New[T]()
	s.Add(members...)
	return s
}

func (s *Set[T]) lock() {
	if s.mutex != nil {
		s.mutex.Lock()
	}
}

func (s *Set[T]) unlock() {
	if s.mutex != nil {
		s.mutex.Unlock()
	}
}

func (s *Set[T]) rLock() {
	if s.mutex != nil {
		s.mutex.RLock()
	}
}

func (s *Set[T]) rUnlock() {
	if s.mutex != nil {
		s.mutex.RUnlock()
	}
}

func (s *Set[T]) newReceiver() (newSet Set[T]) {
	if s.mutex != nil {
		newSet = NewConcurrent[T]()
	} else {
		newSet = New[T]()
	}
	return
}

// Add members.
func (s *Set[T]) Add(e ...T) {
	s.rLock()
	for _, element := range e {
		s.members[element] = member
	}
	s.rUnlock()
}

// Remove member.
func (s *Set[T]) Remove(e ...T) {
	s.lock()
	for _, element := range e {
		delete(s.members, element)
	}
	s.unlock()
}

// Contains returns true if this contains e, else false.
func (s *Set[T]) Contains(e T) bool {
	s.rLock()
	_, exists := s.members[e]
	s.rUnlock()
	return exists
}

// Intersect returns the set intersection of this and rhs.
func (s *Set[T]) Intersect(rhs Set[T]) Set[T] {
	s.rLock()
	rhs.rLock()
	newSet := s.newReceiver()
	for k := range rhs.members {
		if _, exists := s.members[k]; exists {
			newSet.members[k] = member
		}
	}
	s.rUnlock()
	rhs.rUnlock()
	return newSet
}

// Union returns the set union of this and rhs.
func (s *Set[T]) Union(rhs Set[T]) Set[T] {
	s.rLock()
	rhs.rLock()
	newSet := s.newReceiver()
	for k := range s.members {
		newSet.members[k] = member
	}
	for k := range rhs.members {
		newSet.members[k] = member
	}
	s.rUnlock()
	rhs.rUnlock()
	return newSet
}

// Size returns the number of members in this.
func (s *Set[T]) Size() int {
	return len(s.members)
}

// Members returns a slice containing all the members of the set. This is a shallow copy.
func (s *Set[T]) Members() []T {
	s.rLock()
	result := make([]T, 0, len(s.members))
	for k := range s.members {
		result = append(result, k)
	}
	s.rUnlock()
	return result
}

//---- The below came from Venice.ai after sharing the above

// Difference returns the set difference of this and rhs.
func (s *Set[T]) Difference(rhs Set[T]) Set[T] {
	s.rLock()
	rhs.rLock()
	newSet := s.newReceiver()
	for k := range s.members {
		if _, exists := rhs.members[k]; !exists {
			newSet.members[k] = member
		}
	}
	s.rUnlock()
	rhs.rUnlock()
	return newSet
}

// SymmetricDifference returns the symmetric set difference of this and rhs.
func (s *Set[T]) SymmetricDifference(rhs Set[T]) Set[T] {
	s.rLock()
	rhs.rLock()
	newSet := s.newReceiver()
	for k := range s.members {
		if _, exists := rhs.members[k]; !exists {
			newSet.members[k] = member
		}
	}
	for k := range rhs.members {
		if _, exists := s.members[k]; !exists {
			newSet.members[k] = member
		}
	}
	s.rUnlock()
	rhs.rUnlock()
	return newSet
}

// IsSubset checks if this is a subset of rhs.
func (s *Set[T]) IsSubset(rhs Set[T]) bool {
	s.rLock()
	rhs.rLock()
	for k := range s.members {
		if _, exists := rhs.members[k]; !exists {
			s.rUnlock()
			rhs.rUnlock()
			return false
		}
	}
	s.rUnlock()
	rhs.rUnlock()
	return true
}

// IsSuperset checks if this is a superset of rhs.
func (s *Set[T]) IsSuperset(rhs Set[T]) bool {
	s.rLock()
	rhs.rLock()
	for k := range rhs.members {
		if _, exists := s.members[k]; !exists {
			s.rUnlock()
			rhs.rUnlock()
			return false
		}
	}
	s.rUnlock()
	rhs.rUnlock()
	return true
}

// IsDisjoint checks if this and rhs are disjoint.
func (s *Set[T]) IsDisjoint(rhs Set[T]) bool {
	s.rLock()
	rhs.rLock()

	for k := range s.members {
		if _, exists := rhs.members[k]; exists {
			s.rUnlock()
			rhs.rUnlock()
			return false
		}
	}
	s.rUnlock()
	rhs.rUnlock()
	return true
}
