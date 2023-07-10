package domain

import (
	"fmt"
	"strconv"
	"time"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/md5"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/pinyin"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/mqbo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/msgfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func UserInit(orgId int64, corpId string, outUserId string, sourceChannel string) (*bo.BaseUserInfoBo, errs.SystemErrorInfo) {
	return nil, nil
}

func GetOutUserName(orgId, userId int64, sourceChannel string) (string, errs.SystemErrorInfo) {
	//用户外部信息
	userOutInfo := &po.PpmOrgUserOutInfo{}
	err := mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcUserId:        userId,
	}, userOutInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return "", nil
		} else {
			log.Error(err)
			return "", errs.MysqlOperateError
		}
	}

	return userOutInfo.Name, nil
}

func GetUserBo(userId int64) (*bo.UserInfoBo, bool, errs.SystemErrorInfo) {
	user := &po.PpmOrgUser{}
	err1 := mysql.SelectOneByCond(user.TableName(), db.Cond{
		consts.TcId:       userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, user)
	if err1 != nil {
		log.Error(strs.ObjectToString(err1))
		return nil, false, errs.BuildSystemErrorInfo(errs.UserNotFoundError, err1)
	}

	userInfo := &bo.UserInfoBo{}
	err1 = copyer.Copy(user, userInfo)
	if err1 != nil {
		return nil, false, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	//passwordSet := false
	//if user.Password != consts.BlankString {
	//	passwordSet = true
	//}
	////是否需要绑定提醒
	//if userInfo.Mobile != consts.BlankString {
	//	userInfo.RemindBindPhone = consts.AppIsNotRemind
	//}
	return userInfo, false, nil
}

func GetUserInfo(orgId int64, userId int64, sourceChannel string) (*bo.UserInfoBo, bool, errs.SystemErrorInfo) {
	if orgId == 0 {
		log.Error(errs.OrgNotInitError)
		return nil, false, errs.OrgNotInitError
	}
	baseUserInfo, err := GetBaseUserInfo(orgId, userId)
	if err != nil {
		log.Error(err)
		return nil, false, errs.BuildSystemErrorInfo(errs.UserNotInitError)
	}

	ownerInfo, passwordSet, err := GetUserBo(userId)
	if err != nil {
		log.Error(err)
		return nil, false, err
	}
	//用户的sourceChannel从当前组织获取
	orgInfo, orgInfoErr := GetOrgBoById(orgId)
	if orgInfoErr != nil {
		log.Error(orgInfoErr)
		return nil, false, orgInfoErr
	}
	ownerInfo.SourceChannel = orgInfo.SourceChannel

	//获取第三方名称
	if ownerInfo.SourceChannel != "" && orgId > 0 {
		outName, err := GetOutUserName(orgId, userId, ownerInfo.SourceChannel)
		if err != nil {
			log.Error(err)
			return nil, false, err
		}
		ownerInfo.ThirdName = outName
	}
	//部分属性覆盖
	ownerInfo.EmplID = &baseUserInfo.OutUserId
	ownerInfo.OrgID = orgId

	if orgId > 0 {
		baseOrgInfo, err := GetBaseOrgInfo(orgId)
		if err != nil {
			log.Error(err)
			return nil, false, errs.BuildSystemErrorInfo(errs.OrgNotInitError)
		}
		ownerInfo.OrgName = baseOrgInfo.OrgName
	}

	globalUser, err := GetGlobalUserByUserId(userId)
	if err != nil {
		log.Errorf("[GetUserInfo] GetGlobalUserByUserId userId:%v, err:%v", userId, err)
		return nil, false, err
	}
	ownerInfo.Mobile = globalUser.Mobile
	if globalUser.Password != consts.BlankString {
		passwordSet = true
	}
	//是否需要绑定提醒
	if globalUser.Mobile != consts.BlankString {
		ownerInfo.RemindBindPhone = consts.AppIsNotRemind
	} else {
		ownerInfo.RemindBindPhone = consts.AppIsRemind
	}

	//这里先默认写死
	ownerInfo.Rimanente = 10
	ownerInfo.Level = 1
	ownerInfo.LevelName = "试用级别"
	ownerInfo.SourceChannel = sourceChannel // 飞书用户 现在可以创建本地组织，导致user表的sourceChannel和org的不一致
	// 如果第三方用户id为空，则证明是用户是本地用户
	if ownerInfo.EmplID == nil || *ownerInfo.EmplID == "" {
		ownerInfo.SourceChannel = consts.AppSourcePlatformLarkXYJH2019
	}

	return ownerInfo, passwordSet, nil
}

func GetOutUserInfoListBySourceChannel(sourceChannel string, page int, size int) ([]bo.UserOutInfoBo, errs.SystemErrorInfo) {
	userOutInfos := &[]po.PpmOrgUserOutInfo{}
	_, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOutInfo, db.Cond{
		consts.TcSourceChannel: sourceChannel,
		consts.TcStatus:        consts.AppStatusEnable,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, nil, page, size, "org_id asc", userOutInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	userOutInfoBos := &[]bo.UserOutInfoBo{}
	err = copyer.Copy(userOutInfos, userOutInfoBos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *userOutInfoBos, nil
}

// 批量查询用户外部信息
func GetOutUserInfoListByUserIds(idList []int64) ([]bo.UserOutInfoBo, errs.SystemErrorInfo) {
	userOutInfos := &[]po.PpmOrgUserOutInfo{}
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcUserId:   db.In(idList),
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, userOutInfos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	userOutInfoBos := &[]bo.UserOutInfoBo{}
	err = copyer.Copy(userOutInfos, userOutInfoBos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *userOutInfoBos, nil
}

func GetOutUserInfoListByOrgId(orgId int64) ([]*po.PpmOrgUserOutInfo, errs.SystemErrorInfo) {
	userOutInfos := make([]*po.PpmOrgUserOutInfo, 0, 100)
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userOutInfos)
	if err != nil {
		log.Errorf("[GetOutUserInfoListByOrgId] orgId:%v, err:%v", orgId, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return userOutInfos, nil
}

// 批量获取用户的 out info
//func GetOutUserInfo(outUserId, sourceChannel string) (*bo.UserOutInfoBo, errs.SystemErrorInfo) {
//	userOutInfo := &po.PpmOrgUserOutInfo{}
//	err := mysql.SelectOneByCond(userOutInfo.TableName(), db.Cond{
//		consts.TcOutUserId:     outUserId,
//		consts.TcSourceChannel: sourceChannel,
//		consts.TcIsDelete:      consts.AppIsNoDelete}, userOutInfo)
//	if err != nil {
//		log.Error(err)
//		return nil, errs.BuildSystemErrorInfo(errs.UserOutInfoNotExist)
//	}
//	userOutInfoBo := &bo.UserOutInfoBo{}
//	_ = copyer.Copy(userOutInfo, userOutInfoBo)
//	return userOutInfoBo, nil
//}

func GetUserInfoListByOrg(orgId int64) ([]bo.SimpleUserInfoBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	userInfos := &[]po.PpmOrgUser{}

	userAlias := "user."
	orgAlias := "org."
	selectErr := conn.Select("user.name", "user.id", "user.status").
		From(consts.TableUser+" user", consts.TableUserOrganization+" org").
		Where(db.Cond{
			userAlias + consts.TcId:       db.Raw(orgAlias + consts.TcUserId),
			orgAlias + consts.TcIsDelete:  consts.AppIsNoDelete,
			userAlias + consts.TcIsDelete: consts.AppIsNoDelete,
			orgAlias + consts.TcStatus:    consts.AppIsInitStatus,
			orgAlias + consts.TcOrgId:     orgId,
		}).
		All(userInfos)
	if selectErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
	}
	userInfoBos := &[]bo.SimpleUserInfoBo{}
	copyErr := copyer.Copy(userInfos, userInfoBos)

	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *userInfoBos, nil
}

//如果err不等于空，说明用户未注册
func GetUserInfoByMobile(phoneNumber string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	globalUser, err := dao.GetGlobalUser().GetGlobalUserByMobile(phoneNumber)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.UserNotExist
		}
		log.Errorf("[GetUserInfoByMobile] GetGlobalUserByMobile phone:%v, err:%v", phoneNumber, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return GetUserInfoByGlobalUser(globalUser)
}

func GetUserInfoByGlobalUser(globalUser *po.PpmOrgGlobalUser) (*bo.UserInfoBo, errs.SystemErrorInfo) {
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
	userBo.LoginName = globalUser.Mobile

	return userBo, nil
}

func GetCurrentUserInfoByUserId(userId int64) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	globalUser, err := GetGlobalUserByUserId(userId)
	if err != nil {
		log.Errorf("[GetCurrentUserInfoByUserId] GetGlobalUserByUserId userId:%v, err:%v", userId, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if globalUser.Id == 0 {
		userPo, err := GetUserInfoByUserId(userId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		userInfoBo := &bo.UserInfoBo{}
		copyErr := copyer.Copy(userPo, userInfoBo)
		if copyErr != nil {
			log.Error(copyErr)
		}
		return userInfoBo, nil
	}

	return GetUserInfoByGlobalUser(globalUser)
}

// DetectUserInfoInUser 检查是否已经在 user 表中存在的记录、是否已经和组织关联
// 返回值 needRegister 表示需要重新开始注册；
// 返回值 needRelate 表示只需要与组织进行关联
// 新版账户体系下，登录名即手机号（包含地区代码），因而 loginPhones 即 loginNames，而实际上表中的 login_name 不直接使用。
// 由于统一账户逻辑，一个手机号可能绑定多个user，这个时候要选一个user去跟orgId绑定，默认选第一个，如果有一个user绑定了该org，则globalUser下的所有user都要排除
func DetectUserInfoInUser(orgId int64, loginPhones []string) (needRegister []string, needRelateUserIdsMap, needResetCheckStatusUserIdsMap map[int64]string, retErr errs.SystemErrorInfo) {
	loginNames := loginPhones
	if len(loginNames) < 1 {
		return
	}
	phoneToUserIdsMap, err := getUserIdsMapByPhones(loginNames)
	if err != nil {
		log.Errorf("[GetExistUserByPhones] getUserIdsByPhones, phoneNumbers:%v, err:%v ", loginNames, err)
		retErr = err
		return
	}

	maybeToRelateUserIds := make([]int64, 0)
	userIdToPhoneMap := make(map[int64]string, len(phoneToUserIdsMap))
	for _, loginPhone := range loginNames {
		if ids, exist := phoneToUserIdsMap[loginPhone]; !exist {
			needRegister = append(needRegister, loginPhone)
		} else {
			maybeToRelateUserIds = append(maybeToRelateUserIds, ids...)
			for _, id := range ids {
				userIdToPhoneMap[id] = loginPhone
			}
		}
	}

	// 查询 maybeToRelate 是否存在关联
	if len(maybeToRelateUserIds) > 0 {
		userOrgPoList := make([]po.PpmOrgUserOrganization, 0)
		oriErr := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcUserId:   db.In(maybeToRelateUserIds),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, &userOrgPoList)
		if oriErr != nil {
			log.Error(oriErr)
			retErr = errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
			return
		}

		// 排除己绑定了org的对应下的所有user
		excludeUserIds := make(map[int64]bool, len(maybeToRelateUserIds))
		// 审核不通过的user需要重置审核状态
		needResetCheckStatusIds := make(map[int64]bool, len(maybeToRelateUserIds))
		for _, organization := range userOrgPoList {
			for _, userId := range phoneToUserIdsMap[userIdToPhoneMap[organization.UserId]] {
				excludeUserIds[userId] = true
			}
			if organization.CheckStatus == consts.AppCheckStatusFail {
				needResetCheckStatusIds[organization.UserId] = true
			}
		}

		// 获取剩下的需要绑定的用户手机
		needRelatePhone := make(map[string]bool, len(maybeToRelateUserIds))
		for _, tmpUserId := range maybeToRelateUserIds {
			if !excludeUserIds[tmpUserId] && !needResetCheckStatusIds[tmpUserId] {
				needRelatePhone[userIdToPhoneMap[tmpUserId]] = true
			}
		}

		needRelateUserIdsMap = make(map[int64]string, len(needRelatePhone))
		for phone := range needRelatePhone {
			if len(phoneToUserIdsMap[phone]) > 0 {
				needRelateUserIdsMap[phoneToUserIdsMap[phone][0]] = phone
			}
		}

		needResetCheckStatusUserIdsMap = make(map[int64]string)
		for uId := range needResetCheckStatusIds {
			if phone, ok := userIdToPhoneMap[uId]; ok {
				needResetCheckStatusUserIdsMap[uId] = phone
			}
		}
	}
	return
}

// GetExistUserByPhones 通过多个手机号（即登录名）查询已经存在的部分的用户
// 通过手机号判断用户是否存在，并且是否关联对应的组织
func GetExistUserByPhones(orgId int64, phoneNumbers []string) ([]bo.UserInfoBo, errs.SystemErrorInfo) {
	userBoList := make([]bo.UserInfoBo, 0)
	userPoList := make([]po.PpmOrgUser, 0)
	if len(phoneNumbers) < 1 {
		return userBoList, nil
	}

	userIdsMap, err := getUserIdsMapByPhones(phoneNumbers)
	if err != nil {
		log.Errorf("[GetExistUserByPhones] getUserIdsByPhones, phoneNumbers:%v, err:%v ", phoneNumbers, err)
		return nil, errs.UserNotExist
	}
	userIds := make([]int64, 0, len(userIdsMap))
	for _, int64s := range userIdsMap {
		userIds = append(userIds, int64s...)
	}

	// 查询 org 和 user_org 表
	userOrgPoList := make([]po.PpmOrgUserOrganization, 0)
	oriErr := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:       orgId,
		consts.TcUserId:      db.In(userIds),
		consts.TcCheckStatus: db.In([]int{consts.AppCheckStatusSuccess, consts.AppCheckStatusWait}), //审核不通过的用户也算和当前组织无关，可重新邀请
		consts.TcIsDelete:    consts.AppIsNoDelete,
	}, &userOrgPoList)
	if oriErr != nil {
		log.Error(oriErr)
		return nil, errs.BuildSystemErrorInfo(errs.UserNotExist, oriErr)
	}
	existUserIds := make([]int64, 0, len(userOrgPoList))
	for _, item := range userOrgPoList {
		existUserIds = append(existUserIds, item.UserId)
	}

	dbErr := mysql.SelectAllByCond(consts.TableUser, db.Cond{
		consts.TcId:       db.In(existUserIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userPoList)
	if dbErr != nil {
		log.Errorf("[GetExistUserByPhones] err:%v", dbErr)
		return nil, errs.UserNotExist
	}
	copyErr := copyer.Copy(userPoList, &userBoList)
	if copyErr != nil {
		log.Error(copyErr)
		return userBoList, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	// 手机号已经不在user表了，所以要重新赋值一下
	userIdToPhoneMap := make(map[int64]string, len(userIdsMap))
	for phone, int64s := range userIdsMap {
		for _, id := range int64s {
			userIdToPhoneMap[id] = phone
		}
	}
	for i := range userBoList {
		userBoList[i].Mobile = userIdToPhoneMap[userBoList[i].ID]
	}

	return userBoList, nil
}

// 由于现在一个手机号已经不能确定是哪个userId了，所以先拿出所有userIds
func getUserIdsMapByPhones(phoneNumbers []string) (map[string][]int64, errs.SystemErrorInfo) {
	globalUsers, err := dao.GetGlobalUser().GetGlobalUsersByMobiles(phoneNumbers)
	if err != nil {
		log.Errorf("[GetExistUserByPhones] GetGlobalUserIdsByMobiles, phones:%v, err:%v ", phoneNumbers, err)
		return nil, errs.UserNotExist
	}
	if len(globalUsers) <= 0 {
		return map[string][]int64{}, nil
	}
	globalUserIds := make([]int64, 0, len(globalUsers))
	phonesMap := make(map[int64]string, len(globalUsers))
	for _, user := range globalUsers {
		globalUserIds = append(globalUserIds, user.Id)
		phonesMap[user.Id] = user.Mobile
	}

	relations, err := dao.GetGlobalUserRelation().GetRelationsByGlobalUserIds(globalUserIds)
	if err != nil {
		log.Errorf("[GetExistUserByPhones] GetUserIdsByGlobalUserIds, globalUserIds:%v, err:%v ", globalUserIds, err)
		return nil, errs.UserNotExist
	}
	userIdsMap := make(map[string][]int64, len(relations))
	for _, relation := range relations {
		phone := phonesMap[relation.GlobalUserId]
		userIdsMap[phone] = append(userIdsMap[phone], relation.UserId)
	}

	return userIdsMap, nil
}

func GetUserInfoByEmail(email string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	userPo := &po.PpmOrgUser{}
	err := mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcEmail:    email,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, userPo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.UserNotExist
		}
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

func GetUserInfoById(userId int64) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	userPo := &po.PpmOrgUser{}
	err := mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcId:       userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, userPo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.UserNotExist
		}
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

//如果err不等于空，说明用户未注册
func GetUserInfoByLoginNameAndPwd(loginName string, pwd string) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	userPo := &po.PpmOrgUser{}
	err := mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcLoginName: loginName,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}, userPo)
	if err != nil {
		log.Error(err)
		return nil, errs.UserNotExist
	}
	inputPwd := pwd
	salt := userPo.PasswordSalt
	pwd = util.PwdEncrypt(pwd, salt)

	if userPo.Password != pwd {
		// 尝试使用无码的校验方案
		if !CheckIsLessCodeAccount(userPo.Password, userPo.PasswordSalt, loginName, inputPwd) {
			return nil, errs.PwdLoginUsrOrPwdNotMatch
		}
	}
	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

// 校验是否是无码系统的用户体系。因为无码系统的校验方式和极星的不一致，因此需要多一次校验尝试。
func CheckIsLessCodeAccount(userPwd, salt, loginName, inputPwd string) bool {
	encodedPwd := util.PwdEncryptForLesscodeAccoutLogin(loginName+inputPwd, salt)
	return encodedPwd == userPwd
}

//loginName为允许为账号，邮箱，手机号
func CheckLoginNameIsExist(addressType int, loginName string) errs.SystemErrorInfo {
	var err error
	switch addressType {
	case consts.ContactAddressTypeEmail:
		userPo := &po.PpmOrgUser{}
		err = mysql.SelectOneByCond(consts.TableUser, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcEmail:    loginName,
		}, &userPo)
	case consts.ContactAddressTypeMobile:
		_, err = dao.GetGlobalUser().GetGlobalUserByMobile(loginName)
	}

	if err != nil {
		if err == db.ErrNoMoreRows {
			return errs.UserNotExist
		} else {
			log.Error(err)
			return errs.MysqlOperateError
		}
	}

	return nil
}

func UserRegister(regInfo bo.UserSMSRegisterInfo, tx ...sqlbuilder.Tx) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	loginName := regInfo.PhoneNumber
	if regInfo.AccountName != "" {
		loginName = regInfo.AccountName
	}
	if loginName == "" {
		loginName = regInfo.Email
	}
	if loginName == "" {
		return nil, errs.UserRegisterError
	}

	//注册时对手机号加锁
	lockKey := consts.UserRegisterNameLock + loginName
	uid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uid)
	if lockErr != nil {
		log.Error(lockErr)
		return nil, errs.UserRegisterError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		log.Error("注册失败")
		return nil, errs.UserRegisterError
	}

	userPo, err := assemblyUserRegisterUserInfo(regInfo)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	userBo := &bo.UserInfoBo{}
	var mysqlErr error
	if len(tx) > 0 {
		userBo.GlobalUserId, mysqlErr = registerUserToDb(userPo, regInfo, tx[0])
	} else {
		mysqlErr = mysql.TransX(func(tx sqlbuilder.Tx) error {
			userBo.GlobalUserId, mysqlErr = registerUserToDb(userPo, regInfo, tx)
			return mysqlErr
		})
	}
	if mysqlErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.UserRegisterError, mysqlErr)
	}

	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

func registerUserToDb(userPo *po.PpmOrgUser, regInfo bo.UserSMSRegisterInfo, tx sqlbuilder.Tx) (int64, error) {
	//插入用户
	err1 := mysql.TransInsert(tx, userPo)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.UserRegisterError)
	}

	//即使注册时插入失败，查看时也会做二次check并插入
	err := insertUserConfig(userPo.OrgId, userPo.Id, tx)
	if err != nil {
		log.Error(err)
	}

	// 创建全局user，并要关联
	if len(regInfo.PhoneNumber) > 0 {
		globalUserId, err := createGlobalUserAndRelation(regInfo.PhoneNumber, regInfo.Password, userPo.Id, tx)
		if err != nil {
			log.Errorf("[UserRegister] createGlobalUserAndRelation:%v", err)
			return 0, err
		}
		return globalUserId, nil
	}

	return 0, nil
}

// 创建global用户，关联user
func createGlobalUserAndRelation(mobile, password string, userId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	globalUserId, err := idfacade.ApplyPrimaryIdRelaxed(po.TableNamePpmOrgGlobalUser)
	if err != nil {
		log.Errorf("[createGlobalUserAndRelation] ApplyPrimaryIdRelaxed:%v", err)
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}
	relationId, err := idfacade.ApplyPrimaryIdRelaxed(po.TableNamePpmOrgGlobalUserRelation)
	if err != nil {
		log.Errorf("[createGlobalUserAndRelation] ApplyPrimaryIdRelaxed:%v", err)
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}
	salt := ""
	pwd := ""
	if password != "" {
		salt = md5.Md5V(uuid.NewUuid())
		pwd = util.PwdEncrypt(password, salt)
	}

	err2 := dao.GetGlobalUser().Create(&po.PpmOrgGlobalUser{
		Id:              globalUserId,
		Mobile:          mobile,
		LastLoginUserId: userId,
		Password:        pwd,
		PasswordSalt:    salt,
	}, tx)
	if err2 != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	err2 = dao.GetGlobalUserRelation().CreateRelations([]*po.PpmOrgGlobalUserRelation{
		{
			Id:           relationId,
			GlobalUserId: globalUserId,
			UserId:       userId,
		},
	}, tx)
	if err2 != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	return globalUserId, nil
}

func assemblyUserRegisterUserInfo(regInfo bo.UserSMSRegisterInfo) (*po.PpmOrgUser, errs.SystemErrorInfo) {
	loginName := regInfo.PhoneNumber
	if regInfo.AccountName != "" {
		loginName = regInfo.AccountName
	}
	email := regInfo.Email
	sourceChannel := regInfo.SourceChannel
	sourcePlatform := regInfo.SourcePlatform
	name := regInfo.Name

	userId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUser)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}
	pwd := ""
	salt := ""
	if regInfo.Password != "" {
		salt = md5.Md5V(uuid.NewUuid())
		pwd = util.PwdEncrypt(regInfo.Password, salt)

	}
	userPo := &po.PpmOrgUser{
		Id:                 userId,
		OrgId:              regInfo.OrgId,
		Name:               name,
		NamePinyin:         pinyin.ConvertToPinyin(name),
		Avatar:             "",
		LoginName:          loginName, //
		LoginNameEditCount: 0,
		Email:              email,
		//Mobile:             phoneNumber,  // 不再这里维护手机号了，要在globalUser表中维护
		SourceChannel:  sourceChannel,
		SourcePlatform: sourcePlatform,
		//SourceObjId:,
		Creator:      userId,
		Updator:      userId,
		Password:     pwd,
		PasswordSalt: salt,
	}
	if regInfo.MobileRegion != "" {
		userPo.MobileRegion = regInfo.MobileRegion
	}

	return userPo, nil
}

//增加组织成员，添加关联以及加入顶级部门
//inCheck：是否需要被审核
func AddOrgMember(orgId, userId int64, operatorId int64, inCheck bool, inDisabled bool) errs.SystemErrorInfo {
	userOrgRelation, err := GetUserOrganizationNewestRelation(orgId, userId)
	//关联不存在或者已删除，或者审核未通过，此时允许新增关联
	if (err != nil && err.Code() == errs.UserOrgNotRelation.Code()) || userOrgRelation.IsDelete == consts.AppIsDeleted || userOrgRelation.CheckStatus == consts.AppCheckStatusFail {
		log.Infof("用户%d和组织%d需要做关联", userId, orgId)
		//上锁
		lockKey := fmt.Sprintf("%s%d:%d", consts.UserAndOrgRelationLockKey, orgId, userId)
		lockUuid := uuid.NewUuid()

		suc, lockErr := cache.TryGetDistributedLock(lockKey, lockUuid)
		if lockErr != nil {
			log.Error(lockErr)
			return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
		}
		if suc {
			defer func() {
				if _, err := cache.ReleaseDistributedLock(lockKey, lockUuid); err != nil {
					log.Error(err)
				}
			}()
			//二次check
			userOrgRelation, err := GetUserOrganizationNewestRelation(orgId, userId)
			if (err != nil && err.Code() == errs.UserOrgNotRelation.Code()) || userOrgRelation.IsDelete == consts.AppIsDeleted || userOrgRelation.CheckStatus == consts.AppCheckStatusFail {
				//组织用户做关联
				log.Infof("用户%d和组织%d开始关联", userId, orgId)
				err = AddUserOrgRelation(orgId, userId, false, inCheck, inDisabled)
				//判断关联是否失败
				if err != nil {
					log.Error(err)
					return err
				}
			}
		}
	}
	clearErr := ClearBaseUserInfo(orgId, userId)
	if clearErr != nil {
		log.Error(clearErr)
		return clearErr
	}
	return nil
}

//添加用户和组织关联
func AddUserOrgRelation(orgId, userId int64, inUsed bool, inCheck bool, inDisabled bool) errs.SystemErrorInfo {
	userOrgId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOrganization)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}

	useStatus := consts.AppStatusDisabled
	if inUsed {
		useStatus = consts.AppStatusEnable
	}

	checkStatus := consts.AppCheckStatusSuccess
	status := consts.AppStatusEnable
	if inCheck {
		checkStatus = consts.AppCheckStatusWait
		status = consts.AppStatusDisabled
	}

	if inDisabled {
		status = consts.AppStatusDisabled
	}

	userOrgPo := po.PpmOrgUserOrganization{
		Id:          userOrgId,
		OrgId:       orgId,
		UserId:      userId,
		CheckStatus: checkStatus,
		UseStatus:   useStatus,
		Status:      status,
		Creator:     userId,
		Updator:     userId,
	}

	if status == consts.AppStatusDisabled {
		userOrgPo.StatusChangeTime = time.Now()
	}

	err1 := mysql.Insert(&userOrgPo)
	if err != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func UpdateUserInfo(userId int64, upd mysql.Upd) errs.SystemErrorInfo {
	err := mysql.UpdateSmart(consts.TableUser, userId, upd)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

// 根据 outUserId 更新一个用户的 out 信息
func UpdateUserOutInfo(orgId int64, outUserId string, upd mysql.Upd) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableUserOutInfo, db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcOutUserId: outUserId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}, upd)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func UpdateUserDefaultOrg(userId, orgId int64) (int64, errs.SystemErrorInfo) {
	orgUserId, err := UpdateGlobalUserLastLoginInfoWithCheckGlobal(userId, orgId)
	if err != nil {
		return userId, err
	}
	updateUserInfoErr := UpdateUserInfo(orgUserId, mysql.Upd{
		consts.TcOrgId: orgId,
	})
	if updateUserInfoErr != nil {
		log.Error(updateUserInfoErr)
		return userId, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	return orgUserId, nil
}

//用户登录时回调
func UserLoginHook(userId int64, orgId int64, sourceChannel string) errs.SystemErrorInfo {
	//更新用户最后登录时间
	err := UpdateUserInfo(userId, mysql.Upd{
		consts.TcLastLoginTime: date.FormatTime(types.NowTime()),
	})
	if err != nil {
		log.Error(err)
	}
	if orgId != 0 {
		//更新使用状态
		_, err1 := mysql.UpdateSmartWithCond(consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcUseStatus: consts.AppStatusEnable,
		})
		if err1 != nil {
			log.Error(err1)
			err = errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}
	addLoginRecordErr := AddLoginRecord(orgId, userId, sourceChannel)
	if addLoginRecordErr != nil {
		log.Error(addLoginRecordErr)
		err = addLoginRecordErr
	}
	return err
}

func BatchGetUserDetailInfo(userIds []int64) ([]*bo.UserInfoBo, errs.SystemErrorInfo) {
	userIds = slice.SliceUniqueInt64(userIds)
	var po []*po.PpmOrgUser
	err := mysql.SelectAllByCond(consts.TableUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(userIds),
	}, &po)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	var bo []*bo.UserInfoBo
	copyErr := copyer.Copy(po, &bo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bo, nil
}

// BatchGetUserDetailInfoWithMobile 批量获取用户信息，并查询手机号码
func BatchGetUserDetailInfoWithMobile(userIds []int64) ([]*bo.UserInfoBo, errs.SystemErrorInfo) {
	userBos, err := BatchGetUserDetailInfo(userIds)
	if err != nil {
		return nil, err
	}
	userIdsMap, err := GetUserIdsMobileMap(userIds)
	if err != nil {
		return nil, err
	}

	for i := range userBos {
		userBos[i].Mobile = userIdsMap[userBos[i].ID]
	}

	return userBos, nil
}

func CreateUser(orgId, userId int64, regInfo orgvo.CreateUserReq) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	loginName := regInfo.PhoneNumber

	regInfo.DepartmentIds = slice.SliceUniqueInt64(regInfo.DepartmentIds)
	if len(regInfo.DepartmentIds) > 0 {
		//查询部门是否有效
		departmentPos := &[]po.PpmOrgDepartment{}
		departmentErr := mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcId:       db.In(regInfo.DepartmentIds),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, departmentPos)
		if departmentErr != nil {
			log.Error(departmentErr)
			return nil, errs.MysqlOperateError
		}

		if len(*departmentPos) != len(regInfo.DepartmentIds) {
			return nil, errs.DepartmentNotExist
		}
	}

	//注册时对手机号加锁
	lockKey := consts.UserBindLoginNameLock + loginName
	uid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uid)
	if lockErr != nil {
		log.Error(lockErr)
		return nil, errs.UserRegisterError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		log.Error("注册失败")
		return nil, errs.UserRegisterError
	}

	//检测用户名是否存在
	err := CheckPhoneAndEmail(regInfo.PhoneNumber, regInfo.Email)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//插入用户
	salt := md5.Md5V(uuid.NewUuid())
	//默认密码取手机号后六位
	targetPassword := str.Substr(regInfo.PhoneNumber, -6, 6)
	//用户id
	userPoId, userPoIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUser)
	if userPoIdErr != nil {
		log.Error(userPoIdErr)
		return nil, userPoIdErr
	}
	userPo := &po.PpmOrgUser{
		Id:                 userPoId,
		OrgId:              0,
		Name:               regInfo.Name,
		NamePinyin:         pinyin.ConvertToPinyin(regInfo.Name),
		Avatar:             "",
		LoginName:          regInfo.PhoneNumber, //
		LoginNameEditCount: 0,
		Email:              regInfo.Email,
		Mobile:             regInfo.PhoneNumber,
		Creator:            userId,
		Updator:            userId,
		PasswordSalt:       salt,
		Password:           util.PwdEncrypt(targetPassword, salt),
	}
	//用户配置id
	userConfigId, userConfigIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserConfig)
	if userConfigIdErr != nil {
		log.Error(userConfigIdErr)
		return nil, userConfigIdErr
	}
	//用户组织关联id
	userOrgId, userOrgIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserOrganization)
	if userOrgIdErr != nil {
		log.Error(userOrgIdErr)
		return nil, userOrgIdErr
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//用户表
		err1 := mysql.TransInsert(tx, userPo)
		if err1 != nil {
			log.Error(err1)
			return err1
		}

		//用户配置表
		userConfig := &po.PpmOrgUserConfig{}
		userConfig.Id = userConfigId
		userConfig.OrgId = orgId
		userConfig.UserId = userId
		userConfig.Creator = userId
		userConfig.Updator = userId
		err2 := mysql.TransInsert(tx, userConfig)
		if err2 != nil {
			log.Error(err2)
			return err2
		}

		//用户组织表
		err3 := mysql.TransInsert(tx, &po.PpmOrgUserOrganization{
			Id:          userOrgId,
			OrgId:       orgId,
			UserId:      userPo.Id,
			CheckStatus: consts.AppCheckStatusSuccess,
			UseStatus:   consts.AppStatusDisabled,
			Status:      consts.AppStatusEnable,
			Creator:     userId,
		})
		if err3 != nil {
			log.Error(err3)
			return err3
		}

		//用户部门表
		if len(regInfo.DepartmentIds) > 0 {
			userDeptIds, userDeptIdsErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(regInfo.DepartmentIds))
			if userDeptIdsErr != nil {
				log.Error(userDeptIdsErr)
				return userDeptIdsErr
			}
			var userDepartment []interface{}
			for i, id := range regInfo.DepartmentIds {
				userDepartment = append(userDepartment, po.PpmOrgUserDepartment{
					Id:           userDeptIds.Ids[i].Id,
					OrgId:        orgId,
					UserId:       userPo.Id,
					DepartmentId: id,
					Creator:      userId,
				})
			}

			err4 := mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, userDepartment)
			if err4 != nil {
				log.Error(err4)
				return err4
			}
		}
		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	userBo := &bo.UserInfoBo{}
	copyErr := copyer.Copy(userPo, userBo)
	if copyErr != nil {
		log.Error(copyErr)
	}
	return userBo, nil
}

func CheckPhoneAndEmail(phoneNumber, email string) errs.SystemErrorInfo {
	_, err := dao.GetGlobalUser().GetGlobalUserByMobile(phoneNumber)
	if err != nil {
		if err != db.ErrNoMoreRows {
			log.Errorf("[CheckPhoneAndEmail] phoneNumber:%v, err:%v", phoneNumber, err)
			return errs.MysqlOperateError
		}
	} else {
		return errs.MobileAlreadyBindOtherAccountError
	}

	count1, err1 := mysql.SelectCountByCond(consts.TableUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcEmail:    email,
	})
	if err1 != nil {
		log.Error(err1)
		return errs.MysqlOperateError
	}
	if count1 > 0 {
		return errs.EmailAlreadyBindByOtherAccountError
	}

	return nil
}

// "同步成员、部门信息"任务推送到消息队列
func PushSyncMemberDept(paramReqVo orgvo.SyncUserInfoFromFeiShuReqVo) errs.SystemErrorInfo {
	paramReqVoWrapper := mqbo.PushSyncMemberDeptBo{
		SyncUserInfoFromFeiShuReqVo: paramReqVo,
	}
	message, err := json.ToJson(paramReqVoWrapper)
	if err != nil {
		log.Error(err)
	}
	//这里key使用项目id，保证同一项目下导入的任务顺序的有效性
	mqMessage := &model.MqMessage{
		Topic:          config.GetMqSyncMemberDeptTopicConfig().Topic,
		Keys:           strconv.FormatInt(paramReqVo.OrgId, 10),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, paramReqVo.OrgId)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
	log.Infof("mq消息推送：同步成员、部门信息，message has push. ")
	return msgErr
}

func GetUserOutInfoByUserIdAndOrgId(userId int64, orgId int64, sourceChannel string) (*bo.UserOutInfoBo, errs.SystemErrorInfo) {
	info := &po.PpmOrgUserOutInfo{}
	err := mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         db.In([]int64{orgId, 0}),
		consts.TcSourceChannel: sourceChannel,
		consts.TcUserId:        userId,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.NotDisbandCurrentSourceChannel
		} else {
			return nil, errs.MysqlOperateError
		}
	}

	infoBo := &bo.UserOutInfoBo{}
	copyer.Copy(info, infoBo)
	return infoBo, nil
}

func GetUserOutInfosByUserIds(userIds []int64, sourceChannel string) (*bo.UserOutInfoBo, error) {
	info := &po.PpmOrgUserOutInfo{}
	err := mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcSourceChannel: sourceChannel,
		consts.TcUserId:        db.In(userIds),
	}, info)
	if err != nil {
		return nil, err
	}

	infoBo := &bo.UserOutInfoBo{}
	copyer.Copy(info, infoBo)
	return infoBo, nil
}

func GetUserOutInfoByUserIdByOrgIds(userId int64, orgIds []int64, sourceChannel string) (*bo.UserOutInfoBo, errs.SystemErrorInfo) {
	info := &po.PpmOrgUserOutInfo{}
	err := mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcOrgId:         db.In(orgIds),
		consts.TcSourceChannel: sourceChannel,
		consts.TcUserId:        userId,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.NotDisbandCurrentSourceChannel
		} else {
			return nil, errs.MysqlOperateError
		}
	}

	infoBo := &bo.UserOutInfoBo{}
	copyer.Copy(info, infoBo)
	return infoBo, nil
}

// GenNameByPhone 通过手机号生成默认的用户姓名
func GenDefaultNameByPhone(phone string) string {
	suffix := str.Substr(phone, -4, 4)
	return fmt.Sprintf("%s%s", "用户", suffix)
}

func SetVersionConfig(orgId, userId int64, versionInfoVisible bool) errs.SystemErrorInfo {
	// 模板预览的组织不展示版本信息弹窗
	if orgId == consts.PreviewTplOrgId {
		return nil
	}
	userConfigPo := po.PpmOrgUserConfig{}
	err := mysql.SelectOneByCond(userConfigPo.TableName(), db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userConfigPo)
	if err != nil {
		log.Errorf("[SetVersionConfig]查询错误 err:%v, orgId:%v, userId:%v", err, orgId, userId)
		return errs.MysqlOperateError
	}
	userConfigExt := bo.UserConfigExt{}
	if userConfigPo.Ext != "" {
		errJson := json.FromJson(userConfigPo.Ext, &userConfigExt)
		if errJson != nil {
			return errs.BuildSystemErrorInfo(errs.JSONConvertError, errJson)
		}
	}
	versionConfig := bo.VersionConfig{}
	versionConfig.VersionInfoVisible = versionInfoVisible
	userConfigExt.Version = &versionConfig

	// 更新ext
	newExtJson := json.ToJsonIgnoreError(userConfigExt)
	updateErr := updateUserConfig(orgId, userId, newExtJson)
	if updateErr != nil {
		log.Error(updateErr)
		return updateErr
	}

	// 清除user config 缓存
	cacheErr := clearCacheUserConfig(orgId, userId, versionInfoVisible)
	if cacheErr != nil {
		log.Error(cacheErr)
		return cacheErr
	}

	return nil
}

func GetCardConfig(orgId, userId int64) (*orgvo.VersionResp, errs.SystemErrorInfo) {
	key, err := util.ParseCacheKey(sconsts.CacheCardConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err != nil {
		return nil, err
	}
	cardConfigJson, err2 := cache.Get(key)
	if err2 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err2)
	}
	cardInfo := orgvo.VersionResp{}
	if cardConfigJson != "" {
		err3 := json.FromJson(cardConfigJson, &cardInfo)
		if err3 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		userConfigPo := po.PpmOrgUserConfig{}
		err4 := mysql.SelectOneByCond(userConfigPo.TableName(), db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, &userConfigPo)
		if err4 != nil {
			log.Errorf("[GetCardConfig]查询失败 err:%v, orgId:%v, userId:%v", err4, orgId, userId)
			return nil, errs.MysqlOperateError
		}
		userConfigExt := bo.UserConfigExt{}
		jsonErr := json.FromJson(userConfigPo.Ext, &userConfigExt)
		if jsonErr != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		if userConfigExt.Version == nil {
			userConfigExt.Version = &bo.VersionConfig{VersionInfoVisible: false}
		}

		cardData := orgvo.GetVersionData{
			VersionInfoVisible: userConfigExt.Version.VersionInfoVisible,
		}

		cardInfo.Version = cardData

		cardConfigStr := json.ToJsonIgnoreError(userConfigExt)
		cacheSetErr := cache.SetEx(key, cardConfigStr, consts.CacheBaseExpire)
		if cacheSetErr != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, cacheSetErr)
		}
	}

	return &cardInfo, nil
}

func updateUserConfig(orgId, userId int64, ext string) errs.SystemErrorInfo {
	upds := mysql.Upd{
		consts.TcExt: ext,
	}
	_, err := mysql.UpdateSmartWithCond(consts.TableUserConfig, db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: userId,
	}, upds)
	if err != nil {
		log.Errorf("[updateVersionConfig]更新错误 err:%v, orgId:%v, userId:%v", err, orgId, userId)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return nil
}

func clearCacheUserConfig(orgId, userId int64, versionInfoVisible bool) errs.SystemErrorInfo {
	cacheKey, err := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	// 获取用户缓存，将ext中的version信息删除
	configExt, err2 := cache.Get(cacheKey)
	if err2 != nil {
		log.Error(err2)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err2)
	}
	// 新用户会弹出新手知道和版本信息，
	// 如果先点击版本信息，那么这时候就没有缓存信息，直接返回
	if configExt == "" {
		return nil
	}
	configBo := bo.UserConfigBo{}
	errJson := json.FromJson(configExt, &configBo)
	if errJson != nil {
		log.Error(errJson)
		return errs.BuildSystemErrorInfo(errs.JSONConvertError, errJson)
	}
	ext := bo.UserConfigExt{}
	// 缓存中没有ext的信息，说明版本信息弹窗和新手指导都需要展示
	if configBo.Ext == "" {
		// VersionInfoVisible: true 表示需要展示, 和新手指引的flag意思相反
		ext.Version = &bo.VersionConfig{VersionInfoVisible: true}
		configBo.Ext = json.ToJsonIgnoreError(ext.Version)
	}

	errJson2 := json.FromJson(configBo.Ext, &ext)
	if errJson2 != nil {
		log.Error(errJson2)
		return errs.BuildSystemErrorInfo(errs.JSONConvertError, errJson2)
	}
	versionConfig := bo.VersionConfig{}
	versionConfig.VersionInfoVisible = versionInfoVisible
	ext.Version = &versionConfig

	configBo.Ext = json.ToJsonIgnoreError(ext)
	configJson := json.ToJsonIgnoreError(configBo)
	cacheSet := cache.Set(cacheKey, configJson)
	if cacheSet != nil {
		log.Error(cacheSet)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, cacheSet)
	}
	return nil
}

// CheckIsAdmin 检查是否是组织的管理员（包含组织拥有者，超管，子管理员）
func CheckIsAdmin(orgId, curUserId int64) bool {
	if curUserId == 0 {
		return false
	}
	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, curUserId)
	if manageAuthInfoResp.Failure() {
		log.Errorf("[CheckIsAdmin] orgId:%v, curUserId:%v, error:%v", orgId, curUserId, manageAuthInfoResp.Error())
		return false
	}

	return CheckIsAdminByAuthRespData(manageAuthInfoResp.Data)
}

// CheckIsAdminByAuthRespData 通过接口返回的数据，判断是否是组织的管理员
func CheckIsAdminByAuthRespData(authRespData uservo.GetUserAuthorityData) bool {
	isSysAdmin := authRespData.IsSysAdmin
	isOrgOwner := authRespData.IsOrgOwner
	isSubAdmin := authRespData.IsSubAdmin
	manageApps := authRespData.ManageApps
	isCanManageAllApps := isSysAdmin || isOrgOwner || (isSubAdmin && len(manageApps) > 0 && manageApps[0] == -1)
	return isCanManageAllApps
}

func SetUserViewLocation(req orgvo.SaveViewLocationReqVo) errs.SystemErrorInfo {
	userConfigPo := po.PpmOrgUserConfig{}
	err := mysql.SelectOneByCond(userConfigPo.TableName(), db.Cond{
		consts.TcOrgId:    req.OrgId,
		consts.TcUserId:   req.UserId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userConfigPo)
	if err != nil {
		log.Errorf("[SetUserViewLocation]查询错误userConfigPo err:%v, orgId:%v, userId:%v", err, req.OrgId, req.UserId)
		return errs.MysqlOperateError
	}
	locationOld := []bo.UserViewLocation{}
	locationNew := []bo.UserViewLocation{}
	locationConfig := userConfigPo.UserViewLocationConfig
	if locationConfig != nil {
		err = json.FromJson(*locationConfig, &locationOld)
		if err != nil {
			log.Errorf("[SetUserViewLocation] json err:%v", err)
			return errs.JSONConvertError
		}

		for _, v := range locationOld {
			if v.AppId == req.Input.AppId {
				continue
			}
			locationNew = append(locationNew, v)
		}
	}

	locationNew = append(locationNew, bo.UserViewLocation{
		AppId:       req.Input.AppId,
		IterationId: req.Input.IterationId,
		TableId:     req.Input.TableId,
		ProjectId:   req.Input.ProjectId,
		ViewId:      req.Input.ViewId,
		MenuId:      req.Input.MenuId,
		DashboardId: req.Input.DashboardId,
	})

	jsonStr := json.ToJsonIgnoreError(locationNew)
	upds := mysql.Upd{
		consts.TcUserViewLocation: jsonStr,
	}
	_, err = mysql.UpdateSmartWithCond(userConfigPo.TableName(), db.Cond{
		consts.TcOrgId:  req.OrgId,
		consts.TcUserId: req.UserId,
	}, upds)
	if err != nil {
		log.Errorf("[SetUserViewLocation] update err:%v, config:%v", err, jsonStr)
		return errs.MysqlOperateError
	}
	// 清除缓存 清除用户的上一次浏览的位置缓存
	errCache := clearUserViewLocationConfig(req.OrgId, req.UserId)
	if errCache != nil {
		log.Errorf("[SetUserViewLocation] clear cache err:%v", errCache)
		return errs.RedisOperateError
	}
	return nil
}

func clearUserViewLocationConfig(orgId int64, userId int64) errs.SystemErrorInfo {
	key, err := util.ParseCacheKey(sconsts.CacheUserLocationConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err != nil {
		return err
	}
	_, err2 := cache.Del(key)
	if err2 != nil {
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err2)
	}
	return nil
}

func GetUserViewLocation(orgId int64, userId int64) ([]*bo.UserViewLocation, errs.SystemErrorInfo) {
	key, err := util.ParseCacheKey(sconsts.CacheUserLocationConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err != nil {
		return nil, err
	}

	locationJson, err2 := cache.Get(key)
	if err2 != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err2)
	}

	locationConfigBo := []*bo.UserViewLocation{}
	if locationJson != "" {
		err3 := json.FromJson(locationJson, &locationConfigBo)
		if err3 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
	} else {
		userConfigPo := po.PpmOrgUserConfig{}
		err4 := mysql.SelectOneByCond(userConfigPo.TableName(), db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcUserId:   userId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, &userConfigPo)

		if err4 != nil {
			if err4 != db.ErrNoMoreRows {
				log.Errorf("[GetUserViewLocation] 查询userConfig err:%v, orgId:%d, userId:%d", err4, orgId, userId)
				return nil, errs.MysqlOperateError
			}
		}

		if userConfigPo.UserViewLocationConfig != nil {
			err5 := json.FromJson(*userConfigPo.UserViewLocationConfig, &locationConfigBo)
			if err5 != nil {
				return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
			}
		}

		err5 := cache.SetEx(key, json.ToJsonIgnoreError(locationConfigBo), consts.CacheBaseExpire)
		if err5 != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err5)
		}
	}

	return locationConfigBo, nil
}

func GetNewUserGuideInfo(orgId int64, userConfigBo bo.UserConfigBo,
	userExtraDataMap map[string]interface{}) errs.SystemErrorInfo {
	if orgId == consts.PreviewTplOrgId {
		// 预览模板组织视为已经展示过，不再展示“新手指引” 和 版本信息弹窗
		userExtraDataMap["visitedNewUserGuide"] = true
		versionConfig := bo.VersionConfig{VersionInfoVisible: false}
		userExtraDataMap["version"] = versionConfig
	} else {
		if userConfigBo.Ext == "" {
			userExtraDataMap["visitedNewUserGuide"] = false

			// 有新手指引就不展示版本信息弹窗了, false代表不展示版本信息
			versionConfig := bo.VersionConfig{VersionInfoVisible: false}
			userExtraDataMap["version"] = versionConfig
		} else {
			configExt := bo.UserConfigExt{}
			if err := json.FromJson(userConfigBo.Ext, &configExt); err != nil {
				log.Error(err)
				return errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
			}
			userExtraDataMap["visitedNewUserGuide"] = configExt.VisitedNewUserGuide

			userExtraDataMap["version"] = configExt.Version
			if configExt.Version == nil {
				vConfig := bo.VersionConfig{VersionInfoVisible: true}
				configExt.Version = &vConfig
				userExtraDataMap["version"] = vConfig
			}
		}
	}

	return nil
}

func GetActivity20221111Info(orgId int64, userConfigBo bo.UserConfigBo, payActivity11Flag int, userExtraDataMap map[string]interface{}) errs.SystemErrorInfo {
	// 查询活动总开关缓存，如果没有key 表明活动结束，直接返回 2（表明不需要展示弹窗），并且需要更新一下表数据 活动key为2
	//                  如果有key，查询相应的缓存key，为空 表示需要展示弹窗，不为空，直接返回相应的提示

	if orgId == consts.PreviewTplOrgId || CheckIsPrivateDeploy() {
		// 预览模板组织，私有化部署 不需要展示付费引导弹窗
		userExtraDataMap[consts.ActivityFlag] = consts.ActivityFinished
		return nil
	}

	// 查询活动总开关缓存
	switchCache, errCache := cache.Get(sconsts.CacheActivity20221111Switch)
	if errCache != nil {
		log.Error(errCache)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, errCache)
	}

	if switchCache == "" || payActivity11Flag == consts.ActivityNotFinished {
		// 活动结束 或者 已经付费了
		userExtraDataMap[consts.ActivityFlag] = consts.ActivityFinished

	} else {
		// 活动中
		if userConfigBo.Ext != "" {
			userConfigExt := bo.UserConfigExt{}
			errJson := json.FromJson(userConfigBo.Ext, &userConfigExt)
			if errJson != nil {
				return errs.BuildSystemErrorInfo(errs.JSONConvertError, errJson)
			}
			if userConfigExt.Activity20221111 == 0 {
				userConfigExt.Activity20221111 = consts.ActivityNotFinished
			}
			userExtraDataMap[consts.ActivityFlag] = userConfigExt.Activity20221111
		} else {
			// 表中没有活动标识，又在活动中，视为需要展示弹窗
			userExtraDataMap[consts.ActivityFlag] = consts.ActivityNotFinished
		}
	}

	return nil
}

func GetRemindPopUp(userConfigBo *bo.UserConfigBo, remindBindPhone int) int {
	userConfigExt := &bo.UserConfigExt{}
	json.FromJson(userConfigBo.Ext, userConfigExt)

	if userConfigExt.RemindPopUp != 0 {
		return userConfigExt.RemindPopUp
	}
	if remindBindPhone == consts.AppIsNotRemind {
		// 绑定了手机号
		return consts.NotNeedRemindPopUp
	} else {
		return consts.NeedRemindPopUp
	}
}

func updateCacheUserConfig(orgId, userId int64, userConfig string) errs.SystemErrorInfo {
	key, err := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err != nil {
		log.Error(err)
		return err
	}
	errCache := cache.Set(key, userConfig)
	if errCache != nil {
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, errCache)
	}
	return nil
}

func SetActivity11Info(orgId, userId int64, flag int) errs.SystemErrorInfo {
	if orgId == consts.PreviewTplOrgId || CheckIsPrivateDeploy() {
		return nil
	}

	userConfigPo := po.PpmOrgUserConfig{}
	err := mysql.SelectOneByCond(userConfigPo.TableName(), db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userConfigPo)
	if err != nil {
		log.Errorf("[SetActivity11Info]查询错误 err:%v, orgId:%v, userId:%v", err, orgId, userId)
		return errs.MysqlOperateError
	}
	userConfigExt := bo.UserConfigExt{}
	if userConfigPo.Ext != "" {
		errJson := json.FromJson(userConfigPo.Ext, &userConfigExt)
		if errJson != nil {
			return errs.BuildSystemErrorInfo(errs.JSONConvertError, errJson)
		}
	}

	userConfigExt.Activity20221111 = flag

	// 更新ext
	newExtJson := json.ToJsonIgnoreError(userConfigExt)
	updateErr := updateUserConfig(orgId, userId, newExtJson)
	if updateErr != nil {
		log.Error(updateErr)
		return updateErr
	}

	// 清除缓存
	key, err2 := util.ParseCacheKey(sconsts.CacheUserConfig, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:  orgId,
		consts.CacheKeyUserIdConstName: userId,
	})
	if err2 != nil {
		log.Error(err2)
		return err2
	}
	_, errCache := cache.Del(key)
	if errCache != nil {
		log.Error(errCache)
		return errs.RedisOperateError
	}

	return nil
}

// CheckIsFsPlatformAdminByOpenId 检查是否是 fs 的应用管理员
//func CheckIsFsPlatformAdminByOpenId(outOrgId string, openId string) (bool, errs.SystemErrorInfo) {
//	if outOrgId == "" {
//		return false, nil
//	}
//	tenant, err := feishu.GetTenant(outOrgId)
//	if err != nil {
//		log.Errorf("[PersonalInfo] err: %v", err)
//		return false, err
//	}
//	outAdmins, oriErr := tenant.AdminUserList()
//	if oriErr != nil {
//		log.Errorf("[PersonalInfo] err: %v", oriErr)
//		return false, errs.BuildSystemErrorInfo(errs.FeiShuOpenApiCallError, oriErr)
//	}
//	adminOpenIdsMap := make(map[string]struct{}, 5)
//	for _, one := range outAdmins.Data.UserList {
//		adminOpenIdsMap[one.OpenId] = struct{}{}
//	}
//	if _, ok := adminOpenIdsMap[openId]; ok {
//		return true, nil
//	}
//
//	return false, nil
//}

//func GetFsPlatformAdmin(orgId, userId int64, outOrgId, openId string) (bool, errs.SystemErrorInfo) {
//	// 先获取缓存
//	key, errParse := util.ParseCacheKey(sconsts.CacheFsPlatformAdmin, map[string]interface{}{
//		consts.CacheKeyOrgIdConstName:  orgId,
//		consts.CacheKeyUserIdConstName: userId,
//	})
//	if errParse != nil {
//		return false, errParse
//	}
//	platformAdmin, errCache := cache.Get(key)
//	if errCache != nil {
//		log.Error(errCache)
//		return false, errs.RedisOperateError
//	}
//	if platformAdmin != "" {
//		if platformAdmin == fmt.Sprintf("%d", consts.IsPlatformAdmin) {
//			return true, nil
//		} else {
//			return false, nil
//		}
//
//	} else {
//		// 重新请求飞书接口获取 是否是飞书管理员
//		isAdmin, err := CheckIsFsPlatformAdminByOpenId(outOrgId, openId)
//		if err != nil {
//			log.Errorf("[GetFsPlatformAdmin] CheckIsFsPlatformAdminByOpenId, err:%v", err)
//			return false, err
//		}
//
//		isAdminCacheStr := fmt.Sprintf("%d", consts.NotPlatformAdmin)
//		if isAdmin {
//			isAdminCacheStr = fmt.Sprintf("%d", consts.IsPlatformAdmin)
//		}
//
//		errCache = cache.SetEx(key, isAdminCacheStr, consts.CacheExpire1Day)
//		if errCache != nil {
//			log.Error(errCache)
//			return false, errs.RedisOperateError
//		}
//
//		return isAdmin, nil
//	}
//
//}

func SaveUserViewLocation(req orgvo.SaveViewLocationReqVo) errs.SystemErrorInfo {
	params := req.Input
	viewConfig := po.PpmOrgUserViewLocation{}
	conds := db.Cond{
		consts.TcOrgId:    req.OrgId,
		consts.TcUserId:   req.UserId,
		consts.TcAppId:    params.AppId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	appId, parseError := strconv.ParseInt(params.AppId, 10, 64)
	if parseError != nil {
		log.Errorf("[UpdateUserViewLocation] parse err:%v", parseError)
		return errs.TypeConvertError
	}
	err := mysql.SelectOneByCond(consts.TableOrgUserLocation, conds, &viewConfig)
	if err != nil {
		if err == db.ErrNoMoreRows {
			// 没有就插入数据
			idResp, errSys := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgUserLocation)
			if errSys != nil {
				log.Errorf("[UpdateUserViewLocation] idfacade err:%v", errSys)
				return errSys
			}
			locationConfig := po.PpmOrgUserViewLocation{
				Id:      idResp,
				OrgId:   req.OrgId,
				UserId:  req.UserId,
				AppId:   appId,
				Config:  json.ToJsonIgnoreError(params),
				Creator: req.UserId,
				Updator: req.UserId,
			}
			err = mysql.Insert(&locationConfig)
			if err != nil {
				log.Errorf("[UpdateUserViewLocation] insert err:%v", err)
				return errs.MysqlOperateError
			}
			return nil

		} else {
			log.Errorf("[UpdateUserViewLocation] select err:%v, orgId:%v, userId:%v, appId:%v", err, req.OrgId, req.UserId, params.AppId)
			return errs.MysqlOperateError
		}
	}
	// 更新
	configStr := viewConfig.Config
	inputStr := json.ToJsonIgnoreError(params)
	if configStr == inputStr {
		return nil
	}
	_, err = mysql.UpdateSmartWithCond(consts.TableOrgUserLocation, conds, mysql.Upd{
		consts.TcConfig: inputStr,
	})
	if err != nil {
		log.Errorf("[UpdateUserViewLocation] update err:%v, orgId:%v, userId:%v, appId:%v", err, req.OrgId, req.UserId, params.AppId)
		return errs.MysqlOperateError
	}
	return nil
}

func GetViewLocationList(orgId, userId int64) ([]*orgvo.UserLastViewLocationData, errs.SystemErrorInfo) {
	pos := []po.PpmOrgUserViewLocation{}
	err := mysql.SelectAllByCond(consts.TableOrgUserLocation, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &pos)
	if err != nil {
		log.Errorf("[GetViewLocationList] SelectAllByCond err:%v, orgId:%v, userId:%v", err, orgId, userId)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	ret := make([]*orgvo.UserLastViewLocationData, 0, len(pos))

	for _, v := range pos {
		locationData := orgvo.UserLastViewLocationData{}
		if v.Config != "" {
			errJson := json.FromJson(v.Config, &locationData)
			if errJson != nil {
				log.Errorf("[GetViewLocationList] json err:%v,orgId:%v, userId:%v", err, orgId, userId)
				return nil, errs.JSONConvertError
			}
			ret = append(ret, &locationData)
		}
	}
	return ret, nil
}

func DeleteUserLocationWithAppId(orgId, userId int64, appId int64) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableOrgUserLocation, db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: userId,
		consts.TcAppId:  appId,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	if err != nil {
		log.Errorf("[DeleteUserLocationWithAppId] err:%v", err)
		return errs.MysqlOperateError
	}
	return nil
}

func CheckIsActivity11() (bool, errs.SystemErrorInfo) {
	activityFlag, errCache := cache.Get(sconsts.CacheActivity20221111Switch)
	if errCache != nil {
		return false, errs.RedisOperateError
	}
	if activityFlag == "" {
		// 活动结束
		return false, nil
	}

	return true, nil
}

func GetOrgSuperAdminInfo(orgId int64) ([]*orgvo.GetOrgSuperAdminInfoData, errs.SystemErrorInfo) {
	res := []*orgvo.GetOrgSuperAdminInfoData{}
	sysManageGroup, err := GetSysManageGroup(orgId)
	if err != nil {
		log.Errorf("[GetOrgSuperAdminInfo] GetSysManageGroup err:%v, orgId:%v", err, orgId)
		return nil, errs.MysqlOperateError
	}
	adminUserIds := []int64{}
	if sysManageGroup.UserIds != nil {
		errJson := json.FromJson(*sysManageGroup.UserIds, &adminUserIds)
		if errJson != nil {
			return nil, errs.JSONConvertError
		}
	}
	userInfoBatch, errSys := GetBaseUserInfoBatch(orgId, adminUserIds)
	if errSys != nil {
		log.Errorf("[GetOrgSuperAdminInfo] GetBaseUserInfoBatch err:%v, orgId:%v", errSys, orgId)
		return nil, errSys
	}
	for _, user := range userInfoBatch {
		res = append(res, &orgvo.GetOrgSuperAdminInfoData{
			UserId: user.UserId,
			OpenId: user.OutUserId,
		})
	}
	return res, nil
}

func UpdateUserToSysManageGroup(orgId int64, updateUserIds []int64, updateType int) errs.SystemErrorInfo {
	group, err := GetSysManageGroup(orgId)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	var userIds []int64
	jsonErr := json.FromJson(*group.UserIds, &userIds)
	if jsonErr != nil {
		log.Error(jsonErr)
		return errs.JSONConvertError
	}

	newUserIds := []int64{}
	if updateType == consts.AddType {
		newUserIds = append(userIds, updateUserIds...)
	} else if updateType == consts.DelType {
		newUserIds = businees.DifferenceInt64Set(userIds, updateUserIds)
	}

	newUserIds = slice.SliceUniqueInt64(newUserIds)
	valueStr := json.ToJsonIgnoreError(newUserIds)
	_, err = mysql.UpdateSmartWithCond(consts.TableLcPerManageGroup, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcLangCode: consts.ManageGroupSys,
	}, mysql.Upd{
		consts.TcUserIds: valueStr,
		consts.TcVersion: db.Raw("version + 1"),
	})
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	return nil
}
