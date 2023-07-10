package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

//func (PostGreeter) TodoFilter(reqVo *projectvo.TodoFilterReq) *projectvo.TodoFilterResp {
//	todos, err := service.TodoFilter(reqVo)
//	return &projectvo.TodoFilterResp{Err: vo.NewErr(err), Data: todos}
//}

func (PostGreeter) TodoUrge(reqVo *projectvo.TodoUrgeReq) *vo.CommonRespVo {
	err := service.TodoUrge(reqVo)
	return &vo.CommonRespVo{Err: vo.NewErr(err)}
}
