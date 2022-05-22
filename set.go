// set provides set operations over comparable types. These operations are not thread safe. The implementation is
// a map.

package set

type void struct{}

var member void

type Set[T comparable] map[T]void

// New creates a new set of type T. Type T must implement comparable.
func New[T comparable]() Set[T] {
	return Set[T]{}
}

// Add element e to this.
func (s Set[T]) Add(e T) bool {
	_, exists := s[e]
	if exists {
		return false
	}
	s[e] = member
	return true
}

// Remove element e from this.
func (s Set[T]) Remove(e T) bool {
	_, exists := s[e]
	if exists {
		delete(s, e)
		return true
	}
	return false
}

// Contains returns true if this contains element e.
func (s Set[T]) Contains(e T) bool {
	_, exists := s[e]
	if exists {
		return true
	}
	return false
}

// Intersect returns the intersection of this and rhs.
func (s Set[T]) Intersect(rhs Set[T]) Set[T] {
	newSet := New[T]()

	for k := range rhs {
		if s.Contains(k) {
			newSet.Add(k)
		}
	}
	return newSet
}

// Union returns the union of this and rhs.
func (s Set[T]) Union(rhs Set[T]) Set[T] {
	newSet := New[T]()

	for k := range s {
		newSet.Add(k)
	}

	for k := range rhs {
		if !newSet.Contains(k) {
			newSet.Add(k)
		}
	}

	return newSet
}

// Size returns the number of elements in this.
func (s Set[T]) Size() int {
	return len(s)
}
