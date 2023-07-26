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

	maxBranchSize float64
	maxLeafSize   float64
	gravy         float64
}

type SphereEntry struct {
	ID     uint64
	Sphere maths.Sphere

	parent     int32
	firstChild int32
	next       int32
	flags      data.Bitfield1[uint8]
}

func NewSphereTree(center maths.Vector3, maxBranchSize, maxLeafSize, gravy float64) *SphereTree {
	st := &SphereTree{
		spheres:       data.FreeList[SphereEntry]{FirstFree: -1},
		maxBranchSize: maxBranchSize,
		maxLeafSize:   maxLeafSize,
		gravy:         gravy,
	}

	branchRoot := SphereEntry{Sphere: maths.Sphere{Center: center, Radius: math.MaxFloat64}, firstChild: -1, next: -1}
	branchRoot.flags.Set(SPFRootNode)
	st.spheres.Insert(branchRoot)

	return st
}

func (st *SphereTree) Insert(id uint64, sphere maths.Sphere) int {
	entry := SphereEntry{ID: id, Sphere: sphere, firstChild: -1, next: -1}
	i := st.spheres.Insert(entry)
	st.queueIntegrate(i)
	return i
}

func (st *SphereTree) Remove(entryID int) {
	entry := st.spheres.Get(entryID)
	st.removeChild(int(entry.parent), entryID)
	st.spheres.Erase(entryID)
}

func (st *SphereTree) Move(entryID int, sphere maths.Sphere) {
	entry := st.spheres.Get(entryID)
	entry.Sphere = sphere
	st.spheres.Set(entryID, entry)

	parentID := entry.parent
	parent := st.spheres.Get(int(parentID))
	if parent.Sphere.ContainsSphere(entry.Sphere) {
		return
	}

	// entry has broken out of Sphere
	// so detach it and queue it for integration
	// and recompute its old parent
	st.removeChild(int(parentID), entryID)
	st.queueIntegrate(entryID)
	st.queueRecompute(int(parentID))
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

func (st *SphereTree) Recompute() {
	for {
		recomputeCandidateID, ok := st.recomputeFIFO.Pop()
		if !ok {
			break
		}
		st.recompute(recomputeCandidateID)
	}
}

func (st *SphereTree) Walk(f func(s maths.Sphere, level int)) {
	root := st.spheres.Get(0)

	branchID := root.firstChild
	for {
		if branchID == -1 {
			break
		}
		branch := st.spheres.Get(int(branchID))
		leafID := branch.firstChild
		for {
			if leafID == -1 {
				break
			}
			leaf := st.spheres.Get(int(leafID))
			f(leaf.Sphere, 1)

			leafID = leaf.next
		}
		f(branch.Sphere, 0)
		branchID = branch.next
	}
}
func (st *SphereTree) Scan(entries *[]SphereEntry, selectionSphere maths.Sphere) {
	root := st.spheres.Get(0)
	branchID := root.firstChild

	for {
		if branchID == -1 {
			break
		}
		branch := st.spheres.Get(int(branchID))
		if selectionSphere.IntersectsSphere(branch.Sphere) {
			leafID := branch.firstChild
			for {
				if leafID == -1 {
					break
				}

				leaf := st.spheres.Get(int(leafID))
				if selectionSphere.IntersectsSphere(leaf.Sphere) {
					entryID := leaf.firstChild
					for {
						if entryID == -1 {
							break
						}
						entry := st.spheres.Get(int(entryID))
						if selectionSphere.IntersectsSphere(entry.Sphere) {
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

func (st *SphereTree) queueIntegrate(entryID int) {
	entry := st.spheres.Get(entryID)

	if entry.flags.Has(SPFIntegrate) {
		return
	}

	entry.flags.Set(SPFIntegrate)
	st.spheres.Set(entryID, entry)
	st.integrateFIFO.Push(entryID)
}

func (st *SphereTree) queueRecompute(entryID int) {
	entry := st.spheres.Get(entryID)

	if entry.flags.Has(SPFRootNode) {
		return
	}
	if entry.flags.Has(SPFRecompute) {
		return
	}

	entry.flags.Set(SPFRecompute)
	st.spheres.Set(entryID, entry)
	st.recomputeFIFO.Push(entryID)
}

func (st *SphereTree) integrate(entryID int) {
	entry := st.spheres.Get(entryID)
	entry.flags.Clear(SPFIntegrate)
	st.spheres.Set(entryID, entry)

	// look through all super spheres for candidate parent
	// look for super spheres that fully contain the candidate or find the closet one
	root := st.spheres.Get(0) // root
	sphere := entry.Sphere
	sphere.Radius += st.gravy
	branchContainsID, branchNearestID, branchNearestDist := st.findContainsOrNearest(root, sphere)

	// if a branch contains us then look for a leaf to hold us
	if branchContainsID != -1 {
		branchContains := st.spheres.Get(branchContainsID)
		leafContainsID, leafNearestID, leafNearestDist := st.findContainsOrNearest(branchContains, entry.Sphere)

		// if a leaf fully contains us, then just add
		if leafContainsID >= 0 {
			st.addChild(leafContainsID, entryID)
			return
		}

		// if there's a leaf that can grow to contain us
		if leafNearestID != -1 {
			leafNearest := st.spheres.Get(leafNearestID)
			newSize := leafNearestDist + leafNearest.Sphere.Radius
			if newSize <= st.maxLeafSize {
				leafNearest.Sphere.Radius = newSize + st.gravy
				st.spheres.Set(leafNearestID, leafNearest)
				st.addChild(leafNearestID, entryID)
				return
			}
		}

		// make a new leaf
		newSphere := entry.Sphere
		newSphere.Radius += st.gravy
		newLeafID := st.spheres.Insert(SphereEntry{Sphere: newSphere, firstChild: -1, next: -1})
		st.addChild(branchContainsID, newLeafID)
		st.addChild(newLeafID, entryID)
		return
	}

	// check to see if the branchNearestID Sphere can grow to contain us
	if branchNearestID != -1 {
		branchNearest := st.spheres.Get(branchNearestID)
		newSize := branchNearestDist + branchNearest.Sphere.Radius
		if newSize <= st.maxBranchSize {
			branchNearest.Sphere.Radius = newSize + st.gravy
			st.spheres.Set(branchNearestID, branchNearest)

			leafContainsID, leafNearestID, leafNearestDist := st.findContainsOrNearest(branchNearest, entry.Sphere)

			if leafContainsID != -1 {
				// this shouldn't be possible
				//st.addChild(leafContainsID, entryID)
				//return
			}

			if leafNearestID != -1 {
				leafNearest := st.spheres.Get(leafNearestID)
				newSize := leafNearestDist + leafNearest.Sphere.Radius
				if newSize <= st.maxLeafSize {
					leafNearest.Sphere.Radius = newSize + st.gravy
					st.spheres.Set(leafNearestID, leafNearest)
					st.addChild(leafNearestID, entryID)
					st.queueRecompute(branchNearestID)
					return
				}
			}

			// make a new leaf
			newSphere := entry.Sphere
			newSphere.Radius += st.gravy
			leafID := st.spheres.Insert(SphereEntry{Sphere: newSphere, firstChild: -1, next: -1})
			st.addChild(branchNearestID, leafID)
			st.addChild(leafID, entryID)
			st.queueRecompute(branchNearestID)
			return
		}
	}

	// we'll have to make a new branch
	newSphere := entry.Sphere
	newSphere.Radius += st.gravy
	branchID := st.spheres.Insert(SphereEntry{Sphere: newSphere, firstChild: -1, next: -1})
	leafID := st.spheres.Insert(SphereEntry{Sphere: newSphere, firstChild: -1, next: -1})

	st.addChild(leafID, entryID)
	st.addChild(branchID, leafID)
	st.addChild(0, branchID)
}

func (st *SphereTree) findContainsOrNearest(rootSphere SphereEntry, sphere maths.Sphere) (int, int, float64) {
	containsID := -1
	nearestID := -1
	nearestDist := math.MaxFloat64
	superSphereIndex := rootSphere.firstChild
	for {
		if superSphereIndex < 0 {
			break
		}
		superSphere := st.spheres.Get(int(superSphereIndex))
		if superSphere.Sphere.ContainsSphere(sphere) {
			// TODO: choose nearestID that can contain us
			containsID = int(superSphereIndex)
			break
		}

		// TODO: Verify this
		dist := superSphere.Sphere.Center.Distance(sphere.Center) + sphere.Radius - superSphere.Sphere.Radius // RINSE
		if dist < nearestDist {
			nearestID = int(superSphereIndex)
			nearestDist = dist
		}

		superSphereIndex = superSphere.next
	}
	return containsID, nearestID, nearestDist
}

func (st *SphereTree) recompute(superSphereID int) {
	superSphere := st.spheres.Get(superSphereID)

	if superSphere.firstChild == -1 {
		st.removeChild(int(superSphere.parent), superSphereID)
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
		child := st.spheres.Get(int(childIndex))
		total = total.Add(child.Sphere.Center)
		childIndex = child.next
	}
	recip := 1.0 / float64(childCount)
	oldCenter := superSphere.Sphere.Center
	superSphere.Sphere.Center = total.Multiply(recip)

	newRadius := 0.0
	childIndex = superSphere.firstChild
	for {
		if childIndex == -1 {
			break
		}
		child := st.spheres.Get(int(childIndex))
		radius := superSphere.Sphere.Center.Distance(child.Sphere.Center) + child.Sphere.Radius
		if radius > newRadius {
			newRadius = radius
			if newRadius+st.gravy > superSphere.Sphere.Radius {
				superSphere.Sphere.Center = oldCenter
				superSphere.flags.Clear(SPFRecompute)
				st.spheres.Set(superSphereID, superSphere)
				return
			}
		}
		childIndex = child.next
	}
	superSphere.Sphere.Radius = newRadius + st.gravy
	superSphere.flags.Clear(SPFRecompute)
	st.spheres.Set(superSphereID, superSphere)
}

func (st *SphereTree) addChild(parentID, sphereID int) {
	parent := st.spheres.Get(parentID)
	entry := st.spheres.Get(sphereID)

	entry.parent = int32(parentID)
	entry.next = parent.firstChild
	parent.firstChild = int32(sphereID)

	st.spheres.Set(parentID, parent)
	st.spheres.Set(sphereID, entry)

	st.queueRecompute(parentID)
}

func (st *SphereTree) removeChild(parentID, entryID int) {
	parent := st.spheres.Get(parentID)
	entry := st.spheres.Get(entryID)

	childIndex := int(parent.firstChild)
	if childIndex == entryID {
		parent.firstChild = entry.next
	} else {
		for {
			if childIndex == -1 {
				break
			}
			child := st.spheres.Get(childIndex)
			if int(child.next) == entryID {
				child.next = entry.next
				st.spheres.Set(childIndex, child)
				break
			}
			childIndex = int(child.next)
		}
	}

	st.spheres.Set(parentID, parent)
	st.spheres.Set(entryID, entry)

	st.queueRecompute(parentID)
}
