package orgvo

import (
	"time"

	"github.com/star-table/common/core/types"

	"github.com/star-table/polaris-backend/common/model/vo/projectvo"

	tableV1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type UserOrganizationListReqVo struct {
	UserId int64 `json:"userId"`
}

type UserOrganizationListRespVo struct {
	vo.Err
	UserOrganizationListResp *vo.UserOrganizationListResp `json:"data"`
}

type SwitchUserOrganizationReqVo struct {
	UserId int64  `json:"userId"`
	OrgId  int64  `json:"orgId"`
	Token  string `json:"userToken"`
}

type SwitchUserOrganizationRespVo struct {
	vo.Err
	OrgId int64 `json:"data"`
}

type UpdateOrganizationSettingReqVo struct {
	//入参
	Input vo.UpdateOrganizationSettingsReq `json:"input"`
	//用户Id
	UserId int64 `json:"userId"`
}

type UpdateOrgRemarkSettingReqVo struct {
	OrgId  int64                             `json:"orgId"`
	UserId int64                             `json:"userId"`
	Input  SaveOrgSummaryTableAppIdReqVoData `json:"input"`
}

type UpdateOrganizationSettingRespVo struct {
	vo.Err
	OrgId int64 `json:"data"`
}

type OrganizationInfoReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type OrganizationInfoRespVo struct {
	vo.Err
	OrganizationInfo *vo.OrganizationInfoResp `json:"data"`
}

type GetOutOrgInfoByOutOrgIdReqVo struct {
	OrgId    int64  `json:"orgId"`
	UserId   int64  `json:"userId"`
	OutOrgId string `json:"outOrgId"`
}

type GetOutOrgInfoByOutOrgIdRespVo struct {
	vo.Err
	Data *OutOrgInfo `json:"data"`
}

type LcCreateProIssueTableReqVo struct {
	OrgId     int64 `json:"orgId"`
	UserId    int64 `json:"userId"`
	ExtendsId int64 `json:"extendsId"`
}

type LcCreateProIssueTableRespVo struct {
	vo.Err
	Data *LcCreateProIssueTableReqVoData `json:"data"`
}

type LcCreateProIssueTableReqVoData struct {
}

type GetOutOrgInfoByOrgIdBatchReqVo struct {
	OrgIds []int64 `json:"orgIds"`
}

type GetOutOrgInfoByOrgIdBatchRespVo struct {
	vo.Err
	Data []*OutOrgInfo `json:"data"`
}

type CheckSpecificScopeReqVo struct {
	OrgId     int64  `json:"orgId"`
	UserId    int64  `json:"userId"`
	PowerFlag string `json:"powerFlag"`
}

type CheckSpecificScopeRespVo struct {
	vo.Err
	Data *CheckSpecificScopeRespData `json:"data"`
}

type CheckSpecificScopeRespData struct {
	HasPower bool `json:"hasPower"`
}

type ApplyScopesReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type ApplyScopesRespVo struct {
	vo.Err
	Data *ApplyScopesRespData `json:"data"`
}

type ApplyScopesRespData struct {
	ThirdCode int64  `json:"thirdCode"`
	ThirdMsg  string `json:"thirdMsg"`
}

type OutOrgInfo struct {
	Id              int64  `json:"id"`
	OrgId           int64  `json:"orgId"`
	OutOrgId        string `json:"outOrgId"`
	SourceChannel   string `json:"sourceChannel"`
	SourcePlatform  string `json:"sourcePlatform"`
	TenantCode      string `json:"tenantCode"`
	Name            string `json:"name"`
	Industry        string `json:"industry"`
	IsAuthenticated int    `json:"isAuthenticated"`
	AuthTicket      string `json:"authTicket"`
	AuthLevel       string `json:"authLevel"`
	Status          int    `json:"status"`
}

type ScheduleOrganizationPageListReqVo struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type ScheduleOrganizationPageListRespVo struct {
	vo.Err
	ScheduleOrganizationPageListResp *ScheduleOrganizationPageListResp `json:"data"`
}

type ScheduleOrganizationPageListResp struct {
	Total                        int64                            `json:"total"`
	ScheduleOrganizationListResp *[]*ScheduleOrganizationListResp `json:"list"`
}

type ScheduleOrganizationListResp struct {
	OrgId                      int64  `json:"orgId"`
	ProjectDailyReportSendTime string `json:"projectDailyReportSendTime"`
}

type GetOrgIdListByPageReqVo struct {
	Page  int                         `json:"page"`
	Size  int                         `json:"size"`
	Input GetOrgIdListByPageReqVoData `json:"input"`
}

// 预留，如果要增加筛选条件
type GetOrgIdListByPageReqVoData struct {
	OrgIds []int64 `json:"orgIds"`
}

type GetOrgIdListBySourceChannelReqVo struct {
	Page  int                             `json:"page"`
	Size  int                             `json:"size"`
	Input GetOrgIdListBySourceChannelData `json:"input"`
}

type GetOrgIdListBySourceChannelData struct {
	//是否付费1否2是 其余默认为全部
	IsPaid         int      `json:"isPaid"`
	SourceChannels []string `json:"sourceChannels"`
}
type GetOrgIdListBySourceChannelRespVo struct {
	Data GetOrgIdListBySourceChannelRespData `json:"data"`
	vo.Err
}

type GetOrgIdListBySourceChannelRespData struct {
	OrgIds []int64 `json:"orgIds"`
}

type GetOrgConfigReq struct {
	OrgId int64 `json:"orgId"`
}

type GetOrgConfigResp struct {
	Data *vo.OrgConfig `json:"data"`
	vo.Err
}

type GetOrgConfigRestResp struct {
	vo.Err
	Data *OrgConfig `json:"data"`
}

type GetOrgFunctionConfigReq struct {
	OrgId int64 `json:"orgId"`
}

type GetOrgFunctionConfigResp struct {
	Data *vo.FunctionConfigResp `json:"data"`
	vo.Err
}

type GetFunctionArrByOrgResp struct {
	vo.Err
	Data GetFunctionArrByOrgRespData `json:"data"`
}

type GetFunctionArrByOrgRespData struct {
	Functions []FunctionLimitObj `json:"functions"`
}

type UpdateOrgFunctionConfigReq struct {
	OrgId         int64                    `json:"orgId"`
	UserId        int64                    `json:"userId"`
	SourceChannel string                   `json:"sourceChannel"`
	Input         UpdateFunctionConfigData `json:"input"`
}

type UpdateFunctionConfigData struct {
	Level         int64     `json:"level"`
	BuyType       string    `json:"buyType"`
	PricePlanType string    `json:"pricePlanType"`
	PayTime       time.Time `json:"payTime"`
	TrailDays     int       `json:"trailDays"`  //试用天数
	ExpireDays    int       `json:"expireDays"` //有效天数
	EndDate       time.Time `json:"endDate"`    //有效期限
	Seats         int       `json:"seats"`      //人数限制
}

type TransferOrgOwnerReq struct {
	OrgId      int64 `json:"orgId"`
	UserId     int64 `json:"userId"`
	NewOwnerId int64 `json:"newOwnerId"`
}

type GetAppTicketReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetAppTicketResp struct {
	vo.Err
	Data *vo.GetAppTicketResp `json:"data"`
}

type UpdateOrgBasicShowSettingReq struct {
	OrgId  int64                           `json:"orgId"`
	UserId int64                           `json:"userId"`
	Params vo.UpdateOrgBasicShowSettingReq `json:"params"`
}

type CheckAndSetSuperAdminReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type MigrateIssueToLcReq struct {
	Data MigrateIssueToLcData `json:"data"`
}

type MigrateIssueToLcData struct {
	OrgIds []int64 `json:"orgIds"`
}

type GetOrgColumnsReq struct {
	OrgId  int64                          `json:"orgId"`
	UserId int64                          `json:"userId"`
	Input  *tableV1.ReadOrgColumnsRequest `json:"input"`
}

type GetOrgColumnsRespVo struct {
	vo.Err
	Data *ReadOrgColumnsReply `json:"data"`
}

type ReadOrgColumnsReply struct {
	Columns []*projectvo.TableColumnData `json:"columns"`
}

type CreateOrgColumnReq struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Input  *CreateOrgColumnRequest `json:"input"`
}

type CreateOrgColumnRequest struct {
	Column *projectvo.TableColumnData `json:"column"`
}

type CreateOrgColumnRespVo struct {
	vo.Err
	Data *tableV1.CreateOrgColumnReply
}

type DeleteOrgColumnReq struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Input  *DeleteOrgColumnRequest `json:"input"`
}

type DeleteOrgColumnRequest struct {
	ColumnId string `json:"columnId"`
}

type DeleteOrgColumnRespVo struct {
	vo.Err
	Data *tableV1.DeleteOrgColumnReply `json:"data"`
}

type InitOrgColumnsReq struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Input  *InitOrgColumnsRequest `json:"input"`
}

type InitOrgColumnsRequest struct {
	Columns []interface{} `json:"columns"`
}

type OrgColumnsData struct {
	Name              string      `json:"name"`
	Label             string      `json:"label"`
	AliasTitle        string      `json:"aliasTitle"`
	Description       string      `json:"description"`
	IsSys             bool        `json:"isSys"`
	IsOrg             bool        `json:"isOrg"`
	Writable          bool        `json:"writable"`
	Editable          bool        `json:"editable"`
	Unique            bool        `json:"unique"`
	UniquePreHandler  string      `json:"uniquePreHandler"`
	SensitiveStrategy string      `json:"sensitiveStrategy"`
	SensitiveFlag     int32       `json:"sensitiveFlag"`
	Field             ColumnField `json:"field"`
}

type ColumnField struct {
	Type       int32       `json:"type"`
	CustomType string      `json:"customType"`
	DataType   int32       `json:"dataType"`
	Props      interface{} `json:"props"`
}

type InitOrgColumnRespVo struct {
	vo.Err
	Data *InitOrgColumnsReply `json:"data"`
}

type InitOrgColumnsReply struct{}

type OrgConfig struct {
	// id
	ID int64 `json:"id"`
	// 组织id
	OrgID int64 `json:"orgId"`
	// 付费级别1通用免费，2标准版
	PayLevel int `json:"payLevel"`
	// 付费开始时间
	PayStartTime types.Time `json:"payStartTime"`
	// 付费结束时间
	PayEndTime types.Time `json:"payEndTime"`
	// 付费级别实际（1免费2标准3试用）
	PayLevelTrue int `json:"payLevelTrue"`
	// 创建时间
	CreateTime types.Time `json:"createTime"`
	// 组织人数
	OrgMemberNumber int64 `json:"orgMemberNumber"`
	// 是否是灰度企业
	IsGrayLevel bool `json:"isGrayLevel"`
	// 汇总表id
	SummaryAppID string `json:"summaryAppId"`
	// 展示基础设置
	BasicShowSetting *BasicShowSetting `json:"basicShowSetting"`
	// 企业自定义logo
	Logo string `json:"logo"`
	// 组织支持的功能信息
	Functions []FunctionLimitObj `json:"functions"`
	// 部署方式：public：公共；private：私有化部署
	// 注意：“移动”部署时，需要叮嘱运维修改移动那边的 nacos 配置上的 RunMode 为 3
	AppDeployType string `json:"appDeployType"`
	// 距离使用付费产品截止日期的剩余天数
	RemainDays uint `json:"remainDays"`
	// 是否是试用；true 试用；false 非试用
	IsEvaluate bool `json:"isEvaluate"`
	// 付费过双11 6.6折的 团队标识  2022-11-11  活动期间 付费为1，没有付费为2，活动结束为0，字段不序列化
	IsPayActivity11 int    `json:"isPayActivity11,omitempty"`
	OutOrgId        string `json:"outOrgId"`
}

type BasicShowSetting struct {
	// 工作台
	WorkBenchShow bool `json:"workBenchShow"`
	// 侧边栏
	SideBarShow bool `json:"sideBarShow"`
	// 镜像统计
	MirrorStat bool `json:"mirrorStat"`
}

type SendCardToAdminForUpgradeReqVo struct {
	OrgId         int64  `json:"orgId"`
	UserId        int64  `json:"userId"`
	SourceChannel string `json:"sourceChannel"`
	// Input         SendCardToAdminForUpgradeReqData `json:"input"`
}

type SendCardToAdminForUpgradeReqData struct {
	TargetPayLevel int `json:"targetPayLevel"`
}

type SetInteractiveCardCallbackUrl struct {
}

type PrepareTransferOrgReq struct {
	OrgId  int64                    `json:"orgId"`
	UserId int64                    `json:"userId"`
	Input  *PrepareTransferOrgInput `json:"input"`
}

type PrepareTransferOrgInput struct {
	OriginOrgId int64 `json:"originOrgId"`
	ToOrgId     int64 `json:"toOrgId"`
}

type PrepareTransferOrgRespVo struct {
	vo.Err
	Data *PrepareTransferOrgData `json:"data"`
}

type PrepareTransferOrgData struct {
	MatchUsers    []*UserInfo `json:"matchUsers"`
	NotMatchUsers []*UserInfo `json:"notNatchUsers"`
}

type GetCorpAgentIdReq struct {
	CorpId string `json:"corp_id"`
}

type GetCorpAgentIdResp struct {
	AgentId string `json:"agent_id"`
}

type GetOrgInfoReq struct {
	OrgId int64 `json:"orgId"`
}

type GetOrgInfoResp struct {
	vo.Err
	Data *bo.OrganizationBo `json:"data"`
}

type PayBaseLevelData struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Duration int64  `json:"duration"`
}
