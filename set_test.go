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
