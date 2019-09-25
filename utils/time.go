package utils

import (
	"time"
)

// 判断时间是否为空
func IsTimeEmpty(t time.Time) bool {
	if !t.IsZero() {
		return false
	}
	return true
}

// 转换为自定义格式
func FormatDateString(timestamp uint32, format string) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(int64(timestamp), 0)
	if IsStringEmpty(format) {
		return tm.Format("2006-01-02")
	}
	return tm.Format(format)
}

// 获取时间，格式yyyy-MM-dd HH:mm:ss
func FormatDateTimeString(timestamp uint32, format string) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(int64(timestamp), 0)
	if IsStringEmpty(format) {
		return tm.Format("2006-01-02 15:04:00")
	}
	return tm.Format(format)
}

// 时间转字符串，格式yyyy-MM-dd HH:mm:ss
func TimeToString(t time.Time) string {
	if IsTimeEmpty(t) {
		t = time.Now()
	}
	return t.Format("2006-01-02 15:04:05")
}

// 字符串转时间
func StringToTime(str string) time.Time {
	if IsStringEmpty(str) {
		return time.Now()
	}
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2017-06-20 18:16:15", local)
	return t
}

// 解析字符串时间为系统格式
func ParseTime(times string) int64 {
	if "" == times {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02 15:04", times, loc)
	return parse.Unix()
}

// 解析字符串日期为系统格式
func ParseDate(dates string) int64 {
	if "" == dates {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02", dates, loc)
	return parse.Unix()
}

// 判断两个日期是否相等
func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
