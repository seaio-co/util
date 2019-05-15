// 关于 integer 的数学工具包测试
package math

import (
	"math/big"
	"testing"
)

func Test_FirstBitSet(t *testing.T) {

	temp := new(big.Int).SetBytes([]byte("abc"))
	s := FirstBitSet(temp)
	t.Log(s)
}
