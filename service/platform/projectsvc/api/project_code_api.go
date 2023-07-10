package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) ConvertCode(reqVo projectvo.ConvertCodeReqVo) projectvo.ConvertCodeRespVo {
	res, err := service.ConvertCode(reqVo.OrgId, reqVo.Input)
	return projectvo.ConvertCodeRespVo{Err: vo.NewErr(err), ConvertCode: res}
}
