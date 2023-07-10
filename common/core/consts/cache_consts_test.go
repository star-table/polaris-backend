package consts

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConstCacheKey(t *testing.T) {

	convey.Convey("Test CacheKey", t, func() {

		convey.Convey("CacheKeyPrefix", func() {
			convey.So(CacheKeyPrefix, convey.ShouldEqual, "polaris:")
		})

		convey.Convey("CacheKeyOfSys", func() {
			convey.So(CacheKeyOfSys, convey.ShouldEqual, "sys:")
		})

		convey.Convey("CacheKeyOfOrg", func() {
			convey.So(CacheKeyOfOrg, convey.ShouldEqual, "org_{{."+CacheKeyOrgIdConstName+"}}:")
		})
		convey.Convey("CacheKeyOfUser", func() {
			convey.So(CacheKeyOfUser, convey.ShouldEqual, "user_{{."+CacheKeyUserIdConstName+"}}:")
		})

		convey.Convey("CacheKeyOfOutUser", func() {
			convey.So(CacheKeyOfOutUser, convey.ShouldEqual, "outuser_{{."+CacheKeyOutUserIdConstName+"}}:")
		})

		convey.Convey("CacheKeyOfProject", func() {
			convey.So(CacheKeyOfProject, convey.ShouldEqual, "project_{{."+CacheKeyProjectIdConstName+"}}:")
		})

		convey.Convey("CacheKeyOfProcess", func() {
			convey.So(CacheKeyOfProcess, convey.ShouldEqual, "process_{{."+CacheKeyProcessIdConstName+"}}:")
		})

		convey.Convey("CacheKeyOfRole", func() {
			convey.So(CacheKeyOfRole, convey.ShouldEqual, "role_{{."+CacheKeyRoleIdConstName+"}}:")
		})

		convey.Convey("CacheKeyOfSourceChannel", func() {
			convey.So(CacheKeyOfSourceChannel, convey.ShouldEqual, "source_channel_{{."+CacheKeySourceChannelConstName+"}}:")
		})
	})
}

func TestConstCacheKeyName(t *testing.T) {

	convey.Convey("Test CacheKeyName", t, func() {

		convey.Convey("CacheKeyOrgIdConstName", func() {
			convey.So(CacheKeyOrgIdConstName, convey.ShouldEqual, "orgId")
		})

		convey.Convey("CacheKeyUserIdConstName", func() {
			convey.So(CacheKeyUserIdConstName, convey.ShouldEqual, "userId")
		})

		convey.Convey("CacheKeyOutUserIdConstName", func() {
			convey.So(CacheKeyOutUserIdConstName, convey.ShouldEqual, "outUserId")
		})
		convey.Convey("CacheKeyProjectIdConstName", func() {
			convey.So(CacheKeyProjectIdConstName, convey.ShouldEqual, "projectId")
		})

		convey.Convey("CacheKeyIssueIdConstName", func() {
			convey.So(CacheKeyIssueIdConstName, convey.ShouldEqual, "issueId")
		})

		convey.Convey("CacheKeyObjectCodeConstName", func() {
			convey.So(CacheKeyObjectCodeConstName, convey.ShouldEqual, "objectCode")
		})

		convey.Convey("CacheKeyProcessIdConstName", func() {
			convey.So(CacheKeyProcessIdConstName, convey.ShouldEqual, "processId")
		})

		convey.Convey("CacheKeyRoleIdConstName", func() {
			convey.So(CacheKeyRoleIdConstName, convey.ShouldEqual, "roleId")
		})

		convey.Convey("CacheKeySourceChannelConstName", func() {
			convey.So(CacheKeySourceChannelConstName, convey.ShouldEqual, "sourceChannel")
		})

		convey.Convey("CacheKeyYearConstName", func() {
			convey.So(CacheKeyYearConstName, convey.ShouldEqual, "year")
		})

		convey.Convey("CacheKeyMonthConstName", func() {
			convey.So(CacheKeyMonthConstName, convey.ShouldEqual, "month")
		})

		convey.Convey("CacheKeyDayConstName", func() {
			convey.So(CacheKeyDayConstName, convey.ShouldEqual, "day")
		})
	})
}

//系统缓存
func TestConstSystemCache(t *testing.T) {

	convey.Convey("Test SystemCache", t, func() {
		//DingTalk Suite Ticket
		//convey.Convey("CacheSuiteTicketKey", func() {
		//	convey.So(CacheSuiteTicketKey, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfSys+AppOrgDingtalkChannel+":suite_ticket")
		//})
		////用户token
		//convey.Convey("CacheUserTokenExpire", func() {
		//	convey.So(CacheUserTokenExpire, convey.ShouldEqual, 60*60*24*15)
		//})
		//
		//convey.Convey("CacheUserToken", func() {
		//	convey.So(CacheUserToken, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfSys+"user:token:")
		//})
		////对象id缓存key前缀
		//convey.Convey("CacheObjectIdPreKey", func() {
		//	convey.So(CacheObjectIdPreKey, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfSys+"object_id:")
		//})
		////角色操作列表
		//convey.Convey("CacheRoleOperationList", func() {
		//	convey.So(CacheRoleOperationList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfSys+"role_operation_list")
		//})

	})
}

//组织相关缓存
func TestConstOrganizationCacheKey(t *testing.T) {

	convey.Convey("Test Organization", t, func() {
		//用户配置缓存
		//convey.Convey("CacheUserConfig", func() {
		//	convey.So(CacheUserConfig, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+CacheKeyOfUser+"config")
		//})
		//用户基础信息缓存key
		//convey.Convey("CacheBaseUserInfo", func() {
		//	convey.So(CacheBaseUserInfo, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+CacheKeyOfUser+"baseinfo")
		//})
		//组织基础信息
		//convey.Convey("CacheBaseOrgInfo", func() {
		//	convey.So(CacheBaseOrgInfo, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"baseinfo")
		//})
		//项目信息
		//convey.Convey("CacheBaseProjectInfo", func() {
		//	convey.So(CacheBaseProjectInfo, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+CacheKeyOfProject+"baseinfo")
		//})
		//项目对象类型缓存
		//convey.Convey("CacheProjectObjectTypeList", func() {
		//	convey.So(CacheProjectObjectTypeList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"project_object_type")
		//})
		//流程缓存
		//convey.Convey("CacheProcessList", func() {
		//	convey.So(CacheProcessList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"process_list")
		//})
		//流程状态列表
		//convey.Convey("CacheProcessStatusList", func() {
		//	convey.So(CacheProcessStatusList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+CacheKeyOfProcess+"process_status_list")
		//})
		//状态列表
		//convey.Convey("CacheStatusList", func() {
		//	convey.So(CacheStatusList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"process_status_list")
		//})
		////流程步骤列表
		//convey.Convey("CacheProcessStepList", func() {
		//	convey.So(CacheProcessStepList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+CacheKeyOfProcess+"process_step_list")
		//})
		//
		////优先级列表
		//convey.Convey("CachePriorityList", func() {
		//	convey.So(CachePriorityList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"priority_list")
		//})
		////项目类型
		//convey.Convey("CacheProjectTypeList", func() {
		//	convey.So(CacheProjectTypeList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"project_type_list")
		//})
		////项目类型与项目对象类型关联缓存
		//convey.Convey("CacheKeyCacheProjectTypeProjectObjectTypeListOfSourceChannel", func() {
		//	convey.So(CacheProjectTypeProjectObjectTypeList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"project_type_project_object_type")
		//})
		//
		////角色列表
		//convey.Convey("CacheRoleList", func() {
		//	convey.So(CacheRoleList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"role_list")
		//})
		//
		////用户角色列表
		//convey.Convey("CacheUserRoleList", func() {
		//	convey.So(CacheUserRoleList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+CacheKeyOfUser+"user_role_list")
		//})
		//
		////角色权限列表
		//convey.Convey("CacheRolePermissionOperationList", func() {
		//	convey.So(CacheRolePermissionOperationList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+CacheKeyOfRole+"role_permission_list")
		//})
		//
		////获取外部用户id关联的内部用户id
		//convey.Convey("CacheDingTalkOutUserIdRelationId", func() {
		//	convey.So(CacheDingTalkOutUserIdRelationId, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+CacheKeyOfSourceChannel+CacheKeyOfOutUser+"user_id")
		//})
		//
		////项目状态缓存
		//convey.Convey("CacheProjectProcessStatus", func() {
		//	convey.So(CacheProjectProcessStatus, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+CacheKeyOfProject+"process_status")
		//})
		//
		////任务来源列表
		//convey.Convey("CacheIssueSourceList", func() {
		//	convey.So(CacheIssueSourceList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"issue_source_list")
		//})
		//
		////任务类型列表
		//convey.Convey("CacheIssueObjectTypeList", func() {
		//	convey.So(CacheIssueObjectTypeList, convey.ShouldEqual, CacheKeyPrefix+CacheKeyOfOrg+"issue_object_type_list")
		//})
	})
}

func TestGetCacheBaseExpire(t *testing.T) {
	convey.Convey("Test GetCacheBaseExpire", t, func() {
		expire := GetCacheBaseExpire()
		t.Log(expire)
		expire = GetCacheBaseExpire()
		t.Log(expire)
	})
}
