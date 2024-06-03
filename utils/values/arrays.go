package values

import (
	"math/rand"
	"slices"
	"sort"
)

// ReverseArray reverses the array. The original array will not be modified.
//
// example:
//
//	ReverseArray([]int{1, 2, 3}) => []int{3, 2, 1}
//	ReverseArray([]string{"a", "b", "c"}) => []string{"c", "b", "a"}
func ReverseArray[T any](array []T) []T {
	length := len(array)
	if length <= 1 {
		return array
	}

	result := make([]T, length)
	for i, item := range array {
		result[length-i-1] = item
	}

	return result
}

// ContainsArray checks if the target is in the array.
//
// example:
//
//	ContainsArray([]int{1, 2, 3}, 2) => true
//	ContainsArray([]string{"a", "b", "c"}, "d") => false
func ContainsArray[T comparable](array []T, target T) bool {
	return slices.Contains(array, target)
}

// UniqueArray returns a new array with unique elements.
//
// example:
//
//	UniqueArray([]int{1, 2, 2, 3}) => []int{1, 2, 3}
//	UniqueArray([]string{"a", "b", "b", "c"}) => []string{"a", "b", "c"}
func UniqueArray[T comparable](array []T) []T {
	mp := make(map[T]struct{})
	result := make([]T, 0, len(array))
	for _, item := range array {
		if _, ok := mp[item]; !ok {
			mp[item] = struct{}{}
			result = append(result, item)
		}
	}

	return slices.Clip(result)
}

// RemoveArray removes the elements from the array from [start, end).
//
// example:
//
//	RemoveArray([]int{1, 2, 3}, 1, 2) => []int{1}
//	RemoveArray([]string{"a", "b", "c"}, 0, 1) => []string{"c"}
func RemoveArray[T comparable](array []T, start, end int) []T {
	return slices.Delete(array, start, end)
}

// IndexArray returns the index of the target in the array.
//
// example:
//
//	IndexArray([]int{1, 2, 3}, 2) => 1
//	IndexArray([]string{"a", "b", "c"}, "d") => -1
func IndexArray[T comparable](array []T, target T) int {
	return slices.Index(array, target)
}

// ShuffleArray shuffles the array from [start, end).
//
// example:
//
//	ShuffleArray([]int{1, 2, 3}, 0, 3) => []int{3, 1, 2}
//	ShuffleArray([]string{"a", "b", "c"}, 0, 3) => []string{"c", "a", "b"}
func ShuffleArray[T any](array []T, start, end int) []T {
	rand.Shuffle(end-start, func(i, j int) {
		array[start+i], array[start+j] = array[start+j], array[start+i]
	})

	return array
}

// MergeArrays concatenates the arrays.
//
// example:
//
//	MergeArrays([]int{1, 2}, []int{3, 4}) => []int{1, 2, 3, 4}
//	MergeArrays([]string{"a", "b"}, []string{"c", "d"}) => []string{"a", "b", "c", "d"}
func MergeArrays[T any](arrays ...[]T) []T {
	var result []T
	for _, array := range arrays {
		result = append(result, array...)
	}

	return result
}

// FilterArray filters the array with the filter function.
//
// example:
//
//	FilterArray([]int{1, 2, 3}, func(i int) bool { return i > 1 }) => []int{2, 3}
//	FilterArray([]string{"a", "b", "c"}, func(s string) bool { return s != "a" }) => []string{"b", "c"}
func FilterArray[T any](array []T, filter func(T) bool) []T {
	result := make([]T, 0, len(array))
	for _, item := range array {
		if filter(item) {
			result = append(result, item)
		}
	}

	return result
}

// SortArray sorts the array with the compare function. The original array will not be modified.
//
// example:
//
//	SortArray([]int{3, 1, 2}, func(i, j int) bool { return i < j }) => []int{1, 2, 3}
//	SortArray([]string{"c", "a", "b"}, func(i, j string) bool { return i < j }) => []string{"a", "b", "c"}
func SortArray[T any](array []T, cmp func(T, T) bool) []T {
	sorted := make([]T, len(array))
	copy(sorted, array)
	sort.SliceStable(sorted, func(i, j int) bool {
		return cmp(sorted[i], sorted[j])
	})

	return sorted
}

func LastOfArray[T any](array []T) T {
	if len(array) == 0 {
		return Nil[T]()
	}

	return array[len(array)-1]
}

func TopNArray[T any](array []T, percentN int, cmp func(a, b T) bool) (result T) {
	if len(array) == 0 {
		return Nil[T]()
	}

	tempArray := make([]T, len(array))
	copy(tempArray, array)

	n := len(tempArray) * percentN / 100
	if n <= 0 {
		n = 1
	}

	// quick select algorithm
	left, right := 0, len(tempArray)-1
	for left < right {
		pivotIndex := left + (right-left)/2
		tempArray[pivotIndex], tempArray[right] = tempArray[right], tempArray[pivotIndex]

		// partition
		storeIndex := left
		for i := left; i < right; i++ {
			if cmp(tempArray[i], tempArray[right]) {
				tempArray[storeIndex], tempArray[i] = tempArray[i], tempArray[storeIndex]
				storeIndex++
			}
		}

		tempArray[storeIndex], tempArray[right] = tempArray[right], tempArray[storeIndex]

		if storeIndex == n-1 {
			return tempArray[storeIndex]
		} else if storeIndex < n-1 {
			left = storeIndex + 1
		} else {
			right = storeIndex - 1
		}
	}

	return tempArray[left]
}
