package domain

import (
	"fmt"
	"strings"
	"time"

	"github.com/star-table/polaris-backend/service/platform/orgsvc/dao"

	"github.com/star-table/polaris-backend/common/model/vo/orgvo"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

//用来返回用户组织列表
func GetUserOrganizationIdList(userId int64) (*[]bo.PpmOrgUserOrganizationBo, errs.SystemErrorInfo) {
	UserOrganizationPo := &[]po.PpmOrgUserOrganization{}
	UserOrganizationBo := &[]bo.PpmOrgUserOrganizationBo{}

	userIds, err := dao.GetGlobalUserRelation().GetUserIdsByUserId(userId)
	if err != nil {
		log.Errorf("[GetUserOrganizationIdList] GetUserIdsByUserId userId:%v, err:%v", userIds, err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err = mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcUserId:      db.In(userIds),
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		//consts.TcStatus:   consts.AppStatusEnable,
	}, UserOrganizationPo)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	_ = copyer.Copy(UserOrganizationPo, UserOrganizationBo)

	return UserOrganizationBo, nil
}

// 根据uid列表和orgId，从而查到这个orgId属于哪个userId
func getUserOrgInfoByUserIdsAndOrgId(userIds []int64, orgId int64) (*po.PpmOrgUserOrganization, error) {
	userOrgInfo := &po.PpmOrgUserOrganization{}
	err := mysql.SelectOneByCond(consts.TableUserOrganization, db.Cond{
		consts.TcUserId:   db.In(userIds),
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, userOrgInfo)
	if err != nil {
		return nil, err
	}

	return userOrgInfo, nil
}

// 判断这些userId是否有重复的orgId
func checkUserIdsOrgIsRepeated(userIds []int64) (bool, errs.SystemErrorInfo) {
	userOrgInfos := make([]*po.PpmOrgUserOrganization, 0, len(userIds))
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Errorf("[checkUserIdsOrgIsRepeated] GetConnect err: %v", err)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	err = conn.Collection(consts.TableUserOrganization).Find(db.Cond{
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}).Select(consts.TcOrgId).All(&userOrgInfos)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return false, nil
		}
		log.Errorf("[checkUserIdsOrgIsRepeated] GetConnect err: %v", err)
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	idCountMap := make(map[int64]int, len(userOrgInfos))
	for _, info := range userOrgInfos {
		idCountMap[info.OrgId] += 1
		if idCountMap[info.OrgId] > 1 {
			return true, nil
		}
	}

	return false, nil
}

//用来获取用户最新的组织关系
func GetUserOrganizationNewestRelation(orgId, userId int64) (*bo.PpmOrgUserOrganizationBo, errs.SystemErrorInfo) {
	UserOrganizationPo := &po.PpmOrgUserOrganization{}
	UserOrganizationBo := &bo.PpmOrgUserOrganizationBo{}

	conn, err := mysql.GetConnect()
	if err != nil {
		return nil, errs.MysqlOperateError
	}
	err = conn.Collection(consts.TableUserOrganization).Find(db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcUserId: userId,
	}).OrderBy("id desc").Limit(1).One(UserOrganizationPo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.BuildSystemErrorInfo(errs.UserOrgNotRelation)
		} else {
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}
	_ = copyer.Copy(UserOrganizationPo, UserOrganizationBo)
	return UserOrganizationBo, nil
}

func handleGetOrganizationUserListCond(input *vo.OrgUserListReq, cond db.Cond) {
	if len(input.CheckStatus) > 0 {
		cond["o."+consts.TcCheckStatus] = db.In(input.CheckStatus)
	}
	if input.Status != nil && *input.Status != 0 {
		cond["o."+consts.TcStatus] = *input.Status
	}
	if input.UseStatus != nil && *input.UseStatus != 0 {
		cond["o."+consts.TcUseStatus] = *input.UseStatus
	}
	if input.Name != nil && *input.Name != "" {
		cond["u."+consts.TcName] = db.Like("%" + *input.Name + "%")
	}
	if input.Email != nil && *input.Email != "" {
		cond["u."+consts.TcEmail] = db.Eq(*input.Email)
	}
	if input.Mobile != nil && *input.Mobile != "" {
		cond["u."+consts.TcMobile] = db.Eq(*input.Mobile)
	}
}

func GetOrganizationUserList(orgId int64, page, size int, input *vo.OrgUserListReq, allUserHaveRoleIds []int64) (uint64, []bo.PpmOrgUserOrganizationBo, errs.SystemErrorInfo) {

	conn, err := mysql.GetConnect()
	if err != nil {
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	cond := db.Cond{
		"o." + consts.TcOrgId:    orgId,
		"o." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"o." + consts.TcUserId:   db.Raw("u." + consts.TcId),
	}
	if input != nil {
		handleGetOrganizationUserListCond(input, cond)
	}
	total := &po.PpmOrgUserOrganizationCount{}
	totalErr := conn.Select(db.Raw("count(*) as total")).From("ppm_org_user_organization o", "ppm_org_user u").Where(cond).One(total)
	if totalErr != nil {
		log.Error(totalErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, totalErr)
	}
	count := total.Total

	//count, err := mysql.SelectCountByCond(consts.TableUserOrganization, cond)
	//if err != nil {
	//	return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	orgUserPo := &[]po.PpmOrgUserOrganization{}
	//默认是审核时间升序，创建时间降序
	order := db.Raw("o.check_status asc, o.create_time desc")
	if len(allUserHaveRoleIds) > 0 {
		idStr := strings.Replace(strings.Trim(fmt.Sprint(allUserHaveRoleIds), "[]"), " ", ",", -1)
		order = db.Raw("FIELD(o.user_id," + idStr + ") desc, o.check_status asc, o.create_time desc")
	}

	//err = mysql.SelectAllByCondWithNumAndOrder(consts.TableUserOrganization, cond, nil, page, size, order, orgUserPo)
	//if err != nil {
	//	return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}

	selectErr := conn.Select(db.Raw("o.*")).From("ppm_org_user_organization o", "ppm_org_user u").Where(cond).Offset((page - 1) * size).Limit(size).
		OrderBy(order).All(orgUserPo)
	if selectErr != nil {
		log.Error(selectErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
	}

	orgUserBo := &[]bo.PpmOrgUserOrganizationBo{}
	copyErr := copyer.Copy(orgUserPo, orgUserBo)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return count, *orgUserBo, nil
}

// GetOrgIdListBySourceChannel 获取组织 id 列表
// sourceChannel 为 "-1" 时，表示忽略此筛选条件
// isPaid 为 -1 时，表示忽略此筛选条件
func GetOrgIdListBySourceChannel(sourceChannels []string, page int, size int, isPaid int) ([]int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	orgOutInfos := &[]po.PpmOrgOrganizationOutInfo{}
	cond := db.Cond{
		"o." + consts.TcIsDelete: consts.AppIsNoDelete,
		"o." + consts.TcStatus:   consts.AppStatusEnable,
		"c." + consts.TcIsDelete: db.In([]int{consts.AppIsDeleteUndefined, consts.AppIsNoDelete}), // 本地组织is_delete是0
		"o." + consts.TcOrgId:    db.Raw("c." + consts.TcOrgId),
	}
	// sourceChannels 为空时，表示忽略此筛选条件
	if len(sourceChannels) > 0 {
		cond["o."+consts.TcSourceChannel] = db.In(sourceChannels)
	}
	if ok, _ := slice.Contain([]int{1, 2}, isPaid); ok {
		if isPaid == 1 {
			cond["c."+consts.TcPayLevel] = 1
		} else {
			cond["c."+consts.TcPayLevel] = db.NotEq(1)
			cond["c."+consts.TcPayEndTime] = db.Gte(date.Format(time.Now()))
		}
	}
	mid := conn.Select(db.Raw("o.org_id")).From("ppm_org_organization_out_info o", "ppm_orc_config c").Where(cond)
	if page > 0 && size > 0 {
		mid = mid.Offset((page - 1) * size).Limit(size)
	}

	selectErr := mid.All(orgOutInfos)
	if selectErr != nil {
		log.Error(selectErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
	}
	orgIds := make([]int64, 0)

	for _, outInfo := range *orgOutInfos {
		orgIds = append(orgIds, outInfo.OrgId)
	}

	return orgIds, nil
}

func GetOrgIdListByPage(input *orgvo.GetOrgIdListByPageReqVoData, page, size int) ([]int64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	cond := db.Cond{}
	mid := conn.Select(db.Raw("id")).From(consts.TableOrganization).Where(cond).OrderBy("id desc")
	if page > 0 && size > 0 {
		mid = mid.Offset((page - 1) * size).Limit(size)
	} else {
		mid = mid.Offset(0).Limit(20)
	}
	orgs := make([]po.PpmOrgOrganization, 0)
	selectErr := mid.All(&orgs)
	if selectErr != nil {
		log.Error(selectErr)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
	}

	orgIds := make([]int64, 0)
	for _, tmpOrg := range orgs {
		orgIds = append(orgIds, tmpOrg.Id)
	}

	return orgIds, nil
}

func GetOutOrgInfoByOrgIdBatch(orgIds []int64) ([]*po.PpmOrgOrganizationOutInfo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	cond := db.Cond{
		consts.TcOrgId:    db.In(orgIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}

	outInfo := &po.PpmOrgOrganizationOutInfo{}
	var result []*po.PpmOrgOrganizationOutInfo
	err = conn.SelectFrom(outInfo.TableName()).Where(cond).All(&result)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return result, nil
}

func ResetUserCheckStatusWait(orgId, userId int64) errs.SystemErrorInfo {
	userOrgPoList := make([]*po.PpmOrgUserOrganization, 0)
	err := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, &userOrgPoList)
	if err != nil {
		log.Errorf("[GetUserOrgListByUserId] err:%v, orgId:%v, userId:%v", err, orgId, userId)
		return errs.MysqlOperateError
	}
	if len(userOrgPoList) == 0 {
		return nil
	}
	if userOrgPoList[0].CheckStatus == consts.AppCheckStatusFail {
		// 重置成 待审核
		_, err = mysql.UpdateSmartWithCond(consts.TableUserOrganization, db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcUserId:      userId,
			consts.TcCheckStatus: consts.AppCheckStatusFail,
			consts.TcIsDelete:    consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcCheckStatus: consts.AppCheckStatusWait,
		})
		if err != nil {
			log.Error(err)
			return errs.MysqlOperateError
		}
	}
	return nil
}
