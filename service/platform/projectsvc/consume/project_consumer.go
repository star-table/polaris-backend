package consume

//import (
//	"github.com/star-table/common/core/config"
//	"github.com/star-table/common/core/errors"
//	"github.com/star-table/common/core/model"
//	"github.com/star-table/common/core/util/copyer"
//	"github.com/star-table/common/core/util/json"
//	"github.com/star-table/common/library/mq"
//	"github.com/star-table/polaris-backend/common/core/errs"
//	"github.com/star-table/polaris-backend/common/core/util/stack"
//	"github.com/star-table/polaris-backend/common/model/bo"
//	"github.com/star-table/polaris-backend/facade/msgfacade"
//	"github.com/star-table/polaris-backend/service/platform/projectsvc/notice"
//)
//
//func ProjectTrendsAndNoticeConsume() {
//
//	log.Infof("mq消息-项目动态消费者启动成功")
//
//	projectTrendsTopicConfig := config.GetMqProjectTrendsTopicConfig()
//
//	client := *mq.GetMQClient()
//	_ = client.ConsumeMessage(projectTrendsTopicConfig.Topic, projectTrendsTopicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
//		log.Infof("mq消息-项目动态-信息详情 topic %s, value %s", message.Topic, message.Body)
//
//		projectTrendsBo := &bo.ProjectTrendsBo{}
//		projectMemberChangeBo := &bo.ProjectMemberChangeBo{}
//
//		err := json.FromJson(message.Body, projectTrendsBo)
//		if err != nil {
//			log.Error(err)
//			return errs.BuildSystemErrorInfo(errs.JSONConvertError)
//		}
//		copyer.Copy(projectTrendsBo, projectMemberChangeBo)
//
//		log.Infof("projectMemberChangeBo %v", projectMemberChangeBo)
//
//		//处理通知
//		go func() {
//			defer func() {
//				if r := recover(); r != nil {
//					log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
//				}
//			}()
//			notice.ProjectMemberChangeNotice(*projectMemberChangeBo)
//		}()
//		//go func() {
//		//	defer func() {
//		//		if r := recover(); r != nil {
//		//			log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
//		//		}
//		//	}()
//		//	trendsfacade.AddProjectTrends(trendsvo.AddProjectTrendsReqVo{ProjectTrendsBo: *projectTrendsBo})
//		//}()
//
//		return nil
//	}, func(message *model.MqMessageExt) {
//		//失败的消息处理回调
//		mqMessage := message.MqMessage
//
//		log.Infof("mq消息消费失败-项目动态-信息详情 topic %s, value %s", message.Topic, message.Body)
//
//		projectTrendsBo := &bo.ProjectTrendsBo{}
//		err := json.FromJson(mqMessage.Body, projectTrendsBo)
//		if err != nil {
//			log.Error(err)
//			return
//		}
//
//		pushType := int(projectTrendsBo.PushType)
//		orgId := projectTrendsBo.OrgId
//
//		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, pushType, orgId)
//		if msgErr != nil {
//			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
//		}
//	})
//
//}
