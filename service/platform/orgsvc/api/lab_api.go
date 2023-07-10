package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) SetLabConfig(reqVo orgvo.SetLabReqVo) orgvo.SetLabRespVo {
	resp, err := service.SetLabConfig(reqVo.OrgId, reqVo.UserId, reqVo.Input)
	return orgvo.SetLabRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (GetGreeter) GetLabConfig(reqVo orgvo.GetLabReqVo) orgvo.GetLabRespVo {
	resp, err := service.GetLabConfig(reqVo)
	return orgvo.GetLabRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}
