package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) GetShareViewInfo(req *projectvo.GetShareViewInfoReq) *projectvo.GetShareViewInfoResp {
	info, err := service.GetShareViewInfo(req)
	return &projectvo.GetShareViewInfoResp{
		Err:  vo.NewErr(err),
		Data: info,
	}
}

func (PostGreeter) GetShareViewInfoByKey(req *projectvo.GetShareViewInfoByKeyReq) *projectvo.GetShareViewInfoResp {
	info, err := service.GetShareViewInfoByKey(req)
	return &projectvo.GetShareViewInfoResp{
		Err:  vo.NewErr(err),
		Data: info,
	}
}

func (PostGreeter) CreateShareView(req *projectvo.CreateShareViewReq) *projectvo.GetShareViewInfoResp {
	info, err := service.CreateShareView(req)
	return &projectvo.GetShareViewInfoResp{
		Err:  vo.NewErr(err),
		Data: info,
	}
}

func (PostGreeter) DeleteShareView(req *projectvo.DeleteShareViewReq) *vo.CommonRespVo {
	err := service.DeleteShareView(req)
	return &vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) ResetShareKey(req *projectvo.ResetShareKeyReq) *projectvo.GetShareViewInfoResp {
	info, err := service.ResetShareKey(req)
	return &projectvo.GetShareViewInfoResp{
		Err:  vo.NewErr(err),
		Data: info,
	}
}

func (PostGreeter) UpdateSharePassword(req *projectvo.UpdateSharePasswordReq) *vo.CommonRespVo {
	err := service.UpdateSharePassword(req)
	return &vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) UpdateShareConfig(req *projectvo.UpdateShareConfigReq) *vo.CommonRespVo {
	err := service.UpdateShareConfig(req)
	return &vo.CommonRespVo{
		Err: vo.NewErr(err),
	}
}

func (PostGreeter) CheckShareViewPassword(req *projectvo.CheckShareViewPasswordReq) *projectvo.CheckShareViewPasswordResp {
	info, err := service.CheckShareViewPassword(req.Input)
	return &projectvo.CheckShareViewPasswordResp{
		Err:  vo.NewErr(err),
		Data: info,
	}
}
