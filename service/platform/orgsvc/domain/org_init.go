package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/idvo"
	"github.com/star-table/polaris-backend/common/model/vo/lc_table"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
	"github.com/star-table/polaris-backend/facade/appfacade"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func OrgInit(corpId string, permanentCode string, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	return 0, nil
}

func dealOrgOutInfoIdAndorgId(orgOutInfo *po.PpmOrgOrganizationOutInfo, org *po.PpmOrgOrganization) (int64, int64, errs.SystemErrorInfo) {
	orgOutInfoIdVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: orgOutInfo.TableName()})
	if orgOutInfoIdVo.Failure() {
		return int64(0), int64(0), orgOutInfoIdVo.Error()
	}

	orgIdVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: org.TableName()})
	if orgIdVo.Failure() {
		return int64(0), int64(0), orgIdVo.Error()
	}

	return orgOutInfoIdVo.Id, orgIdVo.Id, nil
}

func OrgOwnerInit(orgId int64, owner, creator int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	org := &po.PpmOrgOrganization{}
	org.Id = orgId
	org.Owner = owner
	org.Creator = creator
	err := mysql.TransUpdate(tx, org)
	if err != nil {
		log.Error(strs.ObjectToString(err))
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}

//组装外部信息
//func AssemblyOrgOutInfo(corpId string, orgOutInfo *po.PpmOrgOrganizationOutInfo) errs.SystemErrorInfo {
//	//获取企业授权信息
//	authInfo, err := GetCorpAuthInfo(corpId)
//	if err != nil {
//		return err
//	}
//
//	authCorpInfo := authInfo.AuthCorpInfo
//
//	orgOutInfo.Name = authCorpInfo.CorpName
//	orgOutInfo.Industry = authCorpInfo.Industry
//
//	isAuth := 0
//	if authCorpInfo.IsAuthenticated {
//		isAuth = 1
//	}
//	orgOutInfo.IsAuthenticated = isAuth
//	orgOutInfo.AuthLevel = strconv.FormatInt(authCorpInfo.AuthLevel, 10)
//
//	return nil
//}

//func AssemblyOrg(corpId string, org *po.PpmOrgOrganization) errs.SystemErrorInfo {
//	authInfo, err := GetCorpAuthInfo(corpId)
//	if err != nil {
//		return err
//	}
//
//	authCorpInfo := authInfo.AuthCorpInfo
//
//	org.Name = authCorpInfo.CorpName
//	org.LogoUrl = authCorpInfo.CorpLogoUrl
//	org.Address = authCorpInfo.CorpProvince + authCorpInfo.CorpCity
//
//	isAuth := 0
//	if authCorpInfo.IsAuthenticated {
//		isAuth = 1
//	}
//	org.IsAuthenticated = isAuth
//
//	return nil
//}

//func GetCorpAuthInfo(corpId string) (sdk.GetAuthInfoResp, errs.SystemErrorInfo) {
//	suiteTicket, err := GetSuiteTicket()
//	if err != nil {
//		return sdk.GetAuthInfoResp{}, err
//	}
//	//创建企业对象
//	corpProxy := dingtalk.GetSDKProxy().CreateCorp(corpId, suiteTicket)
//	resp, err2 := corpProxy.GetAuthInfo()
//	if err2 != nil {
//		return sdk.GetAuthInfoResp{}, errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err2)
//	}
//
//	if resp.ErrCode != 0 {
//		return sdk.GetAuthInfoResp{}, errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, errors.New(resp.ErrMsg))
//	}
//	resp, err2 = corpProxy.GetAuthInfo()
//	if err2 != nil {
//		return resp, errs.BuildSystemErrorInfo(errs.DingTalkOpenApiCallError, err2)
//	}
//	return resp, nil
//}

// CreateIssueSummaryTable 初始化企业后的执行逻辑
// 调用无码的 app 服务创建任务**汇总表**（创建项目时触发）
func CreateIssueSummaryTable(orgId, opUserId int64, appName string, remarkObj *orgvo.OrgRemarkConfigType, pkgId int64) (int64, errs.SystemErrorInfo) {
	configJson := GetCreateTableConfigForSummaryIssue(remarkObj)
	appId := int64(0)
	// 4表示极星项目
	appType := consts.LcAppTypeForSummaryTable
	resp := appfacade.CreateLessCodeApp(&permissionvo.CreateLessCodeAppReq{
		OrgId:        &orgId,
		AppType:      &appType,
		Name:         &appName,
		UserId:       &opUserId,
		PkgId:        pkgId,
		Config:       configJson,
		AddAllMember: true, // 组织下，所有成员可见
		ProjectId:    -1,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return appId, resp.Error()
	}
	appId = resp.Data.Id
	return appId, nil
}

// GetCreateTableConfigForSummaryIssue 汇总表的配置（新的表头定义）
func GetCreateTableConfigForSummaryIssue(remarkObj *orgvo.OrgRemarkConfigType) string {
	falseFlag := true
	fields := []interface{}{
		lc_helper.GetLcCtTextArea(consts.BasicFieldTitle, "标题", "Title", true, true, false, true),
		//lc_helper.GetLcCtInputFull(consts.BasicFieldTitle, "标题", "Title", false, false, true, true, false, false, true),
		lc_helper.GetLcCtSelect(consts.BasicFieldProjectId, "所属项目", "Project", "select", nil, true, true, false, false, true),
		lc_helper.GetLcCtSelect(consts.BasicFieldIterationId, "迭代", "Sprint", "select", lc_helper.GetDefaultSelectOptionsForIterationId(), true, true, false, false, true),
		lc_helper.GetLcCtTextArea(consts.BasicFieldParentId, "父任务ID", "Parent ID", true, true, true, false),
		//lc_helper.GetLcCtInputFull(consts.BasicFieldParentId, "父任务ID", "Parent ID", false, false, true, true, true, true, false),
		//lc_helper.GetLcCtInputFull(consts.BasicFieldCode, "编号", "ID Number", false, false, true, true, true, false, false),
		lc_helper.GetLcCtTextArea(consts.BasicFieldCode, "编号", "ID Number", true, true, false, false),
		lc_table.LcCommonField{
			Label:    "描述",
			EnLabel:  "Description",
			Name:     consts.BasicFieldRemark,
			Editable: &falseFlag,
			Writable: true,
			Field: lc_table.LcFieldData{
				Type:  "richtext",
				Props: lc_table.LcProps{PushMsg: false},
			},
		},
		lc_helper.GetLcCtMember(consts.BasicFieldOwnerId, "负责人", "Owners", true, true, false, 1, true, true),
		lc_helper.GetLcCtSelect(consts.BasicFieldProjectObjectTypeId, "任务栏", "Task Bar", "select", lc_helper.GetDefaultSelectOptionsForTaskBar(), true, true, true, false, false),
		lc_helper.GetLcCtGroupSelect(consts.BasicFieldIssueStatus, "任务状态", "groupSelect", lc_helper.GetDefaultGroupSelectForIssueStatus(), true, true),
		lc_helper.GetLcCtDatepicker(consts.BasicFieldPlanStartTime, "开始时间", "Start Date", true, true),
		lc_helper.GetLcCtDatepicker(consts.BasicFieldPlanEndTime, "截止时间", "Due Date", true, true),
		lc_helper.GetLcCtMember(consts.BasicFieldFollowerIds, "关注人", "Collaborators", true, true, true, 0, true, true),
		lc_helper.GetLcCtMember(consts.BasicFieldAuditorIds, "确认人", "Operators", true, true, true, 0, true, true),
	}
	config := lc_helper.NewLcTableConfig(fields)
	configJson := json.ToJsonIgnoreError(config)
	return configJson
}

// 将汇总表id 存入组织表的字段中
// 将数据以 json 配置的方式存储在 remark 中。
func SaveOrgSummaryTableAppId(orgId, userId int64, summaryAppId int64) (bool, errs.SystemErrorInfo) {
	orgBo, err := GetOrgBoById(orgId)
	if err != nil {
		return false, err
	}
	upds := mysql.Upd{
		consts.TcRemark: AssemblyOrgRemarkConfigForAppId(orgBo, summaryAppId),
	}
	_, oriErr := mysql.UpdateSmartWithCond(consts.TableOrganization, db.Cond{
		consts.TcId: orgId,
	}, upds)
	if oriErr != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	return true, nil
}

func SaveOrgSomeTableAppId(orgId, userId int64, input orgvo.SaveOrgSummaryTableAppIdReqVoData, tx ...sqlbuilder.Tx) (bool, errs.SystemErrorInfo) {
	orgBo, err := GetOrgBoById(orgId, tx...)
	if err != nil {
		return false, err
	}
	upds := mysql.Upd{
		consts.TcRemark: AssemblyOrgRemarkConfigForSomeAppId(orgBo, input),
	}
	var oriErr error
	if len(tx) > 0 {
		_, oriErr = mysql.TransUpdateSmartWithCond(tx[0], consts.TableOrganization, db.Cond{
			consts.TcId: orgId,
		}, upds)
	} else {
		_, oriErr = mysql.UpdateSmartWithCond(consts.TableOrganization, db.Cond{
			consts.TcId: orgId,
		}, upds)
	}

	if oriErr != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	return true, nil
}

func AssemblyOrgRemarkConfigForAppId(orgBo *bo.OrganizationBo, value int64) string {
	remark := orgBo.Remark
	defaultVal := &orgvo.OrgRemarkConfigType{}
	if len(remark) < 1 {
		defaultVal.OrgSummaryTableAppId = value
	} else {
		_ = json.FromJson(orgBo.Remark, defaultVal)
		defaultVal.OrgSummaryTableAppId = value
	}

	return json.ToJsonIgnoreError(defaultVal)
}

func AssemblyOrgRemarkConfigForSomeAppId(orgBo *bo.OrganizationBo, valueObj orgvo.SaveOrgSummaryTableAppIdReqVoData) string {
	remark := orgBo.Remark
	defaultVal := &orgvo.OrgRemarkConfigType{}
	if len(remark) < 1 {
		// do nothing
	} else {
		_ = json.FromJson(orgBo.Remark, defaultVal)
	}
	if err := copyer.Copy(valueObj, &defaultVal); err != nil {
		log.Error(err)
	}

	// 特殊一点的，单独赋值
	defaultVal.OrgSummaryTableAppId = valueObj.AppId

	return json.ToJsonIgnoreError(defaultVal)
}

// 创建一个目录用于存放表/表单
func CreateLcFolder(orgId, opUserId int64, appName string) (appId int64, info errs.SystemErrorInfo) {
	// 4表示极星项目
	appType := consts.LcAppTypeForFolder
	resp := appfacade.CreateLessCodeApp(&permissionvo.CreateLessCodeAppReq{
		OrgId:   &orgId,
		AppType: &appType,
		Name:    &appName,
		UserId:  &opUserId,
		Config:  "",
		PkgId:   0,
	})
	if resp.Failure() {
		log.Error(resp.Error())
		return appId, resp.Error()
	}
	appId = resp.Data.Id
	return appId, nil
}
