package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
)

func (PostGreeter) ProjectChatList(req projectvo.ProjectChatListReqVo) projectvo.ProjectChatListRespVo {
	return projectvo.ProjectChatListRespVo{
		Err:  vo.NewErr(nil),
		List: nil,
	}
}

func (PostGreeter) UnrelatedChatList(req projectvo.UnrelatedChatListReqVo) projectvo.ProjectChatListRespVo {
	return projectvo.ProjectChatListRespVo{
		Err:  vo.NewErr(nil),
		List: nil,
	}
}

func (PostGreeter) AddProjectChat(req projectvo.UpdateRelateChatReqVo) vo.CommonRespVo {
	return vo.CommonRespVo{
		Err:  vo.NewErr(nil),
		Void: nil,
	}
}

func (PostGreeter) DisbandProjectChat(req projectvo.UpdateRelateChatReqVo) vo.CommonRespVo {
	return vo.CommonRespVo{
		Err:  vo.NewErr(nil),
		Void: nil,
	}
}

func (PostGreeter) FsChatDisbandCallback(req projectvo.FsChatDisbandCallbackReq) vo.CommonRespVo {
	return vo.CommonRespVo{
		Err:  vo.NewErr(nil),
		Void: nil,
	}
}

func (PostGreeter) GetProjectMainChatId(req projectvo.GetProjectMainChatIdReq) projectvo.GetProjectMainChatIdResp {
	return projectvo.GetProjectMainChatIdResp{
		Err:    vo.NewErr(nil),
		ChatId: "",
	}
}

func (PostGreeter) CheckIsShowProChatIcon(req projectvo.CheckIsShowProChatIconReq) projectvo.CheckIsShowProChatIconResp {
	return projectvo.CheckIsShowProChatIconResp{
		Err: vo.NewErr(nil),
		Data: projectvo.CheckShowProChatIconRespData{
			IsShow: false,
		},
	}
}

func (PostGreeter) UpdateFsProjectChatPushSettings(req projectvo.UpdateFsProjectChatPushSettingsReq) vo.CommonRespVo {
	return vo.CommonRespVo{
		Err:  vo.NewErr(nil),
		Void: &vo.Void{ID: 1},
	}
}

func (PostGreeter) GetFsProjectChatPushSettings(req projectvo.GetFsProjectChatPushSettingsReq) projectvo.GetFsProjectChatPushSettingsResp {
	return projectvo.GetFsProjectChatPushSettingsResp{
		Err:  vo.NewErr(nil),
		Data: &vo.GetFsProjectChatPushSettingsResp{},
	}
}

func (PostGreeter) DeleteChatCallback(req projectvo.DeleteChatReq) vo.CommonRespVo {
	return vo.CommonRespVo{
		Err:  vo.NewErr(nil),
		Void: nil,
	}
}
