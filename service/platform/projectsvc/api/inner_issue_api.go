package api

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
	"github.com/spf13/cast"
)

func (PostGreeter) InnerIssueFilter(reqVo *projectvo.InnerIssueFilterReq) string {
	appId := cast.ToString(reqVo.Input.AppId)
	tableId := cast.ToString(reqVo.Input.TableId)
	_, projectId, err := domain.GetProjectIdByAppId(appId)
	if err != nil {
		return json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"})
	}
	res, err := service.LcHomeIssuesForProject(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, &projectvo.HomeIssueInfoReq{
		MenuAppID:     &appId,
		ProjectID:     &projectId,
		TableID:       &tableId,
		FilterColumns: reqVo.Input.Columns,
		LessConds:     reqVo.Input.Condition,
		LessOrder:     reqVo.Input.Orders,
	}, true)
	if err != nil {
		return json.ToJsonIgnoreError(&projectvo.LcHomeIssuesRespVo{Err: vo.NewErr(err), Data: "{}"})
	}
	return res.Data
}

func (PostGreeter) InnerIssueCreate(reqVo *projectvo.InnerIssueCreateReq) *projectvo.LcDataListRespVo {
	_, projectId, err := domain.GetProjectIdByAppId(cast.ToString(reqVo.Input.AppId))
	if err != nil {
		return &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: nil}
	}
	req := &projectvo.BatchCreateIssueReqVo{
		OrgId:     reqVo.OrgId,
		UserId:    reqVo.UserId,
		AppId:     reqVo.Input.AppId,
		ProjectId: projectId,
		TableId:   reqVo.Input.TableId,
		Data:      reqVo.Input.Data,
	}
	res, userDept, relateData, err := service.SyncBatchCreateIssue(req, true, reqVo.Input.TriggerBy)
	return &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: res, UserDept: userDept, RelateData: relateData}
}

func (PostGreeter) InnerIssueCreateByCopy(reqVo *projectvo.InnerIssueCreateByCopyReq) *projectvo.LcDataListRespVo {
	_, projectId, err := domain.GetProjectIdByAppId(cast.ToString(reqVo.Input.AppId))
	if err != nil {
		return &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: nil}
	}
	res, userDept, relateData, err := service.CopyIssueBatchWithData(reqVo.OrgId, reqVo.UserId, projectId, cast.ToInt64(reqVo.Input.TableId),
		reqVo.Input.IssueIds, reqVo.Input.Data, reqVo.Input.TriggerBy, true, reqVo.Input.IsStaticCopy, true)
	return &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: res, UserDept: userDept, RelateData: relateData}
}

func (PostGreeter) InnerIssueUpdate(reqVo *projectvo.InnerIssueUpdateReq) *vo.VoidErr {
	_, projectId, err := domain.GetProjectIdByAppId(cast.ToString(reqVo.Input.AppId))
	if err != nil {
		return &vo.VoidErr{Err: vo.NewErr(err)}
	}
	req := &projectvo.BatchUpdateIssueReqInnerVo{
		OrgId:     reqVo.OrgId,
		UserId:    reqVo.UserId,
		AppId:     reqVo.Input.AppId,
		ProjectId: projectId,
		TableId:   reqVo.Input.TableId,
		Data:      reqVo.Input.Data,
	}
	err = service.SyncBatchUpdateIssue(req, true, reqVo.Input.TriggerBy)
	return &vo.VoidErr{Err: vo.NewErr(err)}
}
