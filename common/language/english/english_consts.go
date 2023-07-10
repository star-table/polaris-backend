package english

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/util/lang"
)

var ProjectObjectTypeLang = map[string]string{
	consts.ProjectObjectTypeLangCodeTask:      "Task",
	consts.ProjectObjectTypeLangCodeDemand:    "Demand",
	consts.ProjectObjectTypeLangCodeBug:       "Bug",
	consts.ProjectObjectTypeLangCodeIteration: "Iteration",
}

//敏捷项目状态对应关系
var StatusLang = map[string]string{
	consts.StatusLangCodeProjectNotStart:     "Pending",
	consts.StatusLangCodeProjectRunning:      "In Progress",
	consts.StatusLangCodeProjectComplete:     "Completed",
	consts.StatusLangCodeIterationNotStart:   "Pending",
	consts.StatusLangCodeIterationRunning:    "In Progress",
	consts.StatusLangCodeIterationComplete:   "Completed",
	consts.StatusLangCodeIssueNotStart:       "Pending",
	consts.StatusLangCodeIssueWaitEvaluate:   "WaitEvaluate",
	consts.StatusLangCodeIssueWaitConfirmBug: "WaitConfirmBug",
	consts.StatusLangCodeIssueConfirmedBug:   "ConfirmedBug",
	consts.StatusLangCodeIssueReOpen:         "ReOpen",
	consts.StatusLangCodeIssueEvaluated:      "Evaluated",
	consts.StatusLangCodeIssuePlanning:       "Planning",
	consts.StatusLangCodeIssueDesign:         "Design",
	consts.StatusLangCodeIssueProcessing:     "Processing",
	consts.StatusLangCodeIssueDevelopment:    "Development",
	consts.StatusLangCodeIssueWaitTest:       "WaitTest",
	consts.StatusLangCodeIssueTesting:        "Testing",
	consts.StatusLangCodeIssueWaitRelease:    "WaitRelease",
	consts.StatusLangCodeIssueWait:           "Wait",
	consts.StatusLangCodeIssueRepair:         "Repair",
	consts.StatusLangCodeIssueSuccess:        "Success",
	consts.StatusLangCodeIssueFail:           "Fail",
	consts.StatusLangCodeIssueReleased:       "Released",
	consts.StatusLangCodeIssueClosed:         "Closed",
	consts.StatusLangCodeIssueComplete:       "Completed",
	consts.StatusLangCodeIssueConfirmed:      "Confirmed",
}

//权限
var PermissionLang = map[string]string{
	consts.PermissionSysSys:               " System Management",
	consts.PermissionSysDic:               "NewData Dictionary Management",
	consts.PermissionSysSource:            "Source Channel Management",
	consts.PermissionSysPayLevel:          "Order Level Management",
	consts.PermissionOrgOrg:               "Organization-related Authority",
	consts.PermissionOrgConfig:            "Organizational Setting Management",
	consts.PermissionOrgMessageConfig:     "System Message Settings",
	consts.PermissionOrgUser:              "Organize User Management",
	consts.PermissionOrgTeam:              "Organize Team Management",
	consts.PermissionOrgRoleGroup:         "Organizational Role Group Management",
	consts.PermissionOrgRole:              "Organizational Role Group Management",
	consts.PermissionOrgProject:           "Organize Project Management",
	consts.PermissionOrgProjectType:       "Project Type Management",
	consts.PermissionOrgIssueSource:       "Issue Source Management",
	consts.PermissionOrgProjectObjectType: "Taskbar Management",
	consts.PermissionOrgPriority:          "Priority Management",
	consts.PermissionOrgProcessStatus:     "Process State Management",
	consts.PermissionOrgProcess:           "process Management ",
	consts.PermissionOrgProcessStep:       "Process Steps Management",
	consts.PermissionProPro:               "Project Related Authority",
	consts.PermissionProConfig:            "Project Setup Management",
	consts.PermissionProBan:               "Project Panel Management",
	consts.PermissionProIteration:         "Project Iteration Management",
	consts.PermissionProIssue:             "Project Problem Management",
	consts.PermissionProIssue2:            "Project Feature Management",
	consts.PermissionProIssue3:            "Project Demand Management",
	consts.PermissionProIssue4:            "Project Task Management",
	consts.PermissionProIssue5:            "Project Defect Management",
	consts.PermissionProIssue6:            "Project Test Task Management",
	consts.PermissionProComment:           "Comments Management",
	consts.PermissionProProjectVersion:    "Project Version Management",
	consts.PermissionProProjectModule:     "Project Module Management",
	consts.PermissionProRole:              "Project Role Management",
	consts.PermissionProTest:              "Project Test Management",
	consts.PermissionProTestTestApp:       "Test Application Management",
	consts.PermissionProTestTestDevice:    "Test Equipment Management",
	consts.PermissionProTestTestReport:    "Test Report Management",
	consts.PermissionProFile:              "File Management ",
	consts.PermissionProTag:               "Label Management",
	consts.PermissionProAttachment:        "Attachment Management",
	consts.PermissionProMember:            "Project Member Management",
	consts.PermissionOrgDepartment:        "Organizational Structure Management",
}

var OperationLang = map[string]string{
	consts.OperationSysDicView:                 "View",
	consts.OperationSysDicCreate:               "Create",
	consts.OperationSysDicModify:               "Modify",
	consts.OperationSysDicDelete:               "Delete",
	consts.OperationSysSourceView:              "View",
	consts.OperationSysSourceCreate:            "Create",
	consts.OperationSysSourceModify:            "Modify",
	consts.OperationSysSourceDelete:            "Delete",
	consts.OperationSysPayLevelView:            "View",
	consts.OperationSysPayLevelCreate:          "Create",
	consts.OperationSysPayLevelModify:          "Modify",
	consts.OperationSysPayLevelDelete:          "Delete",
	consts.OperationOrgConfigView:              "View",
	consts.OperationOrgConfigModify:            "Modify",
	consts.OperationOrgConfigTransfer:          "Transfer",
	consts.OperationOrgMessageConfigView:       "View",
	consts.OperationOrgMessageConfigModify:     "Modify",
	consts.OperationOrgUserView:                "View",
	consts.OperationOrgUserModifyStatus:        "ModifyStatus",
	consts.OperationOrgUserInvite:              "Invite",
	consts.OperationOrgTeamView:                "View",
	consts.OperationOrgTeamCreate:              "Create",
	consts.OperationOrgTeamModify:              "Modify",
	consts.OperationOrgTeamDelete:              "Delete",
	consts.OperationOrgTeamModifyStatus:        "ModifyStatus",
	consts.OperationOrgTeamBind:                "Bind",
	consts.OperationOrgRoleGroupView:           "View",
	consts.OperationOrgRoleGroupCreate:         "Create",
	consts.OperationOrgRoleGroupModify:         "Modify",
	consts.OperationOrgRoleGroupDelete:         "Delete",
	consts.OperationOrgRoleView:                "View",
	consts.OperationOrgRoleCreate:              "Create",
	consts.OperationOrgRoleModify:              "Modify",
	consts.OperationOrgRoleDelete:              "Delete",
	consts.OperationOrgRoleBind:                "Bind",
	consts.OperationOrgProjectView:             "View",
	consts.OperationOrgProjectCreate:           "Create",
	consts.OperationOrgProjectModify:           "Modify",
	consts.OperationOrgProjectDelete:           "Delete",
	consts.OperationOrgProjectAttention:        "Attention",
	consts.OperationOrgProjectFiling:           "Filing",
	consts.OperationOrgProjectTypeView:         "View",
	consts.OperationOrgProjectTypeModify:       "Modify",
	consts.OperationOrgProjectTypeCreate:       "Create",
	consts.OperationOrgProjectTypeDelete:       "Delete",
	consts.OperationOrgIssueSourceView:         "View",
	consts.OperationOrgIssueSourceModify:       "Modify",
	consts.OperationOrgIssueSourceCreate:       "Create",
	consts.OperationOrgIssueSourceDelete:       "Delete",
	consts.OperationOrgProjectObjectTypeView:   "View",
	consts.OperationOrgProjectObjectTypeModify: "Modify",
	consts.OperationOrgProjectObjectTypeCreate: "Create",
	consts.OperationOrgProjectObjectTypeDelete: "Delete",
	consts.OperationOrgPriorityView:            "View",
	consts.OperationOrgPriorityModify:          "Modify",
	consts.OperationOrgPriorityCreate:          "Create",
	consts.OperationOrgPriorityDelete:          "Delete",
	consts.OperationOrgProcessStatusView:       "View",
	consts.OperationOrgProcessStatusModify:     "Modify",
	consts.OperationOrgProcessStatusCreate:     "Create",
	consts.OperationOrgProcessStatusDelete:     "Delete",
	consts.OperationOrgProcessView:             "View",
	consts.OperationOrgProcessModify:           "Modify",
	consts.OperationOrgProcessCreate:           "Create",
	consts.OperationOrgProcessDelete:           "Delete",
	consts.OperationOrgProcessStepView:         "View",
	consts.OperationOrgProcessStepModify:       "Modify",
	consts.OperationOrgProcessStepCreate:       "Create",
	consts.OperationOrgProcessStepDelete:       "Delete",
	consts.OperationProConfigView:              "View",
	consts.OperationProConfigModify:            "Modify",
	consts.OperationProConfigFiling:            "Filing",
	consts.OperationProConfigModifyStatus:      "ModifyStatus",
	consts.OperationProBanView:                 "View",
	consts.OperationProIterationView:           "View",
	consts.OperationProIterationModify:         "Modify",
	consts.OperationProIterationCreate:         "Create",
	consts.OperationProIterationDelete:         "Delete",
	consts.OperationProIterationModifyStatus:   "ModifyStatus",
	consts.OperationProIterationBind:           "Bind",
	consts.OperationProIterationAttention:      "Attention",
	consts.OperationProIssue2View:              "View",
	consts.OperationProIssue2Modify:            "Modify",
	consts.OperationProIssue2Create:            "Create",
	consts.OperationProIssue2Delete:            "Delete",
	consts.OperationProIssue2ModifyStatus:      "ModifyStatus",
	consts.OperationProIssue2Comment:           "Comment",
	consts.OperationProIssue2Attention:         "Attention",
	consts.OperationProIssue3View:              "View",
	consts.OperationProIssue3Modify:            "Modify",
	consts.OperationProIssue3Create:            "Create",
	consts.OperationProIssue3Delete:            "Delete",
	consts.OperationProIssue3ModifyStatus:      "ModifyStatus",
	consts.OperationProIssue3Comment:           "Comment",
	consts.OperationProIssue3Attention:         "Attention",
	consts.OperationProIssue4View:              "View",
	consts.OperationProIssue4Modify:            "Modify",
	consts.OperationProIssue4Create:            "Create",
	consts.OperationProIssue4Delete:            "Delete",
	consts.OperationProIssue4ModifyStatus:      "ModifyStatus",
	consts.OperationProIssue4Comment:           "Comment",
	consts.OperationProIssue4Attention:         "Attention",
	consts.OperationProIssue5View:              "View",
	consts.OperationProIssue5Modify:            "Modify",
	consts.OperationProIssue5Create:            "Create",
	consts.OperationProIssue5Delete:            "Delete",
	consts.OperationProIssue5ModifyStatus:      "ModifyStatus",
	consts.OperationProIssue5Comment:           "Comment",
	consts.OperationProIssue5Attention:         "Attention",
	consts.OperationProCommentModify:           "Modify",
	consts.OperationProCommentDelete:           "Delete",
	consts.OperationProProjectVersionView:      "View",
	consts.OperationProProjectVersionModify:    "Modify",
	consts.OperationProProjectVersionCreate:    "Create",
	consts.OperationProProjectVersionDelete:    "Delete",
	consts.OperationProProjectModuleView:       "View",
	consts.OperationProProjectModuleModify:     "Modify",
	consts.OperationProProjectModuleCreate:     "Create",
	consts.OperationProProjectModuleDelete:     "Delete",
	consts.OperationProRoleView:                "View",
	consts.OperationProRoleModify:              "Modify",
	consts.OperationProRoleBind:                "Bind",
	consts.OperationProTestTestAppView:         "View",
	consts.OperationProTestTestAppCreate:       "Create",
	consts.OperationProTestTestAppModify:       "Modify",
	consts.OperationProTestTestAppDelete:       "Delete",
	consts.OperationProTestTestDeviceView:      "View",
	consts.OperationProTestTestDeviceCreate:    "Create",
	consts.OperationProTestTestDeviceModify:    "Modify",
	consts.OperationProTestTestDeviceDelete:    "Delete",
	consts.OperationProTestTestReportView:      "View",
	consts.OperationProTestTestReportCreate:    "Create",
	consts.OperationProTestTestReportModify:    "Modify",
	consts.OperationProTestTestReportDelete:    "Delete",
	consts.OperationProFileUpload:              "Upload",
	consts.OperationProFileDownload:            "Download",
	consts.OperationProFileModify:              "Modify",
	consts.OperationProFileDelete:              "Delete",
	consts.OperationProFileCreateFolder:        "CreateFolder",
	consts.OperationProFileModifyFolder:        "ModifyFolder",
	consts.OperationProFileDeleteFolder:        "DeleteFolder",
	consts.OperationProTagCreate:               "Create",
	consts.OperationProTagDelete:               "Delete",
	consts.OperationProTagRemove:               "Remove",
	consts.OperationProAttachmentUpload:        "Upload",
	consts.OperationProAttachmentDownload:      "Download",
	consts.OperationProAttachmentDelete:        "Delete",
	consts.OperationProMemberBind:              "Bind",
	consts.OperationProMemberUnbind:            "Unbind",
	consts.OperationProMemberCreate:            "Create",
	consts.OperationProMemberDelete:            "Delete",
	consts.OperationProMemberModifyOperation:   "ModifyOperation",
	consts.OperationProTagModify:               "Modify",
	consts.OperationOrgDepartmentCreate:        "Create",
	consts.OperationOrgDepartmentModify:        "Modify",
	consts.OperationOrgDepartmentDelete:        "Delete",
	consts.OperationOrgUserWatch:               "Watch",
	consts.OperationOrgUserModifyDepartment:    "ModifyDepartment",
	consts.OperationProConfigModifyField:       "ModifyField",
	consts.OperationOrgProjectModifyField:      "ModifyField",
}

var TrendsLang = map[string]string{
	consts.Code:                  "code",
	consts.Title:                 "title",
	consts.Owner:                 "owner",
	consts.Status:                "status",
	consts.PlanStartTime:         "plan start time",
	consts.PlanEndTime:           "end time",
	consts.PlanWorkHour:          "plan work hour",
	consts.Priority:              "priority",
	consts.Source:                "source",
	consts.IssueObjectTypeId:     "type",
	consts.Remark:                "remark",
	consts.PublicStatus:          "public status",
	consts.ProjectNotice:         "project notice",
	consts.ProjectResourcePath:   "project cover",
	consts.ProjectResourceName:   "project resource name",
	consts.ProjectResourceFolder: "resource folder",
	consts.ProjectFolderName:     "folder name",
	consts.ProjectFolderParentId: "folder parent id",
	consts.ProjectObjectType:     "taskbar",
	consts.IssuePropertyId:       "property",
	consts.Iteration:             "iteration",
	consts.Name:                  "name",
	consts.OrgField:              "organization field",
	consts.Value:                 "value",
	consts.Description:           "description",
	consts.WorkId:                "performer",
	consts.StartTime:             "start time",
	consts.WorkTime:              "work time",
	consts.WorkContent:           "job content",
}

var ProRolesLang = map[string]string{
	"负责人":  "Admin",
	"项目成员": "Member",
}

var ProObjectTypeLang = map[string]string{}

var ProjectTypeLang = map[string]string{
	consts.ProjectTypeLangCodeNormalTask: "Common Project",
	consts.ProjectTypeLangCodeAgile:      "Agile Project",
}

var WordDictionary = map[string]string{
	"未分配": "Unallocated",
	"待确认": "Unconfirmed",
	"已完成": "Completed",
}

func WordTransLate(data string) string {
	if lang.IsEnglish() {
		if name, ok := WordDictionary[data]; ok {
			return name
		}
	}

	return data
}
