package consts

import "github.com/star-table/common/core/util/json"

const (
	PermissionForPro = `{
    "Attachment":[
        "Upload",
        "Download",
        "Delete"
    ],
    "File":[
        "Upload",
        "Download",
        "Modify",
        "Delete",
        "CreateFolder",
        "ModifyFolder",
        "DeleteFolder"
    ],
    "Issue":[
        "Modify",
        "Bind",
        "Unbind",
        "Create",
        "Delete",
        "ModifyStatus",
        "Comment",
        "Import",
        "Export"
    ],
    "Iteration":[
        "Modify",
        "Create",
        "Delete",
        "ModifyStatus",
        "Bind",
        "Unbind"
    ],
    "Member":[
        "Modify",
        "Bind",
        "Unbind"
    ],
    "ProConfig":[
        "Delete",
        "Modify",
        "Bind",
        "Unbind",
        "Filing",
        "UnFiling",
        "ModifyStatus",
        "ModifyField"
    ],
    "ProjectObjectType":[
        "Modify",
        "Create",
        "Delete"
    ],
    "Role":[
        "Modify",
        "Create",
        "Delete",
        "ModifyPermission"
    ],
    "Tag":[
        "Create",
        "Delete",
        "Remove",
        "Modify"
    ],
    "View":[
        "ManagePrivate",
        "ManagePublic"
    ],
    "MenuPermissionPro":[
        "Iteration",
        "Demand",
        "Issue",
        "Bug",
        "IterationOverview",
        "Plan",
        "File",
        "ProOverview",
        "WorkHour",
        "Statistics",
        "Gantt",
        "Setting",
        "Trash",
        "ProName",
        "Collection",
        "ProStatus",
        "ProMember",
        "GroupChat",
        "MoreOperation"
    ]
}`

	PermissionDefaultOperationForPro = `{
    "Attachment":[
        "Upload",
        "Download",
        "Delete"
    ],
    "File":[
        "Upload",
        "Download",
        "Modify",
        "Delete",
        "CreateFolder",
        "ModifyFolder",
        "DeleteFolder"
    ],
    "Issue":[
        "Modify",
        "Bind",
        "Unbind",
        "Create",
        "Delete",
        "ModifyStatus",
        "Comment"
    ],
    "Iteration":[
        "Modify",
        "Create",
        "Delete",
        "ModifyStatus",
        "Bind",
        "Unbind"
    ],
    "Member":[
        "Modify",
        "Bind",
        "Unbind"
    ],
    "ProConfig":[
        "Delete",
        "Modify",
        "Bind",
        "Unbind",
        "Filing",
        "UnFiling",
        "ModifyStatus",
        "ModifyField"
    ],
    "ProjectObjectType":[
        "Modify",
        "Create",
        "Delete"
    ],
    "Role":[
        "Modify",
        "Create",
        "Delete",
        "ModifyPermission"
    ],
    "Tag":[
        "Create",
        "Delete",
        "Remove",
        "Modify"
    ],
    "View":[
        "ManagePrivate"
    ],
    "MenuPermissionPro":[
        "Iteration",
        "Demand",
        "Issue",
        "Bug",
        "IterationOverview",
        "Plan",
        "File",
        "ProOverview",
        "WorkHour",
        "Statistics",
        "Gantt",
        "Setting",
        "Trash",
        "ProName",
        "Collection",
        "ProStatus",
        "ProMember",
        "GroupChat",
        "MoreOperation"
    ]
}`

	PermissionForNoProIssue = `{
    "Attachment":[
        "Upload",
        "Download",
        "Delete"
    ],
    "File":[
        "Upload",
        "Download",
        "Modify",
        "Delete",
        "CreateFolder",
        "ModifyFolder",
        "DeleteFolder"
    ],
    "Issue":[
        "Modify",
        "Bind",
        "Unbind",
        "Create",
        "Delete",
        "ModifyStatus",
        "Comment"
    ],
    "Tag":[
        "Create",
        "Delete",
        "Remove",
        "Modify"
    ],
    "MenuPermissionPro":[
        "Iteration",
        "Demand",
        "Issue",
        "Bug",
        "IterationOverview",
        "Plan",
        "File",
        "ProOverview",
        "WorkHour",
        "Statistics",
        "Gantt",
        "Setting",
        "Trash",
        "ProName",
        "Collection",
        "ProStatus",
        "ProMember",
        "GroupChat",
        "MoreOperation"
    ]
}`

	PermissionForOrg = `{
    "Department":[
        "Create",
        "Modify",
        "Delete"
    ],
    "MessageConfig":[
        "Modify"
    ],
    "OrgConfig":[
        "Modify",
        "Transfer",
        "ModifyField",
        "TplSaveAs",
        "TplDelete"
    ],
    "Project":[
        "Create",
        "Attention",
        "UnAttention",
        "ModifyField",
        "Manage"
    ],
    "Role":[
        "Create",
        "Modify",
        "Delete"
    ],
    "RoleGroup":[
        "Create",
        "Modify",
        "Delete"
    ],
    "AdminGroup":[
        "View",
        "Create",
        "Modify",
        "Delete"
    ],
    "Team":[
        "Create",
        "Modify",
        "Delete",
        "ModifyStatus",
        "Bind",
        "Unbind"
    ],
    "User":[
        "ModifyStatus",
        "ModifyUserAdminGroup",
        "Bind",
        "Unbind",
        "Watch",
        "ModifyUserDept",
        "ModifyDepartment"
    ],
    "InviteUser":[
        "Invite"
    ],
	"AddUser":[
		"Add"
	],
	"PersonInfo":[
		"Manage"
	],
    "MenuPermissionOrg":[
        "Workspace",
        "Issue",
        "Project",
        "PolarisTpl",
        "Member",
        "WorkHour",
        "Setting",
        "Trash",
		"Trend",
		"CreateButton"
    ]
}`
)

// GetPermissionForOrg 获取管理组织的权限组以及对应的操作项，其中包含了**所有**的权限项
// 这是转换后，提供给前端使用的。
func GetPermissionForOrg() map[string]interface{} {
	allJson := PermissionForOrg
	res := make(map[string][]string, 0)
	json.FromJson(allJson, &res)
	res1 := make(map[string]interface{}, 0)
	for k, item := range res {
		res1[k] = item
	}
	return res1
}

// 对于无项目（未分配项目）的任务，其参与者具有更改任务的权限
func GetPermissionForNoProIssue() map[string]interface{} {
	allJson := PermissionForNoProIssue
	res := make(map[string][]string, 0)
	json.FromJson(allJson, &res)
	res1 := make(map[string]interface{}, 0)
	for k, item := range res {
		res1[k] = item
	}
	return res1
}

// GetPermissionForPro 获取**项目**所有的权限操作组以及操作项。一般是超管、普通管理员、应用的管理员具备
func GetPermissionForPro() map[string]interface{} {
	allJson := PermissionForPro
	res := make(map[string][]string, 0)
	json.FromJson(allJson, &res)
	res1 := make(map[string]interface{}, 0)
	for k, item := range res {
		res1[k] = item
	}
	return res1
}

// GetPermissionDefaultOperationForPro 获取项目“默认角色”的权限项
func GetPermissionDefaultOperationForPro() map[string]interface{} {
	allJson := PermissionDefaultOperationForPro
	res := make(map[string][]string, 0)
	json.FromJson(allJson, &res)
	res1 := make(map[string]interface{}, 0)
	for k, item := range res {
		res1[k] = item
	}
	return res1
}
