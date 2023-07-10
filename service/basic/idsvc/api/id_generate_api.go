package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/idvo"
	"github.com/star-table/polaris-backend/service/basic/idsvc/service"
)

func (PostGreeter) ApplyPrimaryId(req idvo.ApplyPrimaryIdReqVo) idvo.ApplyPrimaryIdRespVo {
	id, err := service.ApplyPrimaryId(req.Code)
	return idvo.ApplyPrimaryIdRespVo{Id: id, Err: vo.NewErr(err)}
}

func (PostGreeter) ApplyCode(req idvo.ApplyCodeReqVo) idvo.ApplyCodeRespVo {
	code, err := service.ApplyCode(req.OrgId, req.Code, req.PreCode)
	return idvo.ApplyCodeRespVo{Code: code, Err: vo.NewErr(err)}
}

func (PostGreeter) ApplyMultipleId(req idvo.ApplyMultipleIdReqVo) idvo.ApplyMultipleIdRespVo {
	idCodes, err := service.ApplyMultipleIdExtra(req.OrgId, req.Code, req.PreCode, req.Count)
	return idvo.ApplyMultipleIdRespVo{IdCodes: idCodes, Err: vo.NewErr(err)}
}

func (PostGreeter) ApplyMultiplePrimaryId(req idvo.ApplyMultiplePrimaryIdReqVo) idvo.ApplyMultipleIdRespVo {
	idCodes, err := service.ApplyMultiplePrimaryId(req.Code, req.Count)
	return idvo.ApplyMultipleIdRespVo{IdCodes: idCodes, Err: vo.NewErr(err)}
}

func (PostGreeter) ApplyMultiplePrimaryIdByCodes(req idvo.ApplyMultiplePrimaryIdByCodesReqVo) idvo.ApplyMultipleIdRespByCodesVo {
	idCodes, err := service.ApplyMultiplePrimaryIdByCodes(req)
	return idvo.ApplyMultipleIdRespByCodesVo{CodesIds: idCodes, Err: vo.NewErr(err)}
}
