package appauth

import "testing"

func TestGetAppAuthData_IsAdmin(t *testing.T) {
	a := GetAppAuthData{
		AppAuth:              nil,
		LangCode:             nil,
		AppId:                0,
		AppOwner:             false,
		HasAppRootPermission: false,
		OptAuth:              nil,
		OrgOwner:             false,
		SysAdmin:             true,
	}
	t.Log(a.IsAdmin())
}