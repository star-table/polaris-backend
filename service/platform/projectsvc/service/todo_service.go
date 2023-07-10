package service

import (
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

func CreateTodoHook(req *projectvo.InnerCreateTodoHookReq) errs.SystemErrorInfo {
	return nil
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//type CreateTodoAuditContext struct {
//	Req       *projectvo.InnerCreateTodoAuditReq
//	Operators map[int64]*bo.TodoResult
//}

//func CreateTodoAudit(req *projectvo.InnerCreateTodoAuditReq) errs.SystemErrorInfo {
//	ctx := &CreateTodoAuditContext{
//		Req: req,
//	}
//	log.Infof("[CreateTodoAudit] req: %v", json.ToJsonIgnoreError(req))
//
//	// 1. 检查参数 组装数据
//	errSys := ctx.prepare()
//	if errSys != nil {
//		return errSys
//	}
//
//	// 2. 保存数据
//	errSys = ctx.saveToDB()
//	if errSys != nil {
//		return errSys
//	}
//
//	return nil
//}
//
//func (ctx *CreateTodoAuditContext) prepare() errs.SystemErrorInfo {
//	ctx.Operators = make(map[int64]*bo.TodoResult)
//
//	if ctx.Req.Input.TriggerBy.WorkflowId == 0 || ctx.Req.Input.TriggerBy.ExecutionId == 0 {
//		log.Errorf("[CreateTodoAudit] cannot create todo, invalid param. workflowId: %d, executionId: %d",
//			ctx.Req.Input.TriggerBy.WorkflowId, ctx.Req.Input.TriggerBy.ExecutionId)
//		return errs.ParamError
//	}
//
//	if ctx.Req.Input.AppId == 0 || ctx.Req.Input.TableId == 0 || ctx.Req.Input.DataId == 0 {
//		log.Errorf("[CreateTodoAudit] cannot create todo, invalid param. appId: %d, tableId: %d, dataId: %d",
//			ctx.Req.Input.AppId, ctx.Req.Input.TableId, ctx.Req.Input.DataId)
//		return errs.ParamError
//	}
//
//	if ctx.Req.Input.TriggerUserId == 0 {
//		log.Error("[CreateTodoAudit] cannot create todo, no trigger.")
//		return errs.ParamError
//	}
//
//	// 获取操作人
//	for _, op := range ctx.Req.Input.Operators {
//		switch op.Type {
//		// 成员
//		case consts.OperatorTypeUser:
//			ids, _, _, _ := businees.LcInterfaceToIds(consts.LcCustomFieldUserType, op.Ids, true, true)
//			for _, userId := range ids {
//				ctx.Operators[userId] = &bo.TodoResult{
//					Op: consts.OpInit,
//				}
//			}
//		default:
//			log.Errorf("[CreateTodoAudit] unknown operator type: %v, req: %v", op.Type, json.ToJsonIgnoreError(ctx.Req))
//		}
//	}
//	if len(ctx.Operators) == 0 {
//		log.Error("[CreateTodoAudit] cannot create todo, no operators.")
//		return errs.ParamError
//	}
//
//	return nil
//}
//
//func (ctx *CreateTodoAuditContext) saveToDB() errs.SystemErrorInfo {
//	now := time.Now()
//	todo := &bo.Todo{}
//	todo.Id = snowflake.Id()
//	todo.OrgId = ctx.Req.OrgId
//	todo.AppId = ctx.Req.Input.AppId
//	todo.TableId = ctx.Req.Input.TableId
//	todo.DataId = ctx.Req.Input.DataId
//	todo.WorkflowId = ctx.Req.Input.TriggerBy.WorkflowId
//	todo.WorkflowName = ctx.Req.Input.TriggerBy.WorkflowName
//	todo.ExecutionId = ctx.Req.Input.TriggerBy.ExecutionId
//	todo.TriggerUserId = ctx.Req.Input.TriggerUserId
//	todo.AllowWithdrawByTrigger = ctx.Req.Input.AllowWithdrawByTrigger
//	todo.AllowUrgeByTrigger = ctx.Req.Input.AllowUrgeByTrigger
//	todo.Type = consts.TodoTypeAudit
//	todo.Status = consts.TodoStatusUnFinished
//	todo.Parameters = ctx.Req.Input.Parameters
//	todo.Operators = ctx.Operators
//	todo.Creator = ctx.Req.UserId
//	todo.CreatedAt = now
//	todo.Updater = ctx.Req.UserId
//	todo.UpdatedAt = now
//
//	// save to db
//	errSys := domain.InsertTodo(todo)
//	if errSys != nil {
//		log.Errorf("[CreateTodoAudit] InsertTodo error: %v", errSys)
//		return errSys
//	}
//	return nil
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//type CreateTodoFillInContext struct {
//	Req       *projectvo.InnerCreateTodoFillInReq
//	Operators map[int64]*bo.TodoResult
//}
//
//func CreateTodoFillIn(req *projectvo.InnerCreateTodoFillInReq) errs.SystemErrorInfo {
//	ctx := &CreateTodoFillInContext{
//		Req: req,
//	}
//	log.Infof("[CreateTodoFillIn] req: %v", json.ToJsonIgnoreError(req))
//
//	// 1. 检查参数 组装数据
//	errSys := ctx.prepare()
//	if errSys != nil {
//		return errSys
//	}
//
//	// 2. 保存数据
//	errSys = ctx.saveToDB()
//	if errSys != nil {
//		return errSys
//	}
//
//	return nil
//}
//
//func (ctx *CreateTodoFillInContext) prepare() errs.SystemErrorInfo {
//	ctx.Operators = make(map[int64]*bo.TodoResult)
//
//	if ctx.Req.Input.TriggerBy.WorkflowId == 0 || ctx.Req.Input.TriggerBy.ExecutionId == 0 {
//		log.Errorf("[CreateTodoFillIn] cannot create todo, invalid param. workflowId: %d, executionId: %d",
//			ctx.Req.Input.TriggerBy.WorkflowId, ctx.Req.Input.TriggerBy.ExecutionId)
//		return errs.ParamError
//	}
//
//	if ctx.Req.Input.AppId == 0 || ctx.Req.Input.TableId == 0 || ctx.Req.Input.DataId == 0 {
//		log.Errorf("[CreateTodoFillIn] cannot create todo, invalid param. appId: %d, tableId: %d, dataId: %d",
//			ctx.Req.Input.AppId, ctx.Req.Input.TableId, ctx.Req.Input.DataId)
//		return errs.ParamError
//	}
//
//	if ctx.Req.Input.TriggerUserId == 0 {
//		log.Error("[CreateTodoFillIn] cannot create todo, no trigger.")
//		return errs.ParamError
//	}
//
//	// 获取操作人
//	for _, op := range ctx.Req.Input.Operators {
//		switch op.Type {
//		// 成员
//		case consts.OperatorTypeUser:
//			ids, _, _, _ := businees.LcInterfaceToIds(consts.LcCustomFieldUserType, op.Ids, true, true)
//			for _, userId := range ids {
//				ctx.Operators[userId] = &bo.TodoResult{
//					Op: consts.OpInit,
//				}
//			}
//		default:
//			log.Errorf("[CreateTodoFillIn] unknown operator type: %v, req: %v", op.Type, json.ToJsonIgnoreError(ctx.Req))
//		}
//	}
//	if len(ctx.Operators) == 0 {
//		log.Error("[CreateTodoFillIn] cannot create todo, no operators.")
//		return errs.ParamError
//	}
//
//	return nil
//}
//
//func (ctx *CreateTodoFillInContext) saveToDB() errs.SystemErrorInfo {
//	now := time.Now()
//	todo := &bo.Todo{}
//	todo.Id = snowflake.Id()
//	todo.OrgId = ctx.Req.OrgId
//	todo.AppId = ctx.Req.Input.AppId
//	todo.TableId = ctx.Req.Input.TableId
//	todo.DataId = ctx.Req.Input.DataId
//	todo.WorkflowId = ctx.Req.Input.TriggerBy.WorkflowId
//	todo.WorkflowName = ctx.Req.Input.TriggerBy.WorkflowName
//	todo.ExecutionId = ctx.Req.Input.TriggerBy.ExecutionId
//	todo.TriggerUserId = ctx.Req.Input.TriggerUserId
//	todo.AllowWithdrawByTrigger = ctx.Req.Input.AllowWithdrawByTrigger
//	todo.AllowUrgeByTrigger = ctx.Req.Input.AllowUrgeByTrigger
//	todo.Type = consts.TodoTypeFillIn
//	todo.Status = consts.TodoStatusUnFinished
//	todo.Parameters = ctx.Req.Input.Parameters
//	todo.Operators = ctx.Operators
//	todo.Creator = ctx.Req.UserId
//	todo.CreatedAt = now
//	todo.Updater = ctx.Req.UserId
//	todo.UpdatedAt = now
//
//	// save to db
//	errSys := domain.InsertTodo(todo)
//	if errSys != nil {
//		log.Errorf("[CreateTodoFillIn] InsertTodo error: %v", errSys)
//		return errSys
//	}
//
//	return nil
//}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

//func TodoFilter(req *projectvo.TodoFilterReq) ([]*bo.Todo, errs.SystemErrorInfo) {
//	return domain.TodoFilter(req.OrgId, req.UserId, req.FilterType, req.Page, req.Size)
//}

func TodoUrge(req *projectvo.TodoUrgeReq) errs.SystemErrorInfo {
	return domain.TodoUrge(req.OrgId, req.UserId, req.Input.TodoId, req.Input.Msg)
}
