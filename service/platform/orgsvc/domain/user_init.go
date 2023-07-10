package domain
//
//import (
//	"strconv"
//	"strings"
//
//	sdk_const "gitea.bjx.cloud/allstar/platform-sdk/consts"
//
//	"github.com/star-table/common/core/util/slice"
//	"github.com/star-table/common/core/util/strs"
//	"github.com/star-table/common/library/db/mysql"
//	"github.com/star-table/common/sdk/dingtalk"
//	"github.com/star-table/polaris-backend/common/core/consts"
//	"github.com/star-table/polaris-backend/common/core/errs"
//	"github.com/star-table/polaris-backend/common/core/util/pinyin"
//	dingtalk2 "github.com/star-table/polaris-backend/common/extra/dingtalk"
//	"github.com/star-table/polaris-backend/facade/idfacade"
//	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
//	"github.com/pkg/errors"
//	"github.com/polaris-team/dingtalk-sdk-golang/sdk"
//	"upper.io/db.v3"
//	"upper.io/db.v3/lib/sqlbuilder"
//)
//
////组织初始化时-用户初始化
//func UserInitByOrg(userId string, corpId string, orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
//	user := &po.PpmOrgUser{}
//	userOutInfo := &po.PpmOrgUserOutInfo{}
//
//	userDetail, err := GetUserDetail(userId, corpId)
//
//	if err != nil {
//		return 0, err
//	}
//
//	unionId := userDetail.UnionId
//
//	registered := true
//	registeredUserOutInfo := &po.PpmOrgUserOutInfo{}
//	//判断用户是否注册过
//	err2 := tx.Collection(userOutInfo.TableName()).Find(db.Cond{consts.TcOutUserId: unionId, consts.TcSourceChannel: sdk_const.SourceChannelDingTalk}).One(registeredUserOutInfo)
//	if err2 != nil {
//		registered = false
//		log.Info("当前用户没有注册过")
//	}
//
//	err2 = tx.Collection(userOutInfo.TableName()).Find(db.Cond{consts.TcOrgId: orgId, consts.TcOutOrgUserId: userId, consts.TcOutUserId: unionId, consts.TcSourceChannel: sdk_const.SourceChannelDingTalk}).One(userOutInfo)
//	if err2 == nil {
//
//		userOutInfoId, err2 := restoreCancellateUser(userOutInfo, userDetail, tx, err)
//
//		if err2 != nil {
//			return 0, err2
//		}
//		return userOutInfoId, nil
//
//	} else {
//		log.Info("user insert")
//
//		userConfig := &po.PpmOrgUserConfig{}
//		userOutInfoId, userNativeId, userConfigId, err := assemblyUserInfo(userOutInfo, user)
//
//		if err != nil {
//			return 0, err
//		}
//
//		userOrg := &po.PpmOrgUserOrganization{}
//		userOrgId, err := idfacade.ApplyPrimaryIdRelaxed(userOrg.TableName())
//
//		registeredError := assemblyRegisteredInfo(registered, orgId, userConfigId, &userNativeId, user, userDetail,
//			userConfig, tx, registeredUserOutInfo)
//
//		if registeredError != nil {
//			return 0, registeredError
//		}
//
//		userOutInfo.Id = userOutInfoId
//		userOutInfo.UserId = userNativeId
//		userOutInfo.OrgId = orgId
//		userOutInfo.OutOrgUserId = userId
//		userOutInfo.OutUserId = unionId
//		userOutInfo.IsDelete = consts.AppIsNoDelete
//		userOutInfo.Status = consts.AppStatusEnable
//		userOutInfo.SourceChannel = sdk_const.SourceChannelDingTalk
//		AssemblyUserOutInfo(userOutInfo, *userDetail)
//
//		err2 = mysql.TransInsert(tx, userOutInfo)
//		if err2 != nil {
//			log.Error(strs.ObjectToString(err2))
//			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
//		}
//		log.Info("初始化用户外部信息成功")
//
//		userOrg.Id = userOrgId
//		userOrg.OrgId = orgId
//		userOrg.UserId = userNativeId
//		userOrg.CheckStatus = consts.AppCheckStatusSuccess
//		userOrg.IsDelete = consts.AppIsNoDelete
//		err2 = mysql.TransInsert(tx, userOrg)
//		if err2 != nil {
//			log.Error("组织用户关联失败" + strs.ObjectToString(err2))
//			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
//		}
//		log.Info("初始化用户组织关联成功")
//
//		//初始化部门用户
//		initDepartmentUserErr := initDepartmentUser(userNativeId, userDetail, tx, orgId)
//		if initDepartmentUserErr != nil {
//			log.Error("部门用户初始化失败" + strs.ObjectToString(initDepartmentUserErr))
//			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, initDepartmentUserErr)
//		}
//		log.Info("部门用户初始化成功")
//
//		//
//		//if isRoot{
//		//	org.Id = orgId
//		//	org.Owner = userNativeId
//		//	err = mysql.TransUpdate(tx, org)
//		//	if err != nil{
//		//		log.Error(err)
//		//		return 0, err
//		//	}
//		//}
//		return userNativeId, nil
//	}
//}
//
//func initDepartmentUser(userNativeId int64, userDetail *sdk.GetUserDetailResp, tx sqlbuilder.Tx, orgId int64) errs.SystemErrorInfo {
//	//暂不考虑后期增加了部门导致的部门不存在
//	//1.获取外部部门与内部部门的id对应关系
//	deptRelationList, err := GetOutDeptAndInnerDept(orgId, &tx)
//	if err != nil {
//		return err
//	}
//	IsLeaderInDepts := dealLeaderString(userDetail.IsLeaderInDepts)
//
//	//该用户为负责人的所有部门
//	leaderDepts := []string{}
//	for k, v := range IsLeaderInDepts {
//		if v == "true" {
//			leaderDepts = append(leaderDepts, k)
//		}
//	}
//	departmentUserList := []interface{}{}
//	for _, v := range userDetail.Department {
//		deptId, ok := deptRelationList[strconv.FormatInt(v, 10)]
//		if !ok {
//			continue
//		}
//
//		id, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserDepartment)
//		if err != nil {
//			return err
//		}
//		isLeader := 2 //默认不是负责人
//		if bol, _ := slice.Contain(leaderDepts, strconv.FormatInt(v, 10)); bol {
//			isLeader = 1
//		}
//		departmentUserList = append(departmentUserList, po.PpmOrgUserDepartment{
//			Id:           id,
//			OrgId:        orgId,
//			UserId:       userNativeId,
//			DepartmentId: deptId,
//			IsLeader:     isLeader,
//		})
//	}
//	if len(deptRelationList) == 0 {
//		return nil
//	}
//	insertErr := mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, departmentUserList)
//	if insertErr != nil {
//		log.Error(insertErr)
//		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, insertErr)
//	}
//
//	return nil
//}
//
//func dealLeaderString(str string) map[string]string {
//	res := map[string]string{}
//	str = strings.ReplaceAll(str, " ", "")
//	str = str[1 : len(str)-1]
//	strArr := strings.Split(str, ",")
//	if len(strArr) == 0 {
//		return res
//	}
//	for _, v := range strArr {
//		kv := strings.Split(v, ":")
//		if len(kv) >= 2 {
//			res[kv[0]] = kv[1]
//		}
//	}
//
//	return res
//}
//
////组装用户信息
//func assemblyUserInfo(userOutInfo *po.PpmOrgUserOutInfo, user *po.PpmOrgUser) (int64, int64, int64, errs.SystemErrorInfo) {
//
//	userOutInfoId, err := idfacade.ApplyPrimaryIdRelaxed(userOutInfo.TableName())
//	if err != nil {
//		log.Error(strs.ObjectToString(err))
//		return 0, 0, 0, err
//	}
//	userNativeId, err := idfacade.ApplyPrimaryIdRelaxed(user.TableName())
//	if err != nil {
//		log.Error(strs.ObjectToString(err))
//		return 0, 0, 0, err
//	}
//	userConfigId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableUserConfig)
//	if err != nil {
//		log.Error(strs.ObjectToString(err))
//		return 0, 0, 0, err
//	}
//
//	return userOutInfoId, userNativeId, userConfigId, nil
//}
//
////恢复注销用户
//func restoreCancellateUser(userOutInfo *po.PpmOrgUserOutInfo, userDetail *sdk.GetUserDetailResp, tx sqlbuilder.Tx, err errs.SystemErrorInfo) (int64, errs.SystemErrorInfo) {
//	log.Info("user update ")
//	//用户之前被注销掉，恢复状态
//	if userOutInfo.IsDelete == consts.AppIsDeleted {
//		userOutInfo.IsDelete = consts.AppIsNoDelete
//		userOutInfo.SourceChannel = sdk_const.SourceChannelDingTalk
//		AssemblyUserOutInfo(userOutInfo, *userDetail)
//
//		err2 := mysql.TransUpdate(tx, userOutInfo)
//		if err2 != nil {
//			log.Error(strs.ObjectToString(err2))
//			return 0, err
//		}
//	}
//	return userOutInfo.Id, nil
//}
//
//func assemblyRegisteredInfo(registered bool, orgId, userConfigId int64, userNativeId *int64, user *po.PpmOrgUser, userDetail *sdk.GetUserDetailResp,
//	userConfig *po.PpmOrgUserConfig, tx sqlbuilder.Tx, registeredUserOutInfo *po.PpmOrgUserOutInfo) errs.SystemErrorInfo {
//	if !registered {
//		user.Id = *userNativeId
//		user.OrgId = orgId
//		user.SourceChannel = sdk_const.SourceChannelDingTalk
//		AssemblyUser(user, *userDetail)
//
//		err2 := mysql.TransInsert(tx, user)
//		if err2 != nil {
//			log.Error(strs.ObjectToString(err2))
//			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
//		}
//		log.Info("初始化用户成功")
//
//		userConfig.Id = userConfigId
//		userConfig.OrgId = orgId
//		userConfig.UserId = *userNativeId
//
//		err2 = mysql.TransInsert(tx, userConfig)
//		if err2 != nil {
//			log.Error(strs.ObjectToString(err2))
//			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
//		}
//		log.Info("初始化用户配置成功")
//	} else {
//		*userNativeId = registeredUserOutInfo.Id
//	}
//
//	return nil
//}
//
//func GetUserDetail(userId, corpId string) (*sdk.GetUserDetailResp, errs.SystemErrorInfo) {
//	suiteTicket, err := dingtalk2.GetSuiteTicket()
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.SuiteTicketError, err)
//	}
//
//	dingtalkClient, err := dingtalk.GetDingTalkClient(corpId, suiteTicket)
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.DingTalkClientError, err)
//	}
//
//	lang := consts.AppSourceChannelDingTalkDefaultLang
//	userDetailResp, err := dingtalkClient.GetUserDetail(userId, &lang)
//	if err != nil {
//		return nil, errs.BuildSystemErrorInfo(errs.DingTalkGetUserInfoError, err)
//	}
//
//	if userDetailResp.ErrCode != 0 {
//		return nil, errs.BuildSystemErrorInfo(errs.DingTalkGetUserInfoError, errors.New(userDetailResp.ErrMsg))
//	}
//	return &userDetailResp, nil
//}
//
//func AssemblyUserOutInfo(userOutInfo *po.PpmOrgUserOutInfo, userDetailResp sdk.GetUserDetailResp) {
//	isActive := 0
//	if userDetailResp.Active {
//		isActive = 1
//	}
//
//	userOutInfo.Name = userDetailResp.Name
//	userOutInfo.Avatar = userDetailResp.Avatar
//	userOutInfo.IsActive = isActive
//	userOutInfo.JobNumber = userDetailResp.JobNumber
//}
//
//func AssemblyUser(user *po.PpmOrgUser, userDetailResp sdk.GetUserDetailResp) {
//	user.Name = userDetailResp.Name
//	user.NamePinyin = pinyin.ConvertToPinyin(user.Name)
//	user.LoginName = user.Name
//	if user.Avatar == "" {
//		user.Avatar = userDetailResp.Avatar
//	}
//}
