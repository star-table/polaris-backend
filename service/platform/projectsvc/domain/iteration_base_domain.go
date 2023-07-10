package domain

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/status"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func GetIterationInfoBo(iterationBo bo.IterationBo) (*bo.IterationInfoBo, errs.SystemErrorInfo) {
	orgId := iterationBo.OrgId

	statusInfo := consts.GetIterationStatusById(iterationBo.Status)
	if statusInfo == nil {
		statusInfo = &status.StatusInfoBo{}
	}

	projectInfo, err := GetHomeProjectInfoBo(orgId, iterationBo.ProjectId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectNotExist, err)
	}

	//负责人信息
	ownerId := iterationBo.Owner
	ownerBaseInfo, err := orgfacade.GetDingTalkBaseUserInfoRelaxed(orgId, ownerId)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError, err)
	}
	ownerInfo := &bo.UserIDInfoBo{
		UserID: ownerBaseInfo.UserId,
		Name:   ownerBaseInfo.Name,
		Avatar: ownerBaseInfo.Avatar,
		EmplID: ownerBaseInfo.OutUserId,
	}

	return &bo.IterationInfoBo{
		Iteration:  iterationBo,
		Project:    projectInfo,
		Status:     statusInfo,
		Owner:      ownerInfo,
		NextStatus: nil, // 这个字段没有用上
	}, nil
}
