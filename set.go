// Package set provides a generic hash set for comparable types.
// A Set can be created in two modes:
//   - concurrent: safe for use by multiple goroutines
//   - non-concurrent: faster, but not safe for concurrent access
//
// The zero value of Set is not ready for use; construct sets with New,
// NewWithInitializer, NewConcurrent, or NewConcurrentWithInitializer.
//
// Unless otherwise stated, methods do not modify their arguments and
// return new sets when producing results (e.g., Union, Intersect).
package set

import (
	"context"
	"sync"
	"time"
)

type void struct{} // empty element value for map

var voidValue void

// Set is a generic collection of unique elements of a comparable type.
// Internally, it uses a map[T]void to represent membership.
//
// Concurrency:
//   - Sets constructed via NewConcurrent or NewConcurrentWithInitializer
//     synchronize method calls using an RWMutex.
//   - Sets constructed via New or NewWithInitializer are not synchronized and
//     must not be accessed from multiple goroutines without external
//     synchronization.
type Set[T comparable] struct {
	members map[T]void
	mutex   *sync.RWMutex // if nil, the set is non-concurrent and performs no locking
	trigger chan struct{} // used to indicate transition to empty set
}

// NewConcurrent returns a set that is concurrent safe.
func NewConcurrent[T comparable]() Set[T] {
	return Set[T]{members: make(map[T]void), mutex: new(sync.RWMutex)}
}

// NewConcurrentWithInitializer returns a set that is concurrent safe and contains the provided initial set of members.
func NewConcurrentWithInitializer[T comparable](members ...T) Set[T] {
	s := Set[T]{members: make(map[T]void), mutex: new(sync.RWMutex)}
	for _, v := range members {
		s.members[v] = voidValue
	}
	if len(s.members) > 0 {
		s.trigger = make(chan struct{})
	}
	return s
}

// New returns a set that is not concurrent safe.
func New[T comparable]() Set[T] {
	return Set[T]{members: make(map[T]void)}
}

// NewWithInitializer returns a set that is not concurrent safe and contains the provided initial set of members.
func NewWithInitializer[T comparable](members ...T) Set[T] {
	s := Set[T]{members: make(map[T]void)}
	for _, v := range members {
		s.members[v] = voidValue
	}
	return s
}

// Add inserts e into the set.
// It returns true if the set was modified (e was not already present), or false otherwise.
func (s *Set[T]) Add(e T) bool {
	if s.mutex != nil {
		s.mutex.Lock()
		defer s.mutex.Unlock()

		if len(s.members) == 0 {
			s.trigger = make(chan struct{})
		}
	}

	_, exists := s.members[e]
	if exists {
		return false
	}
	s.members[e] = voidValue
	return true
}

// Remove deletes e from the set.
// It returns true if the set was modified (e was present), or false otherwise.
func (s *Set[T]) Remove(e T) bool {
	var beforeSize int

	if s.mutex != nil {
		s.mutex.Lock()
		beforeSize = len(s.members)
		defer s.mutex.Unlock()
	}
	_, exists := s.members[e]
	if exists {
		delete(s.members, e)

		if s.mutex != nil && beforeSize == 1 && len(s.members) == 0 {
			close(s.trigger)
		}

		return true
	}
	return false
}

// Contains reports whether e is a voidValue of the set.
func (s *Set[T]) Contains(e T) bool {
	if s.mutex != nil {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
	}
	_, exists := s.members[e]
	return exists
}

// Intersect returns a new set containing the elements common to s and rhs.
// The returned set inherits the concurrency mode of the receiver (s).
func (s *Set[T]) Intersect(rhs Set[T]) Set[T] {
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

	if rhs.mutex != nil {
		rhs.mutex.RLock()
		defer rhs.mutex.RUnlock()
	}
	for k := range rhs.members {
		_, exists := s.members[k]
		if exists {
			newSet.members[k] = voidValue
		}
	}
	return newSet
}

// Union returns a new set containing all elements present in either s or rhs.
// The returned set inherits the concurrency mode of the receiver (s).
func (s *Set[T]) Union(rhs Set[T]) Set[T] {
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
		newSet.members[k] = voidValue
	}

	if rhs.mutex != nil {
		rhs.mutex.RLock()
		defer rhs.mutex.RUnlock()
	}
	for k := range rhs.members {
		newSet.members[k] = voidValue
	}

	return newSet
}

// Size returns the number of elements currently in the set.
func (s *Set[T]) Size() int {
	if s.mutex != nil {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
	}
	return len(s.members)
}

// IsEmpty returns true if the set contains no elements, otherwise false.
func (s *Set[T]) IsEmpty() bool {
	return s.Size() == 0
}

// Members returns a snapshot slice containing all members of the set.
// The order of elements in the returned slice is unspecified.
func (s *Set[T]) Members() []T {
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

// Clear removes all elements from the set while maintaining its concurrent/non-concurrent state.
func (s *Set[T]) Clear() {
	var beforeSize int

	if s.mutex != nil {
		s.mutex.Lock()
		beforeSize = len(s.members)
		defer s.mutex.Unlock()
	}
	clear(s.members)

	if s.mutex != nil && beforeSize > 0 {
		close(s.trigger)
	}
}

// Clone returns a new set containing all elements from the original set.
// The returned set inherits the concurrency mode of the receiver.
func (s *Set[T]) Clone() Set[T] {
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
		newSet.members[k] = voidValue
	}
	return newSet
}

// Difference returns a new set containing elements present in s but not in rhs.
// The returned set inherits the concurrency mode of the receiver (s).
func (s *Set[T]) Difference(rhs Set[T]) Set[T] {
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

	if rhs.mutex != nil {
		rhs.mutex.RLock()
		defer rhs.mutex.RUnlock()
	}

	for k := range s.members {
		if _, exists := rhs.members[k]; !exists {
			newSet.members[k] = voidValue
		}
	}
	return newSet
}

// IsSubset returns true if all elements in s are present in rhs.
func (s *Set[T]) IsSubset(rhs Set[T]) bool {
	if s.mutex != nil {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
	}

	if rhs.mutex != nil {
		rhs.mutex.RLock()
		defer rhs.mutex.RUnlock()
	}

	for k := range s.members {
		if _, exists := rhs.members[k]; !exists {
			return false
		}
	}
	return true
}

// Equal returns true if s and rhs contain exactly the same elements.
func (s *Set[T]) Equal(rhs Set[T]) bool {
	if s.mutex != nil {
		s.mutex.RLock()
		defer s.mutex.RUnlock()
	}

	if rhs.mutex != nil {
		rhs.mutex.RLock()
		defer rhs.mutex.RUnlock()
	}

	if len(s.members) != len(rhs.members) {
		return false
	}

	for k := range s.members {
		if _, exists := rhs.members[k]; !exists {
			return false
		}
	}
	return true
}

// WaitForEmptyWithTimeout waits until the set becomes empty or the specified timeout elapses.
// Returns true if the set becomes empty, false if the timeout occurs first.
// The set may receive new elements before this method returns. Negative timeout duration is converted to a duration of 0.
func (s *Set[T]) WaitForEmptyWithTimeout(timeout time.Duration) bool {

	if timeout < 0 {
		timeout = 0
	}

	if s.mutex == nil {
		panic("cannot wait for empty on a non-concurrent set")
	}

	s.mutex.RLock()

	if len(s.members) == 0 {
		s.mutex.RUnlock()
		return true
	}

	// Capture state before unlocking.
	capturedChannel := s.trigger

	// Release mutex.
	s.mutex.RUnlock()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	select {
	case <-capturedChannel: // Captured lock.
	case <-ctx.Done(): // timeout
	}

	// Are we empty now?
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.members) == 0
}
