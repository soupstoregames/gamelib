package space

import (
	"github.com/soupstoregames/gamelib/data"
	"github.com/soupstoregames/gamelib/maths"
	"math"
)

// CircleTree is a space partitioning data structure that contains circles inside larger super circles
type CircleTree struct {
	circles       data.FreeList[CircleEntry]
	integrateFIFO data.Queue[int]
	recomputeFIFO data.Queue[int]

	maxBranchSize float64
	maxLeafSize   float64
	gravy         float64
}

const (
	SPFRootNode = iota
	SPFIntegrate
	SPFRecompute
)

type CircleEntry struct {
	ID     uint64
	Circle maths.Circle

	parent     int32
	firstChild int32
	next       int32
	flags      data.Bitfield1[uint8]
}

func NewCircleTree(center maths.Vector2, maxBranchSize, maxLeafSize, gravy float64) *CircleTree {
	st := &CircleTree{
		circles:       data.FreeList[CircleEntry]{FirstFree: -1},
		maxBranchSize: maxBranchSize,
		maxLeafSize:   maxLeafSize,
		gravy:         gravy,
	}

	branchRoot := CircleEntry{Circle: maths.Circle{Center: center, Radius: math.MaxFloat64}, firstChild: -1, next: -1}
	branchRoot.flags.Set(SPFRootNode)
	st.circles.Insert(branchRoot)

	return st
}

func (st *CircleTree) Insert(id uint64, Circle maths.Circle) int {
	entry := CircleEntry{ID: id, Circle: Circle, firstChild: -1, next: -1}
	i := st.circles.Insert(entry)
	st.queueIntegrate(i)
	return i
}

func (st *CircleTree) Remove(entryID int) {
	entry := st.circles.Get(entryID)
	st.removeChild(int(entry.parent), entryID)
	st.circles.Erase(entryID)
}

func (st *CircleTree) Move(entryID int, Circle maths.Circle) {
	entry := st.circles.Get(entryID)
	entry.Circle = Circle
	st.circles.Set(entryID, entry)

	parentID := entry.parent
	parent := st.circles.Get(int(parentID))
	if parent.Circle.ContainsCircle(entry.Circle) {
		return
	}

	// entry has broken out of Circle
	// so detach it and queue it for integration
	// and recompute its old parent
	st.removeChild(int(parentID), entryID)
	st.queueIntegrate(entryID)
	st.queueRecompute(int(parentID))
}

func (st *CircleTree) Integrate() {
	for {
		integrateCandidateID, ok := st.integrateFIFO.Pop()
		if !ok {
			break
		}
		st.integrate(integrateCandidateID)
	}
}

func (st *CircleTree) Recompute() {
	for {
		recomputeCandidateID, ok := st.recomputeFIFO.Pop()
		if !ok {
			break
		}
		st.recompute(recomputeCandidateID)
	}
}

func (st *CircleTree) Walk(f func(s maths.Circle, level int)) {
	root := st.circles.Get(0)

	branchID := root.firstChild
	for {
		if branchID == -1 {
			break
		}
		branch := st.circles.Get(int(branchID))
		leafID := branch.firstChild
		for {
			if leafID == -1 {
				break
			}
			leaf := st.circles.Get(int(leafID))
			f(leaf.Circle, 1)

			leafID = leaf.next
		}
		f(branch.Circle, 0)
		branchID = branch.next
	}
}
func (st *CircleTree) Scan(entries *[]CircleEntry, selectionCircle maths.Circle) {
	root := st.circles.Get(0)
	branchID := root.firstChild

	for {
		if branchID == -1 {
			break
		}
		branch := st.circles.Get(int(branchID))
		if selectionCircle.IntersectsCircle(branch.Circle) {
			leafID := branch.firstChild
			for {
				if leafID == -1 {
					break
				}

				leaf := st.circles.Get(int(leafID))
				if selectionCircle.IntersectsCircle(leaf.Circle) {
					entryID := leaf.firstChild
					for {
						if entryID == -1 {
							break
						}
						entry := st.circles.Get(int(entryID))
						if selectionCircle.IntersectsCircle(entry.Circle) {
							*entries = append(*entries, entry)
						}
						entryID = entry.next
					}
				}

				leafID = leaf.next
			}
		}

		branchID = branch.next
	}
}

func (st *CircleTree) queueIntegrate(entryID int) {
	entry := st.circles.Get(entryID)

	if entry.flags.Has(SPFIntegrate) {
		return
	}

	entry.flags.Set(SPFIntegrate)
	st.circles.Set(entryID, entry)
	st.integrateFIFO.Push(entryID)
}

func (st *CircleTree) queueRecompute(entryID int) {
	entry := st.circles.Get(entryID)

	if entry.flags.Has(SPFRootNode) {
		return
	}
	if entry.flags.Has(SPFRecompute) {
		return
	}

	entry.flags.Set(SPFRecompute)
	st.circles.Set(entryID, entry)
	st.recomputeFIFO.Push(entryID)
}

func (st *CircleTree) integrate(entryID int) {
	entry := st.circles.Get(entryID)
	entry.flags.Clear(SPFIntegrate)
	st.circles.Set(entryID, entry)

	// look through all super circles for candidate parent
	// look for super circles that fully contain the candidate or find the closet one
	root := st.circles.Get(0) // root
	Circle := entry.Circle
	Circle.Radius += st.gravy
	branchContainsID, branchNearestID, branchNearestDist := st.findContainsOrNearest(root, Circle)

	// if a branch contains us then look for a leaf to hold us
	if branchContainsID != -1 {
		branchContains := st.circles.Get(branchContainsID)
		leafContainsID, leafNearestID, leafNearestDist := st.findContainsOrNearest(branchContains, entry.Circle)

		// if a leaf fully contains us, then just add
		if leafContainsID >= 0 {
			st.addChild(leafContainsID, entryID)
			return
		}

		// if there's a leaf that can grow to contain us
		if leafNearestID != -1 {
			leafNearest := st.circles.Get(leafNearestID)
			newSize := leafNearestDist + leafNearest.Circle.Radius
			if newSize <= st.maxLeafSize {
				leafNearest.Circle.Radius = newSize + st.gravy
				st.circles.Set(leafNearestID, leafNearest)
				st.addChild(leafNearestID, entryID)
				return
			}
		}

		// make a new leaf
		newCircle := entry.Circle
		newCircle.Radius += st.gravy
		newLeafID := st.circles.Insert(CircleEntry{Circle: newCircle, firstChild: -1, next: -1})
		st.addChild(branchContainsID, newLeafID)
		st.addChild(newLeafID, entryID)
		return
	}

	// check to see if the branchNearestID Circle can grow to contain us
	if branchNearestID != -1 {
		branchNearest := st.circles.Get(branchNearestID)
		newSize := branchNearestDist + branchNearest.Circle.Radius
		if newSize <= st.maxBranchSize {
			branchNearest.Circle.Radius = newSize + st.gravy
			st.circles.Set(branchNearestID, branchNearest)

			leafContainsID, leafNearestID, leafNearestDist := st.findContainsOrNearest(branchNearest, entry.Circle)

			if leafContainsID != -1 {
				// this shouldn't be possible
				//st.addChild(leafContainsID, entryID)
				//return
			}

			if leafNearestID != -1 {
				leafNearest := st.circles.Get(leafNearestID)
				newSize := leafNearestDist + leafNearest.Circle.Radius
				if newSize <= st.maxLeafSize {
					leafNearest.Circle.Radius = newSize + st.gravy
					st.circles.Set(leafNearestID, leafNearest)
					st.addChild(leafNearestID, entryID)
					st.queueRecompute(branchNearestID)
					return
				}
			}

			// make a new leaf
			newCircle := entry.Circle
			newCircle.Radius += st.gravy
			leafID := st.circles.Insert(CircleEntry{Circle: newCircle, firstChild: -1, next: -1})
			st.addChild(branchNearestID, leafID)
			st.addChild(leafID, entryID)
			st.queueRecompute(branchNearestID)
			return
		}
	}

	// we'll have to make a new branch
	newCircle := entry.Circle
	newCircle.Radius += st.gravy
	branchID := st.circles.Insert(CircleEntry{Circle: newCircle, firstChild: -1, next: -1})
	leafID := st.circles.Insert(CircleEntry{Circle: newCircle, firstChild: -1, next: -1})

	st.addChild(leafID, entryID)
	st.addChild(branchID, leafID)
	st.addChild(0, branchID)
}

func (st *CircleTree) findContainsOrNearest(rootCircle CircleEntry, Circle maths.Circle) (int, int, float64) {
	containsID := -1
	nearestID := -1
	nearestDist := math.MaxFloat64
	superCircleIndex := rootCircle.firstChild
	for {
		if superCircleIndex < 0 {
			break
		}
		superCircle := st.circles.Get(int(superCircleIndex))
		if superCircle.Circle.ContainsCircle(Circle) {
			// TODO: choose nearestID that can contain us
			containsID = int(superCircleIndex)
			break
		}

		// TODO: Verify this
		dist := superCircle.Circle.Center.Distance(Circle.Center) + Circle.Radius - superCircle.Circle.Radius // RINSE
		if dist < nearestDist {
			nearestID = int(superCircleIndex)
			nearestDist = dist
		}

		superCircleIndex = superCircle.next
	}
	return containsID, nearestID, nearestDist
}

func (st *CircleTree) recompute(superCircleID int) {
	superCircle := st.circles.Get(superCircleID)

	if superCircle.firstChild == -1 {
		st.removeChild(int(superCircle.parent), superCircleID)
		st.circles.Erase(superCircleID)
		return
	}

	var childCount int
	var total maths.Vector2
	childIndex := superCircle.firstChild
	for {
		if childIndex == -1 {
			break
		}

		childCount++
		child := st.circles.Get(int(childIndex))
		total = total.Add(child.Circle.Center)
		childIndex = child.next
	}
	recip := 1.0 / float64(childCount)
	oldCenter := superCircle.Circle.Center
	superCircle.Circle.Center = total.Multiply(recip)

	newRadius := 0.0
	childIndex = superCircle.firstChild
	for {
		if childIndex == -1 {
			break
		}
		child := st.circles.Get(int(childIndex))
		radius := superCircle.Circle.Center.Distance(child.Circle.Center) + child.Circle.Radius
		if radius > newRadius {
			newRadius = radius
			if newRadius+st.gravy > superCircle.Circle.Radius {
				superCircle.Circle.Center = oldCenter
				superCircle.flags.Clear(SPFRecompute)
				st.circles.Set(superCircleID, superCircle)
				return
			}
		}
		childIndex = child.next
	}
	superCircle.Circle.Radius = newRadius + st.gravy
	superCircle.flags.Clear(SPFRecompute)
	st.circles.Set(superCircleID, superCircle)
}

func (st *CircleTree) addChild(parentID, CircleID int) {
	parent := st.circles.Get(parentID)
	entry := st.circles.Get(CircleID)

	entry.parent = int32(parentID)
	entry.next = parent.firstChild
	parent.firstChild = int32(CircleID)

	st.circles.Set(parentID, parent)
	st.circles.Set(CircleID, entry)

	st.queueRecompute(parentID)
}

func (st *CircleTree) removeChild(parentID, entryID int) {
	parent := st.circles.Get(parentID)
	entry := st.circles.Get(entryID)

	childIndex := int(parent.firstChild)
	if childIndex == entryID {
		parent.firstChild = entry.next
	} else {
		for {
			if childIndex == -1 {
				break
			}
			child := st.circles.Get(childIndex)
			if int(child.next) == entryID {
				child.next = entry.next
				st.circles.Set(childIndex, child)
				break
			}
			childIndex = int(child.next)
		}
	}

	st.circles.Set(parentID, parent)
	st.circles.Set(entryID, entry)

	st.queueRecompute(parentID)
}
