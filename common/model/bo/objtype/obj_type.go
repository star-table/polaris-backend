package objtype

var (
	/**
	系统
	*/
	Sys = "Sys"
	/**
	数据字典
	*/
	Dic = "Dic"
	/**
	注册来源渠道
	*/
	SourceChannel = "SourceChannel"
	/**
	订购级别
	*/
	PayLevel = "PayLevel"
	/**
	组织
	*/
	Org = "Org"
	/**
	组织系统设置
	*/
	OrgConfig = "OrgConfig"
	/**
	用户
	*/
	User = "User"
	/***
	团队
	*/
	Team = "Team"
	/**
	角色组
	*/
	RoleGroup = "RoleGroup"

	/**
	角色
	*/
	Role = "Role"

	/**
	权限项
	*/
	Permission = "Permission"

	/**
	操作
	*/
	Operation = "Operation"

	/**
	权限操作
	*/
	PermissionOperation = "PermissionOperation"

	/**
	资源
	*/
	Resource = "Resource"
	/**
	项目
	*/
	Project = "Project"

	/**
	面板
	*/
	Ban = "Ban"

	/**
	迭代
	*/
	Iteration = "Iteration"

	/**
	问题
	*/
	Issue = "Issue"

	/**
	优先级
	*/
	Priority = "Priority"

	/**
	流程状态
	*/
	ProcessStatus = "ProcessStatus"

	/**
	动态
	*/
	Trends = "Trends"

	/**
	评论
	*/
	Comment = "Comment"

	/**
	消息设置
	*/
	MessageConfig = "MessageConfig"

	/**
	项目版本
	*/
	ProjectVersion = "ProjectVersion"

	/**
	项目模块
	*/
	ProjectModule = "ProjectModule"

	/**
	问题来源
	*/
	IssueSource = "IssueSource"

	/**
	应用
	*/
	AppInfo = "AppInfo"

	/**
	用户设置
	*/
	UserConfig = "UserConfig"

	/**
	项目对象类型
	*/
	ProjectObjectType = "ProjectObjectType"

	/**
	流程
	*/
	Process = "Process"

	/**
	项目类型
	*/
	ProjectType = "ProjectType"

	/**
	流程步骤
	*/
	ProcessStep = "ProcessStep"
)

var (
	ObjTypes = map[string]ObjTypeInfoBo{
		/**
		系统
		*/
		Sys: ObjTypeInfoBo{
			LangCode: "ObjectType.Sys",
			Name:     "系统",
		},
		/**
		  数据字典
		*/
		Dic: ObjTypeInfoBo{
			LangCode: "ObjectType.Dic",
			Name:     "数据字典",
		},
		/**
		注册来源渠道
		*/
		SourceChannel: ObjTypeInfoBo{
			LangCode: "ObjectType.SourceChannel",
			Name:     "注册来源渠道",
		},
		/**
		订购级别
		*/
		PayLevel: ObjTypeInfoBo{
			LangCode: "ObjectType.PayLevel",
			Name:     "订购级别",
		},
		/**
		组织
		*/
		Org: ObjTypeInfoBo{
			LangCode: "ObjectType.Org",
			Name:     "组织",
		},
		/**
		组织系统设置
		*/
		OrgConfig: ObjTypeInfoBo{
			LangCode: "ObjectType.OrgConfig",
			Name:     "组织系统设置",
		},
		/**
		用户
		*/
		User: ObjTypeInfoBo{
			LangCode: "ObjectType.User",
			Name:     "用户",
		},
		/***
		团队
		*/
		Team: ObjTypeInfoBo{
			LangCode: "ObjectType.Team",
			Name:     "团队",
		},
		/**
		角色组
		*/
		RoleGroup: ObjTypeInfoBo{
			LangCode: "ObjectType.角色组",
			Name:     "RoleGroup",
		},

		/**
		角色
		*/
		Role: ObjTypeInfoBo{
			LangCode: "ObjectType.Role",
			Name:     "角色",
		},

		/**
		权限项
		*/
		Permission: ObjTypeInfoBo{
			LangCode: "ObjectType.Permission",
			Name:     "权限项",
		},

		/**
		操作
		*/
		Operation: ObjTypeInfoBo{
			LangCode: "ObjectType.Operation",
			Name:     "操作",
		},

		/**
		权限操作
		*/
		PermissionOperation: ObjTypeInfoBo{
			LangCode: "ObjectType.PermissionOperation",
			Name:     "权限操作",
		},

		/**
		资源
		*/
		Resource: ObjTypeInfoBo{
			LangCode: "ObjectType.Resource",
			Name:     "资源",
		},
		/**
		项目
		*/
		Project: ObjTypeInfoBo{
			LangCode: "ObjectType.Project",
			Name:     "项目",
		},

		/**
		面板
		*/
		Ban: ObjTypeInfoBo{
			LangCode: "ObjectType.Ban",
			Name:     "面板",
		},

		/**
		迭代
		*/
		Iteration: ObjTypeInfoBo{
			LangCode: "ObjectType.Iteration",
			Name:     "迭代",
		},

		/**
		问题
		*/
		Issue: ObjTypeInfoBo{
			LangCode: "ObjectType.Issue",
			Name:     "问题",
		},

		/**
		优先级
		*/
		Priority: ObjTypeInfoBo{
			LangCode: "ObjectType.Priority",
			Name:     "优先级",
		},

		/**
		流程状态
		*/
		ProcessStatus: ObjTypeInfoBo{
			LangCode: "ObjectType.ProcessStatus",
			Name:     "流程状态",
		},

		/**
		动态
		*/
		Trends: ObjTypeInfoBo{
			LangCode: "ObjectType.Trends",
			Name:     "动态",
		},

		/**
		评论
		*/
		Comment: ObjTypeInfoBo{
			LangCode: "ObjectType.Comment",
			Name:     "评论",
		},

		/**
		消息设置
		*/
		MessageConfig: ObjTypeInfoBo{
			LangCode: "ObjectType.MessageConfig",
			Name:     "消息设置",
		},

		/**
		项目版本
		*/
		ProjectVersion: ObjTypeInfoBo{
			LangCode: "ObjectType.ProjectVersion",
			Name:     "项目版本",
		},

		/**
		项目模块
		*/
		ProjectModule: ObjTypeInfoBo{
			LangCode: "ObjectType.ProjectModule",
			Name:     "项目模块",
		},

		/**
		问题来源
		*/
		IssueSource: ObjTypeInfoBo{
			LangCode: "ObjectType.IssueSource",
			Name:     "问题来源",
		},

		/**
		应用
		*/
		AppInfo: ObjTypeInfoBo{
			LangCode: "ObjectType.AppInfo",
			Name:     "应用",
		},

		/**
		用户设置
		*/
		UserConfig: ObjTypeInfoBo{
			LangCode: "ObjectType.UserConfig",
			Name:     "用户设置",
		},

		/**
		项目对象类型
		*/
		ProjectObjectType: ObjTypeInfoBo{
			LangCode: "ObjectType.ProjectObjectType",
			Name:     "项目对象类型",
		},

		/**
		流程
		*/
		Process: ObjTypeInfoBo{
			LangCode: "ObjectType.Process",
			Name:     "流程",
		},

		/**
		项目类型
		*/
		ProjectType: ObjTypeInfoBo{
			LangCode: "ObjectType.ProjectType",
			Name:     "项目类型",
		},
		/**
		流程步骤
		*/
		ProcessStep: ObjTypeInfoBo{
			LangCode: "ObjectType.ProcessStep",
			Name:     "流程步骤",
		},
	}
)

type ObjTypeInfoBo struct {
	LangCode string
	Name     string
}
