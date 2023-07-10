package consts

const (
	// 1催办负责人；2催办审批
	UrgeOwner    = 1
	UrgeAuditors = 2

	DefaultTableName = "任务"

	// 群聊机器人推送配置，该项代表“任务更新”选项勾选后，勾选任务的所有列
	GroupChatAllIssueColumnFlag = "__allIssueColumnSelect"

	// 创建任务的来源，目前有用户界面创建、任务导入、任务复制等几种来源
	CreateIssueSourceNormal = "normalCreate" // 普通的用户界面创建来源
	CreateIssueSourceImport = "import"       // 导入
	CreateIssueSourceCopy   = "copy"         // 复制

	// 未分配项目的名称
	//NotAssignProjectName = "未放入项目"

	// 飞书卡片任务更新内容描述，一行不可超过 50 个字符
	GrouChatIssueChangeDescLimitPerLine = 50
	// 卡片中评论文本内容限制 150 **字节**
	GroupChatIssueCommentLimit = 150
	// 卡片中评论文本内容限制 50 **字符**
	GroupChatIssueCommentCutLimit = 50

	// 列名/别名限制6个字符
	ColumnAliasLimit = 6
	// 截断后的展示占位
	CardIssueColumnNameOverflow     = "…"      // 列名超出限制后，截断占位补全。
	CardIssueChangeDescTextOverflow = "......" // cell 值的截断占位补全。在飞书卡片中，两个点占一个中文字符
	// 列明后的符号，冒号或者不需要冒号
	CardSymbolAfterColumnName = ""
	// 更新 cell 时，如果有多个 target，拼接连接符
	CardTargetNameGlue = "、"

	CardIssueChangeDescForNotSupportShow = "*资源不支持显示*"

	FsCardNotUnderStand = "我没有理解你的意思，你可以通过以下方式与项目交互：\n"
	FsCardBindMultiPro  = "已将2个项目关联到本群聊\n" +
		" 您关联的所有项目都将推送动态至本群聊。\n" +
		"您还可以发送 @极星协作 设置 ，重新进行推送设置。"

	FsCardNoBindPro = "Hi~您好，我是极星协作的机器人。\n" +
		"我可以为您推送关注的项目动态。快来对我设置吧^ ^ ~" + ""

	FsCardHelpOfBind0ProForGroupChat = "项目相关动态会同步到这个群里，点击下面的按钮调整动态通知范围："

	FsCardProjectSettings = "项目相关的动态会同步到这个群里，点击下面的按钮调整动态通知范围："

	FsCardBotNotificationText   = "点击下面的按钮调整通知范围："
	FsCardBotNotificationButton = "通知范围设置"

	FsCardUpgradeNoticeTitle   = "你的同事申请升级/续费极星协作"
	FsCardUpgradeNoticeContent = "你的同事 %s 正在使用极星，为了使用更多功能或保证继续使用，你可以在应用管理后台升级/续费极星协作。"

	FsCardReplyIssueChatIns = "本群为特定任务群聊，暂不支持快捷操作。你是否想要查看当前记录的详情？"

	FsCardSubscribeNotice = "您已退订该消息，后续【%s】类消息将不再主动推送，若要重新订阅，可回复** %s **进入通知设置"

	FsCardHelpInfoForInstructionHelp = "@极星协作 项目任务\n" +
		"@极星协作 项目进展\n" +
		"@极星协作 设置\n" +
		"也可以 @负责人姓名 来查询其负责任务情况或指派任务给该负责人\n" +
		"@极星协作 @负责人姓名\n" +
		"@极星协作 @负责人姓名 任务标题"
	// 飞书卡片 - 个人机器人创建任务时- 用户创建任务是达到上限文案提示
	FsCardTextIssueMaxLimit = "标准版用户可创建记录上限为1000条（包括已删除记录），为不影响使用，可点击下方配置按钮前往升级版本。"

	// 工时卡片中工时时间展示的时间 format eg: 2006-01-02 15:04:05
	FsCardWorkHourTimeRangeFormat         = "01月02号"
	FsCardWorkHourTimeRangeFormatWithYear = "2006年01月02号"
	// 协作人的标识
	UserTypeOfCollaborate = "collaborate"
	// 个人推送开关配置，开关标识
	PersonalPushConfigForMyDailyReport     = "myDailyReport"     // 个人日报
	PersonalPushConfigForProDailyReport    = "proDailyReport"    // 项目日报
	PersonalPushConfigForOwnerRange        = "ownerRange"        // 我负责的
	PersonalPushConfigForCollaborateRange  = "collaborateRange"  // 我协作的
	PersonalPushConfigForIssueUpdate       = "issueUpdate"       // 任务被更新
	PersonalPushConfigForIssueBeComment    = "issueBeComment"    // 任务被评论
	PersonalPushConfigForRelateIssueTrends = "relateIssueTrends" // 关联任务动态
	PersonalPushConfigForRemindMessage     = "remindMessage"     // 任务提醒，被添加到记录

	CardDefaultOwnerNameForUpdateIssue        = "无" // 更新记录如果有 无
	CardDefaultIssueParentTitleForUpdateIssue = "无"
	CardDefaultIssueProjectName               = "未放入项目"
	CardDefaultRelationIssueTitle             = "未命名记录"

	// 卡片推送订阅开关值
	CardSubscribeSwitchOn  = 1
	CardSubscribeSwitchOff = 2

	CardButtonTextForViewInsideApp  = "应用内查看"
	CardButtonTextForViewDetail     = "查看详情"
	CardButtonTextForPersonalReport = "去处理"

	// 飞书卡片左对齐的字符数
	//CardAlineCharacter = "%-15s\t%s"
	// 卡片要推送的字段名
	CardElementTitle        = "标题"
	CardElementCode         = "编号"
	CardElementProject      = "项目"
	CardElementProjectTable = "项目/表"
	CardElementAppTable     = "应用/表"
	CardElementParent       = "父记录"
	CardElementOwner        = "负责人"
	CardElementOperator     = "操作人"
	CardElementChatTopic    = "讨论主题"
	//CardElementUrge         = "催办："
	CardElementReply                 = "你的回复:"
	CardElementUpdateInfo            = "更新详情："
	CardElementComment               = "评论:"
	CardElementEvaluate              = "评价:"
	CardElementAttachment            = "附件:"
	CardElementWithoutPermissionTips = "*注：暂无部分字段的编辑权限*"
	CardTitlePlanEndRemind           = "⏰ 即将到达截止时间"
	CardTitlePlanStartRemind         = "⏰ 任务即将开始"
	CardTitleAuditAlready            = "审批已处理，无需再操作"
	CardTitleUrgeIssue               = "👨🏻‍💻%s催办了记录"
	CardUrgePeople                   = "%s催办:"
	CardDoubleQuotationMarks         = "“%s”"
	CardReplyUrgeTitle               = "%s回复了催办"
	CardReplyBeUrged                 = "%s回复:"
	CardAuditorTitle                 = "✅ %s 提交了审批"
	CardTodoAuditTitle               = "✅ %s 请您审批"
	CardTodoFillInTitle              = "📝 %s 请您填写"
	CardIssueAddMember               = "%s 将你添加为 %s"
	CardPlanEndTime                  = "计划截止时间"
	CardWorkHourAddDesc              = "添加了%s：%s %s小时 %s"
	CardWorkHourDelDesc              = "删除了%s：~~%s %s小时 %s~~"
	CardWorkHourModify               = "修改了%s：~~%s~~ **→** %s"
	CardDocAndImageAddDesc           = "添加了 “%s”"
	CardDocAndImageDelDesc           = "删除了 ~~“%s”~~"
	CardRelatingAddDesc              = "添加了：%s"
	CardRelatingDelDesc              = "删除了：~~%s~~"
	CardBaRelatingLinkToAddDesc      = "添加了（前置）：%s"
	CardBaRelatingLinkToDelDesc      = "删除了（前置）：~~%s~~"
	CardBaRelatingLinkFromAddDesc    = "添加了（后置）：%s"
	CardBaRelatingLinkFromDelDesc    = "删除了（后置）：~~%s~~"
	CardCommonChangeDesc             = "~~%s~~ → %s"
	CardCompleteIssue                = "完成记录"
	CardShare                        = " 分享了记录"
	CardUnSubscribe                  = "点击退订"
	CardReplySelectFirst             = "现在忙，晚点处理"
	CardAuditPass                    = "通过"
	CardAuditReject                  = "驳回"
	CardReplyUrgeSelf                = "你回复了催办"
	CardModifyDeadline               = "修改截止时间"
	CardCompleteTask                 = "完成任务"
	CardIssueCreateErr               = "任务创建失败"
	CardProTrendsSettings            = "项目动态设置"
	CardCronTrigger                  = "工作流【%s】"
	CardDataEvent                    = "%s 触发工作流【%s】"

	// 飞书卡片标题颜色
	// https://open.feishu.cn/document/ukTMukTMukTM/ukTNwUjL5UDM14SO1ATN
	FsCardTitleColorBlue   = "blue"
	FsCardTitleColorOrange = "orange"

	// 飞书按钮颜色
	FsCardButtonColorDefault = "default"
	FsCardButtonColorPrimary = "primary"
	FsCardButtonColorDanger  = "danger"

	// 企微卡片标题颜色
	WxCardTitleColorInfo    = "info"
	WxCardTitleColorComment = "comment"
	WxCardTitleColorWarning = "warning"

	FsCardButtonTextForViewInsideApp = "应用内查看"
	FsCardButtonTextForViewDetail    = "查看详情"
	FsCardButtonTextForGotoSetting   = "前往配置"

	CardTablePro = "%s/%s"

	// 制表符
	FsCard1Tab           = "\t"
	FsCard2Tab           = "\t\t"
	FsCard3Tab           = "\t\t\t"
	FsCard4Tab           = "\t\t\t\t"
	FsCardCommentPattern = "<at id=\\w+></at>"

	MarkdownBr         = "\n"
	MarkdownBold       = "**%s**"
	MarkdownLink       = "[%s](%s)"
	MarkdownTitle      = "### %s"
	MarkdownColorTitle = "### <font color=\"%s\">%s</font>"
	MarkdownColor      = "<font color=\"%s\">%s</font>"
	MarkDownDel        = "~~"

	CardPaddingChBlank = "　"
	CardPaddingEnBlank = " "
	CardKeyChLength    = 7
	CardKeyEnLength    = 3

	CardLevelInfo    = 0
	CardLevelWarning = 1

	// 飞书卡片按钮，链接打开方式
	// 值详情参考 https://open.feishu.cn/document/uAjLw4CM/uYjL24iN/applink-protocol/supported-protocol/open-a-gadget
	FsCardAppLinkOpenTypeSideBar   = "sidebar-semi" // 侧边栏打开
	FsCardAppLinkOpenTypeAppCenter = "appCenter"    // 工作台打开
)

const (
	DailyProjectReportTitle  = "项目日报"
	DailyFinishCountTitle    = "今日完成"
	DailyRemainingCountTitle = "剩余未完成"

	DailyPersonalIssueReportTitle = "个人日报"
	DailyDueOfTodayTitle          = "我的今日截止"
	BeAboutToOverdueTitle         = "我的即将逾期"
	DailyOverdueCountTitle        = "我的已逾期"
	PendingTitle                  = "我的待完成"
)

func GetTabCharacter(str string) string {
	i := len([]rune(str))
	if i == 7 {
		// 返回一个\t
		return FsCard1Tab
	}
	if i >= 5 && i < 7 {
		// 返回2个\t
		return FsCard2Tab
	}
	switch i {
	case 3:
		return FsCard3Tab
	case 4:
		return FsCard3Tab
	}
	if i <= 2 {
		// 返回4个\t
		return FsCard4Tab
	}
	return FsCard1Tab
}

func GetUrgeMsgOptions() []string {
	return []string{"现在忙，晚点处理", "抱歉，现在不方便处理", "暂不处理", "忽略"}
}
