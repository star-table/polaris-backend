package service

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/star-table/polaris-backend/service/platform/resourcesvc/domain"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/random"
	"github.com/star-table/common/sdk/oss"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func GetOssSignURL(req resourcevo.OssApplySignURLReqVo) (*vo.OssApplySignURLResp, errs.SystemErrorInfo) {
	oc := config.GetOSSConfig()
	if oc == nil {
		log.Error("oss缺少配置")
		return nil, errs.OssError
	}

	input := req.Input
	orgId := req.OrgId
	//userId := req.UserId

	ossKey := str.ParseOssKey(input.URL)
	ossKeyInfo := util.GetOssKeyInfo(ossKey)

	log.Infof("oss key %s, oss key info %s", ossKey, json.ToJsonIgnoreError(ossKeyInfo))
	if ossKeyInfo.OrgId != orgId {
		return nil, errs.NoOperationPermissions
	}
	//TODO 权限

	signUrl, signErr := oss.GetObjectUrl(ossKey, 60*60)
	if signErr != nil {
		log.Error(signErr)
		return nil, errs.BuildSystemErrorInfo(errs.OssError, signErr)
	}

	_, path := str.UrlParse(signUrl)
	//使用了自定义域名，要替换
	signUrl = util.JointUrl(oc.EndPoint, path)

	return &vo.OssApplySignURLResp{
		SignURL: signUrl,
	}, nil
}

func GetOssPostPolicy(req resourcevo.GetOssPostPolicyReqVo) (*vo.OssPostPolicyResp, errs.SystemErrorInfo) {
	oc := config.GetOSSConfig()
	if oc == nil {
		log.Error("oss缺少配置")
		return nil, errs.OssError
	}

	input := req.Input
	orgId := req.OrgId
	userId := req.UserId

	projectId := int64(0)
	issueId := int64(0)
	resourceFolderId := int64(0)
	columnId := ""
	tableId := "0"
	appId := "0"
	if input.AppID != nil {
		appId = *input.AppID
	}
	if input.ProjectID != nil {
		projectId = *input.ProjectID
	}
	if input.IssueID != nil {
		issueId = *input.IssueID
	}
	if input.FolderID != nil {
		resourceFolderId = *input.FolderID
	}
	if input.ColumnID != nil {
		columnId = *input.ColumnID
	}
	if input.TableID != nil {
		tableId = *input.TableID
	}
	var policyConfig *config.OSSPolicyInfo = nil

	policyType := input.PolicyType

	toDay := time.Now()
	year, month, day := toDay.Date()

	//文件名
	fileName := random.RandomFileName()

	//定义callback
	callbackJson := ""
	switch policyType {
	case consts.OssPolicyTypeProjectCover:
		c := config.GetProjectCoverPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})

		policyConfig = &c
	case consts.OssPolicyTypeIssueResource, consts.OssPolicyTypeLesscodeResource:
		c := config.GetIssueResourcePolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyIssueIdConstName:   issueId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c

		callbackBody := fmt.Sprintf("{\"userId\":%d,\"orgId\":%d,\"type\":%d,\"projectId\":%d,\"host\":\"%s\",\"path\":\"%s\",\"filename\":\"%s\",\"issueId\":%d,\"columnId\":\"%s\",\"tableId\":\"%s\",\"appId\":\"%s\",\"bucket\":\"%s\",\"size\":%s,\"format\":%s,\"object\":%s,\"realName\":%s}",
			userId, orgId, policyType, projectId, oc.EndPoint, c.Dir, fileName, issueId, columnId, tableId, appId, oc.BucketName, "${size}", "${imageInfo.format}", "${object}", "${x:filename}")
		callbackBo := bo.OssCallbackBo{
			CallbackBody:     callbackBody,
			CallbackUrl:      policyConfig.CallbackUrl,
			CallbackBodyType: consts.OssCallbackBodyTypeApplicationJson,
		}

		callbackJson = json.ToJsonIgnoreError(callbackBo)
		callbackJson = base64.StdEncoding.EncodeToString([]byte(callbackJson))
	case consts.OssPolicyTypeIssueInputFile:
		c := config.GetIssueInputFilePolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c
	case consts.OssPolicyTypeProjectResource:
		c := config.GetProjectResourcePolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c

		callbackBody := fmt.Sprintf("{\"userId\":%d,\"orgId\":%d,\"type\":%d,\"host\":\"%s\",\"path\":\"%s\",\"filename\":\"%s\",\"folderId\":%d,\"projectId\":%d,\"bucket\":\"%s\",\"size\":%s,\"format\":%s,\"object\":%s,\"realName\":%s}", userId, orgId, policyType, oc.EndPoint, c.Dir, fileName, resourceFolderId, projectId, oc.BucketName, "${size}", "${imageInfo.format}", "${object}", "${x:filename}")
		callbackBo := bo.OssCallbackBo{
			CallbackBody:     callbackBody,
			CallbackUrl:      policyConfig.CallbackUrl,
			CallbackBodyType: consts.OssCallbackBodyTypeApplicationJson,
		}

		callbackJson = json.ToJsonIgnoreError(callbackBo)
		callbackJson = base64.StdEncoding.EncodeToString([]byte(callbackJson))
	case consts.OssPolicyTypeCompatTest:
		c := config.GetCompatTestPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName: orgId,
			consts.CacheKeyYearConstName:  year,
			consts.CacheKeyMonthConstName: int(month),
			consts.CacheKeyDayConstName:   day,
		})
		policyConfig = &c
	case consts.OssPolicyTypeUserAvatar:
		c := config.GetUserAvatarPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c
	case consts.OssPolicyTypeFeedback:
		c := config.GetFeedbackPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName: orgId,
			consts.CacheKeyYearConstName:  year,
			consts.CacheKeyMonthConstName: int(month),
			consts.CacheKeyDayConstName:   day,
		})
		policyConfig = &c
	case consts.OssPolicyTypeIssueRemark:
		c := config.GetFeedbackPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyIssueIdConstName:   issueId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c
	case consts.OssPolicyTypeImportMembers:
		c := config.GetFeedbackPolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName: orgId,
			consts.CacheKeyYearConstName:  year,
			consts.CacheKeyMonthConstName: int(month),
			consts.CacheKeyDayConstName:   day,
		})
		policyConfig = &c
	case consts.OssPolicyTypeCommentAttachments:
		c := config.GetIssueResourcePolicyConfig()
		c.Dir, _ = util.ParseCacheKey(c.Dir, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
			consts.CacheKeyIssueIdConstName:   issueId,
			consts.CacheKeyYearConstName:      year,
			consts.CacheKeyMonthConstName:     int(month),
			consts.CacheKeyDayConstName:       day,
		})
		policyConfig = &c
		callbackBody := fmt.Sprintf("{\"userId\":%d,\"orgId\":%d,\"type\":%d,\"projectId\":%d,\"host\":\"%s\",\"path\":\"%s\",\"filename\":\"%s\",\"issueId\":%d,\"bucket\":\"%s\",\"size\":%s,\"format\":%s,\"object\":%s,\"realName\":%s}",
			userId, orgId, policyType, projectId, oc.EndPoint, c.Dir, fileName, issueId, oc.BucketName, "${size}", "${imageInfo.format}", "${object}", "${x:filename}")
		callbackBo := bo.OssCallbackBo{
			CallbackBody:     callbackBody,
			CallbackUrl:      policyConfig.CallbackUrl,
			CallbackBodyType: consts.OssCallbackBodyTypeApplicationJson,
		}

		callbackJson = json.ToJsonIgnoreError(callbackBo)
		callbackJson = base64.StdEncoding.EncodeToString([]byte(callbackJson))
	}

	if policyConfig == nil {
		return nil, errs.BuildSystemErrorInfo(errs.OssPolicyTypeError)
	}

	respVo := orgfacade.GetOrgConfig(orgvo.GetOrgConfigReq{OrgId: orgId})
	if respVo.Failure() {
		log.Error(respVo.Message)
	} else {
		// 由前端统一拦截
		// if businees.CheckIsPaidVer(respVo.NewData.PayLevel) {
		// 	policyConfig.MaxFileSize = consts.MaxFileSizeStandard
		// }
	}

	// 由前端统一拦截
	// 任务评论单张附件限制10M
	// if policyType == consts.OssPolicyTypeCommentAttachments {
	// 	policyConfig.MaxFileSize = consts.MaxIssueAttachmentSize
	// }

	//policyBo := oss.PostPolicyWithCallback(policyConfig.Dir, policyConfig.Expire, policyConfig.MaxFileSize, callbackJson)
	// 暂时先把文件最大值改为1G，不去拿配置上指定的最大文件大小。付费版本拦截需求...
	// 根据付费版本确定单文件限制。如果没有，则默认限制 1G
	limitNum, err := domain.GetPayFunctionLimitResourceNum(orgId, consts.FunctionAttachSizeLimit)
	if err != nil {
		log.Errorf("[GetOssPostPolicy] err: %v", err)
	}
	if limitNum <= 0 {
		limitNum = consts.MaxFileSizeStandard
	} else {
		// 转换成字节
		// 存在 functionLevel 中的配置是以 M 为单位
		limitNum *= 1024 * 1024
	}
	policyBo := oss.PostPolicy(policyConfig.Dir, policyConfig.Expire, int64(limitNum))

	//oss绑定自定义域名，所以覆盖host
	policyBo.Host = oc.EndPoint

	resp := &vo.OssPostPolicyResp{}
	copyErr := copyer.Copy(policyBo, resp)
	if copyErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	resp.FileName = fileName
	resp.Callback = callbackJson
	resp.MaxFileSize = int64(limitNum)
	log.Infof("[GetOssPostPolicy] resp: %v", json.ToJsonIgnoreError(resp))
	return resp, nil
}
