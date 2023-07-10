package domain

import (
	msgPb "github.com/star-table/interface/golang/msg/v1"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade/common/report"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/spf13/cast"

	"upper.io/db.v3/lib/sqlbuilder"

	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

func SetUserPassword(globalUserId int64, password string, salt string, operatorId int64) errs.SystemErrorInfo {
	log.Infof("密码: %s, Slat: %s", password, salt)
	_, err := mysql.UpdateSmartWithCond(po.TableNamePpmOrgGlobalUser, db.Cond{
		consts.TcId:       globalUserId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcPassword:     password,
		consts.TcPasswordSalt: salt,
	})
	if err != nil {
		log.Error(err)
		return errs.SetUserPasswordError
	}
	return nil
}

func UnbindUserName(userId int64, addressType int) errs.SystemErrorInfo {
	upd := mysql.Upd{
		consts.TcUpdator: userId,
	}
	if addressType == consts.ContactAddressTypeEmail {
		upd[consts.TcEmail] = consts.BlankString
	} else if addressType == consts.ContactAddressTypeMobile {
		upd[consts.TcMobile] = consts.BlankString
	}

	_, err := mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
		consts.TcId:       userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)
	if err != nil {
		log.Error(err)
		return errs.UnBindLoginNameFail
	}
	return nil
}

func BindUserName(orgId, userId int64, addressType int, username string) errs.SystemErrorInfo {
	lockKey := consts.UserBindLoginNameLock + username
	uid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uid)
	if lockErr != nil {
		log.Error(lockErr)
		return errs.BindLoginNameFail
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	} else {
		return errs.BindLoginNameFail
	}

	//判断是否被其他账户绑定过
	switch addressType {
	case consts.ContactAddressTypeEmail:
		return bindEmail(userId, username)
	case consts.ContactAddressTypeMobile:
		return bindMobile(orgId, userId, username)
	default:
		return errs.NotSupportedContactAddressType
	}
}

// bindMobile 绑定手机号码
func bindMobile(orgId, userId int64, newMobile string) errs.SystemErrorInfo {
	// 老手机绑定关系
	oldRelations, err := GetGlobalUserRelationsByUserId(userId)
	if err != nil {
		log.Errorf("[BindLoginName] GetGlobalUserRelationsByUserId userId:%v, err:%v", userId, err)
		return err
	}

	if oldRelations.Mobile == newMobile {
		return errs.CanNotBindSameMobile
	}

	// 新手机绑定关系
	newRelations, err := GetGlobalUserRelationsByMobile(newMobile)
	if err != nil {
		log.Errorf("[BindLoginName] GetGlobalUserRelationsByMobile mobile:%v, err:%v", newMobile, err)
		return err
	}

	// 如果要绑定的手机有绑定user，则要检查下是否有相同org，如果有，则不允绑定
	if newRelations.GlobalUserId > 0 {
		var userIds []int64
		if len(oldRelations.BindUserIds) > 0 {
			userIds = append(newRelations.BindUserIds, oldRelations.BindUserIds...)
		} else {
			userIds = append(newRelations.BindUserIds, userId)
		}
		isRepeated, err := checkUserIdsOrgIsRepeated(userIds)
		if err != nil {
			return err
		}
		if isRepeated {
			return errs.SameUserInOrg
		}

		// 如果user绑定了全局用户
		if oldRelations.GlobalUserId > 0 {
			err := replaceMobileByOld(newRelations, oldRelations)
			if err != nil {
				return err
			}
		} else {
			// 没有绑定过手机的第三方用户
			err := firstBindOldMobile(userId, newRelations)
			if err != nil {
				return err
			}
		}
	} else {
		// 绑定新手机，如果老的存在，则直接替换就好了
		if oldRelations.GlobalUserId > 0 {
			return replaceMobileByNew(newRelations, oldRelations)
		} else {
			// user并没有绑定手机号，则一定是第三方账号，这个时候新手机号要创建globalUser和本地账号，并与当前第三方账号关联
			err := firstBindNewMobile(userId, newMobile)
			if err != nil {
				return err
			}
		}
	}

	if err2 := dao.GetGlobalUser().UpdateGlobalLastLoginInfo(newRelations.GlobalUserId, userId, orgId); err2 != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	// 上报事件
	e := &commonvo.UserEvent{
		OrgId:  orgId,
		UserId: userId,
		New:    newMobile,
	}

	openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
	openTraceIdStr := cast.ToString(openTraceId)

	report.ReportUserEvent(msgPb.EventType_UserBindMobile, openTraceIdStr, e)

	return nil
}

// firstBindNewMobile 创建一个本地账号，以及全局账号，并且绑定在一起，发生在一个第三方账号第一次绑定一个新手机号
func firstBindNewMobile(userId int64, mobile string) errs.SystemErrorInfo {
	relationId, err := idfacade.ApplyPrimaryIdRelaxed(po.TableNamePpmOrgGlobalUserRelation)
	if err != nil {
		log.Errorf("[firstBindMobile] ApplyPrimaryIdRelaxed err:%v", err)
		return err
	}

	dbErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		userBo, err := UserRegister(bo.UserSMSRegisterInfo{
			PhoneNumber:    mobile,
			SourceChannel:  consts.AppSourceChannelWeb,
			SourcePlatform: consts.AppSourceChannelWeb,
			Name:           "用户_" + str.Last(mobile, 4),
		}, tx)
		if err != nil {
			return err
		}

		relations := []*po.PpmOrgGlobalUserRelation{
			{Id: relationId, UserId: userId, GlobalUserId: userBo.GlobalUserId},
		}
		err2 := dao.GetGlobalUserRelation().CreateRelations(relations, tx)
		if err2 != nil {
			log.Errorf("[firstBindNewMobile] CreateRelations, relations:%v, err :%v", relations, err2)
			return err2
		}

		return nil
	})
	if dbErr != nil {
		log.Errorf("[firstBindMobile] err :%v", dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return nil
}

// firstBindOldMobile 第三方账号第一次绑定手机号，而且手机是老手机
func firstBindOldMobile(userId int64, newRelations *bo.GlobalUserRelations) errs.SystemErrorInfo {
	relationId, err := idfacade.ApplyPrimaryIdRelaxed(po.TableNamePpmOrgGlobalUserRelation)
	if err != nil {
		log.Errorf("[firstBindMobile] ApplyPrimaryIdRelaxed err:%v", err)
		return err
	}

	relations := []*po.PpmOrgGlobalUserRelation{
		{Id: relationId, UserId: userId, GlobalUserId: newRelations.GlobalUserId},
	}
	dbErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		return dao.GetGlobalUserRelation().CreateRelations(relations, tx)
	})
	if dbErr != nil {
		log.Errorf("[firstBindOldMobile] err :%v", dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return nil
}

// replaceMobileByNew 替换手机号码，而且需要绑定的手机号为新号
func replaceMobileByNew(newRelations, oldRelations *bo.GlobalUserRelations) errs.SystemErrorInfo {
	// 如果新手机没有注册过，则可以直接替换即可
	dbErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := dao.GetGlobalUser().UpdateMobile(oldRelations.GlobalUserId, newRelations.Mobile, tx)
		if err != nil {
			log.Errorf("[replaceMobileByNew] updateMobile err :%v", err)
			return err
		}

		return nil
	})
	if dbErr != nil {
		log.Errorf("[replaceMobileByNew] err :%v", dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return nil
}

// replaceMobileByOld 替换手机号码，而且需要绑定的手机号码已经绑定了其他用户
func replaceMobileByOld(newRelations, oldRelations *bo.GlobalUserRelations) errs.SystemErrorInfo {
	dbErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := dao.GetGlobalUserRelation().UpdateUsersToNewGlobalId(oldRelations.GlobalUserId, newRelations.GlobalUserId, tx)
		if err != nil {
			log.Errorf("[replaceMobileByOld] UpdateUsersToNewGlobalId from:%v, to:%v, err :%v", oldRelations.GlobalUserId, newRelations.GlobalUserId, err)
			return err
		}

		return dao.GetGlobalUser().Delete(oldRelations.GlobalUserId, tx)
	})
	if dbErr != nil {
		log.Errorf("[replaceMobileByOld] err :%v", dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return nil
}

func bindEmail(userId int64, username string) errs.SystemErrorInfo {
	err := CheckLoginNameIsExist(consts.ContactAddressTypeEmail, username)
	if err != nil {
		if err.Code() != errs.UserNotExist.Code() {
			return err
		}
	} else {
		return errs.EmailAlreadyBindByOtherAccountError
	}

	upd := mysql.Upd{
		consts.TcUpdator: userId,
		consts.TcEmail:   username,
	}

	_, dbErr := mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
		consts.TcId:       userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)
	if dbErr != nil {
		log.Error(dbErr)
		return errs.UnBindLoginNameFail
	}
	return nil
}

func RetrievePasswordByAccountName(orgId int64, loginName, password string) errs.SystemErrorInfo {
	userPo := &po.PpmOrgUser{}
	err := mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcLoginName: loginName,
	}, userPo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			log.Error(err)
			return errs.UserNotExist
		}
		return errs.MysqlOperateError
	}
	salt := userPo.PasswordSalt
	newPassword := util.PwdEncrypt(password, userPo.PasswordSalt)
	_, err = mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       userPo.Id,
	}, mysql.Upd{
		consts.TcPassword:     newPassword,
		consts.TcPasswordSalt: salt,
	})
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	return nil
}
