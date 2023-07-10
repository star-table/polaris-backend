package projectvo

import (
	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
)

type CreateIssueReqVo struct {
	CreateIssue   vo.CreateIssueReq `json:"createIssue"`
	UserId        int64             `json:"userId"`
	OrgId         int64             `json:"orgId"`
	SourceChannel string            `json:"sourceChannel"`
	InputAppId    int64             `json:"inputAppId"`
	TableId       int64             `json:"tableId"`
}

type IssueRespVo struct {
	vo.Err
	Issue      *vo.Issue                         `json:"data"`
	Data       []map[string]interface{}          `json:"lessdata"`
	UserDepts  map[string]*uservo.MemberDept     `json:"userDepts,omitempty"`
	RelateData map[string]map[string]interface{} `json:"relateData,omitempty"`
}

type UpdateIssueReqVo struct {
	Input         vo.UpdateIssueReq `json:"input"`
	UserId        int64             `json:"userId"`
	OrgId         int64             `json:"orgId"`
	InputAppId    int64             `json:"inputAppId"`
	SourceChannel string            `json:"sourceChannel"`
}

type UpdateIssueRespVo struct {
	vo.Err
	UpdateIssue *vo.UpdateIssueResp `json:"data"`
}

type DeleteIssueRespVo struct {
	vo.Err
}

type IssueListRespVo struct {
	vo.Err
	HomeIssueInfo *vo.HomeIssueInfoResp `json:"data"`
}

type IssueInfoReqVo struct {
	IssueID int64 `json:"issueId"`
	UserId  int64 `json:"userId"`
	OrgId   int64 `json:"orgId"`
	// 是否包含已删除的。0、1：包含，即使删除/无权限的任务详情也会正常返回。2不包含，如果查询删除/无权限的任务（如：隐私项目中未参与的任务）信息，则返回异常信息。
	IncludeDeletedStatus int `json:"includeDeletedStatus"`
}

type GetIssueRestInfosReqVo struct {
	Page          int                  `json:"page"`
	Size          int                  `json:"size"`
	Input         *vo.IssueRestInfoReq `json:"input"`
	OrgId         int64                `json:"orgId"`
	SourceChannel string               `json:"source_channel"`
}

type GetIssueRestInfosRespVo struct {
	vo.Err
	GetIssueRestInfos *vo.IssueRestInfoResp `json:"data"`
}

type DeleteIssueReqVo struct {
	Input         vo.DeleteIssueReq `json:"input"`
	UserId        int64             `json:"userId"`
	OrgId         int64             `json:"orgId"`
	SourceChannel string            `json:"sourceChannel"`
}

type DeleteIssueBatchReqVo struct {
	Input         vo.DeleteIssueBatchReq `json:"input"`
	UserId        int64                  `json:"userId"`
	OrgId         int64                  `json:"orgId"`
	SourceChannel string                 `json:"sourceChannel"`
	InputAppId    int64                  `json:"inputAppId"`
}

type DeleteIssueBatchRespVo struct {
	Data *vo.DeleteIssueBatchResp `json:"data"`
	vo.Err
}

type UpdateIssueStatusReqVo struct {
	Input         vo.UpdateIssueStatusReq `json:"input"`
	UserId        int64                   `json:"userId"`
	OrgId         int64                   `json:"orgId"`
	SourceChannel string                  `json:"sourceChannel"`
	InputAppId    int64                   `json:"inputAppId"`
}

type MoveIssueReqVo struct {
	Input         vo.MoveIssueReq `json:"input"`
	UserId        int64           `json:"userId"`
	OrgId         int64           `json:"orgId"`
	SourceChannel string          `json:"sourceChannel"`
	InputAppId    int64           `json:"inputAppId"`
}

type NewbieGuideIssuesInitReqVo struct {
	OrgId     int64 `json:"orgId"`
	AppId     int64 `json:"appId"`
	ProjectId int64 `json:"projectId"`
	TableId   int64 `json:"tableId"`
	UserId    int64 `json:"userId"`
}

type GetIssueInfoReqVo struct {
	UserId  int64 `json:"userId"`
	OrgId   int64 `json:"orgId"`
	IssueId int64 `json:"issueId"`
}

type GetIssueInfoRespVo struct {
	vo.Err
	IssueInfo *bo.IssueBo `json:"data"`
}

type IssueInfoListReqVo struct {
	UserId   int64   `json:"userId"`
	OrgId    int64   `json:"orgId"`
	IssueIds []int64 `json:"issueIds"` // 若 issueIds 有值，则优先使用 issueIds。
}

type IssueInfoListReqVoInput struct {
}

type IssueInfoListRespVo struct {
	vo.Err
	IssueInfos []vo.Issue `json:"data"`
}

type IssueInfoListByDataIdsReqVo struct {
	UserId  int64   `json:"userId"`
	OrgId   int64   `json:"orgId"`
	DataIds []int64 `json:"dataIds"`
}

type IssueInfoListByDataIdsRespVo struct {
	vo.Err
	IssueInfos []vo.Issue `json:"data"`
}

type UpdateIssueSortReqVo struct {
	Input         vo.UpdateIssueSortReq `json:"input"`
	UserId        int64                 `json:"userId"`
	OrgId         int64                 `json:"orgId"`
	SourceChannel string                `json:"sourceChannel"`
}

type DailyProjectIssueStatisticsCountRespVo struct {
	//完成数量
	DailyFinishCount int `json:"dailyFinishCount"`
	//剩余未完成
	RemainingCount int `json:"dailyRemainingCount"`
	//逾期任务数量
	OverdueCount int `json:"dailyOverdueCount"`
}

type DailyProjectIssueStatisticsRespVo struct {
	vo.Err
	DailyProjectIssueStatisticsCountRespVo *DailyProjectIssueStatisticsCountRespVo `json:"data"`
}

type DailyProjectIssueStatisticsReqVo struct {
	ProjectId int64 `json:"projectId"`
	OrgId     int64 `json:"orgId"`
}

type CopyIssueReqVo struct {
	Input         *vo.LessCopyIssueReq `json:"input"`
	UserId        int64                `json:"userId"`
	OrgId         int64                `json:"orgId"`
	SourceChannel string               `json:"sourceChannel"`
}

type CopyIssueBatchReqVo struct {
	Input         vo.CopyIssueBatchReq `json:"input"`
	UserId        int64                `json:"userId"`
	OrgId         int64                `json:"orgId"`
	SourceChannel string               `json:"sourceChannel"`
}

type CopyIssueBatchRespVo struct {
	Data int64 `json:"data"`
	vo.Err
}

type MoveIssueBatchReqVo struct {
	Input         vo.MoveIssueBatchReq `json:"input"`
	UserId        int64                `json:"userId"`
	OrgId         int64                `json:"orgId"`
	SourceChannel string               `json:"sourceChannel"`
}

type MoveIssueBatchRespVo struct {
	vo.Err
	Data *vo.MoveIssueBatchResp `json:"data"`
}

type ArchiveIssueReqVo struct {
	UserId        int64   `json:"userId"`
	OrgId         int64   `json:"orgId"`
	SourceChannel string  `json:"sourceChannel"`
	IssueIds      []int64 `json:"issueIds"`
}

type ArchiveIssueRespVo struct {
	vo.Err
	Data *vo.ArchiveIssueBatchResp `json:"data"`
}

type CancelArchiveIssueReqVo struct {
	UserId        int64                    `json:"userId"`
	OrgId         int64                    `json:"orgId"`
	SourceChannel string                   `json:"sourceChannel"`
	Input         vo.CancelArchiveIssueReq `json:"input"`
}

type CancelArchiveIssueRespVo struct {
	vo.Err
	Data *vo.CancelArchiveIssueBatchResp `json:"data"`
}

type ConvertIssueToParentReq struct {
	UserId int64                      `json:"userId"`
	OrgId  int64                      `json:"orgId"`
	Input  vo.ConvertIssueToParentReq `json:"input"`
}

type ChangeParentIssueReq struct {
	UserId int64                   `json:"userId"`
	OrgId  int64                   `json:"orgId"`
	Input  vo.ChangeParentIssueReq `json:"input"`
}

// 协作人
type Collaborators struct {
	// 确认者。可选
	AuditorIds []string `json:"auditorIds"`
	// 。可选
	OwnerID []string `json:"ownerId"`
	// 关联者
	FollowerIds []string `json:"followerIds"`
}

// 建立任务协作人请求结构体
type CreateIssueCollaboratorUsersReq struct {
	// 应用id
	AppID int64 `json:"appId"`
	// 组织id
	OrgID int64 `json:"orgId"`
	// table id
	TableID int64 `json:"tableId"`
	// 项目类型id
	ProjectTypeID int64 `json:"projectTypeId"`
	// 协作人ids
	Collaborators *Collaborators `json:"collaborators"`
}

type CreateIssueCollaboratorUsersReqVo struct {
	Input  CreateIssueCollaboratorUsersReq `json:"input"`
	UserId int64                           `json:"userId"`
	OrgId  int64                           `json:"orgId"`
}

// 新增任务时，需要用到的辅助数据可以放在这里
type GetSomeInfoForCreateIssueBo struct {
	ProjectId int64
	TableId   int64
	Project   *bo.ProjectBo
	ColumnMap map[string]lc_table.LcCommonField `json:"columnMap"` // 表头信息
}

type UpdateChatIssueVo struct {
	IssueId   int64
	LcData    map[string]interface{}
	OldLcData map[string]interface{}
}

type IssueCardShareReqVo struct {
	UserId int64             `json:"userId"`
	OrgId  int64             `json:"orgId"`
	Input  IssueCardShareReq `json:"input"`
}

type IssueCardShareReq struct {
	IssueId int64    `json:"issueId"`
	ChatIds []string `json:"chatIds"`
	OpenIds []string `json:"openIds"`
	DeptIds []string `json:"deptIds"`
}

type IssueCardShareResp struct {
	vo.Err
	Data bool `json:"data"`
}

type IssueRowListReq struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Input  *tablePb.ListRawRequest `json:"input"`
}

type IssueRowListResp struct {
	vo.Err
	Data []map[string]interface{} `json:"data"`
}

type LcDataListRespVo struct {
	vo.Err
	Data       []map[string]interface{}          `json:"data"`
	UserDept   map[string]*uservo.MemberDept     `json:"userDept,omitempty"`
	RelateData map[string]map[string]interface{} `json:"relateData,omitempty"`
}
