package values

import (
	"math/rand"
	"slices"
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
