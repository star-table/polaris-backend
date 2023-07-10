package issue_view

import (
	"context"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
	"strconv"
	"testing"
	"time"
)

func TestCreateTaskView(t *testing.T) {
	convey.Convey("TestGetIssueWorkHoursList", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1213)
		curUserId := int64(1624)
		// 一些值
		projectId := int64(7966)
		isPrivate := true
		remark := "test remark content"
		viewName := "view-" + strconv.FormatInt(time.Now().Unix(), 10)
		resp, err := CreateTaskView(orgId, curUserId, &vo.CreateIssueViewReq{
			ProjectID: &projectId,
			Config:    "{}",
			IsPrivate: &isPrivate,
			Remark:    &remark,
			ViewName:  viewName,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestGetTaskViewList(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1242)
		userId := int64(3307)
		projectObjectTypeID := int64(11395)
		projectID := int64(8411)
		size := 2000
		resp, err := GetTaskViewList(orgId, userId, &vo.GetIssueViewListReq{
			ProjectID: &projectID,
			ProjectObjectTypeID:  &projectObjectTypeID,
			Size: &size,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}
func TestDeleteTaskView(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1213)
		userId := int64(1624)
		viewId := int64(1007)
		resp, err := DeleteTaskView(orgId, userId, viewId)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestGetOptionList(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1213)
		projectId := int64(8161)
		projectObjectTypeId := int64(0)
		resp, err := GetOptionList(orgId, projectvo.GetOptionListReqInput{
			ProjectId:           projectId,
			ProjectObjectTypeId: projectObjectTypeId,
			FieldName:           -4,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}
