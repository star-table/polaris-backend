package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) SaveMenu(req projectvo.SaveMenuReqVo) projectvo.SaveMenuRespVo {
	res, err := service.SaveMenu(req.OrgId, req.Input)
	return projectvo.SaveMenuRespVo{
		Err:   vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetMenu(req projectvo.GetMenuReqVo) projectvo.GetMenuRespVo {
	res, err := service.GetMenu(req.OrgId, req.AppId)
	return projectvo.GetMenuRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
