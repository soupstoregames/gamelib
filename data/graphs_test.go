package data

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUndirectedGraph_Connect(t *testing.T) {
	graph := NewUndirectedGraph[uint64]()
	graph.Insert(0)
	graph.Insert(1)
	graph.Insert(2)
	graph.Insert(3)

	graph.Connect(0, 1)
	graph.Connect(1, 2)
	graph.Connect(2, 3)

	assert.ElementsMatch(t, []int{1}, graph.Get(0).Adjacent)
	assert.ElementsMatch(t, []int{0, 2}, graph.Get(1).Adjacent)
	assert.ElementsMatch(t, []int{1, 3}, graph.Get(2).Adjacent)
	assert.ElementsMatch(t, []int{2}, graph.Get(3).Adjacent)
}

func TestUndirectedGraph_Merge(t *testing.T) {
	graph := NewUndirectedGraph[uint64]()
	graph.Insert(0)
	graph.Insert(1)
	graph.Insert(2)
	graph.Insert(3)
	graph.Insert(4)

	graph.Connect(0, 1)
	graph.Connect(1, 2)
	graph.Connect(2, 3)
	graph.Connect(2, 4)
	graph.Connect(3, 4)

	graph.Merge(1, 2)

	assert.ElementsMatch(t, []int{1}, graph.Get(0).Adjacent)
	assert.ElementsMatch(t, []int{0, 3, 4}, graph.Get(1).Adjacent)
	assert.ElementsMatch(t, []int{1, 4}, graph.Get(3).Adjacent)
}

func TestUndirectedGraph_Consolidate(t *testing.T) {
	graph := NewUndirectedGraph[uint64]()
	graph.Insert(0)
	graph.Insert(1)
	graph.Insert(2)
	graph.Insert(3)

	graph.Connect(0, 1)
	graph.Connect(1, 2)
	graph.Connect(1, 3)

	/*
		         3
		        ||
			0 = 1 = 2
	*/

	graph.Remove(2)
	changes := graph.Consolidate()

	t.Log(changes)
}
