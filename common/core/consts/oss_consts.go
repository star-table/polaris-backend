package consts

//oss callback 请求类型
const (
	//JSON
	OssCallbackBodyTypeApplicationJson = "application/json"
)

//oss policy类型
const (
	//项目封面
	OssPolicyTypeProjectCover = 1
	//任务附件
	OssPolicyTypeIssueResource = 2
	//导入文件
	OssPolicyTypeIssueInputFile = 3
	//项目资源
	OssPolicyTypeProjectResource = 4
	//兼容测试
	OssPolicyTypeCompatTest = 5
	//头像
	OssPolicyTypeUserAvatar = 6
	//反馈资源
	OssPolicyTypeFeedback = 7
	//任务备注
	OssPolicyTypeIssueRemark = 8
	//任务附件（无码）
	OssPolicyTypeLesscodeResource = 9

	// excel 导入成员
	OssPolicyTypeImportMembers = 11

	// 评论添加附件
	OssPolicyTypeCommentAttachments = 12
)

//oss key segment
const (
	OssKeySegmentOrg     = "org_"
	OssKeySegmentProject = "project_"
	OssKeySegmentIssue   = "issue_"
)
