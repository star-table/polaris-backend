package orgvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type CacheUserInfoVo struct {
	vo.Err

	CacheInfo bo.CacheUserInfoBo `json:"data"`
}
