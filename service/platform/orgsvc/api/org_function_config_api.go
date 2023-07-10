package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

// GetFunctionConfig 获取当前组织支持的功能 key 列表
func (PostGreeter) GetFunctionConfig(req orgvo.GetOrgFunctionConfigReq) orgvo.GetOrgFunctionConfigResp {
	res, err := service.GetFunctionKeysByOrg(req.OrgId)
	return orgvo.GetOrgFunctionConfigResp{
		Data: res,
		Err:  vo.NewErr(err),
	}
}

// GetFunctionObjArrByOrg 获取当前组织支持的功能列表信息
func (PostGreeter) GetFunctionObjArrByOrg(req orgvo.GetOrgFunctionConfigReq) orgvo.GetFunctionArrByOrgResp {
	res, err := service.GetFunctionsByOrg(req.OrgId)
	return orgvo.GetFunctionArrByOrgResp{
		Data: orgvo.GetFunctionArrByOrgRespData{
			Functions: res,
		},
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) UpdateOrgFunctionConfig(req orgvo.UpdateOrgFunctionConfigReq) vo.CommonRespVo {
	err := service.UpdateOrgFunctionConfig(req.OrgId, req.SourceChannel, req.Input.Level, req.Input.BuyType, req.Input.PricePlanType, req.Input.PayTime, req.Input.Seats, req.Input.ExpireDays, req.Input.EndDate, req.Input.TrailDays)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: &vo.Void{ID: req.OrgId},
	}
}
