package orgvo

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Inner User Infos
////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type InnerUserInfo struct {
	Id     int64  `json:"id"`
	OrgID  int64  `json:"orgId"`
	Name   string `json:"name"`
	Mail   string `json:"mail"`
	Phone  string `json:"phone"`
	Avatar string `json:"avatar"`
}

type InnerUserInfosInput struct {
	Ids interface{} `json:"ids"` // 对应无码的成员部门字段 有几种可能的结构
}

type InnerUserInfosReq struct {
	OrgId int64                `json:"orgId"`
	Input *InnerUserInfosInput `json:"input"`
}
