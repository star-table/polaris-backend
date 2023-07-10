package routes

import (
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	_ "github.com/star-table/polaris-backend/app/docs"
	"github.com/star-table/polaris-backend/app/openapi"
	openapiV2 "github.com/star-table/polaris-backend/app/openapi/v2"
	"github.com/star-table/polaris-backend/app/server/handler"
	"github.com/star-table/polaris-backend/app/server/handler/lab_handler"
	"github.com/star-table/polaris-backend/app/server/handler/org_handler"
	"github.com/star-table/polaris-backend/app/server/handler/other_handler"
	"github.com/star-table/polaris-backend/app/server/handler/project_handler"
	"github.com/star-table/polaris-backend/app/server/handler/user_handler"
	"github.com/star-table/polaris-backend/app/server/middlewares"
	"github.com/star-table/polaris-backend/common/extra/gin/mid"
	"github.com/star-table/polaris-backend/common/extra/trace/gin2micro"
	trace "github.com/star-table/polaris-backend/common/extra/trace/jaeger"
	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	log = logger.GetDefaultLogger()
)

func SetMiddlewares(r *gin.Engine, env string) {
	applicationName := config.GetApplication().Name
	sentryConfig := config.GetSentryConfig()
	sentryDsn := ""
	if sentryConfig != nil {
		sentryDsn = sentryConfig.Dsn
	}

	if config.GetJaegerConfig() != nil {
		t, io, err := trace.NewTracer(config.GetJaegerConfig())
		if err != nil {
			log.Infof("err %v", err)
		}
		defer func() {
			if err := io.Close(); err != nil {
				log.Errorf("err %v", err)
			}
		}()
		opentracing.SetGlobalTracer(t)

		r.Use(gin2micro.TracerWrapper)
	}

	r.Use(mid.SentryMiddleware(applicationName, env, sentryDsn))
	r.Use(mid.StartTrace())
	r.Use(mid.GinContextToContextMiddleware())
	r.Use(mid.CorsMiddleware())
	r.Use(mid.AuthMiddleware())
}

func SetRoutes(r *gin.Engine) {
	//1.获取验证码
	r.POST("/api/task/captcha", handler.CaptchaGetHandler())

	//2.获取验证码图片
	r.GET("/api/task/captcha/:captchaId", handler.CaptchaShowHandler())

	r.POST("/task", handler.GraphqlHandler())
	r.GET("/task/health", handler.HeartbeatHandler())
	r.GET("/", handler.PlaygroundHandler())
	r.POST("/callback", handler.DingTalkCallBackHandler())

	r.POST("/api/task", handler.GraphqlHandler())
	r.GET("/api/task/health", handler.HeartbeatHandler())
	r.GET("/admin/pg", handler.PlaygroundHandler())
	r.POST("/api/callback", handler.DingTalkCallBackHandler())

	r.POST("/api/task/upload", handler.FileUploadHandler())
	r.POST("/api/task/uploadUpdateFile", handler.FileUploadAndUpdateResourceInfoHandler())

	r.GET("/api/task/read/*path", handler.FileReadHandler())
	r.POST("/api/task/importData", handler.ImportDataHandler())

	// 无需路由前缀、无需登录的路由
	r.GET("/api/rest/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// restful 接口路由定义
	polarisRestful := r.Group("/api/rest", middlewares.PolarisRestAuth)
	{
		// 任务视图
		polarisRestful.POST("/issueView/create", handler.IssueViewHandles.CreateTaskView)
		polarisRestful.PUT("/issueView/:viewId", handler.IssueViewHandles.UpdateTaskView)
		polarisRestful.POST("/issueView/list", handler.IssueViewHandles.GetTaskViewList)
		polarisRestful.DELETE("/issueView/:viewId", handler.IssueViewHandles.DeleteTaskView)
		//任务列表
		polarisRestful.POST("/issue/filter", handler.IssueListHandler.Filter)
		// 给任务添加/更新所属项目
		//polarisRestful.POST("/issue/add-belong-project", handler.IssueListHandler.AddBelongProject)
		//polarisRestful.POST("/issue/update-belong-project", handler.IssueListHandler.UpdateBelongProject)

		// 筛选时 获取可选值列表
		polarisRestful.POST("/filter/optionList", handler.IssueViewHandles.GetOptionList)

		// 项目
		// 项目成员列表
		polarisRestful.POST("/project/projectUserList", project_handler.ProjectHandler.GetProjectUserList)
		//表头编辑
		//polarisRestful.POST("/project/:projectId/saveForm", project_handler.ProjectHandler.SaveFormHeader)
		//获取表头
		//polarisRestful.GET("/project/:projectId/config", project_handler.ProjectHandler.GetFormConfig)
		//获取表头（批量）
		//polarisRestful.POST("/project/configs", project_handler.ProjectHandler.GetFormConfigBatch)

		//任务（表单数据）
		//任务列表
		polarisRestful.POST("/project/:projectId/issue/filter", handler.IssueListHandler.Filter)
		//任务详情
		polarisRestful.POST("/project/:projectId/issue/filterForIssue", handler.IssueListHandler.FilterForIssue)
		//创建
		polarisRestful.POST("/project/:projectId/values", handler.IssueListHandler.CreateIssueNew)
		//更新
		//polarisRestful.PUT("/project/:projectId/values", handler.IssueListHandler.UpdateIssue)
		polarisRestful.PUT("/project/:projectId/values", handler.IssueListHandler.UpdateIssueNew)
		//批量更新
		polarisRestful.POST("/project/:projectId/batchValues", handler.IssueListHandler.BatchUpdateIssue)
		//批量审批/驳回
		polarisRestful.POST("/project/:projectId/batchAudit", handler.IssueListHandler.BatchAuditIssue)
		//批量催办
		polarisRestful.POST("/project/:projectId/batchUrge", handler.IssueListHandler.BatchUrgeIssue)
		//删除
		polarisRestful.POST("/project/:projectId/values/delete", handler.IssueListHandler.DeleteIssue)
		// 复制
		polarisRestful.POST("/project/:projectId/values/copy", handler.IssueListHandler.CopyIssue)
		// 发起任务群聊讨论，如果群不存在，则创建之
		polarisRestful.POST("/project/issue/startChat", handler.IssueListHandler.StartChat)
		// 导出同名的用户和部门列表 excel
		polarisRestful.POST("/project/export-same-user-dept", project_handler.ProjectHandler.ExportSameNameUserDept)
		// 查询一些异步任务的进度/进度条
		polarisRestful.POST("/project/task/queryProcess", project_handler.ProjectHandler.QueryProcessForAsyncTask)
		// 检查项目页是否展示项目群聊 icon
		polarisRestful.POST("/project/checkIsShowProChatIcon", project_handler.ProjectHandler.CheckIsShowProChatIcon)
		// 组织下项目等资源统计信息，用于未付费资源限制
		polarisRestful.POST("/project/payLimitNum", project_handler.ProjectHandler.PayLimitNum)

		// 自动化工作流
		//polarisRestful.POST("/automation/workflows", handler.AutomationHandler.GetWorkflowsByApp)
		//polarisRestful.POST("/automation/workflow", handler.AutomationHandler.GetWorkflow)
		//polarisRestful.POST("/automation/createWorkflow", handler.AutomationHandler.CreateWorkflow)
		//polarisRestful.POST("/automation/updateWorkflow", handler.AutomationHandler.UpdateWorkflow)
		//polarisRestful.POST("/automation/deleteWorkflow", handler.AutomationHandler.DeleteWorkflow)
		//polarisRestful.POST("/automation/activateWorkflow", handler.AutomationHandler.ActivateWorkflow)
		//polarisRestful.POST("/automation/deactivateWorkflow", handler.AutomationHandler.DeactiveWorkflow)
		//polarisRestful.POST("/automation/executions", handler.AutomationHandler.GetExecutions)
		//polarisRestful.POST("/automation/execution", handler.AutomationHandler.GetExecution)

		// 待办
		//polarisRestful.POST("/todo/filter", handler.TodoHandler.Filter)
		polarisRestful.POST("/todo/urge", handler.TodoHandler.Urge)

		// 个人任务统计
		polarisRestful.GET("/issue/stat", handler.IssueListHandler.IssueStat)

		// 分享卡片
		//polarisRestful.POST("/issue/share", handler.IssueListHandler.IssueCardShare)

		// 根据DataId查询任务基本信息
		polarisRestful.POST("/project/issues/filterSimpleByDataIds", handler.IssueListHandler.GetIssueListSimpleByDataIds)
		// 根据TableId查询任务基本信息
		polarisRestful.POST("/project/issues/filterSimpleByTableIds", handler.IssueListHandler.GetIssueListSimpleByTableIds)

		// 获取项目成员id列表
		polarisRestful.POST("/project/getProjectMemberIdList", project_handler.ProjectHandler.GetProjectMemberIds)

		// 成员
		// 邀请-通过多个手机号邀请成员
		polarisRestful.POST("/users/invite/with-phones", org_handler.InviteWithPhones)
		// 邀请-下载邀请成员的 excel 模板
		polarisRestful.POST("/users/invite/export-template", org_handler.ExportInviteTemplate)
		// 邀请-excel 数据导入成员
		polarisRestful.POST("/users/invite/import-members", org_handler.ImportMembers)
		// 设定用户已经浏览过“用户指引”
		polarisRestful.POST("/user/setUserGuideStatus", user_handler.UserHandler.SetVisitUserGuideStatus)

		//视图统计
		polarisRestful.POST("/view/stat", handler.IssueListHandler.ViewStat)

		//镜像视图统计
		polarisRestful.POST("/mirrors/stat", handler.IssueListHandler.MirrorsStat)

		// 实验室开关
		polarisRestful.POST("/lab/config", lab_handler.LabHandler.SetLab)
		polarisRestful.GET("/lab/config", lab_handler.LabHandler.GetLab)

		// 获取组织配置信息
		polarisRestful.POST("/org/getOrgConfig", org_handler.OrgHandler.GetOrgConfig)

		// 获取用户信息
		polarisRestful.POST("/user/getUserInfo", user_handler.UserHandler.GetUserInfo)
		// 设定用户已经浏览过版本信息了
		polarisRestful.POST("/user/setUserVersionVisible", user_handler.UserHandler.SetVersionVisible)
		polarisRestful.POST("/user/setUserActivity", user_handler.UserHandler.SetUserActivity)

		// 保存用户上一次浏览视图的位置
		polarisRestful.POST("/user/setUserViewLocation", user_handler.UserHandler.SetUserViewLocation)
		// 获取用户上一次浏览视图的位置
		polarisRestful.GET("/user/userViewLocation", user_handler.UserHandler.GetUserViewLocation)
		// 获取同名用户列表
		polarisRestful.POST("/user/getSameNameUserOrDept", user_handler.UserHandler.GetSameNameUserOrDept)

		// 组织转换相关
		//polarisRestful.POST("/org/transferOrg", org_handler.OrgHandler.TransferOrg)
		//polarisRestful.POST("/org/prepareTransferOrg", org_handler.OrgHandler.PrepareTransferOrg)

		// 调用无码组织字段, 包装一下
		polarisRestful.POST("/read/org/columns", org_handler.ColumnHandler.GetOrgColumns)
		polarisRestful.POST("/org/column/create", org_handler.ColumnHandler.CreateOrgColumn)
		polarisRestful.POST("/org/column/delete", org_handler.ColumnHandler.DeleteOrgColumn)

		// 给组织管理员发送卡片提醒，有成员申请升级极星
		polarisRestful.POST("/org/sendCardToAdminForUpgrade", org_handler.OrgHandler.SendCardToAdminForUpgrade)

		// 普通表头字段操作，代替saveForm
		polarisRestful.POST("/column/create", project_handler.ProjectHandler.CreateColumn)
		polarisRestful.POST("/column/copy", project_handler.ProjectHandler.CopyColumn)
		polarisRestful.POST("/column/update", project_handler.ProjectHandler.UpdateColumn)
		polarisRestful.POST("/column/description/update", project_handler.ProjectHandler.UpdateColumnDescription)
		polarisRestful.POST("/column/delete", project_handler.ProjectHandler.DeleteColumn)

		// 表相关的
		polarisRestful.POST("/table/create", project_handler.ProjectHandler.CreateTable)
		polarisRestful.POST("/table/rename", project_handler.ProjectHandler.RenameTable)
		polarisRestful.POST("/table/delete", project_handler.ProjectHandler.DeleteTable)
		polarisRestful.POST("/read/table", project_handler.ProjectHandler.GetOneTable)
		polarisRestful.POST("/read/tables", project_handler.ProjectHandler.GetTableList)
		polarisRestful.POST("/read/tablesByApps", project_handler.ProjectHandler.GetTablesByApps)
		//polarisRestful.POST("/read/filterTablesByOrg", project_handler.ProjectHandler.GetTableListByOrg)
		// 表-自动排期
		polarisRestful.POST("/table/autoSchedule", project_handler.ProjectHandler.SetTableAutoSchedule)

		// 获取表头
		polarisRestful.POST("/project/:projectId/columns", project_handler.ProjectHandler.GetTableColumns)
		polarisRestful.POST("/project/tables/columns", project_handler.ProjectHandler.GetTablesColumns)

		// 菜单相关的接口
		polarisRestful.POST("/project/menu/getConfig", project_handler.ProjectHandler.GetMenu)
		polarisRestful.POST("/project/menu", project_handler.ProjectHandler.SaveMenu)

		// 获取协作人
		polarisRestful.POST("/collaborator/collaborators", project_handler.ProjectHandler.GetCollaborators)

		// 在线office接口
		polarisRestful.GET("/office/config", handler.GetOfficeConfigHandler())

		polarisRestful.POST("/sdk/dingJsApiSign", org_handler.DingTalkHandlers.GetJsApiSign)
		// 企业微信
		polarisRestful.POST("/sdk/weiXinJsApiSign", org_handler.WeiXinHandlers.GetJsApiSign)

		// 给动态列表选择项目成员用的
		polarisRestful.POST("/project/filterForTrends", project_handler.ProjectHandler.FilterForTrendsMembers)

		// 项目统计信息
		polarisRestful.POST("/project/statistics", project_handler.ProjectHandler.FilterProjectStatistics)
	}

	//开放api
	open := r.Group("/openapi")
	{
		// 项目
		//open.GET("/projects/:projectId/stats", openapi.IssueStatusTypeStat)
		open.POST("/projects", openapi.CreateProject)
		open.POST("/projects/filter", openapi.Projects)
		open.GET("/projects/:projectId", openapi.ProjectInfo)
		open.PUT("/projects/:projectId", openapi.UpdateProject)
		open.DELETE("/projects/:projectId", openapi.DeleteProject)
		// 一些属性列表查询
		// 项目的优先级列表
		open.GET("/projects/:projectId/priorities", openapi.GetPriorityList)
		// 项目的任务（对象）类型。如：需求、任务、缺陷
		open.GET("/projects/:projectId/object-types", openapi.GetProjectObjectTypeList)
		// 项目的迭代列表
		open.GET("/projects/:projectId/iterations", openapi.GetIterationList)
		// 获取任务需求来源列表
		open.GET("/projects/:projectId/issue-sources", openapi.GetIssueSourceList)
		// 获取缺陷严重程度列表
		open.GET("/projects/:projectId/properties", openapi.GetPropertyList)

		// 任务
		open.POST("/issues", openapi.CreateIssue)
		open.PUT("/issues/:issueId", openapi.UpdateIssue)
		//open.GET("/issues/:issueId", openapi.IssueInfo)
		open.POST("/issues/filter", openapi.IssueList)
		open.DELETE("/issues/:issueId", openapi.DeleteIssue)
		// 评论任务
		open.POST("/issues/:issueId/comment", openapi.CreateIssueComment)
		//镜像视图统计
		//open.POST("/mirrors/stat", openapi.MirrorsStat)
		// 任务评论列表
		open.POST("/issues/:issueId/comment/filter", openapi.IssueCommentsList)

		// 通讯录
		open.GET("/users", openapi.UserList)
	}

	//v2新增接口
	openV2 := r.Group("/openapi/v2")
	{
		openV2.POST("/issues", openapiV2.CreateIssue)

		openV2.POST("/issues/filter", openapiV2.IssueList)

		openV2.GET("/issues/:issueId", openapiV2.IssueInfo)
		openV2.PUT("/issues/:issueId", openapiV2.UpdateIssue)

		// 表接口
		// 通过AppId获取Tables
		openV2.POST("/apps/:appId", openapi.TablesByAppId)

		// 通过TableId获取表头字段
		openV2.POST("/meta/columns", openapi.TablesBaseColumns)

		// 一些属性列表查询

		// 获取任务需求来源列表
		openV2.GET("/projects/:projectId/issue-sources", openapiV2.GetIssueSourceList)

		// 项目的任务状态列表
		openV2.POST("/projects/issue-statuses", openapi.GetIssueStatusList)
	}

	// **无需登录**的极星接口，一些特殊的用处。谨慎调用。
	bjxNoLoginRest := r.Group("/api/rest")
	{
		bjxNoLoginRest.POST("/schedule/add-script-task", other_handler.OtherHandler.AddScriptTask)
		// 钉钉登陆相关
		bjxNoLoginRest.POST("/login/dingAuthCode", org_handler.DingTalkHandlers.AuthCode)

		bjxNoLoginRest.POST("/login/weiXinAuthCode", org_handler.WeiXinHandlers.AuthCode)
		bjxNoLoginRest.POST("/login/getAgentId", org_handler.WeiXinHandlers.GetAgentId)
		bjxNoLoginRest.GET("/redirect/install", org_handler.WeiXinHandlers.RedirectInstall)
		bjxNoLoginRest.GET("/redirect/register", org_handler.WeiXinHandlers.RedirectRegister)

		bjxNoLoginRest.POST("/login/personWeiXinQrCode", org_handler.WeiXinHandlers.GetQrCode)
		bjxNoLoginRest.POST("/login/checkPersonWeiXinQrCode", org_handler.WeiXinHandlers.CheckQrCodeScan)
		bjxNoLoginRest.POST("/login/personWeiXinBind", org_handler.WeiXinHandlers.PersonWeiXinBind)
		bjxNoLoginRest.POST("/login/checkMobileHasBind", org_handler.WeiXinHandlers.CheckMobileHasBind)

		//bjxNoLoginRest.POST("/login/exchangeShareToken", org_handler.OrgHandler.ExchangeShareToken)
		//bjxNoLoginRest.POST("/projectsvc/getShareViewInfoByKey", project_handler.ProjectHandler.GetShareViewInfoByKey)
	}

	// 取代graphql的接口
	newApi := r.Group("/newApi/rest", middlewares.PolarisRestAuth)
	{
		newApi.Any("/*any", handler.CommonHandler.Handler)
	}

}
