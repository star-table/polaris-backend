package domain

import (
	"strconv"

	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/model/vo"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetUserConfigInfoBatch(orgId int64, userIds []int64) ([]bo.UserConfigBo, errs.SystemErrorInfo) {
	userConfigBoArr := make([]bo.UserConfigBo, 0, len(userIds))
	if len(userIds) < 1 {
		return userConfigBoArr, nil
	}
	list := make([]po.PpmOrgUserConfig, 0)
	dbErr := mysql.SelectAllByCond(consts.TableUserConfig, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &list)
	if dbErr != nil {
		log.Errorf("[GetUserConfigInfoBatch] err: %v, orgId: %d", dbErr, orgId)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}
	_ = copyer.Copy(list, &userConfigBoArr)

	return userConfigBoArr, nil
}

func UpdateUserConfig(orgId, operatorId int64, userConfigBo bo.UserConfigBo) errs.SystemErrorInfo {
	upd := mysql.Upd{}
	if util.IsBool(userConfigBo.DailyReportMessageStatus) {
		upd[consts.TcDailyReportMessageStatus] = userConfigBo.DailyReportMessageStatus
	}
	if util.IsBool(userConfigBo.OwnerRangeStatus) {
		upd[consts.TcOwnerRangeStatus] = userConfigBo.OwnerRangeStatus
	}
	if util.IsBool(userConfigBo.ParticipantRangeStatus) {
		upd[consts.TcParticipantRangeStatus] = userConfigBo.ParticipantRangeStatus
	}
	if util.IsBool(userConfigBo.AttentionRangeStatus) {
		upd[consts.TcAttentionRangeStatus] = userConfigBo.AttentionRangeStatus
	}
	if util.IsBool(userConfigBo.CreateRangeStatus) {
		upd[consts.TcCreateRangeStatus] = userConfigBo.CreateRangeStatus
	}
	if util.IsBool(userConfigBo.RemindMessageStatus) {
		upd[consts.TcRemindMessageStatus] = userConfigBo.RemindMessageStatus
	}
	if util.IsBool(userConfigBo.CommentAtMessageStatus) {
		upd[consts.TcCommentAtMessageStatus] = userConfigBo.CommentAtMessageStatus
	}
	if util.IsBool(userConfigBo.ModifyMessageStatus) {
		upd[consts.TcModifyMessageStatus] = userConfigBo.ModifyMessageStatus
	}
	if util.IsBool(userConfigBo.RelationMessageStatus) {
		upd[consts.TcRelationMessageStatus] = userConfigBo.RelationMessageStatus
	}
	if util.IsBool(userConfigBo.CollaborateMessageStatus) {
		upd[consts.TcCollaborateMessageStatus] = userConfigBo.CollaborateMessageStatus
	}
	if util.IsBool(userConfigBo.DailyProjectReportMessageStatus) {
		upd[consts.TcDailyProjectReportMessageStatus] = userConfigBo.DailyProjectReportMessageStatus
	}
	if userConfigBo.Ext != "" {
		upd[consts.TcExt] = userConfigBo.Ext
	}
	upd[consts.TcRemindExpiring] = userConfigBo.RemindExpiring

	//更新人必填
	upd[consts.TcUpdator] = operatorId

	_, err := mysql.UpdateSmartWithCond(consts.TableUserConfig, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   operatorId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)

	if err != nil {
		//配置更新失败
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateUserPcConfig(orgId, operatorId int64, userConfigBo bo.UserConfigBo) errs.SystemErrorInfo {
	upd := mysql.Upd{}
	if util.IsBool(userConfigBo.PcNoticeOpenStatus) {
		upd[consts.TcPcNoticeOpenStatus] = userConfigBo.PcNoticeOpenStatus
	}
	if util.IsBool(userConfigBo.PcIssueRemindMessageStatus) {
		upd[consts.TcPcIssueRemindMessageStatus] = userConfigBo.PcIssueRemindMessageStatus
	}
	if util.IsBool(userConfigBo.PcOrgMessageStatus) {
		upd[consts.TcPcOrgMessageStatus] = userConfigBo.PcOrgMessageStatus
	}
	if util.IsBool(userConfigBo.PcProjectMessageStatus) {
		upd[consts.TcPcProjectMessageStatus] = userConfigBo.PcProjectMessageStatus
	}
	if util.IsBool(userConfigBo.PcCommentAtMessageStatus) {
		upd[consts.TcPcCommentAtMessageStatus] = userConfigBo.PcCommentAtMessageStatus
	}
	_, err := mysql.UpdateSmartWithCond(consts.TableUserConfig, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   operatorId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, upd)
	if err != nil {
		//配置更新失败
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func UpdateUserDefaultProjectIdConfig(orgId, operatorId int64, userConfigBo bo.UserConfigBo, defaultProjectId int64) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableUserConfig, db.Cond{
		consts.TcId:       userConfigBo.ID,
		consts.TcOrgId:    orgId,
		consts.TcUserId:   operatorId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcDefaultProjectId: defaultProjectId,
	})

	if err != nil {
		//配置更新失败
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

func InsertUserConfig(orgId, userId int64) (*bo.UserConfigBo, errs.SystemErrorInfo) {
	userConfig := &bo.UserConfigBo{}
	//如果不存在就插入
	userIdStr := strconv.FormatInt(userId, 10)
	lockKey := consts.AddUserConfigLock + userIdStr
	suc, err := cache.TryGetDistributedLock(lockKey, userIdStr)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	if suc {
		defer cache.ReleaseDistributedLock(lockKey, userIdStr)

		userConfig, err = getUserConfigInfo(orgId, userId)
		if err != nil {

			inserUserConfigError := insertUserConfig(orgId, userId)

			if inserUserConfigError != nil {

				return nil, inserUserConfigError
			}

		}
	} else {
		userConfig, err = getUserConfigInfo(orgId, userId)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.UserConfigNotExist)
		}
	}

	userConfigBo := &bo.UserConfigBo{}
	err2 := copyer.Copy(userConfig, userConfigBo)
	if err2 != nil {
		log.Error(err2)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err2)
	}
	return userConfigBo, nil
}

//用户配置不存在插入用户配置
func insertUserConfig(orgId, userId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {

	userConfig := &po.PpmOrgUserConfig{}
	userConfigId, err := idfacade.ApplyPrimaryIdRelaxed(userConfig.TableName())
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.ApplyIdError)
	}
	userConfig.Id = userConfigId
	userConfig.OrgId = orgId
	userConfig.UserId = userId
	userConfig.Creator = userId
	userConfig.Updator = userId
	userConfig.IsDelete = consts.AppIsNoDelete
	var err2 error
	if len(tx) > 0 {
		err2 = mysql.TransInsert(tx[0], userConfig)
	} else {
		err2 = mysql.Insert(userConfig)
	}

	if err2 != nil {
		log.Error(err2)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	return nil
}

func SetRemindPopUp(orgId, userId int64, updateFields []string) (*vo.Void, errs.SystemErrorInfo) {
	remindPopUp := consts.NeedRemindPopUp
	if ok, errSlice := slice.Contain(updateFields, consts.RemindPopUp); errSlice == nil && ok {
		remindPopUp = consts.NotNeedRemindPopUp
	}
	// 个人中心 不再提醒弹窗
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId:   userId,
	}
	userConfigPos := []*po.PpmOrgUserConfig{}
	err := mysql.SelectAllByCond(consts.TableUserConfig, cond, &userConfigPos)
	if err != nil {
		log.Errorf("[SetRemindPopUp] err:%v, orgId:%v, userId:%v", err, orgId, userId)
		return nil, errs.MysqlOperateError
	}

	for _, userConfigPo := range userConfigPos {
		userConfigExt := bo.UserConfigExt{}
		if userConfigPo.Ext != "" {
			errJson := json.FromJson(userConfigPo.Ext, &userConfigExt)
			if errJson != nil {
				return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, errJson)
			}
		}
		userConfigExt.RemindPopUp = remindPopUp
		extStr := json.ToJsonIgnoreError(userConfigExt)

		_, err = mysql.UpdateSmartWithCond(consts.TableUserConfig, cond, mysql.Upd{
			consts.TcExt: extStr,
		})
		if err != nil {
			log.Errorf("[SetRemindPopUp] err:%v, orgId:%v, userId:%v", err, orgId, userId)
			return nil, errs.MysqlOperateError
		}

		// 清除缓存
		errCache := DeleteUserConfigInfo(userConfigPo.OrgId, userId)
		if errCache != nil {
			log.Errorf("[SetRemindPopUp] delCache err:%v, orgId:%v, userId:%v", errCache, orgId, userId)
			return nil, errCache
		}
	}

	return &vo.Void{
		ID: userId,
	}, nil
}
