package domain

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/cache"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	sconsts "github.com/star-table/polaris-backend/service/platform/projectsvc/consts"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

func GetProjectAuthBoBatch(orgId int64, projectIds []int64) ([]bo.ProjectAuthBo, errs.SystemErrorInfo) {
	keys := make([]interface{}, len(projectIds))
	for i, projectId := range projectIds {
		key, _ := util.ParseCacheKey(sconsts.CacheBaseProjectInfo, map[string]interface{}{
			consts.CacheKeyOrgIdConstName:     orgId,
			consts.CacheKeyProjectIdConstName: projectId,
		})
		keys[i] = key
	}
	resultList := make([]string, 0)
	if len(keys) > 0 {
		list, err := cache.MGet(keys...)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}
		resultList = list
	}
	projectAuthBos := make([]bo.ProjectAuthBo, 0)
	validProjectIds := make([]int64, 0)
	for _, projectInfoJson := range resultList {
		projectCacheInfo := &bo.ProjectAuthBo{}

		err := json.FromJson(projectInfoJson, projectCacheInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		projectAuthBos = append(projectAuthBos, *projectCacheInfo)
		validProjectIds = append(validProjectIds, projectCacheInfo.Id)
	}
	//找不存在的
	if len(projectIds) != len(validProjectIds) {
		for _, projectId := range projectIds {
			exist, _ := slice.Contain(validProjectIds, projectId)
			if !exist {
				if projectId == 0 {
					continue
				}
				projectAuthBo, err := LoadProjectAuthBo(orgId, projectId)
				if err != nil {
					log.Errorf("[GetProjectAuthBoBatch] orgId:%v, projectId:%v, error:%v", orgId, projectId, err)
					continue
				}
				projectAuthBos = append(projectAuthBos, *projectAuthBo)
			}
		}
	}
	return projectAuthBos, nil
}

func LoadProjectAuthBo(orgId int64, projectId int64) (*bo.ProjectAuthBo, errs.SystemErrorInfo) {
	key, projectInfoJson, err5 := getCacheProjectInfo(orgId, projectId)

	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}

	if projectInfoJson != "" {
		projectCacheInfo := &bo.ProjectAuthBo{}
		err := json.FromJson(projectInfoJson, projectCacheInfo)
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}

		if projectCacheInfo.AppId == int64(0) {
			//补偿获取appId
			project := &po.PpmProProject{}
			err := mysql.SelectOneByCond(project.TableName(), db.Cond{
				consts.TcId:       projectId,
				consts.TcOrgId:    orgId,
				consts.TcIsDelete: consts.AppIsNoDelete,
			}, project)
			if err != nil {
				log.Error(err)
				return nil, errs.ProjectNotExist
			}
			projectCacheInfo.AppId = project.AppId
		}
		return projectCacheInfo, nil
	} else {
		if projectId == 0 {
			return nil, errs.ProjectNotExist
		}
		project := &po.PpmProProject{}
		err := mysql.SelectOneByCond(project.TableName(), db.Cond{
			consts.TcId:       projectId,
			consts.TcOrgId:    orgId,
			consts.TcIsDelete: consts.AppIsNoDelete,
		}, project)
		if err != nil {
			log.Errorf("[LoadProjectAuthBo] orgId:%v, projectId:%v, error:%v", orgId, projectId, err)
			return nil, errs.ProjectNotExist
		}

		projectMemberInfos := &[]po.PpmProProjectRelation{}
		err = mysql.SelectAllByCond(consts.TableProjectRelation, db.Cond{
			consts.TcProjectId:    projectId,
			consts.TcOrgId:        orgId,
			consts.TcStatus:       consts.AppStatusEnable,
			consts.TcIsDelete:     consts.AppIsNoDelete,
			consts.TcRelationType: db.In([]int{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeDepartmentParticipant, consts.IssueRelationTypeOwner}),
		}, projectMemberInfos)

		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
		}
		participants := &[]int64{}
		followers := &[]int64{}
		participantDepartments := &[]int64{}
		ownerIds := &[]int64{}

		dealProjectMemberInfo(projectMemberInfos, followers, participants, participantDepartments, ownerIds)

		projectCacheInfo := &bo.ProjectAuthBo{
			Id:                     project.Id,
			AppId:                  project.AppId,
			Name:                   project.Name,
			Creator:                project.Creator,
			Owner:                  project.Owner,
			Participants:           *participants,
			ParticipantDepartments: *participantDepartments,
			Followers:              *followers,
			Status:                 project.Status,
			IsFilling:              project.IsFiling,
			PublicStatus:           project.PublicStatus,
			ProjectType:            project.ProjectTypeId,
			OwnerIds:               *ownerIds,
		}

		projectInfoJson, err = json.ToJson(projectCacheInfo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		err = cache.SetEx(key, projectInfoJson, consts.GetCacheBaseExpire())
		if err != nil {
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
		}
		return projectCacheInfo, nil
	}
}

func dealProjectMemberInfo(projectMemberInfos *[]po.PpmProProjectRelation, followers, participants, participantDepartments, ownerIds *[]int64) {
	for _, projectMemberInfo := range *projectMemberInfos {
		//if projectMemberInfo.RelationType == consts.IssueRelationTypeFollower {
		//	*followers = append(*followers, projectMemberInfo.RelationId)
		if projectMemberInfo.RelationType == consts.IssueRelationTypeParticipant {
			*participants = append(*participants, projectMemberInfo.RelationId)
		} else if projectMemberInfo.RelationType == consts.IssueRelationTypeDepartmentParticipant {
			*participantDepartments = append(*participantDepartments, projectMemberInfo.RelationId)
		} else if projectMemberInfo.RelationType == consts.IssueRelationTypeOwner {
			*ownerIds = append(*ownerIds, projectMemberInfo.RelationId)
		}
	}
}

func getCacheProjectInfo(orgId, projectId int64) (string, string, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseProjectInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return "", "", err5
	}

	projectInfoJson, err := cache.Get(key)
	if err != nil {
		return key, "", errs.BuildSystemErrorInfo(errs.RedisOperateError, err)
	}

	return key, projectInfoJson, nil
}

func RefreshProjectAuthBo(orgId int64, projectId int64) error {
	key, err5 := util.ParseCacheKey(sconsts.CacheBaseProjectInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return err5
	}
	mod, err := cache.Del(key)
	if err != nil {
		log.Errorf("RefreshProjectInfo 发生异常 %v", err)
		return err
	}
	if mod == 0 {
		log.Info("RefreshProjectInfo 刷新失败，key不存在")
	}

	return nil
}

func GetProjectCalendarInfo(orgId, projectId int64) (*bo.CacheProjectCalendarInfoBo, errs.SystemErrorInfo) {
	key, err5 := util.ParseCacheKey(sconsts.CacheProjectCalendarInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})
	if err5 != nil {
		log.Error(err5)
		return nil, err5
	}
	infoJson, err := cache.Get(key)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}
	infoBo := &bo.CacheProjectCalendarInfoBo{}
	if infoJson != "" {
		err = json.FromJson(infoJson, infoBo)
		if err != nil {
			log.Error(err)
			return nil, errs.BuildSystemErrorInfo(errs.JSONConvertError)
		}
		return infoBo, nil
	} else {
		//查看项目是否支持导出日历
		detail, err := GetProjectDetailByProjectIdBo(projectId, orgId)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		//查看项目日历对应的id
		relation, err := GetProjectRelationByType(projectId, []int64{consts.IssueRelationTypeCalendar})
		if err != nil {
			log.Error(err)
			return nil, err
		}
		var calendarId string
		for _, v := range *relation {
			calendarId = v.RelationCode
		}
		infoBo.IsSyncOutCalendar = detail.IsSyncOutCalendar
		infoBo.CalendarId = calendarId
		infoBo.Creator = detail.Creator

		setErr := cache.SetEx(key, json.ToJsonIgnoreError(infoBo), consts.GetCacheBaseExpire())
		if setErr != nil {
			log.Error(setErr)
			return nil, errs.BuildSystemErrorInfo(errs.RedisOperateError)
		}

		return infoBo, nil
	}
}

func DeleteProjectCalendarInfo(orgId, projectId int64) errs.SystemErrorInfo {
	key, err := util.ParseCacheKey(sconsts.CacheProjectCalendarInfo, map[string]interface{}{
		consts.CacheKeyOrgIdConstName:     orgId,
		consts.CacheKeyProjectIdConstName: projectId,
	})

	if err != nil {
		log.Error(err)
		return err
	}

	_, err1 := cache.Del(key)

	if err1 != nil {
		log.Error(err1)
		return errs.BuildSystemErrorInfo(errs.RedisOperateError)
	}

	return nil
}
