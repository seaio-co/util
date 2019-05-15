// Package memory 该文件主要是关于虚拟机内存的一些操作
package memory

import "fmt"

// Memory 内存用于一些内存操作（MLOAD,MSTORE,MSTORE8）及合约调用的参数拷贝（CALL，CALLCODE）
// 内存数据结构，维护了一个byte数组，MLOAD，MSTORE读取存入的时候都要指定位置及长度才能准确的读写
type Memory struct {
	store []byte
}

// NewMemory 创建内存对象
func NewMemory() *Memory {
	return &Memory{}
}

// Set 从指定位置设置内存一定大小的值
// offset	偏移量
// size		大小
// value	内容
func (m *Memory) Set(offset, size uint64, value []byte) {
	// 参数检查
	if size > uint64(len(m.store)) {
		panic("INVALID memory: store empty")
	}

	if size > 0 {
		copy(m.store[offset:offset+size], value)
	}
}

// SetByte 存放字节数据
func (m *Memory) SetByte(offset uint64, value byte) {
	m.store[offset] = value
}

// Resize 重新分配内存大小
func (m *Memory) Resize(size uint64) {
	if uint64(m.Len()) < size {
		m.store = append(m.store, make([]byte, size-uint64(m.Len()))...)
	}
}

// Get 返回内存中一定偏移量/固定大小的内容,返回一个新的切片
func (m *Memory) Get(offset, size int64) (cpy []byte) {
	if size == 0 {
		return nil
	}

	if len(m.store) > int(offset) {
		cpy = make([]byte, size)
		copy(cpy, m.store[offset:offset+size])

		return
	}

	return
}

// GetPtr 返回内存中一定偏移量/固定大小的内容
func (m *Memory) GetPtr(offset, size int64) []byte {
	if size == 0 {
		return nil
	}

	if len(m.store) > int(offset) {
		return m.store[offset : offset+size]
	}

	return nil
}

// Len 返回内存内容的长度
func (m *Memory) Len() int {
	return len(m.store)
}

// Data 返回内存内容
func (m *Memory) Data() []byte {
	return m.store
}

// Print 打印内存中的数据
func (m *Memory) Print() {
	fmt.Printf("### mem %d bytes ###\n", len(m.store))
	if len(m.store) > 0 {
		addr := 0
		for i := 0; i+32 <= len(m.store); i += 32 {
			fmt.Printf("%03d: % x\n", addr, m.store[i:i+32])
			addr++
		}
	} else {
		fmt.Println("-- empty --")
	}
	fmt.Println("####################")
}
