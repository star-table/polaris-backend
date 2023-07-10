package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
)

// 项目群聊，响应用户指令，项目任务
func (PostGreeter) HandleGroupChatUserInsProIssue(reqVo projectvo.HandleGroupChatUserInsProIssueReqVo) projectvo.BoolRespVo {
	isOk, err := service.HandleGroupChatUserInsProIssue(reqVo.Input)
	return projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	}
}

// 项目群聊，响应用户指令，项目进展
func (PostGreeter) HandleGroupChatUserInsProProgress(reqVo projectvo.HandleGroupChatUserInsProProgressReqVo) projectvo.BoolRespVo {
	isOk, err := service.HandleGroupChatUserInsProProgress(reqVo.Input)
	return projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	}
}

// 项目群聊，响应用户指令，项目设置
func (PostGreeter) HandleGroupChatUserInsProjectSettings(reqVo projectvo.HandleGroupChatUserInsProjectSettingsReq) projectvo.BoolRespVo {
	isOk, err := service.HandleGroupChatUserInsProjectSettings(reqVo.OpenChatId, reqVo.SourceChannel)
	return projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	}
}

// 项目群聊，响应用户指令：@用户姓名
func (PostGreeter) HandleGroupChatUserInsAtUserName(reqVo projectvo.HandleGroupChatUserInsAtUserNameReqVo) projectvo.BoolRespVo {
	isOk, err := service.HandleGroupChatAtUserName(reqVo.Input)
	return projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	}
}

// 项目群聊，响应用户指令:@用户姓名 任务标题1
func (PostGreeter) HandleGroupChatUserInsAtUserNameWithIssueTitle(reqVo projectvo.HandleGroupChatUserInsAtUserNameWithIssueTitleReqVo) projectvo.BoolRespVo {
	isOk, err := service.HandleGroupChatAtUserNameWithIssueTitle(reqVo.Input)
	return projectvo.BoolRespVo{
		Err: vo.NewErr(err),
		Data: &vo.BoolResp{
			IsTrue: isOk,
		},
	}
}

// GetProjectIdByChatId 通过 chatId 获取群聊对应的项目id
func (PostGreeter) GetProjectIdByChatId(reqVo projectvo.GetProjectIdByChatIdReqVo) projectvo.GetProjectIdByChatIdRespVo {
	resp, err := service.GetProjectIdByChatId(reqVo)
	return projectvo.GetProjectIdByChatIdRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}

// GetProjectIdsByChatId 通过 chatId 获取群聊对应的项目。因为现在一个群聊可以关联多个项目 id todo
func (PostGreeter) GetProjectIdsByChatId(reqVo projectvo.GetProjectIdsByChatIdReqVo) projectvo.GetProjectIdsByChatIdRespVo {
	resp, err := service.GetProjectIdsByChatId(reqVo.OrgId, reqVo)
	return projectvo.GetProjectIdsByChatIdRespVo{
		Err:  vo.NewErr(err),
		Data: resp,
	}
}
