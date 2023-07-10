package service

const LuckyTagFinishTime = "2021-01-31 23:59:59"

////飞书福袋活动
//func LuckyTag(orgId int64) errs.SystemErrorInfo {
//	if time.Now().Format(consts.AppTimeFormat) > LuckyTagFinishTime {
//		return nil
//	}
//
//	//飞书组织在白名单里面的直接设置为标准版
//	if !domain.GetFeishuLuckyTagOrg(orgId) {
//		return nil
//	}
//
//	orgConfig := &po.PpmOrcConfig{}
//	err := mysql.SelectOneByCond(consts.TableOrgConfig, db.Cond{
//		consts.TcIsDelete: consts.AppIsNoDelete,
//		consts.TcOrgId:    orgId,
//	}, orgConfig)
//	if err != nil {
//		if err == db.ErrNoMoreRows {
//			return errs.OrgConfigNotExist
//		} else {
//			log.Error(err)
//			return errs.MysqlOperateError
//		}
//	}
//
//	//如果本来就是标准版且没有过期那沿用原来的
//	if orgConfig.PayLevel == consts.PayLevelStandard && orgConfig.PayEndTime.Format(consts.AppTimeFormat) > LuckyTagFinishTime {
//		return nil
//	}
//
//	_, updateErr := mysql.UpdateSmartWithCond(consts.TableOrgConfig, db.Cond{
//		consts.TcOrgId: orgId,
//	}, mysql.Upd{
//		consts.TcPayLevel:     consts.PayLevelStandard,
//		consts.TcPayStartTime: time.Now(),
//		consts.TcPayEndTime:   LuckyTagFinishTime,
//	})
//	if updateErr != nil {
//		log.Error(updateErr)
//		return errs.MysqlOperateError
//	}
//
//	clearErr := domain.ClearOrgConfig(orgId)
//	if clearErr != nil {
//		log.Error(clearErr)
//		return clearErr
//	}
//	return nil
//}
