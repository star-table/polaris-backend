package consts

var (
	// 应用对应的权限组
	LcAppPermissionGroupOpsReadJson  = `["Permission.Org.ProjectObjectType-View","Permission.Pro.Config-View","Permission.Pro.Config-View","Permission.Pro.Iteration-View","Permission.Pro.Issue.4-View","Permission.Pro.Role-View","Permission.Pro.File-Download","Permission.Pro.File-View","Permission.Pro.Attachment-View","Permission.Pro.Member-View", "hasRead"]`
	LcAppPermissionGroupOpsWriteJson = `["Permission.Org.ProjectObjectType-Modify","Permission.Org.ProjectObjectType-Create","Permission.Org.ProjectObjectType-Delete","Permission.Pro.Config-View","Permission.Pro.Config-Modify,Bind,Unbind","Permission.Pro.Config-Filing,UnFiling","Permission.Pro.Config-ModifyStatus","Permission.Pro.Config-ModifyField","Permission.Pro.Iteration-Modify","Permission.Pro.Iteration-Create","Permission.Pro.Iteration-Delete","Permission.Pro.Iteration-ModifyStatus","Permission.Pro.Iteration-Bind,Unbind","Permission.Pro.Issue.4-Modify,Bind,Unbind","Permission.Pro.Issue.4-Create","Permission.Pro.Issue.4-Delete","Permission.Pro.Issue.4-ModifyStatus","Permission.Pro.Issue.4-Comment","Permission.Pro.Role-Modify","Permission.Pro.Role-Create","Permission.Pro.Role-Delete","Permission.Pro.Role-ModifyPermission","Permission.Pro.File-Download","Permission.Pro.File-Modify","Permission.Pro.File-Delete","Permission.Pro.File-CreateFolder","Permission.Pro.File-ModifyFolder","Permission.Pro.File-DeleteFolder","Permission.Pro.Tag-Delete","Permission.Pro.Tag-Remove","Permission.Pro.Tag-Modify","Permission.Pro.Attachment-Download","Permission.Pro.Attachment-Delete","Permission.Pro.Member-Bind","Permission.Pro.Member-Unbind", "hasRead", "hasCreate", "hasCopy", "hasUpdate", "hasDelete", "hasImport", "hasExport", "hasShare"]`
	LcAppPermissionGroupOpsAdminJson = `["Permission.Pro.Automation-Manage","Permission.Org.ProjectObjectType-Modify", "Permission.Org.ProjectObjectType-Create", "Permission.Org.ProjectObjectType-Delete", "Permission.Pro.Config-View", "Permission.Pro.Config-Modify,Bind,Unbind", "Permission.Pro.Config-Filing,UnFiling", "Permission.Pro.Config-ModifyStatus", "Permission.Pro.Config-ModifyField", "Permission.Pro.Iteration-Modify", "Permission.Pro.Iteration-Create", "Permission.Pro.Iteration-Delete", "Permission.Pro.Iteration-ModifyStatus", "Permission.Pro.Iteration-Bind,Unbind", "Permission.Pro.Issue.4-Modify,Bind,Unbind", "Permission.Pro.Issue.4-Create", "Permission.Pro.Issue.4-Delete", "Permission.Pro.Issue.4-ModifyStatus", "Permission.Pro.Issue.4-Comment", "Permission.Pro.Role-Modify", "Permission.Pro.Role-Create", "Permission.Pro.Role-Delete", "Permission.Pro.Role-ModifyPermission", "Permission.Pro.File-Download", "Permission.Pro.File-Modify", "Permission.Pro.File-Delete", "Permission.Pro.File-CreateFolder", "Permission.Pro.File-ModifyFolder", "Permission.Pro.File-DeleteFolder", "Permission.Pro.Tag-Delete", "Permission.Pro.Tag-Remove", "Permission.Pro.Tag-Modify", "Permission.Pro.Attachment-Download", "Permission.Pro.Attachment-Delete", "Permission.Pro.Member-Bind", "Permission.Pro.Member-Unbind", "hasRead", "hasCreate", "hasCopy", "hasUpdate", "hasDelete", "hasImport", "hasExport", "hasInvite", "hasShare", "hasEditMember"]`

	// 权限组的 langCode 值
	// 包含：查看权限-`-2`，编辑权限-`-3`，管理员-`-4`。详情可以咨询无码系统开发者 @刘千源
	AppPermissionGroupLangCodeForRead  = -2
	AppPermissionGroupLangCodeForWrite = -3
	AppPermissionGroupLangCodeForAdmin = -4

	// 无码自定义字段中，部门类型的字段的值的 type 字段值
	LcCustomFieldDeptType = "D_"
	// 用户类型标识
	LcCustomFieldUserType = "U_"

	// 项目管理员权限组的langCode
	AppProjectPermissionAdmin = "41"
)
