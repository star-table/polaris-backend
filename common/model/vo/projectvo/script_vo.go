package projectvo

type InitAppForProjectReq struct {
	ProjectId int64 `json:"projectId"`
	AppId     int64 `json:"appId"`
	OrgId     int64 `json:"orgId"`
}

type ScriptCommonReq struct {
	OrgId int64 `json:"orgId"`
}

type ScheduleOrgMenuReq struct {
	OrgId int64 `json:"orgId"`
}

type ScheduleOrgMenuOption struct {
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	LinkUrl string `json:"linkUrl"`
}

type SyncOrgDefaultViewMirrorReq struct {
	OrgIds []int64 `json:"orgIds"`
	// OrgIds 为空时， StartPage 起作用。表示从多少页的偏移量查询组织开始跑脚本
	StartPage int `json:"startPage"`
	// 为 true，表示为汇总表增加“全部成员”可见性
	NeedUpdateSummaryAppVisibilityToAll bool `json:"needUpdateSummaryAppVisibilityToAll"`
}

type AddTrashMenuReq struct {
	Input AddTrashMenuReqInput `json:"input"`
}

type AddTrashMenuReqInput struct {
	OrgIds []int64 `json:"orgIds"`
	// OrgIds 为空时， StartPage 起作用。表示从多少页的偏移量查询组织开始跑脚本
	StartPage int `json:"startPage"`
}

type AppendNewMenuToOrgRoleReq struct {
	Input AppendNewMenuToOrgRoleReqInput `json:"input"`
}

type AppendNewMenuToOrgRoleReqInput struct {
	OrgIds   []int64  `json:"orgIds"`
	PrmItems []string `json:"prmItems"` // permission items 待追加的权限项、菜单权限项
	// OrgIds 为空时， StartPage 起作用。表示从多少页的偏移量查询组织开始跑脚本
	StartPage int `json:"startPage"`
}

type FixIssueWrongAppIdsReq struct {
	Input FixIssueWrongAppIdsReqData `json:"input"`
}

type FixIssueWrongAppIdsReqData struct {
	OrgIds []int64 `json:"orgIds"`
	// OrgIds 为空时， StartPage 起作用。表示从多少页的偏移量查询组织开始跑脚本
	StartPage int `json:"startPage"`
}
