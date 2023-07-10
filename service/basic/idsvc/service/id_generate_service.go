package service

import (
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/times"
	uuid2 "github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/model/bo"
	sconsts "github.com/star-table/polaris-backend/service/basic/idsvc/consts"
	"github.com/star-table/polaris-backend/service/basic/idsvc/po"
	"strconv"
	"upper.io/db.v3"
)

const IdCodeSpan = "-"

const IdMaxMultipleCount = 100

var log = logger.GetDefaultLogger()

func ApplyPrimaryId(code string) (int64, errs.SystemErrorInfo) {
	//return ApplyId(0, code, "")
	id, err := ApplyMultipleId(0, code, "", 1)
	if err != nil {
		return 0, err
	}

	if id != nil && len(id.Ids) > 0 {
		return id.Ids[0].Id, err
	}
	return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError)

}

func ApplyCode(orgId int64, code string, preCode string) (string, errs.SystemErrorInfo) {
	//return ApplyId(0, code, "")
	id, err := ApplyMultipleId(orgId, code, preCode, 1)
	if err != nil {
		return "", err
	}

	if id != nil && len(id.Ids) > 0 {
		return id.Ids[0].Code, err
	}
	return "", errs.BuildSystemErrorInfo(errs.ApplyIdError)

}

func ApplyId(orgId int64, code string, preCode string) (*bo.IdCodes, errs.SystemErrorInfo) {
	//return ApplyId(0, code, "")
	return ApplyMultipleId(orgId, code, preCode, 1)
}

//返回id集合的指针地址和错误信息
func ApplyMultipleId(orgId int64, code string, preCode string, count int64) (*bo.IdCodes, errs.SystemErrorInfo) {
	// 参数校验
	paramErrors := checkApplyMultipleIdParam(orgId, code, count)

	if paramErrors != nil {
		return nil, paramErrors
	}

	realCode := code
	if orgId > 0 {
		realCode += IdCodeSpan
	}

	idKey, idObjKey := buildIdCacheKey(orgId, realCode)

	existErrors := checkCacheExist(orgId, realCode, idKey)

	if existErrors != nil {
		return nil, existErrors
	}

	//cache自增id
	newId, err2 := cache.Incrby(idKey, count)
	if err2 != nil {
		log.Error(idKey + strs.ObjectToString(err2))
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err2)
	}

	idInfos := []bo.IdCodeInfo{}

	for i := int64(1); i <= count; i++ {
		idInfo := bo.IdCodeInfo{
			Id: newId - count + i,
		}
		idInfo.Code = realCode + strconv.FormatInt(idInfo.Id, 10)
		idInfos = append(idInfos, idInfo)
	}
	idCodes := bo.IdCodes{
		OrgId:   orgId,
		Code:    code,
		PreCode: preCode,
		Ids:     idInfos,
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()
		applyErr := checkApplyIdThreshold(orgId, realCode, idObjKey, newId)
		if applyErr != nil{
			log.Error(applyErr)
		}
	}()

	return &idCodes, nil
}

//返回结果和判断是否等待
func waitOrReturn(tt bool, err6 errs.SystemErrorInfo, idCodes bo.IdCodes, i int) (bool, *bo.IdCodes, errs.SystemErrorInfo) {

	if err6 != nil {
		log.Error(err6)
		return true, nil, err6
	}
	if tt == true {
		return true, &idCodes, nil
	} else if tt == false && i >= 8 {
		// 扩容等待过长
		return true, nil, errs.BuildSystemErrorInfo(errs.ApplyIdError)
	} else {
		// 循环等待
		times.SleepMillisecond(100)
	}

	return false, nil, err6
}

//参数校验
func checkApplyMultipleIdParam(orgId int64, code string, count int64) errs.SystemErrorInfo {
	// 参数校验
	if orgId < 0 || "" == code {
		return errs.BuildSystemErrorInfo(errs.ParamError)
	}
	if count > IdMaxMultipleCount {
		return errs.BuildSystemErrorInfo(errs.ApplyIdCountTooMany)
	}

	return nil
}

//判断缓存是否存在
func checkCacheExist(orgId int64, realCode string, idKey string) errs.SystemErrorInfo {
	t, err := cache.Exist(idKey)

	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}
	if t == false {
		// id缓存不存在
		err2 := buildIdObject2Cache(orgId, realCode, 0)
		if err2 != nil {
			log.Error(err2)
			return err2
		}
	}

	return nil
}

func checkApplyIdThreshold(orgId int64, code string, idObjKey string, newId int64) errs.SystemErrorInfo {
	t, err := cache.Get(idObjKey)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}
	if t == "" {
		buildIdObject2Cache(orgId, code, 0)
		return nil
	}

	objCache := &bo.ObjectIdCache{}
	err2 := json.FromJson(t, objCache)
	if err2 != nil {
		log.Error(err)
		cache.Del(idObjKey)
		return errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}

	if newId > objCache.Threshold {
		// 需要扩展下阶段id
		buildIdObject2Cache(orgId, code, objCache.Threshold)
	}

	return nil
}

func buildIdCacheKey(orgId int64, code string) (string, string) {
	key := strconv.FormatInt(orgId, 10) + IdCodeSpan + code
	idKey := sconsts.CacheObjectIdPreKey + key
	idObjKey := sconsts.CacheObjectIdPreKey + "obj:" + key
	return idKey, idObjKey
}

func buildIdObject2Cache(orgId int64, code string, thresholdOldValue int64) errs.SystemErrorInfo {

	idKey, idObjKey := buildIdCacheKey(orgId, code)

	uuid := uuid2.NewUuid()

	lockKey := consts.IdServiceIdKey + idKey
	tp, err := cache.TryGetDistributedLock(lockKey, uuid)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	}
	if tp == false {
		return errs.BuildSystemErrorInfo(errs.GetDistributedLockError)
	}
	defer cache.ReleaseDistributedLock(lockKey, uuid)

	t, err := cache.Get(idObjKey)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}

	if thresholdOldValue == 0 && t != "" {
		// 初始化id策略，且已初始化过
		return nil
	}

	if t != "" {
		objCache := &bo.ObjectIdCache{}
		err2 := json.FromJson(t, objCache)
		if err2 != nil {
			log.Error(err)
			//格式不正确，清除缓存
			cache.Del(idObjKey)
			return errs.BuildSystemErrorInfo(errs.ApplyIdError)
		}

		if objCache.Threshold != thresholdOldValue {
			// 已完成延续下阶段，返回
			return nil
		}
	}

	objectIdObj, err5 := getObjectIdInfoByDb(orgId, code)

	if err5 != nil {
		log.Error(err5)
		return err5
	}

	if thresholdOldValue == 0 {
		// 初始化阶段需要写入maxId
		cache.Set(idKey, strconv.FormatInt(objectIdObj.MaxId, 10))
	}

	threshold := config.GetParameters().IdBufferThreshold

	objIdCache := bo.ObjectIdCache{
		OrgId: orgId,
		Code:  objectIdObj.Code,
		MaxId: objectIdObj.MaxId + int64(objectIdObj.Step),
	}
	objIdCache.Threshold = objectIdObj.MaxId + int64(float64(objectIdObj.Step)*threshold)

	jsonstr, err3 := json.ToJson(objIdCache)
	if err3 != nil {
		log.Error(err3)
		return errs.BuildSystemErrorInfo(errs.JSONConvertError, err3)
	}
	err4 := cache.Set(idObjKey, jsonstr)
	if err4 != nil {
		log.Error(err4)
		return errs.BuildSystemErrorInfo(errs.JSONConvertError, err4)
	}
	return nil
}

func getObjectIdInfoByDb(orgId int64, code string) (*po.PpmBasObjectId, errs.SystemErrorInfo) {
	objectId := &po.PpmBasObjectId{}

	cond := db.Cond{
		consts.TcOrgId: orgId,
		consts.TcCode:  code,
	}

	err := mysql.SelectOneByCond(objectId.TableName(), cond, objectId)
	if err != nil {
		// id表未查到对象，需要新建，该表为自增id，降低分布式锁时间
		initMaxId := consts.DefaultObjectIdMaxId + consts.DefaultObjectIdStep
		if orgId > 0 {
			// orgId 大于 0，编号从1开始
			initMaxId = consts.DefaultObjectIdStep
		}
		objectId = &po.PpmBasObjectId{
			OrgId:    orgId,
			Code:     code,
			MaxId:    initMaxId,
			Step:     consts.DefaultObjectIdStep, //之后从字典获取
			IsDelete: consts.AppIsNoDelete,
		}
		id, err2 := mysql.InsertReturnId(objectId)
		if err2 != nil {
			log.Error(err2)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		objectId.Id = id.(int64)

	} else {
		// 找到对象，需要更新最大值
		objectId.MaxId = objectId.MaxId + int64(objectId.Step)
		err = mysql.Update(objectId)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}

	objectId.MaxId = objectId.MaxId - int64(objectId.Step)

	return objectId, nil
}
