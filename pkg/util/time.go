package util

import (
	"time"
)

// 1小时的时间戳毫秒
const Hour1 = 3600000

// 2小时的时间戳毫秒
const Hour2 = 7200000

// 1天的时间戳毫秒
const Day1 = 86400000

// 7天的时间戳毫秒
const Day7 = 604800000

// 14天的时间戳毫秒
const Day14 = 1209600000

// 获取传入的时间所在月份的第一天，即某月第一天的0点。如传入time.Now(), 返回当前月份的第一天0点时间。
func GetFirstDateOfMonth(d time.Time) time.Time {
	d = d.AddDate(0, 0, -d.Day()+1)
	return GetZeroTime(d)
}

// 获取传入的时间所在月份的最后一天，即某月最后一天的0点。如传入time.Now(), 返回当前月份的最后一天0点时间。
func GetLastDateOfMonth(d time.Time) time.Time {
	return GetFirstDateOfMonth(d).AddDate(0, 1, -1)
}

// 获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// 获取今天午夜0点的时间戳毫秒
func Today0AMMs() int64 {
	t := time.Now()
	tm1 := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())

	return tm1.Unix() * 1000
}

// 获取当前时间戳毫秒
func NowTimeMs() int64 {
	return time.Now().UnixNano() / 1000000
}

// 获取昨天午夜0点的时间戳毫秒
func Yesterday0AMMs() int64 {
	return Today0AMMs() - Day1
}

// 获取明天午夜0点的时间戳毫秒
func Tomorrow0AMMs() int64 {
	return Today0AMMs() + Day1
}

// 获取1周前午夜0点的时间戳毫秒
// 也就是获取7天前午夜0点的时间戳毫秒
func WeekAgo0AMMs() int64 {
	return Today0AMMs() - Day7
}

// 获取2周前午夜0点的时间戳毫秒
// 也就是获取14天前午夜0点的时间戳毫秒
func TwoWeekAgo0AMMs() int64 {
	return Today0AMMs() - Day14
}

// 时间戳毫秒 转 字符串
// fmt.Printf("%v", util.StampToString(1611072000000, "2006-01-02 15:04:05"))
func StampToString(stamp int64, format string) string {
	return time.Unix(stamp/1000, 0).Format(format)
}

// 当前时间戳毫秒 转 字符串
// fmt.Printf("%v", util.NowToString("2006年01月02日 15:04"))
func NowToString(format string) string {
	return StampToString(time.Now().Unix()*1000, format)
}

// 字符串 转 时间戳毫秒
// fmt.Printf("%v", util.StringToStamp("2006-01-02 15:04:05"))
func StringToStamp(str string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")                   //设置时区
	tt, _ := time.ParseInLocation("2006-01-02 15:04:05", str, loc) // 2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
	return tt.Unix() * 1000
}

// 给开始、结束毫秒时间，按天数算出时间的差值
func GetDifferTimeToDay(startTime, endTime int64) (diffDays int64) {
	if startTime == endTime {
		diffDays = 1
		return
	}

	diffDays = (endTime - startTime) / 86400000

	return
}
