package api

import (
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
)

//飞书初始化调用
func (PostGreeter) InitOrg(reqVo orgvo.InitOrgReqVo) orgvo.OrgInitRespVo {
	respVo := orgvo.OrgInitRespVo{}
	return respVo
}

//发送飞书帮助信息
func (PostGreeter) SendFeishuMemberHelpMsg(reqVo orgvo.SendFeishuMemberHelpMsgReqVo) vo.VoidErr {
	return vo.VoidErr{Err: vo.NewErr(nil)}
}

func (GetGreeter) ScheduleOrgUseMobileAndEmail(reqVo orgvo.ScheduleOrgUseMobileAndEmailReqVo) vo.VoidErr {
	return vo.VoidErr{Err: vo.NewErr(nil)}
}
