package dashboardfacade

import (
	"fmt"

	"github.com/star-table/polaris-backend/facade"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
)

var log = logger.GetDefaultLogger()

func GetDashboardList(orgId, userId int64, appIds []int64) *projectvo.DashboardsRespVo {
	respVo := &projectvo.DashboardsRespVo{}
	reqUrl := fmt.Sprintf("%s/dashboard/inner/api/v3/dashboards/list", config.GetPreUrl(consts.ServiceDashboard))
	fullUrl := reqUrl + fmt.Sprintf("?orgId=%d&userId=%d", orgId, userId)
	err := facade.Request(consts.HttpMethodPost, fullUrl, nil, nil, appIds, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
