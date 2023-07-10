package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

//func (PostGreeter) InnerCreateTodoAudit(reqVo *projectvo.InnerCreateTodoAuditReq) *vo.VoidErr {
//	return &vo.VoidErr{Err: vo.NewErr(service.CreateTodoAudit(reqVo))}
//}
//
//func (PostGreeter) InnerCreateTodoFillIn(reqVo *projectvo.InnerCreateTodoFillInReq) *vo.VoidErr {
//	return &vo.VoidErr{Err: vo.NewErr(service.CreateTodoFillIn(reqVo))}
//}

func (PostGreeter) InnerCreateTodoHook(reqVo *projectvo.InnerCreateTodoHookReq) *vo.VoidErr {
	return &vo.VoidErr{Err: vo.NewErr(service.CreateTodoHook(reqVo))}
}
