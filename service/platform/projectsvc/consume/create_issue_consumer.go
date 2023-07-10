package consume

import (
	"github.com/star-table/common/core/config"
	consts2 "github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/mq"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/model/bo/mqbo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/msgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/service"
	"github.com/jtolds/gls"
)

func BatchCreateIssueConsume() {
	log.Infof("mq消息-批量创建任务消费者启动成功")

	topic := config.GetBatchCreateIssueTopicConfig()

	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(topic.Topic, topic.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		batch := &mqbo.PushBatchCreateIssue{}
		err := json.FromJson(message.Body, batch)
		if err != nil {
			log.Errorf("[BatchCreateIssueConsume] json解析失败, err: %v", err)
			return errs.JSONConvertError
		}
		threadlocal.Mgr.SetValues(gls.Values{consts2.TraceIdKey: batch.TraceId}, func() {
			log.Infof("[BatchCreateIssueConsume] topic: %s, partition: %d, orgId: %d, appId: %d, projectId: %d, tableId: %d, data count: %v",
				message.Topic, message.Partition, batch.Req.OrgId, batch.Req.AppId, batch.Req.ProjectId, batch.Req.TableId, len(batch.Req.Data))

			res, _, _, errSys := service.BatchCreateIssue(batch.Req, batch.WithoutAuth, batch.TriggerBy, batch.AsyncTaskId, batch.Count)
			if errSys != nil {
				log.Errorf("[BatchCreateIssueConsume] service.BatchCreateIssue err: %v", errSys)
				return
			}

			log.Infof("[BatchCreateIssueConsume] 创建任务成功: %s", json.ToJsonIgnoreError(res))
		})

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("[ImportIssueConsume] mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		createIssueReq := &projectvo.CreateIssueReqVo{}
		err := json.FromJson(message.Body, createIssueReq)
		if err != nil {
			log.Error(err)
		}

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, createIssueReq.OrgId)
		if msgErr != nil {
			log.Errorf("[ImportIssueConsume] mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage),
				msgErr)
		}
	})
}
