package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

func GetDingCoolApp() DingCoolApInterface {
	return dingCoolAppDao
}

type DingCoolApInterface interface {
	Get(conversationId string) (*po.PpmOrgDingCoolApp, errs.SystemErrorInfo)

	GetByProjectId(orgId, projectId int64) ([]*po.PpmOrgDingCoolApp, errs.SystemErrorInfo)

	CreateOrUpdate(m *po.PpmOrgDingCoolApp) errs.SystemErrorInfo

	Delete(conversationId string) errs.SystemErrorInfo

	DeleteByProjectId(orgId, projectId int64) errs.SystemErrorInfo
}

type dingCoolApp struct{}

var dingCoolAppDao = &dingCoolApp{}

func (g dingCoolApp) Get(conversationId string) (*po.PpmOrgDingCoolApp, errs.SystemErrorInfo) {
	temp := &po.PpmOrgDingCoolApp{}
	err := mysql.SelectOneByCond(po.TableNamePpmOrgDingCoolApp, db.Cond{consts.TcConversationId: conversationId}, temp)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return temp, nil
}

func (g dingCoolApp) GetByProjectId(orgId, projectId int64) ([]*po.PpmOrgDingCoolApp, errs.SystemErrorInfo) {
	temp := make([]*po.PpmOrgDingCoolApp, 0, 1)
	err := mysql.SelectAllByCond(po.TableNamePpmOrgDingCoolApp, db.Cond{consts.TcOrgId: orgId, consts.TcProjectId: projectId}, &temp)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return temp, nil
}

func (g dingCoolApp) CreateOrUpdate(m *po.PpmOrgDingCoolApp) errs.SystemErrorInfo {
	temp := &po.PpmOrgDingCoolApp{}
	err := mysql.SelectOneByCond(po.TableNamePpmOrgDingCoolApp, db.Cond{consts.TcConversationId: m.ConversationId}, temp)
	if err != nil {
		if err == db.ErrNoMoreRows {
			err = mysql.Insert(m)
			if err != nil {
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}
	} else {
		upd := mysql.Upd{
			consts.TcOrgId:     m.OrgId,
			consts.TcProjectId: m.ProjectId,
			consts.TcAppId:     m.AppId,
		}
		if m.TopTraceId != "" {
			upd[consts.TcTopTraceId] = m.TopTraceId
		}
		_, err = mysql.UpdateSmartWithCond(po.TableNamePpmOrgDingCoolApp, db.Cond{
			consts.TcConversationId: m.ConversationId,
		}, upd)
		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}

	return nil
}

func (g dingCoolApp) Delete(conversationId string) errs.SystemErrorInfo {
	conn, err := mysql.GetConnect()
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	_, err = conn.DeleteFrom(po.TableNamePpmOrgDingCoolApp).Where(consts.TcConversationId+" = ?", conversationId).Exec()
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func (g dingCoolApp) DeleteByProjectId(orgId, projectId int64) errs.SystemErrorInfo {
	conn, err := mysql.GetConnect()
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	_, err = conn.DeleteFrom(po.TableNamePpmOrgDingCoolApp).Where(consts.TcOrgId+" = ? and "+consts.TcProjectId+" = ? ", orgId, projectId).Exec()
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}
