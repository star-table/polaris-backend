package commonvo

import "github.com/star-table/polaris-backend/common/model/vo"

type UploadOssByFsImageKeyReq struct {
	OrgId int64 `json:"orgId"`
	ImageKey string `json:"imageKey"`
	IsApp bool `json:"isApp"`
}

type UploadOssByFsImageKeyResp struct {
	vo.Err
	Url string `json:"url"`
}
