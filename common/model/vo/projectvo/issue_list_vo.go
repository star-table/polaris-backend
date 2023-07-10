package projectvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
)

type IssueDetailReqVo struct {
	OrgId   int64 `json:"orgId"`
	UserId  int64 `json:"userId"`
	AppId   int64 `json:"appId"`
	TableId int64 `json:"tableId"`
	IssueId int64 `json:"issueId"`
}

type IssueDetailInfo struct {
	Project   *ProjectMetaInfo              `json:"project"`
	Table     *TableSimpleInfo              `json:"table"`
	Data      []map[string]interface{}      `json:"data"`
	UserDepts map[string]*uservo.MemberDept `json:"userDepts"`
}

type IssueDetailRespVo struct {
	vo.Err
	Data *IssueDetailInfo `json:"data"`
}

type HomeIssuesReqVo struct {
	Page   int               `json:"page"`
	Size   int               `json:"size"`
	Input  *HomeIssueInfoReq `json:"input"`
	UserId int64             `json:"userId"`
	OrgId  int64             `json:"orgId"`
}

type HomeIssuesRespVo struct {
	vo.Err
	HomeIssueInfo *vo.HomeIssueInfoResp `json:"data"`
}

type LcViewStatReqVo struct {
	UserId int64 `json:"userId"`
	OrgId  int64 `json:"orgId"`
}

type LcViewStatVo struct {
	Id    int64  `json:"id,string"`
	Name  string `json:"name"`
	Total int64  `json:"total"`
}

type LcViewStatRespVo struct {
	vo.Err
	Data []*LcViewStatVo `json:"data"`
}

type LcHomeIssuesRespVo struct {
	vo.Err
	Data string `json:"data"`
}

type LcHomeIssueInfoResp struct {
	// 总数量
	Total int64 `json:"total"`
	// 实际总数量
	ActualTotal int64 `json:"actualTotal"`
	// 数据列表
	List      []map[string]interface{}      `json:"list"`
	UserDepts map[string]*uservo.MemberDept `json:"userDepts"`
	// 不支持的引用列，由于数据太多
	UnSupportRefColumnIds    []string `json:"unSupportRefColumnIds"`
	LastUpdateTime           string   `json:"lastUpdateTime"`
	LastPermissionUpdateTime string   `json:"lastPermissionUpdateTime"`
	// 是否是刷新全部数据
	IsRefreshAll bool `json:"isRefreshAll"`
}

type IssueReportReqVo struct {
	ReportType int64 `json:"reportType"`
	UserId     int64 `json:"userId"`
	OrgId      int64 `json:"orgId"`
}

type IssueReportRespVo struct {
	vo.Err
	IssueReport *vo.IssueReportResp `json:"data"`
}

type IssueReportDetailReqVo struct {
	ShareID string `json:"shareId"`
}

type IssueReportDetailRespVo struct {
	vo.Err
	IssueReportDetail *vo.IssueReportResp `json:"data"`
}

type IssueStatusTypeStatRespVo struct {
	vo.Err
	IssueStatusTypeStat *vo.IssueStatusTypeStatResp `json:"data"`
}

type IssueStatusTypeStatDetailRespVo struct {
	vo.Err
	IssueStatusTypeStatDetail *vo.IssueStatusTypeStatDetailResp `json:"data"`
}

type IssueStatusTypeStatReqVo struct {
	Input  *vo.IssueStatusTypeStatReq `json:"input"`
	UserId int64                      `json:"userId"`
	OrgId  int64                      `json:"orgId"`
}

type GetIssueStatusStatReqVo struct {
	Input bo.IssueStatusStatCondBo `json:"input"`
}

type GetIssueStatusStatRespVo struct {
	vo.Err
	Data GetIssueStatusStatRespData `json:"data"`
}

type GetIssueStatusStatRespData struct {
	List []bo.IssueStatusStatBo `json:"list"`
}

type IssueStatusTypeStatDetailReqVo struct {
	Input *vo.IssueStatusTypeStatReq `json:"input"`
}

type GetSimpleIssueInfoBatchReqVo struct {
	OrgId int64   `json:"orgId"`
	Ids   []int64 `json:"ids"`
}

type GetSimpleIssueInfoBatchRespVo struct {
	vo.Err
	Data *[]vo.Issue `json:"data"`
}

type GetLcIssueInfoBatchReqVo struct {
	OrgId    int64   `json:"orgId"`
	IssueIds []int64 `json:"issueIds"`
}

type GetLcIssueInfoBatchRespVo struct {
	vo.Err
	Data []*bo.IssueBo `json:"data"`
}

type GetIssueRemindInfoListReqVo struct {
	Page  int                           `json:"page"`
	Size  int                           `json:"size"`
	Input GetIssueRemindInfoListReqData `json:"input"`
}

type GetIssueRemindInfoListRespVo struct {
	Data *GetIssueRemindInfoListRespData `json:"data"`
	vo.Err
}

type GetIssueRemindInfoListRespData struct {
	Total int64                  `json:"total"`
	List  []bo.IssueRemindInfoBo `json:"issueRemindInfoList"`
}

type GetIssueRemindInfoListReqData struct {
	//计划截止时间开始范围
	BeforePlanEndTime *string `json:"beforePlanEndTime"`
	//计划截止时间结束范围
	AfterPlanEndTime *string `json:"afterPlanEndTime"`
	//计划开始时间开始范围
	BeforePlanStartTime *string `json:"beforePlanStartTime"`
	//计划开始时间结束范围
	AfterPlanStartTime *string `json:"afterPlanStartTime"`
}

type IssueListStatReq struct {
	UserId int64               `json:"userId"`
	OrgId  int64               `json:"orgId"`
	Input  vo.IssueListStatReq `json:"input"`
}

type IssueListStatResp struct {
	vo.Err
	Data *vo.IssueListStatResp `json:"data"`
}

type IssueViewLimitInfoForExportQuery struct {
	IsLimited     bool    `json:"isLimited"`
	AllowIssueIds []int64 `json:"allowIssueIds"`
	Page          int     `json:"page"` // 查询时分页数据
	Size          int     `json:"size"`
}

// 首页的任务列表请求结构体
type HomeIssueInfoReq struct {
	ProjectID *int64 `json:"projectId"`
	// 无码格式
	LessConds *vo.LessCondsData `json:"condition"`
	// 无码格式排序
	LessOrder []*vo.LessOrder `json:"orders"`
	// 镜像应用id
	MenuAppID *string `json:"menuAppId"`
	// 表Id
	TableID *string `json:"tableId"`
	// 需要查询的无码字段（id不需要传，默认会查）
	FilterColumns []string `json:"filterColumns"`
	// 是否全刷，如果是，忽略删除的和回收站的，如果是增量获取，则要获取到删除、移动的数据
	IsRefreshAll bool `json:"isRefreshAll"`
	// 最后更新时间
	LastUpdateTime string `json:"lastUpdateTime"`
	// 最后的权限更新时间
	LastPermissionUpdateTime string `json:"lastPermissionUpdateTime"`
	// 是否需要全部的数据，因为默认forAll的情况过滤了空应用数据，导致前后置显示不出来
	IsNeedAll bool `json:"isNeedAll"`

	Page *int `json:"page"`
	Size *int `json:"size"`
}
