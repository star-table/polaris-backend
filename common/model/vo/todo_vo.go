package vo

//type TodoFilterReq struct {
//	FilterType int `json:"filterType"`
//	Page       int `json:"page"`
//	Size       int `json:"size"`
//}

type TodoUrgeReq struct {
	TodoId int64  `json:"todoId,string"`
	Msg    string `json:"msg"`
}
