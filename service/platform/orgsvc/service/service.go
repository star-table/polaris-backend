package service

import (
	"context"
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/gin-gonic/gin"
)

var log = *logger.GetDefaultLogger()

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GetCtxParameters(ctx context.Context, key string) (string, error) {
	gc, err := GinContextFromContext(ctx)
	if err != nil {
		return "", err
	}
	v := gc.GetString(key)
	return v, nil
}

func GetCurrentUser(ctx context.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	return GetCurrentUserWithCond(ctx, true, true)
}

func GetCurrentUserWithoutOrgVerify(ctx context.Context) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	return GetCurrentUserWithCond(ctx, false, true)
}

//payVerify 兼容部分接口不需要验证用户可用性
func GetCurrentUserWithCond(ctx context.Context, orgVerify bool, payVerify bool) (*bo.CacheUserInfoBo, errs.SystemErrorInfo) {
	token, err := GetCtxParameters(ctx, consts.AppHeaderTokenName)

	// 模板预览需求-token判断 认证处理
	if token == consts.PreviewTplToken {
		return &bo.CacheUserInfoBo{
			OutUserId:     "",
			SourceChannel: "",
			UserId:        consts.PreviewTplUserId,
			CorpId:        "",
			OrgId:         consts.PreviewTplOrgId,
		}, nil
	} else if token == consts.PreviewTplTokenForWrite {
		return &bo.CacheUserInfoBo{
			OutUserId:     "",
			SourceChannel: "",
			UserId:        consts.PreviewOrWriteTplUserId,
			CorpId:        "",
			OrgId:         consts.PreviewTplOrgId,
		}, nil
	}

	if err != nil || token == "" {
		return nil, errs.BuildSystemErrorInfo(errs.TokenNotExist)
	} else {
		cacheUserInfoJson, err := cache.Get(sconsts.CacheUserToken + token)
		if err != nil {
			log.Error(strs.ObjectToString(err))
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		if cacheUserInfoJson == "" {
			log.Error("token失效")
			return nil, errs.BuildSystemErrorInfo(errs.TokenExpires)
		}
		cacheUserInfo := &bo.CacheUserInfoBo{}
		err = json.FromJson(cacheUserInfoJson, cacheUserInfo)
		_, _ = cache.Expire(sconsts.CacheUserToken+token, consts.CacheUserTokenExpire)
		if err != nil {
			log.Error(strs.ObjectToString(err))
			return nil, errs.BuildSystemErrorInfo(errs.TokenExpires)
		}
		//判断用户组织状态
		userName := ""
		if cacheUserInfo.OrgId != 0 && orgVerify {
			baseUserInfo, err := GetBaseUserInfo(cacheUserInfo.OrgId, cacheUserInfo.UserId)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			err = baseUserInfoOrgStatusCheck(*baseUserInfo)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			userName = baseUserInfo.Name
		}

		if orgVerify && payVerify && !domain.GetWhiteListVipOrg(cacheUserInfo.OrgId) && config.GetApplication().RunType == 0 {
			//查看有没有外部组织
			outOrgInfo, outOrgInfoErr := domain.GetOrgOutInfoWithoutLocal(cacheUserInfo.OrgId)
			if outOrgInfoErr != nil {
				if outOrgInfoErr == errs.OrgOutInfoNotExist {
					//没有外部信息就不用判断是否付费
					return cacheUserInfo, nil
				} else {
					log.Error(outOrgInfoErr)
					return nil, outOrgInfoErr
				}
			}
			if cacheUserInfo.OutUserId == "" {
				userOutInfo, err := domain.GetUserOutInfoByUserIdAndOrgId(cacheUserInfo.UserId, cacheUserInfo.OrgId, cacheUserInfo.SourceChannel)
				if err != nil {
					if err == errs.NotDisbandCurrentSourceChannel {
						return cacheUserInfo, nil
					} else {
						log.Error(err)
						return nil, err
					}
				}
				cacheUserInfo.OutUserId = userOutInfo.OutUserId
			}
			checkErr := GetUserPay(cacheUserInfo.OrgId, cacheUserInfo.UserId, outOrgInfo.OutOrgId, cacheUserInfo.OutUserId, outOrgInfo.SourceChannel, userName)
			if checkErr != nil {
				log.Error(checkErr)
				return nil, checkErr
			}
		}
		return cacheUserInfo, nil
	}
}

//获取用户是否在授权范围
func GetUserPay(orgId int64, userId int64, orgOrgId, outUserId, sourceChannel, userName string) errs.SystemErrorInfo {
	orgConfig, err := domain.GetOrgConfig(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	if orgConfig.PayLevel == consts.PayLevelStandard {
		return nil
	}
	// 私有化部署时 无需校验是否在付费范围内
	appDeployType := domain.GetAppDeployType(config.GetConfig().Application.RunMode)
	isPrivateDeploy := appDeployType == "private"
	if isPrivateDeploy {
		return nil
	}

	//如果是付费的话要判断用户是否在付费范围内
	//return domain.CheckUserPay(orgId, userId, outUserId)
	return nil
}

//用户信息所在组织状态监测
func baseUserInfoOrgStatusCheck(baseUserInfo bo.BaseUserInfoBo) errs.SystemErrorInfo {
	orgInfoBo, err := domain.GetBaseOrgInfo(baseUserInfo.OrgId)
	if err != nil {
		log.Errorf("[baseUserInfoOrgStatusCheck] GetBaseOrgInfo err:%v, orgId:%v", err, baseUserInfo.OrgId)
		return err
	}
	if baseUserInfo.OrgUserStatus != consts.AppStatusEnable {
		if orgInfoBo.OutOrgId != "" {
			remindInfo := bo.OrgAndUserInfo{
				OrgName:       orgInfoBo.OrgName,
				UserName:      baseUserInfo.Name,
				SourceChannel: orgInfoBo.SourceChannel,
			}
			return errs.BuildSystemErrorInfoWithMessage(errs.OrgUserInvalid, json.ToJsonIgnoreError(remindInfo))
		}
		return errs.OrgUserUnabled
	}
	if baseUserInfo.OrgUserCheckStatus != consts.AppCheckStatusSuccess {
		return errs.OrgUserCheckStatusUnabled
	}
	if baseUserInfo.OrgUserIsDelete == consts.AppIsDeleted {
		return errs.OrgUserDeleted
	}
	return nil
}

func UpdateCacheUserInfoOrgId(token string, orgId, userId int64, updOutInfo bool) errs.SystemErrorInfo {
	cacheUserInfoJson, err := cache.Get(sconsts.CacheUserToken + token)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if cacheUserInfoJson == "" {
		return errs.BuildSystemErrorInfo(errs.TokenExpires)
	}
	cacheUserInfo := &bo.CacheUserInfoBo{}
	_ = json.FromJson(cacheUserInfoJson, cacheUserInfo)

	if updOutInfo {
		//查询组织集成信息
		outInfo, _ := domain.GetOrgOutInfoWithoutLocal(orgId)
		if outInfo != nil {
			cacheUserInfo.SourceChannel = outInfo.SourceChannel
			cacheUserInfo.CorpId = outInfo.OutOrgId
			// 查询用户在当前企业的外部信息
			userOutInfo, busErr := domain.GetUserOutInfoByUserIdByOrgIds(userId, []int64{orgId}, cacheUserInfo.SourceChannel)
			if busErr != nil && busErr != errs.NotDisbandCurrentSourceChannel {
				log.Error(busErr)
				return busErr
			}
			if userOutInfo != nil {
				cacheUserInfo.OutUserId = userOutInfo.OutUserId
			}
		} else {
			cacheUserInfo.SourceChannel = consts.AppSourceChannelWeb
			cacheUserInfo.CorpId = ""
			cacheUserInfo.OutUserId = ""
		}
	}

	//更新缓存用户的orgId
	cacheUserInfo.OrgId = orgId
	cacheUserInfo.UserId = userId
	cacheUserInfoJson = json.ToJsonIgnoreError(cacheUserInfo)
	err = cache.SetEx(sconsts.CacheUserToken+token, cacheUserInfoJson, consts.CacheUserTokenExpire)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}

func UpdateCacheUserFsInfo(token string, tenantKey, openId, sourceChannel string) errs.SystemErrorInfo {
	cacheUserInfoJson, err := cache.Get(sconsts.CacheUserToken + token)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if cacheUserInfoJson == "" {
		return errs.BuildSystemErrorInfo(errs.TokenExpires)
	}
	cacheUserInfo := &bo.CacheUserInfoBo{}
	_ = json.FromJson(cacheUserInfoJson, cacheUserInfo)

	cacheUserInfo.SourceChannel = sourceChannel
	cacheUserInfo.OutUserId = openId
	cacheUserInfo.CorpId = tenantKey
	cacheUserInfoJson = json.ToJsonIgnoreError(cacheUserInfo)
	err = cache.SetEx(sconsts.CacheUserToken+token, cacheUserInfoJson, consts.CacheUserTokenExpire)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	return nil
}
