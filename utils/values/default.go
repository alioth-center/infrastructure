package values

func Nil[T any]() T {
	var nilValue T
	return nilValue
}

func NilFunction[T any]() func() T {
	return func() T {
		return Nil[T]()
	}
}

func Ptr[T any](val T) *T {
	return &val
}
