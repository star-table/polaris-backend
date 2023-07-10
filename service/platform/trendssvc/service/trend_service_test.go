package service

import (
	"context"
	"fmt"
	"testing"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/trendssvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestTrendList(t *testing.T) {
	//type Exte struct {
	//	ObjName string
	//}
	//a := "{\"issueType\":\"T\",\"objName\":\"测试任务\",\"ccc\":\"aa\"}"
	//ex := map[string]interface{}{
	//	"issueType":"T",
	//	"objName":"测试",
	//}
	////ex := Ext{
	////	ObjName:"ceshi",
	////}
	//b := &Exte{}
	//c := &Exte{}
	//d := &map[string]interface{}{}
	//fmt.Println(json.FromJson(json.ToJsonIgnoreError(ex), b), b)
	//fmt.Println(json.FromJson(a, c), c)
	//fmt.Println(json.FromJson(a, d), d, (*d)["objName"])
	old := "{\"participant\":[1081,1064,1069,1077,1078,1079,1080], \"name\":\"ss\"}"
	new := "{\"participant\":[1080,1069,1079,1077,1078], \"name\":\"ssb\"}"
	old1 := &map[string]interface{}{}
	new1 := &map[string]interface{}{}
	json.FromJson(old, old1)
	json.FromJson(new, new1)
	fmt.Println(slice.JsonCompare(*old1, *new1))
}

func TestTrendList2(t *testing.T) {
	convey.Convey("UnreadNoticeCount", t, test.StartUp(func(ctx context.Context) {
		//operId := int64(1042)
		//objType := "Issue"
		OrderType := int(2)
		page := int64(1)
		size := int64(1)
		chooseType := int(5)
		res, err := TrendList(1417, 2252, &vo.TrendReq{
			LastTrendID: nil,
			ObjType:     nil,
			ObjID:       nil,
			OperID:      nil,
			StartTime:   nil,
			EndTime:     nil,
			Type:        &chooseType,
			Page:        &page,
			Size:        &size,
			OrderType:   &OrderType,
		})
		t.Log(json.ToJsonIgnoreError(res))
		t.Log(err)
	}))

}

func TestGetProjectLatestUpdateTime(t *testing.T) {
	convey.Convey("UnreadNoticeCount", t, test.StartUp(func(ctx context.Context) {
		res, err := GetProjectLatestUpdateTime(1111, []int64{1511, 1374})
		t.Log(err)
		t.Log(json.ToJsonIgnoreError(res))
	}))
}
