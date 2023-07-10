package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) CreateColumn(req projectvo.CreateColumnReqVo) projectvo.CreateColumnRespVo {
	res, err := service.CreateColumn(req)
	return projectvo.CreateColumnRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) CopyColumn(req projectvo.CopyColumnReqVo) projectvo.CopyColumnRespVo {
	res, err := service.CopyColumn(req)
	return projectvo.CopyColumnRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) UpdateColumn(req projectvo.UpdateColumnReqVo) projectvo.UpdateColumnRespVo {
	res, err := service.UpdateColumn(req)
	return projectvo.UpdateColumnRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) DeleteColumn(req projectvo.DeleteColumnReqVo) projectvo.DeleteColumnRespVo {
	res, err := service.DeleteColumn(req)
	return projectvo.DeleteColumnRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) UpdateColumnDescription(req projectvo.UpdateColumnDescriptionReqVo) projectvo.UpdateColumnDescriptionRespVo {
	res, err := service.UpdateColumnDesc(req)
	return projectvo.UpdateColumnDescriptionRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}
