package utils

import "sort"

func Find[T comparable](slice []T, needle T) (int, bool) {
	if len(slice) == 0 {
		return -1, false
	}
	for i := range slice {
		if slice[i] == needle {
			return i, true
		}
	}
	return -1, false
}

func BinarySearch[T any](slice []T, needle T, compare func(i, j T) bool) (int, bool) {
	if len(slice) == 0 {
		return 0, false
	}
	i := sort.Search(len(slice), func(i int) bool {
		return compare(slice[i], needle)
	})
	return i, i < len(slice) && i > 0
}

func RemoveAt[T any](s []T, index int) []T {
	sliceLen := len(s)
	sliceLastIndex := sliceLen - 1

	if index != sliceLastIndex {
		s[index] = s[sliceLastIndex]
	}

	return s[:sliceLastIndex]
}

func RemoveAtPreserveOrder[T any](slice []T, index int) []T {
	return append(slice[:index], slice[index+1:]...)
}
