package consume

func ImportIssueConsume() {

	//log.Infof("mq消息-任务动态消费者启动成功")
	//
	//importIssueTopicConfig := config.GetMqImportIssueTopicConfig()
	//
	//client := *mq.GetMQClient()
	//_ = client.ConsumeMessage(importIssueTopicConfig.Topic, importIssueTopicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
	//		}
	//	}()
	//
	//	createIssueBo := &mqbo.PushCreateIssueBo{}
	//	err := json.FromJson(message.Body, createIssueBo)
	//	if err != nil {
	//		log.Infof("[ImportIssueConsume] mq消息-动态-信息详情 topic %s, value %s", message.Topic, message.Body)
	//		log.Errorf("[ImportIssueConsume] err: %v", err)
	//		return errs.JSONConvertError
	//	}
	//	threadlocal.Mgr.SetValues(gls.Values{consts2.TraceIdKey: createIssueBo.TraceId}, func() {
	//		log.Infof("[ImportIssueConsume] mq消息-动态-信息详情 topic: %s, body: %s", message.Topic, message.Body)
	//		log.Infof("[ImportIssueConsume][###] mq消息-动态-信息详情 topic: %s, partition: %d, orgId: %d, projectId: %d, issueId: %d",
	//			message.Topic, message.Partition, createIssueBo.CreateIssueReqVo.OrgId, createIssueBo.CreateIssueReqVo.CreateIssue.ProjectID, createIssueBo.IssueId)
	//		issueVo, _, respErr := service.CreateIssueWithId(createIssueBo.CreateIssueReqVo, createIssueBo.IssueId, false)
	//		if respErr != nil {
	//			// 如果异常，则将异常标记存入异步任务缓存中
	//			domain.UpdateAsyncTaskWithError(createIssueBo.CreateIssueReqVo.OrgId,
	//				createIssueBo.CreateIssueReqVo.CreateIssue.ProjectID,
	//				cast.ToInt64(createIssueBo.CreateIssueReqVo.CreateIssue.TableID),
	//				createIssueBo.CreateIssueReqVo.UserId, respErr, false, 1, true)
	//			log.Errorf("[ImportIssueConsume] service.CreateIssueWithId err: %v, issueId: %d", respErr, createIssueBo.IssueId)
	//			return
	//		} else {
	//			// 处理完，更新 redis 中的异步任务进度
	//			domain.UpdateAsyncTaskWithSucc(&createIssueBo.CreateIssueReqVo, 1, true)
	//		}
	//
	//		log.Infof("[ImportIssueConsume] 创建任务成功 %s", json.ToJsonIgnoreError(issueVo))
	//	})
	//
	//	return nil
	//}, func(message *model.MqMessageExt) {
	//	//失败的消息处理回调
	//	mqMessage := message.MqMessage
	//
	//	log.Infof("[ImportIssueConsume] mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)
	//
	//	createIssueReq := &projectvo.CreateIssueReqVo{}
	//	err := json.FromJson(message.Body, createIssueReq)
	//	if err != nil {
	//		log.Error(err)
	//	}
	//
	//	msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, createIssueReq.OrgId)
	//	if msgErr != nil {
	//		log.Errorf("[ImportIssueConsume] mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage),
	//			msgErr)
	//	}
	//})
}
