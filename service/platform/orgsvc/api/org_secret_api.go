package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (GetGreeter) OpenAPIAuth(req orgvo.OpenAPIAuthReq) orgvo.OpenAPIAuthResp{
	data, err := service.OpenAPIAuth(req)
	return orgvo.OpenAPIAuthResp{Data: data, Err: vo.NewErr(err)}
}

func (PostGreeter)GetAppTicket(req orgvo.GetAppTicketReq) orgvo.GetAppTicketResp {
	data, err := service.GetAppTicket(req)
	return orgvo.GetAppTicketResp{
		Err:  vo.NewErr(err),
		Data: data,
	}
}
