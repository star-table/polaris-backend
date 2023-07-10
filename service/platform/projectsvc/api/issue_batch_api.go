package api

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) BatchCreateIssue(reqVo *projectvo.BatchCreateIssueReqVo) *projectvo.BatchCreateIssueRespVo {
	data, userDepts, relateData, err := service.BatchCreateIssue(reqVo, false, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	}, "", 0)
	return &projectvo.BatchCreateIssueRespVo{Data: data, UserDepts: userDepts, RelateData: relateData, Err: vo.NewErr(err)}
}

func (PostGreeter) BatchUpdateIssue(reqVo *projectvo.BatchUpdateIssueReqVo) *vo.VoidErr {
	err := service.BatchUpdateIssue(&projectvo.BatchUpdateIssueReqInnerVo{
		OrgId:        reqVo.OrgId,
		UserId:       reqVo.UserId,
		AppId:        reqVo.Input.AppId,
		ProjectId:    reqVo.Input.ProjectId,
		TableId:      reqVo.Input.TableId,
		Data:         reqVo.Input.Data,
		BeforeDataId: reqVo.Input.BeforeDataId,
		AfterDataId:  reqVo.Input.AfterDataId,
		TodoId:       reqVo.Input.TodoId,
		TodoOp:       reqVo.Input.TodoOp,
		TodoMsg:      reqVo.Input.TodoMsg,
	}, false, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	})
	return &vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) BatchAuditIssue(reqVo *projectvo.BatchAuditIssueReqVo) *vo.DataRespVo {
	issueIds, err := service.BatchAuditIssue(reqVo)
	return &vo.DataRespVo{Data: issueIds, Err: vo.NewErr(err)}
}

func (PostGreeter) BatchUrgeIssue(reqVo *projectvo.BatchUrgeIssueReqVo) *vo.DataRespVo {
	issueIds, err := service.BatchUrgeIssue(reqVo)
	return &vo.DataRespVo{Data: issueIds, Err: vo.NewErr(err)}
}
