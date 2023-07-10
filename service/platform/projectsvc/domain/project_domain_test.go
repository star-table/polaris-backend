package domain

import (
	"context"
	"testing"

	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

//func TestGetProjectList(t *testing.T) {
//	convey.Convey("UnreadNoticeCount", t, test.StartUp(func(ctx context.Context) {
//		order := "create_time  desc"
//		//order1 := "if(1=1,sleep(5),1)"
//		order2 := "id asc"
//		resp, _, err := GetProjectList(3306, db.Cond{
//			consts.TcOrgId:    1242,
//			consts.TcId:       8543,
//			consts.TcIsDelete: 2,
//		}, nil, []*string{&order, nil, &order2}, 10, 1, nil)
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		t.Log(json.ToJsonIgnoreError(resp))
//	}))
//
//}

//func TestIssueCondStatusListAssembly(t *testing.T) {
//	convey.Convey("UnreadNoticeCount", t, test.StartUp(func(ctx context.Context) {
//		t.Log(IssueCondStatusListAssembly(1004, []int{1, 5}))
//	}))
//}

//func TestGetProjectInfo(t *testing.T) {
//	convey.Convey("GetProjectSimple", t, test.StartUp(func(ctx context.Context) {
//		projectBo, err := GetProjectSimple(12769, 1629)
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		t.Log(json.ToJsonIgnoreError(projectBo))
//	}))
//}

func TestGetAllMyProjectIdsWithDeptIds(t *testing.T) {
	convey.Convey("GetProjectSimple", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetAllMyProjectIdsWithDeptIds(2373, 29528, []int64{}, false))
	}))
}
