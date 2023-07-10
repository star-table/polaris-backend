package vo

type HomeIssuesRestReq struct {
	Condition     *LessCondsData `json:"condition"`
	Orders        []*LessOrder   `json:"orders"`
	FilterColumns []string       `json:"filterColumns"`
	MenuAppId     *string        `json:"menuAppId"`
	TableId       *string        `json:"tableId"`
	Page          *int           `json:"page"`
	Size          *int           `json:"size"`
}

type IssueDetailRestReq struct {
	AppId   *string `json:"appId"`
	TableId *string `json:"tableId"`
	IssueId int     `json:"issueId"`
}

type AppViewConfig struct {
	RealCondition *LessCondsData `json:"realCondition"`
	Orders        []*LessOrder   `json:"orders"`
	//ProjectObjectTypeId int64          `json:"projectObjectTypeId"`
	TableId int64 `json:"tableId"`
}

type LessCreateIssueReq struct {
	BeforeId     *int64                   `json:"beforeId"`
	AfterId      *int64                   `json:"afterId"`
	BeforeDataId *string                  `json:"beforeDataId"`
	AfterDataId  *string                  `json:"afterDataId"`
	Asc          *bool                    `json:"asc"`
	Form         []map[string]interface{} `json:"form"`
	MenuAppId    string                   `json:"menuAppId"`
	TableId      string                   `json:"tableId"`
	ProjectId    int64                    `json:"projectId"`
}

type LessDeleteIssueBatchReq struct {
	AppValueIds []int64 `json:"appValueIds"`
	MenuAppId   string  `json:"menuAppId"`
	TableId     string  `json:"tableId"`
}

type LessCopyIssueReq struct {
	OldIssueId   int64    `json:"oldIssueId"`
	ProjectId    int64    `json:"projectId"`
	TableId      int64    `json:"tableId,string"`
	BeforeDataId int64    `json:"beforeDataId,string"`
	AfterDataId  int64    `json:"afterDataId,string"`
	ChooseField  []string `json:"chooseField"`
	ChildrenIds  []int64  `json:"childrenIds"`
	Title        *string  `json:"title"`
	IsStaticCopy bool     `json:"isStaticCopy"`
}

type LessCopyIssueResp struct {
	Data []map[string]interface{} `json:"data"`
}

type AuditResult struct {
}

type LessUpdateIssueReq struct {
	Form         []map[string]interface{} `json:"form"`
	MenuAppId    int64                    `json:"menuAppId,string"`
	BeforeDataId int64                    `json:"beforeDataId,string"`
	AfterDataId  int64                    `json:"afterDataId,string"`
	TodoId       int64                    `json:"todoId,string"`
	TodoOp       int                      `json:"todoOp"`
	TodoMsg      string                   `json:"todoMsg"`
}

type LessBatchUpdateIssueReq struct {
	AppId   string                   `json:"appId"`
	TableId string                   `json:"tableId"`
	Data    []map[string]interface{} `json:"data"` // map[string]interface{}中id必传，和values保持一致，传issueId
}

// 首页任务列表响应结构体
type LessIssueFilterResp struct {
	// 总数量
	Total int64 `json:"total"`
	// 实际总数量
	ActualTotal int64 `json:"actualTotal"`
	// 首页任务列表
	List []map[string]interface{} `json:"list"`
}

type LessBatchAuditIssueReq struct {
	IssueIds    []int64                 `json:"issueIds"`
	AuditStatus int                     `json:"auditStatus"`
	Message     string                  `json:"message"`     // 评论
	Attachments []*AttachmentSimpleInfo `json:"attachments"` // 附件
}

type LessBatchUrgeIssueReq struct {
	IssueIds []int64 `json:"issueIds"`
	Message  string  `json:"message"` // 评论
}

type SysCommonField struct {
	Title               string         `json:"title"`
	ProjectId           int64          `json:"projectId"`
	Code                string         `json:"code"`
	OwnerId             []LcMemberInfo `json:"ownerId"`
	Status              int64          `json:"status"`
	ParentId            int64          `json:"parentId"`
	PlanStartTime       interface{}    `json:"planStartTime"`
	PlanEndTime         interface{}    `json:"planEndTime"`
	ProjectObjectTypeId int64          `json:"projectObjectTypeId"`
	AuditStatus         int            `json:"auditStatus"`
	FollowerIds         []LcMemberInfo `json:"followerIds"`
	AuditorIds          []LcMemberInfo `json:"auditorIds"`
	Tag                 interface{}    `json:"tag"`
}

type LcMemberInfo struct {
	Avatar   string `json:"avatar"`
	Id       string `json:"id"`
	IsDelete int    `json:"isDelete"`
	Name     string `json:"name"`
	Status   int    `json:"status"`
	Type     string `json:"type"`
}

type GetFormConfigReq struct {
	//通用项目可以传0
	ProjectObjectTypeId int64 `json:"projectObjectTypeId"`
}

type GetFormConfigResp struct {
	Fields       []map[string]interface{} `json:"fields"`
	Columns      []map[string]interface{} `json:"columns"`
	FieldOrders  interface{}              `json:"fieldOrders"`
	CustomConfig map[string][]interface{} `json:"customConfig"`
}

type SaveFormHeaderData struct {
	ProjectId                 int64                    `json:"projectId"`
	Name                      string                   `json:"name"`
	Config                    []interface{}            `json:"config"`
	FieldOrders               []string                 `json:"fieldOrders"`
	ViewOrders                []string                 `json:"viewOrders"`
	BaseFields                []string                 `json:"baseFields"`
	CustomConfig              map[string][]interface{} `json:"customConfig"`
	ProjectObjectTypeId       int64                    `json:"projectObjectTypeId"`
	IsUpdateIssueStatus       *bool                    `json:"isUpdateIssueStatus"`
	MenuAppId                 string                   `json:"menuAppId"`
	IsUpdateProjectObjectType *bool                    `json:"isUpdateProjectObjectType"`
}

type SaveFormHeaderRespData struct {
	AppId     int64  `json:"appId"`
	Config    string `json:"config"`
	Drafted   bool   `json:"drafted"`
	ExtendsId int64  `json:"extendsId"`
	IsExt     bool   `json:"isExt"`
	OrgId     int64  `json:"orgId"`
	Type      int    `json:"type"`
	UserId    int64  `json:"userId"`
}

type MirrorCountReq struct {
	AppIds []string `json:"appIds"`
}

type MirrorsStatResp struct {
	DataStat map[int64]int64 `json:"dataStat"`
}

type CreateColumnData struct {
}
