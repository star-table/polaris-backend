package service

import (
	"context"
	"fmt"
	"net/url"
	"testing"
	"time"

	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/config"
	consts1 "github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/appvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/appfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain/lc_pro_domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/dgrijalva/jwt-go"
	"github.com/getsentry/sentry-go"
	"github.com/smartystreets/goconvey/convey"
)

func TestGetProjectRelation(t *testing.T) {
	convey.Convey("Test GetProjectRelation", t, test.StartUp(func(ctx context.Context) {
		t.Log(GetProjectRelation(1001, []int64{1, 2, 3}))
	}))
}
func TestArchiveProject(t *testing.T) {
	convey.Convey("Test ArchiveProject", t, test.StartUp(func(ctx context.Context) {
		cacheUserInfo, _ := orgfacade.GetCurrentUserRelaxed(ctx)

		cacheUserInfoJson, _ := json.ToJson(cacheUserInfo)

		log.Info("缓存用户信息" + cacheUserInfoJson)

		if cacheUserInfo == nil {
			cacheUserInfo = &bo.CacheUserInfoBo{OutUserId: "aFAt7VhhZ2zcE8mdFFWWPAiEiE", SourceChannel: "dingtalk", UserId: int64(1001), CorpId: "1", OrgId: 1001}

		}

		t.Log(ArchiveProject(cacheUserInfo.OrgId, cacheUserInfo.UserId, 1))
		t.Log(CancelArchivedProject(cacheUserInfo.OrgId, cacheUserInfo.UserId, 1))
	}))
}

func TestOrgProjectMembers(t *testing.T) {

	convey.Convey("Test OrgProjectMembers", t, test.StartUp(func(ctx context.Context) {

		resp, err := OrgProjectMembers(projectvo.OrgProjectMemberReqVo{
			ProjectId: 1001,
			OrgId:     1001,
			UserId:    1001,
		})

		convey.ShouldBeNil(err)
		convey.ShouldNotBeNil(resp)
	}))
}

func TestProjects(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		// statusType := 3
		res, _ := Projects(projectvo.ProjectsRepVo{OrgId: 2585, UserId: 29572, Page: -1, Size: 20, ProjectExtraBody: projectvo.ProjectExtraBody{
			Input: &vo.ProjectsReq{
				IsFiling: intToPtr(3),
			},
		}})
		t.Log(json.ToJsonIgnoreError(res))
	}))
}

func TestUpdateProject(t *testing.T) {
	convey.Convey("Test ArchiveProject", t, test.StartUp(func(ctx context.Context) {
		a := 1 // 文小兰 3307   苏13 3306
		res, err := UpdateProject(projectvo.UpdateProjectReqVo{
			OrgId:         1588,
			UserId:        24391,
			SourceChannel: "fs",
			Input: vo.UpdateProjectReq{
				ID:             10643,
				IsCreateFsChat: &a,
				UpdateFields:   []string{"isCreateFsChat"},
			},
		})
		t.Log(res)
		t.Log(err)
	}))
}

func TestProjectInfo(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(2373)
		resp, err := ProjectInfo(orgId, 29611, vo.ProjectInfoReq{ProjectID: 13924})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestStarProject(t *testing.T) {
	convey.Convey("Test StarProject", t, test.StartUp(func(ctx context.Context) {
		t.Log(StarProject(projectvo.ProjectIdReqVo{ProjectId: 1006, SourceChannel: "", UserId: 1007, OrgId: 1003}))
	}))
}

func TestJudgeProjectFiling2(t *testing.T) {
	a := "[{\"id\":\"1\",\"value\":\"1\"},{\"name\":1}]"
	b := &[]map[string]interface{}{}
	c := json.FromJson(a, b)
	fmt.Println(json.ToJsonIgnoreError(b))
	fmt.Println(c)
}

func TestDoSentry(t *testing.T) {
	convey.Convey("Test TestExportData", t, test.StartUp(func(ctx context.Context) {
		sentryConfig := config.GetSentryConfig()
		sentryDsn := ""
		if sentryConfig != nil {
			sentryDsn = sentryConfig.Dsn
		}
		dsn := sentryDsn
		applicationName := config.GetApplication().Name
		env := "dev"
		connectSentryErr := sentry.Init(sentry.ClientOptions{
			Dsn:         dsn,
			ServerName:  applicationName,
			Environment: env,
		})
		errMsg := "thisIsTestErrorForSentry"
		stackInfo := "【极星警告】test error stack info. "
		event := sentry.NewEvent()
		event.EventID = sentry.EventID(errMsg)
		event.Message = errMsg + "\n" + stackInfo
		event.Environment = env
		event.ServerName = applicationName
		event.Tags[consts1.TraceIdKey] = threadlocal.GetTraceId()
		if connectSentryErr == nil {
			sentry.CaptureEvent(event)
			sentry.Flush(time.Second * 5)
		} else {
			fmt.Println(connectSentryErr)
		}
		fmt.Println(event.Message)
	}))
}

func TestGetProjectIdByAppId(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		result, err := GetProjectIdByAppId(projectvo.GetProjectIdByAppIdReqVo{
			AppId: "12325412121",
		})
		if err != nil {
			t.Error(err)
		}
		t.Log(json.ToJsonIgnoreError(result))
	}))
}

func TestJsonDecodeToUserInfo(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		// {\"code\":0,\"message\":\"OK\",\"data\":{\"outUserId\":\"\",\"sourceChannel\":\"\",\"userId\":3,\"corpId\":\"\",\"orgId\":202505060582035456}}\n
		json1 := `{"code":0,"message":"OK","data":{"outUserId":"","sourceChannel":"","userId":3,"corpId":"","orgId":202505060582035456}}\n`
		respVo := &orgvo.CacheUserInfoVo{}
		_ = json.FromJson(json1, respVo)
		t.Log(respVo.CacheInfo.OrgId)
		cacheUserInfo := respVo.CacheInfo
		resp := projectfacade.Projects(projectvo.ProjectsRepVo{
			Page: 1,
			Size: 10,
			ProjectExtraBody: projectvo.ProjectExtraBody{
				Params: map[string]interface{}{},
				Order:  make([]*string, 0),
				Input:  &vo.ProjectsReq{},
			},
			UserId:        cacheUserInfo.UserId,
			OrgId:         cacheUserInfo.OrgId,
			SourceChannel: cacheUserInfo.SourceChannel,
		})
		if resp.Failure() {
			t.Error(resp.Err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestIfToInt64(t *testing.T) {
	json1 := `{"outUserId":"","sourceChannel":"","userId":3,"corpId":"","orgId":202505060582035456}`
	m1 := make(map[string]interface{}, 0)
	json.FromJson(json1, &m1)
	for k, v := range m1 {
		fmt.Printf("%s-%s\n", k, v)
	}
	t.Log(json.ToJsonIgnoreError(m1))
}

func TestCreateLessCodeApp(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(1403)
		userId := int64(27707)
		name := "appNameqwewqeq 123"
		s := fmt.Sprintf("?appType=%d&name=%s&orgId=%d&userId=%d", 4, url.QueryEscape(name), &orgId, &userId)
		t.Log(s)
	}))
}

func TestCreateLessCodeAppForPro(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(202505060582035456)
		userId := int64(3)
		extendsId := int64(1380077674080432129)
		newAppId, err := domain.CreateAppInLessCode(orgId, userId, "项目应用表", extendsId, nil, 1, "{}", 1)
		if err != nil {
			log.Error(err)
			return
		}
		// newAppId 表示项目对应的应用。
		t.Log(newAppId) // 1380431404755775490
		// 基于该 appId 创建任务数据
	}))
}

func TestFormCreatePriority(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(202505060582035456)
		userId := int64(3)
		appId := int64(1380431404755775490)
		projectId := int64(8216)
		resp := appfacade.FormCreatePriority(&appvo.FormCreatePriorityReq{
			AppId:  appId,
			OrgId:  orgId,
			UserId: userId,
			Form: []*appvo.FormCreatePriorityReqForm{
				{
					ProjectIds: []int64{projectId},
					LangCode:   "p0",
					Name:       "最高",
					Type:       "",
					Sort:       0,
					BgStyle:    "",
					FontStyle:  "",
					IsDefault:  "",
					Remark:     "",
				},
			},
		})
		if resp.Failure() {
			t.Error(resp.Error())
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestFormCreateIssue(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(202505060582035456)
		userId := int64(3)
		appId := int64(1380431404755775490)
		resp := appfacade.FormCreateIssue(&appvo.FormCreateIssueReq{
			AppId:  appId,
			OrgId:  orgId,
			UserId: userId,
			Form: []*appvo.FormCreateIssueReqForm{
				{
					Title:     "this is first issue title003",
					Status:    "todo",
					StartTime: "2021-04-09 10:19:10",
					EndTime:   "2021-04-09 10:19:10",
					Priority:  0,
					Owner:     userId,
					Followers: []int64{},
					Tags:      []int64{},
				},
			},
		})
		if resp.Failure() {
			t.Error(resp.Error())
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

//// 暂时无需通过 dsl api 进行数据的增删改。
//func TestFormCreateIssueWithDsl(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		orgId := int64(202505060582035456)
//		// userId := int64(3)
//		appId := int64(1380431404755775490)
//		formId, err := domain.GetFormIdByAppId(orgId, appId)
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		var tb = WrapperTableName(orgId, formId)
//		exec := LcDsL.NewExecutor()
//		exec.Save(tb, map[string]interface{}{
//			"id":   333,
//			"data": []int{},
//		})
//		url := ""
//		LcNet.Post(url, exec, "application/json")
//
//	}))
//}

func TestUpdateAppPermissionGroupOptAuth(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		appId := int64(1384338246158446593)
		if err := lc_pro_domain.UpdateOpForAppPermissionGroup(appId); err != nil {
			t.Log(err)
			return
		}
		t.Log("ok")
	}))
}

func TestOpenPriorityList(t *testing.T) {
	convey.Convey("TestOpenPriorityList", t, test.StartUp(func(ctx context.Context) {
		resp, err := OpenPriorityList(1242)
		if err != nil {
			t.Error(err)

		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestOpenGetIterationList(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp, err := OpenGetIterationList(projectvo.OpenGetIterationListReqVo{
			OrgId:     1242,
			ProjectId: 8271,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestOpenGetDemandSourceList(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp, err := OpenGetDemandSourceList(projectvo.OpenGetDemandSourceListReqVo{
			OrgId:     1242,
			ProjectId: 8271,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestOpenGetPropertyList(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp, err := OpenGetPropertyList(projectvo.OpenGetPropertyListReqVo{
			OrgId:     1242,
			ProjectId: 8271,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestOpenGetToken(t *testing.T) {
	key := "c10ec58587de4830b39c0511a847c521"
	secret := "4f6d623f7e9243dd9570f05d5d3a7ab3"
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := make(jwt.MapClaims)
		claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
		claims["iat"] = time.Now().Unix()
		claims["_key"] = key
		token.Claims = claims
		accessToken, _ := token.SignedString([]byte(secret))
		fmt.Println(accessToken)
	}))
}

func TestGetProcessStatus1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp, err := domain.GetProjectStatus(1574, 10590)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(resp))
	}))
}

func TestDeleteLcApp(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		appIds := []int64{
			1466397049493749761, 1467833219842932737,
		}
		orgId := int64(1574)
		// appId := int64(1465590454346113026)
		for _, appId := range appIds {
			lcResp := appfacade.DeleteLessCodeApp(&appvo.DeleteLessCodeAppReq{
				AppId:  appId,
				OrgId:  orgId,
				UserId: 24370,
			})
			if lcResp.Failure() {
				log.Error(lcResp.Error())
				return
			}
			t.Log(json.ToJsonIgnoreError(lcResp))
		}
	}))
}

func TestGetColumnConfigForIssueStatus(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		//t.Log(domain.GetAppDefaultColumnConfig())
		t.Log(json.ToJsonIgnoreError(domain.GetDefaultColumnForTaskBar()))
	}))
}

func TestQueryIssueStat1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		statData, err := domain.GetIssueFinishedStatForProjects(1081, 1012, 1511315533113843713, []int64{1099})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(statData))
	}))
}

func TestCreateApp1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		summaryAppId := int64(1510963474333515778)
		newAppId, err := domain.CreateAppInLessCode(1070, 1012, "test-su-p005-1", summaryAppId, nil, 1088, "{}", 1)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(newAppId)
	}))
}

func TestImplode1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		finishedTypeIds := []int64{consts.StatusTypeComplete}
		finishedTypeIdsStr := "'" + str.Int64Implode(finishedTypeIds, "','") + "'"
		t.Log(finishedTypeIdsStr)
	}))
}

func TestGetTableSchemas1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		resp := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
			OrgId:  1061,
			UserId: 1021,
			Input: &tablePb.ReadTableSchemasRequest{
				TableIds:          []int64{1510238783935614976},
				ColumnIds:         nil,
				IsNeedDescription: false,
			},
		})
		if resp.Failure() {
			t.Error(resp.Error())
			return
		}
		t.Log(resp)
	}))
}

func TestGetIssueStatInfo(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		res, err := domain.GetIssueFinishedStatForProjects(1888, 24370, 1479328877254860801, []int64{11812})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}

func TestGetProjectInfoByOrgIds(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		orgIds := []int64{2373}
		projects, err := GetProjectInfoByOrgIds(orgIds)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(projects))
	}))
}

func TestQueryAsyncTaskInfo1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		res, err := QueryProcessForAsyncTask(2373, &projectvo.QueryProcessForAsyncTaskReqVoData{
			TaskId: "imp_1567474958894989314_t1567474959624704000",
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}
