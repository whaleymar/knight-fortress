package math

type Ordered interface {
	~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~string
}

type Signed interface {
	~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64
}

func Clamp[T Ordered](val, minVal, maxVal T) T {
	if val < minVal {
		return minVal
	} else if val > maxVal {
		return maxVal
	}
	return val
}

func Between[T Ordered](val, minVal, maxVal T) bool {
	return val >= minVal && val <= maxVal
}

func Sign[T Signed](val T) T {
	if val < 0 {
		return -1
	}
	return 1
}
