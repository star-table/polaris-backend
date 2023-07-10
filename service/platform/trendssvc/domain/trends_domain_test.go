package domain

import (
	"context"
	"fmt"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/objtype"
	"github.com/star-table/polaris-backend/common/model/bo/operation"
	"github.com/star-table/polaris-backend/service/platform/trendssvc/po"
	"github.com/star-table/polaris-backend/service/platform/trendssvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateTrends(t *testing.T) {

	convey.Convey("CreateTrends", t, test.StartUp(func(ctx context.Context) {

		orgId := int64(1000)
		projectId := int64(1)

		trendsBo := &bo.TrendsBo{
			OrgId:           orgId,
			Module1:         objtype.Sys,
			Module2Id:       projectId,
			Module2:         objtype.Project,
			Module3Id:       int64(11),
			Module3:         objtype.Issue,
			OperCode:        operation.Create,
			OperObjId:       int64(11),
			OperObjType:     objtype.Issue,
			OperObjProperty: "",
			RelationObjId:   int64(0),
			RelationType:    "",
			Ext:             "",
			Creator:         int64(1222),
			//CreateTime:      types.NowTime(),
		}
		bil, err := CreateTrends(trendsBo)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(*bil)
	}))
}

func TestQueryTrends(t *testing.T) {

	convey.Convey("queryTrends", t,  test.StartUp(func(ctx context.Context) {

		objId := int64(1)
		condBo := &bo.TrendsQueryCondBo{
			OrgId:   1000,
			ObjType: &objtype.Project,
			ObjId:   &objId,
		}

		obj, err := QueryTrends(condBo)

		fmt.Println(err)
		fmt.Println(json.ToJson(obj))
	}))
}

func TestBo(t *testing.T) {
	bo := bo.TrendsBo{}
	t.Log(bo.CreateTime)

	po := &po.PpmTreTrends{}
	t.Log(po.CreateTime)
	t.Log(po.CreateTime.IsZero())

	copyer.Copy(bo, po)

	t.Log(po.CreateTime)
	t.Log(po.CreateTime.IsZero())

}
