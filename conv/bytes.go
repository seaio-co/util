package conv

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
)

// ToHex 返回以“0x”为前缀的b的十六进制表示。对于空片，返回值是“0x0”。
func ToHex(b []byte) string {
	hex := Bytes2Hex(b)
	if len(hex) == 0 {
		hex = "0"
	}
	return "0x" + hex
}

// ToHexArray 创建一个基于[]字节的十六进制字符串数组
func ToHexArray(b [][]byte) []string {
	r := make([]string, len(b))
	for i := range b {
		r[i] = ToHex(b[i])
	}
	return r
}

// FromHex 返回由十六进制字符串s. s表示的字节，s可以以“0x”为前缀。
func FromHex(s string) []byte {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	h, _ := Hex2Bytes(s)
	return h
}

// hasHexPrefix 验证str以“0x”或“0X”开头。
func HasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

// 验证是否为有效的十六进制字符。
func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isHex验证每个字节是否是有效的十六进制字符串。
func IsHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}

// Bytes2Hex 返回d的十六进制编码。
func Bytes2Hex(d []byte) string {
	return hex.EncodeToString(d)
}

// Hex2Bytes 返回十六进制字符串str所代表的字节。
func Hex2Bytes(str string) ([]byte, error) {
	h, err := hex.DecodeString(str)
	return h, err
}

// Hex2BytesFixed 返回指定长度的字节。
func Hex2BytesFixed(str string, flen int) []byte {
	h, _ := hex.DecodeString(str)
	if len(h) == flen {
		return h
	}
	if len(h) > flen {
		return h[len(h)-flen:]
	}
	hh := make([]byte, flen)
	copy(hh[flen-len(h):flen], h)
	return hh
}

// CopyBytes 返回所提供字节的精确副本。
func CopyBytes(b []byte) (copiedBytes []byte) {
	if b == nil {
		return nil
	}
	copiedBytes = make([]byte, len(b))
	copy(copiedBytes, b)

	return
}

// Int2bytes int 转 bytes
func Int2bytes(i int) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(i))
	return b[:]
}

// Bytes2int bytes转int
func Bytes2int(b []byte) int {
	return int(binary.BigEndian.Uint64(b))
}

//Int82bytes int8类型转为 Bytes
func Int82bytes(i int8) []byte {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], uint16(i))
	return b[:]
}

//Bytes2int8 bytes转int8类型
func Bytes2int8(b []byte) int8 {
	return int8(binary.BigEndian.Uint16(b))
}

//int162bytes int16转bytes
func int162bytes(i int16) []byte {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], uint16(i))
	return b[:]
}

//Bytes2int16 bytes转int16
func Bytes2int16(b []byte) int16 {
	return int16(binary.BigEndian.Uint16(b))
}

//Int322bytes int32转bytes
func Int322bytes(i int32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(i))
	return b[:]
}

//Bytes2int32 bytes转int32
func Bytes2int32(b []byte) int32 {
	return int32(binary.BigEndian.Uint32(b))
}

// Int642bytes int64转Bytes
func Int642bytes(i int64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(i))
	return b[:]
}

// Bytes2int64 bytes转int64
func Bytes2int64(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

// Uint2bytes uint转bytes
func Uint2bytes(u uint) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], uint64(u))
	return b[:]
}

// Bytes2uint bytes转uint
func Bytes2uint(b []byte) uint {
	return uint(binary.BigEndian.Uint64(b))
}

// Uint82bytes uint8转bytes
func Uint82bytes(u uint8) []byte {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], uint16(u))
	return b[:]
}

// Bytes2uint8 bytes转uint8
func Bytes2uint8(b []byte) uint8 {
	return uint8(binary.BigEndian.Uint16(b))
}

// Uint162bytes uint16转bytes
func Uint162bytes(u uint16) []byte {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], u)
	return b[:]
}

// Bytes2uint16 bytes转uint16
func Bytes2uint16(b []byte) uint16 {
	return binary.BigEndian.Uint16(b)
}

// Uint322bytes uint32转bytes
func Uint322bytes(u uint32) []byte {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], u)
	return b[:]
}

// Bytes2uint32 bytes转uint32
func Bytes2uint32(b []byte) uint32 {
	return binary.BigEndian.Uint32(b)
}

// Uint642bytes uint64转bytes
func Uint642bytes(u uint64) []byte {
	var b [8]byte
	binary.BigEndian.PutUint64(b[:], u)
	return b[:]
}

// Bytes2uint64 bytes转uint64
func Bytes2uint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Bool2bytes bool转bytes
func Bool2bytes(b bool) []byte {
	if b {
		return []byte{1}
	}
	return []byte{0}
}

// Bytes2bool ytes转bool
func Bytes2bool(b []byte) bool {
	if len(b) > 0 && b[0] == 1 {
		return true
	}
	return false
}

// Error2bytes error转bytes
func Error2bytes(e error) []byte {
	if e == nil {
		return nil
	}
	return []byte(e.Error())
}

// Bytes2error bytes转error
func Bytes2error(b []byte) error {
	if len(b) == 0 {
		return nil
	}
	return errors.New(string(b))
}

// Rune2bytes rune转bytes
func Rune2bytes(r rune) []byte {
	return []byte(string([]rune{r}))
}

// Bytes2rune bytes 转rune
func Bytes2rune(b []byte) rune {
	return []rune(string(b))[0]
}

// Float642bytes float64转bytes
func Float642bytes(f float64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], math.Float64bits(f))
	return buf[:]
}

// Bytes2float64 bytes转float
func Bytes2float64(b []byte) float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(b))
}

// Float322bytes float32转bytes
func Float322bytes(f float32) []byte {
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[:], math.Float32bits(f))
	return buf[:]
}

// Bytes2float32 bytes  转 float32
func Bytes2float32(b []byte) float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(b))
}

// GetData 根据开始和大小从数据中返回一个切片，并以零填充。 此功能是溢出安全的
func GetData(data []byte, start uint64, size uint64) []byte {
	length := uint64(len(data))
	if start > length {
		start = length
	}
	end := start + size
	if end > length {
		end = length
	}
	return RightPadBytes(data[start:end], int(size))
}

// RightPadBytes 右padbytes 0 -pad片向右直到长度l。如如果L长度小于切片长度则直接输出，如果L长度大于数组长度则右边补0至L长度
func RightPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded, slice)

	return padded
}

// LeftPadBytes left - padbytes向左切片至长度l。如果L长度小于切片长度则直接输出，如果L长度大于数组长度则左边补0至L长度
func LeftPadBytes(slice []byte, l int) []byte {
	if l <= len(slice) {
		return slice
	}

	padded := make([]byte, l)
	copy(padded[l-len(slice):], slice)

	return padded
}
