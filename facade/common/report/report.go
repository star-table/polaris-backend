package report

import (
	"fmt"
	"time"

	msgPb "github.com/star-table/interface/golang/msg/v1"
	pushPb "github.com/star-table/interface/golang/push/v1"
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/mq"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/convert"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade/pushfacade"
	"github.com/spf13/cast"
)

var log = logger.GetDefaultLogger()

func formatKafkaKey(orgId, appId int64) string {
	if appId > 0 {
		// 保证同一个APP内的事件有序
		return fmt.Sprintf("%d_%d", orgId, appId)
	} else {
		return cast.ToString(orgId)
	}
}

func ReportDataEvent(eventType msgPb.EventType, traceId string, dataEvent *commonvo.DataEvent, mqttFlag ...bool) {
	e := &commonvo.Event{}
	e.Category = msgPb.EventCategory_Data.String()
	e.Type = eventType.String()
	e.Timestamp = time.Now().UnixNano()
	e.TraceId = traceId
	e.Payload = dataEvent

	body, err := json.Marshal(e)
	if err != nil {
		log.Errorf("[ReportDataEvent] Marshal err: %v", err)
		return
	}

	reportEventToKafka(dataEvent.OrgId, dataEvent.AppId, body)
	if len(mqttFlag) == 0 || mqttFlag[0] {
		reportEventToMqtt(dataEvent.OrgId, dataEvent.ProjectId, body, eventType)
	}
}

func ReportOrgEvent(eventType msgPb.EventType, traceId string, orgEvent *commonvo.OrgEvent, mqttFlag ...bool) {
	e := &commonvo.Event{}
	e.Category = msgPb.EventCategory_Org.String()
	e.Type = eventType.String()
	e.Timestamp = time.Now().UnixNano()
	e.TraceId = traceId
	e.Payload = orgEvent

	body, err := json.Marshal(e)
	if err != nil {
		log.Errorf("[ReportDataEvent] Marshal err: %v", err)
		return
	}

	reportEventToKafka(orgEvent.OrgId, 0, body)
	if len(mqttFlag) == 0 || mqttFlag[0] {
		reportEventToMqtt(orgEvent.OrgId, 0, body, eventType)
	}
}

func ReportAppEvent(eventType msgPb.EventType, traceId string, appEvent *commonvo.AppEvent) {
	e := &commonvo.Event{}
	e.Category = msgPb.EventCategory_App.String()
	e.Type = eventType.String()
	e.Timestamp = time.Now().UnixNano()
	e.TraceId = traceId
	e.Payload = appEvent

	body, err := json.Marshal(e)
	if err != nil {
		log.Errorf("[ReportAppEvent] Marshal err: %v", err)
		return
	}

	reportEventToKafka(appEvent.OrgId, appEvent.AppId, body)
	reportEventToMqtt(appEvent.OrgId, appEvent.ProjectId, body, eventType)
}

func ReportTableEvent(eventType msgPb.EventType, traceId string, tableEvent *commonvo.TableEvent) {
	e := &commonvo.Event{}
	e.Category = msgPb.EventCategory_Table.String()
	e.Type = eventType.String()
	e.Timestamp = time.Now().UnixNano()
	e.TraceId = traceId
	e.Payload = tableEvent

	body, err := json.Marshal(e)
	if err != nil {
		log.Errorf("[ReportTableEvent] Marshal err: %v", err)
		return
	}

	reportEventToKafka(tableEvent.OrgId, tableEvent.AppId, body)
	reportEventToMqtt(tableEvent.OrgId, tableEvent.ProjectId, body, eventType)
}

func ReportUserEvent(eventType msgPb.EventType, traceId string, userEvent *commonvo.UserEvent, mqttFlag ...bool) {
	e := &commonvo.Event{}
	e.Category = msgPb.EventCategory_User.String()
	e.Type = eventType.String()
	e.Timestamp = time.Now().UnixNano()
	e.TraceId = traceId
	e.Payload = userEvent

	body, err := json.Marshal(e)
	if err != nil {
		log.Errorf("[ReportUserEvent] Marshal err: %v", err)
		return
	}

	reportEventToKafka(userEvent.OrgId, 0, body)
	if len(mqttFlag) == 0 || mqttFlag[0] {
		reportEventToMqtt(userEvent.OrgId, 0, body, eventType)
	}
}

/////////////////////////////////////////////////////////////////////////////////////////////////////

func reportEventToKafka(orgId, appId int64, body []byte) errs.SystemErrorInfo {
	topic := config.GetMQ().Topics.LogEvent.Topic

	client := *mq.GetMQClient()
	if client == nil {
		log.Errorf("[reportEventToKafka] GetMQClient failed, orgId: %v, appId: %v, msg: %q",
			orgId, appId, body)
		return errs.KafkaMqSendMsgError
	}

	repushTimes := 3
	reconsumeTimes := 0

	msg := &model.MqMessage{}
	msg.Topic = topic
	msg.Keys = formatKafkaKey(orgId, appId)
	msg.Body = convert.UnsafeBytesToString(body)
	msg.RePushTimes = &repushTimes
	msg.ReconsumeTimes = &reconsumeTimes
	_, sysErr := client.PushMessage(msg)
	if sysErr != nil {
		log.Errorf("[reportEventToKafka] PushMessage failed: %v, orgId: %v, appId: %v, msg: %q",
			sysErr, orgId, appId, body)
		return sysErr
	}
	return nil
}

func reportEventToMqtt(orgId, projectId int64, body []byte, eventType msgPb.EventType) errs.SystemErrorInfo {
	channel := ""
	if projectId == 0 || eventType == msgPb.EventType_DataMoved { // 数据移动可能会跨app，因此推送到组织级别
		channel = util.GetMQTTOrgChannel(orgId)
	} else {
		channel = util.GetMQTTProjectChannel(orgId, projectId)
	}

	resp := pushfacade.PushMqtt(&pushPb.PushMqttReq{
		Channel: channel,
		Body:    body,
	})
	if resp.Failure() {
		log.Errorf("[reportEventToMqtt] PushMqtt failed: %v, orgId: %v, projectId: %v, msg: %q",
			resp.Err.Error(), orgId, projectId, body)
		return resp.Err.Error()
	}
	return nil
}
