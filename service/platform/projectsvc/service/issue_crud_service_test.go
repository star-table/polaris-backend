package service

import (
	"context"
	"fmt"
	"sort"
	"testing"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

//func TestCreateIssue(t *testing.T) {
//
//	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
//
//		fmt.Printf("redis的配置.....%v", config.RedisConfig{})
//
//		fmt.Printf("数据库的的配置.....%v", config.RedisConfig{})
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息：" + cacheUserInfoJson)
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "ou_c980b4120a9156f98dc33642a10a642a", SourceChannel: "fs", UserId: int64(2355), CorpId: "2eb972ddc1cf1652", OrgId: 1417}
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("Test CreateProject", func() {
//
//			rand.Seed(time.Now().Unix())
//			intn := rand.Intn(10000)
//
//			var preCode = "alanTest" + strconv.Itoa(intn)
//
//			startTime := types.NowTime()
//
//			endTime, _ := time.Parse(consts.AppTimeFormat, "2021-01-06 12:20:22")
//
//			planEndTime := types.Time(endTime)
//
//			fmt.Println("时间", endTime)
//
//			var remark = "哈哈哈"
//
//			projectName := "alan测试项目" + strconv.Itoa(intn)
//
//			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1417, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}
//
//			reqVo := projectvo.CreateProjectReqVo{
//				Input:  input,
//				UserId: cacheUserInfo.UserId,
//				OrgId:  cacheUserInfo.OrgId,
//			}
//
//			project := projectfacade.CreateProject(reqVo)
//
//			if project.Failure() {
//				log.Info("创建Issue中..........创建项目失败")
//				return
//			}
//
//			log.Infof("创建issueUnittest中的 project %v", project)
//
//			convey.Convey("Test CreateIssue", func() {
//
//				projectId := project.Project.ID
//
//				rand.Seed(time.Now().Unix())
//				intn := rand.Intn(10000)
//				//projectObjectId
//				typeId := int64(intn)
//
//				currentUserId := cacheUserInfo.UserId
//				orgId := cacheUserInfo.OrgId
//
//				convey.Convey("createIssue insert Project proecss", func() {
//
//					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())
//
//					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
//						Id:                  processId,
//						OrgId:               orgId,
//						ProjectId:           projectId,
//						ProjectObjectTypeId: typeId,
//						ProcessId:           1120, //这个是一开始默认的不能写死
//						Creator:             currentUserId,
//						CreateTime:          time.Now(),
//						Updator:             currentUserId,
//						UpdateTime:          time.Now(),
//						Version:             1,
//					}
//
//					_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						err4 := mysql.TransInsert(tx, projectObjectTypeProcess)
//						if err4 != nil {
//							log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						}
//						return nil
//					})
//
//					convey.Convey("insertPriority.... ", func() {
//
//						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
//							OrgID: orgId,
//							Type:  consts.PriorityTypeIssue})
//
//						convey.Convey("mock createIssue input", func() {
//
//							var remark = "这是我创建的第一个issue"
//							var issueObjectId = int64(1)
//
//							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//							if err != nil {
//								log.Error(err)
//								return
//							}
//
//							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, OwnerID: []int64{currentUserId}, Remark: &remark,
//								IssueObjectID: &issueObjectId}
//
//							reqVo := projectvo.CreateIssueReqVo{
//								CreateIssue: input,
//								OrgId:       cacheUserInfo.OrgId,
//								UserId:      cacheUserInfo.UserId,
//							}
//
//							convey.Convey("createIssue done", func() {
//
//								issue, err := CreateIssue(reqVo)
//								convey.So(issue, convey.ShouldNotBeNil)
//								convey.So(err, convey.ShouldBeNil)
//							})
//						})
//					})
//				})
//			})
//		})
//	}))
//}

//func TestCreateProject2(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		//projectObjectTypeId := int64(11407)
//		resp, err := CreateIssue(projectvo.CreateIssueReqVo{
//			CreateIssue: vo.CreateIssueReq{
//				ProjectID:  8161,
//				Title:      "manyPro-002",
//				PriorityID: 1793,
//				//TypeID:         &projectObjectTypeId,
//				OwnerID:        []int64{1624},
//				ParticipantIds: nil,
//				FollowerIds:    []int64{1624},
//			},
//			UserId:        1624,
//			OrgId:         1213,
//			SourceChannel: "fs",
//		})
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		log.Error(json.ToJsonIgnoreError(resp))
//	}))
//}

//func TestCreateChildIssue(t *testing.T) {
//
//	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息" + cacheUserInfoJson)
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
//
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//		cacheUserInfoJson, err := json.ToJson(cacheUserInfo)
//
//		if err == nil {
//
//		}
//
//		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("Test CreateProject", func() {
//
//			rand.Seed(time.Now().Unix())
//			intn := rand.Intn(10000)
//
//			var preCode = "alanTest" + strconv.Itoa(intn)
//
//			startTime := types.NowTime()
//
//			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")
//
//			planEndTime := types.Time(endTime)
//
//			fmt.Println("时间", endTime)
//
//			var remark = "哈哈哈"
//
//			projectName := "alan测试项目" + strconv.Itoa(intn)
//
//			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}
//
//			project, err := CreateProject(projectvo.CreateProjectReqVo{
//				OrgId:  cacheUserInfo.OrgId,
//				UserId: cacheUserInfo.UserId,
//				Input:  input,
//			})
//
//			if err != nil {
//				fmt.Println("错误.......")
//			}
//
//			fmt.Printf("项目%+v", project)
//
//			convey.Convey("Test CreateIssue", func() {
//
//				var projectId = int64(1204)
//
//				rand.Seed(time.Now().Unix())
//				intn := rand.Intn(10000)
//				//projectObjectId
//				typeId := int64(intn)
//
//				currentUserId := cacheUserInfo.UserId
//				orgId := cacheUserInfo.OrgId
//
//				convey.Convey("createIssue insert Project proecss", func() {
//
//					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())
//
//					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
//						Id:                  processId,
//						OrgId:               orgId,
//						ProjectId:           projectId,
//						ProjectObjectTypeId: typeId,
//						ProcessId:           1120, //这个是一开始默认的不能写死
//						Creator:             currentUserId,
//						CreateTime:          time.Now(),
//						Updator:             currentUserId,
//						UpdateTime:          time.Now(),
//						Version:             1,
//					}
//
//					_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						err4 := mysql.TransInsert(tx, projectObjectTypeProcess)
//						if err4 != nil {
//							log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						}
//						return nil
//					})
//
//					v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
//						OrgID: orgId,
//						Type:  consts.PriorityTypeIssue})
//
//					convey.Convey("mock createIssue input", func() {
//
//						var remark = "这是我创建的第一个issue"
//						var issueObjectId = int64(1)
//						children := make([]*vo.IssueChildren, 0)
//
//						firstChild := &vo.IssueChildren{
//							Title: "alan子issue",
//							// 负责人
//							OwnerID: []int64{currentUserId},
//							// 优先级
//							PriorityID: v.ID}
//
//						children = append(children, firstChild)
//
//						cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//						if err != nil {
//							log.Error(err)
//							return
//						}
//
//						input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, OwnerID: []int64{currentUserId}, Remark: &remark,
//							IssueObjectID: &issueObjectId, Children: children}
//
//						reqVo := projectvo.CreateIssueReqVo{
//							CreateIssue: input,
//							OrgId:       cacheUserInfo.OrgId,
//							UserId:      cacheUserInfo.UserId,
//						}
//
//						convey.Convey("createIssue done", func() {
//
//							issue, err := CreateIssue(reqVo)
//							convey.So(issue, convey.ShouldNotBeNil)
//							convey.So(err, convey.ShouldBeNil)
//						})
//					})
//				})
//			})
//		})
//	}))
//}

//更新issue
//func TestUpdateIssue(t *testing.T) {
//
//	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
//
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("Test CreateProject", func() {
//
//			rand.Seed(time.Now().Unix())
//			intn := rand.Intn(10000)
//
//			var preCode = "alanTest" + strconv.Itoa(intn)
//
//			startTime := types.NowTime()
//
//			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")
//
//			planEndTime := types.Time(endTime)
//
//			fmt.Println("时间", endTime)
//
//			var remark = "哈哈哈"
//
//			projectName := "alan测试项目" + strconv.Itoa(intn)
//
//			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}
//
//			project, err := CreateProject(projectvo.CreateProjectReqVo{
//				OrgId:  cacheUserInfo.OrgId,
//				UserId: cacheUserInfo.UserId,
//				Input:  input,
//			})
//
//			if err != nil {
//				fmt.Println("错误.......")
//			}
//
//			fmt.Printf("项目%+v", project)
//
//			convey.Convey("Test CreateIssue", func() {
//
//				var projectId = int64(1204)
//
//				rand.Seed(time.Now().Unix())
//				intn := rand.Intn(10000)
//				//projectObjectId
//				typeId := int64(intn)
//
//				currentUserId := cacheUserInfo.UserId
//				orgId := cacheUserInfo.OrgId
//
//				convey.Convey("createIssue insert Project proecss", func() {
//
//					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())
//
//					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
//						Id:                  processId,
//						OrgId:               orgId,
//						ProjectId:           projectId,
//						ProjectObjectTypeId: typeId,
//						ProcessId:           1120, //这个是一开始默认的不能写死
//						Creator:             currentUserId,
//						CreateTime:          time.Now(),
//						Updator:             currentUserId,
//						UpdateTime:          time.Now(),
//						Version:             1,
//					}
//
//					_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						err4 := mysql.TransInsert(tx, projectObjectTypeProcess)
//						if err4 != nil {
//							log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						}
//						return nil
//					})
//
//					convey.Convey("insertPriority.... ", func() {
//
//						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
//						//
//						//priority := bo.PriorityBo{
//						//	Id:    id,
//						//	OrgId: orgId,
//						//	Type:  consts.PriorityTypeIssue}
//						//
//						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						//	err4 := mysql.TransInsert(tx, &priority)
//						//	if err4 != nil {
//						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						//	}
//						//	return nil
//						//})
//						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
//							OrgID: orgId,
//							Type:  consts.PriorityTypeIssue})
//
//						convey.Convey("mock createiteration input", func() {
//
//							var iterationName = "alan创建的迭代" + strconv.Itoa(intn)
//							planStartTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
//							planEndTime, err := time.Parse(consts.AppTimeFormat, "2020-09-26 12:30:22")
//
//							iterationInput := vo.CreateIterationReq{
//								ProjectID:     projectId,
//								Name:          iterationName,
//								Owner:         currentUserId,
//								PlanStartTime: types.Time(planStartTime),
//								PlanEndTime:   types.Time(planEndTime)}
//
//							modelsVoid, err := CreateIteration(cacheUserInfo.OrgId, cacheUserInfo.UserId, iterationInput)
//
//							convey.So(modelsVoid, convey.ShouldNotBeNil)
//							convey.So(err, convey.ShouldBeNil)
//
//							var remark = "这是我创建的第一个issue"
//							var issueObjectId = int64(1)
//
//							children := make([]*vo.IssueChildren, 0)
//
//							firstChild := &vo.IssueChildren{
//								Title: "alan子issue",
//								// 负责人
//								OwnerID: []int64{currentUserId},
//								// 优先级
//								PriorityID: v.ID}
//
//							children = append(children, firstChild)
//							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//							if err != nil {
//								log.Error(err)
//								return
//							}
//							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: []int64{currentUserId}, Remark: &remark,
//								IssueObjectID: &issueObjectId, Children: children}
//
//							reqVo := projectvo.CreateIssueReqVo{
//								CreateIssue: input,
//								OrgId:       cacheUserInfo.OrgId,
//								UserId:      cacheUserInfo.UserId,
//							}
//
//							convey.Convey("createIssue done", func() {
//
//								issue, err := CreateIssue(reqVo)
//								convey.So(issue, convey.ShouldNotBeNil)
//								convey.So(err, convey.ShouldBeNil)
//
//								convey.Convey("mock updateIssue req", func() {
//
//									var hour = 10
//									randomStr := strconv.Itoa(intn)
//									var updateRemark = "这是我更新的第一个issue" + randomStr
//
//									participantIds := make([]int64, 1)
//									participantIds[0] = 1065
//									followerIds := make([]int64, 1)
//									followerIds[0] = 1070
//									updateFields := make([]string, 0)
//
//									updateFields = append(updateFields, "remark")
//									updateFields = append(updateFields, "title")
//									updateFields = append(updateFields, "planEndTime")
//									updateFields = append(updateFields, "planStartTime")
//									updateFields = append(updateFields, "planWorkHour")
//									updateFields = append(updateFields, "priorityId")
//									updateFields = append(updateFields, "iterationId")
//									updateFields = append(updateFields, "sourceId")
//									updateFields = append(updateFields, "issueObjectTypeId")
//									updateFields = append(updateFields, "ownerId")
//									updateFields = append(updateFields, "followerIds")
//									updateFields = append(updateFields, "participantIds")
//									updateFields = append(updateFields, "hour")
//
//									updateInput := vo.UpdateIssueReq{
//										ID:    issue.ID,
//										Title: &issue.Title,
//										//OwnerID:           []int64{issue.Owner},
//										PriorityID:        &issue.PriorityID,
//										PlanStartTime:     &issue.PlanStartTime,
//										PlanEndTime:       &issue.PlanEndTime,
//										PlanWorkHour:      &hour,
//										Remark:            &updateRemark,
//										IterationID:       &modelsVoid.ID,
//										SourceID:          &issue.SourceID,
//										IssueObjectTypeID: &issue.IssueObjectTypeID,
//										ParticipantIds:    participantIds,
//										FollowerIds:       followerIds,
//										UpdateFields:      updateFields}
//
//									req := projectvo.UpdateIssueReqVo{
//										Input:  updateInput,
//										UserId: cacheUserInfo.UserId,
//										OrgId:  cacheUserInfo.OrgId,
//									}
//
//									convey.Convey(" updateIssue done", func() {
//
//										updateIssueResp, err := UpdateIssue(req)
//										convey.So(updateIssueResp, convey.ShouldNotBeNil)
//										convey.So(err, convey.ShouldBeNil)
//									})
//								})
//							})
//						})
//					})
//				})
//			})
//		})
//	}))
//}

//func TestUpdateIssueForOwner(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		endTime := types.Time(time.Date(2021, 01, 05, 0, 0, 0, 0, time.Local))
//		userId := int64(1012)
//		orgId := int64(1004)
//		newOwnerId := int64(1007)
//		req := projectvo.UpdateIssueReqVo{
//			Input: vo.UpdateIssueReq{
//				ID:           1030,
//				OwnerID:      []int64{newOwnerId},
//				PlanEndTime:  &endTime,
//				UpdateFields: []string{"planEndTime"},
//			},
//			UserId: userId,
//			OrgId:  orgId,
//		}
//		updateIssueResp, err := UpdateIssue(req)
//		if err != nil {
//			t.Error(err)
//		}
//		log.Info(updateIssueResp)
//	}))
//}

//func TestUpdateCalendar(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		planStartTime := types.Time(time.Date(2021, 01, 04, 0, 0, 0, 0, time.Local))
//		planEndTime := types.Time(time.Date(2021, 01, 05, 0, 0, 0, 0, time.Local))
//		updateIssueBo := bo.IssueUpdateBo{
//			IssueBo: bo.IssueBo{
//				Id:            38062,
//				PlanStartTime: planStartTime,
//				PlanEndTime:   planEndTime,
//			},
//			NewIssueBo: bo.IssueBo{
//				//Owner:         1624,
//				PlanStartTime: planStartTime,
//				PlanEndTime:   planEndTime,
//			},
//		}
//		beforeFollowers := []int64{}
//		afterFollowers := []int64{}
//		domain.UpdateCalendarEvent(updateIssueBo, 38062, 1213, beforeFollowers, afterFollowers)
//	}))
//}

func TestGetDifMemberIds(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		before := []int64{111}
		after := []int64{}
		deleteIds, add := util.GetDifMemberIds(before, after)
		t.Log(deleteIds)
		t.Log(add)
		res := slice.SliceUniqueInt64(append(deleteIds, add...))
		t.Log(res)
	}))
}

////删除迭代 有子任务不可以删除的
//func TestDeleteIssueExistChildren(t *testing.T) {
//
//	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("Test CreateProject", func() {
//
//			rand.Seed(time.Now().Unix())
//			intn := rand.Intn(10000)
//
//			var preCode = "alanTest" + strconv.Itoa(intn)
//
//			startTime := types.NowTime()
//
//			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")
//
//			planEndTime := types.Time(endTime)
//
//			fmt.Println("时间", endTime)
//
//			var remark = "哈哈哈"
//
//			projectName := "alan测试项目" + strconv.Itoa(intn)
//
//			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}
//
//			project, err := CreateProject(projectvo.CreateProjectReqVo{
//				OrgId:  cacheUserInfo.OrgId,
//				UserId: cacheUserInfo.UserId,
//				Input:  input,
//			})
//
//			if err != nil {
//				fmt.Println("错误.......")
//			}
//
//			fmt.Printf("项目%+v", project)
//
//			convey.Convey("Test CreateIssue", func() {
//
//				var projectId = int64(1204)
//
//				rand.Seed(time.Now().Unix())
//				intn := rand.Intn(10000)
//				//projectObjectId
//				typeId := int64(intn)
//
//				currentUserId := cacheUserInfo.UserId
//				orgId := cacheUserInfo.OrgId
//
//				convey.Convey("createIssue insert Project proecss", func() {
//
//					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())
//
//					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
//						Id:                  processId,
//						OrgId:               orgId,
//						ProjectId:           projectId,
//						ProjectObjectTypeId: typeId,
//						ProcessId:           1120, //这个是一开始默认的不能写死
//						Creator:             currentUserId,
//						CreateTime:          time.Now(),
//						Updator:             currentUserId,
//						UpdateTime:          time.Now(),
//						Version:             1,
//					}
//
//					_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						err4 := mysql.TransInsert(tx, projectObjectTypeProcess)
//						if err4 != nil {
//							log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						}
//						return nil
//					})
//
//					convey.Convey("insertPriority.... ", func() {
//
//						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
//						//
//						//priority := bo.PriorityBo{
//						//	Id:    id,
//						//	OrgId: orgId,
//						//	Type:  consts.PriorityTypeIssue}
//						//
//						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						//	err4 := mysql.TransInsert(tx, &priority)
//						//	if err4 != nil {
//						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						//	}
//						//	return nil
//						//})
//						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
//							OrgID: orgId,
//							Type:  consts.PriorityTypeIssue})
//
//						convey.Convey("mock createIssue input", func() {
//
//							var remark = "来删除的issue含字"
//							var issueObjectId = int64(1)
//
//							children := make([]*vo.IssueChildren, 0)
//
//							firstChild := &vo.IssueChildren{
//								Title: "alan子issue",
//								// 负责人
//								OwnerID: []int64{currentUserId},
//								// 优先级
//								PriorityID: v.ID}
//
//							children = append(children, firstChild)
//
//							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//							if err != nil {
//								log.Error(err)
//								return
//							}
//
//							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, OwnerID: []int64{currentUserId}, Remark: &remark,
//								IssueObjectID: &issueObjectId, Children: children}
//
//							reqVo := projectvo.CreateIssueReqVo{
//								CreateIssue: input,
//								OrgId:       cacheUserInfo.OrgId,
//								UserId:      cacheUserInfo.UserId,
//							}
//
//							convey.Convey("createIssue done", func() {
//
//								issue, err := CreateIssue(reqVo)
//								convey.So(issue, convey.ShouldNotBeNil)
//								convey.So(err, convey.ShouldBeNil)
//
//								convey.Convey("mock deleteIssue req", func() {
//
//									deleteInput := vo.DeleteIssueReq{
//										ID: issue.ID}
//
//									req := projectvo.DeleteIssueReqVo{
//										Input:  deleteInput,
//										OrgId:  cacheUserInfo.OrgId,
//										UserId: cacheUserInfo.UserId,
//									}
//
//									convey.Convey(" DeleteIssue done", func() {
//
//										deleteIssueResp, err := DeleteIssue(req)
//										convey.So(deleteIssueResp, convey.ShouldBeNil)
//										convey.So(err, convey.ShouldNotBeNil)
//
//									})
//								})
//
//							})
//						})
//					})
//				})
//			})
//		})
//	}))
//}

////删除迭代 有子任务不可以删除的
//func TestDeleteIssueNotExistChildren(t *testing.T) {
//
//	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
//
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("Test CreateProject", func() {
//
//			rand.Seed(time.Now().Unix())
//			intn := rand.Intn(10000)
//
//			var preCode = "alanTest" + strconv.Itoa(intn)
//
//			startTime := types.NowTime()
//
//			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")
//
//			planEndTime := types.Time(endTime)
//
//			fmt.Println("时间", endTime)
//
//			var remark = "哈哈哈"
//
//			projectName := "alan测试项目" + strconv.Itoa(intn)
//
//			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}
//
//			project, err := CreateProject(projectvo.CreateProjectReqVo{
//				OrgId:  cacheUserInfo.OrgId,
//				UserId: cacheUserInfo.UserId,
//				Input:  input,
//			})
//
//			if err != nil {
//				fmt.Println("错误.......")
//			}
//
//			fmt.Printf("项目%+v", project)
//
//			convey.Convey("Test CreateIssue", func() {
//
//				var projectId = int64(1204)
//
//				rand.Seed(time.Now().Unix())
//				intn := rand.Intn(10000)
//				//projectObjectId
//				typeId := int64(intn)
//
//				currentUserId := cacheUserInfo.UserId
//				orgId := cacheUserInfo.OrgId
//
//				convey.Convey("createIssue insert Project proecss", func() {
//
//					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())
//
//					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
//						Id:                  processId,
//						OrgId:               orgId,
//						ProjectId:           projectId,
//						ProjectObjectTypeId: typeId,
//						ProcessId:           1120, //这个是一开始默认的不能写死
//						Creator:             currentUserId,
//						CreateTime:          time.Now(),
//						Updator:             currentUserId,
//						UpdateTime:          time.Now(),
//						Version:             1,
//					}
//
//					_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						err4 := mysql.TransInsert(tx, projectObjectTypeProcess)
//						if err4 != nil {
//							log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						}
//						return nil
//					})
//
//					convey.Convey("insertPriority.... ", func() {
//
//						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
//						//
//						//priority := bo.PriorityBo{
//						//	Id:    id,
//						//	OrgId: orgId,
//						//	Type:  consts.PriorityTypeIssue}
//						//
//						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						//	err4 := mysql.TransInsert(tx, &priority)
//						//	if err4 != nil {
//						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						//	}
//						//	return nil
//						//})
//
//						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
//							OrgID: orgId,
//							Type:  consts.PriorityTypeIssue})
//
//						convey.Convey("mock createIssue input", func() {
//
//							var remark = "来删除的issue"
//							var issueObjectId = int64(1)
//
//							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//							if err != nil {
//								log.Error(err)
//								return
//							}
//
//							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, OwnerID: []int64{currentUserId}, Remark: &remark,
//								IssueObjectID: &issueObjectId}
//
//							reqVo := projectvo.CreateIssueReqVo{
//								CreateIssue: input,
//								OrgId:       cacheUserInfo.OrgId,
//								UserId:      cacheUserInfo.UserId,
//							}
//
//							convey.Convey("createIssue done", func() {
//
//								issue, err := CreateIssue(reqVo)
//								convey.So(issue, convey.ShouldNotBeNil)
//								convey.So(err, convey.ShouldBeNil)
//
//								convey.Convey("mock deleteIssue req", func() {
//
//									deleteInput := vo.DeleteIssueReq{
//										ID: issue.ID}
//
//									req := projectvo.DeleteIssueReqVo{
//										Input:  deleteInput,
//										OrgId:  cacheUserInfo.OrgId,
//										UserId: cacheUserInfo.UserId,
//									}
//
//									convey.Convey(" DeleteIssue done", func() {
//
//										deleteIssueResp, err := DeleteIssue(req)
//										convey.So(deleteIssueResp, convey.ShouldNotBeNil)
//										convey.So(err, convey.ShouldBeNil)
//									})
//								})
//
//							})
//						})
//					})
//				})
//			})
//		})
//	}))
//}

////问题详情测试用例
//func TestIssueInfo(t *testing.T) {
//
//	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("Test CreateProject", func() {
//
//			rand.Seed(time.Now().Unix())
//			intn := rand.Intn(10000)
//
//			var preCode = "alanTest" + strconv.Itoa(intn)
//
//			startTime := types.NowTime()
//
//			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")
//
//			planEndTime := types.Time(endTime)
//
//			fmt.Println("时间", endTime)
//
//			var remark = "哈哈哈"
//
//			projectName := "alan测试项目" + strconv.Itoa(intn)
//
//			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}
//
//			project, err := CreateProject(projectvo.CreateProjectReqVo{
//				OrgId:  cacheUserInfo.OrgId,
//				UserId: cacheUserInfo.UserId,
//				Input:  input,
//			})
//
//			if err != nil {
//				fmt.Println("错误.......")
//			}
//
//			fmt.Printf("项目%+v", project)
//
//			convey.Convey("Test CreateIssue", func() {
//
//				var projectId = int64(1204)
//
//				rand.Seed(time.Now().Unix())
//				intn := rand.Intn(10000)
//				//projectObjectId
//				typeId := int64(intn)
//
//				currentUserId := cacheUserInfo.UserId
//				orgId := cacheUserInfo.OrgId
//
//				convey.Convey("createIssue insert Project proecss", func() {
//
//					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())
//
//					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
//						Id:                  processId,
//						OrgId:               orgId,
//						ProjectId:           projectId,
//						ProjectObjectTypeId: typeId,
//						ProcessId:           1120, //这个是一开始默认的不能写死
//						Creator:             currentUserId,
//						CreateTime:          time.Now(),
//						Updator:             currentUserId,
//						UpdateTime:          time.Now(),
//						Version:             1,
//					}
//
//					_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						err4 := mysql.TransInsert(tx, projectObjectTypeProcess)
//						if err4 != nil {
//							log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						}
//						return nil
//					})
//
//					convey.Convey("insertPriority.... ", func() {
//
//						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
//						//
//						//priority := bo.PriorityBo{
//						//	Id:    id,
//						//	OrgId: orgId,
//						//	Type:  consts.PriorityTypeIssue}
//						//
//						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						//	err4 := mysql.TransInsert(tx, &priority)
//						//	if err4 != nil {
//						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						//	}
//						//	return nil
//						//})
//						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
//							OrgID: orgId,
//							Type:  consts.PriorityTypeIssue})
//
//						convey.Convey("mock createIssue input", func() {
//
//							var remark = "来删除的issue"
//							var issueObjectId = int64(1)
//
//							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//							if err != nil {
//								log.Error(err)
//								return
//							}
//
//							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, OwnerID: []int64{currentUserId}, Remark: &remark,
//								IssueObjectID: &issueObjectId}
//
//							reqVo := projectvo.CreateIssueReqVo{
//								CreateIssue: input,
//								OrgId:       cacheUserInfo.OrgId,
//								UserId:      cacheUserInfo.UserId,
//							}
//
//							convey.Convey("createIssue done", func() {
//
//								issue, err := CreateIssue(reqVo)
//								convey.So(issue, convey.ShouldNotBeNil)
//								convey.So(err, convey.ShouldBeNil)
//
//								convey.Convey("Test IssueInfo", t, func() {
//									convey.Convey("IssueInfo mock req", func() {
//
//										resp, err := IssueInfo(cacheUserInfo.OrgId, cacheUserInfo.UserId, issue.ID, 1)
//										convey.So(resp, convey.ShouldNotBeNil)
//
//										convey.So(err, convey.ShouldBeNil)
//									})
//								})
//							})
//						})
//					})
//				})
//			})
//		})
//	}))
//}

//func TestGetIssueRestInfos(t *testing.T) {
//
//	convey.Convey("Test GetIssueRestInfos", t, test.StartUp(func(ctx context.Context) {
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息" + cacheUserInfoJson)
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
//
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("GetIssueRestInfos mock req", func() {
//			input := &vo.IssueRestInfoReq{}
//			page := 1
//			size := 20
//			resp, err := GetIssueRestInfos(cacheUserInfo.OrgId, "", page, size, input)
//
//			convey.So(resp, convey.ShouldNotBeNil)
//
//			convey.So(err, convey.ShouldBeNil)
//		})
//	}))
//}

//func TestUpdateIssueStatus(t *testing.T) {
//
//	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息" + strs.ObjectToString(cacheUserInfoJson))
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
//
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("Test CreateProject", func() {
//
//			rand.Seed(time.Now().Unix())
//			intn := rand.Intn(10000)
//
//			var preCode = "alanTest" + strconv.Itoa(intn)
//
//			startTime := types.NowTime()
//
//			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")
//
//			planEndTime := types.Time(endTime)
//
//			fmt.Println("时间", endTime)
//
//			var remark = "哈哈哈"
//
//			projectName := "alan测试项目" + strconv.Itoa(intn)
//
//			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}
//
//			project, err := CreateProject(projectvo.CreateProjectReqVo{
//				OrgId:  cacheUserInfo.OrgId,
//				UserId: cacheUserInfo.UserId,
//				Input:  input,
//			})
//
//			if err != nil {
//				fmt.Println("错误.......")
//			}
//
//			fmt.Printf("项目%+v", project)
//
//			convey.Convey("Test CreateIssue", func() {
//
//				var projectId = int64(1204)
//
//				rand.Seed(time.Now().Unix())
//				intn := rand.Intn(10000)
//				//projectObjectId
//				typeId := int64(intn)
//
//				currentUserId := cacheUserInfo.UserId
//				orgId := cacheUserInfo.OrgId
//
//				convey.Convey("createIssue insert Project proecss", func() {
//
//					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())
//
//					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
//						Id:                  processId,
//						OrgId:               orgId,
//						ProjectId:           projectId,
//						TableId: typeId,
//						ProcessId:           1120, //这个是一开始默认的不能写死
//						Creator:             currentUserId,
//						CreateTime:          time.Now(),
//						Updator:             currentUserId,
//						UpdateTime:          time.Now(),
//						Version:             1,
//					}
//
//					_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
//						err4 := mysql.TransInsert(tx, projectObjectTypeProcess)
//						if err4 != nil {
//							log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
//							return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
//						}
//						return nil
//					})
//
//					convey.Convey("insertPriority.... ", func() {
//						v, _ := CreatePriority(cacheUserInfo.UserId, vo.CreatePriorityReq{
//							OrgID: orgId,
//							Type:  consts.PriorityTypeIssue})
//
//						convey.Convey("mock createIssue input", func() {
//
//							var remark = "这是我创建的第一个issue"
//							var issueObjectId = int64(1)
//
//							cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
//							if err != nil {
//								log.Error(err)
//								return
//							}
//
//							input := vo.CreateIssueReq{ProjectID: projectId, Title: "alan创建的issue", PriorityID: v.ID, TypeID: &typeId, OwnerID: []int64{currentUserId}, Remark: &remark,
//								IssueObjectID: &issueObjectId}
//
//							reqVo := projectvo.CreateIssueReqVo{
//								CreateIssue: input,
//								OrgId:       cacheUserInfo.OrgId,
//								UserId:      cacheUserInfo.UserId,
//							}
//
//							//convey.Convey("createIssue done", func() {
//							//
//							//	issue, err := CreateIssue(reqVo)
//							//	convey.So(issue, convey.ShouldNotBeNil)
//							//	convey.So(err, convey.ShouldBeNil)
//							//
//							//	req := processvo.GetProcessStatusListByCategoryReqVo{
//							//		OrgId:    orgId,
//							//		Category: consts.ProcessStatusCategoryIssue}
//							//
//							//	//category := processfacade.GetProcessStatusListByCategory(req)
//							//	//
//							//	//bos := category.CacheProcessStatusBoList
//							//
//							//	//processStatusList := &[]po.PpmPrsProcessStatus{}
//							//	//err = mysql.SelectAllByCond(processStatus.TableName(), db.Cond{
//							//	//	consts.TcOrgId:    orgId,
//							//	//	consts.TcIsDelete: consts.AppIsNoDelete,
//							//	//	consts.TcStatus:   consts.AppStatusEnable,
//							//	//	consts.TcCategory: consts.ProcessStatusCategoryIssue}, processStatusList)
//							//	//if err != nil {
//							//	//	return
//							//	//}
//							//
//							//	var nextStatusId int64
//							//
//							//	issueBo, err := domain.GetIssueBo(orgId, issue.ID)
//							//
//							//	for _, status := range bos {
//							//		if issueBo.Status != status.StatusId {
//							//			nextStatusId = status.StatusId
//							//			break
//							//		}
//							//	}
//							//
//							//	convey.Convey("UpdateIssueStatus mock req", func() {
//							//
//							//		input := vo.UpdateIssueStatusReq{
//							//			ID:           issue.ID,
//							//			NextStatusID: &nextStatusId}
//							//
//							//		resp, err := UpdateIssueStatus(projectvo.UpdateIssueStatusReqVo{
//							//			OrgId:  cacheUserInfo.OrgId,
//							//			UserId: cacheUserInfo.UserId,
//							//			Input:  input,
//							//		})
//							//
//							//		convey.So(resp, convey.ShouldNotBeNil)
//							//
//							//		convey.So(err, convey.ShouldBeNil)
//							//	})
//							//})
//						})
//					})
//				})
//
//			})
//		})
//	}))
//}

//func TestIssueReport(t *testing.T) {
//
//	convey.Convey("Test IssueReport", t, test.StartUp(func(ctx context.Context) {
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息" + cacheUserInfoJson)
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
//
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("IssueReport mock req", func() {
//
//			resp, err := IssueReport(cacheUserInfo.OrgId, cacheUserInfo.UserId, consts.DailyReport)
//
//			convey.So(resp, convey.ShouldNotBeNil)
//
//			convey.So(err, convey.ShouldBeNil)
//
//		})
//	}))
//}

//func TestIssueReportDetail(t *testing.T) {
//	convey.Convey("Test IssueReportDetail", t, test.StartUp(func(ctx context.Context) {
//
//		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)
//
//		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)
//
//		log.Info("缓存用户信息" + cacheUserInfoJson)
//
//		if cacheUserInfo == nil {
//			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}
//
//		}
//
//		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)
//
//		convey.Convey("IssueReportDetail mock req", func() {
//
//			resp, _ := IssueReport(cacheUserInfo.OrgId, cacheUserInfo.UserId, consts.DailyReport)
//			//str := resp.ShareID  因为每次状态要变
//			str := "201"
//			aesStr, _ := encrypt.AesEncrypt(str)
//
//			resp, err2 := IssueReportDetail(aesStr)
//
//			convey.So(resp, convey.ShouldNotBeNil)
//
//			convey.So(err2, convey.ShouldBeNil)
//		})
//	}))
//}

//func TestUpdateIssueProjectObjectType(t *testing.T) {
//	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
//		t.Log(MoveIssue(1574, 24362, vo.UpdateIssueProjectObjectTypeReq{ID: 76474, ProjectObjectTypeID: 0}, 0))
//	}))
//}

//func TestCreateIssue11(t *testing.T) {
//	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
//		finishId := int64(26)
//		startId := int64(7)
//		child := &vo.IssueChildren{
//			Title:    "111",
//			OwnerID:  nil,
//			StatusID: &startId,
//		}
//		children := []*vo.IssueChildren{}
//		children = append(children, child)
//		t.Log(CreateIssue(projectvo.CreateIssueReqVo{
//			CreateIssue: vo.CreateIssueReq{
//				ProjectID:  1384,
//				Title:      "111",
//				PriorityID: 1400,
//				//TypeID:     nil,
//				OwnerID:  []int64{1293},
//				StatusID: &finishId,
//				Children: children,
//			},
//			UserId:        1293,
//			OrgId:         1113,
//			SourceChannel: "",
//		}))
//	}))
//}

func TestCreateIssue122(t *testing.T) {
	arr := []int{6, 33, 4, 5}
	sort.Ints(arr)
	t.Log(arr)
}

//func TestCopyIssue(t *testing.T) {
//	convey.Convey("Test CopyIssue", t, test.StartUp(func(ctx context.Context) {
//		t.Log(CopyIssue(2571, 1368436, vo.CopyIssueReq{
//			OldIssueID:  10093399,
//			Title:       "测试2",
//			ProjectID:   13757,
//			IterationID: 0,
//			StatusID:    0,
//			TableID:     "1572900881987276800",
//			ChildrenIds: []int64{10093400},
//			ChooseField: []string{"ownerId", "planStartTime", "planEndTime", "followerIds", "remark", "issueStatus", "auditorIds", "_field_priority", "relating", "children"},
//		}))
//	}))
//}

func TestUpdateIssueProjectObjectTypeBatch(t *testing.T) {
	//convey.Convey("Test CopyIssue", t, test.StartUp(func(ctx context.Context) {
	//	//statusId := int64(7)
	//	res, err := MoveIssueBatch(1629, 1042, vo.UpdateIssueProjectObjectTypeBatchReq{
	//		Ids:          []int64{76016, 76017, 76018, 76019},
	//		OldProjectID: 9788,
	//		//ProjectObjectTypeID: 14393,
	//		//StatusID:            &statusId,
	//	})
	//	t.Log(json.ToJsonIgnoreError(res))
	//	t.Log(err)
	//}))

}

func TestIssueInfo2(t *testing.T) {
	//convey.Convey("TestIssueInfo2", t, test.StartUp(func(ctx context.Context) {
	//	resp, err := IssueInfo(2373, 29611, 10084446, 2)
	//	if err != nil {
	//		t.Error(err)
	//		return
	//	}
	//	t.Log(json.ToJsonIgnoreError(resp))
	//}))
}

//func TestUrgeIssue(t *testing.T) {
//	convey.Convey("TestUpdateIssue", t, test.StartUp(func(ctx context.Context) {
//		orgId := int64(1242)
//		userId := int64(3306)
//		text := "看看进度啥样了。"
//		isOk, err := UrgeIssue(orgId, userId, "fs", &vo.UrgeIssueReq{
//			IssueID:            39304,
//			IsNeedAtIssueOwner: false,
//			UrgeText:           &text,
//		})
//		if err != nil {
//			t.Error(err)
//		}
//		t.Log(isOk)
//	}))
//}

//func TestUpdateIssue2(t *testing.T) {
//	convey.Convey("TestUpdateIssue", t, test.StartUp(func(ctx context.Context) {
//		t.Log(UpdateIssue(projectvo.UpdateIssueReqVo{
//			Input: vo.UpdateIssueReq{
//				ID:              53994,
//				FollowerIds:     []int64{1885},
//				FollowerDeptIds: []int64{1885},
//				UpdateFields:    []string{"followerIds"},
//			},
//			UserId:        1883,
//			OrgId:         1242,
//			SourceChannel: "",
//		}))
//	}))
//}

//func TestUpdateIssue3(t *testing.T) {
//	convey.Convey("TestUpdateIssue", t, test.StartUp(func(ctx context.Context) {
//		resp, err := UpdateIssue(projectvo.UpdateIssueReqVo{
//			Input: vo.UpdateIssueReq{
//				ID: 74649,
//				//UpdateColumns:[]string{"followerIds"},
//				//FollowerIds: []int64{3306},
//				LessUpdateIssueReq: map[string]interface{}{
//					"id":                   74649,
//					"_field_1635495062038": "aa",
//				},
//			},
//			UserId:        1042,
//			OrgId:         1574,
//			SourceChannel: "fs",
//		})
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		t.Log(json.ToJsonIgnoreError(resp))
//	}))
//}

//func TestGetIssueRestInfos2(t *testing.T) {
//	convey.Convey("TestUpdateIssue", t, test.StartUp(func(ctx context.Context) {
//		parentId := int64(54110)
//		res, err := GetIssueRestInfos(1045, "", 1, 10, &vo.IssueRestInfoReq{ParentID: &parentId})
//		t.Log(json.ToJsonIgnoreError(res))
//		t.Log(err)
//	}))
//}

func TestDeleteIssueBatch1(t *testing.T) {
	convey.Convey("TestUpdateIssue", t, test.StartUp(func(ctx context.Context) {
		deleteIds := []int64{
			84707,
		}
		resp, err := DeleteIssueBatch(projectvo.DeleteIssueBatchReqVo{
			Input: vo.DeleteIssueBatchReq{
				ProjectID: 10809,
				Ids:       deleteIds,
			},
			UserId:        24423,
			OrgId:         1574,
			SourceChannel: "",
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestCheckHasPermission1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		err := domain.AuthProject(1574, 24362, 10584, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify)
		if err != nil {
			t.Log(err)
			return
		}
		t.Log("ok, end")
	}))
}

func TestApplyProjectTemplateInner(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		inputJson := "{\"projectDetail\":{\"isEnableWorkHours\":2,\"isSyncOutCalendar\":2,\"notice\":\"\"},\"projectInfo\":{\"isFiling\":2,\"name\":\"xxx\",\"owner\":1152,\"projectTypeId\":1,\"publicStatus\":1,\"publishStatus\":1,\"remark\":\"\",\"resourceId\":6936,\"status\":2},\"projectIssueTemplates\":[{\"id\":1484428161637736450,\"sourceId\":0,\"code\":\"$12156-2\",\"ownerChangeTime\":\"2022-01-21 15:30:27\",\"followerIds\":[\"U_1152\"],\"planWorkHour\":0,\"remark\":\"\",\"title\":\"2\",\"delFlag\":2,\"ownerId\":[\"U_1152\"],\"orgId\":1886,\"priorityId\":0,\"path\":\"0,104826,\",\"issueStatus\":26,\"collaborators\":{\"followerIds\":[\"U_1152\"],\"auditorIds\":[\"U_1042\"],\"ownerId\":[\"U_1152\"]},\"startTime\":\"1970-01-01 00:00:00\",\"auditorIds\":[\"U_1042\"],\"moduleId\":0,\"iterationId\":0,\"propertyId\":0,\"order\":16580607,\"owner\":1152,\"creator\":\"1152\",\"issueId\":104827,\"recycleFlag\":2,\"planStartTime\":\"1970-01-01 00:00:00\",\"isDelete\":2,\"updateTime\":\"2022-01-21 16:02:05\",\"sort\":6869837445,\"projectObjectTypeId\":19700,\"version\":1,\"parentId\":104826,\"versionId\":0,\"appIds\":[\"1484428121678643201\"],\"createTime\":\"2022-01-21 15:30:27\",\"isFiling\":0,\"issueObjectTypeId\":0,\"updator\":\"1042\",\"auditStatus\":1,\"planEndTime\":\"1970-01-01 00:00:00\",\"endTime\":\"2022-01-21 16:02:04\",\"projectId\":12156,\"status\":1},{\"id\":1484428184219869186,\"sourceId\":0,\"code\":\"$12156-3\",\"ownerChangeTime\":\"2022-01-21 15:30:32\",\"followerIds\":[\"U_1152\"],\"planWorkHour\":0,\"remark\":\"\",\"title\":\"3\",\"delFlag\":2,\"ownerId\":[\"U_1152\"],\"orgId\":1886,\"priorityId\":0,\"path\":\"0,104826,104827,\",\"issueStatus\":26,\"collaborators\":{\"followerIds\":[\"U_1152\"],\"auditorIds\":[\"U_1042\"],\"ownerId\":[\"U_1152\"]},\"startTime\":\"1970-01-01 00:00:00\",\"auditorIds\":[\"U_1042\"],\"moduleId\":0,\"iterationId\":0,\"propertyId\":0,\"order\":16646143,\"owner\":1152,\"creator\":\"1152\",\"issueId\":104828,\"recycleFlag\":2,\"planStartTime\":\"1970-01-01 00:00:00\",\"isDelete\":2,\"updateTime\":\"2022-01-21 16:02:46\",\"sort\":6869902980,\"projectObjectTypeId\":19700,\"version\":1,\"parentId\":104827,\"versionId\":0,\"appIds\":[\"1484428121678643201\"],\"createTime\":\"2022-01-21 15:30:32\",\"isFiling\":0,\"issueObjectTypeId\":0,\"updator\":\"1042\",\"auditStatus\":1,\"planEndTime\":\"1970-01-01 00:00:00\",\"endTime\":\"2022-01-21 16:02:45\",\"projectId\":12156,\"status\":1},{\"id\":1484428138682310657,\"sourceId\":0,\"code\":\"$12156-1\",\"ownerChangeTime\":\"2022-01-21 15:30:21\",\"followerIds\":[\"U_1152\"],\"planWorkHour\":0,\"remark\":\"\",\"title\":\"1\",\"delFlag\":2,\"ownerId\":[\"U_1152\"],\"orgId\":1886,\"priorityId\":0,\"path\":\"0,\",\"issueStatus\":26,\"collaborators\":{\"followerIds\":[\"U_1152\"],\"auditorIds\":[\"U_1042\"],\"ownerId\":[\"U_1152\"]},\"startTime\":\"0001-01-01 00:00:00\",\"auditorIds\":[\"U_1042\"],\"moduleId\":0,\"iterationId\":0,\"propertyId\":0,\"order\":16515071,\"_field_priority\":6600,\"owner\":1152,\"creator\":\"1152\",\"issueId\":104826,\"recycleFlag\":2,\"planStartTime\":\"1970-01-01 00:00:00\",\"isDelete\":2,\"updateTime\":\"2022-01-21 15:32:51\",\"sort\":6869771910,\"projectObjectTypeId\":19700,\"version\":1,\"parentId\":0,\"versionId\":0,\"appIds\":[\"1484428121678643201\"],\"createTime\":\"2022-01-21 15:30:21\",\"isFiling\":0,\"issueObjectTypeId\":0,\"updator\":\"1042\",\"auditStatus\":3,\"planEndTime\":\"1970-01-01 00:00:00\",\"endTime\":\"2022-01-21 15:30:58\",\"projectId\":12156,\"status\":1}],\"projectObjectTypeTemplates\":[{\"id\":19700,\"langCode\":\"Project.ObjectType.Task\",\"name\":\"任务\",\"objectType\":2,\"processId\":5,\"sort\":4}],\"templateFlag\":1}"
		input := projectvo.GetProjectTemplateData{}
		t.Log(json.FromJson(inputJson, &input))
		t.Log(ApplyProjectTemplateInner(&projectvo.ApplyProjectTemplateReq{
			OrgId:  1886,
			UserId: 1152,
			AppId:  1,
			Input:  input,
		}))
	}))
}

//func TestCopyIssueBatch1(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		res, err := CopyIssueBatch(2501, 24370, vo.CopyIssueBatchReq{
//			//OldProjectID: 13296,
//			OldIssueIds: []int64{147692, 147691},
//			ProjectID:   13333,
//			//IterationID:  2331,
//			//StatusID:    0,
//			ChooseField: []string{"children", "followerIds", "issueDetail", "priorityId", "tags"},
//			TableID:     "1520243956003115009",
//		}, false)
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		t.Log("--end--", res)
//	}))
//}

// 制表符补齐规则：https://blog.csdn.net/u012142460/article/details/94733824
func TestStrLen1(t *testing.T) {
	s1 := "标题1"
	s2 := "标题12"
	s3 := "标12"
	s4 := "…"
	msgText := fmt.Sprintf("s1 len: %d, s2 len: %d, s3 len: %d, s4 len: %d", len(s1), len(s2), len(s3), len(s4))
	t.Log(msgText)
}

func TestPaddingTabLen1(t *testing.T) {
	s1 := "标题1"
	tabStr1 := domain.GetTabStrForColumnName(s1)
	t.Log(tabStr1)
}

//func TestGetTableHelperInfo1(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		helperInfo, err := domain.GetSomeInfoForCreateIssue(2373, 1552488242861838336)
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		t.Log(json.ToJsonIgnoreError(helperInfo))
//	}))
//}
