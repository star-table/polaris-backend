package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) GetBaseOrgInfo(reqVo orgvo.GetBaseOrgInfoReqVo) orgvo.GetBaseOrgInfoRespVo {
	res, err := service.GetBaseOrgInfo(reqVo.OrgId)
	return orgvo.GetBaseOrgInfoRespVo{Err: vo.NewErr(err), BaseOrgInfo: res}
}

func (GetGreeter) GetBaseOrgInfoByOutOrgId(reqVo orgvo.GetBaseOrgInfoByOutOrgIdReqVo) orgvo.GetBaseOrgInfoByOutOrgIdRespVo {
	res, err := service.GetBaseOrgInfoByOutOrgId(reqVo.OutOrgId)
	return orgvo.GetBaseOrgInfoByOutOrgIdRespVo{Err: vo.NewErr(err), BaseOrgInfo: res}
}

func (PostGreeter) GetOrgOutInfoByOutOrgIdBatch(reqVo orgvo.GetOrgOutInfoByOutOrgIdBatchReqVo) orgvo.GetOrgOutInfoByOutOrgIdBatchRespVo {
	orgBoList, err := service.GetOrgOutInfoByOutOrgIdBatch(reqVo.Input)
	return orgvo.GetOrgOutInfoByOutOrgIdBatchRespVo{Err: vo.NewErr(err), Data: orgBoList}
}

func (GetGreeter) GetBaseUserInfoByEmpId(reqVo orgvo.GetBaseUserInfoByEmpIdReqVo) orgvo.GetBaseUserInfoByEmpIdRespVo {
	res, err := service.GetBaseUserInfoByEmpId(reqVo.OrgId, reqVo.EmpId)
	return orgvo.GetBaseUserInfoByEmpIdRespVo{Err: vo.NewErr(err), BaseUserInfo: res}
}

func (PostGreeter) GetBaseUserInfoByEmpIdBatch(reqVo orgvo.GetBaseUserInfoByEmpIdBatchReqVo) orgvo.GetBaseUserInfoByEmpIdBatchRespVo {
	res, err := service.GetBaseUserInfoByEmpIdBatch(reqVo.OrgId, reqVo.Input)
	return orgvo.GetBaseUserInfoByEmpIdBatchRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) GetUserConfigInfo(reqVo orgvo.GetUserConfigInfoReqVo) orgvo.GetUserConfigInfoRespVo {
	res, err := service.GetUserConfigInfo(reqVo.OrgId, reqVo.UserId)
	return orgvo.GetUserConfigInfoRespVo{Err: vo.NewErr(err), UserConfigInfo: res}
}

func (PostGreeter) GetUserConfigInfoBatch(reqVo orgvo.GetUserConfigInfoBatchReqVo) orgvo.GetUserConfigInfoBatchRespVo {
	boList, err := service.GetUserConfigInfoBatch(reqVo.OrgId, &reqVo.Input)
	return orgvo.GetUserConfigInfoBatchRespVo{Err: vo.NewErr(err), Data: boList}
}

func (GetGreeter) GetBaseUserInfo(reqVo orgvo.GetBaseUserInfoReqVo) orgvo.GetBaseUserInfoRespVo {
	res, err := service.GetBaseUserInfo(reqVo.OrgId, reqVo.UserId)
	return orgvo.GetBaseUserInfoRespVo{Err: vo.NewErr(err), BaseUserInfo: res}
}

func (GetGreeter) GetDingTalkBaseUserInfo(reqVo orgvo.GetDingTalkBaseUserInfoReqVo) orgvo.GetBaseUserInfoRespVo {
	res, err := service.GetDingTalkBaseUserInfo(reqVo.OrgId, reqVo.UserId)
	return orgvo.GetBaseUserInfoRespVo{Err: vo.NewErr(err), BaseUserInfo: res}
}

func (PostGreeter) GetBaseUserInfoBatch(reqVo orgvo.GetBaseUserInfoBatchReqVo) orgvo.GetBaseUserInfoBatchRespVo {
	res, err := service.GetBaseUserInfoBatch(reqVo.OrgId, reqVo.UserIds)
	return orgvo.GetBaseUserInfoBatchRespVo{Err: vo.NewErr(err), BaseUserInfos: res}
}

func (PostGreeter) GetShareUrl(req orgvo.GetShareUrlReq) orgvo.GetShareUrlResp {
	res, err := service.GetShareUrl(req.Key)
	return orgvo.GetShareUrlResp{
		Err: vo.NewErr(err),
		Url: res,
	}
}

func (PostGreeter) ClearOrgUsersPayCache(reqVo orgvo.GetBaseOrgInfoReqVo) vo.CommonRespVo {
	err := service.ClearOrgUsersPayCache(reqVo.OrgId)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: 0},
	}
}

func (PostGreeter) SetShareUrl(req orgvo.SetShareUrlReq) vo.CommonRespVo {
	err := service.SetShareUrl(req.Data.Key, req.Data.Value)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: 0},
	}
}
