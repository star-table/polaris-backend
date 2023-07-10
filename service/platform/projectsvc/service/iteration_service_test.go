package service

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestCreateIteration(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("Test CreateProject", func() {

			rand.Seed(time.Now().Unix())
			intn := rand.Intn(10000)

			var preCode = "alanTest" + strconv.Itoa(intn)

			startTime := types.NowTime()

			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")

			planEndTime := types.Time(endTime)

			fmt.Println("时间", endTime)

			var remark = "哈哈哈"

			projectName := "alan测试项目" + strconv.Itoa(intn)

			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}

			project, err := CreateProject(projectvo.CreateProjectReqVo{
				OrgId:         cacheUserInfo.OrgId,
				UserId:        cacheUserInfo.UserId,
				Input:         input,
				SourceChannel: "",
			})

			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("Test CreateIteration ", func() {
				//
				//testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())
				//
				//fmt.Println("unitteset Mysql配置json:", testMysqlJson)
				//
				//testRedisJson, _ := json.ToJson(config.GetRedisConfig())
				//
				//fmt.Println("unitteset redis配置json:", testRedisJson)
				//
				////添加token操作
				//ginCtx := gin.Context{}
				//ginCtx.Set(consts.AppHeaderTokenName, "abc")
				////获得一个顶级上下文
				//ctx := context.Background()
				////返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
				//ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)

				//var testConst =consts.ProcessStatusCategoryIteration

				var projectId = project.ID

				//projectObj
				//typeId := int64(intn)

				cacheUserInfo, errCur := orgfacade.GetCurrentUserRelaxed(ctx)

				fmt.Println(errCur)

				currentUserId := cacheUserInfo.UserId
				orgId := cacheUserInfo.OrgId

				convey.Convey("createIssue insert Project proecss ", func() {
					convey.Convey("insertPriority....  ", func() {
						//PpmPrsPriority
						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}

						//_, _ = CreatePriority(cacheUserInfo.OrgId, vo.CreatePriorityReq{
						//	OrgID: orgId,
						//	Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIteration input ", func() {

							var iterationName = "alan创建的迭代" + strconv.Itoa(intn)
							planStartTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
							planEndTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

							iterationInput := vo.CreateIterationReq{
								ProjectID:     projectId,
								Name:          iterationName,
								Owner:         currentUserId,
								PlanStartTime: types.Time(planStartTime),
								PlanEndTime:   types.Time(planEndTime)}

							modelsVoid, err := CreateIteration(orgId, currentUserId, iterationInput)

							convey.So(modelsVoid, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)
						})
					})
				})
			})

		})
	}))
}

func TestUpdateIteration(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("Test CreateProject", func() {

			rand.Seed(time.Now().Unix())
			intn := rand.Intn(10000)

			var preCode = "alanTest" + strconv.Itoa(intn)

			startTime := types.NowTime()

			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")

			planEndTime := types.Time(endTime)

			fmt.Println("时间", endTime)

			var remark = "哈哈哈"

			projectName := "alan测试项目" + strconv.Itoa(intn)

			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}

			project, err := CreateProject(projectvo.CreateProjectReqVo{
				OrgId:  cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input:  input,
			})

			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("Test CreateIteration ", func() {
				//
				//testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())
				//
				//fmt.Println("unitteset Mysql配置json:", testMysqlJson)
				//
				//testRedisJson, _ := json.ToJson(config.GetRedisConfig())
				//
				//fmt.Println("unitteset redis配置json:", testRedisJson)
				//
				////添加token操作
				//ginCtx := gin.Context{}
				//ginCtx.Set(consts.AppHeaderTokenName, "abc")
				////获得一个顶级上下文
				//ctx := context.Background()
				////返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
				//ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)

				//var testConst =consts.ProcessStatusCategoryIteration

				var projectId = project.ID

				//projectObj
				//typeId := int64(intn)

				cacheUserInfo, errCur := orgfacade.GetCurrentUserRelaxed(ctx)

				fmt.Println(errCur)

				currentUserId := cacheUserInfo.UserId
				orgId := cacheUserInfo.OrgId

				convey.Convey("createIssue insert Project proecss ", func() {

					convey.Convey("insertPriority....  ", func() {
						//PpmPrsPriority
						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}
						//
						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						//	err4 := mysql.TransInsert(tx, &priority)
						//	if err4 != nil {
						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						//	}
						//	return nil
						//})

						//_, _ = CreatePriority(cacheUserInfo.OrgId, vo.CreatePriorityReq{
						//	OrgID: orgId,
						//	Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIteration input ", func() {

							var iterationName = "alan创建的迭代" + strconv.Itoa(intn)
							planStartTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
							planEndTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

							iterationInput := vo.CreateIterationReq{
								ProjectID:     projectId,
								Name:          iterationName,
								Owner:         currentUserId,
								PlanStartTime: types.Time(planStartTime),
								PlanEndTime:   types.Time(planEndTime)}

							modelsVoid, err := CreateIteration(orgId, currentUserId, iterationInput)

							convey.So(modelsVoid, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)

							convey.Convey("update iteration", func() {

								itoa := strconv.Itoa(intn)
								updateName := "更改issue" + itoa

								updateFields := make([]string, 0)

								updateFields = append(updateFields, "name")
								updateFields = append(updateFields, "owner")
								updateFields = append(updateFields, "planStartTime")
								updateFields = append(updateFields, "planEndTime")

								pointerPlanStartTime := types.Time(planStartTime)
								pointerPlanEndTime := types.Time(planEndTime)

								input := vo.UpdateIterationReq{
									// 主键
									ID: modelsVoid.ID,
									// 名称
									Name: &updateName,
									// 负责人
									Owner: &currentUserId,
									// 计划开始时间
									PlanStartTime: &pointerPlanStartTime,
									// 计划结束时间
									PlanEndTime: &pointerPlanEndTime,
									// 变动的字段列表
									UpdateFields: updateFields}

								updateVoid, err := UpdateIteration(orgId, currentUserId, input)
								convey.So(updateVoid, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)
							})
						})
					})
				})
			})
		})
	}))
}

func TestDeleteIteration(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("Test DeleteIteration CreateProject", func() {

			rand.Seed(time.Now().Unix())
			intn := rand.Intn(10000)

			var preCode = "alanTest" + strconv.Itoa(intn)

			startTime := types.NowTime()

			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")

			planEndTime := types.Time(endTime)

			fmt.Println("时间", endTime)

			var remark = "哈哈哈"

			projectName := "alan测试项目" + strconv.Itoa(intn)

			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}

			project, err := CreateProject(projectvo.CreateProjectReqVo{
				OrgId:  cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input:  input,
			})

			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("Test CreateIteration ", func() {
				//
				//testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())
				//
				//fmt.Println("unitteset Mysql配置json:", testMysqlJson)
				//
				//testRedisJson, _ := json.ToJson(config.GetRedisConfig())
				//
				//fmt.Println("unitteset redis配置json:", testRedisJson)
				//
				////添加token操作
				//ginCtx := gin.Context{}
				//ginCtx.Set(consts.AppHeaderTokenName, "abc")
				////获得一个顶级上下文
				//ctx := context.Background()
				////返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
				//ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)

				//var testConst =consts.ProcessStatusCategoryIteration

				var projectId = project.ID

				//projectObj
				//typeId := int64(intn)

				cacheUserInfo, errCur := orgfacade.GetCurrentUserRelaxed(ctx)

				fmt.Println(errCur)

				currentUserId := cacheUserInfo.UserId
				orgId := cacheUserInfo.OrgId

				convey.Convey("createIssue insert Project proecss ", func() {
					convey.Convey("insertPriority....  ", func() {
						//PpmPrsPriority
						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}
						//
						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						//	err4 := mysql.TransInsert(tx, &priority)
						//	if err4 != nil {
						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						//	}
						//	return nil
						//})

						//_, _ = CreatePriority(cacheUserInfo.OrgId, vo.CreatePriorityReq{
						//	OrgID: orgId,
						//	Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIteration input ", func() {

							var iterationName = "alan创建的迭代" + strconv.Itoa(intn)
							planStartTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
							planEndTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

							iterationInput := vo.CreateIterationReq{
								ProjectID:     projectId,
								Name:          iterationName,
								Owner:         currentUserId,
								PlanStartTime: types.Time(planStartTime),
								PlanEndTime:   types.Time(planEndTime)}

							modelsVoid, err := CreateIteration(orgId, currentUserId, iterationInput)

							convey.So(modelsVoid, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)

							convey.Convey("daleteIssue done", func() {

								input := vo.DeleteIterationReq{
									ID: modelsVoid.ID}

								deleteVoid, err := DeleteIteration(orgId, currentUserId, input)

								convey.So(deleteVoid, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)
							})
						})
					})
				})
			})
		})
	}))
}

func TestIterationInfo(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("Test CreateProject", func() {

			rand.Seed(time.Now().Unix())
			intn := rand.Intn(10000)

			var preCode = "alanTest" + strconv.Itoa(intn)

			startTime := types.NowTime()

			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")

			planEndTime := types.Time(endTime)

			fmt.Println("时间", endTime)

			var remark = "哈哈哈"

			projectName := "alan测试项目" + strconv.Itoa(intn)

			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}

			project, err := CreateProject(projectvo.CreateProjectReqVo{
				OrgId:  cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input:  input,
			})

			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("Test CreateIteration ", func() {
				//
				//testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())
				//
				//fmt.Println("unitteset Mysql配置json:", testMysqlJson)
				//
				//testRedisJson, _ := json.ToJson(config.GetRedisConfig())
				//
				//fmt.Println("unitteset redis配置json:", testRedisJson)
				//
				////添加token操作
				//ginCtx := gin.Context{}
				//ginCtx.Set(consts.AppHeaderTokenName, "abc")
				////获得一个顶级上下文
				//ctx := context.Background()
				////返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
				//ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)

				//var testConst =consts.ProcessStatusCategoryIteration

				var projectId = project.ID

				//projectObj
				//typeId := int64(intn)

				cacheUserInfo, errCur := orgfacade.GetCurrentUserRelaxed(ctx)

				fmt.Println(errCur)

				currentUserId := cacheUserInfo.UserId
				orgId := cacheUserInfo.OrgId

				convey.Convey("createIssue insert Project proecss ", func() {
					convey.Convey("insertPriority....  ", func() {
						//PpmPrsPriority
						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}
						//
						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						//	err4 := mysql.TransInsert(tx, &priority)
						//	if err4 != nil {
						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						//	}
						//	return nil
						//})
						//_, _ = CreatePriority(cacheUserInfo.OrgId, vo.CreatePriorityReq{
						//	OrgID: orgId,
						//	Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIteration input ", func() {

							var iterationName = "alan创建的迭代" + strconv.Itoa(intn)
							planStartTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
							planEndTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

							iterationInput := vo.CreateIterationReq{
								ProjectID:     projectId,
								Name:          iterationName,
								Owner:         currentUserId,
								PlanStartTime: types.Time(planStartTime),
								PlanEndTime:   types.Time(planEndTime)}

							modelsVoid, err := CreateIteration(orgId, currentUserId, iterationInput)

							convey.So(modelsVoid, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)

							convey.Convey("IterationInfo done", func() {

								input := vo.IterationInfoReq{
									ID: modelsVoid.ID}

								resp, err := IterationInfo(orgId, input)

								convey.So(resp, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)
							})
						})
					})
				})
			})

		})
	}))
}

func TestIterationStatusTypeStat(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("Test CreateProject", func() {

			rand.Seed(time.Now().Unix())
			intn := rand.Intn(10000)

			var preCode = "alanTest" + strconv.Itoa(intn)

			startTime := types.NowTime()

			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")

			planEndTime := types.Time(endTime)

			fmt.Println("时间", endTime)

			var remark = "哈哈哈"

			projectName := "alan测试项目" + strconv.Itoa(intn)

			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}

			project, err := CreateProject(projectvo.CreateProjectReqVo{
				OrgId:  cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input:  input,
			})

			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("Test CreateIteration ", func() {
				//
				//testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())
				//
				//fmt.Println("unitteset Mysql配置json:", testMysqlJson)
				//
				//testRedisJson, _ := json.ToJson(config.GetRedisConfig())
				//
				//fmt.Println("unitteset redis配置json:", testRedisJson)
				//
				////添加token操作
				//ginCtx := gin.Context{}
				//ginCtx.Set(consts.AppHeaderTokenName, "abc")
				////获得一个顶级上下文
				//ctx := context.Background()
				////返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
				//ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)

				//var testConst =consts.ProcessStatusCategoryIteration

				var projectId = project.ID

				//projectObj
				//typeId := int64(intn)

				cacheUserInfo, errCur := orgfacade.GetCurrentUserRelaxed(ctx)

				fmt.Println(errCur)

				currentUserId := cacheUserInfo.UserId
				orgId := cacheUserInfo.OrgId

				convey.Convey("createIssue insert Project proecss ", func() {
					convey.Convey("insertPriority....  ", func() {
						//PpmPrsPriority
						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}
						//
						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						//	err4 := mysql.TransInsert(tx, &priority)
						//	if err4 != nil {
						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						//	}
						//	return nil
						//})

						//_, _ = CreatePriority(cacheUserInfo.OrgId, vo.CreatePriorityReq{
						//	OrgID: orgId,
						//	Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIteration input ", func() {

							var iterationName = "alan创建的迭代" + strconv.Itoa(intn)
							planStartTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
							planEndTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

							iterationInput := vo.CreateIterationReq{
								ProjectID:     projectId,
								Name:          iterationName,
								Owner:         currentUserId,
								PlanStartTime: types.Time(planStartTime),
								PlanEndTime:   types.Time(planEndTime)}

							modelsVoid, err := CreateIteration(orgId, currentUserId, iterationInput)

							convey.So(modelsVoid, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)

							convey.Convey("Test IterationStatusTypeStat", func() {

								input := vo.IterationStatusTypeStatReq{
									ProjectID: &projectId}

								resp, err := IterationStatusTypeStat(orgId, &input)

								convey.So(resp, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)

							})
						})
					})
				})
			})

		})
	}))
}

//func TestIterationIssueRelate(t *testing.T) {
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
//			convey.Convey("Test CreateIteration ", func() {
//				//
//				//testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())
//				//
//				//fmt.Println("unitteset Mysql配置json:", testMysqlJson)
//				//
//				//testRedisJson, _ := json.ToJson(config.GetRedisConfig())
//				//
//				//fmt.Println("unitteset redis配置json:", testRedisJson)
//				//
//				////添加token操作
//				//ginCtx := gin.Context{}
//				//ginCtx.Set(consts.AppHeaderTokenName, "abc")
//				////获得一个顶级上下文
//				//ctx := context.Background()
//				////返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
//				//ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)
//
//				//var testConst =consts.ProcessStatusCategoryIteration
//
//				var projectId = project.ID
//
//				//projectObj
//				//typeId := int64(intn)
//
//				cacheUserInfo, errCur := orgfacade.GetCurrentUserRelaxed(ctx)
//
//				fmt.Println(errCur)
//
//				currentUserId := cacheUserInfo.UserId
//				orgId := cacheUserInfo.OrgId
//
//				convey.Convey("createIssue insert Project proecss ", func() {
//
//					processId, _ := idfacade.ApplyPrimaryIdRelaxed((&po.PpmPrsProjectObjectTypeProcess{}).TableName())
//
//					projectObjectTypeProcess := &po.PpmPrsProjectObjectTypeProcess{
//						Id:                  processId,
//						OrgId:               orgId,
//						ProjectId:           projectId,
//						ProjectObjectTypeId: 1117,
//						ProcessId:           1118, //这个是一开始默认的不能写死  1118代表迭代
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
//					convey.Convey("insertPriority....  ", func() {
//						//PpmPrsPriority
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
//						_, _ = CreatePriority(cacheUserInfo.OrgId, vo.CreatePriorityReq{
//							OrgID: orgId,
//							Type:  consts.PriorityTypeIssue})
//
//						convey.Convey("mock createIssue input ", func() {
//
//							var iterationName = "alan创建的迭代" + strconv.Itoa(intn)
//							planStartTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
//							planEndTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")
//
//							iterationInput := vo.CreateIterationReq{
//								ProjectID:     projectId,
//								Name:          iterationName,
//								Owner:         currentUserId,
//								PlanStartTime: types.Time(planStartTime),
//								PlanEndTime:   types.Time(planEndTime)}
//
//							modelsVoid, err := CreateIteration(orgId, currentUserId, iterationInput)
//
//							convey.So(modelsVoid, convey.ShouldNotBeNil)
//							convey.So(err, convey.ShouldBeNil)
//
//							convey.Convey("Test IterationIssueRelate", func() {
//
//								addIssueIds := make([]int64, 0)
//								delIssueIds := make([]int64, 0)
//
//								input := vo.IterationIssueRealtionReq{
//									IterationID: modelsVoid.ID,
//									AddIssueIds: addIssueIds,
//									DelIssueIds: delIssueIds}
//
//								resp, err := IterationIssueRelate(orgId, currentUserId, input)
//
//								convey.So(resp, convey.ShouldNotBeNil)
//								convey.So(err, convey.ShouldBeNil)
//							})
//						})
//					})
//				})
//			})
//
//		})
//	}))
//}

func TestUpdateIterationStatus(t *testing.T) {

	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {

		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1070), CorpId: "1", OrgId: 17}

		}

		cache.Set("polaris:sys:user:token:abc", cacheUserInfoJson)

		convey.Convey("Test CreateProject", func() {

			rand.Seed(time.Now().Unix())
			intn := rand.Intn(10000)

			var preCode = "alanTest" + strconv.Itoa(intn)

			startTime := types.NowTime()

			endTime, _ := time.Parse(consts.AppTimeFormat, "2020-10-09 12:20:22")

			planEndTime := types.Time(endTime)

			fmt.Println("时间", endTime)

			var remark = "哈哈哈"

			projectName := "alan测试项目" + strconv.Itoa(intn)

			input := vo.CreateProjectReq{Name: projectName, PreCode: &preCode, Owner: 1070, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark}

			project, err := CreateProject(projectvo.CreateProjectReqVo{
				OrgId:  cacheUserInfo.OrgId,
				UserId: cacheUserInfo.UserId,
				Input:  input,
			})

			if err != nil {
				fmt.Println("错误.......")
			}

			fmt.Printf("项目%+v", project)

			convey.Convey("Test CreateIteration ", func() {
				//
				//testMysqlJson, _ := json.ToJson(config.GetMysqlConfig())
				//
				//fmt.Println("unitteset Mysql配置json:", testMysqlJson)
				//
				//testRedisJson, _ := json.ToJson(config.GetRedisConfig())
				//
				//fmt.Println("unitteset redis配置json:", testRedisJson)
				//
				////添加token操作
				//ginCtx := gin.Context{}
				//ginCtx.Set(consts.AppHeaderTokenName, "abc")
				////获得一个顶级上下文
				//ctx := context.Background()
				////返回父上下文  在父上下文中设置 key/value  这边是 GinContextKey:ginCtx
				//ctx = context.WithValue(ctx, "GinContextKey", &ginCtx)

				//var testConst =consts.ProcessStatusCategoryIteration

				var projectId = project.ID

				//projectObj
				//typeId := int64(intn)

				cacheUserInfo, errCur := orgfacade.GetCurrentUserRelaxed(ctx)

				fmt.Println(errCur)

				currentUserId := cacheUserInfo.UserId
				orgId := cacheUserInfo.OrgId

				convey.Convey("createIssue insert Project proecss ", func() {
					convey.Convey("insertPriority....  ", func() {
						//PpmPrsPriority
						//id, _ := idfacade.ApplyPrimaryIdRelaxed((&bo.PriorityBo{}).TableName())
						//
						//priority := bo.PriorityBo{
						//	Id:    id,
						//	OrgId: orgId,
						//	Type:  consts.PriorityTypeIssue}
						//
						//_ = mysql.TransX(func(tx sqlbuilder.Tx) error {
						//	err4 := mysql.TransInsert(tx, &priority)
						//	if err4 != nil {
						//		log.Errorf(consts.Mysql_TransInsert_error_printf, err4)
						//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err4)
						//	}
						//	return nil
						//})

						//_, _ = CreatePriority(cacheUserInfo.OrgId, vo.CreatePriorityReq{
						//	OrgID: orgId,
						//	Type:  consts.PriorityTypeIssue})

						convey.Convey("mock createIssue input ", func() {

							var iterationName = "alan创建的迭代" + strconv.Itoa(intn)
							planStartTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:20:22")
							planEndTime, err := time.Parse(consts.AppTimeFormat, "2019-09-26 12:30:22")

							iterationInput := vo.CreateIterationReq{
								ProjectID:     projectId,
								Name:          iterationName,
								Owner:         currentUserId,
								PlanStartTime: types.Time(planStartTime),
								PlanEndTime:   types.Time(planEndTime)}

							modelsVoid, err := CreateIteration(orgId, currentUserId, iterationInput)

							convey.So(modelsVoid, convey.ShouldNotBeNil)
							convey.So(err, convey.ShouldBeNil)

							convey.Convey("Test IterationStatusTypeStat", func() {

								input := vo.UpdateIterationStatusReq{
									// 迭代id
									ID: modelsVoid.ID,
									// 要更新的状态id  ppm_prs_process_status表中初始化的 category=2(迭代类型的process)
									NextStatusID: 5}

								resp, err := UpdateIterationStatus(orgId, currentUserId, input)

								convey.So(resp, convey.ShouldNotBeNil)
								convey.So(err, convey.ShouldBeNil)
							})
						})
					})
				})
			})

		})
	}))
}

func TestIterationInfo2(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		res, err := IterationInfo(1872, vo.IterationInfoReq{ID: 2308})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}

func TestTime(t *testing.T) {
	now := types.NowTime()
	t.Log(now)
	t.Log(time.Time(now).String())
	formatTime, _ := time.ParseInLocation(consts.AppTimeFormat, now.String(), time.Local)
	t.Log(formatTime.Equal(time.Time(now)))
}

func TestUpdateIteration2(t *testing.T) {
	convey.Convey("TestUpdateIteration2", t, test.StartUp(func(ctx context.Context) {
		owner := int64(1284)
		t.Log(UpdateIteration(1111, 1284, vo.UpdateIterationReq{ID: 1008, Owner: &owner, UpdateFields: []string{"owner"}}))

	}))
}

func TestUpdateIterationStatus2(t *testing.T) {
	convey.Convey("TestUpdateIteration2", t, test.StartUp(func(ctx context.Context) {
		t.Log(UpdateIterationStatus(1872, 24370, vo.UpdateIterationStatusReq{
			ID:                  2369,
			NextStatusID:        3,
			BeforeStatusEndTime: date.ParseTime("2022-04-02 00:00:00"),
			NextStatusStartTime: date.ParseTime("2022-05-20 00:00:00"),
		}))
	}))
}

func TestUpdateIterationSort(t *testing.T) {
	convey.Convey("TestUpdateIteration2", t, test.StartUp(func(ctx context.Context) {
		t.Log(UpdateIterationSort(1108, 1250, vo.UpdateIterationSortReq{
			IterationID: 1001,
			BeforeID:    1002,
		}))
	}))
}
