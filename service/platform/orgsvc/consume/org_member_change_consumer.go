package consume

import (
	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/errors"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/mq"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/stack"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/msgfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
)

var log = logger.GetDefaultLogger()

//暂时忽略处理的组织
var ignoreOrgIds = []int64{
	16317,
}

func OrgMemberChangeConsume() {

	log.Infof("mq消息-组织成员变动消费者启动成功")

	orgMemberChangeTopicConfig := config.GetMqOrgMemberChangeConfig()

	client := *mq.GetMQClient()
	_ = client.ConsumeMessage(orgMemberChangeTopicConfig.Topic, orgMemberChangeTopicConfig.GroupId, func(message *model.MqMessageExt) errors.SystemErrorInfo {
		defer func() {
			if r := recover(); r != nil {
				log.Error(errs.BuildSystemErrorInfoWithPanicRecover(r, stack.GetStack()))
			}
		}()

		log.Infof("mq消息-组织成员变动消费信息 topic %s, value %s", message.Topic, message.Body)

		orgMemberChange := &bo.OrgMemberChangeBo{}
		err := json.FromJson(message.Body, orgMemberChange)
		if err != nil {
			log.Error(err)
			return errs.JSONConvertError
		}

		orgId := orgMemberChange.OrgId

		var businessErr errs.SystemErrorInfo = nil

		changeType := orgMemberChange.ChangeType
		//业务处理
		switch changeType {
		//禁用
		case consts.OrgMemberChangeTypeDisable:
			////暂时先不禁用
			//return nil
			businessErr = domain.ModifyOrgMemberStatus(orgId, []int64{orgMemberChange.UserId}, consts.AppStatusHidden, 0)
		//启用
		case consts.OrgMemberChangeTypeEnable:
			businessErr = domain.ModifyOrgMemberStatus(orgId, []int64{orgMemberChange.UserId}, consts.AppStatusEnable, 0)
		//添加用户
		case consts.OrgMemberChangeTypeAdd, consts.OrgMemberChangeTypeAddDisable:
			err := addUser(orgMemberChange)
			if err != nil {
				return err
			}
		case consts.OrgMemberChangeTypeRemove:
			businessErr = domain.RemoveOrgMember(orgId, []int64{orgMemberChange.UserId}, 0)
		case consts.OrgMemberChangeTypeUpdate: // 更新用户信息的回调处理
			err := updateUser(orgMemberChange)
			if err != nil {
				log.Error(err)
				return err
			}
		}

		if businessErr != nil {
			log.Error(businessErr)
		}

		//在并发操作时，有几率更新失败，所以忽略异常
		return nil
	}, func(message *model.MqMessageExt) {
		//失败的消息处理回调
		mqMessage := message.MqMessage

		log.Infof("mq消息消费失败-动态-信息详情 topic %s, value %s", message.Topic, message.Body)

		orgMemberChange := &bo.OrgMemberChangeBo{}
		err := json.FromJson(message.Body, orgMemberChange)
		if err != nil {
			log.Error(err)
			return
		}

		msgErr := msgfacade.InsertMqConsumeFailMsgRelaxed(mqMessage, 0, orgMemberChange.OrgId)
		if msgErr != nil {
			log.Errorf("mq消息消费失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
		}
	})
}

func addUser(orgMemberChange *bo.OrgMemberChangeBo) errors.SystemErrorInfo {

	return nil
}

func updateUser(orgMemberChange *bo.OrgMemberChangeBo) errs.SystemErrorInfo {
	return nil
}
