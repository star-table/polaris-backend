package commonvo

import (
	pushPb "github.com/star-table/interface/golang/push/v1"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
)

type CardField struct {
	Level int
	Key   string
	Value string // 如果value为空则整行展示key即可
}

type CardDiv struct {
	Fields []*CardField
}

type CardMeta struct {
	IsWide           bool
	Level            int
	Title            string
	Divs             []*CardDiv
	ActionMarkdowns  []string
	FsActionElements []interface{}
}

type Card struct {
	Title            string
	MarkdownElements []string
	FsCard           interface{}
}

type CardElements struct {
	MarkdownElements []string
	FsElements       []interface{}
}

// 微信卡片
type WxCard struct {
	Content string
}

// 钉钉卡片
type DtCard struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PushCard struct {
	OrgId         int64    `json:"orgId"`
	OutOrgId      string   `json:"outOrgId"`
	SourceChannel string   `json:"sourceChannel"`
	OpenIds       []string `json:"openIds"`
	DeptIds       []string `json:"deptIds"`
	ChatIds       []string `json:"chatIds"`
	CardMsg       *pushPb.TemplateCard
}

type GenerateCard struct {
	OrgId             int64                                 `json:"orgId"`
	CardTitle         string                                `json:"cardTitle"`
	CardLevel         pushPb.CardElementLevel               `json:"cardLevel"`
	ParentTitle       string                                `json:"parentTitle"`
	ProjectName       string                                `json:"projectName"`
	TableName         string                                `json:"tableName"`
	IssueTitle        string                                `json:"issueTitle"`
	ProjectId         int64                                 `json:"projectId"`
	ProjectTypeId     int64                                 `json:"projectTypeId"`
	ParentId          int64                                 `json:"parentId"`
	OwnerInfos        []*bo.BaseUserInfoBo                  `json:"ownerInfos"`
	Links             *projectvo.IssueLinks                 `json:"links"`
	IsNeedUnsubscribe bool                                  `json:"isNeedUnsubscribe"`
	SpecialDiv        []*pushPb.CardTextDiv                 `json:"specialDiv"`
	SpecialAction     []*pushPb.CardActionModule            `json:"specialAction"`
	TableColumnInfo   map[string]*projectvo.TableColumnData `json:"tableColumnInfo"`
}
