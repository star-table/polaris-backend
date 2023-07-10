package domain

import (
	"strconv"

	"github.com/star-table/common/core/util/date"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"upper.io/db.v3"
)

//var getProcessError = "proxies.GetProcessStatusId: %v\n"

//func IssueCondRelatedTypeAssembly(issueCond db.Cond, relatedType int, relatedUserId int64, orgId int64) {
//	switch relatedType {
//	case 1:
//		issueCond["creator"] = relatedUserId
//	case 2:
//		issueCond[consts.TcOwner] = relatedUserId
//	case 3, 4:
//		rt := consts.IssueRelationTypeParticipant
//		if relatedType == 4 {
//			rt = consts.IssueRelationTypeFollower
//		}
//		issueCond[consts.TcId] = db.In(db.Raw("select ir.issue_id as id from ppm_pri_issue_relation ir where ir.is_delete = 2 and ir.relation_id = ? and ir.relation_type = ?", relatedUserId, rt))
//	case 5:
//		issueCond[consts.TcId] = db.In(db.Raw("select ir.issue_id as id from ppm_pri_issue_relation ir where ir.is_delete = 2 and ir.relation_id = ? and ir.relation_type = ?", relatedUserId, consts.IssueRelationTypeAuditor))
//		issueCond[consts.TcProjectId+" "] = db.In(db.Raw("select id from ppm_pro_project where project_type_id = 1 and org_id in ?", []int64{0, orgId}))
//	case 6:
//		issueCond[consts.TcId] = db.In(db.Raw("select ir.issue_id as id from ppm_pri_issue_relation ir where ir.is_delete = 2 and ir.relation_id = ? and ir.relation_type = ? and ir.status in ?", relatedUserId, consts.IssueRelationTypeAuditor, []int{consts.AuditStatusNotView, consts.AuditStatusView}))
//		issueCond[consts.TcAuditStatus+" "] = consts.AuditStatusNotView
//		issueCond[consts.TcStatus+" "] = int64(26)
//		issueCond[consts.TcProjectId+" "] = db.In(db.Raw("select id from ppm_pro_project where project_type_id = 1 and org_id in ?", []int64{0, orgId}))
//	}
//}

//func IssueCondNoRelatedTypeAssembly(issueCond db.Cond, relatedUserId int64, orgId int64, isAdmin bool) {
//	args := []interface{}{orgId}
//	sql := "select p.id as id from ppm_pro_project p where (p.org_id = ? and p.is_delete = 2"
//	//增加项目限制
//	if !isAdmin {
//		deptInfo := orgfacade.GetUserDeptIds(orgvo.GetUserDeptIdsReq{
//			OrgId:  orgId,
//			UserId: relatedUserId,
//		})
//		if deptInfo.Failure() {
//			log.Error(deptInfo.Error())
//			return
//		}
//		deptInfo.NewData.DeptIds = append(deptInfo.NewData.DeptIds, 0)
//		//把任务项目id为0的也带出来
//		sql += " and (p.public_status = 1 or p.id in (SELECT DISTINCT pr.project_id FROM ppm_pro_project_relation pr WHERE ((pr.relation_id = ? AND relation_type in (1,2)) or (pr.relation_id in ? and pr.relation_type = 25)) AND pr.is_delete = 2))"
//		args = append(args, relatedUserId, deptInfo.NewData.DeptIds)
//	}
//	sql += ") or p.id = 0"
//
//	issueCond[consts.TcProjectId] = db.In(db.Raw(sql, args...))
//}

//func IssueCondTagId(issueCond db.Cond, orgId int64, tagIds []int64) {
//	//查询tag
//	issueCond[consts.TcId+" IN"] = db.In(db.Raw("select it.issue_id from ppm_pri_issue_tag it where it.org_id = ? and it.tag_id in ? and it.is_delete = 2", orgId, tagIds))
//}

//func IssueCondResourceId(issueCond db.Cond, orgId, resourceId int64) {
//	//查询tag
//	issueCond[consts.TcId+" IN "] = db.In(db.Raw("select ir.issue_id from ppm_pri_issue_relation ir where ir.org_id = ? and ir.relation_id = ? and ir.relation_type = ? and ir.is_delete = 2", orgId, resourceId, consts.IssueRelationTypeResource))
//}

//func IssueCondFiling(issueCond db.Cond, orgId int64, isFiling int) {
//	if isFiling == 3 {
//		return
//	} else {
//		if isFiling == 2 {
//			//未归档把0塞进去
//			issueCond[consts.TcProjectId+" IN"] = db.In(db.Raw("select id from ppm_pro_project where is_filing = ? and org_id in ? and is_delete = 2", isFiling, []int64{0, orgId}))
//		} else {
//			issueCond[consts.TcProjectId+" IN"] = db.In(db.Raw("select id from ppm_pro_project where is_filing = ? and org_id = ? and is_delete = 2", isFiling, orgId))
//		}
//	}
//}

//func IssueCondOrderBy(orgId int64, orderByType int, projectInfo *bo.Project, projectObjectTypeId *int64, issueIds []int64) interface{} {
//	var orderBy interface{} = nil
//	switch orderByType {
//	case 1:
//		//项目分组
//		orderBy = db.Raw("project_id desc, id desc")
//	case 2:
//		priorities, err := GetPriorityListByType(orgId, consts.PriorityTypeIssue)
//		if err != nil {
//			log.Error(err)
//			orderBy = db.Raw("(select sort from ppm_prs_priority p where p.id = priority_id) asc, plan_end_time asc, id desc")
//		} else {
//			bo.SortPriorityBo(*priorities)
//			orderBySort := ""
//			for _, priority := range *priorities {
//				orderBySort += fmt.Sprintf(",%d", priority.Id)
//			}
//			orderBy = db.Raw("FIELD(priority_id" + orderBySort + ")")
//		}
//	case 3:
//		//创建时间降序
//		orderBy = db.Raw("create_time desc, id desc")
//	case 4:
//		//更新时间降序
//		orderBy = db.Raw("update_time desc, id desc")
//	case 5:
//		//按开始时间最早
//		orderBy = db.Raw("if (plan_start_time > '1970-02-01 00:00:00',1,0) desc, plan_start_time asc, id desc")
//	case 6:
//		//按开始时间最晚
//		orderBy = db.Raw("plan_start_time desc, id desc")
//	case 8:
//		//按截止时间最近
//		orderBy = db.Raw("if (plan_end_time > '1970-02-01 00:00:00',1,0) desc, plan_end_time desc, id desc")
//	case 9:
//		//按创建时间最早
//		orderBy = db.Raw("create_time asc, id desc")
//	case 10:
//		//按照sort排序正序
//		orderBy = db.Raw("sort asc, id desc")
//	case 11:
//		//按照sort倒序
//		orderBy = db.Raw("sort desc, id desc")
//	case 12:
//		//按截止时间正序
//		orderBy = db.Raw("if (plan_end_time > '1970-02-01 00:00:00',1,0) desc, plan_end_time asc, id desc")
//	case 13:
//		//优先级正序
//		priorities, err := GetPriorityListByType(orgId, consts.PriorityTypeIssue)
//		if err != nil {
//			log.Error(err)
//			orderBy = db.Raw("(select sort from ppm_prs_priority p where p.id = priority_id) asc, plan_end_time asc, id desc")
//		} else {
//			bo.SortPriorityBo(*priorities)
//			orderBySort := ""
//			for _, priority := range *priorities {
//				orderBySort += fmt.Sprintf(",%d", priority.Id)
//			}
//			orderBy = db.Raw("FIELD(priority_id" + orderBySort + ")")
//		}
//		//orderBy = db.Raw("priority_id desc")
//	case 14:
//		//优先级倒序
//		priorities, err := GetPriorityListByType(orgId, consts.PriorityTypeIssue)
//		if err != nil {
//			log.Error(err)
//			orderBy = db.Raw("(select sort from ppm_prs_priority p where p.id = priority_id) desc, plan_end_time asc, id desc")
//		} else {
//			bo.SortDescPriorityBo(*priorities)
//			orderBySort := ""
//			for _, priority := range *priorities {
//				orderBySort += fmt.Sprintf(",%d", priority.Id)
//			}
//			orderBy = db.Raw("FIELD(priority_id" + orderBySort + ")")
//		}
//		//orderBy = db.Raw("priority_id asc")
//	case 15:
//		//负责人正序
//		orderBy = db.Raw("owner asc")
//	case 16:
//		//负责人倒序
//		orderBy = db.Raw("owner desc")
//	case 17:
//		//编号正序
//		orderBy = db.Raw("code asc")
//	case 18:
//		//编号倒序
//		orderBy = db.Raw("code desc")
//	case 19:
//		//标题正序
//		orderBy = db.Raw("title asc")
//	case 20:
//		//标题倒序
//		orderBy = db.Raw("title desc")
//	case 21:
//		//按照状态正序（只有传入了项目id才生效，敏捷要传入任务栏id）
//		if projectInfo != nil && projectInfo.Id > 0 {
//			if projectInfo.ProjectTypeId == consts.ProjectTypeNormalId {
//				orderBy = db.Raw("status asc")
//			} else {
//				if projectObjectTypeId != nil {
//					allStatus, err := GetIssueAllStatus(orgId, []int64{projectInfo.Id}, []int64{*projectObjectTypeId})
//					if err != nil {
//						log.Error(err)
//					}
//					if _, ok := allStatus[*projectObjectTypeId]; ok {
//						orderBySort := ""
//						for _, status := range allStatus[*projectObjectTypeId] {
//							orderBySort += fmt.Sprintf(",%d", status.ID)
//						}
//						orderBy = db.Raw("FIELD(status" + orderBySort + ")")
//					}
//				}
//			}
//		}
//	case 22:
//		//按照状态倒序（只有传入了项目id才生效，敏捷要传入任务栏id）
//		if projectInfo != nil && projectInfo.Id > 0 {
//			if projectInfo.ProjectTypeId == consts.ProjectTypeNormalId {
//				orderBy = db.Raw("status desc")
//			} else {
//				if projectObjectTypeId != nil {
//					allStatus, err := GetIssueAllStatus(orgId, []int64{projectInfo.Id}, []int64{*projectObjectTypeId})
//					if err != nil {
//						log.Error(err)
//					}
//					if _, ok := allStatus[*projectObjectTypeId]; ok {
//						length := len(allStatus[*projectObjectTypeId])
//						newIds := make([]int64, length)
//						for i, status := range allStatus[*projectObjectTypeId] {
//							newIds[length-1-i] = status.ID
//						}
//						orderBySort := ""
//						for _, id := range newIds {
//							orderBySort += fmt.Sprintf(",%d", id)
//						}
//						orderBy = db.Raw("FIELD(status" + orderBySort + ")")
//					}
//				}
//			}
//		}
//	case 23:
//		//按照完成时间倒序
//		orderBy = db.Raw("end_time desc, id desc")
//	case 24:
//		//按照传入id排序
//		if len(issueIds) > 0 {
//			orderBySort := ""
//			for _, id := range issueIds {
//				orderBySort += fmt.Sprintf(",%d", id)
//			}
//			orderBy = db.Raw("FIELD(id" + orderBySort + ")")
//		}
//	case 25:
//		//按照父任务正序
//		orderBy = db.Raw("parent_id asc")
//	case 26:
//		//按照父任务降序
//		orderBy = db.Raw("parent_id desc")
//	}
//	return orderBy
//}

//状态类型,-1:待确定 1:未完成，2：已完成，3：未开始，4：进行中 5: 已逾期
//func IssueCondStatusAssembly(issueCond db.Cond, orgId int64, status int) errs.SystemErrorInfo {
//	var statusIds []int64 = nil
//
//	if status == consts.IssueStatusUnfinished {
//		notStartedIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeNotStart)
//		if err != nil {
//			log.Errorf(getProcessError, err)
//			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//		}
//		processingIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeRunning)
//		if err != nil {
//			log.Errorf(getProcessError, err)
//			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//		}
//		statusIds = append(notStartedIds, processingIds...)
//		issueCond[consts.TcStatus] = db.In(statusIds)
//	} else if status == consts.IssueStatusFinished {
//		finishedId, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
//		if err != nil {
//			log.Errorf(getProcessError, err)
//			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//		}
//		statusIds = finishedId
//		issueCond[consts.TcStatus] = db.In(statusIds)
//		issueCond[consts.TcAuditStatus] = consts.AuditStatusPass
//	} else if status == consts.IssueStatusNotStarted {
//		notStartedIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeNotStart)
//		if err != nil {
//			log.Errorf(getProcessError, err)
//			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//		}
//		statusIds = notStartedIds
//		issueCond[consts.TcStatus] = db.In(statusIds)
//	} else if status == consts.IssueStatusProcessing {
//		processingIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeRunning)
//		if err != nil {
//			log.Errorf(getProcessError, err)
//			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//		}
//		statusIds = processingIds
//		issueCond[consts.TcStatus] = db.In(statusIds)
//	} else if status == consts.IssueStatusOverdue {
//		//已经逾期筛选条件，并且未完成
//		nowTime := time.Now()
//		//逾期
//		issueCond[consts.TcPlanEndTime] = db.Between(consts.BlankElasticityTime, date.Format(nowTime))
//		//未完成
//		finishedId, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
//		if err != nil {
//			log.Errorf(getProcessError, err)
//			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//		}
//		issueCond[consts.TcStatus] = db.NotIn(finishedId)
//	} else if status == consts.IssueStatusNone {
//		//待确认
//		finishedId, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
//		if err != nil {
//			log.Errorf(getProcessError, err)
//			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//		}
//		statusIds = finishedId
//		issueCond[consts.TcStatus] = db.In(statusIds)
//		issueCond[consts.TcAuditStatus] = consts.AuditStatusNotView
//		issueCond[consts.TcProjectId+"   "] = db.In(db.Raw("select id from ppm_pro_project where project_type_id = 1 and org_id in ?", []int64{0, orgId}))
//	}
//	return nil
//}

//// IssueCondStatusAssemblyForLc 状态类型,1:未完成，2：已完成，3：未开始，4：进行中 5: 已逾期。拼接查询无码的条件 todo
//// 任务状态改造版本，拼接查询无码的条件。改造后，查询状态需要通过所属表/项目。
//func IssueCondStatusAssemblyForLc(orgId, projectId int64, tableIds []int64, issueCond db.Cond, status int) errs.SystemErrorInfo {
//	var err errs.SystemErrorInfo
//	proAppId := int64(0)
//	// 一般是有传入的 projectId，如果没有则意味着任务位于“空项目”中
//	if projectId > 0 {
//		proAppId, err = GetAppIdFromProjectId(orgId, projectId)
//	} else {
//		// 查询组织的空项目对应的 proAppId
//		proAppId, err = GetEmptyProAppId(orgId)
//		if err != nil {
//			log.Errorf("[IssueCondStatusAssemblyForLc] orgId: %d, err: %v", orgId, err)
//			return err
//		}
//	}
//	log.Infof("proAppId: %d", proAppId)
//	// 查询项目下，所有的表对应的任务状态列
//	switch status {
//	case 1:
//	case 2:
//	case 3:
//	case 4:
//	case 5:
//
//	}
//	return nil
//}

//// GetEmptyProAppId 获取组织的空项目对应的 appId todo
//func GetEmptyProAppId(orgId int64) (int64, errs.SystemErrorInfo) {
//	orgResp := orgfacade.GetOrgBoListByPage(orgvo.GetOrgIdListByPageReqVo{
//		Page: 1,
//		Size: 1,
//		Input: orgvo.GetOrgIdListByPageReqVoData{
//			OrgIds: []int64{orgId},
//		},
//	})
//	if orgResp.Failure() {
//		log.Errorf("[GetEmptyProAppId] orgId: %d, err: %v", orgId, orgResp.Error())
//		return 0, orgResp.Error()
//	}
//	if len(orgResp.NewData.List) < 1 {
//		err := errs.OrgNotExist
//		log.Errorf("[GetEmptyProAppId] orgId: %d, err: %v", orgId, err)
//		return 0, err
//	}
//	org := orgResp.NewData.List[0]
//	orgRemarkObj := orgvo.OrgRemarkConfigType{}
//	if err := json.FromJson(org.Remark, &orgRemarkObj); err != nil {
//		log.Errorf("[GetEmptyProAppId] fromJson err: %v", err)
//		return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
//	}
//
//	return orgRemarkObj.EmptyProjectAppId, nil
//}

//状态类型,-1:待确定,1:未完成，2：已完成，3：未开始，4：进行中 5: 已逾期
//func IssueCondStatusListAssembly(orgId int64, statusList []int) (*db.Union, errs.SystemErrorInfo) {
//	var issueUnion db.Union
//	for _, status := range statusList {
//		var statusIds []int64 = nil
//		if status == consts.IssueStatusUnfinished {
//			notStartedIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeNotStart)
//			if err != nil {
//				log.Errorf(getProcessError, err)
//				return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//			}
//			processingIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeRunning)
//			if err != nil {
//				log.Errorf(getProcessError, err)
//				return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//			}
//
//			statusIds = append(append(statusIds, notStartedIds...), processingIds...)
//			issueUnion = *(issueUnion.Or(
//				db.Cond{
//					consts.TcStatus: db.In(statusIds),
//				}))
//		} else if status == consts.IssueStatusFinished {
//			finishedId, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
//			if err != nil {
//				log.Errorf(getProcessError, err)
//				return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//			}
//			statusIds = append(statusIds, finishedId...)
//			issueUnion = *(issueUnion.Or(db.And(
//				db.Cond{
//					consts.TcStatus:      db.In(statusIds),
//					consts.TcAuditStatus: consts.AuditStatusPass,
//				})))
//		} else if status == consts.IssueStatusNotStarted {
//			notStartedIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeNotStart)
//			if err != nil {
//				log.Errorf(getProcessError, err)
//				return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//			}
//			statusIds = append(statusIds, notStartedIds...)
//			issueUnion = *(issueUnion.Or(
//				db.Cond{
//					consts.TcStatus: db.In(statusIds),
//				}))
//		} else if status == consts.IssueStatusProcessing {
//			processingIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeRunning)
//			if err != nil {
//				log.Errorf(getProcessError, err)
//				return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//			}
//			statusIds = append(statusIds, processingIds...)
//			issueUnion = *(issueUnion.Or(
//				db.Cond{
//					consts.TcStatus: db.In(statusIds),
//				}))
//		} else if status == consts.IssueStatusOverdue {
//			notStartedIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeNotStart)
//			if err != nil {
//				log.Errorf(getProcessError, err)
//				return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//			}
//			processingIds, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeRunning)
//			if err != nil {
//				log.Errorf(getProcessError, err)
//				return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//			}
//			statusIds = append(append(statusIds, notStartedIds...), processingIds...)
//			nowTime := time.Now()
//			issueUnion = *(issueUnion.Or(db.And(
//				db.Cond{
//					consts.TcStatus:      db.In(statusIds),
//					consts.TcPlanEndTime: db.Between(consts.BlankElasticityTime, date.Format(nowTime)),
//				})))
//
//		} else if status == consts.IssueStatusNone {
//			finishedId, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
//			if err != nil {
//				log.Errorf(getProcessError, err)
//				return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//			}
//			statusIds = append(statusIds, finishedId...)
//			issueUnion = *(issueUnion.Or(db.And(db.Cond{
//				consts.TcStatus:            db.In(statusIds),
//				consts.TcAuditStatus:       consts.AuditStatusNotView,
//				consts.TcProjectId + "   ": db.In(db.Raw("select id from ppm_pro_project where project_type_id = 1 and org_id in ?", []int64{0, orgId})),
//			})))
//
//		}
//	}
//
//	return &issueUnion, nil
//}

//func IssueCondRelationMemberAssembly(queryCond db.Cond, input *vo.HomeIssueInfoReq) {
//	selectOthersIssueSql := "select ir.issue_id as id from ppm_pri_issue_relation ir where ir.is_delete = 2"
//	memberIdsList := &[]interface{}{}
//	selectOthersIssueSqlIsJoint := false
//	selectOthersIssueSqlJointTag := " and "
//
//	//可以在基本条件封装时做处理
//	//if input.OwnerIds != nil{
//	//	*memberIdsList = append(*memberIdsList, input.OwnerIds)
//	//	if selectOthersIssueSqlIsJoint{
//	//		selectOthersIssueSqlJointTag = " or "
//	//	}else{
//	//		selectOthersIssueSqlIsJoint = !selectOthersIssueSqlIsJoint
//	//	}
//	//	selectOthersIssueSql += selectOthersIssueSqlJointTag + "(ir.relation_id in ? and ir.relation_type = 1)"
//	//}
//	if input.ParticipantIds != nil && len(input.ParticipantIds) > 0 {
//		*memberIdsList = append(*memberIdsList, input.ParticipantIds)
//		if selectOthersIssueSqlIsJoint {
//			selectOthersIssueSqlJointTag = " or "
//		} else {
//			selectOthersIssueSqlIsJoint = !selectOthersIssueSqlIsJoint
//		}
//		selectOthersIssueSql += selectOthersIssueSqlJointTag + "(ir.relation_id in ? and ir.relation_type = 2)"
//	}
//	if input.FollowerIds != nil && len(input.FollowerIds) > 0 {
//		*memberIdsList = append(*memberIdsList, input.FollowerIds)
//		if selectOthersIssueSqlIsJoint {
//			selectOthersIssueSqlJointTag = " or "
//		} else {
//			selectOthersIssueSqlIsJoint = !selectOthersIssueSqlIsJoint
//		}
//		selectOthersIssueSql += selectOthersIssueSqlJointTag + "(ir.relation_id in ? and ir.relation_type = 3)"
//	}
//	if selectOthersIssueSqlIsJoint {
//		queryCond[consts.TcId] = db.In(db.Raw(selectOthersIssueSql, *memberIdsList...))
//	}
//}

//func IssueCondCombinedCondAssembly(queryCond db.Cond, input *vo.HomeIssueInfoReq, currentUserId int64, orgId int64) errs.SystemErrorInfo {
//
//	//if input.CombinedType != nil {
//	//	todayTimeQuantum := times.GetTodayTimeQuantum()
//	//	switch *input.CombinedType {
//	//	//1: 今日指派给我
//	//	case 1:
//	//		queryCond[consts.TcOwnerChangeTime] = db.Between(date.Format(todayTimeQuantum[0]), date.Format(todayTimeQuantum[1]))
//	//		queryCond[" "+consts.TcOwner] = currentUserId
//	//	//2: 最近截止，展示已逾期任务和预计时间小于后天凌晨的任务
//	//	case 2:
//	//		tomorrowTime := todayTimeQuantum[1].Add(time.Duration(60*60*24) * time.Second)
//	//		queryCond[consts.TcPlanEndTime] = db.Between(consts.BlankElasticityTime, date.Format(tomorrowTime))
//	//		queryCond[" "+consts.TcOwner] = currentUserId
//	//	//3: 今日逾期
//	//	case 3:
//	//		nowTime := time.Now()
//	//		todayBegin := date.GetZeroTime(nowTime)
//	//		todayEnd := date.GetZeroTime(nowTime).Add((86400 - 1) * time.Second)
//	//		//今日到期并且尚未完成
//	//		queryCond[consts.TcPlanEndTime] = db.Between(date.Format(todayBegin), date.Format(todayEnd))
//	//		//未完成
//	//		allStatusGroup, err := GetIssueAllStatus(orgId, []int64{}, []int64{tableId})
//	//		//finishedId, err := consts.GetStatusIdsByCategory(orgId, consts.ProcessStatusCategoryIssue, consts.StatusTypeComplete)
//	//		if err != nil {
//	//			log.Errorf(getProcessError, err)
//	//			return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
//	//		}
//	//		if allStatus, ok := allStatusGroup[tableId]; ok {
//	//			finishedId := GetIssueAllStatusIdsByTypes(allStatus, []int{consts.StatusTypeNotStart, consts.StatusTypeRunning})
//	//			queryCond[consts.TcStatus] = db.NotIn(finishedId)
//	//		}
//	//	//4:逾期完成
//	//	case 4:
//	//		//逾期
//	//		queryCond[consts.TcPlanEndTime+" "] = db.Gt(consts.BlankElasticityTime)
//	//		queryCond[consts.TcEndTime] = db.Gt(db.Raw(consts.TcPlanEndTime))
//	//	//5: 即将逾期:预计时间小于后天凌晨的任务
//	//	case 5:
//	//		tomorrowTime := todayTimeQuantum[1].Add(time.Duration(60*60*24) * time.Second)
//	//		queryCond[consts.TcPlanEndTime] = db.Between(date.Format(time.Now()), date.Format(tomorrowTime))
//	//	//6：今日创建的任务
//	//	case 6:
//	//		queryCond[consts.TcCreateTime] = db.Between(date.Format(todayTimeQuantum[0]), date.Format(todayTimeQuantum[1]))
//	//	//7：今日完成的任务
//	//	case 7:
//	//		queryCond[consts.TcEndTime] = db.Between(date.Format(todayTimeQuantum[0]), date.Format(todayTimeQuantum[1]))
//	//		//8:今日添加我为关注人的任务
//	//	case 8:
//	//		queryCond[consts.TcId] = db.In(db.Raw("select ir.issue_id as id from ppm_pri_issue_relation ir where ir.is_delete = 2 and ir.relation_id = ? and ir.relation_type = ? and create_time between ? and ?",
//	//			currentUserId, consts.IssueRelationTypeFollower, date.Format(todayTimeQuantum[0]), date.Format(todayTimeQuantum[1])))
//	//	case 9:
//	//		//今日添加我为审批人的任务
//	//		queryCond[consts.TcId] = db.In(db.Raw("select ir.issue_id as id from ppm_pri_issue_relation ir where ir.is_delete = 2 and ir.relation_id = ? and ir.relation_type = ? and create_time between ? and ?",
//	//			currentUserId, consts.IssueRelationTypeAuditor, date.Format(todayTimeQuantum[0]), date.Format(todayTimeQuantum[1])))
//	//		queryCond[consts.TcProjectId+"  "] = db.In(db.Raw("select id from ppm_pro_project where project_type_id = 1 and org_id in ?", []int64{0, orgId}))
//	//	case 10:
//	//		//10:今日分配给我审批，待我审批的（审批人是我，我还没有审批的）
//	//		queryCond[consts.TcId] = db.In(db.Raw("select ir.issue_id as id from ppm_pri_issue_relation ir where ir.is_delete = 2 and ir.relation_id = ? and ir.relation_type = ? and status in ? and create_time between ? and ?",
//	//			currentUserId, consts.IssueRelationTypeAuditor, []int{consts.AuditStatusNotView, consts.AuditStatusView}, date.Format(todayTimeQuantum[0]), date.Format(todayTimeQuantum[1])))
//	//		queryCond[consts.TcProjectId+"  "] = db.In(db.Raw("select id from ppm_pro_project where project_type_id = 1 and org_id in ?", []int64{0, orgId}))
//	//		queryCond[consts.TcAuditStatus+"  "] = consts.AuditStatusNotView
//	//	}
//	//}
//
//	return nil
//}

//增量查询条件封装
func IssueCondLastUpdateTimeCondAssembly(queryCond db.Cond, input *vo.HomeIssueInfoReq) {
	if input.LastUpdateTime != nil {
		//删除is_delete条件，因为增量查询要将删除的变动也查出来
		delete(queryCond, consts.TcIsDelete)
		queryCond[consts.TcUpdateTime] = db.Gte(date.FormatTime(*input.LastUpdateTime))
	}
}

//// SetIssueAttr 从无码查询任务数据后，将一些值 set 到 bo obj 的属性中
//func SetIssueAttr(issueBoList []bo.IssueBo) {
//	for i, _ := range issueBoList {
//		SetIssueAttrOfTableId(&issueBoList[i])
//	}
//}

//// SetIssueAttrOfTableId 将 lcData 中的 tableId set 到 bo 对象中
//func SetIssueAttrOfTableId(issueBo *bo.IssueBo) {
//	var oriErr error
//	defaultTableId := int64(0)
//	if valIf, ok := issueBo.LessData[consts.BasicFieldTableId]; ok {
//		defaultTableId, oriErr = cast.ToInt64E(valIf)
//		if oriErr != nil {
//			log.Infof("[SetIssueAttrOfTableId] err: %v", oriErr)
//		}
//	}
//	issueBo.TableId = defaultTableId
//}

// ConvertInterfaceIntoInt64 将存在 interface{} 中的数值取出并转为 int64
func ConvertInterfaceIntoInt64(interfaceVal interface{}) int64 {
	var oriErr error
	defaultValInt64 := int64(0)
	if val, ok := interfaceVal.(float64); ok {
		//理论上 map 解析json 会将int转为float64
		valStr := strconv.FormatFloat(val, 'f', -1, 64)
		defaultValInt64, oriErr = strconv.ParseInt(valStr, 10, 64)
		if oriErr != nil {
			log.Infof("[ConvertInterfaceIntoInt64] err: %v", oriErr)
		}
	} else if val, ok := interfaceVal.(int64); ok {
		defaultValInt64 = val
	} else if val, ok := interfaceVal.(string); ok {
		valInt64, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			log.Errorf("[ConvertInterfaceIntoInt64] 无法解析的整型数 from string。val: %v", interfaceVal)
		} else {
			defaultValInt64 = valInt64
		}
	} else {
		log.Errorf("[ConvertInterfaceIntoInt64] 无法解析的整型数。val: %v", interfaceVal)
	}

	return defaultValInt64
}
