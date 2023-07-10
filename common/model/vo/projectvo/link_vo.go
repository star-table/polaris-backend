package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

// GetIssueLinks
type GetIssueLinksReqVo struct {
	SourceChannel string `json:"sourceChannel"`
	OrgId         int64  `json:"orgId"`
	IssueId       int64  `json:"issueId"`
}

type IssueLinks struct {
	Link        string `json:"link"`        // 应用内查看链接
	SideBarLink string `json:"sideBarLink"` // 侧边栏链接
}

type GetIssueLinksRespVo struct {
	vo.Err
	Data *IssueLinks `json:"data"`
}
