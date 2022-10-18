package data

import (
	"github.com/soupstoregames/gamelib/maths"
	"github.com/soupstoregames/gamelib/utils"
)

type GraphNode[T any] struct {
	Element  T
	Adjacent []int
}

// UndirectedGraph is a graph of nodes with two-way connections.
type UndirectedGraph[T any] struct {
	data FreeList[GraphNode[T]]
}

func NewUndirectedGraph[T any]() *UndirectedGraph[T] {
	freeList := NewFreeList[GraphNode[T]]()

	return &UndirectedGraph[T]{
		data: freeList,
	}
}

func (u *UndirectedGraph[T]) Insert(element T) int {
	return u.data.Insert(GraphNode[T]{
		Element: element,
	})
}

func (u *UndirectedGraph[T]) Set(idx int, element T) {
	node := u.data.Get(idx)
	node.Element = element
	u.data.Set(idx, node)
}

func (u *UndirectedGraph[T]) Remove(id int) {
	adj := u.data.Get(id).Adjacent
	for i := range adj {
		u.Disconnect(id, adj[i])
	}
	u.data.Erase(id)
}

func (u *UndirectedGraph[T]) Get(id int) GraphNode[T] {
	return u.data.Get(id)
}

func (u *UndirectedGraph[T]) Clear() {
	u.data.Clear()
	u.data.Truncate(0)
}

func (u *UndirectedGraph[T]) Connect(a, b int) {
	aNode := u.data.Get(a)
	if _, found := utils.Find(aNode.Adjacent, b); !found {
		aNode.Adjacent = append(aNode.Adjacent, b)
		u.data.Set(a, aNode)
	}

	bNode := u.data.Get(b)
	if _, found := utils.Find(bNode.Adjacent, a); !found {
		bNode.Adjacent = append(bNode.Adjacent, a)
		u.data.Set(b, bNode)
	}
}

func (u *UndirectedGraph[T]) Disconnect(a, b int) {
	aNode := u.data.Get(a)
	if i, found := utils.Find(aNode.Adjacent, b); found {
		aNode.Adjacent = utils.RemoveAt(aNode.Adjacent, i)
		u.data.Set(a, aNode)
	}

	bNode := u.data.Get(b)
	if i, found := utils.Find(bNode.Adjacent, a); found {
		bNode.Adjacent = utils.RemoveAt(bNode.Adjacent, i)
		u.data.Set(b, bNode)
	}
}

// Merge adds b's connections into a and removes b
func (u *UndirectedGraph[T]) Merge(a, b int) {
	u.Disconnect(a, b)

	bConnections := u.Get(b).Adjacent
	for _, n := range bConnections {
		u.Disconnect(b, n)
		u.Connect(a, n)
	}

	u.data.Erase(b)
}

func (u *UndirectedGraph[T]) Consolidate() []maths.Tuple2[int] {
	changes, newLength := u.data.Consolidate()

	// update the graph connections to point to the new positions
	for i := newLength - 1; i >= 0; i-- {
		node := u.data.Get(i)
		for _, c := range changes {
			idx, ok := utils.Find(node.Adjacent, c.A)
			if ok {
				node.Adjacent[idx] = c.B
			}
		}
		u.data.Set(i, node)
	}

	u.data.Truncate(newLength)

	return changes
}

func (u *UndirectedGraph[T]) Len() int {
	return u.data.Len()
}
