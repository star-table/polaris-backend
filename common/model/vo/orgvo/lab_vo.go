package orgvo

import "github.com/star-table/polaris-backend/common/model/vo"

type SetLabReqVo struct {
	OrgId  int64     `json:"orgId"`
	UserId int64     `json:"userId"`
	Input  SetLabReq `json:"input"`
}

type SetLabReq struct {
	WorkBenchShow bool `json:"workBenchShow"`
	ProOverview   bool `json:"proOverview"`
	SideBarShow   bool `json:"sideBarShow"`
	EmptyApp      bool `json:"emptyApp"`
	//AutomationSwitch bool `json:"automationSwitch"`
	//SubmitButton  bool `json:"submitButton"`
}

type SetLabRespVo struct {
	vo.Err
	Data *SetLabResp `json:"data"`
}

type SetLabResp struct {
	OrgId int64 `json:"orgId"`
}

type GetLabReqVo struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetLabRespVo struct {
	vo.Err
	Data *GetLabResp `json:"data"`
}

type GetLabResp struct {
	//WorkBenchShow    bool `json:"workBenchShow"`
	ProOverview bool `json:"proOverview"`
	SideBarShow bool `json:"sideBarShow"`
	EmptyApp    bool `json:"emptyApp"`
	//AutomationSwitch bool `json:"automationSwitch"`
	//SubmitButton  bool `json:"submitButton"`
}
