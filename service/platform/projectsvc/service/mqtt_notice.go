package service

func PushAddProjectNotice(orgId, projectId int64, operatorId int64) {
	//projectInfo, err := ProjectInfo(orgId, operatorId, vo.ProjectInfoReq{ProjectID: projectId})
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	////推送refresh
	//respVo := msgfacade.PushMsgToMqtt(msgvo.PushMsgToMqttReqVo{
	//	Msg: bo.MQTTDataRefreshNotice{
	//		OrgId:       orgId,
	//		Action:      consts.MQTTDataRefreshActionAdd,
	//		Type:        consts.MQTTDataRefreshTypePro,
	//		OperationId: operatorId,
	//		GlobalRefresh: []bo.MQTTGlobalRefresh{
	//			{
	//				ObjectId:    projectId,
	//				ObjectValue: projectInfo,
	//			},
	//		},
	//	}})
	//if respVo.Failure() {
	//	log.Error(respVo.Error())
	//}
}

func PushModifyProjectNotice(orgId, projectId int64, operatorId int64) {
	//projectInfo, err := ProjectInfo(orgId, operatorId, vo.ProjectInfoReq{ProjectID: projectId})
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//respVo := msgfacade.PushMsgToMqtt(msgvo.PushMsgToMqttReqVo{
	//	Msg: bo.MQTTDataRefreshNotice{
	//		OrgId:       orgId,
	//		Action:      consts.MQTTDataRefreshActionModify,
	//		Type:        consts.MQTTDataRefreshTypePro,
	//		OperationId: operatorId,
	//		GlobalRefresh: []bo.MQTTGlobalRefresh{
	//			{
	//				ObjectId:    projectId,
	//				ObjectValue: projectInfo,
	//			},
	//		},
	//	}})
	//if respVo.Failure() {
	//	log.Error(respVo.Error())
	//}
	//
	//appId, _ := domain.GetAppIdFromProjectId(orgId, projectId)
	////推送refresh
	//respVo = msgfacade.PushMsgToMqtt(msgvo.PushMsgToMqttReqVo{
	//	Msg: bo.MQTTDataRefreshNotice{
	//		OrgId:       orgId,
	//		ProjectId:   projectId,
	//		AppId:       appId,
	//		Action:      consts.MQTTDataRefreshActionModify,
	//		Type:        consts.MQTTDataRefreshTypePro,
	//		OperationId: operatorId,
	//		GlobalRefresh: []bo.MQTTGlobalRefresh{
	//			{
	//				ObjectId:    projectId,
	//				ObjectValue: projectInfo,
	//			},
	//		},
	//	}})
	//if respVo.Failure() {
	//	log.Error(respVo.Error())
	//}
}

func PushDeleteProjectNotice(orgId, projectId int64, operatorId int64) {
	//projectInfo, err := ProjectInfo(orgId, operatorId, vo.ProjectInfoReq{ProjectID: projectId})
	//if err != nil {
	//	log.Error(err)
	//	return
	//}
	//
	//respVo := msgfacade.PushMsgToMqtt(msgvo.PushMsgToMqttReqVo{
	//	Msg: bo.MQTTDataRefreshNotice{
	//		OrgId:       orgId,
	//		Action:      consts.MQTTDataRefreshActionModify,
	//		Type:        consts.MQTTDataRefreshTypePro,
	//		OperationId: operatorId,
	//		GlobalRefresh: []bo.MQTTGlobalRefresh{
	//			{
	//				ObjectId:    projectId,
	//				ObjectValue: projectInfo,
	//			},
	//		},
	//	}})
	//if respVo.Failure() {
	//	log.Error(respVo.Error())
	//}
	//
	//appId, _ := domain.GetAppIdFromProjectId(orgId, projectId)
	////推送refresh
	//respVo = msgfacade.PushMsgToMqtt(msgvo.PushMsgToMqttReqVo{
	//	Msg: bo.MQTTDataRefreshNotice{
	//		OrgId:       orgId,
	//		ProjectId:   projectId,
	//		AppId:       appId,
	//		Action:      consts.MQTTDataRefreshActionDel,
	//		Type:        consts.MQTTDataRefreshTypePro,
	//		OperationId: operatorId,
	//		GlobalRefresh: []bo.MQTTGlobalRefresh{
	//			{
	//				ObjectId:    projectId,
	//				ObjectValue: projectInfo,
	//			},
	//		},
	//	}})
	//if respVo.Failure() {
	//	log.Error(respVo.Error())
	//}
}
