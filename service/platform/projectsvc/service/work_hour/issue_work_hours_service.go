package work_hour

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/spf13/cast"

	msgPb "github.com/star-table/interface/golang/msg/v1"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/polaris-backend/facade/common/report"

	tablev1 "github.com/star-table/interface/golang/table/v1"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	types2 "github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	slice "github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/core/util/format"
	int64Slice "github.com/star-table/polaris-backend/common/core/util/slice/int64"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/commonfacade"
	"github.com/star-table/polaris-backend/facade/formfacade"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	consts2 "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"github.com/cznic/sortutil"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = *logger.GetDefaultLogger()

// GetIssueWorkHoursInfo 获取一个任务的工时信息
// 包含预估工时、实际工时记录列表等信息
func GetIssueWorkHoursInfo(orgId, currentUserId int64, param vo.GetIssueWorkHoursInfoReq) (*vo.GetIssueWorkHoursInfoResp, errs.SystemErrorInfo) {
	defaultResp := &vo.GetIssueWorkHoursInfoResp{
		SimplePredictWorkHour: &vo.OneWorkHourRecord{
			Worker: &vo.WorkHourWorker{},
		},
		PredictWorkHourList: []*vo.OneWorkHourRecord{},
		ActualWorkHourList:  []*vo.OneActualWorkHourRecord{},
		ActualNeedTimeTotal: "0",
	}
	// 校验是否开启工时功能
	issues, err := domain.GetIssueInfosLc(orgId, currentUserId, []int64{param.IssueID})
	if err != nil {
		log.Error(err)
		return defaultResp, err
	}
	if len(issues) < 1 {
		err = errs.IssueNotExist
		log.Error(err)
		return defaultResp, err
	}
	issue := issues[0]

	if issue.ProjectId <= 0 {
		log.Infof("[GetIssueWorkHoursInfo] 无归属项目的issueId:%d", param.IssueID)
		return defaultResp, nil
	}

	projectInfo, err := domain.GetProjectSimple(orgId, issue.ProjectId)
	if err != nil {
		log.Errorf("[GetIssueWorkHoursInfo] issueId: %d, GetProject err: %v", issue.Id, err)
		return defaultResp, err
	}

	// 查询工时记录
	allList, err := domain.GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
		OrgId:     orgId,
		ProjectId: issue.ProjectId,
		IssueId:   param.IssueID,
		IsDelete:  consts.AppIsNoDelete,
		Page:      1,
		Size:      10000, // 只是为了取出所有满足条件的记录，所以用较大的整数。
	})
	if err != nil {
		log.Errorf("[GetIssueWorkHoursInfo] issueId: %d, GetIssueWorkHoursList err: %v", issue.Id, err)
		return defaultResp, err
	}
	// 获取工作者id，查询对应的名称，头像
	var userIds []int64
	usersMap := map[int64]vo.WorkHourWorker{}
	for _, one := range allList {
		if one.WorkerId == 0 {
			continue
		}
		userIds = append(userIds, one.WorkerId, one.Creator)
	}
	userIds = slice.SliceUniqueInt64(userIds)
	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
	if err != nil {
		log.Error(err)
		return defaultResp, err
	}
	for _, user := range userInfos {
		usersMap[user.UserId] = vo.WorkHourWorker{
			UserID: user.UserId,
			Name:   user.Name,
			Avatar: user.Avatar,
		}
	}
	var predictList []*vo.OneWorkHourRecord
	var actualList []*vo.OneActualWorkHourRecord
	simplePredictWorkHour := &vo.OneWorkHourRecord{
		Worker: &vo.WorkHourWorker{},
	}
	totalNeedTimeInt := uint32(0)
	// 查询 currentUserId 对应的协作者角色，检查角色是否有编辑工时列的权限（字段权限）
	couldWrite, err := domain.CheckHasWriteForWorkHourColumn(orgId, currentUserId, projectInfo.AppId, issue)
	// 如果是没权限，这还可以往下执行
	if err != nil &&
		err.Code() != errs.NoOperationPermissionForProject.Code() &&
		err.Code() != errs.NoOperationPermissionForIssueUpdate.Code() &&
		err.Code() != errs.ProjectIsArchivedWhenModifyIssue.Code() {
		log.Errorf("[GetIssueWorkHoursInfo] issueId: %d, CheckHasWriteForWorkHourColumn err: %v", issue.Id, err)
		return defaultResp, err
	}
	// 查询项目负责人、任务负责人，参与者，关注人，用于检验是否有编辑权限
	//owners, _, _, err := GetPermissionUserIdsForWorkHour(orgId, param.IssueID)
	for _, one := range allList {
		worker := vo.WorkHourWorker{}
		if tmpWorker, ok := usersMap[one.WorkerId]; ok {
			worker = tmpWorker
		}
		record := vo.OneWorkHourRecord{
			ID:        int64(one.Id),
			Type:      int64(one.Type),
			Worker:    &worker,
			NeedTime:  format.FormatNeedTimeIntoString(int64(one.NeedTime)),
			StartTime: int64(one.StartTime),
			EndTime:   int64(one.EndTime),
			Desc:      one.Desc,
			IsEnable:  0,
		}
		// 如果是工时协作人，且有工时列的编辑权限，则可以编辑
		if couldWrite {
			record.IsEnable = 1
		} else if currentUserId == worker.UserID {
			// 如果当前用户是工时执行者本人，可以编辑
			record.IsEnable = 1
		}
		switch one.Type {
		case consts2.WorkHourTypeTotalPredict: // 总预估工时
			simplePredictWorkHour = &record
		case consts2.WorkHourTypeSubPredict: // 子预估工时
			predictList = append(predictList, &record)
		case consts2.WorkHourTypeActual:
			tmpCreator := vo.WorkHourWorker{}
			if user, ok := usersMap[one.Creator]; ok {
				tmpCreator = user
			}
			actualWorkHour := &vo.OneActualWorkHourRecord{
				ID:          record.ID,
				Type:        record.Type,
				Worker:      record.Worker,
				NeedTime:    record.NeedTime,
				StartTime:   record.StartTime,
				EndTime:     record.EndTime,
				CreatorName: tmpCreator.Name,
				CreateTime:  one.CreateTime.Unix(),
				Desc:        record.Desc,
				IsEnable:    record.IsEnable,
			}
			actualList = append(actualList, actualWorkHour)
			totalNeedTimeInt += one.NeedTime
		}
	}
	defaultResp.ActualNeedTimeTotal = format.FormatNeedTimeIntoString(int64(totalNeedTimeInt))
	defaultResp.SimplePredictWorkHour = simplePredictWorkHour
	if len(predictList) > 0 {
		defaultResp.PredictWorkHourList = predictList
	}
	if len(actualList) > 0 {
		defaultResp.ActualWorkHourList = actualList
	}

	return defaultResp, nil
}

// 介于如果根据不同情况去更新有可能会导致数据不准，所以每次获取最新数据库的值进行更新
func saveWorkHourInfoToLess(orgId, userId int64, issue *bo.IssueBo, tx sqlbuilder.Tx) error {
	projectInfo, appIdErr := domain.GetProjectSimple(orgId, issue.ProjectId, tx)
	if appIdErr != nil {
		log.Error(appIdErr)
		return appIdErr
	}

	// 查询工时记录
	allList, err := domain.GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
		OrgId:     orgId,
		ProjectId: issue.ProjectId,
		IssueId:   issue.Id,
		IsDelete:  consts.AppIsNoDelete,
		Page:      1,
		Size:      10000, // 只是为了取出所有满足条件的记录，所以用较大的整数。
	}, tx)
	if err != nil {
		log.Error(err)
		return err
	}

	predictHour := "0"
	totalNeedTimeInt := uint32(0)
	userIds := make([]string, 0, len(allList))
	for _, one := range allList {
		if one.WorkerId > 0 {
			userIds = append(userIds, fmt.Sprintf("U_%d", one.WorkerId))
		}

		switch one.Type {
		case consts2.WorkHourTypeTotalPredict: // 总预估工时
			predictHour = format.FormatNeedTimeIntoString(int64(one.NeedTime))
		case consts2.WorkHourTypeSubPredict: // 子预估工时
		case consts2.WorkHourTypeActual:
			totalNeedTimeInt += one.NeedTime
		}
	}
	actualNeedTimeTotal := format.FormatNeedTimeIntoString(int64(totalNeedTimeInt))
	userIds = slice.SliceUniqueString(userIds)
	lcUpd := map[string]interface{}{
		"issueId": issue.Id,
		consts.ProBasicFieldWorkHour: &formvo.LessWorkHour{
			PlanHour:        predictHour,
			ActualHour:      actualNeedTimeTotal,
			CollaboratorIds: userIds,
		},
	}
	lcResp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
		AppId:   projectInfo.AppId,
		OrgId:   orgId,
		UserId:  userId,
		TableId: issue.TableId,
		Form:    []map[string]interface{}{lcUpd},
	})
	if lcResp.Failure() {
		log.Errorf("[saveWorkHourInfoToLess] issue.Id: %d, LessUpdateIssue err: %v", issue.Id, lcResp.Error())
		return lcResp.Error()
	}

	return nil
}

func saveWorkHourInfoToLessWithIssueId(orgId, userId, issueId int64, tx sqlbuilder.Tx) error {
	// 查询 issue 信息
	issues, err := domain.GetIssueInfosLc(orgId, userId, []int64{issueId})
	if err != nil || len(issues) < 1 {
		log.Errorf("[saveWorkHourInfoToLessWithIssueId] GetIssueInfoBasic id:%v, err:%v", issueId, err)
		return errs.IssueNotExist
	}
	issue := issues[0]
	return saveWorkHourInfoToLess(orgId, userId, issue, tx)
}

// 工时记录列表
func GetIssueWorkHoursList(orgId, currentUserId int64, param vo.GetIssueWorkHoursListReq) (*vo.GetIssueWorkHoursListResp, errs.SystemErrorInfo) {
	// 过滤条件赋值，没有则用控制替代
	filterParam := bo.IssueWorkHoursBoListCondBo{
		Ids:      nil,
		OrgId:    orgId,
		IssueId:  param.IssueID,
		Type:     uint8(param.Type),
		IsDelete: consts.AppIsNoDelete,
		Page:     int(*param.Page),
		Size:     int(*param.Size),
	}
	list, err := domain.GetIssueWorkHoursList(filterParam)
	if err != nil {
		return nil, err
	}
	resList := &[]*vo.IssueWorkHours{}
	oriErr := copyer.Copy(list, resList)
	if oriErr != nil {
		log.Error(oriErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	resp := &vo.GetIssueWorkHoursListResp{
		Total: int64(len(*resList)),
		List:  *resList,
	}
	return resp, nil
}

func CheckNeedTimeValid(needTime string) errs.SystemErrorInfo {
	if !format.VerifyFloat1(needTime) {
		return errs.WorkHourNeedTimeInvalid
	}
	f, _ := strconv.ParseFloat(needTime, 64)
	if f < 0 {
		return errs.WorkHourNeedTimeInvalid
	}
	return nil
}

// 校验工时执行者id是否合法
func CheckWorkerIdValid(workerId int64) bool {
	if workerId < 1 {
		return false
	}
	return true
}

// CreateIssueWorkHours 创建新的工时记录，2种情况
// 1.issue 负责人创建预估工时或者实际工时
// 2.issue 参与人，创建自己的实际工时
func CreateIssueWorkHours(orgId, currentUserId int64, param vo.CreateIssueWorkHoursReq, needPush bool, tx ...sqlbuilder.Tx) (*vo.BoolResp, errs.SystemErrorInfo) {
	returnResp := &vo.BoolResp{IsTrue: false}
	reqVo := projectvo.CreateIssueWorkHoursReqVo{
		Input:  param,
		UserId: currentUserId,
		OrgId:  orgId,
	}
	if err := CheckNeedTimeValid(param.NeedTime); err != nil {
		return returnResp, err
	}
	oldIssue, err := domain.GetIssueInfoLc(orgId, currentUserId, param.IssueID)
	if err != nil {
		log.Error(err)
		return returnResp, err
	}
	if ok := CheckWorkerIdValid(param.WorkerID); !ok {
		return returnResp, errs.UserNotExist
	}
	if param.Desc == nil {
		defaultDesc := ""
		param.Desc = &defaultDesc
	}
	if param.EndTime == nil {
		defaultEndTime := int64(0)
		param.EndTime = &defaultEndTime
	}
	if ok, err := CheckPowerForCreateIssueWorkHours(orgId, currentUserId, oldIssue); !ok {
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return returnResp, errs.DenyCreateIssueWorkHours
	}
	// 如果新增的是总预估工时，校验是否存在
	if param.Type == consts2.WorkHourTypeTotalPredict {
		list, err := domain.GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
			OrgId:    orgId,
			IssueId:  param.IssueID,
			Type:     consts2.WorkHourTypeTotalPredict,
			IsDelete: consts.AppIsNoDelete,
		})
		if err != nil {
			log.Error(err)
			return returnResp, err
		}
		if len(list) > 0 {
			return returnResp, errs.SimplePredictIssueWorkHourExist
		}
	}
	// 如果是新增子预估工时，必须存在 workerId
	// 不允许新增子预估工时，新增子预估工时需要通过 CreateMultiIssueWorkHours 函数
	if param.Type == consts2.WorkHourTypeSubPredict {
		log.Error("不允许新增子预估工时，新增子预估工时需要通过 CreateMultiIssueWorkHours 函数")
		return returnResp, errs.ParamError
	}
	workHourNewId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableIssueWorkHours)
	if err != nil {
		log.Error(err)
		return returnResp, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	if reqVo.Input.ProjectID == nil {
		reqVo.Input.ProjectID = &oldIssue.ProjectId
	}
	var oriErr error
	if tx != nil && len(tx) > 0 {
		oriErr = domain.CreateIssueWorkHours(workHourNewId, reqVo, tx[0])
	} else {
		oriErr = mysql.TransX(func(tx sqlbuilder.Tx) error {
			err := domain.CreateIssueWorkHours(workHourNewId, reqVo, tx)
			if err != nil {
				return err
			}

			return saveWorkHourInfoToLess(orgId, currentUserId, oldIssue, tx)
		})
	}
	if oriErr != nil {
		log.Error(oriErr)
		return returnResp, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	// 增加动态
	dataFormat1 := vo.OneWorkHourRecord{}
	// copyer.Copy(param, &dataFormat1)
	dataFormat1.ID = workHourNewId
	dataFormat1.Type = param.Type
	dataFormat1.Worker = nil
	dataFormat1.NeedTime = param.NeedTime
	dataFormat1.StartTime = param.StartTime
	dataFormat1.EndTime = *param.EndTime
	dataFormat1.Desc = *param.Desc
	userInfoMap, err := GetWorkerInfos(orgId, []int64{param.WorkerID})
	if err != nil {
		log.Error(err)
		return returnResp, err
	}
	if val, ok := userInfoMap[param.WorkerID]; ok {
		dataFormat1.Worker = val
	}
	// 总预估工时 无需新增动态
	if param.Type != consts2.WorkHourTypeTotalPredict {
		TriggerCreateWorkHourTrend(orgId, currentUserId, param.IssueID, dataFormat1)
	}
	returnResp.IsTrue = true

	if needPush {
		asyn.Execute(func() {
			PushMqttForWorkHour(orgId, currentUserId, param.IssueID)
		})
	}

	return returnResp, nil
}

// CreateMultiIssueWorkHours 新增详细版预估工时（多个子预估工时）信息。新建多个子预估工时时，会生成一个总预估工时记录
func CreateMultiIssueWorkHours(orgId, currentUserId int64, param vo.CreateMultiIssueWorkHoursReq) (*vo.BoolResp, errs.SystemErrorInfo) {
	returnRes := &vo.BoolResp{
		IsTrue: false,
	}
	if param.TotalIssueWorkHourRecord == nil {
		log.Error("参数错误。新增详细预估工时时，需要有总预估工时数据。")
		return returnRes, errs.ParamError
	}
	// 查询 issue 信息
	issues, err := domain.GetIssueInfosLc(orgId, currentUserId, []int64{param.IssueID})
	if err != nil || len(issues) < 1 {
		return returnRes, errs.IssueNotExist
	}
	issue := issues[0]
	// 权限校验
	if ok, err := CheckPowerForCreateIssueWorkHours(orgId, currentUserId, issue); !ok {
		if err != nil {
			log.Error(err)
			return nil, err
		}
		return nil, errs.DenyCreateIssueWorkHours
	}
	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueWorkHours, len(param.PredictWorkHourList))
	if err != nil {
		log.Error(err)
		return returnRes, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	newDataList := []*po.PpmPriIssueWorkHours{}
	workerIds := []int64{}
	for index, v := range param.PredictWorkHourList {
		if ok := CheckWorkerIdValid(v.WorkerID); !ok {
			return returnRes, errs.UserNotExist
		}
		if err := CheckNeedTimeValid(v.NeedTime); err != nil {
			return returnRes, err
		}
		needTimeInt := format.FormatNeedTimeIntoSecondNumber(v.NeedTime)
		if needTimeInt > consts2.MaxWorkHour {
			return returnRes, errs.WorkHourMaxNeedTime
		}
		if *v.EndTime > 0 && !format.CheckTimeRangeTimeNumIsValid(needTimeInt, v.StartTime, *v.EndTime) {
			return returnRes, errs.WorkHourTimeRangeForNeedTimeInvalid
		}
		oneNewData := &po.PpmPriIssueWorkHours{
			Id:                uint64(ids.Ids[index].Id),
			OrgId:             orgId,
			ProjectId:         issue.ProjectId,
			IssueId:           param.IssueID,
			Type:              consts2.WorkHourTypeSubPredict,
			WorkerId:          v.WorkerID,
			NeedTime:          uint32(needTimeInt),
			RemainTimeCalType: consts2.RemainTimeCalTypeDefault,
			RemainTime:        0,
			StartTime:         uint32(v.StartTime),
			EndTime:           uint32(*v.EndTime),
			Desc:              "",
			Creator:           currentUserId,
			CreateTime:        time.Time{},
			Updator:           currentUserId,
			UpdateTime:        time.Time{},
			Version:           1,
			IsDelete:          consts.AppIsNoDelete,
		}
		newDataList = append(newDataList, oneNewData)
		if v.WorkerID > 0 {
			workerIds = append(workerIds, v.WorkerID)
		}
	}

	// 校验总预估工时
	if !format.VerifyFloat1(param.TotalIssueWorkHourRecord.NeedTime) {
		return returnRes, errs.WorkHourNeedTimeInvalid
	}
	var oriErr error
	oriErr = mysql.TransX(func(tx sqlbuilder.Tx) error {
		if param.TotalIssueWorkHourRecord.StartTime <= 0 {
			param.TotalIssueWorkHourRecord.StartTime = time.Time(issue.PlanStartTime).Unix()
		}
		// 新增一条总预估工时记录
		initDesc := ""
		ownerIds := []int64{}
		//for _, own := range issue.OwnerId {
		//	id, errParse := strconv.ParseInt(own[2:], 10, 64)
		//	if errParse != nil {
		//		log.Error(errParse)
		//		return errs.TypeConvertError
		//	}
		//	ownerIds = append(ownerIds, id)
		//}
		if param.TotalIssueWorkHourRecord.WorkerID > 0 {
			ownerIds = append(ownerIds, param.TotalIssueWorkHourRecord.WorkerID)
		}
		totalWorkHourWorkId := int64(0)
		if len(ownerIds) > 0 {
			totalWorkHourWorkId = ownerIds[0]
		}
		_, err = CreateIssueWorkHours(orgId, currentUserId, vo.CreateIssueWorkHoursReq{
			ProjectID: &issue.ProjectId,
			IssueID:   param.IssueID,
			Type:      consts2.WorkHourTypeTotalPredict,
			// 旧版；总预估工时的执行者是任务负责人；新版：总预估工时执行人为 0
			WorkerID:  totalWorkHourWorkId,
			NeedTime:  param.TotalIssueWorkHourRecord.NeedTime,
			StartTime: param.TotalIssueWorkHourRecord.StartTime,
			EndTime:   param.TotalIssueWorkHourRecord.EndTime,
			Desc:      &initDesc,
		}, false, tx)
		if err != nil {
			log.Error(err)
			return err
		}
		if len(newDataList) > 0 {
			err = domain.CreateMultiIssueWorkHours(newDataList, tx)
			if err != nil {
				log.Error(err)
				return err
			}
		}

		return saveWorkHourInfoToLess(orgId, currentUserId, issue, tx)
	})
	if oriErr != nil {
		log.Error(oriErr)
		return returnRes, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	dataListFormat1 := []vo.OneWorkHourRecord{}
	// 赋值到新的结构体
	for _, val := range newDataList {
		tmpVal := vo.OneWorkHourRecord{
			ID:        int64(val.Id),
			Type:      int64(val.Type),
			Worker:    nil,
			NeedTime:  format.FormatNeedTimeIntoString(int64(val.NeedTime)),
			StartTime: int64(val.StartTime),
			EndTime:   int64(val.EndTime),
			Desc:      val.Desc,
			IsEnable:  -1,
		}
		dataListFormat1 = append(dataListFormat1, tmpVal)
	}
	// 查询 worker 信息，用于生成动态
	userInfoMap, err := GetWorkerInfos(orgId, workerIds)
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	//for index, oneWorkHour := range dataListFormat1 {
	//	if val, ok := userInfoMap[newDataList[index].WorkerId]; ok {
	//		oneWorkHour.Worker = val
	//	} else {
	//		log.Error("【异常】创建工时的执行者的信息找不到。")
	//		return returnRes, errs.UserNotFoundError
	//	}
	//	// 增加多条时，触发新增动态记录
	//	TriggerCreateWorkHourTrend(orgId, currentUserId, param.IssueID, oneWorkHour)
	//}

	returnRes.IsTrue = true

	asyn.Execute(func() {
		for index, oneWorkHour := range dataListFormat1 {
			if val, ok := userInfoMap[newDataList[index].WorkerId]; ok {
				oneWorkHour.Worker = val
			} else {
				log.Error("【异常】创建工时的执行者的信息找不到。")
				return
			}
			// 增加多条时，触发新增动态记录
			TriggerCreateWorkHourTrend(orgId, currentUserId, param.IssueID, oneWorkHour)
		}
	})

	asyn.Execute(func() {
		PushMqttForWorkHour(orgId, currentUserId, param.IssueID)
	})

	return returnRes, nil
}

// AddUserIntoProParticipant 将一个用户加为项目参与者
func AddUserIntoProParticipant(orgId int64, opUid, targetUid int64, projectId int64) errs.SystemErrorInfo {
	// 检查是否是项目参与人
	proUserIds, oriErr1 := domain.GetProjectParticipantIds(orgId, projectId)
	if oriErr1 != nil {
		log.Error(oriErr1)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr1)
	}
	if hasIt, _ := slice.Contain(proUserIds, targetUid); hasIt {
		return nil
	}
	needAddUidsForPro := []int64{targetUid}
	if len(needAddUidsForPro) > 0 {
		memberIdStrArr := make([]string, 0)
		for _, oneId := range needAddUidsForPro {
			if oneId < 1 {
				continue
			}
			memberIdStrArr = append(memberIdStrArr, fmt.Sprintf("U_%d", oneId))
		}
		// 校验成员 id 的合法性
		verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, needAddUidsForPro)
		if !verifyOrgUserFlag {
			err := errs.BuildSystemErrorInfoWithMessage(errs.VerifyOrgError, "[AddUserIntoPro] 用户、组织校验失败")
			log.Error(err)
			return err
		}
		// 如果要新增项目成员，需要有此权限才行
		perErr := domain.AuthProject(orgId, opUid, projectId, consts.RoleOperationPathOrgProProConfig, consts.OperationProConfigModify)
		if perErr != nil {
			log.Error(perErr)
			return errs.AddUserToProButNoPower
		}

		memberIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(needAddUidsForPro))
		if idErr != nil {
			log.Error(idErr)
			return idErr
		}
		memberAdd := make([]interface{}, 0)
		for i, uid := range needAddUidsForPro {
			memberAdd = append(memberAdd, po.PpmProProjectRelation{
				Id:           memberIds.Ids[i].Id,
				OrgId:        orgId,
				ProjectId:    projectId,
				RelationId:   uid,
				RelationType: consts.IssueRelationTypeParticipant,
				Creator:      opUid,
				CreateTime:   time.Now(),
				IsDelete:     consts.AppIsNoDelete,
				Status:       consts.ProjectMemberEffective,
				Updator:      opUid,
				UpdateTime:   time.Now(),
				Version:      1,
			})
		}
		err := domain.PaginationInsert(memberAdd, &po.PpmProProjectRelation{})
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}
	return nil
}

// UpdateIssueWorkHours 编辑工时记录
// needUpdateTotal 是否需要更新总预估工时。当一个请求中，只更新一个工时记录时，传 true。
// needTrend 是否需要产生动态
func UpdateIssueWorkHours(orgId, currentUserId int64, param vo.UpdateIssueWorkHoursReq, needCheckPermission, needUpdateTotal, needTrend, needPush bool, tx ...sqlbuilder.Tx) (*vo.BoolResp, errs.SystemErrorInfo) {
	returnRes := &vo.BoolResp{
		IsTrue: false,
	}
	if ok := CheckWorkerIdValid(param.WorkerID); !ok {
		return returnRes, errs.UserNotExist
	}
	if param.Desc == nil {
		defaultDesc := ""
		param.Desc = &defaultDesc
	}
	if err := CheckNeedTimeValid(param.NeedTime); err != nil {
		return returnRes, err
	}
	needTimeInt := format.FormatNeedTimeIntoSecondNumber(param.NeedTime)
	if needTimeInt > consts2.MaxWorkHour {
		return returnRes, errs.WorkHourMaxNeedTime
	}
	paramBo := bo.UpdateOneIssueWorkHoursBo{
		Id:        uint64(param.IssueWorkHoursID),
		OrgId:     orgId,
		WorkerId:  param.WorkerID,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		NeedTime:  uint32(needTimeInt),
		Desc:      *param.Desc,
		IsDelete:  consts.AppIsNoDelete,
	}
	// 必须有截止时间，才会校验起止时间的区间的时间是否合法。
	if param.EndTime > 0 && !format.CheckTimeRangeTimeNumIsValid(needTimeInt, param.StartTime, param.EndTime) {
		return returnRes, errs.WorkHourTimeRangeForNeedTimeInvalid
	}

	// 工时记录是否存在
	oldIssueWorkHour, err := domain.GetOneIssueWorkHoursById(int64(paramBo.Id))
	if oldIssueWorkHour == nil {
		return returnRes, errs.IssueWorkHourNotExist
	}
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	// 检查是否需要更新，如果没有变更，则无需更新
	needUpdateWorkHour := true
	if format.TransformTimeStampToDate(int64(oldIssueWorkHour.StartTime)) == format.TransformTimeStampToDate(paramBo.StartTime) &&
		format.TransformTimeStampToDate(int64(oldIssueWorkHour.EndTime)) == format.TransformTimeStampToDate(paramBo.EndTime) &&
		oldIssueWorkHour.WorkerId == paramBo.WorkerId &&
		oldIssueWorkHour.NeedTime == paramBo.NeedTime &&
		oldIssueWorkHour.Desc == paramBo.Desc {
		log.Infof("[UpdateIssueWorkHours] 该工时没有更改任何信息，无需执行更新。workHourId: %d", oldIssueWorkHour.Id)
		needUpdateWorkHour = false
	}
	// 校验是否有权限进行编辑操作
	if needCheckPermission {
		hasPower, err := CheckHasWorkHoursUpdatePower(orgId, bo.CheckWorkHoursUpdatePowerBo{
			CurrentUserId:    currentUserId,
			IssueId:          oldIssueWorkHour.IssueId,
			IssueWorkHoursId: paramBo.Id,
			WorkerId:         oldIssueWorkHour.WorkerId,
		})
		if err != nil {
			return returnRes, err
		}
		if !hasPower {
			return returnRes, errs.DenyUpdateIssueWorkHours
		}
	}
	// 触发生成动态 start
	oldValue := vo.OneWorkHourRecord{
		ID:        int64(oldIssueWorkHour.Id),
		Type:      int64(oldIssueWorkHour.Type),
		Worker:    nil,
		NeedTime:  format.FormatNeedTimeIntoString(int64(oldIssueWorkHour.NeedTime)),
		StartTime: int64(oldIssueWorkHour.StartTime),
		EndTime:   int64(oldIssueWorkHour.EndTime),
		Desc:      oldIssueWorkHour.Desc,
	}
	oldValue.StartTime = int64(oldIssueWorkHour.StartTime)
	oldValue.EndTime = int64(oldIssueWorkHour.EndTime)
	newValue := vo.OneWorkHourRecord{
		ID:        param.IssueWorkHoursID,
		Type:      int64(oldIssueWorkHour.Type), // 不支持更新 type，因此还是之前的旧 type 值。
		Worker:    nil,
		NeedTime:  param.NeedTime,
		StartTime: param.StartTime,
		EndTime:   param.EndTime,
		Desc:      *param.Desc,
	}
	// 查询 worker
	var workerIds []int64
	if param.WorkerID == 0 {
		workerIds = []int64{oldIssueWorkHour.WorkerId}
		param.WorkerID = oldIssueWorkHour.WorkerId
	} else {
		workerIds = slice.SliceUniqueInt64([]int64{oldIssueWorkHour.WorkerId, param.WorkerID})
	}
	userInfoMap, err := GetWorkerInfos(orgId, workerIds)
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	if val, ok := userInfoMap[oldIssueWorkHour.WorkerId]; ok {
		oldValue.Worker = val
	} else {
		log.Error("[UpdateIssueWorkHours] 更新前的工时执行者的信息找不到。")
		return returnRes, errs.UserNotFoundError
	}
	if val, ok := userInfoMap[param.WorkerID]; ok {
		newValue.Worker = val
	} else {
		log.Error("[UpdateIssueWorkHours] 更新后的工时执行者的信息找不到。")
		return returnRes, errs.UserNotFoundError
	}
	// 执行更新
	// 如果更新的总预估工时、实际工时，则直接更新
	// 如果更新的是子预估工时，则还需更改总预估工时
	if tx != nil && len(tx) > 0 {
		err = domain.UpdateIssueWorkHours(paramBo, tx[0])
	} else {
		oriErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
			err = domain.UpdateIssueWorkHours(paramBo, tx)
			if err != nil {
				return err
			}
			return saveWorkHourInfoToLessWithIssueId(orgId, currentUserId, oldIssueWorkHour.IssueId, tx)
		})
		if oriErr != nil {
			log.Errorf("[UpdateIssueWorkHours] error:%v", oriErr)
			return returnRes, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
		}
	}
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	// 需要更新总预估工时，并且是在更新子预估工时的情况下，才会更新总预估工时
	if needUpdateTotal && oldIssueWorkHour.Type == consts2.WorkHourTypeSubPredict {
		// 查询总预估工时
		predictTotalNeedTime, err := domain.GetSumWorkHourSumNeedTime(db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIssueId:  oldIssueWorkHour.IssueId,
			consts.TcType:     consts2.WorkHourTypeSubPredict,
			consts.TcIsDelete: consts.AppIsNoDelete,
		})
		if err != nil {
			log.Error(err)
			return returnRes, err
		}
		// 更新总预估工时记录（简单工时）
		condition := db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIssueId:  oldIssueWorkHour.IssueId,
			consts.TcType:     consts2.WorkHourTypeTotalPredict,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}
		err = domain.UpdateIssueWorkHourByCond(condition, bo.UpdateOneIssueWorkHoursBo{
			NeedTime: uint32(predictTotalNeedTime),
		})
		if err != nil {
			log.Error(err)
			return returnRes, err
		}
	}
	// 一般更新时，需要产生动态。
	// 但还需兼容，在编辑预估工时时，会产生子预估工时动态，而**不会产生**总预估工时的动态。
	if needTrend && needUpdateWorkHour {
		TriggerUpdateWorkHourTrend(orgId, currentUserId, oldIssueWorkHour.IssueId, oldValue, newValue)
	}
	returnRes.IsTrue = true

	if needPush {
		asyn.Execute(func() {
			PushMqttForWorkHour(orgId, currentUserId, oldIssueWorkHour.IssueId)
		})
	}

	return returnRes, nil
}

// UpdateMultiIssueWorkHourWithDelete 编辑多条预估工时，支持 diff 后的删除操作。
func UpdateMultiIssueWorkHourWithDelete(orgId, currentUserId int64, param vo.UpdateMultiIssueWorkHoursReq) (*vo.BoolResp, errs.SystemErrorInfo) {
	returnRes := &vo.BoolResp{
		IsTrue: false,
	}
	if param.TotalIssueWorkHourRecord == nil {
		log.Error("缺少总预估工时记录。")
		return returnRes, errs.ParamError
	}
	// 检查是否开启了工时功能
	issues, err := domain.GetIssueInfosLc(orgId, currentUserId, []int64{param.IssueID})
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	if len(issues) < 1 {
		return returnRes, errs.IssueNotExist
	}
	issue := issues[0]
	if issue.ProjectId == 0 {
		err = errs.DenyUpdateWorkHourForIssueHashNoPro
		// 暂不支持向没有项目归属的任务更新工时
		log.Error(err)
		return returnRes, err
	}

	insertData, updateData, deleteIds, err := GetUpdateData(orgId, currentUserId, param, issue)
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	var oriErr error
	oriErr = mysql.TransX(func(tx sqlbuilder.Tx) error {
		// 更新处理。UpdateIssueWorkHours 中会产生动态
		for _, oneWorkHour := range updateData {
			if _, err := UpdateIssueWorkHours(orgId, currentUserId, *oneWorkHour, false, false, true, false, tx); err != nil {
				log.Error(err)
				return err
			}
		}
		// 新增处理
		err = domain.CreateMultiIssueWorkHours(insertData, tx)
		if err != nil {
			log.Error(err)
			return err
		}
		// 根据多条新增的数据产生动态
		if err := generateTrendForCreate(orgId, currentUserId, param.IssueID, insertData); err != nil {
			log.Error(err)
			return err
		}
		// 删除。DeleteIssueWorkHours 中会产生动态。
		if len(deleteIds) > 0 {
			err = domain.UpdateIssueWorkHourByCond(db.Cond{
				consts.TcId: db.In(deleteIds),
			}, bo.UpdateOneIssueWorkHoursBo{
				IsDelete: consts.AppIsDeleted,
			}, tx)
			if err != nil {
				return err
			}
			err = generateTrendForDelete(orgId, currentUserId, param.IssueID, deleteIds)
			if err != nil {
				return err
			}
		}
		// 总预估工时更新。
		// 一般更新时，需要产生动态。
		// 但还需兼容，在编辑预估工时时，会产生子预估工时动态。有子预估工时时，总预估工时的动态无需产生。
		needTrend1 := false
		if len(param.IssueWorkHourRecords) < 1 {
			needTrend1 = true
		}
		if _, err := UpdateIssueWorkHours(orgId, currentUserId, vo.UpdateIssueWorkHoursReq{
			IssueWorkHoursID: param.TotalIssueWorkHourRecord.ID,
			NeedTime:         param.TotalIssueWorkHourRecord.NeedTime,
			WorkerID:         param.TotalIssueWorkHourRecord.WorkerID,
			StartTime:        param.TotalIssueWorkHourRecord.StartTime,
			EndTime:          param.TotalIssueWorkHourRecord.EndTime,
			Desc:             param.TotalIssueWorkHourRecord.Desc,
		}, false, false, needTrend1, false, tx); err != nil {
			log.Error(err)
			return err
		}

		return saveWorkHourInfoToLess(orgId, currentUserId, issue, tx)
	})
	if oriErr != nil {
		log.Error(oriErr)
		return returnRes, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}

	asyn.Execute(func() {
		PushMqttForWorkHour(orgId, currentUserId, issue.Id)
	})

	returnRes.IsTrue = true

	return returnRes, nil
}

// 对新增的多条工时，生成动态
func generateTrendForCreate(orgId, opUserId, issueId int64, newDataList []*po.PpmPriIssueWorkHours) errs.SystemErrorInfo {
	if len(newDataList) < 1 {
		return nil
	}
	workerIds := make([]int64, len(newDataList))
	for i, item := range newDataList {
		workerIds[i] = item.WorkerId
	}
	workerIds = slice.SliceUniqueInt64(workerIds)
	userInfoMap, err := GetWorkerInfos(orgId, workerIds)
	if err != nil {
		return err
	}
	dataListFormat1 := []*vo.OneWorkHourRecord{}
	// 赋值到新的结构体
	for _, val := range newDataList {
		tmpVal := &vo.OneWorkHourRecord{
			ID:        int64(val.Id),
			Type:      int64(val.Type),
			Worker:    nil,
			NeedTime:  format.FormatNeedTimeIntoString(int64(val.NeedTime)),
			StartTime: int64(val.StartTime),
			EndTime:   int64(val.EndTime),
			Desc:      val.Desc,
			IsEnable:  -1,
		}
		dataListFormat1 = append(dataListFormat1, tmpVal)
	}
	for i, item := range dataListFormat1 {
		tmpWorkerId := newDataList[i].WorkerId
		if worker, ok := userInfoMap[tmpWorkerId]; ok {
			item.Worker = worker
			TriggerCreateWorkHourTrend(orgId, opUserId, issueId, *item)
		}
	}
	return nil
}

// 对删除多条工时，生成动态
func generateTrendForDelete(orgId, opUserId, issueId int64, ids []int64) errs.SystemErrorInfo {
	// 通过 ids 查询具体的工时信息
	ids2 := []uint64{}
	copyer.Copy(ids, &ids2)
	list, err := domain.GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
		OrgId: orgId,
		Ids:   ids2,
		Page:  int(1),
		Size:  len(ids),
	})
	if len(list) < 1 {
		return errs.IssueWorkHourNotExist
	}
	if err != nil {
		return err
	}
	workerIds := make([]int64, len(list))
	for i, item := range list {
		workerIds[i] = item.WorkerId
	}
	workerIds = slice.SliceUniqueInt64(workerIds)
	userInfoMap, err := GetWorkerInfos(orgId, workerIds)
	if err != nil {
		return err
	}
	// dataListFormat1 := []*vo.OneWorkHourRecord{}
	// copyer.Copy(list, &dataListFormat1)
	dataListFormat1 := []*vo.OneWorkHourRecord{}
	// 赋值到新的结构体
	for _, val := range list {
		tmpVal := &vo.OneWorkHourRecord{
			ID:        int64(val.Id),
			Type:      int64(val.Type),
			Worker:    nil,
			NeedTime:  format.FormatNeedTimeIntoString(int64(val.NeedTime)),
			StartTime: int64(val.StartTime),
			EndTime:   int64(val.EndTime),
			Desc:      val.Desc,
			IsEnable:  -1,
		}
		dataListFormat1 = append(dataListFormat1, tmpVal)
	}
	for i, item := range dataListFormat1 {
		tmpWorkerId := list[i].WorkerId
		if worker, ok := userInfoMap[tmpWorkerId]; ok {
			item.Worker = worker
			TriggerDeleteWorkHourTrend(orgId, opUserId, issueId, *item)
		}
	}
	return nil
}

// 摘取出需要更新的一些信息。编辑预估工时时调用！
func GetUpdateData(orgId, currentUserId int64, param vo.UpdateMultiIssueWorkHoursReq, issue *bo.IssueBo) (newData []*po.PpmPriIssueWorkHours, updateData []*vo.UpdateIssueWorkHoursReq, deleteIds []int64, err errs.SystemErrorInfo) {
	// 查询表中现有的工时记录id
	poList, err := domain.GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
		OrgId:   orgId,
		IssueId: param.IssueID,
		Type:    consts2.WorkHourTypeSubPredict,
		//Types:    []uint8{consts2.WorkHourTypeSubPredict, consts2.WorkHourTypeTotalPredict},
		IsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		log.Error(err)
		return
	}
	workHourMap := map[int64]po.PpmPriIssueWorkHours{}
	oldIds := []int64{}
	for _, workHourPo := range poList {
		oldIds = append(oldIds, int64(workHourPo.Id))
		workHourMap[int64(workHourPo.Id)] = workHourPo
	}
	updateIds := []int64{}
	defaultDesc := ""
	for _, oneWorkHour := range param.IssueWorkHourRecords {
		if ok := CheckWorkerIdValid(oneWorkHour.WorkerID); !ok {
			err = errs.UserNotExist
			return
		}
		if oneWorkHour.Desc == nil {
			oneWorkHour.Desc = &defaultDesc
		}
		if tmpErr := CheckNeedTimeValid(oneWorkHour.NeedTime); tmpErr != nil {
			err = tmpErr
			return
		}
		needTimeInt := format.FormatNeedTimeIntoSecondNumber(oneWorkHour.NeedTime)
		if needTimeInt > consts2.MaxWorkHour {
			err = errs.WorkHourMaxNeedTime
			return
		}
		if oneWorkHour.EndTime > 0 && !format.CheckTimeRangeTimeNumIsValid(needTimeInt, oneWorkHour.StartTime, oneWorkHour.EndTime) {
			err = errs.WorkHourTimeRangeForNeedTimeInvalid
			return
		}
		if oneWorkHour.ID == 0 {
			newData = append(newData, &po.PpmPriIssueWorkHours{
				Id:                0, // uint64(ids.Ids[index].Id)
				OrgId:             orgId,
				ProjectId:         issue.ProjectId,
				IssueId:           param.IssueID,
				Type:              consts2.WorkHourTypeSubPredict,
				WorkerId:          oneWorkHour.WorkerID,
				NeedTime:          uint32(needTimeInt),
				RemainTimeCalType: consts2.RemainTimeCalTypeDefault,
				RemainTime:        0,
				StartTime:         uint32(oneWorkHour.StartTime),
				EndTime:           uint32(oneWorkHour.EndTime),
				Desc:              "",
				Creator:           currentUserId,
				CreateTime:        time.Time{},
				Updator:           currentUserId,
				UpdateTime:        time.Time{},
				Version:           1,
				IsDelete:          consts.AppIsNoDelete,
			})
		} else {
			updateIssueWorkHour := workHourMap[oneWorkHour.ID]
			updateIds = append(updateIds, oneWorkHour.ID)
			// 更新时 过滤掉需要没有变化的工时参数
			if oneWorkHour.Type == int64(updateIssueWorkHour.Type) {
				if oneWorkHour.WorkerID == updateIssueWorkHour.WorkerId &&
					format.FormatNeedTimeIntoSecondNumber(oneWorkHour.NeedTime) == int64(updateIssueWorkHour.NeedTime) &&
					format.TransformTimeStampToDate(oneWorkHour.StartTime) == format.TransformTimeStampToDate(int64(updateIssueWorkHour.StartTime)) &&
					format.TransformTimeStampToDate(oneWorkHour.EndTime) == format.TransformTimeStampToDate(int64(updateIssueWorkHour.EndTime)) &&
					*oneWorkHour.Desc == updateIssueWorkHour.Desc {
					continue
				}
				tmpData := &vo.UpdateIssueWorkHoursReq{
					IssueWorkHoursID: oneWorkHour.ID,
					NeedTime:         oneWorkHour.NeedTime,
					WorkerID:         oneWorkHour.WorkerID,
					StartTime:        oneWorkHour.StartTime,
					EndTime:          oneWorkHour.EndTime,
					Desc:             oneWorkHour.Desc,
				}
				updateData = append(updateData, tmpData)
			}
		}
	}
	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableIssueWorkHours, len(newData))
	if err != nil {
		log.Error(err)
		return
	}
	for index := 0; index < len(newData); index++ {
		newData[index].Id = uint64(ids.Ids[index].Id)
	}
	// 在 oldIds 中，但不在 updateIds 中的 id，就是要删除的数据id
	deleteIds = int64Slice.ArrayDiff(oldIds, updateIds)
	return
}

// DeleteIssueWorkHours 删除工时记录
// 只有自己和负责人类型的人可以删除工时
// 删除子预估工时，会削减总预估工时的时长
func DeleteIssueWorkHours(orgId, currentUserId int64, param vo.DeleteIssueWorkHoursReq, needUpdateTotal, needCheckPermission bool) (*vo.BoolResp, errs.SystemErrorInfo) {
	returnRes := &vo.BoolResp{
		IsTrue: false,
	}
	// 删除可以直接调用更新的方法
	updateData := bo.UpdateOneIssueWorkHoursBo{
		Id:       uint64(param.IssueWorkHoursID),
		IsDelete: consts.AppIsDeleted,
	}
	// 工时记录是否存在
	oneIssueWorkHours, err := domain.GetOneIssueWorkHoursById(int64(updateData.Id))
	if oneIssueWorkHours == nil {
		return returnRes, errs.IssueWorkHourNotExist
	}
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	// 如果是删除总工时记录，则检查它是否存在总预估工时和实际预估工时，如果存在，则不允许删除
	if oneIssueWorkHours.Type == consts2.WorkHourTypeTotalPredict {
		list, err := domain.GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
			OrgId:    orgId,
			IssueId:  oneIssueWorkHours.IssueId,
			Types:    []uint8{consts2.WorkHourTypeSubPredict, consts2.WorkHourTypeActual},
			IsDelete: consts.AppIsNoDelete,
		})
		if err != nil {
			log.Error(err)
			return returnRes, err
		}
		if len(list) > 0 {
			log.Error(errs.WorkHourHasSubRecordDisableDel.Error())
			return returnRes, errs.WorkHourHasSubRecordDisableDel
		}
	}
	// 校验是否有权限进行编辑操作
	if needCheckPermission {
		hasPower, err := CheckHasWorkHoursUpdatePower(orgId, bo.CheckWorkHoursUpdatePowerBo{
			CurrentUserId:    currentUserId,
			IssueId:          oneIssueWorkHours.IssueId,
			IssueWorkHoursId: updateData.Id,
			WorkerId:         oneIssueWorkHours.WorkerId,
		})
		if err != nil {
			return returnRes, err
		}
		if !hasPower {
			return returnRes, errs.DenyUpdateIssueWorkHours
		}
	}
	// 触发生成动态 start
	oldValue := vo.OneWorkHourRecord{}
	oldValue.ID = int64(oneIssueWorkHours.Id)
	oldValue.Type = int64(oneIssueWorkHours.Type)
	oldValue.Worker = nil
	oldValue.NeedTime = format.FormatNeedTimeIntoString(int64(oneIssueWorkHours.NeedTime))
	oldValue.StartTime = int64(oneIssueWorkHours.StartTime)
	oldValue.EndTime = int64(oneIssueWorkHours.EndTime)
	oldValue.Desc = oneIssueWorkHours.Desc
	// 查询 worker
	userInfoMap, err := GetWorkerInfos(orgId, []int64{oneIssueWorkHours.WorkerId})
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	if val, ok := userInfoMap[oneIssueWorkHours.WorkerId]; ok {
		oldValue.Worker = val
	}
	// 查询总预估时长，用于更新总预估时长
	// 更新实际工时时，无需更新总预估工时
	totalPredictWorkHourNeedTime := uint32(0)
	if needUpdateTotal && oneIssueWorkHours.Type != consts2.WorkHourTypeActual {
		tmpList, err := domain.GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
			OrgId:    orgId,
			IssueId:  oneIssueWorkHours.IssueId,
			Type:     consts2.WorkHourTypeTotalPredict,
			IsDelete: consts.AppIsNoDelete,
		})
		if err != nil {
			log.Error(err)
			return returnRes, err
		}
		if len(tmpList) < 1 {
			log.Error(err)
			return returnRes, errs.IssueWorkHourNotExist
		}
		totalPredictWorkHourNeedTime = tmpList[0].NeedTime
	}

	// 执行删除更新
	var oriErr error
	oriErr = mysql.TransX(func(tx sqlbuilder.Tx) error {
		err = domain.UpdateIssueWorkHours(updateData, tx)
		if err != nil {
			log.Error(err)
			return err
		}
		// 删除子预估工时，需更新总预估工时。
		if needUpdateTotal && oneIssueWorkHours.Type == consts2.WorkHourTypeSubPredict {
			// 更新总预估时长
			condition := db.Cond{
				consts.TcOrgId:    orgId,
				consts.TcIssueId:  oneIssueWorkHours.IssueId,
				consts.TcType:     consts2.WorkHourTypeTotalPredict,
				consts.TcIsDelete: consts.AppIsNoDelete,
			}
			err := domain.UpdateIssueWorkHourByCond(condition, bo.UpdateOneIssueWorkHoursBo{
				NeedTime: totalPredictWorkHourNeedTime - oneIssueWorkHours.NeedTime,
			}, tx)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		return saveWorkHourInfoToLessWithIssueId(orgId, currentUserId, oneIssueWorkHours.IssueId, tx)
	})
	if oriErr != nil {
		log.Error(oriErr)
		return returnRes, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	TriggerDeleteWorkHourTrend(orgId, currentUserId, oneIssueWorkHours.IssueId, oldValue)
	returnRes.IsTrue = true

	asyn.Execute(func() {
		PushMqttForWorkHour(orgId, currentUserId, oneIssueWorkHours.IssueId)
	})

	return returnRes, nil
}

// CheckIsEnableWorkHour 引入工时字段后，默认就是开启的。增加工时列后，就能添加工时记录。
func CheckIsEnableWorkHour(orgId, currentUserId int64, param vo.CheckIsEnableWorkHourReq) (*vo.CheckIsEnableWorkHourResp, errs.SystemErrorInfo) {
	returnRes := &vo.CheckIsEnableWorkHourResp{
		IsEnable: true,
	}
	// “未放入项目”默认不开启工时。
	if param.ProjectID == int64(0) {
		returnRes.IsEnable = false
		return returnRes, nil
	}

	return returnRes, nil
}

// 启用/关闭工时记录功能
// 通过更新项目详情中的 is_enable_work_hours 判断是否开启工时功能
// 只有项目负责人可以开启/关闭工时功能
func DisOrEnableIssueWorkHours(orgId, currentUserId int64, param vo.DisOrEnableIssueWorkHoursReq) (*vo.BoolResp, errs.SystemErrorInfo) {
	returnRes := &vo.BoolResp{
		IsTrue: false,
	}
	project, err := domain.GetProject(orgId, param.ProjectID)
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	// 检查是否有项目的管理权限。有权创建项目，则意味着是管理员。
	err = domain.AuthProject(orgId, currentUserId, param.ProjectID, consts.RoleOperationPathOrgProIssueT, consts.OperationProConfigModify)
	if err != nil {
		// err 不为空时，则继续下一步的判断是否有权限
		if project.Owner != currentUserId {
			log.Error("无权启用、关闭工时功能。DenyEnableFuncIssueWorkHours。")
			return returnRes, errs.DenyEnableFuncIssueWorkHours
		}
	}

	if project == nil {
		log.Error("DisOrEnableIssueWorkHours： 项目不存在。")
		return returnRes, errs.ProjectNotExist
	}
	projectDetail, err := domain.GetProjectDetailByProjectIdBo(param.ProjectID, orgId)
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	err = domain.UpdateProjectDetail(&bo.ProjectDetailBo{
		Id:                projectDetail.Id,
		OrgId:             orgId,
		ProjectId:         projectDetail.ProjectId,
		IsEnableWorkHours: int(param.Enable),
		Updator:           currentUserId,
	})
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	// 生成动态
	workHourTrendTypeStr := ""
	if param.Enable == 1 {
		workHourTrendTypeStr = "enable"
	} else {
		workHourTrendTypeStr = "disable"
	}
	TriggerEnableOrDisableWorkHourTrend(orgId, currentUserId, project, workHourTrendTypeStr)
	returnRes.IsTrue = true
	return returnRes, nil
}

// CheckHasWorkHoursUpdatePower 检查是否有权限编辑工时记录
// 工时执行人、项目负责人，以及有工时字段权限角色的协作人能对工时进行编辑
func CheckHasWorkHoursUpdatePower(orgId int64, paramBo bo.CheckWorkHoursUpdatePowerBo) (bool, errs.SystemErrorInfo) {
	powerUserIds := []int64{}
	// 查询 issue 的 owner
	issues, err := domain.GetIssueInfosLc(orgId, 0, []int64{paramBo.IssueId})
	if err != nil || len(issues) < 1 {
		return false, err
	}
	if len(issues) < 1 {
		return false, errs.IssueNotExist
	}
	issue := issues[0]

	// 检查是否有项目的管理权限。有权创建项目，则意味着是管理员。
	currentUserId := paramBo.CurrentUserId
	err = domain.AuthProject(orgId, currentUserId, issue.ProjectId, consts.RoleOperationPathOrgProIssueT, consts.OperationOrgProjectCreate)
	if err == nil {
		// 如果 err 为 nil，则表示有权限
		return true, nil
	}

	// 有权限编辑的人
	powerUserIds = append(powerUserIds, paramBo.WorkerId)
	// 检查是否是项目的负责人
	projectBos, err := domain.GetProjectRelationByType(issue.ProjectId, []int64{consts.IssueRelationTypeOwner})
	if err != nil {
		return false, err
	}
	for _, oneRelate := range *projectBos {
		powerUserIds = append(powerUserIds, oneRelate.RelationId)
	}
	// 任务协作人，并且协作者角色中有“工时”列的编辑权限，才可创建工时。
	projectInfo, err := domain.GetProjectSimple(orgId, issue.ProjectId)
	if err != nil {
		log.Errorf("[CheckPowerForCreateIssueWorkHours] issueId: %d, GetProject err: %v", issue.Id, err)
		return false, err
	}
	couldWriteWorkHour, err := domain.CheckHasWriteForWorkHourColumn(orgId, currentUserId, projectInfo.AppId, issue)
	if err != nil {
		log.Errorf("[CheckPowerForCreateIssueWorkHours] issueId: %d, CheckHasWriteForWorkHourColumn err: %v", issue.Id, err)
		return false, err
	}

	// 任务的创建人，也可以编辑、删除工时。
	powerUserIds = append(powerUserIds, issue.Creator)
	// 是否具备权限
	ok, _ := slice.Contain(powerUserIds, paramBo.CurrentUserId)
	if couldWriteWorkHour || ok {
		return true, nil
	} else {
		return false, errs.DenyUpdateIssueWorkHours
	}
}

// CheckPowerForCreateIssueWorkHours 创建工时记录时，检查是否有权限
// 1.先检查当前用户是否是项目负责人、任务负责人、关注人
// 2.检查当前人是否是任务的相关人
func CheckPowerForCreateIssueWorkHours(orgId, currentUserId int64, issue *bo.IssueBo) (bool, errs.SystemErrorInfo) {
	var err errs.SystemErrorInfo
	// 检查是否是项目的负责人
	ownerIds, errSys := businees.LcMemberToUserIdsWithError(issue.OwnerId)
	if errSys != nil {
		log.Errorf("原始的数据:%v, err:%v", issue.OwnerId, err)
		return false, errSys
	}
	if ok, _ := slice.Contain(ownerIds, currentUserId); ok {
		return true, nil
	}

	// 检查是否有项目的管理权限。有权创建项目，则意味着是管理员。
	err = domain.AuthProject(orgId, currentUserId, issue.ProjectId, consts.RoleOperationPathOrgProIssueT, consts.OperationOrgProjectCreate)
	if err == nil {
		// 如果 err 为 nil，则表示有权限
		return true, nil
	}

	projectBos, err := domain.GetProjectRelationByType(issue.ProjectId, []int64{consts.IssueRelationTypeOwner})
	if err != nil {
		return false, err
	}
	relateUsers := []int64{}
	for _, oneRelate := range *projectBos {
		relateUsers = append(relateUsers, oneRelate.RelationId)
	}
	if ok, _ := slice.Contain(relateUsers, currentUserId); ok {
		return true, nil
	}
	// 2.任务协作人，并且协作者角色中有“工时”列的编辑权限，才可创建工时。
	projectInfo, err := domain.GetProjectSimple(orgId, issue.ProjectId)
	if err != nil {
		log.Errorf("[CheckPowerForCreateIssueWorkHours] issueId: %d, GetProject err: %v", issue.Id, err)
		return false, err
	}
	couldWriteWorkHour, err := domain.CheckHasWriteForWorkHourColumn(orgId, currentUserId, projectInfo.AppId, issue)
	if err != nil {
		log.Errorf("[CheckPowerForCreateIssueWorkHours] issueId: %d, CheckHasWriteForWorkHourColumn err: %v", issue.Id, err)
		return false, err
	}
	relateUsers = []int64{}
	// 任务的创建人，也可以创建工时。
	relateUsers = append(relateUsers, issue.Creator)
	ok, _ := slice.Contain(relateUsers, currentUserId)
	if couldWriteWorkHour || ok {
		return true, nil
	}

	return false, errs.DenyCreateIssueWorkHours
}

// 引入工时字段后，就不维护工时开关了，可以通过是否有工时列来判断是否开启工时。
func checkWorkHourIsEnable(orgId, projectId int64) (err errs.SystemErrorInfo) {
	return nil
}

// 通过 issueId 获取项目负责人、任务负责人，参与人
func GetPermissionUserIdsForWorkHour(orgId, userId, issueId int64) (owners, participant, followers []int64, err errs.SystemErrorInfo) {
	// 获取任务信息
	issues, err := domain.GetIssueInfosLc(orgId, userId, []int64{issueId})
	if err != nil {
		log.Error(err)
		return
	}
	if len(issues) < 1 {
		err = errs.IssueNotExist
		log.Error(err)
		return
	}
	issue := issues[0]
	ownerIds := []int64{}
	for _, own := range issue.OwnerId {
		ownerId, errParse := strconv.ParseInt(own[2:], 10, 64)
		if errParse != nil {
			log.Error(errParse)
			return
		}
		ownerIds = append(ownerIds, ownerId)
	}
	owners = append(owners, ownerIds...)
	projectBos, err := domain.GetProjectRelationByType(issue.ProjectId, []int64{consts.IssueRelationTypeOwner})
	if err != nil {
		log.Error(err)
		return
	}
	for _, oneRelate := range *projectBos {
		owners = append(owners, oneRelate.RelationId)
	}
	// 任务的创建人，也可以创建、更新工时。
	owners = append(owners, issue.Creator)
	// 查询任务参与人
	issueRelateBos, err := domain.GetIssueRelationByRelateTypeList(orgId, issueId, []int{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeFollower})
	if err != nil {
		log.Error(err)
		return
	}
	for _, oneRelate := range issueRelateBos {
		switch oneRelate.RelationType {
		case consts.IssueRelationTypeParticipant:
			participant = append(participant, oneRelate.RelationId)
		case consts.IssueRelationTypeFollower:
			followers = append(followers, oneRelate.RelationId)
		}
	}
	return
}

// 工时统计
// startTime、endTime 如果不传，则会默认最近 7 天内
func GetWorkHourStatistic(orgId int64, param vo.GetWorkHourStatisticReq) (*vo.GetWorkHourStatisticResp, errs.SystemErrorInfo) {
	returnData := &vo.GetWorkHourStatisticResp{
		GroupStatisticList: make([]*vo.OnePersonWorkHourStatisticInfo, 0),
		Total:              0,
		Summary:            &vo.GetWorkHourStatisticSummary{},
	}
	// 参数校验和校准。
	if err := CheckAndFixGetStatisticDataParam(&param); err != nil {
		log.Errorf("[GetWorkHourStatistic] err: %v", err)
		return returnData, err
	}
	projectIds := []int64{}
	if len(param.ProjectIds) > 0 {
		copyer.Copy(param.ProjectIds, &projectIds)
	}
	// 防止任务id太多，先查询一下有工时记录的任务id，以减小查询范围。
	limitIssueIds, err := domain.GetHasWorkHourRecordIssueIds(orgId, projectIds)
	if err != nil {
		log.Errorf("[GetWorkHourStatistic] err: %v", err)
		return returnData, err
	}
	if len(param.IssueIds) < 1 && len(limitIssueIds) > 0 {
		param.IssueIds = make([]*int64, 0)
		copyer.Copy(limitIssueIds, &param.IssueIds)
	}
	// 参数一些初始化
	// 通过**筛选条件**，获取符合条件的任务 ids
	// 考虑一下，如果任务id过多怎么办？待优化
	issueIds, err := GetIssueIdsByStatisticInput(orgId, param)
	if err != nil {
		return returnData, err
	}
	// 做个预警，如果超过 10w 条，则需改进。
	if len(issueIds) > 100000 {
		errMsg := "工时统计-查询到的任务数量已超过 100000 条。请开发者注意！"
		log.Errorf("[GetWorkHourStatistic] err: %s", errMsg)
		// 发送日志到钉钉
		commonfacade.DingTalkInfo(commonvo.DingTalkInfoReqVo{
			LogData: map[string]interface{}{
				"app": config.GetApplication().Name,
				"env": config.GetEnv(),
				"msg": errMsg,
				"err": nil,
			},
		})
	}
	// 根据筛选条件查询符合条件的工时执行人。
	workerIds := []int64{}
	totalWorkerIds := []int64{}
	totalWorkerIds, err = GetWorkerIdsForStatistic(orgId, param)
	if err != nil {
		return returnData, err
	}
	// 处理分页。按工时执行人分页。
	totalWorkerIds = int64Slice.ArrayUnique(totalWorkerIds)
	returnData.Total = int64(len(totalWorkerIds))
	// 找不到要统计的执行人，则直接返回空数组
	if returnData.Total == 0 {
		return returnData, nil
	}
	offset := (*param.Page - 1) * (*param.Size)
	end := offset + (*param.Size)
	max := int64(len(totalWorkerIds))
	if offset < 0 {
		offset = 0
	}
	if end > max {
		end = max
	}
	if offset > max {
		offset = max
	}
	workerIds = totalWorkerIds[offset:end]
	// 先查询时间段内的所有工时记录。
	// 总预估工时和子预估工时，一定会有起止时间，可以按照 StartTimeT1 和 EndTimeT2 进行筛选。
	// 而实际工时记录，因为 end_time 字段为 0，刚好也满足 end_time<EndTimeT2。
	condParam := bo.IssueWorkHoursBoListCondBo{
		OrgId:       orgId,
		ProjectIds:  projectIds,
		IssueIds:    issueIds,
		WorkerIds:   workerIds,
		StartTimeT1: *param.StartTime,
		EndTimeT2:   *param.EndTime,
		IsDelete:    consts.AppIsNoDelete,
		Page:        1,
		Size:        100_0000, // 给个较大的数值，表示查询所有
	}
	list, err := domain.GetIssueWorkHoursStatisticList(condParam)
	if err != nil {
		log.Error(err)
		return returnData, err
	}
	if summary, err := GetStatisticDataForSummary(orgId, param, projectIds, issueIds, totalWorkerIds); err != nil {
		return returnData, err
	} else {
		returnData.Summary = &summary
	}
	// 循环一遍，找出不能参与的总预估工时记录
	cacheMapSubWorkHour := map[int64][]uint64{}
	for _, item := range list {
		if item.Type == consts2.WorkHourTypeSubPredict {
			if _, ok := cacheMapSubWorkHour[item.IssueId]; ok {
				cacheMapSubWorkHour[item.IssueId] = append(cacheMapSubWorkHour[item.IssueId], item.Id)
			} else {
				cacheMapSubWorkHour[item.IssueId] = []uint64{item.Id}
			}
		}
	}
	// 统计详细类预估工时、实际工时，按照 workerId 分类
	groupList := map[int64]*po.OnePersonWorkHourData{}
	for _, workerId := range workerIds {
		groupList[workerId] = &po.OnePersonWorkHourData{
			TotalPredict: []*po.PpmPriIssueWorkHours{},
			SubPredict:   []*po.PpmPriIssueWorkHours{},
			Actual:       []*po.PpmPriIssueWorkHours{},
		}
	}
	tmpOnePersonWorkHourData := &po.OnePersonWorkHourData{}
	for _, record := range list {
		tmpOnePersonWorkHourData = nil
		if val, ok := groupList[record.WorkerId]; ok {
			tmpOnePersonWorkHourData = val
		} else {
			continue
		}
		// 复制对象
		tmpCloneObj := record
		switch record.Type {
		case consts2.WorkHourTypeTotalPredict:
			// 如果预估工时只有总预估工时没有子预估工时，则将其统计为预估工时。
			// 如果预估工时包含了子预估工时，则不能统计为预估工时（因为会重复统计）。
			if _, ok := cacheMapSubWorkHour[record.IssueId]; !ok {
				groupList[record.WorkerId].SubPredict = append(tmpOnePersonWorkHourData.SubPredict, &tmpCloneObj)
			}
		case consts2.WorkHourTypeSubPredict:
			groupList[record.WorkerId].SubPredict = append(tmpOnePersonWorkHourData.SubPredict, &tmpCloneObj)
		case consts2.WorkHourTypeActual:
			groupList[record.WorkerId].Actual = append(tmpOnePersonWorkHourData.Actual, &tmpCloneObj)
		}
	}
	// 批量获取用户信息
	usersNameMap := map[int64]string{}
	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, workerIds)
	if err != nil {
		log.Error(err)
		return returnData, err
	}
	for _, user := range userInfos {
		usersNameMap[user.UserId] = user.Name
	}
	// 数据分组完成后，一个用户一个用户地计算
	groupStatisticData := []*vo.OnePersonWorkHourStatisticInfo{}
	for _, workerId := range workerIds {
		tmpName := ""
		tmpOnePersonWorkHourData = &po.OnePersonWorkHourData{}
		if val, ok := groupList[workerId]; ok {
			tmpOnePersonWorkHourData = val
		} else {
			continue
		}
		if val, ok := usersNameMap[workerId]; ok {
			tmpName = val
		} else {
			tmpName = "匿名"
		}
		actualTotal, workHourList, err := GetStatisticDataForActual(*param.StartTime, *param.EndTime, tmpOnePersonWorkHourData.Actual)
		if err != nil {
			return returnData, err
		}
		predictTotal, err := GetStatisticDataForPredict(tmpOnePersonWorkHourData.SubPredict)
		if err != nil {
			return returnData, err
		}
		workHourListTrans := []*vo.OneDateWorkHour{}
		copyer.Copy(workHourList, &workHourListTrans)
		groupStatisticData = append(groupStatisticData, &vo.OnePersonWorkHourStatisticInfo{
			WorkerID:         workerId,
			Name:             tmpName,
			PredictHourTotal: format.FormatNeedTimeIntoString(predictTotal),
			ActualHourTotal:  format.FormatNeedTimeIntoString(actualTotal),
			DateWorkHourList: workHourListTrans,
		})
	}
	returnData.GroupStatisticList = groupStatisticData

	return returnData, err
}

// 通过查询条件获取要查询的工时执行人。
func GetWorkerIdsForStatistic(orgId int64, param vo.GetWorkHourStatisticReq) (totalWorkerIds []int64, retErr errs.SystemErrorInfo) {
	if len(param.WorkerIds) > 0 {
		copyer.Copy(param.WorkerIds, &totalWorkerIds)
	} else {
		// 没有传入的工时执行人时，进行查询
		padWorkerIds := []int64{}

		// 系统管理员（超管、普通管理员）的工时数据也需要展示
		adminUserIds, err := domain.GetAdminUserIdsOfOrg(orgId, -1)
		if err != nil {
			log.Error(err)
			retErr = err
			return
		}
		padWorkerIds = slice.SliceUniqueInt64(append(padWorkerIds, adminUserIds...))

		cond1 := db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeOwner, consts.IssueRelationTypeDepartmentParticipant}),
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}
		if len(param.ProjectIds) > 0 {
			cond1[consts.TcProjectId] = db.In(param.ProjectIds)
		}
		relationBos, err := domain.GetProjectRelationByCond(cond1)
		if err != nil {
			log.Error(err)
			retErr = err
			return
		}
		deptIds := make([]int64, 0)
		for _, oneRelation := range *relationBos {
			// 如果是部门，则拿出来查询部门人员
			if oneRelation.RelationType == consts.IssueRelationTypeDepartmentParticipant {
				deptIds = append(deptIds, oneRelation.RelationId)
			} else {
				padWorkerIds = append(padWorkerIds, oneRelation.RelationId)
			}
		}
		// 查询部门成员
		if len(deptIds) > 0 {
			userIdsResp := orgfacade.GetUserIdsByDeptIds(&orgvo.GetUserIdsByDeptIdsReq{
				OrgId:   orgId,
				DeptIds: deptIds,
			})
			if userIdsResp.Failure() {
				retErr = userIdsResp.Error()
				log.Error(retErr)
				return
			}
			padWorkerIds = append(padWorkerIds, userIdsResp.Data.UserIds...)
		}

		totalWorkerIds = padWorkerIds
	}
	if *param.ShowResigned == 2 {
		resp := orgfacade.FilterResignedUserIds(orgvo.FilterResignedUserIdsReqVo{
			OrgId:   orgId,
			UserIds: totalWorkerIds,
		})
		if len(resp.Data) < 1 {
			// 没有用户信息，表示查询不到用户信息。
			// 查询不到有效的执行人时，直接返回空数组
			// retErr = errs.BuildSystemErrorInfoWithMessage(errs.UserNotFoundError, " 找不到在职的用户工时信息。")
			totalWorkerIds = make([]int64, 0)
			return
		}
		totalWorkerIds = resp.Data
	}
	var sortPadWorkerIds sortutil.Int64Slice = totalWorkerIds
	sortPadWorkerIds.Sort()
	totalWorkerIds = sortPadWorkerIds
	return
}

// GetIssueIdsByStatisticInput 通过统计的输入参数，获取符合条件的所有任务id
func GetIssueIdsByStatisticInput(orgId int64, param vo.GetWorkHourStatisticReq) ([]int64, errs.SystemErrorInfo) {
	// param 中影响结果的若干个字段：projectIds, issueIds, issueStatus...
	issueIds := []int64{}
	//issueCond := db.Cond{}
	issueCnd := make([]*tablev1.Condition, 0, 6) // 预估一个容量
	issueCnd = append(issueCnd, domain.GetRowsCondition(consts.BasicFieldOrgId, tablev1.ConditionType_equal, orgId, nil))
	if len(param.ProjectIds) > 0 {
		// 通过 projectIds 查询任务
		issueCnd = append(issueCnd, domain.GetRowsCondition(consts.BasicFieldProjectId, tablev1.ConditionType_in, nil, param.ProjectIds))
	}
	// 筛选任务状态
	if len(param.IssueStatus) > 0 {
		issueCnd = append(issueCnd, domain.GetRowsCondition(consts.BasicFieldIssueStatus, tablev1.ConditionType_in, nil, param.IssueStatus))
	}
	// 筛选任务优先级
	if len(param.IssuePriorities) > 0 {
		issueCnd = append(issueCnd, domain.GetRowsCondition(consts.ProBasicFieldPriorityId, tablev1.ConditionType_in, nil, param.IssuePriorities))
	}
	// 筛选任务id
	if len(param.IssueIds) > 0 {
		issueCnd = append(issueCnd, domain.GetRowsCondition(consts.BasicFieldIssueId, tablev1.ConditionType_in, nil, param.IssueIds))
	}
	//issueInfoList, err := domain.GetIssueIdListByCond(orgId, issueCond)
	issueInfoList, err := domain.GetIssueIdList(orgId, issueCnd, 0, 0)
	if err != nil {
		log.Error(err)
		return issueIds, err
	}
	for _, val := range issueInfoList {
		issueIds = append(issueIds, val.Id)
	}
	return issueIds, nil
}

func CheckAndFixGetStatisticDataParam(param *vo.GetWorkHourStatisticReq) errs.SystemErrorInfo {
	// 校验 time range
	zone, _ := time.LoadLocation("Local")
	if param.StartTime == nil || *param.StartTime < 1 || param.EndTime == nil || *param.EndTime < 1 {
		defaultEndTime := time.Now().Unix()
		defaultStartTime := defaultEndTime - 7*24*3600
		t1 := time.Unix(defaultStartTime, 0)
		t2 := time.Unix(defaultEndTime, 0)
		t1Start := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, zone)
		t2End := time.Date(t2.Year(), t2.Month(), t2.Day(), 23, 59, 59, 0, zone)
		defaultStartTime = t1Start.Unix()
		defaultEndTime = t2End.Unix()
		param.StartTime = &defaultStartTime
		param.EndTime = &defaultEndTime
	} else {
		// 设定半年时间。即不能超过 185 天。
		if !format.CheckTimeRangeLimitDayNumValid(*param.StartTime, *param.EndTime, 185) {
			return errs.TimeRangeLimitHalfYearInvalid
		}
		t1 := time.Unix(*param.StartTime, 0)
		t2 := time.Unix(*param.EndTime, 0)
		t1Start := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, zone)
		t2End := time.Date(t2.Year(), t2.Month(), t2.Day(), 23, 59, 59, 0, zone)
		*param.StartTime = t1Start.Unix()
		*param.EndTime = t2End.Unix()
	}
	if param.Page == nil || *param.Page < 1 {
		defaultPage := int64(1)
		param.Page = &defaultPage
	}
	if param.Size == nil || *param.Size < 1 {
		defaultSize := int64(20)
		param.Size = &defaultSize
	}
	// 是否显示已离职人员。 默认不显示。
	if param.ShowResigned == nil {
		defaultShowResigned := 2
		param.ShowResigned = &defaultShowResigned
	}
	return nil
}

// 筛选条件下的数据汇总
func GetStatisticDataForSummary(orgId int64, param vo.GetWorkHourStatisticReq, projectIds, issueIds, totalWorkerIds []int64) (vo.GetWorkHourStatisticSummary, errs.SystemErrorInfo) {
	returnData := vo.GetWorkHourStatisticSummary{
		PredictTotal: "",
		ActualTotal:  "",
	}
	// 先查询时间段内的所有工时记录
	condParam := bo.IssueWorkHoursBoListCondBo{
		OrgId:       orgId,
		ProjectIds:  projectIds,
		IssueIds:    issueIds,
		WorkerIds:   totalWorkerIds,
		StartTimeT1: *param.StartTime,
		EndTimeT2:   *param.EndTime,
		IsDelete:    consts.AppIsNoDelete,
		Page:        1,
		Size:        100_0000, // 给个较大的数值，表示查询所有
	}
	total, err := domain.GetCountOfStatisticWorkHourByCond(condParam)
	if err != nil {
		log.Error(err)
		return returnData, err
	}
	// 统计筛选条件下的总预估工时，总实际工时。
	predictTotal := uint32(0)
	actualTotal := uint32(0)
	// 防止数据太多，做一下分批。
	size := float64(10000)
	totalPage := int(math.Ceil(float64(total) / size))
	for condParam.Page = 1; condParam.Page <= totalPage; condParam.Page += 1 {
		list, err := domain.GetIssueWorkHoursStatisticList(condParam)
		if err != nil {
			log.Error(err)
			return returnData, err
		}
		if len(list) <= 0 {
			break
		}
		cacheMapSubWorkHour := map[int64][]uint64{}
		for _, item := range list {
			if item.Type == consts2.WorkHourTypeSubPredict {
				if _, ok := cacheMapSubWorkHour[item.IssueId]; ok {
					cacheMapSubWorkHour[item.IssueId] = append(cacheMapSubWorkHour[item.IssueId], item.Id)
				} else {
					cacheMapSubWorkHour[item.IssueId] = []uint64{item.Id}
				}
			}
		}
		for _, item := range list {
			switch item.Type {
			case consts2.WorkHourTypeSubPredict:
				// 有子预估工时的，则不将该 issue 对应的总预估工时计算在内。
				if _, ok := cacheMapSubWorkHour[item.IssueId]; ok {
					predictTotal += item.NeedTime
				}
			case consts2.WorkHourTypeTotalPredict:
				if _, ok := cacheMapSubWorkHour[item.IssueId]; !ok {
					predictTotal += item.NeedTime
				}
			case consts2.WorkHourTypeActual:
				actualTotal += item.NeedTime
			}
		}
	}
	returnData.PredictTotal = format.FormatNeedTimeIntoString(int64(predictTotal))
	returnData.ActualTotal = format.FormatNeedTimeIntoString(int64(actualTotal))
	return returnData, nil
}

// 实际工时计算：计算时间区域内，每个日期的工时
func GetStatisticDataForActual(startTime, endTime int64, workHours []*po.PpmPriIssueWorkHours) (total int64, list []vo.OneDateWorkHour, err errs.SystemErrorInfo) {
	total = 0
	diyFormat1 := consts.AppTimeFormat
	// 将入参的工时列表转换为以日期为key的map
	workHourNeedTimeMapByDate := map[string]uint32{}
	for _, record := range workHours {
		time1 := time.Unix(int64(record.StartTime), 0)
		// 因为只需要精确到日期，因此时分秒不用理会，视为 `00:00:00`
		date := time.Date(time1.Year(), time1.Month(), time1.Day(), 0, 0, 0, 0, time.Local).Format(diyFormat1)
		workHourNeedTimeMapByDate[date] += record.NeedTime
		total += int64(record.NeedTime)
	}
	// 通过穷举的日期序列进行数据填充
	dateList := GetDateListByDateRange(startTime, endTime)
	for _, oneDate := range dateList {
		tmpTime, err := time.Parse(diyFormat1, oneDate)
		if err != nil {
			log.Error(err)
			return total, list, errs.DateParseError
		}
		tmpTime = time.Date(tmpTime.Year(), tmpTime.Month(), tmpTime.Day(), 0, 0, 0, 0, time.Local)
		weekDay := int64(tmpTime.Weekday())
		if tmpNeedTimeTotal, ok := workHourNeedTimeMapByDate[oneDate]; ok {
			list = append(list, vo.OneDateWorkHour{
				Date:    oneDate,
				WeekDay: weekDay,
				Time:    format.FormatNeedTimeIntoString(int64(tmpNeedTimeTotal)),
			})
		} else {
			list = append(list, vo.OneDateWorkHour{
				Date:    oneDate,
				WeekDay: weekDay,
				Time:    "0",
			})
		}
	}
	return total, list, nil
}

// 预估工时计算
func GetStatisticDataForPredict(workHours []*po.PpmPriIssueWorkHours) (total int64, err errs.SystemErrorInfo) {
	total = 0
	for _, record := range workHours {
		total += int64(record.NeedTime)
	}
	return total, nil
}

// 通过起始时间生成这段日期内的所有日期穷举。
func GetDateListByDateRange(startTime, endTime int64) []string {
	dateList := []string{}
	if endTime < startTime {
		return dateList
	}
	tmpTimestamp := startTime
	for {
		if tmpTimestamp > endTime {
			break
		}
		tmpTime := time.Unix(tmpTimestamp, 0)
		oneDateItem := time.Date(tmpTime.Year(), tmpTime.Month(), tmpTime.Day(), 0, 0, 0, 0, time.Local).Format(consts.AppTimeFormat)
		dateList = append(dateList, oneDateItem)
		tmpTimestamp += 24 * 3600
	}
	return dateList
}

// TriggerCreateWorkHourTrend 新增工时，触发新增动态，触发工时变更卡片推送
func TriggerCreateWorkHourTrend(orgId, currentUserId, issueId int64, newWorkHourParam vo.OneWorkHourRecord) errs.SystemErrorInfo {
	// 查询 issue 信息
	issues, err := domain.GetIssueInfosLc(orgId, currentUserId, []int64{issueId})
	if err != nil || len(issues) < 1 {
		return errs.IssueNotExist
	}
	issue := issues[0]
	newWorkHourJson := json.ToJsonIgnoreError(newWorkHourParam)
	issueTrendsBo := &bo.IssueTrendsBo{
		PushType:          consts.PushTypeUpdateIssue,
		OrgId:             orgId,
		DataId:            issue.DataId,
		IssueId:           issueId,
		ProjectId:         issue.ProjectId,
		TableId:           issue.TableId,
		IssueTitle:        issue.Title,
		OperatorId:        currentUserId,
		NewValue:          newWorkHourJson,
		BeforeWorkHourIds: []int64{},
		AfterWorkHourIds:  []int64{newWorkHourParam.ID},
		UpdateWorkHour:    true,
	}
	// 查询工时列信息
	columnMap, err := domain.GetTableColumnsMap(orgId, issue.TableId, []string{consts.ProBasicFieldWorkHour})
	if err != nil {
		log.Errorf("[TriggerCreateWorkHourTrend] err: %v, tableId: %d", err, issue.TableId)
		return err
	}
	workHourColumn, ok := columnMap[consts.ProBasicFieldWorkHour]
	if !ok {
		err := errs.TableColumnNotExist
		log.Errorf("[TriggerCreateWorkHourTrend] err: %v", err)
		return err
	}
	changeList := []bo.TrendChangeListBo{
		{
			Field:      consts.ProBasicFieldWorkHour,
			FieldType:  consts.ProBasicFieldWorkHour,
			FieldName:  workHourColumn.Label,
			AliasTitle: workHourColumn.AliasTitle,
			NewValue:   newWorkHourJson,
		},
	}
	issueTrendsBo.Ext.ChangeList = changeList
	domain.PushIssueTrends(issueTrendsBo)

	// fs 卡片推送
	if err := TriggerFsCardForWorkHourChange(issueTrendsBo); err != nil {
		log.Errorf("[TriggerCreateWorkHourTrend] err: %v, issueId: %d", err, issueId)
		return err
	}

	return nil
}

// TriggerFsCardForWorkHourChange 工时的卡片推送
func TriggerFsCardForWorkHourChange(issueTrendsBo *bo.IssueTrendsBo) errs.SystemErrorInfo {
	return nil
}

// TriggerUpdateWorkHourTrend 编辑工时，触发新增动态
func TriggerUpdateWorkHourTrend(orgId, currentUserId, issueId int64, oldWorkHour vo.OneWorkHourRecord, newWorkHour vo.OneWorkHourRecord) error {
	if oldWorkHour.Type == consts2.WorkHourTypeTotalPredict {
		log.Infof("[TriggerUpdateWorkHourTrend] issueId: %d, workHourId: %d", issueId, oldWorkHour.ID)
		return nil
	}
	// 查询 issue 信息
	issues, err := domain.GetIssueInfosLc(orgId, currentUserId, []int64{issueId})
	if err != nil || len(issues) < 1 {
		return errs.IssueNotExist
	}
	issue := issues[0]
	changeList := []bo.TrendChangeListBo{}
	if oldWorkHour.Worker.UserID != newWorkHour.Worker.UserID {
		changeList = append(changeList, bo.TrendChangeListBo{
			FieldName:     consts.WorkId,
			Field:         "worker_id",
			OldValue:      oldWorkHour.Worker.Name,
			NewValue:      newWorkHour.Worker.Name,
			IsForWorkHour: true,
		})
	}
	if oldWorkHour.StartTime != newWorkHour.StartTime {
		changeList = append(changeList, bo.TrendChangeListBo{
			FieldName:     consts.StartTime,
			Field:         "start_time",
			OldValue:      strconv.FormatInt(oldWorkHour.StartTime, 10),
			NewValue:      strconv.FormatInt(newWorkHour.StartTime, 10),
			IsForWorkHour: true,
		})
	}
	if oldWorkHour.EndTime != newWorkHour.EndTime {
		changeList = append(changeList, bo.TrendChangeListBo{
			FieldName:     consts.PlanEndTime,
			Field:         "end_time",
			OldValue:      strconv.FormatInt(oldWorkHour.EndTime, 10),
			NewValue:      strconv.FormatInt(newWorkHour.EndTime, 10),
			IsForWorkHour: true,
		})
	}
	if oldWorkHour.NeedTime != newWorkHour.NeedTime {
		changeList = append(changeList, bo.TrendChangeListBo{
			FieldName:     consts.WorkTime,
			Field:         "need_time",
			OldValue:      oldWorkHour.NeedTime,
			NewValue:      newWorkHour.NeedTime,
			IsForWorkHour: true,
		})
	}
	if oldWorkHour.Desc != newWorkHour.Desc {
		changeList = append(changeList, bo.TrendChangeListBo{
			FieldName:     consts.WorkContent,
			Field:         "desc",
			OldValue:      oldWorkHour.Desc,
			NewValue:      newWorkHour.Desc,
			IsForWorkHour: true,
		})
	}
	issueTrendsBo := &bo.IssueTrendsBo{
		PushType:   consts.PushTypeUpdateIssue,
		OrgId:      orgId,
		DataId:     issue.DataId,
		IssueId:    issueId,
		ProjectId:  issue.ProjectId,
		TableId:    issue.TableId,
		IssueTitle: issue.Title,
		OperatorId: currentUserId,
		OldValue:   json.ToJsonIgnoreError(oldWorkHour),
		NewValue:   json.ToJsonIgnoreError(newWorkHour),
		Ext: bo.TrendExtensionBo{
			ChangeList: changeList,
		},
		// 目前的情况是，一个任务下，更新了 m 个工时的 n 个字段，则会对应产生 m 个动态。单条更新，则直接放入更新的工时 id
		BeforeWorkHourIds: []int64{oldWorkHour.ID},
		AfterWorkHourIds:  []int64{newWorkHour.ID},
		UpdateWorkHour:    true,
	}
	domain.PushIssueTrends(issueTrendsBo)

	// fs 卡片推送
	if err := TriggerFsCardForWorkHourChange(issueTrendsBo); err != nil {
		log.Errorf("[TriggerUpdateWorkHourTrend] err: %v, issueId: %d", err, issueId)
		return err
	}

	return nil
}

// TriggerDeleteWorkHourTrend 删除工时，触发新增动态、飞书卡片推送
func TriggerDeleteWorkHourTrend(orgId, currentUserId, issueId int64, workHour vo.OneWorkHourRecord) errs.SystemErrorInfo {
	// 查询 issue 信息
	issues, err := domain.GetIssueInfosLc(orgId, currentUserId, []int64{issueId})
	if err != nil || len(issues) < 1 {
		return errs.IssueNotExist
	}
	issue := issues[0]
	workHourJson := json.ToJsonIgnoreError(workHour)
	issueTrendsBo := &bo.IssueTrendsBo{
		PushType:          consts.PushTypeUpdateIssue,
		OrgId:             orgId,
		DataId:            issue.DataId,
		IssueId:           issueId,
		ProjectId:         issue.ProjectId,
		TableId:           issue.TableId,
		IssueTitle:        issue.Title,
		OperatorId:        currentUserId,
		OldValue:          workHourJson,
		BeforeWorkHourIds: []int64{workHour.ID},
		AfterWorkHourIds:  []int64{},
		UpdateWorkHour:    true,
	}
	// 查询工时列信息
	columnMap, err := domain.GetTableColumnsMap(orgId, issue.TableId, []string{consts.ProBasicFieldWorkHour})
	if err != nil {
		log.Errorf("[TriggerDeleteWorkHourTrend] err: %v, tableId: %d", err, issue.TableId)
		return err
	}
	workHourColumn, ok := columnMap[consts.ProBasicFieldWorkHour]
	if !ok {
		err := errs.TableColumnNotExist
		log.Errorf("[TriggerDeleteWorkHourTrend] err: %v", err)
		return err
	}
	changeList := []bo.TrendChangeListBo{
		{
			Field:      consts.ProBasicFieldWorkHour,
			FieldType:  consts.ProBasicFieldWorkHour,
			FieldName:  workHourColumn.Label,
			AliasTitle: workHourColumn.AliasTitle,
			OldValue:   workHourJson,
		},
	}
	issueTrendsBo.Ext.ChangeList = changeList
	domain.PushIssueTrends(issueTrendsBo)

	// fs 卡片推送
	if err := TriggerFsCardForWorkHourChange(issueTrendsBo); err != nil {
		log.Errorf("[TriggerDeleteWorkHourTrend] err: %v, issueId: %d", err, issueId)
		return err
	}

	return nil
}

// 生成开启工时或者关闭工时功能的项目动态
func TriggerEnableOrDisableWorkHourTrend(orgId, currentUserId int64, project *bo.ProjectBo, disOrEnable string) errs.SystemErrorInfo {
	workHourTrendType := consts.PushTypeEnableWorkHour
	switch disOrEnable {
	case "enable":
		workHourTrendType = consts.PushTypeEnableWorkHour
	case "disable":
		workHourTrendType = consts.PushTypeDisableWorkHour
	}
	trendsBo := bo.ProjectTrendsBo{
		PushType:   workHourTrendType,
		OrgId:      orgId,
		ProjectId:  project.Id,
		OperatorId: currentUserId,
		Ext: bo.TrendExtensionBo{
			IssueType: "",
			ObjName:   project.Name,
		},
	}
	domain.PushProjectTrends(trendsBo)
	return nil
}

// 通过工时中的 worker_id 获取对应工时执行者的信息
func GetWorkerInfos(orgId int64, workerIds []int64) (map[int64]*vo.WorkHourWorker, errs.SystemErrorInfo) {
	userInfoMapById := map[int64]*vo.WorkHourWorker{}
	if len(workerIds) < 1 {
		return userInfoMapById, nil
	}
	// 查询工时执行人的信息
	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, workerIds)
	if err != nil {
		log.Error(err)
		return userInfoMapById, err
	}
	for _, user := range userInfos {
		userInfoMapById[user.UserId] = &vo.WorkHourWorker{
			UserID: user.UserId,
			Name:   user.Name,
			Avatar: user.Avatar,
		}
	}
	return userInfoMapById, nil
}

// 当任务的负责人变更之后，对应的工时的执行人也发生变更
func ChangeWorkerIdWhenChangedIssueOwner(orgId, issueId, oldOwner, newOwner int64) errs.SystemErrorInfo {
	// 如果修改前后的负责人不变，则无需执行下面的逻辑
	if oldOwner == newOwner || oldOwner < 1 || newOwner < 1 {
		return nil
	}
	condition := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  issueId,
		consts.TcWorkerId: oldOwner,
		consts.TcType:     consts2.WorkHourTypeTotalPredict,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	// 检查工时是否存在
	list, err := domain.GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
		OrgId:     orgId,
		IssueId:   issueId,
		Type:      consts2.WorkHourTypeTotalPredict,
		WorkerIds: []int64{oldOwner},
		IsDelete:  consts.AppIsNoDelete,
	})
	// 工时不存在，则直接返回 nil，无需做任何操作。
	if len(list) < 1 {
		return nil
	}
	if err != nil {
		return err
	}
	if len(list) < 1 {
		return errs.IssueWorkHourNotExist
	}
	err = domain.UpdateIssueWorkHourByCond(condition, bo.UpdateOneIssueWorkHoursBo{
		WorkerId: newOwner,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// 更新任务时，触发更新总预估工时的某个字段。
func TriggerWorkHourWhenChangedPlanStartOrEndTime(orgId, issueId int64, updateField string, oldValue, newValue types2.Time) errs.SystemErrorInfo {
	oldValue1, newValue1 := int64(0), int64(0)
	newUpdateBo := bo.UpdateOneIssueWorkHoursBo{}
	timezoneOffset := int64(8 * 3600)
	switch updateField {
	case consts.TcStartTime:
		oldValue1 = time.Time(oldValue).Unix()
		newValue1 = time.Time(newValue).Unix()
		// 入参的 types.Time 丢失了时区，默认为当前的时区。转为时间戳时，会多出 8h，因此减去偏移量。
		newUpdateBo.StartTime = newValue1 - timezoneOffset
	case consts.TcEndTime:
		oldValue1 = time.Time(oldValue).Unix()
		newValue1 = time.Time(newValue).Unix()
		newUpdateBo.EndTime = newValue1 - timezoneOffset
	}
	// 如果修改前后的值不变，则无需执行下面的逻辑。
	if oldValue1 == newValue1 || newValue1 < 1 {
		return nil
	}
	condition := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIssueId:  issueId,
		consts.TcType:     consts2.WorkHourTypeTotalPredict,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	// 检查工时是否存在
	list, err := domain.GetIssueWorkHoursList(bo.IssueWorkHoursBoListCondBo{
		OrgId:    orgId,
		IssueId:  issueId,
		Type:     consts2.WorkHourTypeTotalPredict,
		IsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		return err
	}
	if len(list) < 1 {
		// 如果工时不存在，则无需更新，直接返回 nil
		return nil
	}
	err = domain.UpdateIssueWorkHourByCond(condition, newUpdateBo)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// 检查员工是否是项目成员
func CheckIsIssueMember(orgId int64, param vo.CheckIsIssueMemberReq) (*vo.BoolResp, errs.SystemErrorInfo) {
	returnRes := &vo.BoolResp{
		IsTrue: false,
	}
	// 2.检查当前人是否是任务的负责人、关注人
	issueRelateBos, err := domain.GetIssueRelationByRelateTypeList(orgId, param.IssueID, []int{consts.IssueRelationTypeOwner, consts.IssueRelationTypeFollower})
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	relateUserIds := []int64{}
	for _, oneRelate := range issueRelateBos {
		relateUserIds = append(relateUserIds, oneRelate.RelationId)
	}
	if ok, _ := slice.Contain(relateUserIds, param.UserID); ok {
		returnRes.IsTrue = true
	}
	return returnRes, nil
}

// 检查是否是任务相关人员（创建人，关注人，负责人）
func CheckIsIssueRelatedPeople(orgId int64, param vo.CheckIsIssueMemberReq) (*vo.BoolResp, errs.SystemErrorInfo) {
	returnRes := &vo.BoolResp{
		IsTrue: false,
	}
	info, infoErr := domain.GetIssueInfoLc(orgId, 0, param.IssueID)
	if infoErr != nil {
		log.Error(infoErr)
		return returnRes, infoErr
	}
	if info.Creator == param.UserID {
		returnRes.IsTrue = true
		return returnRes, nil
	}
	// 2.检查当前人是否是任务的负责人、关注人
	issueRelateBos, err := domain.GetIssueRelationByRelateTypeList(orgId, param.IssueID, []int{consts.IssueRelationTypeOwner, consts.IssueRelationTypeFollower})
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	relateUserIds := []int64{}
	for _, oneRelate := range issueRelateBos {
		relateUserIds = append(relateUserIds, oneRelate.RelationId)
	}
	if ok, _ := slice.Contain(relateUserIds, param.UserID); ok {
		returnRes.IsTrue = true
	}
	return returnRes, nil
}

// 将用户加入为任务的关注者
func SetUserJoinIssue(orgId, currentUserId int64, param vo.SetUserJoinIssueReq) (*vo.BoolResp, errs.SystemErrorInfo) {
	returnRes := &vo.BoolResp{
		IsTrue: false,
	}
	// 查询是否已经是负责人或关注者
	resp, _ := CheckIsIssueMember(orgId, vo.CheckIsIssueMemberReq{
		IssueID: param.IssueID,
		UserID:  param.UserID,
	})
	// 如果将一个用户设置为任务的关注人，而他已经是关注人，此时直接返回而非报错，因为结果都是“该用户是关注人”。
	if resp.IsTrue == true {
		// errMsg := "该员工已经是任务参与者。"
		return resp, nil
	}
	// 查询任务
	issues, err := domain.GetIssueInfosLc(orgId, currentUserId, []int64{param.IssueID})
	if err != nil {
		log.Error(err)
		return returnRes, err
	}
	if len(issues) < 1 {
		log.Error("任务不存在")
		return returnRes, errs.IssueNotExist
	}
	issue := issues[0]

	// 先将其加为项目参与人
	err = AddUserIntoProParticipant(orgId, currentUserId, param.UserID, issue.ProjectId)
	if err != nil {
		log.Error(err)
		return returnRes, err
	}

	// 将目标用户加为任务关注人
	_, err = domain.UpdateIssueRelationSingle(currentUserId, issue, consts.IssueRelationTypeFollower, []int64{param.UserID})
	if err != nil {
		log.Error(err)
		return returnRes, err
	}

	//同步到无码
	followerIds, err3 := domain.GetIssueRelationIdsByRelateType(orgId, param.IssueID, consts.IssueRelationTypeFollower)
	if err3 != nil {
		log.Error(err3)
		return nil, err3
	}
	// 关注人id 在无码中存储的形式是形如 U_23214 的字符串数组
	lcFollowerIds := make([]string, 0)
	for _, tmpId := range *followerIds {
		lcFollowerIds = append(lcFollowerIds, fmt.Sprintf("U_%d", tmpId))
	}
	appId, appIdErr := domain.GetAppIdFromProjectId(orgId, issue.ProjectId)
	if appIdErr != nil {
		log.Error(appIdErr)
		return nil, appIdErr
	}
	lcUpd := map[string]interface{}{
		"issueId":                    param.IssueID,
		consts.BasicFieldFollowerIds: lcFollowerIds,
	}
	lcResp := formfacade.LessUpdateIssue(formvo.LessUpdateIssueReq{
		AppId:   appId,
		OrgId:   orgId,
		UserId:  currentUserId,
		TableId: issue.TableId,
		Form:    []map[string]interface{}{lcUpd},
	})
	if lcResp.Failure() {
		log.Error(lcResp.Error())
		return nil, lcResp.Error()
	}
	returnRes.IsTrue = true

	return returnRes, nil
}

func PushMqttForWorkHour(orgId, userId, issueId int64) {
	issue, err := domain.GetIssueInfoLc(orgId, userId, issueId)
	if err != nil {
		log.Error(err)
		return
	}

	newData := map[string]interface{}{
		consts.BasicFieldId:       issueId,
		consts.BasicFieldDataId:   cast.ToString(issue.DataId),
		consts.BasicFieldIssueId:  issueId,
		consts.BasicFieldWorkHour: issue.LessData[consts.BasicFieldWorkHour],
	}
	e := &commonvo.DataEvent{
		OrgId:     orgId,
		AppId:     issue.AppId,
		ProjectId: issue.ProjectId,
		TableId:   issue.TableId,
		DataId:    issue.DataId,
		IssueId:   issue.Id,
		UserId:    userId,
		New:       newData,
	}

	openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
	openTraceIdStr := cast.ToString(openTraceId)

	// 上报事件
	report.ReportDataEvent(msgPb.EventType_DataWorkHourUpdated, openTraceIdStr, e)
}
