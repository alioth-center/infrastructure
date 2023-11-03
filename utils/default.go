package utils

func NilValue[T any]() T {
	var nilValue T
	return nilValue
}

func NilFunction[T any]() func() T {
	return func() T {
		return NilValue[T]()
	}
}
