package service

import (
	"context"
	"fmt"
	"testing"

	tablePb "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"

	"github.com/star-table/polaris-backend/common/model/vo"

	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestBatchCreateIssue(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		var data []map[string]interface{}
		documents := make(map[string]*bo.Attachments)
		documents["1008050"] = &bo.Attachments{
			Id:           1008050,
			Name:         "v2-58e1965db5115b7b78481ce5494e65de_720w.jpeg",
			Path:         "https://attachments.startable.cn/org_2571/project_13757/issue_10093031/resource/2022/9/27/fa0f9abe260b476ebfc947ea60480e091664281657351.jpeg",
			ResourceType: 2,
			RecycleFlag:  2,
		}

		data1 := map[string]interface{}{
			consts.BasicFieldTitle:               "titletitle",                     // 标题
			consts.BasicFieldOwnerId:             []string{"U_1368436", "U_29600"}, // 负责人
			consts.BasicFieldIssueStatus:         26,                               // 任务状态
			consts.BasicFieldProjectObjectTypeId: 1,                                // 任务栏
			consts.BasicFieldPlanStartTime:       "2022-09-20 10:00:00",            // 计划开始时间
			consts.BasicFieldPlanEndTime:         "2022-09-20 11:00:00",            // 计划结束时间
			consts.BasicFieldRemark:              "remark",                         // 描述
			consts.BasicFieldRemarkDetail:        "remark detail",                  // 描述详情
			consts.BasicFieldFollowerIds:         []string{"U_1368436", "U_29600"}, // 关注人
			consts.BasicFieldAuditorIds:          []string{"U_1368436", "U_29600"}, // 确认人
			consts.BasicFieldAuditStatus:         1,
			consts.BasicFieldAuditStatusDetail: map[string]int{
				"1368436": 2,
				"29600":   3,
			},
			consts.BasicFieldPriority: 11224, // 优先级
			//"_field_1664272984785":               2,                                // 分组单选
			//consts.BasicFieldDocument:            documents,                        // 附件
			consts.BasicFieldRelating: bo.RelatingIssue{
				LinkTo:   []string{"1573159497357164545"},
				LinkFrom: []string{},
			},
		}
		data2 := map[string]interface{}{ // 子任务
			consts.BasicFieldTitle:       "测试2",
			consts.BasicFieldAuditorIds:  []string{"U_111"},
			consts.BasicFieldIssueStatus: 26,
			consts.TempFieldChildren: []map[string]interface{}{
				data1,
			},
			//consts.BasicFieldParentId: 10093400,
		}
		data = append(data, data2)

		//data = append(data, map[string]interface{}{
		//	consts.BasicFieldTitle:       "测试3",
		//	consts.BasicFieldAuditorIds:  []string{"U_111", "U_29600"},
		//	consts.BasicFieldIssueStatus: 26,
		//})
		result, _, _, err := BatchCreateIssue(&projectvo.BatchCreateIssueReqVo{
			OrgId:     2571,
			UserId:    1368436,
			AppId:     1527557851545960449,
			ProjectId: 13757,
			TableId:   1572900881987276800,
			Data:      data,
		}, true, &projectvo.TriggerBy{
			TriggerBy: consts.TriggerByNormal,
		}, "", 0)
		if err != nil {
			t.Log(err)
		} else {
			t.Log(result)
		}
	}))
}

func TestBatchCreateIssue2(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		var data []map[string]interface{}
		documents := make(map[string]*bo.Attachments)
		documents["1008050"] = &bo.Attachments{
			Id:           1008050,
			Name:         "v2-58e1965db5115b7b78481ce5494e65de_720w.jpeg",
			Path:         "https://attachments.startable.cn/org_2571/project_13757/issue_10093031/resource/2022/9/27/fa0f9abe260b476ebfc947ea60480e091664281657351.jpeg",
			ResourceType: 2,
			RecycleFlag:  2,
		}

		data1 := map[string]interface{}{
			consts.BasicFieldTitle:               "titletitle",          // 标题
			consts.BasicFieldOwnerId:             []string{"U_29612"},   // 负责人
			consts.BasicFieldIssueStatus:         7,                     // 任务状态
			consts.BasicFieldProjectObjectTypeId: 1,                     // 任务栏
			consts.BasicFieldPlanStartTime:       "2022-09-20 10:00:00", // 计划开始时间
			consts.BasicFieldPlanEndTime:         "2022-09-20 11:00:00", // 计划结束时间
			consts.BasicFieldRemark:              "remark",              // 描述
			consts.BasicFieldRemarkDetail:        "remark detail",       // 描述详情
			consts.BasicFieldFollowerIds:         []string{"U_29612"},   // 关注人
			consts.BasicFieldAuditorIds:          []string{"U_29612"},   // 确认人
			//consts.BasicFieldPriority:            11224,                 // 优先级
			consts.BasicFieldDocument: documents, // 附件
			"xxxxxx":                  1,
			//consts.BasicFieldRelating: bo.RelatingIssue{
			//	LinkTo:   []string{"1573159497357164545"},
			//	LinkFrom: []string{},
			//},
		}
		data2 := map[string]interface{}{ // 子任务
			consts.BasicFieldTitle:       "测试2",
			consts.BasicFieldAuditorIds:  []string{"U_111"},
			consts.BasicFieldIssueStatus: 26,
			consts.TempFieldChildren: []map[string]interface{}{
				data1,
			},
		}
		data = append(data, data2)

		//data = append(data, map[string]interface{}{
		//	consts.BasicFieldTitle:       "测试3",
		//	consts.BasicFieldAuditorIds:  []string{"U_111", "U_29600"},
		//	consts.BasicFieldIssueStatus: 26,
		//})
		result, _, _, err := BatchCreateIssue(&projectvo.BatchCreateIssueReqVo{
			OrgId:     2373,
			UserId:    29612,
			AppId:     1537674994417500161,
			ProjectId: 13945,
			TableId:   1537674995105271808,
			Data:      data,
		}, true, &projectvo.TriggerBy{
			TriggerBy: consts.TriggerByNormal,
		}, "", 0)
		if err != nil {
			t.Log(err)
		} else {
			t.Log(result)
		}
	}))
}

func TestBatchCreateIssue3(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		var data []map[string]interface{}
		data2 := map[string]interface{}{ // 子任务
			consts.BasicFieldTitle: "测试2",
		}
		data = append(data, data2)

		//data = append(data, map[string]interface{}{
		//	consts.BasicFieldTitle:       "测试3",
		//	consts.BasicFieldAuditorIds:  []string{"U_111", "U_29600"},
		//	consts.BasicFieldIssueStatus: 26,
		//})
		result, _, _, err := BatchCreateIssue(&projectvo.BatchCreateIssueReqVo{
			OrgId:     2578,
			UserId:    29572,
			AppId:     1582283822806958081,
			ProjectId: 61490,
			TableId:   1582283823167574016,
			Data:      data,
		}, true, &projectvo.TriggerBy{
			TriggerBy: consts.TriggerByNormal,
		}, "", 0)
		if err != nil {
			t.Log(err)
		} else {
			t.Log(result)
		}
	}))
}

func TestChangeParentIssue(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		ChangeParentIssue(2571, 1368436, vo.ChangeParentIssueReq{
			IssueID:  10108670,
			ParentID: 10108672,
		})
	}))
}

func TestGetRows(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(domain.GetRawRows(47964, 29612, &tablePb.ListRawRequest{
			DbType: tablePb.DbType_slave3,
			//FilterColumns: []string{
			//	lc_helper.ConvertToFilterColumn(consts.BasicFieldPath),
			//	lc_helper.ConvertToFilterColumn("8Ha4"),
			//	lc_helper.ConvertToFilterColumn("8Ha4sdfsdfsd"),
			//	lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
			//	lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
			//},
			Condition: &tablePb.Condition{
				Type: tablePb.ConditionType_and,
				Conditions: domain.GetNoRecycleCondition(
					//domain.GetRowsCondition(consts.BasicFieldIssueId, tablePb.ConditionType_equal, 10120261, nil),
					domain.GetRowsCondition(consts.BasicFieldTableId, tablePb.ConditionType_equal, "1599612991676878848", nil),
					//domain.GetRowsCondition(consts.BasicFieldTableId, tablePb.ConditionType_un_equal, "11233", nil),
				),
			},
		}))
	}))
}
