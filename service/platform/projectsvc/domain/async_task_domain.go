package domain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/star-table/common/core/util/json"

	"github.com/star-table/polaris-backend/common/core/util/asyn"

	"github.com/shopspring/decimal"

	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/spf13/cast"
)

// GenAsyncTaskIdForImport 为导入任务的异步任务生成task id
func GenAsyncTaskIdForImport(appId, tableId int64) string {
	return fmt.Sprintf("imp_%d_t%d", appId, tableId)
}

// GenAsyncTaskIdForApplyTemplate 为应用模板的异步任务生成task id
func GenAsyncTaskIdForApplyTemplate(appId int64) string {
	return fmt.Sprintf("%d", appId)
}

// CreateAsyncTask 将新的异步任务信息存入 redis
func CreateAsyncTask(orgId int64, totalCount int64, taskId string, params map[string]string) errs.SystemErrorInfo {
	cacheKey, _ := util.ParseCacheKey(consts.CacheKeyOfAsyncTaskInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:       orgId,
		consts.CacheKeyAsyncTaskIdConstName: taskId,
	})
	values := map[string]string{
		consts.AsyncTaskHashPartKeyOfProcessed: "0",
		consts.AsyncTaskHashPartKeyOfTotal:     cast.ToString(totalCount),
		consts.AsyncTaskHashPartKeyOfStartTime: cast.ToString(time.Now().Unix()),
		consts.AsyncTaskHashPartKeyOfCover:     "",
		consts.AsyncTaskHashPartKeyOfCardSend:  "0",
	}
	for k, v := range params {
		values[k] = v
	}
	log.Infof("[CreateAsyncTask] orgId:%d, taskId:%s, total:%d, %v", orgId, taskId, totalCount, json.ToJsonIgnoreError(values))
	cache.Del(cacheKey)
	err := cache.HMSet(cacheKey, values)
	if err != nil {
		log.Errorf("[CreateAsyncTask] HMSet err: %v, orgId: %d, taskId: %s", err, orgId, taskId)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	_, err = cache.Expire(cacheKey, consts.CacheExpire1Day)
	if err != nil {
		log.Errorf("[CreateAsyncTask] Expire err: %v, cacheKey: %s", err, cacheKey)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}

// GetAsyncTask 获取异步任务的进度信息
func GetAsyncTask(orgId int64, taskId string) (*projectvo.AsyncTask, errs.SystemErrorInfo) {
	cacheKey, _ := util.ParseCacheKey(consts.CacheKeyOfAsyncTaskInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:       orgId,
		consts.CacheKeyAsyncTaskIdConstName: taskId,
	})
	// polaris:orgsvc:org_2373:asyncTask:id_tid_asyncTaskId202207131734
	infoMap, err := cache.HMGet(cacheKey,
		consts.AsyncTaskHashPartKeyOfStartTime,
		consts.AsyncTaskHashPartKeyOfTotal,
		consts.AsyncTaskHashPartKeyOfProcessed,
		consts.AsyncTaskHashPartKeyOfFailed,
		consts.AsyncTaskHashPartKeyOfCover,
		consts.AsyncTaskHashPartKeyOfErrCode,
		consts.AsyncTaskHashPartKeyOfCardSend,
		consts.AsyncTaskHashPartKeyOfTableIds)
	if err != nil {
		log.Errorf("[GetAsyncTask] err: %v, orgId: %d, taskId: %s", err, orgId, taskId)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}
	taskExist := false
	totalNum, processNum := 0, 0
	defaultCover, hasSendCardFlag := "", ""
	numParseMap := map[string]int{
		consts.AsyncTaskHashPartKeyOfTotal:     0,
		consts.AsyncTaskHashPartKeyOfProcessed: 0,
		consts.AsyncTaskHashPartKeyOfFailed:    0,
		//consts.AsyncTaskHashPartKeyOfParentTotalCount: 0,
		//consts.AsyncTaskHashPartKeyOfParentSucCount:   0,
		//consts.AsyncTaskHashPartKeyOfParentErrCount:   0,
	}
	zeroVal := "0"
	tableIdsStr := "[]"
	for hKey, _ := range numParseMap {
		if val, ok := infoMap[hKey]; ok {
			if val == nil {
				val = &zeroVal
			}
			tmpNum, err := strconv.Atoi(*val)
			if err != nil {
				log.Errorf("[GetAsyncTask] Atoi totalNum err: %v, orgId: %d, taskId: %s, hKey: %s, val: %s", err,
					orgId, taskId, hKey, *val)
				return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, err)
			}
			numParseMap[hKey] = tmpNum

			if hKey == consts.AsyncTaskHashPartKeyOfTotal {
				taskExist = true
				totalNum = tmpNum
			}
		}
	}
	processNum = numParseMap[consts.AsyncTaskHashPartKeyOfProcessed]
	if processNum > totalNum {
		numParseMap[consts.AsyncTaskHashPartKeyOfProcessed] = totalNum
	}

	if val, ok := infoMap[consts.AsyncTaskHashPartKeyOfCover]; ok && val != nil {
		defaultCover = *val
	}
	if val, ok := infoMap[consts.AsyncTaskHashPartKeyOfTableIds]; ok && val != nil {
		tableIdsStr = *val
	}
	if val, ok := infoMap[consts.AsyncTaskHashPartKeyOfCardSend]; ok && val != nil {
		hasSendCardFlag = *val
	}
	startTimeStr := strconv.FormatInt(time.Now().Unix(), 10)
	if val, ok := infoMap[consts.AsyncTaskHashPartKeyOfStartTime]; ok && val != nil {
		startTimeStr = *val
	}
	if !taskExist {
		return nil, errs.AsyncTaskNotExist
	}
	defaultErrCode := -1
	if val, ok := infoMap[consts.AsyncTaskHashPartKeyOfErrCode]; ok && val != nil {
		defaultErrCode, err = strconv.Atoi(*val)
		if err != nil {
			log.Errorf("[GetAsyncTask] Atoi defaultErrCode err: %v, orgId: %d, taskId: %s, val: %s", err, orgId, taskId, *val)
			return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, err)
		}
	}
	percentVal := calcAsyncTaskPercentage(processNum, totalNum)
	startTimeInt, err := strconv.ParseInt(startTimeStr, 10, 64)
	if err != nil {
		log.Errorf("[GetAsyncTask] ParseInt err: %v, orgId: %d, taskId: %s, startTimeStr: %s", err, orgId, taskId, startTimeStr)
		return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, err)
	}
	res := &projectvo.AsyncTask{
		Total:          totalNum,
		Processed:      processNum,
		PercentVal:     percentVal,
		Cover:          defaultCover,
		Failed:         numParseMap[consts.AsyncTaskHashPartKeyOfFailed],
		ErrCode:        defaultErrCode,
		StartTimestamp: startTimeInt,
		CardSend:       hasSendCardFlag,
		TableIds:       tableIdsStr,
		//ParentTotalCount: numParseMap[consts.AsyncTaskHashPartKeyOfParentTotalCount],
		//ParentSucCount:   numParseMap[consts.AsyncTaskHashPartKeyOfParentSucCount],
		//ParentErrCount:   numParseMap[consts.AsyncTaskHashPartKeyOfParentErrCount],
	}

	return res, nil
}

func calcAsyncTaskPercentage(numerator, denominator int) float64 {
	if denominator == 0 {
		return -1
	}
	d1 := decimal.NewFromFloat(float64(numerator)).Div(decimal.NewFromFloat(float64(denominator)))
	f, _ := d1.Round(4).Mul(decimal.NewFromFloat(100)).Float64()
	if f > 100 {
		f = 100
	}

	return f
}

func checkAsyncTaskIsRunning(orgId int64, taskId string) bool {
	taskInfo, err := GetAsyncTask(orgId, taskId)
	if err != nil {
		log.Errorf("[checkAsyncTaskIsRunning] GetAsyncTask err: %v, orgId: %d, taskId: %v", err, orgId, taskId)
		return false
	}
	if taskInfo.Total > 0 && taskInfo.PercentVal >= 100 {
		return false
	}
	if taskInfo.Processed+taskInfo.Failed >= taskInfo.Total {
		return false
	}
	return true
}

// CheckAsyncTaskIsRunning 检查是否有异步任务（异步批量创建任务）在执行
func CheckAsyncTaskIsRunning(orgId, appId, tableId int64) bool {
	// 检查是否有导入任务
	if tableId > 0 {
		taskId := GenAsyncTaskIdForImport(appId, tableId)
		if checkAsyncTaskIsRunning(orgId, taskId) {
			return true
		}
	} else {
		// 检查羡慕下所有的表是否有导入任务
		tablesByApp, err := GetTableListMapByProAppIds(orgId, []int64{appId})
		if err != nil {
			log.Errorf("[CheckAsyncTaskIsRunning] GetTableListMapByProAppIds err: %v, appId: %d", err, appId)
			return false
		}
		if tables, ok := tablesByApp[appId]; ok {
			for _, table := range tables {
				taskId := GenAsyncTaskIdForImport(appId, table.TableId)
				if checkAsyncTaskIsRunning(orgId, taskId) {
					return true
				}
			}
		}
	}

	// 检查是否有应用模板任务
	taskId := GenAsyncTaskIdForApplyTemplate(appId)
	if checkAsyncTaskIsRunning(orgId, taskId) {
		return true
	}

	return false
}

///////////////////////////////////////////////////////////////////////

//// SetAsyncTaskIdToIssue 向 issue 信息中，赋值异步任务的 id，用于后续区分其是否是某一次异步任务的子任务
//// idSeeder 异步任务的 id 值提供者，如果不存在则调用 GenAsyncTaskIdForIssue 生成
//// opUserId 操作人，当前操作用户
//// 返回异步任务总数量
//func SetAsyncTaskIdToIssue(issues *[]vo.CreateIssueReq, idSeeder string, opUserId int64, isApplyTpl bool) (allCnt int, parentCnt int) {
//	asyncTaskId := idSeeder
//	count := 0
//	for i, issue := range *issues {
//		if issue.ExtraInfo == nil {
//			(*issues)[i].ExtraInfo = make(map[string]interface{}, 3)
//		}
//		(*issues)[i].ExtraInfo[consts.CacheKeyAsyncTaskIdConstName] = asyncTaskId
//		(*issues)[i].ExtraInfo[consts.HelperFieldOperateUserId] = opUserId
//		(*issues)[i].ExtraInfo[consts.HelperFieldIsApplyTpl] = isApplyTpl
//		count += 1
//		parentCnt += 1
//
//		SetAsyncTaskIdToChildIssue((*issues)[i].Children, asyncTaskId, opUserId, isApplyTpl, &count)
//	}
//	allCnt = count
//
//	return
//}

//// SetAsyncTaskIdToChildIssue 向子任务中追加 extraInfo 信息
//func SetAsyncTaskIdToChildIssue(children []*vo.IssueChildren, asyncTaskId string, opUserId int64, isApplyTpl bool, count *int) {
//	if children == nil || len(children) < 1 {
//		return
//	}
//	for _, child := range children {
//		if child.ExtraInfo == nil {
//			child.ExtraInfo = make(map[string]interface{}, 3)
//		}
//		child.ExtraInfo[consts.CacheKeyAsyncTaskIdConstName] = asyncTaskId
//		child.ExtraInfo[consts.HelperFieldOperateUserId] = opUserId
//		child.ExtraInfo[consts.HelperFieldIsApplyTpl] = isApplyTpl
//		*count += 1
//
//		if len(child.Children) > 0 {
//			SetAsyncTaskIdToChildIssue(child.Children, asyncTaskId, opUserId, isApplyTpl, count)
//		}
//	}
//}

// UpdateAsyncTaskWithSucc 异步任务完成时，向任务信息中更新进度
func UpdateAsyncTaskWithSucc(orgId, projectId, tableId, userId int64, taskId string, isApplyTemplate bool, count int64) {
	if taskId == "" {
		//log.Infof("[UpdateAsyncTaskWithSucc] get taskId err taskId: %s", taskId)
		return
	}
	log.Infof("[UpdateAsyncTaskWithSucc] orgId:%d, projectId:%d, tableId:%d, userId:%d, taskId:%s, count:%d",
		orgId, projectId, tableId, userId, taskId, count)

	lockKey := consts.CreateIssueAsyncTaskLockKey + taskId
	uid := uuid.NewUuid()
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("[UpdateAsyncTaskWithSucc] TryGetDistributedLock err: %v", err)
		return
	} else if !suc {
		log.Errorf("[UpdateAsyncTaskWithSucc] TryGetDistributedLock failed err: %v", errs.GetDistributedLockError)
		return
	}
	defer func() {
		if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
			log.Error(err)
		}
	}()

	errSys := updateAsyncTaskProcess(orgId, taskId, consts.AsyncTaskHashPartKeyOfProcessed, count)
	if errSys != nil {
		log.Errorf("[UpdateAsyncTaskWithSucc] updateAsyncTaskProcess err: %v, taskId: %s", errSys, taskId)
		return
	}

	// 检查数量是否已经处理完（无论成功或失败）
	taskInfo, errSys := GetAsyncTask(orgId, taskId)
	if errSys != nil {
		log.Errorf("[UpdateAsyncTaskWithSucc] GetAsyncTask err: %v, taskId: %s", errSys, taskId)
		return
	}
	if taskInfo.Failed+taskInfo.Processed >= taskInfo.Total {
		asyn.Execute(func() {
			if errSys = AsyncTaskCompleteSendCard(orgId, projectId, tableId, userId, taskId, isApplyTemplate, taskInfo); errSys != nil {
				log.Errorf("[UpdateAsyncTaskWithSucc] AsyncTaskCompleteSendCard err: %v, taskId: %s", err, taskId)
			}
		})
	}
}

// UpdateAsyncTaskWithError 异步任务异常时，将异常 code 放入缓存信息中，并更新缓存中失败的计数
func UpdateAsyncTaskWithError(orgId, projectId, tableId, userId int64, taskId string, isApplyTemplate bool, errCode errs.SystemErrorInfo, count int64) {
	if taskId == "" {
		//log.Infof("[UpdateAsyncTaskWithError] get taskId err taskId: %s", taskId)
		return
	}

	log.Infof("[UpdateAsyncTaskWithError] orgId:%d, projectId:%d, tableId:%d, userId:%d, taskId:%s, count:%d, err:%v",
		orgId, projectId, tableId, userId, taskId, count, errCode)

	cacheKey, _ := util.ParseCacheKey(consts.CacheKeyOfAsyncTaskInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:       orgId,
		consts.CacheKeyAsyncTaskIdConstName: taskId,
	})
	if exist, err := cache.Exist(cacheKey); err != nil {
		log.Errorf("[UpdateAsyncTaskWithError] Exist err: %v, orgId: %d, taskId: %s", err, orgId, taskId)
		return
	} else if !exist {
		log.Errorf("[UpdateAsyncTaskWithError] async task not exist. taskId: %s", taskId)
		return
	}

	// 保存错误码
	err := cache.HSet(cacheKey, consts.AsyncTaskHashPartKeyOfErrCode, strconv.Itoa(errCode.Code()))
	if err != nil {
		log.Errorf("[UpdateAsyncTaskWithError] HSet err: %v, orgId: %d, taskId: %s", err, orgId, taskId)
		return
	}

	lockKey := consts.CreateIssueAsyncTaskLockKey + taskId
	uid := uuid.NewUuid()
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("[UpdateAsyncTaskWithError] TryGetDistributedLock err: %v", err)
		return
	} else if !suc {
		log.Errorf("[UpdateAsyncTaskWithError] TryGetDistributedLock failed err: %v", errs.GetDistributedLockError)
		return
	}
	defer func() {
		if _, err = cache.ReleaseDistributedLock(lockKey, uid); err != nil {
			log.Error(err)
		}
	}()

	// 记录出错的数量
	errSys := updateAsyncTaskProcess(orgId, taskId, consts.AsyncTaskHashPartKeyOfFailed, count)
	if errSys != nil {
		log.Errorf("[UpdateAsyncTaskWithError] updateAsyncTaskProcess err: %v, taskId: %s", errSys, taskId)
		return
	}

	// 检查任务是否完成，完成则发通知卡片
	taskInfo, errSys := GetAsyncTask(orgId, taskId)
	if errSys != nil {
		log.Errorf("[UpdateAsyncTaskWithError] GetAsyncTask err: %v, taskId: %s", errSys, taskId)
		return
	}
	if taskInfo.Failed+taskInfo.Processed >= taskInfo.Total {
		asyn.Execute(func() {
			if errSys = AsyncTaskCompleteSendCard(orgId, projectId, tableId, userId, taskId, isApplyTemplate, taskInfo); errSys != nil {
				log.Errorf("[UpdateAsyncTaskWithError] AsyncTaskCompleteSendCard err: %v, taskId: %s", errSys, taskId)
			}
		})
	}

	return
}

// updateAsyncTaskProcess 异步任务处理完状态统计递增
func updateAsyncTaskProcess(orgId int64, taskId string, key string, count int64) errs.SystemErrorInfo {
	cacheKey, _ := util.ParseCacheKey(consts.CacheKeyOfAsyncTaskInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:       orgId,
		consts.CacheKeyAsyncTaskIdConstName: taskId,
	})
	_, err := cache.HINCRBY(cacheKey, key, count)
	if err != nil {
		log.Errorf("[updateAsyncTaskProcess] err: %v, orgId: %d, taskId: %s", err, orgId, taskId)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}

func ClearAsyncTaskCacheForCreateIssue(orgId int64, taskId string) errs.SystemErrorInfo {
	cacheKey, _ := util.ParseCacheKey(consts.CacheKeyOfAsyncTaskInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:       orgId,
		consts.CacheKeyAsyncTaskIdConstName: taskId,
	})
	_, err := cache.Del(cacheKey)
	if err != nil {
		log.Errorf("[ClearAsyncTaskCacheForCreateIssue] Del err: %v", err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}

	return nil
}

// AsyncTaskCompleteSendCard 异步任务完成后，发送卡片
func AsyncTaskCompleteSendCard(orgId, projectId, tableId, userId int64, taskId string, isApplyTpl bool, taskInfo *projectvo.AsyncTask) errs.SystemErrorInfo {

	return nil
}

// SetAsyncTaskHasSendCard 异步任务处理完成，发送卡片后，标记已发送卡片
func SetAsyncTaskHasSendCard(orgId int64, taskId string) errs.SystemErrorInfo {
	cacheKey, _ := util.ParseCacheKey(consts.CacheKeyOfAsyncTaskInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:       orgId,
		consts.CacheKeyAsyncTaskIdConstName: taskId,
	})
	err := cache.HSet(cacheKey, consts.AsyncTaskHashPartKeyOfCardSend, "1")
	if err != nil {
		log.Errorf("[SetAsyncTaskHasSendCard] HSet err: %v, orgId: %d, taskId: %s", err, orgId, taskId)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}
