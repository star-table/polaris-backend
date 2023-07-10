package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

func OrcConfigPageList(page int, size int) (*[]*po.PpmOrcConfig, int64, error) {

	//organizationPo := &[]*po.ScheduleOrganizationListPo{}

	conn, err := mysql.GetConnect()
	//paginator := conn.Select(db.Raw("o.id as id ,oc.project_daily_report_send_time as project_daily_report_send_time ")).From("ppm_org_organization as o").
	//	Join("ppm_orc_config as oc").On("oc.org_id = o.id").Where(db.Cond{
	//	"oc.is_delete": consts.AppIsNoDelete,
	//	"oc.status":    consts.AppStatusEnable,
	//	"o.is_delete":  consts.AppIsNoDelete,
	//	"o.status":     consts.AppStatusEnable,
	//}).Paginate(uint(size)).Page(uint(page))
	//
	//err = paginator.All(organizationPo)
	//
	//count, err := paginator.TotalEntries()
	//if err != nil {
	//	return nil, int64(count), errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	orcConfigPo := &[]*po.PpmOrcConfig{}

	if err != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	mid := conn.Collection(consts.TableOrgConfig).Find(db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	})

	count, err := mid.TotalEntries()

	if err != nil {
		return nil, int64(count), errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	if size > 0 && page > 0 {
		err = mid.Paginate(uint(size)).Page(uint(page)).All(orcConfigPo)
	} else {
		err = mid.All(orcConfigPo)
	}

	return orcConfigPo, int64(count), nil

}
