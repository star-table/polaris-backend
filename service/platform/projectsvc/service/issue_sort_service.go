package service

//func UpdateIssueSort(reqVo projectvo.UpdateIssueSortReqVo) (*vo.Void, errs.SystemErrorInfo) {
//	currentUserId := reqVo.UserId
//	orgId := reqVo.OrgId
//	reqInput := reqVo.Input
//	issueId := reqInput.ID
//
//	issueBo, err := domain.GetIssueBo(orgId, issueId)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
//	}
//
//	//issueStatus := issueBo.Status
//
//	err = domain.AuthIssue(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
//	}
//
//	//if reqInput.ProjectObjectTypeID != nil && *reqInput.ProjectObjectTypeID != issueBo.TableId {
//	//	//更新项目对象类型(相比sort优先)
//	//	err1 := domain.MoveIssue(orgId, currentUserId, reqInput.FromProjectID, reqInput.FromProjectObjectTypeID, *issueBo, *reqInput.ProjectObjectTypeID, nil, nil, false, true, issueBo.TableId)
//	//
//	//	if err1 != nil {
//	//		log.Error(err1)
//	//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
//	//	}
//	//}
//
//	//if reqInput.StatusID != nil && *reqInput.StatusID != issueStatus {
//	//	//更新任务状态（相比sort优先，对于敏捷任务）
//	//	//needModifyChildStatus := int(1)
//	//	_, updateStatusErr := UpdateIssueStatus(projectvo.UpdateIssueStatusReqVo{
//	//		Input: vo.UpdateIssueStatusReq{
//	//			ID:           issueId,
//	//			NextStatusID: reqInput.StatusID,
//	//			//NeedModifyChildStatus:&needModifyChildStatus,
//	//		},
//	//		UserId:        currentUserId,
//	//		OrgId:         orgId,
//	//		SourceChannel: reqVo.SourceChannel,
//	//	})
//	//	if updateStatusErr != nil {
//	//		log.Error(updateStatusErr)
//	//		return nil, updateStatusErr
//	//	}
//	//	issueStatus = *reqInput.StatusID
//	//}
//
//	//排序走无码
//	if reqInput.BeforeDataID != nil || reqInput.AfterDataID != nil {
//		orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
//			OrgId:  orgId,
//			UserId: 0,
//		})
//		if orgInfoResp.Failure() {
//			log.Error(orgInfoResp.Error())
//			return nil, orgInfoResp.Error()
//		}
//		orgRemarkObj := &orgvo.OrgRemarkConfigType{}
//		oriErr := json.FromJson(orgInfoResp.OrganizationInfo.Remark, orgRemarkObj)
//		if oriErr != nil {
//			log.Error(oriErr)
//			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
//		}
//		lessReqParam := formvo.LessIssueListReq{
//			Condition: vo.LessCondsData{
//				Type:   "in",
//				Values: []interface{}{issueId},
//				Column: "issueId",
//				Left:   nil,
//				Right:  nil,
//				Conds:  nil,
//			},
//			AppId:   orgRemarkObj.OrgSummaryTableAppId,
//			OrgId:   orgId,
//			UserId:  currentUserId,
//			Page:    0,
//			Size:    0,
//			TableId: issueBo.TableId,
//		}
//
//		lessResp := formfacade.LessIssueList(lessReqParam)
//		if lessResp.Failure() {
//			log.Error(lessResp.Error())
//			return nil, lessResp.Error()
//		}
//
//		if len(lessResp.Data.List) > 0 {
//			appId, appIdErr := domain.GetAppIdFromProjectId(orgId, issueBo.ProjectId)
//			if appIdErr != nil {
//				log.Error(appIdErr)
//				return nil, appIdErr
//			}
//			dataId, ok := util.InterfaceToInt64(lessResp.Data.List[0]["id"])
//			if !ok {
//				//没有取到值，放行
//				return &vo.Void{ID: issueId}, nil
//			}
//
//			params := formvo.LessMoveIssueReq{
//				AppId:  appId,
//				DataId: dataId,
//				OrgId:  orgId,
//				UserId: currentUserId,
//			}
//			//处理数据
//			if reqInput.BeforeDataID != nil {
//				id, err := strconv.ParseInt(*reqInput.BeforeDataID, 10, 64)
//				if err != nil {
//					log.Error(err)
//					return nil, errs.ReqParamsValidateError
//				}
//				params.BeforeId = &id
//			}
//			if reqInput.AfterDataID != nil {
//				id, err := strconv.ParseInt(*reqInput.AfterDataID, 10, 64)
//				if err != nil {
//					log.Error(err)
//					return nil, errs.ReqParamsValidateError
//				}
//				params.AfterId = &id
//			}
//
//			if reqInput.Asc != nil {
//				params.Asc = *reqInput.Asc
//			}
//			moveResp := formfacade.LessMoveIssue(params)
//			if moveResp.Failure() {
//				log.Error(moveResp.Error())
//				return nil, moveResp.Error()
//			}
//
//			//推送可以后面再搞
//			asyn.Execute(func() {
//				PushModifyIssueNotice(issueBo.OrgId, issueBo.ProjectId, issueBo.Id, currentUserId)
//			})
//		}
//	}
//	return &vo.Void{ID: issueId}, nil
//}
