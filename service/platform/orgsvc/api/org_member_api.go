package api

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) UpdateOrgMemberStatus(reqVo orgvo.UpdateOrgMemberStatusReq) vo.CommonRespVo {
	status := reqVo.Input.Status
	if status != 1 && status != 2 {
		log.Errorf("修改组织%d成员状态，状态不在正常范围内 %d", reqVo.OrgId, status)
		return vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError)), Void: nil}
	}
	res, err := service.UpdateOrgMemberStatus(reqVo)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

func (PostGreeter) UpdateOrgMemberCheckStatus(reqVo orgvo.UpdateOrgMemberCheckStatusReq) vo.CommonRespVo {
	status := reqVo.Input.CheckStatus
	if status != 1 && status != 2 && status != 3 {
		log.Errorf("修改组织%d成员审核状态，状态不在正常范围内 %d", reqVo.OrgId, status)
		return vo.CommonRespVo{Err: vo.NewErr(errs.BuildSystemErrorInfo(errs.ReqParamsValidateError)), Void: nil}
	}
	res, err := service.UpdateOrgMemberCheckStatus(reqVo)
	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
}

//func (PostGreeter) RemoveOrgMember(reqVo orgvo.RemoveOrgMemberReq) vo.CommonRespVo {
//	res, err := service.RemoveOrgMember(reqVo)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}

func (PostGreeter) OrgUserList(reqVo orgvo.OrgUserListReq) orgvo.OrgUserListResp {
	res, err := service.OrgUserList(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	return orgvo.OrgUserListResp{Err: vo.NewErr(err), Data: res}
}

func (PostGreeter) GetOrgUserInfoListBySourceChannel(reqVo orgvo.GetOrgUserInfoListBySourceChannelReq) orgvo.GetOrgUserInfoListBySourceChannelResp {
	res, err := service.GetOrgUserInfoListBySourceChannel(reqVo)
	return orgvo.GetOrgUserInfoListBySourceChannelResp{Data: res, Err: vo.NewErr(err)}
}
