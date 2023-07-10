package appfacade

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/star-table/polaris-backend/facade"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/appvo"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
)

// CreateLessCodeApp 创建管项目对应的应用
func CreateLessCodeApp(req *permissionvo.CreateLessCodeAppReq) *permissionvo.CreateLessCodeAppResp {
	respVo := &permissionvo.CreateLessCodeAppResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps", config.GetPreUrl(consts.ServiceNameLcApp))
	fullUrl := reqUrl + fmt.Sprintf("?appType=%d&name=%s&orgId=%d&userId=%d", *req.AppType, url.QueryEscape(*req.Name), *req.OrgId, *req.UserId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, fullUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// UpdateLessCodeApp 更新项目对应的应用
func UpdateLessCodeApp(req *appvo.UpdateLessCodeAppReq) *formvo.LessCommonIssueResp {
	respVo := &formvo.LessCommonIssueResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps", config.GetPreUrl(consts.ServiceNameLcApp))
	fullUrl := reqUrl + fmt.Sprintf("?name=%s&orgId=%d&userId=%d", url.QueryEscape(req.Name), req.OrgId, req.UserId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPut, fullUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// DeleteLessCodeApp 删除项目对应的应用
func DeleteLessCodeApp(req *appvo.DeleteLessCodeAppReq) *formvo.LessCommonIssueResp {
	respVo := &formvo.LessCommonIssueResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps", config.GetPreUrl(consts.ServiceNameLcApp))
	fullUrl := reqUrl + fmt.Sprintf("?orgId=%d&userId=%d", req.OrgId, req.UserId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodDelete, fullUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetAppInfoByAppId 获取应用信息
func GetAppInfoByAppId(req *appvo.GetAppInfoByAppIdReq) *appvo.GetAppInfoByAppIdResp {
	respVo := &appvo.GetAppInfoByAppIdResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/get-app-info", config.GetPreUrl(consts.ServiceNameLcApp))
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// StarApp 收藏应用
func StarApp(req appvo.CancelStarAppReq) *formvo.LessCommonIssueResp {
	respVo := &formvo.LessCommonIssueResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/%d/stars", config.GetPreUrl(consts.ServiceNameLcApp), req.AppId)
	fullUrl := reqUrl + fmt.Sprintf("?orgId=%d&userId=%d", req.OrgId, req.UserId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, fullUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// CancelStarApp 取消应用收藏
func CancelStarApp(req appvo.CancelStarAppReq) *formvo.LessCommonIssueResp {
	respVo := &formvo.LessCommonIssueResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/%d/stars", config.GetPreUrl(consts.ServiceNameLcApp), req.AppId)
	fullUrl := reqUrl + fmt.Sprintf("?orgId=%d&userId=%d", req.OrgId, req.UserId)
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodDelete, fullUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// CreateAppViews 批量创建视图
func CreateAppViews(req appvo.CreateAppViewReq) *appvo.CreateAppViewResp {
	respVo := &appvo.CreateAppViewResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/views/batch", config.GetPreUrl(consts.ServiceNameLcApp))
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req.Reqs, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetAppViewList 获取应用视图列表
func GetAppViewList(req *appvo.GetAppViewListReq) *appvo.GetAppViewListResp {
	respVo := &appvo.GetAppViewListResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/views", config.GetPreUrl(consts.ServiceNameLcApp))
	queryParams := map[string]interface{}{
		"appId":         req.AppId,
		"includePublic": true,
		"orgId":         req.OrgId,
	}
	if req.UserId > 0 {
		queryParams["owner"] = req.UserId
	}
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// AddAppMembers 增加项目参与人 todo 改方法暂未实现！
func AddAppMembers(req appvo.AddAppMembersReq) *appvo.AddAppMembersResp {
	respVo := &appvo.AddAppMembersResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/add-project-relation", config.GetPreUrl(consts.ServiceNameLcApp))
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, req.Input, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// GetAppInfoList 获取应用列表
func GetAppInfoList(req appvo.GetAppInfoListReq) *appvo.GetAppInfoListResp {
	respVo := &appvo.GetAppInfoListResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/get-app-info-list", config.GetPreUrl(consts.ServiceNameLcApp))
	queryParams := map[string]interface{}{
		"appIds": strings.Replace(strings.Trim(fmt.Sprint(req.AppIds), "[]"), " ", ",", -1),
		"orgId":  req.OrgId,
	}
	err := facade.Request(consts.HttpMethodGet, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// AddAppRelation 增加应用关联关系，可用于增加应用的可见性（例如：对组织所有成员可见）
// 基于 app 服务：`/app/inner/api/v1/apps/add-app-relation`
func AddAppRelation(req appvo.AddAppRelationReq) *appvo.AddAppRelationResp {
	respVo := &appvo.AddAppRelationResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/add-app-relation", config.GetPreUrl(consts.ServiceNameLcApp))
	queryParams := map[string]interface{}{}
	err := facade.Request(consts.HttpMethodPut, reqUrl, queryParams, nil, req, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// IsAppMember 判断是否是应用成员
func IsAppMember(req appvo.IsAppMemberReq) *formvo.LessCommonIssueResp {
	respVo := &formvo.LessCommonIssueResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/is-project-member", config.GetPreUrl(consts.ServiceNameLcApp))
	queryParams := map[string]interface{}{
		"appId":  req.AppId,
		"orgId":  req.OrgId,
		"userId": req.UserId,
	}
	err := facade.Request(consts.HttpMethodPut, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}

// ApplyTemplate 应用模板
func ApplyTemplate(req appvo.ApplyTemplateReq) *appvo.ApplyTemplateResp {
	respVo := &appvo.ApplyTemplateResp{}
	reqUrl := fmt.Sprintf("%s/app/inner/api/v1/apps/apply-template", config.GetPreUrl(consts.ServiceNameLcApp))
	queryParams := map[string]interface{}{
		"orgId":         req.OrgId,
		"userId":        req.UserId,
		"templateId":    req.TemplateId,
		"isNewbieGuide": req.IsNewbieGuide,
	}
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return respVo
}
