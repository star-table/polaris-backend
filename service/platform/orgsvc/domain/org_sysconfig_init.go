package domain

import (
	"github.com/star-table/polaris-backend/common/core/consts"
)

const OrcConfigSql = consts.TemplateDirPrefix + "ppm_orc_config.template"

//func OrgSysConfigInit(tx sqlbuilder.Tx, orgId int64) errs.SystemErrorInfo {
//	sysConfig := &po.PpmOrcConfig{}
//
//	payLevel := &po.PpmBasPayLevel{}
//	err := mysql.SelectById(payLevel.TableName(), 1, payLevel)
//	if err != nil {
//		log.Error(err)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//
//	respVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: sysConfig.TableName()})
//	if respVo.Failure() {
//		log.Error(respVo.Message)
//		return respVo.Error()
//	}
//
//	contextMap := map[string]interface{}{}
//	contextMap["Id"] = respVo.Id
//	contextMap["OrgId"] = orgId
//	contextMap["TimeZone"] = "Asia/Shanghai"
//	contextMap["TimeDifference"] = "+08:00"
//	contextMap["PayLevel"] = 1
//	contextMap["PayStartTime"] = time.Now().Format(consts.AppTimeFormat)
//	contextMap["PayEndTime"] = time.Now().Add(time.Duration(payLevel.Duration) * time.Second).Format(consts.AppTimeFormat)
//	contextMap["Language"] = "zh-CN"
//	contextMap["RemindSendTime"] = "09:00:00"
//	contextMap["DatetimeFormat"] = "yyyy-MM-dd HH:mm:ss"
//	contextMap["PasswordLength"] = 6
//	contextMap["PasswordRule"] = 1
//	contextMap["MaxLoginFailCount"] = 0
//	contextMap["Status"] = 1
//	insertErr := util.ReadAndWrite(OrcConfigSql, contextMap, tx)
//	if insertErr != nil {
//		return errs.BuildSystemErrorInfo(errs.BaseDomainError, insertErr)
//	}
//
//	//sysConfig.Id = respVo.Id
//	//sysConfig.OrgId = orgId
//	//sysConfig.TimeZone = "Asia/Shanghai"
//	//sysConfig.TimeDifference = "+08:00"
//	//sysConfig.PayLevel = 1
//	//sysConfig.PayStartTime = time.Now()
//	//sysConfig.PayEndTime = time.Now().Add(time.Duration(payLevel.Duration) * time.Second)
//	//sysConfig.Language = "zh-CN"
//	//sysConfig.RemindSendTime = "09:00:00"
//	//sysConfig.DatetimeFormat = "yyyy-MM-dd HH:mm:ss"
//	//sysConfig.PasswordLength = 6
//	//sysConfig.PasswordRule = 1
//	//sysConfig.MaxLoginFailCount = 0
//	//sysConfig.Status = 1
//	//sysConfig.IsDelete = consts.AppIsNoDelete
//	//
//	//err = mysql.TransInsert(tx, sysConfig)
//	//if err != nil {
//	//	log.Error(err)
//	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	//}
//	return nil
//}
