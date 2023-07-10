package status

// 项目/迭代/任务的状态数据结构
type StatusInfoBo struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	BgStyle     string `json:"bgStyle"`
	FontStyle   string `json:"fontStyle"`
	Type        int    `json:"type"`
	Sort        int    `json:"sort"`
}

type TaskBarInfoBo struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}
