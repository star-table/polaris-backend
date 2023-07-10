package formvo

import (
	"time"

	tablev1 "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/polaris-backend/common/model/vo/uservo"

	"github.com/star-table/polaris-backend/common/model/vo/datacenter"

	"github.com/star-table/polaris-backend/common/model/vo"
)

type FormCreateOneResp struct {
	vo.Err
	Timestamp interface{}              `json:"timestamp"`
	Data      []map[string]interface{} `json:"data"`
}

type LessCreateIssueReq struct {
	AppId          int64                    `json:"appId"`
	OrgId          int64                    `json:"orgId"`
	UserId         int64                    `json:"userId"`
	Asc            bool                     `json:"asc"`
	BeforeId       int64                    `json:"beforeId"`
	AfterId        int64                    `json:"afterId"`
	RedirectIds    []int64                  `json:"redirectIds"`
	Import         bool                     `json:"import"`
	CreateTemplate bool                     `json:"createTemplate"`
	Form           []map[string]interface{} `json:"form"`
	TableId        int64                    `json:"tableId"`
}

type LessUpdateIssueReq struct {
	AppId   int64                    `json:"appId"`
	OrgId   int64                    `json:"orgId"`
	UserId  int64                    `json:"userId"`
	TableId int64                    `json:"tableId"`
	Form    []map[string]interface{} `json:"form"`
}

type LessUpdateIssueBatchReq struct {
	OrgId     int64            `json:"orgId"`
	AppId     int64            `json:"appId"`
	UserId    int64            `json:"userId"`
	Condition vo.LessCondsData `json:"condition"`
	Sets      []datacenter.Set `json:"sets"`
}

type LessDeleteIssueReq struct {
	AppId       int64   `json:"appId"`
	OrgId       int64   `json:"orgId"`
	UserId      int64   `json:"userId"`
	AppValueIds []int64 `json:"appValueIds"`
	IssueIds    []int64 `json:"issueIds"`
}

type LessCommonIssueResp struct {
	vo.Err
	Timestamp interface{} `json:"timestamp"`
	Data      bool        `json:"data"`
}

type LessUpdateIssueResp struct {
	vo.Err
	Timestamp interface{}              `json:"timestamp"`
	Data      []map[string]interface{} `json:"data"`
}

type LessIssueListResp struct {
	vo.Err
	Data      *tablev1.ListReply            `json:"data"`
	List      []map[string]interface{}      `json:"list"`
	UserDepts map[string]*uservo.MemberDept `json:"userDepts"`
}

type LessIssueRawListResp struct {
	vo.Err
	Timestamp interface{}              `json:"timestamp"`
	Data      []map[string]interface{} `json:"data"`
}

type LessFilterStatResp struct {
	vo.Err
	Timestamp interface{} `json:"timestamp"`
	Data      int64       `json:"data"`
}

type LessFilterCustomStatResp struct {
	vo.Err
	Timestamp interface{}              `json:"timestamp"`
	Data      []map[string]interface{} `json:"data"`
}

type GetFormByAppIdResp struct {
	vo.Err
	Timestamp interface{}        `json:"timestamp"`
	Data      GetFormByAppIdData `json:"data"`
}

type GetFormByAppIdData struct {
	ID         int64     `json:"id,string"`
	OrgID      int64     `json:"orgId,string"`
	AppID      int64     `json:"appId,string"`
	Type       int       `json:"type"`
	Config     string    `json:"config"`
	Status     int       `json:"status"`
	Creator    int64     `json:"creator,string"`
	CreateTime time.Time `json:"createTime"`
	Updator    int64     `json:"updator,string"`
	UpdateTime time.Time `json:"updateTime"`
}

type LessIssueListData struct {
	Total int64                    `json:"total"`
	List  []map[string]interface{} `json:"list"`
}

type Condition struct {
	Type   string      `json:"type"`
	Value  interface{} `json:"value"`
	Column string      `json:"column"`
	Left   interface{} `json:"left"`
	Right  interface{} `json:"right"`
	Conds  []Condition `json:"conds"`
}

type LessIssueListReq struct {
	Condition      vo.LessCondsData `json:"condition"`
	Orders         []*vo.LessOrder  `json:"orders"`
	RedirectIds    []int64          `json:"redirectIds"`
	AppId          int64            `json:"appId"`
	OrgId          int64            `json:"orgId"`
	UserId         int64            `json:"userId"`
	Page           int64            `json:"page"`
	Size           int64            `json:"size"`
	Columns        []string         `json:"columns"`
	FilterColumns  []string         `json:"filterColumns"`
	Groups         []string         `json:"groups"`
	TableId        int64            `json:"tableId"`
	Export         bool             `json:"export"`              // 是否是导出，导出不限制size
	NeedTotal      bool             `json:"needTotal,omitempty"` // 是否需要总数，大部分情况都不需要总数
	NeedRefColumn  bool             `json:"needRefColumn"`       // 是否需要引用列
	AggNoLimit     bool             `json:"aggNoLimit"`          // 引用列是否限制返回合并个数
	NeedDeleteData bool             `json:"needDeleteData"`      // 需要删除的数据
	NeedChangeId   bool             `json:"-"`
}

func (l *LessIssueListReq) Copy() *LessIssueListReq {
	return &LessIssueListReq{
		Condition:      l.Condition,
		Orders:         l.Orders,
		RedirectIds:    l.RedirectIds,
		AppId:          l.AppId,
		OrgId:          l.OrgId,
		UserId:         l.UserId,
		Page:           l.Page,
		Size:           l.Size,
		Columns:        l.Columns,
		FilterColumns:  l.FilterColumns,
		Groups:         l.Groups,
		TableId:        l.TableId,
		Export:         l.Export,
		NeedTotal:      l.NeedTotal,
		NeedRefColumn:  l.NeedRefColumn,
		NeedDeleteData: l.NeedDeleteData,
	}
}

type GetFormByAppIdReq struct {
	OrgID int64 `json:"orgId"`
	AppID int64 `json:"appId"`
}

type LessRecycleIssueReq struct {
	AppId    int64   `json:"appId"`
	OrgId    int64   `json:"orgId"`
	UserId   int64   `json:"userId"`
	IssueIds []int64 `json:"issueIds"`
	DataIds  []int64 `json:"dataIds"`
	TableId  int64   `json:"tableId"`
}

type LessRecoverIssueReq struct {
	AppId       int64   `json:"appId"`
	OrgId       int64   `json:"orgId"`
	UserId      int64   `json:"userId"`
	AppValueIds []int64 `json:"appValueIds"`
	IssueIds    []int64 `json:"issueIds"`
	TableId     int64   `json:"tableId"`
}

type LessSaveFormReq struct {
	AppId     int64  `json:"appId"`
	OrgId     int64  `json:"orgId"`
	UserId    int64  `json:"userId"`
	ExtendsId int64  `json:"extendsId"`
	PkgId     int64  `json:"pkgId"`
	Config    string `json:"config"`
	Type      int64  `json:"type"`
}

type LessFormConfigData struct {
	Fields       []interface{}            `json:"fields"`
	BaseFields   []string                 `json:"baseFields"`
	CustomConfig map[string][]interface{} `json:"customConfig"`
}

type LessBaseFormSaveReq struct {
	OrgId   int64         `json:"orgId"`
	UserId  int64         `json:"userId"`
	Added   []interface{} `json:"added"`
	Deleted []string      `json:"deleted"`
}

type LessRecycleIssueResp struct {
	vo.Err
	Timestamp interface{} `json:"timestamp"`
	Data      int64       `json:"data"`
}

type LessRecoverIssueResp struct {
	vo.Err
	Data []map[string]interface{} `json:"data"`
}

type GetFormConfigReq struct {
	OrgID int64 `json:"orgId"`
	AppID int64 `json:"appId"`
}

type LessMoveIssueReq struct {
	AppId    int64  `json:"appId"`
	BeforeId *int64 `json:"beforeId"`
	AfterId  *int64 `json:"afterId"`
	Asc      bool   `json:"asc"`
	DataId   int64  `json:"dataId"`
	OrgId    int64  `json:"orgId"`
	UserId   int64  `json:"userId"`
}

type LessMoveIssueResp struct {
	vo.Err
	Timestamp interface{} `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type GetFormConfigBatchReq struct {
	OrgID  int64   `json:"orgId"`
	AppIds []int64 `json:"appIds"`
}

type LessWorkHour struct {
	PlanHour        string   `json:"planHour"`
	ActualHour      string   `json:"actualHour"`
	CollaboratorIds []string `json:"collaboratorIds"`
}
