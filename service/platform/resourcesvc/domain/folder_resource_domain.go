package domain

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/dao"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/po"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func DeleteMidTable(folderIds []int64, orgId, userId int64, recycleVersionId int64, tx ...sqlbuilder.Tx) ([]int64, errs.SystemErrorInfo) {
	resources := &[]po.PpmResFolderResource{}
	err1 := mysql.SelectAllByCond(consts.TableFolderResource, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    orgId,
		consts.TcFolderId: db.In(folderIds),
	}, resources)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.MysqlOperateError
	}

	upd := mysql.Upd{}
	upd[consts.TcIsDelete] = consts.AppIsDeleted
	upd[consts.TcUpdator] = userId
	upd[consts.TcVersion] = recycleVersionId
	upd[consts.TcUpdateTime] = time.Now()
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcFolderId: db.In(folderIds),
		consts.TcOrgId:    orgId,
	}
	err := dao.UpdateMidTableByCond(cond, upd, tx...)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	res := []int64{}
	for _, resource := range *resources {
		res = append(res, resource.ResourceId)
	}

	return res, nil
}

// 通过resourceId查找folderId是否被删除
func CheckResourceFolderIsDelete(orgId, projectId, resourceId int64, versionId int) bool {
	pos := po.PpmResFolder{}
	cond := db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcProjectId: projectId,
		consts.TcIsDelete:  consts.AppIsDeleted,
		consts.TcId:        db.In(db.Raw("select folder_id from ppm_res_folder_resource where is_delete = 1 and org_id = ? and resource_id = ? and version = ?", orgId, resourceId, versionId)),
	}
	err := mysql.SelectOneByCond(consts.TableFolder, cond, &pos)
	if err != nil {
		return false
	}
	return true
}
