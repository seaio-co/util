常用的数学相关的工具类

## 目录
 序号 | 名称 | 描述 
---|---|---
 1 | big.go | big的一些常用方法
 &nbsp; | `BigPow` | 返回一个指向big.Int类型的指针地址的指针
 &nbsp; | `BigMax` | 返回较大的一个指针
 &nbsp; | `BigMin` |返回较小的一个指针
 &nbsp; | `PaddedBigBytes` | 将一个大整数编码为一个大端字节切片
 &nbsp; | `Byte` | 返回位置n的字节
 &nbsp; | `U256` | U256编码为256位2的补码
 &nbsp; | `S256` | S256将x解释为一个2的补码
 2| integer.go | integer的一些常用方法
 &nbsp; | `SafeSub` | 返回减法结果和是否发生溢出
 &nbsp; | `SafeAdd` | 返回加法结果和是否发生溢出
 &nbsp; | `SafeMul` | 返回乘法结果和是否发生溢出
 