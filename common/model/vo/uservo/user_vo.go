package uservo

import (
	"time"

	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"

	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type GetUserAuthorityResp struct {
	vo.Err
	Data GetUserAuthorityData `json:"data"`
}

type GetUserAuthorityData struct {
	OrgId                int64    `json:"orgId"`
	UserId               int64    `json:"userId"`
	IsOrgOwner           bool     `json:"isOrgOwner"`
	IsSysAdmin           bool     `json:"isSysAdmin"`
	IsSubAdmin           bool     `json:"isSubAdmin"`
	OptAuth              []string `json:"optAuth"`
	HasDeptOptAuth       bool     `json:"hasDeptOptAuth"`
	HasRoleDeptAuth      bool     `json:"hasRoleDeptAuth"`
	HasAppPackageOptAuth bool     `json:"hasAppPackageOptAuth"`
	// 目前以下字段极星中尚未使用，为了避免误解或歧义，暂时注释，如有需要需到 usercenter 对应的方法中进行实现。
	RefDeptIds []int64 `json:"refDeptIds"`
	//RefRoleIds           []int64  `json:"refRoleIds"`
	//ManageAppPackages    []int64  `json:"manageAppPackages"`
	ManageApps []int64 `json:"manageApps"`
	//ManageRoles          []int64  `json:"manageRoles"`
	//ManageDepts          []int64  `json:"manageDepts"`
}

type GetUserAuthorityReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type InitDefaultManageGroupResp struct {
	vo.Err
	Data *InitDefaultManageGroupRespData `json:"data"`
}

type InitDefaultManageGroupRespData struct {
	// 系统管理组id
	SysGroupID int64 `json:"sysGroupId,string"`
}

type AddUserToSysManageGroupReq struct {
	UserIds []int64 `json:"userIds"`
}

type AddUserToSysManageGroupResp struct {
	vo.Err
	Data bool `json:"data"`
}

type AddScriptTaskReq struct {
	AppId    string      `json:"appId"`
	TaskId   string      `json:"taskId"`
	TaskName string      `json:"taskName"`
	CronSpec string      `json:"cronSpec"`
	ParamMap interface{} `json:"paramMap"` // 使用时，使用 map[string]string 类型
}

type AddScriptTaskResp struct {
	vo.Err
	Data vo.BoolResp `json:"data"`
}

type GetUsersCouldManageResp struct {
	vo.Err
	Data GetUsersCouldManageRespData `json:"data"`
}

type GetUsersCouldManageRespData struct {
	List []bo.SimpleUserInfoBo `json:"list"`
}

type CheckAndSetSuperAdminResp struct {
	vo.Err
	Data bool `json:"data"`
}

type GetAllDeptByIdsResp struct {
	vo.Err
	Data []GetAllDeptByIdsRespDataUser `json:"data"`
}

type GetAllDeptByIdsRespDataUser struct {
	Id            int64  `json:"id"`            // Id 部门ID
	OrgId         int64  `json:"orgId"`         // OrgId 组织ID
	Name          string `json:"name"`          // Name 部门名称
	Code          string `json:"code"`          // Code
	ParentId      int64  `json:"parentId"`      // ParentId 父部门ID
	Path          string `json:"path"`          // Path
	Sort          int    `json:"sort"`          // Sort
	IsHide        int    `json:"isHide"`        // IsHide
	SourceChannel string `json:"sourceChannel"` // SourceChannel
	Status        int    `json:"status"`        // Status
	IsDelete      int    `json:"isDelete"`      // IsDelete
}

type GetAllUserByIdsResp struct {
	vo.Err
	Data []GetAllUserByIdsRespDataUser `json:"data"`
}

type GetAllUserByIdsRespDataUser struct {
	Id          int64     `json:"id"`          // Id 用户ID
	LoginName   string    `json:"loginName"`   // LoginName
	Name        string    `json:"name"`        // Name 姓名
	NamePy      string    `json:"namePy"`      // NamePy 姓名拼音
	Avatar      string    `json:"avatar"`      // Avatar 用户头像
	Email       string    `json:"email"`       // Email 邮箱
	PhoneRegion string    `json:"phoneRegion"` // PhoneRegion 手机区号
	PhoneNumber string    `json:"phoneNumber"` // PhoneNumber 手机
	Status      int       `json:"status"`      // Status 状态1启用2禁用3离职
	Creator     int64     `json:"creator"`     // Creator 组织成员 创建者
	CreateTime  time.Time `json:"createTime"`  // CreateTime 创建时间
	Updator     int64     `json:"updator"`     // Updator 组织成员 最后修改者ID
	UpdateTime  time.Time `json:"updateTime"`  // UpdateTime 创建时间
	IsDelete    int       `json:"isDelete"`    // IsDelete
	Type        int       `json:"type"`        // Type 用户类型，1：内部，2：外部
	UserBindDeptAndRoleResp
}

type PerGroup struct {
	Id        string `json:"id,omitempty"`
	HasDelete bool   `json:"hasDelete,omitempty"`
	HasEdit   bool   `json:"hasEdit,omitempty"`
	LangCode  string `json:"langCode,omitempty"`
	Name      string `json:"name,omitempty"`
	ReadOnly  int32  `json:"readOnly,omitempty"`
	Remake    string `json:"remake,omitempty"`
}

type MemberDept struct {
	Id        int64                       `json:"id,string"`           // Id ID
	Name      string                      `json:"name"`                // Name 名字
	Avatar    string                      `json:"avatar"`              // Avatar 头像
	Type      string                      `json:"type"`                // Type "U_" "D_"
	Status    int                         `json:"status"`              // Status 状态1启用2禁用3离职
	IsDelete  int                         `json:"isDelete"`            // IsDelete
	PerGroups []*permissionvo.AppRoleInfo `json:"perGroups,omitempty"` // 权限组
}

type UserBindDeptAndRoleResp struct {
	DeptList []UserDeptBindData `json:"deptList"` // DeptList 部门列表
}

type UserDeptBindData struct {
	DepartmentId             int64  `json:"departmentId"`             // DepartmentId 部门ID
	DepartmentName           string `json:"departmentName"`           // DepartmentName 部门名称
	OutOrgDepartmentId       string `json:"outOrgDepartmentId"`       // OutOrgDepartmentId
	OutOrgDepartmentCode     string `json:"outOrgDepartmentCode"`     // OutOrgDepartmentCode
	OutOrgDepartmentParentId string `json:"outOrgDepartmentParentId"` // OutOrgDepartmentParentId
	IsLeader                 int    `json:"isLeader"`                 // IsLeader 是否是部门负责人 1是 2否
	PositionId               int64  `json:"positionId"`               // PositionId 注意此处是org内的职级ID
	PositionName             string `json:"positionName"`             // PositionName 职级名称
	PositionLevel            int    `json:"positionLevel"`            // PositionLevel 职级等级
}

type GetFullDeptNameArrByIdsResp struct {
	vo.Err
	Data map[int64][]string `json:"data"`
}

type GetDeptListResp struct {
	vo.Err
	Data []GetDeptListRespDataItem `json:"data"`
}

type GetDeptListRespDataItem struct {
	Id            int64  `json:"id"`
	OrgId         int64  `json:"orgId"`
	Name          string `json:"name"`
	Code          string `json:"code"`
	ParentId      int64  `json:"parentId"`
	Path          string `json:"path"`
	Sort          int    `json:"sort"`
	IsHide        int    `json:"isHide"`
	SourceChannel string `json:"sourceChannel"`
	Status        int    `json:"status"`
	IsDelete      int    `json:"isDelete"` // IsDelete
}

type AddNewMenuToRoleReqInput struct {
	OrgIds []int64 `json:"orgIds"`
	// 待追加的权限项标识
	PrmItems []string `json:"prmItems"`
}

type BoolResp struct {
	vo.Err
	Data bool `json:"data"`
}

type GetCommAdminMangeAppsResp struct {
	vo.Err
	Data []*GetCommAdminMangeAppsData `json:"data"`
}

type GetCommAdminMangeAppsData struct {
	UserId int64   `json:"userId"`
	AppIds []int64 `json:"appIds"` // [-1] 表示可以管理全部应用
}

type GetUserIdsByDeptIdsReq struct {
	// 组织id
	OrgId int64 `json:"orgId"`
	// 部门id
	DeptIds []int64 `json:"deptIds"`
}

type GetUserIdsByDeptIdsData struct {
	UserIds []int64 `json:"userIds"`
}

type GetUserIdsByDeptIdsResp struct {
	vo.Err
	Data *GetUserIdsByDeptIdsData `json:"data"`
}
