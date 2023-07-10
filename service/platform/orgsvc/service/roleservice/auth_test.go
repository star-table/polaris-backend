package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/star-table/common/core/util/slice"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/operation"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/test"
	"github.com/smartystreets/goconvey/convey"
)

func TestAuthenticate(t *testing.T) {

	convey.Convey("Test AppendIterationStat", t, test.StartUp(func(ctx context.Context) {
		t.Log(json.ToJsonIgnoreError(config.GetMysqlConfig()))

		orgId := int64(1001)
		projectId := int64(1)
		issueId := int64(1083)
		projectAuthInfo := &bo.ProjectAuthBo{
			Id:    projectId,
			Owner: 1016,
		}
		issueAuthInfo := &bo.IssueAuthBo{
			Id:    issueId,
			Owner: []int64{1016},
		}

		err := Authenticate(orgId, 1029, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/T", operation.Delete, nil)
		t.Log(err)
		err = Authenticate(orgId, 1004, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/T", operation.Delete, nil)
		t.Log(err)
		err = Authenticate(orgId, 1031, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/T", operation.Delete, nil)
		t.Log(err)
		err = Authenticate(orgId, 1031, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/T", operation.Modify, nil)
		t.Log(err)
		err = Authenticate(orgId, 1031, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/ProConfig", operation.ModifyStatus, nil)
		t.Log(err)
		err = Authenticate(orgId, 1031, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}", operation.ModifyStatus, nil)
		t.Log(err)
		err = Authenticate(orgId, 1031, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/0/Test/TestApp", operation.Modify, nil)
		t.Log(err)

		for i := 0; i < 2; i++ {
			if i == 0 {
				go func() {
					err = Authenticate(orgId, 1031, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/0/Test/TestDevice", operation.Modify, nil)
					t.Log(err)
				}()
			} else {
				go func() {
					err = Authenticate(orgId, 1031, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/0/Test/TestReport", operation.Modify, nil)
					t.Log(err)
				}()
			}
		}

	}))
}

func TestAuthenticate3(t *testing.T) {
	convey.Convey("Test AppendIterationStat", t, test.StartUp(func(ctx context.Context) {
		orgId := int64(1222)
		projectId := int64(1513)
		issueId := int64(7079)
		projectAuthInfo := &bo.ProjectAuthBo{
			Id:           projectId,
			Owner:        1464,
			Participants: []int64{1464, 1477, 1466, 1458},
			Followers:    []int64{1464, 1466, 1465},
		}
		issueAuthInfo := &bo.IssueAuthBo{
			Id:    issueId,
			Owner: []int64{1464},
		}
		err := Authenticate(orgId, 1458, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/4", operation.ModifyStatus, nil)
		t.Log(err)

	}))
}

func TestOperationMath(t *testing.T) {

	res, err := regexp.MatchString("Modify", "(View)|(Modify)|(ModifyStatus)|(Bind)|(Unbind)")
	if err != nil {
		panic(err)
	}
	t.Log(res)

	res, err = regexp.MatchString("Modify", "(View)|(ModifyStatus)|(Bind)|(Unbind)")
	if err != nil {
		panic(err)
	}
	t.Log(res)

	res, err = regexp.MatchString("Modify", "(View)|(Bind)|(Unbind)")
	if err != nil {
		panic(err)
	}
	t.Log(res)

	res, err = regexp.MatchString(`View`, `.*?`)
	if err != nil {
		panic(err)
	}
	t.Log(res)
	res, err = regexp.MatchString("View", "View")
	if err != nil {
		panic(err)
	}
	t.Log(res)
	//e,_:= regexp.Compile(`(View)|(Modify)|(ModifyStatus)|(Bind)|(Unbind)`)
	//mat := e.MatchString("View1")
	//t.Log(mat)
}

func TestAuthenticate2(t *testing.T) {

	convey.Convey("Test AppendIterationStat", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("加载配置", func() {

			convey.Convey("mock info", func() {
				//初始化参数
				projectAuthInfo := &bo.ProjectAuthBo{
					Id:    1,
					Owner: 1016,
				}
				issueAuthInfo := &bo.IssueAuthBo{
					Id:    1004,
					Owner: []int64{100234},
				}
				convey.Convey("mock info", func() {
					err := Authenticate(100233, 1039, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/T", operation.Delete, nil)

					convey.So(err, convey.ShouldNotBeNil)
				})
			})
		})
	}))
}

//项目负责人开发所有权限
func TestAuthenticateProjectOwn(t *testing.T) {

	convey.Convey("Test AppendIterationStat", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("加载配置", func() {

			convey.Convey("mock info", func() {
				//初始化参数
				projectAuthInfo := &bo.ProjectAuthBo{
					Id:    1,
					Owner: 1016,
				}
				issueAuthInfo := &bo.IssueAuthBo{
					Id:    1004,
					Owner: []int64{100234},
				}
				convey.Convey("mock info", func() {
					err := Authenticate(100233, 1077, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/T", operation.Delete, nil)

					convey.So(err, convey.ShouldBeNil)
				})
			})
		})
	}))
}

//项目负责人开发所有权限
func TestAuthenticateOperationGetRoleOperationListError(t *testing.T) {

	convey.Convey("Test AppendIterationStat", t, test.StartUp(func(ctx context.Context) {
		convey.Convey("加载配置", func() {

			convey.Convey("mock info", func() {
				//初始化参数
				projectAuthInfo := &bo.ProjectAuthBo{
					Id:    1045,
					Owner: 1010,
				}
				issueAuthInfo := &bo.IssueAuthBo{
					Id:      1281,
					Owner:   []int64{1010},
					Creator: 1023,
				}
				convey.Convey("mock info", func() {
					err := Authenticate(1004, 1023, projectAuthInfo, issueAuthInfo, "/Org/{org_id}/Pro/{pro_id}/Issue/4", operation.Comment, nil)
					t.Log(err)
					//convey.So(err, convey.ShouldBeNil)
				})
			})
		})
	}))
}

func TestOperationMath2(t *testing.T) {

	convey.Convey("单元测试,正则匹配1", t, func() {

		result, err := regexp.MatchString("Modify", "(View)|(Modify)|(ModifyStatus)|(Bind)|(Unbind)")

		convey.So(result, convey.ShouldEqual, true)

		convey.So(err, convey.ShouldEqual, nil)

		convey.Convey("单元测试,正则匹配2", func() {

			result, err = regexp.MatchString(`View`, `.*?`)

			convey.So(result, convey.ShouldEqual, false)

			convey.So(err, convey.ShouldEqual, nil)

			convey.Convey("单元测试,正则匹配3", func() {

				result, err = regexp.MatchString("View", "View")

				convey.So(result, convey.ShouldEqual, true)

				convey.So(err, convey.ShouldEqual, nil)
			})
		})
	})
}

func TestAuthenticate4(t *testing.T) {
	convey.Convey("Test AppendIterationStat", t, test.StartUp(func(ctx context.Context) {
		t.Log(Authenticate(1323, 1608, nil, nil, "/Org/{org_id}/Pro/{pro_id}/Member", "Bind", nil))
	}))
}

func TestAuthenticate5(t *testing.T) {
	convey.Convey("Test AppendIterationStat", t, test.StartUp(func(ctx context.Context) {
		t.Log(Authenticate(
			1058,
			1867,
			&bo.ProjectAuthBo{Id: 1726, Status: 2, PublicStatus: 1, IsFilling: 2, ProjectType: 1, Participants: []int64{1867}}, nil, "/Org/{org_id}/Pro/{pro_id}/Issue/4", "Create", nil))
	}))
}

func TestAuthContain1(t *testing.T) {
	optAuthRespData := []string{
		"Permission.Pro.Tag-Create",
		"hasRead",
		"hasCreate",
		"hasCopy",
		"hasUpdate",
		"hasDelete",
		"hasImport",
		"hasExport",
		"hasInvite",
		"hasShare",
		"hasEditMember",
	}
	opList := make([]string, 0)
	for _, item := range optAuthRespData {
		infos := strings.Split(item, "-")
		if len(infos) > 1 {
			opPrev := infos[0]
			opSuffixArr := strings.Split(infos[1], ",")
			for _, oneSuffix := range opSuffixArr {
				opList = append(opList, fmt.Sprintf("%s.%s", opPrev, oneSuffix))
			}
		}
	}
	operation := "Permission.Pro.Tag.Create"
	access, _ := slice.Contain(opList, operation)
	t.Log(access)
}
