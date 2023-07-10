package service

import (
	"strconv"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/id/snowflake"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func Departments(page uint, size uint, params *vo.DepartmentListReq, orgId int64) (*vo.DepartmentList, errs.SystemErrorInfo) {

	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcStatus:   consts.AppStatusEnable,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcIsHide:   consts.AppIsNotHiding, //默认只查询非隐藏部门
	}

	if params != nil {
		//查询父部门的子部门信息
		if params.ParentID != nil {
			cond[consts.TcParentId] = params.ParentID
		}
		//名称
		if params.Name != nil {
			cond[consts.TcName] = db.Like("%" + *params.Name + "%")
		}
		//查询顶级部门
		if params.IsTop != nil && *params.IsTop == 1 {
			cond[consts.TcParentId] = 0
		}
		//展示隐藏的部门
		if params.ShowHiding != nil && *params.ShowHiding == 1 {
			delete(cond, consts.TcIsHide)
		}
		if params.DepartmentIds != nil && len(params.DepartmentIds) != 0 {
			cond[consts.TcId] = db.In(params.DepartmentIds)
		}
	}

	departmentBos, total, err := domain.GetDepartmentBoList(page, size, cond)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.DepartmentDomainError, err)
	}

	resultList := &[]*vo.Department{}
	copyErr := copyer.Copy(departmentBos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &vo.DepartmentList{
		Total: total,
		List:  *resultList,
	}, nil
}

func DepartmentMembers(params vo.DepartmentMemberListReq, orgId int64) ([]*vo.DepartmentMemberInfo, errs.SystemErrorInfo) {
	//currentUserInfo, err := GetCurrentUser(ctx)
	//if err != nil {
	//	return nil, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	//}
	//orgId := currentUserInfo.OrgId
	departmentId := params.DepartmentID

	//departmentBo, err := domain.GetDepartmentBoWithOrg(departmentId, orgId)
	//if err != nil{
	//	log.Error(err)
	//	return nil, errs.BuildSystemErrorInfo(errs.DepartmentNotExist)
	//}

	userIdInfoBoList, err := domain.GetDepartmentMembers(orgId, departmentId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.DepartmentDomainError, err)
	}

	userIdInfoVos := &[]*vo.DepartmentMemberInfo{}
	err1 := copyer.Copy(userIdInfoBoList, userIdInfoVos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return *userIdInfoVos, nil
}

func CreateDepartment(params orgvo.CreateDepartmentReq, orgId int64, userId int64) (*vo.Void, errs.SystemErrorInfo) {
	//检验部门名称
	if !format.VerifyDepartmentName(params.Name) {
		return nil, errs.DepartmentNameInvalid
	}

	//查看部门名称是否重复
	isExist, err := mysql.IsExistByCond(consts.TableDepartment, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcName:     params.Name,
	})
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}
	if isExist {
		return nil, errs.DepartmentExistAlready
	}

	//查询父部门
	var path string
	if params.ParentID != 0 {
		parentInfo, err := domain.GetDepartmentInfo(orgId, params.ParentID)
		if err != nil {
			log.Error(err)
			return nil, err
		}

		path = parentInfo.Path + "," + strconv.FormatInt(params.ParentID, 10)
	}

	//用户id
	id, deptIdErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableDepartment)
	if deptIdErr != nil {
		log.Error(deptIdErr)
		return nil, deptIdErr
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//新建部门
		insertErr := mysql.TransInsert(tx, &po.PpmOrgDepartment{
			Id:       id,
			OrgId:    orgId,
			Name:     params.Name,
			ParentId: params.ParentID,
			Sort:     0,
			Creator:  userId,
			Path:     path,
		})
		if insertErr != nil {
			log.Error(insertErr)
			return insertErr
		}
		//添加部门主管
		if params.LeaderIds != nil && len(*params.LeaderIds) > 0 {
			var departmentUserList []interface{}
			for _, leaderId := range *params.LeaderIds {
				departmentUserList = append(departmentUserList, po.PpmOrgUserDepartment{
					Id:           id,
					OrgId:        orgId,
					UserId:       leaderId,
					DepartmentId: id,
					IsLeader:     1,
					Creator:      userId,
				})
			}
			err := mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, departmentUserList)
			if err != nil {
				log.Error(err)
				return err
			}
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	return &vo.Void{ID: id}, nil
}

func UpdateDepartment(params orgvo.UpdateDepartmentReq, orgId, userId int64) (*vo.Void, errs.SystemErrorInfo) {
	info, err := domain.GetDepartmentInfo(orgId, params.DepartmentId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if params.Name != nil && *params.Name != "" && *params.Name != info.Name {
		//检验部门名称
		if !format.VerifyDepartmentName(*params.Name) {
			return nil, errs.DepartmentNameInvalid
		}

		_, updateErr := mysql.UpdateSmartWithCond(consts.TableDepartment, db.Cond{
			consts.TcId: params.DepartmentId,
		}, mysql.Upd{
			consts.TcName:    params.Name,
			consts.TcUpdator: userId,
		})
		if updateErr != nil {
			log.Error(updateErr)
			return nil, errs.MysqlOperateError
		}
	}

	result := &vo.Void{ID: params.DepartmentId}
	if params.LeaderIds == nil {
		return result, nil
	}

	//查看之前部门的主管
	oldLeaders := &[]po.PpmOrgUserDepartment{}
	leaderErr := mysql.SelectAllByCond(consts.TableUserDepartment, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcDepartmentId: params.DepartmentId,
		//consts.TcIsLeader:1,
	}, oldLeaders)
	if leaderErr != nil {
		log.Error(leaderErr)
		return nil, errs.MysqlOperateError
	}

	var oldLeaderIds, oldUserIds []int64
	for _, department := range *oldLeaders {
		if department.IsLeader == 1 {
			oldLeaderIds = append(oldLeaderIds, department.UserId)
		}
		oldUserIds = append(oldUserIds, department.UserId)
	}

	deleteIds, addIds := util.GetDifMemberIds(oldLeaderIds, *params.LeaderIds)
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		if len(deleteIds) > 0 {
			_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
				consts.TcUserId:       db.In(deleteIds),
				consts.TcOrgId:        orgId,
				consts.TcDepartmentId: params.DepartmentId,
			}, mysql.Upd{
				consts.TcIsLeader: 2,
				consts.TcUpdator:  userId,
			})
			if err != nil {
				log.Error(err)
				return err
			}
		}

		if len(addIds) > 0 {
			updUserIds := []int64{}
			addUserIds := []int64{}
			addRelations := []interface{}{}
			for _, id := range addIds {
				if ok, _ := slice.Contain(oldUserIds, id); ok {
					//如果本来就是部门用户
					updUserIds = append(updUserIds, id)
				} else {
					//如果没有要新增
					addUserIds = append(addUserIds, id)
				}
			}
			if len(addUserIds) > 0 {
				idResp, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableUserDepartment, len(addUserIds))
				if idErr != nil {
					log.Error(idErr)
					return idErr
				}
				for i, id := range addIds {
					addRelations = append(addRelations, po.PpmOrgUserDepartment{
						Id:           idResp.Ids[i].Id,
						OrgId:        orgId,
						UserId:       id,
						DepartmentId: params.DepartmentId,
						IsLeader:     1,
						Creator:      userId,
					})
				}
			}
			if len(updUserIds) > 0 {
				_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
					consts.TcUserId:       db.In(updUserIds),
					consts.TcOrgId:        orgId,
					consts.TcDepartmentId: params.DepartmentId,
				}, mysql.Upd{
					consts.TcIsLeader: 1,
					consts.TcUpdator:  userId,
				})
				if err != nil {
					log.Error(err)
					return err
				}
			}
			if len(addRelations) > 0 {
				err := mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, addRelations)
				if err != nil {
					log.Error(err)
					return err
				}
			}
		}
		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	return result, nil
}

func DeleteDepartment(params orgvo.DeleteDepartmentReq, orgId, userId int64) (*vo.Void, errs.SystemErrorInfo) {
	info, err := domain.GetDepartmentInfo(orgId, params.DepartmentId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result := &vo.Void{ID: params.DepartmentId}

	allDepartmentIds := []int64{params.DepartmentId}
	//找部门下的子部门
	childrenDepartments := &[]po.PpmOrgDepartment{}
	childrenErr := mysql.SelectAllByCond(consts.TableDepartment, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcPath:     db.Like(info.Path + "," + strconv.FormatInt(params.DepartmentId, 10) + "%"),
	}, childrenDepartments)
	if childrenErr != nil {
		log.Error(childrenErr)
		return nil, errs.MysqlOperateError
	}
	for _, department := range *childrenDepartments {
		allDepartmentIds = append(allDepartmentIds, department.Id)
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除部门
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableDepartment, db.Cond{
			consts.TcId:    db.In(allDepartmentIds),
			consts.TcOrgId: orgId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  userId,
		})
		if err != nil {
			log.Error(err)
			return err
		}
		//删除部门用户关系
		_, err1 := mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
			consts.TcDepartmentId: db.In(allDepartmentIds),
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcOrgId:        orgId,
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  userId,
		})
		if err1 != nil {
			log.Error(err1)
			return err
		}

		return nil
	})
	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}

	return result, nil
}

//已废弃，使用usercenter
func AllocateDepartment(params orgvo.AllocateDepartmentReq, orgId, userId int64) errs.SystemErrorInfo {
	if len(params.UserIds) == 0 {
		return nil
	}
	//查看用户是否是有效用户
	orgUsers := &[]po.PpmOrgUserOrganization{}
	err := mysql.SelectAllByCond(consts.TableUserOrganization, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcUserId:   db.In(params.UserIds),
	}, orgUsers)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}
	params.UserIds = slice.SliceUniqueInt64(params.UserIds)
	if len(*orgUsers) != len(params.UserIds) {
		return errs.UserNotExist
	}

	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		//删除旧有部门关联
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableUserDepartment, db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
			consts.TcUserId:   db.In(params.UserIds),
		}, mysql.Upd{
			consts.TcIsDelete: consts.AppIsDeleted,
			consts.TcUpdator:  userId,
		})
		if err != nil {
			log.Error(err)
			return errs.MysqlOperateError
		}

		//添加新的关联关系
		var departmentUser []interface{}
		for _, id := range params.UserIds {
			for _, departmentId := range params.DepartmentIds {
				departmentUser = append(departmentUser, po.PpmOrgUserDepartment{
					Id:           snowflake.Id(),
					OrgId:        orgId,
					UserId:       id,
					DepartmentId: departmentId,
					Creator:      userId,
				})
			}
		}
		insertErr := mysql.TransBatchInsert(tx, &po.PpmOrgUserDepartment{}, departmentUser)
		if insertErr != nil {
			log.Error(insertErr)
			return insertErr
		}

		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return errs.MysqlOperateError
	}

	return nil
}

func SetUserDepartmentLevel(orgId, userId int64, params orgvo.SetUserDepartmentLevelReq) errs.SystemErrorInfo {
	if ok, _ := slice.Contain([]int{consts.UserTeamRelationTypeLeader, consts.UserTeamRelationTypeMember}, params.IsLeader); !ok {
		return errs.ParamError
	}
	//查看用户是否存在
	_, _, infoErr := domain.GetUserInfo(orgId, params.UserId, "")
	if infoErr != nil {
		log.Error(infoErr)
		return infoErr
	}
	//查询部门用户是否存在
	departmentInfo := &po.PpmOrgUserDepartment{}
	err := mysql.SelectOneByCond(consts.TableUserDepartment, db.Cond{
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcOrgId:        orgId,
		consts.TcUserId:       params.UserId,
		consts.TcDepartmentId: params.DepartmentId,
	}, departmentInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			//判断用户是否存在
			id, idErr := idfacade.ApplyPrimaryIdRelaxed(consts.TableDepartment)
			if idErr != nil {
				log.Error(idErr)
				return idErr
			}
			insertErr := mysql.Insert(&po.PpmOrgUserDepartment{
				Id:           id,
				OrgId:        orgId,
				UserId:       params.UserId,
				DepartmentId: params.DepartmentId,
				IsLeader:     params.IsLeader,
				Creator:      userId,
			})
			if insertErr != nil {
				log.Error(insertErr)
				return errs.MysqlOperateError
			}
		} else {
			log.Error(err)
			return errs.MysqlOperateError
		}
	}
	if departmentInfo.IsLeader == params.IsLeader {
		return nil
	}

	_, updErr := mysql.UpdateSmartWithCond(consts.TableUserDepartment, db.Cond{
		consts.TcId: departmentInfo.Id,
	}, mysql.Upd{
		consts.TcIsLeader: params.IsLeader,
		consts.TcUpdator:  userId,
	})
	if updErr != nil {
		log.Error(updErr)
		return errs.MysqlOperateError
	}

	return nil
}

func DepartmentMembersList(orgId int64, sourceChannel string, name *string, page, size int, userIds []int64, ignoreDelete bool, projectId, relationType int64) (*vo.DepartmentMembersListResp, errs.SystemErrorInfo) {
	total, userIdInfoBoList, err := domain.GetDepartmentMembersPaginate(orgId, sourceChannel, name, page, size, userIds, ignoreDelete, projectId, relationType)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	userIdInfoVos := &[]*vo.DepartmentMemberInfo{}
	err1 := copyer.Copy(userIdInfoBoList, userIdInfoVos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return &vo.DepartmentMembersListResp{
		Total: total,
		List:  *userIdInfoVos,
	}, nil
}

func VerifyDepartments(orgId int64, departmentIds []int64) bool {
	departmentIds = slice.SliceUniqueInt64(departmentIds)
	total, err := mysql.SelectCountByCond(consts.TableDepartment, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcId:       db.In(departmentIds),
		consts.TcStatus:   consts.AppStatusEnable,
	})
	if err != nil {
		log.Error(err)
		return false
	}

	if int(total) != len(departmentIds) {
		return false
	}

	return true
}

func GetDeptByIds(orgId int64, deptIds []int64) ([]*vo.Department, errs.SystemErrorInfo) {
	departmentBos, err := domain.GetDepartmentBoListByIds(orgId, deptIds)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.DepartmentDomainError, err)
	}

	resultList := &[]*vo.Department{}
	copyErr := copyer.Copy(departmentBos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return *resultList, nil
}
