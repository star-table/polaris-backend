package othervo

type AddScriptTaskReq struct {
	OrgId    int64       `json:"orgId"`
	UserId   int64       `json:"userId"`
	AppId    string      `json:"appId"`
	TaskId   string      `json:"taskId"`
	TaskName string      `json:"taskName"`
	CronSpec string      `json:"cronSpec"`
	ParamMap interface{} `json:"paramMap"` // 使用时，使用 map[string]string 类型
}
