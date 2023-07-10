package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) OpenOrgUserList(reqVo orgvo.OpenOrgUserListReqVo) orgvo.OpenOrgUserListRespVo {
	res, err := service.OrgUserList(reqVo.OrgId, reqVo.UserId, reqVo.Page, reqVo.Size, reqVo.Input)
	return orgvo.OpenOrgUserListRespVo{Err: vo.NewErr(err), Data: res}
}
