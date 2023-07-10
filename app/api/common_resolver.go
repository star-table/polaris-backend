package api

import (
	"context"

	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade/commonfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
)

func (r *queryResolver) AreaLinkageList(ctx context.Context, input vo.AreaLinkageListReq) (*vo.AreaLinkageListResp, error) {
	_, err := orgfacade.GetCurrentUserRelaxed(ctx)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	respVo := commonfacade.AreaLinkageList(commonvo.AreaLinkageListReqVo{
		Input: input,
	})

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.AreaLinkageListResp, nil

}

func (r *queryResolver) IndustryList(ctx context.Context) (*vo.IndustryListResp, error) {

	//_, err := orgfacade.GetCurrentUserRelaxed(ctx)
	//
	//if err != nil {
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}

	respVo := commonfacade.IndustryList()

	if respVo.Failure() {
		log.Error(respVo.Message)
		return nil, respVo.Error()
	}
	return respVo.IndustryList, nil
}
