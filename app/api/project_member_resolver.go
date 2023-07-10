package api

import (
	"context"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
)

func (r *queryResolver) ProjectMemberIDList(ctx context.Context, input vo.ProjectMemberIDListReq) (*vo.ProjectMemberIDListResp, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(ctx)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}
	if input.IncludeAdmin == nil {
		defaultVal := 0
		input.IncludeAdmin = &defaultVal
	}

	resp := projectfacade.ProjectMemberIdList(projectvo.ProjectMemberIdListReq{
		OrgId:     cacheUserInfo.OrgId,
		ProjectId: input.ProjectID,
		Data: &projectvo.ProjectMemberIdListReqData{
			IncludeAdmin: *input.IncludeAdmin,
		},
	})

	return resp.Data, resp.Error()
}
