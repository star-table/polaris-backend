package domain

import (
	"strconv"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/permissionfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

// 获取项目信息以及 关联项目的管理员（不包含系统管理员）
func GetProject(orgId int64, projectId int64) (*bo.ProjectBo, errs.SystemErrorInfo) {
	project := &po.PpmProProject{}
	err := mysql.SelectOneByCond(project.TableName(), db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcId:       projectId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, project)
	if err != nil {
		log.Error(err)
		return nil, errs.ProjectNotExist
	}
	projectBo := &bo.ProjectBo{}
	err1 := copyer.Copy(project, projectBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	//负责人集合
	//relationPos := &[]po.PpmProProjectRelation{}
	//relationErr := mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
	//	consts.TcOrgId:        orgId,
	//	consts.TcIsDelete:     consts.AppIsNoDelete,
	//	consts.TcRelationType: consts.IssueRelationTypeOwner,
	//	consts.TcProjectId:    projectId,
	//}, relationPos)
	//if relationErr != nil {
	//	log.Error(relationErr)
	//	return nil, errs.MysqlOperateError
	//}
	//for _, relation := range *relationPos {
	//	projectBo.OwnerIds = append(projectBo.OwnerIds, relation.RelationId)
	//}
	// 从无码的管理员角色中找出关联的项目管理员
	extraAdminUserIds, busiErr := GetProAdminUserIdsForLessCodeApp(orgId, 0, projectBo.AppId)
	if busiErr != nil {
		log.Error(busiErr)
		return nil, busiErr
	}
	projectBo.OwnerIds = append(projectBo.OwnerIds, extraAdminUserIds...)
	// 如果从上面的接口中找不到创建项目的owner，就把project 中的owner放上去
	projectBo.OwnerIds = append(projectBo.OwnerIds, projectBo.Owner)
	projectBo.OwnerIds = slice.SliceUniqueInt64(projectBo.OwnerIds)
	return projectBo, nil
}

func GetProjectByAppId(orgId int64, proAppId int64) (*bo.ProjectBo, errs.SystemErrorInfo) {
	project := &po.PpmProProject{}
	err := mysql.SelectOneByCond(project.TableName(), db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcAppId:    proAppId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}, project)
	if err != nil {
		log.Error(err)
		return nil, errs.ProjectNotExist
	}
	projectBo := &bo.ProjectBo{}
	err1 := copyer.Copy(project, projectBo)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	// 从无码的管理员角色中找出关联的项目管理员
	extraAdminUserIds, busiErr := GetProAdminUserIdsForLessCodeApp(orgId, 0, projectBo.AppId)
	if busiErr != nil {
		log.Error(busiErr)
		return nil, busiErr
	}
	projectBo.OwnerIds = append(projectBo.OwnerIds, extraAdminUserIds...)
	// 如果从上面的接口中找不到创建项目的owner，就把project 中的owner放上去
	projectBo.OwnerIds = append(projectBo.OwnerIds, projectBo.Owner)
	projectBo.OwnerIds = slice.SliceUniqueInt64(projectBo.OwnerIds)

	return projectBo, nil
}

// 查询无码应用管理员角色对应的成员 ids (不包含普通管理员、系统管理员、组织拥有者)
func GetProAdminUserIdsForLessCodeApp(orgId, userId, appId int64) ([]int64, errs.SystemErrorInfo) {
	adminUserIds := make([]int64, 0)
	resp := permissionfacade.GetUserGroupMappings(orgId, userId, appId)
	if resp.Failure() {
		log.Errorf("[GetProAdminUserIdsForLessCodeApp] error:%v, orgId:%d, userId:%d, appId:%d",
			resp.Error(), orgId, userId, appId)
		return nil, resp.Error()
	}
	// 从 user list 找到管理员角色
	userGroupMappings := resp.Data.UserGroupMappings
	for uIdStr, group := range userGroupMappings {
		uId, err := strconv.ParseInt(uIdStr, 10, 64)
		if err != nil {
			log.Errorf("[GetProAdminUserIdsForLessCodeApp] error:%v, orgId:%d, userId:%d, appId:%d",
				resp.Error(), orgId, userId, appId)
			return nil, errs.TypeConvertError
		}
		for _, role := range group {
			if role.LangCode == consts.AppProjectPermissionAdmin {
				adminUserIds = append(adminUserIds, uId)
				break
			}
		}
	}

	return adminUserIds, nil
}

func GetProjectBoList(orgId int64, ids []int64) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(ids),
	}, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.ProjectBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return *bos, nil
}

//通过项目类型langCode获取项目列表
func GetProjectBoListByProjectTypeLangCode(orgId int64, projectTypeLangCode *string) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if projectTypeLangCode != nil {
		projectType, err3 := GetProjectTypeByLangCode(orgId, *projectTypeLangCode)
		if err3 != nil {
			log.Error(err3)
			return nil, err3
		}
		cond[consts.TcProjectTypeId] = projectType.Id
	}

	pos := &[]po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.ProjectBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return *bos, nil
}

func GetProjectInfoByOrgIds(orgIds []int64) ([]bo.ProjectBo, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcOrgId:    db.In(orgIds),
		consts.TcIsFiling: consts.AppIsNotFilling,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}

	pos := &[]po.PpmProProject{}
	err := mysql.SelectAllByCond(consts.TableProject, cond, pos)

	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	bos := &[]bo.ProjectBo{}
	err1 := copyer.Copy(pos, bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return *bos, nil
}

// GetProjectDetailsByOrgId 通过组织获取项目详情列表
func GetProjectDetailsByOrgId(orgId int64) ([]*bo.ProjectDetailBo, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcOrgId:    db.In(orgId),
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	pos := &[]po.PpmProProjectDetail{}
	err := mysql.SelectAllByCond(consts.TableProjectDetail, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	bos := make([]*bo.ProjectDetailBo, 0)
	err1 := copyer.Copy(pos, &bos)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	return bos, nil
}
