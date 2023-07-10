package appauth

import (
	"github.com/star-table/common/core/util/slice"
)

const (
	Read                   = 1 //可见
	Write                  = 2 //可编辑
	Masking                = 4 //脱敏
	ReadAndWrite           = 3 //可见，可编辑
	ReadAndMasking         = 5 //可见，脱敏
	WriteAndMasking        = 6 //可编辑，脱敏
	ReadAndWriteAndMasking = 7 //可见，可编辑，脱敏
)

type GetAppAuthData struct {
	AppAuth              []string                  `json:"appAuth"`
	LangCode             *string                   `json:"langCode"`
	AppId                int64                     `json:"appId,string"`
	AppOwner             bool                      `json:"appOwner"`
	OrgOwner             bool                      `json:"orgOwner"`
	SysAdmin             bool                      `json:"sysAdmin"`
	HasAppRootPermission bool                      `json:"hasAppRootPermission"`
	FieldAuth            map[string]map[string]int `json:"fieldAuth"`
	OptAuth              []string                  `json:"optAuth"`
	TableAuth            []string                  `json:"tableAuth"`
}

func (auth GetAppAuthData) HasOrgRootPermission() bool {
	return auth.SysAdmin || auth.OrgOwner
}
func (auth GetAppAuthData) IsAdmin() bool {
	return auth.SysAdmin
}

func (auth GetAppAuthData) HasAllFieldAuthOfTable(tableIdStr string) bool {
	if auth.HasAppRootPermission {
		return true
	}
	if auth.FieldAuth == nil || len(auth.FieldAuth) == 0 {
		return true
	}

	var fa map[string]int
	var ok bool
	if fa, ok = auth.FieldAuth[tableIdStr]; !ok {
		return true
	}

	// 如果字段权限为空，则拥有所有权限
	if fa == nil || len(fa) == 0 {
		return true
	}

	return false
}

//查看权限
func (auth GetAppAuthData) HasFieldViewAuth(tableIdStr string, field string) bool {
	if auth.HasAppRootPermission {
		return true
	}
	if auth.FieldAuth == nil || len(auth.FieldAuth) == 0 {
		return true
	}

	var fa map[string]int
	var ok bool
	if fa, ok = auth.FieldAuth[tableIdStr]; !ok {
		return true
	}

	// 如果字段权限为空，则拥有所有权限
	if fa == nil || len(fa) == 0 {
		return true
	}

	allRead, ok := fa["hasRead"]
	if ok && allRead == 1 {
		return true
	}

	//任务栏就始终展示吧
	if ok, _ := slice.Contain([]string{"id", "issueId", "creator", "updator", "createTime", "updateTime", "appId", "projectId", "auditStatus", "parentId", "auditorIds", "projectObjectTypeId"}, field); ok {
		return true
	}
	if code, ok := fa[field]; ok {
		return code&1 == 1
	}
	return false
}

//编辑权限
func (auth GetAppAuthData) HasFieldWriteAuth(tableIdStr string, field string) bool {
	if auth.HasAppRootPermission {
		return true
	}
	if auth.FieldAuth == nil || len(auth.FieldAuth) == 0 {
		return true
	}

	var fa map[string]int
	var ok bool
	if fa, ok = auth.FieldAuth[tableIdStr]; !ok {
		return true
	}

	// 如果字段权限为空，则拥有所有权限
	if fa == nil || len(fa) == 0 {
		return true
	}

	allWrite, ok := fa["hasWrite"]
	if ok && allWrite == 1 {
		return true
	}

	if code, ok := fa[field]; ok {
		return code>>1&1 == 1
	}
	return false
}
