时间相关的工具包

## 目录
 序号 | 名称 | 描述 
---|---|---
 1 | time.go | 时间的一些常用方法
 &nbsp; | `CurrentSecond` | 当前时间戳(秒级)
 &nbsp; | `CurrentMilliSecond` | 当前时间戳(毫秒级)
 &nbsp; | `GetEarlyYearUnix` | 获取年初时间(秒级）
 &nbsp; | `GetEarlyMonthUnix` | 获取当前月月初时间(秒级）
 &nbsp; | `GetZeroHourUnix` | 获取当日零时时间(秒级）
 &nbsp; | `GetNowHourUnix` | 获取当前小时时间
 &nbsp; | `GetUnixToFormatString` | 按照给定格式输出给定的时间戳 
 &nbsp; | `GetUnixToString` | 获取给定时间戳的年月日
 &nbsp; | `GetUnixToHourString` | 获取指定格式的给定时间戳的小时和分钟数
 &nbsp; | `GetUnixToMonth` | 获取指定时间戳对应的月份
 &nbsp; | `GetUnixToDay` | 获取指定时间戳对应的月份
 &nbsp; | `GetUnixToDayTime` | 获取指定时间戳对应的月份.日期
 &nbsp; | `GetStringToFormatUnix` | 将时间字符串转换为时间戳
 
 ## 单元测试 
   
   序号 | 名称 | 说明
   ---|---|---
  1 | time_test.go  | 对bytes常用方法进行功能测试
  &nbsp; | `Test_CurrentSecond` | 测试当前时间戳(秒级)
  &nbsp; | `Test_CurrentMilliSecond` | 测试当前时间戳(毫秒级)
  &nbsp; | `Test_GetEarlyMonthUnix` | 测试获取当前月月初时间(秒级）
  &nbsp; | `Test_GetEarlyYearUnix` | 测试获取年初时间(秒级）
  &nbsp; | `Test_GetZeroHourUnix` | 测试获取当日零时时间(秒级）
  &nbsp; | `Test_GetUnixToFormatString` | 测试按照给定格式输出给定的时间戳
  &nbsp; | `Test_GetUnixToString` | 测试获取给定时间戳的年月日
  &nbsp; | `Test_GetUnixToHourString` | 测试获取指定格式的给定时间戳的小时和分钟数
  &nbsp; | `Test_GetUnixToMonth` | 测试获取指定时间戳对应的月份
  &nbsp; | `Test_GetUnixToDayTime` | 测试获取指定时间戳对应的月份.日期
  &nbsp; | `Test_GetStringToFormatUnix` | 测试将时间字符串转换为时间戳
 