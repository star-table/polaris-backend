package service

import (
	"time"

	"github.com/spf13/cast"

	msgPb "github.com/star-table/interface/golang/msg/v1"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade/common/report"

	"github.com/star-table/polaris-backend/common/core/util/asyn"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"

	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"upper.io/db.v3"
)

func GetFunctionKeysByOrg(orgId int64) (*vo.FunctionConfigResp, errs.SystemErrorInfo) {
	bos, err := domain.GetOrgPayFunction(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	funcKeys := domain.GetFunctionKeyListByFunctions(bos)
	result := &vo.FunctionConfigResp{FunctionCodes: funcKeys}

	return result, nil
}

func GetFunctionsByOrg(orgId int64) ([]orgvo.FunctionLimitObj, errs.SystemErrorInfo) {
	bos, err := domain.GetOrgPayFunction(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	newFuncArr := make([]orgvo.FunctionLimitObj, 0, len(bos))
	copyer.Copy(bos, &newFuncArr)

	return newFuncArr, nil
}

func UpdateOrgFunctionConfig(orgId int64, sourceChannel string, level int64, buyType string, pricePlanType string, payTime time.Time, seats int, expireDays int, endDate time.Time, trailDays int) errs.SystemErrorInfo {
	log.Infof("[UpdateOrgFunctionConfig] sourceChannel:%v, orgId:%v, buyType:%v, pricePlanType:%v, payTime:%v, endDate:%v, payLeve:%v",
		sourceChannel, orgId, buyType, pricePlanType, payTime, endDate, level)
	//查看当前等级是否存在

	errSet := domain.ResetOrgPayNum(orgId)
	if errSet != nil {
		log.Errorf("[ResetOrgPayNum] orgId:%v, err:%v", orgId, errSet)
	}

	levelIsExist, existErr := domain.PayLevelIsExist(level)
	if existErr != nil {
		log.Error(existErr)
		return existErr
	}
	if !levelIsExist {
		return errs.PayLevelNotExist
	}
	//获取当前组织等级
	orgConfig, err := domain.GetOrgConfig(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	upd := mysql.Upd{}
	if int64(orgConfig.PayLevel) != level {
		upd[consts.TcPayLevel] = level
	}
	if seats != 0 {
		upd[consts.TcSeats] = seats //目前只用于定向方案
	}
	//支付当天的零点
	if len(upd) > 0 {
		_, updateErr := mysql.UpdateSmartWithCond(consts.TableOrgConfig, db.Cond{
			consts.TcOrgId: orgId,
		}, upd)
		if updateErr != nil {
			log.Error(updateErr)
			return errs.MysqlOperateError
		}
	}

	clearErr := domain.ClearOrgConfig(orgId)
	if clearErr != nil {
		log.Error(clearErr)
		return clearErr
	}

	// 上报事件
	asyn.Execute(func() {
		orgConfig, err := domain.GetOrgConfigRich(orgId)
		if err != nil {
			return
		}
		orgEvent := &commonvo.OrgEvent{}
		orgEvent.OrgId = orgId
		orgEvent.New = orgConfig

		openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		openTraceIdStr := cast.ToString(openTraceId)

		report.ReportOrgEvent(msgPb.EventType_OrgConfigUpdated, openTraceIdStr, orgEvent, true)
	})

	return nil
}
