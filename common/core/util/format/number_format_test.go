package format

import (
	"fmt"
	"github.com/magiconair/properties/assert"
	"testing"
	"time"
)

func TestVerifyFloat1(t *testing.T) {
	var isOk bool
	isOk = VerifyFloat1("xxxxx.xx")
	assert.Equal(t, isOk, false)
	isOk = VerifyFloat1("1.1")
	assert.Equal(t, isOk, true)
	isOk = VerifyFloat1("1.11")
	assert.Equal(t, isOk, false)
	isOk = VerifyFloat1("211")
	assert.Equal(t, isOk, true)
}

func TestName(t *testing.T) {
	h1 := 24 * 3600
	tt1, _ := time.Parse("2006-01-02 00:00:00", "2020-10-01 00:00:00")
	tt2, _ := time.Parse("2006-01-02 00:00:00", "2020-10-01 00:00:00")
	t1 := tt1.Unix()
	t2 := tt2.Unix()
	// 一天内，时间是 24 小时减 1 秒
	isValid := CheckTimeRangeTimeNumIsValid(int64(h1), t1, t2)
	if isValid != true {
		t.Error("异常 001")
	}
	t.Logf("%v\n", isValid)
	h1 = 24 * 3600 + 1
	tt1, _ = time.Parse("2006-01-02 00:00:00", "2020-10-01 00:00:00")
	tt2, _ = time.Parse("2006-01-02 00:00:00", "2020-10-01 00:00:00")
	t1 = tt1.Unix()
	t2 = tt2.Unix()
	// 一天内，时间是 24 小时减 1 秒
	isValid = CheckTimeRangeTimeNumIsValid(int64(h1), t1, t2)
	if isValid != false {
		t.Error("异常 002")
	}
}

func TestFormatNeedTimeIntoString(t *testing.T) {
	n1 := FormatNeedTimeIntoString(22100)
	fmt.Printf("%v\n", n1)
	n1 = FormatNeedTimeIntoString(7200)
	fmt.Printf("%v\n", n1)
}

func TestCheckTimeRangeLimitDayNumValid(t *testing.T) {
	tt1, _ := time.Parse("2006-01-02 00:00:00", "2020-01-01 00:00:00")
	tt2, _ := time.Parse("2006-01-02 00:00:00", "2020-07-01 00:00:00")
	isValid := CheckTimeRangeLimitDayNumValid(tt1.Unix(), tt2.Unix(), 185)
	if !isValid {
		t.Error(isValid)
	}
	tt1, _ = time.Parse("2006-01-02 00:00:00", "2020-10-01 00:00:00")
	tt2, _ = time.Parse("2006-01-02 00:00:00", "2020-10-02 00:00:00")
	if !CheckTimeRangeLimitDayNumValid(tt1.Unix(), tt2.Unix(), 2) {
		t.Error(isValid)
	}
}

func TestTransformTimeStampToDate(t *testing.T) {
	stampToDate := TransformTimeStampToDate(1660187336)
	t.Log(stampToDate)
}
