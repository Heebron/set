package set

import (
	"testing"
)

func TestSet_String(t *testing.T) {

	stringSet := New[string]()

	stringSet.Add("hello")

	if stringSet.Contains("goodbye") {
		t.Fail()
	}

	stringSet.Add("goodbye")

	if !stringSet.Contains("goodbye") {
		t.Fail()
	}
}

func TestSet_Int(t *testing.T) {

	stringSet := New[int]()

	if stringSet.Size() != 0 {
		t.Fail()
	}

	stringSet.Add(1)

	if stringSet.Contains(2) {
		t.Fail()
	}

	stringSet.Add(2)

	if !stringSet.Contains(2) {
		t.Fail()
	}
}

func Test_Intersect(t *testing.T) {
	set1 := New[int]()
	set1.Add(1, 2, 3)
	set2 := New[int]()
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)
	intersect := set1.Intersect(set2)
	if intersect.Size() != 2 {
		t.Errorf("Intersect should return a set with 2 elements")
	}
	if !intersect.Contains(2) || !intersect.Contains(3) {
		t.Errorf("Intersect should contain elements 2 and 3")
	}
}

func TestSet_Union(t *testing.T) {
	setA := New[int]()
	setB := New[int]()

	setA.Add(1, 2, 3)

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

// --- all tests below were created by venice.ai

func TestNew(t *testing.T) {
	set := New[int]()
	if set.Size() != 0 {
		t.Errorf("New set should be empty")
	}
}

func TestNewConcurrent(t *testing.T) {
	set := NewConcurrent[int]()
	if set.Size() != 0 {
		t.Errorf("New concurrent set should be empty")
	}
}

func TestNewWithInitializer(t *testing.T) {
	set := NewWithInitializer[int](1, 2, 3)
	if set.Size() != 3 {
		t.Errorf("New set with initializer should have 3 elements")
	}
}

func TestNewConcurrentWithInitializer(t *testing.T) {
	set := NewConcurrentWithInitializer[int](1, 2, 3)
	if set.Size() != 3 {
		t.Errorf("New concurrent set with initializer should have 3 elements")
	}
}

func TestContains(t *testing.T) {
	set := New[int]()
	set.Add(1)
	if !set.Contains(1) {
		t.Errorf("Contains should return true for an existing element")
	}
	if set.Contains(2) {
		t.Errorf("Contains should return false for a non-existing element")
	}
}

func TestIntersect(t *testing.T) {
	set1 := New[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := New[int]()
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)
	intersect := set1.Intersect(set2)
	if intersect.Size() != 2 {
		t.Errorf("Intersect should return a set with 2 elements")
	}
	if !intersect.Contains(2) || !intersect.Contains(3) {
		t.Errorf("Intersect should contain elements 2 and 3")
	}
}

func TestUnion(t *testing.T) {
	set1 := New[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := New[int]()
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)
	union := set1.Union(set2)
	if union.Size() != 4 {
		t.Errorf("Union should return a set with 4 elements")
	}
	if !union.Contains(1) || !union.Contains(2) || !union.Contains(3) || !union.Contains(4) {
		t.Errorf("Union should contain elements 1, 2, 3, and 4")
	}
}

func TestSize(t *testing.T) {
	set := New[int]()
	if set.Size() != 0 {
		t.Errorf("Empty set should have size 0")
	}
	set.Add(1)
	if set.Size() != 1 {
		t.Errorf("Set with 1 element should have size 1")
	}
	set.Add(2)
	if set.Size() != 2 {
		t.Errorf("Set with 2 elements should have size 2")
	}
}

func TestMembersAI(t *testing.T) {
	set := New[int]()
	set.Add(1)
	set.Add(2)
	members := set.Members()
	if len(members) != 2 {
		t.Errorf("Members should return a slice with 2 elements")
	}
	if !contains(members, 1) || !contains(members, 2) {
		t.Errorf("Members should contain elements 1 and 2")
	}
}

func TestDifference(t *testing.T) {
	set1 := New[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := New[int]()
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)
	difference := set1.Difference(set2)
	if difference.Size() != 1 {
		t.Errorf("Difference should return a set with 1 element")
	}
	if !difference.Contains(1) {
		t.Errorf("Difference should contain element 1")
	}
}

func TestSymmetricDifference(t *testing.T) {
	set1 := New[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)
	set2 := New[int]()
	set2.Add(2)
	set2.Add(3)
	set2.Add(4)
	symDiff := set1.SymmetricDifference(set2)
	if symDiff.Size() != 2 {
		t.Errorf("Symmetric difference should return a set with 2 elements")
	}
	if !symDiff.Contains(1) || !symDiff.Contains(4) {
		t.Errorf("Symmetric difference should contain elements 1 and 4")
	}
}

func TestIsSubset(t *testing.T) {
	set1 := New[int]()
	set1.Add(1)
	set1.Add(2)

	set2 := New[int]()
	set2.Add(1)
	set2.Add(2)
	set2.Add(3)

	if !set1.IsSubset(set2) {
		t.Errorf("Set1 should be a subset of set2")
	}

	if set2.IsSubset(set1) {
		t.Errorf("Set2 should not be a subset of set1")
	}
}

func TestIsSuperset(t *testing.T) {
	set1 := New[int]()
	set1.Add(1)
	set1.Add(2)
	set1.Add(3)

	set2 := New[int]()
	set2.Add(1)
	set2.Add(2)

	if !set1.IsSuperset(set2) {
		t.Errorf("Set1 should be a superset of set2")
	}

	if set2.IsSuperset(set1) {
		t.Errorf("Set2 should not be a superset of set1")
	}
}

func TestIsDisjoint(t *testing.T) {
	set1 := New[int]()
	set1.Add(1)
	set1.Add(2)
	set2 := New[int]()
	set2.Add(3)
	set2.Add(4)
	if !set1.IsDisjoint(set2) {
		t.Errorf("Set1 and set2 should be disjoint")
	}
	set2.Add(1)
	if set1.IsDisjoint(set2) {
		t.Errorf("Set1 and set2 should not be disjoint")
	}
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
