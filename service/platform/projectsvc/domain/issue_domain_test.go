package domain

import (
	"context"
	"fmt"
	"testing"

	tablePb "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/polaris-backend/common/extra/lc_helper"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

//func TestCreateIssue(t *testing.T) {
//
//	config.LoadConfig("/Users/tree/work/08_all_star/01_src/go/polaris-backend/polaris-server/configs", "application")
//
//	pid, _ := idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue")
//
//	orgId := int64(1000)
//	projectId := int64(1)
//
//	code, _ := idfacade.ApplyCode(orgId, "PC", "")
//
//	issueBo := &bo.IssueBo{
//		Id:                  pid,
//		OrgId:               orgId,
//		Code:                code,
//		ProjectId:           projectId,
//		TableId: 1234,
//		Title:               "1AAA",
//		Owner:               345,
//		PriorityId:          234,
//		SourceId:            345,
//		IssueObjectTypeId:   456,
//		PlanStartTime:       types.NowTime(),
//		PlanEndTime:         types.NowTime(),
//		StartTime:           types.NowTime(),
//		EndTime:             types.NowTime(),
//		PlanWorkHour:        4,
//		IterationId:         0,
//		VersionId:           0,
//		ModuleId:            0,
//		ParentId:            0,
//		Status:              567,
//		Creator:             345,
//		CreateTime:          types.NowTime(),
//		Updator:             345,
//		UpdateTime:          types.NowTime(),
//		Version:             1,
//	}
//
//	pdid, _ := idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue_detail")
//
//	issueDetailBo := &bo.IssueDetailBo{
//		Id:         pdid,
//		OrgId:      orgId,
//		IssueId:    pid,
//		ProjectId:  projectId,
//		StoryPoint: 0,
//		Tags:       "tags",
//		Remark:     consts.TcRemark,
//		Status:     123,
//		Creator:    345,
//		CreateTime: types.NowTime(),
//		Updator:    345,
//		UpdateTime: types.NowTime(),
//		Version:    1,
//	}
//
//	issueBo.IssueDetailBo = *issueDetailBo
//
//	prid, _ := idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue_relation")
//	issueBo.OwnerInfo = &bo.IssueUserBo{
//		IssueRelationBo: bo.IssueRelationBo{
//			Id:           prid,
//			OrgId:        orgId,
//			IssueId:      pid,
//			RelationId:   345,
//			RelationType: 2,
//			Creator:      345,
//			CreateTime:   types.NowTime(),
//			Updator:      345,
//			UpdateTime:   types.NowTime(),
//			Version:      1,
//		},
//	}
//
//	relationBos := make([]bo.IssueUserBo, 10)
//
//	for i := int64(0); i < 10; i++ {
//		prid, _ = idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue_relation")
//		relationBos[i] = bo.IssueUserBo{
//			IssueRelationBo: bo.IssueRelationBo{
//				Id:           prid,
//				OrgId:        orgId,
//				IssueId:      pid,
//				RelationId:   i + 10,
//				RelationType: 2,
//				Creator:      345,
//				CreateTime:   types.NowTime(),
//				Updator:      345,
//				UpdateTime:   types.NowTime(),
//				Version:      1,
//			},
//		}
//	}
//
//	relationBos2 := make([]bo.IssueUserBo, 10)
//
//	for i := int64(0); i < 10; i++ {
//		prid, _ = idfacade.ApplyPrimaryIdRelaxed("ppm_pri_issue_relation")
//		relationBos2[i] = bo.IssueUserBo{
//			IssueRelationBo: bo.IssueRelationBo{
//				Id:           prid,
//				OrgId:        orgId,
//				IssueId:      pid,
//				RelationId:   i + 20,
//				RelationType: 2,
//				Creator:      345,
//				CreateTime:   types.NowTime(),
//				Updator:      345,
//				UpdateTime:   types.NowTime(),
//				Version:      1,
//			},
//		}
//	}
//
//	issueBo.FollowerInfos = &relationBos2
//	issueBo.ParticipantInfos = &relationBos
//
//	err3 := CreateIssue(issueBo)
//
//	fmt.Println(err3)
//	fmt.Println(issueBo.Id)
//}

//func TestGetIssueBo(t *testing.T) {
//	convey.Convey("获取任务bo", t, tests.StartUp(func() {
//		convey.Convey("获取任务bo", func() {
//			issueBo, _ := GetIssueBo(1000, 1795)
//			fmt.Println(issueBo.PlanStartTime)
//			fmt.Println(issueBo.PlanEndTime)
//
//			newIssueBo := &bo.IssueBo{}
//			_ = copyer.Copy(issueBo, newIssueBo)
//
//			fmt.Println(newIssueBo.PlanStartTime)
//			fmt.Println(newIssueBo.PlanEndTime)
//		})
//	}))

func TestConvertIssueTagBosToMapGroupByIssueId(t *testing.T) {

	issueTagBos := []bo.IssueTagBo{
		{
			IssueId: 1,
			TagId:   1,
			TagName: "hello",
		},
		{
			IssueId: 1,
			TagId:   2,
			TagName: "world",
		},
		{
			IssueId: 2,
			TagId:   1,
			TagName: "hello",
		},
		{
			IssueId: 2,
			TagId:   2,
			TagName: "world",
		},
	}
	lm := ConvertIssueTagBosToMapGroupByIssueId(issueTagBos)
	t.Log(json.ToJsonIgnoreError(lm))
}

func TestGetCalendarInfo(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetCalendarInfo(1082, 1744))
	}))

}

func TestGetRunningIterationStat(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetRunningIterationStat([]int64{1374}, 1111))
	}))
}

//func TestGetIssueInfoByProject(t *testing.T) {
//	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
//		t.Log(GetIssueInfoByProject([]int64{1374}, 1111, 1284))
//	}))
//}

//func TestGetIssueAndChildren(t *testing.T) {
//	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
//		t.Log(GetIssueAndChildren(1004, []int64{4572}, nil))
//	}))
//}

func TestUpdateIssueProjectObjectType(t *testing.T) {
	a := mysql.Upd{
		"a": "1",
		"b": "2",
	}

	b := a
	b["c"] = "3"

	t.Log(json.ToJsonIgnoreError(b))
	t.Log(json.ToJsonIgnoreError(a))
}

func TestGetissueInfo(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		issueBos, err := GetIssueInfosLc(1000023, 0, []int64{5846})
		if err != nil {
			t.Log(err)
		}
		t.Log(issueBos)
	}))
}

//func TestGetIssueIdListByConditions(t *testing.T) {
//	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
//		issueBos, err := GetIssueIdListByConditions(1000023, []*vo.LessCondsData{
//			consts.WrapperLessDataAndCondition(consts.ConditionIn, consts.BasicFieldId, nil, []int64{5846}),
//		}, false)
//		if err != nil {
//			t.Log(err)
//		}
//		t.Log(issueBos[0])
//	}))
//}

func TestConvertInterfaceIntoInt64(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		var val interface{}
		val = "opt_12"
		intoInt64 := ConvertInterfaceIntoInt64(val)
		t.Log(intoInt64)
	}))
}

func TestGetIssueInfosMapLc(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		orgId := 2373

		infosMapLc, err := GetIssueInfosMapLc(int64(orgId), 0, nil, nil, 1, 2000)
		if err != nil {
			t.Error(err)
		}
		for _, issueMap := range infosMapLc {
			colIMap := issueMap["collaborators"]
			if col, ok := colIMap.(map[string]interface{}); ok {
				for _, c := range col {
					if cI, ok2 := c.([]interface{}); ok2 {
						if len(cI) == 0 {
							continue
						}
					}
				}
			}
			fmt.Println(colIMap)
		}
		fmt.Println(infosMapLc)
	}))
}

func TestGetIssueInfosMapLcByIssueIds(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		resp, err := GetIssueInfosMapLcByIssueIds(2373, 29610, []int64{10093481}, lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle))
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
	}))
}

func TestGetRawRows(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		_, err := GetRawRows(14396, 29610, &tablePb.ListRawRequest{
			DbType: tablePb.DbType_slave3,
			FilterColumns: []string{
				lc_helper.ConvertToFilterColumn(consts.BasicFieldId),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldCode),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldOwnerId),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldAuditStatus),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldAuditStatusDetail),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueStatus),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueStatusType),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldPlanStartTime),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldPlanEndTime),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldAppId),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldRemark),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldFollowerIds),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldAuditorIds),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldParentId),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldRelating),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldBaRelating),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldUpdateTime),
				lc_helper.ConvertToFilterColumn(consts.BasicFieldCreator),
			},
			Condition: &tablePb.Condition{Type: tablePb.ConditionType_and, Conditions: GetNoRecycleCondition(
				//GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_in, nil, []int64{16396, 26094, 26785, 27874, 28192, 29619, 29688, 29713, 29783, 29924, 31049, 31074, 31346, 31553, 31555, 31565, 31717, 31769, 32016, 32116, 32180, 32193, 32284, 32288, 32303, 32305, 32389, 32452, 32605, 33527, 33528, 33532, 33533, 33534, 33613, 33647, 33648, 33649, 34254, 34490, 34682, 34683, 34686, 34726, 34737, 34738, 34740, 34755, 34838, 34843, 34882, 34887, 34891, 34892, 34928, 35066, 35088, 35089, 35211, 35214, 35220, 35246, 35248, 35977, 35978, 36098, 36101, 36424, 36612, 36613, 36614, 36615, 36618, 36700, 36850, 37303, 37707, 38239, 38591, 39175, 39519, 39562, 41030, 41041, 41132, 41215, 41274, 41361, 41363, 41624, 42409, 42643, 42647, 43223, 43383, 43389, 43503, 43689, 43919, 44484, 44492, 44875, 44882, 45260, 45333, 45334, 45335, 45508, 45510, 45605, 45633, 45658, 45659, 45789, 45828, 45844, 45882, 46584, 46606, 47219, 47276, 47456, 48059, 48523, 50173, 51577, 53707, 53709, 54127, 55039, 58142, 58211, 58409, 58600, 58779, 58901, 58902, 59474, 59484, 59494, 60374, 60420, 60428, 60429, 60430, 60453, 60455, 60487, 60493, 60537, 60540, 60542, 60543, 60545, 60547, 60549, 60580, 60587, 60588, 60708, 60775, 60778, 61199, 61201, 61357, 61365, 61497, 61774, 62181, 62449, 63278, 63305, 63310, 63800, 64185, 64187, 64225, 64378, 64738, 65023, 65361, 65427, 65487, 65496, 65499, 65510, 65514, 65524, 65591, 65623, 65637, 65726, 65822, 65924, 65932, 66122, 66362, 66771, 67264, 67267, 67293, 67533, 67919, 67941, 67966, 67976, 67978, 68135, 68151, 68260, 68266, 68268, 68757, 68823, 69066, 69088, 70165, 70194, 70788, 70890, 71074, 71241, 71327, 71727, 72702, 72707, 72708, 73711, 73717, 73719, 74080, 74193, 74258, 74898, 75642, 76407, 78304, 78429, 78572, 78965, 79703, 79802, 81220, 81338, 83024, 83032, 83035, 83364, 83381, 83818, 83822, 83825, 83831, 83832, 83851, 83900, 84009, 84156, 84262, 84263, 84590, 84598, 84676, 84693, 84715, 84774, 85074, 85085, 85157, 85172, 85776, 85787, 87444, 87445, 88011, 90484, 90485, 90496, 92426, 92427, 92498, 92637, 92723, 92724, 92725, 92833, 93549, 94760, 94774, 95219, 95949, 96347, 96352, 96519, 96620, 96906, 97228, 97237, 97312, 97319, 97538, 97617, 97678, 97679, 97681, 98453, 98517, 98651, 98703, 99711, 100430, 100576, 0}),
				GetRowsCondition(consts.BasicFieldTableId, tablePb.ConditionType_equal, 1592758586046418946, nil),
				//GetRowsCondition(consts.BasicFieldUpdateTime, tablePb.ConditionType_gte, "2022-11-22 15:33:46", nil),
			)},
			//Orders: []*tablePb.Order{{Asc: true, Column: lc_helper.ConvertToFilterColumn(consts.BasicFieldUpdateTime)}},
			Page: 1,
			Size: 10000,
		})
		if err != nil {
			fmt.Println(err)
		}
	}))
}
