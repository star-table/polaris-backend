package officefacade

import (
	"fmt"

	"github.com/star-table/polaris-backend/facade"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/model/vo/lc_office"

	"github.com/star-table/polaris-backend/common/core/consts"
)

func GetOfficeConfig(orgId, userId int64) lc_office.GetOfficeConfigRespVo {
	respVo := &lc_office.GetOfficeConfigRespVo{}
	reqUrl := fmt.Sprintf("%s/api/config", config.GetPreUrl("officesvc"))
	err := facade.RequestWithCommonHeader(orgId, userId, consts.HttpMethodGet, reqUrl, nil, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}
