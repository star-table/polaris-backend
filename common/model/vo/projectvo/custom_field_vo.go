package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type CreateCustomFieldReq struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Data   vo.CreateCustomFieldReq `json:"data"`
}

type DeleteCustomFieldReq struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Data   vo.DeleteCustomFieldReq `json:"data"`
}

type UseOrgCustomFieldReq struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Data   vo.UseOrgCustomFieldReq `json:"data"`
}

type UpdateCustomFieldReq struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Data   vo.UpdateCustomFieldReq `json:"data"`
}

type ChangeProjectCustomFieldStatusReq struct {
	OrgId  int64                                `json:"orgId"`
	UserId int64                                `json:"userId"`
	Data   vo.ChangeProjectCustomFieldStatusReq `json:"data"`
}

type CustomFieldListReq struct {
	OrgId  int64                 `json:"orgId"`
	UserId int64                 `json:"userId"`
	Data   vo.CustomFieldListReq `json:"data"`
}

type CustomFieldListResp struct {
	vo.Err
	Data *vo.CustomFieldListResp `json:"data"`
}

type UpdateIssueCustomFieldReq struct {
	OrgId  int64                        `json:"orgId"`
	UserId int64                        `json:"userId"`
	Data   vo.UpdateIssueCustomFieldReq `json:"data"`
}

type ProjectFieldViewReq struct {
	OrgId  int64                  `json:"orgId"`
	UserId int64                  `json:"userId"`
	Params vo.ProjectFieldViewReq `json:"params"`
}

type ProjectFieldViewResp struct {
	vo.Err
	Data *vo.ProjectFieldViewResp `json:"data"`
}

type UpdateProjectFieldViewReq struct {
	OrgId  int64                        `json:"orgId"`
	UserId int64                        `json:"userId"`
	Params vo.UpdateProjectFieldViewReq `json:"params"`
}

type SimpleNameInfo struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type GetSimpleCustomFieldInfoReq struct {
	OrgId    int64   `json:"orgId"`
	FieldIds []int64 `json:"fieldIds"`
}

type GetSimpleCustomFieldInfoResp struct {
	vo.Err
	Data []SimpleNameInfo `json:"data"`
}

type CreateCustomFieldResp struct {
	vo.Err
	Data *vo.CustomField `json:"data"`
}
