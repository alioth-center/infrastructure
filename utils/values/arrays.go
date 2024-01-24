package values

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
