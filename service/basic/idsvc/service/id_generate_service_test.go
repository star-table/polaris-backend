package service

import (
	"context"
	"fmt"
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo/idvo"
	"github.com/star-table/polaris-backend/service/basic/idsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
	"time"
)

func TestApplyId(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		for i := 0; i < 1000; i++ {
			go Apply(t)
		}
		time.Sleep(time.Duration(10) * time.Second)
		t.Log("end")
	}))
}

func Apply(t *testing.T) {
	//id, err := ApplyId(0, "test1", "")
	//t.Logf("id is %d", id)
	//if err != nil {
	//	t.Log(err)
	//}
}

func applyId(a int) {

	for i := 0; i < 10; i++ {
		//id, err := ApplyMultipleId(100, "C6", "TT", 1)
		id, err := ApplyMultipleId(0, "C6", "", 1)
		if err != nil {
			fmt.Println("error")
			fmt.Println(err)
		}
		fmt.Println(a, id.Ids)
	}
}

func TestApply2(t *testing.T) {

	//config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\configs", "application")
	//id, err := ApplyId2(100, "test1", "CC-T11-", 1)
	//fmt.Println(id)
	//if err != nil {
	//	t.Log(err)
	//}
	//
	//id2, err := ApplyId2(100, "test1", "CC-T11-", 9)
	//fmt.Println(id2)
	//if err != nil {
	//	t.Log(err)
	//}

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(config.GetEnv())
		fmt.Println(config.GetMysqlConfig().Database)

		for i := 0; i < 1; i++ {
			for i := 0; i < 1; i++ {
				go applyId(i)
			}

			select {}
		}
	}))

}

func TestApply3(t *testing.T) {

	//config.LoadConfig("F:\\workspace-golang-polaris\\polaris-backend\\polaris-server\\configs", "application")
	//id, err := ApplyId2(100, "test1", "CC-T11-", 1)
	//fmt.Println(id)
	//if err != nil {
	//	t.Log(err)
	//}
	//
	//id2, err := ApplyId2(100, "test1", "CC-T11-", 9)
	//fmt.Println(id2)
	//if err != nil {
	//	t.Log(err)
	//}

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		applyId(1)

	}))

	//select {}

}

func TestApplyMultiplePrimaryIdByCodes(t *testing.T) {
	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		req := idvo.ApplyMultiplePrimaryIdByCodesReqVo{
			CodeInfos: []idvo.ApplyMultiplePrimaryIdReqVo{
				idvo.ApplyMultiplePrimaryIdReqVo{Count: 2, Code: "aaa"},
				idvo.ApplyMultiplePrimaryIdReqVo{Count: 2, Code: "bbb"},
				idvo.ApplyMultiplePrimaryIdReqVo{Count: 2, Code: "ccc"},
			},
		}

		ids, err := ApplyMultiplePrimaryIdByCodes(req)
		fmt.Println(json.ToJson(ids))

		convey.ShouldBeNil(err)
		convey.ShouldEqual(len(ids.IdCodes), 3)
		convey.ShouldEqual(ids.IdCodes[0], 2)
		convey.ShouldEqual(ids.IdCodes[1], 2)
		convey.ShouldEqual(ids.IdCodes[2], 2)

	}))
}
