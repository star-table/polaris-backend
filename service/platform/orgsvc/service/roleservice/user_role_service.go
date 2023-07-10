package service

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	lang2 "github.com/star-table/polaris-backend/common/core/util/lang"
	"github.com/star-table/polaris-backend/common/language/english"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	orgdomain "github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	domain "github.com/star-table/polaris-backend/service/platform/orgsvc/domain/roledomain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

func GetOrgRoleUser(orgId int64, projectId int64) ([]rolevo.RoleUser, errs.SystemErrorInfo) {
	roleBo, err := domain.GetOrgRoleUser(orgId, projectId, nil)
	if err != nil {
		return nil, err
	}
	roleVo := &[]rolevo.RoleUser{}
	copyErr := copyer.Copy(roleBo, roleVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return *roleVo, nil
}

func GetOrgAdminUser(orgId int64) ([]int64, errs.SystemErrorInfo) {
	roleBos, err := domain.GetOrgRoleUser(orgId, 0, []string{consts.RoleGroupOrgAdmin, consts.RoleGroupOrgManager})
	if err != nil {
		return nil, err
	}
	userIds := make([]int64, len(*roleBos))
	for i, roleBo := range *roleBos {
		userIds[i] = roleBo.UserId
	}
	userIds = slice.SliceUniqueInt64(userIds)
	return userIds, nil
}

// 批量获取多个组织下的超管+管理员用户id
func GetOrgAdminUserBatch(orgIds []int64) (result map[int64][]int64, busiErr errs.SystemErrorInfo) {
	// orgId => [userId]
	// 初始化一下。接收端需判断是否只有一个元素，并且 key 为 0。
	result = map[int64][]int64{}
	roleBos, err := domain.GetOrgRoleUserBatch(orgIds, 0, []string{consts.RoleGroupOrgAdmin, consts.RoleGroupOrgManager})
	if err != nil {
		busiErr = err
		return
	}
	for _, item := range *roleBos {
		if _, ok := result[item.OrgId]; ok {
			result[item.OrgId] = append(result[item.OrgId], item.UserId)
		} else {
			result[item.OrgId] = []int64{item.UserId}
		}
	}
	return
}

func UpdateUserOrgRole(req rolevo.UpdateUserOrgRoleReqVo) (*vo.Void, errs.SystemErrorInfo) {
	orgId := req.OrgId
	operatorId := req.CurrentUserId
	targetUserId := req.UserId

	//这里根据id查角色
	role, err := domain.GetRole(0, 0, req.RoleId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if role.OrgId == 0 {
		//如果是全局角色，只能修改默认角色(暂时特殊角色只能修改为组织成员)
		if role.LangCode != consts.RoleGroupSpecialMember {
			log.Error("目标特殊角色只能是组织成员")
			return nil, errs.NoOperationPermissions
		}
	} else {
		if role.OrgId != orgId || role.LangCode == consts.RoleGroupOrgAdmin {
			log.Error("不能修改为超级管理员")
			return nil, errs.NoOperationPermissions
		}
	}

	if req.ProjectId == nil || *req.ProjectId == 0 {
		//判断组织权限
		authErr := AuthOrgRole(req.OrgId, req.CurrentUserId, consts.RoleOperationPathOrgUser, consts.RoleOperationBind)
		if authErr != nil {
			log.Error(authErr)
			return nil, authErr
		}

		// 组织角色逻辑判断，当前用户必须角色属性高于被修改角色（目前只有超管可以修改角色的权限）
		targetUserAdminFlag, err := GetUserAdminFlag(orgId, targetUserId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		if targetUserAdminFlag.IsAdmin {
			log.Error("超管的角色不允许修改")
			return nil, errs.OrgAdminRoleCannotModify
		}

		operatorAdminFlag, err := GetUserAdminFlag(orgId, operatorId)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		if targetUserAdminFlag.IsManager && !operatorAdminFlag.IsAdmin {
			log.Error("没有权限修改，因为只有超管才能修改管理员的角色")
			return nil, errs.OrgManagerRoleCannotModify
		}
	} else {
		if role.ProjectId != 0 && role.ProjectId != *req.ProjectId {
			log.Error("角色不属于该项目")
			return nil, errs.NoOperationPermissions
		}
		//判断项目权限
		authResp := projectfacade.AuthProjectPermission(projectvo.AuthProjectPermissionReqVo{
			Input: projectvo.AuthProjectPermissionReqData{
				OrgId:      req.OrgId,
				UserId:     req.CurrentUserId,
				ProjectId:  *req.ProjectId,
				Path:       consts.RoleOperationPathOrgProRole,
				Operation:  consts.RoleOperationBind,
				AuthFiling: true,
			},
		})
		if authResp.Failure() {
			log.Error(authResp.Message)
			return nil, authResp.Error()
		}
	}

	id, updErr := domain.UpdateUserOrgRole(role, orgId, req.CurrentUserId, []int64{req.UserId}, req.ProjectId)
	if updErr != nil {
		log.Error(updErr)
		return nil, updErr
	}

	return &vo.Void{
		ID: id,
	}, nil
}

func UpdateOrgAdmin(userId int64, orgId int64, oldOwnerId, newOwnerId int64) (*vo.Void, errs.SystemErrorInfo) {
	//将新人变成超级管理员
	adminRole, err := GetRoleByLangCode(orgId, consts.RoleGroupOrgAdmin)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	_, err = domain.UpdateUserOrgRole(adminRole, orgId, userId, []int64{newOwnerId}, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//将旧超管变成组织成员
	orgMemberRole, err := GetRoleByLangCode(0, consts.RoleGroupSpecialMember)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	_, err = domain.UpdateUserOrgRole(orgMemberRole, orgId, userId, []int64{oldOwnerId}, nil)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.Void{ID: newOwnerId}, nil
}

func GetUserAdminFlag(orgId, userId int64) (*bo.UserAdminFlagBo, errs.SystemErrorInfo) {
	//userRole, err := GetUserRoleList(orgId, userId, 0)
	//
	//if err != nil {
	//	log.Error(err)
	//	return nil, err
	//}
	//adminRole, err := GetRoleByLangCode(orgId, consts.RoleGroupOrgAdmin)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
	//}
	//managerRole, err := GetRoleByLangCode(orgId, consts.RoleGroupOrgManager)
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
	//}
	//isAdmin := false
	//isManager := false
	//for _, v := range *userRole {
	//	if adminRole != nil && v.RoleId == adminRole.Id {
	//		isAdmin = true
	//	}
	//	if managerRole != nil && v.RoleId == managerRole.Id {
	//		isManager = true
	//	}
	//}
	res := &bo.UserAdminFlagBo{
		IsAdmin:         false,
		IsManager:       false,
		IsPlatformAdmin: false,
	}
	// 融合极星
	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
	if manageAuthInfoResp.Failure() {
		log.Error(manageAuthInfoResp.Message)
		return res, manageAuthInfoResp.Error()
	}
	manageAuthInfo := manageAuthInfoResp.Data
	res.IsAdmin = manageAuthInfo.IsSysAdmin
	res.IsManager = manageAuthInfo.IsSubAdmin

	return res, nil
}

//获取组织角色列表
func GetOrgRoleList(orgId int64) ([]*vo.Role, errs.SystemErrorInfo) {
	groupRole, err := GetRoleListByGroup(orgId, consts.RoleGroupOrg, 0)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	orgMember, err := GetRoleByLangCode(0, consts.RoleGroupSpecialMember)
	if err != nil {
		return nil, err
	}
	if orgMember != nil {
		//替换下名称（奇怪吧）
		remark := orgMember.Remark
		name := orgMember.Name
		orgMember.Name = remark
		orgMember.Remark = name

		if len(groupRole) > 2 {
			//把角色按照 超管/管理员/成员的排序
			groupRole = append(groupRole[:2], append([]bo.RoleBo{*orgMember}, groupRole[2:]...)...)
		} else {
			groupRole = append(groupRole, *orgMember)
		}
	}

	resVo := &[]*vo.Role{}
	copyErr := copyer.Copy(groupRole, resVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	lang := lang2.GetLang()
	isOtherLang := lang2.IsEnglish()
	if isOtherLang {
		if tmpMap, ok := consts.LANG_ROLE_NAME_MAP[lang]; ok {
			for index, item := range *resVo {
				if tmpVal1, ok2 := tmpMap[item.Name]; ok2 {
					(*resVo)[index].Name = tmpVal1
				}
			}
		}
	}

	return *resVo, nil
}

//获取项目所有角色
func GetProjectRoleList(orgId int64, projectId int64) ([]*vo.Role, errs.SystemErrorInfo) {
	groupRole, err := GetRoleListByGroup(orgId, consts.RoleGroupPro, projectId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//负责人
	ownerRole, err := GetRoleByLangCode(0, consts.RoleGroupSpecialOwner)
	if err != nil {
		return nil, err
	}
	newGroup := []bo.RoleBo{}
	if ownerRole != nil {
		//替换下名称（奇怪吧）
		remark := ownerRole.Remark
		name := ownerRole.Name
		ownerRole.Name = remark
		ownerRole.Remark = name
		newGroup = append(newGroup, *ownerRole)
	}

	for _, roleBo := range groupRole {
		//该项目的角色和项目成员
		if roleBo.ProjectId == projectId || roleBo.LangCode == consts.RoleGroupProMember {
			newGroup = append(newGroup, roleBo)
		}
	}

	resVo := &[]*vo.Role{}
	copyErr := copyer.Copy(newGroup, resVo)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	// 多语言适配
	if lang2.IsEnglish() {
		for index, item := range *resVo {
			if tmpVal, ok2 := english.ProRolesLang[item.Name]; ok2 {
				(*resVo)[index].Name = tmpVal
			}
		}
	}

	return *resVo, nil
}

func UpdateUserOrgRoleBatch(req rolevo.UpdateUserOrgRoleBatchReqVo) (*vo.Void, errs.SystemErrorInfo) {
	orgId := req.OrgId
	operatorId := req.CurrentUserId
	targetUserIds := req.Input.UserIds
	roleId := req.Input.RoleId
	if len(targetUserIds) == 0 {
		return &vo.Void{ID: 0}, nil
	}

	_, err := orgdomain.GetBaseUserInfoBatch(orgId, targetUserIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//这里根据id查角色
	role, err := domain.GetRole(0, 0, roleId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if role.OrgId == 0 {
		//如果是全局角色，只能修改默认角色(暂时特殊角色只能修改为组织成员)
		if role.LangCode != consts.RoleGroupSpecialMember {
			log.Error("目标特殊角色只能是组织成员")
			return nil, errs.NoOperationPermissions
		}
	} else {
		if role.OrgId != orgId || role.LangCode == consts.RoleGroupOrgAdmin {
			log.Error("不能修改为超级管理员")
			return nil, errs.NoOperationPermissions
		}
	}

	//判断组织权限
	authErr := AuthOrgRole(req.OrgId, req.CurrentUserId, consts.RoleOperationPathOrgUser, consts.RoleOperationBind)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	// 组织角色逻辑判断，当前用户必须角色属性高于被修改角色（目前只有超管可以修改角色的权限
	adminRole, err := GetRoleByLangCode(orgId, consts.RoleGroupOrgAdmin)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}
	adminUsers, err := GetRoleUserIds(orgId, adminRole.Id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	managerRole, err := GetRoleByLangCode(orgId, consts.RoleGroupOrgManager)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
	}
	managerUsers, err := GetRoleUserIds(orgId, managerRole.Id)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	for _, user := range adminUsers {
		if ok, _ := slice.Contain(targetUserIds, user); ok {
			log.Error("超管的角色不允许修改")
			return nil, errs.OrgAdminRoleCannotModify
		}
	}

	operatorAdminFlag, err := GetUserAdminFlag(orgId, operatorId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if !operatorAdminFlag.IsAdmin {
		for _, user := range managerUsers {
			if ok, _ := slice.Contain(targetUserIds, user); ok {
				log.Error("没有权限修改，因为只有超管才能修改管理员的角色")
				return nil, errs.OrgManagerRoleCannotModify
			}
		}
	}

	id, updErr := domain.UpdateUserOrgRole(role, orgId, req.CurrentUserId, targetUserIds, nil)
	if updErr != nil {
		log.Error(updErr)
		return nil, updErr
	}

	return &vo.Void{
		ID: id,
	}, nil
}

func GetRoleUserIds(orgId, roleId int64) ([]int64, errs.SystemErrorInfo) {
	pos := &[]po.PpmRolRoleUser{}
	infoErr := mysql.SelectAllByCond(consts.TableRoleUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcRoleId:   roleId,
	}, pos)
	if infoErr != nil {
		log.Error(infoErr)
		return nil, errs.MysqlOperateError
	}

	res := []int64{}
	for _, user := range *pos {
		res = append(res, user.UserId)
	}

	return slice.SliceUniqueInt64(res), nil
}
