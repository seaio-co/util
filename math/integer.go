package math

// 整数限制的值
const (
	MaxInt8   = 1<<7 - 1
	MinInt8   = -1 << 7
	MaxInt16  = 1<<15 - 1
	MinInt16  = -1 << 15
	MaxInt32  = 1<<31 - 1
	MinInt32  = -1 << 31
	MaxInt64  = 1<<63 - 1
	MinInt64  = -1 << 63
	MaxUint8  = 1<<8 - 1
	MaxUint16 = 1<<16 - 1
	MaxUint32 = 1<<32 - 1
	MaxUint64 = 1<<64 - 1
)

// SafeSub 返回减法结果和是否发生溢出
func SafeSub(x, y uint64) (uint64, bool) {
	return x - y, x < y
}

// SafeAdd 返回加法结果和是否发生溢出
func SafeAdd(x, y uint64) (uint64, bool) {
	return x + y, y > MaxUint64-x
}

// SafeMul 返回乘法结果和是否发生溢出。
func SafeMul(x, y uint64) (uint64, bool) {
	if x == 0 || y == 0 {
		return 0, false
	}
	return x * y, y > MaxUint64/x
}

func IntMax(a, b int) int {
	if b > a {
		return b
	}
	return a
}

// IntMin returns the minimum of the params
func IntMin(a, b int) int {
	if b < a {
		return b
	}
	return a
}

// Int32Max returns the maximum of the params
func Int32Max(a, b int32) int32 {
	if b > a {
		return b
	}
	return a
}

// Int32Min returns the minimum of the params
func Int32Min(a, b int32) int32 {
	if b < a {
		return b
	}
	return a
}

// Int64Max returns the maximum of the params
func Int64Max(a, b int64) int64 {
	if b > a {
		return b
	}
	return a
}

// Int64Min returns the minimum of the params
func Int64Min(a, b int64) int64 {
	if b < a {
		return b
	}
	return a
}

// RoundToInt32 rounds floats into integer numbers.
func RoundToInt32(a float64) int32 {
	if a < 0 {
		return int32(a - 0.5)
	}
	return int32(a + 0.5)
}