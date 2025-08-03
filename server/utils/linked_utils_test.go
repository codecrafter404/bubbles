package utils

import (
	"fmt"
	"testing"
)

func TestResolveDependencyNoLoop(t *testing.T) {
	node := []LinkedNode{
		{Id: 1, DependsOn: pointer(2)},
		{Id: 3, DependsOn: nil},
		{Id: 2, DependsOn: pointer(3)},
	}
	tree, successful := LinkedNode{Id: 0, DependsOn: pointer(1)}.ResolveDependency(node, []LinkedNode{})
	if !successful {
		t.Error("Should be successful")
	}

	expected := []LinkedNode{
		{Id: 1, DependsOn: pointer(2)},
		{Id: 2, DependsOn: pointer(3)},
		{Id: 3, DependsOn: nil},
	}
	for y, x := range tree {
		if x.Id != expected[y].Id || (x.DependsOn != nil && expected[y].DependsOn != nil && *x.DependsOn != *expected[y].DependsOn) {
			t.Error(fmt.Sprintf("[%d] %#+v != %+v", y, x, expected[y]))
		}
	}
}
func TestResolveDependencyLoop(t *testing.T) {
	node := []LinkedNode{
		{Id: 1, DependsOn: pointer(2)},
		{Id: 3, DependsOn: pointer(0)},
		{Id: 2, DependsOn: pointer(3)},
	}
	_, successful := LinkedNode{Id: 0, DependsOn: pointer(1)}.ResolveDependency(node, []LinkedNode{})
	if successful {
		t.Error("Shouldn't be successful")
	}
}
func TestResolveDependencyMalformed(t *testing.T) {
	node := []LinkedNode{
		{Id: 1, DependsOn: pointer(2)},
		{Id: 3, DependsOn: pointer(42)},
		{Id: 2, DependsOn: pointer(3)},
	}
	_, successful := LinkedNode{Id: 0, DependsOn: pointer(1)}.ResolveDependency(node, []LinkedNode{})
	if successful {
		t.Error("Shouldn't be successful")
	}
}

func pointer(value int) *int {
	return &value
}
