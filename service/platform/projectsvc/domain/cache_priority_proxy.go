package domain

//func GetPriorityList(orgId int64) (*[]bo.PriorityBo, errs.SystemErrorInfo) {
//	key, err5 := util.ParseCacheKey(sconsts.CachePriorityList, map[string]interface{}{
//		consts.CacheKeyOrgIdConstName: orgId,
//	})
//	if err5 != nil {
//		log.Error(err5)
//		return nil, err5
//	}
//	priorityListJson, err := cache.Get(key)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
//	}
//	priorityListPo := &[]po2.PpmPrsPriority{}
//	priorityListBo := &[]bo.PriorityBo{}
//	if priorityListJson != "" {
//		err = json.FromJson(priorityListJson, priorityListBo)
//		if err != nil {
//			log.Error(err)
//			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
//		}
//		return priorityListBo, nil
//	} else {
//		err := mysql.SelectAllByCondWithNumAndOrder(consts.TablePriority, db.Cond{
//			consts.TcOrgId:    orgId,
//			consts.TcIsDelete: consts.AppIsNoDelete,
//		}, nil, 0, 0, "sort asc", priorityListPo)
//		if err != nil {
//			log.Error(err)
//			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//		}
//		_ = copyer.Copy(priorityListPo, priorityListBo)
//		priorityListJson, err := json.ToJson(priorityListBo)
//		if err != nil {
//			log.Error(err)
//			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
//		}
//		err = cache.SetEx(key, priorityListJson, consts.GetCacheBaseExpire())
//		if err != nil {
//			log.Error(err)
//			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
//		}
//		return priorityListBo, nil
//	}
//}
//
//func GetPriorityListByType(orgId int64, typ int) (*[]bo.PriorityBo, errs.SystemErrorInfo) {
//	list, err := GetPriorityList(orgId)
//	if err != nil {
//		return nil, err
//	}
//	priorityList := &[]bo.PriorityBo{}
//	for _, priority := range *list {
//		if priority.Type == typ {
//			*priorityList = append(*priorityList, priority)
//		}
//	}
//	// 多语言适配
//	lang := lang2.GetLang()
//	isOtherLang := lang2.IsEnglish()
//	if isOtherLang {
//		otherLanguageMap := make(map[string]string, 0)
//		if tmpMap, ok1 := consts.LANG_PRIORITIES_MAP[lang]; ok1 {
//			otherLanguageMap = tmpMap
//		}
//		for index, item := range *priorityList {
//			if tmpVal, ok2 := otherLanguageMap[item.LangCode]; ok2 {
//				(*priorityList)[index].Name = tmpVal
//			}
//		}
//	}
//
//	return priorityList, nil
//}
//
//func GetPriorityById(orgId int64, id int64) (*bo.PriorityBo, errs.SystemErrorInfo) {
//	list, err := GetPriorityList(orgId)
//	if err != nil {
//		return nil, err
//	}
//	// 如果 id 为 0，则默认返回”普通优先级“
//	needDefaultPriority := false
//	if id == 0 {
//		needDefaultPriority = true
//	}
//	for _, priority := range *list {
//		if needDefaultPriority && priority.LangCode == consts.PriorityIssueCommon {
//			return &priority, nil
//		}
//		if priority.Id == id {
//			return &priority, nil
//		}
//	}
//
//	return nil, errs.BuildSystemErrorInfo(errs.PriorityNotExist)
//}
//
//func GetDefaultPriorityId(orgId int64, typ int) (int64, errs.SystemErrorInfo) {
//	list, err := GetPriorityList(orgId)
//	if err != nil {
//		return 0, err
//	}
//	for _, priorityBo := range *list {
//		if priorityBo.Type == typ && priorityBo.IsDefault == consts.APPIsDefault {
//			return priorityBo.Id, nil
//		}
//	}
//	return 0, errs.BuildSystemErrorInfo(errs.PriorityNotExist)
//}
