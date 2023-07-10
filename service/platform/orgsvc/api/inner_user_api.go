package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) InnerUserInfo(req *orgvo.InnerUserInfosReq) *vo.DataRespVo {
	userInfos, err := service.InnerGetUserInfos(req.OrgId, req.Input.Ids)
	return &vo.DataRespVo{Err: vo.NewErr(err), Data: userInfos}
}
