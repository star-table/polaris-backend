package api

import (
	"fmt"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/msgvo"
	"github.com/star-table/polaris-backend/service/basic/msgsvc/service"
)

func (PostGreeter) PushMsgToMq(msg msgvo.PushMsgToMqReqVo) vo.VoidErr{
	err := service.PushMsgToMq(msg)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (PostGreeter) InsertMqConsumeFailMsg(msg msgvo.InsertMqConsumeFailMsgReqVo) vo.VoidErr{
	err := service.InsertMqConsumeFailMsg(msg)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

func (GetGreeter) GetFailMsgList(req msgvo.GetFailMsgListReqVo) msgvo.GetFailMsgListRespVo{
	res, err := service.GetFailMsgList(req)
	return msgvo.GetFailMsgListRespVo{Data:msgvo.GetFailMsgListRespData{MsgList:res,}, Err: vo.NewErr(err)}
}

func (PostGreeter) UpdateMsgStatus(req msgvo.UpdateMsgStatusReqVo) vo.VoidErr{
	err := service.UpdateMsgStatus(req)
	return vo.VoidErr{Err: vo.NewErr(err)}
}

// WriteSomeFailedMsg 将一些场景下的任务处理失败的消息记录到表中，使用 msgType 进行区分查询。
func (PostGreeter) WriteSomeFailedMsg(msg msgvo.WriteSomeFailedMsgReqVo) vo.VoidErr {
	var busiErr errs.SystemErrorInfo
	err := service.WriteSomeFailedMsg(msg)
	if err != nil {
		// 处理错误，打印错误栈，并返回错误。
		errMsg := fmt.Sprintf("%+v", err)
		log.Error(errMsg)
		busiErr = errs.BuildSystemErrorInfo(errs.ServerError, err)
	}
	return vo.VoidErr{Err: vo.NewErr(busiErr)}
}
