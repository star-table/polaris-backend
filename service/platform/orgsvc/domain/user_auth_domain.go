package domain

//func DingAuth(corpId, openId string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
//	//获取组织信息
//	orgInfo, err := GetOrgInfoByOutOrgId(corpId, sdk_const.SourceChannelDingTalk)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.OrgNotInitError)
//	}
//	//获取用户信息
//	baseUserInfo, err := GetBaseUserInfoByEmpId(orgInfo.OrgId, openId)
//	if err != nil {
//		//这里做用户初始化的兜底
//		lockKey := consts.InitUserLock + sdk_const.SourceChannelDingTalk + openId
//		suc, err := cache.TryGetDistributedLock(lockKey, openId)
//		log.Infof("准备获取分布式锁 %v", suc)
//		if err != nil {
//			log.Error(err)
//			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
//		}
//		if suc {
//			log.Infof("获取分布式锁成功 %v", suc)
//			defer func() {
//				if _, lockErr := cache.ReleaseDistributedLock(lockKey, openId); lockErr != nil {
//					log.Error(lockErr)
//				}
//			}()
//
//			//double check
//			baseUserInfo, err = GetBaseUserInfoByEmpId(orgInfo.OrgId, openId)
//			if err != nil {
//				err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
//					_, err := InitDingTalkUser(orgInfo.OrgId, corpId, openId, tx)
//					if err != nil {
//						log.Error(err)
//						return err
//					}
//					return nil
//				})
//				if err1 != nil {
//					log.Error(err1)
//					return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//				}
//			}
//		}
//		if baseUserInfo == nil {
//			baseUserInfo, err = GetBaseUserInfoByEmpId(orgInfo.OrgId, openId)
//			if err != nil {
//				log.Error(err)
//				return nil, errs.UserInitError
//			}
//		}
//	}
//	return baseUserInfo, nil
//}
