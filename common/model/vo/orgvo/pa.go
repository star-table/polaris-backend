package orgvo

type PAReportMsg struct {
	Secret       string `json:"secret"`
	UserCount    int64  `json:"userCount"`
	DeptCount    int64  `json:"deptCount"`
	ProjectCount int64  `json:"projectCount"`
	IssueCount   int64  `json:"issueCount"`
	OutSideIp    string `json:"outSideIp"`
}

type PAReportMsgReqVo struct {
	Body PAReportMsg `json:"body"`
}
