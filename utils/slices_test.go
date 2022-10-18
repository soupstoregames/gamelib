package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBinarySearch(t *testing.T) {
	s := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	_, found := BinarySearch(s, 10, func(i, j uint64) bool { return i > j })
	assert.False(t, found)

	i, found := BinarySearch(s, 3, func(i, j uint64) bool { return i > j })
	assert.Equal(t, 4, i)
	assert.True(t, found)
}

func TestRemoveAt(t *testing.T) {
	s := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s = RemoveAt(s, 9)
	s = RemoveAt(s, 0)
	s = RemoveAt(s, 4)

	assert.Equal(t, []uint64{8, 1, 2, 3, 7, 5, 6}, s)
}

func TestRemoveAtPreserveOrder(t *testing.T) {
	s := []uint64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	s = RemoveAtPreserveOrder(s, 9)
	s = RemoveAtPreserveOrder(s, 0)
	s = RemoveAtPreserveOrder(s, 4)

	assert.Equal(t, []uint64{1, 2, 3, 4, 6, 7, 8}, s)
}
