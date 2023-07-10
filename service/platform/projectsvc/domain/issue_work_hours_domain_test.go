package domain

import (
	"context"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

// 工时记录列表
func TestGetIssueWorkHoursList(t *testing.T) {
	convey.Convey("TestGetIssueWorkHoursList", t, test.StartUp(func(ctx context.Context) {
		param := bo.IssueWorkHoursBoListCondBo{OrgId: 1004}
		list, err := GetIssueWorkHoursList(param)
		if err != nil {
			t.Errorf("%v\n", err)
			return
		}
		t.Logf("%#v\n", list)
	}))
}

//// 新增工时记录
//func TestCreateIssueHours(t *testing.T) {
//	convey.Convey("TestCreateIssueHours", t, test.StartUp(func(ctx context.Context) {
//		orgId := int64(1004)
//		userId := int64(1023)
//		desc := "this is some work content desc"
//		projectId := 1007
//		issueId := 1030
//		startTime := int64(time.Now().Second())
//		endTime := int64(0)
//		input := vo.CreateIssueWorkHoursReq{
//			ProjectID: int64(projectId),
//			IssueID:   int64(issueId),
//			Type:      1,// 枚举记录类型：1预估工时记录，2实际工时记录
//			WorkerID:  0,
//			NeedTime:  62,// 分钟
//			StartTime: startTime,
//			EndTime:   &endTime,
//			Title:     "test-title",
//			Desc:      &desc,
//		}
//		reqVo := projectvo.CreateIssueWorkHoursReqVo{
//			Input:  input,
//			UserId: userId,
//			OrgId:  orgId,
//		}
//		err := CreateIssueWorkHours(reqVo)
//		if err != nil {
//			t.Errorf("%v\n", err)
//			return
//		}
//		t.Logf("end-ok\n")
//	}))
//}
