package domain

import (
	"math"
	"strings"
	"time"

	"github.com/spf13/cast"

	msgPb "github.com/star-table/interface/golang/msg/v1"
	"github.com/star-table/common/core/threadlocal"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/facade/common/report"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/ordervo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetOrgBoList() ([]bo.OrganizationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrgOrganization{}
	err := mysql.SelectAllByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
	}, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.OrganizationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return *bos, nil
}

func GetOrgBoListByIds(orgIds []int64) (*[]bo.OrganizationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrgOrganization{}
	err := mysql.SelectAllByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcId:       db.In(orgIds),
	}, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.OrganizationBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return bos, nil
}

func GetOrgBoById(orgId int64, tx ...sqlbuilder.Tx) (*bo.OrganizationBo, errs.SystemErrorInfo) {
	po := &po.PpmOrgOrganization{}
	var err error
	if len(tx) > 0 {
		err = mysql.TransSelectOneByCond(tx[0], consts.TableOrganization, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
			consts.TcId:       orgId,
		}, po)
	} else {
		err = mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcStatus:   consts.AppStatusEnable,
			consts.TcId:       orgId,
		}, po)
	}

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bo := &bo.OrganizationBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return bo, nil
}

func GetOrgBoByCode(code string) (*bo.OrganizationBo, errs.SystemErrorInfo) {
	po := &po.PpmOrgOrganization{}
	err := mysql.SelectOneByCond(consts.TableOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcCode:     code,
	}, po)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bo := &bo.OrganizationBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func ScheduleOrganizationPageList(size int, page int) (*[]*bo.ScheduleOrganizationListBo, int64, errs.SystemErrorInfo) {

	pos, count, err := dao.OrcConfigPageList(size, page)

	if err != nil {
		return nil, int64(count), errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	bos := &[]*bo.ScheduleOrganizationListBo{}

	err = copyer.Copy(pos, bos)

	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return bos, count, nil

}

//校验当前用户是否有效
func VerifyOrg(orgId int64, userId int64) bool {
	baseUserInfo, err := GetBaseUserInfo(orgId, userId)
	if err != nil {
		log.Error(err)
		return false
	}
	//未被移除且通过审核，认为当前用户有效
	return baseUserInfo.OrgUserIsDelete == consts.AppIsNoDelete && baseUserInfo.OrgUserCheckStatus == consts.AppCheckStatusSuccess
}

func VerifyOrgUsers(orgId int64, userIds []int64) bool {
	userIds = slice.SliceUniqueInt64(userIds)
	baseUserInfos, err := GetBaseUserInfoBatch(orgId, userIds)
	if err != nil {
		log.Error(err)
		return false
	}
	if len(baseUserInfos) != len(userIds) {
		log.Error("部分用户无效")
		return false
	}
	for _, userInfo := range baseUserInfos {
		//如果存在待审核或者已移除的，则不允许更新
		if userInfo.OrgUserIsDelete != consts.AppIsNoDelete || userInfo.OrgUserCheckStatus != consts.AppCheckStatusSuccess {
			return false
		}
	}
	return true
}

func VerifyOrgUsersReturnValid(orgId int64, userIds []int64) []int64 {
	userIds = slice.SliceUniqueInt64(userIds)
	baseUserInfos, err := GetBaseUserInfoBatch(orgId, userIds)
	if err != nil {
		log.Error(err)
		return nil
	}

	var validUserIds []int64
	for _, userInfo := range baseUserInfos {
		// 如果存在待审核或者已移除的，则不允许更新
		if userInfo.OrgUserIsDelete != consts.AppIsNoDelete ||
			userInfo.OrgUserCheckStatus != consts.AppCheckStatusSuccess {
		} else {
			validUserIds = append(validUserIds, userInfo.UserId)
		}
	}
	return validUserIds
}

func GetOrgByOutOrgId(outOrgId string) (*bo.OrganizationBo, errs.SystemErrorInfo) {
	outInfo := &po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOutOrgId: outOrgId,
	}, outInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.OrgOutInfoNotExist
		}
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	orgId := outInfo.OrgId
	orgPo := &po.PpmOrgOrganization{}
	err = mysql.SelectById(consts.TableOrganization, orgId, orgPo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	orgBo := &bo.OrganizationBo{}
	_ = copyer.Copy(orgPo, orgBo)
	return orgBo, nil
}

func CreateOrg(createOrgBo bo.CreateOrgBo, creatorId int64, sourceChannel, sourcePlatform string, outOrgId string) (int64, errs.SystemErrorInfo) {
	orgId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrganization)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}

	orgName := strings.TrimSpace(createOrgBo.OrgName)
	isOrgNameRight := format.VerifyOrgNameFormat(orgName)
	if !isOrgNameRight {
		return 0, errs.OrgNameLenError
	}

	//组织
	org := &po.PpmOrgOrganization{}
	org.Id = orgId
	org.Status = consts.AppStatusEnable
	org.IsDelete = consts.AppIsNoDelete
	org.Creator = creatorId
	org.Owner = creatorId
	org.Updator = creatorId
	org.Name = orgName
	org.SourceChannel = sourceChannel
	org.SourcePlatform = sourcePlatform
	if createOrgBo.IndustryID != nil {
		org.IndustryId = *createOrgBo.IndustryID
	}
	if createOrgBo.Scale != nil {
		org.Scale = *createOrgBo.Scale
	}
	err1 := mysql.Insert(org)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	orgOutId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrganizationOutInfo)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}
	//外部组织信息
	orgOutInfo := &po.PpmOrgOrganizationOutInfo{
		Id:             orgOutId,
		OrgId:          orgId,
		OutOrgId:       outOrgId,
		SourceChannel:  sourceChannel,
		SourcePlatform: sourcePlatform,
		Name:           createOrgBo.OrgName,
		Creator:        creatorId,
		Updator:        creatorId,
	}

	err1 = mysql.Insert(orgOutInfo)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	//组织配置
	orgConfigId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgConfig)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}
	//组织配置信息
	//payLevel := &po.PpmBasPayLevel{}
	//payLevelErr := mysql.SelectById(payLevel.TableName(), 1, payLevel)
	//if payLevelErr != nil {
	//	log.Errorf("[CreateOrg] payLevel SelectById failed:%v", payLevelErr)
	//	return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}
	orgConfig := &po.PpmOrcConfig{
		Id:           orgConfigId,
		OrgId:        orgId,
		PayStartTime: time.Now(),
		PayEndTime:   time.Now().Add(time.Duration(0) * time.Second),
	}
	if CheckIsPrivateDeploy() {
		orgConfig.PayLevel = consts.PayLevelPrivateDeploy
		orgConfig.PayEndTime = time.Now().AddDate(100, 0, 0)
	}
	err1 = mysql.Insert(orgConfig)
	if err1 != nil {
		log.Error(err1)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	// 事件上报
	asyn.Execute(func() {
		e := &commonvo.OrgEvent{}
		e.OrgId = org.Id
		e.New = org

		openTraceId, _ := threadlocal.Mgr.GetValue(consts.JaegerContextTraceKey)
		openTraceIdStr := cast.ToString(openTraceId)

		report.ReportOrgEvent(msgPb.EventType_OrgInited, openTraceIdStr, e)
	})

	return orgId, nil
}

func GetOrgInfoByOutOrgId(outOrgId string, sourceChannel string) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	outOrgInfo := &po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcOutOrgId:      outOrgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, outOrgInfo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.OrgOutInfoNotExist)
	}
	orgInfo, err1 := GetOrgBoById(outOrgInfo.OrgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.OrgNotExist)
	}
	return &bo.BaseOrgInfoBo{
		OrgId:         orgInfo.Id,
		OrgName:       orgInfo.Name,
		OutOrgId:      outOrgId,
		SourceChannel: sourceChannel,
		OrgOwnerId:    orgInfo.Owner,
	}, nil
}

func UpdateOrg(updateBo bo.UpdateOrganizationBo) errs.SystemErrorInfo {

	organizationBo := updateBo.Bo
	upds := updateBo.OrganizationUpdateCond

	_, err := mysql.UpdateSmartWithCond(consts.TableOrganization, db.Cond{
		consts.TcId: organizationBo.Id,
	}, upds)

	if err != nil {
		log.Errorf("mysql.TransUpdateSmart: %q\n", err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	err = ClearCacheBaseOrgInfo(organizationBo.Id)

	if err != nil {
		log.Errorf("redis err: %q\n", err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}

// UpdateOrgWithTx 在事务内更新组织
func UpdateOrgWithTx(updateBo bo.UpdateOrganizationBo, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	var oriErr error
	var err errs.SystemErrorInfo
	organizationBo := updateBo.Bo
	upd := updateBo.OrganizationUpdateCond

	_, oriErr = mysql.TransUpdateSmartWithCond(tx, consts.TableOrganization, db.Cond{
		consts.TcId: organizationBo.Id,
	}, upd)
	if oriErr != nil {
		log.Errorf("[UpdateOrgWithTx] mysql.TransUpdateSmartWithCond err: %v\n", oriErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	err = ClearCacheBaseOrgInfo(organizationBo.Id)
	if err != nil {
		log.Errorf("[UpdateOrgWithTx] ClearCacheBaseOrgInfo err: %v\n", err)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return nil
}

func JudgeUserIsAdmin(outOrgId string, outUserId string, sourceChannel string) bool {
	return false
}

func GetSimpleOrgOutInfo(orgId int64, sourceChannel string) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	outOrgInfo := &po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcIsDelete:      consts.AppIsNoDelete,
		consts.TcSourceChannel: sourceChannel,
	}, outOrgInfo)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.OrgOutInfoNotExist)
	}
	orgInfo, err1 := GetOrgBoById(outOrgInfo.OrgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.OrgNotExist)
	}
	return &bo.BaseOrgInfoBo{
		OrgId:      orgInfo.Id,
		OrgName:    orgInfo.Name,
		OutOrgId:   outOrgInfo.OutOrgId,
		OrgOwnerId: orgInfo.Owner,
	}, nil
}

func GetOrgOutInfo(orgId int64) (*bo.OrgOutInfoBo, errs.SystemErrorInfo) {
	outOrgInfo := &po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOutOrgId: db.NotEq(""),
	}, outOrgInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.BuildSystemErrorInfo(errs.OrgOutInfoNotExist, err)
		}
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bo := &bo.OrgOutInfoBo{}
	_ = copyer.Copy(outOrgInfo, bo)
	return bo, nil
}

func GetOrgOutInfoByOutOrgId(orgId int64, outOrgId string) (*po.PpmOrgOrganizationOutInfo, errs.SystemErrorInfo) {
	outOrgInfo := &po.PpmOrgOrganizationOutInfo{}
	var cond db.Cond
	if orgId > 0 {
		cond = db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}
	} else {
		cond = db.Cond{
			consts.TcOutOrgId: outOrgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}
	}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, cond, outOrgInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.BuildSystemErrorInfo(errs.OrgOutInfoNotExist, err)
		}
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return outOrgInfo, nil
}

func GetOrgOutInfoWithoutLocal(orgId int64) (*bo.OrgOutInfoBo, errs.SystemErrorInfo) {
	outOrgInfo := &po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, outOrgInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.OrgOutInfoNotExist
		}
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	bo := &bo.OrgOutInfoBo{}
	_ = copyer.Copy(outOrgInfo, bo)
	return bo, nil
}

func GetOrgOutInfoByTenantKey(tenantKey string) (*bo.OrgOutInfoBo, errs.SystemErrorInfo) {
	outOrgInfo := &po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcOutOrgId: tenantKey,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, outOrgInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.OrgOutInfoNotExist
		}
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	bo := &bo.OrgOutInfoBo{}
	_ = copyer.Copy(outOrgInfo, bo)
	return bo, nil
}

func UpdateOrgOutInfoTenantCode(orgId int64, sourceChannel, tenantCode string) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcOrgId:         orgId,
		consts.TcSourceChannel: sourceChannel,
		consts.TcIsDelete:      consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcTenantCode: tenantCode,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

// UpdateOrgOutInfoPermanentCode 企微重新安装后需要更新code
func UpdateOrgOutInfoPermanentCode(orgId int64, tenantCode string) errs.SystemErrorInfo {
	_, err := mysql.UpdateSmartWithCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcAuthTicket: tenantCode,
	})
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func GetOrgSourceChannel(orgId int64, sourceChannel string) string {
	info := po.PpmOrgOrganizationOutInfo{}
	err := mysql.SelectOneByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete: db.In([]int{0, 2}),
		consts.TcOrgId:    orgId,
	}, &info)
	if err != nil {
		return sourceChannel
	}

	return info.SourceChannel
}

// GetOutOrgListByOutOrgIdsAndSource 通过 outOrgIds 查询外部组织信息
func GetOutOrgListByOutOrgIdsAndSource(outOrgIds []string) ([]po.PpmOrgOrganizationOutInfo, errs.SystemErrorInfo) {
	list := make([]po.PpmOrgOrganizationOutInfo, 0)
	if len(outOrgIds) < 1 {
		return list, nil
	}
	dbErr := mysql.SelectAllByCond(consts.TableOrganizationOutInfo, db.Cond{
		consts.TcIsDelete: db.In([]int{0, 2}),
		consts.TcOutOrgId: db.In(outOrgIds),
	}, &list)
	if dbErr != nil {
		if dbErr == db.ErrNoMoreRows {
			return list, nil
		}
		log.Errorf("[GetOutOrgListByOutOrgIds] err: %v", dbErr)
		return list, errs.BuildSystemErrorInfo(errs.MysqlOperateError, dbErr)
	}

	return list, nil
}

func GetFunctionKeyListByFunctions(functions []ordervo.FunctionLimitObj) []string {
	funcKeys := make([]string, 0, len(functions))
	for _, item := range functions {
		funcKeys = append(funcKeys, item.Key)
	}

	return funcKeys
}

func GetOrgMemberCount(orgId int64) (uint64, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	cond := db.Cond{
		"o." + consts.TcOrgId:    orgId,
		"o." + consts.TcIsDelete: consts.AppIsNoDelete,
		"u." + consts.TcIsDelete: consts.AppIsNoDelete,
		"o." + consts.TcUserId:   db.Raw("u." + consts.TcId),
	}
	total := &po.PpmOrgUserOrganizationCount{}
	totalErr := conn.Select(db.Raw("count(*) as total")).From("ppm_org_user_organization o", "ppm_org_user u").Where(cond).One(total)
	if totalErr != nil {
		log.Error(totalErr)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, totalErr)
	}
	count := total.Total

	return count, nil
}

// GetRemainDays 根据用户使用付费产品的截止时间，计算剩余天数。不足一天的按一天计算。
func GetRemainDays(payEndTime time.Time) uint {
	diffDays := payEndTime.Sub(time.Now())
	diff := diffDays.Hours() / float64(24)
	diff = math.Ceil(diff)
	if diff <= 0 {
		diff = 0
	}

	return uint(diff)
}

// GetAppDeployType 获取应用部署类型。public：saas 公共模式；private: 私有化部署
func GetAppDeployType(runMode int) string {
	defaultVal := "unknown"
	switch runMode {
	case 1, 2:
		defaultVal = "public"
	case 3, 4:
		defaultVal = "private"
	}

	return defaultVal
}

// CheckIsPrivateDeploy 检查是否是私有化部署
func CheckIsPrivateDeploy() bool {
	appDeployType := GetAppDeployType(config.GetConfig().Application.RunMode)
	return appDeployType == "private"
}

func GetNewbieGuideTemplateId(sourceChannel string) int64 {
	return 0
}
