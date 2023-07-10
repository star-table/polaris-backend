package domain

import (
	"fmt"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	int64Slice "github.com/star-table/polaris-backend/common/core/util/slice/int64"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"strconv"
	"strings"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func GetDepartmentBoList(page uint, size uint, cond db.Cond) (*[]bo.DepartmentBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectDepartmentByPage(cond, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.DepartmentBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func GetDepartmentBoListByIds(orgId int64, deptIds []int64) (*[]bo.DepartmentBo, errs.SystemErrorInfo) {
	pos, err := dao.SelectDepartment(db.Cond{
		consts.TcOrgId: orgId,
		consts.TcId:    db.In(deptIds),
	})
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.DepartmentBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, nil
}

func GetDepartmentBoWithOrg(id int64, orgId int64) (*bo.DepartmentBo, errs.SystemErrorInfo) {
	po, err := dao.SelectDepartmentByIdAndOrg(id, orgId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.TargetNotExist)
	}
	bo := &bo.DepartmentBo{}
	err1 := copyer.Copy(po, bo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bo, nil
}

func JudgeDepartmentIsExist(orgId int64, name string) (bool, errs.SystemErrorInfo) {
	exist, err := mysql.IsExistByCond(consts.TableDepartment, db.Cond{
		consts.TcName:     name,
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	})
	if err != nil {
		return false, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return exist, nil
}

func GetDepartmentMembers(orgId int64, departmentId *int64) ([]bo.DepartmentMemberInfoBo, errs.SystemErrorInfo) {
	userIdInfoBoList := make([]bo.DepartmentMemberInfoBo, 0)
	if departmentId == nil {
		//不传部门id就是全部
		pos, err := GetUserOrgInfos(orgId)
		if err != nil {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
		userIds := make([]int64, 0)
		for _, info := range pos {
			userId := info.UserId
			userIdExist, _ := slice.Contain(userIds, userId)
			if !userIdExist {
				userIds = append(userIds, userId)
			}
		}
		userInfos, err1 := GetBaseUserInfoBatch(orgId, userIds)
		if err1 != nil {
			log.Error(err1)
			return nil, err1
		}

		cond := db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcUserId:   userIds,
		}
		userDepartmentList, dbErr := dao.SelectUserDepartment(cond)
		if dbErr != nil {
			log.Error(dbErr)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
		userDepartmentMap := map[int64]int64{}
		for _, department := range *userDepartmentList {
			userDepartmentMap[department.UserId] = department.DepartmentId
		}

		for _, info := range userInfos {
			temp := bo.DepartmentMemberInfoBo{
				UserID:        info.UserId,
				Name:          info.Name,
				NamePy:        info.NamePy,
				Avatar:        info.Avatar,
				EmplID:        info.OutUserId,
				OrgUserStatus: info.OrgUserStatus,
			}
			if _, ok := userDepartmentMap[info.UserId]; ok {
				temp.DepartmentID = userDepartmentMap[info.UserId]
			} else {
				temp.DepartmentID = 0
			}
			userIdInfoBoList = append(userIdInfoBoList, temp)
		}

	} else {
		cond := db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}

		userDepartmentList, err := dao.SelectUserDepartment(cond)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}

		userIds := make([]int64, 0)
		for _, userDepartmentInfo := range *userDepartmentList {
			userId := userDepartmentInfo.UserId
			userIdExist, _ := slice.Contain(userIds, userId)
			if !userIdExist {
				userIds = append(userIds, userId)
			}
		}
		userInfos, err1 := GetBaseUserInfoBatch(orgId, userIds)
		if err1 != nil {
			log.Error(err1)
			return nil, err1
		}
		userMap := maps.NewMap("UserId", userInfos)

		for _, userDepartment := range *userDepartmentList {
			if userCacheInfo, ok := userMap[userDepartment.UserId]; ok {
				baseUserInfo := userCacheInfo.(bo.BaseUserInfoBo)
				userIdInfoBoList = append(userIdInfoBoList, bo.DepartmentMemberInfoBo{
					UserID:        baseUserInfo.UserId,
					Name:          baseUserInfo.Name,
					NamePy:        baseUserInfo.NamePy,
					Avatar:        baseUserInfo.Avatar,
					EmplID:        baseUserInfo.OutUserId,
					DepartmentID:  userDepartment.DepartmentId,
					OrgUserStatus: baseUserInfo.OrgUserStatus,
				})
			} else {
				log.Errorf("GetDepartmentMembers: 查询不到部门%d下的用户%d信息", userDepartment.DepartmentId, userDepartment.UserId)
			}
		}
	}
	return userIdInfoBoList, nil
}

func GetUserOrgInfos(orgId int64) ([]*po.PpmOrgUserOrganization, errs.SystemErrorInfo) {
	pos := []*po.PpmOrgUserOrganization{}
	err := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete:    consts.AppIsNoDelete,
		consts.TcOrgId:       orgId,
		consts.TcCheckStatus: consts.AppCheckStatusSuccess,
	}, &pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	return pos, nil
}

func GetTopDepartmentInfoList(orgId int64) ([]bo.DepartmentBo, errs.SystemErrorInfo) {
	departmentInfo := &[]po.PpmOrgDepartment{}
	err := mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcStatus:   consts.AppIsInitStatus,
		consts.TcParentId: 0,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, departmentInfo)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	departmentInfoBo := &[]bo.DepartmentBo{}
	err1 := copyer.Copy(departmentInfo, departmentInfoBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *departmentInfoBo, nil
}

func GetTopDepartmentInfo(orgId int64) (*bo.DepartmentBo, errs.SystemErrorInfo) {
	departmentInfoList, err := GetTopDepartmentInfoList(orgId)
	if err != nil {
		log.Errorf("获取部门信息错误 %v", err)
		return nil, err
	}
	if len(departmentInfoList) == 0 {
		log.Errorf("组织%d下不存在顶级部门", orgId)
		return nil, errs.BuildSystemErrorInfo(errs.TopDepartmentNotExist)
	}
	departmentInfo := departmentInfoList[0]
	return &departmentInfo, nil
}

func BoundOrgMemberToTopDepartment(orgId int64, userIds []int64, operatorId int64) (int64, errs.SystemErrorInfo) {
	departmentInfo, err := GetTopDepartmentInfoList(orgId)
	if err != nil {
		log.Error("获取部门信息错误 " + strs.ObjectToString(err))
		return 0, err
	}
	var departmentId int64
	for _, v := range departmentInfo {
		departmentId = v.Id
		break
	}
	return departmentId, BoundDepartmentUser(orgId, userIds, departmentId, operatorId, false)
}

//解绑部门用户，解绑当前用户所在的所有部门
func UnBoundDepartmentUser(orgId int64, userIds []int64, operatorId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	//查询已有的绑定关系
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(userIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	_, err := dao.UpdateUserDepartmentByCond(cond, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcUpdator:  operatorId,
	}, tx)
	if err != nil {
		log.Error(err)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

//绑定部门用户，带分布式锁
func BoundDepartmentUser(orgId int64, userIds []int64, departmentId, operatorId int64, isLeaderFlag bool) errs.SystemErrorInfo {
	isLeader := 2
	if isLeaderFlag {
		isLeader = 1
	}

	//先上锁
	lockKey := consts.UserAndDepartmentRelationLockKey + fmt.Sprintf("%d:%d", orgId, departmentId)
	uuid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uuid)
	if lockErr != nil {
		log.Error(lockErr)
		return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	if !suc {
		log.Errorf("绑定用户时没有获取到锁 orgId %d departmentId %d", orgId, departmentId)
		return errs.BuildSystemErrorInfo(errs.TryDistributedLockError)
	}
	defer func() {
		if _, err := cache.ReleaseDistributedLock(lockKey, uuid); err != nil {
			log.Error(err)
		}
	}()
	//查询已有的绑定关系
	cond := db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcDepartmentId: departmentId,
		consts.TcUserId:       db.In(userIds),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}
	userDepartmentList, dbErr := dao.SelectUserDepartment(cond)
	if dbErr != nil {
		log.Error(dbErr)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	notRelationUserIds := make([]int64, 0)
	alreadyRelationUserIdMap := map[int64]bool{}
	if userDepartmentList != nil && len(*userDepartmentList) > 0 {
		for _, userDepartment := range *userDepartmentList {
			alreadyRelationUserIdMap[userDepartment.UserId] = true
		}
	}
	//获取没有关联关系的用户
	for _, userId := range userIds {
		if _, ok := alreadyRelationUserIdMap[userId]; !ok {
			notRelationUserIds = append(notRelationUserIds, userId)
			alreadyRelationUserIdMap[userId] = true
		}
	}
	if len(notRelationUserIds) == 0 {
		return nil
	}

	userIdsLen := len(notRelationUserIds)

	ids, err := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, userIdsLen)
	if err != nil {
		log.Error(err)
		return err
	}

	userDepartments := make([]po.PpmOrgUserDepartment, userIdsLen)
	for i, userId := range notRelationUserIds {
		userDepartments[i] = po.PpmOrgUserDepartment{
			Id:           ids.Ids[i].Id,
			OrgId:        orgId,
			UserId:       userId,
			DepartmentId: departmentId,
			IsLeader:     isLeader,
			Creator:      operatorId,
			Updator:      operatorId,
		}
	}

	err1 := dao.InsertUserDepartmentBatch(userDepartments)
	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	return nil
}

func GetDepartmentInfo(orgId int64, departmentId int64) (*bo.DepartmentBo, errs.SystemErrorInfo) {
	info := &po.PpmOrgDepartment{}
	err := mysql.SelectOneByCond(consts.TableDepartment, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcId:       departmentId,
	}, info)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.DepartmentNotExist
		} else {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
	}

	result := &bo.DepartmentBo{}
	_ = copyer.Copy(info, result)
	return result, nil
}

func GetChildrenDepartmentIds(orgId int64, departmentIds []int64) ([]int64, errs.SystemErrorInfo) {
	//获取部门信息
	departmentInfo := &[]po.PpmOrgDepartment{}
	err := mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       db.In(departmentIds),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, departmentInfo)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	allDepartmentIds := []int64{}
	for _, department := range *departmentInfo {
		childrenInfo := &[]po.PpmOrgDepartment{}
		err := mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcPath:     db.Like(department.Path + "," + strconv.FormatInt(department.Id, 10)),
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, departmentInfo)
		if err != nil {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}

		for _, orgDepartment := range *childrenInfo {
			allDepartmentIds = append(allDepartmentIds, orgDepartment.Id)
		}
	}

	allDepartmentIds = append(allDepartmentIds, departmentIds...)

	return slice.SliceUniqueInt64(allDepartmentIds), nil
}

func GetUserDepartmentInfo(orgId int64, userIds []int64) ([]bo.UserDepartmentInfo, errs.SystemErrorInfo) {
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.UserDepartmentInfo{}
	err1 := conn.Select(db.Raw("u.user_id,u.department_id,u.is_leader,d.name")).From("ppm_org_user_department u", "ppm_org_department d").Where(db.Cond{
		"u." + consts.TcDepartmentId: db.Raw("d.id"),
		"u." + consts.TcIsDelete:     consts.AppIsNoDelete,
		"d." + consts.TcIsDelete:     consts.AppIsNoDelete,
		"u." + consts.TcUserId:       db.In(userIds),
		"u." + consts.TcOrgId:        orgId,
		"d." + consts.TcOrgId:        orgId,
	}).All(bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.MysqlOperateError
	}

	return *bos, nil
}

// 获取员工部门id的映射，员工id => 部门id
func GetUserDepartmentIdMap(orgId int64, userIds []int64) (map[int64]int64, error) {
	userDeptIdMap := map[int64]int64{}
	userDepts, err := GetUserDepartmentInfo(orgId, userIds)
	if err != nil {
		return userDeptIdMap, err
	}
	for _, oneInfo := range userDepts {
		userDeptIdMap[oneInfo.UserId] = oneInfo.DepartmentId
	}

	return userDeptIdMap, nil
}

// 通过 department ids 删除部门信息
func DeleteDepartmentByOrgId(orgId int64) error {
	_, err := mysql.UpdateSmartWithCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	return err
}

// 通过 out department ids 删除部门 out 信息
func DeleteDepartmentOutInfoByOrgId(orgId int64) error {
	_, err := mysql.UpdateSmartWithCond(consts.TableDepartmentOutInfo, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	return err
}

// 通过 department ids 删除部门信息
func DeleteDepartmentByDepartmentIds(orgId int64, departmentIds []string) error {
	if len(departmentIds) < 1 {
		return errs.InputParamEmpty
	}
	_, err := mysql.UpdateSmartWithCond(consts.TableDepartment, db.Cond{
		consts.TcOrgId:        orgId,
		consts.TcDepartmentId: db.In(departmentIds),
		consts.TcIsDelete:     consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	return err
}

// 通过 department ids 删除用户与部门的关联关系
func DeleteUserDepartmentByOrgId(orgId int64) error {
	_, err := mysql.UpdateSmartWithCond(consts.TableUserDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
	})
	return err
}

// 通过 out department ids 获取星书方的 department ids
func GetDepartmentIdsByOutDepartmentIds(orgId int64, outDepartmentIds []string) ([]int64, errs.SystemErrorInfo) {
	if len(outDepartmentIds) < 1 {
		return nil, errs.InputParamEmpty
	}
	var list []struct {
		Id                 int64  `db:"id,omitempty" json:"id"`
		OrgId              int64  `db:"org_id,omitempty" json:"orgId"`
		DepartmentId       int64  `db:"department_id,omitempty" json:"departmentId"`
		OutOrgDepartmentId string `db:"out_org_department_id,omitempty" json:"outOrgDepartmentId"`
		IsDelete           string `db:"is_delete,omitempty" json:"isDelete"`
	}
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	err = conn.Select(db.Raw("do.id,do.department_id,do.out_org_department_id")).From("ppm_org_department_out_info do").Where(db.Cond{
		"do." + consts.TcOutOrgDepartmentId: db.In(outDepartmentIds),
		"do." + consts.TcOrgId:              orgId,
		"do." + consts.TcIsDelete:           consts.AppIsNoDelete,
	}).All(&list)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	var deptIds []int64
	for _, oneInfo := range list {
		deptIds = append(deptIds, oneInfo.DepartmentId)
	}
	return deptIds, nil
}

// 通过多个 outDeptId，查询对应的我方部门id
func GetOutDeptIdMapToDeptId(orgId int64, outDeptIds []string) (map[string]int64, error) {
	var mapOutDeptIdToDeptId = map[string]int64{}
	var list []struct {
		Id                 int64  `db:"id,omitempty" json:"id"`
		OrgId              int64  `db:"org_id,omitempty" json:"orgId"`
		DepartmentId       int64  `db:"department_id,omitempty" json:"departmentId"`
		OutOrgDepartmentId string `db:"out_org_department_id,omitempty" json:"outOrgDepartmentId"`
		IsDelete           string `db:"is_delete,omitempty" json:"isDelete"`
	}
	conn, err := mysql.GetConnect()
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	err = conn.Select(db.Raw("do.id,do.department_id,do.out_org_department_id, do.is_delete")).From("ppm_org_department_out_info do").Where(db.Cond{
		"do." + consts.TcOutOrgDepartmentId: db.In(outDeptIds),
		"do." + consts.TcOrgId:              orgId,
		"do." + consts.TcIsDelete:           consts.AppIsNoDelete,
	}).All(&list)
	for _, oneOutDept := range list {
		mapOutDeptIdToDeptId[oneOutDept.OutOrgDepartmentId] = oneOutDept.DepartmentId
	}

	return mapOutDeptIdToDeptId, nil
}

func GetDepartmentMembersPaginate(orgId int64, sourceChannel string, name *string, page, size int, userIds []int64, ignoreDelete bool, projectId, relationType int64) (int64, []bo.DepartmentMemberInfoBo, errs.SystemErrorInfo) {
	userIdInfoBoList := make([]bo.DepartmentMemberInfoBo, 0)
	total := int64(0)
	pos := &[]po.PpmOrgUserOrganization{}
	excludeUserIds := []int64{}
	if projectId != 0 {
		projectUserResp := projectfacade.GetProjectRelationUserIds(projectvo.GetProjectRelationUserIdsReq{
			ProjectId:    projectId,
			RelationType: &relationType,
		})
		if projectUserResp.Failure() {
			log.Error(projectUserResp.Error())
			return 0, nil, projectUserResp.Error()
		}
		excludeUserIds = projectUserResp.UserIds
	}
	if name == nil || *name == "" {
		cond := db.Cond{
			consts.TcOrgId:       orgId,
			consts.TcCheckStatus: consts.AppCheckStatusSuccess,
		}
		if userIds != nil && len(userIds) > 0 {
			cond[consts.TcUserId] = db.In(userIds)
		}
		if len(excludeUserIds) > 0 {
			cond[consts.TcUserId+" "] = db.NotIn(excludeUserIds)
		}
		if !ignoreDelete {
			cond[consts.TcIsDelete] = consts.AppIsNoDelete
		}
		count, err := mysql.SelectAllByCondWithPageAndOrder(consts.TableUserOrganization, cond, nil, page, size, nil, pos)
		if err != nil {
			log.Error(err)
			return 0, nil, errs.MysqlOperateError
		}
		total = int64(count)
	} else {
		conn, err := mysql.GetConnect()
		if err != nil {
			log.Error(err)
			return 0, nil, errs.MysqlOperateError
		}
		namePy := strings.ToLower(*name)
		cond := db.Cond{
			"o." + consts.TcOrgId:       orgId,
			"o." + consts.TcCheckStatus: consts.AppCheckStatusSuccess,
			"o." + consts.TcUserId:      db.Raw("u." + consts.TcId),
		}
		if userIds != nil && len(userIds) > 0 {
			cond[" o."+consts.TcUserId] = db.In(userIds)
		}
		if len(excludeUserIds) > 0 {
			cond[" o."+consts.TcUserId+" "] = db.NotIn(excludeUserIds)
		}
		if !ignoreDelete {
			cond["o."+consts.TcIsDelete] = consts.AppIsNoDelete
		}
		mid := conn.Select(db.Raw("o.user_id")).From("ppm_org_user_organization as o", "ppm_org_user as u").Where(cond).And(db.Or(db.Cond{
			"u." + consts.TcName: db.Like("%" + *name + "%"),
		}, db.Cond{
			db.Raw("lower(u." + consts.TcNamePinyin + ")"): db.Like("%" + namePy + "%"),
		}))
		allPos := &[]po.PpmOrgUserOrganization{}
		allErr := mid.All(allPos)
		if allErr != nil {
			log.Error(allErr)
			return 0, nil, errs.MysqlOperateError
		}
		total = int64(len(*allPos))
		if page > 0 && size > 0 {
			mid = mid.Offset((page - 1) * size).Limit(size)
		}
		needErr := mid.All(pos)
		if needErr != nil {
			log.Error(needErr)
			return 0, nil, errs.MysqlOperateError
		}
	}

	trulyUserIds := make([]int64, 0)
	for _, info := range *pos {
		userId := info.UserId
		userIdExist, _ := slice.Contain(trulyUserIds, userId)
		if !userIdExist {
			trulyUserIds = append(trulyUserIds, userId)
		}
	}
	userInfos, err1 := GetBaseUserInfoBatch(orgId, trulyUserIds)
	if err1 != nil {
		log.Error(err1)
		return 0, nil, err1
	}

	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId:   db.In(trulyUserIds),
	}
	userDepartmentList, err := dao.SelectUserDepartment(cond)
	if err != nil {
		log.Error(err)
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}
	userDepartmentMap := map[int64]int64{}
	for _, department := range *userDepartmentList {
		userDepartmentMap[department.UserId] = department.DepartmentId
	}

	for _, info := range userInfos {
		temp := bo.DepartmentMemberInfoBo{
			UserID:        info.UserId,
			Name:          info.Name,
			NamePy:        info.NamePy,
			Avatar:        info.Avatar,
			EmplID:        info.OutUserId,
			OrgUserStatus: info.OrgUserStatus,
		}
		if _, ok := userDepartmentMap[info.UserId]; ok {
			temp.DepartmentID = userDepartmentMap[info.UserId]
		} else {
			temp.DepartmentID = 0
		}
		userIdInfoBoList = append(userIdInfoBoList, temp)
	}

	return total, userIdInfoBoList, nil
}

func GetUserDepartmentIdsByUserIds(orgId int64, userIds []int64) ([]int64, errs.SystemErrorInfo) {
	pos := []*po.PpmOrgUserDepartment{}
	err := mysql.SelectAllByCond(consts.TableUserDepartment, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcUserId:   db.In(userIds),
	}, &pos)
	if err != nil {
		log.Errorf("[GetUserDepartmentIdsByUserIds] err:%v, orgId:%v, userIds:%v", err, orgId, userIds)
		return nil, errs.MysqlOperateError
	}
	deptIds := make([]int64, 0, len(pos))
	for _, dept := range pos {
		deptIds = append(deptIds, dept.DepartmentId)
	}
	return deptIds, nil
}

func BindUserDepartments(orgId, userId int64, departmentIds []int64) errs.SystemErrorInfo {
	idResp, errSys := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(departmentIds))
	if errSys != nil {
		log.Errorf("[BindUserDepartments] idfacade err:%v, orgId:%v, deptIds:%v", errSys, orgId, departmentIds)
		return errSys
	}
	pos := []po.PpmOrgUserDepartment{}
	for i, id := range departmentIds {
		pos = append(pos, po.PpmOrgUserDepartment{
			Id:           idResp.Ids[i].Id,
			OrgId:        orgId,
			UserId:       userId,
			DepartmentId: id,
		})
	}
	err := mysql.BatchInsert(&po.PpmOrgUserDepartment{}, slice.ToSlice(pos))
	if err != nil {
		log.Errorf("[BindUserDepartments] BatchInsert err:%v, orgId:%v, deptIds:%v", err, orgId, departmentIds)
		return errs.MysqlOperateError
	}

	return nil
}

func ChangeUserDept(orgId, operatorId int64, userId int64, deptIds []int64) errs.SystemErrorInfo {
	if len(deptIds) > 0 {
		deptIds = slice.SliceUniqueInt64(deptIds)
		deptList, err := dao.SelectDepartment(db.Cond{
			consts.TcOrgId:    orgId,
			consts.TcId:       db.In(deptIds),
			consts.TcIsDelete: consts.AppIsNoDelete,
		})
		if err != nil {
			log.Errorf("[ChangeUserDept] GetDepartmentBoListByIds err:%v, orgId:%v, deptIds:%v", err, orgId, deptIds)
			return errs.MysqlOperateError
		}
		if len(*deptList) != len(deptIds) {
			return errs.DepartmentNotExist
		}
	}
	// 查询原关系
	oldDeptIds, errSys := GetUserDepartmentIdsByUserIds(orgId, []int64{userId})
	if errSys != nil {
		log.Errorf("[ChangeUserDept] GetUserDepartmentIdsByUserIds err:%v, orgId:%v, deptIds:%v", errSys, orgId, deptIds)
		return errSys
	}

	_, addDeptIds, delDeptIds := int64Slice.CompareSliceAddDelInt64(deptIds, oldDeptIds)
	if len(addDeptIds) == 0 && len(delDeptIds) == 0 {
		return nil
	}

	addBindList := []interface{}{}
	if len(addDeptIds) > 0 {
		idResp, errSys := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(addDeptIds))
		if errSys != nil {
			log.Errorf("[ChangeUserDept] idfacade err:%v, orgId:%v, deptIds:%v", errSys, orgId, deptIds)
			return errSys
		}
		for i, deptId := range addDeptIds {
			addBindList = append(addBindList, po.PpmOrgUserDepartment{
				Id:           idResp.Ids[i].Id,
				OrgId:        orgId,
				UserId:       userId,
				DepartmentId: deptId,
			})
		}
	}

	delBindDeptIds := delDeptIds
	dbErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(delBindDeptIds) > 0 {
			_, dbErr := mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
				consts.TcIsDelete:     consts.AppIsNoDelete,
				consts.TcDepartmentId: db.In(delBindDeptIds),
				consts.TcUserId:       userId,
				consts.TcOrgId:        orgId,
			}, mysql.Upd{
				consts.TcIsDelete: consts.AppIsDeleted,
				consts.TcUpdator:  operatorId,
			})
			if dbErr != nil {
				log.Errorf("[ChangeUserDept] err:%v, orgId:%v", dbErr, orgId)
				return dbErr
			}
		}
		if len(addBindList) > 0 {
			dbErr := mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, addBindList)
			if dbErr != nil {
				log.Errorf("[ChangeUserDept] err:%v, orgId:%v", dbErr, orgId)
				return dbErr
			}
		}
		return nil
	})

	if dbErr != nil {
		log.Errorf("[ChangeUserDept] update err:%v, orgId:%v", dbErr, orgId)
		return errs.MysqlOperateError
	}
	return nil
}
