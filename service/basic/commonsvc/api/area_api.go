package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/service"
)

func (PostGreeter) AreaLinkageList(req commonvo.AreaLinkageListReqVo) commonvo.AreaLinkageListRespVo {
	res, err := service.AreaLinkageList(req.Input)
	return commonvo.AreaLinkageListRespVo{Err: vo.NewErr(err), AreaLinkageListResp: res}
}

func (PostGreeter) AreaInfo(req commonvo.AreaInfoReqVo) commonvo.AreaInfoRespVo {

	res, err := service.OrgAreaInfo(req)

	return commonvo.AreaInfoRespVo{Err: vo.NewErr(err), AreaInfoResp: res}

}
