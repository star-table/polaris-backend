package bo

type InitOrgBo struct {
	//外部组织id
	OutOrgId string `json:"outOrgId"`
	//组织负责人外部id
	OutOrgOwnerId string `json:"outOrgOwnerId"`
	//组织名称
	OrgName string `json:"orgName"`
	//企业编号
	TenantCode string `json:"tenantCode"`
	//组织logo
	OrgLogo string `json:"orgLogo"`
	//行业
	Industry string `json:"industry"`
	//是否认证
	IsAuthenticated bool `json:"isAuthenticated"`
	//认证等级
	AuthLevel int `json:"authLevel"`
	//省份
	OrgProvince string `json:"corpProvince"`
	//城市
	OrgCity string `json:"corpCity"`
	//来源
	SourceChannel string `json:"sourceChannel"`
	//永久授权码 dingding的时候才会有
	PermanentCode string `json:"permanentCode"`
}
