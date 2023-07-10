package orgfacade

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	"github.com/star-table/polaris-backend/facade/userfacade"
)

func GetUserAdminFlagRelaxed(orgId int64, userId int64) (*bo.UserAdminFlagBo, errs.SystemErrorInfo) {
	resBo := &bo.UserAdminFlagBo{}
	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
	if manageAuthInfoResp.Failure() {
		log.Error(manageAuthInfoResp.Message)
		return resBo, manageAuthInfoResp.Error()
	}
	resBo.IsAdmin = manageAuthInfoResp.Data.IsSysAdmin

	return resBo, nil
}

func AuthenticateRelaxed(orgId int64, userId int64, projectAuthInfo *bo.ProjectAuthBo, issueAuthInfo *bo.IssueAuthBo, path string, operation string, authFields []string) errs.SystemErrorInfo {
	respVo := Authenticate(rolevo.AuthenticateReqVo{
		OrgId:  orgId,
		UserId: userId,
		AuthInfoReqVo: rolevo.AuthenticateAuthInfoReqVo{
			ProjectAuthInfo: projectAuthInfo,
			IssueAuthInfo:   issueAuthInfo,
			Fields:          authFields,
		},
		Path:      path,
		Operation: operation,
	})
	if respVo.Failure() {
		return respVo.Error()
	}

	return nil
}

func RoleUserRelationRelaxed(orgId, userId, roleId int64) errs.SystemErrorInfo {
	respVo := RoleUserRelation(rolevo.RoleUserRelationReqVo{
		OrgId:  orgId,
		UserId: userId,
		RoleId: roleId,
	})
	if respVo.Failure() {
		return respVo.Error()
	}

	return nil
}
