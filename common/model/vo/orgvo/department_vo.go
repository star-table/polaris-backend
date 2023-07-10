package orgvo

import "github.com/star-table/polaris-backend/common/model/vo"

type DepartmentsReqVo struct {
	Page          *int
	Size          *int
	Params        *vo.DepartmentListReq
	CurrentUserId int64 `json:"userId"`
	OrgId         int64 `json:"orgId"`
}

type DepartmentsRespVo struct {
	DepartmentList *vo.DepartmentList `json:"data"`
	vo.Err
}

type DepartmentMembersReqVo struct {
	CurrentUserId int64                      `json:"userId"`
	OrgId         int64                      `json:"orgId"`
	Params        vo.DepartmentMemberListReq `json:"params"`
}

type DepartmentMembersRespVo struct {
	DepartmentMemberInfos []*vo.DepartmentMemberInfo `json:"data"`
	vo.Err
}

type CreateDepartmentReqVo struct {
	Data          CreateDepartmentReq `json:"data"`
	CurrentUserId int64               `json:"userId"`
	OrgId         int64               `json:"orgId"`
}

type CreateDepartmentReq struct {
	// 父部门id（根目录为0）
	ParentID int64 `json:"parentId"`
	// 部门名称
	Name string `json:"name"`
	// 部门主管(选填)
	LeaderIds *[]int64 `json:"leaderIds"`
}

type UpdateDepartmentReqVo struct {
	Data          UpdateDepartmentReq `json:"data"`
	CurrentUserId int64               `json:"userId"`
	OrgId         int64               `json:"orgId"`
}

type UpdateDepartmentReq struct {
	//部门id
	DepartmentId int64 `json:"departmentId"`
	//部门名称（选填）
	Name *string `json:"name"`
	// 部门主管(选填，不传表示不更新，传空数组则表示取消主管)
	LeaderIds *[]int64 `json:"leaderIds"`
}

type DeleteDepartmentReqVo struct {
	Data          DeleteDepartmentReq `json:"data"`
	CurrentUserId int64               `json:"userId"`
	OrgId         int64               `json:"orgId"`
}

type DeleteDepartmentReq struct {
	//部门id
	DepartmentId int64 `json:"departmentId"`
}

type AllocateDepartmentReqVo struct {
	Data          AllocateDepartmentReq `json:"data"`
	CurrentUserId int64                 `json:"userId"`
	OrgId         int64                 `json:"orgId"`
}

type AllocateDepartmentReq struct {
	//用户id
	UserIds []int64 `json:"userIds"`
	//部门id
	DepartmentIds []int64 `json:"departmentIds"`
}

type SetUserDepartmentLevelReqVo struct {
	Data          SetUserDepartmentLevelReq `json:"data"`
	CurrentUserId int64                     `json:"userId"`
	OrgId         int64                     `json:"orgId"`
}

type SetUserDepartmentLevelReq struct {
	//用户id
	UserId int64 `json:"userId"`
	//是否是部门主管(1是2否)
	IsLeader int `json:"isLeader"`
	//部门id
	DepartmentId int64 `json:"departmentId"`
}

type DepartmentMembersListReq struct {
	OrgId         int64                        `json:"orgId"`
	SourceChannel string                       `json:"sourceChannel"`
	Params        *vo.DepartmentMembersListReq `json:"params"`
	Page          *int                         `json:"page"`
	Size          *int                         `json:"size"`
	//是否忽略离职的，不传默认为否
	IgnoreDelete bool `json:"ignoreDelete"`
}

type DepartmentMembersListResp struct {
	vo.Err
	Data *vo.DepartmentMembersListResp `json:"data"`
}

type VerifyDepartmentsReq struct {
	OrgId         int64   `json:"orgId"`
	DepartmentIds []int64 `json:"departmentIds"`
}

type GetDeptByIdsReq struct {
	OrgId   int64   `json:"orgId"`
	DeptIds []int64 `json:"deptIds"`
}

type GetDeptByIdsResp struct {
	List []*vo.Department `json:"data"`
	vo.Err
}

type OrgRemarkConfigType struct {
	OrgSummaryTableAppId int64 `json:"orgSummaryTableAppId"`
	TagAppId             int64 `json:"tagAppId"`
	PriorityAppId        int64 `json:"priorityAppId"`
	IssueBarAppId        int64 `json:"issueBarAppId"`
	SceneAppId           int64 `json:"sceneAppId"`
	IssueStatusAppId     int64 `json:"issueStatusAppId"`
	EmptyProjectAppId    int64 `json:"emptyProjectAppId"`
	// 项目 form 对应的 appId，一个组织一个
	ProjectFormAppId int64 `json:"projectFormAppId"`
	// “项目视图”所放的目录 id
	ProjectFolderAppId int64 `json:"projectFolderAppId"`
	// “任务视图”所放的目录 id
	IssueFolderAppId int64 `json:"issueFolderAppId"`
}

type SimpleDeptForCustomField struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	IsDelete int    `json:"isDelete"`
	Status   int    `json:"status"`
	Type     string `json:"type"` // 无码自定义字段，部门的 type，值为：`D_`
}
