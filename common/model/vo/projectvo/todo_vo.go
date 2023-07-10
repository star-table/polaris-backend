package projectvo

//// TodoFilterReq .
//type TodoFilterReq struct {
//	OrgId      int64 `json:"orgId"`
//	UserId     int64 `json:"userId"`
//	Page       int   `json:"page"`
//	Size       int   `json:"size"`
//	FilterType int   `json:"filterType"`
//}
//
//// TodoFilterResp .
//type TodoFilterResp struct {
//	vo.Err
//	Data []*bo.Todo `json:"data"`
//}

type TodoUrgeInput struct {
	TodoId int64  `json:"todoId"`
	Msg    string `json:"msg"`
}

// TodoUrgeReq .
type TodoUrgeReq struct {
	OrgId  int64          `json:"orgId"`
	UserId int64          `json:"userId"`
	Input  *TodoUrgeInput `json:"input"`
}
