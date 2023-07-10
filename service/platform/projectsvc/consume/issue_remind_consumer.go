package consume

import (
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/mq"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/msgfacade"
)

func IssueRemindConsumer() {
	log.Infof("[IssueRemindConsumer] mq消息-任务提醒通知消费者启动成功")

	if config.GetMQ() == nil {
		log.Error("mq未配置")
		return
	}
	issueRemindConfig := config.GetMQ().Topics.IssueRemind

	if issueRemindConfig.Topic == "" {
		log.Error("[IssueRemindConsumer] mq issueRemind 未配置")
		return
	}

	client := *mq.GetMQClient()

	_ = client.ConsumeMessage(issueRemindConfig.Topic, issueRemindConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		//获取消息实体
		msgBo := &bo.IssueRemindMsg{}
		err := json.FromJson(message.Body, msgBo)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		log.Infof("[IssueRemindConsumer] mq消息-任务提醒通知-信息详情 %s", message.Body)
		// 任务提醒
		IssueRemind(msgBo)

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		msgBo := &bo.IssueRemindMqBo{}
		err := json.FromJson(message.Body, msgBo)
		if err != nil {
			log.Error(err)
			return
		}
		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, int(msgBo.PushType), 0)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})

}

func IssueRemind(msg *bo.IssueRemindMsg) {
}
