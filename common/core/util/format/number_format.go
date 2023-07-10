package format

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"time"
)

// 验证，小数点后只允许一位
func VerifyFloat1(input string) bool {
	reg := regexp.MustCompile(NumFloat1)
	return reg.MatchString(input)
}

// 前端输入的小时数转换为秒数
func FormatNeedTimeIntoSecondNumber(hourStr string) int64 {
	f, _ := strconv.ParseFloat(hourStr, 64)
	seconds := f * 3600
	return int64(seconds)
}

// 数据库中的数值转为浮点数的小时数，再转为字符串，传给前端。
// 当浮点数为 `0.0` 时，强制为 `0`
func FormatNeedTimeIntoString(secondNum int64) string {
	floatHour := float64(secondNum) / 3600
	if floatHour == 0 || floatHour == 0.0{
		return "0"
	} else {
		return FormatFloatWithoutZero(floatHour, 1)
	}
}

// 保留n位小数，并去除小数点后无效的 0。
// 主要逻辑就是先乘，trunc之后再除回去，就达到了保留N位小数的效果
// 来自于 https://www.jianshu.com/p/ddec820bc4a4
func FormatFloatWithoutZero(num float64, decimal int) string {
	// 默认乘1
	d := float64(1)
	if decimal > 0 {
		// 10的N次方
		d = math.Pow10(decimal)
	}
	// math.trunc作用就是返回浮点数的整数部分
	// 再除回去，小数点后无效的0也就不存在了
	return strconv.FormatFloat(math.Trunc(num*d)/d, 'f', -1, 64)
}

// 检查时间段间的时间，是否合法。不能超过。
func CheckTimeRangeTimeNumIsValid(secondsNum int64, startTime, endTime int64) bool {
	t1 := time.Unix(startTime, 0)
	t2 := time.Unix(endTime, 0)
	t1Start := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2End := time.Date(t2.Year(), t2.Month(), t2.Day(), 23, 59, 59, 0, time.Local)
	// 起止时间转为 xx-xx-xx 00:00:00 ~ xx-xx-xx 23:59:59
	// 但是要计算间隔，间隔值是包含当天最后一秒。所以需 +1。
	diff := (t2End.Unix() + 1) - t1Start.Unix()
	return secondsNum <= diff
}

// 检查起止日期内，天数是否超过限制
func CheckTimeRangeLimitDayNumValid(startTime, endTime, limitDay int64) bool {
	limitSeconds := limitDay * 24 * 3600
	// 起止时间转为 xx-xx-xx 00:00:00 ~ xx-xx-xx 23:59:59
	t1 := time.Unix(startTime, 0)
	t2 := time.Unix(endTime, 0)
	t1Start := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, time.Local)
	t2End := time.Date(t2.Year(), t2.Month(), t2.Day(), 23, 59, 59, 0, time.Local)
	diff := (t2End.Unix() + 1) - t1Start.Unix()
	if diff > limitSeconds {
		return false
	}
	return true
}

// 浮点数保留 n 位有效小数。
func Round(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}


// 将时间戳转换为年月日字符串
func TransformTimeStampToDate(timeStamp int64) string  {
	format := time.Unix(timeStamp, 0).Format("2006-01-02")
	return format
}
