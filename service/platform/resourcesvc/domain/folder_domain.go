package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/dao"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/po"
	"strings"
	"time"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func CreateFolder(input bo.CreateFolderBo) (int64, errs.SystemErrorInfo) {
	folderId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableFolder)
	if err != nil {
		log.Error(err)
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	folderPo := po.PpmResFolder{
		Id:        folderId,
		OrgId:     input.OrgId,
		ProjectId: input.ProjectId,
		Name:      input.Name,
		ParentId:  input.ParentId,
		FileType:  input.FileType,
		Creator:   input.UserId,
		Updator:   input.UserId,
	}
	err0 := dao.InsertFolder(folderPo)
	if err0 != nil {
		log.Error(err0)
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err0)
	}
	return folderId, nil
}

func UpdateFolder(folderId int64, input bo.UpdateFolderBo, tx ...sqlbuilder.Tx) (mysql.Upd, errs.SystemErrorInfo) {
	upd := mysql.Upd{}
	if util.FieldInUpdate(input.UpdateFields, "parentId") && input.ParentID != nil {
		err := CheckFolderIds([]int64{*input.ParentID}, input.ProjectID, input.OrgId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		upd[consts.TcParentId] = *input.ParentID
	} else if util.FieldInUpdate(input.UpdateFields, "name") && input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		isNameRight := format.VerifyFolderNameFormat(name)
		if !isNameRight {
			return nil,errs.InvalidFolderNameError
		}
		upd[consts.TcName] = name
	}
	if len(upd) != 0 {
		upd[consts.TcUpdator] = input.UserId
		upd[consts.TcUpdateTime] = time.Now()
		cond := db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcId:       folderId,
		}
		err := dao.UpdateFolderByCond(cond, upd)
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	return upd, nil
}

func CheckFolderIds(folderIds []int64, projectId, orgId int64) errs.SystemErrorInfo {
	isExist, err := dao.FolderIdIsExist(folderIds, projectId, orgId)
	if err != nil {
		log.Error(err)
		return err
	}
	if !isExist {
		log.Error(errs.InvalidFolderIdsError)
		return errs.InvalidFolderIdsError
	}
	return nil
}
func DeleteFolder(folderIds []int64, userId int64, recycleVerisonId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	upd := mysql.Upd{}
	upd[consts.TcIsDelete] = consts.AppIsDeleted
	upd[consts.TcUpdator] = userId
	upd[consts.TcVersion] = recycleVerisonId
	upd[consts.TcUpdateTime] = time.Now()
	cond := db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcId:       db.In(folderIds),
	}
	err := dao.UpdateFolderByCond(cond, upd, tx...)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func GetFolder(parentId *int64, projectId int64, page bo.PageBo) (*[]po.PpmResFolder, uint64, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcIsDelete:  consts.AppIsNoDelete,
		consts.TcProjectId: projectId,
	}
	if parentId != nil {
		cond[consts.TcParentId] = *parentId
	}
	folderPos, total, err := dao.SelectFolderByPage(cond, page)
	if err != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return folderPos, total, nil
}

func GetFolderById(folderIds []int64) ([]bo.FolderBo, errs.SystemErrorInfo) {
	resourceEntities := &[]po.PpmResFolder{}
	err := mysql.SelectAllByCond((&po.PpmResFolder{}).TableName(), db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcId:       db.In(folderIds),
	}, resourceEntities)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.FolderBo{}
	_ = copyer.Copy(resourceEntities, bos)

	return *bos, nil
}

func RecoverFolder(orgId, userId, projectId, folderId int64, recycleId int64) (*bo.FolderBo, errs.SystemErrorInfo) {
	info := &po.PpmResFolder{}
	err := mysql.SelectOneByCond(consts.TableFolder, db.Cond{
		consts.TcIsDelete:consts.AppIsDeleted,
		consts.TcProjectId:projectId,
		consts.TcId:folderId,
	}, info)

	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.RecycleObjectNotExist
		} else {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
	}

	//取文件夹包含的所有文件夹
	allRelateFolderIds, relationErr := GetAllRelateFolderIds(orgId, projectId, []int64{folderId}, recycleId)
	if relationErr != nil {
		log.Error(relationErr)
		return nil, relationErr
	}
	transErr := mysql.TransX(func(tx sqlbuilder.Tx) error {
		_, err := mysql.TransUpdateSmartWithCond(tx, consts.TableFolder, db.Cond{
			consts.TcId:db.In(allRelateFolderIds),
			consts.TcVersion:recycleId,
		}, mysql.Upd{
			consts.TcIsDelete:consts.AppIsNoDelete,
			consts.TcUpdator:userId,
		})
		if err != nil {
			log.Error(err)
			return err
		}

		_, err2 := mysql.TransUpdateSmartWithCond(tx, consts.TableFolderResource, db.Cond{
			consts.TcOrgId:orgId,
			consts.TcFolderId:db.In(allRelateFolderIds),
			consts.TcVersion:recycleId,
		}, mysql.Upd{
			consts.TcIsDelete:consts.AppIsNoDelete,
			consts.TcUpdator:userId,
		})
		if err2 != nil {
			log.Error(err2)
			return err2
		}

		//查找文件夹里的文件，恢复文件关联关系
		resourceFolder := &[]po.PpmResFolderResource{}
		resErr := mysql.SelectAllByCond(consts.TableFolderResource, db.Cond{
			consts.TcOrgId:orgId,
			consts.TcFolderId:db.In(allRelateFolderIds),
			consts.TcVersion:recycleId,
		}, resourceFolder)
		if resErr != nil {
			log.Error(resErr)
			return resErr
		}

		resourceIds := []int64{}
		for _, resource := range *resourceFolder {
			resourceIds = append(resourceIds, resource.ResourceId)
		}
		if len(resourceIds) > 0 {
			_, err3 := mysql.TransUpdateSmartWithCond(tx, consts.TableResourceRelation, db.Cond{
				consts.TcOrgId:orgId,
				consts.TcProjectId:projectId,
				consts.TcResourceId:db.In(resourceIds),
				consts.TcVersion:recycleId,
			}, mysql.Upd{
				consts.TcIsDelete:consts.AppIsNoDelete,
				consts.TcUpdator:userId,
			})
			if err3 != nil {
				log.Error(err3)
				return err3
			}
		}

		return nil
	})

	if transErr != nil {
		log.Error(transErr)
		return nil, errs.MysqlOperateError
	}
	bo := &bo.FolderBo{}
	_ = copyer.Copy(info, bo)

	return bo, nil
}

//recycleId 如果传了版本id，表示是查询已删除的文件夹
func GetAllRelateFolderIds(orgId, project int64, folderIds []int64, recycleId int64) ([]int64, errs.SystemErrorInfo) {
	treeMap, err := GetFolderTree(orgId, project, recycleId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	res := []int64{}
	for _, id := range folderIds {
		currTree, ok := treeMap[id]
		if !ok {
			continue
		}
		lookGet(currTree, &res)
	}

	return slice.SliceUniqueInt64(res), nil
}

func lookGet(tree *bo.FolderTreeBo, result *[]int64) {
	if tree.Id != 0 {
		*result = append(*result, tree.Id)
	}

	for _, child := range tree.Children {
		lookGet(child, result)
	}
}

//获取某些目录下的所有目录
func GetFolderTree(orgId, projectId int64, recycleId int64) (map[int64]*bo.FolderTreeBo, errs.SystemErrorInfo) {
	pos := &[]po.PpmResFolder{}
	cond := db.Cond{
		consts.TcOrgId:orgId,
		consts.TcProjectId:projectId,
		consts.TcIsDelete:consts.AppIsNoDelete,
	}
	if recycleId != 0 {
		//如果传了版本id，表示是查询已删除的文件夹
		cond[consts.TcVersion] = recycleId
		cond[consts.TcIsDelete] = consts.AppIsDeleted
	}
	err := mysql.SelectAllByCond(consts.TableFolder, cond, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	treeMap := map[int64]*bo.FolderTreeBo{}
	treeMap[0] = &bo.FolderTreeBo{
		Id:       0,
		ParentId: -1,
	}

	for _, folder := range *pos {
		treeMap[folder.Id] = &bo.FolderTreeBo{
			Id:       folder.Id,
			ParentId: folder.ParentId,
		}
	}

	for _, treeBo := range treeMap {
		parent, ok := treeMap[treeBo.ParentId]
		if ok {
			parent.Children = append(parent.Children, treeBo)
		}
	}

	return treeMap, nil
}
