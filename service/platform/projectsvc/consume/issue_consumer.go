package consume

import (
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/mq"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/msgfacade"
)

var log = logger.GetDefaultLogger()

func IssueTrendsAndNoticeConsume() {

	log.Infof("[IssueTrendsAndNoticeConsume] mq消息-任务动态消费者启动成功")

	issueTrendsTopicConfig := config.GetMqIssueTrendsTopicConfig()

	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(issueTrendsTopicConfig.Topic, issueTrendsTopicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		log.Infof("[IssueTrendsAndNoticeConsume] mq消息-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		issueTrendsBo := &bo.IssueTrendsBo{}
		issueNoticeBo := &bo.IssueNoticeBo{}
		err := json.FromJson(message.Body, issueTrendsBo)
		if err != nil {
			log.Errorf("[IssueTrendsAndNoticeConsume] err: %v", err)
			return errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
		}
		err = copyer.Copy(issueTrendsBo, issueNoticeBo)
		if err != nil {
			log.Errorf("[IssueTrendsAndNoticeConsume] err: %v", err)
			return errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
		}

		log.Infof("[IssueTrendsAndNoticeConsume] issueTrendsBo: %v", json.ToJsonIgnoreError(issueTrendsBo))

		if issueTrendsBo.PushType == consts.PushTypeIssueComment || issueTrendsBo.PushType == consts.PushTypeIssueRemarkRemind {
			ext := issueTrendsBo.Ext
			//处理评论通知
			go func() {
				defer func() {
					if r := recover(); r != nil {
						log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
					}
				}()

				content := ext.CommentBo.Content
				//content = consts.TruncateText(consts.GetCardCommentRealWords(content), consts.GrouChatIssueChangeDescLimitPerLine)
				if ext.ResourceInfo != nil && len(ext.ResourceInfo) > 0 {
					content += " [附件]"
				}
				if issueTrendsBo.PushType == consts.PushTypeIssueRemarkRemind {
					content = ext.Remark
				}
			}()

		}

		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		issueTrendsBo := &bo.IssueTrendsBo{}
		err := json.FromJson(mqMessage.Body, issueTrendsBo)
		if err != nil {
			log.Error(err)
			return
		}

		pushType := int(issueTrendsBo.PushType)
		orgId := issueTrendsBo.OrgId

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, pushType, orgId)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}
