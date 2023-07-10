package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) GetInviteCode(req orgvo.GetInviteCodeReqVo) orgvo.GetInviteCodeRespVo {
	res, err := service.GetInviteCode(req.CurrentUserId, req.OrgId, req.SourcePlatform)
	return orgvo.GetInviteCodeRespVo{Err: vo.NewErr(err), Data: res}
}

func (GetGreeter) GetInviteInfo(req orgvo.GetInviteInfoReqVo) orgvo.GetInviteInfoRespVo {
	res, err := service.GetInviteInfo(req.InviteCode)
	return orgvo.GetInviteInfoRespVo{Err: vo.NewErr(err), Data: res}
}

// 通过多个手机号邀请
func (PostGreeter) InviteUserByPhones(req orgvo.InviteUserByPhonesReqVo) orgvo.InviteUserByPhonesRespVo {
	res, err := service.InviteUserByPhones(req.OrgId, req.UserId, req.Input)
	return orgvo.InviteUserByPhonesRespVo{Err: vo.NewErr(err), Data: res}
}

// 下载导入模板
func (PostGreeter) ExportInviteTemplate(req orgvo.ExportInviteTemplateReqVo) orgvo.ExportInviteTemplateRespVo {
	res, err := service.ExportInviteTemplate(req.OrgId, req.UserId)
	return orgvo.ExportInviteTemplateRespVo{Err: vo.NewErr(err), Data: res}
}

// 批量导入成员
func (PostGreeter) ImportMembers(req orgvo.ImportMembersReqVo) orgvo.ImportMembersRespVo {
	res, err := service.ImportMembers(req.OrgId, req.UserId, req.Input)
	return orgvo.ImportMembersRespVo{Err: vo.NewErr(err), Data: res}
}

// AddUser 直接添加用户
func (PostGreeter) AddUser(req orgvo.AddUserReqVo) vo.CommonRespVo {
	err := service.AddUser(&req)
	return vo.CommonRespVo{Err: vo.NewErr(err)}
}
