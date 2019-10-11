package time

import (
	"strconv"
	"time"
)

// CurrentSecond 当前时间戳(秒级) timestamp
func CurrentSecond() int64 {
	return time.Now().Unix()
}

// CurrentMilliSecond 当前时间戳(毫秒级)
func CurrentMilliSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

// GetEarlyYearUnix 获取年初时间(秒级）
func GetEarlyYearUnix() int64 {
	now := time.Now()
	tm := time.Date(now.Year(), 1, 1, 0, 0, 0, 0, now.Location())
	return tm.Unix()
}

// GetEarlyMonthUnix 获取当前月月初时间(秒级）
func GetEarlyMonthUnix() int64 {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return tm.Unix()
}

// GetZeroHourUnix 获取当日零时时间(秒级）
func GetZeroHourUnix() int64 {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	return tm.Unix()
}

// GetNowHourUnix 获取当前小时时间(秒级）
func GetNowHourUnix() int64 {
	now := time.Now()
	tm := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), 0, 0, 0, now.Location())
	return tm.Unix()
}

// MonthMap 12个月份用两位数值表示
var MonthMap = map[string]string{
	"January":   "01",
	"February":  "02",
	"March":     "03",
	"April":     "04",
	"May":       "05",
	"June":      "06",
	"July":      "07",
	"August":    "08",
	"September": "09",
	"October":   "10",
	"November":  "11",
	"December":  "12",
}

// MonthIntMap 12个月份
var MonthIntMap = map[string]int{
	"January":   1,
	"February":  2,
	"March":     3,
	"April":     4,
	"May":       5,
	"June":      6,
	"July":      7,
	"August":    8,
	"September": 9,
	"October":   10,
	"November":  11,
	"December":  12,
}

// GetUnixToFormatString 按照给定格式输出给定的时间戳 例：1550116063  2006-01-02 15:04 输出 2019-02-14 11:47 (时间转换的模板，golang里面只能是 "2006-01-02 15:04:05" （go的诞生时间）)
func GetUnixToFormatString(timestamp int64, f string) string {
	tm := time.Unix(timestamp, 0)
	return tm.Format(f)
}

// GetUnixToString 获取给定时间戳的年月日 如下格式输出
func GetUnixToString(timestamp int64) string {
	return GetUnixToFormatString(timestamp, "2006-01-02")
}

// GetUnixToHourString 获取指定格式的给定时间戳的小时和分钟数
func GetUnixToHourString(timestamp int64) string {
	return GetUnixToFormatString(timestamp, "15:04")
}

// GetUnixToMonth 获取指定时间戳对应的月份
func GetUnixToMonth(timestamp int64) string {
	tm := time.Unix(timestamp, 0)
	return MonthMap[tm.Month().String()]
}

// GetUnixToDay 获取指定时间戳对应的月份
func GetUnixToDay(timestamp int64) int {
	tm := time.Unix(timestamp, 0)
	return tm.Day()
}

// GetUnixToDayTime 获取指定时间戳对应的月份.日期 例：02.14
func GetUnixToDayTime(timestamp int64) string {
	month := GetUnixToMonth(timestamp)
	day := GetUnixToDay(timestamp)
	d := month + "." + strconv.Itoa(day)
	return d
}

// GetStringToFormatUnix 将时间字符串转换为时间戳
func GetStringToFormatUnix(timestr string) int64 {
	timeTemplate := "2006-01-02 15:04:05"                               //常规类型
	stamp, _ := time.ParseInLocation(timeTemplate, timestr, time.Local) //使用parseInLocation将字符串格式化返回本地时区时间
	return stamp.Unix()
}

// Slight modification of the RFC3339Nano but it right pads all zeros and drops the time zone info
const SortableTimeFormat = "2006-01-02T15:04:05.000000000"

// Formats a time.Time into a []byte that can be sorted
func FormatTimeBytes(t time.Time) []byte {
	return []byte(t.UTC().Round(0).Format(SortableTimeFormat))
}

// Parses a []byte encoded using FormatTimeKey back into a time.Time
func ParseTimeBytes(bz []byte) (time.Time, error) {
	str := string(bz)
	t, err := time.Parse(SortableTimeFormat, str)
	if err != nil {
		return t, err
	}
	return t.UTC().Round(0), nil
}
