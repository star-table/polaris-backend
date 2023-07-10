package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

func (PostGreeter) GetIssueLinks(reqVo projectvo.GetIssueLinksReqVo) projectvo.GetIssueLinksRespVo {
	return projectvo.GetIssueLinksRespVo{
		Err:  vo.NewErr(nil),
		Data: domain.GetIssueLinks(reqVo.SourceChannel, reqVo.OrgId, reqVo.IssueId),
	}
}
