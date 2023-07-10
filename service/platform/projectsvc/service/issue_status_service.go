package service

//func UpdateIssueStatus(reqVo projectvo.UpdateIssueStatusReqVo) (*vo.Issue, errs.SystemErrorInfo) {
//	orgId := reqVo.OrgId
//	input := reqVo.Input
//	currentUserId := reqVo.UserId
//	issueBo, err := domain.GetIssueBo(orgId, input.ID)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	//[权限迭代20211026] consts.OperationProIssue4ModifyStatus 被删除，对应的鉴权由字段权限控制。
//	err = domain.AuthIssueWithAppId(orgId, currentUserId, *issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify, reqVo.InputAppId, consts.BasicFieldIssueStatus)
//	if err != nil {
//		log.Error(err)
//		if err.Code() == errs.NoOperationPermissionForIssue.Code() {
//			// 更新任务信息时，没有权限时，产品需要使用该提示文案
//			return nil, errs.NoOperationPermissionForIssueUpdate
//		}
//		return nil, err
//	}
//	return UpdateIssueStatusWithoutAuth(reqVo)
//}
//
//func UpdateIssueStatusWithoutAuth(reqVo projectvo.UpdateIssueStatusReqVo) (*vo.Issue, errs.SystemErrorInfo) {
//	orgId := reqVo.OrgId
//	currentUserId := reqVo.UserId
//	input := reqVo.Input
//	sourceChannel := reqVo.SourceChannel
//
//	issueBo, err := domain.GetIssueBo(orgId, input.ID)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	var err1 errs.SystemErrorInfo = nil
//	if input.NextStatusID != nil && *input.NextStatusID > 0 {
//		err1 = domain.UpdateIssueStatus(*issueBo, currentUserId, *input.NextStatusID, sourceChannel)
//	} else if input.NextStatusType != nil && *input.NextStatusType > 0 {
//		err1 = domain.UpdateIssueStatusByStatusType(*issueBo, currentUserId, *input.NextStatusType, sourceChannel)
//	} else {
//		log.Error("要更新的状态无效")
//		return nil, errs.BuildSystemErrorInfo(errs.IssueStatusUpdateError)
//	}
//
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
//	}
//
//	result := &vo.Issue{}
//	copyErr := copyer.Copy(issueBo, result)
//	if copyErr != nil {
//		log.Errorf("copyer.Copy(): %q\n", copyErr)
//	}
//
//	//asyn.Execute(func() {
//	//	PushModifyIssueNotice(issueBo.OrgId, issueBo.ProjectId, issueBo.Id, currentUserId)
//	//})
//	return result, nil
//}
//
//func GetIssueAllStatusNew(orgId int64, projectIds, tableIds []int64) (map[int64][]status.StatusInfoBo, errs.SystemErrorInfo) {
//	return domain.GetIssueAllStatus(orgId, projectIds, tableIds)
//}
