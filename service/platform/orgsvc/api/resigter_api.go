/**
2 * @Author: Nico
3 * @Date: 2020/1/31 11:17
4 */
package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func (PostGreeter) UserRegister(req orgvo.UserRegisterReqVo) orgvo.UserRegisterRespVo {
	res, err := service.UserRegister(req)
	return orgvo.UserRegisterRespVo{Data: res, Err: vo.NewErr(err)}
}
