package service

import (
	"fmt"
	"strconv"

	automationPb "github.com/star-table/interface/golang/automation/v1"
	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"github.com/spf13/cast"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetProjectTemplateInner(req projectvo.GetProjectTemplateReq) (*projectvo.GetProjectTemplateData, errs.SystemErrorInfo) {
	projectObjectTypeTemplates := make([]projectvo.ProjectObjectTypeTemplate, 0)

	project, err := domain.GetProject(req.OrgId, req.ProjectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	projectDetail, projectDetailErr := domain.GetProjectDetailByProjectIdBo(req.ProjectId, req.OrgId)
	if projectDetailErr != nil {
		log.Error(projectDetailErr)
		return nil, projectDetailErr
	}

	//projectObjectTypeProcesses, err := domain.GetProjectObjectTypeProcessByIds(projectObjectTypeIds)
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//
	//projectAllStatusRespVo := processfacade.GetProjectAllStatus(processvo.GetProjectAllStatusReqVo{
	//	OrgId:      req.OrgId,
	//	ProjectIds: []int64{req.ProjectId},
	//})
	//if projectAllStatusRespVo.Failure() {
	//	log.Error(projectAllStatusRespVo.Error())
	//	return nil, projectAllStatusRespVo.Error()
	//}
	//projectAllStatusList := projectAllStatusRespVo.ProcessStatusBoList[req.ProjectId]
	//processStatusMap := map[int64][]bo.CacheProcessStatusBo{}
	//for _, status := range projectAllStatusList {
	//	processStatusMap[status.ProcessId] = append(processStatusMap[status.ProcessId], status)
	//}

	//for _, projectObjectTypeProcess := range projectObjectTypeProcesses {
	//	projectObjectType, ok := projectObjectTypeMap[projectObjectTypeProcess.TableId]
	//	if !ok {
	//		continue
	//	}
	//	projectObjectTypeTemplate := projectvo.ProjectObjectTypeTemplate{
	//		ID:         projectObjectType.Id,
	//		Name:       projectObjectType.Name,
	//		ProcessId:  projectObjectTypeProcess.ProcessId,
	//		LangCode:   projectObjectType.LangCode,
	//		Sort:       projectObjectType.Sort,
	//		ObjectType: projectObjectType.ObjectType,
	//	}
	//
	//	//敏捷项目传递状态
	//	if project.ProjectTypeId == consts.ProjectTypeAgileId {
	//		if status, ok := processStatusMap[projectObjectTypeProcess.ProcessId]; ok {
	//			statusTemplates := make([]projectvo.ProjectStatusTemplate, 0)
	//			for _, s := range status {
	//				isInitStatus := 2
	//				if s.IsInit {
	//					isInitStatus = 1
	//				}
	//				statusTemplates = append(statusTemplates, projectvo.ProjectStatusTemplate{
	//					ID:           s.StatusId,
	//					Name:         s.Name,
	//					Type:         s.StatusType,
	//					BgStyle:      s.BgStyle,
	//					FontStyle:    s.FontStyle,
	//					IsInitStatus: isInitStatus,
	//					Sort:         s.Sort,
	//					Category:     s.Category,
	//				})
	//			}
	//			projectObjectTypeTemplate.Status = statusTemplates
	//		}
	//	}
	//
	//	projectObjectTypeTemplates = append(projectObjectTypeTemplates, projectObjectTypeTemplate)
	//}

	res := &projectvo.GetProjectTemplateData{
		ProjectObjectTypeTemplates: projectObjectTypeTemplates,
		ProjectIterationTemplates:  []projectvo.ProjectIterationTemplate{},
		//任务模板走无码查询
		//ProjectIssueTemplates:      projectIssueTemplates,
		ProjectInfo: projectvo.ProjectInfoData{
			AppId:         project.AppId,
			OrgId:         project.OrgId,
			Name:          project.Name,
			Status:        project.Status,
			PublishStatus: project.PublicStatus,
			ProjectTypeId: project.ProjectTypeId,
			Owner:         project.Owner,
			PublicStatus:  project.PublicStatus,
			ResourceId:    project.ResourceId,
			IsFiling:      project.IsFiling,
			Remark:        project.Remark,
		},
		ProjectDetail: projectvo.ProjectDetailData{
			Notice:            projectDetail.Notice,
			IsEnableWorkHours: projectDetail.IsEnableWorkHours,
			IsSyncOutCalendar: projectDetail.IsSyncOutCalendar,
		},
	}
	//迭代处理
	if project.ProjectTypeId == consts.ProjectTypeAgileId {
		projectIterationTemplates, projectIterationTemplatesErr := domain.GetProjectIterationInfo(req.OrgId, req.ProjectId)
		if projectIterationTemplatesErr != nil {
			log.Error(projectIterationTemplatesErr)
			return nil, projectIterationTemplatesErr
		}
		res.ProjectIterationTemplates = projectIterationTemplates
	}

	return res, nil
}

func ChangeProjectChatMember(req projectvo.ChangeProjectMemberReq) errs.SystemErrorInfo {
	return nil
}

func ApplyProjectTemplateInner(req *projectvo.ApplyProjectTemplateReq) (*projectvo.ApplyProjectTemplateData, errs.SystemErrorInfo) {
	input := req.Input
	projectInfo := input.ProjectInfo
	projectDetail := input.ProjectDetail
	isUploadTemplate := req.OrgId != projectInfo.OrgId

	authFunctionErr := domain.AuthPayProjectNum(req.OrgId, consts.FunctionProjectCreate)
	if authFunctionErr != nil {
		log.Error(authFunctionErr)
		return nil, authFunctionErr
	}

	orgResp := orgfacade.GetBaseOrgInfo(orgvo.GetBaseOrgInfoReqVo{OrgId: req.OrgId})
	if orgResp.Failure() {
		log.Errorf("[ApplyProjectTemplateInner] GetBaseOrgInfo failed: %v, orgId: %d", orgResp.Error(), req.OrgId)
		return nil, orgResp.Error()
	}

	//获取原table与现table的映射关系
	oldAppId := projectInfo.AppId
	newAppId := req.AppId
	originTableList, originTableListErr := domain.GetAppTableList(projectInfo.OrgId, projectInfo.AppId)
	if originTableListErr != nil {
		log.Errorf("[ApplyProjectTemplateInner] GetAppTableList failed: %v, orgId: %d, appId: %d", originTableListErr, projectInfo.OrgId, projectInfo.AppId)
		return nil, originTableListErr
	}
	curTableList, curTableListErr := domain.GetAppTableList(req.OrgId, req.AppId)
	if curTableListErr != nil {
		log.Errorf("[ApplyProjectTemplateInner] GetAppTableList failed: %v, orgId: %d, appId: %d", curTableListErr, req.OrgId, req.AppId)
		return nil, curTableListErr
	}
	tableMap := make(map[string]string, len(originTableList))
	curLen := len(curTableList)
	for i, table := range originTableList {
		if i < curLen {
			tableMap[strconv.FormatInt(table.TableId, 10)] = strconv.FormatInt(curTableList[i].TableId, 10)
		} else {
			tableMap[strconv.FormatInt(table.TableId, 10)] = "0"
		}
	}

	// 获取原来的菜单列表，将现在的tableId换上去
	menuConfigOrigin, menuConfigErr := domain.GetMenu(projectInfo.OrgId, projectInfo.AppId)
	if menuConfigErr != nil {
		log.Errorf("[ApplyProjectTemplateInner] domain.GetMenu failed: %v, orgId: %d, appId: %d", menuConfigErr, projectInfo.OrgId, projectInfo.AppId)
		return nil, menuConfigErr
	}

	menuListOrigin := []bo.MenuConfig{}
	menuMap := menuConfigOrigin.Config
	err2 := copyer.Copy(menuMap["menuList"], &menuListOrigin)
	if err2 != nil {
		log.Errorf("[ApplyProjectTemplateInner] copy异常")
		return nil, errs.ObjectCopyError
	}
	curMenuMap := make(map[string][]bo.MenuConfig, 0)
	curMenuList := []bo.MenuConfig{}
	for _, menu := range menuListOrigin {
		curTableIdStr := menu.Id
		if newTableIdStr, ok := tableMap[menu.Id]; ok {
			curTableIdStr = newTableIdStr
		}
		curMenuList = append(curMenuList, bo.MenuConfig{
			Id:   curTableIdStr,
			Name: menu.Name,
		})
	}
	curMenuMap["menuList"] = curMenuList

	// 将现在整个菜单排序配置复制到新的应用
	menuJsonStr := json.ToJsonIgnoreError(curMenuMap)
	_, errMenu := domain.SaveMenu(req.OrgId, req.AppId, menuJsonStr)
	if errMenu != nil {
		log.Errorf("[ApplyProjectTemplateInner] domain.SaveMenu failed: %v, orgId: %d, appId: %d", errMenu, req.OrgId, req.AppId)
		return nil, errMenu
	}

	// 创建项目
	projectId, projectIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableProject)
	if projectIdErr != nil {
		log.Error(projectIdErr)
		return nil, projectIdErr
	}
	// 项目关联表（负责人）
	relationId, relationIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
	if relationIdErr != nil {
		log.Error(relationIdErr)
		return nil, relationIdErr
	}
	// 如果是新手项目，把全体成员加为项目成员
	var relationIdForDept0 int64
	if req.Input.IsNewbieGuide {
		relationIdForDept0, relationIdErr = idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
		if relationIdErr != nil {
			log.Error(relationIdErr)
			return nil, relationIdErr
		}
	}

	projectDetailId, projectDetailIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectDetail)
	if projectDetailIdErr != nil {
		log.Error(projectDetailIdErr)
		return nil, projectDetailIdErr
	}

	projectTypeId, status, err := domain.GetTypeAndStatus(projectInfo.ProjectTypeId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	projectPo := &po.PpmProProject{
		Id:            projectId,
		AppId:         req.AppId,
		OrgId:         req.OrgId,
		Name:          projectInfo.Name,
		PreCode:       "",
		Owner:         req.UserId,
		Status:        status,
		ProjectTypeId: projectTypeId,
		PublicStatus:  projectInfo.PublicStatus,
		ResourceId:    projectInfo.ResourceId,
		IsFiling:      2,
		Remark:        projectInfo.Remark,
		Creator:       req.UserId,
		Updator:       req.UserId,
	}
	if req.Input.TemplateFlag != nil && *req.Input.TemplateFlag == 1 {
		projectPo.TemplateFlag = *req.Input.TemplateFlag
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 创建项目表
		createProjectErr := mysql.TransInsert(tx, projectPo)
		if createProjectErr != nil {
			log.Error(createProjectErr)
			return createProjectErr
		}

		// 项目关联（负责人）
		createProjectRelationErr := mysql.TransInsert(tx, &po.PpmProProjectRelation{
			Id:           relationId,
			OrgId:        req.OrgId,
			ProjectId:    projectId,
			RelationId:   req.UserId,
			RelationType: consts.IssueRelationTypeOwner,
			Creator:      req.UserId,
			Updator:      req.UserId,
		})
		if createProjectRelationErr != nil {
			log.Error(createProjectRelationErr)
			return createProjectRelationErr
		}

		// 项目关联（全体成员）——只针对新手项目
		if input.IsNewbieGuide {
			createProjectRelationErr = mysql.TransInsert(tx, &po.PpmProProjectRelation{
				Id:           relationIdForDept0,
				OrgId:        req.OrgId,
				ProjectId:    projectId,
				RelationId:   0,
				RelationType: consts.IssueRelationTypeDepartmentParticipant,
				Creator:      req.UserId,
				Updator:      req.UserId,
			})
			if createProjectRelationErr != nil {
				log.Error(createProjectRelationErr)
				return createProjectRelationErr
			}
		}

		// 创建项目详情表
		createProjectDetailErr := mysql.TransInsert(tx, &po.PpmProProjectDetail{
			Id:                projectDetailId,
			OrgId:             req.OrgId,
			ProjectId:         projectId,
			Notice:            projectDetail.Notice,
			IsEnableWorkHours: projectDetail.IsEnableWorkHours,
			IsSyncOutCalendar: projectDetail.IsSyncOutCalendar,
			Creator:           req.UserId,
			Updator:           req.UserId,
		})
		if createProjectDetailErr != nil {
			log.Error(createProjectDetailErr)
			return createProjectDetailErr
		}

		// 创建项目状态关联
		err = domain.InitProjectStatus(req.OrgId, projectId, req.UserId, tx)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
		}
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	// 创建迭代
	iterationMap, iterationErr := domain.CreateIterationFromTemplate(req.OrgId, req.UserId, projectId, input.ProjectIterationTemplates)
	if iterationErr != nil {
		log.Error(iterationErr)
		return nil, iterationErr
	}

	// 创建工作流
	workflowReq := &automationPb.ApplyTemplateReq{
		IsCreate: req.Input.IsCreateTemplate,
		OrgIdMapping: &automationPb.NumberMapping{
			Old: req.Input.FromOrgId,
			New: req.OrgId,
		},
		AppIdMapping: &automationPb.StringMapping{
			Old: cast.ToString(oldAppId),
			New: cast.ToString(newAppId),
		},
		TableIdMappings: []*automationPb.StringMapping{},
	}
	for oldId, newId := range tableMap {
		workflowReq.TableIdMappings = append(workflowReq.TableIdMappings, &automationPb.StringMapping{
			Old: oldId,
			New: newId,
		})
	}
	for oldId, newId := range iterationMap {
		workflowReq.IterationIdMappings = append(workflowReq.IterationIdMappings, &automationPb.StringMapping{
			Old: cast.ToString(oldId),
			New: cast.ToString(newId),
		})
	}

	// 创建任务
	if req.Input.NeedData {
		reply, err := domain.GetRawRows(req.Input.ProjectInfo.OrgId, req.UserId, &tablePb.ListRawRequest{
			Orders: []*tablePb.Order{{
				Column: lc_helper.ConvertToFilterColumn(consts.BasicFieldOrder),
				Asc:    false,
			}},
			Condition: domain.GetRowsAndCondition(domain.GetNoRecycleCondition(
				domain.GetRowsCondition(consts.BasicFieldAppId, tablePb.ConditionType_equal, req.Input.ProjectInfo.AppId, nil),
			)...),
		})
		if err != nil {
			log.Errorf("[ApplyProjectTemplateInner] GetRawRows err:%v", err)
			return nil, err
		}
		if len(reply.Data) > 0 {
			authFunctionErr = domain.AuthPayTask(req.OrgId, consts.FunctionTaskLimit, len(reply.Data))
			if authFunctionErr != nil {
				// 超出限制不用报错，直接不创建任务即可
				log.Error(authFunctionErr)
			} else {
				err = handleApplyTemplateData(req, projectId, projectPo.PreCode, reply.Data, iterationMap, tableMap, isUploadTemplate)
				if err != nil {
					log.Error(err)
					return nil, err
				}
			}
		}
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		projectTrendsBo := bo.ProjectTrendsBo{
			PushType:   consts.PushTypeCreateProject,
			OrgId:      req.OrgId,
			ProjectId:  projectId,
			OperatorId: req.UserId,
			NewValue:   json.ToJsonIgnoreError(projectPo),
			Ext:        ext,
		}
		domain.PushProjectTrends(projectTrendsBo)
		PushAddProjectNotice(req.OrgId, projectId, req.UserId)
	})
	asyn.Execute(func() {
		if req.Input.TemplateFlag == nil || *req.Input.TemplateFlag != 1 {
			//日历
		}
	})
	return &projectvo.ApplyProjectTemplateData{
		ProjectId: projectId,
		//ProjectObjectTypeMap: projectObjectTypeMap,
		//IterationMap: iterationMap,
		//StatusMap: statusMap,
	}, nil
}

func handleApplyTemplateData(input *projectvo.ApplyProjectTemplateReq, projectId int64, projectCode string,
	allTemplateData []map[string]interface{}, iterationMapping map[int64]int64, tableIdMapping map[string]string,
	isUploadTemplate bool) errs.SystemErrorInfo {
	orgId := input.OrgId
	userId := input.UserId
	appId := input.AppId

	// 过滤不合法的任务
	filteredData := make([]map[string]interface{}, 0)
	for _, d := range allTemplateData {
		// 兼容一些模板中有一些任务属于项目，但不属于任何表的任务（用户看不到，可以通过任务模块看到），会导致创建多出几条意料外的任务。
		issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
		tableId := cast.ToInt64(d[consts.BasicFieldTableId])
		if issueId == 0 || tableId == 0 {
			continue
		}
		filteredData = append(filteredData, d)
	}
	allTemplateData = filteredData

	// 生成issue ids
	issueIds, errSys := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssue, len(allTemplateData))
	if errSys != nil {
		log.Errorf("[BatchCreateIssue] ApplyMultiplePrimaryIdRelaxed failed, count: %v, err: %v",
			len(allTemplateData), errSys)
		return errSys
	}

	// 生成issue codes
	preCode := consts.NoProjectPreCode
	if projectCode != "" {
		preCode = projectCode
	} else {
		preCode = fmt.Sprintf("$%d", projectId)
	}
	issueCodes, errSys := idfacade.ApplyMultipleIdRelaxed(orgId, preCode, "", int64(len(allTemplateData)))
	if errSys != nil {
		log.Errorf("[BatchCreateIssue] ApplyMultiplePrimaryIdRelaxed failed, count: %v, err: %v",
			len(allTemplateData), errSys)
		return errSys
	}

	// 找出父子任务关系
	parentChildrenMapping := map[int64][]int64{}                  // 父任务->子任务列表
	allData := map[int64]map[string]interface{}{}                 // 任务id->任务数据
	allParentData := make(map[string][]map[string]interface{}, 0) // oldTableId->父任务
	dataByTable := make(map[string][]map[string]interface{}, 0)   // newTableId->待创建的任务数据
	issueIdMapping := make(map[int64]int64, 0)                    // oldIssueId->newIssueId
	for i, d := range allTemplateData {
		issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
		parentId := cast.ToInt64(d[consts.BasicFieldParentId])
		tableId := cast.ToString(d[consts.BasicFieldTableId])

		allData[issueId] = d

		// 赋予新issueId/code
		newIssueId := issueIds.Ids[i].Id
		d[consts.TempFieldIssueId] = newIssueId
		issueIdMapping[issueId] = newIssueId
		d[consts.TempFieldCode] = issueCodes.Ids[i].Code

		if parentId != 0 {
			// 子任务
			parentChildrenMapping[parentId] = append(parentChildrenMapping[parentId], issueId)
		} else {
			// 父任务
			allParentData[tableId] = append(allParentData[tableId], d)
		}
	}

	var tableColumns map[string]*projectvo.TableColumnData
	for oldTableIdStr, newTableIdStr := range tableIdMapping {
		oldTableId := cast.ToInt64(oldTableIdStr)
		newTableId := cast.ToInt64(newTableIdStr)

		// 获取表头
		tableColumns, errSys = domain.GetTableColumnsMap(orgId, newTableId, nil, true)
		if errSys != nil {
			log.Errorf("[handleApplyTemplateData] GetTableColumnsMap failed, org:%d table:%d, err: %v",
				orgId, newTableId, errSys)
			return errSys
		}

		// 替换关联/前后置/单向关联ID
		replacingRelatingIds(oldTableId, tableColumns, allData, issueIdMapping, input.Input.IsCreateTemplate, isUploadTemplate)
	}

	// 以表为粒度收集数据
	var totalCount int64 // 所有表加起来的任务总数（包括子任务）
	for oldTableId, newTableId := range tableIdMapping {
		data := make([]map[string]interface{}, 0)

		// 处理父任务
		for _, d := range allParentData[oldTableId] {
			issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
			parentData := buildIssueDataByCopy(d, iterationMapping, nil)

			if childrenIds, ok := parentChildrenMapping[issueId]; ok {
				// 处理子任务（递归）
				count, err := buildIssueChildDataByCopy(parentData, childrenIds, allData, parentChildrenMapping, iterationMapping, nil)
				if err != nil {
					log.Errorf("[handleApplyTemplateData] buildIssueChildDataByCopy err: %v", err)
					return err
				}
				totalCount += count
			}

			totalCount += 1
			data = append(data, parentData)
		}

		dataByTable[newTableId] = data
	}

	tableIds := make([]string, 0, len(tableIdMapping))
	for _, tableId := range tableIdMapping {
		tableIds = append(tableIds, tableId)
	}

	// 创建异步任务
	asyncTaskId := domain.GenAsyncTaskIdForApplyTemplate(appId)
	err := domain.CreateAsyncTask(orgId, totalCount, asyncTaskId, map[string]string{
		consts.AsyncTaskHashPartKeyOfCover:    input.Input.TemplateInfo.Icon,
		consts.AsyncTaskHashPartKeyOfTableIds: json.ToJsonIgnoreError(tableIds),
	})
	if err != nil {
		log.Errorf("[handleApplyTemplateData] CreateAsyncTask err: %v", err)
		return err
	}

	// 按表格为单位的异步批量创建
	for _, newTableId := range tableIdMapping {
		tableId := cast.ToInt64(newTableId)
		data := dataByTable[newTableId]
		if len(data) > 0 {
			req := &projectvo.BatchCreateIssueReqVo{
				OrgId:         orgId,
				UserId:        userId,
				AppId:         appId,
				ProjectId:     projectId,
				TableId:       tableId,
				Data:          data,
				IsIdGenerated: true,
			}
			AsyncBatchCreateIssue(req, true, &projectvo.TriggerBy{
				TriggerBy:        consts.TriggerByApplyTemplate,
				IsCreateTemplate: input.Input.IsCreateTemplate,
			}, asyncTaskId)
		}
	}
	return nil
}

// 删除项目
func DeleteProjectBatchInner(orgId, userId int64, projectIds []int64) errs.SystemErrorInfo {
	if len(projectIds) == 0 {
		return nil
	}
	updateErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除项目
		_, updateProjectErr := mysql.TransUpdateSmartWithCond(tx, consts.TableProject, db.Cond{
			consts.TcId: db.In(projectIds),
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  userId,
		})
		if updateProjectErr != nil {
			log.Error(updateProjectErr)
			return updateProjectErr
		}
		//删除项目下的任务
		//_, updateIssueErr := mysql.TransUpdateSmartWithCond(tx, consts.TableIssue, db.Cond{
		//	consts.TcIsDelete:  consts.AppIsNoDelete,
		//	consts.TcProjectId: db.In(projectIds),
		//	consts.TcOrgId:     orgId,
		//}, mysql.Upd{
		//	consts.TcIsDelete: consts.AppIsDeleted,
		//})
		//if updateIssueErr != nil {
		//	log.Error(updateIssueErr)
		//	return updateIssueErr
		//}
		//删除迭代
		_, updateIterationErr := mysql.TransUpdateSmartWithCond(tx, consts.TableIteration, db.Cond{
			consts.TcIsDelete:  consts.AppIsNoDelete,
			consts.TcProjectId: db.In(projectIds),
			consts.TcOrgId:     orgId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if updateIterationErr != nil {
			log.Error(updateIterationErr)
			return updateIterationErr
		}

		return nil
	})
	if updateErr != nil {
		return errs.DeleteProjectErr
	}

	for _, id := range projectIds {
		err1 := domain.RefreshProjectAuthBo(orgId, id)
		if err1 != nil {
			log.Error(err1)
			return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
		}
	}

	return nil
}
