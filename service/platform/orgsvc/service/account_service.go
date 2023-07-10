package service

import (
	"strings"

	"github.com/star-table/polaris-backend/service/platform/orgsvc/dao"

	"github.com/star-table/common/core/config"

	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"

	"github.com/star-table/common/core/util/md5"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
)

func SetPassword(req orgvo.SetPasswordReqVo) errs.SystemErrorInfo {
	userId := req.UserId

	targetPassword := strings.TrimSpace(req.Input.Password)
	pwdLen := len(targetPassword)
	if pwdLen < 6 || pwdLen > 16 {
		return errs.PwdFormatError
	}

	suc := format.VerifyPwdFormat(targetPassword)
	if !suc {
		return errs.PwdFormatError
	}

	globalUser, err := domain.GetGlobalUserByUserId(userId)
	if err != nil {
		log.Error(err)
		return err
	}

	if strings.TrimSpace(globalUser.Password) != consts.BlankString {
		return errs.PwdAlreadySettingsErr
	}

	salt := md5.Md5V(uuid.NewUuid())
	pwd := util.PwdEncrypt(targetPassword, salt)
	err = domain.SetUserPassword(globalUser.Id, pwd, salt, userId)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func ResetPassword(req orgvo.ResetPasswordReqVo) errs.SystemErrorInfo {
	userId := req.UserId
	input := req.Input

	globalUser, err := domain.GetGlobalUserByUserId(userId)
	if err != nil {
		log.Error(err)
		return err
	}

	if strings.TrimSpace(globalUser.Password) == consts.BlankString {
		return errs.PasswordNotSetError
	}

	salt := globalUser.PasswordSalt

	targetPassword := strings.TrimSpace(input.NewPassword)
	pwdLen := len(targetPassword)
	if pwdLen < 6 || pwdLen > 18 {
		return errs.PwdFormatError
	}

	suc := format.VerifyPwdFormat(targetPassword)
	if !suc {
		return errs.PwdFormatError
	}

	currentPwd := util.PwdEncrypt(input.CurrentPassword, salt)
	if currentPwd != globalUser.Password {
		return errs.PasswordNotMatchError
	}

	newPassword := util.PwdEncrypt(targetPassword, salt)
	err = domain.SetUserPassword(globalUser.Id, newPassword, salt, userId)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func RetrievePassword(req orgvo.RetrievePasswordReqVo) errs.SystemErrorInfo {
	input := req.Input
	needAuthCode := false
	authCode := ""
	if input.AuthCode != nil {
		authCode = *input.AuthCode
		needAuthCode = true
	}
	username := input.Username
	password := input.NewPassword

	targetPassword := strings.TrimSpace(password)
	suc := format.VerifyPwdFormat(targetPassword)
	if !suc {
		return errs.PwdFormatError
	}

	globalUser, err := domain.GetGlobalUserByMobile(username)
	if err != nil {
		if err == db.ErrNoMoreRows {
			// 走另一个分支，查询user表信息
			return domain.RetrievePasswordByAccountName(req.OrgId, username, targetPassword)
		} else {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}

	addressType := consts.ContactAddressTypeMobile
	if format.VerifyEmailFormat(username) {
		addressType = consts.ContactAddressTypeEmail
	}

	runMode := config.GetConfig().Application.RunMode
	if (runMode == 1 || runMode == 2) && needAuthCode {
		err1 := domain.AuthCodeVerify(consts.AuthCodeTypeRetrievePwd, addressType, username, authCode)
		if err1 != nil {
			log.Error(err1)
			return err1
		}
	}

	salt := globalUser.PasswordSalt
	newPassword := util.PwdEncrypt(targetPassword, salt)
	err2 := domain.SetUserPassword(globalUser.Id, newPassword, salt, 0)
	if err2 != nil {
		log.Error(err2)
		return err2
	}

	return nil
}

func UnbindLoginName(req orgvo.UnbindLoginNameReqVo) errs.SystemErrorInfo {
	userId := req.UserId
	input := req.Input

	addressType := input.AddressType
	authCode := input.AuthCode

	userBo, _, err := domain.GetUserBo(userId)
	if err != nil {
		log.Error(err)
		return err
	}

	username := ""
	if addressType == consts.ContactAddressTypeEmail {
		if strings.TrimSpace(userBo.Email) == consts.BlankString {
			return errs.EmailNotBindError
		}
		username = userBo.Email
	} else if addressType == consts.ContactAddressTypeMobile {
		// 禁止解绑手机号，因统一账户依赖手机号，解绑将导致数据丢失
		return errs.NotSupportedContactAddressType
		//if strings.TrimSpace(userBo.Mobile) == consts.BlankString {
		//	return errs.MobileNotBindError
		//}
		//username = userBo.Mobile
	} else {
		return errs.NotSupportedContactAddressType
	}

	//if strings.TrimSpace(userBo.Email) == consts.BlankString || strings.TrimSpace(userBo.Mobile) == consts.BlankString {
	//	return errs.HaveNoContract
	//}

	err1 := domain.AuthCodeVerify(consts.AuthCodeTypeUnBind, addressType, username, authCode)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	err = domain.UnbindUserName(userId, addressType)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func BindLoginName(req orgvo.BindLoginNameReqVo) errs.SystemErrorInfo {
	userId := req.UserId
	input := req.Input

	username := input.Address
	addressType := input.AddressType
	authCode := input.AuthCode

	switch addressType {
	case consts.ContactAddressTypeEmail:
		userBo, _, err := domain.GetUserBo(userId)
		if err != nil {
			log.Error(err)
			return err
		}
		if strings.TrimSpace(userBo.Email) != consts.BlankString {
			return errs.EmailAlreadyBindError
		}
	case consts.ContactAddressTypeMobile:
	default:
		return errs.NotSupportedContactAddressType
	}

	err1 := domain.AuthCodeVerify(consts.AuthCodeTypeBind, addressType, username, authCode)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	err := domain.BindUserName(req.OrgId, userId, addressType, username)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func CheckLoginName(req orgvo.CheckLoginNameReqVo) errs.SystemErrorInfo {
	input := req.Input

	addressType := input.AddressType
	address := input.Address
	if addressType == consts.ContactAddressTypeEmail {
		_, err := domain.GetUserInfoByEmail(address)
		if err != nil {
			if err.Code() == errs.UserNotExist.Code() {
				return errs.EmailNotBindAccountError
			}
		}
	}

	if addressType == consts.ContactAddressTypeMobile {
		_, err := domain.GetUserInfoByMobile(address)
		if err != nil {
			if err.Code() == errs.UserNotExist.Code() {
				return errs.MobileNotBindAccountError
			}
		}
	}

	return nil
}

func VerifyOldName(req orgvo.UnbindLoginNameReqVo) errs.SystemErrorInfo {
	userId := req.UserId
	input := req.Input

	addressType := input.AddressType
	authCode := input.AuthCode

	userBo, _, err := domain.GetUserBo(userId)
	if err != nil {
		log.Error(err)
		return err
	}

	username := ""
	if addressType == consts.ContactAddressTypeEmail {
		if strings.TrimSpace(userBo.Email) == consts.BlankString {
			return errs.EmailNotBindError
		}
		username = userBo.Email
	} else if addressType == consts.ContactAddressTypeMobile {
		globalUser, err := domain.GetGlobalUserByUserId(userId)
		if err != nil {
			log.Errorf("[VerifyOldName] GetGlobalUserByUserId userId:%v, err:%v", userId, err)
			return err
		}
		if strings.TrimSpace(globalUser.Mobile) == consts.BlankString {
			return errs.MobileNotBindError
		}
		username = globalUser.Mobile
	} else {
		return errs.NotSupportedContactAddressType
	}

	err1 := domain.AuthCodeVerify(consts.AuthCodeTypeUnBind, addressType, username, authCode)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	//将行为记录到缓存，在有效期内可用
	cacheErr := domain.SetChangeLoginNameSign(req.OrgId, req.UserId, input.AddressType)
	if cacheErr != nil {
		log.Error(cacheErr)
		return cacheErr
	}

	return nil
}

func ChangeLoginName(req orgvo.BindLoginNameReqVo) errs.SystemErrorInfo {
	userId := req.UserId
	input := req.Input

	username := input.Address
	addressType := input.AddressType
	authCode := input.AuthCode

	//userBo, _, err := domain.GetUserBo(userId)
	//if err != nil{
	//	log.Error(err)
	//	return err
	//}

	if ok, _ := slice.Contain([]int{consts.ContactAddressTypeEmail, consts.ContactAddressTypeMobile}, addressType); !ok {
		return errs.NotSupportedContactAddressType
	}
	cacheErr := domain.GetChangeLoginNameSign(req.OrgId, req.UserId, req.Input.AddressType)
	if cacheErr != nil {
		log.Error(cacheErr)
		return cacheErr
	}

	err1 := domain.AuthCodeVerify(consts.AuthCodeTypeBind, addressType, username, authCode)
	if err1 != nil {
		log.Error(err1)
		return err1
	}

	err := domain.BindUserName(req.OrgId, userId, addressType, username)
	if err != nil {
		log.Error(err)
		return err
	}

	clearErr := domain.ClearChangeLoginNameSign(req.OrgId, req.UserId, req.Input.AddressType)
	if clearErr != nil {
		log.Error(clearErr)
		return clearErr
	}
	return nil
}

func DisbandThirdAccount(orgId, userId int64, sourceChannel string) errs.SystemErrorInfo {
	//判断当前企业是否绑定了该平台（绑定平台的组织不允许自己私自解绑）
	outPo := po.OutOrgInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, &outPo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			//判断该用户是否绑定了当前第三方平台
			outUserPo := po.PpmOrgUserOutInfo{}
			userErr := mysql.SelectOneByCond(consts.TableUserOutInfo, db.Cond{
				consts.TcIsDelete:      consts.AppIsNoDelete,
				consts.TcOrgId:         0, //个人绑定外部信息都是0
				consts.TcUserId:        userId,
				consts.TcSourceChannel: sourceChannel,
			}, &outUserPo)
			if userErr != nil {
				if userErr == db.ErrNoMoreRows {
					return errs.NotDisbandCurrentSourceChannel
				} else {
					log.Error(userErr)
					return errs.MysqlOperateError
				}
			}

			_, updateErr := mysql.UpdateSmartWithCond(consts.TableUserOutInfo, db.Cond{
				consts.TcId: outUserPo.Id,
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcUpdator:  userId,
			})
			if updateErr != nil {
				log.Error(updateErr)
				return errs.MysqlOperateError
			}
		} else {
			log.Error(err)
			return errs.MysqlOperateError
		}
	} else {
		return errs.DisbandThirdAccountError
	}

	return nil
}

func ThirdAccountBindList(orgId, userId int64) (*orgvo.ThirdAccountListData, errs.SystemErrorInfo) {
	userIds, err := dao.GetGlobalUserRelation().GetUserIdsByUserId(userId)
	isBindMobile := false
	if err == nil {
		isBindMobile = true
	}

	result := &orgvo.ThirdAccountListData{IsBindMobile: isBindMobile, List: []*orgvo.ThirdAccountData{
		{
			BindPlatform: consts.AppSourcePlatformPersonWeixin,
		},
	}}

	//判断该用户是否绑定了当前第三方平台
	thirdList, err := dao.GetThirdLogin().GetThirdLoginInfoByUserIds(userIds)
	if err != nil {
		return nil, errs.MysqlOperateError
	}

	for _, data := range result.List {
		for _, info := range thirdList {
			if info.SourcePlatform == data.BindPlatform {
				data.BindStatus = 1
				data.BindTime = info.CreateTime.Unix()
				break
			}
		}
	}

	return result, nil
}
