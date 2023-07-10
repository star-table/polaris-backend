package domain

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/dashboardfacade"
)

func GetAppDashboards(orgId, appId int64) ([]*projectvo.DashboardInfo, errs.SystemErrorInfo) {
	resp := dashboardfacade.GetDashboardList(orgId, 0, []int64{appId})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}
	return resp.Data, nil
}
