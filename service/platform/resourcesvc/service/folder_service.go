package service

import (
	"strings"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/domain"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func CreateFolder(input bo.CreateFolderBo) (*vo.Void, errs.SystemErrorInfo) {
	orgId := input.OrgId
	projectId := input.ProjectId
	parentId := input.ParentId
	input.Name = strings.TrimSpace(input.Name)
	isNameRight := format.VerifyFolderNameFormat(input.Name)
	if !isNameRight {
		return nil, errs.InvalidFolderNameError
	}
	err := domain.CheckFolderIds([]int64{parentId}, projectId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	folderId, err := domain.CreateFolder(input)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	return &vo.Void{ID: folderId}, nil
}

func DeleteFolder(input bo.DeleteFolderBo) (*resourcevo.DeleteFolderData, errs.SystemErrorInfo) {
	orgId := input.OrgId
	userId := input.UserId
	folderIds := input.FolderIds
	projectId := input.ProjectId
	err := domain.CheckFolderIds(folderIds, projectId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	oldBos, err := domain.GetFolderById(folderIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	folderNames := []string{}
	for _, value := range oldBos {
		folderNames = append(folderNames, value.Name)
	}
	resp1 := projectfacade.AddRecycleRecord(projectvo.AddRecycleRecordReqVo{
		UserId: userId,
		OrgId:  orgId,
		Input: projectvo.AddRecycleRecordDataVo{
			ProjectId:    projectId,
			RelationIds:  folderIds,
			RelationType: consts.RecycleTypeFolder,
		},
	})
	if resp1.Failure() {
		log.Error(resp1.Error())
		return nil, resp1.Error()
	}
	//真正删除要先获取所有的文件夹及子文件夹
	allRelateFolderIds, relationErr := domain.GetAllRelateFolderIds(orgId, projectId, folderIds, 0)
	if relationErr != nil {
		log.Error(relationErr)
		return nil, relationErr
	}
	err0 := mysql.TransX(func(tx sqlbuilder.Tx) error {
		err = domain.DeleteFolder(allRelateFolderIds, userId, resp1.Data, tx)
		if err != nil {
			log.Error(err)
			return err
		}
		resourceIds, err := domain.DeleteMidTable(allRelateFolderIds, orgId, userId, resp1.Data, tx)
		if err != nil {
			log.Error(err)
			return err
		}

		if len(resourceIds) > 0 {
			relationErr := domain.DeleteResourceRelation(orgId, userId, projectId, nil, resourceIds, []int{consts.OssPolicyTypeProjectResource}, resp1.Data, "", false, tx)
			if relationErr != nil {
				log.Error(relationErr)
				return relationErr
			}
		}
		return nil
	})
	if err0 != nil {
		log.Error(err0)
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err0)
	}
	resp := &resourcevo.DeleteFolderData{
		FolderIds:   folderIds,
		FolderNames: folderNames,
	}
	return resp, nil
}

func UpdateFolder(input bo.UpdateFolderBo) (*resourcevo.UpdateFolderData, errs.SystemErrorInfo) {
	orgId := input.OrgId
	folderId := input.FolderID
	updateFields := input.UpdateFields
	projectId := input.ProjectID
	if updateFields == nil || len(updateFields) == 0 {
		return nil, errs.UpdateFiledIsEmpty
	}
	resp := &resourcevo.UpdateFolderData{}
	err := domain.CheckFolderIds([]int64{folderId}, projectId, orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	bos, err := domain.GetFolderById([]int64{folderId})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	upd, err := domain.UpdateFolder(folderId, input)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	oldBo := bos[0]
	resp.FolderId = folderId
	resp.FolderName = &oldBo.Name
	resp.UpdateFields = input.UpdateFields
	blank := ""
	if _, ok := upd[consts.TcParentId]; ok {
		if oldBo.ParentId == 0 {
			resp.OldValue = &blank
		} else {
			oldParentBos, err := domain.GetFolderById([]int64{oldBo.ParentId})
			if err != nil {
				log.Error(err)
				return nil, err
			}
			resp.OldValue = &(oldParentBos[0].Name)
		}
		newParentId := upd[consts.TcParentId].(int64)
		if newParentId == 0 {
			resp.NewValue = &blank
		} else {
			newParentBos, err := domain.GetFolderById([]int64{upd[consts.TcParentId].(int64)})
			if err != nil {
				log.Error(err)
				return nil, err
			}
			resp.NewValue = &(newParentBos[0].Name)
		}
	} else if _, ok := upd[consts.TcName]; ok {
		resp.OldValue = &oldBo.Name
		nv := upd[consts.TcName].(string)
		resp.NewValue = &nv
	}
	return resp, nil
}

func GetFolder(input bo.GetFolderBo) (*vo.FolderList, errs.SystemErrorInfo) {
	parentId := input.ParentId
	orgId := input.OrgId
	projectId := input.ProjectId
	var folderPos *[]po.PpmResFolder
	var total uint64
	var err errs.SystemErrorInfo
	if parentId != nil {
		err = domain.CheckFolderIds([]int64{*parentId}, projectId, orgId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		folderPos, total, err = domain.GetFolder(parentId, projectId, bo.PageBo{
			Size:  input.Size,
			Page:  input.Page,
			Order: "id desc"})
		if err != nil {
			log.Error(err)
			return nil, err
		}
	} else {
		folderPos, total, err = domain.GetFolder(nil, projectId, bo.PageBo{
			Size:  -1,
			Page:  1,
			Order: "id desc"})
		if err != nil {
			log.Error(err)
			return nil, err
		}
	}
	folderVos := &[]*vo.Folder{}
	copyErr := copyer.Copy(folderPos, folderVos)
	if copyErr != nil {
		log.Error(copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	creatorIds := make([]int64, 0)
	for _, value := range *folderVos {
		creatorIds = append(creatorIds, value.Creator)
	}
	ownerMap, err := domain.GetBaseUserInfoMap(orgId, creatorIds)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	for i, resourceVo := range *folderVos {
		if ownerInfoInterface, ok := ownerMap[resourceVo.Creator]; ok {
			ownerInfo := ownerInfoInterface.(bo.BaseUserInfoBo)
			(*folderVos)[i].CreatorName = ownerInfo.Name
		} else {
			log.Errorf("用户 %d 信息不存在，组织id %d", resourceVo.Creator, orgId)
		}
	}
	return &vo.FolderList{
		List:  *folderVos,
		Total: int64(total),
	}, nil

}

func RecoverFolder(orgId, userId, projectId, folderId int64, recycleId int64) (*bo.FolderBo, errs.SystemErrorInfo) {
	return domain.RecoverFolder(orgId, userId, projectId, folderId, recycleId)
}

func GetFolderInfoBasic(orgId int64, folderIds []int64) ([]vo.Folder, errs.SystemErrorInfo) {
	info := &[]po.PpmResFolder{}
	err := mysql.SelectAllByCond(consts.TableFolder, db.Cond{
		consts.TcOrgId: orgId,
		consts.TcId:    db.In(folderIds),
	}, info)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	res := &[]vo.Folder{}
	_ = copyer.Copy(info, res)

	return *res, nil
}

func CompleteDeleteFolder(orgId, userId int64, projectId, folderId int64, recycleId int64) errs.SystemErrorInfo {
	//取文件夹包含的所有文件夹
	allRelateFolderIds, relationErr := domain.GetAllRelateFolderIds(orgId, projectId, []int64{folderId}, recycleId)
	if relationErr != nil {
		log.Error(relationErr)
		return relationErr
	}
	//查找文件夹里本次跟随删除的文件
	resourceFolder := &[]po.PpmResFolderResource{}
	resErr := mysql.SelectAllByCond(consts.TableFolderResource, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcFolderId: db.In(allRelateFolderIds),
		consts.TcVersion:  recycleId,
	}, resourceFolder)
	if resErr != nil {
		log.Error(resErr)
		return errs.MysqlOperateError
	}

	resourceIds := []int64{}
	for _, resource := range *resourceFolder {
		resourceIds = append(resourceIds, resource.ResourceId)
	}

	if len(resourceIds) == 0 {
		return nil
	}

	deleteErr := CompleteDeleteResource(orgId, userId, resourceIds)
	if deleteErr != nil {
		log.Error(deleteErr)
		return deleteErr
	}

	return nil
}
