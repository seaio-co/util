package time

import (
	"fmt"
	"testing"
)

//测试当前时间戳(秒级)
func Test_CurrentSecond(t *testing.T) {
	num := CurrentSecond()
	t.Log(num)
}

//测试当前时间戳(毫秒级)
func Test_CurrentMilliSecond(t *testing.T) {
	num := CurrentMilliSecond()
	t.Log(num)
}

//测试获取当前月月初时间(秒级）
func Test_GetEarlyMonthUnix(t *testing.T) {
	s := GetEarlyMonthUnix()
	t.Log(s)
}

//测试获取年初时间(秒级）
func Test_GetEarlyYearUnix(t *testing.T) {
	s := GetEarlyYearUnix()
	t.Log(s)
}

//测试获取当日零时时间(秒级）
func Test_GetZeroHourUnix(t *testing.T) {
	t.Log(GetZeroHourUnix())
}

//测试按照给定格式输出给定的时间戳
func Test_GetUnixToFormatString(t *testing.T) {
	s := GetUnixToFormatString(1550116063, "2006-01-02 15:04:05")
	t.Log(s)
}

///测试获取给定时间戳的年月日
func Test_GetUnixToString(t *testing.T) {
	fmt.Println(GetUnixToString(1550116063))
}

//测试获取指定格式的给定时间戳的小时和分钟数
func Test_GetUnixToHourString(t *testing.T) {
	fmt.Println(GetUnixToHourString(1550116063))
}

//测试获取指定时间戳对应的月份
func Test_GetUnixToMonth(t *testing.T) {
	fmt.Println(GetUnixToMonth(1550116063))
}

//测试获取指定时间戳对应的月份.日期
func Test_GetUnixToDayTime(t *testing.T) {
	fmt.Println(GetUnixToDayTime(1550116063))
}

//测试将时间字符串转换为时间戳
func Test_GetStringToFormatUnix(t *testing.T) {
	t1 := "2019-01-08 13:50:30" //外部传入的时间字符串
	s := GetStringToFormatUnix(t1)
	t.Log(s)
}
