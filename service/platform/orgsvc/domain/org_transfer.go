package domain

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

// GetTransferOrgUserMatchInfo 获取两个组织里面的用户分别都绑定同一个手机的map，newUserId->originUserId
func GetTransferOrgUserMatchInfo(originOrgId, newOrgId int64) (*bo.TransferMatchInfo, errs.SystemErrorInfo) {
	matchInfo := &bo.TransferMatchInfo{}
	originUsers, err := GetUserOrgInfos(originOrgId)
	if err != nil {
		return matchInfo, err
	}
	newUsers, err := GetUserOrgInfos(newOrgId)
	if err != nil {
		return matchInfo, err
	}
	if len(originUsers) == 0 || len(newUsers) == 0 {
		return matchInfo, nil
	}

	originIdsMap := make(map[int64]bool, len(originUsers))
	allUserIds := make([]int64, 0, len(originUsers)+len(newUsers))
	for _, user := range originUsers {
		matchInfo.OriginUserIds = append(matchInfo.OriginUserIds, user.UserId)
		allUserIds = append(allUserIds, user.UserId)
		originIdsMap[user.UserId] = true
	}
	for _, user := range newUsers {
		matchInfo.NewUserIds = append(matchInfo.NewUserIds, user.UserId)
		allUserIds = append(allUserIds, user.UserId)
	}

	// 获取绑定关系，如果一个globalId下有两个userId，则证明新老绑定在了一起，可以替换第三方信息
	globalIdToUserIdsMap, dbErr := dao.GetGlobalUserRelation().GetGlobalUserIdsMapByUserIds(allUserIds)
	if dbErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	matchInfo.NewToOriginUserIdMap = make(map[int64]int64, len(globalIdToUserIdsMap))
	for _, userIds := range globalIdToUserIdsMap {
		if len(userIds) == 2 {
			// 判断下这两个userId中哪个是原始的，哪个是新组织的
			if originIdsMap[userIds[0]] {
				matchInfo.NewToOriginUserIdMap[userIds[1]] = userIds[0]
			} else if originIdsMap[userIds[1]] {
				matchInfo.NewToOriginUserIdMap[userIds[0]] = userIds[1]
			}
		}
	}

	return matchInfo, nil
}

// TransferOrgToOtherPlatform 需要将originOrgId相关的第三方信息都替换成newOrgId下的
// 比如钉钉的信息替换飞书的信息，这个时候其实是客户从飞书转到了钉钉，替换完后，飞书再次同步就会变成一个新组织
func TransferOrgToOtherPlatform(originOrgId, newOrgId int64) errs.SystemErrorInfo {
	transferMatchInfo, err := GetTransferOrgUserMatchInfo(originOrgId, newOrgId)
	if err != nil {
		return err
	}
	if len(transferMatchInfo.NewToOriginUserIdMap) == 0 {
		return nil
	}

	transferOrgInfo, err := getTransferOrgInfo(originOrgId, newOrgId, transferMatchInfo.NewToOriginUserIdMap)
	if err != nil {
		log.Errorf("[TransferOrgToOtherPlatform] getTransferOrgInfo originOrgId:%v, newOrgId:%v, err:%v", originOrgId, newOrgId, err)
		return err
	}

	transferUserInfo, err := getTransferUserInfo(originOrgId, newOrgId, transferMatchInfo.NewToOriginUserIdMap)
	if err != nil {
		log.Errorf("[TransferOrgToOtherPlatform] getTransferUserInfo originOrgId:%v, newOrgId:%v, err:%v", originOrgId, newOrgId, err)
		return err
	}

	transferDepartmentInfo, err := getTransferDepartmentInfo(originOrgId, newOrgId, transferMatchInfo.NewToOriginUserIdMap)
	if err != nil {
		log.Errorf("[TransferOrgToOtherPlatform] getTransferOrgInfo originOrgId:%v, newOrgId:%v, err:%v", originOrgId, newOrgId, err)
		return err
	}

	dbErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err := saveTransferOrgInfo(originOrgId, newOrgId, transferOrgInfo, tx)
		if err != nil {
			return err
		}

		err = saveTransferUserInfo(originOrgId, newOrgId, transferUserInfo, transferOrgInfo, transferMatchInfo, tx)
		if err != nil {
			return err
		}

		return saveTransferDepartmentInfo(originOrgId, newOrgId, transferDepartmentInfo, tx)
	})
	if dbErr != nil {
		log.Errorf("[TransferOrgToOtherPlatform] saveToDb originOrgId:%v, newOrgId:%v, err:%v", originOrgId, newOrgId, dbErr)
		return errs.MysqlOperateError
	}

	return nil
}

// getTransferUserInfo 组织转移平台的时候，获取需要更新和新增的用户相关信息
func getTransferUserInfo(originOrgId, newOrgId int64, newToOriginUserIdMap map[int64]int64) (*bo.TransferUserInfo, errs.SystemErrorInfo) {
	userOutInfos, err := GetOutUserInfoListByOrgId(newOrgId)
	if err != nil {
		return nil, err
	}
	idCodes, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserOutInfo, len(userOutInfos))
	if err != nil {
		return nil, err
	}

	transferInfo := &bo.TransferUserInfo{}
	for i, info := range userOutInfos {
		info.Id = idCodes.Ids[i].Id
		info.OrgId = originOrgId
		if newToOriginUserIdMap[info.UserId] != 0 {
			info.UserId = newToOriginUserIdMap[info.UserId]
		}
		transferInfo.MatchUsersOutInfos = append(transferInfo.MatchUsersOutInfos, info)
	}

	return transferInfo, nil
}

func saveTransferUserInfo(originOrgId, newOrgId int64, transferUserInfo *bo.TransferUserInfo, transferOrgInfo *bo.TransferOrgInfo,
	transferMatchInfo *bo.TransferMatchInfo, tx sqlbuilder.Tx) error {
	// 将用户信息中的source_channel和source_platform更改为新的平台
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUser, db.Cond{consts.TcId: db.In(transferMatchInfo.OriginUserIds)},
		transferOrgInfo.ChannelUpdateInfo)
	if err != nil {
		return err
	}

	// 将两个org的userOutInfo全部删掉
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableUserOutInfo, db.Cond{consts.TcOrgId: db.In([]int64{originOrgId, newOrgId})},
		mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
	if err != nil {
		return err
	}

	// 将新组织的用户重新插入，这是为了批量操作，如果update话需要一条条update
	err = mysql.TransBatchInsert(tx, &po.PpmOrgUserOutInfo{}, transferUserInfo.MatchUsersOutInfos)

	return err
}

func getTransferOrgInfo(originOrgId, newOrgId int64, newToOriginUserIdMap map[int64]int64) (*bo.TransferOrgInfo, errs.SystemErrorInfo) {
	transferInfo := &bo.TransferOrgInfo{}

	outInfo, err := GetOrgOutInfo(newOrgId)
	if err != nil {
		return nil, err
	}
	transferInfo.ChannelUpdateInfo = mysql.Upd{
		consts.TcSourceChannel:  outInfo.SourceChannel,
		consts.TcSourcePlatForm: outInfo.SourcePlatform,
	}

	userOrgInfos, err := GetUserOrgInfos(newOrgId)
	if err != nil {
		return nil, err
	}
	idCodes, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserOrganization, len(userOrgInfos))
	if err != nil {
		return nil, err
	}

	for i, info := range userOrgInfos {
		info.Id = idCodes.Ids[i].Id
		info.OrgId = originOrgId
		if newToOriginUserIdMap[info.UserId] != 0 {
			info.UserId = newToOriginUserIdMap[info.UserId]
		}
		transferInfo.OrgUserInfos = append(transferInfo.OrgUserInfos, info)
	}

	return transferInfo, nil
}

func saveTransferOrgInfo(originOrgId, newOrgId int64, transferInfo *bo.TransferOrgInfo, tx sqlbuilder.Tx) error {
	// 更新原orgInfo的sourceChannel和SourcePlatForm
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableOrganization, db.Cond{consts.TcId: originOrgId},
		transferInfo.ChannelUpdateInfo)
	if err != nil {
		return err
	}

	// 删除新的org
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableOrganization, db.Cond{consts.TcId: newOrgId},
		mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
	if err != nil {
		return err
	}

	// 删除老org的orgOutInfo
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableOrganizationOutInfo, db.Cond{consts.TcOrgId: originOrgId},
		mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
	if err != nil {
		return err
	}

	// orgOutInfo中将新的orgId，替换为老的orgId，这样orgId还和原来一样，但是变成了新平台信息
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableOrganizationOutInfo, db.Cond{consts.TcOrgId: newOrgId},
		mysql.Upd{consts.TcOrgId: originOrgId})

	// 删除两个org绑定的用户
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableUserOrganization, db.Cond{consts.TcOrgId: db.In([]int64{originOrgId, newOrgId})},
		mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
	if err != nil {
		return err
	}

	// 插入org与user的关联关系
	err = mysql.TransBatchInsert(tx, &po.PpmOrgUserOrganization{}, transferInfo.OrgUserInfos)

	return err
}

func getTransferDepartmentInfo(originOrgId, newOrgId int64, newToOriginUserIdMap map[int64]int64) (*bo.TransferDepartmentInfo, errs.SystemErrorInfo) {
	userDepartmentInfos, dbErr := dao.SelectUserDepartment(db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    newOrgId,
	})
	if dbErr != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	idCodes, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(*userDepartmentInfos))
	if err != nil {
		return nil, err
	}

	transferInfo := &bo.TransferDepartmentInfo{}
	for i, userDepartmentInfo := range *userDepartmentInfos {
		userDepartmentInfo.Id = idCodes.Ids[i].Id
		userDepartmentInfo.OrgId = originOrgId
		if newToOriginUserIdMap[userDepartmentInfo.UserId] != 0 {
			userDepartmentInfo.UserId = newToOriginUserIdMap[userDepartmentInfo.UserId]
		}
		transferInfo.MatchUserDepartmentInfos = append(transferInfo.MatchUserDepartmentInfos, userDepartmentInfo)
	}

	return transferInfo, nil
}

func saveTransferDepartmentInfo(originOrgId, newOrgId int64, transferInfo *bo.TransferDepartmentInfo, tx sqlbuilder.Tx) error {
	if len(transferInfo.MatchUserDepartmentInfos) == 0 {
		return nil
	}

	// 删除老org的ppm_org_department
	_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{consts.TcOrgId: originOrgId},
		mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
	if err != nil {
		return err
	}

	// 将新部门的orgId都替换成老的orgId
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{consts.TcOrgId: newOrgId},
		mysql.Upd{consts.TcOrgId: originOrgId})
	if err != nil {
		return err
	}

	// 删除老org的ppm_org_department_out_info
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableDepartmentOutInfo, db.Cond{consts.TcOrgId: originOrgId},
		mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
	if err != nil {
		return err
	}

	// 将新部门的orgId都替换成老的orgId
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableDepartmentOutInfo, db.Cond{consts.TcOrgId: newOrgId},
		mysql.Upd{consts.TcOrgId: originOrgId})
	if err != nil {
		return err
	}

	// 删除两个org的ppm_org_user_department
	_, err = mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{consts.TcOrgId: db.In([]int64{originOrgId, newOrgId})},
		mysql.Upd{consts.TcIsDelete: consts.AppIsDeleted})
	if err != nil {
		return err
	}

	// 插入新的用户与部门的关联关系
	err = mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, transferInfo.MatchUserDepartmentInfos)

	return err
}
