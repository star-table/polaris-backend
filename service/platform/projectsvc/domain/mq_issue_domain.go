package domain

import (
	"strconv"
	"time"

	"github.com/star-table/polaris-backend/common/core/consts"

	"github.com/spf13/cast"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/bo/mqbo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/facade/msgfacade"
	"github.com/star-table/polaris-backend/facade/trendsfacade"
)

func PushIssueTrends(issueTrendsBo *bo.IssueTrendsBo) {
	issueTrendsBo.OperateTime = time.Now()
	mqKeys := strconv.FormatInt(issueTrendsBo.OrgId, 10)
	//动态改成同步的
	resp := trendsfacade.AddIssueTrends(trendsvo.AddIssueTrendsReqVo{IssueTrendsBo: *issueTrendsBo, Key: mqKeys})
	if resp.Failure() {
		log.Error(resp.Message)
	}
}

func PushIssueThirdPlatformNotice(issueTrendsBo *bo.IssueTrendsBo) {
	issueTrendsBo.OperateTime = time.Now()
	mqKeys := strconv.FormatInt(issueTrendsBo.OrgId, 10)

	orgId := issueTrendsBo.OrgId
	pushType := int(issueTrendsBo.PushType)
	issueTrendsBo.OperateTime = time.Now()
	message, err := json.ToJson(issueTrendsBo)
	if err != nil {
		log.Error(err)
	}

	mqMessage := &model.MqMessage{
		Topic:          config.GetMqIssueTrendsTopicConfig().Topic,
		Keys:           mqKeys,
		Body:           message,
		DelayTimeLevel: 3,
	}

	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, pushType, orgId)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

func PushBatchCreateIssue(batch *mqbo.PushBatchCreateIssue) {
	message, err := json.ToJson(batch)
	if err != nil {
		log.Error(err)
	}

	//这里key使用项目id，保证同一项目下导入的任务顺序的有效性
	mqMessage := &model.MqMessage{
		Topic:          config.GetBatchCreateIssueTopicConfig().Topic,
		Keys:           cast.ToString(batch.Req.ProjectId),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, int(consts.PushTypeCreateIssue), batch.Req.OrgId)
	if msgErr != nil {
		log.Errorf("[PushBatchCreateIssue] 消息推送失败 消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

func PushCreateIssue(createIssueBo mqbo.PushCreateIssueBo) {
	message, err := json.ToJson(createIssueBo)
	if err != nil {
		log.Error(err)
	}

	reqVo := createIssueBo.CreateIssueReqVo
	//这里key使用项目id，保证同一项目下导入的任务顺序的有效性
	mqMessage := &model.MqMessage{
		Topic:          config.GetMqImportIssueTopicConfig().Topic,
		Keys:           strconv.FormatInt(reqVo.CreateIssue.ProjectID, 10),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, reqVo.OrgId)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

func PushProjectTrends(projectMemberChangeBo bo.ProjectTrendsBo) {
	projectMemberChangeBo.OperateTime = time.Now()
	trendsfacade.AddProjectTrends(trendsvo.AddProjectTrendsReqVo{ProjectTrendsBo: projectMemberChangeBo})
}

func PushProjectThirdPlatformNotice(projectMemberChangeBo bo.ProjectTrendsBo) {
	projectMemberChangeBo.OperateTime = time.Now()

	orgId := projectMemberChangeBo.OrgId
	pushType := int(projectMemberChangeBo.PushType)

	message, err := json.ToJson(projectMemberChangeBo)
	if err != nil {
		log.Error(err)
	}
	mqKeys := strconv.FormatInt(projectMemberChangeBo.OrgId, 10)

	mqMessage := &model.MqMessage{
		Topic:          config.GetMqProjectTrendsTopicConfig().Topic,
		Keys:           mqKeys,
		Body:           message,
		DelayTimeLevel: 3,
	}

	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, pushType, orgId)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}
