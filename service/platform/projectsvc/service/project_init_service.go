package service

import (
	"strconv"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/vo/appvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/appfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

//func ProjectInit(orgId int64, tx sqlbuilder.Tx) (map[string]interface{}, errs.SystemErrorInfo) {
//	//逻辑保留，暂时不需要传递上下文
//	contextMap := map[string]interface{}{}
//
//	//初始化优先级
//	//err := domain.PriorityInit(orgId, tx)
//	//if err != nil {
//	//	log.Error(err)
//	//	return nil, err
//	//}
//
//	return contextMap, nil
//}

func CreateOrgDirectoryAppsAndViews(reqVo projectvo.CreateOrgDirectoryAppsReq) errs.SystemErrorInfo {
	orgId := reqVo.OrgId
	// 创建左侧目录应用
	projectFolderAppId, issueFolderAppId, err := CreateProjectFolderApp(orgId)
	if err != nil {
		log.Errorf("[CreateOrgDirectoryApps] CreateProjectFolderApp err:%v, orgId:%v", err, orgId)
		return err
	}
	log.Infof("[CreateOrgDirectoryApps] 组织 %d 已创建目录应用：%d, %d", orgId, projectFolderAppId, issueFolderAppId)

	// 创建任务视图
	err = CreateSummaryDefaultViewsForIssue(orgId, issueFolderAppId)
	if err != nil {
		log.Errorf("[CreateOrgDirectoryApps] CreateSummaryDefaultViewsForIssue err:%v, orgId:%v", err, orgId)
		return err
	}

	return nil
}

func CreateProjectFolderApp(orgId int64) (int64, int64, errs.SystemErrorInfo) {
	projectFolderId := int64(0)
	issueFolderId := int64(0)
	options := businees.GetMenuSettings()
	for _, option := range options {
		resp := domain.CreateExternalApp(orgId, 0, option.Name, option.Icon, option.LinkUrl)
		if resp.Failure() {
			log.Errorf("[CreateProjectFolderApp] domain.CreateExternalApp err:%v, orgId:%v", resp.Error(), orgId)
			return 0, 0, resp.Error()
		}
		if option.LinkUrl == "/project" {
			projectFolderId = resp.Data.Id
		}
		if option.LinkUrl == "/task" {
			issueFolderId = resp.Data.Id
		}
	}

	updateResp := orgfacade.UpdateOrgRemarkSetting(orgvo.UpdateOrgRemarkSettingReqVo{
		OrgId:  orgId,
		UserId: 0,
		Input: orgvo.SaveOrgSummaryTableAppIdReqVoData{
			ProjectFolderAppId: projectFolderId,
			IssueFolderAppId:   issueFolderId,
		},
	})
	if updateResp.Failure() {
		log.Errorf("[CreateProjectFolderApp] orgfacade.UpdateOrgRemarkSetting err:%v, orgId:%v", updateResp.Error(), orgId)
		return 0, 0, updateResp.Error()
	}

	return projectFolderId, issueFolderId, nil
}

func CreateSummaryDefaultViewsForIssue(orgId, issueFolderAppId int64) errs.SystemErrorInfo {
	orgInfoResp := orgfacade.GetOrgInfo(orgvo.GetOrgInfoReq{OrgId: orgId})
	if orgInfoResp.Failure() {
		log.Errorf("[CreateSummaryDefaultViewsForIssue]orgfacade.GetOrgInfo err:%v, orgId:%v", orgInfoResp.Error(), orgId)
		return orgInfoResp.Error()
	}
	remarkStr := orgInfoResp.Data.Remark
	remarkObj := orgvo.OrgRemarkConfigType{}
	if err := json.FromJson(remarkStr, &remarkObj); err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
	}

	if err := domain.CreateDefaultViewsForIssue(orgId, remarkObj.OrgSummaryTableAppId); err != nil {
		log.Errorf("[CreateSummaryDefaultViewsForIssue]domain.CreateDefaultViewsForIssue err:%v, orgId:%v",
			orgInfoResp.Error(), orgId)
		return err
	}
	// 查询这些视图，为其增加视图镜像（视图镜像是一种特殊的应用）
	// 极星标品：无需增加默认的任务视图镜像。
	// 移动：增加默认的任务视图镜像。
	if util.IsChinaMobileEnv() {
		err := CreateDefaultViewMirror(orgId, remarkObj.OrgSummaryTableAppId, issueFolderAppId, true)
		if err != nil {
			log.Errorf("[CreateSummaryDefaultViewsForIssue]CreateDefaultViewMirror err:%v, orgId:%v",
				orgInfoResp.Error(), orgId)
			return err
		}
	}

	return nil
}

func CreateDefaultViewMirror(orgId, appId, parentId int64, isSummary bool) errs.SystemErrorInfo {
	viewResp := appfacade.GetAppViewList(&appvo.GetAppViewListReq{
		OrgId: orgId,
		AppId: appId,
	})
	if viewResp.Failure() {
		log.Error(viewResp.Error())
		return viewResp.Error()
	}
	views := viewResp.Data
	appType := consts.LcAppTypeForViewMirror
	userId := int64(0)
	for _, view := range views {
		viewMirrorName := view.ViewName
		viewId, _ := strconv.ParseInt(view.ID, 10, 64)
		appIdInt64, _ := strconv.ParseInt(view.AppID, 10, 64)
		req := &permissionvo.CreateLessCodeAppReq{
			OrgId:        &orgId,
			AppType:      &appType,
			Name:         &viewMirrorName,
			UserId:       &userId,
			PkgId:        0,
			ParentId:     parentId,
			Config:       "",
			Hidden:       0,
			MirrorViewId: viewId,
			MirrorAppId:  appIdInt64,
			AddAllMember: true,
		}
		if isSummary {
			req.ProjectId = -1
		}
		createAppResp := appfacade.CreateLessCodeApp(req)
		if createAppResp.Failure() {
			log.Error(createAppResp.Error())
			return createAppResp.Error()
		}
	}
	return nil
}
