关于类型转换的工具类

## 目录
 序号 | 名称 | 描述 
---|---|---
 1 | bytes.go | 关于bytes转换的常用方法
 &nbsp; | `ToHex` | ToHex返回以“0x”为前缀的b的十六进制表示。对于空片，返回值是“0x0”
 &nbsp; | `ToHexArray` | ToHexArray创建一个基于[]字节的十六进制字符串数组
 &nbsp; | `FromHex` | FromHex返回由十六进制字符串
 &nbsp; | `Bytes2Hex` | Bytes2Hex返回十六进制编码
 &nbsp; | `Hex2Bytes` | Hex2Bytes返回十六进制字符串所代表的字节
 &nbsp; | `Hex2BytesFixed` | Hex2BytesFixed返回指定长度的字节
 &nbsp; | `CopyBytes` | CopyBytes返回所提供字节的精确副本
 &nbsp; | `Int2bytes` | int类型 转 bytes类型
 &nbsp; | `Bytes2int` | bytes转int类型
 &nbsp; | `Int82bytes` | int8类型转为 Bytes
 &nbsp; | `Bytes2int8` | bytes转int8类型
 &nbsp; | `int162bytes` | int16转bytes
 &nbsp; | `Bytes2int16` | bytes转int16
 &nbsp; | `Int322bytes` | int32转bytes
 &nbsp; | `Bytes2int32` | bytes转int32
 &nbsp; | `Int642bytes` | int64转Bytes
 &nbsp; | `Bytes2int64` | bytes转int64
 &nbsp; | `Uint2bytes` | uint转bytes
 &nbsp; | `Bytes2uint` | bytes转uint
 &nbsp; | `Uint82bytes` | uint8转bytes
 &nbsp; | `Bytes2uint8` | bytes转uint8
 &nbsp; | `Uint162bytes` | uint16转bytes
 &nbsp; | `Bytes2uint16` | bytes转uint16
 &nbsp; | `Uint322bytes` | uint32转bytes
 &nbsp; | `Bytes2uint32` | bytes转uint32
 &nbsp; | `Uint642bytes` | uint64转bytes
 &nbsp; | `Bytes2uint64` | bytes转uint64
 &nbsp; | `Bool2bytes` | bool转bytes
 &nbsp; | `Bytes2bool` | bytes转bool
 &nbsp; | `Error2bytes` | error转bytes
 &nbsp; | `Bytes2error` | bytes转error
 &nbsp; | `Rune2bytes` | rune转bytes
 &nbsp; | `Bytes2rune` | bytes 转rune
 &nbsp; | `Float642bytes` | float64转bytes 
 &nbsp; | `Bytes2float64` | bytes转float 
 &nbsp; | `Float322bytes` | float32转bytes 
 &nbsp; | `Bytes2float32` | bytes转float32
 &nbsp; | `getData` | 根据开始和大小从数据中返回一个切片，并以零填充。此功能是溢出安全的
 &nbsp; | `RightPadBytes` | 右padbytes 0 -pad片向右直到长度l。如如果L长度小于切片长度则直接输出，如果L长度大于数组长度则右边补0至L长度
 &nbsp; | `LeftPadBytes` | left - padbytes向左切片至长度l。如果L长度小于切片长度则直接输出，如果L长度大于数组长度则左边补0至L长度
 2 | types.go | 关于hash、sign、address等转换的常用方法
  &nbsp; | `Bytes` | Bytes获取基础散列的字节表示
  &nbsp; | `Big` | Big将散列转换为大整数
  &nbsp; | `Hex` | 十六进制将散列转换为十六进制字符串
  &nbsp; | `TerminalString` | TerminalString实现日志。在日志记录期间格式化控制台输出的字符串
  &nbsp; | `String` | String实现了stringer接口，并且在对文件进行完整日志记录时也被日志记录器使用
  &nbsp; | `Format` | 实现了fmt格式。格式化程序，强制按原样格式化字节片，而不需要通过用于日志记录的stringer接口
  &nbsp; | `BytesToHash` | 将byte数组转换为hash
  &nbsp; | `BigToHash` | BigToHash将b的字节表示设置为hash
  &nbsp; | `HexToHash` | 将十六进制的字符串表示为hash
  &nbsp; | `BytesToAddress` | byte转address
  &nbsp; | `StringToAddress` | String转Address返回字节
  &nbsp; | `BigToAddress` | Big转Address返回字节值为b的地址
  &nbsp; | `HexToAddress` | 十六进制字符串转地址
  &nbsp; | `IsHexAddress` | 验证字符串是否可以表示有效的十六进制编码的地址
  &nbsp; | `Bytes` | 字节获取底层地址的字符串表示形式
  &nbsp; | `Big` | Big将地址转换为一个大整数
  &nbsp; | `Hash` | 哈希通过左填充0将地址转换为哈希
  &nbsp; | `Hex` | 十六进制返回地址的十六进制字符串表示形式
  &nbsp; | `BytesToSign` | bytes转sign
  &nbsp; | `HexToSign` | 十六进制字符串转为sign
  
  ## 单元测试 
  
  序号 | 名称 | 说明
  ---|---|---
 1 | bytes_test.go  | 对bytes常用方法进行功能测试
 &nbsp; | `Test_ToHex` | 测试byte数组转换为以0x十六进制字符串
 &nbsp; | `Test_Bytes2Hex` | 测试byte数组转换为十六进制字符串
 &nbsp; | `Test_Hex2BytesFixed` | Hex2BytesFixed返回指定长度的字节
 &nbsp; | `Test_RightPadBytes` | 右padbytes 0 -pad片向右直到长度L
 &nbsp; | `Test_LeftPadBytes` | left - padbytes向左切片至长度L
 2 | types_test.go  | 对hash、sign、address常用方法进行功能测试
 &nbsp; | `Test_BytesToHash` | 测试bytes 转 Hash
 &nbsp; | `Test_HexToHash` | 测试十六进制字符串转hash
 &nbsp; | `Test_BytesToAddress` | 测试 byte转address
 &nbsp; | `Test_StringToAddress` | 测试string转address
 &nbsp; | `Test_BigToAddress` | 测试big.int转地址
 &nbsp; | `Test_IsHexAddress` | 测试是否为有效的十六进制地址.int转地址
 &nbsp; | `Test_BytesToSign` | 测试bytes转 sign
 &nbsp; | `Test_HexToSign` | 测试十六进制转sign
  