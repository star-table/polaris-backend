package commonvo

import (
	"github.com/star-table/polaris-backend/common/model/vo"
)

type IndustryListRespVo struct {
	vo.Err
	IndustryList *vo.IndustryListResp `json:"data"`
}
