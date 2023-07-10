package domain

//func RelateIteration(orgId, projectId, iterationId, operatorId int64, addIssueIds []int64, delIssueIds []int64) errs.SystemErrorInfo {
//	appId, appIdErr := GetAppIdFromProjectId(orgId, projectId)
//	if appIdErr != nil {
//		log.Error(appIdErr)
//		return appIdErr
//	}
//
//	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
//		updForm := []map[string]interface{}{}
//
//		if len(addIssueIds) > 0 {
//			allIssueInfo, err := GetIssueAndChildren(orgId, addIssueIds)
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//			allNeedIds := []int64{}
//			for _, bo := range allIssueInfo {
//				allNeedIds = append(allNeedIds, bo.Id)
//			}
//
//			where := db.And(db.Cond{
//				consts.TcOrgId:     orgId,
//				consts.TcProjectId: projectId,
//				//consts.TcProjectObjectTypeId: db.NotEq(projectObjectType.Id),
//				consts.TcIsDelete: consts.AppIsNoDelete,
//				consts.TcId:       allNeedIds,
//			})
//			_, err1 := tx.Update(consts.TableIssue).Set(mysql.Upd{
//				consts.TcIterationId: iterationId,
//				consts.TcUpdator:     operatorId,
//			}).Where(where).Exec()
//			if err1 != nil {
//				log.Error(err1)
//				return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//			}
//
//			for _, id := range addIssueIds {
//				updForm = append(updForm, map[string]interface{}{
//					"issueId":     id,
//					"iterationId": iterationId,
//					"updator":     operatorId,
//				})
//			}
//		}
//
//		if len(delIssueIds) > 0 {
//			allIssueInfo, err := GetIssueAndChildren(orgId, delIssueIds)
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//			allNeedIds := []int64{}
//			for _, bo := range allIssueInfo {
//				allNeedIds = append(allNeedIds, bo.Id)
//			}
//			where := db.And(db.Cond{
//				consts.TcOrgId:       orgId,
//				consts.TcProjectId:   projectId,
//				consts.TcIterationId: db.NotEq(0),
//				//consts.TcProjectObjectTypeId: db.NotEq(projectObjectType.Id),
//				consts.TcIsDelete: consts.AppIsNoDelete,
//				consts.TcId:       allNeedIds,
//			})
//			_, err1 := tx.Update(consts.TableIssue).Set(mysql.Upd{
//				consts.TcIterationId: 0,
//				consts.TcUpdator:     operatorId,
//			}).Where(where).Exec()
//			if err1 != nil {
//				log.Error(err1)
//				return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//			}
//
//			for _, id := range delIssueIds {
//				updForm = append(updForm, map[string]interface{}{
//					"issueId":     id,
//					"iterationId": 0,
//					"updator":     operatorId,
//				})
//			}
//		}
//
//		//更新无码
//		resp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
//			AppId:  appId,
//			OrgId:  orgId,
//			UserId: operatorId,
//			//TableId:todo
//			Form: updForm,
//		})
//		if resp.Failure() {
//			log.Error(resp.Error())
//			return resp.Error()
//		}
//
//		return nil
//	})
//	if err1 != nil {
//		log.Error(err1)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
//	}
//	return nil
//}
//
//func JudgeIterationIsExist(orgId, id int64) bool {
//	return dao.JudgeIterationIsExist(orgId, id)
//}
