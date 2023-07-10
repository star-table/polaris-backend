package domain

func PushAddOrgMemberNotice(orgId, depId int64, memberIds []int64, operatorId int64) {
	//globalRefreshList := make([]bo.MQTTGlobalRefresh, 0)
	//
	//baseUserInfos, err := GetBaseUserInfoBatch(orgId, memberIds)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//for _, baseUserInfo := range baseUserInfos {
	//	globalRefreshList = append(globalRefreshList, bo.MQTTGlobalRefresh{
	//		ObjectId: baseUserInfo.UserId,
	//		ObjectValue: bo.BaseUserInfoExtBo{
	//			BaseUserInfoBo: baseUserInfo,
	//			DepartmentId:   depId,
	//		},
	//	})
	//}
	//
	////推送refresh
	//respVo := msgfacade.PushMsgToMqtt(msgvo.PushMsgToMqttReqVo{
	//	Msg: bo.MQTTDataRefreshNotice{
	//		OrgId:         orgId,
	//		Action:        consts.MQTTDataRefreshActionAdd,
	//		Type:          consts.MQTTDataRefreshTypeMember,
	//		OperationId:   operatorId,
	//		GlobalRefresh: globalRefreshList,
	//	}})
	//if respVo.Failure() {
	//	log.Error(respVo.Error())
	//}
}

// 推送更新用户信息事件
func PushUpdateOrgMemberNotice(orgId, depId int64, memberIds []int64, operatorId int64) {
	//globalRefreshList := make([]bo.MQTTGlobalRefresh, 0)
	//baseUserInfos, err := GetBaseUserInfoBatch(orgId, memberIds)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//for _, baseUserInfo := range baseUserInfos {
	//	globalRefreshList = append(globalRefreshList, bo.MQTTGlobalRefresh{
	//		ObjectId: baseUserInfo.UserId,
	//		ObjectValue: bo.BaseUserInfoExtBo{
	//			BaseUserInfoBo: baseUserInfo,
	//			DepartmentId:   depId,
	//		},
	//	})
	//}
	//
	//// 推送refresh
	//respVo := msgfacade.PushMsgToMqtt(msgvo.PushMsgToMqttReqVo{
	//	Msg: bo.MQTTDataRefreshNotice{
	//		OrgId:         orgId,
	//		Action:        consts.MQTTDataRefreshActionModify,
	//		Type:          consts.MQTTDataRefreshTypeMember,
	//		OperationId:   operatorId,
	//		GlobalRefresh: globalRefreshList,
	//	}})
	//if respVo.Failure() {
	//	log.Error(respVo.Error())
	//}
}

func PushRemoveOrgMemberNotice(orgId int64, memberIds []int64, operatorId int64) {
	//globalRefreshList := make([]bo.MQTTGlobalRefresh, 0)
	//
	//baseUserInfos, err := GetBaseUserInfoBatch(orgId, memberIds)
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//for _, baseUserInfo := range baseUserInfos {
	//	globalRefreshList = append(globalRefreshList, bo.MQTTGlobalRefresh{
	//		ObjectId:    baseUserInfo.UserId,
	//		ObjectValue: baseUserInfo,
	//	})
	//}
	//
	////推送refresh
	//respVo := msgfacade.PushMsgToMqtt(msgvo.PushMsgToMqttReqVo{
	//	Msg: bo.MQTTDataRefreshNotice{
	//		OrgId:         orgId,
	//		Action:        consts.MQTTDataRefreshActionDel,
	//		Type:          consts.MQTTDataRefreshTypeMember,
	//		OperationId:   operatorId,
	//		GlobalRefresh: globalRefreshList,
	//	}})
	//if respVo.Failure() {
	//	log.Error(respVo.Error())
	//}
}
