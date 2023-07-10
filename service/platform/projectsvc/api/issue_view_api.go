package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service/issue_view"
)

/// 任务视图相关的操作
// 创建视图
func (PostGreeter) CreateTaskView(reqVo *projectvo.CreateTaskViewReqVo) projectvo.CreateTaskViewRespVo {
	res, err := issue_view.CreateTaskView(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
	return projectvo.CreateTaskViewRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) GetTaskViewList(reqVo *projectvo.GetTaskViewListReqVo) projectvo.GetTaskViewListRespVo {
	res, err := issue_view.GetTaskViewList(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
	return projectvo.GetTaskViewListRespVo{
		Err:  vo.NewErr(err),
		Data: res,
	}
}

func (PostGreeter) UpdateTaskView(reqVo *projectvo.UpdateTaskViewReqVo) vo.CommonRespVo {
	res, err := issue_view.UpdateTaskView(reqVo.OrgId, reqVo.UserId, &reqVo.Input)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
}

func (PostGreeter) DeleteTaskView(reqVo *projectvo.DeleteTaskViewReqVo) vo.CommonRespVo {
	res, err := issue_view.DeleteTaskView(reqVo.OrgId, reqVo.UserId, reqVo.Input.ID)
	return vo.CommonRespVo{
		Err:  vo.NewErr(err),
		Void: res,
	}
}

// 筛选时，获取值列表
//func (PostGreeter) GetOptionList(reqVo *projectvo.GetOptionListReqVo) projectvo.GetOptionListRespVo {
//	return projectvo.GetOptionListRespVo{}
//
//	//res, err := issue_view.GetOptionList(reqVo.OrgId, reqVo.Input)
//	//return projectvo.GetOptionListRespVo{
//	//	Err:  vo.NewErr(err),
//	//	NewData: res,
//	//}
//}
