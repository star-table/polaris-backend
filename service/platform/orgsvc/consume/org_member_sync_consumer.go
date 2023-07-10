package consume

import (
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/mq"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/model/bo/mqbo"
	"github.com/star-table/polaris-backend/facade/msgfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/service"
)

func OrgMemberSyncConsumer() {
	log.Infof("mq消息-任务 OrgMemberSyncConsumer 消费者启动成功")

	topicConfig := config.GetMqSyncMemberDeptTopicConfig()

	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(topicConfig.Topic, topicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()
		log.Infof("mq消息-同步成员、部门信息 topic %s, value %s", message.Topic, message.Body)
		messageBo := &mqbo.PushSyncMemberDeptBo{}
		err := json.FromJson(message.Body, messageBo)
		if err != nil {
			log.Error(err)
			return errs.JSONConvertError
		}
		respErr := service.SyncUserInfoFromFeiShu(messageBo.SyncUserInfoFromFeiShuReqVo.OrgId, messageBo.SyncUserInfoFromFeiShuReqVo.CurrentUserId, messageBo.SyncUserInfoFromFeiShuReqVo.SyncUserInfoFromFeiShuReq)
		if respErr != nil {
			log.Error(respErr)
			return respErr
		}
		log.Infof("同步成员、部门信息成功。")
		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		messageBo := &mqbo.PushSyncMemberDeptBo{}
		err := json.FromJson(message.Body, messageBo)
		if err != nil {
			log.Error(err)
		}

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, messageBo.SyncUserInfoFromFeiShuReqVo.OrgId)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}
