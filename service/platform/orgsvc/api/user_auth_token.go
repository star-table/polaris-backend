package api

import (
	"context"

	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) GetCurrentUser(ctx context.Context) orgvo.CacheUserInfoVo {
	info, err := service.GetCurrentUser(ctx)
	res := orgvo.CacheUserInfoVo{Err: vo.NewErr(err)}
	if info != nil {
		res.CacheInfo = *info
	}
	return res
}

func (GetGreeter) GetCurrentUserWithoutOrgVerify(ctx context.Context) orgvo.CacheUserInfoVo {
	info, err := service.GetCurrentUserWithoutOrgVerify(ctx)
	res := orgvo.CacheUserInfoVo{Err: vo.NewErr(err)}
	if info != nil {
		res.CacheInfo = *info
	}
	return res
}

func (GetGreeter) GetCurrentUserWithoutPayVerify(ctx context.Context) orgvo.CacheUserInfoVo {
	info, err := service.GetCurrentUserWithCond(ctx, true, false)
	res := orgvo.CacheUserInfoVo{Err: vo.NewErr(err)}
	if info != nil {
		res.CacheInfo = *info
	}
	return res
}
