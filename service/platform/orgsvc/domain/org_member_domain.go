package domain

import (
	"fmt"
	"time"

	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"

	"github.com/star-table/common/library/cache"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/model"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	int64Slice "github.com/star-table/polaris-backend/common/core/util/slice/int64"
	"github.com/star-table/polaris-backend/facade/msgfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

//修改组织成员状态，企业用户状态, 1可用,2禁用,3隐藏
func ModifyOrgMemberStatus(orgId int64, memberIds []int64, status int, operatorId int64) errs.SystemErrorInfo {
	//组织负责人不允许被修改状态
	orgInfo, err := GetBaseOrgInfo(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	orgOwnerId := orgInfo.OrgOwnerId
	memberIds = filterMemberIds(memberIds, operatorId)
	if orgOwnerId != operatorId {
		memberIds = filterMemberIds(memberIds, orgOwnerId)
	}
	if len(memberIds) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcUserId:      db.In(memberIds),
			consts.TcIsDelete:    consts.AppIsNoDelete,
			consts.TcStatus:      db.NotEq(status),
			consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		}, mysql.Upd{
			consts.TcStatus:           status,
			consts.TcStatusChangerId:  operatorId,
			consts.TcUpdator:          operatorId,
			consts.TcStatusChangeTime: time.Now(),
		})
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		////如果更新数量与预期不符，认为动作失败
		//if modifyCount != int64(len(memberIds)) {
		//	return errs.BuildSystemErrorInfo(errs.UpdateMemberStatusFail)
		//}
		//禁用的用户是否可以在选人界面显示，待定
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}
	//最后将用户信息缓存清掉
	clearErr := ClearBaseUserInfoBatch(orgId, memberIds)
	if clearErr != nil {
		log.Error(clearErr)
	}
	return nil
}

//修改组织成员审核状态, 审核状态,1待审核,2审核通过,3审核不过
func ModifyOrgMemberCheckStatus(orgId int64, memberIds []int64, checkStatus int, operatorId int64, isNeedResetCheckStatus bool) errs.SystemErrorInfo {
	//组织负责人不允许被修改状态
	orgInfo, err := GetBaseOrgInfo(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	orgOwnerId := orgInfo.OrgOwnerId
	memberIds = filterMemberIds(memberIds, operatorId)
	if orgOwnerId != operatorId {
		memberIds = filterMemberIds(memberIds, orgOwnerId)
	}
	if len(memberIds) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	condCheckStatus := consts.AppCheckStatusWait
	isCheckPass := checkStatus == consts.AppCheckStatusSuccess
	if isNeedResetCheckStatus {
		condCheckStatus = consts.AppCheckStatusFail
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		upd := mysql.Upd{
			consts.TcCheckStatus: checkStatus,
			consts.TcAuditorId:   operatorId,
			consts.TcUpdator:     operatorId,
			consts.TcAuditTime:   time.Now(),
		}

		if isCheckPass || isNeedResetCheckStatus {
			upd[consts.TcStatus] = consts.AppStatusEnable
		}
		modifyCount, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcUserId:      db.In(memberIds),
			consts.TcCheckStatus: condCheckStatus,
			consts.TcIsDelete:    consts.AppIsNoDelete,
		}, upd)
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		//如果更新数量与预期不符，认为动作失败
		if modifyCount != int64(len(memberIds)) {
			return errs.BuildSystemErrorInfo(errs.UpdateMemberStatusFail)
		}
		////审核通过，加入部门
		//if isCheckPass {
		//	depId, err := BoundOrgMemberToTopDepartment(orgId, memberIds, operatorId)
		//	if err != nil {
		//		log.Error(err)
		//		return err
		//	}
		//	departmentId = depId
		//}
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}
	////推送消息
	//if isCheckPass{
	//	asyn.Execute(func(){
	//		PushAddOrgMemberNotice(orgId, departmentId, memberIds, operatorId)
	//	})
	//}

	//最后将用户信息缓存清掉
	clearErr := ClearBaseUserInfoBatch(orgId, memberIds)
	if clearErr != nil {
		log.Error(clearErr)
	}
	return nil
}

//修改组织成员
func RemoveOrgMember(orgId int64, memberIds []int64, operatorId int64) errs.SystemErrorInfo {
	//组织负责人不允许被修改状态
	orgInfo, err := GetBaseOrgInfo(orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	orgOwnerId := orgInfo.OrgOwnerId
	memberIds = filterMemberIds(memberIds, operatorId)
	isRemoveOwner := false
	for _, id := range memberIds {
		if id == orgOwnerId {
			isRemoveOwner = true
		}
	}

	if len(memberIds) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if isRemoveOwner {
			firstUser, errGet := getFirstUserInOrg(orgId, memberIds)
			if errGet == nil {
				_, dbErr := mysql.TransUpdateSmartWithCond(tx, consts.TableOrganization, db.Cond{
					consts.TcId: orgId,
				}, mysql.Upd{
					consts.TcOwner: firstUser.UserId,
				})
				if dbErr != nil {
					return dbErr
				}
			}
		}

		modifyCount, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcUserId:      db.In(memberIds),
			consts.TcCheckStatus: consts.AppCheckStatusSuccess,
			consts.TcIsDelete:    consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcAuditorId: operatorId,
			consts.TcUpdator:   operatorId,
			consts.TcIsDelete:  consts.AppIsDeleted,
		})
		if err != nil {
			log.Error(err)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		//如果更新数量与预期不符，认为动作失败
		if modifyCount != int64(len(memberIds)) {
			return errs.BuildSystemErrorInfo(errs.UpdateMemberStatusFail)
		}
		//将用户从组织移除之后 - 将该用户从部门移除
		err = UnBoundDepartmentUser(orgId, memberIds, operatorId, tx)
		if err != nil {
			log.Error(err)
			return err
		}
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, transErr)
	}

	removePayUser(orgId, memberIds)

	asyn.Execute(func() {
		PushRemoveOrgMemberNotice(orgId, memberIds, operatorId)
	})

	// 融合版本，清理用户角色 —— 将用户从所有角色中删除。
	// 如果超管离职，则会将普通管理员提升为超管；如果没有普通管理员，则将一个普通成员设为超管。
	for _, oneUserId := range memberIds {
		resp := userfacade.ClearAdminGroupForOneUser(orgId, operatorId, oneUserId)
		if resp.Failure() {
			errMsg := fmt.Sprintf("[RemoveOrgMember]将用户从角色/管理组中移除，并把其他用户设为超管时异常，err: %v；orgId: %d", resp.Error(), orgId)
			log.Errorf("[RemoveOrgMember] err: %s, orgId: %d, oneUserId: %d", errMsg, orgId, oneUserId)
		}
	}
	//最后将用户信息缓存清掉
	clearErr := ClearBaseUserInfoBatch(orgId, memberIds)
	if clearErr != nil {
		log.Error(clearErr)
	}
	return nil
}

// 获取最早进入组织的人
func getFirstUserInOrg(orgId int64, userIds []int64) (*po.PpmOrgUserOrganization, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		return nil, errs.MysqlOperateError
	}
	userOrgPo := &po.PpmOrgUserOrganization{}
	err = conn.Collection(consts.TableUserOrganization).Find(db.Cond{
		consts.TcOrgId:       orgId,
		consts.TcUserId:      db.NotIn(userIds),
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		consts.TcIsDelete:    consts.AppIsNoDelete,
	}).OrderBy("user_id asc").One(userOrgPo)
	if err != nil {
		return nil, errs.MysqlOperateError
	}

	return userOrgPo, nil
}

func removePayUser(orgId int64, userIds []int64) {
	payInfoJson, errCache := cache.HGet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId))
	if errCache != nil {
		log.Errorf("[removePayUser] cache err:%v", errCache)
	}

	if payInfoJson != "" {
		rangeData := bo.PayRangeData{}
		_ = json.FromJson(payInfoJson, &rangeData)
		count := len(rangeData.UserIds)
		for _, id := range userIds {
			rangeData.UserIds = filterMemberIds(rangeData.UserIds, id)
		}
		rangeData.ScopeNum = rangeData.ScopeNum - (count - len(rangeData.UserIds))
		errCache = cache.HSet(consts.CachePayRangeInfo, fmt.Sprintf("%d", orgId), json.ToJsonIgnoreError(rangeData))
		if errCache != nil {
			log.Errorf("[removePayUser] cache err:%v", errCache)
		}
	}
}

//过滤掉当前操作人
func filterMemberIds(memberIds []int64, operatorId int64) []int64 {
	memberIds = slice.SliceUniqueInt64(memberIds)
	newMemberIds := make([]int64, 0)
	for _, memberId := range memberIds {
		if memberId != operatorId {
			newMemberIds = append(newMemberIds, memberId)
		}
	}
	return newMemberIds
}

//通过渠道获取组织用户信息
func GetOrgUserInfoListBySourceChannel(orgId int64, sourceChannel string, page, size int) ([]bo.OrgUserInfo, int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, 0, errs.MysqlOperateError
	}

	mid := conn.Select(
		"userOrg.user_id as user_id",
		"userOutInfo.out_user_id as out_user_id",
		"userOrg.org_id as org_id",
		"userOrg.status as org_user_status",
		"userOrg.check_status as org_user_check_status",
	).
		From("ppm_org_user_organization userOrg").
		LeftJoin("ppm_org_user_out_info userOutInfo").
		On("userOrg.user_id = userOutInfo.user_id").
		Where(db.Cond{
			"userOrg.org_id":             orgId,
			"userOrg.is_delete":          consts.AppIsNoDelete,
			"userOutInfo.is_delete":      consts.AppIsNoDelete,
			"userOutInfo.source_channel": sourceChannel,
			"userOutInfo.org_id":         orgId,
		})

	result := &[]bo.OrgUserInfo{}
	total := int64(0)
	if size > 0 && page > 0 {
		pageResult := mid.Paginate(uint(size)).Page(uint(page))
		rowSize, err := pageResult.TotalEntries()
		if err != nil {
			log.Error(err)
			return nil, 0, errs.MysqlOperateError
		}
		total = int64(rowSize)
		err = pageResult.All(result)
		if err != nil {
			log.Error(err)
			return nil, 0, errs.MysqlOperateError
		}
	} else {
		err := mid.All(result)
		if err != nil {
			log.Error(err)
			return nil, 0, errs.MysqlOperateError
		}
		total = int64(len(*result))
	}
	return *result, total, nil
}

func FeishuMemberHelpMsgToMq(feishuHelpObjectBo bo.FeishuHelpObjectBo) {
	message, err := json.ToJson(feishuHelpObjectBo)
	if err != nil {
		log.Error(err)
	}
	mqMessage := &model.MqMessage{
		Topic:          config.GetHelpMessagePushToFeiShu().Topic,
		Keys:           uuid.NewUuid(),
		Body:           message,
		DelayTimeLevel: 3,
	}
	msgErr := msgfacade.PushMsgToMqRelaxed(*mqMessage, 0, feishuHelpObjectBo.OrgId)
	if msgErr != nil {
		log.Errorf("mq消息推送失败，入表失败，消息内容：%s, 失败信息: %v", json.ToJsonIgnoreError(mqMessage), msgErr)
	}
}

// GetBaseUserInfoByOpenIdBatch 通过 openId 批量查询用户信息
func GetBaseUserInfoByOpenIdBatch(orgId int64, openIds []string) ([]bo.BaseUserInfoBo,
	errs.SystemErrorInfo) {
	returnArr := make([]bo.BaseUserInfoBo, 0)
	outInfoArr := make([]po.PpmOrgUserOutInfo, 0)
	if len(openIds) < 1 {
		return returnArr, nil
	}
	cond := db.Cond{
		consts.TcOutUserId: db.In(openIds),
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcStatus:    consts.AppStatusEnable,
	}
	// 一些调用方可能会传 0，这里兼容一下。
	if orgId > 0 {
		cond[consts.TcOrgId] = orgId
	}
	err := mysql.SelectAllByCond(consts.TableUserOutInfo, cond, &outInfoArr)
	if err != nil {
		log.Errorf("[GetBaseUserInfoByEmpIdBatch] err: %v, orgId: %d", err, orgId)
		return returnArr, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, outUser := range outInfoArr {
		returnArr = append(returnArr, bo.BaseUserInfoBo{
			UserId:             outUser.UserId,
			OutUserId:          outUser.OutUserId,
			OrgId:              outUser.OrgId,
			OutOrgId:           "",
			Name:               outUser.Name,
			NamePy:             "",
			Avatar:             outUser.Avatar,
			HasOutInfo:         false,
			HasOrgOutInfo:      false,
			OutOrgUserId:       outUser.OutUserId,
			OrgUserIsDelete:    outUser.IsDelete,
			OrgUserStatus:      0,
			OrgUserCheckStatus: 0,
		})
	}

	return returnArr, nil
}

func GetOrgMemberBaseInfoListByUser(orgId, userId int64) (*bo.OrgMemberBaseInfoBo, errs.SystemErrorInfo) {
	conn, dbErr := mysql.GetConnect()
	if dbErr != nil {
		log.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	var baseUserInfo *bo.OrgMemberBaseInfoBo
	conds := db.Cond{
		"user." + consts.TcOrgId:       orgId,
		"userOrg." + consts.TcOrgId:    orgId,
		"user." + consts.TcIsDelete:    consts.AppIsNoDelete,
		"userOrg." + consts.TcIsDelete: consts.AppIsNoDelete,
		"userOrg." + consts.TcUserId:   userId,
		"user." + consts.TcId:          userId,
	}
	dbErr = conn.Select(
		"user.id as user_id",
		"user.login_name",
		"user.name",
		"user.name_pinyin",
		"user.sex",
		"user.birthday",
		"user.email",
		"user.mobile_region",
		"user.mobile",
		"user.avatar as avatar",
		"user.language",
		"userOrg.check_status",
		"userOrg.status as status",
		"userOrg.status_change_time",
		"userOrg.use_status",
		"userOrg.auditor_id",
		"userOrg.audit_time",
		"userOrg.emp_no",
		"userOrg.weibo_ids",
		"userOrg.creator",
		"userOrg.create_time",
		"userOrg.updator",
		"userOrg.update_time",
		"userOrg.is_delete as is_delete",
		"userOrg.type as type").
		From("ppm_org_user_organization userOrg", "ppm_org_user user").
		Where(conds).One(&baseUserInfo)
	if dbErr != nil {
		log.Error(dbErr)
		return nil, errs.MysqlOperateError
	}
	return baseUserInfo, nil
}

func GetOrgMemberInfoByUserId(orgId, userId int64) (*bo.UserInfoBo, errs.SystemErrorInfo) {
	var userInfoBo *bo.UserInfoBo
	userInfoBo, errSys := GetCurrentUserInfoByUserId(userId)
	if errSys != nil {
		if errSys == errs.UserNotExist {
			// 尝试查询user表
			userPo, err := GetUserInfoByUserId(userId)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			userInfoBo = &bo.UserInfoBo{}
			copyErr := copyer.Copy(userPo, userInfoBo)
			if copyErr != nil {
				log.Error(copyErr)
			}
		} else {
			log.Errorf("[GetOrgMemberInfoByUserId] GetCurrentUserInfoByUserId err:%v, userId:%v", errSys, userId)
			return nil, errSys
		}
	}
	userOrgPo := po.PpmOrgUserOrganization{}
	err := mysql.SelectOneByCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId:   userId,
	}, &userOrgPo)
	if err != nil {
		log.Errorf("[GetOrgMemberInfoByUserId] err:%v, orgId:%v, userId:%v", err, orgId, userId)
		return nil, errs.MysqlOperateError
	}
	userInfoBo.UserStatus = userOrgPo.Status
	return userInfoBo, nil
}

func UpdateOrgMemberInfo(orgId, userId int64, userUpd, orgUserUpd mysql.Upd) errs.SystemErrorInfo {
	err := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(userUpd) > 0 {
			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUser, db.Cond{
				consts.TcIsDelete: consts.AppIsNoDelete,
				consts.TcId:       userId,
			}, userUpd)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		if len(orgUserUpd) > 0 {
			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{
				consts.TcOrgId:    orgId,
				consts.TcIsDelete: consts.AppIsNoDelete,
				consts.TcUserId:   userId,
			}, orgUserUpd)
			if err != nil {
				log.Error(err)
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Errorf("[UpdateOrgMemberInfo] err:%v", err)
		return errs.MysqlOperateError
	}
	return nil
}

func GetOrgInfosByUserIds(userIds []int64) ([]*po.PpmOrgOrganization, errs.SystemErrorInfo) {
	userOrgPos := []*po.PpmOrgUserOrganization{}
	orgPos := []*po.PpmOrgOrganization{}
	err := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId:   db.In(userIds),
	}, &userOrgPos)
	if err != nil {
		log.Errorf("[GetOrgInfosByUserIds] err:%v, userIds:%v", err, userIds)
		return nil, errs.MysqlOperateError
	}
	orgIds := []int64{}
	for _, u := range userOrgPos {
		orgIds = append(orgIds, u.OrgId)
	}
	err = mysql.SelectAllByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(orgIds),
	}, &orgPos)
	if err != nil {
		log.Errorf("[GetOrgInfosByUserIds] err:%v, userIds:%v", err, userIds)
		return nil, errs.MysqlOperateError
	}
	return orgPos, nil
}

// userId: orgIds
func GetLocalOrgUserIdMap(userIds []int64) (map[int64][]int64, errs.SystemErrorInfo) {
	userOrgPos := []*po.PpmOrgUserOrganization{}
	outOrgPos := []*po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId:   db.In(userIds),
	}, &userOrgPos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	allOrgIds := []int64{}
	for _, u := range userOrgPos {
		allOrgIds = append(allOrgIds, u.OrgId)
	}
	allOrgIds = slice.SliceUniqueInt64(allOrgIds)
	err = mysql.SelectAllByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    db.In(allOrgIds),
	}, &outOrgPos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	thirdPlatformOrgIds := []int64{}
	for _, o := range outOrgPos {
		thirdPlatformOrgIds = append(thirdPlatformOrgIds, o.OrgId)
	}
	needOrgIds := int64Slice.ArrayDiff(allOrgIds, thirdPlatformOrgIds)
	userOrgIdsMap := make(map[int64][]int64)
	for _, u := range userOrgPos {
		if int64Slice.InArray(u.OrgId, needOrgIds) {
			userOrgIdsMap[u.UserId] = append(userOrgIdsMap[u.UserId], u.OrgId)
		}
	}

	return userOrgIdsMap, nil
}

func UpdateLocalOrgUserNames(orgIds []int64, userId int64, name string) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableUser, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       userId,
		consts.TcOrgId:    db.In(orgIds),
	}, mysql.Upd{consts.TcName: name})
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	return nil
}

func CheckUserAccountByLoginName(orgId int64, loginName string) bool {
	userPo := &po.PpmOrgUser{}
	err := mysql.SelectOneByCond(consts.TableUser, db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcLoginName: loginName,
	}, userPo)
	if err != nil {
		log.Errorf("[CheckUserAccountByLoginName] err:%v, orgId:%v, loginName:%v", err, orgId, loginName)
		return false
	}

	return true
}
