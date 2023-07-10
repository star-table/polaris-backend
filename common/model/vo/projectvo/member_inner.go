package projectvo

type ChangeProjectMemberReq struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Input  ChangeProjectMemberData `json:"input"`
}

type ChangeProjectMemberData struct {
	ProjectId  int64   `json:"projectId"`
	AddUserIds []int64 `json:"addUserIds"`
	AddDeptIds []int64 `json:"addDeptIds"`
	DelUserIds []int64 `json:"delUserIds"`
	DelDeptIds []int64 `json:"delDeptIds"`
}