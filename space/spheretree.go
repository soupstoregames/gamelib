package space

import (
	"github.com/soupstoregames/gamelib/data"
	"github.com/soupstoregames/gamelib/maths"
	"math"
)

// SphereTree is a space partitioning data structure that contains spheres inside larger super spheres
type SphereTree struct {
	spheres       data.FreeList[SphereEntry]
	integrateFIFO data.Queue[int]
	recomputeFIFO data.Queue[int]
	maxSize       float64
	gravy         float64
}

type SphereEntry struct {
	ID         uint64
	sphere     maths.Sphere
	parent     int
	firstChild int
	next       int
}

func NewSphereTree(rootSphere maths.Sphere, maxSize, gravy float64) *SphereTree {
	st := &SphereTree{
		spheres: data.FreeList[SphereEntry]{FirstFree: -1},
		maxSize: maxSize,
		gravy:   gravy,
	}

	st.spheres.Insert(SphereEntry{sphere: rootSphere, firstChild: -1, next: -1})

	return st
}

func (st *SphereTree) Insert(id uint64, sphere maths.Sphere) int {
	entry := SphereEntry{ID: id, sphere: sphere, firstChild: -2, next: -1}
	// insert entry into list
	i := st.spheres.Insert(entry)
	// set the entry to be integrated into the tree, don't do it now
	st.QueueIntegrate(i)
	return i
}

func (st *SphereTree) Remove(entryID int) {
	entry := st.spheres.Get(entryID)
	st.spheres.Erase(entryID)

	// remove this entry from the parent's linked list
	st.removeChild(entry.parent, entryID)
}

func (st *SphereTree) Move(entryID int, sphere maths.Sphere) {
	entry := st.spheres.Get(entryID)
	entry.sphere = sphere
	st.spheres.Set(entryID, entry)

	parentID := entry.parent
	parent := st.spheres.Get(parentID)
	// parent spheres still contains
	if parent.sphere.ContainsSphere(entry.sphere) {
		//st.QueueRecompute(parentID) // TODO: maybe dont
		return
	}

	// entry has broken out of sphere
	// so detach it and queue it for integration
	// and recompute its old parent
	st.removeChild(parentID, entryID)
	st.QueueIntegrate(entryID)
}

func (st *SphereTree) Integrate() {
	for {
		integrateCandidateID, ok := st.integrateFIFO.Pop()
		if !ok {
			break
		}
		st.integrate(integrateCandidateID)
	}
}

func (st *SphereTree) integrate(entryID int) {
	integrateCandidate := st.spheres.Get(entryID)

	containsUs := -1
	nearest := -1
	nearestDist := math.MaxFloat64

	// look through all super spheres for candidate parent
	// look for super spheres that fully contain the candidate or find the closet one
	rootSphere := st.spheres.Get(0) // root
	superSphereIndex := rootSphere.firstChild
	for {
		if superSphereIndex < 0 {
			break
		}
		superSphere := st.spheres.Get(superSphereIndex)
		if superSphere.sphere.ContainsSphere(integrateCandidate.sphere) {
			containsUs = superSphereIndex
			break
		}

		// TODO: Verify this
		dist := superSphere.sphere.Center.Distance(integrateCandidate.sphere.Center) + integrateCandidate.sphere.Radius - superSphere.sphere.Radius
		if dist < nearestDist {
			nearest = superSphereIndex
			nearestDist = dist
		}

		superSphereIndex = superSphere.next
	}

	// if a super sphere contains it, just insert it!
	if containsUs != -1 {
		st.addChild(containsUs, entryID)
		return
	}

	// check to see if the nearest sphere can grow to contain us
	if nearest != -1 {
		parent := st.spheres.Get(nearest)
		newSize := nearestDist + parent.sphere.Radius
		if newSize <= st.maxSize {
			parent.sphere.Radius = newSize + st.gravy
			st.spheres.Set(nearest, parent)
			st.addChild(nearest, entryID)
			return
		}
	}

	// we'll have to make a new super sphere
	parent := SphereEntry{
		sphere:     integrateCandidate.sphere,
		firstChild: -1,
		parent:     0,
	}
	parent.sphere.Radius += st.gravy
	rootSphere = st.spheres.Get(0) // root
	parent.next = rootSphere.firstChild
	parentID := st.spheres.Insert(parent)
	rootSphere.firstChild = parentID
	st.spheres.Set(0, rootSphere)

	st.addChild(parentID, entryID)
}

func (st *SphereTree) Recompute() {
	for {
		recomputeCandidateID, ok := st.recomputeFIFO.Pop()
		if !ok {
			break
		}
		st.recompute(recomputeCandidateID)
	}
}

func (st *SphereTree) recompute(superSphereID int) {
	superSphere := st.spheres.Get(superSphereID)

	if superSphere.firstChild == -1 {
		st.removeChild(0, superSphereID)
		st.spheres.Erase(superSphereID)
		return
	}

	var childCount int
	var total maths.Vector3
	childIndex := superSphere.firstChild
	for {
		if childIndex == -1 {
			break
		}

		childCount++
		child := st.spheres.Get(childIndex)
		total = total.Add(child.sphere.Center)
		childIndex = child.next
	}
	recip := 1.0 / float64(childCount)
	oldCenter := superSphere.sphere.Center
	superSphere.sphere.Center = total.Multiply(recip)

	newRadius := 0.0
	childIndex = superSphere.firstChild
	for {
		if childIndex == -1 {
			break
		}
		child := st.spheres.Get(childIndex)
		radius := superSphere.sphere.Center.Distance(child.sphere.Center) + child.sphere.Radius
		if radius > newRadius {
			newRadius = radius
			if newRadius+st.gravy > superSphere.sphere.Radius {
				superSphere.sphere.Center = oldCenter
				return
			}
		}
		childIndex = child.next
	}
	superSphere.sphere.Radius = newRadius + st.gravy
	st.spheres.Set(superSphereID, superSphere)

	// RINSE: look for circle that can own this one
	root := st.spheres.Get(0)
	possibleParentID := root.firstChild
	for {
		if possibleParentID == 0 {
			panic("no")
		}
		if possibleParentID == -1 {
			break
		}

		if possibleParentID == superSphereID {
			possibleParentID = superSphere.next
			continue
		}

		possibleParent := st.spheres.Get(possibleParentID)
		if possibleParent.sphere.ContainsSphere(superSphere.sphere) {
			// find last child of possible parent

			childIndex := superSphere.firstChild
			for {
				if childIndex == -1 {
					break
				}

				child := st.spheres.Get(childIndex)
				next := child.next
				child.next = possibleParent.firstChild
				child.parent = possibleParentID
				possibleParent.firstChild = childIndex
				st.spheres.Set(childIndex, child)

				childIndex = next
			}

			superSphere.firstChild = -1
			st.spheres.Set(superSphereID, superSphere)
			st.spheres.Set(possibleParentID, possibleParent)
			st.QueueRecompute(superSphereID)
			st.QueueRecompute(possibleParentID)

			break
		}

		possibleParentID = possibleParent.next
	}
}

func (st *SphereTree) Walk(f func(s maths.Sphere, isSuper bool)) {
	root := st.spheres.Get(0)

	childIndex := root.firstChild
	for {
		if childIndex == -1 {
			break
		}
		child := st.spheres.Get(childIndex)
		entryIndex := child.firstChild
		for {
			if entryIndex == -1 {
				break
			}
			entry := st.spheres.Get(entryIndex)
			f(entry.sphere, false)
			entryIndex = entry.next
		}
		childIndex = child.next
	}

	childIndex = root.firstChild
	for {
		if childIndex == -1 {
			break
		}
		child := st.spheres.Get(childIndex)
		f(child.sphere, true)
		childIndex = child.next
	}
}

func (st *SphereTree) QueueIntegrate(entryID int) {
	for i := 0; i < st.integrateFIFO.Len(); i++ {
		if st.integrateFIFO.Peek(i) == entryID {
			return
		}
	}
	st.integrateFIFO.Push(entryID)
}

func (st *SphereTree) QueueRecompute(parentID int) {
	if parentID == 0 {
		return
	}
	for i := 0; i < st.recomputeFIFO.Len(); i++ {
		if st.recomputeFIFO.Peek(i) == parentID {
			return
		}
	}
	st.recomputeFIFO.Push(parentID)
}

func (st *SphereTree) Scan(selectionSphere maths.Sphere) []SphereEntry {
	var results []SphereEntry

	root := st.spheres.Get(0)
	sphereIndex := root.firstChild

	for {
		if sphereIndex == -1 {
			break
		}
		sphere := st.spheres.Get(sphereIndex)
		if selectionSphere.IntersectsSphere(sphere.sphere) {
			childIndex := sphere.firstChild
			for {
				if childIndex == -1 {
					break
				}

				child := st.spheres.Get(childIndex)

				if selectionSphere.IntersectsSphere(child.sphere) {
					results = append(results, child)
				}

				childIndex = child.next
			}
		}

		sphereIndex = sphere.next
	}

	return results
}

func (st *SphereTree) addChild(parentID, sphereID int) {
	parent := st.spheres.Get(parentID)
	entry := st.spheres.Get(sphereID)

	entry.parent = parentID
	entry.next = parent.firstChild
	parent.firstChild = sphereID

	st.spheres.Set(parentID, parent)
	st.spheres.Set(sphereID, entry)

	st.QueueRecompute(parentID)
}

func (st *SphereTree) removeChild(parentID, entryID int) {
	parent := st.spheres.Get(parentID)
	entry := st.spheres.Get(entryID)

	childIndex := parent.firstChild
	if childIndex == entryID {
		parent.firstChild = entry.next
	} else {
		for {
			if childIndex == -1 {
				break
			}
			child := st.spheres.Get(childIndex)
			if child.next == entryID {
				child.next = entry.next
				st.spheres.Set(childIndex, child)
				break
			}
			childIndex = child.next
		}
	}

	st.spheres.Set(parentID, parent)
	st.spheres.Set(entryID, entry)

	st.QueueRecompute(parentID)
}
