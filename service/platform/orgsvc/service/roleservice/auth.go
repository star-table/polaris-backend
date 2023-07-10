package service

import (
	"fmt"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"strings"

	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/businees"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo"
	"github.com/star-table/polaris-backend/common/model/vo/permissionvo/appauth"

	"github.com/star-table/polaris-backend/facade/permissionfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/spf13/cast"
)

func Authenticate(orgId int64, userId int64, projectAuthInfo *bo.ProjectAuthBo, issueAuthInfo *bo.IssueAuthBo, path string, operation string, authFields []string) errs.SystemErrorInfo {
	var err errs.SystemErrorInfo
	//如果是issue负责人，拥有所有权限
	//if projectAuthInfo != nil {
	//	isMember, _ := slice.Contain(projectAuthInfo.Participants, userId) //如果是任务创建人要判断是否是项目成员（是的话就拥有所有权限）
	//	if issueAuthInfo != nil && (userId == issueAuthInfo.Owner || (userId == issueAuthInfo.Creator && isMember)) {
	//		log.Infof("权限校验成功，用户 %d 是任务 %d 的负责人/创建人", userId, issueAuthInfo.Id)
	//		return nil
	//	}
	//}

	projectId := int64(0)
	if projectAuthInfo != nil {
		projectId = projectAuthInfo.Id
		//创建人要判断是否是项目成员，移出了项目就没有权限
		isMember, _ := slice.Contain(projectAuthInfo.Participants, userId)
		isProjectOwner, _ := slice.Contain(projectAuthInfo.OwnerIds, userId)
		if userId == projectAuthInfo.Owner || isProjectOwner || (userId == projectAuthInfo.Creator && isMember) {
			log.Infof("权限校验成功，用户 %d 是项目 %d 的负责人/创建人", userId, projectId)
			return nil
		}
	}

	manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
	if manageAuthInfoResp.Failure() {
		log.Error(manageAuthInfoResp.Message)
		return manageAuthInfoResp.Error()
	}
	manageAuthInfo := manageAuthInfoResp.Data
	isOrgAdmin := domain.CheckIsAdminByAuthRespData(manageAuthInfo)
	if isOrgAdmin {
		log.Infof("权限校验成功，用户 %d 是组织 %d 的超级管理员", userId, orgId)
		return nil
	}

	access := false
	if projectAuthInfo == nil && issueAuthInfo == nil {
		opList := TransferOperationArr(manageAuthInfo.OptAuth)
		access, _ = slice.Contain(opList, operation)
	} else if (projectAuthInfo == nil && issueAuthInfo != nil) || ((projectAuthInfo != nil && projectAuthInfo.Id == 0) && issueAuthInfo != nil) {
		// 如果没有分配项目，则直接任务的负责人、创建人、管理员可以管理该任务。
		access, err = CheckIssueAuth(orgId, userId, issueAuthInfo)
		if err != nil {
			log.Errorf("[Authenticate] CheckIssueAuth err: %v, orgId: %d, userId: %d", err, orgId, userId)
			return err
		}
	} else {
		optAuthResp := &permissionvo.GetAppAuthResp{}
		if issueAuthInfo == nil {
			optAuthResp = permissionfacade.GetAppAuth(orgId, projectAuthInfo.AppId, 0, userId)
			if optAuthResp.Failure() {
				log.Error(optAuthResp.Message)
				return optAuthResp.Error()
			}
		} else {
			optAuthResp = permissionfacade.GetAppAuthWithoutCollaborator(orgId, projectAuthInfo.AppId, userId)
			if optAuthResp.Failure() {
				log.Error(optAuthResp.Message)
				return optAuthResp.Error()
			}
		}

		// 鉴权时，检查是否是系统管理员
		if optAuthResp.Data.HasAppRootPermission {
			return nil
		}

		orgConfig, orgConfigErr := domain.GetOrgConfig(orgId)
		if orgConfigErr != nil {
			log.Error(orgConfigErr)
			return orgConfigErr
		}
		var tableId int64
		if issueAuthInfo != nil {
			tableId = issueAuthInfo.TableId
		}
		fieldAuth := fieldCheck(tableId, orgConfig.PayLevel, authFields, operation, optAuthResp.Data)
		if fieldAuth {
			// 将 optAuthResp.NewData 数据元素转换成单个操作项的数组。
			// 此处的 optAuthResp.NewData 形如：`["Permission.Pro.Tag-Create,Modify,Delete"]`
			// operation 形如：Permission.Pro.Tag.Create
			// 要将 operation 与 optAuthResp.NewData 中的元素比较，需要对 optAuthResp.NewData 进行一些组装和处理。
			opList := TransferOperationArr(optAuthResp.Data.OptAuth)
			access, _ = slice.Contain(opList, operation)
			log.Infof("Authenticate opList: %s", json.ToJsonIgnoreError(opList))
		} else {
			access = fieldAuth
			//return errs.NoOperationPermissionForIssueUpdate
		}
		if !access && issueAuthInfo != nil {
			//如果没有通过再判断任务协作人权限
			//todo 数据权限
			dataId, err := getDataIdByIssueId(orgId, userId, issueAuthInfo.Id)
			if err != nil {
				log.Errorf("[Authenticate] getDataIdByIssueId issueId:%v, err:%v", issueAuthInfo.Id, err)
				return err
			}
			if dataId > 0 {
				optAuthResp := permissionfacade.GetDataAuth(orgId, projectAuthInfo.AppId, userId, dataId)
				if optAuthResp.Failure() {
					log.Error(optAuthResp.Message)
					return optAuthResp.Error()
				}

				log.Infof("[Authenticate] GetDataAuth, tableId:%v, authFields:%v, result:%v", tableId, authFields, json.ToJsonIgnoreError(optAuthResp.Data))

				fieldAuth = fieldCheck(tableId, orgConfig.PayLevel, authFields, operation, optAuthResp.Data)
				if fieldAuth {
					opList := TransferOperationArr(optAuthResp.Data.OptAuth)
					access, _ = slice.Contain(opList, operation)
					log.Infof("Authenticate opList: %s", json.ToJsonIgnoreError(opList))
				} else {
					//数据权限 也没有通过字段权限的判断，那就直接报字段没有权限
					return errs.NoOperationPermissionForIssueUpdate
				}
			}
		}
	}

	if access {
		return nil
	}
	return errs.BuildSystemErrorInfo(errs.NoOperationPermissions)
}

func getDataIdByIssueId(orgId, userId, issueId int64) (int64, errs.SystemErrorInfo) {
	condition := &tablePb.Condition{
		Type: tablePb.ConditionType_and,
		Conditions: []*tablePb.Condition{
			&tablePb.Condition{
				Type:   tablePb.ConditionType_equal,
				Value:  json.ToJsonIgnoreError([]interface{}{orgId}),
				Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
			},
			&tablePb.Condition{
				Type:   tablePb.ConditionType_equal,
				Value:  json.ToJsonIgnoreError([]interface{}{issueId}),
				Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
			},
		},
	}
	reply := projectfacade.GetIssueRowList(projectvo.IssueRowListReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &tablePb.ListRawRequest{
			DbType: tablePb.DbType_slave1,
			FilterColumns: []string{lc_helper.ConvertToFilterColumn(consts.BasicFieldId)},
			Condition:     condition,
		},
	})
	if reply.Failure() {
		log.Errorf("[getDataIdByIssueId] err:%v, orgId:%v, userId:%v, issueId:%v", reply.Error(), orgId, userId, issueId)
		return 0, reply.Error()
	}

	if len(reply.Data) > 0 {
		return cast.ToInt64(reply.Data[0][consts.BasicFieldId]), nil
	}

	return 0, nil
}

func fieldCheck(tableId int64, payLevel int, authFields []string, operation string, appAuth appauth.GetAppAuthData) bool {
	if businees.CheckIsPaidVer(payLevel) {
		tableIdStr := cast.ToString(tableId)
		for _, s := range authFields {
			// 这是系统定义的一些，不需要校验
			// id和issueId在无码里会剔除掉，不会更新，所以不需要校验
			if ok, _ := slice.Contain([]string{"lessUpdateIssueReq", consts.BasicFieldId, consts.BasicFieldIssueId}, s); ok {
				continue
			}
			if operation == consts.OperationProIssue4Create {
				if ok, _ := slice.Contain([]string{consts.BasicFieldOwnerId,
					consts.BasicFieldFollowerIds,
					consts.BasicFieldProjectObjectTypeId,
					consts.BasicFieldIssueStatus,
					consts.BasicFieldPriority,
					consts.BasicFieldCreator,
					consts.BasicFieldUpdator,
					consts.BasicFieldIterationId,
					consts.BasicFieldOwnerId,
					consts.BasicFieldProjectId,
					consts.BasicFieldParentId,
				}, s); ok {
					continue
				}
			}

			if !appAuth.HasFieldWriteAuth(tableIdStr, s) {
				log.Infof("fieldCheck 无权字段：%s", s)
				return false
			}
		}
	}

	return true
}

// CheckIssueAuth 检查对任务的操作是否有权限，任务的负责人、创建人、关注人有权限编辑任务。
func CheckIssueAuth(orgId, userId int64, issueAuthBo *bo.IssueAuthBo) (bool, errs.SystemErrorInfo) {
	isFollower, _ := slice.Contain(issueAuthBo.Followers, userId)
	isOk, _ := slice.Contain(issueAuthBo.Owner, userId)
	if isOk || userId == issueAuthBo.Creator || isFollower {
		return true, nil
	}

	return false, nil
}

// 将形如：`["Permission.Pro.Tag-Create,Modify,Delete"]` 的权限数组转换为 operation（格式是：`Permission.Pro.Tag.Create`） 一致的格式。
func TransferOperationArr(optAuthArr []string) []string {
	opList := []string{}
	for _, item := range optAuthArr {
		infos := strings.Split(item, "-")
		if len(infos) > 1 {
			opPrev := infos[0]
			opSuffixArr := strings.Split(infos[1], ",")
			for _, oneSuffix := range opSuffixArr {
				opList = append(opList, fmt.Sprintf("%s.%s", opPrev, oneSuffix))
			}
		} else {
			opList = append(opList, item)
		}
	}
	return opList
}

//func GetUserRoleIds(orgId, userId, projectId int64, projectAuthInfo *bo.ProjectAuthBo, issueAuthInfo *bo.IssueAuthBo) ([]int64, error) {
//	//开始加载用户所有的角色
//	roleUsers, err1 := GetUserRoleListByProjectId(orgId, userId, projectId)
//	if err1 != nil {
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
//	}
//	noSpecific := len(roleUsers) == 0 && projectId != 0
//
//	//如果没有对项目增加特殊权限，查询通用权限
//	if noSpecific {
//		roleUsers, err1 = GetUserRoleListByProjectId(orgId, userId, 0)
//		if err1 != nil {
//			log.Error(err1)
//			return nil, errs.BuildSystemErrorInfo(errs.CacheProxyError)
//		}
//	}
//
//	roleIds := make([]int64, 0)
//	for _, roleUser := range roleUsers {
//		roleIds = append(roleIds, roleUser)
//	}
//
//	//如果没有对项目做单独的权限设定或者查询的是通用权限，那么增加特殊角色
//	err := assemblyRoleIds(&roleIds, projectId, noSpecific, orgId, userId, projectAuthInfo, issueAuthInfo)
//	if err != nil {
//		log.Error(err)
//		return nil, err
//	}
//
//	return roleIds, nil
//}

//func assemblyRoleIds(roleIds *[]int64, projectId int64, noSpecific bool, orgId, userId int64, projectAuthInfo *bo.ProjectAuthBo, issueAuthInfo *bo.IssueAuthBo) (error error) {
//
//	if projectId == 0 || noSpecific {
//		//如果任务信息不为空，则获取任务相关角色，否则获取项目角色
//		if issueAuthInfo != nil {
//			err := loadRoleByProjectAuthInfo(roleIds, orgId, userId, issueAuthInfo.Creator, issueAuthInfo.Followers, issueAuthInfo.Participants)
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//		}
//		if projectAuthInfo != nil {
//			//加载创建者相关信息,参与者 ,关注者
//			err := loadRoleByProjectAuthInfo(roleIds, orgId, userId, projectAuthInfo.Creator, projectAuthInfo.Followers, projectAuthInfo.Participants)
//			if err != nil {
//				log.Error(err)
//				return err
//			}
//		}
//		//组织成员
//		if len(*roleIds) == 0 {
//			orgMemberRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialMember)
//			if err != nil {
//				log.Error(err)
//				return errs.BuildSystemErrorInfo(errs.CacheProxyError)
//			}
//			*roleIds = append(*roleIds, orgMemberRole.Id)
//		}
//
//		//访客
//		visitorRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialVisitor)
//		if err != nil {
//			log.Error(err)
//			return errs.BuildSystemErrorInfo(errs.CacheProxyError)
//		}
//		*roleIds = append(*roleIds, visitorRole.Id)
//	}
//	return nil
//}

//func loadRoleByProjectAuthInfo(roleIds *[]int64, orgId, userId int64, creator int64, followers []int64, participants []int64) (error error) {
//
//	//加载创建者相关信息
//	creatorIsError, creatorError := loadCreateRole(roleIds, orgId, userId, creator)
//	if creatorIsError {
//		log.Error(creatorIsError)
//		return creatorError
//	}
//
//	if participants != nil && len(participants) > 0 {
//		//参与者
//		participantIsError, ParticipantError := loadParticipantRole(roleIds, orgId, userId, participants)
//
//		if participantIsError {
//			log.Error(participantIsError)
//			return ParticipantError
//		}
//	}
//
//	if followers != nil && len(followers) > 0 {
//		//关注者
//		followIsError, followError := loadFollowRole(roleIds, orgId, userId, followers)
//		if followIsError {
//			log.Error(followIsError)
//			return followError
//		}
//	}
//
//	return nil
//}

////加载创建者相关角色
//func loadCreateRole(roleIds *[]int64, orgId, userId int64, creator int64) (isNil bool, error error) {
//
//	if userId == creator {
//		createRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialCreator)
//		if err != nil {
//			log.Error(err)
//			return true, errs.BuildSystemErrorInfo(errs.CacheProxyError)
//		}
//		*roleIds = append(*roleIds, createRole.Id)
//	}
//	//不需要返回 异常为空
//	return false, nil
//}
//
////判断是否是参与者 和加载参与者权限
//func loadParticipantRole(roleIds *[]int64, orgId, userId int64, participants []int64) (isNil bool, error error) {
//
//	isParticipant, err2 := slice.Contain(participants, userId)
//	if err2 != nil {
//		log.Error(err2)
//		return true, errs.BuildSystemErrorInfo(errs.SystemError)
//	}
//	if isParticipant {
//		//workerRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialWorker)
//		workerRole, err := GetRoleByLangCode(orgId, consts.RoleGroupProMember)
//		if err != nil {
//			log.Error(err)
//			return true, errs.BuildSystemErrorInfo(errs.CacheProxyError)
//		}
//		*roleIds = append(*roleIds, workerRole.Id)
//	}
//	return false, nil
//}
//
////关注着
//func loadFollowRole(roleIds *[]int64, orgId, userId int64, followers []int64) (isNil bool, error error) {
//
//	isFollower, err2 := slice.Contain(followers, userId)
//	if err2 != nil {
//		log.Error(err2)
//		return true, errs.BuildSystemErrorInfo(errs.SystemError)
//	}
//	if isFollower {
//		//attentionRole, err := GetRoleByLangCode(orgId, consts.RoleGroupSpecialAttention)
//		attentionRole, err := GetRoleByLangCode(orgId, consts.RoleGroupProMember)
//		if err != nil {
//			log.Error(err)
//			return true, errs.BuildSystemErrorInfo(errs.CacheProxyError)
//		}
//		*roleIds = append(*roleIds, attentionRole.Id)
//	}
//	return false, nil
//}
