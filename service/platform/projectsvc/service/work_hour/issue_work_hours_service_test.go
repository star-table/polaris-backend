package work_hour

import (
	"context"
	"testing"
	"time"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	consts2 "github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade/commonfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/cznic/sortutil"
	"github.com/smartystreets/goconvey/convey"
)

// 获取工时记录列表
func TestGetIssueWorkHoursList(t *testing.T) {
	convey.Convey("TestGetIssueWorkHoursList", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		userId := int64(1023)
		issueId := int64(1030)
		page := int64(1)
		size := int64(10)
		reqVo := vo.GetIssueWorkHoursListReq{
			Page:    &page,
			Size:    &size,
			IssueID: issueId,
			Type:    1,
		}
		respVo, err := GetIssueWorkHoursList(orgId, userId, reqVo)
		if err != nil {
			log.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(respVo.List))
	}))
}

//// 新增一条工时记录
//func TestCreateIssueWorkHours(t *testing.T) {
//	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
//		orgId := int64(1004)
//		userId := int64(1012)
//		issueId := int64(1030)
//		desc := "work hours desc " + string(time.Now().Format(consts2.AppTimeFormat))
//		endTime := int64(0)
//		fmt.Printf("场景1：总预估工时已存在--------------------------------\n")
//		resp, err := CreateIssueWorkHours(orgId, userId, vo.CreateIssueWorkHoursReq{
//			IssueID:   issueId,
//			Type:      consts.WorkHourTypeTotalPredict,
//			WorkerID:  0,
//			NeedTime:  "12.8",
//			StartTime: int64(time.Now().Unix()),
//			EndTime:   &endTime,
//			Desc:      &desc,
//		})
//		convey.So(err.Code(), convey.ShouldEqual, errs.SimplePredictIssueWorkHourExist.Code())
//		convey.So(resp.IsTrue, convey.ShouldEqual, false)
//		log.Infof("%v\n", resp)
//		fmt.Printf("场景2：不允许新增子预估工时--------------------------------\n")
//		_, err = CreateIssueWorkHours(orgId, userId, vo.CreateIssueWorkHoursReq{
//			IssueID:   issueId,
//			Type:      consts.WorkHourTypeSubPredict,
//			WorkerID:  1023,
//			NeedTime:  "12.4",
//			StartTime: time.Now().Unix(),
//			EndTime:   &endTime,
//			Desc:      &desc,
//		})
//		convey.So(err.Code(), convey.ShouldEqual, errs.ParamError.Code())
//		fmt.Printf("场景3：无权限增加工时用例--------------------------------\n")
//		userId = 1023
//		resp, err = CreateIssueWorkHours(orgId, userId, vo.CreateIssueWorkHoursReq{
//			IssueID:   issueId,
//			Type:      consts.WorkHourTypeTotalPredict,
//			WorkerID:  0,
//			NeedTime:  "20.3",
//			StartTime: int64(time.Now().Unix()),
//			EndTime:   &endTime,
//			Desc:      &desc,
//		})
//		convey.So(err.Code(), convey.ShouldEqual, errs.DenyCreateIssueWorkHours.Code())
//		fmt.Printf("场景4：正常增加工时--------------------------------\n")
//		userId = 1012
//		issueId = int64(1030)
//		// 先删掉，再新增
//		domain.UpdateIssueWorkHourByCond(db.Cond{
//			consts2.TcIssueId: issueId,
//		}, bo.UpdateOneIssueWorkHoursBo{IsDelete: consts2.AppIsDeleted})
//		resp, err = CreateIssueWorkHours(orgId, userId, vo.CreateIssueWorkHoursReq{
//			IssueID:   issueId,
//			Type:      consts.WorkHourTypeTotalPredict,
//			WorkerID:  userId,
//			NeedTime:  "21.2",
//			StartTime: time.Now().Unix(),
//			EndTime:   &endTime,
//		})
//		if err != nil {
//			t.Error(err)
//		}
//		log.Info(resp)
//	}))
//}

//// 更新工时记录
//func TestUpdateIssueWorkHours(t *testing.T) {
//	convey.Convey("TestCreateIssueWorkHours", t, test.StartUpForSubSvc(func(ctx context.Context) {
//		orgId := int64(2373)
//		userId := int64(29611) // 1023
//		desc := "work hours desc " + string(time.Now().Format(consts2.AppTimeFormat))
//		startTime := time.Now().Unix()
//		endTime := time.Now().Unix() + 3600
//		issuWorkHoursId := int64(5562)
//		resp, err := UpdateIssueWorkHours(orgId, userId, vo.UpdateIssueWorkHoursReq{
//			IssueWorkHoursID:  issuWorkHoursId,
//			WorkerID:          userId,
//			NeedTime:          "2.1",
//			RemainTimeCalType: 1,
//			RemainTime:        0,
//			StartTime:         startTime,
//			EndTime:           endTime,
//			Desc:              &desc,
//		}, true, true, true)
//		if err != nil {
//			log.Errorf("%v\n", err)
//			return
//		}
//		log.Infof("%v\n", resp)
//	}))
//}

//// 删除工时记录
//func TestDeleteIssueWorkHours(t *testing.T) {
//	convey.Convey("TestDeleteIssueWorkHours", t, test.StartUpForSubSvc(func(ctx context.Context) {
//		orgId := int64(2373)
//		userId := int64(29611)
//		// 场景1：删除一个子预估工时
//		issueWorkHoursId := int64(5572)
//		_ = domain.UpdateIssueWorkHours(bo.UpdateOneIssueWorkHoursBo{
//			Id:       uint64(issueWorkHoursId),
//			IsDelete: 2,
//		})
//		_, err := DeleteIssueWorkHours(orgId, userId, vo.DeleteIssueWorkHoursReq{
//			IssueWorkHoursID: issueWorkHoursId,
//		}, true, true)
//		if err != nil {
//			log.Errorf("%v\n", err)
//			return
//		}
//		// 场景2：删除一个已经删除的，提示找不到记录
//		issueWorkHoursId = 5573
//		_, err = DeleteIssueWorkHours(orgId, userId, vo.DeleteIssueWorkHoursReq{
//			IssueWorkHoursID: issueWorkHoursId,
//		}, true, true)
//		convey.ShouldEqual(err.Code(), 10203) // 10203：db操作出现异常
//		log.Info("end...")
//	}))
//}

// 一个项目中，开启、关闭工时功能
func TestDisOrEnableIssueWorkHours(t *testing.T) {
	convey.Convey("DisOrEnableIssueWorkHours", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		userId := int64(1012) // 1012
		projectId := int64(1007)
		resp, err := DisOrEnableIssueWorkHours(orgId, userId, vo.DisOrEnableIssueWorkHoursReq{
			ProjectID: projectId,
			Enable:    1, // 1开启；2关闭
		})
		if err != nil {
			log.Errorf("%v\n", err)
			return
		}
		log.Infof("%v\n", resp)
	}))
}

func TestSliceContain(t *testing.T) {
	convey.Convey("DisOrEnableIssueWorkHours", t, test.StartUpForSubSvc(func(ctx context.Context) {
		participantIds := []int64{1, 2}
		currentUserId1 := int64(2)
		currentUserId2 := int64(3)
		if ok, _ := slice.Contain(participantIds, currentUserId1); true {
			convey.So(ok, convey.ShouldEqual, true)
		}
		if ok, _ := slice.Contain(participantIds, currentUserId2); true {
			convey.So(ok, convey.ShouldEqual, false)
		}
	}))
}

// 获取一个任务的工时信息
func TestGetIssueWorkHoursInfo(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1146)
		userId := int64(1010)  //  项目负责人 1010
		issueId := int64(6056) // 1030 // 任务负责人 1023
		resp, err := GetIssueWorkHoursInfo(orgId, userId, vo.GetIssueWorkHoursInfoReq{IssueID: issueId})
		convey.ShouldEqual(err, nil)
		convey.ShouldNotEqual(resp.SimplePredictWorkHour, nil)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info(resp.PredictWorkHourList)
		if len(resp.PredictWorkHourList) > 0 {
			log.Info(*resp.PredictWorkHourList[0])
		}
		log.Info(resp.ActualWorkHourList)
		if len(resp.ActualWorkHourList) > 0 {
			log.Info(*resp.ActualWorkHourList[0])
		}
		log.Info(resp.SimplePredictWorkHour)
	}))
}

//// 创建多个子预估工时
//func TestCreateMultiIssueWorkHours(t *testing.T) {
//	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
//		orgId := int64(1574)
//		userId := int64(24370) // 1023
//		issueId := int64(60862)
//		startTime := time.Now().Unix()
//		endTime := int64(1628935200)
//		newDataList := []*vo.NewPredicateWorkHour{
//			{
//				WorkerID:  24362,
//				NeedTime:  "2.1",
//				StartTime: 1628816400,
//				EndTime:   &endTime,
//			},
//		}
//		resp, err := CreateMultiIssueWorkHours(orgId, userId, vo.CreateMultiIssueWorkHoursReq{
//			IssueID: issueId,
//			TotalIssueWorkHourRecord: &vo.NewPredicateWorkHour{
//				WorkerID:  userId,
//				NeedTime:  "2.1",
//				StartTime: startTime,
//				EndTime:   &endTime,
//			},
//			PredictWorkHourList: newDataList,
//		})
//		if err != nil {
//			log.Error(err)
//			return
//		}
//		log.Info(resp)
//	}))
//}

// 工时统计的查询
func TestGetWorkHourStatistic(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1872)
		// projectId := int64(10809)
		//workerId := int64(1007)
		page := int64(-1)
		size := int64(20)
		t1 := time.Date(2022, 04, 27, 0, 0, 0, 0, time.Local).Unix()
		t2 := time.Date(2022, 05, 03, 0, 0, 0, 0, time.Local).Unix()
		showResigned := int(2)
		resp, err := GetWorkHourStatistic(orgId, vo.GetWorkHourStatisticReq{
			//WorkerIds: []*int64{&workerId},
			// ProjectIds:   []*int64{int642Ptr(13446)},
			StartTime:    &t1,
			EndTime:      &t2,
			ShowResigned: &showResigned,
			Page:         &page,
			Size:         &size,
		})
		if err != nil {
			log.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

// 生成日期序列
func TestGetDateListByDateRange(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		ts1, _ := time.Parse(consts2.AppTimeFormat, "2020-10-20 00:00:00")
		ts2, _ := time.Parse(consts2.AppTimeFormat, "2020-10-23 00:00:00")
		dateList := GetDateListByDateRange(ts1.Unix(), ts2.Unix())
		log.Info(dateList)
	}))
}

func TestTriggerCreateWorkHourTrend(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		userId := int64(1012)
		issueId := int64(1030)
		startTime := int64(time.Now().Unix())
		endTime := int64(time.Now().Unix() + 3600*12)
		desc := ""
		// 查询工时执行人的信息
		userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, []int64{userId})
		if err != nil {
			log.Error(err)
			return
		}
		userInfoMapById := map[int64]*vo.WorkHourWorker{}
		for _, user := range userInfos {
			userInfoMapById[user.UserId] = &vo.WorkHourWorker{
				UserID: user.UserId,
				Name:   user.Name,
				Avatar: user.Avatar,
			}
		}
		worker := &vo.WorkHourWorker{}
		if val, ok := userInfoMapById[userId]; ok {
			worker.UserID = val.UserID
			worker.Name = val.Name
			worker.Avatar = val.Avatar
		}
		err = TriggerCreateWorkHourTrend(orgId, userId, issueId, vo.OneWorkHourRecord{
			ID:        1030,
			Type:      consts.WorkHourTypeSubPredict,
			Worker:    worker,
			NeedTime:  "2",
			StartTime: int64(startTime),
			EndTime:   endTime,
			Desc:      desc,
		})
		if err != nil {
			log.Error(err)
		}
	}))
}

func TestGenerateTrendForCreate(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		userId := int64(1012)
		issueId := int64(1030)
		startTime := int64(time.Now().Unix())
		endTime := int64(time.Now().Unix() + 3600*12)
		err := generateTrendForCreate(orgId, userId, issueId, []*po.PpmPriIssueWorkHours{
			{
				Id:                uint64(1102),
				OrgId:             orgId,
				IssueId:           issueId,
				Type:              consts.WorkHourTypeSubPredict,
				WorkerId:          userId,
				NeedTime:          21,
				RemainTimeCalType: consts.RemainTimeCalTypeDefault,
				RemainTime:        0,
				StartTime:         uint32(startTime),
				EndTime:           uint32(endTime),
				Desc:              "",
				Creator:           userId,
				CreateTime:        time.Time{},
				Updator:           userId,
				UpdateTime:        time.Time{},
				Version:           0,
				IsDelete:          consts2.AppIsNoDelete,
			},
		})
		convey.ShouldEqual(err, nil)
		if err != nil {
			log.Error(err)
		}
	}))
}

// 更新工时生成动态，删除工时生成动态
func TestTriggerUpdateWorkHourTrend(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		userId := int64(1012)
		issueId := int64(1030)
		workHourId := int64(1003)
		// 工时记录是否存在
		oneIssueWorkHour, err := domain.GetOneIssueWorkHoursById(workHourId)
		if err != nil {
			log.Error(err)
			return
		}
		editValue := oneIssueWorkHour
		editValue.NeedTime = 100
		editValue.StartTime = editValue.StartTime + 3600
		editValue.EndTime = uint32(time.Now().Unix() + 1)
		oldValue := vo.OneWorkHourRecord{}
		newValue := vo.OneWorkHourRecord{}
		copyer.Copy(oneIssueWorkHour, &oldValue)
		copyer.Copy(editValue, &newValue)
		// 查询工时执行人的信息
		userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, []int64{userId})
		if err != nil {
			log.Error(err)
			return
		}
		userInfoMapById := map[int64]*vo.WorkHourWorker{}
		for _, user := range userInfos {
			userInfoMapById[user.UserId] = &vo.WorkHourWorker{
				UserID: user.UserId,
				Name:   user.Name,
				Avatar: user.Avatar,
			}
		}
		worker := &vo.WorkHourWorker{}
		if val, ok := userInfoMapById[userId]; ok {
			worker.UserID = val.UserID
			worker.Name = val.Name
			worker.Avatar = val.Avatar
		}
		oldValue.Worker = worker
		newValue.Worker = worker
		oriErr := TriggerUpdateWorkHourTrend(orgId, userId, issueId, oldValue, newValue)
		if oriErr != nil {
			log.Error(oriErr)
			return
		}
		log.Info("更新工时生成动态 complete...")
		oriErr = TriggerDeleteWorkHourTrend(orgId, userId, issueId, oldValue)
		if oriErr != nil {
			log.Error(oriErr)
			return
		}
		log.Info("删除工时生成动态 complete...")

		log.Info("executed complete...")
	}))
}

// 测试变更任务负责人，简单预估工时执行人，则也发生变化
func TestChangeWorkerIdWhenChangedIssueOwner(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		issueId := int64(1030)
		oldWorkerId := int64(1012)
		newWorkerId := int64(1023)
		err := ChangeWorkerIdWhenChangedIssueOwner(orgId, issueId, oldWorkerId, newWorkerId)
		if err != nil {
			log.Error(err)
			return
		}
		log.Info("end...")
	}))
}

func TestTriggerWorkHourWhenChangedPlanStartOrEndTime(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		issueId := int64(1030)
		t1Str := "2020-11-10 11:00:00"
		t1 := types.NowTime()
		t1.UnmarshalJSON([]byte(t1Str))
		t2Str := "2020-11-10 11:22:00"
		t2 := types.NowTime()
		t2.UnmarshalJSON([]byte(t2Str))
		err := TriggerWorkHourWhenChangedPlanStartOrEndTime(orgId, issueId, "start_time", t1, t2)
		if err != nil {
			t.Error(err)
		}
		oldValue := types.Time(time.Date(2020, 11, 11, 0, 0, 0, 0, time.Local))
		newValue := types.Time(time.Date(2020, 11, 11, 10, 0, 0, 0, time.Local))
		err = TriggerWorkHourWhenChangedPlanStartOrEndTime(orgId, issueId, "end_time", oldValue, newValue)
		if err != nil {
			t.Error(err)
		}
	}))
}

func TestSort(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		totalWorkerIds := []int64{12, 89, 17, 8, 99, 172, 65, 176}
		var sortPadWorkerIds sortutil.Int64Slice = totalWorkerIds
		sortPadWorkerIds.Sort()
		totalWorkerIds = sortPadWorkerIds
		t.Log(json.ToJsonIgnoreError(totalWorkerIds))
	}))
}

func TestGetIssueIdsByStatisticInput(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		projectIds := []int64{1005, 1006, 1007}
		projectIds1 := []*int64{}
		copyer.Copy(projectIds, &projectIds1)
		resp, err := GetIssueIdsByStatisticInput(orgId, vo.GetWorkHourStatisticReq{
			ProjectIds:   []*int64{},
			EndTime:      int642Ptr(1651593600),
			ShowResigned: int2Ptr(2),
			Page:         int642Ptr(1),
			Size:         int642Ptr(20),
		})
		convey.ShouldNotEqual(err, nil)
		if err != nil {
			log.Error(err)
			return
		}
		t.Log(resp)
	}))
}

func int2Ptr(val int) *int {
	return &val
}

func int642Ptr(val int64) *int64 {
	return &val
}

// 检查用户是否是关注者
func TestCheckIsIssueMember(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1574)
		//userId := int64(1012)
		issueId := int64(59255)
		// 场景1：不是关注者
		resp, err := CheckIsIssueMember(orgId, vo.CheckIsIssueMemberReq{
			IssueID: issueId,
			UserID:  int64(24370),
		})
		convey.ShouldEqual(resp.IsTrue, false)
		convey.ShouldEqual(err, nil)
		//// 场景1：是关注者
		//resp, err = CheckIsIssueMember(orgId, vo.CheckIsIssueMemberReq{
		//	IssueID: issueId,
		//	UserID:  userId,
		//})
		//convey.ShouldEqual(err, nil)
		//convey.ShouldEqual(resp.IsTrue, true)
		//if err != nil {
		//	log.Error(err)
		//	return
		//}
		//t.DataEvent(resp)
	}))
}

func TestSetUserJoinIssueReq(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1449)
		userId := int64(23214)
		newUserId := int64(23217)
		issueId := int64(68287)
		resp, err := SetUserJoinIssue(orgId, userId, vo.SetUserJoinIssueReq{
			IssueID: issueId,
			UserID:  newUserId,
		})
		if err != nil {
			log.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

//func TestUpdateMultiIssueWorkHourWithDelete(t *testing.T) {
//	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
//		orgId := int64(1004)
//		userId := int64(1012)
//		issueId := int64(1030)
//		startTime := time.Now().Unix()
//		endTime := time.Now().Unix() + 3600*22
//		desc := ""
//		resp, err := UpdateMultiIssueWorkHourWithDelete(orgId, userId, vo.UpdateMultiIssueWorkHoursReq{
//			IssueID: issueId,
//			TotalIssueWorkHourRecord: &vo.UpdateOneMultiWorkHourRecord{
//				ID:        1111,
//				Type:      consts.WorkHourTypeTotalPredict,
//				WorkerID:  userId,
//				NeedTime:  "21",
//				StartTime: startTime,
//				EndTime:   endTime,
//				Desc:      &desc,
//			},
//			IssueWorkHourRecords: []*vo.UpdateOneMultiWorkHourRecord{
//				{
//					ID:        0,
//					Type:      consts.WorkHourTypeSubPredict,
//					WorkerID:  userId,
//					NeedTime:  "1.5",
//					StartTime: startTime,
//					EndTime:   endTime,
//					Desc:      &desc,
//				}, {
//					ID:        0,
//					Type:      consts.WorkHourTypeSubPredict,
//					WorkerID:  userId,
//					NeedTime:  "1.6",
//					StartTime: startTime,
//					EndTime:   endTime,
//					Desc:      &desc,
//				}, {
//					ID:        1109,
//					Type:      consts.WorkHourTypeSubPredict,
//					WorkerID:  userId,
//					NeedTime:  "1.3",
//					StartTime: startTime,
//					EndTime:   endTime,
//					Desc:      &desc,
//				},
//			},
//		}, true)
//		convey.ShouldNotEqual(err, nil)
//		convey.ShouldEqual(resp.IsTrue, true)
//		if err != nil {
//			log.Error(err)
//			return
//		}
//		t.Log(resp)
//	}))
//}

func TestGenerateTrendForDelete(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		userId := int64(1012)
		issueId := int64(1030)
		err := generateTrendForDelete(orgId, userId, issueId, []int64{1095, 1096})
		convey.ShouldEqual(err, nil)
		if err != nil {
			log.Error(err)
		}
	}))
}

func TestCheckIsEnableWorkHour(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		userId := int64(1012)
		projectId := int64(1007)
		resp, err := CheckIsEnableWorkHour(orgId, userId, vo.CheckIsEnableWorkHourReq{ProjectID: projectId})
		convey.ShouldEqual(err, nil)
		if err != nil {
			t.Error(err)
		}
		if resp.IsEnable != true {
			t.Error(resp)
		}
		log.Info(resp)
	}))
}

// 变更工时的项目id
func TestUpdateWorkHourProjectId(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		oldProjectId := int64(1008)
		newProjectId := int64(1007)
		issueIds := []int64{1030}
		err := domain.UpdateWorkHourProjectId(orgId, oldProjectId, newProjectId, issueIds, false)
		if err != nil {
			t.Error(err)
		}
	}))
}

// 向钉钉群发送告警
func TestSendDingTalkLog(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		resp := commonfacade.DingTalkInfo(commonvo.DingTalkInfoReqVo{
			LogData: map[string]interface{}{
				"app": config.GetApplication().Name,
				"env": config.GetEnv(),
				"msg": "test by suhanyu. 【警告日志】更换任务的负责人时，对应的总预估工时执行者也随着发生变化，处理异常。",
				"err": errs.FileNotExist,
			},
		})
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestExportWorkHourStatistic(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1004)
		projectId := int64(1007)
		page := int64(-1)
		size := int64(5)
		t1 := time.Date(2020, 11, 20, 0, 0, 0, 0, time.Local).Unix()
		t2 := time.Date(2020, 11, 26, 0, 0, 0, 0, time.Local).Unix()
		showResigned := int(2)
		resp, err := ExportWorkHourStatistic(orgId, vo.GetWorkHourStatisticReq{
			ProjectIds:   []*int64{&projectId},
			ShowResigned: &showResigned,
			StartTime:    &t1,
			EndTime:      &t2,
			Page:         &page,
			Size:         &size,
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestGetIssuesWorkHourList(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(1004)
		projectId := int64(1007)
		issueIds := []int64{1030}
		list, err := domain.GetIssuesWorkHourList(orgId, []int64{projectId}, issueIds)
		if err != nil {
			t.Error(err)
		}
		t.Log(json.ToJsonIgnoreError(list))
	}))
}

func TestSetWorkHourMap(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(1004)
		projectId := int64(1007)
		issueIds := []int64{1030}
		workHourMap := map[int64]bo.IssueWorkHour{}
		domain.SetWorkHourMap(orgId, []int64{projectId}, issueIds, workHourMap)

		t.Log(json.ToJsonIgnoreError(workHourMap))
	}))
}

func TestCheckEditWorkHourPer1(t *testing.T) {
	convey.Convey("Test", t, test.StartUpForSubSvc(func(ctx context.Context) {
		orgId := int64(1146)
		issueId := int64(8858)
		issues, err := domain.GetIssueInfosLc(orgId, 0, []int64{issueId})
		if err != nil {
			t.Error(err)
			return
		}
		if len(issues) < 1 {
			err = errs.IssueNotExist
			t.Error(err)
			return
		}
		issue := issues[0]

		res, err := domain.CheckHasWriteForWorkHourColumn(orgId, 1012, 1519192571513749505, issue)
		if err != nil {
			log.Error(err)
			return
		}
		t.Log(res)
		t.Log("--end--")
	}))
}

//func TestCreateMultiIssueWorkHours2(t *testing.T) {
//	convey.Convey("test", t, test.StartUpForSubSvc(func(ctx context.Context) {
//		projectId := int64(60503)
//		hours, err := CreateIssueWorkHours(2373, 29612, vo.CreateIssueWorkHoursReq{
//			ProjectID: &projectId,
//			IssueID:   10030286,
//			Type:      2,
//			WorkerID:  29612,
//			NeedTime:  "1",
//		})
//		if err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(hours)
//	}))
//}
