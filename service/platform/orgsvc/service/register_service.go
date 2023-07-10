/**
2 * @Author: Nico
3 * @Date: 2020/1/31 11:20
4 */
package service

import (
	"strings"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/random"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
)

func UserRegister(req orgvo.UserRegisterReqVo) (*vo.UserRegisterResp, errs.SystemErrorInfo) {
	input := req.Input
	username := input.UserName
	registerType := input.RegisterType
	authCode := input.AuthCode

	name := ""
	if input.Name != nil {
		name = strings.TrimSpace(*input.Name)
		//检测姓名是否合法
		isNameRight := format.VerifyUserNameFormat(name)
		if !isNameRight {
			return nil, errs.UserNameLenError
		}
	}

	//校验验证码，账号密码登录暂时不需要验证码
	addressType := -1
	if registerType == consts.RegisterTypeSMSCode {
		addressType = consts.ContactAddressTypeMobile
	} else if registerType == consts.RegisterTypeMail {
		addressType = consts.ContactAddressTypeEmail
	}

	if addressType > 0 {
		if authCode == nil {
			return nil, errs.AuthCodeIsNull
		}
		//验证码是否正确
		err1 := domain.AuthCodeVerify(consts.AuthCodeTypeRegister, addressType, username, *authCode)
		if err1 != nil {
			log.Error(err1)
			return nil, err1
		}
	}

	//检测用户名是否存在
	err := domain.CheckLoginNameIsExist(addressType, username)
	if err == nil {
		return nil, errs.AccountAlreadyBindError
	} else {
		if err.Code() != errs.UserNotExist.Code() {
			log.Errorf("[UserRegister] CheckLoginNameIsExist, username:%v, err:%v", username, err)
			return nil, err
		}
	}

	var registerUserBo bo.UserInfoBo
	//注册
	switch registerType {
	case consts.RegisterTypeMail:
		userBo, err := domain.UserRegister(bo.UserSMSRegisterInfo{
			Email:          username,
			SourceChannel:  input.SourceChannel,
			SourcePlatform: input.SourcePlatform,
			Name:           name,
		})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		registerUserBo = *userBo
	case consts.RegisterTypeSMSCode:
		userBo, err := domain.UserRegister(bo.UserSMSRegisterInfo{
			PhoneNumber:    username,
			SourceChannel:  input.SourceChannel,
			SourcePlatform: input.SourcePlatform,
			Name:           name,
			//InviteCode:     inviteCode,
		})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		registerUserBo = *userBo
	default:
		//暂时不支持的注册方式
		return nil, errs.NotSupportedRegisterType
	}

	token := random.GenTokenWithOrgId(registerUserBo.OrgID, registerUserBo.ID)
	//自动登录
	cacheUserBo := &bo.CacheUserInfoBo{
		UserId:        registerUserBo.ID,
		SourceChannel: input.SourceChannel,
	}
	cacheErr := cache.SetEx(sconsts.CacheUserToken+token, json.ToJsonIgnoreError(cacheUserBo), consts.CacheUserTokenExpire)
	if cacheErr != nil {
		log.Error(cacheErr)
	}

	return &vo.UserRegisterResp{
		Token: token,
	}, nil
}
