package orgvo

import (
	"time"

	"github.com/star-table/common/core/types"

	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type PersonalInfoRespVo struct {
	vo.Err
	PersonalInfo *vo.PersonalInfo `json:"data"`
}

type PersonalInfoRestRespVo struct {
	vo.Err
	Data *PersonalInfo `json:"data"`
}

type PersonalInfoReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	UserId        int64  `json:"userId"`
	OrgId         int64  `json:"orgId"`
}

type GetUserIdsReqVo struct {
	SourceChannel string       `json:"sourceChannel"`
	EmpIdsBody    EmpIdsBodyVo `json:"empIdsBody"`
	CorpId        string       `json:"corpId"`
	OrgId         int64        `json:"orgId"`
}

type EmpIdsBodyVo struct {
	EmpIds []string `json:"empIds"`
}

type GetUserIdsRespVo struct {
	vo.Err
	GetUserIds []*vo.UserIDInfo `json:"data"`
}

type GetUserIdReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	EmpId         string `json:"empId"`
	CorpId        string `json:"corpId"`
	OrgId         int64  `json:"orgId"`
}

type GetUserIdRespVo struct {
	vo.Err
	GetUserId *vo.UserIDInfo `json:"data"`
}

type UserConfigInfoRespVo struct {
	vo.Err
	UserConfigInfo *vo.UserConfig `json:"data"`
}

type UserConfigInfoReqVo struct {
	UserId int64 `json:"userId"`
	OrgId  int64 `json:"orgId"`
}

type UpdateUserConfigRespVo struct {
	vo.Err
	UpdateUserConfig *vo.UpdateUserConfigResp `json:"data"`
}

type UpdateUserConfigReqVo struct {
	UpdateUserConfigReq vo.UpdateUserConfigReq `json:"updateUserConfigReq"`
	OrgId               int64                  `json:"orgId"`
	UserId              int64                  `json:"userId"`
}

type UpdateUserPcConfigReqVo struct {
	UpdateUserPcConfigReq vo.UpdateUserPcConfigReq `json:"updateUserConfigPcReq"`
	OrgId                 int64                    `json:"orgId"`
	UserId                int64                    `json:"userId"`
}

type UpdateUserInfoReqVo struct {
	UpdateUserInfoReq vo.UpdateUserInfoReq `json:"updateUserInfoReq"`
	OrgId             int64                `json:"orgId"`
	UserId            int64                `json:"userId"`
}

type UpdateUserDefaultProjectIdConfigReqVo struct {
	UpdateUserDefaultProjectIdConfigReq vo.UpdateUserDefaultProjectConfigReq `json:"updateUserDefaultProjectIdConfigReq"`
	OrgId                               int64                                `json:"orgId"`
	UserId                              int64                                `json:"userId"`
}

type VerifyOrgReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type VerifyOrgUsersReqVo struct {
	OrgId int64                 `json:"orgId"`
	Input VerifyOrgUsersReqData `json:"input"`
}

type VerifyOrgUsersReqData struct {
	UserIds []int64 `json:"userIds"`
}

type GetUserInfoReqVo struct {
	OrgId         int64  `json:"orgId"`
	UserId        int64  `json:"userId"`
	SourceChannel string `json:"sourceChannel"`
}

type GetUserInfoRespVo struct {
	vo.Err
	UserInfo *bo.UserInfoBo `json:"data"`
}

type GetOrgBoListRespVo struct {
	vo.Err
	OrganizationBoList []bo.OrganizationBo `json:"data"`
}

type GetOrgBoListByPageRespVo struct {
	vo.Err
	Data GetOrgBoListByPageRespData `json:"data"`
}

type GetOrgBoListByPageRespData struct {
	List  []bo.OrganizationBo `json:"list"`
	Total int64               `json:"total"`
	Page  int                 `json:"page"`
	Size  int                 `json:"size"`
}

type CreateOrgRespVo struct {
	vo.Err
	Data CreateOrgRespVoData `json:"data"`
}

type CreateOrgRespVoData struct {
	OrgId int64 `json:"orgId"`
}

type GetOutUserInfoListBySourceChannelReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	Page          int    `json:"page"`
	Size          int    `json:"size"`
}

type GetOutUserInfoListBySourceChannelRespVo struct {
	vo.Err
	UserOutInfoBoList []bo.UserOutInfoBo `json:"data"`
}

type GetOutUserInfoListByUserIdsReqVo struct {
	UserIds []int64 `json:"userIds"`
}

type GetOutUserInfoListByUserIdsRespVo struct {
	vo.Err
	UserOutInfoBoList []bo.UserOutInfoBo `json:"data"`
}

type GetUserInfoListReqVo struct {
	OrgId int64 `json:"orgId"`
}

type GetUserInfoListRespVo struct {
	vo.Err
	SimpleUserInfo []bo.SimpleUserInfoBo `json:"data"`
}

type InitExistOrgReqVo struct {
	Input  vo.InitExistOrgReq `json:"input"`
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
}

type SaveOrgSummaryTableAppIdReqVo struct {
	Input  SaveOrgSummaryTableAppIdReqVoData `json:"input"`
	OrgId  int64                             `json:"orgId"`
	UserId int64                             `json:"userId"`
}

type VoidRespVo struct {
	vo.Err
	Data VoidRespVoData `json:"data"`
}

type VoidRespVoData struct {
	IsOk bool `json:"isOk"`
}

type SaveOrgSummaryTableAppIdReqVoData struct {
	AppId             int64 `json:"appId"` // 对应于 OrgSummaryTableAppId
	TagAppId          int64 `json:"tagAppId"`
	PriorityAppId     int64 `json:"priorityAppId"`
	IssueBarAppId     int64 `json:"issueBarAppId"`
	SceneAppId        int64 `json:"sceneAppId"`
	IssueStatusAppId  int64 `json:"issueStatusAppId"`
	EmptyProjectAppId int64 `json:"emptyProjectAppId"`
	// 项目迁移到无码中需要的，用于存储组织的项目数据
	ProjectFormAppId int64 `json:"projectFormAppId"`
	// 存放项目视图的目录 id
	ProjectFolderAppId int64 `json:"projectFolderAppId"`
	// 存放任务视图的目录 id
	IssueFolderAppId int64 `json:"issueFolderAppId"`
}

type InitExistOrgRespVo struct {
	vo.Err
	Data *InitExistOrgRespData `json:"data"`
}

type InitExistOrgRespData struct {
	IsOk    bool   `json:"isOk"`
	StrInfo string `json:"strInfo"`
}

type CreateOrgReqVo struct {
	Data   CreateOrgReqVoData `json:"data"`
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
}

type CreateOrgReqVoData struct {
	CreatorId        int64           `json:"creatorId"`
	CreateOrgReq     vo.CreateOrgReq `json:"createOrgReq"`
	UserToken        string          `json:"userToken"`
	ImportSampleData int             `json:"importSampleData"`
}

type GetUserInfoByUserIdsReqVo struct {
	UserIds []int64 `json:"userIds"`
	OrgId   int64   `json:"orgId"`
}

type GetUserInfoByUserIdsListRespVo struct {
	vo.Err
	GetUserInfoByUserIdsRespVo *[]GetUserInfoByUserIdsRespVo `json:"data"`
}

type GetUserInfoByUserIdsRespVo struct {
	UserId        int64  `json:"userId"`
	OutUserId     string `json:"outUserId"` //有可能为空
	OrgId         int64  `json:"orgId"`
	OutOrgId      string `json:"outOrgId"` //有可能为空
	Name          string `json:"name"`
	Avatar        string `json:"avatar"`
	HasOutInfo    bool   `json:"hasOutInfo"`
	HasOrgOutInfo bool   `json:"hasOrgOutInfo"`

	OrgUserIsDelete    int `json:"orgUserIsDelete"`    //是否被组织移除
	OrgUserStatus      int `json:"orgUserStatus"`      //用户组织状态
	OrgUserCheckStatus int `json:"orgUserCheckStatus"` //用户组织审核状态
}

type ExportAddressListReqVo struct {
	Data          ExportAddressListReq `json:"data"`
	CurrentUserId int64                `json:"userId"`
	OrgId         int64                `json:"orgId"`
}

type ExportAddressListReq struct {
	//搜索字段
	SearchCode *string `json:"searchCode"`
	//是否已分配部门（1已分配2未分配，默认全部）
	IsAllocate *int `json:"isAllocate"`
	//是否禁用（1启用2禁用，默认全部）
	Status *int `json:"status"`
	//角色id
	RoleId *int64 `json:"roleId"`
	//部门id
	DepartmentId *int64 `json:"departmentId"`
	//导出字段(name,mobile,email,department,role,isLeader,statusChangeTime,createTime)
	ExportField []string `json:"exportField"`
}

type ExportAddressListRespVo struct {
	vo.Err
	Data ExportAddressListResp `json:"data"`
}

type ExportAddressListResp struct {
	Url string `json:"url"`
}

type UserListResp struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	List []*UserInfo `json:"list"`
}

type UserInfo struct {
	// id
	UserID int64 `json:"userId"`
	// 姓名
	Name string `json:"name"`
	// 姓名拼音
	NamePy string `json:"namePy"`
	// 用户头像
	Avatar string `json:"avatar"`
	// 邮箱
	Email string `json:"email"`
	// 手机
	PhoneNumber string `json:"phoneNumber"`
	// 用户部门信息
	DepartmentList []UserDepartmentData `json:"departmentList"`
	// 角色信息
	RoleList []UserRoleData `json:"roleList"`
	// 创建时间
	CreateTime time.Time `json:"createTime"`
	// 状态1启用2禁用
	Status int `json:"status"`
	// 禁用时间
	StatusChangeTime time.Time `json:"statusChangeTime"`
	// 是否是组织创建人
	IsCreator bool `json:"isCreator"`
}

type UserDepartmentData struct {
	//部门id
	DepartmentId int64 `json:"departmentId"`
	//是否是主管1是2否
	IsLeader int `json:"isLeader"`
	//部门名称
	DeparmentName string `json:"deparmentName"`
}

type UserRoleData struct {
	//角色id
	RoleId int64 `json:"roleId"`
	//角色名称
	RoleName string `json:"roleName"`
}

type UserListReqVo struct {
	Data          UserListReq `json:"data"`
	CurrentUserId int64       `json:"userId"`
	OrgId         int64       `json:"orgId"`
}

type UserListReq struct {
	//搜索字段
	SearchCode *string `json:"searchCode"`
	//是否已分配部门（1已分配2未分配，默认全部）
	IsAllocate *int `json:"isAllocate"`
	//是否禁用（1启用2禁用，默认全部）
	Status *int `json:"status"`
	//角色id
	RoleId *int64 `json:"roleId"`
	//部门id
	DepartmentId *int64 `json:"departmentId"`
	// 页码（选填，不填为全部）
	Page int `json:"page"`
	// 数量（选填，不填为全部）
	Size int `json:"size"`
}

type EmptyUserReqVo struct {
	Data          EmptyUserReq `json:"data"`
	CurrentUserId int64        `json:"userId"`
	OrgId         int64        `json:"orgId"`
}

//清空用户
type EmptyUserReq struct {
	// 状态,  1可用,2禁用
	Status int `json:"status"`
}

type CreateUserReqVo struct {
	Data          CreateUserReq `json:"data"`
	CurrentUserId int64         `json:"userId"`
	OrgId         int64         `json:"orgId"`
}

type CreateUserReq struct {
	// 手机号(必填)
	PhoneNumber string `json:"phoneNumber"`
	// 邮箱(必填)
	Email string `json:"email"`
	// 姓名(必填)
	Name string `json:"name"`
	// 部门id
	DepartmentIds []int64 `json:"departmentIds"`
	// 角色id
	RoleIds []int64 `json:"roleIds"`
	// 状态（1启用2禁用, 默认启用）
	Status *int `json:"status"`
}

type UpdateUserReqVo struct {
	Data          UpdateUserReq `json:"data"`
	CurrentUserId int64         `json:"userId"`
	OrgId         int64         `json:"orgId"`
}

type UpdateUserReq struct {
	// 用户id(必填)
	UserId int64 `json:"userId"`
	// 手机号(选填)
	PhoneNumber *string `json:"phoneNumber"`
	// 邮箱(选填)
	Email *string `json:"email"`
	// 姓名(选填)
	Name *string `json:"name"`
	// 部门id（选填）
	DepartmentIds *[]int64 `json:"departmentIds"`
	// 角色id（选填）
	RoleIds *[]int64 `json:"roleIds"`
	// 状态（1启用2禁用 选填）
	Status *int `json:"status"`
}

type UserListRespVo struct {
	vo.Err
	Data *UserListResp `json:"data"`
}

type GetUserListWithCreateTimeRangeReqVo struct {
	Data *GetUserListWithCreateTimeRangeReq `json:"data"`
}

type GetUserListWithCreateTimeRangeReq struct {
	SourceChannel string `json:"sourceChannel"`
	CreateTime1   string `json:"crateTime1"`
	CreateTime2   string `json:"crateTime2"`
	Page          int    `json:"page"`
	Size          int    `json:"size"`
}

type GetUserListWithCreateTimeRangeRespVo struct {
	vo.Err
	Data *GetUserListWithCreateTimeRangeResp `json:"data"`
}

type GetUserListWithCreateTimeRangeResp struct {
	UserList []*GetUserListWithCreateTimeRangeItem `json:"userList"`
}

type GetUserListWithCreateTimeRangeItem struct {
	// 用户id(必填)
	UserId int64 `json:"userId"`
	// 姓名(选填)
	Name string `json:"name"`
	// fs openId
	OpenId string `json:"openId"`
	// fs openId
	OrgId int64 `json:"orgId"`
	// fs out orgId 关联的外部的组织id
	OutOrgId string `json:"outOrgId"`
	// 用户的创建时间
	CreateTime string `json:"createTime"`
}

type SearchUserReqVo struct {
	Data          SearchUserReq `json:"data"`
	CurrentUserId int64         `json:"userId"`
	OrgId         int64         `json:"orgId"`
}

type SearchUserReq struct {
	Email string `json:"email"`
}

type SearchUserRespVo struct {
	vo.Err
	Data *SearchUserResp `json:"data"`
}

type SearchUserResp struct {
	//搜索结果（1可邀请2已邀请3已添加）
	Status   int       `json:"status"`
	UserInfo *UserInfo `json:"userInfo"`
}

type InviteUserReqVo struct {
	CurrentUserId int64         `json:"userId"`
	OrgId         int64         `json:"orgId"`
	Data          InviteUserReq `json:"data"`
}

type InviteUserListReqVo struct {
	CurrentUserId int64             `json:"userId"`
	OrgId         int64             `json:"orgId"`
	Data          InviteUserListReq `json:"data"`
}

type InviteUserListReq struct {
	// 页码（选填，不填为全部）
	Page int `json:"page"`
	// 数量（选填，不填为全部）
	Size int `json:"size"`
}

type InviteUserReq struct {
	Data []*vo.InviteUserData `json:"data"`
}

type InviteUserData struct {
	//邮箱
	Email string `json:"email"`
	//姓名（再次邀请时不用传了）
	Name string `json:"name"`
}

type InviteUserRespVo struct {
	vo.Err
	Data *InviteUserResp `json:"data"`
}

type InviteUserResp struct {
	//成功的邮箱
	SuccessEmail []string `json:"successEmail"`
	//已邀请的邮箱
	InvitedEmail []string `json:"invitedEmail"`
	//已经是用户的邮箱
	IsUserEmail []string `json:"isUserEmail"`
	//不符合规范的邮箱
	InvalidEmail []string `json:"invalidEmail"`
}

type InviteUserListRespVo struct {
	vo.Err
	Data *InviteUserListResp `json:"data"`
}

type InviteUserListResp struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	List []InviteUserInfo `json:"list"`
}

type InviteUserInfo struct {
	//用户id
	Id int64 `json:"id"`
	//名称
	Name string `json:"name"`
	//邮箱
	Email string `json:"email"`
	//邀请时间
	InviteTime time.Time `json:"inviteTime"`
	//是否24h内已邀请
	IsInvitedRecent bool `json:"isInvitedRecent"`
}

type RemoveInviteUserReqVo struct {
	CurrentUserId int64               `json:"userId"`
	OrgId         int64               `json:"orgId"`
	Data          RemoveInviteUserReq `json:"data"`
}

type RemoveInviteUserReq struct {
	//要移除的id
	Ids []int64 `json:"ids"`
	//是否删除全部(1是，其余为否)
	IsAll int `json:"isAll"`
}

type UserStatReqVo struct {
	CurrentUserId int64 `json:"userId"`
	OrgId         int64 `json:"orgId"`
}

type UserStatRespVo struct {
	vo.Err
	Data *UserStatResp `json:"data"`
}

type UserStatResp struct {
	//所有成员数量
	AllCount int64 `json:"allCount"`
	//未分配成员数量
	UnallocatedCount int64 `json:"unallocatedCount"`
	//未接受邀请成员数量
	UnreceivedCount int64 `json:"unreceivedCount"`
	//已禁用成员数量
	ForbiddenCount int64 `json:"forbiddenCount"`
}

// 从飞书同步用户、部门信息
type SyncUserInfoFromFeiShuReqVo struct {
	CurrentUserId             int64                        `json:"userId"`
	OrgId                     int64                        `json:"orgId"`
	SyncUserInfoFromFeiShuReq vo.SyncUserInfoFromFeiShuReq `json:"syncUserInfoFromFeiShuReq"`
}

type GetUserInfoByFeishuTenantKeyData struct {
	OrgId  int64
	UserId int64
}

type GetUserInfoByFeishuTenantKeyResp struct {
	vo.Err
	Data GetUserInfoByFeishuTenantKeyData `json:"data"`
}

type GetUserInfoByFeishuTenantKeyReq struct {
	TenantKey string `json:"tenantKey"`
	OpenId    string `json:"openId"`
}

type GetOrgUserIdsByEmIdsResp struct {
	vo.Err
	Data map[string]int64 `json:"data"`
}

type GetOrgUserIdsByEmIdsReq struct {
	OrgId         int64    `json:"orgId"`
	SourceChannel string   `json:"sourceChannel"`
	EmpIds        []string `json:"empIds"`
}

type FilterResignedUserIdsReqVo struct {
	OrgId   int64   `json:"orgId"`
	UserIds []int64 `json:"userIds"`
}

type FilterResignedUserIdsResp struct {
	vo.Err
	Data []int64 `json:"data"`
}

type GetOrgUserIdsReq struct {
	OrgId int64 `json:"orgId"`
}

type GetOrgUserIdsResp struct {
	vo.Err
	Data []int64 `json:"data"`
}

type GetPayRemindReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetPayRemindResp struct {
	vo.Err
	Data *vo.GetPayRemindResp `json:"data"`
}

type GetOrgUsersInfoByEmIdsReq struct {
	OrgId         int64    `json:"orgId"`
	SourceChannel string   `json:"sourceChannel"`
	EmpIds        []string `json:"empIds"`
}

type GetOrgUserAndDeptCountReq struct {
	OrgId int64 `json:"orgId"`
}

type GetOrgUserAndDeptCountResp struct {
	vo.Err
	Data *GetOrgUserAndDeptCountRespData `json:"data"`
}

type GetOrgUserAndDeptCountRespData struct {
	UserCount uint64 `json:"userCount"`
	DeptCount uint64 `json:"deptCount"`
}

type SetVisitUserGuideStatusReq struct {
	OrgId  int64  `json:"orgId"`
	UserId int64  `json:"userId"`
	Flag   string `json:"flag"`
}

type SetVersionReq struct {
	OrgId  int64          `json:"orgId"`
	UserId int64          `json:"userId"`
	Input  SetVersionData `json:"input"`
}

type SetVersionData struct {
	//Version string `json:"version"`
	VersionInfoVisible bool `json:"versionInfoVisible"`
}

type SetUserActivityReq struct {
	OrgId  int64           `json:"orgId"`
	UserId int64           `json:"userId"`
	Input  SetActivityData `json:"input"`
}

type SetActivityData struct {
	ActivityFlag int `json:"activityFlag"`
}

type GetVersionReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetVersionRespVo struct {
	vo.Err
	Data *VersionResp `json:"data"`
}

type VersionResp struct {
	Version GetVersionData `json:"version"`
}

type GetVersionData struct {
	VersionInfoVisible bool `json:"versionInfoVisible"`
	//Version string `json:"version"`
}

type SaveViewLocationReqVo struct {
	OrgId  int64                `json:"orgId"`
	UserId int64                `json:"userId"`
	Input  *SaveViewLocationReq `json:"input"`
}

type SaveViewLocationReq struct {
	AppId       string `json:"appId"`
	ProjectId   int64  `json:"projectId"`
	TableId     string `json:"tableId"`
	IterationId int64  `json:"iterationId"`
	ViewId      string `json:"viewId"`
	MenuId      string `json:"menuId"`
	DashboardId string `json:"dashboardId"`
}

type SaveViewLocationRespVo struct {
	vo.Err
	Data bool `json:"data"`
}

type GetViewLocationReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetViewLocationRespVo struct {
	vo.Err
	Data []*UserLastViewLocationData `json:"data"`
}

type UpdateUserViewLocationReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
	AppId  int64 `json:"appId"`
}

type UserLastViewLocationData struct {
	AppId       string `json:"appId"`
	ProjectId   int64  `json:"projectId"`
	TableId     string `json:"tableId,omitempty"`
	ViewId      string `json:"viewId,omitempty"`
	IterationId int64  `json:"iterationId,omitempty"`
	MenuId      string `json:"menuId,omitempty"`
	DashboardId string `json:"dashboardId,omitempty"`
}

type GetUserOrDeptSameNameListReqVo struct {
	OrgId  int64                              `json:"orgId"`
	UserId int64                              `json:"userId"`
	Input  GetUserOrDeptSameNameListReqVoData `json:"input"`
}

type GetUserOrDeptSameNameListReqVoData struct {
	DataType string `json:"dataType"`
}

type PersonalInfo struct {
	// 主键
	ID int64 `json:"id"`
	// 工号
	EmplID *string `json:"emplId"`
	// 组织id
	OrgID int64 `json:"orgId"`
	// 组织名称
	OrgName string `json:"orgName"`
	// 组织code
	OrgCode string `json:"orgCode"`
	// 名称
	Name string `json:"name"`
	// 第三方名称
	ThirdName string `json:"thirdName"`
	// 登录名
	LoginName string `json:"loginName"`
	// 登录名编辑次数
	LoginNameEditCount int `json:"loginNameEditCount"`
	// 邮箱
	Email string `json:"email"`
	// 电话
	Mobile string `json:"mobile"`
	// 生日
	Birthday types.Time `json:"birthday"`
	// 性别
	Sex int `json:"sex"`
	// 剩余使用时长
	Rimanente int `json:"rimanente"`
	// 付费等级
	Level int `json:"level"`
	// 付费等级名
	LevelName string `json:"levelName"`
	// 头像
	Avatar string `json:"avatar"`
	// 来源
	SourceChannel string `json:"sourceChannel"`
	// 语言
	Language string `json:"language"`
	// 座右铭
	Motto string `json:"motto"`
	// 上次登录ip
	LastLoginIP string `json:"lastLoginIp"`
	// 上次登录时间
	LastLoginTime types.Time `json:"lastLoginTime"`
	// 登录失败次数
	LoginFailCount int `json:"loginFailCount"`
	// 创建时间
	CreateTime types.Time `json:"createTime"`
	// 更新时间
	UpdateTime types.Time `json:"updateTime"`
	// 密码是否设置过(1已设置0未设置)
	PasswordSet int `json:"passwordSet"`
	// 是否需要提醒（1需要2不需要）
	RemindBindPhone int `json:"remindBindPhone"`
	// 是否是超管
	IsAdmin bool `json:"isAdmin"`
	// 是否是管理员
	IsManager bool `json:"isManager"`
	// 是否是三方平台管理员，如：飞书管理员
	IsPlatformAdmin bool `json:"isPlatformAdmin"`
	// 权限
	Functions []FunctionLimitObj `json:"functions"`
	// 一些额外数据，如：观看新手指引的状态
	ExtraDataMap map[string]interface{} `json:"extraDataMap"`
	// 个人中心弹窗是否弹出
	RemindPopUp int `json:"remindPopUp"` // 1 弹出， 2 不弹出
}

type FunctionLimitObj struct {
	Name     string              `json:"name"`
	Key      string              `json:"key"`
	HasLimit bool                `json:"hasLimit"`
	Limit    []FunctionLimitItem `json:"limit"`
}

type FunctionLimitItem struct {
	Typ  string `json:"typ"`
	Num  int    `json:"num"`
	Unit string `json:"unit"`
}

type NormalAdminAppIds struct {
	UserId int64   `json:"userId"`
	AppIds []int64 `json:"appIds"`
}

type CreateOrgMemberReqVo struct {
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
	Input  CreateOrgMemberReq `json:"input"`
}

type CreateOrgMemberReq struct {
	// 账号type  1:手机号  2:账号名
	AccountType int `json:"accountType"`
	// 账号名
	AccountName string `json:"accountName"`
	// 手机区号
	PhoneRegion string `json:"phoneRegion"`
	// 手机号
	PhoneNumber string `json:"phoneNumber"`
	// 邮箱
	Email string `json:"email"`
	// 密码
	Password string `json:"password"`
	// 昵称
	Name string `json:"name"`
	// 部门id
	DepartmentIds []int64 `json:"departmentIds"`
	// 角色id
	RoleGroupIds []int64 `json:"roleGroupIds"`
	// 状态（1启用2禁用3离职, 默认启用） 忽视0
	Status int `json:"status"`
}

type UpdateOrgMemberReqVo struct {
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
	Input  UpdateOrgMemberReq `json:"input"`
}

// UpdateOrgMemberReq 更新组织成员信息
type UpdateOrgMemberReq struct {
	// 账号type  1:手机号  2:账号名
	AccountType int `json:"accountType"`
	// UserId
	UserId int64 `json:"userId"`
	// 手机区号
	PhoneRegion string `json:"phoneRegion"`
	// 手机号
	PhoneNumber string `json:"phoneNumber"`
	// 密码
	Password string `json:"password"`
	// 邮箱
	Email string `json:"email"`
	// 姓名
	Name string `json:"name"`
	// 部门/职级列表
	DeptAndPositions []DeptAndPositionReq `json:"deptAndPositions"`
	// 用户所在的部门 ids. polaris 没有职级,因此无法使用 DeptAndPositions 字段
	DepartmentIds []int64 `json:"departmentIds"`
	// 管理组。一个用户可以在多个管理组中。
	RoleGroupIds []int64 `json:"roleGroupIds"`
	// 状态（1启用2禁用3离职 选填）
	Status int `json:"status"`
	// 头像。可选
	Avatar string `json:"avatar"`
}

type DeptAndPositionReq struct {
	// DepartmentId 部门ID(必填)
	DepartmentId int64 `json:"departmentId"`
	// PositionId 职级ID(必填 注意是Org内的OrgPositionId)
	PositionId int64 `json:"positionId"`
}

type ChangeUserManageGroupReq struct {
	UserId         int64   `json:"userId"`
	ManageGroupIds []int64 `json:"manageGroupIds"`
}

type GetOrgSuperAdminInfoReq struct {
	OrgId int64 `json:"orgId"`
}

type GetOrgSuperAdminInfoResp struct {
	vo.Err
	Data []*GetOrgSuperAdminInfoData `json:"data"`
}

type GetOrgSuperAdminInfoData struct {
	UserId int64  `json:"userId"`
	OpenId string `json:"openId"`
}

type UpdateUserToSysManageGroupReq struct {
	OrgId int64                          `json:"orgId"`
	Input UpdateUserToSysManageGroupData `json:"input"`
}

type UpdateUserToSysManageGroupData struct {
	UserIds    []int64 `json:"userIds"`
	UpdateType int     `json:"updateType"`
}

type ExchangeShareTokenData struct {
	ShareKey string
	Password string
}

type ExchangeShareTokenReq struct {
	Input ExchangeShareTokenData `json:"input"`
}
