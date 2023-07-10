package bo

// 任务标签结构体
type IssueTagReqBo struct {
	// 标签id
	Id int64 `json:"id"`
	// 标签名称
	Name string `json:"name"`
}