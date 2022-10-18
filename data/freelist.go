package data

type freeListEntry[T any] struct {
	element  T
	nextFree int
}

// FreeList is a structure that holds any kind of object.
// When an entry is erased, it marks that slot as free so that new entries can be filled in to
// the existing allocated memory.
type FreeList[T any] struct {
	data      []freeListEntry[T]
	FirstFree int
}

func NewFreeList[T any]() FreeList[T] {
	return FreeList[T]{
		FirstFree: -1,
	}
}

func (f *FreeList[T]) Insert(element T) int {
	// if there is a free entry
	if f.FirstFree != -1 {
		// get the index of the free entry
		index := f.FirstFree
		// set the next free entry from this index
		f.FirstFree = f.data[index].nextFree
		f.data[index].nextFree = 0
		// set the data at current free index
		f.data[index].element = element
		return index
	}
	// insert at the end of the list
	f.data = append(f.data, freeListEntry[T]{
		element:  element,
		nextFree: 0,
	})
	return len(f.data) - 1
}

func (f *FreeList[T]) Set(n int, element T) {
	f.data[n].element = element
}

func (f *FreeList[T]) Erase(n int) {
	// set the old next free index to this node
	f.data[n].nextFree = f.FirstFree
	// set the first free index to this nodes next free Raw
	f.FirstFree = n
}

func (f *FreeList[T]) Clear() {
	f.data = f.data[:0]
	f.FirstFree = -1
}

func (f *FreeList[T]) Get(n int) T {
	return f.data[n].element
}

func (f *FreeList[T]) Len() int {
	return len(f.data)
}
