package projectvo

import (
	"github.com/star-table/polaris-backend/common/model/vo"
)

type GetProjectTemplateReq struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type GetProjectTemplateResp struct {
	vo.Err

	Data *GetProjectTemplateData `json:"data"`
}

type AuthCreateProjectReq struct {
	OrgId int64 `json:"orgId"`
}

type GetProjectTemplateData struct {
	TemplateFlag               *int                        `json:"templateFlag"`
	ProjectObjectTypeTemplates []ProjectObjectTypeTemplate `json:"projectObjectTypeTemplates"`
	ProjectIterationTemplates  []ProjectIterationTemplate  `json:"projectIterationTemplates"`
	//ProjectIssueTemplates      []map[string]interface{}    `json:"projectIssueTemplates"`
	ProjectInfo      ProjectInfoData    `json:"projectInfo"`
	ProjectDetail    ProjectDetailData  `json:"projectDetail"`
	TemplateInfo     SimpleTemplateInfo `json:"templateInfo"`
	IsNewbieGuide    bool               `json:"isNewbieGuide"`
	IsCreateTemplate bool               `json:"isCreateTemplate"`
	NeedData         bool               `json:"needData"`
	FromOrgId        int64              `json:"fromOrgId"`
}

type SimpleTemplateInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type ProjectInfoData struct {
	AppId         int64  `json:"appId"`
	OrgId         int64  `json:"orgId"`
	Status        int64  `json:"status"`
	Name          string `json:"name"`
	PublishStatus int    `json:"publishStatus"`
	ProjectTypeId int64  `json:"projectTypeId"`
	Owner         int64  `json:"owner"`
	PublicStatus  int    `json:"publicStatus"`
	ResourceId    int64  `json:"resourceId"`
	IsFiling      int    `json:"isFiling"`
	Remark        string `json:"remark"`
}

type ProjectDetailData struct {
	Notice            string `json:"notice"`
	IsEnableWorkHours int    `json:"isEnableWorkHours"`
	IsSyncOutCalendar int    `json:"isSyncOutCalendar"`
}

type ProjectIterationTemplate struct {
	ID              int64             `json:"id"`
	Name            string            `json:"name"`
	StartTime       string            `json:"startTime"`
	EndTime         string            `json:"endTime"`
	Sort            int64             `json:"sort"`
	Owner           int64             `json:"owner"`
	StatusType      int               `json:"status"` //这里存statusType，因为每个项目的迭代状态是不一样的
	IterationStatus []IterationStatus `json:"iterationStatus"`
}

type IterationStatus struct {
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	PlanStartTime string `json:"planStartTime"`
	PlanEndTime   string `json:"planEndTime"`
	StatusType    int    `json:"statusType"` //这里存statusType，因为每个项目的迭代状态是不一样的
}

type ProjectObjectTypeTemplate struct {
	ID         int64                   `json:"id"`
	Name       string                  `json:"name"`
	ProcessId  int64                   `json:"processId"`
	LangCode   string                  `json:"langCode"`
	Sort       int                     `json:"sort"`
	ObjectType int                     `json:"objectType"`
	Status     []ProjectStatusTemplate `json:"status"`
}

type ProjectStatusTemplate struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Type         int    `json:"type"`
	FontStyle    string `json:"fontStyle"`
	BgStyle      string `json:"bgStyle"`
	IsInitStatus int    `json:"isInitStatus"`
	Sort         int    `json:"sort"`
	Category     int    `json:"category"`
}

type ApplyProjectTemplateReq struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	AppId  int64                  `json:"appId"`
	Input  GetProjectTemplateData `json:"input"`
}

type ApplyProjectTemplateResp struct {
	vo.Err
	Data *ApplyProjectTemplateData `json:"data"`
}

type ApplyProjectTemplateData struct {
	ProjectId            int64           `json:"projectId"`
	ProjectObjectTypeMap map[int64]int64 `json:"projectObjectTypeMap"`
	IterationMap         map[int64]int64 `json:"iterationMap"`
	StatusMap            map[int64]int64 `json:"statusMap"`
}

type DeleteProjectInnerReq struct {
	OrgId      int64   `json:"orgId"`
	UserId     int64   `json:"userId"`
	ProjectIds []int64 `json:"projectIds"`
}
