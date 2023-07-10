package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/service"
)

func (GetGreeter) IndustryList() commonvo.IndustryListRespVo {
	res, err := service.IndustryList()
	return commonvo.IndustryListRespVo{Err: vo.NewErr(err), IndustryList: res}
}
