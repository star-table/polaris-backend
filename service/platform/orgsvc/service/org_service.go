package service

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/core/util/pinyin"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/uservo"
	"github.com/star-table/polaris-backend/facade/appfacade"
	"github.com/star-table/polaris-backend/facade/commonfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	sconsts "github.com/star-table/polaris-backend/service/platform/orgsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	roledomain "github.com/star-table/polaris-backend/service/platform/orgsvc/domain/roledomain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetOrgBoList() ([]bo.OrganizationBo, errs.SystemErrorInfo) {
	return domain.GetOrgBoList()
}

func GetBaseOrgInfoByOutOrgId(outOrgId string) (*bo.BaseOrgInfoBo, errs.SystemErrorInfo) {
	return domain.GetBaseOrgInfoByOutOrgId(outOrgId)
}

// GetOrgOutInfoByOutOrgIdBatch 批量获取外部组织信息
func GetOrgOutInfoByOutOrgIdBatch(input orgvo.GetOrgOutInfoByOutOrgIdBatchReqInput) ([]bo.BaseOrgOutInfoBo, errs.SystemErrorInfo) {
	resList := make([]bo.BaseOrgOutInfoBo, 0, len(input.OutOrgIds))
	outOrgArr, err := domain.GetOutOrgListByOutOrgIdsAndSource(input.OutOrgIds)
	if err != nil {
		log.Errorf("[GetBaseOrgInfoByOutOrgIdBatch] err: %v", err)
		return resList, err
	}
	copyer.Copy(outOrgArr, &resList)

	return resList, nil
}

func CreateOrgBase(req orgvo.CreateOrgReqVo, sourceChannel, sourcePlatform string) (int64, errs.SystemErrorInfo) {
	if req.Data.CreateOrgReq.CodeToken != nil && *req.Data.CreateOrgReq.CodeToken != "" {
		//如果传入了codeToken，则需要绑定飞书已有团队
		codeInfoJson, redisErr := cache.Get(sconsts.CacheFsAuthCodeToken + *req.Data.CreateOrgReq.CodeToken)
		if redisErr != nil {
			log.Error(redisErr)
			return 0, errs.RedisOperateError
		}
		if codeInfoJson == "" {
			return 0, errs.CodeTokenInvalid
		}
		codeInfo := vo.FeiShuAuthCodeResp{}
		_ = json.FromJson(codeInfoJson, &codeInfo)
		tenantKey := codeInfo.TenantKey

		cacheKey := fmt.Sprintf("%s%s", sconsts.CacheFsOrgInit, tenantKey)
		defer cache.Del(cacheKey)
		cache.SetEx(cacheKey, tenantKey, 60)

		//加锁
		lockKey := consts.NewFeiShuCorpInitKey + tenantKey
		uuid := uuid.NewUuid()
		suc, err := cache.TryGetDistributedLock(lockKey, uuid)
		log.Infof("准备获取分布式锁 %v", suc)
		if err != nil {
			log.Error(err)
			return 0, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		if suc {
			log.Infof("获取分布式锁成功 %v", suc)
			defer func() {
				if _, lockErr := cache.ReleaseDistributedLock(lockKey, uuid); lockErr != nil {
					log.Error(lockErr)
				}
			}()
			//如果已经有可用团队，直接正常进去(防止一个人停在注册页面，重复创建)
			orgOutInfo, _ := domain.GetOrgOutInfoByTenantKey(tenantKey)
			if orgOutInfo != nil {
				return orgOutInfo.OrgId, nil
			}

			orgId, createErr := CreateOrg(req, sourceChannel, sourcePlatform)
			if createErr != nil {
				log.Error(createErr)
				return 0, createErr
			}
			return orgId, nil
		}
	} else {
		return CreateOrg(req, sourceChannel, sourcePlatform)
	}

	return 0, nil
}

func CreateOrg(req orgvo.CreateOrgReqVo, sourceChannel, sourcePlatform string) (int64, errs.SystemErrorInfo) {
	creatorId := req.Data.CreatorId
	createReqInfo := req.Data.CreateOrgReq
	if req.Data.CreateOrgReq.CreatorName != nil {
		if !format.VerifyUserNameFormat(strings.Trim(*(req.Data.CreateOrgReq.CreatorName), " ")) {
			return 0, errs.BuildSystemErrorInfo(errs.UserNameLenError)
		}
	}
	userInfoBo, _, err := domain.GetUserBo(creatorId)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}
	orgId, err := domain.CreateOrg(bo.CreateOrgBo{
		OrgName:    createReqInfo.OrgName,
		IndustryID: createReqInfo.IndustryID,
		Scale:      createReqInfo.Scale,
	}, creatorId, sourceChannel, sourcePlatform, "")
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return 0, err
	}

	//初始化组织相关资源
	err = CreateOrgRelationResource(orgId, creatorId, sourceChannel, sourcePlatform, req.Data.CreateOrgReq.OrgName)
	if err != nil {
		log.Error(err)
		return orgId, err
	}

	//用户和组织关联
	err = domain.AddUserOrgRelation(orgId, creatorId, true, false, false)
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return orgId, err
	}
	upd := mysql.Upd{}
	if userInfoBo.OrgID == 0 {
		//更新用户的orgId
		upd[consts.TcOrgId] = orgId

		err = domain.UpdateGlobalUserLastLoginInfo(userInfoBo.ID, orgId)
		if err != nil {
			log.Errorf("[CreateOrg] UpdateGlobalUserLastLoginInfo userId:%v, orgId:%v, err:%v", userInfoBo.ID, orgId, err)
			return orgId, err
		}
	}
	if createReqInfo.CreatorName != nil {
		creatorName := strings.Trim(*(req.Data.CreateOrgReq.CreatorName), " ")
		//更新用户的名称
		upd[consts.TcName] = creatorName
		upd[consts.TcNamePinyin] = pinyin.ConvertToPinyin(creatorName)
	}
	if len(upd) > 0 {
		err = domain.UpdateUserInfo(creatorId, upd)
		if err != nil {
			log.Info(strs.ObjectToString(err))
			return orgId, err
		}
	}

	// 同步无码组织配置和组织字段
	_, saveErr := SaveOrgRemarkAndOrgFields(orgId, creatorId, &orgvo.OrgRemarkConfigType{
		OrgSummaryTableAppId: 0,
		TagAppId:             0,
		PriorityAppId:        0,
		IssueBarAppId:        0,
		SceneAppId:           0,
		IssueStatusAppId:     0,
		EmptyProjectAppId:    0,
	})
	if saveErr != nil {
		log.Error(saveErr)
		return 0, saveErr
	}

	// 创建组织时，创建项目视图及其视图镜像
	// 极星标品：创建目录应用、创建任务视图
	//if viewResp := projectfacade.SyncOrgDefaultViewMirror(projectvo.SyncOrgDefaultViewMirrorReq{
	//	OrgIds:                              []int64{orgId},
	//	StartPage:                           1,
	//	NeedUpdateSummaryAppVisibilityToAll: true,
	//}); viewResp.Failure() {
	//	log.Error(viewResp.Error())
	//	return orgId, viewResp.Error()
	//}
	viewResp := projectfacade.CreateOrgDirectoryAppsAndViews(projectvo.CreateOrgDirectoryAppsReq{OrgId: orgId})
	if viewResp.Failure() {
		log.Errorf("[CreateOrg] projectfacade.CreateOrgDirectoryAppsAndViews err:%v, orgId:%v", viewResp.Error(), orgId)
		return orgId, viewResp.Error()
	}

	//刷新用户缓存
	userToken := req.Data.UserToken
	// 创建组织后，必须 updOutInfo，否则后续会话会使用旧组织的 outOrg info 导致一系列问题
	err = UpdateCacheUserInfoOrgId(userToken, orgId, req.UserId, true)
	if err != nil {
		log.Info(strs.ObjectToString(err))
	}

	err = domain.ClearBaseUserInfo(orgId, creatorId)
	if err != nil {
		log.Error(err)
	}

	return orgId, nil
}

//// 初始化一个已存在的组织。
//// 目前看只需初始化优先级、ocr_config 等
//func InitExistOrg(orgId, userId int64, input vo.InitExistOrgReq) (*orgvo.InitExistOrgRespData, errs.SystemErrorInfo) {
//	// 需要向 ppm_orc_config 中增加配置
//	if input.NeedOcrConfig == 1 {
//		orgConfigId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableOrgConfig)
//		if err != nil {
//			log.Info(strs.ObjectToString(err))
//			return nil, err
//		}
//		//组织配置信息
//		orgConfig := &po.PpmOrcConfig{
//			Id:    orgConfigId,
//			OrgId: orgId,
//		}
//		err1 := mysql.Insert(orgConfig)
//		if err1 != nil {
//			log.Error(err1)
//			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err1)
//		}
//	}
//	//优先级初始化
//	if input.NeedPriority == 1 {
//		priorityInfo := projectfacade.ProjectInit(projectvo.ProjectInitReqVo{OrgId: orgId})
//		if priorityInfo.Failure() {
//			return nil, priorityInfo.Error()
//		}
//	}
//	// 将组织设置为付费用户
//	if err := lc_org_domain.SetOrgPaid(orgId, *input.NeedSetToPaid); err != nil {
//		log.Error(err)
//		return nil, err
//	}
//	// 为组织新建汇总表
//	if input.NeedSummaryTable == 1 {
//		orgBo, err := domain.GetOrgBoById(orgId)
//		orgRemarkJson := orgBo.Remark
//		orgRemarkObj := &orgvo.OrgRemarkConfigType{}
//		if len(orgRemarkJson) > 0 {
//			oriErr := json.FromJson(orgRemarkJson, orgRemarkObj)
//			if oriErr != nil {
//				return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
//			}
//		}
//
//		// 将汇总表的 appId 存入组织属性（remark）中
//		_, _, err = SaveOrgSomeTableAppId(orgId, userId, orgRemarkObj)
//		if err != nil {
//			return nil, err
//		}
//	}
//
//	return &orgvo.InitExistOrgRespData{
//		IsOk:    true,
//		StrInfo: QueryForInitExistOrg(orgId, userId, &input),
//	}, nil
//}

// 利用初始化接口查询一些信息
func QueryForInitExistOrg(orgId, userId int64, input *vo.InitExistOrgReq) string {
	info := make(map[string]interface{})
	if *input.QueryAuthTokenInfo == 1 {
		info["QueryAuthTokenInfo"] = map[string]interface{}{
			"orgId":  strconv.FormatInt(orgId, 10),
			"userId": strconv.FormatInt(userId, 10),
		}
		info["OrgConfig"] = ""
	}
	return json.ToJsonIgnoreError(info)
}

func SaveOrgSummaryTableAppId(orgId, userId int64, input orgvo.SaveOrgSummaryTableAppIdReqVoData) (bool, errs.SystemErrorInfo) {
	return domain.SaveOrgSummaryTableAppId(orgId, userId, input.AppId)
}

// SaveOrgSomeTableAppId 保存组织下的一些应用表的 appId。如果这些应用不存在，则创建
// 如果项目 form 的 appId 为 0，则创建保存项目数据的 project form
func SaveOrgSomeTableAppId(orgId, userId int64, orgRemarkObj *orgvo.OrgRemarkConfigType, tx ...sqlbuilder.Tx) (*orgvo.SaveOrgSummaryTableAppIdReqVoData, bool, errs.SystemErrorInfo) {
	input := orgvo.SaveOrgSummaryTableAppIdReqVoData{
		AppId:             orgRemarkObj.OrgSummaryTableAppId,
		TagAppId:          orgRemarkObj.TagAppId,
		PriorityAppId:     orgRemarkObj.PriorityAppId,
		IssueBarAppId:     orgRemarkObj.IssueBarAppId,
		SceneAppId:        orgRemarkObj.SceneAppId,
		IssueStatusAppId:  orgRemarkObj.IssueStatusAppId,
		EmptyProjectAppId: orgRemarkObj.EmptyProjectAppId,
		ProjectFormAppId:  orgRemarkObj.ProjectFormAppId,
	}
	// 用一个标识记录是否变化，发生变化了，才会更新
	hasChanged := false
	var err errs.SystemErrorInfo
	if input.AppId < 1 {
		input.AppId, err = domain.CreateIssueSummaryTable(orgId, userId, "企业汇总表", orgRemarkObj, 0)
		if err != nil {
			log.Error(err)
			return nil, false, err
		}
		hasChanged = true
	}

	if input.EmptyProjectAppId < 1 {
		appType := consts.LcAppTypeForPolaris
		appName := "空项目"
		resp := appfacade.CreateLessCodeApp(&permissionvo.CreateLessCodeAppReq{
			OrgId:     &orgId,
			AppType:   &appType,
			Name:      &appName,
			UserId:    &userId,
			ExtendsId: input.AppId,
			PkgId:     0,
			Config:    "",
			Hidden:    1,
		})
		if resp.Failure() {
			log.Error(resp.Error())
			return nil, false, resp.Error()
		}
		input.EmptyProjectAppId = resp.Data.Id
		hasChanged = true
	}
	if input.ProjectFormAppId < 1 {
		orgBo, err := domain.GetOrgBoById(orgId, tx...)
		if err != nil {
			log.Error(err)
			return nil, false, err
		}
		projectFormAppId, err := CreateOrgProjectForm(*orgBo, orgRemarkObj)
		if err != nil {
			log.Error(err)
			return nil, false, err
		}
		input.ProjectFormAppId = projectFormAppId
		hasChanged = true
	}

	if hasChanged {
		ok, err := domain.SaveOrgSomeTableAppId(orgId, userId, input, tx...)
		if err != nil {
			log.Error(err)
			return &input, ok, err
		}
	}

	return &input, true, err
}

func bindingUserAndRole(isRoot, isAdmin bool, orgId, userId int64, roleInitResp *bo.RoleInitResp) errs.SystemErrorInfo {
	if isRoot {
		err2 := roledomain.AddRoleUserRelation(orgId, userId, roleInitResp.OrgSuperAdminRoleId)
		if err2 != nil {
			log.Error(err2)
			return err2
		}
	}
	if isAdmin {
		err2 := roledomain.AddRoleUserRelation(orgId, userId, roleInitResp.OrgNormalAdminRoleId)
		if err2 != nil {
			log.Error(err2)
			return err2
		}
	}
	return nil
}

func CreateOrgRelationResource(orgId int64, creatorId int64, sourceChannel, sourcePlatform string, orgName string) errs.SystemErrorInfo {
	//初始化角色、优先级
	//权限、角色初始化
	//roleInitResp, err := service.RoleInit(orgId)
	//if err != nil {
	//	log.Errorf("Rollback： %v", err)
	//	return err
	//}
	//log.Info("权限、角色初始化成功")
	//
	////为用户绑定超级管理员
	//err = bindingUserAndRole(true, false, orgId, creatorId, roleInitResp)
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}
	//log.Info("组织超级管理员角色赋予成功")

	// 初始化管理组
	initResp := userfacade.InitDefaultManageGroup(orgId, creatorId)
	if initResp.Failure() {
		log.Error(initResp.Message)
		return initResp.Error()
	}

	// 系统管理组填充人员
	saveManageGroupResp := userfacade.AddUserToSysManageGroup(orgId, creatorId, uservo.AddUserToSysManageGroupReq{UserIds: []int64{creatorId}})
	if saveManageGroupResp.Failure() {
		log.Error(saveManageGroupResp.Message)
		return saveManageGroupResp.Error()
	}

	//优先级初始化
	//priorityInfo := projectfacade.ProjectInit(projectvo.ProjectInitReqVo{OrgId: orgId})
	//if priorityInfo.Failure() {
	//	return errs.BuildSystemErrorInfo(errs.BaseDomainError, priorityInfo.Error())
	//}
	//log.Info("优先级初始化成功")

	////部门初始化
	//departmentId, departmentErr := domain.LarkDepartmentInit(orgId, sourceChannel, sourcePlatform, orgName, creatorId)
	//if departmentErr != nil {
	//	return errs.BuildSystemErrorInfo(errs.BaseDomainError, departmentErr)
	//}
	//log.Info("部门初始化成功")
	//
	////用户和顶级部门绑定
	//err = domain.BoundDepartmentUser(orgId, []int64{creatorId}, departmentId, creatorId, true)
	//if err != nil {
	//	log.Error(err)
	//	return err
	//}
	return nil
}

func UserOrganizationList(userId int64) (*vo.UserOrganizationListResp, errs.SystemErrorInfo) {

	organizationBo, err := domain.GetUserOrganizationIdList(userId)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	orgIds := []int64{}
	enableOrgIds := []int64{}
	enableStatus := consts.AppStatusEnable
	disabledStatus := consts.AppStatusDisabled
	for _, value := range *organizationBo {
		orgIds = append(orgIds, value.OrgId)
		if value.Status == enableStatus {
			enableOrgIds = append(enableOrgIds, value.OrgId)
		}
	}

	bos, err := domain.GetOrgBoListByIds(orgIds)

	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	resultList := &[]*vo.UserOrganization{}
	copyErr := copyer.Copy(bos, resultList)

	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	functionMap := map[int64][]string{}
	for _, id := range orgIds {
		functions, err := domain.GetOrgPayFunction(id)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		funcKeys := domain.GetFunctionKeyListByFunctions(functions)
		functionMap[id] = funcKeys
	}

	orgSysAdminIds, err := domain.GetOrgSysAdmin(orgIds)
	if err != nil {
		return nil, err
	}

	for k, v := range *resultList {
		if ok, _ := slice.Contain(enableOrgIds, v.ID); ok {
			(*resultList)[k].OrgIsEnabled = &enableStatus
		} else {
			(*resultList)[k].OrgIsEnabled = &disabledStatus
		}
		if functions, ok := functionMap[v.ID]; ok {
			(*resultList)[k].Functions = functions
		} else {
			(*resultList)[k].Functions = []string{}
		}
		isAdmin, _ := slice.Contain(orgSysAdminIds[v.ID], userId)
		v.IsAdmin = isAdmin
	}

	return &vo.UserOrganizationListResp{
		List: *resultList,
	}, nil
}

func SwitchUserOrganization(orgId, userId int64, token string) errs.SystemErrorInfo {
	orgUserInfo, orgErr := domain.GetOrgUserInfo(userId, orgId)
	if orgErr != nil {
		if orgErr != db.ErrNoMoreRows {
			log.Errorf("[SwitchUserOrganization] GetOrgUserInfo userId:%v,orgId:%v, err:%v", userId, orgId, orgErr)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, orgErr)
		}
	} else {
		userId = orgUserInfo.UserId
	}

	//监测可用性
	baseUserInfo, err := GetBaseUserInfo(orgId, userId)
	if err != nil {
		log.Error(err)
		return err
	}

	err = baseUserInfoOrgStatusCheck(*baseUserInfo)
	if err != nil {
		log.Error(err)
		return err
	}

	//修改用户默认组织
	orgUserId, updateUserInfoErr := domain.UpdateUserDefaultOrg(userId, orgId)
	if updateUserInfoErr != nil {
		log.Error(updateUserInfoErr)
	}

	//更改用户缓存的orgId
	err = UpdateCacheUserInfoOrgId(token, orgId, orgUserId, true)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

// 获取组织信息
func OrganizationInfo(req orgvo.OrganizationInfoReqVo) (*vo.OrganizationInfoResp, errs.SystemErrorInfo) {

	bo, err := domain.GetOrgBoById(req.OrgId)

	if err != nil {
		return nil, err
	}

	//跨服务查询
	resp := commonfacade.AreaInfo(commonvo.AreaInfoReqVo{
		IndustryID: bo.IndustryId,
		CountryID:  bo.CountryId,
		ProvinceID: bo.ProvinceId,
		CityID:     bo.CityId,
	})

	if resp.Failure() {
		log.Error(resp.Message)
		return nil, resp.Error()
	}
	ownerInfo, err := domain.GetBaseUserInfo(req.OrgId, bo.Owner)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	ownerResp := vo.UserIDInfo{
		UserID:     bo.Owner,
		Name:       ownerInfo.Name,
		Avatar:     ownerInfo.Avatar,
		EmplID:     ownerInfo.OutOrgUserId,
		IsDeleted:  ownerInfo.OrgUserIsDelete == consts.AppIsDeleted,
		IsDisabled: ownerInfo.OrgUserStatus == consts.AppStatusDisabled,
	}

	infoResp := vo.OrganizationInfoResp{
		OrgID:         bo.Id,
		OrgName:       bo.Name,
		Code:          bo.Code,
		WebSite:       bo.WebSite,
		IndustryID:    bo.IndustryId,
		IndustryName:  resp.AreaInfoResp.IndustryName,
		Scale:         bo.Scale,
		CountryID:     bo.CountryId,
		CountryCname:  resp.AreaInfoResp.CountryCname,
		ProvinceID:    bo.ProvinceId,
		ProvinceCname: resp.AreaInfoResp.ProvinceCname,
		CityID:        bo.CityId,
		CityCname:     resp.AreaInfoResp.CityCname,
		Address:       bo.Address,
		LogoURL:       bo.LogoUrl,
		Owner:         bo.Owner,
		OwnerInfo:     &ownerResp,
		Remark:        bo.Remark,
	}

	outOrgInfo, outOrgInfoErr := domain.GetOrgOutInfo(bo.Id)
	if outOrgInfoErr != nil {
		log.Infof("[OrganizationInfo] orgId: %d, domain.GetOrgOutInfo err: %v", bo.Id, outOrgInfoErr)
	} else if outOrgInfo.TenantCode == "" {
	}
	return &infoResp, nil
}

func GetOrgOutInfoByOutOrgId(orgId int64, outOrgId string) (*orgvo.OutOrgInfo, errs.SystemErrorInfo) {
	outOrgPo, err := domain.GetOrgOutInfoByOutOrgId(orgId, outOrgId)
	if err != nil {
		return nil, err
	}
	outOrgBo := &orgvo.OutOrgInfo{}
	_ = copyer.Copy(outOrgPo, outOrgBo)
	return outOrgBo, nil
}

func UpdateOrgRemarkSetting(input orgvo.UpdateOrgRemarkSettingReqVo) (bool, errs.SystemErrorInfo) {
	// Owns转让的成员需要判断是否在这个组织里面 暂定
	orgBo, err := domain.GetOrgBoById(input.OrgId)
	if err != nil {
		log.Error(err)
		return false, err
	}
	orgRemarkJson := orgBo.Remark
	orgRemarkObj := &orgvo.OrgRemarkConfigType{}
	if len(orgRemarkJson) > 0 {
		oriErr := json.FromJson(orgRemarkJson, orgRemarkObj)
		if oriErr != nil {
			log.Errorf("[UpdateOrgRemarkSetting]失败，组织id:%d,原因:%s", orgBo.Id, oriErr)
			return false, errs.BuildSystemErrorInfo(errs.JSONConvertError, oriErr)
		}
	}
	newRemarkBo := orgvo.SaveOrgSummaryTableAppIdReqVoData{
		AppId:              orgRemarkObj.OrgSummaryTableAppId,
		TagAppId:           orgRemarkObj.TagAppId,
		PriorityAppId:      orgRemarkObj.PriorityAppId,
		IssueBarAppId:      orgRemarkObj.IssueBarAppId,
		SceneAppId:         orgRemarkObj.SceneAppId,
		IssueStatusAppId:   orgRemarkObj.IssueStatusAppId,
		EmptyProjectAppId:  orgRemarkObj.EmptyProjectAppId,
		ProjectFormAppId:   orgRemarkObj.ProjectFormAppId,
		ProjectFolderAppId: orgRemarkObj.ProjectFolderAppId,
		IssueFolderAppId:   orgRemarkObj.IssueFolderAppId,
	}

	if input.Input.AppId > 0 {
		newRemarkBo.AppId = input.Input.AppId
	}
	if input.Input.EmptyProjectAppId > 0 {
		newRemarkBo.ProjectFolderAppId = input.Input.EmptyProjectAppId
	}
	if input.Input.ProjectFormAppId > 0 {
		newRemarkBo.ProjectFormAppId = input.Input.ProjectFormAppId
	}
	if input.Input.ProjectFolderAppId > 0 {
		newRemarkBo.ProjectFolderAppId = input.Input.ProjectFolderAppId
	}
	if input.Input.IssueFolderAppId > 0 {
		newRemarkBo.IssueFolderAppId = input.Input.IssueFolderAppId
	}

	_, err = domain.SaveOrgSomeTableAppId(input.OrgId, input.UserId, newRemarkBo)
	if err != nil {
		log.Error(err)
		return false, err
	}

	return true, nil
}

// 对于自己创建的组织，暂时不支持转让
//
// 对于加入的企业只有查看全，无操作权
//
// 暂时只做基本设置
func UpdateOrganizationSetting(req orgvo.UpdateOrganizationSettingReqVo) (int64, errs.SystemErrorInfo) {

	input := req.Input
	// Owns转让的成员需要判断是否在这个组织里面 暂定
	organizationBo, err := domain.GetOrgBoById(input.OrgID)

	if err != nil {
		return 0, err
	}
	// 针对更新组织负责人，如果不是所有者，则不可以更改。这个接口不会对负责人进行修改。
	// 更新组织普通信息，则需要对”团队设置“项进行权限判断
	authErr := AuthOrgRole(input.OrgID, req.UserId, consts.RoleOperationPathOrgOrgConfig, consts.OperationOrgConfigModify)
	if authErr != nil {
		log.Error(authErr)
		return 0, authErr
	}

	//更改Own的接口拆开来,  这一期暂时也不做
	updateOrgBo, err := assemblyOrganizationBo(input, req.UserId, organizationBo)

	if err != nil {
		return 0, err
	}

	err = domain.UpdateOrg(*updateOrgBo)

	if err != nil {
		return 0, err
	}
	return input.OrgID, nil
}

func assemblyOwner(userId int64, input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "owner") && input.Owner != nil {
		orgId := input.OrgID
		organizationBo, err := domain.GetOrgBoById(orgId)
		if err != nil {
			log.Error(err)
			return err
		}
		if organizationBo.Owner == *input.Owner {
			//负责人没有变动
			return nil
		}
		//不是所有者不可以更改信息
		if organizationBo.Owner != userId {
			return errs.BuildSystemErrorInfo(errs.OrgOwnTransferError)
		}
		//查看新用户是否存在
		_, userErr := domain.GetBaseUserInfo(orgId, *input.Owner)
		if userErr != nil {
			log.Error(userErr)
			return userErr
		}

		userfacade.AddUserToSysManageGroup(orgId, userId, uservo.AddUserToSysManageGroupReq{
			UserIds: []int64{*input.Owner},
		})
		//更新用户角色 要更新无码的管理组
		//_, err = service.UpdateOrgAdmin(orgId, userId, organizationBo.Owner, *input.Owner)
		//if err != nil {
		//	log.Error(err)
		//	return err
		//}
		//修改组织拥有者
		(*upd)[consts.TcOwner] = *input.Owner
	}

	return nil
}

func assemblyOrganizationBo(input vo.UpdateOrganizationSettingsReq, userId int64, orgOrganization *bo.OrganizationBo) (*bo.UpdateOrganizationBo, errs.SystemErrorInfo) {
	//公用初始化
	orgBo := bo.OrganizationBo{Id: input.OrgID}

	upd := &mysql.Upd{}
	//名字
	err := assemblyOrgName(input, upd)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//网址
	err = assemblyCode(input, upd)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//行业
	assemblyIndustryID(input, upd)
	//组织规模
	assemblyScaleID(input, upd)
	//所在国家
	assemblyCountryID(input, upd)
	// 所在省份
	assemblyProvince(input, upd)
	// 所在城市
	assemblyCity(input, upd)
	// 组织地址
	err = assemblyAddress(input, upd)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	// 组织logo地址
	err = assemblyLogoUrl(input, upd)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//组织负责人
	err = assemblyOwner(userId, input, upd)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//sourceChannel
	orgBo.SourceChannel = orgOrganization.SourceChannel

	return &bo.UpdateOrganizationBo{
		Bo:                     orgBo,
		OrganizationUpdateCond: *upd,
	}, nil
}

func assemblyOrgName(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {
	if NeedUpdate(input.UpdateFields, "orgName") {
		orgName := strings.TrimSpace(input.OrgName)
		isOrgNameRight := format.VerifyOrgNameFormat(orgName)
		if !isOrgNameRight {
			return errs.OrgNameLenError
		}

		(*upd)[consts.TcName] = orgName
	}
	return nil
}

// 网址
func assemblyCode(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "code") {
		if input.Code != nil {
			orgCode := *input.Code
			orgCode = strings.TrimSpace(orgCode)

			//判断当前组织有没有设置过code
			organizationBo, err := domain.GetOrgBoById(input.OrgID)
			if err != nil {
				log.Error(err)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
			if organizationBo.Code != consts.BlankString {
				return errs.BuildSystemErrorInfo(errs.OrgCodeAlreadySetError)
			}

			//orgCodeLen := strs.Len(orgCode)
			////判断长度
			//if orgCodeLen > sconsts.OrgCodeLength || orgCodeLen < 1 {
			//	return errs.BuildSystemErrorInfo(errs.OrgCodeLenError)
			//}
			isOrgCodeRight := format.VerifyOrgCodeFormat(orgCode)
			if !isOrgCodeRight {
				return errs.OrgCodeLenError
			}

			_, err = domain.GetOrgBoByCode(orgCode)
			//查不到才能更改
			if err != nil {
				(*upd)[consts.TcCode] = orgCode
			} else {
				return errs.BuildSystemErrorInfo(errs.OrgCodeExistError)
			}
		}
	}

	return nil
}

// 组织行业Id
func assemblyIndustryID(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "industryId") {

		if input.IndustryID != nil {
			(*upd)[consts.TcIndustryId] = *input.IndustryID
		} else {
			(*upd)[consts.TcIndustryId] = 0
		}
	}
}

// 组织规模
func assemblyScaleID(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "scale") {

		if input.Scale != nil {
			(*upd)[consts.TcScale] = *input.Scale
		} else {
			(*upd)[consts.TcScale] = 0
		}
	}
}

// 所在国家
func assemblyCountryID(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "countryId") {

		if input.CountryID != nil {
			(*upd)[consts.TcCountryId] = *input.CountryID
		} else {
			(*upd)[consts.TcCountryId] = 0
		}
	}
}

// 省份
func assemblyProvince(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "provinceId") {

		if input.ProvinceID != nil {
			(*upd)[consts.TcProvinceId] = *input.ProvinceID
		} else {
			(*upd)[consts.TcProvinceId] = 0
		}
	}
}

// 城市
func assemblyCity(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) {

	if NeedUpdate(input.UpdateFields, "cityId") {

		if input.CityID != nil {
			(*upd)[consts.TcCityId] = *input.CityID
		} else {
			(*upd)[consts.TcCityId] = 0
		}
	}
}

// 地址
func assemblyAddress(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "address") {

		if input.Address != nil {
			//len := strs.Len(*input.Address)
			//if len > 256 {
			//	return errs.BuildSystemErrorInfo(errs.OrgAddressLenError)
			//}
			isAdressRight := format.VerifyOrgAdressFormat(*input.Address)
			if !isAdressRight {
				return errs.OrgAddressLenError
			}

			(*upd)[consts.TcAddress] = *input.Address
		}
	}
	return nil
}

// Logo
func assemblyLogoUrl(input vo.UpdateOrganizationSettingsReq, upd *mysql.Upd) errs.SystemErrorInfo {

	if NeedUpdate(input.UpdateFields, "logoUrl") {
		if input.LogoURL != nil {
			logoLen := strs.Len(*input.LogoURL)
			if logoLen > 512 {
				return errs.BuildSystemErrorInfo(errs.OrgLogoLenError)
			}

			(*upd)[consts.TcLogoUrl] = *input.LogoURL
		}
	}
	return nil
}

func ScheduleOrganizationPageList(reqVo orgvo.ScheduleOrganizationPageListReqVo) (*orgvo.ScheduleOrganizationPageListResp, errs.SystemErrorInfo) {

	page := reqVo.Page
	size := reqVo.Size

	bos, count, err := domain.ScheduleOrganizationPageList(page, size)

	if err != nil {
		return nil, err
	}

	list := &[]*orgvo.ScheduleOrganizationListResp{}

	copyErr := copyer.Copy(bos, list)

	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &orgvo.ScheduleOrganizationPageListResp{
		Total:                        count,
		ScheduleOrganizationListResp: list,
	}, nil

}

// 通过来源获取组织id列表
func GetOrgIdListBySourceChannel(sourceChannels []string, page int, size int, isPaid int) ([]int64, errs.SystemErrorInfo) {
	return domain.GetOrgIdListBySourceChannel(sourceChannels, page, size, isPaid)
}

// GetOrgIdListByPage 分页获取组织 id 列表
func GetOrgIdListByPage(input orgvo.GetOrgIdListByPageReqVoData, page int, size int) ([]int64, errs.SystemErrorInfo) {
	return domain.GetOrgIdListByPage(&input, page, size)
}

// GetOrgBoListByPage 分页获取组织列表
func GetOrgBoListByPage(input orgvo.GetOrgIdListByPageReqVoData, page int, size int) (orgvo.GetOrgBoListByPageRespData, errs.SystemErrorInfo) {
	result := orgvo.GetOrgBoListByPageRespData{
		List:  make([]bo.OrganizationBo, 0),
		Total: 0,
		Page:  0,
		Size:  0,
	}
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	cond := db.Cond{}
	if len(input.OrgIds) > 0 {
		cond[consts.TcId] = db.In(input.OrgIds)
	}
	totalRes := po.PpmOrgUserOrganizationCount{}
	if err := conn.Select(db.Raw("count(1) as total")).From(consts.TableOrganization).Where(cond).One(&totalRes); err != nil {
		log.Error(err)
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	result.Total = int64(totalRes.Total)
	mid := conn.Select(db.Raw("*")).From(consts.TableOrganization).Where(cond).OrderBy("id desc")
	if page > 0 && size > 0 {
		mid = mid.Offset((page - 1) * size).Limit(size)
	} else {
		mid = mid.Offset(0).Limit(20)
	}
	orgs := make([]po.PpmOrgOrganization, 0)
	selectErr := mid.All(&orgs)
	if selectErr != nil {
		log.Error(selectErr)
		return result, errs.BuildSystemErrorInfo(errs.MysqlOperateError, selectErr)
	}
	orgBoList := make([]bo.OrganizationBo, 0)
	if copyErr := copyer.Copy(orgs, &orgBoList); copyErr != nil {
		log.Error(copyErr)
		return result, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	result.Page = page
	result.Size = size
	result.List = orgBoList

	return result, nil
}

// GetOrgConfig 获取组织配置
func GetOrgConfig(orgId int64) (*orgvo.OrgConfig, errs.SystemErrorInfo) {
	return domain.GetOrgConfigRich(orgId)
}

func GetOutOrgInfoByOrgIdBatch(orgIds []int64) ([]*orgvo.OutOrgInfo, errs.SystemErrorInfo) {
	list, err := domain.GetOutOrgInfoByOrgIdBatch(orgIds)
	if err != nil {
		return nil, err
	}
	var outOrgInfoList []*orgvo.OutOrgInfo
	copyer.Copy(list, &outOrgInfoList)
	return outOrgInfoList, nil
}

func GetOrgInfoBo(orgId int64) (*bo.OrganizationBo, errs.SystemErrorInfo) {
	info, infoErr := domain.GetOrgBoById(orgId)
	if infoErr != nil {
		log.Error(infoErr)
		return nil, infoErr
	}

	return info, nil
}

//func UpdateOrgBasicShowSetting(orgId int64, userId int64, params vo.UpdateOrgBasicShowSettingReq) (*vo.Void, errs.SystemErrorInfo) {
//	// 1.系统权限判断
//	authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgTeam, consts.OperationOrgConfigModify)
//	if authErr != nil {
//		log.Error(authErr)
//		return nil, authErr
//	}
//
//	_, err := mysql.UpdateSmartWithCond(consts.TableOrgConfig, db.Cond{
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, mysql.Upd{
//		consts.TcBasicShowSetting: json.ToJsonIgnoreError(params),
//	})
//
//	if err != nil {
//		log.Error(err)
//		return nil, errs.MysqlOperateError
//	}
//
//	clearErr := domain.ClearOrgConfig(orgId)
//	if clearErr != nil {
//		log.Error(clearErr)
//	}
//
//	return &vo.Void{ID: orgId}, nil
//}

// CheckAndSetSuperAdmin 检查组织拥有者是否是超管，如果不是，则设置为超管。
func CheckAndSetSuperAdmin(orgId int64) (bool, errs.SystemErrorInfo) {
	resp := userfacade.CheckAndSetSuperAdmin(orgId)
	if resp.Failure() {
		log.Error(resp.Error())
		return false, resp.Error()
	}
	return resp.Data, nil
}

func StopThirdIntegration(orgId int64, currentUserId int64, sourceChannel string) errs.SystemErrorInfo {
	authErr := AuthOrgRole(orgId, currentUserId, consts.RoleOperationPathOrgTeam, consts.OperationOrgConfigModify)
	if authErr != nil {
		log.Error(authErr)
		return authErr
	}

	_, outInfoErr := domain.GetBaseOrgOutInfo(orgId)
	if outInfoErr != nil {
		log.Error(outInfoErr)
		return outInfoErr
	}

	//删除外部信息
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除组织外部信息
		_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableOrganizationOutInfo, db.Cond{
			consts.TcOrgId:         orgId,
			consts.TcSourceChannel: sourceChannel,
			consts.TcIsDelete:      consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  currentUserId,
		})
		if err1 != nil {
			log.Error(err1)
			return err1
		}
		//删除人员外部信息
		_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableUserOutInfo, db.Cond{
			consts.TcOrgId:         orgId,
			consts.TcSourceChannel: sourceChannel,
			consts.TcIsDelete:      consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  currentUserId,
		})
		if err2 != nil {
			log.Error(err2)
			return err2
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return errs.MysqlOperateError
	}

	//删除缓存
	cacheErr := domain.ClearCacheBaseOrgInfo(orgId)
	if cacheErr != nil {
		log.Errorf("redis err: %q\n", cacheErr)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError, cacheErr)
	}

	return nil
}
