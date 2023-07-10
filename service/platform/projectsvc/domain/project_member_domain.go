package domain

import (
	"time"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/asyn"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/rolevo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"gopkg.in/fatih/set.v0"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func HandleProjectMember(orgId int64, currentUserId int64, owner int64, projectId int64, memberIds []int64, followerIds []int64, isAllMember *bool, departmentIds []int64, ownerIds []int64) ([]interface{}, []int64, errs.SystemErrorInfo) {
	//插入项目成员
	memberEntities := []interface{}{}
	addedMemberIds := []int64{}

	//1.负责人
	allOwner := []int64{}
	if owner != int64(0) {
		allOwner = append(allOwner, owner)
	}
	if ownerIds != nil && len(ownerIds) > 0 {
		allOwner = append(allOwner, ownerIds...)
	}
	allOwner = slice.SliceUniqueInt64(allOwner)
	if len(allOwner) != 0 {
		verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, allOwner)
		if !verifyOrgUserFlag {
			log.Error("存在用户组织校验失败")
			return memberEntities, nil, errs.VerifyOrgError
		}
		memberRelationIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(allOwner))
		if idErr != nil {
			log.Error(idErr)
			return memberEntities, nil, idErr
		}
		for i, v := range allOwner {
			memberEntities = append(memberEntities, po.PpmProProjectRelation{
				Id:           memberRelationIds.Ids[i].Id,
				OrgId:        orgId,
				ProjectId:    projectId,
				RelationId:   v,
				RelationType: consts.IssueRelationTypeOwner,
				Creator:      currentUserId,
				CreateTime:   time.Now(),
				IsDelete:     consts.AppIsNoDelete,
				Status:       consts.ProjectMemberEffective,
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
				Version:      1,
			})
			addedMemberIds = append(addedMemberIds, v)
		}
	}

	//2.项目成员
	if isAllMember != nil && *isAllMember == true {
		//全选，获取所有成员
		tempMemberIds := []int64{}
		resp := orgfacade.GetOrgUserIds(orgvo.GetOrgUserIdsReq{
			OrgId: orgId,
		})

		if resp.Failure() {
			log.Error(resp.Error())
			return nil, nil, resp.Error()
		}

		for _, info := range resp.Data {
			tempMemberIds = append(tempMemberIds, info)
		}
		memberIds = tempMemberIds
	} else {
		//默认创建者也是项目成员
		if owner != currentUserId {
			if bool, _ := slice.Contain(memberIds, currentUserId); !bool {
				memberIds = append(memberIds, currentUserId)
			}
		}
	}
	memberIds = slice.SliceUniqueInt64(memberIds)

	if len(memberIds) != 0 {
		verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, memberIds)
		if !verifyOrgUserFlag {
			log.Error("存在用户组织校验失败")
			return memberEntities, nil, errs.VerifyOrgError
		}
		memberRelationIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(memberIds))
		if idErr != nil {
			log.Error(idErr)
			return memberEntities, nil, idErr
		}
		for i, v := range memberIds {
			memberEntities = append(memberEntities, po.PpmProProjectRelation{
				Id:           memberRelationIds.Ids[i].Id,
				OrgId:        orgId,
				ProjectId:    projectId,
				RelationId:   v,
				RelationType: consts.IssueRelationTypeParticipant,
				Creator:      currentUserId,
				CreateTime:   time.Now(),
				IsDelete:     consts.AppIsNoDelete,
				Status:       consts.ProjectMemberEffective,
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
				Version:      1,
			})
			addedMemberIds = append(addedMemberIds, v)
		}
	}

	//成员部门处理
	if len(departmentIds) > 0 {
		departmentIds = slice.SliceUniqueInt64(departmentIds)
		if ok, _ := slice.Contain(departmentIds, int64(0)); ok {
			departmentIds = []int64{0}
		} else {
			verifyDepartment := orgfacade.VerifyDepartments(orgvo.VerifyDepartmentsReq{DepartmentIds: departmentIds, OrgId: orgId})
			if !verifyDepartment.IsTrue {
				log.Errorf("存在无效部门, 组织id:【%d】,部门：【%s】", orgId, json.ToJsonIgnoreError(departmentIds))
				return memberEntities, nil, errs.DepartmentNotExist
			}
		}

		memberRelationIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(departmentIds))
		if idErr != nil {
			log.Error(idErr)
			return memberEntities, nil, idErr
		}
		for i, v := range departmentIds {
			memberEntities = append(memberEntities, po.PpmProProjectRelation{
				Id:           memberRelationIds.Ids[i].Id,
				OrgId:        orgId,
				ProjectId:    projectId,
				RelationId:   v,
				RelationType: consts.IssueRelationTypeDepartmentParticipant,
				Creator:      currentUserId,
				Updator:      currentUserId,
				UpdateTime:   time.Now(),
			})
		}

		//todo 成员（为了日历和群聊服务）
	}

	return memberEntities, addedMemberIds, nil
}

//我参与的
func GetParticipantMembers(orgId, currentUserId int64) ([]int64, errs.SystemErrorInfo) {
	projectIdsNeed := []int64{}
	memberEntities := &[]*po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond((&po.PpmProProjectRelation{}).TableName(), db.Cond{
		consts.TcIsDelete:     db.Eq(consts.AppIsNoDelete),
		consts.TcRelationType: db.Eq(consts.IssueRelationTypeParticipant),
		consts.TcRelationId:   db.Eq(currentUserId),
		consts.TcOrgId:        orgId,
	}, memberEntities)
	if err != nil {
		return projectIdsNeed, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, v := range *memberEntities {
		projectIdsNeed = append(projectIdsNeed, v.ProjectId)
	}

	return projectIdsNeed, nil
}

//我参与的和我负责的
func GetParticipantMembersAndOwner(orgId, currentUserId int64) ([]int64, errs.SystemErrorInfo) {
	projectIdsNeed := []int64{}
	memberEntities := &[]*po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond((&po.PpmProProjectRelation{}).TableName(), db.Cond{
		consts.TcIsDelete:     db.Eq(consts.AppIsNoDelete),
		consts.TcRelationType: db.In([]int{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeOwner}),
		consts.TcRelationId:   db.Eq(currentUserId),
		consts.TcOrgId:        orgId,
	}, memberEntities)
	if err != nil {
		return projectIdsNeed, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	for _, v := range *memberEntities {
		projectIdsNeed = append(projectIdsNeed, v.ProjectId)
	}

	return projectIdsNeed, nil
}

func GetProjectMemberInfo(projectIds []int64, orgId int64, creatorIds []int64) (map[int64][]bo.UserIDInfoBo, map[int64][]bo.UserIDInfoBo, map[int64][]bo.UserIDInfoBo, map[int64]bo.UserIDInfoBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	relatedInfo := &[]bo.RelationInfoTypeBo{}
	err1 := conn.Select("relation_id", "relation_type", "project_id").From("ppm_pro_project_relation").
		Where(db.Cond{
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcProjectId:    db.In(projectIds),
			consts.TcStatus:       1,
			consts.TcOrgId:        orgId,
			consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant}),
		}).All(relatedInfo)
	if err1 != nil {
		log.Error(err1)
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	creatorInfo := map[int64]bo.UserIDInfoBo{}

	ownerInfo := map[int64][]bo.UserIDInfoBo{}
	participantInfo := map[int64][]bo.UserIDInfoBo{}
	followerInfo := map[int64][]bo.UserIDInfoBo{}
	allRelationIds := []int64{}
	for _, v := range *relatedInfo {
		allRelationIds = append(allRelationIds, v.RelationId)
	}

	creatorIds = slice.SliceUniqueInt64(creatorIds)

	allRelationIds = append(allRelationIds, creatorIds...)
	allRelationIds = slice.SliceUniqueInt64(allRelationIds)
	allUserInfo, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, allRelationIds)
	if err != nil {
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	userInfoById := map[int64]bo.BaseUserInfoBo{}
	for _, v := range allUserInfo {
		userInfoById[v.UserId] = v
	}
	for _, v := range creatorIds {
		userInfo, ok := userInfoById[v]
		if !ok {
			continue
		}
		temp := bo.UserIDInfoBo{}
		temp.Id = userInfo.UserId
		temp.Name = userInfo.Name
		temp.NamePy = userInfo.NamePy
		temp.Avatar = userInfo.Avatar
		temp.UserID = userInfo.UserId
		temp.EmplID = userInfo.OutUserId
		temp.IsDeleted = userInfo.OrgUserIsDelete == consts.AppIsDeleted
		temp.IsDisabled = userInfo.OrgUserStatus == consts.AppStatusDisabled
		creatorInfo[v] = temp
	}
	for _, v := range *relatedInfo {
		userInfo, ok := userInfoById[v.RelationId]
		if !ok {
			continue
		}
		temp := bo.UserIDInfoBo{}
		if userInfo.OrgUserIsDelete == consts.AppIsDeleted {
			continue
		}

		temp.Id = userInfo.UserId
		temp.Name = userInfo.Name
		temp.NamePy = userInfo.NamePy
		temp.Avatar = userInfo.Avatar
		temp.UserID = userInfo.UserId
		temp.EmplID = userInfo.OutUserId
		temp.IsDeleted = userInfo.OrgUserIsDelete == consts.AppIsDeleted
		temp.IsDisabled = userInfo.OrgUserStatus == consts.AppStatusDisabled
		if v.RelationType == consts.IssueRelationTypeOwner {
			ownerInfo[v.ProjectId] = append(ownerInfo[v.ProjectId], temp)
		} else if v.RelationType == consts.IssueRelationTypeParticipant {
			participantInfo[v.ProjectId] = append(participantInfo[v.ProjectId], temp)
			//} else if v.RelationType == consts.IssueRelationTypeFollower {
			//	followerInfo[v.ProjectId] = append(followerInfo[v.ProjectId], temp)
		}
	}

	//项目成员展示去重
	uniqueParticipantInfo := map[int64][]bo.UserIDInfoBo{}
	for i, bos := range participantInfo {
		relationIdsForProject := []int64{}
		for _, infoBo := range bos {
			if ok, _ := slice.Contain(relationIdsForProject, infoBo.UserID); !ok {
				uniqueParticipantInfo[i] = append(uniqueParticipantInfo[i], infoBo)
				relationIdsForProject = append(relationIdsForProject, infoBo.UserID)
			}
		}
	}
	return ownerInfo, uniqueParticipantInfo, followerInfo, creatorInfo, nil
}

func JudgeIsProjectMember(currentUserId, orgId, projectId int64) (po.PpmProProjectRelation, errs.SystemErrorInfo) {
	member := &po.PpmProProjectRelation{}
	err := mysql.SelectOneByCond(member.TableName(), db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationId:   currentUserId,
		consts.TcProjectId:    projectId,
		consts.TcOrgId:        orgId,
		consts.TcRelationType: consts.IssueRelationTypeParticipant,
	}, member)
	if err != nil {
		return *member, errs.BuildSystemErrorInfo(errs.NotProjectParticipant)
	}

	return *member, nil
}

func GetChangeMembersAndDeleteOld(tx sqlbuilder.Tx, input bo.UpdateProjectBo, orgId int64, oldOwnerIds []int64, updPoint *mysql.Upd, currentUserId int64) ([]int64, []int64, set.Interface, set.Interface, errs.SystemErrorInfo) {
	upd := *updPoint
	memberAdd := []interface{}{}

	//负责人更新(有ownerIds就取ownerIds，没有就取owner)
	delOwners := []int64{}
	addOwners := []int64{}
	afterOwners := oldOwnerIds
	if util.FieldInUpdate(input.UpdateFields, "ownerIds") && input.OwnerIds != nil && len(input.OwnerIds) > 0 {
		delOwners, addOwners = util.GetDifMemberIds(oldOwnerIds, input.OwnerIds)
		afterOwners = input.OwnerIds
	} else if util.FieldInUpdate(input.UpdateFields, "owner") && input.Owner != nil {
		afterOwners = []int64{*input.Owner}
		delOwners, addOwners = util.GetDifMemberIds(oldOwnerIds, afterOwners)
		upd[consts.TcOwner] = *input.Owner
	}
	if len(afterOwners) == 0 {
		return nil, nil, nil, nil, errs.NeedProjectOwner
	}
	projectOwner := oldOwnerIds[0]
	if len(delOwners) != 0 || len(addOwners) != 0 {
		upd[consts.TcOwner] = afterOwners[0]
		projectOwner = afterOwners[0]
		if len(addOwners) > 0 {
			verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, addOwners)
			if !verifyOrgUserFlag {
				log.Error("存在用户组织校验失败")
				return nil, nil, nil, nil, errs.VerifyOrgError
			}
			addOwners = slice.SliceUniqueInt64(addOwners)
			memberIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(addOwners))
			if idErr != nil {
				log.Error(idErr)
				return nil, nil, nil, nil, idErr
			}
			for i, val := range addOwners {
				memberAdd = append(memberAdd, po.PpmProProjectRelation{
					Id:           memberIds.Ids[i].Id,
					OrgId:        orgId,
					ProjectId:    input.ID,
					RelationId:   val,
					RelationType: consts.IssueRelationTypeOwner,
					Creator:      currentUserId,
					CreateTime:   time.Now(),
					IsDelete:     consts.AppIsNoDelete,
					Status:       consts.ProjectMemberEffective,
					Updator:      currentUserId,
					UpdateTime:   time.Now(),
					Version:      1,
				})
			}
		}
		if len(delOwners) > 0 {
			//删除旧有负责人
			err := tx.Collection((&po.PpmProProjectRelation{}).TableName()).Find(db.Cond{
				consts.TcRelationId:   db.In(delOwners),
				consts.TcOrgId:        db.Eq(orgId),
				consts.TcProjectId:    db.Eq(input.ID),
				consts.TcRelationType: consts.IssueRelationTypeOwner,
			}).Update(map[string]int{
				consts.TcIsDelete: consts.AppIsDeleted,
			})
			if err != nil {
				return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}
	}

	oldMembers := set.New(set.ThreadSafe)
	thisMembers := set.New(set.ThreadSafe)

	//成员更新
	if err := assemblyMembers(input, orgId, &oldMembers, &thisMembers, tx); err != nil {
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	//需要删除的
	delMembers := set.Difference(oldMembers, thisMembers)
	if len(delMembers.List()) != 0 {
		err := DeleteRelationByDeleteMember(tx, delMembers.List(), projectOwner, input.ID, orgId, currentUserId)
		if err != nil {
			return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
		}
	}

	addMembers := set.Difference(thisMembers, oldMembers)
	addMemberArr := []int64{}
	for _, i2 := range addMembers.List() {
		val, ok := i2.(int64)
		if !ok {
			continue
		}
		addMemberArr = append(addMemberArr, val)
	}

	verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, addMemberArr)
	if !verifyOrgUserFlag {
		log.Error("存在用户组织校验失败")
		return nil, nil, nil, nil, errs.VerifyOrgError
	}
	addMemberArr = slice.SliceUniqueInt64(addMemberArr)
	memberIds, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableProjectRelation, len(addMemberArr))
	if idErr != nil {
		log.Error(idErr)
		return nil, nil, nil, nil, idErr
	}
	for i, val := range addMemberArr {
		memberAdd = append(memberAdd, po.PpmProProjectRelation{
			Id:           memberIds.Ids[i].Id,
			OrgId:        orgId,
			ProjectId:    input.ID,
			RelationId:   val,
			RelationType: consts.IssueRelationTypeParticipant,
			Creator:      currentUserId,
			CreateTime:   time.Now(),
			IsDelete:     consts.AppIsNoDelete,
			Status:       consts.ProjectMemberEffective,
			Updator:      currentUserId,
			UpdateTime:   time.Now(),
			Version:      1,
		})
	}

	err := PaginationInsert(memberAdd, &po.PpmProProjectRelation{}, tx)
	if err != nil {
		log.Error(err)
		return nil, nil, nil, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return oldOwnerIds, afterOwners, oldMembers, thisMembers, nil
}

// 项目的负责人、关注人 id。
func GetProjectParticipantIds(orgId, projectId int64) ([]int64, error) {
	uids := make([]int64, 0)
	memberEntities := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: db.In([]int{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeOwner}),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, memberEntities)
	if err != nil {
		return uids, err
	}
	for _, item := range *memberEntities {
		uids = append(uids, item.RelationId)
	}
	return uids, nil
}

func assemblyMembers(input bo.UpdateProjectBo, orgId int64, oldMembers, thisMembers *set.Interface, tx sqlbuilder.Tx) error {

	if util.FieldInUpdate(input.UpdateFields, "memberIds") && input.MemberIds != nil {
		memberEntities := &[]po.PpmProProjectRelation{}
		err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcProjectId:    input.ID,
			consts.TcRelationType: consts.IssueRelationTypeParticipant,
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}, memberEntities)
		if err != nil {
			return err
		}
		for _, v := range *memberEntities {
			(*oldMembers).Add(v.RelationId)
		}

		input.MemberIds = slice.SliceUniqueInt64(input.MemberIds)
		for _, v := range input.MemberIds {
			(*thisMembers).Add(v)
		}

		delMembers := set.Difference(*oldMembers, *thisMembers)
		if len(delMembers.List()) != 0 {
			err := tx.Collection((&po.PpmProProjectRelation{}).TableName()).Find(db.Cond{
				consts.TcRelationId:   db.In(delMembers.List()),
				consts.TcOrgId:        db.Eq(orgId),
				consts.TcProjectId:    db.Eq(input.ID),
				consts.TcRelationType: consts.IssueRelationTypeParticipant,
			}).Update(map[string]int{
				consts.TcIsDelete: consts.AppIsDeleted,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 分组插入新数据
func PaginationInsert(list []interface{}, domainObj mysql.Domain, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	isTx := tx != nil && len(tx) > 0
	totalSize := len(list)
	batch := 1000
	offset := 0
	for {
		limit := offset + batch
		if totalSize < limit {
			limit = totalSize
		}
		oneBatch := list[offset:limit]
		if isTx {
			batchInsert := mysql.TransBatchInsert(tx[0], domainObj, oneBatch)
			if batchInsert != nil {
				log.Error(batchInsert)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
			}
		} else {
			batchInsert := mysql.BatchInsert(domainObj, oneBatch)
			if batchInsert != nil {
				log.Error(batchInsert)
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, batchInsert)
			}
		}

		if totalSize <= limit {
			break
		}
		offset += batch
	}
	return nil
}

func JudgeIsStarProject(projectId, currentUserId, orgId int64) (bool, errs.SystemErrorInfo) {
	isExist, err := mysql.IsExistByCond(consts.TableProjectRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcRelationId:   currentUserId,
		consts.TcStatus:       1,
		consts.TcRelationType: consts.IssueRelationTypeStar,
		consts.TcProjectId:    projectId,
	})
	if err != nil {
		return isExist, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return isExist, nil
}

func JudgeIsMember(projectId, currentUserId, orgId int64) (bool, errs.SystemErrorInfo) {
	isExist, err := mysql.IsExistByCond(consts.TableProjectRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcRelationId:   currentUserId,
		consts.TcStatus:       1,
		consts.TcRelationType: db.In(consts.MemberRelationTypeList),
		consts.TcProjectId:    projectId,
	})
	if err != nil {
		return isExist, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return isExist, nil
}

func AddMember(projectId, orgId, userId, currentUserId int64, relationType int) (bool, errs.SystemErrorInfo) {
	memberId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	err1 := mysql.Insert(&po.PpmProProjectRelation{
		Id:           memberId,
		IsDelete:     consts.AppIsNoDelete,
		OrgId:        orgId,
		ProjectId:    projectId,
		Status:       1,
		RelationType: relationType,
		RelationId:   userId,
		Creator:      currentUserId,
		CreateTime:   time.Now(),
		Updator:      currentUserId,
		UpdateTime:   time.Now(),
	})
	if err1 != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return true, nil
}

func DeleteMember(projectId, orgId, userId, currentUserId int64, relationType int) (bool, errs.SystemErrorInfo) {
	_, err := mysql.UpdateSmartWithCond(consts.TableProjectRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcRelationId:   userId,
		consts.TcStatus:       1,
		consts.TcRelationType: relationType,
		consts.TcProjectId:    projectId,
	}, mysql.Upd{
		consts.TcIsDelete:   consts.AppIsDeleted,
		consts.TcUpdator:    currentUserId,
		consts.TcUpdateTime: time.Now(),
	})
	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return true, nil
}

func RemoveProjectMember(orgId, userId int64, input vo.RemoveProjectMemberReq, sourceChannel string) errs.SystemErrorInfo {
	if len(input.MemberIds) == 0 && len(input.MemberForDepartmentID) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	projectId := input.ProjectID

	projectInfo, infoErr := GetProjectSimple(orgId, projectId)
	if infoErr != nil {
		log.Error(infoErr)
		return infoErr
	}
	//负责人不能被移除
	if ok, _ := slice.Contain(input.MemberIds, projectInfo.Owner); ok {
		return errs.CannotRemoveProjectOwner
	}
	//删除项目关联
	if len(input.MemberIds) > 0 {
		_, updateErr := mysql.UpdateSmartWithCond(consts.TableProjectRelation, db.Cond{
			consts.TcProjectId:    projectId,
			consts.TcRelationId:   db.In(input.MemberIds),
			consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeParticipant}),
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcUpdator:  userId,
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if updateErr != nil {
			log.Error(updateErr)
			return errs.MysqlOperateError
		}
		//删除任务关注人
		_, err := mysql.UpdateSmartWithCond(consts.TableIssueRelation, db.Cond{
			consts.TcOrgId:        orgId,
			consts.TcProjectId:    projectId,
			consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeFollower}),
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcRelationId:   db.In(input.MemberIds),
		}, mysql.Upd{
			consts.TcIsDelete:   consts.AppIsDeleted,
			consts.TcUpdator:    userId,
			consts.TcUpdateTime: time.Now(),
		})

		if err != nil {
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
	}
	if len(input.MemberForDepartmentID) > 0 {
		_, updateErr := mysql.UpdateSmartWithCond(consts.TableProjectRelation, db.Cond{
			consts.TcProjectId:    projectId,
			consts.TcRelationId:   db.In(input.MemberForDepartmentID),
			consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeDepartmentParticipant}),
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}, mysql.Upd{
			consts.TcUpdator:  userId,
			consts.TcIsDelete: consts.AppIsDeleted,
		})
		if updateErr != nil {
			log.Error(updateErr)
			return errs.MysqlOperateError
		}
	}

	//清掉用户的角色
	roleErr := orgfacade.RemoveRoleUserRelation(rolevo.RemoveRoleUserRelationReqVo{
		OrgId:      orgId,
		UserIds:    input.MemberIds,
		OperatorId: userId,
	})
	if roleErr.Failure() {
		log.Error(roleErr.Message)
	}
	//清掉部门的角色 todo
	if len(input.MemberForDepartmentID) > 0 {
		//roleErr := rolefacade.RemoveRoleDepartmentRelation(rolevo.RemoveRoleDepartmentRelationReqVo{
		//	OrgId:      orgId,
		//	DeptIds:    input.MemberForDepartmentID,
		//	UserId: userId,
		//})
		//if roleErr.Failure() {
		//	log.Error(roleErr.Message)
		//}
	}

	//最后将用户信息缓存清掉
	clearErr := orgfacade.ClearUserRoleList(rolevo.ClearUserRoleReqVo{
		ProjectId: projectId,
		OrgId:     orgId,
		UserIds:   input.MemberIds,
	})
	if clearErr.Failure() {
		log.Error(clearErr.Error())
		return clearErr.Error()
	}

	refreshProjectAuthErr := RefreshProjectAuthBo(orgId, projectId)
	if refreshProjectAuthErr != nil {
		log.Error(refreshProjectAuthErr)
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:            consts.PushTypeUpdateProjectMembers,
			OrgId:               orgId,
			ProjectId:           projectId,
			OperatorId:          userId,
			BeforeChangeMembers: input.MemberIds,
			AfterChangeMembers:  []int64{},
			Ext:                 ext,
		})
	})

	return nil
}

func AddProjectMember(orgId, userId int64, input vo.RemoveProjectMemberReq, sourceChannel string) errs.SystemErrorInfo {
	if len(input.MemberIds) == 0 && len(input.MemberForDepartmentID) == 0 {
		return errs.UpdateMemberIdsIsEmptyError
	}

	projectId := input.ProjectID
	//判断项目是否存在
	projectInfo, infoErr := GetProjectSimple(orgId, projectId)
	if infoErr != nil {
		log.Error(infoErr)
		return infoErr
	}

	addMemberIds := []int64{}
	if input.MemberIds != nil && len(input.MemberIds) > 0 {
		verifyOrgUserFlag := orgfacade.VerifyOrgUsersRelaxed(orgId, input.MemberIds)
		if !verifyOrgUserFlag {
			log.Error("存在用户组织校验失败")
			return errs.VerifyOrgError
		}
		addIds, updateProjectRelationErr := UpdateProjectRelationWithRelationTypes(userId, orgId, projectId, consts.MemberRelationTypeList, consts.IssueRelationTypeParticipant, input.MemberIds)
		if updateProjectRelationErr != nil {
			log.Error(updateProjectRelationErr)
			return updateProjectRelationErr
		}
		addMemberIds = addIds
	}
	if input.MemberForDepartmentID != nil && len(input.MemberForDepartmentID) > 0 {
		//预先查询已有的关联
		projectRelations := &[]po.PpmProProjectRelation{}
		err5 := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
			consts.TcProjectId:    projectId,
			consts.TcRelationType: consts.IssueRelationTypeDepartmentParticipant,
			consts.TcIsDelete:     consts.AppIsNoDelete,
		}, projectRelations)
		if err5 != nil {
			log.Error(err5)
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		if len(*projectRelations) > 0 {
			for _, relation := range *projectRelations {
				if relation.RelationId == int64(0) {
					//如果项目部门里面包含了所有，就不需要处理别的了
					return nil
				}
			}
		}
	}

	refreshProjectAuthErr := RefreshProjectAuthBo(orgId, projectId)
	if refreshProjectAuthErr != nil {
		log.Error(refreshProjectAuthErr)
	}

	asyn.Execute(func() {
		ext := bo.TrendExtensionBo{}
		ext.ObjName = projectInfo.Name
		PushProjectTrends(bo.ProjectTrendsBo{
			PushType:            consts.PushTypeUpdateProjectMembers,
			OrgId:               orgId,
			ProjectId:           projectId,
			OperatorId:          userId,
			BeforeChangeMembers: []int64{},
			AfterChangeMembers:  addMemberIds,
			Ext:                 ext,
		})
	})

	return nil
}

func GetProjectAllMember(orgId, projectId int64, page, size int) (int64, []bo.ProjectRelationBo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return 0, nil, errs.MysqlOperateError
	}

	cond := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
		//consts.TcRelationId:db.NotEq(0),
		consts.TcRelationType: db.In([]int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant, consts.IssueRelationTypeDepartmentParticipant}),
	}

	pos := &[]po.PpmProProjectRelation{}
	//获取所有成员（最小的relation_id代表最高的用户角色（负责人），最早的创建时间表示加入时间，随机挑选一名操作人）
	//项目部门放在前面
	mid := conn.Select(db.Raw("relation_id, relation_type, create_time, creator")).
		From(consts.TableProjectRelation).
		Where(cond).OrderBy("relation_type asc")

	selectErr := mid.All(pos)
	if selectErr != nil {
		log.Error(selectErr)
		return 0, nil, errs.MysqlOperateError
	}

	bos := &[]bo.ProjectRelationBo{}
	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	if len(*bos) == 0 {
		return 0, *bos, nil
	}
	userIds := []int64{}
	deptIds := []int64{}
	//排列顺序 重新组合（部门->负责人->成员）
	userMap := map[int64]bo.ProjectRelationBo{}
	deptMap := map[int64]bo.ProjectRelationBo{}
	ownerId := int64(0)
	for _, relationBo := range *bos {
		if relationBo.RelationType != consts.IssueRelationTypeDepartmentParticipant {
			if _, ok := userMap[relationBo.RelationId]; !ok {
				userIds = append(userIds, relationBo.RelationId)
				//一个是为了去重，一个是为了把负责人放到第一个
				userMap[relationBo.RelationId] = relationBo
			}
			if relationBo.RelationType == consts.IssueRelationTypeOwner {
				ownerId = relationBo.RelationId
			}
		} else {
			deptIds = append(deptIds, relationBo.RelationId)
			deptMap[relationBo.RelationId] = relationBo
		}
	}

	res := []bo.ProjectRelationBo{}
	if len(deptIds) > 0 {
		departmentsInfo := orgfacade.Departments(orgvo.DepartmentsReqVo{
			Page: nil,
			Size: nil,
			Params: &vo.DepartmentListReq{
				DepartmentIds: deptIds,
			},
			OrgId: orgId,
		})
		if departmentsInfo.Failure() {
			log.Error(departmentsInfo.Error())
			return 0, nil, departmentsInfo.Error()
		}
		if ok, _ := slice.Contain(deptIds, int64(0)); ok {
			if relationInfo, ok := deptMap[int64(0)]; ok {
				res = append(res, relationInfo)
			}
		}
		for _, department := range departmentsInfo.DepartmentList.List {
			if relationInfo, ok := deptMap[department.ID]; ok {
				res = append(res, relationInfo)
			}
		}
	}
	if len(userIds) > 0 {
		deletedUserIds := []int64{}
		userInfos := orgfacade.GetUserInfoByUserIds(orgvo.GetUserInfoByUserIdsReqVo{
			UserIds: userIds,
			OrgId:   orgId,
		})
		if userInfos.Failure() {
			log.Error(userInfos.Error())
			return 0, nil, userInfos.Error()
		}
		if userInfos.GetUserInfoByUserIdsRespVo != nil {
			for _, respVo := range *userInfos.GetUserInfoByUserIdsRespVo {
				if respVo.OrgUserIsDelete == 1 {
					deletedUserIds = append(deletedUserIds, respVo.UserId)
				}
			}
		}
		//负责人放到上面
		if ok, _ := slice.Contain(deletedUserIds, ownerId); !ok {
			res = append(res, userMap[ownerId])
		}
		for i, relationBo := range userMap {
			if i == ownerId {
				continue
			}
			if ok, _ := slice.Contain(deletedUserIds, i); !ok {
				res = append(res, relationBo)
			}
		}
	}

	count := len(res)
	if page > 0 && size > 0 {
		offset := (page - 1) * size
		end := offset + size
		if offset > count {
			res = []bo.ProjectRelationBo{}
		} else {
			if end > count {
				end = count
			}
			res = res[offset:end]
		}
	}

	return int64(count), res, nil
}

func GetProjectMembers(orgId, projectId int64, relationTypes []int64) ([]bo.ProjectRelationBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProjectRelation{}
	cond := db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcRelationType: db.In(relationTypes),
	}
	if projectId > 0 {
		cond[consts.TcProjectId] = projectId
	}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.ProjectRelationBo{}
	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}

	return *bos, nil
}

func GetProjectMemberDepartmentsInfo(orgId, projectId int64) ([]bo.DepartmentSimpleInfoBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProjectRelation{}
	err := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcProjectId:    projectId,
		consts.TcRelationType: consts.IssueRelationTypeDepartmentParticipant,
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	departmentIds := []int64{}
	for _, relation := range *pos {
		departmentIds = append(departmentIds, relation.RelationId)
	}
	res := []bo.DepartmentSimpleInfoBo{}
	if len(departmentIds) == 0 {
		return res, nil
	}
	departmentIds = slice.SliceUniqueInt64(departmentIds)

	departmentsInfo := orgfacade.Departments(orgvo.DepartmentsReqVo{
		Page: nil,
		Size: nil,
		Params: &vo.DepartmentListReq{
			DepartmentIds: departmentIds,
		},
		OrgId: orgId,
	})
	if departmentsInfo.Failure() {
		log.Error(departmentsInfo.Error())
		return nil, departmentsInfo.Error()
	}

	if departmentsInfo.DepartmentList != nil {
		for _, department := range departmentsInfo.DepartmentList.List {
			res = append(res, bo.DepartmentSimpleInfoBo{
				ID:   department.ID,
				Name: department.Name,
			})
		}
	}
	deptUserCountResp := orgfacade.GetUserCountByDeptIds(&orgvo.GetUserCountByDeptIdsReq{
		OrgId:   orgId,
		DeptIds: departmentIds,
	})
	if deptUserCountResp.Failure() {
		log.Error(deptUserCountResp.Error())
		return nil, deptUserCountResp.Error()
	}
	if ok, _ := slice.Contain(departmentIds, int64(0)); ok {
		res = append(res, bo.DepartmentSimpleInfoBo{
			ID:   0,
			Name: "全部",
		})
	}
	if deptUserCountResp.Data != nil {
		userCount := deptUserCountResp.Data.UserCount
		for i, re := range res {
			if count, ok := userCount[re.ID]; ok {
				res[i].UserCount = int64(count)
			}
		}
	}

	return res, nil
}
