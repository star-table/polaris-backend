package service

import (
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/random"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
)

func UserLogin(req vo.UserLoginReq) (*vo.UserLoginResp, errs.SystemErrorInfo) {
	loginType := req.LoginType
	var userBo *bo.UserInfoBo = nil
	var err errs.SystemErrorInfo = nil
	// 1,3验证码登陆， 2密码登陆，4密码注册，私有化部署可用
	if loginType == 1 || loginType == 3 {
		userBo, err = UserAuthCodeLogin(req)
	} else if loginType == 2 || loginType == 4 {
		userBo, err = UserPwdLogin(req)
	}
	if err != nil {
		if err == errs.AccountNotBelongToCurrentFs {
			return &vo.UserLoginResp{
				Token:       "",
				UserID:      0,
				OrgID:       0,
				OrgName:     "",
				OrgCode:     "",
				Name:        "",
				Avatar:      "",
				NeedInitOrg: false,
				NotFsMobile: true,
			}, nil
		}
		log.Error(err)
		return nil, err
	}
	if userBo == nil {
		//暂时不支持的登录类型
		return nil, errs.BuildSystemErrorInfo(errs.UnSupportLoginType)
	}

	return userLogin(userBo, req.SourceChannel, consts.TokenTypeNormal)
}

func userLogin(userBo *bo.UserInfoBo, sourceChannel string, tokenType int) (*vo.UserLoginResp, errs.SystemErrorInfo) {
	token := random.GenTokenWithOrgId(userBo.OrgID, userBo.ID)
	res := &vo.UserLoginResp{
		Token:       token,
		UserID:      userBo.ID,
		OrgID:       userBo.OrgID,
		Name:        userBo.Name,
		Avatar:      userBo.Avatar,
		NeedInitOrg: userBo.OrgID == 0,
	}

	if userBo.OrgID != 0 {
		organizationBo, err := domain.GetOrgBoById(userBo.OrgID)
		if err != nil {
			log.Error(err)
		} else {
			res.OrgCode = organizationBo.Code
			res.OrgName = organizationBo.Name
		}
	}

	cacheUserBo := &bo.CacheUserInfoBo{
		UserId:        userBo.ID,
		SourceChannel: domain.GetOrgSourceChannel(userBo.OrgID, sourceChannel),
		OrgId:         userBo.OrgID,
		TokenType:     tokenType,
	}
	cacheErr := cache.SetEx(sconsts.CacheUserToken+res.Token, json.ToJsonIgnoreError(cacheUserBo), consts.CacheUserTokenExpire)
	if cacheErr != nil {
		log.Error(cacheErr)
	}

	//更新上次登录时间
	if tokenType != consts.TokenTypeShare {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
				}
			}()
			if err := domain.UserLoginHook(userBo.ID, userBo.OrgID, sourceChannel); err != nil {
				log.Error(err)
			}
		}()
	}

	return res, nil
}

func FsValidateMobileForLogin(mobile string, codeToken *string) errs.SystemErrorInfo {
	if codeToken == nil {
		return errs.ParamError
	}
	mobile = util.GetMobile(mobile)
	codeInfoJson, redisErr := cache.Get(sconsts.CacheFsAuthCodeToken + *codeToken)
	if redisErr != nil {
		log.Error(redisErr)
		return errs.RedisOperateError
	}
	if codeInfoJson == "" {
		return errs.CodeTokenInvalid
	}
	codeInfo := vo.FeiShuAuthCodeResp{}
	_ = json.FromJson(codeInfoJson, &codeInfo)

	return nil
}

func UserAuthCodeLogin(req vo.UserLoginReq) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	loginName := req.LoginName
	authCode := req.AuthCode
	loginType := req.LoginType

	addressType := consts.ContactAddressTypeEmail
	if loginType == consts.LoginTypeSMSCode {
		addressType = consts.ContactAddressTypeMobile
	}

	name := "用户_" + str.Last(loginName, 4)
	inviteCode := ""
	if req.Name != nil {
		name = *req.Name
	}
	if req.InviteCode != nil {
		inviteCode = *req.InviteCode
	}
	if authCode == nil {
		log.Error("验证码不能为空")
		return nil, errs.AuthCodeIsNull
	}

	//台湾跳过短信验证 && config.GetEnv() != consts.RunEnvStag
	if config.GetEnv() != consts.RunEnvProdTw {
		err1 := domain.AuthCodeVerify(consts.AuthCodeTypeLogin, addressType, loginName, *authCode)
		if err1 != nil {
			log.Error(err1)
			return nil, err1
		}
	}

	var (
		loginUserBo *bo.UserInfoBo
		isRegister  = false
		err         errs.SystemErrorInfo
	)

	switch loginType {
	case consts.LoginTypeSMSCode:
		loginUserBo, isRegister, err = loginMobile(req, name, inviteCode)
	case consts.LoginTypeMail:
		loginUserBo, isRegister, err = loginEmail(req, name, inviteCode)
	default:
		return nil, errs.UnSupportLoginType
	}
	if err != nil {
		return nil, err
	}

	err = userAlreadyRegisterHandleHook(req.InviteCode, loginUserBo, isRegister)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return loginUserBo, nil
}

func loginMobile(req vo.UserLoginReq, name, inviteCode string) (*bo.UserInfoBo, bool, errs.SystemErrorInfo) {
	//做登录和自动注册逻辑
	userBo, err := domain.GetUserInfoByMobile(req.LoginName)
	if err != nil {
		if err == errs.UserNotExist {
			if req.CodeToken != nil && *req.CodeToken != "" {
				err := FsValidateMobileForLogin(req.LoginName, req.CodeToken)
				if err != nil {
					log.Error(err)
					return nil, false, err
				}
			}
			log.Infof("用户%s未注册，开始注册....", req.LoginName)
			userBo, err = domain.UserRegister(bo.UserSMSRegisterInfo{
				PhoneNumber:    req.LoginName,
				SourceChannel:  req.SourceChannel,
				SourcePlatform: req.SourcePlatform,
				Name:           name,
				InviteCode:     inviteCode,
			})
			if err != nil {
				log.Error(err)
				return nil, false, err
			}
			return userBo, true, nil
		} else {
			log.Error(err)
			return nil, false, errs.AccountNotRegister
		}
	}

	return userBo, false, nil
}

func loginEmail(req vo.UserLoginReq, name, inviteCode string) (*bo.UserInfoBo, bool, errs.SystemErrorInfo) {
	userBo, err := domain.GetUserInfoByEmail(req.LoginName)
	if err != nil {
		if err == errs.UserNotExist {
			log.Infof("用户%s未注册，开始注册....", req.LoginName)
			userBo, err = domain.UserRegister(bo.UserSMSRegisterInfo{
				Email:          req.LoginName,
				SourceChannel:  req.SourceChannel,
				SourcePlatform: req.SourcePlatform,
				Name:           name,
				InviteCode:     inviteCode,
			})
			if err != nil {
				log.Error(err)
				return nil, false, err
			}
			return userBo, true, nil
		} else {
			log.Error(err)
			return nil, false, errs.AccountNotRegister
		}
	}

	return userBo, false, nil
}

func userAlreadyRegisterHandleHook(inviteCode *string, userBo *bo.UserInfoBo, isRegister bool) errs.SystemErrorInfo {
	//判断邀请逻辑
	if inviteCode != nil {
		inviteCodeInfo, err := domain.GetUserInviteCodeInfo(*inviteCode)
		if err != nil {
			log.Error(err)
			return err
		}
		orgId := inviteCodeInfo.OrgId
		inviterId := inviteCodeInfo.InviterId
		userId := userBo.ID

		// 非注册的时候，检查下组织是否冲突，如果该globalUser下已经有user加入了该组织，则不能再次加入
		if !isRegister {
			isIn, err := domain.CheckUserOrgIsConflict(userId, orgId)
			if err != nil {
				return err
			}
			// 已经有用户在这个orgId里面了，可以直接返回
			if isIn {
				// 查询一下是不是已经绑定过了，但是有可能之前审核不通过，现在重新审核 check_status需要重置为1
				return domain.ResetUserCheckStatusWait(orgId, userId)
			}
		}

		//这里为用户切换组织
		userBo.OrgID = orgId

		//修改用户默认组织
		_, updateUserInfoErr := domain.UpdateUserDefaultOrg(userId, orgId)
		if updateUserInfoErr != nil {
			log.Error(updateUserInfoErr)
		}

		err = domain.AddOrgMember(orgId, userId, inviterId, true, false)
		if err != nil {
			log.Error(err)
			return err
		}

		//增加动态
		asyn.Execute(func() {
			domain.PushOrgTrends(bo.OrgTrendsBo{
				OrgId:         orgId,
				PushType:      consts.PushTypeApplyJoinOrg,
				TargetMembers: []int64{userId},
				SourceChannel: userBo.SourceChannel,
				OperatorId:    inviterId,
			})
		})
	}
	return nil
}

func ImportMemberRegisterHandleHook(importInfo orgvo.ImportUserInfo, userBo bo.UserInfoBo) errs.SystemErrorInfo {
	var err errs.SystemErrorInfo
	orgId := importInfo.OrgId
	operateUid := importInfo.OperateUid
	userId := userBo.ID

	//这里为用户切换组织
	userBo.OrgID = orgId

	//修改用户默认组织
	_, updateUserInfoErr := domain.UpdateUserDefaultOrg(userId, orgId)
	if updateUserInfoErr != nil {
		log.Error(updateUserInfoErr)
	}

	err = domain.AddOrgMember(orgId, userId, operateUid, false, false)
	if err != nil {
		log.Error(err)
		return err
	}

	//增加动态
	asyn.Execute(func() {
		domain.PushOrgTrends(bo.OrgTrendsBo{
			OrgId:         orgId,
			PushType:      consts.PushTypeApplyJoinOrg,
			TargetMembers: []int64{userId},
			SourceChannel: userBo.SourceChannel,
			OperatorId:    operateUid,
		})
	})
	return nil
}

//用户密码登录
func UserPwdLogin(req vo.UserLoginReq) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	if req.Password == nil {
		return nil, errs.PasswordEmptyError
	}
	loginName := req.LoginName
	pwd := *req.Password

	freezeValue, err := domain.GetPwdLoginCodeVerifyFreeze(consts.AuthCodeTypeLogin, consts.ContactAddressTypeMobile, loginName)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if freezeValue != "" {
		return nil, errs.PwdLoginLimitError
	}

	runMode := config.GetConfig().Application.RunMode
	userBo, err := domain.GetUserInfoByPwd(loginName, pwd)
	if err != nil {
		if err == errs.PwdLoginUsrOrPwdNotMatch {
			if runMode != 1 && runMode != 2 && req.LoginType == 4 {
				return nil, err
			}
			c, err := domain.SetPwdLoginCodeVerifyFailTimesIncrement(consts.AuthCodeTypeLogin, consts.ContactAddressTypeMobile, loginName)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			if c > 4 {
				_ = domain.ClearPwdLoginCodeVerifyFailTimesIncrement(consts.AuthCodeTypeLogin, consts.ContactAddressTypeMobile, loginName)
				err = domain.SetPwdLoginCodeVerifyFreeze(consts.AuthCodeTypeLogin, consts.ContactAddressTypeMobile, loginName)
				if err != nil {
					log.Error(err)
					return nil, err
				}
			}
		} else if req.LoginType == 4 && err == errs.AccountNotRegister {
			return registerByPassword(req)
		}
		log.Error(err)
		return nil, err
	}

	return userBo, nil
}

func registerByPassword(req vo.UserLoginReq) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	runMode := config.GetConfig().Application.RunMode
	// 如果是手机号注册，则判断下是否是私有化部署，不是报错
	if runMode == 1 || runMode == 2 {
		return nil, errs.BuildSystemErrorInfo(errs.NotSupportedRegisterType)
	}
	if !format.VerifyPwdFormat(*req.Password) {
		return nil, errs.PwdFormatError
	}
	log.Infof("用户%s未注册，开始注册....", req.LoginName)
	userBo, err := domain.UserRegister(bo.UserSMSRegisterInfo{
		PhoneNumber:    req.LoginName,
		SourceChannel:  req.SourceChannel,
		SourcePlatform: req.SourcePlatform,
		Name:           "用户_" + str.Last(req.LoginName, 4),
		Password:       *req.Password,
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return userBo, nil
}

//用户退出
func UserQuit(req orgvo.UserQuitReqVo) errs.SystemErrorInfo {
	return domain.ClearUserCacheInfo(req.Token)
}

func JoinOrgByInviteCode(req orgvo.JoinOrgByInviteCodeReq) errs.SystemErrorInfo {
	if req.InviteCode == "" {
		return errs.InviteCodeEmpty
	}
	userBo, userBoErr := domain.GetUserInfoById(req.UserId)
	if userBoErr != nil {
		log.Error(userBoErr)
		return userBoErr
	}
	err := userAlreadyRegisterHandleHook(&req.InviteCode, userBo, false)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func ExchangeShareToken(req orgvo.ExchangeShareTokenReq) (*vo.UserLoginResp, errs.SystemErrorInfo) {
	resp := projectfacade.CheckShareViewPassword(&projectvo.CheckShareViewPasswordReq{Input: &projectvo.CheckPasswordData{
		ShareKey: req.Input.ShareKey,
		Password: req.Input.Password,
	}})
	if resp.Failure() {
		return nil, resp.Error()
	}

	userBo, _, err2 := domain.GetUserBo(resp.Data.UserId)
	if err2 != nil {
		return nil, err2
	}
	userBo.OrgID = resp.Data.OrgId

	return userLogin(userBo, "", consts.TokenTypeShare)
}
