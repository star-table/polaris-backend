package domain

import (
	"github.com/star-table/polaris-backend/common/core/errs"
)

//func InsertTodo(todo *bo.Todo) errs.SystemErrorInfo {
//	p := &po.LcTodo{}
//	p.ConvertFromBo(todo)
//	err := mysql.Insert(p)
//	if err != nil {
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	return nil
//}
//
//func GetTodoById(todoId int64) (*bo.Todo, errs.SystemErrorInfo) {
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//
//	res := &po.LcTodo{}
//	s := conn.Select("*").From(consts.TableLcTodo).Where(db.Raw("id=?", todoId))
//	if err = s.One(res); err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//	return res.ConvertToBo(), nil
//}
//
//func TodoFilter(orgId, userId int64, filterType int, page, size int) ([]*bo.Todo, errs.SystemErrorInfo) {
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		log.Errorf("[TodoFilter] mysql error: %v", err)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//
//	var where db.RawValue
//	switch filterType {
//	case consts.TodoFilterTypeWaitingAudit:
//		where = db.Raw(fmt.Sprintf("org_id=%d and type=%d and status=%d and json_extract(operators, '$.\"%d\".op')=%d", orgId, consts.TodoTypeAudit, consts.TodoStatusUnFinished, userId, consts.OpInit))
//	case consts.TodoFilterTypeWaitingFillIn:
//		where = db.Raw(fmt.Sprintf("org_id=%d and type=%d and status=%d and json_extract(operators, '$.\"%d\".op')=%d", orgId, consts.TodoTypeFillIn, consts.TodoStatusUnFinished, userId, consts.OpInit))
//	case consts.TodoFilterTypeSelfTrigger:
//		where = db.Raw(fmt.Sprintf("org_id=%d and trigger_user_id=%d", orgId, userId))
//	case consts.TodoFilterTypeFinished:
//		where = db.Raw(fmt.Sprintf("org_id=%d and status=%d and json_extract(operators, '$.\"%d\".op')<>%d", orgId, consts.TodoStatusFinished, userId, consts.OpInit))
//	default:
//		log.Errorf("[TodoFilter] unknown filterType: %v", filterType)
//		return nil, errs.ParamError
//	}
//
//	var res []*po.LcTodo
//	s := conn.Select("*").From(consts.TableLcTodo).Where(where)
//
//	// 分页查询
//	if page > 0 && size > 0 {
//		s = s.Offset((page - 1) * size).Limit(size)
//	}
//
//	log.Infof("[TodoFilter] %s", s.String())
//	err = s.All(&res)
//	if err != nil {
//		log.Errorf("[TodoFilter] %s, err: %v", s.String(), err)
//		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
//	}
//
//	var todos []*bo.Todo
//	for _, po := range res {
//		todos = append(todos, po.ConvertToBo())
//	}
//	return todos, nil
//}

func TodoUrge(orgId, userId, todoId int64, msg string) errs.SystemErrorInfo {
	return nil
}
