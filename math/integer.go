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
