package date

import (
	"fmt"
	"strconv"
	"time"

	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/carmo-evan/strtotime"
)

//获取某一天的0点时间
func GetZeroTime(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

/**
在时间基础上增减时间
ParseDuration解析一个时间段字符串。一个时间段字符串是一个序列，每个片段包含可选的正负号、十进制数、
可选的小数部分和单位后缀，如"300ms"、"-1.5h"、"2h45m"。合法的单位有"ns"、"us" /"µs"、"ms"、"s"、"m"、"h"。
*/
func CTime(t time.Time, time_str string) time.Time {
	time_part, err := time.ParseDuration(time_str)
	if err != nil {
		return t
	}
	return t.Add(time_part)
}

//组装时间区间
// secondInterval        区间,
// orgTime    			 原始时kl
// symbol 				 增减或者减少符号,
// intervalUnit    		 增加或者减少单位,
func AssemblyDateTime(secondInterval int, orgTime time.Time, symbol string, intervalUnit string) string {
	//当前时间加上时间区间的时间
	plusTimeStr := fmt.Sprintf("%s"+strconv.Itoa(secondInterval)+"%s", symbol, intervalUnit)
	currentTimeAddInterval := CTime(orgTime, plusTimeStr)
	//转换时间区间的字符串时间 =  当前时间+ 时间区间
	dateTime := repairTimeZero(strconv.Itoa(currentTimeAddInterval.Hour())) + ":" + repairTimeZero(strconv.Itoa(currentTimeAddInterval.Minute()))
	return dateTime
}

func repairTimeZero(str string) string {
	if strs.Len(str) >= 2 {
		return str
	}
	return "0" + str
}

func StrToTime(str string) (time.Time, error) {
	result, err := time.Parse(consts.AppTimeFormat, str)
	if err == nil {
		return result, nil
	}
	result, err = time.Parse(consts.AppDateFormat, str)
	if err == nil {
		return result, nil
	}
	result, err = time.Parse(consts.AppTimeFormatYYYYMMDDHHmm, str)
	if err == nil {
		return result, nil
	}
	result, err = time.Parse(consts.AppTimeFormatYYYYMMDDHHmmTimezone, str)
	if err == nil {
		return result, nil
	}
	result, err = time.Parse(consts.AppSystemTimeFormat, str)
	if err == nil {
		return result, nil
	}
	result, err = time.Parse(consts.AppSystemTimeFormat8, str)
	if err == nil {
		return result, nil
	}

	return time.Time{}, err
}

func StrToTimeWithLoc(str string, loc *time.Location) (time.Time, error) {
	if str == "" {
		return consts.BlankTimeObject, nil
	}
	result, err := time.ParseInLocation(consts.AppTimeFormat, str, loc)
	if err == nil {
		return result, nil
	}
	result, err = time.ParseInLocation(consts.AppDateFormat, str, loc)
	if err == nil {
		return result, nil
	}
	result, err = time.ParseInLocation(consts.AppTimeFormatYYYYMMDDHHmm, str, loc)
	if err == nil {
		return result, nil
	}
	result, err = time.ParseInLocation(consts.AppTimeFormatYYYYMMDDHHmmTimezone, str, loc)
	if err == nil {
		return result, nil
	}
	result, err = time.ParseInLocation(consts.AppSystemTimeFormat, str, loc)
	if err == nil {
		return result, nil
	}
	result, err = time.ParseInLocation(consts.AppSystemTimeFormat8, str, loc)
	if err == nil {
		return result, nil
	}

	return time.Time{}, err
}

// SmartStrToTime 使用更强大的 strtotime 库来解析时间
// ref: github.com/carmo-evan/strtotime
func SmartStrToTime(str string) (time.Time, error) {
	defaultTime := time.Time{}
	if ts, err := strtotime.Parse(str, time.Now().Unix()); err != nil {
		return defaultTime, err
	} else {
		return time.Unix(ts, 0).UTC(), nil
	}
}
