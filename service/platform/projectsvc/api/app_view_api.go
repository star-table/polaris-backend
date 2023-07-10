package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

func (PostGreeter) MirrorsStat(reqVo *projectvo.MirrorsStatReq) projectvo.MirrorsStatResp {
	res, err := service.MirrorStat(reqVo.OrgId, reqVo.UserId, reqVo.Input.AppIds)
	return projectvo.MirrorsStatResp{
		Err: vo.NewErr(err),
		Data: vo.MirrorsStatResp{
			DataStat: res,
		},
	}
}
