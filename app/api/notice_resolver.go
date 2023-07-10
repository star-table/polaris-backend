package api

import (
	"context"

	pushPb "github.com/star-table/interface/golang/push/v1"
	"github.com/star-table/polaris-backend/facade/pushfacade"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/trendsfacade"
)

func (r *queryResolver) NoticeList(ctx context.Context, page *int, size *int, params *vo.NoticeListReq) (*vo.NoticeList, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	var defaultPage, defaultSize int
	if page == nil {
		page = &defaultPage
	}
	if size == nil {
		size = &defaultSize
	}

	resp := trendsfacade.NoticeList(trendsvo.NoticeListReqVo{
		UserId: cacheUserInfo.UserId,
		OrgId:  cacheUserInfo.OrgId,
		Page:   *page,
		Size:   *size,
		Input:  params,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}

	return resp.Data, nil
}

func (r *queryResolver) GetMQTTChannelKey(ctx context.Context, input vo.GetMQTTChannelKeyReq) (*vo.GetMQTTChannelKeyResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	var projectId int64
	if input.ProjectID != nil {
		projectId = *input.ProjectID
	}
	resp := pushfacade.GenerateMqttKey(&pushPb.GenerateMqttKeyReq{
		ChannelType: int32(input.ChannelType),
		OrgId:       cacheUserInfo.OrgId,
		ProjectId:   projectId,
	})
	if resp.Failure() {
		return nil, resp.Error()
	}
	var port int
	port = int(resp.Data.Port)
	return &vo.GetMQTTChannelKeyResp{
		Address: resp.Data.Address,
		Host:    resp.Data.Host,
		Port:    &port,
		Channel: resp.Data.Channel,
		Key:     resp.Data.Key,
	}, nil
}
