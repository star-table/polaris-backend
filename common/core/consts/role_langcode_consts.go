package consts

const (
	RoleSysAdmin                 = "Role.Sys.Admin"
	RoleSysManager               = "Role.Sys.Manager"
	RoleSysMember                = "Role.Sys.Member"
	RoleGroupSpecialCreator      = "RoleGroup.Special.Creator"
	RoleGroupSpecialOwner        = "RoleGroup.Special.Owner"
	RoleGroupSpecialWorker       = "RoleGroup.Special.Worker"
	RoleGroupSpecialAttention    = "RoleGroup.Special.Attention"
	RoleGroupSpecialMember       = "RoleGroup.Special.Member"
	RoleGroupSpecialVisitor      = "RoleGroup.Special.Visitor"
	RoleGroupOrgAdmin            = "RoleGroup.Org.Admin"
	RoleGroupOrgManager          = "RoleGroup.Org.Manager"
	RoleGroupProProjectManager   = "RoleGroup.Pro.ProjectManager"
	RoleGroupProTechnicalManager = "RoleGroup.Pro.TechnicalManager"
	RoleGroupProProductManager   = "RoleGroup.Pro.ProductManager"
	RoleGroupProDeveloper        = "RoleGroup.Pro.Developer"
	RoleGroupProTester           = "RoleGroup.Pro.Tester"
	RoleGroupProMember           = "RoleGroup.Pro.Member"
)

const (
	GroupLandCodeRead                   = "-2" // "查看者" "只可查看当前文件夹或者是数据表的信息，不可邀请成员和新成员信息"
	GroupLandCodeEdit                   = "-3" // "编辑者" "只可编辑当前文件夹或者是数据表的信息，不可邀请成员和新成员信息"
	GroupLandCodeFormAdministrator      = "-4" // "管理员" "可直接管理当前文件夹或者是数据表所有权限信息，并可邀请成员和新成员"
	GroupLandCodeDashboardAdministrator = "-5" // "管理员" "可直接管理当前文件夹或者是数据表所有权限信息，并可邀请成员和新成员"
	GroupLandCodeOwner                  = "41" // "管理员" "当前项目的管理员，具有项目的所有管理权限"
	GroupLandCodeProjectMember          = "42" // "编辑者" "当前项目的编辑者，可操作任务、文件、标签、附件"
	GroupLandCodeProjectViewer          = "43" // "查看者" "查看者"
)

const (
	RoleGroupSys     = "RoleGroup.Sys"
	RoleGroupSpecial = "RoleGroup.Special"
	RoleGroupOrg     = "RoleGroup.Org"
	RoleGroupPro     = "RoleGroup.Pro"
)

const (
	RoleOperationView             = "View"
	RoleOperationModify           = "Modify"
	RoleOperationDelete           = "Delete"
	RoleOperationCreate           = "Create"
	RoleOperationCheck            = "Check"
	RoleOperationInvite           = "Invite"
	RoleOperationBind             = "Bind"
	RoleOperationUnbind           = "Unbind"
	RoleOperationAttention        = "Attention"
	RoleOperationUnAttention      = "UnAttention"
	RoleOperationModifyStatus     = "ModifyStatus"
	RoleOperationComment          = "Comment"
	RoleOperationTransfer         = "Transfer"
	RoleOperationInit             = "Init"
	RoleOperationDrop             = "Drop"
	RoleOperationFiling           = "Filing"
	RoleOperationUnFiling         = "UnFiling"
	RoleOperationUpload           = "Upload"
	RoleOperationDownload         = "Download"
	RoleOperationRemove           = "Remove"
	RoleOperationModifyPermission = "ModifyPermission"
	RoleOperationCreateFolder     = "CreateFolder"
	RoleOperationModifyFolder     = "ModifyFolder"
	RoleOperationDeleteFolder     = "DeleteFolder"
	RoleOperationModifyDepartment = "ModifyDepartment"
	RoleOperationAdd              = "Add"
	RoleOperationModifyField      = "ModifyField"
	RoleOperationEnableWorkHour   = "EnableWorkHour"
	RoleOperationDisableWorkHour  = "DisableWorkHour"
	RoleOperationWithdraw         = "Withdraw"
	RoleOperationStartIssueChat   = "StartIssueChat" // 发起任务群聊
)

const (
	RoleOperationPathSys                  = "/Sys"
	RoleOperationPathSysDic               = "/Sys/Dic"
	RoleOperationPathSysSource            = "/Sys/Source"
	RoleOperationPathSysPayLevel          = "/Sys/PayLevel"
	RoleOperationPathOrgProRole           = "/Org/{org_id}/Pro/{pro_id}/Role"
	RoleOperationPathOrgProProjectVersion = "/Org/{org_id}/Pro/{pro_id}/ProjectVersion"
	RoleOperationPathOrgProProjectModule  = "/Org/{org_id}/Pro/{pro_id}/ProjectModule"
	RoleOperationPathOrgPro               = "/Org/{org_id}/Pro/{pro_id}"
	RoleOperationPathOrgProIteration      = "/Org/{org_id}/Pro/{pro_id}/Iteration"
	RoleOperationPathOrgProIssueTT        = "/Org/{org_id}/Pro/{pro_id}/Issue/6"
	RoleOperationPathOrgProIssueT         = "/Org/{org_id}/Pro/{pro_id}/Issue/4"
	RoleOperationPathOrgProIssueF         = "/Org/{org_id}/Pro/{pro_id}/Issue/2"
	RoleOperationPathOrgProIssueD         = "/Org/{org_id}/Pro/{pro_id}/Issue/3"
	RoleOperationPathOrgProIssueB         = "/Org/{org_id}/Pro/{pro_id}/Issue/5"
	RoleOperationPathOrgProIssue          = "/Org/{org_id}/Pro/{pro_id}/Issue"
	RoleOperationPathOrgProProConfig      = "/Org/{org_id}/Pro/{pro_id}/ProConfig"
	RoleOperationPathOrgProComment        = "/Org/{org_id}/Pro/{pro_id}/Comment"
	RoleOperationPathOrgProBan            = "/Org/{org_id}/Pro/{pro_id}/Ban"
	RoleOperationPathOrgUser              = "/Org/{org_id}/User"
	RoleOperationPathOrgTeam              = "/Org/{org_id}/Team"
	RoleOperationPathOrgRoleGroup         = "/Org/{org_id}/RoleGroup"
	RoleOperationPathOrgRole              = "/Org/{org_id}/Role"
	RoleOperationPathOrgProjectType       = "/Org/{org_id}/ProjectType"
	RoleOperationPathOrgProjectObjectType = "/Org/{org_id}/ProjectObjectType"
	RoleOperationPathOrgProject           = "/Org/{org_id}/Project"
	RoleOperationPathOrgProcessStep       = "/Org/{org_id}/ProcessStep"
	RoleOperationPathOrgProcessStatus     = "/Org/{org_id}/ProcessStatus"
	RoleOperationPathOrgProcess           = "/Org/{org_id}/Process"
	RoleOperationPathOrgPriority          = "/Org/{org_id}/Priority"
	RoleOperationPathOrg                  = "/Org/{org_id}"
	RoleOperationPathOrgMessageConfig     = "/Org/{org_id}/MessageConfig"
	RoleOperationPathOrgIssueSource       = "/Org/{org_id}/IssueSource"
	RoleOperationPathOrgOrgConfig         = "/Org/{org_id}/OrgConfig"
	RoleOperationPathOrgProMember         = "/Org/{org_id}/Pro/{pro_id}/Member"
	RoleOperationPathOrgProFile           = "/Org/{org_id}/Pro/{pro_id}/File"
	RoleOperationPathOrgProTag            = "/Org/{org_id}/Pro/{pro_id}/Tag"
	RoleOperationPathOrgProAttachment     = "/Org/{org_id}/Pro/{pro_id}/Attachment"
)

//权限项
const (
	PermissionSysSys               = "Permission.Sys.Sys"
	PermissionSysDic               = "Permission.Sys.Dic"
	PermissionSysSource            = "Permission.Sys.Source"
	PermissionSysPayLevel          = "Permission.Sys.PayLevel"
	PermissionOrgOrg               = "Permission.Org.Org"
	PermissionOrgConfig            = "Permission.Org.Config"
	PermissionOrgMessageConfig     = "Permission.Org.MessageConfig"
	PermissionOrgUser              = "Permission.Org.User"
	PermissionOrgTeam              = "Permission.Org.Team"
	PermissionOrgRoleGroup         = "Permission.Org.RoleGroup"
	PermissionOrgRole              = "Permission.Org.Role"
	PermissionOrgProject           = "Permission.Org.Project"
	PermissionOrgProjectType       = "Permission.Org.ProjectType"
	PermissionOrgIssueSource       = "Permission.Org.IssueSource"
	PermissionOrgProjectObjectType = "Permission.Org.ProjectObjectType"
	PermissionOrgPriority          = "Permission.Org.Priority"
	PermissionOrgProcessStatus     = "Permission.Org.ProcessStatus"
	PermissionOrgProcess           = "Permission.Org.Process"
	PermissionOrgProcessStep       = "Permission.Org.ProcessStep"
	PermissionProPro               = "Permission.Pro.Pro"
	PermissionProConfig            = "Permission.Pro.Config"
	PermissionProBan               = "Permission.Pro.Ban"
	PermissionProIteration         = "Permission.Pro.Iteration"
	PermissionProIssue             = "Permission.Pro.Issue"
	PermissionProIssue2            = "Permission.Pro.Issue.2"
	PermissionProIssue3            = "Permission.Pro.Issue.3"
	PermissionProIssue4            = "Permission.Pro.Issue.4"
	PermissionProIssue5            = "Permission.Pro.Issue.5"
	PermissionProIssue6            = "Permission.Pro.Issue.6"
	PermissionProComment           = "Permission.Pro.Comment"
	PermissionProProjectVersion    = "Permission.Pro.ProjectVersion"
	PermissionProProjectModule     = "Permission.Pro.ProjectModule"
	PermissionProRole              = "Permission.Pro.Role"
	PermissionProTest              = "Permission.Pro.Test"
	PermissionProTestTestApp       = "Permission.Pro.Test.TestApp"
	PermissionProTestTestDevice    = "Permission.Pro.Test.TestDevice"
	PermissionProTestTestReport    = "Permission.Pro.Test.TestReport"
	PermissionProFile              = "Permission.Pro.File"
	PermissionProTag               = "Permission.Pro.Tag"
	PermissionProAttachment        = "Permission.Pro.Attachment"
	PermissionProMember            = "Permission.Pro.Member"
	PermissionOrgDepartment        = "Permission.Org.Department"
)

//权限操作项
const (
	OperationSysDicView                 = "PermissionOperation.Sys.Dic.View"
	OperationSysDicCreate               = "PermissionOperation.Sys.Dic.Create"
	OperationSysDicModify               = "PermissionOperation.Sys.Dic.Modify"
	OperationSysDicDelete               = "PermissionOperation.Sys.Dic.Delete"
	OperationSysSourceView              = "PermissionOperation.Sys.Source.View"
	OperationSysSourceCreate            = "PermissionOperation.Sys.Source.Create"
	OperationSysSourceModify            = "PermissionOperation.Sys.Source.Modify"
	OperationSysSourceDelete            = "PermissionOperation.Sys.Source.Delete"
	OperationSysPayLevelView            = "PermissionOperation.Sys.PayLevel.View"
	OperationSysPayLevelCreate          = "PermissionOperation.Sys.PayLevel.Create"
	OperationSysPayLevelModify          = "PermissionOperation.Sys.PayLevel.Modify"
	OperationSysPayLevelDelete          = "PermissionOperation.Sys.PayLevel.Delete"
	OperationOrgConfigView              = "Permission.Org.Config.View"
	OperationOrgConfigModify            = "Permission.Org.Config.Modify"
	OperationOrgConfigTransfer          = "Permission.Org.Config.Transfer"
	OperationOrgConfigModifyField       = "Permission.Org.Config.ModifyField"
	OperationOrgMessageConfigView       = "Permission.Org.MessageConfig.View"
	OperationOrgMessageConfigModify     = "Permission.Org.MessageConfig.Modify"
	OperationOrgUserView                = "Permission.Org.User.View"
	OperationOrgUserModifyStatus        = "Permission.Org.User.ModifyStatus"
	OperationOrgUserInvite              = "Permission.Org.User.Invite"
	OperationOrgInviteUserInvite        = "Permission.Org.InviteUser.Invite" // 新的分组 产品：子龙
	OperationOrgTeamView                = "Permission.Org.Team.View"
	OperationOrgTeamCreate              = "Permission.Org.Team.Create"
	OperationOrgTeamModify              = "Permission.Org.Team.Modify"
	OperationOrgTeamDelete              = "Permission.Org.Team.Delete"
	OperationOrgTeamModifyStatus        = "Permission.Org.Team.ModifyStatus"
	OperationOrgTeamBind                = "Permission.Org.Team.Bind"
	OperationOrgRoleGroupView           = "Permission.Org.RoleGroup.View"
	OperationOrgRoleGroupCreate         = "Permission.Org.RoleGroup.Create"
	OperationOrgRoleGroupModify         = "Permission.Org.RoleGroup.Modify"
	OperationOrgRoleGroupDelete         = "Permission.Org.RoleGroup.Delete"
	OperationOrgRoleView                = "Permission.Org.Role.View"
	OperationOrgRoleCreate              = "Permission.Org.Role.Create"
	OperationOrgRoleModify              = "Permission.Org.Role.Modify"
	OperationOrgRoleDelete              = "Permission.Org.Role.Delete"
	OperationOrgRoleBind                = "Permission.Org.Role.Bind"
	OperationOrgProjectView             = "Permission.Org.Project.View"
	OperationOrgProjectCreate           = "Permission.Org.Project.Create"
	OperationOrgProjectModify           = "Permission.Org.Project.Modify"
	OperationOrgProjectDelete           = "Permission.Org.Project.Delete"
	OperationOrgProjectAttention        = "Permission.Org.Project.Attention"
	OperationOrgProjectFiling           = "Permission.Org.Project.Filing"
	OperationOrgProjectTypeView         = "Permission.Org.ProjectType.View"
	OperationOrgProjectTypeModify       = "Permission.Org.ProjectType.Modify"
	OperationOrgProjectTypeCreate       = "Permission.Org.ProjectType.Create"
	OperationOrgProjectTypeDelete       = "Permission.Org.ProjectType.Delete"
	OperationOrgIssueSourceView         = "Permission.Org.IssueSource.View"
	OperationOrgIssueSourceModify       = "Permission.Org.IssueSource.Modify"
	OperationOrgIssueSourceCreate       = "Permission.Org.IssueSource.Create"
	OperationOrgIssueSourceDelete       = "Permission.Org.IssueSource.Delete"
	OperationOrgProjectObjectTypeView   = "Permission.Org.ProjectObjectType.View"
	OperationOrgProjectObjectTypeModify = "Permission.Org.ProjectObjectType.Modify"
	OperationOrgProjectObjectTypeCreate = "Permission.Org.ProjectObjectType.Create"
	OperationOrgProjectObjectTypeDelete = "Permission.Org.ProjectObjectType.Delete"
	OperationOrgPriorityView            = "Permission.Org.Priority.View"
	OperationOrgPriorityModify          = "Permission.Org.Priority.Modify"
	OperationOrgPriorityCreate          = "Permission.Org.Priority.Create"
	OperationOrgPriorityDelete          = "Permission.Org.Priority.Delete"
	OperationOrgProcessStatusView       = "Permission.Org.ProcessStatus.View"
	OperationOrgProcessStatusModify     = "Permission.Org.ProcessStatus.Modify"
	OperationOrgProcessStatusCreate     = "Permission.Org.ProcessStatus.Create"
	OperationOrgProcessStatusDelete     = "Permission.Org.ProcessStatus.Delete"
	OperationOrgProcessView             = "Permission.Org.Process.View"
	OperationOrgProcessModify           = "Permission.Org.Process.Modify"
	OperationOrgProcessCreate           = "Permission.Org.Process.Create"
	OperationOrgProcessDelete           = "Permission.Org.Process.Delete"
	OperationOrgProcessStepView         = "Permission.Org.ProcessStep.View"
	OperationOrgProcessStepModify       = "Permission.Org.ProcessStep.Modify"
	OperationOrgProcessStepCreate       = "Permission.Org.ProcessStep.Create"
	OperationOrgProcessStepDelete       = "Permission.Org.ProcessStep.Delete"
	OperationProConfigView              = "Permission.Pro.Config.View"
	OperationProConfigModify            = "Permission.Pro.Config.Modify"
	OperationProConfigFiling            = "Permission.Pro.Config.Filing"
	OperationProConfigUnFiling          = "Permission.Pro.Config.UnFiling"
	OperationProConfigModifyStatus      = "Permission.Pro.Config.ModifyStatus"
	OperationProConfigDelete            = "Permission.Pro.Config.Delete"
	OperationProBanView                 = "Permission.Pro.Ban.View"
	OperationProIterationView           = "Permission.Pro.Iteration.View"
	OperationProIterationModify         = "Permission.Pro.Iteration.Modify"
	OperationProIterationCreate         = "Permission.Pro.Iteration.Create"
	OperationProIterationDelete         = "Permission.Pro.Iteration.Delete"
	OperationProIterationModifyStatus   = "Permission.Pro.Iteration.ModifyStatus"
	OperationProIterationBind           = "Permission.Pro.Iteration.Bind"
	OperationProIterationAttention      = "Permission.Pro.Iteration.Attention"
	// OperationProProjectObjectTypeModify 管理表权限（子表）（原：任务对象类型）。但考虑到任务栏无需特定的权限管理（可通过字段权限可控制）。
	OperationProProjectObjectTypeModify = "Permission.Pro.ProjectObjectType.Modify"
	OperationProProjectObjectTypeCreate = "Permission.Pro.ProjectObjectType.Create"
	OperationProProjectObjectTypeDelete = "Permission.Pro.ProjectObjectType.Delete"
	// 项目下，表的管理权限操作项
	OperationProProjectTableCreate     = "Permission.Pro.TableManage.Create"
	OperationProProjectTableModify     = "Permission.Pro.TableManage.Modify"
	OperationProProjectTableDelete     = "Permission.Pro.TableManage.Delete"
	OperationProIssue2View             = "Permission.Pro.Issue.2.View"
	OperationProIssue2Modify           = "Permission.Pro.Issue.2.Modify"
	OperationProIssue2Create           = "Permission.Pro.Issue.2.Create"
	OperationProIssue2Delete           = "Permission.Pro.Issue.2.Delete"
	OperationProIssue2ModifyStatus     = "Permission.Pro.Issue.2.ModifyStatus"
	OperationProIssue2Comment          = "Permission.Pro.Issue.2.Comment"
	OperationProIssue2Attention        = "Permission.Pro.Issue.2.Attention"
	OperationProIssue3View             = "Permission.Pro.Issue.3.View"
	OperationProIssue3Modify           = "Permission.Pro.Issue.3.Modify"
	OperationProIssue3Create           = "Permission.Pro.Issue.3.Create"
	OperationProIssue3Delete           = "Permission.Pro.Issue.3.Delete"
	OperationProIssue3ModifyStatus     = "Permission.Pro.Issue.3.ModifyStatus"
	OperationProIssue3Comment          = "Permission.Pro.Issue.3.Comment"
	OperationProIssue3Attention        = "Permission.Pro.Issue.3.Attention"
	OperationProIssue4View             = "Permission.Pro.Issue.4.View"
	OperationProIssue4Modify           = "Permission.Pro.Issue.4.Modify"
	OperationProIssue4Create           = "Permission.Pro.Issue.4.Create"
	OperationProIssue4Delete           = "Permission.Pro.Issue.4.Delete"
	OperationProIssue4ModifyStatus     = "Permission.Pro.Issue.4.ModifyStatus"
	OperationProIssue4Comment          = "Permission.Pro.Issue.4.Comment"
	OperationProIssue4Import           = "Permission.Pro.Issue.4.Import"
	OperationProIssue4Export           = "Permission.Pro.Issue.4.Export"
	OperationProIssue4Attention        = "Permission.Pro.Issue.4.Attention"
	OperationProIssue5View             = "Permission.Pro.Issue.5.View"
	OperationProIssue5Modify           = "Permission.Pro.Issue.5.Modify"
	OperationProIssue5Create           = "Permission.Pro.Issue.5.Create"
	OperationProIssue5Delete           = "Permission.Pro.Issue.5.Delete"
	OperationProIssue5ModifyStatus     = "Permission.Pro.Issue.5.ModifyStatus"
	OperationProIssue5Comment          = "Permission.Pro.Issue.5.Comment"
	OperationProIssue5Attention        = "Permission.Pro.Issue.5.Attention"
	OperationProTestIssue5View         = "Permission.Pro.Issue.5.View"
	OperationProTestIssue5Modify       = "Permission.Pro.Issue.5.Modify"
	OperationProTestIssue5Create       = "Permission.Pro.Issue.5.Create"
	OperationProTestIssue5Delete       = "Permission.Pro.Issue.5.Delete"
	OperationProTestIssue5ModifyStatus = "Permission.Pro.Issue.5.ModifyStatus"
	OperationProTestIssue5Comment      = "Permission.Pro.Issue.5.Comment"
	OperationProTestIssue5Attention    = "Permission.Pro.Issue.5.Attention"
	OperationProCommentModify          = "Permission.Pro.Comment.Modify"
	OperationProCommentDelete          = "Permission.Pro.Comment.Delete"
	OperationProProjectVersionView     = "Permission.Pro.ProjectVersion.View"
	OperationProProjectVersionModify   = "Permission.Pro.ProjectVersion.Modify"
	OperationProProjectVersionCreate   = "Permission.Pro.ProjectVersion.Create"
	OperationProProjectVersionDelete   = "Permission.Pro.ProjectVersion.Delete"
	OperationProProjectModuleView      = "Permission.Pro.ProjectModule.View"
	OperationProProjectModuleModify    = "Permission.Pro.ProjectModule.Modify"
	OperationProProjectModuleCreate    = "Permission.Pro.ProjectModule.Create"
	OperationProProjectModuleDelete    = "Permission.Pro.ProjectModule.Delete"
	OperationProRoleView               = "Permission.Pro.Role.View"
	OperationProRoleModify             = "Permission.Pro.Role.Modify"
	OperationProRoleBind               = "Permission.Pro.Role.Bind"
	OperationProTestTestAppView        = "Permission.Pro.Test.TestApp.View"
	OperationProTestTestAppCreate      = "Permission.Pro.Test.TestApp.Create"
	OperationProTestTestAppModify      = "Permission.Pro.Test.TestApp.Modify"
	OperationProTestTestAppDelete      = "Permission.Pro.Test.TestApp.Delete"
	OperationProTestTestDeviceView     = "Permission.Pro.Test.TestDevice.View"
	OperationProTestTestDeviceCreate   = "Permission.Pro.Test.TestDevice.Create"
	OperationProTestTestDeviceModify   = "Permission.Pro.Test.TestDevice.Modify"
	OperationProTestTestDeviceDelete   = "Permission.Pro.Test.TestDevice.Delete"
	OperationProTestTestReportView     = "Permission.Pro.Test.TestReport.View"
	OperationProTestTestReportCreate   = "Permission.Pro.Test.TestReport.Create"
	OperationProTestTestReportModify   = "Permission.Pro.Test.TestReport.Modify"
	OperationProTestTestReportDelete   = "Permission.Pro.Test.TestReport.Delete"
	OperationProFileView               = "Permission.Pro.File.View"
	OperationProFileUpload             = "Permission.Pro.File.Upload"
	OperationProFileDownload           = "Permission.Pro.File.Download"
	OperationProFileModify             = "Permission.Pro.File.Modify"
	OperationProFileDelete             = "Permission.Pro.File.Delete"
	OperationProFileCreateFolder       = "Permission.Pro.File.CreateFolder"
	OperationProFileModifyFolder       = "Permission.Pro.File.ModifyFolder"
	OperationProFileDeleteFolder       = "Permission.Pro.File.DeleteFolder"
	OperationProTagCreate              = "Permission.Pro.Tag.Create"
	OperationProTagDelete              = "Permission.Pro.Tag.Delete"
	OperationProTagRemove              = "Permission.Pro.Tag.Remove"
	OperationProAttachmentView         = "Permission.Pro.Attachment.View"
	OperationProAttachmentUpload       = "Permission.Pro.Attachment.Upload"
	OperationProAttachmentDownload     = "Permission.Pro.Attachment.Download"
	OperationProAttachmentDelete       = "Permission.Pro.Attachment.Delete"
	OperationProMemberBind             = "Permission.Pro.Member.Bind"
	OperationProMemberUnbind           = "Permission.Pro.Member.Unbind"
	OperationProMemberCreate           = "Permission.Pro.Member.Create"
	OperationProMemberDelete           = "Permission.Pro.Member.Delete"
	OperationProMemberModifyOperation  = "Permission.Pro.Member.ModifyPermission"
	OperationProTagModify              = "Permission.Pro.Tag.Modify"
	OperationOrgDepartmentCreate       = "Permission.Org.Department.Create"
	OperationOrgDepartmentModify       = "Permission.Org.Department.Modify"
	OperationOrgDepartmentDelete       = "Permission.Org.Department.Delete"
	OperationOrgUserBind               = "Permission.Org.User.Bind"
	OperationOrgUserUnbind             = "Permission.Org.User.Unbind"
	OperationOrgUserWatch              = "Permission.Org.User.Watch"
	OperationOrgUserModifyDepartment   = "Permission.Org.User.ModifyDepartment"
	OperationProConfigModifyField      = "Permission.Pro.Config.ModifyField"  // 项目的自定义字段
	OperationOrgProjectModifyField     = "Permission.Org.Project.ModifyField" // 组织的自定义字段
)

const (
	ManageGroupSys            = "ManageGroup.Sys"             // 超级管理员 1
	ManageGroupSubNormalAdmin = "ManageGroup.Sub.NormalAdmin" // bjx：普通管理员   3 // 这个数字和创建管理组入参的 groupType 对应
	ManageGroupSubNormalUser  = "ManageGroup.Sub.NormalUser"  // bjx：团队成员  4
	ManageGroupSubUserCustom  = "ManageGroup.Sub.UserCustom"  // 用户创建的管理组 6
)
