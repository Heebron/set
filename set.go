// set provides set operations over comparable types. These operations are not thread safe.

package set

type void struct{}

var member void

type Set[T comparable] map[T]void

func New[T comparable]() Set[T] {
	return Set[T]{}
}

func (s Set[T]) Add(e T) bool {
	_, exists := s[e]
	if exists {
		return false
	}
	s[e] = member
	return true
}

func (s Set[T]) Remove(e T) bool {
	_, exists := s[e]
	if exists {
		delete(s, e)
		return true
	}
	return false
}

func (s Set[T]) Contains(e T) bool {
	_, exists := s[e]
	if exists {
		return true
	}
	return false
}

func (s Set[T]) Intersect(rhs Set[T]) Set[T] {
	newSet := New[T]()

	for k, _ := range rhs {
		if s.Contains(k) {
			newSet.Add(k)
		}
	}
	return newSet
}

func (s Set[T]) Union(rhs Set[T]) Set[T] {
	newSet := New[T]()

	for k, _ := range s {
		newSet.Add(k)
	}

	for k, _ := range rhs {
		if !newSet.Contains(k) {
			newSet.Add(k)
		}
	}

	return newSet
}

func (s Set[T]) Size() int {
	return len(s)
}
