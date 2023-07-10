package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/service"
)

func (PostGreeter)UploadOssByFsImageKey(req commonvo.UploadOssByFsImageKeyReq) commonvo.UploadOssByFsImageKeyResp {
	res, err := service.UploadOssByFsImageKey(req.OrgId, req.ImageKey, req.IsApp)
	return commonvo.UploadOssByFsImageKeyResp{
		Err: vo.NewErr(err),
		Url: res,
	}
}
