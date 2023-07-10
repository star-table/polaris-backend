package projectvo

type AuthProjectPermissionReqVo struct {
	Input AuthProjectPermissionReqData `json:"input"`
}

type AuthProjectPermissionReqData struct {
	OrgId int64 `json:"orgId"`
	UserId int64 `json:"userId"`
	ProjectId int64 `json:"projectId"`
	Path string `json:"path"`
	Operation string `json:"operation"`
	AuthFiling bool `json:"authFiling"`
}
