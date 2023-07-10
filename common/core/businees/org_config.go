package businees

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
)

// OrgInitConfig 极星中的每个组织都需要进行特定的初始化，其中包括：创建汇总表、空项目应用、文件夹、项目数据表单（这些都是一种应用）等，
// 应用对应的唯一 id（appId）皆以 json 的方式存储在 organization 表中的 remark 中。`ConfigObj` 中是各个应用的 appId 值。
type OrgInitConfig struct {
	ConfigObj *orgvo.OrgRemarkConfigType
}

func NewOrgConfig(orgRemarkJson string) (*OrgInitConfig, errs.SystemErrorInfo) {
	orgConfig := OrgInitConfig{}
	orgRemarkObj := orgvo.OrgRemarkConfigType{}
	oriErr := json.FromJson(orgRemarkJson, &orgRemarkObj)
	if oriErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.OrgInitError, oriErr)
	}
	orgConfig.ConfigObj = &orgRemarkObj

	return &orgConfig, nil
}

func (oc *OrgInitConfig) GetConfigObj() *orgvo.OrgRemarkConfigType {
	return oc.ConfigObj
}

// GetSummaryAppId 获取汇总表 appId
func (oc *OrgInitConfig) GetSummaryAppId() int64 {
	return oc.ConfigObj.OrgSummaryTableAppId
}

// GetEmptyProAppId 获取“空项目” appId
func (oc *OrgInitConfig) GetEmptyProAppId() int64 {
	return oc.ConfigObj.EmptyProjectAppId
}

// GetProFormAppId 获取存储项目数据的 form 的 appId
func (oc *OrgInitConfig) GetProFormAppId() int64 {
	return oc.ConfigObj.ProjectFormAppId
}

// CheckIsPaidVer 检查是否是付费版本（极星）
func CheckIsPaidVer(payLevel int) bool {
	return payLevel == consts.PayLevelFlagship || payLevel == consts.PayLevelEnterprise ||
		payLevel == consts.PayLevelDouble11Activity
}

// GetOrgPayVersionName 获取组织的付费版本
func GetOrgPayVersionName(level int) string {
	defaultDesc := "未知版本"
	switch level {
	case consts.PayLevelStandard, consts.PayLevelEnterprise, consts.PayLevelFlagship, consts.PayLevelDouble11Activity:
		return consts.PayLevelDescMap[level]
	default:
		return defaultDesc
	}
}

func GetMenuSettings() []projectvo.PolarisMenuOption {
	options := []projectvo.PolarisMenuOption{
		{
			Name:    "设置",
			Icon:    "set_ico,set_ico_light",
			LinkUrl: "/setting/organization/basic",
		},
		{
			Name:    "回收站",
			Icon:    "dustbin,dustbin_light",
			LinkUrl: "/trash/:appId",
		},
		{
			Name:    "工时",
			Icon:    "panel_ico,panel_ico_light",
			LinkUrl: "/statistical",
		},
		{
			Name:    "动态",
			Icon:    "dynamic_ico,dynamic_ico_light",
			LinkUrl: "/dynamic",
		},
		{
			Name:    "成员",
			Icon:    "member_ico,member_ico_light",
			LinkUrl: "/address/org/members",
		},
		{
			Name:    "极星模板",
			Icon:    "template_icon,template_icon_light",
			LinkUrl: "/template",
		},
		{
			Name:    "项目",
			Icon:    "project_ico,project_ico_light",
			LinkUrl: "/project",
		},
		{
			Name:    "任务",
			Icon:    "task_ico,task_ico_light",
			LinkUrl: "/task",
		},
		{
			Name:    "工作台",
			Icon:    "work_ico,work_ico_light",
			LinkUrl: "/overview",
		},
	}
	return options
}
