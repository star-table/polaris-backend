package consts

import (
	"github.com/star-table/polaris-backend/common/core/consts"
)

var (
	//用户配置缓存
	CacheBaseProjectInfo = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "baseinfo"
	//用户基础信息缓存key
	CacheProjectObjectTypeList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "project_object_type"
	//优先级列表
	CachePriorityList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "priority_list"
	//项目类型
	CacheProjectTypeList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "project_type_list"
	//项目类型与项目对象类型关联缓存
	CacheProjectTypeProjectObjectTypeList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "project_type_project_object_type"
	//项目状态缓存
	CacheProjectProcessStatus = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "process_status"
	//任务来源列表
	CacheIssueSourceList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "issue_source_list"
	//任务性质列表
	CacheIssuePropertyList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "issue_property_list"
	//任务类型列表
	CacheIssueObjectTypeList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "issue_object_type_list"
	//项目日历信息缓存
	CacheProjectCalendarInfo = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "project_calendar_info"
	//项目标签缓存
	CacheProjectTagInfo = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "tag_stat_info"
	//项目主群聊id缓存
	CacheProjectMainChatId = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + "chat_id"
	//项目类型分类缓存
	CacheProjectTypeCategoryList = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "project_type_category_list"
	//项目对象类型关联流程缓存
	CacheProjectObjectTypeProcessHash = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + "project_object_type_process_hash"
	//飞书群聊项目设置
	CacheFsPushSettings = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.
				CacheKeyOfProject + "push_settings_2022"
	//任务审批催办缓存
	CacheUrgeIssueAudit = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + consts.CacheKeyOfIssue + "urge_issue_audit"
	//任务催办缓存
	CacheUrgeIssue = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + consts.CacheKeyOfIssue + "urge_issue"

	// 项目顶部菜单列表缓存key
	CacheProjectMenuConfig = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfApp + "project_menu_config"

	// 缓存记录 任务数目
	CacheDeletedIssueNums = consts.CacheKeyPrefix + consts.ProjectsvcApplicationName + consts.CacheKeyOfSys + "deleted_issue_num"
)
