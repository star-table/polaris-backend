package projectvo

import (
	tableV1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"

	permissionV1 "github.com/star-table/interface/golang/permission/v1"

	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/status"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type ProjectsRepVo struct {
	Page             int              `json:"page"`
	Size             int              `json:"size"`
	ProjectExtraBody ProjectExtraBody `json:"projectExtraBody"`
	OrgId            int64            `json:"orgId"`
	UserId           int64            `json:"userId"`
	SourceChannel    string           `json:"sourceChannel"`
}

type ProjectExtraBody struct {
	Params               map[string]interface{} `json:"params"`
	Order                []*string              `json:"order"`
	Input                *vo.ProjectsReq        `json:"input"`
	NoNeedRedundancyInfo bool                   `json:"noNeedRedundancyInfo"`
	ProjectTypeIds       []int64                `json:"projectTypeIds"`
}

//项目列表
type ProjectsRespVo struct {
	vo.Err
	ProjectList *vo.ProjectList `json:"data"`
}

//项目统计入参
type ProjectDayStatsReqVo struct {
	Page   uint                  `json:"page"`
	Size   uint                  `json:"size"`
	Params *vo.ProjectDayStatReq `json:"params"`
	OrgId  int64                 `json:"orgId"`
}

//项目统计出参
type ProjectDayStatsRespVo struct {
	vo.Err
	ProjectDayStatList *vo.ProjectDayStatList `json:"data"`
}

type CreateProjectReqVo struct {
	Input         vo.CreateProjectReq `json:"input"`
	OrgId         int64               `json:"orgId"`
	UserId        int64               `json:"userId"`
	Version       string              `json:"version"`
	SourceChannel string              `json:"sourceChannel"`
}

type UpdateProjectReqVo struct {
	Input         vo.UpdateProjectReq `json:"input"`
	OrgId         int64               `json:"orgId"`
	UserId        int64               `json:"userId"`
	SourceChannel string              `json:"sourceChannel"`
}

type ProjectRespVo struct {
	vo.Err
	Project *vo.Project `json:"data"`
}

type ProjectIdReqVo struct {
	ProjectId     int64  `json:"projectId"`
	OrgId         int64  `json:"orgId"`
	UserId        int64  `json:"userId"`
	SourceChannel string `json:"sourceChannel"`
}

type GetCacheProjectInfoReqVo struct {
	ProjectId int64 `json:"projectId"`
	OrgId     int64 `json:"orgId"`
}

type GetCacheProjectInfoRespVo struct {
	vo.Err

	ProjectCacheBo *bo.ProjectAuthBo `json:"data"`
}

type QuitProjectRespVo struct {
	vo.Err
	QuitProject *vo.QuitResult `json:"data"`
}

type OperateProjectRespVo struct {
	vo.Err
	OperateProject *vo.OperateProjectResp `json:"data"`
}

type ProjectStatisticsRespVo struct {
	vo.Err
	ProjectStatistics *vo.ProjectStatisticsResp `json:"data"`
}

type ProjectInfoRespVo struct {
	vo.Err
	ProjectInfo *vo.ProjectInfo `json:"data"`
}

type UpdateProjectStatusReqVo struct {
	Input         vo.UpdateProjectStatusReq `json:"input"`
	OrgId         int64                     `json:"orgId"`
	UserId        int64                     `json:"userId"`
	SourceChannel string                    `json:"sourceChannel"`
}

type ProjectInfoReqVo struct {
	Input         vo.ProjectInfoReq `json:"input"`
	OrgId         int64             `json:"orgId"`
	UserId        int64             `json:"userId"`
	SourceChannel string            `json:"sourceChannel"`
}

type GetProjectProcessIdReqVo struct {
	OrgId               int64 `json:"orgId"`
	ProjectId           int64 `json:"projectId"`
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
}

type GetProjectProcessIdRespVo struct {
	ProcessId int64 `json:"data"`

	vo.Err
}

type GetProjectBoListByProjectTypeLangCodeReqVo struct {
	OrgId               int64   `json:"orgId"`
	ProjectTypeLangCode *string `json:"projectTypeLangCode"`
}

type GetProjectBoListByProjectTypeLangCodeRespVo struct {
	ProjectBoList []bo.ProjectBo `json:"data"`
	vo.Err
}

type AppendProjectDayStatReqVo struct {
	ProjectBo bo.ProjectBo `json:"projectBo"`
	Date      string       `json:"date"`
}

type ProjectObjectTypesReqVo struct {
	Page   uint                      `json:"page"`
	Size   uint                      `json:"size"`
	Params *vo.ProjectObjectTypesReq `json:"params"`
	OrgId  int64                     `json:"orgId"`
}

type ProjectObjectTypeWithProjectVo struct {
	ProjectId int64 `json:"projectId"`
	OrgId     int64 `json:"orgId"`
}

type GetSimpleProjectInfoReqVo struct {
	OrgId int64   `json:"orgId"`
	Ids   []int64 `json:"ids"`
}

type GetSimpleProjectInfoRespVo struct {
	vo.Err
	Data *[]vo.Project `json:"data"`
}

type GetProjectDetailsRespVo struct {
	vo.Err
	Data []bo.ProjectDetailBo `json:"data"`
}

type ProjectRelationList struct {
	Id           int64 `json:"id"`
	RelationType int   `json:"relationType"`
	RelationId   int64 `json:"relationId"`
}

type GetProjectRelationReqVo struct {
	ProjectId    int64   `json:"projectId"`
	RelationType []int64 `json:"relationType"`
}

type GetProjectRelationBatchReqVo struct {
	OrgId int64                        `json:"orgId"`
	Data  *GetProjectRelationBatchData `json:"data"`
}

type GetProjectRelationBatchData struct {
	ProjectIds    []int64 `json:"projectIds"`
	RelationTypes []int   `json:"relationTypes"`
}

type GetProjectRelationRespVo struct {
	vo.Err
	Data []ProjectRelationList `json:"data"`
}

type GetProjectRelationBatchRespVo struct {
	vo.Err
	Data []bo.ProjectRelationBo `json:"data"`
}

type GetProjectInfoListByOrgIdsReqVo struct {
	OrgIds []int64 `json:"orgIds"`
}

type GetProjectInfoListByOrgIdsListRespVo struct {
	vo.Err
	ProjectInfoListByOrgIdsRespVo []GetProjectInfoListByOrgIdsRespVo `json:"data"`
}

type GetProjectInfoListByOrgIdsRespVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
	Owner     int64 `json:"Owner"`
}

type GetProjectIdByChatIdRespVo struct {
	vo.Err
	Data *GetProjectIdByChatIdResp `json:"data"`
}

type GetProjectIdByChatIdResp struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type GetProjectIdsByChatIdRespVo struct {
	vo.Err
	Data *GetProjectIdsByChatIdRespData `json:"data"`
}

type GetProjectIdsByChatIdRespData struct {
	ProjectIds []int64 `json:"projectIds"`
}

type OrgProjectListReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetProjectIdByChatIdReqVo struct {
	OpenChatId string `json:"openChatId"`
}

type GetProjectIdsByChatIdReqVo struct {
	OrgId      int64  `json:"orgId"`
	OpenChatId string `json:"openChatId"`
}

type GetProjectIdByAppIdReqVo struct {
	AppId string `json:"openChatId"`
}

type GetProjectIdByAppIdResp struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type GetProjectIdByAppIdRespVo struct {
	vo.Err
	Data *GetProjectIdByAppIdResp `json:"data"`
}

type OrgProjectMemberReqVo struct {
	OrgId     int64 `json:"orgId"`
	UserId    int64 `json:"userId"`
	ProjectId int64 `json:"projectId"`
}

type OrgProjectMemberListRespVo struct {
	vo.Err
	OrgProjectMemberRespVo *OrgProjectMemberRespVo `json:"data"`
}

type OrgProjectMemberRespVo struct {
	Owner        OrgProjectMemberVo   `json:"owner"`
	Participants []OrgProjectMemberVo `json:"participants"`
	Follower     []OrgProjectMemberVo `json:"follower"`
	AllMembers   []OrgProjectMemberVo `json:"allMembers"`
}

type OrgProjectMemberVo struct {
	UserId        int64  `json:"userId"`
	OutUserId     string `json:"outUserId"` //有可能为空
	OrgId         int64  `json:"orgId"`
	OutOrgId      string `json:"outOrgId"` //有可能为空
	Name          string `json:"name"`
	NamePy        string `json:"namePy"`
	Avatar        string `json:"avatar"`
	HasOutInfo    bool   `json:"hasOutInfo"`
	HasOrgOutInfo bool   `json:"hasOrgOutInfo"`

	OrgUserIsDelete    int `json:"orgUserIsDelete"`    //是否被组织移除
	OrgUserStatus      int `json:"orgUserStatus"`      //用户组织状态
	OrgUserCheckStatus int `json:"orgUserCheckStatus"` //用户组织审核状态
}

type RemoveProjectMemberReqVo struct {
	Input         vo.RemoveProjectMemberReq `json:"input"`
	OrgId         int64                     `json:"orgId"`
	UserId        int64                     `json:"userId"`
	SourceChannel string                    `json:"sourceChannel"`
}

type ProjectUserListReq struct {
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
	OrgId int64                 `json:"orgId"`
	Input vo.ProjectUserListReq `json:"input"`
}

type ProjectUserListRespVo struct {
	vo.Err
	Data *vo.ProjectUserListResp `json:"data"`
}

// 融合版本的项目成员列表
type ProjectUserListForFuseRespVo struct {
	vo.Err
	Data *vo.ProjectUserListResp `json:"data"`
}

type PayLimitNumReq struct {
	OrgId int64
}

type PayLimitNumResp struct {
	vo.Err
	Data *vo.PayLimitNumResp `json:"data"`
}

type PayLimitNumRespData struct {
	// 项目数量
	ProjectNum int64 `json:"projectNum"`
	// 任务数量
	IssueNum int64 `json:"issueNum"`
	// 文件大小
	FileSize int64 `json:"fileSize"`
	// 现有的仪表盘数量
	DashNum int64 `json:"dashNum"`
}

type PayLimitNumForRestResp struct {
	vo.Err
	Data *PayLimitNumRespData `json:"data"`
}

type OrgProjectMemberListReq struct {
	OrgId         int64                      `json:"orgId"`
	SourceChannel string                     `json:"sourceChannel"`
	Page          *int                       `json:"page"`
	Size          *int                       `json:"size"`
	Params        vo.OrgProjectMemberListReq `json:"params"`
}

type OrgProjectMemberListResp struct {
	vo.Err
	Data *vo.OrgProjectMemberListResp `json:"data"`
}

type GetProjectRelationUserIdsReq struct {
	ProjectId    int64  `json:"projectId"`
	RelationType *int64 `json:"relationType"`
}

type GetProjectRelationUserIdsResp struct {
	vo.Err
	UserIds []int64 `json:"userIds"`
}

type ProjectMemberIdListReq struct {
	OrgId     int64                       `json:"orgId"`
	ProjectId int64                       `json:"projectId"`
	Data      *ProjectMemberIdListReqData `json:"data"`
}

type ProjectMemberIdListReqData struct {
	IncludeAdmin int `json:"includeAdmin"`
}

type ProjectMemberIdListResp struct {
	vo.Err
	Data *vo.ProjectMemberIDListResp `json:"data"`
}

type GetSimpleProjectsByOrgIdReq struct {
	OrgId int64 `json:"orgId"`
}

type GetSimpleProjectsByOrgIdResp struct {
	vo.Err
	Data []bo.ProjectBo `json:"data"`
}

type GetOrgIssueAndProjectCountReq struct {
	OrgId int64 `json:"orgId"`
}

type GetOrgIssueAndProjectCountResp struct {
	vo.Err
	Data *GetOrgIssueAndProjectCountRespData `json:"data"`
}

type GetOrgIssueAndProjectCountRespData struct {
	IssueCount   uint64 `json:"issueCount"`
	ProjectCount uint64 `json:"projectCount"`
}

type OpenPriorityListReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type OpenGetIterationListReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type OpenGetDemandSourceListReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type OpenGetPropertyListReqVo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
}

type OpenGetIssueStatusListReqVo struct {
	OrgId               int64 `json:"orgId"`
	ProjectId           int64 `json:"projectId"`
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
}

type OpenGetIterationListRespVo struct {
	vo.Err
	Data *OpenGetIterationListRespVoData `json:"data"`
}

type OpenGetIterationListRespVoData struct {
	List []*OpenGetIterationListRespVoDataItem `json:"list"`
}

type OpenGetIterationListRespVoDataItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type OpenSomeAttrListRespVo struct {
	vo.Err
	Data *OpenSomeAttrListRespVoData `json:"data"`
}

type OpenSomeAttrListRespVoData struct {
	List []*OpenSomeAttrListRespVoDataItem `json:"list"`
}

type OpenSomeAttrListRespVoDataItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GetTableColumnsReq struct {
	OrgId  int64                            `json:"orgId"`
	UserId int64                            `json:"userId"`
	Input  *tableV1.ReadTableSchemasRequest `json:"input"`
}

type GetTableStatusReq struct {
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
	Input  GetTableStatusData `json:"input"`
}

type GetTableStatusData struct {
	TableId int64 `json:"tableId"`
}

type GetTableStatusResp struct {
	vo.Err
	Data []status.StatusInfoBo `json:"data"`
}

type GetTableColumnsReqInput struct {
	TableIds  []int64  `json:"tableIds"`
	ColumnIds []string `json:"columnIds"`
}

type GetTableColumnsRespVo struct {
	vo.Err
	*tableV1.ReadTableSchemasReply
}

type GetTablesOfAppReq struct {
	OrgId  int64                            `json:"orgId"`
	UserId int64                            `json:"userId"`
	Input  *tableV1.ReadTablesByAppsRequest `json:"input"`
}

type GetTablesOfAppRespVo struct {
	vo.Err
	Data *tableV1.ReadTablesByAppsReply `json:"data"`
}

type GetTablesByOrgReq struct {
	OrgId  int64                         `json:"orgId"`
	UserId int64                         `json:"userId"`
	Input  *tableV1.ReadOrgTablesRequest `json:"input"`
}

type GetTablesByOrgRespVo struct {
	vo.Err
	Data *TableData `json:"data"`
}

type CreateTableReq struct {
	OrgId  int64            `json:"orgId"`
	UserId int64            `json:"userId"`
	Input  *CreateTableData `json:"input"`
}

type CreateTableData struct {
	AppId                   int64         `protobuf:"varint,1,opt,name=appId,proto3" json:"appId,omitempty"`
	AppType                 int32         `protobuf:"varint,2,opt,name=appType,proto3" json:"appType,omitempty"`
	Name                    string        `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	BasicColumns            []string      `protobuf:"bytes,4,rep,name=basicColumns,proto3" json:"basicColumns,omitempty"`
	IsNeedStoreTable        bool          `protobuf:"varint,5,opt,name=isNeedStoreTable,proto3" json:"isNeedStoreTable,omitempty"`              // 是否需要pg数据表，目前来说，极星都存在汇总表内，所以创建project的时候其实只需要创建表头，不需要数据表
	IsNeedColumn            bool          `protobuf:"varint,6,opt,name=isNeedColumn,proto3" json:"isNeedColumn,omitempty"`                      // 是否需要表头，有些app其实不需要表头数据，只是创建了一个表记录，不创建表头数据
	Columns                 []interface{} `protobuf:"bytes,7,rep,name=columns,proto3" json:"columns,omitempty"`                                 // 如果传了使用这，不传从汇总表合并
	NotNeedSummeryColumnIds []string      `protobuf:"bytes,8,rep,name=notNeedSummeryColumnIds,proto3" json:"notNeedSummeryColumnIds,omitempty"` // 不需要的汇总表字段
	BindSummeryAppId        int64         `protobuf:"varint,9,opt,name=bindSummeryAppId,proto3" json:"bindSummeryAppId,omitempty"`              // 绑定的汇总表appId
}

type CreateTableRespVo struct {
	vo.Err
	Data *CreateTableReply `json:"data"`
}

type CreateTableReply struct {
	AppId int64          `json:"appId,string"`
	Table *TableMetaData `json:"table"`
}

type RenameTableReq struct {
	OrgId  int64                       `json:"orgId"`
	UserId int64                       `json:"userId"`
	Input  *tableV1.RenameTableRequest `json:"input"`
}

type RenameTableResp struct {
	vo.Err
	Data *TableMetaData `json:"data"`
}

type DeleteTableReq struct {
	OrgId  int64                       `json:"orgId"`
	UserId int64                       `json:"userId"`
	Input  *tableV1.DeleteTableRequest `json:"input"`
}

type DeleteTableResp struct {
	vo.Err
	Data *TableMetaData `json:"data"`
}

type SetAutoScheduleReq struct {
	OrgId  int64                           `json:"orgId"`
	UserId int64                           `json:"userId"`
	Input  *tableV1.SetAutoScheduleRequest `json:"input"`
}

type SetAutoScheduleResp struct {
	vo.Err
	Data *TableAutoSchedule `json:"data"`
}

type GetTableInfoReq struct {
	OrgId  int64                     `json:"orgId"`
	UserId int64                     `json:"userId"`
	Input  *tableV1.ReadTableRequest `json:"input"`
}

type GetTableInfoResp struct {
	vo.Err
	Data *ReadTableReply `json:"data"`
}

type ReadTableReply struct {
	Table *TableMetaData `json:"table"`
}

type GetTablesReqVo struct {
	OrgId  int64                      `json:"orgId"`
	UserId int64                      `json:"userId"`
	Input  *tableV1.ReadTablesRequest `json:"input"`
}

type GetTablesDataResp struct {
	vo.Err
	Data *TableData `json:"data"`
}

type TableData struct {
	Tables []*TableMetaData `json:"tables"`
}

type TableMetaData struct {
	AppId            int64  `json:"appId,string"`
	TableId          int64  `json:"tableId,string"`
	Name             string `json:"name"`
	AutoScheduleFlag int32  `json:"autoScheduleFlag"`
}

type TableAutoSchedule struct {
	TableId          int64 `json:"tableId,string"`
	AutoScheduleFlag int32 `json:"autoScheduleFlag"`
}

type GetSummeryTableIdReqVo struct {
	OrgId  int64                              `json:"orgId"`
	UserId int64                              `json:"userId"`
	Input  *tableV1.ReadSummeryTableIdRequest `json:"input"`
}

type GetSummeryTableIdRespVo struct {
	vo.Err
	Data *tableV1.ReadSummeryTableIdReply `json:"data"`
}

type CreateColumnReqVo struct {
	SourceChannel string               `json:"sourceChannel"`
	OrgId         int64                `json:"orgId"`
	UserId        int64                `json:"userId"`
	Input         *CreateColumnRequest `json:"input"`
}

type CreateColumnRequest struct {
	ProjectId         int64            `json:"projectId"`
	AppId             int64            `json:"appId"`
	TableId           int64            `json:"tableId"`
	ActiveViewId      int64            `json:"activeViewId"`
	SourceOrgColumnId string           `json:"sourceOrgColumnId"` // 来源的团队字段的列名，用于记录对应关系
	Column            *TableColumnData `json:"column"`
}

type CreateColumnRespVo struct {
	vo.Err
	Data *CreateColumnReply `json:"data"`
}

type CreateColumnReply struct {
	AppId   int64            `json:"appId"`
	TableId int64            `json:"tableId"`
	Column  *TableColumnData `json:"column"`
}

type CopyColumnReqVo struct {
	OrgId  int64              `json:"orgId"`
	UserId int64              `json:"userId"`
	Input  *CopyColumnRequest `json:"input"`
}

type CopyColumnRequest struct {
	ProjectId   int64  `json:"projectId"`
	AppId       int64  `json:"appId"`
	TableId     int64  `json:"tableId"`
	SrcColumnId string `json:"srcColumnId"`
}

type CopyColumnRespVo struct {
	vo.Err
	Data *tableV1.CopyColumnReply `json:"data"`
}

type UpdateColumnReqVo struct {
	OrgId         int64                   `json:"orgId"`
	UserId        int64                   `json:"userId"`
	SourceChannel string                  `json:"sourceChannel"`
	Input         *UpdateColumnReqVoInput `json:"input"`
}

// UpdateColumnReqVoInput 类型参考 `*tableV1.UpdateColumnRequest`
type UpdateColumnReqVoInput struct {
	ProjectId int64            `json:"projectId"`
	AppId     int64            `json:"appId"`
	TableId   int64            `json:"tableId"`
	Column    *TableColumnData `json:"column"`
}

type UpdateColumnRespVo struct {
	vo.Err
	Data *UpdateColumnReply `json:"data"`
}

type UpdateColumnReply struct {
	AppId   int64            `json:"appId"`
	TableId int64            `json:"tableId"`
	Column  *TableColumnData `json:"column"`
}

type UpdateColumnDescriptionReqVo struct {
	OrgId         int64                           `json:"orgId"`
	UserId        int64                           `json:"userId"`
	SourceChannel string                          `json:"sourceChannel"`
	Input         *UpdateColumnDescriptionRequest `json:"input"`
}

type UpdateColumnDescriptionRequest struct {
	ProjectId   int64  `json:"projectId"`
	AppId       int64  `json:"appId"`
	TableId     int64  `json:"tableId"`
	ColumnId    string `json:"columnId"`
	Description string `json:"description"`
}

type UpdateColumnDescriptionRespVo struct {
	vo.Err
	Data *tableV1.UpdateColumnDescriptionReply `json:"data"`
}

type DeleteColumnReqVo struct {
	OrgId         int64                `json:"orgId"`
	UserId        int64                `json:"userId"`
	SourceChannel string               `json:"sourceChannel"`
	Input         *DeleteColumnRequest `json:"input"`
}

type DeleteColumnRequest struct {
	ProjectId int64  `json:"projectId"`
	AppId     int64  `json:"appId"`
	TableId   int64  `json:"tableId"`
	ColumnId  string `json:"columnId"`
}

type DeleteColumnRespVo struct {
	vo.Err
	Data *tableV1.DeleteColumnReply `json:"data"`
}

type ReadTablesByAppsReqVo struct {
	OrgId  int64                            `json:"orgId"`
	UserId int64                            `json:"userId"`
	Input  *tableV1.ReadTablesByAppsRequest `json:"input"`
}

type ReadTablesByAppsRespVo struct {
	vo.Err
	Data *ReadTablesByAppsData `json:"data"`
}

type ReadTablesByAppsData struct {
	AppsTables []*AppTables `json:"appsTables"`
}

type AppTables struct {
	AppId  int64            `json:"appId,string"`
	Tables []*TableMetaData `json:"tables,omitempty"`
}

type GetTableSchemasByAppIdReq struct {
	OrgId  int64                                   `json:"orgId"`
	UserId int64                                   `json:"userId"`
	Input  *tableV1.ReadTableSchemasByAppIdRequest `json:"input"`
}

type GetTableSchemasByAppIdRespVo struct {
	vo.Err
	Data *TablesColumnsRespData `json:"data"`
}

type GetTableSchemasByOrgIdReq struct {
	OrgId     int64    `json:"orgId"`
	UserId    int64    `json:"userId"`
	ColumnIds []string `json:"columnIds"`
}

type GetTableSchemasByOrgIdRespVo struct {
	vo.Err
	Data *TablesColumnsRespData `json:"data"`
}

type CreateSummeryTableReq struct {
	OrgId  int64                              `json:"orgId"`
	UserId int64                              `json:"userId"`
	Input  *tableV1.CreateSummeryTableRequest `json:"input"`
}

type CreateSummeryTableRespVo struct {
	vo.Err
	Data *tableV1.CreateSummeryTableReply `json:"data"`
}

type GetMenuReq struct {
	AppId int64 `json:"appId"`
}

type GetMenuReqVo struct {
	OrgId int64 `json:"orgId"`
	AppId int64 `json:"appId"`
}

type MenuData struct {
	AppId      int64                  `json:"appId,string"`
	Config     map[string]interface{} `json:"config"`
	Tables     []*TableMetaData       `json:"tables"`
	Dashboards []*DashboardInfo       `json:"dashboards"`
}

type GetMenuRespVo struct {
	vo.Err
	Data MenuData `json:"data"`
}

type SaveMenuData struct {
	AppId  int64                  `json:"appId"`
	Config map[string]interface{} `json:"config"`
}

type SaveMenuReqVo struct {
	OrgId int64        `json:"orgId"`
	Input SaveMenuData `json:"input"`
}

type SaveMenuRespVo struct {
	vo.Err
	Data SaveMenuResp `json:"data"`
}

type SaveMenuResp struct {
	AppId int64 `json:"appId,string"`
}

type DashboardsRespVo struct {
	vo.Err
	Data []*DashboardInfo `json:"data"`
}

type DashboardInfo struct {
	Id     int64  `json:"id"`
	Name   string `json:"name"`
	Type   int    `json:"type"`
	Config string `json:"config"`
}

type GetCollaboratorsRequest struct {
	OrgId int64 `json:"orgId"`
	AppId int64 `json:"appId,string"`
}

type GetCollaboratorsResp struct {
	vo.Err
	Data *permissionV1.GetCollaboratorsResponse `json:"data"`
}

type RecycleAttachmentResp struct {
	vo.Err
	Data *tableV1.RecycleAttachmentReply
}

type RecoverAttachmentResp struct {
	vo.Err
	Data *tableV1.RecoverAttachmentReply
}

type QueryProcessForAsyncTaskReqVo struct {
	OrgId  int64                             `json:"orgId"`
	UserId int64                             `json:"userId"`
	Input  QueryProcessForAsyncTaskReqVoData `json:"input"`
}

type QueryProcessForAsyncTaskReqVoData struct {
	TaskId string `json:"taskId"`
}

type QueryProcessForAsyncTaskRespVo struct {
	vo.Err
	Data AsyncTask `json:"data"`
}

type AsyncTask struct {
	Total          int     `json:"total"`
	ErrCode        int     `json:"errCode"` // 出错时的错误码
	StartTimestamp int64   `json:"startTimestamp"`
	Processed      int     `json:"processed"` // 处理成功的数量
	Failed         int     `json:"failed"`    // 处理失败的数量
	PercentVal     float64 `json:"percentVal"`
	Cover          string  `json:"cover"`
	CardSend       string  `json:"cardSend"` // 卡片是否发送 `"1"`:发送
	TableIds       string  `json:"tableIds"`
	//ParentTotalCount int     `json:"parentTotalCount"` // 父任务计数
	//ParentSucCount   int     `json:"parentSucCount"`   // 父任务处理成功计数
	//ParentErrCount   int     `json:"parentErrCount"`   // 父任务处理失败计数
}

type GetProjectMemberIdsReqVo struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Input  GetProjectMemberIdsReq `json:"input"`
}

type GetProjectMemberIdsReq struct {
	ProjectId    int64 `json:"projectId"`
	IncludeAdmin int   `json:"includeAdmin"` // 是否需要包含额外的管理员 id。可选，1表示需要 0不需要。默认0。
}

type GetProjectMemberIdsResp struct {
	vo.Err
	Data *GetProjectMemberIdsData `json:"data"`
}

type GetProjectMemberIdsData struct {
	DepartmentIds []int64 `json:"departmentIds"` // 部门id
	UserIds       []int64 `json:"userIds"`       // 人员id
}

type GetTrendListMembersReqVo struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Input  GetTrendListMembersReq `json:"input"`
}

type GetTrendListMembersReq struct {
	ProjectId int64 `json:"projectId"`
}

type GetTrendListMembersResp struct {
	vo.Err
	Data *GetTrendListMembersData `json:"data"`
}

type GetTrendListMembersData struct {
	Users []UserInfo `json:"users"`
}

type ListRowsReply struct {
	vo.Err
	Data      *tableV1.ListReply
	UserDepts map[string]*uservo.MemberDept `json:"userDepts"`
}

type ListRawRowsReply struct {
	vo.Err
	Data *tableV1.ListRawReply
}

type DeleteReply struct {
	vo.Err
	Data *tableV1.DeleteReply
}

type CheckIsAppCollaboratorReply struct {
	vo.Err
	Data *tableV1.CheckIsAppCollaboratorReply
}

type GetUserAppCollaboratorRolesReply struct {
	vo.Err
	Data *tableV1.GetUserAppCollaboratorRolesReply
}

type GetAppCollaboratorRolesReply struct {
	vo.Err
	Data *tableV1.GetAppCollaboratorRolesReply
}

type GetDataCollaboratorsReply struct {
	vo.Err
	Data *tableV1.GetDataCollaboratorsReply
}

type CreateOrgDirectoryAppsReq struct {
	OrgId int64 `json:"orgId"`
	// 为 true，表示为汇总表增加“全部成员”可见性
	//NeedUpdateSummaryAppVisibilityToAll bool `json:"needUpdateSummaryAppVisibilityToAll"`
}

type PolarisMenuOption struct {
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	LinkUrl string `json:"linkUrl"`
}

type GetProjectStatisticsReqVo struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Input  GetProjectStatisticsReq `json:"input"`
}

type GetProjectStatisticsReq struct {
	ExtraInfoReq []GetProjectStatistics `json:"params"`
}

type GetProjectStatistics struct {
	ProjectId  int64 `json:"projectId"`
	ResourceId int64 `json:"resourceId"`
}

type GetProjectStatisticsResp struct {
	vo.Err
	Data []*GetProjectStatisticsData `json:"data"`
}

type GetProjectStatisticsData struct {
	ProjectId         int64  `json:"projectId"`
	AllIssues         int64  `json:"allIssues"`
	FinishIssues      int64  `json:"finishIssues"`
	OverdueIssues     int64  `json:"overdueIssues"`
	UnfinishedIssues  int64  `json:"unfinishedIssues"`
	CoverResourceId   int64  `json:"resourceId"`
	CoverResourcePath string `json:"resourcePath"`
}
