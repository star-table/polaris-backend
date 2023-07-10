package service

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/star-table/polaris-backend/facade/idfacade"

	"github.com/star-table/polaris-backend/common/model/vo/uservo"

	"github.com/star-table/polaris-backend/common/core/util/jsonx"

	msgPb "github.com/star-table/interface/golang/msg/v1"
	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/common/report"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/spf13/cast"
)

func JudgeProjectFiling(orgId, projectId int64) errs.SystemErrorInfo {
	projectInfo, err := domain.LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	if projectInfo.IsFilling == consts.ProjectIsFiling {
		return errs.ProjectIsFilingYet
	}

	return nil
}

func JudgeProjectFilingByIssueId(orgId, issueId int64) errs.SystemErrorInfo {
	issueBo, err := domain.GetIssueInfoLc(orgId, 0, issueId)
	if err != nil {
		log.Error(err)
		return err
	}
	projectId := issueBo.ProjectId
	if projectId != 0 {
		projectInfo, err := domain.LoadProjectAuthBo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return err
		}
		if projectInfo.IsFilling == consts.ProjectIsFiling {
			return errs.ProjectIsFilingYet
		}
	}

	return nil
}

func DeleteIssueWithoutAuth(reqVo projectvo.DeleteIssueReqVo) (*vo.Issue, errs.SystemErrorInfo) {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	sourceChannel := reqVo.SourceChannel

	issueBos, err2 := domain.GetIssueInfosLc(orgId, reqVo.UserId, []int64{input.ID})
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}
	if len(issueBos) < 1 {
		log.Errorf("[DeleteIssueWithoutAuth] not found issue issueId:%v", input.ID)
		return nil, errs.IssueNotExist
	}
	issueBo := issueBos[0]
	err := domain.AuthIssue(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Delete)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	childIds, err3 := domain.DeleteIssue(issueBo, currentUserId, sourceChannel, input.TakeChildren)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	// 事件上报
	asyn.Execute(func() {
		// 拿任务信息
		var deletedIssueIds []int64
		deletedIssueIds = append(deletedIssueIds, input.ID)
		deletedIssueIds = append(deletedIssueIds, childIds...)
		columnIds := []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAppId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
		}
		data, errSys := domain.GetIssueInfosMapLcByIssueIds(orgId, currentUserId, deletedIssueIds, columnIds...)
		if errSys != nil {
			return
		}

		openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		openTraceIdStr := cast.ToString(openTraceId)

		for _, d := range data {
			dataId := cast.ToInt64(d[consts.BasicFieldId])
			issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
			appId := cast.ToInt64(d[consts.BasicFieldAppId])
			projectId := cast.ToInt64(d[consts.BasicFieldProjectId])
			tableId := cast.ToInt64(d[consts.BasicFieldTableId])
			e := &commonvo.DataEvent{
				OrgId:     orgId,
				AppId:     appId,
				ProjectId: projectId,
				TableId:   tableId,
				DataId:    dataId,
				IssueId:   issueId,
				UserId:    currentUserId,
			}
			report.ReportDataEvent(msgPb.EventType_DataDeleted, openTraceIdStr, e)
		}

	})

	result := &vo.Issue{}
	copyErr := copyer.Copy(issueBo, result)
	if copyErr != nil {
		log.Errorf("copyer.Copy(): %q\n", copyErr)
	}
	result.IssueIds = append(childIds, input.ID)
	return result, nil
}

func DeleteIssue(reqVo projectvo.DeleteIssueReqVo) (*vo.Issue, errs.SystemErrorInfo) {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId

	issueBo, err2 := domain.GetIssueInfoLc(orgId, 0, input.ID)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err2)
	}

	err := domain.AuthIssue(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Delete)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}
	return DeleteIssueWithoutAuth(reqVo)
}

func DeleteIssueBatch(reqVo projectvo.DeleteIssueBatchReqVo) (*vo.DeleteIssueBatchResp, errs.SystemErrorInfo) {
	input := reqVo.Input
	orgId := reqVo.OrgId
	currentUserId := reqVo.UserId
	sourceChannel := reqVo.SourceChannel

	//查找目标任务以及所有的子任务
	issueAndChildrenIds, err := domain.GetIssueAndChildrenIds(orgId, input.Ids)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	issueAndChildren, err := domain.GetIssueInfosLc(orgId, currentUserId, issueAndChildrenIds)
	if err != nil {
		log.Errorf("[DeleteIssueBatch] GetIssueInfosLc err: %v", err)
		return nil, errs.IssueNotExist
	}

	allNeedAuth := []*bo.IssueBo{}
	allNeedAuthIssueIds := make([]int64, 0)
	for _, child := range issueAndChildren {
		if ok, _ := slice.Contain(input.Ids, child.Id); ok {
			if child.ParentId == 0 {
				allNeedAuth = append(allNeedAuth, child)
				allNeedAuthIssueIds = append(allNeedAuthIssueIds, child.Id)
			} else if ok1, _ := slice.Contain(input.Ids, child.ParentId); !ok1 {
				//找出需要判断权限的(只要父任务不在所选范围内，就表示他是父任务)
				allNeedAuth = append(allNeedAuth, child)
				allNeedAuthIssueIds = append(allNeedAuthIssueIds, child.Id)
			}
		}
	}

	notAuthPassIssues := []*bo.IssueBo{}
	noAuthIds := []int64{}
	trulyIssueBos := []*bo.IssueBo{}

	var notAuthPath []string
	for _, issueBo := range allNeedAuth {
		//权限判断
		err := domain.AuthIssueWithAppId(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Delete, reqVo.InputAppId)
		if err != nil {
			//log.Error(err)
			notAuthPassIssues = append(notAuthPassIssues, issueBo)
			noAuthIds = append(noAuthIds, issueBo.Id)
			notAuthPath = append(notAuthPath, fmt.Sprintf("%s%d,", issueBo.Path, issueBo.Id))
		} else {
			trulyIssueBos = append(trulyIssueBos, issueBo)
		}
	}

	//移除所有不通过权限校验的以及他们的子任务
	allAuthedIssues := make([]*bo.IssueBo, 0)
	if len(notAuthPath) > 0 {
		for _, child := range issueAndChildren {
			if ok, _ := slice.Contain(noAuthIds, child.Id); ok {
				continue
			}
			pass := true
			if ok, _ := slice.Contain(noAuthIds, child.ParentId); ok {
				for _, s := range notAuthPath {
					if strings.Contains(child.Path, s) {
						pass = false
						break
					}
				}
			}

			if pass {
				allAuthedIssues = append(allAuthedIssues, child)
			}
		}
	} else {
		allAuthedIssues = issueAndChildren
	}
	tableId, tableIdErr := strconv.ParseInt(input.TableID, 10, 64)
	if tableIdErr != nil {
		log.Errorf("[DeleteIssueBatch] tableId参数错误：%v, inputTableId:%s", tableIdErr, input.TableID)
		return nil, errs.ParamError
	}
	deletedIssueIds, err3 := domain.DeleteIssueBatch(orgId, allAuthedIssues, input.Ids, currentUserId, sourceChannel, input.ProjectID, tableId)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
	}

	// 事件上报
	asyn.Execute(func() {
		// 拿任务信息
		columnIds := []string{
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTitle),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldAppId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
			lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
		}
		data, errSys := domain.GetIssueInfosMapLcByIssueIds(orgId, currentUserId, deletedIssueIds, columnIds...)
		if errSys != nil {
			return
		}

		issueProjectId := int64(0)
		openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		openTraceIdStr := cast.ToString(openTraceId)

		for _, d := range data {
			dataId := cast.ToInt64(d[consts.BasicFieldId])
			issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
			appId := cast.ToInt64(d[consts.BasicFieldAppId])
			projectId := cast.ToInt64(d[consts.BasicFieldProjectId])
			tableId := cast.ToInt64(d[consts.BasicFieldTableId])
			e := &commonvo.DataEvent{
				OrgId:     orgId,
				AppId:     appId,
				ProjectId: projectId,
				TableId:   tableId,
				DataId:    dataId,
				IssueId:   issueId,
				UserId:    currentUserId,
			}
			issueProjectId = projectId
			report.ReportDataEvent(msgPb.EventType_DataDeleted, openTraceIdStr, e)
		}

		updateDingTopCard(orgId, issueProjectId)
	})

	resBo := bo.DeleteIssueBatchBo{
		NoAuthIssues:         notAuthPassIssues,
		RemainChildrenIssues: []*bo.IssueBo{},
	}

	res := &vo.DeleteIssueBatchResp{}
	for _, child := range allAuthedIssues {
		if ok, _ := slice.Contain(input.Ids, child.Id); ok {
			resBo.SuccessIssues = append(resBo.SuccessIssues, child)
		}
	}

	_ = copyer.Copy(resBo, res)

	return res, nil
}

func MoveIssue(orgId, userId, appId int64, input vo.MoveIssueReq) (*vo.Void, errs.SystemErrorInfo) {
	//获取原本的Issue
	issueBo, errSys := GetIssueInfo(orgId, userId, input.ID)
	if errSys != nil {
		log.Errorf("[MoveIssue] GetIssueInfo err: %v issueId: %v", errSys, input.ID)
		return nil, errSys
	}

	//查找目标任务以及所有的子任务
	allMoveIssueIds, errSys := domain.GetIssueAndChildrenIds(orgId, input.ChildrenIds)
	if errSys != nil {
		log.Errorf("[MoveIssue] GetIssueAndChildrenIds err: %v", errSys)
		return nil, errSys
	}
	allMoveIssueIds = append(allMoveIssueIds, input.ID)

	issueAndChildrenIds, errSys := domain.GetIssueAndChildrenIds(orgId, []int64{input.ID})
	if errSys != nil {
		log.Errorf("[MoveIssue] GetIssueAndChildrenIds err: %v", errSys)
		return nil, errSys
	}

	batchInput := vo.MoveIssueBatchReq{
		Ids:           allMoveIssueIds,
		FromProjectID: issueBo.ProjectId,
		FromTableID:   cast.ToString(issueBo.TableId),
		ProjectID:     input.ProjectID,
		TableID:       input.TableID,
		ChooseField:   input.ChooseField,
	}
	issueTitleMap := map[int64]string{}
	issueTitleMap[input.ID] = input.Title
	_, errSys = moveIssueBatch(orgId, userId, batchInput, issueAndChildrenIds, issueTitleMap)
	return &vo.Void{ID: input.ID}, errSys
}

func MoveIssueBatch(orgId, userId int64, input vo.MoveIssueBatchReq) (*vo.MoveIssueBatchResp, errs.SystemErrorInfo) {
	//查找目标任务以及所有的子任务
	issueAndChildrenIds, errSys := domain.GetIssueAndChildrenIds(orgId, input.Ids)
	if errSys != nil {
		log.Errorf("[MoveIssueBatch] GetIssueAndChildrenIds err: %v", errSys)
		return nil, errSys
	}

	return moveIssueBatch(orgId, userId, input, issueAndChildrenIds, nil)
}

func moveIssueBatch(orgId, userId int64, input vo.MoveIssueBatchReq, issueAndChildrenIds []int64, issueTitleMap map[int64]string) (*vo.MoveIssueBatchResp, errs.SystemErrorInfo) {
	log.Infof("[MoveIssueBatch]: %v", json.ToJsonIgnoreError(input))
	res := &vo.MoveIssueBatchResp{}

	fromProjectId := input.FromProjectID
	fromTableId := cast.ToInt64(input.FromTableID)
	targetProjectId := input.ProjectID
	targetTableId := cast.ToInt64(input.TableID)
	if fromProjectId == targetProjectId &&
		fromTableId == targetTableId {
		return res, nil
	}

	issueAndChildrenIssues, errSys := domain.GetIssueInfosLc(orgId, userId, issueAndChildrenIds)
	if errSys != nil {
		log.Errorf("[MoveIssueBatch] GetIssueInfosLc err: %v", errSys)
		return nil, errSys
	}

	var allNeedAuthIssues []*bo.IssueBo
	var noAuthPath []string
	var noAuthIssueIds []int64
	var noAuthIssues []*bo.IssueBo
	var allAuthedIssues []*bo.IssueBo

	for _, issue := range issueAndChildrenIssues {
		//如果父任务为0，或者父任务不在传入的任务id里面，这些任务将作为父任务进行权限校验
		if issue.ParentId == int64(0) {
			allNeedAuthIssues = append(allNeedAuthIssues, issue)
		} else {
			if ok, _ := slice.Contain(input.Ids, issue.ParentId); !ok {
				allNeedAuthIssues = append(allNeedAuthIssues, issue)
			}
		}
	}

	// 权限判断
	if fromProjectId != targetProjectId && targetProjectId > 0 {
		errSys = AuthProject(orgId, userId, targetProjectId, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create)
		if errSys != nil {
			log.Errorf("[MoveIssueBatch] AuthProject err: %v", errSys)
			return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, errSys)
		}
	}
	// TODO: 优化成批量
	for _, issueBo := range allNeedAuthIssues {
		errSys = domain.AuthIssueWithAppId(orgId, userId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify, 0)
		if errSys != nil {
			log.Infof("[MoveIssueBatch] AuthIssueWithAppId err: %v", errSys)
			noAuthIssues = append(noAuthIssues, issueBo)
			noAuthIssueIds = append(noAuthIssueIds, issueBo.Id)
			noAuthPath = append(noAuthPath, fmt.Sprintf("%s%d,", issueBo.Path, issueBo.Id))
		}
	}

	// 移除所有不通过权限校验的以及他们的子任务
	if len(noAuthPath) > 0 {
		for _, issue := range issueAndChildrenIssues {
			if ok, _ := slice.Contain(noAuthIssueIds, issue.Id); ok {
				continue
			}
			pass := true
			for _, s := range noAuthPath {
				if strings.Contains(issue.Path, s) {
					pass = false
					break
				}
			}
			if pass {
				allAuthedIssues = append(allAuthedIssues, issue)
			}
		}
	} else {
		allAuthedIssues = issueAndChildrenIssues
	}

	//更新任务的项目对象类型
	allMoveIds, err := domain.MoveIssueProTableBatch(orgId, userId, allAuthedIssues, input.Ids, fromTableId, fromProjectId, targetTableId, targetProjectId, input.ChooseField, issueTitleMap)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err)
	}

	resBo := bo.MoveIssueBatchBo{
		NoAuthIssues:         noAuthIssues,
		RemainChildrenIssues: []*bo.IssueBo{},
		ChildrenIssues:       []*bo.IssueBo{},
	}
	for _, issue := range issueAndChildrenIssues {
		if ok, _ := slice.Contain(allMoveIds, issue.Id); ok {
			resBo.SuccessIssues = append(resBo.SuccessIssues, issue)
		}
	}
	_ = copyer.Copy(resBo, res)

	// 上报事件
	asyn.Execute(func() {
		ReportMoveEvents(orgId, userId, issueAndChildrenIssues, issueAndChildrenIds, fromTableId, targetTableId)
	})

	return res, nil
}

func GetIssueInfo(orgId, userId int64, id int64) (*bo.IssueBo, errs.SystemErrorInfo) {
	bos, err := domain.GetIssueInfosLc(orgId, userId, []int64{id})
	if err != nil {
		return nil, err
	}
	if len(bos) == 0 {
		return nil, errs.IssueNotExist
	}
	return bos[0], nil
}

func GetIssueInfoList(orgId, userId int64, ids []int64) ([]vo.Issue, errs.SystemErrorInfo) {
	result, err := domain.GetIssueInfosLc(orgId, userId, ids)
	if err != nil {
		return nil, err
	}
	issueResp := make([]vo.Issue, 0, len(result))
	err1 := copyer.Copy(result, &issueResp)
	if err1 != nil {
		log.Errorf("copyer.Copy: %q\n", err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}
	// 赋值负责人 ids
	for i, issue := range result {
		issueResp[i].Owners = issue.OwnerIdI64
	}

	return issueResp, nil
}

func GetIssueInfoListByDataIds(orgId, userId int64, dataIds []int64) ([]vo.Issue, errs.SystemErrorInfo) {
	result, err := domain.GetIssueInfosLcByDataIds(orgId, userId, dataIds)
	if err != nil {
		return nil, err
	}
	issueResp := &[]vo.Issue{}
	err1 := copyer.Copy(result, issueResp)
	if err1 != nil {
		log.Errorf("copyer.Copy: %q\n", err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err1)
	}

	return *issueResp, nil
}

// replacingRelatingIds 替换关联/前后置/单向关联ID
func replacingRelatingIds(tableId int64,
	tableColumns map[string]*projectvo.TableColumnData,
	allData map[int64]map[string]interface{}, issueIdMapping map[int64]int64,
	isCreateTemplate, isUploadTemplate bool) {

	// 替换关联/前后置/单向关联ID
	relatingColumnIds := []string{consts.BasicFieldRelating, consts.BasicFieldBaRelating}
	if tableColumns != nil {
		for columnId, column := range tableColumns {
			if columnId != consts.BasicFieldRelating && columnId != consts.BasicFieldBaRelating {
				if column.Field.Type == consts.LcColumnFieldTypeRelating || column.Field.Type == consts.LcColumnFieldTypeSingleRelating {
					relatingColumnIds = append(relatingColumnIds, columnId)
				}
			}
		}
	}

	for issueId, d := range allData {
		issueTableId := cast.ToInt64(d[consts.BasicFieldTableId])
		if tableId != issueTableId {
			continue
		}
		for _, columnId := range relatingColumnIds {
			if relatingI, ok := d[columnId]; ok && relatingI != nil {
				relating := &bo.RelatingIssue{}
				if err := jsonx.Copy(relatingI, relating); err != nil {
					log.Errorf("[ReplacingRelatingIds] parse relating column failed, issueId: %v, columnId: %v, value: %v, err: %v",
						issueId, columnId, json.ToJsonIgnoreError(relatingI), err)
					delete(d, columnId)
				} else {
					var linkTo, linkFrom []string
					for _, id := range relating.LinkTo {
						if newIssueId, ok := issueIdMapping[cast.ToInt64(id)]; ok {
							linkTo = append(linkTo, cast.ToString(newIssueId))
						} else if !isUploadTemplate {
							linkTo = append(linkTo, id)
						}
					}
					for _, id := range relating.LinkFrom {
						if newIssueId, ok := issueIdMapping[cast.ToInt64(id)]; ok {
							linkFrom = append(linkFrom, cast.ToString(newIssueId))
						} else if !isUploadTemplate && !isCreateTemplate {
							linkFrom = append(linkFrom, id)
						}
					}
					relating.LinkTo = linkTo
					relating.LinkFrom = linkFrom
					d[columnId] = relating
				}
			}
		}
	}
}

// getCopyIssueDeleteColumnIds 复制任务的时候，必须清空以便重新设置正确值的字段
func getCopyIssueDeleteColumnIds() []string {
	// 暂不支持工时
	deleteColumnIds := []string{consts.ProBasicFieldWorkHour}
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldOrgId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldAppId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldProjectId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldTableId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldIssueId)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldCode)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldOrder)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldRecycleFlag)
	deleteColumnIds = append(deleteColumnIds, consts.BasicFieldCollaborators)
	return deleteColumnIds
}

// buildIssueDataByCopy 通过复制一条任务数据构造一条新任务数据
// data 复制数据源
// iterationMapping 迭代值映射mapping。传nil则保留原值(同一个项目中复制可以沿用迭代值，跨项目则需要映射)
// keepTableColumns 复制时保留的列，其余的列不会被复制。传nil则全量复制(同一张表中复制任务可以全量复制，应用模板时由于目标表格复制了原表的表头，也可以全量复制）
func buildIssueDataByCopy(data map[string]interface{},
	iterationMapping map[int64]int64,
	keepTableColumns map[string]struct{}) map[string]interface{} {
	newData := make(map[string]interface{})
	copyer.Copy(data, &newData)

	// 处理迭代（迭代没有无码化，导致这里需要映射）
	if iterationMapping != nil {
		if iterationIdI, ok := newData[consts.BasicFieldIterationId]; ok {
			if iterationId, ok := iterationMapping[cast.ToInt64(iterationIdI)]; ok {
				newData[consts.BasicFieldIterationId] = iterationId
			} else {
				delete(newData, consts.BasicFieldIterationId)
			}
		}
	}

	// 删除表头里没有的字段
	if keepTableColumns != nil {
		for k, _ := range newData {
			if k != consts.TempFieldIssueId && k != consts.TempFieldCode {
				if _, ok := keepTableColumns[k]; !ok {
					delete(newData, k)
				}
			}
		}
	}

	// 删除确定目前不支持复制的字段，以及必须需重新设置值的元信息
	deleteColumnIds := getCopyIssueDeleteColumnIds()
	for _, id := range deleteColumnIds {
		delete(newData, id)
	}

	newData[consts.TempFieldChildren] = []map[string]interface{}{}
	newData[consts.TempFieldChildrenCount] = 0
	return newData
}

// buildIssueChildDataByCopy 复制一条任务的子任务，来构造新任务的子任务
// parentData 需复制子任务的父任务
// childrenIds 需复制的子任务id
// allData 任务数据id->数据
// parentChildrenMapping 父任务id->子任务id mapping
// iterationMapping 迭代值映射mapping。传nil则保留原值(同一个项目中复制可以沿用迭代值，跨项目则需要映射)
// keepTableColumns 复制时保留的列，其余的列不会被复制。传nil则全量复制(同一张表中复制任务可以全量复制，应用模板时由于目标表格复制了原表的表头，也可以全量复制）
func buildIssueChildDataByCopy(parentData map[string]interface{},
	childrenIds []int64,
	allData map[int64]map[string]interface{},
	parentChildrenMapping map[int64][]int64,
	iterationMapping map[int64]int64,
	keepTableColumns map[string]struct{}) (int64, errs.SystemErrorInfo) {
	var count int64
	children := parentData[consts.TempFieldChildren].([]map[string]interface{})

	for _, childId := range childrenIds {
		data, ok := allData[childId]
		if !ok {
			continue
		}
		newData := buildIssueDataByCopy(data, iterationMapping, keepTableColumns)

		children = append(children, newData)
		count += 1

		if cIds, ok := parentChildrenMapping[childId]; ok {
			ct, err := buildIssueChildDataByCopy(newData, cIds, allData, parentChildrenMapping, iterationMapping, keepTableColumns)
			if err != nil {
				return 0, err
			}
			count += ct
		}
	}
	parentData[consts.TempFieldChildren] = children
	parentData[consts.TempFieldChildrenCount] = count

	return count, nil
}

func CopyIssue(orgId int64, userId int64, input *vo.LessCopyIssueReq) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	// 判断权限 - 在目标项目创建任务的权限
	errSys := domain.AuthProject(orgId, userId, input.ProjectId, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create)
	if errSys != nil {
		log.Errorf("[CopyIssue] AuthProject err: %v", errSys)
		return nil, nil, nil, errSys
	}

	// 获取所有子任务id + 子任务的子任务id
	allChildrenIds, errSys := domain.GetIssueAndChildrenIds(orgId, input.ChildrenIds)
	if errSys != nil {
		log.Errorf("[CopyIssue] GetIssueAndChildrenIds err: %v", errSys)
		return nil, nil, nil, errSys
	}
	// 加上复制的任务本身的id，即是所有需复制任务的id
	allIds := append(allChildrenIds, input.OldIssueId)

	// 判断权限 - 付费任务数
	errSys = domain.AuthPayTask(orgId, consts.FunctionTaskLimit, len(allIds))
	if errSys != nil {
		log.Errorf("[CopyIssue] AuthPayTask err: %v", errSys)
		return nil, nil, nil, errSys
	}

	// 获取所有原始任务数据
	allCopyData, errSys := domain.GetIssueInfosMapLcByIssueIds(orgId, userId, allIds)
	if errSys != nil {
		log.Errorf("[CopyIssue] GetIssueInfosMapLcByIssueIds err: %v", errSys)
		return nil, nil, nil, errs.IssueNotExist
	}

	// 选中的任务标题改成传入的值
	if input.Title != nil {
		for _, d := range allCopyData {
			issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
			if issueId == input.OldIssueId {
				d[consts.BasicFieldTitle] = *input.Title
			}
		}
	}

	return copyIssueBatchInner(orgId, userId, input.ProjectId, cast.ToInt64(input.TableId), allIds, allCopyData,
		input.ChooseField, &projectvo.TriggerBy{
			TriggerBy: consts.TriggerByCopy,
		}, true, input.IsStaticCopy, input.BeforeDataId, input.AfterDataId)
}

func copyIssueBatchInner(orgId int64, userId int64, newProjectId, newTableId int64, allIssueIds []int64,
	allCopyData []map[string]interface{}, chooseField []string, triggerBy *projectvo.TriggerBy,
	isSyncCreate, isStaticCopy bool, beforeDataId, afterDataId int64) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	// 获取project
	project, errSys := domain.GetProjectSimple(orgId, newProjectId)
	if errSys != nil {
		log.Errorf("[CopyIssueBatch] GetProjectSimple err: %v", errSys)
		return nil, nil, nil, errSys
	}

	// 所有数据按order降序排列
	sort.Slice(allCopyData, func(i, j int) bool {
		orderI := cast.ToInt64(allCopyData[i][consts.BasicFieldOrder])
		orderJ := cast.ToInt64(allCopyData[j][consts.BasicFieldOrder])
		return orderI > orderJ
	})

	// 生成issue ids
	issueIds, errSys := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssue, len(allCopyData))
	if errSys != nil {
		log.Errorf("[BatchCreateIssue] ApplyMultiplePrimaryIdRelaxed failed, count: %v, err: %v",
			len(allCopyData), errSys)
		return nil, nil, nil, errSys
	}

	// 生成issue codes
	preCode := consts.NoProjectPreCode
	if project.PreCode != "" {
		preCode = project.PreCode
	} else {
		preCode = fmt.Sprintf("$%d", newProjectId)
	}
	issueCodes, errSys := idfacade.ApplyMultipleIdRelaxed(orgId, preCode, "", int64(len(allCopyData)))
	if errSys != nil {
		log.Errorf("[BatchCreateIssue] ApplyMultiplePrimaryIdRelaxed failed, count: %v, err: %v",
			len(allCopyData), errSys)
		return nil, nil, nil, errSys
	}

	// 找出父子任务关系
	parentChildrenMapping := map[int64][]int64{}       // 父任务->子任务列表
	allData := map[int64]map[string]interface{}{}      // 任务id->任务数据
	allParentData := make([]map[string]interface{}, 0) // 父任务
	data := make([]map[string]interface{}, 0)          // 待创建的任务数据
	issueIdMapping := make(map[int64]int64, 0)         // oldIssueId->newIssueId
	isSameProject := false
	isSameTable := false
	var oldTableId int64
	for i, d := range allCopyData {
		issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
		parentId := cast.ToInt64(d[consts.BasicFieldParentId])
		projectId := cast.ToInt64(d[consts.BasicFieldProjectId])
		tableId := cast.ToInt64(d[consts.BasicFieldTableId])
		oldTableId = tableId

		allData[issueId] = d

		if i == 0 {
			isSameProject = projectId == newProjectId
			isSameTable = tableId == newTableId
		}

		// 赋予新issueId/code
		newIssueId := issueIds.Ids[i].Id
		d[consts.TempFieldIssueId] = newIssueId
		issueIdMapping[issueId] = newIssueId
		d[consts.TempFieldCode] = issueCodes.Ids[i].Code

		// 与父任务断开的子任务也视为父任务
		if has, _ := slice.Contain(allIssueIds, parentId); !has {
			// 父任务
			allParentData = append(allParentData, d)

			if !isSameTable {
				d[consts.BasicFieldParentId] = 0
				d[consts.BasicFieldPath] = "0,"
			}
		} else {
			// 子任务
			parentChildrenMapping[parentId] = append(parentChildrenMapping[parentId], issueId)
		}
	}

	// 动态复制
	// 替换关联/前后置/单向关联ID
	if !isStaticCopy && oldTableId > 0 {
		// 获取表头
		tableColumns, errSys := domain.GetTableColumnsMap(orgId, oldTableId, nil, true)
		if errSys != nil {
			log.Errorf("[CopyIssue] GetTableColumnsMap failed, org:%d table:%d, err: %v",
				orgId, oldTableId, errSys)
			return nil, nil, nil, errSys
		}

		replacingRelatingIds(oldTableId, tableColumns, allData, issueIdMapping, false, false)
	}

	// 要复制的列
	var keepTableColumns map[string]struct{}
	if !isSameTable {
		keepTableColumns = make(map[string]struct{})
		if len(chooseField) == 1 && chooseField[0] == "*" {
			keepTableColumns[consts.BasicFieldTitle] = struct{}{} // 标题必须复制
			keepTableColumns[consts.BasicFieldRemark] = struct{}{}
			keepTableColumns[consts.BasicFieldRemarkDetail] = struct{}{}
			keepTableColumns[consts.BasicFieldAuditorIds] = struct{}{}
			keepTableColumns[consts.BasicFieldOwnerId] = struct{}{}
			keepTableColumns[consts.BasicFieldFollowerIds] = struct{}{}
			keepTableColumns[consts.BasicFieldPlanStartTime] = struct{}{}
			keepTableColumns[consts.BasicFieldPlanEndTime] = struct{}{}
			keepTableColumns[consts.BasicFieldRelating] = struct{}{}
			keepTableColumns[consts.BasicFieldBaRelating] = struct{}{}
			if isSameProject {
				keepTableColumns[consts.BasicFieldIterationId] = struct{}{}
			}
		} else {
			keepTableColumns[consts.BasicFieldTitle] = struct{}{} // 标题必须复制
			for _, field := range chooseField {
				switch field {
				//case consts.BasicFieldIssueStatus:
				//	// 任务状态只支持本表复制
				//	if isSameTable {
				//		copyIssueStatus = true
				//		keepTableColumns[consts.BasicFieldIssueStatus] = struct{}{}
				//		keepTableColumns[consts.BasicFieldIssueStatusType] = struct{}{}
				//	}
				//case consts.BasicFieldPriority:
				//	// 优先级只支持本表复制
				//	if isSameTable {
				//		keepTableColumns[field] = struct{}{}
				//	}
				case consts.BasicFieldIterationId:
					// 迭代只支持本项目复制
					if isSameProject {
						keepTableColumns[field] = struct{}{}
					}
				case consts.BasicFieldRemark:
					keepTableColumns[field] = struct{}{}
					keepTableColumns[consts.BasicFieldRemarkDetail] = struct{}{}
				case consts.BasicFieldAuditorIds,
					consts.BasicFieldOwnerId,
					consts.BasicFieldFollowerIds,
					consts.BasicFieldPlanStartTime,
					consts.BasicFieldPlanEndTime,
					consts.BasicFieldRelating,
					consts.BasicFieldBaRelating:
					keepTableColumns[field] = struct{}{}
				case consts.BasicFieldWorkHour:
					// TODO: 工时 - 暂不能支持无码化复制
				}
			}
		}
	}

	var totalCount int64 // 任务总数（包括子任务）

	// 处理父任务
	for _, d := range allParentData {
		issueId := cast.ToInt64(d[consts.BasicFieldIssueId])
		parentData := buildIssueDataByCopy(d, nil, keepTableColumns)

		if childrenIds, ok := parentChildrenMapping[issueId]; ok {
			// 处理子任务（递归）
			count, err := buildIssueChildDataByCopy(parentData, childrenIds, allData, parentChildrenMapping, nil, keepTableColumns)
			if err != nil {
				log.Errorf("[CopyIssue] buildIssueChildDataByCopy err: %v", err)
				return nil, nil, nil, err
			}
			totalCount += count
		}

		totalCount += 1
		data = append(data, parentData)
	}

	//// TODO: 工时 - 暂不能支持无码化复制
	//if util.FieldInUpdate(input.ChooseField, consts.BasicFieldWorkHour) {
	//	// 查出原来工时记录，传给createIssue创建
	//	issueWorkHours, errSys := domain.GetCopyIssueWorkHours(orgId, input.ProjectID, input.OldIssueID)
	//	if errSys != nil {
	//		log.Errorf("[copyIssue]GetCopyIssueWorkHours err:%v", errSys)
	//		return nil, errSys
	//	}
	//	createIssueReq.LessCreateIssueReq[consts.CopyIssueWorkHour] = issueWorkHours
	//	// 无码更新值
	//	createIssueReq.LessCreateIssueReq[consts.BasicFieldWorkHour] = data[consts.BasicFieldWorkHour]
	//}

	if !isSyncCreate {
		// 异步批量创建
		req := &projectvo.BatchCreateIssueReqVo{
			OrgId:         orgId,
			UserId:        userId,
			AppId:         project.AppId,
			ProjectId:     project.Id,
			TableId:       newTableId,
			Data:          data,
			BeforeDataId:  beforeDataId,
			AfterDataId:   afterDataId,
			IsIdGenerated: true,
		}
		AsyncBatchCreateIssue(req, true, triggerBy, "")
		return make([]map[string]interface{}, totalCount, totalCount), nil, nil, nil
	} else {
		// 同步批量创建
		req := &projectvo.BatchCreateIssueReqVo{
			OrgId:         orgId,
			UserId:        userId,
			AppId:         project.AppId,
			ProjectId:     project.Id,
			TableId:       newTableId,
			Data:          data,
			BeforeDataId:  beforeDataId,
			AfterDataId:   afterDataId,
			IsIdGenerated: true,
		}
		return SyncBatchCreateIssue(req, true, triggerBy)
	}
}

func CopyIssueBatch(orgId int64, userId int64, input *vo.CopyIssueBatchReq, triggerBy *projectvo.TriggerBy) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	if len(input.OldIssueIds) == 0 {
		return nil, nil, nil, errs.IssueIdsNotBeenChosen
	}

	// 判断权限 - 在目标项目创建任务的权限
	errSys := domain.AuthProject(orgId, userId, input.ProjectID, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create)
	if errSys != nil {
		log.Errorf("[CopyIssueBatch] AuthProject err: %v", errSys)
		return nil, nil, nil, errSys
	}

	// 判断付费限制 - 付费任务数
	errSys = domain.AuthPayTask(orgId, consts.FunctionTaskLimit, len(input.OldIssueIds))
	if errSys != nil {
		log.Errorf("[CopyIssueBatch] AuthPayTask err: %v", errSys)
		return nil, nil, nil, errSys
	}

	// 获取所有原始任务数据
	allCopyData, err := domain.GetIssueInfosMapLcByIssueIds(orgId, userId, input.OldIssueIds)
	if err != nil {
		log.Errorf("[CopyIssueBatch] GetIssueInfosMapLcByIssueIds err: %v", err)
		return nil, nil, nil, errs.IssueNotExist
	}

	return copyIssueBatchInner(orgId, userId, input.ProjectID, cast.ToInt64(input.TableID), input.OldIssueIds, allCopyData,
		input.ChooseField, triggerBy, false, input.IsStaticCopy, 0, 0)
}

func CopyIssueBatchWithData(orgId int64, userId int64, newProjectId, newTableId int64, allIssueIds []int64,
	allCopyData []map[string]interface{}, triggerBy *projectvo.TriggerBy, isSyncCreate, isStaticCopy, isInnerSuper bool) (
	[]map[string]interface{}, map[string]*uservo.MemberDept, map[string]map[string]interface{}, errs.SystemErrorInfo) {
	if len(allIssueIds) == 0 {
		return nil, nil, nil, errs.IssueIdsNotBeenChosen
	}

	// 判断权限 - 在目标项目创建任务的权限
	if !isInnerSuper {
		errSys := domain.AuthProject(orgId, userId, newProjectId, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Create)
		if errSys != nil {
			log.Errorf("[CopyIssueBatchWithData] AuthProject err: %v", errSys)
			return nil, nil, nil, errSys
		}
	}

	// 判断付费限制 - 付费任务数
	errSys := domain.AuthPayTask(orgId, consts.FunctionTaskLimit, len(allIssueIds))
	if errSys != nil {
		log.Errorf("[CopyIssueBatch] AuthPayTask err: %v", errSys)
		return nil, nil, nil, errSys
	}

	return copyIssueBatchInner(orgId, userId, newProjectId, newTableId, allIssueIds, allCopyData,
		[]string{"*"}, triggerBy, isSyncCreate, isStaticCopy, 0, 0)
}

//func ArchiveIssue(orgId, currentUserId int64, sourceChannel string, issueIds []int64) (*vo.ArchiveIssueBatchResp, errs.SystemErrorInfo) {
//	//查找目标任务以及所有的子任务
//	filingType := consts.ProjectIsNotFiling
//	issueAndChildren, err := domain.GetIssueAndChildren(orgId, issueIds, &filingType)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	remainChildrenIssueIds := []int64{}
//	//传进来的父任务id
//	allParentIds := []int64{}
//	for _, child := range issueAndChildren {
//		if ok, _ := slice.Contain(issueIds, child.Id); !ok {
//			//不存在的必然是子任务
//			for _, id := range issueIds {
//				if strings.Contains(child.Path, fmt.Sprintf(",%d,", id)) {
//					remainChildrenIssueIds = append(remainChildrenIssueIds, id)
//				}
//			}
//		}
//		if child.ParentId == 0 {
//			if ok, _ := slice.Contain(allParentIds, child.Id); !ok {
//				allParentIds = append(allParentIds, child.Id)
//			}
//		}
//	}
//	remainChildrenIssueIds = slice.SliceUniqueInt64(remainChildrenIssueIds)
//
//	remainChildrenInfo := []bo.IssueBo{}
//	allNeedAuth := []bo.IssueBo{}
//	for _, child := range issueAndChildren {
//		if ok, _ := slice.Contain(remainChildrenIssueIds, child.Id); ok {
//			remainChildrenInfo = append(remainChildrenInfo, child)
//		} else {
//			if child.ParentId == 0 {
//				allNeedAuth = append(allNeedAuth, child)
//			} else if ok1, _ := slice.Contain(remainChildrenIssueIds, child.ParentId); !ok1 {
//				//找出需要判断权限的(只要父任务不在所选范围内，就表示他是父任务)
//				if ok1, _ := slice.Contain(issueIds, child.ParentId); !ok1 {
//					allNeedAuth = append(allNeedAuth, child)
//				}
//			}
//		}
//	}
//
//	notAuthPassIssues := []bo.IssueBo{}
//	noAuthIds := []int64{}
//	trulyIssueBos := []bo.IssueBo{}
//
//	for _, issueBo := range allNeedAuth {
//		// 目前 consts.OperationProIssue4ModifyStatus 权限项已被去除，对应的鉴权由字段权限控制。
//		//权限判断
//		err := domain.AuthIssue(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify, consts.BasicFieldIssueStatus)
//		if err != nil {
//			log.Error(err)
//			notAuthPassIssues = append(notAuthPassIssues, issueBo)
//			noAuthIds = append(noAuthIds, issueBo.Id)
//		} else {
//			trulyIssueBos = append(trulyIssueBos, issueBo)
//		}
//
//		trulyIssueBos = append(trulyIssueBos, issueBo)
//	}
//
//	//设置归档
//	childIds, err3 := domain.ArchiveIssue(orgId, trulyIssueBos, currentUserId, sourceChannel)
//	if err3 != nil {
//		log.Error(err3)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err3)
//	}
//
//	asyn.Execute(func() {
//		for _, issueBo := range trulyIssueBos {
//			children := []int64{}
//			if _, ok := childIds[issueBo.Id]; ok {
//				children = childIds[issueBo.Id]
//			}
//			PushArchiveIssueNotice(issueBo, children, currentUserId, consts.MQTTDataRefreshActionArchive)
//		}
//	})
//
//	resBo := bo.ArchiveIssueBo{
//		NoAuthIssues: notAuthPassIssues,
//	}
//
//	res := &vo.ArchiveIssueBatchResp{}
//	for _, child := range issueAndChildren {
//		if ok, _ := slice.Contain(noAuthIds, child.Id); ok {
//			continue
//		}
//		if ok, _ := slice.Contain(noAuthIds, child.ParentId); ok {
//			continue
//		}
//		resBo.SuccessIssues = append(resBo.SuccessIssues, child)
//	}
//
//	_ = copyer.Copy(resBo, res)
//
//	return res, nil
//}

//func CancelArchiveIssue(orgId, currentUserId int64, sourceChannel string, issueIds []int64, projectId int64) (*vo.CancelArchiveIssueBatchResp, errs.SystemErrorInfo) {
//	//获取项目信息
//	projectInfo, err := domain.GetProject(orgId, projectId)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	//父任务可以单独取消归档，子任务不能单独在未归档行列中
//	issueInfos, err := domain.GetIssueInfoList(issueIds)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	parentIds := []int64{}
//	for _, info := range issueInfos {
//		if info.ParentId != 0 {
//			for _, s := range strings.Split(info.Path[0:len(info.Path)-1], ",") {
//				id, err := strconv.ParseInt(s, 10, 64)
//				if err != nil {
//					log.Error(err)
//					return nil, errs.SystemError
//				}
//				if id != 0 {
//					parentIds = append(parentIds, id)
//				}
//			}
//		}
//	}
//
//	parentIds = slice.SliceUniqueInt64(parentIds)
//	needConfirmIds, _ := util.GetDifMemberIds(parentIds, issueIds)
//
//	needConfirmInfos, err1 := domain.GetIssueInfoList(needConfirmIds)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, err1
//	}
//
//	//处于归档且此次没有取消归档的父任务
//	filingParentIds := []int64{}
//	//所有需要确认的且没有被删除的父任务
//	allExistConfirmParentIds := []int64{}
//	for _, info := range needConfirmInfos {
//		allExistConfirmParentIds = append(allExistConfirmParentIds, info.Id)
//		if info.IsFiling == consts.ProjectIsFiling {
//			filingParentIds = append(filingParentIds, info.Id)
//		}
//	}
//
//	deleteParentIds, _ := util.GetDifMemberIds(needConfirmIds, allExistConfirmParentIds)
//
//	var successIds []int64
//	var successIssueInfos, parentDeleteInfo, parentFilingInfo []bo.IssueBo
//	for _, info := range issueInfos {
//		if info.ParentId == 0 {
//			successIds = append(successIds, info.Id)
//			successIssueInfos = append(successIssueInfos, info)
//		} else {
//			for _, s := range strings.Split(info.Path[0:len(info.Path)-1], ",") {
//				id, err := strconv.ParseInt(s, 10, 64)
//				if err != nil {
//					log.Error(err)
//					return nil, errs.SystemError
//				}
//				if id == 0 {
//					continue
//				}
//				if ok, _ := slice.Contain(deleteParentIds, id); ok {
//					parentDeleteInfo = append(parentDeleteInfo, info)
//					break
//				} else if ok, _ := slice.Contain(filingParentIds, id); ok {
//					parentFilingInfo = append(parentFilingInfo, info)
//					break
//				} else {
//					successIds = append(successIds, id)
//					successIssueInfos = append(successIssueInfos, info)
//					break
//				}
//			}
//
//		}
//	}
//	var notAuthPassIssues, trulyIssueBos []bo.IssueBo
//	var trulyIds []int64
//	for _, issueBo := range successIssueInfos {
//		// 目前 consts.OperationProIssue4ModifyStatus 权限项已被去除，因此这里先注释一下。（该接口调用没有被用上）
//		//权限判断
//		//err := domain.AuthIssue(orgId, currentUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4ModifyStatus)
//		//if err != nil {
//		//	log.Error(err)
//		//	notAuthPassIssues = append(notAuthPassIssues, issueBo)
//		//} else {
//		//	trulyIssueBos = append(trulyIssueBos, issueBo)
//		//	trulyIds = append(trulyIds, issueBo.Id)
//		//}
//
//		trulyIssueBos = append(trulyIssueBos, issueBo)
//		trulyIds = append(trulyIds, issueBo.Id)
//	}
//
//	if len(trulyIds) > 0 {
//		_, updErr := mysql.UpdateSmartWithCond(consts.TableIssue, db.Cond{
//			consts.TcOrgId: orgId,
//			consts.TcId:    db.In(trulyIds),
//		}, mysql.Upd{
//			consts.TcIsFiling: consts.ProjectIsNotFiling,
//			consts.TcUpdator:  currentUserId,
//		})
//		if updErr != nil {
//			log.Error(updErr)
//			return nil, errs.MysqlOperateError
//		}
//
//		//判断项目是否是敏捷项目（如果是敏捷项目要判断任务状态是否依旧存在，通用项目则判断任务栏是否存在）
//		if projectInfo.ProjectTypeId != consts.ProjectTypeNormalId {
//			//敏捷项目判断任务状态
//			allStatus := consts.IterationStatusList
//			var allStatusIds, needUpdateStatusIds []int64
//			for _, status := range allStatus {
//				allStatusIds = append(allStatusIds, status.ID)
//			}
//			defaultStatusId := allStatusIds[0]
//			for _, issueBo := range trulyIssueBos {
//				if ok, _ := slice.Contain(allStatusIds, issueBo.Status); !ok {
//					needUpdateStatusIds = append(needUpdateStatusIds, issueBo.Id)
//				}
//			}
//			if len(needUpdateStatusIds) > 0 {
//				_, err := mysql.UpdateSmartWithCond(consts.TableIssue, db.Cond{
//					consts.TcOrgId: orgId,
//					consts.TcId:    db.In(needUpdateStatusIds),
//				}, mysql.Upd{
//					consts.TcStatus: defaultStatusId,
//				})
//				if err != nil {
//					log.Error(err)
//					return nil, errs.MysqlOperateError
//				}
//			}
//
//		} else {
//			//通用项目判断任务栏
//			allProjectObjectTypes, err := ProjectObjectTypesWithProject(orgId, projectId)
//			if err != nil {
//				log.Error(err)
//				return nil, err
//			}
//			var allProjectObjectTypeIds, needUpdateTypeIds []int64
//			for _, objectType := range allProjectObjectTypes.List {
//				allProjectObjectTypeIds = append(allProjectObjectTypeIds, objectType.ID)
//			}
//			defaultTypeId := allProjectObjectTypeIds[0]
//			for _, issueBo := range trulyIssueBos {
//				if ok, _ := slice.Contain(allProjectObjectTypeIds, issueBo.ProjectObjectTypeId); !ok {
//					needUpdateTypeIds = append(needUpdateTypeIds, issueBo.Id)
//				}
//			}
//			if len(needUpdateTypeIds) > 0 {
//				_, err := mysql.UpdateSmartWithCond(consts.TableIssue, db.Cond{
//					consts.TcOrgId: orgId,
//					consts.TcId:    db.In(needUpdateTypeIds),
//				}, mysql.Upd{
//					consts.TcProjectObjectTypeId: defaultTypeId,
//				})
//				if err != nil {
//					log.Error(err)
//					return nil, errs.MysqlOperateError
//				}
//			}
//		}
//	}
//
//	asyn.Execute(func() {
//		for _, issueBo := range trulyIssueBos {
//			issueTrendsBo := bo.IssueTrendsBo{
//				PushType:      consts.PushTypeCancelArchiveIssue,
//				OrgId:         orgId,
//				UserId:    currentUserId,
//				IssueId:       issueBo.Id,
//				ParentIssueId: issueBo.ParentId,
//				ProjectId:     issueBo.ProjectId,
//				PriorityId:    issueBo.PriorityId,
//				ParentId:      issueBo.ParentId,
//
//				IssueTitle:    issueBo.Title,
//				IssueStatusId: issueBo.Status,
//				SourceChannel: sourceChannel,
//				TableId:       issueBo.TableId,
//			}
//
//			asyn.Execute(func() {
//				domain.PushIssueTrends(issueTrendsBo)
//			})
//			asyn.Execute(func() {
//				domain.PushIssueThirdPlatformNotice(issueTrendsBo)
//			})
//		}
//
//		for _, issueBo := range trulyIssueBos {
//			PushArchiveIssueNotice(issueBo, nil, currentUserId, consts.MQTTDataRefreshActionCancelArchive)
//		}
//	})
//
//	res := &bo.CancelArchiveIssueBo{
//		ParentDeletedIssues: parentDeleteInfo,
//		ParentFilingIssues:  parentFilingInfo,
//		NoAuthIssues:        notAuthPassIssues,
//		SuccessIssues:       trulyIssueBos,
//	}
//
//	resPo := &vo.CancelArchiveIssueBatchResp{}
//	_ = copyer.Copy(res, resPo)
//
//	return resPo, nil
//}

// ViewAuditIssue 查看审核任务
func ViewAuditIssue(orgId, userId, issueId int64) errs.SystemErrorInfo {
	issueInfos, sysErr := domain.GetIssueInfosLc(orgId, userId, []int64{issueId})
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}
	if len(issueInfos) < 1 {
		log.Errorf("[ViewAuditIssue] not found issue issueId: %v", issueId)
		return errs.IssueNotExist
	}
	issueInfo := issueInfos[0]

	// 判断当前操作人是否是任务确认人
	if ok, _ := slice.Contain(issueInfo.AuditorIdsI64, userId); !ok {
		return nil
	}

	summaryAppId, sysErr := domain.GetOrgSummaryAppId(orgId)
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}

	auditStatus, ok := issueInfo.AuditStatusDetail[cast.ToString(userId)]
	if !ok || auditStatus == consts.AuditStatusNotView {
		// 更新无码
		sysErr = domain.UpdateLcIssueAuditStatusDetailByUser(orgId, summaryAppId, userId, issueId, consts.AuditStatusView)
		if sysErr != nil {
			log.Error(sysErr)
			return sysErr
		}
	}

	return nil
}

// WithdrawIssue 撤回审核任务
func WithdrawIssue(orgId, userId int64, issueId int64) errs.SystemErrorInfo {
	issueInfos, sysErr := domain.GetIssueInfosLc(orgId, userId, []int64{issueId})
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}
	if len(issueInfos) < 1 {
		log.Errorf("[WithdrawIssue] not found issue issueId: %v", issueId)
		return errs.IssueNotExist
	}
	issueInfo := issueInfos[0]

	// 判断当前操作人是否是任务负责人
	if ok, _ := slice.Contain(issueInfo.OwnerIdI64, userId); !ok {
		return errs.OnlyOwnerCanWithdrawIssue
	}

	// 未完成的任务无需撤回
	if issueInfo.IssueStatusType != consts.StatusTypeComplete {
		return errs.NotFinishIssue
	}

	// 已经通过审核了
	if issueInfo.AuditStatus == consts.AuditStatusPass {
		return errs.IssueIsAuditPass
	}

	allStatusMap, sysErr := domain.GetIssueAllStatus(orgId, []int64{issueInfo.ProjectId}, []int64{issueInfo.TableId})
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}
	issueStatusId := issueInfo.Status
	if allStatus, ok := allStatusMap[issueInfo.TableId]; ok {
		for _, statusInfo := range allStatus {
			if statusInfo.Type == consts.StatusTypeNotStart {
				issueStatusId = statusInfo.ID
			}
		}
	} else {
		return errs.BuildSystemErrorInfo(errs.IssueStatusNotExist)
	}

	//transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
	//// 更新无码
	//sysErr = domain.UpdateLcIssueAuditStatusDetailByUser(orgId, appId, userId, issueId, consts.AuditStatusNotView)
	//if sysErr != nil {
	//	log.Error(sysErr)
	//	return sysErr
	//}

	//_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableIssueRelation, db.Cond{
	//	consts.TcRelationType: consts.IssueRelationTypeAuditor,
	//	consts.TcOrgId:        orgId,
	//	consts.TcIssueId:      issueId,
	//}, mysql.Upd{
	//	consts.TcStatus:  consts.AuditStatusNotView,
	//	consts.TcUpdator: userId,
	//})
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}

	// 清空审批状态
	d := make(map[string]interface{})
	d[consts.BasicFieldId] = issueId
	d[consts.BasicFieldAuditStatus] = consts.AuditStatusNotView
	d[consts.BasicFieldAuditStatusDetail] = map[string]int{}
	d[consts.BasicFieldIssueStatus] = issueStatusId
	d[consts.BasicFieldIssueStatusType] = consts.StatusTypeNotStart //撤回到未完成状态

	// 更新无码
	pushType := consts.PushTypeWithdrawIssue
	sysErr = BatchUpdateIssue(&projectvo.BatchUpdateIssueReqInnerVo{
		OrgId:         orgId,
		UserId:        userId,
		AppId:         issueInfo.AppId,
		ProjectId:     issueInfo.ProjectId,
		TableId:       issueInfo.TableId,
		Data:          []map[string]interface{}{d},
		TrendPushType: &pushType,
	}, true, &projectvo.TriggerBy{
		TriggerBy: consts.TriggerByNormal,
	})
	if sysErr != nil {
		log.Error(sysErr)
		return sysErr
	}

	//return nil
	//})
	//if transErr != nil {
	//	log.Error(transErr)
	//	return errs.MysqlOperateError
	//}

	////动态
	//asyn.Execute(func() {
	//	issueTrendsBo := bo.IssueTrendsBo{
	//		PushType:      consts.PushTypeWithdrawIssue,
	//		OrgId:         orgId,
	//		UserId:    userId,
	//		IssueId:       issueId,
	//		ParentIssueId: issueInfo.ParentId,
	//		ProjectId:     issueInfo.ProjectId,
	//		PriorityId:    issueInfo.PriorityId,
	//		ParentId:      issueInfo.ParentId,
	//
	//		IssueTitle:    issueInfo.Title,
	//		IssueStatusId: issueStatusId,
	//		TableId:       issueInfo.TableId,
	//	}
	//
	//	asyn.Execute(func() {
	//		domain.PushIssueTrends(issueTrendsBo)
	//	})
	//
	//	asyn.Execute(func() {
	//		PushModifyIssueNotice(issueInfo.OrgId, issueInfo.ProjectId, issueInfo.Id, userId)
	//	})
	//})

	return nil
}

func ReportMoveEvents(orgId, userId int64, oldIssues []*bo.IssueBo, issueIds []int64, fromTableId, targetTableId int64) {
	var errSys errs.SystemErrorInfo

	// 获取表头
	var oldTableColumns, newTableColumns map[string]*projectvo.TableColumnData
	oldTableColumns, errSys = domain.GetTableColumnsMap(orgId, fromTableId, nil, true)
	if errSys != nil {
		log.Errorf("[ReportMoves] 获取表头失败 org:%d table:%d, err: %v", orgId, fromTableId, errSys)
		return
	}
	if targetTableId != fromTableId {
		newTableColumns, errSys = domain.GetTableColumnsMap(orgId, targetTableId, nil, true)
		if errSys != nil {
			log.Errorf("[ReportMoves] 获取表头失败 org:%d table:%d, err: %v", orgId, targetTableId, errSys)
			return
		}
	} else {
		newTableColumns = oldTableColumns
	}

	var userIds, deptIds []int64

	// issueIds 里面的任务可能移动也可能没移动，但这里统一都拿出来，然后再判断是否有变化
	newDatas, errSys := domain.GetIssueInfosMapLcByIssueIds(orgId, userId, issueIds)
	if errSys != nil {
		log.Errorf("[ReportMoves] GetIssueInfosMapLcByIssueIds org:%d, issues:%q err: %v", orgId, issueIds, errSys)
		return
	}
	oldIssueMap := make(map[int64]*bo.IssueBo)
	for _, issue := range oldIssues {
		oldIssueMap[issue.Id] = issue
	}
	for _, newData := range newDatas {
		uIds, dIds := domain.CollectUserDeptIds(newData, newTableColumns)
		userIds = append(userIds, uIds...)
		deptIds = append(deptIds, dIds...)
	}
	userDepts := domain.AssembleUserDeptsByIds(orgId, userIds, deptIds)

	for _, newData := range newDatas {
		issueId := cast.ToInt64(newData[consts.BasicFieldIssueId])
		if oldIssue, ok := oldIssueMap[issueId]; ok {
			domain.ReportMoveEvent(userId, oldIssue, newData, userDepts)
		}
	}
}

// 变更父任务
func ChangeParentIssue(orgId, userId int64, input vo.ChangeParentIssueReq) (*vo.Void, errs.SystemErrorInfo) {
	//获取原本的Issue
	issueBo, errSys := domain.GetIssueInfoLc(orgId, userId, input.IssueID)
	if errSys != nil {
		log.Error(errSys)
		return nil, errs.BuildSystemErrorInfo(errs.IssueNotExist, errSys)
	}

	errSys = domain.AuthIssue(orgId, userId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify)
	if errSys != nil {
		log.Error(errSys)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, errSys)
	}

	if issueBo.ParentId == input.ParentID {
		return &vo.Void{ID: input.IssueID}, nil
	}

	//获取目标父任务
	parentIssueBo, errSys := domain.GetIssueInfoLc(orgId, userId, input.ParentID)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}

	//先转化为目标任务的子任务
	allIssueIds, oldIssues, errSys := domain.ChangeIssueChildren(userId, issueBo, input.ParentID, parentIssueBo.Path)
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}

	targetProjectId := parentIssueBo.ProjectId
	targetTableId := parentIssueBo.TableId
	if issueBo.ProjectId != targetProjectId || issueBo.TableId != targetTableId {
		errSys = domain.UpdateIssueProjectTable(orgId, userId, issueBo, targetTableId, targetProjectId, true, false)
		if errSys != nil {
			log.Error(errSys)
			return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, errSys)
		}
	}

	// 上报事件
	asyn.Execute(func() {
		ReportMoveEvents(orgId, userId, oldIssues, allIssueIds, issueBo.TableId, targetTableId)
	})

	return &vo.Void{
		ID: input.IssueID,
	}, nil
}

// 转化为父任务
func ConvertIssueToParent(orgId, userId int64, input vo.ConvertIssueToParentReq) (*vo.Void, errs.SystemErrorInfo) {
	//获取原本的Issue
	issueBo, errSys := domain.GetIssueInfoLc(orgId, userId, input.ID)
	if errSys != nil {
		log.Error(errSys)
		return nil, errs.BuildSystemErrorInfo(errs.IssueNotExist, errSys)
	}

	errSys = domain.AuthIssue(orgId, userId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify)
	if errSys != nil {
		log.Error(errSys)
		return nil, errs.BuildSystemErrorInfo(errs.Unauthorized, errSys)
	}

	//先转化为父任务
	allIssueIds, oldIssues, errSys := domain.ChangeIssueChildren(userId, issueBo, 0, "")
	if errSys != nil {
		log.Error(errSys)
		return nil, errSys
	}

	targetProjectId := input.ProjectID
	targetTableId := cast.ToInt64(input.TableID)
	if issueBo.ProjectId != targetProjectId || issueBo.TableId != targetTableId {
		errSys = domain.UpdateIssueProjectTable(orgId, userId, issueBo, targetTableId, targetProjectId, true, false)
		if errSys != nil {
			log.Error(errSys)
			return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, errSys)
		}
	}

	// 上报事件
	asyn.Execute(func() {
		ReportMoveEvents(orgId, userId, oldIssues, allIssueIds, issueBo.TableId, targetTableId)
	})

	return &vo.Void{
		ID: input.ID,
	}, nil
}

func GetIssueRowList(orgId, userId int64, input *tablePb.ListRawRequest) ([]map[string]interface{}, errs.SystemErrorInfo) {
	rawListResp, err := tablefacade.ListRawExpand(orgId, userId, input)
	if err != nil {
		log.Errorf("[GetIssueRowList] ListRawExpand err:%v, orgId:%v, userId:%v, ListRawRequest:%v",
			err, orgId, userId, input)
		return nil, err
	}
	return rawListResp.Data, nil
}
