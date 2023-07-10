package projectvo

import (
	automationPb "github.com/star-table/interface/golang/automation/v1"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type TodoCard struct {
	TodoType     automationPb.TodoType       `json:"todoType"`
	Trigger      *bo.BaseUserInfoBo          `json:"trigger"`
	Owners       []*bo.BaseUserInfoBo        `json:"owners"`
	Issue        *bo.IssueBo                 `json:"issue"`
	Project      *bo.ProjectBo               `json:"project"` // 这里如果projectId为0，appId会补上汇总表id
	Table        *TableMetaData              `json:"table"`
	TableColumns map[string]*TableColumnData `json:"tableColumns"`
	IssueInfoUrl string                      `json:"issueInfoUrl"` // “查看详情”按钮链接
	IssuePcUrl   string                      `json:"issuePcUrl"`   // "PC端查看/应用内查看"按钮链接
}

// 催办的结构体
type UrgeIssueCard struct {
	OrgId           int64
	OperateUserId   int64
	OperateUserName string
	OutOrgId        string
	OpenIds         []string
	ProjectName     string
	TableName       string
	ProjectId       int64
	ProjectTypeId   int64
	IssueId         int64
	ParentTitle     string
	UrgeText        string
	IssueBo         *bo.IssueBo
	OwnerInfos      []*bo.BaseUserInfoBo // 负责人的信息
	IssueLinks      *IssueLinks
	TableColumn     map[string]*TableColumnData
}

// 回复催办的结构体
type UrgeReplyCard struct {
	OrgId           int64
	OutOrgId        string
	OpenIds         []string
	ProjectName     string
	TableName       string
	OperateUserName string
	BeUrgedUserName string
	UrgeText        string
	ReplyMsg        string
	ParentTitle     string
	AppId           int64
	ProjectTypeId   int64
	Issue           vo.Issue
	OwnerInfos      []*bo.BaseUserInfoBo
	IssueLinks      *IssueLinks
	TableColumn     map[string]*TableColumnData
}

type AsyncTaskCard struct {
	Project         *bo.ProjectBo
	Table           *TableMetaData
	User            *bo.BaseUserInfoBo
	IsApplyTemplate bool
	Link            string // “查看详情”按钮链接，进入的是项目的任务列表页
}

type CardBeOverdue struct {
	OrgId          int64
	OwnerId        int64
	IssueId        int64
	IssueLinks     *IssueLinks
	IssueTitle     string
	PlanEndTime    string
	SourceChannel  string
	TableId        string
	AppId          string
	ContentProject string
	OwnerNameStr   string
	ParentTitle    string
	HasPermission  bool
	Tips           string
	TableColumn    map[string]*TableColumnData
}
