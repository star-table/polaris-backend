package domain

import (
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func AddOrgTrends(orgTrendsBo bo.OrgTrendsBo) {
	trendsType := orgTrendsBo.PushType

	var trendsBos []bo.TrendsBo = nil
	var err errs.SystemErrorInfo = nil

	if trendsType == consts.PushTypeApplyJoinOrg {
		trendsBos, err = assemblyApplyJoinOrgTrends(orgTrendsBo)
	} else if trendsType == consts.PushTypeApplicationApproved {
		trendsBos = assemblyApplicationApprovedTrends(orgTrendsBo)
	} else if trendsType == consts.PushTypePromotionToOrgManager {
		trendsBos = assemblyPromotionToOrgManagerTrends(orgTrendsBo)
	}

	if err != nil {
		log.Error(err)
	} else if trendsBos != nil && len(trendsBos) > 0 {
		//插入操作
		for i, _ := range trendsBos {
			trendsBos[i].CreateTime = types.Time(orgTrendsBo.OperateTime)
		}
		err := CreateTrendsBatch(trendsBos)
		if err != nil {
			log.Error(err)
		}
	}

}

//申请加入组织
func assemblyApplyJoinOrgTrends(orgTrendsBo bo.OrgTrendsBo) ([]bo.TrendsBo, errs.SystemErrorInfo) {
	orgId := orgTrendsBo.OrgId
	operatorId := orgTrendsBo.OperatorId
	targetMemberIds := orgTrendsBo.TargetMembers

	if len(targetMemberIds) == 0 {
		return nil, nil
	}

	orgInfo, err := orgfacade.GetBaseOrgInfoRelaxed(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	orgName := orgInfo.OrgName

	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, targetMemberIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	trends := make([]bo.TrendsBo, 0)

	if len(userInfos) > 0 {
		//ext := bo.OrgTrendsBoExt{}
		ext := bo.TrendExtensionBo{}

		simpleUserInfos := make([]bo.SimpleUserInfoBo, len(userInfos))
		for i, userInfo := range userInfos {
			simpleUserInfos[i] = bo.SimpleUserInfoBo{
				Id:     userInfo.UserId,
				Name:   userInfo.Name,
				Avatar: userInfo.Avatar,
			}
		}

		ext.MemberInfo = simpleUserInfos
		ext.ObjName = orgName

		newValue := json.ToJsonIgnoreError(targetMemberIds)

		trendsBo := bo.TrendsBo{
			OrgId:           orgId,
			Module1:         consts.TrendsModuleOrg,
			OperCode:        consts.RoleOperationCheck,
			OperObjId:       orgId,
			OperObjType:     consts.TrendsOperObjectTypeOrg,
			OperObjProperty: consts.BlankString,
			RelationObjId:   0,
			RelationObjType: consts.TrendsOperObjectTypeUser,
			RelationType:    consts.TrendsRelationTypeApplyJoinOrg,
			NewValue:        &newValue,
			OldValue:        nil,
			Ext:             json.ToJsonIgnoreError(ext),
			Creator:         operatorId,
			CreateTime:      types.NowTime(),
		}
		trends = append(trends, trendsBo)
	}
	return nil, nil
}

func assemblyApplicationApprovedTrends(orgTrendsBo bo.OrgTrendsBo) []bo.TrendsBo {
	return nil
}

func assemblyPromotionToOrgManagerTrends(orgTrendsBo bo.OrgTrendsBo) []bo.TrendsBo {
	return nil
}
