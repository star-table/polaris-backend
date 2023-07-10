package permissionvo

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo/appauth"
)

type UpdateLcPermissionGroupReq struct {
	OrgId     int64    `json:"orgId"`
	NewUserId int64    `json:"newUserId"`
	Id        string   `json:"id"`
	Key       string   `json:"key"`
	Values    []string `json:"values"`
}

type UpdateLcPermissionGroupResp struct {
	vo.Err
	Data bool `json:"data"`
}

type UpdateLcAppPermissionGroupOpConfigReq struct {
	AppId     int64  `json:"appId"`
	LangCode  int    `json:"langCode"`
	OptAuth   string `json:"optAuth"`
	TableAuth string `json:"tableAuth"`
}

type UpdateLcAppPermissionGroupOpConfigResp struct {
	vo.Err
	Data bool `json:"data"`
}

type InitDefaultManageGroupReq struct {
	OrgID       int64               `json:"orgId"`
	AuthOptions []OptAuthOptionInfo `json:"authOptions"`
}

type InitDefaultManageGroupResp struct {
	vo.Err
	Data *InitDefaultManageGroupRespData `json:"data"`
}

type InitDefaultManageGroupRespData struct {
	// 系统管理组id
	SysGroupID int64 `json:"sysGroupId,string"`
}

type OptAuthOptionInfo struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Group    string `json:"group"`
	Required bool   `json:"required"`
	IsMenu   bool   `json:"isMenu"`
	Status   int    `json:"status"`
}

type FieldAuthOptionInfo struct {
	Code     int    `json:"code"`
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

type InitAppPermissionReq struct {
	OrgID                      int64                 `json:"orgId"`
	AppPackageID               int64                 `json:"appPackageId"`
	AppID                      int64                 `json:"appId"`
	AppType                    int                   `json:"appType"`
	OptAuthOptions             []OptAuthOptionInfo   `json:"optAuthOptions"`
	FieldAuthOptions           []FieldAuthOptionInfo `json:"fieldAuthOptions"`
	IsExt                      bool                  `json:"isExt"`
	ComponentType              string                `json:"componentType"`
	Creatable                  bool                  `json:"creatable"`
	UserID                     int64                 `json:"userId"`
	Config                     string                `json:"config"`
	DefaultPermissionGroupType int                   `json:"defaultPermissionGroupType"` //默认权限组类型（为空则不初始化），1 初始化表单的权限组，2 初始化仪表盘的权限组，3 初始化极星项目的权限组
}

type InitAppPermissionResp struct {
	vo.Err
	Data bool `json:"data"`
}

type CreateLessCodeAppReq struct {
	AppType      *int    `json:"appType"`
	Name         *string `json:"name"`
	OrgId        *int64  `json:"orgId"`
	UserId       *int64  `json:"userId"`
	PkgId        int64   `json:"pkgId"`
	Config       string  `json:"config"`
	ExtendsId    int64   `json:"extendsId"`
	ProjectId    int64   `json:"projectId"`
	ParentId     int64   `json:"parentId"`
	Hidden       int     `json:"hidden"`
	AuthType     int     `json:"authType"`
	Icon         string  `json:"icon"`
	ExternalApp  int     `json:"externalApp"`
	LinkUrl      string  `json:"linkUrl"`
	MirrorViewId int64   `json:"mirrorViewId"`
	MirrorAppId  int64   `json:"mirrorAppId"`
	AddAllMember bool    `json:"addAllMember"`
	ProjectType  int     `json:"projectType"`
}

type CreateLessCodeAppReqConfig struct {
	Fields string `json:"fields"`
}

type CreateLessCodeAppResp struct {
	vo.Err
	Data *CreateLessCodeAppRespData `json:"data"`
}

type CreateLessCodeAppRespData struct {
	Id      int64  `json:"id,string"`
	Name    string `json:"name"`
	OrgId   string `json:"orgId"`
	Type    int    `json:"type"`
	Creator int64  `json:"creator,string"`
}

// 创建或更新管理组请求结构体
type CreateOrUpdateManageGroupReq struct {
	ID      *int64  `json:"id"`
	UserIds *string `json:"userIds"`
}

// 创建或更新管理组响应结构体
type CreateOrUpdateManageGroupResp struct {
	vo.Err
	Data *CreateOrUpdateManageGroupRespData `json:"data"`
}

type CreateOrUpdateManageGroupRespData struct {
	ID       int64  `json:"id,string"`
	OrgID    int64  `json:"orgId,string"`
	Name     string `json:"name"`
	LangCode string `json:"langCode"`
}

// 管理认证信息响应结构体
type ManageAuthInfoResp struct {
	vo.Err
	Data *ManageAuthInfoRespData `json:"data"`
}

type ManageAuthInfoRespData struct {
	OrgID         int64    `json:"orgId,string"`
	IsAdmin       bool     `json:"admin"` // 无码系统过来时，字段是 admin
	AppPackageIds []string `json:"appPackageIds"`
	AppIds        []string `json:"appIds"`
	OptAuth       []string `json:"optAuth"`
	DeptIds       []string `json:"deptIds"`
	RoleIds       []string `json:"roleIds"`
}

//type GetOptAuthListResp struct {
//	vo.Err
//	NewData *GetOptAuthListRespData `json:"data"`
//}

type GetOptAuthListResp struct {
	vo.Err
	Data []string `json:"data"`
}

type GetOptAuthListRespData struct {
	OptList []string `json:"optList"`
}

type GetLcPermissionGroupTreeReq struct {
	UserId int64 `json:"userId"`
	OrgId  int64 `json:"orgId"`
}

type GetLcPermissionGroupTreeResp struct {
	vo.Err
	Data *GetLcPermissionGroupTreeRespData `json:"data"`
}

type GetLcPermissionGroupTreeRespData struct {
	SysGroup      *GetLcPermissionGroupTreeRespDataItem   `json:"sysGroup"`
	GeneralGroups []*GetLcPermissionGroupTreeRespDataItem `json:"generalGroups"`
}

type GetLcPermissionGroupTreeRespDataItem struct {
	Id       string `json:"id"`
	OrgId    int64  `json:"orgId"`
	Name     string `json:"name"`
	LangCode string `json:"langCode"`
}

type GetLcPermissionAdminGroupDetailResp struct {
	vo.Err
	Data *GetLcPermissionAdminGroupDetailRespData `json:"data"`
}

type GetLcPermissionAdminGroupDetailRespData struct {
	AdminGroup *GetLcPermissionAdminGroupDetailRespDataAdminGroup `json:"adminGroup"`
}

type GetLcPermissionAdminGroupDetailRespDataAdminGroup struct {
	Id      string   `json:"id"`
	Name    string   `json:"name"`
	Type    int      `json:"type"`
	UserIds []string `json:"userIds"`
}

type GetAppAuthResp struct {
	vo.Err
	Data appauth.GetAppAuthData `json:"data"`
}

type GetAppRoleListResp struct {
	vo.Err
	Data []AppRoleInfo `json:"data"`
}

type AppRoleInfo struct {
	Id        int64             `json:"id,string"`
	HasDelete bool              `json:"hasDelete"`
	HasEdit   bool              `json:"hasEdit"`
	LangCode  string            `json:"langCode"`
	Name      string            `json:"name"`
	ReadOnly  string            `json:"readOnly"`
	Remake    string            `json:"remake"`
	Members   AppRoleInfoMember `json:"members"`
}

type AppRoleInfoMember struct {
	Depts interface{}             `json:"depts"`
	Roles interface{}             `json:"roles"`
	Users []AppRoleInfoMemberUser `json:"users"`
}

type AppRoleInfoMemberUser struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GetDataAuthBatchReq struct {
	AppId   int64   `json:"appId"`
	OrgId   int64   `json:"orgId"`
	UserId  int64   `json:"userId"`
	DataIds []int64 `json:"dataIds"`
}

type GetDataAuthBatchResp struct {
	vo.Err
	Data map[int64]appauth.GetAppAuthData `json:"data"`
}

type GetAppAuthBatchReq struct {
	OrgId  int64   `json:"orgId"`
	UserId int64   `json:"userId"`
	AppIds []int64 `json:"appIds"`
}

type GetAppAuthBatchResp struct {
	vo.Err
	Data map[int64]appauth.GetAppAuthData `json:"data"`
}

type GetUserGroupMappingsResp struct {
	vo.Err
	Data GroupMappings `json:"data"`
}

type GroupMappings struct {
	UserGroupMappings map[string][]AppRoleInfo `json:"userGroupMappings"`
	DeptGroupMappings map[string][]AppRoleInfo `json:"deptGroupMappings"`
	RoleGroupMappings map[string][]AppRoleInfo `json:"roleGroupMappings"`
}

type GetPermissionUpdateTimeReq struct {
	AppId int64 `json:"appId"`
}

type GetPermissionUpdateTimeResp struct {
	vo.Err
	Data *GetPermissionUpdateTimeData `json:"data"`
}

type GetPermissionUpdateTimeData struct {
	UpdateTime int64 `json:"updateTime"`
}
