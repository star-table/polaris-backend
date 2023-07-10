package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type ProjectChatListReqVo struct {
	OrgId  int64                 `json:"orgId"`
	UserId int64                 `json:"userId"`
	Input  vo.ProjectChatListReq `json:"input"`
}

type ProjectChatListRespVo struct {
	vo.Err
	List *vo.ChatListResp `json:"list"`
}

type UnrelatedChatListReqVo struct {
	OrgId  int64                   `json:"orgId"`
	UserId int64                   `json:"userId"`
	Input  vo.UnrelatedChatListReq `json:"input"`
}

type UpdateRelateChatReqVo struct {
	OrgId  int64               `json:"orgId"`
	UserId int64               `json:"userId"`
	Input  vo.UpdateRelateChat `json:"input"`
}

type FsChatDisbandCallbackReq struct {
	OrgId  int64  `json:"orgId"`
	ChatId string `json:"chatId"`
}

type GetProjectMainChatIdReq struct {
	OrgId         int64  `json:"orgId"`
	UserId        int64  `json:"userId"`
	SourceChannel string `json:"sourceChannel"`
	ProjectId     int64  `json:"projectId"`
}

type CheckIsShowProChatIconReq struct {
	OrgId         int64                         `json:"orgId"`
	UserId        int64                         `json:"userId"`
	SourceChannel string                        `json:"sourceChannel"`
	Input         CheckIsShowProChatIconReqData `json:"input"`
}

type CheckIsShowProChatIconReqData struct {
	AppId string `json:"appId"`
}

type GetProjectMainChatIdResp struct {
	vo.Err
	ChatId string `json:"chatId"`
}

type CheckIsShowProChatIconResp struct {
	vo.Err
	Data CheckShowProChatIconRespData `json:"data"`
}

type CheckShowProChatIconRespData struct {
	IsShow bool `json:"isShow"`
}

type UpdateFsProjectChatPushSettingsReq struct {
	OrgId         int64                                 `json:"orgId"`
	UserId        int64                                 `json:"userId"`
	SourceChannel string                                `json:"sourceChannel"`
	Input         vo.UpdateFsProjectChatPushSettingsReq `json:"input"`
}

type GetFsProjectChatPushSettingsReq struct {
	OrgId         int64  `json:"orgId"`
	ChatId        string `json:"chatId"`
	ProjectId     int64  `json:"projectId"`
	SourceChannel string `json:"sourceChannel"`
}

type GetFsProjectChatPushSettingsResp struct {
	vo.Err
	Data *vo.GetFsProjectChatPushSettingsResp `json:"data"`
}

type DeleteChatReq struct {
	OutOrgId      string `json:"outOrgId"`
	ProjectId     int64  `json:"projectId"`
	ChatId        string `json:"chatId"`
	SourceChannel string `json:"sourceChannel"`
}
