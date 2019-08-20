package math

// Min return the min one.
func Min(x, y int64) int64 {
	if x <= y {
		return x
	}
	return y
}

// Max return the max one.
func Max(x, y int64) int64 {
	if x >= y {
		return x
	}
	return y
}
