/**
2 * @Author: Nico
3 * @Date: 2020/1/31 11:18
4 */
package orgvo

import "github.com/star-table/polaris-backend/common/model/vo"

type UserRegisterReqVo struct {
	Input vo.UserRegisterReq `json:"input"`
}

type UserRegisterRespVo struct {
	Data *vo.UserRegisterResp `json:"data"`
	vo.Err
}

