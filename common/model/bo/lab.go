package bo

type GetLabConfigBo struct {
	WorkBenchShow bool `json:"workBenchShow"`
	ProOverview   bool `json:"proOverview"`
	SideBarShow   bool `json:"sideBarShow"`
	EmptyApp      bool `json:"emptyApp"`
	//AutomationSwitch bool `json:"automationSwitch"`
	//SubmitButton  bool `json:"submitButton"`
}

type LabConfigBo struct {
	OrgId  int64  `json:"orgId"`
	Config string `json:"config"`
}
