package service

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	selfConsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	domain "github.com/star-table/polaris-backend/service/platform/orgsvc/domain/roledomain"
)

var log = logger.GetDefaultLogger()

func RoleUserRelation(orgId, userId, roleId int64) errs.SystemErrorInfo {
	return domain.AddRoleUserRelation(orgId, userId, roleId)
}

//移除角色和用户关联
func RemoveRoleUserRelation(req rolevo.RemoveRoleUserRelationReqVo) errs.SystemErrorInfo {
	return domain.RemoveRoleUserRelation(req.OrgId, req.UserIds, req.OperatorId)
}

//移除角色和部门关联
func RemoveRoleDepartmentRelation(req rolevo.RemoveRoleDepartmentRelationReqVo) errs.SystemErrorInfo {
	return domain.RemoveRoleDepartmentRelation(req.OrgId, req.DeptIds, req.OperatorId)
}

//创建角色
func CreateRole(orgId, userId int64, input vo.CreateRoleReq) (*vo.Void, errs.SystemErrorInfo) {
	//nameErr := domain.JudgeRoleName(input.Name)
	//if nameErr != nil {
	//	return nil, nameErr
	//}
	isNameRight := format.VerifyRoleNameFormat(input.Name)
	if !isNameRight {
		return nil, errs.RoleNameLenErr
	}
	if ok, _ := slice.Contain([]int{1, 2}, input.RoleGroupType); !ok {
		return nil, errs.OrgRoleGroupNotExist
	}

	if input.RoleGroupType == 1 {
		//判断组织权限
		authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgRole, consts.RoleOperationCreate)
		if authErr != nil {
			log.Error(authErr)
			return nil, authErr
		}
	} else {
		//判断项目权限
		if input.ProjectID != nil {
			authResp := projectfacade.AuthProjectPermission(projectvo.AuthProjectPermissionReqVo{
				Input: projectvo.AuthProjectPermissionReqData{
					OrgId:      orgId,
					UserId:     userId,
					ProjectId:  *input.ProjectID,
					Path:       consts.RoleOperationPathOrgProRole,
					Operation:  consts.RoleOperationCreate,
					AuthFiling: true,
				},
			})
			if authResp.Failure() {
				log.Error(authResp.Message)
				return nil, authResp.Error()
			}
		}
	}

	//获取角色组信息
	roleGroupInfo, err := domain.GetGroupRoleList(orgId)
	if err != nil {
		return nil, err
	}

	var groupLangCode string
	if input.RoleGroupType == 1 {
		groupLangCode = consts.RoleGroupOrg
	} else {
		groupLangCode = consts.RoleGroupPro
	}
	var groupId int64
	for _, v := range *roleGroupInfo {
		if v.LangCode == groupLangCode {
			groupId = v.Id
			break
		}
	}

	if input.RoleGroupType == 2 && input.ProjectID != nil {
		//判断项目是否存在
		projectInfo := projectfacade.ProjectInfo(projectvo.ProjectInfoReqVo{
			Input: vo.ProjectInfoReq{
				ProjectID: *input.ProjectID,
			},
			UserId:userId,
			OrgId: orgId,
		})
		if projectInfo.Failure() {
			log.Error(projectInfo.Error())
			return nil, projectInfo.Error()
		}
	}
	//插入角色
	id, err := domain.CreateRole(orgId, userId, input, groupId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//复制成员角色的权限
	var memberRoleLangCode string
	if input.RoleGroupType == 1 {
		//组织成员
		memberRoleLangCode = consts.RoleGroupSpecialMember
	} else {
		memberRoleLangCode = consts.RoleGroupProMember
	}
	memberRole, err := GetRoleByLangCode(orgId, memberRoleLangCode)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	permissionList, err := GetRolePermissionOperationList(orgId, memberRole.Id, 0)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	copyErr := domain.CopyPermission(permissionList, id, orgId)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, copyErr
	}

	return &vo.Void{
		ID: id,
	}, nil
}

func UpdateRole(orgId int64, userId int64, input vo.UpdateRoleReq) (int64, errs.SystemErrorInfo) {
	roleId := input.RoleID

	role, err := domain.GetRole(orgId, 0, roleId)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	if selfConsts.IsDefaultRole(role.LangCode) {
		log.Error("默认角色不允许编辑")
		return 0, errs.DefaultRoleCantModify
	}

	if role.ProjectId == 0 {
		//判断组织权限
		authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgRole, consts.RoleOperationModify)
		if authErr != nil {
			log.Error(authErr)
			return 0, authErr
		}
	} else {
		//判断项目权限
		authResp := projectfacade.AuthProjectPermission(projectvo.AuthProjectPermissionReqVo{
			Input: projectvo.AuthProjectPermissionReqData{
				OrgId:      orgId,
				UserId:     userId,
				ProjectId:  role.ProjectId,
				Path:       consts.RoleOperationPathOrgProRole,
				Operation:  consts.RoleOperationModify,
				AuthFiling: true,
			},
		})
		if authResp.Failure() {
			log.Error(authResp.Message)
			return 0, authResp.Error()
		}
	}

	updateErr := domain.UpdateRole(orgId, role.ProjectId, roleId, bo.UpdateRoleBo{
		Name: input.Name,
	}, userId)
	if updateErr != nil {
		log.Error(updateErr)
		return 0, updateErr
	}

	return roleId, nil
}

func DelRole(orgId int64, userId int64, input vo.DelRoleReq) errs.SystemErrorInfo {
	roleIds := input.RoleIds
	projectId := int64(0)
	if input.ProjectID != nil {
		projectId = *input.ProjectID
	}

	if projectId == 0 {
		//判断组织权限
		authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgRole, consts.RoleOperationDelete)
		if authErr != nil {
			log.Error(authErr)
			return authErr
		}
	} else {
		//判断项目权限
		authResp := projectfacade.AuthProjectPermission(projectvo.AuthProjectPermissionReqVo{
			Input: projectvo.AuthProjectPermissionReqData{
				OrgId:      orgId,
				UserId:     userId,
				ProjectId:  projectId,
				Path:       consts.RoleOperationPathOrgProRole,
				Operation:  consts.RoleOperationDelete,
				AuthFiling: true,
			},
		})
		if authResp.Failure() {
			log.Error(authResp.Message)
			return authResp.Error()
		}
	}

	//特殊角色不能删除（所以不用查询出来）
	roleList, err := domain.GetRoleList(orgId, projectId, roleIds, false)
	if err != nil {
		log.Error(err)
		return err
	}

	if len(roleList) == 0 {
		log.Error("要删除的角色列表为空")
		//这里不返回错误好了: nico
		return nil
	}

	beDeletedRoleIds := make([]int64, 0)
	for _, role := range roleList {
		if selfConsts.IsDefaultRole(role.LangCode) {
			continue
		}
		beDeletedRoleIds = append(beDeletedRoleIds, role.Id)
	}

	if len(beDeletedRoleIds) == 0 {
		log.Error("要删除的角色列表为空")
		//这里不返回错误好了: nico
		return nil
	}

	deleteErr := domain.DeleteRoles(orgId, beDeletedRoleIds, userId)
	if deleteErr != nil {
		log.Error(deleteErr)
		return deleteErr
	}

	return nil
}
