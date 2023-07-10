package domain

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3/lib/sqlbuilder"
)

func InitProjectStatus(orgId, projectId int64, operatorId int64, tx ...sqlbuilder.Tx) errs.SystemErrorInfo {
	pos := make([]po.PpmProProjectRelation, 0, 10)

	for _, statusBo := range consts.ProjectStatusList {
		id, err3 := idfacade.ApplyPrimaryIdRelaxed(consts.TableProjectRelation)
		if err3 != nil {
			log.Error(err3)
			return errs.BuildSystemErrorInfo(errs.ApplyIdError)
		}
		projectRelationPo := po.PpmProProjectRelation{
			Id:           id,
			OrgId:        orgId,
			ProjectId:    projectId,
			RelationId:   statusBo.ID,
			RelationType: consts.IssueRelationTypeStatus,
			Status:       consts.AppStatusEnable,
			Creator:      operatorId,
			Updator:      operatorId,
			IsDelete:     consts.AppIsNoDelete,
		}
		pos = append(pos, projectRelationPo)
	}

	err2 := dao.InsertProjectRelationBatch(pos, tx...)
	if err2 != nil {
		log.Error(err2)
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError)
	}

	return nil
}
