package api

import (
	"context"
	"fmt"
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/facade/commonfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIndustryList(t *testing.T) {
	convey.Convey("Test 行业列表", t, test.StartUp(func(ctx context.Context) {

		mysqlJson, _ := json.ToJson(config.GetMysqlConfig())
		redisJson, _ := json.ToJson(config.GetRedisConfig())

		fmt.Println("数据库配置:", mysqlJson)
		fmt.Println("redis配置", redisJson)

		orgfacade.GetCurrentUser(ctx)

		resp := commonfacade.IndustryList()

		convey.ShouldBeFalse(resp.Failure())
	}))
}

func TestAA(t *testing.T)  {
	a := "a啊"
	t.Log(len(a))
	t.Log(len(a[:5]))
	t.Log(a[:2])
}
