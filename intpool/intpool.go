// 可以重复利用的big int池

package intpool

import (
	"git.phjr.com/util/stack"
	"math/big"
)

const poolLimit = 256

// IntPool 可以重复利用的big int pool，最大大小为256
type IntPool struct {
	pool *stack.Stack
}

// NewIntPool 返回一个新的big int pool对象
func NewIntPool() *IntPool {
	return &IntPool{pool: stack.Newstack()}
}

// Get 获取一个big int对象
func (p *IntPool) Get() *big.Int {
	if p.pool.Len() > 0 {
		return p.pool.Pop()
	}
	return new(big.Int)
}

// Put 往池里推送一个big int对象
func (p *IntPool) Put(is ...*big.Int) {
	if len(p.pool.Data()) > poolLimit {
		return
	}

	for _, i := range is {

		p.pool.Push(i)
	}
}
