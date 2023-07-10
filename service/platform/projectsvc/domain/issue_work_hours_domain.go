package domain

import (
	"fmt"
	"time"

	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	consts2 "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// 创建一条工时记录
func CreateIssueWorkHours(newId int64, param projectvo.CreateIssueWorkHoursReqVo, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	var oriErr error
	if param.Input.WorkerID <= 0 {
		issueId := param.Input.IssueID
		issues, err := GetIssueInfosLc(param.OrgId, param.UserId, []int64{issueId})
		if err != nil {
			return err
		}
		if len(issues) < 1 {
			return errs.IssueNotExist
		}
		issue := issues[0]
		if param.Input.WorkerID < 1 {
			ownerIds, err2 := businees.LcMemberToUserIdsWithError(issue.OwnerId)
			if err2 != nil {
				log.Errorf("原始的数据:%v, err:%v", issue.OwnerId, err)
				return err2
			}
			param.Input.WorkerID = ownerIds[0] // 先取第一个负责人
		}
	}
	if param.Input.Desc == nil {
		defaultDesc := ""
		param.Input.Desc = &defaultDesc
	}
	if param.Input.EndTime == nil {
		defaultEndTime := int64(0)
		param.Input.EndTime = &defaultEndTime
	}
	needTimeInt := format.FormatNeedTimeIntoSecondNumber(param.Input.NeedTime)
	if *param.Input.EndTime > 0 && !format.CheckTimeRangeTimeNumIsValid(needTimeInt, param.Input.StartTime, *param.Input.EndTime) {
		return errs.WorkHourTimeRangeForNeedTimeInvalid
	}
	if needTimeInt > consts2.MaxWorkHour {
		return errs.WorkHourMaxNeedTime
	}
	newPo := po.PpmPriIssueWorkHours{
		Id:                uint64(newId),
		OrgId:             param.OrgId,
		ProjectId:         *param.Input.ProjectID,
		IssueId:           param.Input.IssueID,
		Type:              int8(param.Input.Type),
		WorkerId:          param.Input.WorkerID,
		NeedTime:          uint32(needTimeInt),
		RemainTimeCalType: consts2.RemainTimeCalTypeAuto,
		RemainTime:        0,
		StartTime:         uint32(param.Input.StartTime),
		EndTime:           uint32(*param.Input.EndTime),
		Desc:              *param.Input.Desc,
		Creator:           param.UserId,
		CreateTime:        time.Time{},
		Updator:           param.UserId,
		UpdateTime:        time.Time{},
		Version:           1, // 暂时没有用到的预留字段
		IsDelete:          consts.AppIsNoDelete,
	}
	if tx != nil && len(tx) > 0 {
		oriErr = dao.InsertIssueWorkHours(newPo, tx[0])
	} else {
		oriErr = dao.InsertIssueWorkHours(newPo)
	}
	if oriErr != nil {
		log.Error(oriErr)
		return errs.MysqlOperateError
	}
	return nil
}

// 新增多个工时
func CreateMultiIssueWorkHours(param []*po.PpmPriIssueWorkHours, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	var oriErr error
	if len(tx) > 0 {
		oriErr = dao.InsertMultiIssueWorkHours(param, tx[0])
	} else {
		oriErr = dao.InsertMultiIssueWorkHours(param)
	}
	if oriErr != nil {
		log.Error(oriErr)
		return errs.MysqlOperateError
	}
	return nil
}

// 针对统计时查询列表的条件组装。
func assemblyConditionForStatisticList(condParam bo.IssueWorkHoursBoListCondBo) db.Cond {
	cond := db.Cond{}
	if condParam.OrgId > 0 {
		cond[consts.TcOrgId] = condParam.OrgId
	}
	if len(condParam.Ids) > 0 {
		cond[consts.TcId] = db.In(condParam.Ids)
	}
	if len(condParam.ProjectIds) > 0 {
		cond[consts.TcProjectId] = db.In(condParam.ProjectIds)
	}
	if condParam.IssueId > 0 {
		// 增加空格，防止键被覆盖
		cond[consts.TcIssueId+" "] = condParam.IssueId
	}
	if condParam.IssueIds != nil && len(condParam.IssueIds) > 0 {
		cond[consts.TcIssueId+"  "] = db.In(condParam.IssueIds)
	}
	// 查询多个 workId 的工时
	if condParam.WorkerIds != nil && len(condParam.WorkerIds) > 0 {
		cond[consts.TcWorkerId] = db.In(condParam.WorkerIds)
	}
	// 如果查询删除了的，则传入参数即可
	if condParam.IsDelete > 0 {
		cond[consts.TcIsDelete] = condParam.IsDelete
	} else {
		// 默认查询未删除的数据
		cond[consts.TcIsDelete] = consts.AppIsNoDelete
	}
	return cond
}

// 针对统计时用的工时列表查询
func GetIssueWorkHoursStatisticList(condParam bo.IssueWorkHoursBoListCondBo) ([]po.PpmPriIssueWorkHours, errs.SystemErrorInfo) {
	cond := assemblyConditionForStatisticList(condParam)
	list := []po.PpmPriIssueWorkHours{}

	conn, err := mysql.GetConnect()
	if err != nil {
		return list, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	conn.SetLogging(true)
	cond2 := fmt.Sprintf("((type=2 and start_time>=%[1]v and start_time<=%[2]v) or (type in(1,3) and start_time>=%[1]v and end_time<=%[2]v))", condParam.StartTimeT1, condParam.EndTimeT2)
	mid := conn.Collection(consts.TableIssueWorkHours).Find().Where(db.Raw(cond2)).And(cond)
	if condParam.Size > 0 && condParam.Page > 0 {
		mid = mid.Page(uint(condParam.Page)).Paginate(uint(condParam.Size))
	}
	err = mid.All(&list)
	if err != nil {
		return list, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return list, nil
}

// 获取工时记录列表
func GetIssueWorkHoursList(condParam bo.IssueWorkHoursBoListCondBo, tx ...sqlbuilder.Tx) ([]po.PpmPriIssueWorkHours, errs.SystemErrorInfo) {
	var (
		err  error
		mid  db.Result
		list = []po.PpmPriIssueWorkHours{}
	)

	cond := GetWorkHourCond1(condParam)
	if len(tx) > 0 {
		mid = tx[0].Collection(consts.TableIssueWorkHours).Find(cond)
	} else {
		var conn sqlbuilder.Database
		conn, err = mysql.GetConnect()
		if err != nil {
			return list, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		mid = conn.Collection(consts.TableIssueWorkHours).Find(cond)
	}

	if condParam.Size > 0 && condParam.Page > 0 {
		mid = mid.Page(uint(condParam.Page)).Paginate(uint(condParam.Size))
	}
	err = mid.All(&list)
	if err != nil {
		return list, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	//oriErr := mysql.SelectAllByCondWithNumAndOrder(consts.TableIssueWorkHours, cond, nil, condParam.Page, condParam.Size, nil, &list)
	//if oriErr != nil {
	//	log.Error(oriErr)
	//	return list, errs.MysqlOperateError
	//}

	return list, nil
}

func GetWorkHourCond1(condParam bo.IssueWorkHoursBoListCondBo) db.Cond {
	cond := db.Cond{}
	if condParam.OrgId > 0 {
		cond[consts.TcOrgId] = condParam.OrgId
	}
	if len(condParam.Ids) > 0 {
		cond[consts.TcId] = db.In(condParam.Ids)
	}
	if len(condParam.ProjectIds) > 0 {
		cond[consts.TcProjectId] = db.In(condParam.ProjectIds)
	}
	if condParam.Type > 0 {
		cond[consts.TcType] = condParam.Type
	}
	if condParam.IssueId > 0 {
		cond[consts.TcIssueId] = condParam.IssueId
	}
	if condParam.IssueIds != nil && len(condParam.IssueIds) > 0 {
		cond[consts.TcIssueId] = db.In(condParam.IssueIds)
	}
	//if len(condParam.RawSqls) > 0 {
	//
	//}
	// 查询多个 workId 的工时
	if condParam.WorkerIds != nil && len(condParam.WorkerIds) > 0 {
		cond[consts.TcWorkerId] = db.In(condParam.WorkerIds)
	}
	if condParam.StartTimeT1 > 0 && condParam.EndTimeT2 > 0 {
		cond[consts.TcStartTime] = db.Gte(condParam.StartTimeT1)
	}
	if condParam.EndTimeT2 > 0 {
		cond[consts.TcEndTime] = db.Lte(condParam.EndTimeT2)
	}
	if len(condParam.Types) > 0 {
		cond[consts.TcType] = db.In(condParam.Types)
	}
	// 如果查询删除了的，则传入参数即可
	if condParam.IsDelete > 0 {
		cond[consts.TcIsDelete] = condParam.IsDelete
	} else {
		// 默认查询未删除的数据
		cond[consts.TcIsDelete] = consts.AppIsNoDelete
	}
	return cond
}

func GetCountOfStatisticWorkHourByCond(condParam bo.IssueWorkHoursBoListCondBo) (uint64, errs.SystemErrorInfo) {
	cnt := uint64(0)
	cond := assemblyConditionForStatisticList(condParam)
	conn, err := mysql.GetConnect()
	if err != nil {
		return cnt, errs.BuildSystemErrorInfoWithMessage(errs.MysqlOperateError, err.Error())
	}
	cond2 := fmt.Sprintf("((type=2 and start_time>=%[1]v and start_time<=%[2]v) or (type in(1,3) and start_time>=%[1]v and end_time<=%[2]v))", condParam.StartTimeT1, condParam.EndTimeT2)
	cnt, err = conn.Collection(consts.TableIssueWorkHours).Find().Where(db.Raw(cond2)).And(cond).Count()
	if err != nil {
		return cnt, errs.BuildSystemErrorInfoWithMessage(errs.MysqlOperateError, err.Error())
	}
	return cnt, nil
}

// 获取一条工时记录
func GetOneIssueWorkHoursById(id int64) (*po.PpmPriIssueWorkHours, errs.SystemErrorInfo) {
	oneIssueWorkHours := &po.PpmPriIssueWorkHours{}
	err := mysql.SelectById(consts.TableIssueWorkHours, id, oneIssueWorkHours)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return oneIssueWorkHours, nil
}

// 根据主键，编辑一条工时记录
func UpdateIssueWorkHours(paramBo bo.UpdateOneIssueWorkHoursBo, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	cond := db.Cond{
		consts.TcId: paramBo.Id,
	}
	upd := mysql.Upd{}
	if paramBo.OrgId > 0 {
		upd[consts.TcOrgId] = paramBo.OrgId
	}
	if paramBo.WorkerId > 0 {
		upd[consts.TcWorkerId] = paramBo.WorkerId
	}
	// NeedTime 可以被清空。
	if &paramBo.NeedTime != nil {
		upd[consts.TcNeedTime] = paramBo.NeedTime
	}
	// 时间不能被清空。
	if &paramBo.StartTime != nil && paramBo.StartTime > 0 {
		upd[consts.TcStartTime] = paramBo.StartTime
	}
	if &paramBo.EndTime != nil && paramBo.StartTime > 0 {
		upd[consts.TcEndTime] = paramBo.EndTime
	}
	if paramBo.RemainTimeCalType > 0 {
		upd[consts.TcRemainTimeCalType] = paramBo.RemainTimeCalType
	}
	if paramBo.RemainTime > 0 {
		upd[consts.TcEndTime] = paramBo.RemainTime
	}
	// 支持清空，或填写内容。
	if &paramBo.Desc != nil {
		upd[consts.TcDesc] = paramBo.Desc
	}
	if paramBo.IsDelete > 0 {
		upd[consts.TcIsDelete] = paramBo.IsDelete
	}
	if len(upd) < 1 {
		return errs.UpdateFiledIsEmpty
	}
	var oriErr error
	if tx != nil && len(tx) > 0 {
		_, oriErr = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueWorkHours, cond, upd)
	} else {
		_, oriErr = mysql.UpdateSmartWithCond(consts.TableIssueWorkHours, cond, upd)
	}
	if oriErr != nil {
		log.Error(oriErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	return nil
}

// DelWorkHours 批量删除任务的工时
func DelWorkHours(orgId int64, issueIds []int64, upd mysql.Upd) errs.SystemErrorInfo {
	if len(issueIds) < 1 {
		return nil
	}
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  db.In(issueIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if _, oriErr := mysql.UpdateSmartWithCond(consts.TableIssueWorkHours, cond, upd); oriErr != nil {
		log.Errorf("[DelWorkHours] orgId: %d, err: %v", orgId, oriErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}

	return nil
}

// RecoverWorkHours 从回收站恢复任务时，恢复对应的工时
func RecoverWorkHours(orgId int64, issueIds []int64, recycleVersion int64) errs.SystemErrorInfo {
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  db.In(issueIds),
		consts.TcVersion:  recycleVersion,
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	upd := mysql.Upd{
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if _, oriErr := mysql.UpdateSmartWithCond(consts.TableIssueWorkHours, cond, upd); oriErr != nil {
		log.Errorf("[RecoverWorkHours] orgId: %d, err: %v", orgId, oriErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}

	return nil
}

// 根据条件，编辑工时记录
func UpdateIssueWorkHourByCond(cond db.Cond, paramBo bo.UpdateOneIssueWorkHoursBo, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	upd := mysql.Upd{}
	if paramBo.ProjectId > 0 {
		upd[consts.TcProjectId] = paramBo.ProjectId
	}
	if paramBo.WorkerId > 0 {
		upd[consts.TcWorkerId] = paramBo.WorkerId
	}
	if paramBo.NeedTime > 0 {
		upd[consts.TcNeedTime] = paramBo.NeedTime
	}
	if paramBo.StartTime > 0 {
		upd[consts.TcStartTime] = paramBo.StartTime
	}
	if paramBo.EndTime > 0 {
		upd[consts.TcEndTime] = paramBo.EndTime
	}
	if paramBo.RemainTimeCalType > 0 {
		upd[consts.TcRemainTimeCalType] = paramBo.RemainTimeCalType
	}
	if paramBo.RemainTime > 0 {
		upd[consts.TcEndTime] = paramBo.RemainTime
	}
	if &paramBo.Desc != nil {
		upd[consts.TcDesc] = paramBo.Desc
	}
	if paramBo.IsDelete > 0 {
		upd[consts.TcIsDelete] = paramBo.IsDelete
	}
	if len(cond) < 1 {
		return errs.ConditionInvalid
	}
	if len(upd) < 1 {
		return errs.UpdateFiledIsEmpty
	}
	var oriErr error
	if tx != nil && len(tx) > 0 {
		_, oriErr = mysql.TransUpdateSmartWithCond(tx[0], consts.TableIssueWorkHours, cond, upd)
	} else {
		_, oriErr = mysql.UpdateSmartWithCond(consts.TableIssueWorkHours, cond, upd)
	}
	if oriErr != nil {
		log.Error(oriErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	return nil
}

// 获取符合条件的工时总数
func GetSumWorkHourSumNeedTime(cond db.Cond) (totalNeedTime int64, err errs.SystemErrorInfo) {
	totalNeedTime = 0
	var list []*struct {
		// Id 自增id
		Id uint64 `db:"id,omitempty" json:"id"`
		// NeedTime 工时的时间，单位分钟
		NeedTime uint32 `db:"need_time,omitempty" json:"needTime"`
	}
	oriErr := mysql.SelectAllByCond(consts.TableIssueWorkHours, cond, &list)
	if oriErr != nil {
		log.Error(oriErr)
		return totalNeedTime, errs.MysqlOperateError
	}
	total := uint32(0)
	for _, val := range list {
		total += val.NeedTime
	}
	return int64(total), nil
}

// 移动任务时，更新任务对应的工时的项目id。
func UpdateWorkHourProjectId(orgId, oldProjectId, newProjectId int64, issueIds []int64, isDelete bool, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	if oldProjectId == newProjectId {
		return nil
	}
	cond1 := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: oldProjectId,
		consts.TcIssueId:   db.In(issueIds),
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}
	params := bo.UpdateOneIssueWorkHoursBo{}
	if isDelete {
		// 如果需要删除工时
		params.IsDelete = consts.AppIsDeleted
	} else {
		params.ProjectId = newProjectId
	}
	return UpdateIssueWorkHourByCond(cond1, params, tx...)
}

// CheckHasWriteForWorkHourColumn 检查某个用户是否可以编辑工时列
func CheckHasWriteForWorkHourColumn(orgId, curUserId, proAppId int64, issueBo *bo.IssueBo) (bool, errs.SystemErrorInfo) {
	err := AuthIssueWithAppId(orgId, curUserId, issueBo, consts.RoleOperationPathOrgProIssueT, consts.OperationProIssue4Modify, proAppId, consts.ProBasicFieldWorkHour)
	if err != nil {
		//log.Errorf("[CheckHasWriteForWorkHourColumn] err: %v", err)
		return false, errs.BuildSystemErrorInfo(errs.Unauthorized, err)
	}

	return true, nil
}

// DeleteWorkHourForIssues 删除一批任务对应的所有工时记录
// 目前，有权删除任务，就视为有权删除任务对应的工时
func DeleteWorkHourForIssues(orgId int64, delIssueIds []int64, recycleVersionId int64) errs.SystemErrorInfo {
	upd := mysql.Upd{
		consts.TcVersion:  recycleVersionId,
		consts.TcIsDelete: consts.AppIsDeleted,
	}
	if err := DelWorkHours(orgId, delIssueIds, upd); err != nil {
		log.Errorf("[DeleteWorkHourForIssues] orgId: %d, err: %v", orgId, err)
		return err
	}

	return nil
}

// GetHasWorkHourRecordIssueIds 查询组织内，有工时记录的任务id
func GetHasWorkHourRecordIssueIds(orgId int64, projectIds []int64) ([]int64, errs.SystemErrorInfo) {
	resIssueIds := make([]int64, 0)
	list, err := GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
		OrgId:      orgId,
		ProjectIds: projectIds,
		IsDelete:   consts.AppIsNoDelete,
		Page:       1,
		Size:       10000,
	})
	if err != nil {
		log.Errorf("[] err: %v", err)
		return resIssueIds, err
	}
	for _, record := range list {
		resIssueIds = append(resIssueIds, record.IssueId)
	}

	return resIssueIds, nil
}

// GetWorkHourListForOneIssue 查询一个任务的现有的工时记录
func GetWorkHourListForOneIssue(orgId int64, issueId int64) ([]po.PpmPriIssueWorkHours, errs.SystemErrorInfo) {
	list, err := GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
		OrgId:    orgId,
		IssueId:  issueId,
		Types:    []uint8{consts2.WorkHourTypeSubPredict, consts2.WorkHourTypeActual},
		IsDelete: consts.AppIsNoDelete,
		Page:     1,
		Size:     2000,
	})
	if err != nil {
		log.Errorf("[GetWorkHourListForOneIssue] err: %v, issueId: %d", err, issueId)
		return nil, err
	}

	return list, nil
}

// 给复制任务用的
func GetCopyIssueWorkHours(orgId, projectId, issueId int64) ([]po.PpmPriIssueWorkHours, errs.SystemErrorInfo) {
	condBo := bo.IssueWorkHoursBoListCondBo{
		OrgId:     orgId,
		ProjectId: projectId,
		IssueId:   issueId,
	}
	workHoursList, err := GetIssueWorkHoursList(condBo)
	if err != nil {
		log.Errorf("[GetCopyIssueWorkHours] GetIssueWorkHoursList err:%v, orgId:%v, projectId:%v, issueId:%v",
			err, orgId, projectId, issueId)
		return nil, err
	}

	return workHoursList, nil
}
