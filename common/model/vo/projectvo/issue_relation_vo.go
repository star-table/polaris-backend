package projectvo

import (
	"github.com/star-table/common/core/types"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
)

type RelatedIssueListRespVo struct {
	vo.Err
	RelatedIssueList *vo.IssueRestInfoResp `json:"relatedIssueList"`
}

//type CreateIssueRelationIssueReqVo struct {
//	Input  vo.UpdateIssueAndIssueRelateReq `json:"input"`
//	UserId int64                           `json:"userId"`
//	OrgId  int64                           `json:"orgId"`
//}

type CreateIssueRelationTagsReqVo struct {
	Input  vo.UpdateIssueTagsReq `json:"input"`
	UserId int64                 `json:"userId"`
	OrgId  int64                 `json:"orgId"`
}

type IssueResourcesReqVo struct {
	Page  uint                     `json:"page"`
	Size  uint                     `json:"size"`
	Input *vo.GetIssueResourcesReq `json:"input"`
	OrgId int64                    `json:"orgId"`
}

type IssueResourcesRespVo struct {
	vo.Err
	IssueResources *vo.ResourceList `json:"issueResources"`
}

type CreateIssueCommentReqVo struct {
	Input  vo.CreateIssueCommentReq `json:"input"`
	UserId int64                    `json:"userId"`
	OrgId  int64                    `json:"orgId"`
}

type CreateIssueResourceReqVo struct {
	Input         vo.CreateIssueResourceReq `json:"input"`
	UserId        int64                     `json:"userId"`
	OrgId         int64                     `json:"orgId"`
	SourceChannel string                    `json:"sourceChannel"`
}

type CreateIssueResourceRespVo struct {
	vo.Err
	Resource bo.ResourceBo `json:"resource"`
}

type DeleteIssueResourceReqVo struct {
	Input  vo.DeleteIssueResourceReq `json:"input"`
	UserId int64                     `json:"userId"`
	OrgId  int64                     `json:"orgId"`
}

type RelatedIssueListReqVo struct {
	Input         vo.RelatedIssueListReq `json:"input"`
	OrgId         int64                  `json:"orgId"`
	SourceChannel string                 `json:"sourceChannel"`
}

type GetIssueMembersReqVo struct {
	IssueId int64 `json:"issueId"`
	OrgId   int64 `json:"orgId"`
}

type GetIssueMembersRespVo struct {
	Data GetIssueMembersRespData `json:"data"`

	vo.Err
}

type GetIssueMembersRespData struct {
	MemberIds      []int64 `json:"memberIds"`
	OwnerId        []int64 `json:"ownerId"`
	ParticipantIds []int64 `json:"participantIds"`
	FollowerIds    []int64 `json:"followerIds"`
}

type GetIssueRelationResourceReqVo struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type GetIssueRelationResourceRespVo struct {
	vo.Err
	Data []bo.IssueRelationBo `json:"data"`
}

type AddIssueAttachmentFsReq struct {
	UserId int64                      `json:"userId"`
	OrgId  int64                      `json:"orgId"`
	Input  vo.AddIssueAttachmentFsReq `json:"input"`
}

type AddIssueAttachmentFsRespVo struct {
	vo.Err
	Data *vo.AddIssueAttachmentFsResp `json:"data"`
}

type AddIssueAttachmentReq struct {
	UserId int64                    `json:"userId"`
	OrgId  int64                    `json:"orgId"`
	Input  vo.AddIssueAttachmentReq `json:"input"`
}

type UpdateIssueBeforeAfterIssuesReq struct {
	OrgId  int64                              `json:"orgId"`
	UserId int64                              `json:"userId"`
	Input  vo.UpdateIssueBeforeAfterIssuesReq `json:"input"`
}

type BeforeAfterIssueListReq struct {
	OrgId         int64  `json:"orgId"`
	SourceChannel string `json:"sourceChannel"`
	IssueId       int64  `json:"issueId"`
}

type BeforeAfterIssueListResp struct {
	vo.Err
	Data *vo.BeforeAfterIssueListResp `json:"data"`
}

type GetIssueListWithConditionsReqVo struct {
	OrgId  int64                      `json:"orgId"`
	UserId int64                      `json:"UserId"`
	Input  IssueListWithConditionsReq `json:"input"`
}

type IssueListWithConditionsReq struct {
	AppId    int64   `json:"appId"`
	TableIds []int64 `json:"tableIds"`
	Page     int32   `json:"page"`
	Size     int32   `json:"size"`
}

type IssueListWithConditionsResp struct {
	vo.Err
	Data []*IssueSimpleInfo `json:"data"`
}

type GetIssueListSimpleByDataIdsReqVo struct {
	OrgId  int64                          `json:"orgId"`
	UserId int64                          `json:"userId"`
	Input  GetIssueListSimpleByDataIdsReq `json:"input"`
}

type GetIssueListSimpleByDataIdsReq struct {
	AppId   int64   `json:"appId"`
	DataIds []int64 `json:"dataIds"`
}

type SimpleIssueListRespVo struct {
	vo.Err
	Data []*IssueSimpleInfo `json:"data"`
}

type IssueSimpleInfo struct {
	DataId          int64      `json:"id,string"`
	IssueId         int64      `json:"issueId"`
	ParentId        int64      `json:"parentId"`
	Title           string     `json:"title"`
	ProjectId       int64      `json:"projectId"`
	TableId         int64      `json:"tableId,string"`
	IssueStatusType int32      `json:"issueStatusType"`
	Owners          []UserInfo `json:"ownerId"`
	FollowerIds     []UserInfo `json:"followerIds"`
	AuditorIds      []UserInfo `json:"auditorIds"`
	PlanStartTime   string     `json:"planStartTime"`
	PlanEndTime     string     `json:"planEndTime"`
}

type UserInfo struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

type BaseInfoBoxForIssueCard struct {
	// OrgInfo           bo.BaseOrgInfoBo  				   `json:"orgInfo"` // 查询组织信息。暂时没用到，先注释，如果需要可解开
	// Project         bo.ProjectBo                      `json:"projectBo"`
	OperateUserId int64             `json:"operateUserId"`
	OperateUser   bo.BaseUserInfoBo `json:"operateUser"`
	//IssueOwnerNameArr  []string                          `json:"issueOwnerNameArr"`
	OwnerInfos         []*bo.BaseUserInfoBo              `json:"ownerInfos"`
	ProjectAuthBo      bo.ProjectAuthBo                  `json:"projectAuthBo"` // 这里如果projectId为0，appId会补上汇总表id
	IssueInfo          bo.IssueBo                        `json:"issueInfo"`
	ParentIssue        bo.IssueBo                        `json:"parentIssue"`
	TableColumnMap     map[string]lc_table.LcCommonField `json:"tableColumnMap"` // 任务所在表的表头信息 map
	IssueTableInfo     TableMetaData                     `json:"issueTableInfo"`
	IssueInfoUrl       string                            `json:"issueInfoUrl"` // “查看详情”按钮链接
	IssuePcUrl         string                            `json:"issuePcUrl"`   // "PC端查看/应用内查看"按钮链接
	CollaboratorIds    []int64                           `json:"collaboratorIds"`
	UnsubscribeCardKey string                            `json:"unsubscribeCardKey"` // 退订的卡片的标识
	// 这个字段只是为了拼凑卡片的时候区分飞书和钉钉、企微，一些可以交互的消息卡片 只有飞书可以使用
	SourceChannel      string                      `json:"sourceChannel"`
	ProjectTableColumn map[string]*TableColumnData `json:"projectTableColumn"`
}

type DataCollaborators struct {
	IssueId             int64                  `json:"issueId"`
	ColumnCollaborators []*ColumnCollaborators `json:"columnCollaborators"`
}

type ColumnCollaborators struct {
	ColumnId string  `json:"columnId"`
	UserIds  []int64 `json:"userIds"`
	DeptIds  []int64 `json:"deptIds"`
}

type FsIssueNoPermissionTipsCard struct {
	OrgId             int64                             `json:"orgId"`
	UserId            int64                             `json:"userId"`
	IssueId           int64                             `json:"issueId"`
	ParentId          int64                             `json:"parentId"`
	ProjectId         int64                             `json:"projectId"`
	ProjectTypeId     int64                             `json:"projectTypeId"`
	AppId             int64                             `json:"appId"`
	TableId           int64                             `json:"tableId"`
	ProjectName       string                            `json:"projectName"`
	TableName         string                            `json:"tableName"`
	OperateUserName   string                            `json:"operateUserName"`
	IssueOwnerName    []string                          `json:"IssueOwnerName"`
	Title             string                            `json:"title"`
	ParentTitle       string                            `json:"parentTitle"`
	ColumnDisplayName string                            `json:"columnDisplayName"`
	CardTitle         string                            `json:"cardTitle"`
	TableColumnMap    map[string]lc_table.LcCommonField `json:"tableColumnMap"`
	PlanStartTime     types.Time                        `json:"planStartTime"`
	PlanEndTime       types.Time                        `json:"planEndTime"`
	IssueLinks        *IssueLinks                       `json:"issueLinks"`
}

type IssueStartChatReqVo struct {
	UserId        int64             `json:"userId"`
	OrgId         int64             `json:"orgId"`
	SourceChannel string            `json:"sourceChannel"`
	Input         IssueStartChatReq `json:"input"`
}

type IssueStartChatReq struct {
	TableId string `json:"tableId"`
	IssueId int64  `json:"issueId"`
	Topic   string `json:"topic"`
}

type IssueStartChatRespVo struct {
	vo.Err
	Data IssueStartChatRespData `json:"data"`
}

type IssueStartChatRespData struct {
	ChatId string `json:"chatId"`
}

type IssueChatTrendObj struct {
	Topic     string `json:"topic"`
	IsNewChat bool   `json:"isNewChat"` // 是否新创建的群聊。如果使用已有的群，则为 false，如果群不存在，创建了群，则为 true。
}

type CardInfoBoxForStartIssueChat struct {
	//OrgInfo            bo.BaseOrgInfoBo                  `json:"orgInfo"` // 查询组织信息。暂时没用到，先注释，如果需要可解开
	OperateUserId         int64                             `json:"operateUserId"`
	OperateUser           *bo.BaseUserInfoBo                `json:"operateUser"`
	IssueOwners           []*bo.BaseUserInfoBo              `json:"issueOwners"`
	ProjectBo             *bo.ProjectBo                     `json:"projectBo"`
	IssueInfo             *bo.IssueBo                       `json:"issueInfo"`
	ParentIssue           *bo.IssueBo                       `json:"parentIssue"`
	TableColumnMap        map[string]lc_table.LcCommonField `json:"tableColumnMap"` // 任务所在表的表头信息 map
	IssueTableInfo        *TableMetaData                    `json:"issueTableInfo"`
	IssueInfoUrl          string                            `json:"issueInfoUrl"`       // “查看详情”按钮链接
	IssuePcUrl            string                            `json:"issuePcUrl"`         // "PC端查看/应用内查看"按钮链接
	CollaborateInfo       []*DataCollaborators              `json:"collaborateInfo"`    // 任务的协作人信息
	UnsubscribeCardKey    string                            `json:"unsubscribeCardKey"` // 退订的卡片的标识
	CardTitleColor        string                            `json:"cardTitleColor"`     // 飞书卡片标题的颜色
	ProjectTableColumnMap map[string]*TableColumnData       `json:"projectTableColumnMap"`
}
