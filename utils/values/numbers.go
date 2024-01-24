package values

import "strconv"

func IntToString[T ~int](raw T) string {
	return strconv.Itoa(int(raw))
}

func Int64ToString[T ~int64](raw T) string {
	return strconv.FormatInt(int64(raw), 10)
}

func UintToString[T ~uint](raw T) string {
	return strconv.Itoa(int(raw))
}

func Uint64ToString[T ~uint64](raw T) string {
	return strconv.FormatUint(uint64(raw), 10)
}

func Float32ToString[T ~float32](raw T) string {
	return strconv.FormatFloat(float64(raw), 'f', -1, 32)
}

func Float64ToString[T ~float64](raw T) string {
	return strconv.FormatFloat(float64(raw), 'f', -1, 64)
}
