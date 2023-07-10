package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/idvo"
)

/**
批量获取id用该方法
*/
func ApplyMultipleIdExtra(orgId int64, code string, preCode string, count int) (*bo.IdCodes, errs.SystemErrorInfo) {
	rc := count / IdMaxMultipleCount
	rcy := count % IdMaxMultipleCount

	idcodes := &bo.IdCodes{}
	idcodeInfos := make([]bo.IdCodeInfo, 0, count)

	for i := 0; i < rc; i++ {
		idcodes, err := ApplyMultipleId(orgId, code, preCode, IdMaxMultipleCount)
		if err != nil {
			return nil, err
		}
		idcodeInfos = append(idcodeInfos, idcodes.Ids...)
	}

	if rcy > 0 {
		idcodes, err := ApplyMultipleId(orgId, code, preCode, int64(rcy))
		if err != nil {
			return nil, err
		}
		idcodeInfos = append(idcodeInfos, idcodes.Ids...)
	}
	idcodes.Ids = idcodeInfos

	return idcodes, nil
}

/**
批量获取id用该方法
*/
func ApplyMultiplePrimaryId(code string, count int) (*bo.IdCodes, errs.SystemErrorInfo) {
	return ApplyMultipleIdExtra(0, code, "", count)
}

/**
根据批量编号批量获取id用该方法
*/
func ApplyMultiplePrimaryIdByCodes(req idvo.ApplyMultiplePrimaryIdByCodesReqVo) (*bo.CodesIds, errs.SystemErrorInfo) {
	codeIds := &bo.CodesIds{
		IdCodes: make([]bo.IdCodes, 0),
	}
	err := errs.BuildSystemErrorInfo(errs.OK)
	for _, v := range req.CodeInfos {
		ids, err1 := ApplyMultipleIdExtra(0, v.Code, "", v.Count)
		if err1 != nil {
			err = err1
		} else {
			codeIds.IdCodes = append(codeIds.IdCodes, *ids)
		}
	}
	return codeIds, err
}
