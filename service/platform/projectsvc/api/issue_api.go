package api

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
	"github.com/spf13/cast"
)

func (PostGreeter) CreateIssue(reqVo *projectvo.CreateIssueReqVo) *projectvo.IssueRespVo {
	var beforeDataId, afterDataId int64
	if reqVo.CreateIssue.BeforeDataID != nil {
		beforeDataId = cast.ToInt64(*reqVo.CreateIssue.BeforeDataID)
	}
	if reqVo.CreateIssue.AfterDataID != nil {
		afterDataId = cast.ToInt64(*reqVo.CreateIssue.AfterDataID)
	}
	req := &projectvo.BatchCreateIssueReqVo{
		OrgId:        reqVo.OrgId,
		UserId:       reqVo.UserId,
		AppId:        reqVo.InputAppId,
		ProjectId:    reqVo.CreateIssue.ProjectID,
		TableId:      cast.ToInt64(reqVo.CreateIssue.TableID),
		BeforeDataId: beforeDataId,
		AfterDataId:  afterDataId,
		Data:         []map[string]interface{}{reqVo.CreateIssue.LessCreateIssueReq},
	}
	data, userDepts, relateData, err := service.BatchCreateIssue(req, false, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	}, "", 0)
	return &projectvo.IssueRespVo{Err: vo.NewErr(err), Data: data, UserDepts: userDepts, RelateData: relateData}
}

func (PostGreeter) DeleteIssue(reqVo projectvo.DeleteIssueReqVo) projectvo.IssueRespVo {
	res, err := service.DeleteIssue(reqVo)
	return projectvo.IssueRespVo{Err: vo.NewErr(err), Issue: res}
}

func (PostGreeter) MoveIssue(reqVo projectvo.MoveIssueReqVo) vo.CommonRespVo {
	res, err := service.MoveIssue(reqVo.OrgId, reqVo.UserId, reqVo.InputAppId, reqVo.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) MoveIssueBatch(req projectvo.MoveIssueBatchReqVo) projectvo.MoveIssueBatchRespVo {
	res, err := service.MoveIssueBatch(req.OrgId, req.UserId, req.Input)
	return projectvo.MoveIssueBatchRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetIssueInfo(reqVo projectvo.GetIssueInfoReqVo) projectvo.GetIssueInfoRespVo {
	res, err := service.GetIssueInfo(reqVo.OrgId, reqVo.UserId, reqVo.IssueId)
	return projectvo.GetIssueInfoRespVo{
		Err:       vo.NewErr(err),
		IssueInfo: res,
	}
}

func (PostGreeter) GetIssueInfoList(reqVo projectvo.IssueInfoListReqVo) projectvo.IssueInfoListRespVo {
	res, err := service.GetIssueInfoList(reqVo.OrgId, reqVo.UserId, reqVo.IssueIds)
	return projectvo.IssueInfoListRespVo{
		Err:        vo.NewErr(err),
		IssueInfos: res,
	}
}

func (PostGreeter) GetIssueInfoByDataIdsList(reqVo projectvo.IssueInfoListByDataIdsReqVo) projectvo.IssueInfoListByDataIdsRespVo {
	res, err := service.GetIssueInfoListByDataIds(reqVo.OrgId, reqVo.UserId, reqVo.DataIds)
	return projectvo.IssueInfoListByDataIdsRespVo{
		Err:        vo.NewErr(err),
		IssueInfos: res,
	}
}

func (PostGreeter) CopyIssue(req projectvo.CopyIssueReqVo) *projectvo.LcDataListRespVo {
	res, userDept, relateData, err := service.CopyIssue(req.OrgId, req.UserId, req.Input)
	return &projectvo.LcDataListRespVo{Err: vo.NewErr(err), Data: res, UserDept: userDept, RelateData: relateData}
}

func (PostGreeter) CopyIssueBatch(req projectvo.CopyIssueBatchReqVo) projectvo.CopyIssueBatchRespVo {
	res, _, _, err := service.CopyIssueBatch(req.OrgId, req.UserId, &req.Input, &projectvo.TriggerBy{TriggerBy: consts.TriggerByCopy})
	return projectvo.CopyIssueBatchRespVo{Data: int64(len(res)), Err: vo.NewErr(err)}
}

func (PostGreeter) DeleteIssueBatch(req projectvo.DeleteIssueBatchReqVo) projectvo.DeleteIssueBatchRespVo {
	res, err := service.DeleteIssueBatch(req)
	return projectvo.DeleteIssueBatchRespVo{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) ConvertIssueToParent(req projectvo.ConvertIssueToParentReq) vo.CommonRespVo {
	res, err := service.ConvertIssueToParent(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) ChangeParentIssue(req projectvo.ChangeParentIssueReq) vo.CommonRespVo {
	res, err := service.ChangeParentIssue(req.OrgId, req.UserId, req.Input)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) AuditIssue(req projectvo.AuditIssueReq) vo.CommonRespVo {
	message := ""
	if req.Params.Comment != nil {
		message = *req.Params.Comment
	}
	issueIds, err := service.BatchAuditIssue(&projectvo.BatchAuditIssueReqVo{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input: &vo.LessBatchAuditIssueReq{
			IssueIds:    []int64{req.Params.IssueID},
			AuditStatus: req.Params.Status,
			Message:     message,
			Attachments: req.Params.Attachments,
		},
	})
	var issueId int64
	if len(issueIds) > 0 {
		issueId = issueIds[0]
	}
	return vo.CommonRespVo{
		Err: vo.NewErr(err),
		Void: &vo.Void{
			ID: issueId,
		},
	}
}

func (PostGreeter) UrgeIssue(req projectvo.UrgeIssueReqVo) vo.BoolRespVo {
	message := ""
	if req.Input.UrgeText != nil {
		message = *req.Input.UrgeText
	}
	issueIds, err := service.BatchUrgeIssue(&projectvo.BatchUrgeIssueReqVo{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input: &vo.LessBatchUrgeIssueReq{
			IssueIds: []int64{req.Input.IssueID},
			Message:  message,
		},
	})
	return vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: len(issueIds) > 0,
	}
}

func (PostGreeter) UrgeAuditIssue(req projectvo.UrgeAuditIssueReq) vo.BoolRespVo {
	issueIds, err := service.BatchUrgeIssue(&projectvo.BatchUrgeIssueReqVo{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input: &vo.LessBatchUrgeIssueReq{
			IssueIds: []int64{req.Input.IssueID},
			Message:  req.Input.UrgeText,
		},
	})
	return vo.BoolRespVo{
		Err:    vo.NewErr(err),
		IsTrue: len(issueIds) > 0,
	}
}

func (PostGreeter) ViewAuditIssue(req projectvo.ViewAuditIssueReq) vo.CommonRespVo {
	err := service.ViewAuditIssue(req.OrgId, req.UserId, req.IssueId)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: req.IssueId},
	}
}

func (PostGreeter) WithdrawIssue(req projectvo.WithdrawIssueReq) vo.CommonRespVo {
	err := service.WithdrawIssue(req.OrgId, req.UserId, req.IssueId)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: req.IssueId},
	}
}

//func (PostGreeter) GetIssueIdsByOrgId(req projectvo.GetIssueIdsByOrgIdReqVo) projectvo.GetIssueIdsByOrgIdRespVo {
//	resp, err := service.GetIssueIdsByOrgId(req.OrgId, req.UserId, req.Input)
//	return projectvo.GetIssueIdsByOrgIdRespVo{
//		Err:  vo.NewErr(err),
//		Data: resp,
//	}
//}

// InsertIssueProRelation 插入任务归属于项目的关联关系
//func (PostGreeter) InsertIssueProRelation(req projectvo.InsertIssueProRelationReqVo) projectvo.BoolRespVo {
//	err := service.InsertIssueProRelation(req.OrgId, req.UserId, req.Input)
//	return projectvo.BoolRespVo{
//		Err: vo.NewErr(err),
//		Data: &vo.BoolResp{
//			IsTrue: err == nil,
//		},
//	}
//}

// GetTableStatusByTableId 查询任务状态 供其他服务调用
func (PostGreeter) GetTableStatusByTableId(req projectvo.GetTableStatusReq) projectvo.GetTableStatusResp {
	resp, err := service.GetTableStatus(req.OrgId, req.Input.TableId)
	return projectvo.GetTableStatusResp{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// IssueStartChat 发起任务的群聊讨论
func (PostGreeter) IssueStartChat(req projectvo.IssueStartChatReqVo) projectvo.IssueStartChatRespVo {
	return projectvo.IssueStartChatRespVo{
		Err:  vo.NewErr(nil),
		Data: projectvo.IssueStartChatRespData{},
	}
}

func (PostGreeter) GetIssueRowList(req projectvo.IssueRowListReq) projectvo.IssueRowListResp {
	res, err := service.GetIssueRowList(req.OrgId, req.UserId, req.Input)
	return projectvo.IssueRowListResp{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
