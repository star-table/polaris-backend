package appvo

import "github.com/star-table/polaris-backend/common/model/vo"

type GetAppInfoByAppIdReq struct {
	AppId int64 `json:"appId"`
	OrgId int64 `json:"orgId"`
}

type GetAppInfoByAppIdResp struct {
	vo.Err
	Data *GetAppInfoByAppIdRespData `json:"data"`
}

type GetAppInfoByAppIdRespData struct {
	Id          int64  `json:"id"`
	DashboardId int64  `json:"dashboard_id"`
	ExtendsId   int64  `json:"extendsId"`
	FormId      int64  `json:"formId"`
	Icon        string `json:"icon"`
	Name        string `json:"name"`
	OrgId       int64  `json:"orgId"`
	Status      int    `json:"status"`
}

type FormCreateIssueReq struct {
	AppId  int64                     `json:"appId"`
	OrgId  int64                     `json:"orgId"`
	UserId int64                     `json:"userId"`
	Form   []*FormCreateIssueReqForm `json:"form"`
}

type FormCreateIssueReqForm struct {
	ProjectIds          []int64 `json:"_field_polaris_project_id"`
	Code                string  `json:"_field_polaris_issue_code"`
	ProjectObjectTypeId int64   `json:"_field_polaris_issue_project_object_type_id"`
	Path                string  `json:"_field_polaris_issue_path"`
	AuditStatus         int     `json:"_field_polaris_issue_audit_status"`
	Title               string  `json:"_field_polaris_issue_title"`
	Status              string  `json:"_field_polaris_issue_status"`
	StartTime           string  `json:"_field_polaris_issue_startTime"`
	EndTime             string  `json:"_field_polaris_issue_endTime"`
	Priority            int64   `json:"_field_polaris_issue_priority"`
	Owner               int64   `json:"_field_polaris_issue_owner"`
	Followers           []int64 `json:"_field_polaris_issue_follower"`
	Tags                []int64 `json:"_field_polaris_issue_tag"`
}

type FormCreateOneResp struct {
	vo.Err
	Timestamp interface{} `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type FormCreatePriorityReq struct {
	AppId  int64                        `json:"appId"`
	OrgId  int64                        `json:"orgId"`
	UserId int64                        `json:"userId"`
	Form   []*FormCreatePriorityReqForm `json:"form"`
}

type FormCreatePriorityReqForm struct {
	ProjectIds []int64 `json:"_field_polaris_project_id"`
	LangCode   string  `json:"_field_polaris_priority_lang_code"`
	Name       string  `json:"_field_polaris_priority_name"`
	Type       string  `json:"_field_polaris_priority_type"`
	Sort       int64   `json:"_field_polaris_priority_sort"`
	BgStyle    string  `json:"_field_polaris_priority_bg_style"`
	FontStyle  string  `json:"_field_polaris_priority_font_style"`
	IsDefault  string  `json:"_field_polaris_priority_is_default"`
	Remark     string  `json:"_field_polaris_priority_remark"`
}

type UpdateLessCodeAppReq struct {
	AppId  int64  `json:"appId"`
	OrgId  int64  `json:"orgId"`
	UserId int64  `json:"userId"`
	Name   string `json:"name"`
}

type DeleteLessCodeAppReq struct {
	AppId  int64 `json:"appId"`
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type CancelStarAppReq struct {
	AppId  int64 `json:"appId"`
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type CreateAppViewReq struct {
	Reqs []CreateAppViewData `json:"reqs"`
}

type CreateAppViewData struct {
	AppId   int64  `json:"appId"`
	Config  string `json:"config"`
	Name    string `json:"name"`
	OrgId   int64  `json:"orgId"`
	OwnerId int64  `json:"ownerId"`
	Type    int    `json:"type"`
}

type CreateAppViewResp struct {
	vo.Err
	Timestamp interface{} `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type ViewEmptyConfig struct {
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
}

type GetAppViewListReq struct {
	OrgId  int64 `json:"orgId"`
	AppId  int64 `json:"appId"`
	UserId int64 `json:"userId"`
}

type GetAppViewListResp struct {
	vo.Err
	Data []GetAppViewListRespDataListItem `json:"data"`
}

type GetAppViewListRespData struct {
	List []GetAppViewListRespDataListItem `json:"list"`
}

type GetAppViewListRespDataListItem struct {
	AppID     string      `json:"appId"`
	ID        string      `json:"id"`
	IsPrivate bool        `json:"isPrivate"`
	OrgID     string      `json:"orgId"`
	Owner     string      `json:"owner"`
	Remark    string      `json:"remark"`
	Sort      string      `json:"sort"`
	Type      int         `json:"type"`
	ViewName  string      `json:"viewName"`
	Config    interface{} `json:"config"` // GetAppViewListRespDataListItemConfig
}

type GetAppViewListRespDataListItemConfig struct {
	Condition           interface{}   `json:"condition"`
	RealCondition       interface{}   `json:"realCondition"`
	HiddenColumnIds     []interface{} `json:"hiddenColumnIds"`
	LessCondition       interface{}   `json:"lessCondition"`
	LessShowCondition   interface{}   `json:"lessShowCondition"`
	Params              interface{}   `json:"params"`
	ProjectObjectTypeID int64         `json:"projectObjectTypeId"`
	TableId             string        `json:"tableId"`
	TableOrder          []string      `json:"tableOrder"`
}

type AddAppMembersReq struct {
	Input *AddAppMembersReqData `json:"input"`
	OrgId int64                 `json:"orgId"`
}

type AddAppMembersReqData struct {
	AppId      int64    `json:"appId"`
	MemberIds  []string `json:"memberIds"`
	PerGroupId int64    `json:"perGroupId"`
}

type AddAppMembersResp struct {
	vo.Err
	Data      interface{} `json:"data"`
	Timestamp interface{} `json:"timestamp"`
}

type GetAppInfoListReq struct {
	OrgId  int64   `json:"orgId"`
	AppIds []int64 `json:"appIds"`
}

type GetAppInfoListResp struct {
	vo.Err
	Data []GetAppInfoListData `json:"data"`
}

type GetAppInfoListData struct {
	Creator      int64  `json:"creator"`
	DashboardId  int64  `json:"dashboardId"`
	ExtendsId    int64  `json:"extendsId"`
	FormId       int64  `json:"formId"`
	Icon         string `json:"icon"`
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	OrgId        int64  `json:"orgId"`
	ParentId     int64  `json:"parentId"`
	PkgId        int64  `json:"pkgId"`
	Status       int64  `json:"status"`
	Type         int64  `json:"type"`
	WorkflowFlag int64  `json:"workflowFlag"`
	MirrorViewId int64  `json:"mirrorViewId"`
	MirrorAppId  int64  `json:"mirrorAppId"`
	ProjectId    int64  `json:"projectId"`
}

type AddAppRelationReq struct {
	AppId        int64 `json:"appId"`
	RelationId   int64 `json:"relationId"`
	RelationType int64 `json:"relationType"`
}

type AddAppRelationResp struct {
	vo.Err
	Data      interface{} `json:"data"`
	Timestamp interface{} `json:"timestamp"`
}

type IsAppMemberReq struct {
	AppId  int64 `json:"appId"`
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type ApplyTemplateReq struct {
	OrgId         int64 `json:"orgId"`
	UserId        int64 `json:"userId"`
	TemplateId    int64 `json:"templateId"`
	IsNewbieGuide bool  `json:"isNewbieGuide"`
}

type ApplyTemplateResp struct {
	vo.Err
	Data      interface{} `json:"data"`
	Timestamp interface{} `json:"timestamp"`
}
