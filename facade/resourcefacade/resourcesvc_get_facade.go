package resourcefacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade"
)

func GetIdByPath(req resourcevo.GetIdByPathReqVo) resourcevo.GetIdByPathRespVo {
	respVo := &resourcevo.GetIdByPathRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getIdByPath", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["resourcePath"] = req.ResourcePath
	queryParams["resourceType"] = req.ResourceType
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}
