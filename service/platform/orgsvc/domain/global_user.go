package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

func GetGlobalUserRelationsByUserId(userId int64) (*bo.GlobalUserRelations, errs.SystemErrorInfo) {
	result := &bo.GlobalUserRelations{}
	globalId, err := dao.GetGlobalUserRelation().GetGlobalUserIdByUserId(userId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return result, nil
		}
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	relations, err := getGlobalUserRelations(globalId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return result, nil
		}
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return relations, nil
}

func GetGlobalUserByMobile(mobile string) (*po.PpmOrgGlobalUser, error) {
	return dao.GetGlobalUser().GetGlobalUserByMobile(mobile)
}

func GetGlobalUserByUserId(userId int64) (*po.PpmOrgGlobalUser, errs.SystemErrorInfo) {
	result := &po.PpmOrgGlobalUser{}
	globalId, err := dao.GetGlobalUserRelation().GetGlobalUserIdByUserId(userId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return result, nil
		}
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	globalUser, err := dao.GetGlobalUser().GetGlobalUserById(globalId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return globalUser, nil
}

func GetGlobalUserRelationsByMobile(mobile string) (*bo.GlobalUserRelations, errs.SystemErrorInfo) {
	result := &bo.GlobalUserRelations{Mobile: mobile}
	globalUser, err := dao.GetGlobalUser().GetGlobalUserByMobile(mobile)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return result, nil
		}
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	result.GlobalUserId = globalUser.Id

	userIds, err := dao.GetGlobalUserRelation().GetUserIdsByGlobalUserId(globalUser.Id)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return result, nil
		}
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	result.BindUserIds = userIds

	return result, nil
}

func getGlobalUserRelations(globalUserId int64) (*bo.GlobalUserRelations, error) {
	userIds, err := dao.GetGlobalUserRelation().GetUserIdsByGlobalUserId(globalUserId)
	if err != nil {
		return nil, err
	}

	globalUser, err := dao.GetGlobalUser().GetGlobalUserById(globalUserId)
	if err != nil {
		return nil, err
	}

	return &bo.GlobalUserRelations{
		GlobalUserId: globalUser.Id,
		Mobile:       globalUser.Mobile,
		BindUserIds:  userIds,
	}, nil
}

// CheckUserOrgIsConflict 检查用户组织是否冲突，先用userId查到对应的globalUserId，然后再通过globalUserId查到绑定的所有组织
func CheckUserOrgIsConflict(userId, orgId int64) (bool, errs.SystemErrorInfo) {
	_, err := GetOrgUserInfo(userId, orgId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return false, nil
		}
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return true, nil
}

// GetOrgUserInfo 这个userId其实有可能并不属于orgId,所以需要先查询这个userId绑定的globalUserId下的所有userIds，看看到底这个org属于哪个userId
func GetOrgUserInfo(userId, orgId int64) (*po.PpmOrgUserOrganization, error) {
	globalUserId, err := dao.GetGlobalUserRelation().GetGlobalUserIdByUserId(userId)
	if err != nil {
		return nil, err
	}

	userIds, err := dao.GetGlobalUserRelation().GetUserIdsByGlobalUserId(globalUserId)
	if err != nil {
		return nil, err
	}

	userOrg, err := getUserOrgInfoByUserIdsAndOrgId(userIds, orgId)
	if err != nil {
		return nil, err
	}

	return userOrg, nil
}

// UpdateGlobalUserLastLoginInfoWithCheckGlobal 不确定这个组织是否是这个uid的，所以需要用globalUserId所绑定的所有userId去查询，看看这个orgId属于哪个userId，再更新
func UpdateGlobalUserLastLoginInfoWithCheckGlobal(userId, orgId int64) (int64, errs.SystemErrorInfo) {
	orgUserInfo, err := GetOrgUserInfo(userId, orgId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return userId, nil
		}
		return userId, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return orgUserInfo.UserId, UpdateGlobalUserLastLoginInfo(orgUserInfo.UserId, orgId)
}

// UpdateGlobalUserLastLoginInfo 创建组织的时候确定这个orgId就是这个用户
func UpdateGlobalUserLastLoginInfo(userId, orgId int64) errs.SystemErrorInfo {
	globalUserId, err := dao.GetGlobalUserRelation().GetGlobalUserIdByUserId(userId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil
		}
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err = dao.GetGlobalUser().UpdateGlobalLastLoginInfo(globalUserId, userId, orgId)
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

// GetUserIdsMobileMap 根据userIds获取手机号map
func GetUserIdsMobileMap(userIds []int64) (map[int64]string, errs.SystemErrorInfo) {
	globalToUserIdsMap, err := dao.GetGlobalUserRelation().GetGlobalUserIdsMapByUserIds(userIds)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return map[int64]string{}, nil
		}
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	if len(globalToUserIdsMap) > 0 {
		globalIds := make([]int64, 0, len(globalToUserIdsMap))
		for globalId := range globalToUserIdsMap {
			globalIds = append(globalIds, globalId)
		}

		mobilesMap, err := dao.GetGlobalUser().GetMobilesMapByIds(globalIds)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}

		result := make(map[int64]string, len(mobilesMap))
		for globalUserId, mobile := range mobilesMap {
			for _, userId := range globalToUserIdsMap[globalUserId] {
				result[userId] = mobile
			}
		}
		return result, nil
	}

	return map[int64]string{}, nil
}

//loginName为允许为账号，邮箱，手机号
func GetUserInfoByPwd(loginName string, pwd string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	globalUser, err := dao.GetGlobalUser().GetGlobalUserByMobile(loginName)
	if err != nil {
		if err == db.ErrNoMoreRows {
			// 可能是账号名登录，再查一遍
			return GetUserInfoByLoginNameAndPwd(loginName, pwd)
		} else {
			return nil, errs.MysqlOperateError
		}
	}
	inputPwd := pwd
	salt := globalUser.PasswordSalt
	pwd = util.PwdEncrypt(pwd, salt)

	if globalUser.Password != pwd {
		// 尝试使用无码的校验方案
		if !CheckIsLessCodeAccount(globalUser.Password, globalUser.PasswordSalt, loginName, inputPwd) {
			return nil, errs.PwdLoginUsrOrPwdNotMatch
		}
	}

	userPo, err := GetUserInfoByUserId(globalUser.LastLoginUserId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.UserNotExist
		}
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	if globalUser.LastLoginOrgId != 0 {
		userPo.OrgId = globalUser.LastLoginOrgId
	}

	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}
