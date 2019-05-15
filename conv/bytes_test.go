package conv

import (
	"testing"
)

//byte数组转换为以0x十六进制字符串
func Test_ToHex(t *testing.T) {
	hex := ToHex([]byte("11"))
	t.Log(hex)
}

//byte数组转换为十六进制字符串
func Test_Bytes2Hex(t *testing.T) {
	h := Bytes2Hex([]byte("11"))
	t.Log(h)
}

// Hex2BytesFixed返回指定长度的字节。
func Test_Hex2BytesFixed(t *testing.T) {
	s := Hex2BytesFixed("hklsla", 6)
	t.Log(s)
}

// 右padbytes 0 -pad片向右直到长度l。
func Test_RightPadBytes(t *testing.T) {
	s := RightPadBytes([]byte("2"), 3)
	t.Log(s)
}

// left - padbytes向左切片至长度l。
func Test_LeftPadBytes(t *testing.T) {
	s := LeftPadBytes([]byte("2"), 3)
	t.Log(s)
}
