package bo

import "github.com/star-table/common/core/types"

type IssueDetailBo struct {
	Id           int64      `json:"id"`
	OrgId        int64      `json:"orgId"`
	IssueId      int64      `json:"issueId"`
	ProjectId    int64      `json:"projectId"`
	StoryPoint   int        `json:"storyPoint"`
	Tags         string     `json:"tags"`
	Remark       string     `json:"remark"`
	RemarkDetail string     `json:"remarkDetail"`
	Status       int64      `json:"status"`
	Creator      int64      `json:"creator"`
	CreateTime   types.Time `json:"createTime"`
	Updator      int64      `json:"updator"`
	UpdateTime   types.Time `json:"updateTime"`
	Version      int        `json:"version"`
}
