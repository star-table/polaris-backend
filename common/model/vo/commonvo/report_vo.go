package commonvo

type Event struct {
	Category  string      `json:"category"`
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp,string"`
	TraceId   string      `json:"traceId"`
	Payload   interface{} `json:"payload"`
}

type OrgEvent struct {
	OrgId int64       `json:"orgId"`
	New   interface{} `json:"new"`
}

type DataEvent struct {
	OrgId          int64                  `json:"orgId"`
	AppId          int64                  `json:"appId,string"`
	ProjectId      int64                  `json:"projectId"`
	TableId        int64                  `json:"tableId,string"`
	DataId         int64                  `json:"dataId,string"`
	IssueId        int64                  `json:"issueId"`
	UserId         int64                  `json:"userId"`
	New            map[string]interface{} `json:"new,omitempty"`
	Old            map[string]interface{} `json:"old,omitempty"`
	Incremental    map[string]interface{} `json:"incremental,omitempty"`
	Decremental    map[string]interface{} `json:"decremental,omitempty"`
	UpdatedColumns []string               `json:"updatedColumns,omitempty"`
	UserDepts      interface{}            `json:"userDepts,omitempty"`
	TriggerBy      string                 `json:"triggerBy"`
}

type TableEvent struct {
	OrgId     int64       `json:"orgId"`
	AppId     int64       `json:"appId,string"`
	ProjectId int64       `json:"projectId"`
	TableId   int64       `json:"tableId,string"`
	UserId    int64       `json:"userId"`
	New       interface{} `json:"new,omitempty"`
	Old       interface{} `json:"old,omitempty"`
}

type AppEvent struct {
	OrgId     int64       `json:"orgId"`
	AppId     int64       `json:"appId,string"`
	ProjectId int64       `json:"projectId"`
	UserId    int64       `json:"userId"`
	App       interface{} `json:"app,omitempty"`
	Project   interface{} `json:"project,omitempty"`
	Chat      interface{} `json:"chat,omitempty"`
}

type UserEvent struct {
	OrgId  int64       `json:"orgId"`
	UserId int64       `json:"userId"`
	New    interface{} `json:"new,omitempty"`
}

type ReportAppEventReq struct {
	EventType int32     `json:"eventType"`
	TraceId   string    `json:"traceId"`
	AppEvent  *AppEvent `json:"appEvent"`
}

type ReportTableEventReq struct {
	EventType  int32       `json:"eventType"`
	TraceId    string      `json:"traceId"`
	TableEvent *TableEvent `json:"tableEvent"`
}
