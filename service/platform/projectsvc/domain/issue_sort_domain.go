package domain

//func UpdateIssueSort(issueBo bo.IssueBo, operatorId int64, beforeId, afterId *int64, issueStatus int64) ([]int64, errs.SystemErrorInfo) {
//	orgId := issueBo.OrgId
//
//	isBefore := false
//	refId := int64(0)
//
//	if beforeId != nil {
//		refId = *beforeId
//		isBefore = true
//	} else if afterId != nil {
//		refId = *afterId
//	}
//	if refId == issueBo.Id {
//		log.Error("要排序的任务和目标任务的id不能一致")
//		return nil, errs.BuildSystemErrorInfo(errs.IssueSortReferenceInvalidError)
//	}
//
//	refIssueBo, err := GetIssueBo(orgId, refId)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
//	}
//
//	projectIsAgile := false
//	projectBo, err := GetProject(orgId, issueBo.ProjectId)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	if projectBo.ProjectTypeId == consts.ProjectTypeAgileId {
//		projectIsAgile = true
//	}
//	relateIds := []int64{}
//	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
//		targetSort := issueBo.Sort
//		if isBefore {
//			cond := db.Cond{
//				consts.TcOrgId:               orgId,
//				consts.TcProjectId:           refIssueBo.ProjectId,
//				consts.TcProjectObjectTypeId: refIssueBo.ProjectObjectTypeId,
//				//consts.TcSort: db.Gt(refIssueBo.Sort),
//			}
//			//如果是敏捷项目，更新同状态下的排序
//			if projectIsAgile {
//				cond[consts.TcStatus] = issueStatus
//			}
//			upd := mysql.Upd{
//				consts.TcUpdator: operatorId,
//			}
//			if issueBo.Sort < refIssueBo.Sort {
//				cond[consts.TcSort] = db.Gt(issueBo.Sort)
//				cond[consts.TcSort+" "] = db.Lte(refIssueBo.Sort)
//				upd[consts.TcSort] = db.Raw("sort - 1")
//				targetSort = refIssueBo.Sort
//
//			} else {
//				cond[consts.TcSort] = db.Lt(issueBo.Sort)
//				cond[consts.TcSort+" "] = db.Gt(refIssueBo.Sort)
//				upd[consts.TcSort] = db.Raw("sort + 1")
//				targetSort = refIssueBo.Sort + 1
//			}
//
//			relateIssue := &[]po.PpmPriIssue{}
//			selectErr := mysql.TransSelectAllByCond(tx, consts.TableIssue, cond, relateIssue)
//			if selectErr != nil {
//				log.Error(selectErr)
//				return selectErr
//			}
//
//			for _, issue := range *relateIssue {
//				relateIds = append(relateIds, issue.Id)
//			}
//
//			_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, cond, upd)
//			if err1 != nil {
//				log.Error(err1)
//				return err1
//			}
//		} else {
//			targetSort = refIssueBo.Sort
//			cond := db.Cond{
//				consts.TcOrgId:               orgId,
//				consts.TcProjectId:           refIssueBo.ProjectId,
//				consts.TcProjectObjectTypeId: refIssueBo.ProjectObjectTypeId,
//				//consts.TcSort: db.Gte(refIssueBo.Sort),
//				//consts.TcSort + " ": db.Lte(issueBo.Sort),
//			}
//			//如果是敏捷项目，更新同状态下的排序
//			if projectIsAgile {
//				cond[consts.TcStatus] = issueStatus
//			}
//			if issueBo.ProjectObjectTypeId != refIssueBo.ProjectObjectTypeId || (projectIsAgile && issueBo.Status != refIssueBo.Status) {
//				cond[consts.TcSort] = db.Gte(refIssueBo.Sort)
//			} else {
//				cond[consts.TcSort] = db.Gte(refIssueBo.Sort)
//				cond[consts.TcSort+" "] = db.Lte(issueBo.Sort)
//			}
//
//			relateIssue := &[]po.PpmPriIssue{}
//			selectErr := mysql.TransSelectAllByCond(tx, consts.TableIssue, cond, relateIssue)
//			if selectErr != nil {
//				log.Error(selectErr)
//				return selectErr
//			}
//
//			for _, issue := range *relateIssue {
//				relateIds = append(relateIds, issue.Id)
//			}
//
//			_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, cond, mysql.Upd{
//				consts.TcSort:    db.Raw("sort + 1"),
//				consts.TcUpdator: operatorId,
//			})
//
//			if err2 != nil {
//				log.Error(err2)
//				return err2
//			}
//		}
//
//		err3 := mysql.TransUpdateSmart(tx, consts.TableIssue, issueBo.Id, mysql.Upd{
//			consts.TcSort:    targetSort,
//			consts.TcUpdator: operatorId,
//		})
//		if err3 != nil {
//			log.Error(err3)
//			return err3
//		}
//
//		if ok, _ := slice.Contain(relateIds, issueBo.Id); !ok {
//			relateIds = append(relateIds, issueBo.Id)
//		}
//
//		return nil
//	})
//	if transErr != nil {
//		log.Error(transErr)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//
//	return relateIds, nil
//}
