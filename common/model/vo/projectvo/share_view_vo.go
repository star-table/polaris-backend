package projectvo

import "github.com/star-table/polaris-backend/common/model/vo"

type CreateShareViewData struct {
	AppId     int64 `json:"appId"`
	ProjectId int64 `json:"projectId"`
	TableId   int64 `json:"tableId"`
	ViewId    int64 `json:"viewId"`
}

type CreateShareViewReq struct {
	OrgId  int64                `json:"orgId"`
	UserId int64                `json:"userId"`
	Input  *CreateShareViewData `json:"input"`
}

type ShareKeyData struct {
	ShareKey string `json:"shareKey"`
}

type GetShareViewInfoByKeyReq struct {
	Input *ShareKeyData `json:"input"`
}

type ShareViewIdData struct {
	ViewId int64 `json:"viewId"`
}

type GetShareViewInfoReq struct {
	OrgId  int64            `json:"orgId"`
	UserId int64            `json:"userId"`
	Input  *ShareViewIdData `json:"input"`
}

type ShareViewInfo struct {
	ShareKey      string `json:"shareKey"`
	Config        string `json:"config"`
	IsSetPassword bool   `json:"isSetPassword"`

	TableId   string `json:"tableId,omitempty"`
	AppId     string `json:"appId,omitempty"`
	ProjectId int64  `json:"projectId,omitempty"`
	ViewId    string `json:"viewId,omitempty"`
}

type GetShareViewInfoResp struct {
	vo.Err
	Data *ShareViewInfo `json:"data"`
}

type DeleteShareViewReq struct {
	OrgId  int64            `json:"orgId"`
	UserId int64            `json:"userId"`
	Input  *ShareViewIdData `json:"input"`
}

type UpdateData struct {
	ViewId   int64  `json:"viewId"`
	Config   string `json:"config"`
	Password string `json:"password"`
}

type UpdateShareConfigReq struct {
	OrgId  int64       `json:"orgId"`
	UserId int64       `json:"userId"`
	Input  *UpdateData `json:"input"`
}

type UpdateSharePasswordReq struct {
	OrgId  int64       `json:"orgId"`
	UserId int64       `json:"userId"`
	Input  *UpdateData `json:"input"`
}

type ResetShareKeyReq struct {
	OrgId  int64            `json:"orgId"`
	UserId int64            `json:"userId"`
	Input  *ShareViewIdData `json:"input"`
}

type CheckPasswordData struct {
	ShareKey string `json:"shareKey"`
	Password string `json:"password"`
}

type CheckShareViewPasswordReq struct {
	Input *CheckPasswordData `json:"input"`
}

type CheckPassword struct {
	IsCorrect bool  `json:"isCorrect"`
	UserId    int64 `json:"userId"`
	OrgId     int64 `json:"orgId"`
}

type CheckShareViewPasswordResp struct {
	vo.Err
	Data *CheckPassword `json:"data"`
}
