package api

import (
	"context"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) GetProjectMainChatID(ctx context.Context, params vo.GetProjectMainChatIDReq) (*vo.GetProjectMainChatIDResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	respVo := projectfacade.GetProjectMainChatId(projectvo.GetProjectMainChatIdReq{
		OrgId:         cacheUserInfo.OrgId,
		UserId:        cacheUserInfo.UserId,
		SourceChannel: cacheUserInfo.SourceChannel,
		ProjectId:     params.ProjectID,
	})

	return &vo.GetProjectMainChatIDResp{ChatID: respVo.ChatId}, respVo.Error()
}

func (r *queryResolver) GetFsProjectChatPushSettings(ctx context.Context, params vo.GetFsProjectChatPushSettingsReq) (*vo.GetFsProjectChatPushSettingsResp, error) {
	return nil, errs.CannotBindChat

}

func (r *mutationResolver) UpdateFsProjectChatPushSettings(ctx context.Context, params vo.UpdateFsProjectChatPushSettingsReq) (*vo.Void, error) {
	return nil, errs.CannotBindChat
}
