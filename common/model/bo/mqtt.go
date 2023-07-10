package bo

type MQTTNoticeBo struct {
	//具体消息类型（跟动态类型一致）
	//NoticeType int `json:"noticeType"`
	//类型：1、提示通知(pc弹框)，2：数据刷新
	Type int `json:"type"`
	//消息体
	Body interface{} `json:"body"`
}

type MQTTRemindNotice struct {
	OperatorId   int64       `json:"operatorId"`
	OperatorName string      `json:"operatorName"`
	OrgID        int64       `json:"orgId"`
	Content      string      `json:"content"`
	Data         interface{} `json:"data"`
}

type MQTTDataRefreshNotice struct {
	//组织id
	OrgId int64 `json:"orgId"`
	//项目id，为0说明跟项目无关
	ProjectId int64 `json:"projectId"`
	//应用id
	AppId int64 `json:"appId"`
	//动作, ADD，新增，DEL，删除，MODIFY，变动
	Action string `json:"action"`
	//操作人
	OperationId int64 `json:"operationId"`
	//类型
	Type string `json:"type"`
	//局部数据
	PartialRefresh []MQTTPartialRefresh `json:"partialRefresh"`
	//全量刷新
	GlobalRefresh []MQTTGlobalRefresh `json:"globalRefresh"`
	//移动
	MoveRefresh MoveRefresh `json:"moveRefresh"`
}

type MoveRefresh struct {
	Old MQTTGlobalRefresh   `json:"old"`
	New []MQTTGlobalRefresh `json:"new"`
}

type MQTTPartialRefresh struct {
	//对象id
	ObjectId int64 `json:"objectId"`
	//field
	Fields []ObjectField `json:"fields"`
}

type MQTTGlobalRefresh struct {
	//对象id
	ObjectId int64 `json:"objectId"`
	//对象内容
	ObjectValue interface{} `json:"objectValue"`
	//新对象内容
	NewObjectValue interface{} `json:"newObjectValue"`
	//子任务id
	ChildrenIssueIds []int64 `json:"childrenIssueIds"`
}

type ObjectField struct {
	Field string      `json:"field"`
	Value interface{} `json:"value"`
	//新对象内容
	NewObjectValue interface{} `json:"newObjectValue"`
}

type MqttIssueInfo struct {
	Info MqttIssueLessData `json:"info"`
}

type MqttIssueLessData struct {
	LessData map[string]interface{} `json:"lessData"`
}
