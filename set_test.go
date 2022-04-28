package set

import (
	"testing"
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
