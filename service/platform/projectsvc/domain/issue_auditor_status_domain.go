package domain

import (
	"fmt"
	"strings"

	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/datacenter"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/facade/formfacade"
)

func UpdateLcIssueAuditStatusDetailByUser(orgId, appId, userId, issueId int64, auditStatus int) errs.SystemErrorInfo {
	req := &formvo.LessUpdateIssueBatchReq{
		OrgId:  orgId,
		AppId:  appId,
		UserId: userId,
		Condition: vo.LessCondsData{
			Type: consts.ConditionAnd,
			Conds: []*vo.LessCondsData{
				{
					Type:   consts.ConditionEqual,
					Value:  orgId,
					Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
				},
				{
					Type:   consts.ConditionEqual,
					Value:  issueId,
					Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
				},
			},
		},
		Sets: []datacenter.Set{
			{
				Column: consts.LcJsonColumn,
				Value: fmt.Sprintf("%s||jsonb_build_object('%s', jsonb_set(coalesce(data->'%s','{}'), '{%d}', '%d'))",
					consts.LcJsonColumn, consts.BasicFieldAuditStatusDetail, consts.BasicFieldAuditStatusDetail, userId, auditStatus),
				Type:            consts.SetTypeJson,
				Action:          consts.SetActionSet,
				WithoutPretreat: true,
			},
		},
	}
	resp := formfacade.LessUpdateIssueBatchRaw(req)
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}
	return nil
}

func DeleteLcIssueAuditStatusDetailByUsers(orgId, appId, userId, issueId int64, userIds ...int64) errs.SystemErrorInfo {
	var f strings.Builder
	for _, id := range userIds {
		f.WriteString(fmt.Sprintf("-'%v'", id))
	}
	req := &formvo.LessUpdateIssueBatchReq{
		OrgId:  orgId,
		AppId:  appId,
		UserId: userId,
		Condition: vo.LessCondsData{
			Type: consts.ConditionAnd,
			Conds: []*vo.LessCondsData{
				{
					Type:   consts.ConditionEqual,
					Value:  orgId,
					Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
				},
				{
					Type:   consts.ConditionEqual,
					Value:  issueId,
					Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
				},
			},
		},
		Sets: []datacenter.Set{
			{
				Column: consts.LcJsonColumn,
				Value: fmt.Sprintf("%s||jsonb_build_object('%s', coalesce(data->'%s','{}')"+f.String()+")",
					consts.LcJsonColumn, consts.BasicFieldAuditStatusDetail, consts.BasicFieldAuditStatusDetail),
				Type:            consts.SetTypeJson,
				Action:          consts.SetActionSet,
				WithoutPretreat: true,
			},
		},
	}
	resp := formfacade.LessUpdateIssueBatchRaw(req)
	if resp.Failure() {
		log.Error(resp.Error())
		return resp.Error()
	}
	return nil
}
