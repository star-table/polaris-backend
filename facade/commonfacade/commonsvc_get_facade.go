package commonfacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade"
)

func IndustryList() commonvo.IndustryListRespVo {
	respVo := &commonvo.IndustryListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/commonsvc/industryList", config.GetPreUrl("commonsvc"))
	err := facade.Request(consts.HttpMethodGet, reqUrl, nil, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}
