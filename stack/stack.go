// 该文件是EVM中关于栈的文件.主要包括栈的定义及其操作

package stack

import (
	"fmt"
	"math/big"
)

// Stack 其中data 是一个数组,数组中每个元素是big.Int类型
type Stack struct {
	data []*big.Int
}

// Newstack 获取新的Stack对象
func Newstack() *Stack {
	return &Stack{data: make([]*big.Int, 0, 1024)} // 指定深度1024
}

// Data 返回栈的数据
func (st *Stack) Data() []*big.Int {
	return st.data
}

// Push 推送一个big int对象
func (st *Stack) Push(d *big.Int) {
	st.data = append(st.data, d)
}

// PushN 推送N个big int对象
func (st *Stack) PushN(ds ...*big.Int) {
	st.data = append(st.data, ds...)
}

// Pop 获取 big int对象
func (st *Stack) Pop() (ret *big.Int) {
	ret = st.data[len(st.data)-1]
	st.data = st.data[:len(st.data)-1]
	return
}

// Len 返回big int长度
func (st *Stack) Len() int {
	return len(st.data)
}

// Swap 交换第len-n个元素和最后一个元素
func (st *Stack) Swap(n int) {
	st.data[st.Len()-n], st.data[st.Len()-1] = st.data[st.Len()-1], st.data[st.Len()-n]
}

// Dup 复制
func (st *Stack) Dup(n int) {
	//st.push(pool.get().Set(st.data[st.len()-n]))
	st.Push(new(big.Int).Set(st.data[st.Len()-n]))
}

// Peek 获取栈的最后一个元素
func (st *Stack) Peek() *big.Int {
	return st.data[st.Len()-1]
}

// Back 返回栈的倒数第n个元素
func (st *Stack) Back(n int) *big.Int {
	return st.data[st.Len()-n-1]
}

// Require 检验栈的深度
func (st *Stack) Require(n int) error {
	if st.Len() < n {
		return fmt.Errorf("stack underflow (%d <=> %d)", len(st.data), n)
	}
	return nil
}

// Print 打印栈里面的内容
func (st *Stack) Print() {
	fmt.Println("### stack ###")
	if len(st.data) > 0 {
		for i, val := range st.data {
			fmt.Printf("%-3d  %v\n", i, val)
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("#############")
}
