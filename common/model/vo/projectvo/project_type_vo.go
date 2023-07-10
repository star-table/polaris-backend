package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type ProjectTypesRespVo struct {
	vo.Err
	ProjectTypes []*vo.ProjectType `json:"data"`
}

type ProjectTypesReqVo struct {
	OrgId int64 `json:"orgId"`
}

type ProjectTypeCategoryVo struct {
	vo.Err
	Data *vo.ProjectTypeCategoryList `json:"data"`
}

type ProjectTypeListReqVo struct {
	OrgId int64 `json:"orgId"`
	CategoryId int64 `json:"categoryId"`
}

type ProjectTypeListRespVo struct {
	vo.Err
	Data *vo.ProjectTypeListResp `json:"data"`
}
