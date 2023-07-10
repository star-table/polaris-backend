package bo

type Collaborator struct {
	AppId     int64  `json:"appId"`
	OrgId     int64  `json:"orgId"`
	TableId   int64  `json:"tableId"`
	FieldName string `json:"fieldName"`
	UserId    int64  `json:"userId"`
	DeptId    int64  `json:"deptId"`
}
