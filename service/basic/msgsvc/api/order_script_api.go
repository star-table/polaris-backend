package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/msgvo"
	"github.com/star-table/polaris-backend/service/basic/msgsvc/service"
)

// FixAddOrderForFeiShu 重跑消费订单异常的订单消息
func (GetGreeter) FixAddOrderForFeiShu(msg msgvo.FixAddOrderForFeiShuReqVo) vo.VoidErr {
	err := service.FixAddOrderForFeiShu(msg.Input)
	return vo.VoidErr{Err: vo.NewErr(err)}
}
