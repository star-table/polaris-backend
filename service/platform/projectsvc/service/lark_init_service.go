package service

//import (
//	"github.com/star-table/common/core/types"
//	"github.com/star-table/common/library/db/mysql"
//	"github.com/star-table/polaris-backend/common/core/consts"
//	"github.com/star-table/polaris-backend/common/core/errs"
//	"github.com/star-table/polaris-backend/common/core/util"
//	"github.com/star-table/polaris-backend/common/model/vo"
//	"github.com/star-table/polaris-backend/common/model/vo/processvo"
//	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
//	"github.com/star-table/polaris-backend/facade/idfacade"
//	"github.com/star-table/polaris-backend/facade/processfacade"
//	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
//	"strconv"
//	"time"
//	"upper.io/db.v3/lib/sqlbuilder"
//)
//
//const larkAppletInitSql = consts.TemplateDirPrefix + "lark_applet_project_issue_init.template"
//
//func DataInitForLarkApplet(orgId, userId int64) errs.SystemErrorInfo {
//	//获取普通项目
//	var projectTypeId int64
//	projectTypes, err := domain.GetProjectTypeList(orgId)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	for _, v := range *projectTypes {
//		if v.LangCode == consts.ProjectTypeLangCodeNormalTask {
//			projectTypeId = v.Id
//		}
//	}
//
//	start := types.NowTime()
//	endTime, _ := time.Parse(consts.AppTimeFormat, "2099-12-12 12:00:00")
//	end := types.Time(endTime)
//	//项目初始化(企业项目)
//	preCode := "QYSLXMYS"
//	remark := "你可参考示例项目，快速熟悉极星协作平台"
//	projectInfo, createErr := CreateProject(projectvo.CreateProjectReqVo{
//		OrgId:         orgId,
//		UserId:        userId,
//		SourceChannel: sdk_const.SourceChannelFeishu,
//		Input: vo.CreateProjectReq{
//			Name:          "企业示例项目演示",
//			PreCode:       &preCode,
//			PublicStatus:  consts.PublicProject,
//			Remark:        &remark,
//			Owner:         userId,
//			MemberIds:     []int64{userId},
//			ResourcePath:  "https://polaris-hd2.oss-cn-shanghai.aliyuncs.com/project/undraw_Projectpicture_update_jjgk.png",
//			ResourceType:  consts.OssResource,
//			PlanStartTime: &start,
//			PlanEndTime:   &end,
//			ProjectTypeID: &projectTypeId,
//		},
//	})
//	if createErr != nil {
//		log.Error(createErr)
//		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, createErr)
//	}
//	log.Info("项目初始化成功")
//
//	objectType := 2
//	demandId, projectObjectTypeErr1 := initProjectObjectType(orgId, userId, projectInfo.ID, objectType, "需求")
//	if projectObjectTypeErr1 != nil {
//		log.Error(projectObjectTypeErr1)
//		return projectObjectTypeErr1
//	}
//	_, projectObjectTypeErr2 := initProjectObjectType(orgId, userId, projectInfo.ID, objectType, "设计")
//	if projectObjectTypeErr2 != nil {
//		log.Error(projectObjectTypeErr2)
//		return projectObjectTypeErr2
//	}
//
//	//项目初始化(私人项目)
//	preCode1 := "SRSLXMYS"
//	remark1 := "你可参考示例项目，快速熟悉极星协作平台"
//	projectInfo1, createErr1 := CreateProject(projectvo.CreateProjectReqVo{
//		OrgId:         orgId,
//		UserId:        userId,
//		SourceChannel: sdk_const.SourceChannelFeishu,
//		Input: vo.CreateProjectReq{
//			Name:          "私人示例项目演示",
//			PreCode:       &preCode1,
//			PublicStatus:  consts.PrivateProject,
//			Remark:        &remark1,
//			Owner:         userId,
//			MemberIds:     []int64{userId},
//			ResourcePath:  "https://polaris-hd2.oss-cn-shanghai.aliyuncs.com/project/undraw_Projectpicture_update_jjgk.png",
//			ResourceType:  consts.OssResource,
//			PlanStartTime: &start,
//			PlanEndTime:   &end,
//			ProjectTypeID: &projectTypeId,
//		},
//	})
//	if createErr1 != nil {
//		log.Error(createErr1)
//		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, createErr1)
//	}
//	log.Info("项目初始化成功")
//
//	privateDemandId, projectObjectTypeErr3 := initProjectObjectType(orgId, userId, projectInfo1.ID, objectType, "需求")
//	if projectObjectTypeErr3 != nil {
//		log.Error(projectObjectTypeErr3)
//		return projectObjectTypeErr3
//	}
//	_, projectObjectTypeErr4 := initProjectObjectType(orgId, userId, projectInfo1.ID, objectType, "设计")
//	if projectObjectTypeErr4 != nil {
//		log.Error(projectObjectTypeErr4)
//		return projectObjectTypeErr4
//	}
//	err1 := mysql.TransX(func(tx sqlbuilder.Tx) error {
//		err := initIssue(orgId, userId, projectInfo.ID, projectInfo1.ID, demandId, privateDemandId, tx)
//		if err != nil {
//			log.Error(err)
//			return err
//		}
//
//		return nil
//	})
//
//	if err1 != nil {
//		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err1)
//	}
//
//	return nil
//}
//
//func initProjectObjectType(orgId, userId int64, projectId int64, objectType int, name string) (int64, errs.SystemErrorInfo) {
//	res, err := CreateProjectObjectType(orgId, userId, vo.CreateProjectObjectTypeReq{
//		ProjectID:  projectId,
//		Name:       name,
//		ObjectType: objectType,
//		BeforeID:   0,
//	})
//	if err != nil {
//		return 0, errs.BuildSystemErrorInfo(errs.BaseDomainError, err)
//	}
//	log.Infof("项目对象类型-%s初始化成功", name)
//	return res.ID, nil
//}
//
//func initIssue(orgId int64, operatorId int64, projectId int64, privateProjectId int64, demandId, privateDemandId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
//	contextMap := map[string]interface{}{}
//	contextMap["OrgId"] = orgId
//	contextMap["UserId"] = operatorId
//	contextMap["ProjectId"] = projectId
//	contextMap["PrivateProjectId"] = privateProjectId
//	contextMap["PrivateObjectTypeDemand"] = privateDemandId
//	contextMap["ObjectTypeDemand"] = demandId
//	contextMap["NowTime"] = types.NowTime()
//
//	statCount := 4
//	//stat id 申请
//	statIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectDayStat, statCount)
//	if err != nil {
//		return err
//	}
//	statIdCount := 1
//	now := time.Now()
//	for _, v := range statIds.Ids {
//		contextMap["StatId"+strconv.Itoa(statIdCount)] = v.Id
//		contextMap["Day"+strconv.Itoa(statIdCount)] = now.AddDate(0, 0, -statIdCount).Format(consts.AppDateFormat)
//		statIdCount++
//	}
//
//	count := 6
//	//issue id 申请
//	issueIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssue, count)
//	if err != nil {
//		return err
//	}
//	issueIdCount := 1
//	for _, v := range issueIds.Ids {
//		contextMap["IssueId"+strconv.Itoa(issueIdCount)] = v.Id
//		issueIdCount++
//	}
//	//issuedetail id 申请
//	issueDetailIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueDetail, count)
//	if err != nil {
//		return err
//	}
//	issueDetailIdCount := 1
//	for _, v := range issueDetailIds.Ids {
//		contextMap["IssueDetailId"+strconv.Itoa(issueDetailIdCount)] = v.Id
//		issueDetailIdCount++
//	}
//
//	//issuerelation id 申请
//	relationIds, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueRelation, count)
//	if err != nil {
//		return err
//	}
//	relationIdCount := 1
//	for _, v := range relationIds.Ids {
//		contextMap["RelationId"+strconv.Itoa(relationIdCount)] = v.Id
//		relationIdCount++
//	}
//
//	//优先级
//	allPriority, err := domain.GetPriorityList(orgId)
//	if err != nil {
//		return err
//	}
//	for _, v := range *allPriority {
//		if v.LangCode == "Priority.Issue.P0" {
//			contextMap["PriorityP0"] = v.Id
//		} else if v.LangCode == "Priority.Issue.P1" {
//			contextMap["PriorityP1"] = v.Id
//		} else if v.LangCode == "Priority.Issue.P3" {
//			contextMap["PriorityP3"] = v.Id
//		} else if v.LangCode == "Priority.Issue.P4" {
//			contextMap["PriorityP4"] = v.Id
//		}
//	}
//
//	//状态（通用项目任务状态）
//	processStatusResp := processfacade.GetProcessStatusListByCategory(processvo.GetProcessStatusListByCategoryReqVo{
//		OrgId:    0,
//		Category: 3,
//	})
//	if processStatusResp.Failure() {
//		return processStatusResp.Error()
//	}
//	for _, v := range processStatusResp.CacheProcessStatusBoList {
//		if v.Name == "未完成" {
//			contextMap["NotStartStatus"] = v.StatusId
//		} else if v.Name == "处理中" {
//			contextMap["RunningStatus"] = v.StatusId
//		}
//	}
//
//	insertErr := util.ReadAndWrite(larkAppletInitSql, contextMap, tx)
//	if insertErr != nil {
//		log.Error(insertErr)
//		return errs.BuildSystemErrorInfo(errs.BaseDomainError, insertErr)
//	}
//
//	return nil
//}
