// Package math 关于 math integer 的数学工具包
package math

import (
	"crypto/rand"
	"github.com/seaio-co/util/conv"
	"math/big"
)

// Various 各种大整数极限值
var (
	tt255     = BigPow(2, 255)
	tt256     = BigPow(2, 256)
	tt256m1   = new(big.Int).Sub(tt256, big.NewInt(1))
	MaxBig256 = new(big.Int).Set(tt256m1)
	tt63      = BigPow(2, 63)
	MaxBig63  = new(big.Int).Sub(tt63, big.NewInt(1))
	Big0      = big.NewInt(0)
	Big1      = big.NewInt(1)
	Big2      = big.NewInt(2)
	Big3      = big.NewInt(3)
	Big32     = big.NewInt(32)
	Big256    = big.NewInt(0xff)
	Big257    = big.NewInt(257)
)

const (
	// 一个big.Word的位数
	wordBits = 32 << (uint64(^big.Word(0)) >> 63)
	// 一个big.Word的字节数
	wordBytes = wordBits / 8
)

// BigPow 返回一个指向big.Int类型的指针地址的指针
func BigPow(a, b int64) *big.Int {
	r := big.NewInt(a)
	return r.Exp(r, big.NewInt(b), nil)
}

// BigMax 返回较大的一个指针
func BigMax(x, y *big.Int) *big.Int {
	if x.Cmp(y) < 0 {
		return y
	}
	return x
}

// BigMin 返回较小的一个指针
func BigMin(x, y *big.Int) *big.Int {
	if x.Cmp(y) > 0 {
		return y
	}
	return x
}

// PaddedBigBytes 将一个大整数编码为一个大端字节切片。
// 这个片长度至少有n个字节
func PaddedBigBytes(bigint *big.Int, n int) []byte {
	if bigint.BitLen()/8 >= n {
		return bigint.Bytes()
	}
	ret := make([]byte, n)
	ReadBits(bigint, ret)
	return ret
}

// bigEndianByteAt 返回位置n的字节，n==0返回最小有效字节
func bigEndianByteAt(bigint *big.Int, n int) byte {
	words := bigint.Bits()
	// 检查字节将驻留在的 word-bucket
	i := n / wordBytes
	if i >= len(words) {
		return byte(0)
	}
	word := words[i]
	// 字节偏移量
	shift := 8 * uint(n%wordBytes)

	return byte(word >> shift)
}

// Byte 返回位置n的字节 例：bigint '5', padlength 32, n=31 => 5
func Byte(bigint *big.Int, padlength, n int) byte {
	if n >= padlength {
		return byte(0)
	}
	return bigEndianByteAt(bigint, padlength-1-n)
}

// ReadBits 将bigint的绝对值编码为大端字节。调用者必须确保buf有足够的空间。如果buf太短，结果将是不完整的。
func ReadBits(bigint *big.Int, buf []byte) {
	i := len(buf)
	for _, d := range bigint.Bits() {
		for j := 0; j < wordBytes && i > 0; j++ {
			i--
			buf[i] = byte(d)
			d >>= 8
		}
	}
}

// U256 编码为256位2的补码。这种操作是破坏性的。
func U256(x *big.Int) *big.Int {
	return x.And(x, tt256m1)
}

// S256 将x解释为一个2的补码。x不能超过256位(如果超过256位，结果是未定义的)，也不能修改。
//   S256(0)        = 0
//   S256(1)        = 1
//   S256(2**255)   = -2**255
//   S256(2**256-1) = -1
func S256(x *big.Int) *big.Int {
	if x.Cmp(tt255) < 0 {
		return x
	}
	return new(big.Int).Sub(x, tt256)
}

// Exp 通过平方实现求幂。
// Exp返回一个新分配的大整数，并且不更改基数或指数。结果被截断为256位。
func Exp(base, exponent *big.Int) *big.Int {
	result := big.NewInt(1)

	for _, word := range exponent.Bits() {
		for i := 0; i < wordBits; i++ {
			if word&1 == 1 {
				U256(result.Mul(result, base))
			}
			U256(base.Mul(base, base))
			word >>= 1
		}
	}
	return result
}

// CalcMemSize 计算某一步需要的内存大小
func CalcMemSize(off, l *big.Int) *big.Int {
	if l.Sign() == 0 {
		return Big0
	}

	return new(big.Int).Add(off, l)
}

// GetDataBig 根据开始位置和大小从数据中返回一个切片,并使用零填充大小.此功能是溢出安全的
func GetDataBig(data []byte, start *big.Int, size *big.Int) []byte {
	dlen := big.NewInt(int64(len(data)))

	s := BigMin(start, dlen)
	e := BigMin(new(big.Int).Add(s, size), dlen)
	return conv.RightPadBytes(data[s.Uint64():e.Uint64()], int(size.Uint64()))
}

// BigUint64 返回已转换为uint64的整数，并返回它是否在进程中溢出
func BigUint64(v *big.Int) (uint64, bool) {
	return v.Uint64(), v.BitLen() > 64
}

// ToWordSize 返回内存扩展所需的word大小.一个word为32个字节
func ToWordSize(size uint64) uint64 {
	if size > MaxUint64-31 {
		return MaxUint64/32 + 1
	}

	return (size + 31) / 32
}

// RandInt 随机数
func RandInt(peerNum int) *big.Int {
	newInt := big.NewInt(int64(peerNum))
	randNum, _ := rand.Int(rand.Reader, newInt)
	return randNum
}
