package tablefacade

import (
	"strings"

	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
	"github.com/star-table/polaris-backend/facade/permissionfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/spf13/cast"
)

func GetDataCollaborateUserIds(orgId, userId, dataId int64) ([]int64, errs.SystemErrorInfo) {
	resp := GetDataCollaborators(orgId, userId, []int64{dataId})
	if resp.Failure() {
		log.Errorf("[GetDataCollaborateUserIds] err: %v, orgId: %d, dataId: %d", resp.Failure(), orgId, dataId)
		return nil, resp.Error()
	}

	var userIds []int64
	for _, c := range resp.Data.Collaborators {
		for _, id := range c.Ids {
			if strings.HasPrefix(id, consts.LcCustomFieldUserType) {
				userIds = append(userIds, cast.ToInt64(id[2:]))
			}
		}
	}
	return userIds, nil
}

func GetAppCollaborators(orgId, userId, appId int64) ([]*uservo.MemberDept, errs.SystemErrorInfo) {
	resp := GetAppCollaboratorRoles(orgId, userId, appId)
	if resp.Failure() {
		log.Errorf("[GetAppCollaborators] GetAppCollaboratorRoles err: %v, orgId: %d, appId: %d", resp.Failure(), orgId, appId)
		return nil, resp.Error()
	}

	log.Infof("[GetAppCollaborators] resp: %v", json.ToJsonIgnoreError(resp))

	var userIds []int64
	userRoles := make(map[int64][]int64)
	for _, c := range resp.Data.CollaboratorRoles {
		if strings.HasPrefix(c.Id, consts.LcCustomFieldUserType) {
			uId := cast.ToInt64(c.Id[2:])
			userIds = append(userIds, uId)
			for _, rId := range c.RoleIds {
				userRoles[uId] = append(userRoles[uId], cast.ToInt64(rId))
			}
		}
	}

	appRolesResp := permissionfacade.GetAppRoleList(orgId, appId)
	if appRolesResp.Failure() {
		log.Errorf("[GetAppCollaborators] GetAppRoleList err: %v, orgId: %d, appId: %d", appRolesResp.Failure(), orgId, appId)
		return nil, appRolesResp.Error()
	}

	var editorRole *permissionvo.AppRoleInfo
	allRoles := make(map[int64]*permissionvo.AppRoleInfo)
	for i, _ := range appRolesResp.Data {
		role := &appRolesResp.Data[i]
		if role.LangCode == consts.GroupLandCodeProjectMember || role.LangCode == consts.GroupLandCodeEdit {
			editorRole = role
		}
		allRoles[role.Id] = role
	}

	usersResp := userfacade.GetAllUserByIds(orgId, userIds)
	if usersResp.Failure() {
		log.Errorf("[GetAppCollaborators] GetAllUserByIds err: %v", resp.Error())
		return nil, resp.Error()
	}

	var result []*uservo.MemberDept
	for _, u := range usersResp.Data {
		us := &uservo.MemberDept{
			Id:       u.Id,
			Name:     u.Name,
			Avatar:   u.Avatar,
			Type:     consts.LcCustomFieldUserType,
			Status:   u.Status,
			IsDelete: u.IsDelete,
		}

		if roleIds, ok := userRoles[u.Id]; ok {
			for _, rId := range roleIds {
				if rId == -1 {
					us.PerGroups = append(us.PerGroups, editorRole)
				} else {
					if role, ok1 := allRoles[rId]; ok1 {
						us.PerGroups = append(us.PerGroups, role)
					}
				}
			}
		}

		result = append(result, us)
	}
	return result, nil
}
