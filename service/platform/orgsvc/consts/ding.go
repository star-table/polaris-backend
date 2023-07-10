package consts

const (
	DingConsole       = "dingtalk://dingtalkclient/action/openapp?corpid=%s&container_type=work_platform&app_id=%v&redirect_type=jump&redirect_url=%s" // 控制台链接
	DingPlatform      = "dingtalk://dingtalkclient/action/open_platform_link?pcLink=%s&mobileLink=%s"                                                  // 双平台半屏打开链接
	DingMobileSide    = "dingtalk://dingtalkclient/action/im_open_hybrid_panel?panelHeight=percent83&hybridType=online&pageUrl=%s"                     // 手机半浮层打开
	DingPcSide        = "dingtalk://dingtalkclient/page/link?pc_slide=true&url=%s"                                                                     // pc侧边栏打开
	BindProjectPath   = "/dd/bindProject?openConversationId=%s&corpId=%s"                                                                              // 绑定项目
	HelpPath          = "https://alidocs.dingtalk.com/i/p/Db6Vz68WpPD1mZ9Nw0oXWkk68PkbMX1J"                                                            // 帮助
	ProjectDetailPath = "/project/%v/task/0?projectTypeId=%v&appId=%v"                                                                                 // 项目详情
	CreateIssuePath   = "/dd/createTask?openConversationId=%s&corpId=%s"                                                                               // 创建任务
	LoginUrl          = "/dd/auth?redirect_url=%v"
)

const (
	FirstCardTemplateId = "17675c8f-fea6-48ca-8eb9-a08963e2fc50"
	TopCardTemplateId   = "1e600bc3-1a1c-4d86-ac40-031161290b6a"
)

const (
	CallBackKeyRefresh = "topRefresh"
	RefreshUrlPath     = "callback/dingtalk/card"
)

const (
	BindTitle        = "绑定项目"
	CreateIssueTitle = "创建任务"
)

const (
	ConversationTypeGroup = 1
)
