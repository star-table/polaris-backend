package domain

import (
	tablePb "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

const ProjectNumLimit = 6
const IterationNumLimit = 2
const TaskNumLimit = 1000

//func AuthIssueWithIssueId(orgId, userId int64, issueId int64, path string, operation string) errs.SystemErrorInfo {
//	issueBo, err1 := GetIssueBo(orgId, issueId)
//	if err1 != nil {
//		log.Error(err1)
//		return err1
//	}
//	return AuthIssue(orgId, userId, issueBo, path, operation)
//}

//验证任务操作权，已归档的项目下的任务不允许编辑
func AuthIssue(orgId, userId int64, issueBo *bo.IssueBo, path string, operation string, authFields ...string) errs.SystemErrorInfo {
	issueAuthBo, err := GetIssueAuthBo(issueBo, userId)
	if err != nil {
		log.Error(err)
		return err
	}
	projectAuthBo := &bo.ProjectAuthBo{
		Id:                     0,
		AppId:                  0,
		Name:                   consts.CardDefaultIssueProjectName,
		Creator:                0,
		Owner:                  0,
		Status:                 1,
		PublicStatus:           consts.PublicProject,
		IsFilling:              consts.AppIsNotFilling,
		Participants:           make([]int64, 0),
		ParticipantDepartments: make([]int64, 0),
		Followers:              make([]int64, 0),
		ProjectType:            consts.ProjectTypeNormalId,
	}
	projectId := issueAuthBo.ProjectId
	if projectId > 0 {
		projectAuthBo, err = LoadProjectAuthBo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	// 暂时不鉴权私有项目，因为协作人的存在，会让私有项目不再“私有”
	//校验私有项目
	//if projectAuthBo.PublicStatus == consts.PrivateProject {
	//	authPrivateProjectErr := AuthPrivateProject(orgId, userId, projectAuthBo)
	//	if authPrivateProjectErr != nil {
	//		//判断当前用户是不是超管
	//		adminFlagBo, err := orgfacade.GetUserAdminFlagRelaxed(orgId, userId)
	//		if err != nil {
	//			log.Error(err)
	//			return authPrivateProjectErr
	//		}
	//
	//		if !adminFlagBo.IsAdmin {
	//			log.Error(authPrivateProjectErr)
	//			return authPrivateProjectErr
	//		} else {
	//			return nil
	//		}
	//	}
	//}
	//校验项目是否归档
	if projectAuthBo.IsFilling == consts.AppIsFilling && operation != consts.RoleOperationView {
		return errs.BuildSystemErrorInfo(errs.ProjectIsArchivedWhenModifyIssue)
	}
	authErr := orgfacade.AuthenticateRelaxed(orgId, userId, projectAuthBo, issueAuthBo, path, operation, authFields)
	if authErr != nil && authErr.Code() == errs.NoOperationPermissions.Code() {
		return errs.NoOperationPermissionForProject
	}

	return authErr
}

//验证任务（有传入appId，可能是镜像id的权限）
func AuthIssueWithAppId(orgId, userId int64, issueBo *bo.IssueBo, path string, operation string, inputAppId int64, authFields ...string) errs.SystemErrorInfo {
	issueAuthBo, err := GetIssueAuthBo(issueBo, userId)
	if err != nil {
		log.Error(err)
		return err
	}
	projectAuthBo := &bo.ProjectAuthBo{
		Id:                     0,
		AppId:                  0,
		Name:                   consts.CardDefaultIssueProjectName,
		Creator:                0,
		Owner:                  0,
		Status:                 1,
		PublicStatus:           consts.PublicProject,
		IsFilling:              consts.AppIsNotFilling,
		Participants:           make([]int64, 0),
		ParticipantDepartments: make([]int64, 0),
		Followers:              make([]int64, 0),
		ProjectType:            consts.ProjectTypeNormalId,
	}
	projectId := issueAuthBo.ProjectId
	if projectId > 0 {
		projectAuthBo, err = LoadProjectAuthBo(orgId, projectId)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	if inputAppId > 0 && projectAuthBo.AppId == 0 {
		projectAuthBo.AppId = inputAppId
	}
	// 暂时不鉴权私有项目，因为协作人的存在，会让私有项目不再“私有”
	//校验私有项目
	//if projectAuthBo.PublicStatus == consts.PrivateProject {
	//	authPrivateProjectErr := AuthPrivateProject(orgId, userId, projectAuthBo)
	//	if authPrivateProjectErr != nil {
	//		//判断当前用户是不是超管
	//		adminFlagBo, err := orgfacade.GetUserAdminFlagRelaxed(orgId, userId)
	//		if err != nil {
	//			log.Error(err)
	//			return authPrivateProjectErr
	//		}
	//		if !adminFlagBo.IsAdmin {
	//			log.Error(authPrivateProjectErr)
	//			return authPrivateProjectErr
	//		} else {
	//			return nil
	//		}
	//	}
	//}
	//校验项目是否归档
	if projectAuthBo.IsFilling == consts.AppIsFilling && operation != consts.RoleOperationView {
		return errs.BuildSystemErrorInfo(errs.ProjectIsArchivedWhenModifyIssue)
	}
	authErr := orgfacade.AuthenticateRelaxed(orgId, userId, projectAuthBo, issueAuthBo, path, operation, authFields)
	if authErr != nil && authErr.Code() == errs.NoOperationPermissions.Code() {
		return errs.NoOperationPermissionForProject
	}

	return authErr
}

//同时会校验项目是否已归档
func AuthProject(orgId, userId, projectId int64, path string, operation string, authFields ...string) errs.SystemErrorInfo {
	return AuthProjectWithCond(orgId, userId, projectId, path, operation, true, false, 0, authFields...)
}

//同时会校验项目是否已归档（有传入appId，可能是镜像id的权限）
func AuthProjectWithAppId(orgId, userId, projectId int64, path string, operation string, inputAppId int64, authFields ...string) errs.SystemErrorInfo {
	return AuthProjectWithCond(orgId, userId, projectId, path, operation, true, false, 0, authFields...)
}

//项目权限校验，不需要检测是否已归档
func AuthProjectWithOutArchivedCheck(orgId, userId, projectId int64, path string, operation string) errs.SystemErrorInfo {
	return AuthProjectWithCond(orgId, userId, projectId, path, operation, false, false, 0)
}

//校验项目,跳过角色权限
func AuthProjectWithOutPermission(orgId, userId, projectId int64, path string, operation string) errs.SystemErrorInfo {
	return AuthProjectWithCond(orgId, userId, projectId, path, operation, false, true, 0)
}

func AuthProjectWithCond(orgId, userId, projectId int64, path string, operation string, authFiling bool, skipAuthPermission bool, inputAppId int64, authFields ...string) errs.SystemErrorInfo {
	projectAuthBo, err := LoadProjectAuthBo(orgId, projectId)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}

	// 暂时不鉴权私有项目，因为协作人的存在，会让私有项目不再“私有”
	//校验私有项目
	//if projectAuthBo.PublicStatus == consts.PrivateProject {
	//	authPrivateProjectErr := AuthPrivateProject(orgId, userId, projectAuthBo)
	//	if authPrivateProjectErr != nil {
	//		manageAuthInfoResp := userfacade.GetUserAuthority(orgId, userId)
	//		if manageAuthInfoResp.Failure() {
	//			log.Error(manageAuthInfoResp.Message)
	//			return manageAuthInfoResp.Error()
	//		}
	//		adminFlagBo := manageAuthInfoResp.NewData
	//
	//		if !adminFlagBo.IsSysAdmin {
	//			log.Error(authPrivateProjectErr)
	//			return authPrivateProjectErr
	//		} else {
	//			return nil
	//		}
	//	}
	//}

	//校验项目是否归档
	if authFiling && projectAuthBo.IsFilling == consts.AppIsFilling && operation != consts.RoleOperationView {
		return errs.BuildSystemErrorInfo(errs.ProjectIsArchivedWhenModifyIssue)
	}

	if !skipAuthPermission {
		if inputAppId != 0 {
			projectAuthBo.AppId = inputAppId
		}
		authErr := orgfacade.AuthenticateRelaxed(orgId, userId, projectAuthBo, nil, path, operation, authFields)
		if authErr != nil {
			log.Error(authErr)
			if authErr.Code() == errs.NoOperationPermissions.Code() {
				return errs.NoOperationPermissionForProject
			} else {
				return authErr
			}
		}
	}
	return nil
}

// func AuthProjectWithCondByAppId(orgId, userId, proAppId int64, operation string, authFields ...string,
// ) errs.SystemErrorInfo {
// 	projectBo, err := GetProjectByAppId(orgId, proAppId)
// 	if err != nil {
// 		log.Errorf("[AuthProjectWithCondByAppId] GetProjectByAppId err: %v", err)
// 		return err
// 	}
// 	err = AuthProjectWithCond(orgId, userId, projectBo.Id, "", operation, false, false,
// 		0, authFields...)
// 	if err != nil {
// 		log.Errorf("[AuthProjectWithCondByAppId] AuthProjectWithCond err: %v", err)
// 		return err
// 	}
//
// 	return nil
// }

// 获取组织的管理员 ids，包含了组织超管、组织普通管理员
func GetAdminUserIdsOfOrg(orgId, appId int64) ([]int64, errs.SystemErrorInfo) {
	adminUserResp := userfacade.GetUsersCouldManage(orgId, appId)
	if adminUserResp.Failure() {
		log.Error(adminUserResp.Error())
		return nil, adminUserResp.Error()
	}
	adminUserIds := make([]int64, 0)
	for _, user := range adminUserResp.Data.List {
		adminUserIds = append(adminUserIds, user.Id)
	}
	return adminUserIds, nil
}

//校验私有项目，只有成员才能编辑
//func AuthPrivateProject(orgId, userId int64, projectAuthBo *bo.ProjectAuthBo) errs.SystemErrorInfo {
//	memberIdMap := map[int64]bool{}
//	memberIdMap[projectAuthBo.Owner] = true
//	if projectAuthBo.Followers != nil {
//		for _, follower := range projectAuthBo.Followers {
//			memberIdMap[follower] = true
//		}
//	}
//	if projectAuthBo.Participants != nil {
//		for _, participant := range projectAuthBo.Participants {
//			memberIdMap[participant] = true
//		}
//	}
//	// 子龙：组织超管和组织普通管理员也视为项目的成员
//	adminUserIds, err := GetAdminUserIdsOfOrg(orgId)
//	if err != nil {
//		log.Error(err)
//		return err
//	}
//	if len(adminUserIds) > 0 {
//		for _, uid := range adminUserIds {
//			memberIdMap[uid] = true
//		}
//	}
//
//	//判断是否在部门里面
//	if len(projectAuthBo.ParticipantDepartments) > 0 {
//		// `0` 表示组织内的全部成员。
//		if hasIt, _ := slice.Contain(projectAuthBo.ParticipantDepartments, int64(0)); hasIt {
//			return nil
//		}
//		deptInfo := orgfacade.GetUserDeptIds(orgvo.GetUserDeptIdsReq{
//			OrgId:  orgId,
//			UserId: userId,
//		})
//		if deptInfo.Failure() {
//			log.Error(deptInfo.Error())
//			return deptInfo.Error()
//		}
//		for _, department := range projectAuthBo.ParticipantDepartments {
//			if ok, _ := slice.Contain(deptInfo.Data.DeptIds, department); ok {
//				return nil
//			}
//		}
//	}
//
//	//最后：如果项目成员不包括当前操作人
//	if _, ok := memberIdMap[userId]; !ok {
//		return errs.BuildSystemErrorInfo(errs.NoPrivateProjectPermissions)
//	}
//
//	return nil
//}

func AuthOrg(orgId, userId int64, path string, operation string) errs.SystemErrorInfo {
	return orgfacade.AuthenticateRelaxed(orgId, userId, nil, nil, path, operation, nil)
}

//func AuthPayFunction(orgId int64, functionCode string) errs.SystemErrorInfo {
//	functionResp := orgfacade.GetFunctionConfig(orgvo.GetOrgFunctionConfigReq{OrgId: orgId})
//	if functionResp.Failure() {
//		log.Error(functionResp.Error())
//		return functionResp.Error()
//	}
//
//	if ok, _ := slice.Contain(functionResp.Data.FunctionCodes, functionCode); ok {
//		return nil
//	} else {
//		return errs.FunctionIsLimitPayLevel
//	}
//}

// GetPayFunctionLimitResourceNum 获取功能限制的资源数量。仅适用于：项目、迭代、任务
func GetPayFunctionLimitResourceNum(orgId int64, functionCode string) (int, errs.SystemErrorInfo) {
	resp := orgfacade.GetFunctionObjArrByOrg(orgvo.GetOrgFunctionConfigReq{
		OrgId: orgId,
	})
	if resp.Failure() {
		log.Errorf("[GetPayFunctionLimitResourceNum] GetFunctionObjArrByOrg err: %v", resp.Error())
		return 0, resp.Error()
	}
	payFunctionMap := GetOrgFunctionInfoMap(resp.Data.Functions)
	// map 中不存在 function 表示不支持此项功能；存在，则验证限制数量
	if payFuncInfo, ok := payFunctionMap[functionCode]; ok {
		if payFuncInfo.HasLimit {
			return payFuncInfo.Limit[0].Num, nil
		} else {
			// 无限制
			return -1, nil
		}
	}

	return 0, nil
}

func AuthPayProjectNum(orgId int64, functionCode string) errs.SystemErrorInfo {
	limitNum, err := GetPayFunctionLimitResourceNum(orgId, functionCode)
	if err != nil {
		log.Errorf("[AuthPayProjectNum] GetPayFunctionLimitResourceNum err: %v, functionCode: %s", err, functionCode)
		return err
	}
	if limitNum == -1 {
		return nil
	}
	// 普通用户只能创建6个项目（包括归档的）
	count, oriErr := mysql.SelectCountByCond(consts.TableProject, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcTemplateFlag: consts.TemplateFalse, //排除模板项目
	})
	if oriErr != nil {
		log.Errorf("[AuthPayProjectNum] err: %v, functionCode: %s", oriErr, functionCode)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	if int(count) >= limitNum {
		return errs.CommonUserCreateProjectLimit
	}

	return nil
}

func AuthPayProjectIteration(orgId int64, functionCode string, projectId int64) errs.SystemErrorInfo {
	limitNum, err := GetPayFunctionLimitResourceNum(orgId, functionCode)
	if err != nil {
		log.Errorf("[AuthPayProjectNum] GetPayFunctionLimitResourceNum err: %v, functionCode: %s", err, functionCode)
		return err
	}
	if limitNum == -1 {
		return nil
	}
	count, oriErr := mysql.SelectCountByCond(consts.TableIteration, db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		consts.TcOrgId:     orgId,
	})
	if oriErr != nil {
		log.Errorf("[AuthPayProjectIteration] err: %v, functionCode: %s", oriErr, functionCode)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	}
	if int(count) >= limitNum {
		return errs.CommonUserCreateIterationLimit
	}

	return nil
}

func AuthPayTask(orgId int64, functionCode string, needAddNum int) errs.SystemErrorInfo {
	limitNum, err := GetPayFunctionLimitResourceNum(orgId, functionCode)
	if err != nil {
		log.Errorf("[AuthPayTask] GetPayFunctionLimitResourceNum err: %v, functionCode: %s", err, functionCode)
		return err
	}
	if limitNum == -1 {
		return nil
	}
	//count, oriErr := mysql.SelectCountByCond(consts.TableIssue, db.Cond{
	//	consts.TcOrgId:     orgId,
	//	consts.TcProjectId: db.NotIn(db.Raw("select id from ppm_pro_project where org_id = ? and is_delete = 2 and template_flag = 1", orgId)),
	//})
	//if oriErr != nil {
	//	log.Errorf("[AuthPayTask] err: %v, functionCode: %s", oriErr, functionCode)
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, oriErr)
	//}

	projectPos := []*po.PpmProProject{}
	errSys := mysql.SelectAllByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcTemplateFlag: consts.TemplateTrue,
	}, &projectPos)
	if errSys != nil {
		log.Errorf("[AuthPayTask] SelectCountByCond err: %v, orgId:%v, functionCode:%v", err, orgId, functionCode)
		return err
	}
	projectIds := make([]int64, 0, len(projectPos))
	for _, pro := range projectPos {
		projectIds = append(projectIds, pro.Id)
	}
	var condition *tablePb.Condition
	if len(projectIds) > 0 {
		condition = &tablePb.Condition{}
		condition.Type = tablePb.ConditionType_and
		condition.Conditions = append(condition.Conditions, GetRowsCondition(consts.BasicFieldProjectId, tablePb.ConditionType_not_in, nil, projectIds))
	}
	count, err := getRowsCountByCondition(orgId, 0, condition)
	if err != nil {
		log.Errorf("[AuthPayTask] getRowsCountByCondition err: %v, orgId:%v, functionCode:%v", err, orgId, functionCode)
		return err
	}

	// 加上被删除的任务
	deletedIssueNum, errCache := GetDeletedIssueNum(orgId)
	if errCache != nil {
		log.Errorf("[AuthPayTask] GetDeletedIssueNum err:%v, orgId:%v", errCache, orgId)
	}

	if int(count)+needAddNum+int(deletedIssueNum) > limitNum {
		return errs.CommonUserCreateTaskLimit
	}

	return nil
}

func GetOrgFunctionInfoMap(functions []orgvo.FunctionLimitObj) map[string]orgvo.FunctionLimitObj {
	funcMap := make(map[string]orgvo.FunctionLimitObj, len(functions))
	for _, function := range functions {
		funcMap[function.Key] = function
	}

	return funcMap
}

func GetNormalAdminManageApps(orgId int64) ([]*orgvo.NormalAdminAppIds, errs.SystemErrorInfo) {
	resp := userfacade.GetNormalAdminManageApps(orgId)
	if resp.Failure() {
		log.Errorf("[GetNormalAdminManageApps] err:%v, orgId:%v", resp.Error(), orgId)
		return nil, resp.Error()
	}
	da := []*orgvo.NormalAdminAppIds{}
	copyer.Copy(resp.Data, &da)
	return da, nil
}
