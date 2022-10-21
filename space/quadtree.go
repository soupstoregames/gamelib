package space

import (
	"github.com/soupstoregames/gamelib/data"
	"github.com/soupstoregames/gamelib/maths"
)

const (
	QuadTreeCapacity = 32
	QuadTreeMaxDepth = 12
)

type QuadTreeEntry struct {
	ID   uint64
	Rect maths.Rectangle
	next int
}

type quadTreeNode struct {
	firstChild int
	count      int
}

// QuadTree is a partitioned space structure for storing rectangles with an associated ID.
// It is used to efficiently search regions of space for elements within.
type QuadTree struct {
	bounds  maths.Rectangle
	nodes   data.FreeList[quadTreeNode]
	entries data.FreeList[QuadTreeEntry]
}

func NewQuadTree(bounds maths.Rectangle) *QuadTree {
	qt := &QuadTree{
		bounds:  bounds,
		nodes:   data.FreeList[quadTreeNode]{FirstFree: -1},
		entries: data.FreeList[QuadTreeEntry]{FirstFree: -1},
	}

	qt.nodes.Insert(quadTreeNode{firstChild: -1})

	return qt
}

func (q *QuadTree) Insert(e QuadTreeEntry) {
	e.next = -1
	q.insert(e, q.bounds, 0, 0)
}

func (q *QuadTree) insert(e QuadTreeEntry, bounds maths.Rectangle, depth int, nodeIndex int) {
	// if the player does not exist within our bounds, do nothing.
	// one of our siblings will accept the player
	if !bounds.Intersects(e.Rect) {
		return
	}

	// get the node
	node := q.nodes.Get(nodeIndex)

	// if this is a branch
	if q.isBranchNode(node) {
		w := bounds.Width / 2
		h := bounds.Height / 2
		q.insert(e, maths.Rectangle{X: bounds.X + w, Y: bounds.Y + h, Width: w, Height: h}, depth+1, node.firstChild)
		q.insert(e, maths.Rectangle{X: bounds.X, Y: bounds.Y + h, Width: w, Height: h}, depth+1, node.firstChild+1)
		q.insert(e, maths.Rectangle{X: bounds.X, Y: bounds.Y, Width: w, Height: h}, depth+1, node.firstChild+2)
		q.insert(e, maths.Rectangle{X: bounds.X + w, Y: bounds.Y, Width: w, Height: h}, depth+1, node.firstChild+3)
	} else {
		// split
		if node.count+1 > QuadTreeCapacity && depth < QuadTreeMaxDepth {
			// save the index of the first Raw element
			currentChild := node.firstChild

			// create a new 4 quad trees and set the first as first index
			node.firstChild = q.nodes.Insert(quadTreeNode{firstChild: -1})
			q.nodes.Insert(quadTreeNode{firstChild: -1})
			q.nodes.Insert(quadTreeNode{firstChild: -1})
			q.nodes.Insert(quadTreeNode{firstChild: -1})

			// split out children into leaf nodes
			w := bounds.Width / 2
			h := bounds.Height / 2

			q.insert(e, maths.Rectangle{X: bounds.X + w, Y: bounds.Y + h, Width: w, Height: h}, depth+1, node.firstChild)
			q.insert(e, maths.Rectangle{X: bounds.X, Y: bounds.Y + h, Width: w, Height: h}, depth+1, node.firstChild+1)
			q.insert(e, maths.Rectangle{X: bounds.X, Y: bounds.Y, Width: w, Height: h}, depth+1, node.firstChild+2)
			q.insert(e, maths.Rectangle{X: bounds.X + w, Y: bounds.Y, Width: w, Height: h}, depth+1, node.firstChild+3)

			// insert this nodes old elements
			for {
				if currentChild == -1 {
					break
				}
				entry := q.entries.Get(currentChild)
				q.insert(entry, maths.Rectangle{X: bounds.X + w, Y: bounds.Y + h, Width: w, Height: h}, depth+1, node.firstChild)
				q.insert(entry, maths.Rectangle{X: bounds.X, Y: bounds.Y + h, Width: w, Height: h}, depth+1, node.firstChild+1)
				q.insert(entry, maths.Rectangle{X: bounds.X, Y: bounds.Y, Width: w, Height: h}, depth+1, node.firstChild+2)
				q.insert(entry, maths.Rectangle{X: bounds.X + w, Y: bounds.Y, Width: w, Height: h}, depth+1, node.firstChild+3)
				q.entries.Erase(currentChild)
				currentChild = entry.next
			}
			node.count = -1
		} else {
			// put element at start of Raw linked list
			index := q.entries.Insert(e)
			entry := q.entries.Get(index)
			entry.next = node.firstChild
			q.entries.Set(index, entry)
			node.firstChild = index
			node.count++
		}

		q.nodes.Set(nodeIndex, node)
	}

	return
}

func (q *QuadTree) Remove(e QuadTreeEntry) {
	q.remove(e, q.bounds, 0)
}

func (q *QuadTree) remove(e QuadTreeEntry, bounds maths.Rectangle, nodeIndex int) {
	// if the entity does not exist within our bounds, do nothing.
	// the entity could never have been in our tree.
	if !bounds.Intersects(e.Rect) {
		return
	}

	// get the node
	node := q.nodes.Get(nodeIndex)

	// find and remove item
	if q.isBranchNode(node) {
		// ask children to remove it
		w := bounds.Width / 2
		h := bounds.Height / 2

		q.remove(e, maths.Rectangle{X: bounds.X + w, Y: bounds.Y + h, Width: w, Height: h}, node.firstChild)
		q.remove(e, maths.Rectangle{X: bounds.X, Y: bounds.Y + h, Width: w, Height: h}, node.firstChild+1)
		q.remove(e, maths.Rectangle{X: bounds.X, Y: bounds.Y, Width: w, Height: h}, node.firstChild+2)
		q.remove(e, maths.Rectangle{X: bounds.X + w, Y: bounds.Y, Width: w, Height: h}, node.firstChild+3)
	} else {
		currentChild := node.firstChild
		lastCheckedIndex := -1
		for {
			if currentChild == -1 {
				break
			}

			// get the child element we're currently checking
			entry := q.entries.Get(currentChild)
			// if the child element has a matching id
			if entry.ID == e.ID {
				// free the element from the list
				q.entries.Erase(currentChild)
				node.count--
				// if the element we removed was the first in the list
				if lastCheckedIndex == -1 {
					// set the new first child to be the second item in the list
					node.firstChild = entry.next
				} else { // otherwise if there was an element before the one that was removed
					// get the last element and update it to point to the removed elements next ptr
					prevEntry := q.entries.Get(lastCheckedIndex)
					prevEntry.next = entry.next
					q.entries.Set(lastCheckedIndex, prevEntry)
				}
				q.nodes.Set(nodeIndex, node)
				return
			}
			// save the index of the last checked element
			lastCheckedIndex = currentChild
			// set the next child to check
			currentChild = entry.next
		}
	}
}

func (q *QuadTree) Scan(results *[]QuadTreeEntry, rect maths.Rectangle) {
	q.scan(results, rect, q.bounds, 0)
}

func (q *QuadTree) scan(results *[]QuadTreeEntry, rect, bounds maths.Rectangle, nodeIndex int) {
	if !bounds.Intersects(rect) {
		return
	}

	node := q.nodes.Get(nodeIndex)

	// find and remove item
	if q.isBranchNode(node) {
		// ask children to search
		w := bounds.Width / 2
		h := bounds.Height / 2

		q.scan(results, rect, maths.Rectangle{X: bounds.X + w, Y: bounds.Y + h, Width: w, Height: h}, node.firstChild)
		q.scan(results, rect, maths.Rectangle{X: bounds.X, Y: bounds.Y + h, Width: w, Height: h}, node.firstChild+1)
		q.scan(results, rect, maths.Rectangle{X: bounds.X, Y: bounds.Y, Width: w, Height: h}, node.firstChild+2)
		q.scan(results, rect, maths.Rectangle{X: bounds.X + w, Y: bounds.Y, Width: w, Height: h}, node.firstChild+3)
	} else {
		currentChild := node.firstChild
		for {
			if currentChild == -1 {
				break
			}

			// get the child element we're currently checking
			entry := q.entries.Get(currentChild)
			// if the child element is in the search rectangle
			if rect.Intersects(entry.Rect) {
				*results = append(*results, entry)
			}
			currentChild = entry.next
		}
	}
}

func (q *QuadTree) CleanUp() {
	var stack []int

	// Only process the root if it's a branch
	if q.isBranchNode(q.nodes.Get(0)) {
		stack = append(stack, 0)
	}

	for {
		if len(stack) == 0 {
			break
		}

		// pop stack
		n := len(stack) - 1
		nodeToProcess := stack[n]
		stack = stack[:n]

		node := q.nodes.Get(nodeToProcess)

		// Loop through the children.
		var numberOfEmptyLeaves int
		for i := 0; i < 4; i++ {
			childIndex := node.firstChild + i
			childNode := q.nodes.Get(childIndex)

			// Increment empty leaf count if the child is an empty
			// leaf. Otherwise if the child is a branch, add it to
			// the stack to be processed in the next iteration.
			if childNode.count == 0 {
				numberOfEmptyLeaves++
			} else if q.isBranchNode(childNode) {
				stack = append(stack, childIndex)
			}
		}

		// If all the children were empty leaves, remove them and
		// make this node the new empty leaf.
		if numberOfEmptyLeaves == 4 {
			// Push all 4 children to the free list.
			for i := 3; i >= 0; i-- {
				q.nodes.Erase(node.firstChild + i)
			}
			// Make this node the new empty leaf.
			node.firstChild = -1
			node.count = 0
			q.nodes.Set(nodeToProcess, node)
		}
	}
}

func (q *QuadTree) Clear() {
	q.nodes.Clear()
	q.entries.Clear()
	q.nodes.Insert(quadTreeNode{firstChild: -1})
}

func (q *QuadTree) isBranchNode(node quadTreeNode) bool {
	return node.count == -1
}

func (q *QuadTree) Walk(f func(i int, r maths.Rectangle, entries []QuadTreeEntry)) {
	q.walk(f, 0, q.bounds, 0)
}

func (q *QuadTree) walk(f func(i int, r maths.Rectangle, entries []QuadTreeEntry), quadrant int, bounds maths.Rectangle, nodeIndex int) {
	node := q.nodes.Get(nodeIndex)

	// find and remove item
	if q.isBranchNode(node) {
		// ask children to search
		w := bounds.Width / 2
		h := bounds.Height / 2

		q.walk(f, 0, maths.Rectangle{X: bounds.X + w, Y: bounds.Y + h, Width: w, Height: h}, node.firstChild)
		q.walk(f, 1, maths.Rectangle{X: bounds.X, Y: bounds.Y + h, Width: w, Height: h}, node.firstChild+1)
		q.walk(f, 2, maths.Rectangle{X: bounds.X, Y: bounds.Y, Width: w, Height: h}, node.firstChild+2)
		q.walk(f, 3, maths.Rectangle{X: bounds.X + w, Y: bounds.Y, Width: w, Height: h}, node.firstChild+3)
	} else {
		var entries []QuadTreeEntry
		currentChild := node.firstChild
		for {
			if currentChild == -1 {
				break
			}

			// get the child element we're currently checking
			entry := q.entries.Get(currentChild)
			// if the child element is in the search rectangle
			entries = append(entries, entry)
			currentChild = entry.next
		}
		f(quadrant, bounds, entries)
	}
}
