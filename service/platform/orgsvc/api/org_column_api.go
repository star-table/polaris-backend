package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) CreateOrgColumn(req orgvo.CreateOrgColumnReq) orgvo.CreateOrgColumnRespVo {
	resp, err := service.CreateOrgColumn(req)
	return orgvo.CreateOrgColumnRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (PostGreeter) GetOrgColumns(req orgvo.GetOrgColumnsReq) orgvo.GetOrgColumnsRespVo {
	resp, err := service.GetOrgColumns(req)
	return orgvo.GetOrgColumnsRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

func (PostGreeter) DeleteOrgColumn(req orgvo.DeleteOrgColumnReq) orgvo.DeleteOrgColumnRespVo {
	resp, err := service.DeleteOrgColumn(req)
	return orgvo.DeleteOrgColumnRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}
