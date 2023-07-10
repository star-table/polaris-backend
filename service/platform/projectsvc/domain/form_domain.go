package domain

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func GetAppIdFromProjectId(orgId, projectId int64) (int64, errs.SystemErrorInfo) {
	if projectId <= 0 {
		//如果未分配项目就分配到默认应用里面去
		orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
			OrgId:  orgId,
			UserId: 0,
		})
		if orgInfoResp.Failure() {
			log.Error(orgInfoResp.Error())
			return 0, orgInfoResp.Error()
		}
		orgRemarkObj := &orgvo.OrgRemarkConfigType{}
		oriErr := json.FromJson(orgInfoResp.OrganizationInfo.Remark, orgRemarkObj)
		if oriErr != nil {
			log.Error(oriErr)
			return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
		}
		if projectId == 0 {
			return orgRemarkObj.EmptyProjectAppId, nil
		} else {
			return orgRemarkObj.OrgSummaryTableAppId, nil
		}
	}
	info, err := GetProjectSimple(orgId, projectId)
	if err != nil {
		log.Errorf("[GetAppIdFromProjectId] GetProjectSimple err: %v", err)
		return 0, err
	}

	return info.AppId, nil
}

func GetProjectFormAppIdBatch(orgIds []int64) (map[int64]int64, errs.SystemErrorInfo) {
	resp := orgfacade.GetOrgBoListByPage(orgvo.GetOrgIdListByPageReqVo{
		Page: 1,
		Size: 10000,
		Input: orgvo.GetOrgIdListByPageReqVoData{
			OrgIds: orgIds,
		},
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}
	proFormAppIdMap := make(map[int64]int64, 0)
	for _, org := range resp.Data.List {
		orgRemarkJson := org.Remark
		orgRemarkObj := &orgvo.OrgRemarkConfigType{}
		if len(orgRemarkJson) > 0 {
			oriErr := json.FromJson(orgRemarkJson, orgRemarkObj)
			if oriErr != nil {
				log.Errorf("[项目视图] 组织 remark 反序列化异常，组织id:%d,原因:%s", org.Id, oriErr)
				return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
			}
		}
		proFormAppIdMap[org.Id] = orgRemarkObj.ProjectFormAppId
	}

	return proFormAppIdMap, nil
}

// 获取保存项目的表单的 appId
func GetProjectFormAppId(orgId int64) (int64, errs.SystemErrorInfo) {
	//如果未分配项目就分配到默认应用里面去
	orgInfoResp := orgfacade.OrganizationInfo(orgvo.OrganizationInfoReqVo{
		OrgId:  orgId,
		UserId: 0,
	})
	if orgInfoResp.Failure() {
		log.Error(orgInfoResp.Error())
		return 0, orgInfoResp.Error()
	}
	orgRemarkObj := &orgvo.OrgRemarkConfigType{}
	oriErr := json.FromJson(orgInfoResp.OrganizationInfo.Remark, orgRemarkObj)
	if oriErr != nil {
		log.Error(oriErr)
		return 0, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
	}

	return orgRemarkObj.ProjectFormAppId, nil
}
