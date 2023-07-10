package date

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"testing"
	"time"
)

func TestGetZeroTime(t *testing.T) {
	now := time.Now()
	t.Log(GetZeroTime(now))
	t.Log(GetZeroTime(now).AddDate(0, 0, 1))
}

func TestDateFormat(t *testing.T){
	dateStr := "2019-11-27 17:51 +0800"
	tim, err := time.Parse(consts.AppTimeFormatYYYYMMDDHHmmTimezone, dateStr)
	t.Log(tim)
	t.Log(err)
}