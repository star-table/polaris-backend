package domain

import (
	"fmt"
	"strings"
	"sync"

	"github.com/star-table/polaris-backend/common/extra/lc_helper"

	tablev1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo/appauth"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/permissionfacade"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/spf13/cast"
)

//func SelectList(cond db.Cond, union []*db.Union, page int, size int, order interface{}, needCount bool) (*[]bo.IssueBo, int64, errs.SystemErrorInfo) {
//	issues := &[]*po.PpmPriIssue{}
//	var total uint64
//
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		log.Error(err)
//		return nil, 0, errs.MysqlOperateError
//	}
//	//conn.SetLogging(true)
//	//defer conn.SetLogging(false)
//
//	//判断排序（如果要按照code排序要另外处理）
//	if order == db.Raw("code asc") || order == db.Raw("code desc") {
//		//conn.SetLogging(true)
//		mid := conn.Select(db.Raw("*, SUBSTRING_INDEX(`code`,'-',1) as code_name, SUBSTRING_INDEX(`code`,'-',-1) + 0 as code_index")).From(consts.TableIssue).Where(cond)
//		if union != nil && len(union) > 0 {
//			for _, d := range union {
//				mid = mid.And(d)
//			}
//		}
//
//		if order == db.Raw("code asc") {
//			mid = mid.OrderBy(db.Raw("code_name asc, code_index asc"))
//		} else {
//			mid = mid.OrderBy(db.Raw("code_name desc, code_index desc"))
//		}
//		if size > 0 && page > 0 {
//			mid1 := mid.Paginate(uint(size)).Page(uint(page))
//
//			if needCount {
//				nowTotal, err := mid1.TotalEntries()
//				if err != nil {
//					return nil, 0, errs.MysqlOperateError
//				}
//				total = nowTotal
//			}
//
//			err = mid1.All(issues)
//			if err != nil {
//				return nil, 0, errs.MysqlOperateError
//			}
//		} else {
//			err := mid.All(issues)
//			if err != nil {
//				return nil, 0, errs.MysqlOperateError
//			}
//			if needCount {
//				total = uint64(len(*issues))
//			}
//		}
//	} else {
//		//conn.SetLogging(true)
//		mid := conn.Collection(consts.TableIssue).Find(cond)
//		if union != nil && len(union) > 0 {
//			for _, d := range union {
//				mid = mid.And(d)
//			}
//		}
//
//		if size > 0 && page > 0 {
//			mid = mid.Page(uint(page)).Paginate(uint(size))
//		}
//		if order != nil && order != "" {
//			mid = mid.OrderBy(order)
//		}
//		if needCount {
//			nowTotal, err := mid.TotalEntries()
//			if err != nil {
//				log.Error(err)
//				return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//			}
//			total = nowTotal
//		}
//
//		err = mid.All(issues)
//		if err != nil {
//			log.Error(err)
//			return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//		}
//	}
//
//	issueBos := &[]bo.IssueBo{}
//	err2 := copyer.Copy(*issues, issueBos)
//	if err2 != nil {
//		log.Error(err2)
//		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err2)
//	}
//
//	return issueBos, int64(total), nil
//}

//func SelectIssueRemindInfoList(issueIdsCondBo bo.SelectIssueIdsCondBo, page int, size int) ([]bo.IssueRemindInfoBo, int64, errs.SystemErrorInfo) {
//	issuePrefix := "i."
//	projectPrefix := "p."
//	//规则，未完成
//	cond := db.Cond{
//		issuePrefix + consts.TcIsDelete: consts.AppIsNoDelete,
//		//通过end_time来筛选未完成的项目
//		issuePrefix + consts.TcEndTime:    db.Lt(consts.BlankElasticityTime),
//		projectPrefix + consts.TcIsDelete: consts.AppIsNoDelete,
//		projectPrefix + consts.TcIsFiling: consts.AppIsNotFilling,
//	}
//
//	//条件组装
//	if issueIdsCondBo.AfterPlanEndTime != nil && issueIdsCondBo.BeforePlanEndTime != nil {
//		cond[issuePrefix+consts.TcPlanEndTime] = db.Between(*issueIdsCondBo.BeforePlanEndTime, *issueIdsCondBo.AfterPlanEndTime)
//	}
//	if issueIdsCondBo.AfterPlanStartTime != nil && issueIdsCondBo.BeforePlanStartTime != nil {
//		cond[issuePrefix+consts.TcPlanStartTime] = db.Between(*issueIdsCondBo.BeforePlanStartTime, *issueIdsCondBo.AfterPlanStartTime)
//	}
//
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//
//	mid := conn.Select(
//		issuePrefix+consts.TcId,
//		issuePrefix+consts.TcPlanEndTime,
//		issuePrefix+consts.TcPlanStartTime,
//		issuePrefix+consts.TcOwner,
//		issuePrefix+consts.TcOrgId,
//		issuePrefix+consts.TcProjectId,
//		issuePrefix+consts.TcTitle,
//		issuePrefix+consts.TcParentId,
//	).From(consts.TableIssue + " i").LeftJoin(consts.TableProject + " p").On("i.project_id = p.id").Where(cond).OrderBy("i.id asc").Paginate(uint(size)).Page(uint(page))
//
//	//查询总数
//	total, err := mid.TotalEntries()
//	if err != nil {
//		log.Error(err)
//		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//
//	issueIdBos := &[]bo.IssueRemindInfoBo{}
//	//总数大于0的话才需要去取数据
//	if total > 0 {
//		err = mid.All(issueIdBos)
//		if err != nil {
//			log.Error(err)
//			return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//		}
//	}
//	return *issueIdBos, int64(total), nil
//}

//// AllIssueForProject 导出任务查询
//// 任务状态改造后，敏捷项目的 projectObjectTypeId 表示的实际是 tableId
//func AllIssueForProject(orgId, projectId int64, isParent bool, tableId int64, iterationId int64,
//	limitInfo projectvo.IssueViewLimitInfoForExportQuery, customColumns []*projectvo.TableColumnData,
//	infoMapObj *projectvo.IssueInfoMapForExportIssue,
//) ([]bo.IssueAndDetailInfoBo, errs.SystemErrorInfo) {
//	// func body start
//	queryLimit := limitInfo.Size
//	var businessErr errs.SystemErrorInfo
//	allowedIssueIds := limitInfo.AllowIssueIds
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//
//	resBo := &[]bo.IssueAndDetailInfoBo{}
//	cond := db.Cond{
//		"i." + consts.TcIsDelete:  consts.AppIsNoDelete,
//		"i." + consts.TcId:        db.Raw("d." + consts.TcIssueId),
//		"i." + consts.TcOrgId:     orgId,
//		"i." + consts.TcProjectId: projectId,
//		"i." + consts.TcTableId:   tableId,
//	}
//	if limitInfo.IsLimited && allowedIssueIds != nil && len(allowedIssueIds) > 0 {
//		// 这里的条件和上方的 `"i." + consts.TcId` key 有重复，因此下面这里加上空格，防止覆盖同 key 条件。
//		cond["  i."+consts.TcId] = db.In(allowedIssueIds)
//	}
//	//if projectObjectTypeId != 0 {
//	//	cond["i."+consts.TcProjectObjectTypeId] = projectObjectTypeId
//	//}
//	if iterationId != 0 {
//		cond["i."+consts.TcIterationId] = iterationId
//	}
//	if !isParent {
//		cond["i."+consts.TcParentId] = db.NotEq(0)
//	} else {
//		cond["i."+consts.TcParentId] = 0
//	}
//	err = conn.Select(db.Raw("i.id, i.code, i.project_object_type_id, i.property_id, i.source_id, i.issue_object_type_id, "+
//		"i.iteration_id, i.title, i.priority_id, i.plan_start_time, i.plan_end_time, i.custom_field, i.end_time, "+
//		"i.parent_id, i.owner, i.status, i.audit_status, i.creator, i.create_time, d.remark")).
//		From("ppm_pri_issue i", "ppm_pri_issue_detail d").Where(cond).OrderBy("i.project_object_type_id").All(resBo)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
//	}
//	// 根据得到的数据，查询它们的自定义字段数据
//	//根据所有符合条件的id去查询无码
//	var allUsefulIssueIds []interface{}
//	issueInfoMap := map[int64]bo.IssueAndDetailInfoBo{}
//	for _, issueBo := range *resBo {
//		allUsefulIssueIds = append(allUsefulIssueIds, issueBo.Id)
//		issueInfoMap[issueBo.Id] = issueBo
//	}
//	allNeedIssueBos := make([]bo.IssueAndDetailInfoBo, 0, len(allUsefulIssueIds))
//	if len(allUsefulIssueIds) < 1 {
//		return allNeedIssueBos, nil
//	}
//
//	orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
//		OrgId:  orgId,
//		UserId: 0,
//	})
//	if orgInfoResp.Failure() {
//		log.Error(orgInfoResp.Error())
//		return nil, orgInfoResp.Error()
//	}
//	orgRemarkObj := &orgvo.OrgRemarkConfigType{}
//	oriErr := json.FromJson(orgInfoResp.OrganizationInfo.Remark, orgRemarkObj)
//	if oriErr != nil {
//		log.Error(oriErr)
//		return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
//	}
//	//暂时只在单项目查询的时候传入，因为综合查询是不需要展示自定义字段的
//	appId, summaryAppId := int64(0), orgRemarkObj.OrgSummaryTableAppId
//	if projectId != 0 {
//		appId, businessErr = GetAppIdFromProjectId(orgId, projectId)
//		if businessErr != nil {
//			log.Error(businessErr)
//			return nil, businessErr
//		}
//	} else {
//		appId = orgRemarkObj.OrgSummaryTableAppId
//	}
//	if len(allUsefulIssueIds) <= queryLimit {
//		// 查询逻辑
//		partIssueBoArr, businessErr := GetIssuesForExport(orgId, projectId, allUsefulIssueIds, appId, summaryAppId, issueInfoMap, customColumns, infoMapObj)
//		if businessErr != nil {
//			log.Errorf("[AllIssueForProject] partIssueBoArr len: %d, GetIssuesForExport err: %v", len(partIssueBoArr), businessErr)
//			return nil, businessErr
//		}
//		allNeedIssueBos = partIssueBoArr
//	} else {
//		// 根据切片分页
//		total := len(allUsefulIssueIds)
//		offset, batch := 0, queryLimit
//		for {
//			limit := offset + batch
//			if limit > total {
//				limit = total
//			}
//			tmpIssueIdsSlice := allUsefulIssueIds[offset:limit]
//			log.Infof("[AllIssueForProject] export query issue this page count: %d", len(tmpIssueIdsSlice))
//			// 查询逻辑
//			partIssueBoArr, businessErr := GetIssuesForExport(orgId, projectId, tmpIssueIdsSlice, appId, summaryAppId, issueInfoMap, customColumns, infoMapObj)
//			if businessErr != nil {
//				log.Errorf("[AllIssueForProject] GetIssuesForExport err: %v", businessErr)
//				return nil, businessErr
//			}
//			allNeedIssueBos = append(allNeedIssueBos, partIssueBoArr...)
//
//			if limit >= total {
//				break
//			}
//			offset += batch
//		}
//	}
//	log.Infof("[AllIssueForProject] all issueId len: %d, export query result len: %d, ", len(allUsefulIssueIds), len(allNeedIssueBos))
//
//	return allNeedIssueBos, nil
//}

//// GetIssuesForExport 导出任务时的任务详情查询
//// appId：汇总表 appId
//func GetIssuesForExport(orgId, projectId int64, issueIdsIf []interface{}, appId, summaryAppId int64,
//	issueInfoMap map[int64]bo.IssueAndDetailInfoBo, customColumns []*projectvo.TableColumnData,
//	infoMapObj *projectvo.IssueInfoMapForExportIssue,
//) ([]bo.IssueAndDetailInfoBo, errs.SystemErrorInfo) {
//	allUsefulIssueIds := issueIdsIf
//	lessReq := vo.LessCondsData{
//		Type:   "in",
//		Values: allUsefulIssueIds,
//		Column: "issueId",
//		Left:   nil,
//		Right:  nil,
//		Conds:  nil,
//	}
//	page, size := 1, 2000
//	lessReqParam := formvo.LessIssueListReq{
//		Condition: lessReq,
//		AppId:     summaryAppId,
//		OrgId:     orgId,
//		Page:      int64(page),
//		Size:      int64(size),
//	}
//	// appIds 脏数据，因此这里先忽略 appIds 的条件查询。而且，issueIds 是精确查询
//	//暂时只在单项目查询的时候传入，因为综合查询是不需要展示自定义字段的
//	//if projectId != 0 {
//	//	lessReqParam.RedirectIds = []int64{appId}
//	//}
//	lessResp := formfacade.LessIssueList(lessReqParam)
//	if lessResp.Failure() {
//		log.Error(lessResp.Error())
//		return nil, lessResp.Error()
//	}
//	//真正符合条件的数据
//	selectedIssuesFormLc := map[int64]bo.IssueAndDetailInfoBo{}
//	for _, m := range lessResp.Data.List {
//		i, ok := m["issueId"]
//		if !ok {
//			continue
//		}
//		issueId := int64(0)
//		if id, ok := i.(int64); ok {
//			issueId = id
//		} else if id, ok1 := i.(int); ok1 {
//			issueId = int64(id)
//		} else if id, ok1 := i.(float64); ok1 {
//			//理论上 map 解析json 会将int转为float64
//			issueIdStr := strconv.FormatFloat(id, 'f', -1, 64)
//			parseId, err := strconv.ParseInt(issueIdStr, 10, 64)
//			if err != nil {
//				log.Error(err)
//				continue
//			} else {
//				issueId = parseId
//			}
//		} else {
//			continue
//		}
//		// 处理自定义字段
//		if issueBo, ok := issueInfoMap[issueId]; ok {
//			//issueBo.CustomField, _ = TransferLcDataIntoExcelString(customColumns, m, infoMapObj, nil)
//			issueBo.TableId = ConvertInterfaceIntoInt64(m[consts.BasicFieldTableId])
//			selectedIssuesFormLc[issueBo.Id] = issueBo
//		}
//	}
//	allNeedIssueBos := make([]bo.IssueAndDetailInfoBo, 0)
//	for _, id := range allUsefulIssueIds {
//		if issueBo, ok := selectedIssuesFormLc[id.(int64)]; ok {
//			allNeedIssueBos = append(allNeedIssueBos, issueBo)
//		}
//	}
//
//	return allNeedIssueBos, nil
//}

//// QueryLcIssueList 从无码接口中查询任务数据
//// @param appId： 汇总表id或者项目的应用id
//func QueryLcIssueList(orgId, appId int64, relateUserId *int64, condParam map[string]interface{}) ([]map[string]interface{}, errs.SystemErrorInfo) {
//	issueArr := make([]map[string]interface{}, 0)
//	condList := make([]*vo.LessCondsData, 0, 5)
//	// appId 为必传
//	if appId < 1 {
//		return issueArr, nil
//	}
//	// 组装筛选条件
//	if val, ok := condParam[consts.BasicFieldAppIds]; ok {
//		condList = append(condList, &vo.LessCondsData{
//			Column: consts.BasicFieldAppIds,
//			Type:   "values_in",
//			Value:  val,
//			Values: []interface{}{val},
//		})
//	}
//	if val, ok := condParam[consts.BasicFieldProjectId]; ok {
//		condList = append(condList, &vo.LessCondsData{
//			Column: consts.BasicFieldProjectId,
//			Type:   "equal",
//			Value:  val,
//		})
//	}
//	if val, ok := condParam[consts.BasicFieldIterationId]; ok {
//		condList = append(condList, &vo.LessCondsData{
//			Column: consts.BasicFieldIterationId,
//			Type:   "equal",
//			Value:  val,
//		})
//	}
//	// recycleFlag 为 2 表示不包含回收站的任务
//	if val, ok := condParam[consts.BasicFieldRecycleFlag]; ok {
//		condList = append(condList, &vo.LessCondsData{
//			Column: consts.BasicFieldRecycleFlag,
//			Type:   "equal",
//			Value:  val,
//		})
//	}
//	// 参考 LcHomeIssuesForProject 函数中查询任务条件
//	if val, ok := condParam[consts.BasicFieldDelFlag]; ok {
//		condList = append(condList, &vo.LessCondsData{
//			Column: consts.BasicFieldDelFlag,
//			Type:   "equal",
//			Value:  val,
//		})
//	}
//
//	// 状态分类。
//	if val, ok := condParam[consts.BasicFieldIssueStatusType]; ok {
//		condList = append(condList, &vo.LessCondsData{
//			Column: consts.BasicFieldIssueStatusType,
//			Type:   "in",
//			Values: val.([]interface{}),
//		})
//	}
//	// 处理 issue relation。例如我负责的、我关注的等筛选
//	// 如： {issueRelationType： 1} 表示查询负责人为 curUserId 的任务（1 表示关联关系为“负责人”）
//	if val, ok := condParam["issueRelationType"]; ok {
//		if relateUserId != nil {
//			tmpCondArr := AssemblyIssueRelationCondForLcQueryIssue(map[int]interface{}{
//				val.(int): *relateUserId,
//			})
//			condList = append(condList, tmpCondArr...)
//		}
//	}
//
//	filterColumns := []string{
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldId),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldProjectId),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldTableId),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldAuditStatus),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueStatus),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueStatusType),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldCreateTime),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldPlanEndTime),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldEndTime),
//		lc_helper.ConvertToFilterColumn(consts.BasicFieldOwnerChangeTime),
//	}
//
//	// 查询
//	lessReqCondition := vo.LessCondsData{
//		Type:  "and",
//		Conds: condList,
//	}
//	for page := int64(1); page < 10; page++ {
//		req := formvo.LessIssueListReq{
//			FilterColumns: filterColumns,
//			Condition:     lessReqCondition,
//			AppId:         appId,
//			OrgId:         orgId,
//			Page:          page,
//			Size:          2000,
//		}
//		lessResp := formfacade.LessIssueList(req)
//		if lessResp.Failure() {
//			log.Errorf("[QueryLcIssueList] orgId: %d, appId: %d, curPape: %d, err: %v", orgId, appId, page, lessResp.Error())
//			return issueArr, lessResp.Error()
//		}
//		if len(lessResp.Data.List) > 0 {
//			issueArr = append(issueArr, lessResp.Data.List...)
//		} else {
//			break
//		}
//	}
//
//	return issueArr, nil
//}

// AssemblyIssueRelationCondForLcQueryIssue 组装查询任务的条件
func AssemblyIssueRelationCondForLcQueryIssue(relateValMap map[int]interface{}) []*tablev1.Condition {
	tmpCondArr := make([]*tablev1.Condition, 0)
	for relateType, val := range relateValMap {
		values := make([]interface{}, 0)
		if trueVal, ok := val.(int64); ok {
			uidStr := fmt.Sprintf("U_%d", trueVal)
			values = append(values, uidStr)
		} else {
			continue
		}

		switch relateType {
		case consts.IssueRelationTypeOwner:
			tmpCondArr = append(tmpCondArr, GetRowsCondition(consts.BasicFieldOwnerId, tablev1.ConditionType_values_in, 0, values))
		case consts.IssueRelationTypeFollower:
			tmpCondArr = append(tmpCondArr, GetRowsCondition(consts.BasicFieldFollowerIds, tablev1.ConditionType_values_in, 0, values))
		case consts.IssueRelationTypeAuditor:
			tmpCondArr = append(tmpCondArr, GetRowsCondition(consts.BasicFieldAuditorIds, tablev1.ConditionType_values_in, 0, values))
		case 6: // 我负责的 + 我关注的
			tmpCondArr = append(tmpCondArr, &tablev1.Condition{
				Type: tablev1.ConditionType_or,
				Conditions: []*tablev1.Condition{
					GetRowsCondition(consts.BasicFieldOwnerId, tablev1.ConditionType_values_in, 0, values),
					GetRowsCondition(consts.BasicFieldFollowerIds, tablev1.ConditionType_values_in, 0, values),
				},
			})
		}
	}

	return tmpCondArr
}

func GetCondForOwner(valArr []interface{}) *vo.LessCondsData {
	return &vo.LessCondsData{
		Column: consts.BasicFieldOwnerId,
		Type:   "values_in",
		Values: valArr,
	}
}

func GetCondForFollower(valArr []interface{}) *vo.LessCondsData {
	return &vo.LessCondsData{
		Column: consts.BasicFieldFollowerIds,
		Type:   "values_in",
		Values: valArr,
	}
}

//func CheckIsAdvancedColumn(field *projectvo.TableColumnField) bool {
//	return CheckIsAdvancedColumnOfMember(field) || CheckIsAdvancedColumnOfDatepicker(field)
//}

// CheckIsAdvancedColumnOfMember 检查当前列是否是成员类型的高级字段（创建人、更新人），如果是高级字段，则可能引用原生值（如：创建人，更新人等）
//func CheckIsAdvancedColumnOfMember(field *projectvo.TableColumnField) bool {
//	propsStr := json.ToJsonIgnoreError(field.Props)
//	propsObj := projectvo.FormConfigColumnFieldMemberProps{}
//	json.FromJson(propsStr, &propsObj)
//
//	return IsUseCreatorForLcColumn(propsObj) || IsUseUpdatorForLcColumn(propsObj)
//}

//func CheckIsAdvancedColumnOfDatepicker(field *projectvo.TableColumnField) bool {
//	propsStr := json.ToJsonIgnoreError(field.Props)
//	propsObj := projectvo.FormConfigColumnFieldDatepickerProps{}
//	json.FromJson(propsStr, &propsObj)
//
//	return IsUseUpdateTimeForLcColumn(propsObj) || IsUseCreateTimeForLcColumn(propsObj)
//}

//// FormatDatetimeForLcColumn 前后端时间格式化库不一致，这部分实现待定 todo
//func FormatDatetimeForLcColumn(timeStr string, propsObj projectvo.FormConfigColumnFieldDatepickerProps) string {
//	if timeStr == "" {
//		return timeStr
//	}
//	// format :
//	// 1.本地格式：`01/02/2006 15:04`
//	// 2.友好格式：`Month 02, 2022 15:04` eg: Jan 01, 2022 21:12
//	// 3.欧洲格式：`01/02/2006 15:04`
//	// 4.国际标准化格式：`2006-01-02 15:04:05`
//	return ""
//}

//// IsUseUpdateTimeForLcColumn 检查是否是“修改时间”类型的列
//func IsUseUpdateTimeForLcColumn(propsObj projectvo.FormConfigColumnFieldDatepickerProps) bool {
//	return propsObj.Datepicker.IsUseUpdateTime || propsObj.IsUseUpdateTime || propsObj.IsUseUpdate
//}
//
//// IsUseCreateTimeForLcColumn 检查是否是“创建时间”类型的列
//func IsUseCreateTimeForLcColumn(propsObj projectvo.FormConfigColumnFieldDatepickerProps) bool {
//	return propsObj.Datepicker.IsUseCreateTime || propsObj.IsUseCreateTime || propsObj.IsUseCreate
//}
//
//// IsUseCreatorForLcColumn 检查是否是“创建人”类型的列
//func IsUseCreatorForLcColumn(propsObj projectvo.FormConfigColumnFieldMemberProps) bool {
//	return propsObj.Member.IsUseCreator || propsObj.IsUseCreator
//}
//
//// IsUseUpdatorForLcColumn 检查是否是“修改人”类型的列
//func IsUseUpdatorForLcColumn(propsObj projectvo.FormConfigColumnFieldMemberProps) bool {
//	return propsObj.Member.IsUseUpdator || propsObj.IsUseUpdator
//}

//func ParseLcDataValueForOneOrMoreMember(columnKey string, lcData map[string]interface{},
//	propsObj projectvo.FormConfigColumnFieldMemberProps, infoMapObj *projectvo.IssueInfoMapForExportIssue,
//	defaultVal interface{},
//) bo.CustomFieldItemFrontendForLc {
//	cellVal := lcData[columnKey]
//	memberInfoObj := bo.CustomFieldItemFrontendForLc{}
//	if cellVal == nil {
//		cellVal = defaultVal
//	}
//	if cellVal == nil {
//		return memberInfoObj
//	}
//	valUserJson := json.ToJsonIgnoreError(cellVal)
//
//	userIdFlagArr := make([]string, 0)
//	// 成员开启多选的情况
//	valUserArr := make([]bo.SimpleUserWithUserIdBo, 0)
//	// 成员未开启多选的情况
//	valUser := bo.SimpleUserWithUserIdBo{}
//	// 成员有两种情况，一种是直接是数组对象 `[{"id": 1001, "name": "ZhangSan1" ...}]`；
//	// 一种是字符串数组 `["U_1011"]`，这种就需要匹配出 id，从而拿到用户信息
//	//isUidFlag := false
//	if propsObj.Multiple {
//		err := json.FromJson(valUserJson, &valUserArr)
//		if err != nil || len(valUserArr) < 1 {
//			// 另一种方式的尝试
//			err := json.FromJson(valUserJson, &userIdFlagArr)
//			if err != nil {
//				log.Errorf("[ParseLcDataValueForOneOrMoreMember] try multi uid flag parse err: %v", err)
//				return memberInfoObj
//			}
//			//isUidFlag = true
//		}
//	} else {
//		err := json.FromJson(valUserJson, &valUser)
//		if err != nil || valUser.Id == 0 {
//			// 另一种方式的尝试
//			// 如果是uid标识，即使单选，也是放在数组里
//			err := json.FromJson(valUserJson, &userIdFlagArr)
//			if err != nil {
//				log.Errorf("[ParseLcDataValueForOneOrMoreMember] try uid flag parse for single err: %v", err)
//				return memberInfoObj
//			}
//			//isUidFlag = true
//		}
//	}
//
//	userNames := make([]string, 0)
//	//if propsObj.Multiple {
//	//	if isUidFlag {
//	//		userMapKeyByUid := ConvertUserInfoMap2KeyByUid(infoMapObj.UserInfoMap)
//	//		userNames = GetUserNamesByUidFlags(userIdFlagArr, userMapKeyByUid)
//	//	} else {
//	//		for _, tmpUser := range valUserArr {
//	//			userNames = append(userNames, tmpUser.Name)
//	//		}
//	//	}
//	//} else {
//	//	if isUidFlag {
//	//		userMapKeyByUid := ConvertUserInfoMap2KeyByUid(infoMapObj.UserInfoMap)
//	//		userNames = GetUserNamesByUidFlags(userIdFlagArr, userMapKeyByUid)
//	//		if len(userNames) > 0 {
//	//			// 单选保证只取一个
//	//			userNames = userNames[:1]
//	//		}
//	//	} else {
//	//		userNames = append(userNames, valUser.Name)
//	//	}
//	//}
//	memberInfoObj.Title = strings.Join(userNames, "||")
//	memberInfoObj.Value = cellVal
//	memberInfoObj.FieldId = columnKey
//
//	return memberInfoObj
//}

//// GetUserNamesByUidFlags 将用户标识数组查询出对应的姓名数组。如：[“U_1001”] -> ["张三"]
//func GetUserNamesByUidFlags(userIdFlagArr []string, userMapKeyByUid map[int64]bo.SimpleUserWithUserIdBo) []string {
//	userNames := make([]string, 0, len(userIdFlagArr))
//	for _, tmpUidFlag := range userIdFlagArr {
//		numberStr := strings.Replace(tmpUidFlag, "U_", "", 1)
//		tmpUid, err := strconv.ParseInt(numberStr, 10, 64)
//		if err != nil {
//			log.Infof("[ParseLcDataValueForOneOrMoreMember] parse uid flag err: %v", err)
//			continue
//		}
//		if user, ok := userMapKeyByUid[tmpUid]; ok {
//			userNames = append(userNames, user.Name)
//		}
//	}
//
//	return userNames
//}

//// ConvertUserInfoMap2KeyByUid 将用户数据map转换为以 uid 为 key 的 map
//func ConvertUserInfoMap2KeyByUid(userMap map[string]bo.SimpleUserWithUserIdBo) map[int64]bo.SimpleUserWithUserIdBo {
//	newUserMap := make(map[int64]bo.SimpleUserWithUserIdBo, len(userMap))
//	for _, user := range userMap {
//		newUserMap[user.UserId] = user
//	}
//
//	return newUserMap
//}

//func GetIssueRelationTypesForProPrivacy() []int {
//	return consts.IssueParticipantTypeList
//}

func CheckIsAdmin(orgId, curUserId, appId int64) bool {
	if curUserId == 0 {
		return false
	}
	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, curUserId)
	if manageAuthInfoResp.Failure() {
		log.Errorf("[CheckIsAdmin] orgId:%v, curUserId:%v, error:%v", orgId, curUserId, manageAuthInfoResp.Error())
		return false
	}
	adminFlag := manageAuthInfoResp.Data
	manageApps := manageAuthInfoResp.Data.ManageApps

	isManageApp, err := slice.Contain(manageApps, appId)
	if err != nil {
		log.Errorf("[CheckIsAdmin] slice.Contain err:%v", err)
		return false
	}

	isSubAdmin := false
	if (len(manageApps) > 0 && manageApps[0] == -1) || isManageApp {
		isSubAdmin = true
	}

	return adminFlag.IsOrgOwner || adminFlag.IsSysAdmin || isSubAdmin
}

// CheckIsProAdmin 检查是否是项目的管理员
func CheckIsProAdmin(orgId, appId, curUserId int64, appAuthInfo *appauth.GetAppAuthData) (bool, errs.SystemErrorInfo) {
	if appId == 0 || curUserId == 0 {
		return false, nil
	}
	if appAuthInfo == nil {
		optAuthResp := permissionfacade.GetAppAuth(orgId, appId, 0, curUserId)
		if optAuthResp.Failure() {
			log.Errorf("[CheckIsProAdmin] orgId:%v, appId:%v, curUserId:%v, error:%v", orgId, appId, curUserId, optAuthResp.Error())
			return false, optAuthResp.Error()
		}
		appAuthInfo = &optAuthResp.Data
	}
	isAdmin := appAuthInfo.HasAppRootPermission || appAuthInfo.SysAdmin ||
		appAuthInfo.OrgOwner || appAuthInfo.AppOwner

	return isAdmin, nil
}

// CheckIsProAdminWithoutCollaborator 检查是否是项目的管理员，不合并协作人角色权限信息
//func CheckIsProAdminWithoutCollaborator(orgId, appId, curUserId int64, appAuthInfo *appauth.GetAppAuthData) (bool, errs.SystemErrorInfo) {
//	if appId == 0 || curUserId == 0 {
//		return false, nil
//	}
//	if appAuthInfo == nil {
//		optAuthResp := permissionfacade.GetAppAuthWithoutCollaborator(orgId, appId, curUserId)
//		if optAuthResp.Failure() {
//			log.Errorf("[CheckIsProAdminWithoutCollaborator] orgId:%v, appId:%v, curUserId:%v, error:%v", orgId, appId, curUserId, optAuthResp.Error())
//			return false, optAuthResp.Error()
//		}
//		appAuthInfo = &optAuthResp.Data
//	}
//	isAdmin := appAuthInfo.HasAppRootPermission || appAuthInfo.SysAdmin ||
//		appAuthInfo.OrgOwner || appAuthInfo.AppOwner
//
//	return isAdmin, nil
//}

//func GetCollaborateIssues(orgId, currentUserId int64, summaryAppId int64, projectIds []int64) ([]int64, errs.SystemErrorInfo) {
//	deptResp := orgfacade.GetUserDeptIds(&orgvo.GetUserDeptIdsReq{
//		OrgId:  orgId,
//		UserId: currentUserId,
//	})
//	if deptResp.Failure() {
//		log.Errorf("[GetCollaborateIssues] orgId:%v, curUserId:%v, failure:%v", orgId, currentUserId, deptResp.Error())
//		return nil, deptResp.Error()
//	}
//	allDeptIds := []int64{0}
//	allDeptIds = append(allDeptIds, deptResp.Data.DeptIds...)
//
//	allValues := []interface{}{fmt.Sprintf("U_%d", currentUserId)}
//	for _, id := range allDeptIds {
//		allValues = append(allValues, fmt.Sprintf("D_%d", id))
//	}
//	req := formvo.LessIssueListReq{
//		Condition: vo.LessCondsData{
//			Type: "and",
//			Conds: []*vo.LessCondsData{
//				&vo.LessCondsData{
//					Type:      "values_in",
//					FieldType: nil,
//					Value:     nil,
//					Values:    allValues,
//					Column:    "jsonb_path_query_array(\"data\"::jsonb, '$.collaborators.*[*]')::jsonb",
//					Left:      nil,
//					Right:     nil,
//					Conds:     nil,
//				},
//			},
//		},
//		Orders:      nil,
//		RedirectIds: nil,
//		AppId:       summaryAppId,
//		OrgId:       orgId,
//		UserId:      currentUserId,
//	}
//	if projectIds != nil && len(projectIds) > 0 {
//		projectIdsInterface := make([]interface{}, 0)
//		copyer.Copy(projectIds, &projectIdsInterface)
//		req.Condition.Conds = append(req.Condition.Conds, &vo.LessCondsData{
//			Type:      "in",
//			FieldType: nil,
//			Value:     nil,
//			Values:    projectIdsInterface,
//			Column:    consts.BasicFieldProjectId,
//			Left:      nil,
//			Right:     nil,
//			Conds:     nil,
//		})
//	}
//	issueResp := formfacade.LessIssueList(req)
//	if issueResp.Failure() {
//		log.Error(issueResp.Error())
//		return nil, issueResp.Error()
//	}
//	allCollaborateIssueIds := []int64{}
//	for _, m := range issueResp.Data.List {
//		i, ok := m["issueId"]
//		if !ok {
//			continue
//		}
//		issueId := int64(0)
//		if id, ok := i.(int64); ok {
//			issueId = id
//		} else if id, ok1 := i.(int); ok1 {
//			issueId = int64(id)
//		} else if id, ok1 := i.(float64); ok1 {
//			//理论上 map 解析json 会将int转为float64
//			issueIdStr := strconv.FormatFloat(id, 'f', -1, 64)
//			parseId, err := strconv.ParseInt(issueIdStr, 10, 64)
//			if err != nil {
//				log.Error(err)
//				continue
//			} else {
//				issueId = parseId
//			}
//		} else {
//			continue
//		}
//
//		allCollaborateIssueIds = append(allCollaborateIssueIds, issueId)
//	}
//
//	return allCollaborateIssueIds, nil
//}

//func GetMyProjectDataQueryConds(orgId, currentUserId int64, deptIds []int64) *vo.LessCondsData {
//	ids, err := GetAllMyProjectIdsWithDeptIds(orgId, currentUserId, deptIds, false)
//	if err != nil {
//		log.Error(err)
//		return nil
//	}
//	myProjectIds := make([]interface{}, len(ids))
//	for i := range ids {
//		myProjectIds[i] = ids[i]
//	}
//	return &vo.LessCondsData{
//		Type:   "in",
//		Values: myProjectIds,
//		Column: consts.BasicFieldProjectId,
//	}
//}

func GetProjectIdsQueryConds(projectIds []int64) *vo.LessCondsData {
	projectIdsInterface := make([]interface{}, len(projectIds))
	for i := range projectIds {
		projectIdsInterface[i] = projectIds[i]
	}
	return &vo.LessCondsData{
		Type:   "in",
		Values: projectIdsInterface,
		Column: consts.BasicFieldProjectId,
	}
}

func GetTableIdsQueryCondsByAppAuths(appAuths map[int64]appauth.GetAppAuthData) *vo.LessCondsData {
	tableIdsInterface := make([]interface{}, 0)
	for _, appAuth := range appAuths {
		for _, tableIdStr := range appAuth.TableAuth {
			tableIdsInterface = append(tableIdsInterface, tableIdStr)
		}
	}
	if len(tableIdsInterface) == 0 {
		return nil
	} else {
		return &vo.LessCondsData{
			Type:   "not_in",
			Values: tableIdsInterface,
			Column: consts.BasicFieldTableId,
		}
	}
}

func GetTableIdsQueryCondsByAppAuth(appAuth *appauth.GetAppAuthData) *vo.LessCondsData {
	if appAuth == nil {
		return nil
	}
	tableIdsInterface := make([]interface{}, 0)
	for _, tableIdStr := range appAuth.TableAuth {
		tableIdsInterface = append(tableIdsInterface, tableIdStr)
	}
	if len(tableIdsInterface) == 0 {
		return nil
	} else {
		return &vo.LessCondsData{
			Type:   "not_in",
			Values: tableIdsInterface,
			Column: consts.BasicFieldTableId,
		}
	}
}

func GetCollaborateDataQueryConds(orgId, currentUserId int64, deptIds []int64) *vo.LessCondsData {
	allValues := []string{fmt.Sprintf("U_%d", currentUserId)}

	// 查询协作人时需要包含D_0，代表协作人是整个组织的情况
	allValues = append(allValues, "D_0")
	for _, id := range deptIds {
		allValues = append(allValues, fmt.Sprintf("D_%d", id))
	}
	return &vo.LessCondsData{
		Type:   "raw_sql",
		Column: "collaborators && ARRAY['" + strings.Join(allValues, "','") + "']",
	}
}

func GetIssueIdDetailQueryConds(issueId int64, path string) vo.LessCondsData {
	// 任务本身
	allValues1 := []interface{}{issueId}

	// 父任务
	allValues2 := []interface{}{}
	ss := strings.Split(path, ",")
	for _, s := range ss {
		id := cast.ToInt64(s)
		if id > 0 {
			allValues2 = append(allValues2, id)
		}
	}

	orCond := vo.LessCondsData{Type: "or"}
	// 任务本身
	orCond.Conds = append(orCond.Conds, &vo.LessCondsData{
		Type:   "in",
		Values: allValues1,
		Column: consts.BasicFieldIssueId,
	})
	// 父任务
	if len(allValues2) > 0 {
		orCond.Conds = append(orCond.Conds, &vo.LessCondsData{
			Type: "and",
			Conds: []*vo.LessCondsData{
				&vo.LessCondsData{
					Type:   "in",
					Values: allValues2,
					Column: consts.BasicFieldIssueId,
				},
				&vo.LessCondsData{
					Type:   "equal",
					Value:  2,
					Column: "recycleFlag",
				},
			},
		})
	}
	// 子任务
	orCond.Conds = append(orCond.Conds, &vo.LessCondsData{
		Type: "and",
		Conds: []*vo.LessCondsData{
			&vo.LessCondsData{
				Type:   "equal",
				Value:  issueId,
				Column: consts.BasicFieldParentId,
			},
			&vo.LessCondsData{
				Type:   "equal",
				Value:  2,
				Column: "recycleFlag",
			},
		},
	})

	return vo.LessCondsData{
		Type:  "and",
		Conds: []*vo.LessCondsData{&orCond},
	}
}

// GetIssueStatusIdsByType 通过状态分类查询对应的任务状态 ids  todo
// proAppId: 项目的 appId
// typ 状态分类: 1未开始；2进行中；3已完成
func GetIssueStatusIdsByType(orgId, proAppId int64, typ int) ([]int64, errs.SystemErrorInfo) {
	resStatusIds := make([]int64, 0)
	statusBoArr, err := GetIssueStatusBoArr(orgId, proAppId) // todo
	if err != nil {
		log.Errorf("[GetIssueStatusIdsByType] err: %v", err)
		return resStatusIds, err
	}
	for _, statusBo := range statusBoArr {
		if statusBo.StatusType == typ {
			resStatusIds = append(resStatusIds, statusBo.StatusId)
		}
	}

	return resStatusIds, nil
}

// GetIssueStatusBoArr 查询项目下任务的状态列的状态列表 todo
// proAppId: 项目的 appId
func GetIssueStatusBoArr(orgId, proAppId int64) ([]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	statusBoArr := make([]bo.CacheProcessStatusBo, 0)
	// todo
	//tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
	//	OrgId:  0,
	//	UserId: 0,
	//	Input: &tableV1.ReadTableSchemasRequest{
	//		TableIds:  nil,
	//		ColumnIds: nil,
	//	},
	//})
	return statusBoArr, nil
}

// GetStatusListByCategoryRelaxed 根据 category 获取对应的状态列表
func GetStatusListByCategoryRelaxed(orgId int64, category int) ([]bo.CacheProcessStatusBo, errs.SystemErrorInfo) {
	statusBoList := make([]bo.CacheProcessStatusBo, 0)
	switch category {
	case consts.ProcessStatusCategoryIteration:
		allStatus := consts.IterationStatusList
		for _, item := range allStatus {
			statusBoList = append(statusBoList, bo.CacheProcessStatusBo{
				StatusId:    item.ID,
				StatusType:  item.Type,
				Category:    consts.ProcessStatusCategoryIteration,
				Name:        item.Name,
				DisplayName: item.DisplayName,
				BgStyle:     item.BgStyle,
				FontStyle:   item.FontStyle,
			})
		}
	case consts.ProcessStatusCategoryProject: // 项目状态
		allStatus := consts.ProjectStatusList
		for _, item := range allStatus {
			statusBoList = append(statusBoList, bo.CacheProcessStatusBo{
				StatusId:    item.ID,
				StatusType:  item.Type,
				Category:    consts.ProcessStatusCategoryProject,
				Name:        item.Name,
				DisplayName: item.DisplayName,
				BgStyle:     item.BgStyle,
				FontStyle:   item.FontStyle,
			})
		}
	case consts.ProcessStatusCategoryIssue: // 任务的状态需要取表头信息才能确定
		return nil, errs.BuildSystemErrorInfoWithMessage(errs.BaseDomainError, "该方法不适用。任务状态的查询，请更换一个方式查询！")
	}
	return statusBoList, nil
}

// GetIssueListByTableIdWithChaId 根据tableId查询表下outChatId不为空的所有任务
func GetIssueListByTableIdWithChaId(orgId int64, tableId int64) ([]map[string]interface{}, errs.SystemErrorInfo) {
	condition := &tablev1.Condition{
		Type: tablev1.ConditionType_and,
		Conditions: []*tablev1.Condition{
			GetRowsCondition(consts.BasicFieldTableId, tablev1.ConditionType_equal, cast.ToString(tableId), nil),
			GetRowsCondition(consts.BasicFieldOutChatId, tablev1.ConditionType_not_null, cast.ToString(tableId), nil),
		},
	}

	issueMaps := []map[string]interface{}{}
	for page := int64(1); ; page++ {
		issueLcDatas, err := GetIssueInfosMapLc(orgId, 0, condition, nil, page, 2000)
		if err != nil {
			log.Errorf("[GetIssueListByTableId]getIssueInfosMapLc err:%v, orgId:%d, tableId:%d", err, orgId, tableId)
			return nil, err
		}
		issueMaps = append(issueMaps, issueLcDatas...)
		if len(issueLcDatas) <= 2000 {
			break
		}
	}
	return issueMaps, nil
}

// GetIssueListWithRefColumn 获取字段的同时获取引用相关字段
func GetIssueListWithRefColumn(req formvo.LessIssueListReq, columns []*projectvo.TableColumnData, appAuthInfo string) (*projectvo.ListRowsReply, []string) {
	req.NeedRefColumn = true
	if req.TableId > 0 && len(columns) == 0 {
		var err errs.SystemErrorInfo
		columns, err = GetTableColumns(req.OrgId, req.UserId, req.TableId, true)
		if err != nil {
			return &projectvo.ListRowsReply{Err: vo.Err{Code: err.Code(), Message: err.Message()}}, nil
		}
	}

	refColumnIdsMap, unSupportColumnIds := getAllRefColumnIds(req.OrgId, req.TableId, columns)
	if len(refColumnIdsMap) > 0 {
		newFilterColumns := make([]string, 0, len(req.FilterColumns)+len(refColumnIdsMap))
		if len(req.FilterColumns) == 0 {
			newFilterColumns = append(newFilterColumns, consts.LcJsonColumn)
			for columnId, _ := range lc_helper.NotJsonColumnMap {
				if columnId != consts.BasicFieldId && columnId != consts.BasicFieldCollaborators {
					newFilterColumns = append(newFilterColumns, lc_helper.ConvertToFilterColumn(columnId))
				}
			}
		} else {
			for _, column := range req.FilterColumns {
				for refColumnId := range refColumnIdsMap {
					if !strings.Contains(column, refColumnId) {
						newFilterColumns = append(newFilterColumns, column)
					}
				}
			}
		}
		for columnId := range refColumnIdsMap {
			newFilterColumns = append(newFilterColumns, columnId)
		}
		req.FilterColumns = newFilterColumns

		return GetRows(&req, appAuthInfo), unSupportColumnIds
	} else {
		req.NeedRefColumn = false
		return GetRows(&req, appAuthInfo), unSupportColumnIds
	}
}

//func getIssueListAsyn(req formvo.LessIssueListReq, refColumnIdsMap map[string]string) *formvo.LessIssueListResp {
//	wait := sync.WaitGroup{}
//	result := make([]*formvo.LessIssueListResp, 0, len(refColumnIdsMap))
//	resultLock := &sync.Mutex{}
//	j := 0
//	for _, id := range refColumnIdsMap {
//		wait.Add(1)
//		tempReq := req.Copy()
//		if j == 0 {
//			if len(tempReq.FilterColumns) > 0 {
//				newFilterColumns := make([]string, 0, len(tempReq.FilterColumns))
//				for _, columnId := range tempReq.FilterColumns {
//					if refColumnIdsMap[columnId] == "" {
//						newFilterColumns = append(newFilterColumns, columnId)
//					}
//				}
//				tempReq.FilterColumns = newFilterColumns
//			} else {
//				tempReq.FilterColumns = append(tempReq.FilterColumns, consts.LcJsonColumn, id)
//			}
//		} else {
//			req.NeedTotal = false
//			tempReq.FilterColumns = []string{consts.BasicFieldId, lc_helper.ConvertToFilterColumn(id)}
//		}
//		j++
//
//		asyn.Execute(func() {
//			defer wait.Done()
//			resp := formfacade.LessIssueList(*tempReq)
//			resultLock.Lock()
//			result = append(result, resp)
//			resultLock.Unlock()
//		})
//	}
//	wait.Wait()
//
//	return mergeIssueListResult(result)
//}

//func mergeIssueListResult(result []*formvo.LessIssueListResp) *formvo.LessIssueListResp {
//	var mergeResult *formvo.LessIssueListResp
//	maxTotal := int64(0)
//	for _, resp := range result {
//		if resp.Data.Total > 0 {
//			maxTotal = resp.Data.Total
//		}
//		if mergeResult == nil {
//			mergeResult = resp
//		} else {
//			m := make(map[interface{}]map[string]interface{})
//			for _, lcData := range resp.Data.List {
//				m[lcData[consts.BasicFieldId]] = lcData
//			}
//
//			for _, m2 := range mergeResult.Data.List {
//				if m[m2[consts.BasicFieldId]] != nil {
//					for s, i := range m[m2[consts.BasicFieldId]] {
//						m2[s] = i
//					}
//				}
//			}
//		}
//		if resp.Failure() {
//			mergeResult.Err = resp.Err
//		}
//	}
//	mergeResult.Data.Total = maxTotal
//
//	return mergeResult
//}

func getAllRefColumnIds(orgId, tableId int64, columns []*projectvo.TableColumnData) (refColumnIdsMap map[string]string, unSupportColumnIds []string) {
	tableIds := make([]int64, 0, 3)
	refColumnIdsMap = make(map[string]string, 3)
	if tableId == 0 {
		return refColumnIdsMap, unSupportColumnIds
	}

	tableIdColumnIdsMap := make(map[int64][]string, 3)
	columnToRefColumnMap := make(map[int64]map[string]string, 3)
	refColumnIds := make([]string, 0, 3)
	for _, column := range columns {
		if column.Field.Type == consts.LcColumnFieldTypeConditionRef {
			props := column.Field.Props[consts.LcColumnFieldTypeConditionRef]
			if m, ok := props.(map[string]interface{}); ok {
				id := cast.ToInt64(m[consts.BasicFieldTableId])
				refColumnId := cast.ToString(m["columnId"])
				if id > 0 && len(refColumnId) > 0 {
					tableIds = append(tableIds, id)
					if columnToRefColumnMap[id] == nil {
						columnToRefColumnMap[id] = map[string]string{}
					}
					columnToRefColumnMap[id][refColumnId] = column.Name
					refColumnIds = append(refColumnIds, refColumnId)
				}
			}
		}
	}

	if len(tableIds) > 0 {
		reply := tablefacade.GetTableColumns(projectvo.GetTableColumnsReq{
			OrgId:  orgId,
			UserId: 0,
			Input: &tablev1.ReadTableSchemasRequest{
				TableIds:  tableIds,
				ColumnIds: refColumnIds,
			},
		})
		if reply.Failure() {
			log.Errorf("[getAllRefColumnIds] GetTableColumns err :%v", reply.Error())
			return refColumnIdsMap, unSupportColumnIds
		}

		existTableIds := make([]int64, 0, len(tableIds))
		for _, table := range reply.Data.Tables {
			for _, column := range table.Columns {
				if columnToRefColumnMap[table.TableId] != nil && len(columnToRefColumnMap[table.TableId][column.Name]) > 0 {
					if tableIdColumnIdsMap[table.TableId] == nil {
						existTableIds = append(existTableIds, table.TableId)
					}
					tableIdColumnIdsMap[table.TableId] = append(tableIdColumnIdsMap[table.TableId], columnToRefColumnMap[table.TableId][column.Name])
				}
			}
		}
		if len(existTableIds) > 0 {
			reqTableIds := make([]int64, 0, len(existTableIds)+1)
			reqTableIds = append(reqTableIds, existTableIds...)
			reqTableIds = append(reqTableIds, tableId)
			reqTableIds = slice.SliceUniqueInt64(reqTableIds)
			result, err := getTablesDataCount(orgId, reqTableIds)
			if err != nil {
				log.Errorf("[getAllRefColumnIds] getTablesDataCount err :%v", err)
				return refColumnIdsMap, unSupportColumnIds
			}

			for _, id := range existTableIds {
				if result[id]*result[tableId] > consts.MaxTowTableColumns {
					unSupportColumnIds = append(unSupportColumnIds, tableIdColumnIdsMap[id]...)
					delete(tableIdColumnIdsMap, id)
				}
			}
		}

	}

	for _, columnIds := range tableIdColumnIdsMap {
		for _, id := range columnIds {
			refColumnIdsMap[id] = id
		}
	}

	return refColumnIdsMap, unSupportColumnIds
}

func getTablesDataCount(orgId int64, tableIds []int64) (map[int64]int64, errs.SystemErrorInfo) {
	wait := sync.WaitGroup{}
	result := make(map[int64]int64, len(tableIds))
	resultLocker := &sync.RWMutex{}
	for _, tableId := range tableIds {
		wait.Add(1)
		conditions := &tablev1.Condition{
			Type: tablev1.ConditionType_and,
			Conditions: GetNoRecycleCondition(
				GetRowsCondition(consts.BasicFieldTableId, tablev1.ConditionType_equal, cast.ToString(tableId), nil),
			),
		}

		asyn.Execute(func() {
			defer wait.Done()
			id, _ := getRowsCountByCondition(orgId, 0, conditions)
			resultLocker.Lock()
			result[tableId] = id
			resultLocker.Unlock()
		})
		wait.Wait()
	}

	return result, nil
}

func getRowsCountByCondition(orgId, userId int64, condition *tablev1.Condition) (int64, errs.SystemErrorInfo) {
	resp, err := GetRawRows(orgId, userId, &tablev1.ListRawRequest{
		FilterColumns: []string{
			" count(*) as \"all\" ",
		},
		Condition: condition,
	})
	if err != nil {
		return 0, err
	}

	if len(resp.Data) > 0 {
		return cast.ToInt64(resp.Data[0]["all"]), nil
	}

	return 0, nil
}
