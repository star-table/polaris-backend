package schedulevo

import "github.com/star-table/polaris-backend/common/model/vo"

type AddTaskReqVo struct {
	Data AddTaskReq `json:"data"`
}

type AddTaskReq struct {
	AppId    string                 `json:"appId"`
	TaskId   string                 `json:"taskId"`
	TaskName string                 `json:"taskName"`
	CronSpec string                 `json:"cronSpec"`
	ParamMap map[string]interface{} `json:"paramMap"`
}

type AddTaskRespVo struct {
	Data *BoolResp `json:"data"`
}

type BoolResp struct {
	// 是否符合期望、确定、ok：true 表示成功、是、确定；false 表示否定、异常
	IsTrue bool `json:"isTrue"`
}

type BoolRespVo struct {
	vo.Err
	Data *vo.BoolResp `json:"data"`
}
