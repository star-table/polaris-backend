package service

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/spf13/cast"
)

func TransferOrg(originOrgId, newOrgId int64) errs.SystemErrorInfo {
	key := consts.TransferOrgLockKey + cast.ToString(originOrgId)
	uuid := uuid.NewUuid()
	suc, err := cache.TryGetDistributedLock(key, uuid)
	if err != nil {
		log.Error("获取锁异常")
		return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	if suc {
		//释放锁
		defer func() {
			if _, e := cache.ReleaseDistributedLock(key, uuid); e != nil {
				log.Error(e)
			}
		}()

		return domain.TransferOrgToOtherPlatform(originOrgId, newOrgId)
	} else {
		return errs.SystemBusy
	}
}

func PrepareTransferOrg(originOrgId, newOrgId int64) (*orgvo.PrepareTransferOrgData, errs.SystemErrorInfo) {
	result := &orgvo.PrepareTransferOrgData{}
	matchInfo, err := domain.GetTransferOrgUserMatchInfo(originOrgId, newOrgId)
	if err != nil {
		return nil, err
	}

	userInfos, err := domain.BatchGetUserDetailInfo(matchInfo.NewUserIds)
	if err != nil {
		return nil, err
	}
	for _, info := range userInfos {
		temp := &orgvo.UserInfo{UserID: info.ID}
		err := copyer.Copy(info, temp)
		if err != nil {
			return nil, errs.JSONConvertError
		}
		if matchInfo.NewToOriginUserIdMap[info.ID] != 0 {
			result.MatchUsers = append(result.MatchUsers, temp)
		} else {
			result.NotMatchUsers = append(result.NotMatchUsers, temp)
		}
	}

	return result, nil
}
