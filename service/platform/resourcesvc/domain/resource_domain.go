package domain

import (
	"fmt"
	"image/jpeg"
	"os"
	"strconv"
	"strings"
	"time"

	"upper.io/db.v3"

	tableV1 "github.com/star-table/interface/golang/table/v1"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/maps"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/core/util/uuid"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/core/util/format"
	"github.com/star-table/polaris-backend/common/core/util/image"
	"github.com/star-table/polaris-backend/common/core/util/str"
	"github.com/star-table/polaris-backend/common/extra/lc_helper"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/datacenter"
	"github.com/star-table/polaris-backend/common/model/vo/formvo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/formfacade"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/tablefacade"
	sconsts "github.com/star-table/polaris-backend/service/platform/resourcesvc/consts"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/dao"
	"github.com/star-table/polaris-backend/service/platform/resourcesvc/po"
	"github.com/disintegration/imaging"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

func CreateResource(createResourceBo bo.CreateResourceBo, tx ...sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	if ok, _ := slice.Contain([]int{consts.LocalResource, consts.OssResource, consts.DingDiskResource}, createResourceBo.Type); !ok {
		return 0, errs.BuildSystemErrorInfo(errs.InvalidResourceType)
	}
	createResourceBo.Name = strings.TrimSpace(createResourceBo.Name)
	isNameRight := format.VerifyResourceNameFormat(createResourceBo.Name)
	if !isNameRight {
		log.Error(json.ToJsonIgnoreError(createResourceBo))
		return 0, errs.InvalidResourceNameError
	}
	//新增folderId相关逻辑,为了保持原有逻辑,所以这里添加if条件分支 2019/12/12
	if createResourceBo.FolderId != nil {
		//if len(createResourceBo.Name) > 15 || createResourceBo.Name == "" {
		//	return 0, errs.InvalidResourceNameError
		//}
		//判断folderId是否存在
		folderIsExist, err := dao.FolderIdIsExist([]int64{*createResourceBo.FolderId}, createResourceBo.ProjectId, createResourceBo.OrgId)
		if err != nil {
			log.Error(err)
			return 0, err
		}
		if !folderIsExist {
			log.Error(errs.FolderIdNotExistError)
			return 0, errs.FolderIdNotExistError
		}
	}
	resourceId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableResource)
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
	}
	host, path := str.UrlParse(createResourceBo.Path)

	suffix := createResourceBo.Suffix
	if suffix == "" {
		suffix = util.ParseFileSuffix(createResourceBo.Name)
	}

	resourceEntity := po.PpmResResource{
		Id:        resourceId,
		OrgId:     createResourceBo.OrgId,
		ProjectId: createResourceBo.ProjectId,
		Host:      host,
		Path:      path,
		Name:      createResourceBo.Name,
		Type:      createResourceBo.Type,
		Suffix:    suffix,
		Bucket:    createResourceBo.Bucket,
		Size:      createResourceBo.Size,
		Md5:       createResourceBo.Md5,
		Creator:   createResourceBo.OperatorId,
		Updator:   createResourceBo.OperatorId,
		IsDelete:  consts.AppIsNoDelete,
	}
	//新增filetype逻辑   2019/12/20
	if createResourceBo.SourceType != nil {
		resourceEntity.SourceType = *createResourceBo.SourceType
	}
	//新增自动检测fileType逻辑 2019/12/30
	suffStr := strings.ToUpper(strings.TrimSpace(suffix))
	if value, ok := consts.FileTypes[suffStr]; ok {
		resourceEntity.FileType = value
	} else {
		resourceEntity.FileType = consts.FileTypeOthers
	}
	compressErr := compressImage(createResourceBo)
	if compressErr != nil {
		log.Error(compressErr)
	}

	err2 := dao.InsertResource(resourceEntity, tx...)
	if err2 != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
	}

	// 插入资源关联(支持先上传附件，再创建任务，此时先跳过创建relation后单独创建)
	if !createResourceBo.SkipCreateRelation {
		resourceRelationId, err := idfacade.ApplyPrimaryIdRelaxed(consts.TableResourceRelation)
		if err != nil {
			return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
		}
		resourceRelationPo := po.PpmResResourceRelation{
			Id:         resourceRelationId,
			OrgId:      createResourceBo.OrgId,
			ProjectId:  createResourceBo.ProjectId,
			IssueId:    createResourceBo.IssueId,
			ResourceId: resourceId,
			Creator:    createResourceBo.OperatorId,
			Updator:    createResourceBo.OperatorId,
		}
		if createResourceBo.SourceType != nil {
			resourceRelationPo.SourceType = *createResourceBo.SourceType
		}
		relationErr := dao.InsertResourceRelation(resourceRelationPo, tx...)
		if relationErr != nil {
			log.Error(relationErr)
			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, relationErr)
		}
	}

	//新增folderId相关逻辑,为了保持原有逻辑,所以这里添加if条件分支 2019/12/12
	if createResourceBo.FolderId != nil {
		midtableId, err1 := idfacade.ApplyPrimaryIdRelaxed(consts.TableFolderResource)
		if err1 != nil {
			log.Error(err1)
			return 0, errs.BuildSystemErrorInfo(errs.ApplyIdError, err1)
		}
		//插入中间表数据
		midtableEntity := po.PpmResFolderResource{
			Id:         midtableId,
			OrgId:      createResourceBo.OrgId,
			ResourceId: resourceId,
			FolderId:   *createResourceBo.FolderId,
			Creator:    createResourceBo.OperatorId,
			Updator:    createResourceBo.OperatorId,
			IsDelete:   consts.AppIsNoDelete,
		}
		err2 := dao.InsertMidTable(midtableEntity, tx...)
		if err2 != nil {
			log.Error(err2)
			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
		}
	}

	cacheErr := ClearCacheResourceSize(createResourceBo.OrgId)
	if cacheErr != nil {
		log.Error(cacheErr)
		return 0, cacheErr
	}
	return resourceId, nil
}

func CheckResourceIds(resourceIds []int64, projectId, orgId int64) errs.SystemErrorInfo {
	isExist, err := dao.ResourceIdIsExist(resourceIds, orgId, projectId)
	if err != nil {
		log.Error(err)
		return err
	}
	if !isExist {
		log.Error(errs.InvalidResourceIdsError)
		return errs.InvalidResourceIdsError
	}
	return nil
}

func CheckRelation(resourceIds []int64, folderId int64, orgId int64) errs.SystemErrorInfo {
	isExist, err := dao.RelationIsExist(resourceIds, folderId, orgId)
	if err != nil {
		return err
	}
	if !isExist {
		return errs.ResouceNotInFolderError
	}
	return nil
}

func DeleteResource(resourceIds []int64, folderId *int64, orgId, appId, userId int64, projectId, issueId int64,
	recycleVersionId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	//仅文件删除关联关系
	if folderId != nil {
		upd := mysql.Upd{}
		upd[consts.TcIsDelete] = consts.AppIsDeleted
		upd[consts.TcUpdator] = userId
		upd[consts.TcVersion] = recycleVersionId
		upd[consts.TcUpdateTime] = time.Now()
		cond := db.Cond{
			consts.TcIsDelete:   consts.AppIsNoDelete,
			consts.TcFolderId:   folderId,
			consts.TcResourceId: db.In(resourceIds),
			consts.TcOrgId:      orgId,
		}
		err := dao.UpdateMidTableByCond(cond, upd, tx)
		if err != nil {
			log.Error(err)
			return err
		}
	}
	//暂时不删除源文件，只删除关联关系
	//upd := mysql.Upd{}
	//upd[consts.TcIsDelete] = consts.AppIsDeleted
	//upd[consts.TcUpdator] = userId
	//upd[consts.TcUpdateTime] = time.Now()
	//cond := db.Cond{
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//	consts.TcId:       db.In(resourceIds),
	//}
	//_, err := dao.UpdateResourceByCond(cond, upd, tx...)
	//if err != nil {
	//	log.Error(err)
	//	return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	//}
	//删除关联关系
	cond1 := db.Cond{
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcResourceId: db.In(resourceIds),
	}
	if projectId != 0 {
		cond1[consts.TcProjectId] = projectId
	}
	if issueId != 0 {
		cond1[consts.TcIssueId] = issueId
	}
	err1 := dao.UpdateResourceRelationByCond(cond1, mysql.Upd{
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcVersion:  recycleVersionId,
		consts.TcUpdator:  userId,
	}, tx)
	if err1 != nil {
		log.Error(err1)
		return errs.MysqlOperateError
	}

	// 如果设置有，则进入回收站，设置下flag
	if recycleVersionId != 0 && appId != 0 {
		cond1[consts.TcIsDelete] = consts.AppIsDeleted
		err1 = recycleResourceToLess(orgId, appId, userId, cond1, tx)
		if err1 != nil {
			return err1
		}
	}

	cacheErr := ClearCacheResourceSize(orgId)
	if cacheErr != nil {
		log.Error(cacheErr)
		return cacheErr
	}

	return nil
}

func GetRelationIssueIdsAndResourceIds(cond db.Cond, tx sqlbuilder.Tx) ([]int64, []int64, errs.SystemErrorInfo) {
	relations, err := dao.SelectResourceRelationByCond(cond, tx)
	if err != nil {
		return nil, nil, errs.MysqlOperateError
	}
	resourceIds := make([]int64, 0, len(relations))
	issueIds := make([]int64, 0, len(relations))
	for _, relation := range relations {
		if relation.IssueId != 0 {
			resourceIds = append(resourceIds, relation.ResourceId)
			issueIds = append(issueIds, relation.IssueId)
		}
	}

	return issueIds, resourceIds, nil
}

func recycleResourceToLess(orgId, appId, userId int64, cond db.Cond, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	issueIds, resourceIds, err := GetRelationIssueIdsAndResourceIds(cond, tx)
	if err != nil {
		return errs.MysqlOperateError
	}
	if len(issueIds) == 0 {
		return nil
	}

	resp := tablefacade.RecycleAttachment(orgId, userId, &tableV1.RecycleAttachmentRequest{AppId: appId, IssueIds: issueIds, ResourceIds: resourceIds})
	if resp.Failure() {
		log.Errorf("[recycleResourceToLess] RecycleAttachment error:%v", resp.Err)
		return errs.SystemError
	}

	return nil
}

func compressImage(createResourceBo bo.CreateResourceBo) errs.SystemErrorInfo {
	if createResourceBo.Type == consts.LocalResource && createResourceBo.DistPath != "" {
		distPath := createResourceBo.DistPath
		suffix := createResourceBo.Suffix
		if _, ok := consts.ImgTypeMap[strings.ToUpper(suffix)]; ok {
			//固定高120
			newImg, err := image.ResizeAuto(distPath, 120, imaging.Lanczos)
			if err != nil {
				log.Error(err)
			} else {
				afterPath := util.GetCompressedPath(distPath, createResourceBo.Type)
				f, err := os.Create(afterPath)
				if err != nil {
					log.Error(err)
				} else {
					defer func() {
						if err := f.Close(); err != nil {
							log.Error(err)
						}
					}()
					imgErr := jpeg.Encode(f, newImg, nil)
					if imgErr != nil {
						log.Error(imgErr)
					}
				}
			}
		}
	}
	return nil
}

// 获取资源信息
func GetResourceByIds(resourceIds []int64) ([]bo.ResourceBo, errs.SystemErrorInfo) {
	resourceEntities := &[]po.PpmResResource{}
	err := mysql.SelectAllByCond((&po.PpmResResource{}).TableName(), db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcId:       db.In(resourceIds),
	}, resourceEntities)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ResourceBo{}
	_ = copyer.Copy(resourceEntities, bos)

	return *bos, nil
}

func GetIdByPath(orgId int64, resourcePath string, resourceType int) (int64, errs.SystemErrorInfo) {
	resourceInfo := &bo.ResourceTypeBo{}
	host, path := str.UrlParse(resourcePath)
	err := mysql.SelectOneByCond(consts.TableResource, db.Cond{
		consts.TcIsDelete: db.Eq(consts.AppIsNoDelete),
		consts.TcPath:     db.Eq(path),
		consts.TcHost:     db.Eq(host),
		consts.TcOrgId:    orgId,
		consts.TcType:     resourceType,
	}, resourceInfo)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return 0, errs.BuildSystemErrorInfo(errs.ResourceNotExist)
		} else {
			return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError)
		}
	}
	return resourceInfo.ID, nil
}

func GetResourceBoList(page uint, size uint, input resourcevo.GetResourceBoListCond) (*[]bo.ResourceBo, int64, errs.SystemErrorInfo) {
	cond := db.Cond{}
	cond[consts.TcOrgId] = input.OrgId
	if input.ResourceIds != nil {
		cond[consts.TcId] = db.In(*input.ResourceIds)
	}
	if input.IsDelete != nil {
		cond[consts.TcIsDelete] = *input.IsDelete
	}
	pos, total, err := dao.SelectResourceByPage(cond, nil, bo.PageBo{
		Page:  int(page),
		Size:  int(size),
		Order: "id desc",
	})
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ResourceBo{}

	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	if len(*bos) == 0 {
		return bos, int64(total), nil

	}
	creatorIds := []int64{}
	for _, resourceBo := range *bos {
		creatorIds = append(creatorIds, resourceBo.Creator)
	}
	baseUserInfos := orgfacade.GetBaseUserInfoBatch(orgvo.GetBaseUserInfoBatchReqVo{
		OrgId:   input.OrgId,
		UserIds: creatorIds,
	})
	if baseUserInfos.Failure() {
		log.Error(baseUserInfos.Error())
		return nil, 0, baseUserInfos.Error()
	}
	userMap := maps.NewMap("UserId", baseUserInfos.BaseUserInfos)
	for i, resourceBo := range *bos {
		if userCacheInfo, ok := userMap[resourceBo.Creator]; ok {
			baseUserInfo := userCacheInfo.(bo.BaseUserInfoBo)
			(*bos)[i].CreatorName = baseUserInfo.Name
		}
	}

	return bos, int64(total), nil
}
func GetResourceBoListByPage(cond db.Cond, union *db.Union, pageBo bo.PageBo) (*[]bo.ResourceBo, int64, errs.SystemErrorInfo) {
	pos, total, err := dao.SelectResourceByPage(cond, union, pageBo)
	if err != nil {
		log.Error(err)
		return nil, 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	bos := &[]bo.ResourceBo{}
	copyErr := copyer.Copy(pos, bos)
	if copyErr != nil {
		return nil, 0, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}
	return bos, int64(total), nil
}

func UpdateResourceInfo(resourceId int64, input bo.UpdateResourceInfoBo, tx ...sqlbuilder.Tx) (*po.PpmResResource, errs.SystemErrorInfo) {
	resourcePo, err := dao.SelectResourceById(resourceId, tx...)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	if util.FieldInUpdate(input.UpdateFields, "fileName") && input.FileName != nil {
		fileName := strings.TrimSpace(*input.FileName)
		isNameRight := format.VerifyResourceNameFormat(fileName)
		if !isNameRight {
			return nil, errs.InvalidResourceNameError
		}
		resourcePo.Name = fileName
	}

	if util.FieldInUpdate(input.UpdateFields, "fileSuffix") && input.FileSuffix != nil {
		resourcePo.Name = *input.FileSuffix
	}

	if util.FieldInUpdate(input.UpdateFields, "fileSize") && input.FileSize > 0 {
		resourcePo.Size = input.FileSize
	}

	resourcePo.Updator = input.UserId
	resourcePo.UpdateTime = time.Now()
	err = dao.UpdateResource(*resourcePo, tx...)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return resourcePo, nil
}

func UpdateResourceFolderId(resourceIds []int64, currentFolderId, targetFolderId, userId, orgId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	//更新关联锁
	lockKey := fmt.Sprintf("%s%d", consts.UpdateResourceFolderLock, targetFolderId)
	uid := uuid.NewUuid()
	suc, lockErr := cache.TryGetDistributedLock(lockKey, uid)
	if lockErr != nil {
		log.Error(lockErr)
		return errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, e := cache.ReleaseDistributedLock(lockKey, uid); e != nil {
				log.Error(e)
			}
		}()
		//查询有没有关联
		hasRelation, err := dao.ResourceFolderHasRelation(resourceIds, targetFolderId, orgId)
		if err != nil {
			log.Error(err)
			return err
		}

		//和目标文件夹已经有关联了就不要再移动了
		if hasRelation {
			return errs.UpdateResourceFolderError
		}

		//修改老的关联
		upd := mysql.Upd{}
		upd[consts.TcFolderId] = targetFolderId
		upd[consts.TcUpdator] = userId
		upd[consts.TcUpdateTime] = time.Now()
		cond := db.Cond{
			consts.TcIsDelete:   consts.AppIsNoDelete,
			consts.TcFolderId:   currentFolderId,
			consts.TcResourceId: db.In(resourceIds),
			consts.TcOrgId:      orgId,
		}
		err = dao.UpdateMidTableByCond(cond, upd, tx...)
		if err != nil {
			log.Error(err)
			return err
		}
	} else {
		return errs.UpdateResourceFolderError
	}
	return nil
}

func GetResourceIdsByFolderId(folderId, orgId int64) (*[]int64, errs.SystemErrorInfo) {
	midtablePos, err := dao.SelectMidTablePoByFolderId(folderId, orgId)
	if err != nil {
		return nil, err
	}
	var resourceIds []int64
	for _, value := range *midtablePos {
		resourceIds = append(resourceIds, value.ResourceId)
	}
	return &resourceIds, nil
}

func GetBaseUserInfoMap(orgId int64, userIds []int64) (map[interface{}]interface{}, errs.SystemErrorInfo) {
	ownerInfos, err := orgfacade.GetBaseUserInfoBatchRelaxed(orgId, userIds)
	if err != nil {
		return nil, err
	}
	ownerMap := maps.NewMap("UserId", ownerInfos)
	return ownerMap, nil
}

func InsertResourceRelation(orgId, projectId, issueId int64, userId int64, resourceIds []int64, resourceIdSourceTypes map[int64]int) errs.SystemErrorInfo {
	//防止重复插入
	uid := uuid.NewUuid()
	projectIdStr := strconv.FormatInt(projectId, 10)
	issueIdStr := strconv.FormatInt(issueId, 10)
	lockKey := consts.AddResourceRelationLock + projectIdStr + ":" + issueIdStr
	suc, err := cache.TryGetDistributedLock(lockKey, uid)
	if err != nil {
		log.Errorf("获取%s锁时异常 %v", lockKey, err)
		return errs.TryDistributedLockError
	}
	if suc {
		defer func() {
			if _, err := cache.ReleaseDistributedLock(lockKey, uid); err != nil {
				log.Error(err)
			}
		}()
	}

	var existRelations []po.PpmResResourceRelation
	err = mysql.SelectAllByCond(consts.TableResourceRelation, db.Cond{
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcOrgId:      orgId,
		consts.TcProjectId:  projectId,
		consts.TcIssueId:    issueId,
		consts.TcResourceId: db.In(resourceIds),
	}, &existRelations)
	if err != nil {
		log.Error(err)
		return errs.MysqlOperateError
	}

	needIds := []int64{}
	if len(existRelations) > 0 {
		existIds := []int64{}
		for _, relation := range existRelations {
			existIds = append(existIds, relation.ResourceId)
		}
		for _, id := range resourceIds {
			if ok, _ := slice.Contain(existIds, id); !ok {
				needIds = append(needIds, id)
			}
		}
	} else {
		needIds = resourceIds
	}
	if len(needIds) == 0 {
		return nil
	}

	ids, idErr := idfacade.ApplyMultiplePrimaryIdRelaxed(consts.TableResourceRelation, len(needIds))
	if idErr != nil {
		log.Error(idErr)
		return idErr
	}

	insertArr := []po.PpmResResourceRelation{}
	for i, id := range needIds {
		insertArr = append(insertArr, po.PpmResResourceRelation{
			Id:         ids.Ids[i].Id,
			OrgId:      orgId,
			ProjectId:  projectId,
			IssueId:    issueId,
			ResourceId: id,
			Creator:    userId,
			Updator:    userId,
			SourceType: resourceIdSourceTypes[id],
		})
	}

	insertErr := mysql.BatchInsert(&po.PpmResResourceRelation{}, slice.ToSlice(insertArr))
	if insertErr != nil {
		log.Error(insertErr)
		return errs.MysqlOperateError
	}

	return nil
}

func GetResourceInfo(orgId, resourceId int64) (*bo.ResourceBo, errs.SystemErrorInfo) {
	info := &po.PpmResResource{}
	err := mysql.SelectOneByCond(consts.TableResource, db.Cond{
		consts.TcId:    resourceId,
		consts.TcOrgId: orgId,
	}, info)

	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.ResourceNotExist
		} else {
			log.Error(err)
			return nil, errs.MysqlOperateError
		}
	}
	bo := &bo.ResourceBo{}
	_ = copyer.Copy(info, bo)
	return bo, nil
}

func CacheResourceSize(orgId int64) (int64, errs.SystemErrorInfo) {
	key := sconsts.CacheOrgResourceSize
	infoJson, err := cache.HGet(key, strconv.FormatInt(orgId, 10))
	if err != nil {
		log.Error(err)
		return 0, errs.RedisOperateError
	}

	if infoJson != "" {
		size, err := strconv.ParseInt(infoJson, 10, 64)
		if err != nil {
			log.Error(err)
			return 0, errs.JSONConvertError
		}
		return size, nil
	} else {
		infoPo := &po.ResourceCount{}
		conn, err1 := mysql.GetConnect()
		if err1 != nil {
			log.Error(err1)
			return 0, errs.MysqlOperateError
		}
		err := conn.Select(db.Raw("sum(size) as total")).From(consts.TableResource).Where(db.Cond{
			consts.TcIsDelete: consts.AppIsNoDelete,
			consts.TcOrgId:    orgId,
			consts.TcType:     consts.OssResource,
			consts.TcId:       db.In(db.Raw("select distinct(resource_id) from ppm_res_resource_relation where org_id = ? and is_delete =2", orgId)),
		}).One(infoPo)
		//err := conn.Select(db.Raw("sum(a.size) as total")).From(db.Raw("(select DISTINCT path, size from ppm_res_resource where type = 2 and is_delete = 2 and org_id = " + strconv.FormatInt(orgId, 10) + ") a")).One(infoPo)
		err = cache.HSet(key, strconv.FormatInt(orgId, 10), strconv.FormatInt(infoPo.Total, 10))
		if err != nil {
			return 0, errs.RedisOperateError
		}

		return infoPo.Total, nil
	}
}

func ClearCacheResourceSize(orgId int64) errs.SystemErrorInfo {
	key := sconsts.CacheOrgResourceSize
	infoJson, err := cache.HGet(key, strconv.FormatInt(orgId, 10))
	if err != nil {
		log.Error(err)
		return errs.RedisOperateError
	}
	if infoJson != "" {
		_, err := cache.HDel(key, strconv.FormatInt(orgId, 10))
		if err != nil {
			log.Error(err)
			return errs.CacheProxyError
		}
	}

	return nil
}

func RecoverLessAttachment(orgId, userId int64, appId int64, recycleId int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	if appId == 0 {
		return nil
	}
	// 查询出columnId和issueId
	cond1 := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsDeleted,
		consts.TcVersion:  recycleId,
	}
	relations, err := dao.SelectResourceRelationByCond(cond1, tx)
	if err != nil {
		log.Errorf("[RecoverLessAttachment] GetResourceRelations err:%v, orgId:%d, recycleId:%d", err, orgId, recycleId)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	if len(relations) < 1 {
		log.Errorf("[RecoverLessAttachment] resource relation not found, orgId:%d, recycleId:%d", orgId, recycleId)
		return nil
	}
	columnIds := []string{}
	issueIds := []int64{}
	resourceIds := []int64{}
	for _, resource := range relations {
		columnIds = append(columnIds, resource.ColumnId)
		issueIds = append(issueIds, resource.IssueId)
		resourceIds = append(resourceIds, resource.ResourceId)
	}

	updateValue := getAttachmentRecycleJson(resourceIds, columnIds, consts.AppIsNoDelete)

	// 把无码附件数据recycleFlag置为2
	conditions := vo.LessCondsData{Type: consts.ConditionAnd}
	conditions.Conds = append(conditions.Conds, &vo.LessCondsData{
		Type:   consts.ConditionIn,
		Values: issueIds,
		Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
	})
	conditions.Conds = append(conditions.Conds, &vo.LessCondsData{
		Type:   consts.ConditionEqual,
		Value:  orgId,
		Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
	})

	sets := []datacenter.Set{
		{
			Column:          consts.LcJsonColumn,
			Value:           updateValue,
			Type:            consts.SetTypeJson,
			Action:          consts.SetActionSet,
			WithoutPretreat: true,
		},
	}
	batchRaw := formfacade.LessUpdateIssueBatchRaw(&formvo.LessUpdateIssueBatchReq{
		OrgId:     orgId,
		AppId:     appId,
		UserId:    userId,
		Condition: conditions,
		Sets:      sets,
	})
	if batchRaw.Failure() {
		log.Errorf("[RecoverLessAttachment] LessUpdateIssueBatchRaw error:%v, orgId:%d, issueIds:%v",
			batchRaw.Error(), orgId, issueIds)
		return batchRaw.Error()
	}

	return nil
}

func RecoverLessAttachments(orgId, userId int64, projectId, appId, recycleVersionId int64, sourceChannel string, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	if appId == 0 {
		return nil
	}
	cond1 := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcVersion:  recycleVersionId,
	}
	relations, err := dao.SelectResourceRelationByCond(cond1, tx)
	if err != nil {
		log.Errorf("[RecoverLessAttachment] GetResourceRelations err:%v, orgId:%d, recycleId:%d", err, orgId, recycleVersionId)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	if len(relations) < 1 {
		log.Errorf("[RecoverLessAttachment] resource relation not found, orgId:%d, recycleId:%d", orgId, recycleVersionId)
		return nil
	}
	issueIds := []int64{}
	columnIds := []string{}
	resourceIds := []int64{}
	filterColumns := []string{
		lc_helper.ConvertToFilterColumn(consts.BasicFieldIssueId),
		lc_helper.ConvertToFilterColumn(consts.BasicFieldId),
	}
	resourceColumnId := map[int64]string{}
	for _, resource := range relations {
		columnIds = append(columnIds, resource.ColumnId)
		issueIds = append(issueIds, resource.IssueId)
		resourceIds = append(resourceIds, resource.ResourceId)
		resourceColumnId[resource.ResourceId] = resource.ColumnId
		filterColumns = append(filterColumns, lc_helper.ConvertToFilterColumn(resource.ColumnId))
	}

	condition := &tableV1.Condition{
		Type: tableV1.ConditionType_and,
		Conditions: []*tableV1.Condition{
			&tableV1.Condition{
				Type:   tableV1.ConditionType_equal,
				Value:  json.ToJsonIgnoreError([]interface{}{orgId}),
				Column: lc_helper.ConvertToCondColumn(consts.BasicFieldOrgId),
			},
			&tableV1.Condition{
				Type:   tableV1.ConditionType_in,
				Value:  json.ToJsonIgnoreError([]interface{}{issueIds}),
				Column: lc_helper.ConvertToCondColumn(consts.BasicFieldIssueId),
			},
		},
	}

	issueMapList := projectfacade.GetIssueRowList(projectvo.IssueRowListReq{
		OrgId:  orgId,
		UserId: userId,
		Input: &tableV1.ListRawRequest{
			DbType:        tableV1.DbType_slave1,
			FilterColumns: filterColumns,
			Condition:     condition,
		},
	})
	if issueMapList.Failure() {
		log.Errorf("[RecoverLessAttachment] GetIssueRowList err:%v, issueIds:%v", issueMapList.Error(), issueIds)
		return issueMapList.Error()
	}

	if len(issueMapList.Data) < 1 {
		return errs.CanNotRecoverDocuments
	}

	formData := []map[string]interface{}{}
	data := map[string]interface{}{}
	for _, issue := range issueMapList.Data {
		attachmentsMap := map[string]*bo.Attachments{}
		for _, columnId := range columnIds {
			if documentMap, ok := issue[columnId]; ok {
				copyer.Copy(documentMap, &attachmentsMap)
			}
		}

		for _, id := range resourceIds {
			idStr := strconv.FormatInt(id, 10)
			if attach, ok := attachmentsMap[idStr]; ok {
				if attach.RecycleFlag == consts.AppIsDeleted {
					attach.RecycleFlag = consts.AppIsNoDelete
				}
			}
			data[resourceColumnId[id]] = attachmentsMap
		}
		data[consts.BasicFieldId] = issue[consts.BasicFieldIssueId]
		formData = append(formData, data)
	}

	resp := projectfacade.BatchUpdateIssue(&projectvo.BatchUpdateIssueReqVo{
		OrgId:  orgId,
		UserId: userId,
		Input: &projectvo.BatchUpdateIssueInput{
			AppId:     appId,
			ProjectId: projectId,
			TableId:   -1,
			Data:      formData,
		},
	})
	if resp.Failure() {
		log.Errorf("[RecoverLessAttachment]BatchUpdateIssue err:%v", resp.Error())
		return resp.Error()
	}

	return nil
}

func getAttachmentRecycleJson(resourceIds []int64, columnIds []string, recycleFlag int) string {
	updateJsons := make([]string, 0, len(columnIds))
	defaultValue := `jsonb_set(%s, '{%s,%d,recycleFlag}', '%d', false)`
	for _, resourceId := range resourceIds {
		dataValue := defaultValue
		if len(columnIds) == 1 {
			dataValue = fmt.Sprintf(defaultValue, consts.LcJsonColumn, columnIds[0], resourceId, recycleFlag)
		} else {
			for i, columnId := range columnIds {
				if i == len(columnIds)-1 {
					dataValue = fmt.Sprintf(dataValue, fmt.Sprintf(defaultValue, consts.LcJsonColumn, columnId, resourceId, recycleFlag))
				} else if i == 0 {
					dataValue = fmt.Sprintf(defaultValue, "%s", columnId, resourceId, recycleFlag)
				} else {
					dataValue = fmt.Sprintf(defaultValue, dataValue, columnId, resourceId, recycleFlag)
				}
			}
		}
		updateJsons = append(updateJsons, dataValue)
	}

	updateJson := updateJsons[0]
	for i := 1; i < len(updateJsons); i++ {
		updateJson = strings.Replace(updateJson, consts.LcJsonColumn, updateJsons[i], 1)
	}

	return updateJson
}

func GetDingSpaceList(outUserId, unionId string, outOrgId string, nextToken string, size int) ([]*resourcevo.SpaceInfo, errs.SystemErrorInfo) {
	return []*resourcevo.SpaceInfo{}, nil
}

func DingFileListById(orgId int64, outOrgId, outUserId, unionId string, spaceId string, dirId string, size int) ([]*resourcevo.DingFileListData, errs.SystemErrorInfo) {
	return []*resourcevo.DingFileListData{}, nil
}

func GetIssueIdsByResource(orgId, projectId int64, resourceIds []int64) ([]resourcevo.GetIssueIdsByResourceIdsResp, errs.SystemErrorInfo) {
	pos := []po.PpmResResourceRelation{}
	err := mysql.SelectAllByCond(consts.TableResourceRelation, db.Cond{
		consts.TcOrgId:      orgId,
		consts.TcProjectId:  projectId,
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcResourceId: db.In(resourceIds),
	}, &pos)
	if err != nil {
		log.Errorf("[GetIssueIdsByResource] select err:%v", err)
		return nil, errs.MysqlOperateError
	}

	resourceMap := map[int64][]int64{}

	for _, item := range pos {
		issueId := item.IssueId
		resourceId := item.ResourceId
		resourceMap[resourceId] = append(resourceMap[resourceId], issueId)
	}

	resp := []resourcevo.GetIssueIdsByResourceIdsResp{}
	for resourceId, issueIds := range resourceMap {
		resp = append(resp, resourcevo.GetIssueIdsByResourceIdsResp{
			ResourceId: resourceId,
			IssueIds:   issueIds,
		})
	}

	return resp, nil

}
