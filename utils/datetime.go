package utils

import (
	"time"
)

/**
 * 转换为自定义格式
 *
 * @param timestamp uint32
 * @param format string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:38
 */
func GetDateFormat(timestamp uint32, format string) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(int64(timestamp), 0)
	return tm.Format(format)
}

/**
 * 获取时间，使用默认格式
 *
 * @param timestamp uint32
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:38
 */
func GetDate(timestamp uint32) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(int64(timestamp), 0)
	return tm.Format("2006-01-02")
}

/**
 * 获取时间，格式yyyy-MM-dd HH:mm
 *
 * @param timestamp uint32
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:39
 */
func GetyyyyMMddHHmm(timestamp uint32) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(int64(timestamp), 0)
	return tm.Format("2006-01-02 15:04")
}

/**
 * 解析字符串时间为系统格式
 *
 * @param times string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:42
 */
func GetTimeParse(times string) int64 {
	if "" == times {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02 15:04", times, loc)
	return parse.Unix()
}

/**
 * 解析字符串日期为系统格式
 *
 * @param dates string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:42
 */
func GetDateParse(dates string) int64 {
	if "" == dates {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02", dates, loc)
	return parse.Unix()
}

/**
 * 判断时间是否为空
 *
 * @param t time.Time
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:43
 */
func IsTimeEmpty(t time.Time) bool {
	if !t.IsZero() {
		return false
	}
	return true
}

/**
 * 时间转字符串，格式yyyy-MM-dd HH:mm:ss
 *
 * @param t time.Time
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:43
 */
func TimeToString(t time.Time) string {
	if IsTimeEmpty(t) {
		t = time.Now()
	}
	format := t.Format("2006-01-02 15:04:05")
	return format
}

/**
 * 字符串转时间
 *
 * @param str string
 * @return
 * @Description
 * @author claer www.bajins.com
 * @date 2019/7/16 16:45
 */
func StringToTime(str string) time.Time {
	if IsStringEmpty(str) {
		return time.Now()
	}
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2017-06-20 18:16:15", local)
	return t
}
