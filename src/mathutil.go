package main

type Ordered interface {
	~float32 | ~float64 | ~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~string
}

func clamp[T Ordered](val, minVal, maxVal T) T {
	if val < minVal {
		return minVal
	} else if val > maxVal {
		return maxVal
	}
	return val
}
