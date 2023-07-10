package bo

type IdCodeInfo struct {
	Id   int64  `json:"id"`
	Code string `json:"code"`
}

type IdCodes struct {
	OrgId   int64  `json:"orgId"`
	Code    string `json:"code"`
	PreCode string `json:"preCode"`

	// 返回的id编号
	Ids []IdCodeInfo `json:"ids"`
}

type CodesIds struct {
	IdCodes []IdCodes `json:"idCodes"`
}

type ObjectIdCache struct {
	OrgId     int64  `json:"orgId"`
	Code      string `json:"code"`
	MaxId     int64  `json:"maxId"`
	Threshold int64  `json:"threshold"`
}
