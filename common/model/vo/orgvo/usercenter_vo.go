package orgvo

import "github.com/star-table/polaris-backend/common/model/vo"

//获取部门人数
type GetUserCountByDeptIdsReq struct {
	//组织id
	OrgId int64 `json:"orgId"`
	//部门id
	DeptIds []int64 `json:"deptIds"`
}

type GetUserCountByDeptIdsResp struct {
	vo.Err
	Data *GetUserCountByDeptIdsData `json:"data"`
}

type GetUserCountByDeptIdsData struct {
	UserCount map[int64]uint64 `json:"userCount"`
}

type GetUserDeptIdsResp struct {
	vo.Err
	Data GetUserDeptIdsData `json:"data"`
}

type GetUserDeptIdsData struct {
	DeptIds []int64 `json:"deptIds"`
}

type GetUserDeptIdsReq struct {
	//组织id
	OrgId int64 `json:"orgId"`
	//用户id
	UserId int64 `json:"userId"`
}

type GetUserDeptIdsBatchResp struct {
	vo.Err
	Data GetUserDeptIdsBatchData `json:"data"`
}

type GetUserDeptIdsBatchData struct {
	Data map[int64][]int64 `json:"data"`
}

type GetUserDeptIdsBatchReq struct {
	//组织id
	OrgId int64 `json:"orgId"`
	//用户id
	UserIds []int64 `json:"userIds"`
}

type GetUserIdsByDeptIdsReq struct {
	// 组织id
	OrgId int64 `json:"orgId"`
	// 部门id
	DeptIds []int64 `json:"deptIds"`
}

type GetUserIdsByDeptIdsResp struct {
	vo.Err
	Data GetUserIdsByDeptIdsData `json:"data"`
}

type GetUserIdsByDeptIdsData struct {
	UserIds []int64 `json:"userIds"`
}

type GetUserDeptIdsWithParentIdReq struct {
	OrgId  int64 `json:"orgId"`
	UserId int64 `json:"userId"`
}

type GetUserDeptIdsWithParentIdResp struct {
	vo.Err
	Data GetUserDeptIdsWithParentIdData `json:"data"`
}

type GetUserDeptIdsWithParentIdData struct {
	DeptIds []int64 `json:"deptIds"`
}
