package utils

import (
	"context"
	"time"
)

// IsTimeEmpty 判断时间是否为空
func IsTimeEmpty(t time.Time) bool {
	if !t.IsZero() {
		return false
	}
	return true
}

// FormatDateString 转换为自定义格式
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

// FormatDateTimeString 获取时间，格式yyyy-MM-dd HH:mm:ss
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

// TimeToString 时间转字符串，格式yyyy-MM-dd HH:mm:ss
func TimeToString(t time.Time) string {
	if IsTimeEmpty(t) {
		t = time.Now()
	}
	return t.Format("2006-01-02 15:04:05")
}

// StringToTime 字符串转时间
func StringToTime(str string) time.Time {
	if IsStringEmpty(str) {
		return time.Now()
	}
	local, _ := time.LoadLocation("Local")
	t, _ := time.ParseInLocation("2006-01-02 15:04:05", "2017-06-20 18:16:15", local)
	return t
}

// ParseTime 解析字符串时间为系统格式
func ParseTime(times string) int64 {
	if "" == times {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02 15:04", times, loc)
	return parse.Unix()
}

// ParseDate 解析字符串日期为系统格式
func ParseDate(dates string) int64 {
	if "" == dates {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02", dates, loc)
	return parse.Unix()
}

// DateEqual 判断两个日期是否相等
func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// GetDateFormat 转换为自定义格式
func GetDateFormat(timestamp uint32, format string) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(int64(timestamp), 0)
	return tm.Format(format)
}

// GetDate 获取时间，使用默认格式
func GetDate(timestamp uint32) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(int64(timestamp), 0)
	return tm.Format("2006-01-02")
}

// GetyyyyMMddHHmm 获取时间，格式yyyy-MM-dd HH:mm
func GetyyyyMMddHHmm(timestamp uint32) string {
	if timestamp <= 0 {
		return ""
	}
	tm := time.Unix(int64(timestamp), 0)
	return tm.Format("2006-01-02 15:04")
}

// GetTimeParse 解析字符串时间为系统格式
func GetTimeParse(times string) int64 {
	if "" == times {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02 15:04", times, loc)
	return parse.Unix()
}

// GetDateParse 解析字符串日期为系统格式
func GetDateParse(dates string) int64 {
	if "" == dates {
		return 0
	}
	loc, _ := time.LoadLocation("Local")
	parse, _ := time.ParseInLocation("2006-01-02", dates, loc)
	return parse.Unix()
}

// SchedulerIntervalsTimer 启动的时候执行一次，不固定某个时间，滚动间隔时间执行
func SchedulerIntervalsTimer(f func(), duration time.Duration) {
	// 定时任务
	ticker := time.NewTicker(duration)
	for {
		go f()
		<-ticker.C
	}
}

// SchedulerIntervalsTimerContext
// 创建一个可以被取消的 context
// ctx, cancel := context.WithCancel(context.Background())
func SchedulerIntervalsTimerContext(ctx context.Context, f func(), duration time.Duration) {
	ticker := time.NewTicker(duration)
	// 在函数退出时，一定要调用 Stop() 来释放资源
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// 等待一个 tick 到达后，再执行任务
			// 这样就避免了立即执行，并且保证了任务不会堆积
			f()
		case <-ctx.Done():
			// 如果外部的 context 发出了取消信号，则退出循环
			return
		}
	}
}

// SchedulerFixedTicker 启动的时候执行一次，固定在每天的某个时间滚动执行
// 首次执行：函数 f 会在 SchedulerFixedTicker 被调用的一瞬间就执行一次。
// 第二次执行：会在下一个午夜0点左右执行。
// 后续执行：从第二次执行开始，每次执行的间隔由传入的 duration 参数决定。
// 如果 duration = 24 * time.Hour：那么它确实会近似于“每天执行一次”。但由于 timer.Reset 存在微小的漂移，长时间运行后，执行时间可能会偏离午夜0点。
// 如果 duration = 1 * time.Hour：那么在第二次执行（午夜0点）之后，它会变成每小时执行一次。
// 如果 duration 是其他值：它就会按该值的间隔执行。
func SchedulerFixedTicker(f func(), duration time.Duration) {
	now := time.Now()
	// 计算下一个时间点
	next := now.Add(duration)
	// 设置目标时间为今天的指定时分秒
	next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
	// 如果目标时间已经过去，则设置为明天的这个时间
	if next.Sub(now) <= 0 {
		next = next.Add(time.Hour * 24)
	}
	// 计算第一次需要等待的时间
	timer := time.NewTimer(next.Sub(now))
	for {
		go f()
		// 等待定时器触发
		<-timer.C
		// Reset 使 ticker 重新开始计时，否则会导致通道堵塞，（本方法返回后再）等待时间段 d 过去后到期。
		// 如果调用时t还在等待中会返回真；如果 t已经到期或者被停止了会返回假
		timer.Reset(duration)
	}
}
