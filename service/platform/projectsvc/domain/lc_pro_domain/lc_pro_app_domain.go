package lc_pro_domain

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
	"github.com/star-table/polaris-backend/facade/permissionfacade"
)

var log = logger.GetDefaultLogger()

// 更新应用对应的权限组的权限项
// 包含：查看权限-`-2`，编辑权限-`-3`，管理员-`-4`
func UpdateOpForAppPermissionGroup(lcAppId int64) errs.SystemErrorInfo {
	if lcAppId < 1 {
		err := errs.LcAppIdInvalid
		log.Error(err)
		return err
	}
	if err := UpdateOpForOneAppPermissionGroup(lcAppId, consts.AppPermissionGroupLangCodeForRead); err != nil {
		log.Error(err)
		return err
	}
	if err := UpdateOpForOneAppPermissionGroup(lcAppId, consts.AppPermissionGroupLangCodeForWrite); err != nil {
		log.Error(err)
		return err
	}
	if err := UpdateOpForOneAppPermissionGroup(lcAppId, consts.AppPermissionGroupLangCodeForAdmin); err != nil {
		log.Error(err)
		return err
	}
	return nil
}

// 更新应用对应的权限组的 optAuth 值
func UpdateOpForOneAppPermissionGroup(appId int64, groupLangCode int) errs.SystemErrorInfo {
	optAuthStr := "[]"
	switch groupLangCode {
	case consts.AppPermissionGroupLangCodeForRead:
		optAuthStr = consts.LcAppPermissionGroupOpsReadJson
	case consts.AppPermissionGroupLangCodeForWrite:
		optAuthStr = consts.LcAppPermissionGroupOpsWriteJson
	case consts.AppPermissionGroupLangCodeForAdmin:
		optAuthStr = consts.LcAppPermissionGroupOpsAdminJson
	}
	resp := permissionfacade.UpdateLcAppPermissionGroupOpConfig(&permissionvo.UpdateLcAppPermissionGroupOpConfigReq{
		AppId:    appId,
		LangCode: groupLangCode,
		OptAuth:  optAuthStr,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}
	return nil
}
