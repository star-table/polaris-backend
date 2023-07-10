package bo

import "github.com/star-table/polaris-backend/common/model/vo"

type AttachmentBo struct {
	vo.Resource
	IssueList []IssueBo
}

type Attachments struct {
	Id           int64  `json:"id"`
	Name         string `json:"name"`
	Path         string `json:"path"`
	Suffix       string `json:"suffix"`
	ResourceType int    `json:"resourceType"`
	RecycleFlag  int    `json:"recycleFlag"`
}

type IssueAttachments struct {
	ColumnId    string  `json:"columnId"`
	ResourceIds []int64 `json:"resourceIds"`
}
