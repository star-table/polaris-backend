package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) IssueSourceList(reqVo projectvo.IssueSourceListReqVo) projectvo.IssueSourceListRespVo {
	res, err := service.IssueSourceList(reqVo.OrgId, reqVo.Page, reqVo.Size, reqVo.Input)
	return projectvo.IssueSourceListRespVo{Err: vo.NewErr(err), IssueSourceList: res}
}

func (PostGreeter) CreateIssueSource(reqVo projectvo.CreateIssueSourceReqVo) vo.CommonRespVo {
	res, err := service.CreateIssueSource(reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateIssueSource(reqVo projectvo.UpdateIssueSourceReqVo) vo.CommonRespVo {
	res, err := service.UpdateIssueSource(reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) DeleteIssueSource(reqVo projectvo.DeleteIssueSourceReqVo) vo.CommonRespVo {
	res, err := service.DeleteIssueSource(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}
