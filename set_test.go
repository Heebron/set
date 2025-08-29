package set

import (
	"sync"
	"testing"
	"time"
)

func TestSet_String(t *testing.T) {

	stringSet := New[string]()

	if !stringSet.Add("hello") {
		t.Fail()
	}

	if stringSet.Contains("goodbye") {
		t.Fail()
	}

	if !stringSet.Add("goodbye") {
		t.Fail()
	}

	if !stringSet.Contains("goodbye") {
		t.Fail()
	}
}

func TestSet_Int(t *testing.T) {

	stringSet := New[int]()

	if stringSet.Size() != 0 {
		t.Fail()
	}

	if !stringSet.Add(1) {
		t.Fail()
	}

	if stringSet.Contains(2) {
		t.Fail()
	}

	if !stringSet.Add(2) {
		t.Fail()
	}

	if !stringSet.Contains(2) {
		t.Fail()
	}
}

func TestSet_Intersect(t *testing.T) {
	setA := New[int]()
	setB := New[int]()

	setA.Add(1)
	setA.Add(2)
	setB.Add(2)
	setB.Add(3)

	setC := setA.Intersect(setB)

	if setC.Size() != 1 {
		t.Fail()
	}

	if setC.Contains(1) {
		t.Fail()
	}

	if setC.Contains(3) {
		t.Fail()
	}

	if !setC.Contains(2) {
		t.Fail()
	}
}

func TestSet_Union(t *testing.T) {
	setA := New[int]()
	setB := New[int]()

	setA.Add(1)
	setA.Add(2)
	setB.Add(2)
	setB.Add(3)

	setC := setA.Union(setB)

	if setC.Size() != 3 {
		t.Fail()
	}

	if !setC.Contains(1) {
		t.Fail()
	}

	if !setC.Contains(3) {
		t.Fail()
	}

	if !setC.Contains(2) {
		t.Fail()
	}
}

func TestMembers(t *testing.T) {
	set := NewWithInitializer("Fred", "Wilma", "Barney")
	m := set.Members()
	if len(m) != set.Size() {
		t.Fail()
	}
	set.Remove("Barney")
	set.Remove("Fred")

	m = set.Members()

	if len(m) != 1 || m[0] != "Wilma" {
		t.Fail()
	}
}

func TestClearingEmptySet(t *testing.T) {
	set := NewConcurrentWithInitializer("Fred", "Wilma", "Barney")
	set.Clear() // should not panic
	set.Clear() // should not panic

	if set.Size() != 0 {
		t.Fail()
	}
}

func TestWaitForEntryTimeout(t *testing.T) {
	set := NewConcurrentWithInitializer("Fred", "Wilma", "Barney")

	if set.WaitForEmptyWithTimeout(time.Millisecond) {
		t.Fail()
	}

	set = NewConcurrentWithInitializer[string]()

	set.Remove("Fred")

	if !set.WaitForEmptyWithTimeout(time.Millisecond) {
		t.Fail()
	}

	set = NewConcurrent[string]()

	if !set.WaitForEmptyWithTimeout(time.Millisecond) {
		t.Fail()
	}

	set.Add("Fred")

	if set.WaitForEmptyWithTimeout(time.Millisecond) {
		t.Fail()
	}

	var result bool
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		result = set.WaitForEmptyWithTimeout(time.Second)
		wg.Done()
	}()

	set.Remove("Fred")
	wg.Wait()
	if !result {
		t.Fail()
	}

	var result2 bool
	wg.Add(2)
	set.Add("Fred")

	go func() {
		result2 = set.WaitForEmptyWithTimeout(time.Second)
		wg.Done()
	}()

	go func() {
		result = set.WaitForEmptyWithTimeout(time.Second)
		wg.Done()
	}()

	set.Remove("Fred")
	wg.Wait()
	if !(result && result2) {
		t.Fail()
	}
}

func TestSet_Difference(t *testing.T) {
	setA := New[int]()
	setB := New[int]()

	setA.Add(1)
	setA.Add(2)
	setB.Add(2)
	setB.Add(3)

	setC := setA.Difference(setB)
	if setC.Size() != 1 {
		t.Fail()
	}
	if !setC.Contains(1) {
		t.Fail()
	}
	if setC.Contains(2) || setC.Contains(3) {
		t.Fail()
	}

	// Difference with empty rhs returns a copy of lhs
	setD := setA.Difference(New[int]())
	if setD.Size() != 2 || !setD.Contains(1) || !setD.Contains(2) {
		t.Fail()
	}

	// Difference with empty lhs is empty
	setE := New[int]()
	setE = setE.Difference(setB)
	if !setE.IsEmpty() {
		t.Fail()
	}
}

func TestSet_IsSubset(t *testing.T) {
	setA := NewWithInitializer(1, 2)
	setB := NewWithInitializer(1, 3, 2)

	if !setA.IsSubset(setB) {
		t.Fail()
	}
	if setB.IsSubset(setA) {
		t.Fail()
	}
	// empty set is subset of any set
	empty := New[int]()
	if !empty.IsSubset(setB) {
		t.Fail()
	}
	// set is subset of itself
	if !setA.IsSubset(setA) {
		t.Fail()
	}
}

func TestSet_Equal(t *testing.T) {
	setA := NewWithInitializer(1, 2, 3)
	setB := NewWithInitializer(3, 2, 1)
	setC := NewWithInitializer(1, 2)
	setD := New[int]()
	setE := New[int]()

	if !setA.Equal(setB) {
		t.Fail()
	}
	if setA.Equal(setC) {
		t.Fail()
	}
	if !setD.Equal(setE) {
		t.Fail()
	}
}

func TestSet_Clone(t *testing.T) {
	orig := NewWithInitializer("a", "b")
	clone := orig.Clone()

	// equal initially
	if !orig.Equal(clone) {
		t.Fail()
	}

	// mutate original; clone should remain unchanged
	orig.Add("c")
	if clone.Contains("c") {
		t.Fail()
	}
	if orig.Equal(clone) {
		t.Fail()
	}

	// mutate clone; original should remain unchanged
	clone.Remove("a")
	if !orig.Contains("a") {
		t.Fail()
	}
}

func TestSet_IsEmpty(t *testing.T) {
	set := New[string]()
	if !set.IsEmpty() {
		t.Fail()
	}
	set.Add("x")
	if set.IsEmpty() {
		t.Fail()
	}
	set.Remove("x")
	if !set.IsEmpty() {
		t.Fail()
	}
	set.Add("y")
	set.Clear()
	if !set.IsEmpty() {
		t.Fail()
	}
}
