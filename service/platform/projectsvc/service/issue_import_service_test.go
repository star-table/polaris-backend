package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/test"
	"github.com/jtolds/gls"
	"github.com/smartystreets/goconvey/convey"
)

func TestExportIssueTemplate(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		threadlocal.Mgr.SetValues(gls.Values{consts.AppHeaderLanguage: "zh"}, func() {
			//url, err := ExportIssueTemplate(2373, 61065, 1564148745195491329)
			url, err := ExportIssueTemplate(2373, 61500, 1582328826073976832)
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(url)
		})
	}))
}

//func Test_BatchImport_Sync(t *testing.T) {
//	convey.Convey("Test login", t, test.StartUp(func(ctx context.Context) {
//		var wg sync.WaitGroup
//		wg.Add(1)
//		for i := 0; i < 1; i++ {
//			go func() {
//				defer wg.Add(-1)
//				issueDetailId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueDetail)
//				t.Log(issueDetailId, err)
//			}()
//		}
//		wg.Wait()
//
//		t.Log("ok")
//	}))
//}

func TestImportIssues2(t *testing.T) {
	convey.Convey("Test TestImportIssues2", t, test.StartUp(func(ctx context.Context) {
		threadlocal.Mgr.SetValues(gls.Values{consts.AppHeaderLanguage: "cn"}, func() {
			//tableId := "1527557852065959936"
			//_, err := ImportIssues(2571, 1368436, vo.ImportIssuesReq{
			//	ProjectID: 13757,
			//	TableID:   tableId,
			//	//URL:                 "D:\\programData\\polaris\\projectsvc\\org_1213\\project_8387\\abc_任务统计.xlsx",
			//	//URL: "C:\\Users\\suhan\\Downloads\\导入模板_任务 (9).xlsx",
			//	URL:     "/Users/helayzhang/Desktop/abc.xlsx",
			tableId := "1592063825702555648"
			_, err := ImportIssues(2585, 29572, vo.ImportIssuesReq{
				ProjectID: 61791,
				TableID:   tableId,
				//URL:                 "D:\\programData\\polaris\\projectsvc\\org_1213\\project_8387\\abc_任务统计.xlsx",
				//URL: "C:\\Users\\suhan\\Downloads\\导入模板_任务 (9).xlsx",
				URL:     "/Users/admin/Downloads/导入模板_需求.xlsx",
				URLType: 2,
			})
			if err != nil {
				t.Error(err)
				return
			}
			t.Log("--end--")
			time.Sleep(time.Second * 10)
		})
	}))
}

func TestExportUserOrDeptSameNameList(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		threadlocal.Mgr.SetValues(gls.Values{consts.AppHeaderLanguage: "en"}, func() {
			url, err := ExportUserOrDeptSameNameList(1888, 0)
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(url)
		})
	}))
}

// 测试任务的导出
func TestExportData(t *testing.T) {
	convey.Convey("Test TestExportData", t, test.StartUp(func(ctx context.Context) {
		threadlocal.Mgr.SetValues(gls.Values{consts.AppHeaderLanguage: "zh"}, func() {
			var (
				orgId       int64  = 2585
				projectId   int64  = 61791
				iterationId int64  = 0
				tableId     string = "1592063825702555648"
			)
			input := projectvo.ExportIssueReqVo{
				UserId: 29572,
				OrgId:  orgId,
				Input: projectvo.ExportIssueReqVoData{
					ProjectId:      projectId,
					TableId:        tableId,
					IterationId:    iterationId,
					IsNeedDocument: true,
				},
			}
			url, err := ExportData(orgId, input)
			if err != nil {
				t.Errorf("%v\n", err)
				return
			}
			t.Log(url)
		})
	}))
}

func TestCopy1(t *testing.T) {
	list := []int{1, 2, 3, 4, 5, 6, 7}
	// 删除2，输出删除位置之前和之后的
	fmt.Printf("%v, %v", list[:1], list[6:])
	//copy(list, list[1:])
}

func delFunc(list *[]int) {
	delIndex := 1
	*list = append((*list)[:delIndex], (*list)[delIndex+1:]...)
}

func TestRecursiveDel(t *testing.T) {
	list := []int{1, 2, 3, 4}
	delIndex := 1
	// delFunc(&list)
	list = append(list[:delIndex], list[delIndex+1:]...)
	fmt.Printf("%v\n", list)

	list1 := []int{1}
	fmt.Printf("list1 slice: %v", list1[:1])
}

type Stu struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestDecodeJson(t *testing.T) {
	json1 := `{
"1300": {
        "title": "2020-11-10 09:00",
        "value": "2020-11-10 09:00",
        "fieldId": 1300
    }}`
	customFieldMap := make(map[string]bo.CustomFieldItemFrontend, 0)
	json.FromJson(json1, &customFieldMap)
	t.Log(json.ToJsonIgnoreError(customFieldMap))
}

func TestFloat1(t *testing.T) {
	f := 0.124
	s1 := fmt.Sprintf("%v", f)
	t.Log(s1)
}

func TestFromJson1(t *testing.T) {
	json1 := ` {
        "title": "文小兰",
        "value": [
            {
                "id": "3307",
                "userId": 3307,
                "name": "文小兰",
                "avatar": "https://s3-fs.pstatp.com/static-resource/v1/b1dfe7f3-2b41-4903-8530-2f33f9d6006g~?image_size=noop&cut_type=&quality=&format=png&sticker_format=.webp"
            }
        ],
        "fieldId": 3773
    }`
	obj := bo.CustomFieldItemFrontend{}
	json.FromJson(json1, &obj)
	t.Log(json.ToJsonIgnoreError(obj))
}

func TestGetBuiltInAndCustomColumns1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		threadlocal.Mgr.SetValues(gls.Values{consts.AppHeaderLanguage: "zh"}, func() {
			param := map[string]string{"usingType": "import"}
			builtInColumns, customColumns, _, err := GetBuiltInAndCustomColumns(1888, 11812, 0, consts.ProjectTypeNormalId, false, param, false)
			if err != nil {
				t.Error(err)
				return
			}
			t.Log(json.ToJsonIgnoreError(builtInColumns))
			t.Log(json.ToJsonIgnoreError(customColumns))
		})
	}))
}

func TestTestAssign(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		page, size := 1, 2000
		t.Log(page, size)
		arr1 := make([]int, 0, 10)
		t.Log(arr1)
		arr1 = append(arr1, 10)
		t.Log(arr1)
	}))
}

func TestCustomEnv(t *testing.T) {
	os.Setenv(consts.RunCustomEnvKey, "0")
	if util.IsChinaMobileEnv() {
		t.Error("err 1")
		return
	}
	os.Setenv(consts.RunCustomEnvKey, consts.CustomEnvChinaMobile)
	if util.IsChinaMobileEnv() {
		t.Log("env is ok!")
		return
	}
}

func TestTrimSpaceForCell1(t *testing.T) {
	input1 := " \n 父 任 务\n\t "
	str := strings.Trim(input1, "\n\r\t ")
	if str != "父 任 务" {
		t.Errorf("remove space not ok origin: [%s], after: [%s]", input1, str)
		return
	}
	t.Log("--end--")
}

//func TestImportAsyncTaskStep1(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		domain.UpdateAsyncTaskWithSucc(&projectvo.CreateIssueReqVo{
//			CreateIssue: vo.CreateIssueReq{
//				ExtraInfo: map[string]interface{}{
//					consts.CacheKeyAsyncTaskIdConstName: "tid_4784f29a64f94612bff9b96164901740",
//					consts.HelperFieldOperateUserId:     int64(29611),
//					consts.HelperFieldIsApplyTpl:        false,
//				},
//				ProjectID: 60452,
//				TableID:   "1555480088131145728",
//			},
//			OrgId: 2373,
//		}, 1, true)
//		t.Log("--end--")
//	}))
//}

//func TestImportAsyncTaskStep2(t *testing.T) {
//	json1 := "{\"createIssue\":{\"projectId\":60452,\"title\":\"su-08101102-t002\",\"priorityId\":0,\"typeId\":0,\"ownerId\":[],\"participantIds\":null,\"followerIds\":[],\"followerDeptIds\":null,\"planStartTime\":\"2022-08-10 10:11:11\",\"planEndTime\":\"0001-01-01 00:00:00\",\"planWorkHour\":null,\"versionId\":null,\"moduleId\":null,\"parentId\":null,\"remark\":null,\"remarkDetail\":null,\"mentionedUserIds\":null,\"iterationId\":0,\"issueObjectId\":null,\"issueSourceId\":null,\"issuePropertyId\":null,\"statusId\":0,\"children\":null,\"tags\":null,\"resourceIds\":null,\"customField\":null,\"auditorIds\":[],\"lessCreateIssueReq\":{\"issueStatus\":0,\"ownerId\":[],\"iterationId\":0,\"title\":\"su-08101102-t002\",\"projectId\":60452,\"createFrom\":\"import\"},\"beforeId\":null,\"afterId\":null,\"beforeDataId\":null,\"afterDataId\":null,\"asc\":null,\"isImport\":null,\"tableId\":\"1555480088131145728\",\"extraInfo\":{\"isApplyTpl\":false,\"asyncTaskId\":\"tid_f03498fd903443ca9e7bab2bb615ddb8\",\"operateUserId\":29611}},\"userId\":29611,\"orgId\":2373,\"sourceChannel\":\"fs\",\"inputAppId\":0,\"tableId\":0}"
//	createObj := projectvo.CreateIssueReqVo{}
//	json.FromJson(json1, &createObj)
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		domain.UpdateAsyncTaskWithSucc(&createObj, 1, true)
//		t.Log("--end--")
//	}))
//}
//
//func TestSetAsyncTaskErrInfo1(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		domain.UpdateAsyncTaskWithError(&projectvo.CreateIssueReqVo{
//			CreateIssue: vo.CreateIssueReq{
//				ExtraInfo: map[string]interface{}{
//					consts.CacheKeyAsyncTaskIdConstName: "1548901510149144577",
//				},
//			},
//			OrgId: 2373,
//		}, errs.ProjectNotExist, 1, true)
//		t.Log("--end--")
//	}))
//}

func TestCheckMap1(t *testing.T) {
	map1 := make(map[string]interface{}, 0)
	isOk := map1["nishi"] == nil
	t.Log(isOk)
}

func TestGetUserOrDeptSameNameList1(t *testing.T) {
	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
		res, err := GetUserOrDeptSameNameList(2373, "user")
		if err != nil {
			t.Error(err)
			return
		}
		t.Log(json.ToJsonIgnoreError(res))
	}))
}

//func TestCheckAsyncTaskIsExecutingForImport(t *testing.T) {
//	convey.Convey("Test", t, test.StartUp(func(ctx context.Context) {
//		isExecuting := domain.CheckAsyncTaskImportIsRunning(2373, 1567474958894989314, 1567474959624704000)
//		t.Log(isExecuting)
//	}))
//}

type Iss0908 struct {
	Title  string `json:"title"`
	Remark string `json:"remark"`
}

func TestJsonDecode1(t *testing.T) {
	m1 := map[string]interface{}{
		"title":  "title 1001",
		"remark": "title 1001",
	}
	obj1 := Iss0908{}
	copyer.Copy(m1, &obj1)
	t.Log(json.ToJsonIgnoreError(obj1))
}

func TestColAxisCal1(t *testing.T) {
	c1 := 'A' + 1
	t.Log(fmt.Sprintf("%c", c1))
}
