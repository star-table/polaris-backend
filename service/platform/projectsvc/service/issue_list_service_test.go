package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/spf13/cast"

	kratosNacos "github.com/go-kratos/kratos/contrib/registry/nacos/v2"
	"github.com/go-kratos/kratos/v2/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"

	"github.com/star-table/polaris-backend/facade/tablefacade"

	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/middleware/validate"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	tablePb "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/polaris-backend/facade/permissionfacade"

	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/magiconair/properties/assert"
	nacosVo "github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/smartystreets/goconvey/convey"
)

func TestCreateProject(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		startTime := types.NowTime()
		endTime, _ := time.Parse(consts.AppTimeFormat, "2023-10-09 12:20:22")
		planEndTime := types.Time(endTime)
		var remark = "哈哈哈"
		projectTypeId := int64(consts.ProjectTypeEmpty)
		input := vo.CreateProjectReq{Name: "我的空项目988", Owner: 29572, PublicStatus: 1, PlanStartTime: &startTime, PlanEndTime: &planEndTime, Remark: &remark, ProjectTypeID: &projectTypeId}

		project, err := CreateProject(projectvo.CreateProjectReqVo{
			OrgId:         2373,
			UserId:        29572,
			Input:         input,
			SourceChannel: "",
		})
		fmt.Println(project, err)
		time.Sleep(time.Second * 10)
	}))
}

//func TestGetIssueIds(t *testing.T) {
//
//	convey.Convey("TestGetIssueIds", t, test.StartUp(func(ctx context.Context) {
//		beforePlanEndTime := "2015-01-01 11:11:11"
//		afterPlanEndTIme := "2020-01-01 11:11:11"
//
//		resp, err := GetIssueRemindInfoList(projectvo.GetIssueRemindInfoListReqVo{
//			Page: 1,
//			Size: 10,
//			Input: projectvo.GetIssueRemindInfoListReqData{
//				BeforePlanEndTime: &beforePlanEndTime,
//				AfterPlanEndTime:  &afterPlanEndTIme,
//			},
//		})
//		t.Log(json.ToJsonIgnoreError(resp))
//		assert.Equal(t, err, nil)
//	}))
//
//}

func TestIssueListStat(t *testing.T) {
	convey.Convey("TestGetIssueIds", t, test.StartUp(func(ctx context.Context) {
		res, err := IssueListStat(1003, 1007, 10125)
		t.Log(json.ToJsonIgnoreError(res))
		assert.Equal(t, err, nil)
	}))
}

func TestMirrorStat(t *testing.T) {
	json.ToJsonIgnoreError(nil)
	convey.Convey("MirrorStat", t, test.StartUp(func(ctx context.Context) {
		discover := newDiscovery()
		tablefacade.InitGrpcClient(discover)

		res, err := MirrorStat(2373, 33988, []string{"1612791244172394498", "1612791194197262338"})
		t.Log(json.ToJsonIgnoreError(res))
		assert.Equal(t, err, nil)
	}))
}

func TestStringSlice(t *testing.T) {
	convey.Convey("TestGetIssueIds", t, test.StartUp(func(ctx context.Context) {
		body1 := `你好hello world....{"key":"some key", "value":"test value..."}`
		body1 = str.Substr(body1, 0, 1)
		if body1 != "你" {
			t.Error("sub str error 001")
			return
		}
		fmt.Println(body1)
	}))
}

func TestGetIssueRemindInfoList(t *testing.T) {
	a := "{\"custom_Field\":\"null\"}"
	b := &po.PpmPriIssue{}
	_ = json.FromJson(a, b)
	c := &bo.IssueBo{}
	_ = copyer.Copy(b, c)
	fmt.Println(json.ToJsonIgnoreError(b))
	fmt.Println(json.ToJsonIgnoreError(c))
}

func TestHandleGroupChatUserInsProIssue(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		var isOk bool
		//isOk, err := HandleGroupChatUserInsProIssue(&projectvo.HandleGroupChatUserInsProIssueReq{
		//	OpenChatId: "oc_ad86229a8fd34a2172baafccda15dd4b",
		//	OpenId:     "ou_febd84ee87595dc83759485d68df43d7",
		//})
		//isOk, err := HandleGroupChatUserInsProProgress(&projectvo.HandleGroupChatUserInsProProgressReq{
		//	OpenChatId: "oc_e066d60a96eda963dd80c2fe7f4d55c8",
		//	OpenId:     "ou_febd84ee87595dc83759485d68df43d7",
		//})
		isOk, err := HandleGroupChatAtUserName(&projectvo.HandleGroupChatUserInsAtUserNameReq{
			OpenChatId:   "oc_e066d60a96eda963dd80c2fe7f4d55c8",
			AtUserOpenId: "ou_febd84ee87595dc83759485d68df43d7",
			OpUserOpenId: "ou_febd84ee87595dc83759485d68df43d7",
		})
		t.Log(isOk, err)
		//HandleGroupChatAtUserNameWithIssueTitle(&projectvo.HandleGroupChatUserInsAtUserNameWithIssueTitleReq{
		//	OpenChatId:    "oc_edcb639b2b721d3e7afe41b3fad11bf1",
		//	AtUserOpenId:  "ou_83776f77689ab1a5f399d5aafa27a0a4",
		//	IssueTitleStr: "指令创建任务001",
		//})
	}))
}

func TestMysql(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		fmt.Println(reflect.TypeOf("aa").Kind().String())
		fmt.Println(reflect.TypeOf(11).Kind().String())
		var a interface{}
		a = []int{1, 2}
		fmt.Println(reflect.TypeOf(a).Kind().String())
	}))
}

//func TestGetIssueIdsByOrgId(t *testing.T) {
//	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
//		resp, err := GetIssueIdsByOrgId(1213, 3306, &projectvo.GetIssueIdsByOrgIdReq{
//			Page: 200,
//			Size: 20,
//		})
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		t.Log(json.ToJsonIgnoreError(resp))
//	}))
//}

//func TestInsertIssueProRelation(t *testing.T) {
//	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
//		err := InsertIssueProRelation(1213, 3306, &projectvo.InsertIssueProRelationReq{
//			IssueIds: []int64{62074},
//		})
//		if err != nil {
//			t.Error(err)
//			return
//		}
//		t.Log("exec end...")
//	}))
//}

func TestGetProject(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		proBo, err := domain.GetProject(1574, 10809)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(proBo))
	}))
}

func TestCheckIsProAdmin(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		isAdmin, err := domain.CheckIsProAdmin(1574, 1458041642576609282, 24370, nil) // 24462
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(isAdmin)

		optAuthResp := permissionfacade.GetAppAuth(1574, 1458041642576609282, 24370, 0)
		if optAuthResp.Failure() {
			t.Error(optAuthResp.Message)
			return
		}
		isAdmin, err = domain.CheckIsProAdmin(1574, 1458041642576609282, 24370, &optAuthResp.Data)
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(isAdmin)
	}))
}

func TestForProject(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		projectId := int64(61958)
		appId := "1598148075753406466"
		tableId := "1598148076478926848"
		fmt.Println(LcHomeIssuesForProject(2373, 29572, 1, 10000, &projectvo.HomeIssueInfoReq{
			ProjectID: &projectId,
			MenuAppID: &appId,
			TableID:   &tableId,
		}, true))
	}))
}

func int64ToPtr(val int64) *int64 {
	return &val
}

func intToPtr(val int) *int {
	return &val
}

func TestIssueStatusTypeStat1(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		//projectId := int64(13647)
		//relateType := consts.IssueRelationTypeOwner
		// 29564: luoDD;29549: suhy
		res, err := IssueStatusTypeStat(2373, 29572, &vo.IssueStatusTypeStatReq{
			ProjectID: int64ToPtr(61501),
			//IterationID: int64ToPtr(61500),
			//RelationType: &relateType,
		})
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}

func TestLcHomeIssuesForAll(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		appId := "1499356343193141249"
		tableId := "0"
		fmt.Println(LcHomeIssuesForAll(2373, 29572, 1, 10000, &projectvo.HomeIssueInfoReq{MenuAppID: &appId, ProjectID: int64ToPtr(0), TableID: &tableId}, false))
	}))
}

func TestLcHomeIssuesForProject(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		appId := "1562352285646946306"
		tableId := "1562352286842228736"
		fmt.Println(LcHomeIssuesForProject(2585, 29572, 1, 10000, &projectvo.HomeIssueInfoReq{MenuAppID: &appId, ProjectID: int64ToPtr(61021), TableID: &tableId}, false))
	}))
}

func TestGrpc(t *testing.T) {
	convey.Convey("test", t, test.StartUp(func(ctx context.Context) {
		conn, err := grpc.DialInsecure(
			context.Background(),
			grpc.WithEndpoint("127.0.0.1:9000"),
			grpc.WithTimeout(15*time.Second),
			grpc.WithMiddleware(
				recovery.Recovery(),
				validate.Validator(),
			),
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		client := tablePb.NewRowsClient(conn)
		reply, err := client.ListRaw(context.Background(), &tablePb.ListRawRequest{
			DbType:    tablePb.DbType_slave3,
			Condition: &tablePb.Condition{Type: tablePb.ConditionType_and, Conditions: domain.GetNoRecycleCondition()},
			Page:      1,
			Size:      1000,
			TableId:   1562352286842228736,
		})
		fmt.Println(reply, err)
	}))
}

func newDiscovery() registry.Discovery {
	sc, cc := getNacosServerAndClientConfig()
	client, err := clients.NewNamingClient(
		nacosVo.NacosClientParam{
			ServerConfigs: sc,
			ClientConfig:  &cc,
		},
	)
	if err != nil {
		log.Error(err)
		panic(err)
		return nil
	}

	r := kratosNacos.New(client)

	return r
}

func getNacosServerAndClientConfig() ([]constant.ServerConfig, constant.ClientConfig) {
	return []constant.ServerConfig{
			*constant.NewServerConfig("172.19.98.19", cast.ToUint64(30048)),
		},
		constant.ClientConfig{
			NamespaceId:         "public", //namespace id
			TimeoutMs:           5000,
			NotLoadCacheAtStart: true,
			LogLevel:            "error",
		}
}
