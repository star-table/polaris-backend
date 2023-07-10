package service

import (
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	lang2 "github.com/star-table/polaris-backend/common/core/util/lang"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/userfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

func RemoveProjectMember(orgId, userId int64, sourceChannel string, input vo.RemoveProjectMemberReq) (*vo.Void, errs.SystemErrorInfo) {
	//校验当前用户是否具有修改删除成员的权限
	authErr := AuthProjectPermission(orgId, userId, input.ProjectID, consts.RoleOperationPathOrgProMember, consts.OperationProMemberUnbind, false)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}
	err := domain.RemoveProjectMember(orgId, userId, input, sourceChannel)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &vo.Void{ID: 0}, nil
}

func AddProjectMember(orgId, userId int64, sourceChannel string, input vo.RemoveProjectMemberReq) (*vo.Void, errs.SystemErrorInfo) {
	//校验当前用户是否具有修改删除成员的权限
	authErr := AuthProjectPermission(orgId, userId, input.ProjectID, consts.RoleOperationPathOrgProMember, consts.OperationProMemberBind, false)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	err := domain.AddProjectMember(orgId, userId, input, sourceChannel)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &vo.Void{ID: 0}, nil
}

func ProjectUserListForFuse(orgId int64, page int, size int, input vo.ProjectUserListReq) (*vo.ProjectUserListResp, errs.SystemErrorInfo) {
	count, bos, err := domain.GetProjectAllMember(orgId, input.ProjectID, page, size)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	res := &vo.ProjectUserListResp{Total: count}

	if len(bos) == 0 {
		return res, nil
	}

	var allRelatedUser, allRelatedDepartment []int64
	for _, bo := range bos {
		if bo.RelationType == consts.IssueRelationTypeDepartmentParticipant {
			allRelatedDepartment = append(allRelatedDepartment, bo.RelationId)
		} else {
			allRelatedUser = append(allRelatedUser, bo.RelationId)
		}
		allRelatedUser = append(allRelatedUser, bo.Creator)
	}
	allRelatedUser = slice.SliceUniqueInt64(allRelatedUser)
	allRelatedDepartment = slice.SliceUniqueInt64(allRelatedDepartment)

	//获取所有相关人员信息
	userInfo := orgfacade.BatchGetUserDetailInfo(orgvo.BatchGetUserInfoReq{UserIds: allRelatedUser})
	if userInfo.Failure() {
		log.Error(userInfo.Error())
		return nil, userInfo.Error()
	}
	userInfoVos := &[]vo.PersonalInfo{}
	copyErr := copyer.Copy(userInfo.Data, userInfoVos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	userInfoMap := maps.NewMap("ID", *userInfoVos)

	departmentsInfo := orgfacade.Departments(orgvo.DepartmentsReqVo{
		Page: nil,
		Size: nil,
		Params: &vo.DepartmentListReq{
			DepartmentIds: allRelatedDepartment,
		},
		OrgId: orgId,
	})
	if departmentsInfo.Failure() {
		log.Error(departmentsInfo.Error())
		return nil, departmentsInfo.Error()
	}
	departmentInfoMap := map[int64]vo.DepartmentSimpleInfo{}
	departmentInfoMap[0] = vo.DepartmentSimpleInfo{
		ID:        0,
		Name:      "全部",
		UserCount: 0, //暂时先不管
	}
	if departmentsInfo.DepartmentList != nil {
		for _, department := range departmentsInfo.DepartmentList.List {
			departmentInfoMap[department.ID] = vo.DepartmentSimpleInfo{
				ID:        department.ID,
				Name:      department.Name,
				UserCount: 0, //暂时先不管
			}
		}
	}

	//特殊角色（项目成员和负责人）
	//roleGroup := orgfacade.GetProjectRoleList(rolevo.GetProjectRoleListReqVo{
	//	OrgId:     orgId,
	//	ProjectId: input.ProjectID,
	//})
	//if roleGroup.Failure() {
	//	log.Error(roleGroup.Error())
	//	return nil, roleGroup.Error()
	//}

	//memberRoleInfo := vo.Role{}
	//ownerRoleInfo := vo.Role{}
	//for _, v := range roleGroup.NewData {
	//	if v.LangCode == consts.RoleGroupProMember {
	//		memberRoleInfo = *v
	//	} else if v.LangCode == consts.RoleGroupSpecialOwner {
	//		ownerRoleInfo = *v
	//	}
	//}

	//查询成员角色
	//roleUserResp := orgfacade.GetOrgRoleUser(rolevo.GetOrgRoleUserReqVo{
	//	OrgId:     orgId,
	//	ProjectId: input.ProjectID,
	//})
	//if roleUserResp.Failure() {
	//	log.Error(roleUserResp.Error())
	//	return nil, roleUserResp.Error()
	//}
	// 极星融合 todo
	//查询部门角色
	//roleDepartmentResp := rolefacade.GetOrgRoleDepartment(rolevo.GetOrgRoleDepartmentReqVo{
	//	OrgId:     orgId,
	//	ProjectId: input.ProjectID,
	//})
	//if roleDepartmentResp.Failure() {
	//	log.Error(roleDepartmentResp.Error())
	//	return nil, roleDepartmentResp.Error()
	//}
	//userRoleMap := map[int64]rolevo.RoleUser{}
	//for _, datum := range roleUserResp.NewData {
	//	if _, ok := userRoleMap[datum.UserId]; !ok {
	//		userRoleMap[datum.UserId] = datum
	//	}
	//}
	// departmentRoleMap := map[int64]rolevo.RoleDepartment{}
	//for _, datum := range roleDepartmentResp.NewData {
	//	if _, ok := departmentRoleMap[datum.DepartmentId]; !ok {
	//		departmentRoleMap[datum.DepartmentId] = datum
	//	}
	//}
	lang := lang2.GetLang()
	isOtherLang := lang2.IsEnglish()
	languageDataMap := make(map[string]string, 0)
	if tmpMap, ok1 := consts.LANG_PRO_ROLE_NAME_MAP[lang]; ok1 {
		languageDataMap = tmpMap
	}
	for _, v := range bos {
		tempInfo := &vo.ProjectUser{}
		tempInfo.Creator = v.Creator
		tempInfo.CreateTime = types.Time(v.CreateTime)
		if _, ok := userInfoMap[v.Creator]; ok {
			user := userInfoMap[v.Creator].(vo.PersonalInfo)
			tempInfo.CreatorInfo = &user
		}
		if v.RelationType == consts.IssueRelationTypeDepartmentParticipant {
			if info, ok := departmentInfoMap[v.RelationId]; ok {
				tempInfo.DepartmentInfo = &info
			}
			tempInfo.Type = 2
		} else {
			if _, ok := userInfoMap[v.RelationId]; ok {
				user := userInfoMap[v.RelationId].(vo.PersonalInfo)
				tempInfo.UserInfo = &user
			}
			tempInfo.Type = 1
		}

		if v.RelationType == consts.IssueRelationTypeOwner {
			//负责人
			//tempInfo.UserRole = &vo.UserRoleInfo{
			//	ID:       ownerRoleInfo.ID,
			//	Name:     ownerRoleInfo.Name,
			//	LangCode: ownerRoleInfo.LangCode,
			//}
			tempInfo.UserRole = &vo.UserRoleInfo{
				ID:       1,
				Name:     "负责人",
				LangCode: "Project.Owner",
			}
		} else if v.RelationType == consts.IssueRelationTypeParticipant {
			//if _, ok := userRoleMap[v.RelationId]; ok {
			//	//特殊角色
			//	role := userRoleMap[v.RelationId]
			//	tempInfo.UserRole = &vo.UserRoleInfo{
			//		ID:       role.RoleId,
			//		Name:     role.RoleName,
			//		LangCode: role.RoleLangCode,
			//	}
			//} else {
			//	//项目成员
			//	tempInfo.UserRole = &vo.UserRoleInfo{
			//		ID:       memberRoleInfo.ID,
			//		Name:     memberRoleInfo.Name,
			//		LangCode: memberRoleInfo.LangCode,
			//	}
			//}
			tempInfo.UserRole = &vo.UserRoleInfo{
				ID:       -1,
				Name:     "成员",
				LangCode: "Project.Member",
			}
		} else if v.RelationType == consts.IssueRelationTypeDepartmentParticipant {
			//if _, ok := departmentRoleMap[v.RelationId]; ok {
			//	//特殊角色
			//	role := departmentRoleMap[v.RelationId]
			//	tempInfo.UserRole = &vo.UserRoleInfo{
			//		ID:       role.RoleId,
			//		Name:     role.RoleName,
			//		LangCode: role.RoleLangCode,
			//	}
			//} else {
			//	//项目成员
			//	tempInfo.UserRole = &vo.UserRoleInfo{
			//		ID:       memberRoleInfo.ID,
			//		Name:     memberRoleInfo.Name,
			//		LangCode: memberRoleInfo.LangCode,
			//	}
			//}
			tempInfo.UserRole = &vo.UserRoleInfo{
				ID:       -2,
				Name:     "其他",
				LangCode: "Project.IssueRelationTypeDepartmentParticipant",
			}
		}
		if isOtherLang {
			if tmpVal, ok2 := languageDataMap[tempInfo.UserRole.Name]; ok2 {
				tempInfo.UserRole.Name = tmpVal
			}
		}
		res.List = append(res.List, tempInfo)
	}

	return res, nil
}

func ProjectUserList(orgId int64, page int, size int, input vo.ProjectUserListReq) (*vo.ProjectUserListResp, errs.SystemErrorInfo) {
	count, bos, err := domain.GetProjectAllMember(orgId, input.ProjectID, page, size)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	res := &vo.ProjectUserListResp{Total: count}

	if len(bos) == 0 {
		return res, nil
	}

	var allRelatedUser, allRelatedDepartment []int64
	for _, bo := range bos {
		if bo.RelationType == consts.IssueRelationTypeDepartmentParticipant {
			allRelatedDepartment = append(allRelatedDepartment, bo.RelationId)
		} else {
			allRelatedUser = append(allRelatedUser, bo.RelationId)
		}
		allRelatedUser = append(allRelatedUser, bo.Creator)
	}
	allRelatedUser = slice.SliceUniqueInt64(allRelatedUser)
	allRelatedDepartment = slice.SliceUniqueInt64(allRelatedDepartment)

	//获取所有相关人员信息
	userInfo := orgfacade.BatchGetUserDetailInfo(orgvo.BatchGetUserInfoReq{UserIds: allRelatedUser})
	if userInfo.Failure() {
		log.Error(userInfo.Error())
		return nil, userInfo.Error()
	}
	userInfoVos := &[]vo.PersonalInfo{}
	copyErr := copyer.Copy(userInfo.Data, userInfoVos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	userInfoMap := maps.NewMap("ID", *userInfoVos)

	departmentsInfo := orgfacade.Departments(orgvo.DepartmentsReqVo{
		Page: nil,
		Size: nil,
		Params: &vo.DepartmentListReq{
			DepartmentIds: allRelatedDepartment,
		},
		OrgId: orgId,
	})
	if departmentsInfo.Failure() {
		log.Error(departmentsInfo.Error())
		return nil, departmentsInfo.Error()
	}
	departmentInfoMap := map[int64]vo.DepartmentSimpleInfo{}
	departmentInfoMap[0] = vo.DepartmentSimpleInfo{
		ID:        0,
		Name:      "全部",
		UserCount: 0, //暂时先不管
	}
	if departmentsInfo.DepartmentList != nil {
		for _, department := range departmentsInfo.DepartmentList.List {
			departmentInfoMap[department.ID] = vo.DepartmentSimpleInfo{
				ID:        department.ID,
				Name:      department.Name,
				UserCount: 0, //暂时先不管
			}
		}
	}

	//特殊角色（项目成员和负责人）
	//roleGroup := orgfacade.GetProjectRoleList(rolevo.GetProjectRoleListReqVo{
	//	OrgId:     orgId,
	//	ProjectId: input.ProjectID,
	//})
	//if roleGroup.Failure() {
	//	log.Error(roleGroup.Error())
	//	return nil, roleGroup.Error()
	//}

	//memberRoleInfo := vo.Role{}
	//ownerRoleInfo := vo.Role{}
	//for _, v := range roleGroup.NewData {
	//	if v.LangCode == consts.RoleGroupProMember {
	//		memberRoleInfo = *v
	//	} else if v.LangCode == consts.RoleGroupSpecialOwner {
	//		ownerRoleInfo = *v
	//	}
	//}

	//查询成员角色
	//roleUserResp := orgfacade.GetOrgRoleUser(rolevo.GetOrgRoleUserReqVo{
	//	OrgId:     orgId,
	//	ProjectId: input.ProjectID,
	//})
	//if roleUserResp.Failure() {
	//	log.Error(roleUserResp.Error())
	//	return nil, roleUserResp.Error()
	//}
	// 极星融合 todo
	//查询部门角色
	//roleDepartmentResp := rolefacade.GetOrgRoleDepartment(rolevo.GetOrgRoleDepartmentReqVo{
	//	OrgId:     orgId,
	//	ProjectId: input.ProjectID,
	//})
	//if roleDepartmentResp.Failure() {
	//	log.Error(roleDepartmentResp.Error())
	//	return nil, roleDepartmentResp.Error()
	//}
	//userRoleMap := map[int64]rolevo.RoleUser{}
	//for _, datum := range roleUserResp.NewData {
	//	if _, ok := userRoleMap[datum.UserId]; !ok {
	//		userRoleMap[datum.UserId] = datum
	//	}
	//}
	// departmentRoleMap := map[int64]rolevo.RoleDepartment{}
	//for _, datum := range roleDepartmentResp.NewData {
	//	if _, ok := departmentRoleMap[datum.DepartmentId]; !ok {
	//		departmentRoleMap[datum.DepartmentId] = datum
	//	}
	//}
	lang := lang2.GetLang()
	isOtherLang := lang2.IsEnglish()
	languageDataMap := make(map[string]string, 0)
	if tmpMap, ok1 := consts.LANG_PRO_ROLE_NAME_MAP[lang]; ok1 {
		languageDataMap = tmpMap
	}
	for _, v := range bos {
		tempInfo := &vo.ProjectUser{}
		tempInfo.Creator = v.Creator
		tempInfo.CreateTime = types.Time(v.CreateTime)
		if _, ok := userInfoMap[v.Creator]; ok {
			user := userInfoMap[v.Creator].(vo.PersonalInfo)
			tempInfo.CreatorInfo = &user
		}
		if v.RelationType == consts.IssueRelationTypeDepartmentParticipant {
			if info, ok := departmentInfoMap[v.RelationId]; ok {
				tempInfo.DepartmentInfo = &info
			}
			tempInfo.Type = 2
		} else {
			if _, ok := userInfoMap[v.RelationId]; ok {
				user := userInfoMap[v.RelationId].(vo.PersonalInfo)
				tempInfo.UserInfo = &user
			}
			tempInfo.Type = 1
		}

		if v.RelationType == consts.IssueRelationTypeOwner {
			//负责人
			//tempInfo.UserRole = &vo.UserRoleInfo{
			//	ID:       ownerRoleInfo.ID,
			//	Name:     ownerRoleInfo.Name,
			//	LangCode: ownerRoleInfo.LangCode,
			//}
			tempInfo.UserRole = &vo.UserRoleInfo{
				ID:       1,
				Name:     "负责人",
				LangCode: "Project.Owner",
			}
		} else if v.RelationType == consts.IssueRelationTypeParticipant {
			//if _, ok := userRoleMap[v.RelationId]; ok {
			//	//特殊角色
			//	role := userRoleMap[v.RelationId]
			//	tempInfo.UserRole = &vo.UserRoleInfo{
			//		ID:       role.RoleId,
			//		Name:     role.RoleName,
			//		LangCode: role.RoleLangCode,
			//	}
			//} else {
			//	//项目成员
			//	tempInfo.UserRole = &vo.UserRoleInfo{
			//		ID:       memberRoleInfo.ID,
			//		Name:     memberRoleInfo.Name,
			//		LangCode: memberRoleInfo.LangCode,
			//	}
			//}
			tempInfo.UserRole = &vo.UserRoleInfo{
				ID:       -1,
				Name:     "成员",
				LangCode: "Project.Member",
			}
		} else if v.RelationType == consts.IssueRelationTypeDepartmentParticipant {
			//if _, ok := departmentRoleMap[v.RelationId]; ok {
			//	//特殊角色
			//	role := departmentRoleMap[v.RelationId]
			//	tempInfo.UserRole = &vo.UserRoleInfo{
			//		ID:       role.RoleId,
			//		Name:     role.RoleName,
			//		LangCode: role.RoleLangCode,
			//	}
			//} else {
			//	//项目成员
			//	tempInfo.UserRole = &vo.UserRoleInfo{
			//		ID:       memberRoleInfo.ID,
			//		Name:     memberRoleInfo.Name,
			//		LangCode: memberRoleInfo.LangCode,
			//	}
			//}
			tempInfo.UserRole = &vo.UserRoleInfo{
				ID:       -2,
				Name:     "其他",
				LangCode: "Project.IssueRelationTypeDepartmentParticipant",
			}
		}
		if isOtherLang {
			if tmpVal, ok2 := languageDataMap[tempInfo.UserRole.Name]; ok2 {
				tempInfo.UserRole.Name = tmpVal
			}
		}
		res.List = append(res.List, tempInfo)
	}

	return res, nil
}

func OrgProjectMemberList(orgId int64, sourceChannel string, page, size *int, input vo.OrgProjectMemberListReq) (*vo.OrgProjectMemberListResp, errs.SystemErrorInfo) {
	distinctUserIds, err := GetProjectRelationUserIds(input.ProjectID, input.RelationType)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if len(distinctUserIds) == 0 {
		return &vo.OrgProjectMemberListResp{
			Total: 0,
			List:  []*vo.OrgProjectMemberInfoResp{},
		}, nil
	}

	//获取组织人员信息
	ignoreDelete := false
	if input.IgnoreDelete != nil {
		ignoreDelete = *input.IgnoreDelete
	}
	resp := orgfacade.DepartmentMembersList(orgvo.DepartmentMembersListReq{
		OrgId:         orgId,
		SourceChannel: sourceChannel,
		Params: &vo.DepartmentMembersListReq{
			Name:    input.Name,
			UserIds: distinctUserIds,
		},
		Page:         page,
		Size:         size,
		IgnoreDelete: ignoreDelete,
	})

	if resp.Failure() {
		log.Error(resp.Error())
		return nil, resp.Error()
	}

	trulyUserIds := []int64{}
	for _, info := range resp.Data.List {
		trulyUserIds = append(trulyUserIds, info.UserID)
	}

	userInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, trulyUserIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	allMembers := &[]*vo.OrgProjectMemberInfoResp{}
	_ = copyer.Copy(userInfos, allMembers)

	return &vo.OrgProjectMemberListResp{
		Total: resp.Data.Total,
		List:  *allMembers,
	}, nil
}

func GetProjectRelationUserIds(projectId int64, relationType *int64) ([]int64, errs.SystemErrorInfo) {
	typeArr := []int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant}
	if relationType != nil {
		if *relationType == 1 {
			//负责人
			typeArr = []int64{consts.IssueRelationTypeOwner}
		} else if *relationType == 2 {
			//参与人
			typeArr = []int64{consts.IssueRelationTypeParticipant}
		}
	}
	relationBo, err := domain.GetProjectRelationByType(projectId, typeArr)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	//去重所有的用户id
	distinctUserIds, _ := DistinctUserIds(relationBo)

	return distinctUserIds, nil
}

func ProjectMemberIdList(orgId int64, projectId int64, input *projectvo.ProjectMemberIdListReqData) (*vo.ProjectMemberIDListResp, errs.SystemErrorInfo) {
	if input == nil {
		input = &projectvo.ProjectMemberIdListReqData{}
	}
	relations, err := domain.GetProjectMembers(orgId, projectId, []int64{consts.IssueRelationTypeDepartmentParticipant, consts.IssueRelationTypeParticipant, consts.IssueRelationTypeOwner})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	appId := int64(-1)
	if projectId > 0 {
		projectBo, err := domain.GetProjectSimple(orgId, projectId)
		if err != nil {
			log.Errorf("[GetProjectMemberIds]GetProjectSimple err:%v, orgId:%v, projectId:%v", err, orgId, projectId)
			return nil, err
		}
		appId = projectBo.AppId
	}
	adminUserIds, err := GetAdminOfProject(orgId, appId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var deptIds, userIds []int64
	for _, relation := range relations {
		if relation.RelationType == consts.IssueRelationTypeDepartmentParticipant {
			deptIds = append(deptIds, relation.RelationId)
		} else {
			userIds = append(userIds, relation.RelationId)
		}
	}

	if input.IncludeAdmin == 1 {
		// 追加上管理员的 id
		userIds = append(userIds, adminUserIds...)
	}

	return &vo.ProjectMemberIDListResp{
		DepartmentIds: slice.SliceUniqueInt64(deptIds),
		UserIds:       slice.SliceUniqueInt64(userIds),
	}, nil
}

// 获取项目的管理员列表
func GetAdminOfProject(orgId, appId int64) ([]int64, errs.SystemErrorInfo) {
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

// GetProjectMemberIds 获取项目成员ids列表，rest api调用
func GetProjectMemberIds(orgId, projectId int64, includeAdmin int) (*projectvo.GetProjectMemberIdsData, errs.SystemErrorInfo) {
	relations, err := domain.GetProjectMembers(orgId, projectId, []int64{consts.IssueRelationTypeDepartmentParticipant, consts.IssueRelationTypeParticipant, consts.IssueRelationTypeOwner})
	if err != nil {
		log.Errorf("[GetProjectMemberIds] GetProjectMembers err:%v, orgId:%v, projectId:%v", err, orgId, projectId)
		return nil, err
	}
	var deptIds, userIds []int64
	for _, relation := range relations {
		if relation.RelationType == consts.IssueRelationTypeDepartmentParticipant {
			deptIds = append(deptIds, relation.RelationId)
		} else {
			userIds = append(userIds, relation.RelationId)
		}
	}

	appId := int64(-1)
	if projectId > 0 {
		projectBo, err := domain.GetProjectSimple(orgId, projectId)
		if err != nil {
			log.Errorf("[GetProjectMemberIds]GetProjectSimple err:%v, orgId:%v, projectId:%v", err, orgId, projectId)
			return nil, err
		}
		appId = projectBo.AppId
	}

	if includeAdmin == consts.IncludeAdmin {
		// 包含管理员
		adminUserIds, err := GetAdminOfProject(orgId, appId)
		if err != nil {
			log.Errorf("[GetProjectMemberIds] GetAdminOfProject err:%v, orgId:%v, projectId:%v", err, orgId, projectId)
			return nil, err
		}
		userIds = append(userIds, adminUserIds...)
	}

	return &projectvo.GetProjectMemberIdsData{
		DepartmentIds: slice.SliceUniqueInt64(deptIds),
		UserIds:       slice.SliceUniqueInt64(userIds),
	}, nil
}

// 超管-选择所有项目或者某个项目都显示所有成员
// 项目管理员-选择管理的项目显示项目所有成员
// 如果是团队成员，只可以选择当前用户
func GetTrendsMembers(orgId, userId, projectId int64) (*projectvo.GetTrendListMembersData, errs.SystemErrorInfo) {
	userIds := []int64{}
	userIds = append(userIds, userId)

	appId := int64(0)
	if projectId > 0 {
		appIdResp, err := domain.GetAppIdFromProjectId(orgId, projectId)
		if err != nil {
			log.Errorf("[GetTrendsMembers] domain.GetAppIdFromProjectId err:%v, orgId:%v, userId:%v, projectId:%v",
				err, orgId, userId, projectId)
			return nil, err
		}
		appId = appIdResp
	}

	projectAdminIds, err := domain.GetProjectAdminIds(orgId, appId)
	if err != nil {
		log.Errorf("[GetTrendsMembers] domain.GetProjectAdminIds err:%v, orgId:%v, userId:%v, projectId:%v",
			err, orgId, userId, projectId)
		return nil, err
	}
	if ok, er := slice.Contain(projectAdminIds, userId); er == nil && ok {
		// 是管理员
		projectMemberData, err := GetProjectMemberIds(orgId, projectId, consts.IncludeAdmin)
		if err != nil {
			log.Errorf("[GetTrendsMembers] GetProjectMemberIds err:%v, orgId:%v, userId:%v, projectId:%v",
				err, orgId, userId, projectId)
			return nil, err
		}
		userIds = append(userIds, projectMemberData.UserIds...)

		if len(projectMemberData.DepartmentIds) > 0 {
			resp := orgfacade.GetUserIdsByDeptIds(&orgvo.GetUserIdsByDeptIdsReq{
				OrgId:   orgId,
				DeptIds: projectMemberData.DepartmentIds,
			})
			if resp.Failure() {
				log.Errorf("[GetTrendsMembers] orgfacade.GetUserIdsByDeptIds err:%v, orgId:%v, userId:%v, projectId:%v",
					resp.Error(), orgId, userId, projectId)
				return nil, resp.Error()
			}
			userIds = append(userIds, resp.Data.UserIds...)
		}
	}

	userIds = slice.SliceUniqueInt64(userIds)

	userInfoBatch := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   orgId,
		UserIds: userIds,
	})

	if userInfoBatch.Failure() {
		log.Errorf("[GetTrendsMembers] GetBaseUserInfoBatch err:%v, orgId:%v, userId:%v, projectId:%v",
			userInfoBatch.Error(), orgId, userId, projectId)
		return nil, userInfoBatch.Error()
	}

	users := []projectvo.UserInfo{}
	for _, u := range userInfoBatch.BaseUserInfos {
		users = append(users, projectvo.UserInfo{
			Id:   u.UserId,
			Name: u.Name,
		})
	}

	return &projectvo.GetTrendListMembersData{Users: users}, nil
}
