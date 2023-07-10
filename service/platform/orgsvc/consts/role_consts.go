package consts

import (
	"strings"

	"github.com/star-table/polaris-backend/common/core/consts"
)

type RolePermission struct {
	RolePermissionId int64
	OperationCodes   string
	RoleLangCodes    []string
}

type RolePermissionOperationDefineInfo struct {
	RoleLangCode string
}

//不允许编辑的默认权限项
var DefaultRoleCodeMap = map[string]bool{
	consts.RoleGroupSpecialMember: true, //组织成员
	consts.RoleGroupOrgAdmin:      true, //组织超管
	consts.RoleGroupOrgManager:    true, //组织管理
	consts.RoleGroupProMember:     true, //项目成员
	consts.RoleGroupSpecialOwner:  true, //负责人
}

var DefaultMemberRoleCodeMap = map[string]bool{
	consts.RoleGroupSpecialMember: true, //组织成员
}

//权限定义
var RolePermissionOperationDefineMap = map[string]RolePermission{
	//测试app
	"/Org/{org_id}/Pro/0/Test/TestApp": {
		RolePermissionId: 35,
		OperationCodes:   "(View)|(Modify)|(Delete)|(Create)",
		RoleLangCodes: []string{
			"RoleGroup.Special.Worker",
			"RoleGroup.Special.Attention",
			"RoleGroup.Special.Member",
			"RoleGroup.Org.Manager",
			"RoleGroup.Pro.ProjectManager",
			"RoleGroup.Pro.TechnicalManager",
			"RoleGroup.Pro.ProductManager",
			"RoleGroup.Pro.Developer",
			"RoleGroup.Pro.Tester",
			"RoleGroup.Pro.Member",
		},
	},
	//测试设备
	"/Org/{org_id}/Pro/0/Test/TestDevice": {
		RolePermissionId: 36,
		OperationCodes:   "(View)|(Modify)|(Delete)|(Create)",
		RoleLangCodes: []string{
			"RoleGroup.Special.Worker",
			"RoleGroup.Special.Attention",
			"RoleGroup.Special.Member",
			"RoleGroup.Org.Manager",
			"RoleGroup.Pro.ProjectManager",
			"RoleGroup.Pro.TechnicalManager",
			"RoleGroup.Pro.ProductManager",
			"RoleGroup.Pro.Developer",
			"RoleGroup.Pro.Tester",
			"RoleGroup.Pro.Member",
		},
	},
	"/Org/{org_id}/Pro/0/Test/TestReport": {
		RolePermissionId: 37,
		OperationCodes:   "(View)|(Modify)|(Delete)|(Create)",
		RoleLangCodes: []string{
			"RoleGroup.Special.Worker",
			"RoleGroup.Special.Attention",
			"RoleGroup.Special.Member",
			"RoleGroup.Org.Manager",
			"RoleGroup.Pro.ProjectManager",
			"RoleGroup.Pro.TechnicalManager",
			"RoleGroup.Pro.ProductManager",
			"RoleGroup.Pro.Developer",
			"RoleGroup.Pro.Tester",
			"RoleGroup.Pro.Member",
		},
	},
}

func IsDefaultRole(roleLangCode string) bool {
	_, ok := DefaultRoleCodeMap[strings.TrimSpace(roleLangCode)]
	return ok
}
