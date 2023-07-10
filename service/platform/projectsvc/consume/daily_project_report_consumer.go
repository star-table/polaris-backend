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
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/msgfacade"
	"github.com/jtolds/gls"
)

func DailyProjectReportMsgConsumer() {
	log.Infof("mq消息-每日项目日报Msg消费者启动成功")

	dailyProjectReportMsgTopicConfig := config.GetMqDailyProjectReportMsgTopicConfig()
	client := *mq.GetMQClient()

	_ = client.ConsumeMessage(dailyProjectReportMsgTopicConfig.Topic, dailyProjectReportMsgTopicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		msgBo := &bo.DailyProjectReportMsgBo{}
		err := json.FromJson(message.Body, msgBo)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}

		if msgBo.SourceChannel == "" || msgBo.OutOrgId == "" || len(msgBo.OpenIds) == 0 {
			return nil
		}

		threadlocal.Mgr.SetValues(gls.Values{consts2.TraceIdKey: msgBo.ScheduleTraceId}, func() {
			log.Infof("[DailyProjectReportMsgConsumer] mq消息-每日项目日报msg-信息详情 topic %s, value %s", message.Topic, message.Body)
			SendCardDailyProjectReport(msgBo)
		})

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		msgBo := &bo.DailyProjectReportMsgBo{}
		err := json.FromJson(message.Body, msgBo)
		if err != nil {
			log.Error(err)
			return
		}
		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, int(msgBo.PushType), msgBo.OrgId)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}

// SendCardDailyProjectReport 发送“项目日报”卡片
func SendCardDailyProjectReport(msg *bo.DailyProjectReportMsgBo) {

}
