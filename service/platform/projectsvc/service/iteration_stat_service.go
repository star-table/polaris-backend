package service

import (
	"time"

	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
	"upper.io/db.v3"
)

func IterationStats(orgId int64, page uint, size uint, params vo.IterationStatReq) (*vo.IterationStatList, errs.SystemErrorInfo) {
	iterationId := params.IterationID

	iterationBo, err1 := domain.GetIterationBoByOrgId(iterationId, orgId)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	cond := db.Cond{
		consts.TcIterationId: iterationBo.Id,
		consts.TcIsDelete:    consts.AppIsNoDelete,
	}
	if params.StartDate == nil && params.EndDate == nil {
		//默认查询十五天
		currentTime := time.Now()
		startDate := types.Time(currentTime.AddDate(0, 0, -14))
		params.StartDate = &startDate
	}
	if params.StartDate != nil && params.EndDate != nil {
		startTime := time.Time(*params.StartDate)
		endTime := time.Time(*params.EndDate)
		if startTime.Before(endTime) {
			startDate := startTime.Format(consts.AppDateFormat)
			endDate := endTime.Format(consts.AppDateFormat)
			cond[consts.TcStatDate] = db.Between(startDate, endDate)
		}
	} else if params.StartDate != nil {
		startDate := time.Time(*params.StartDate).Format(consts.AppDateFormat)
		cond[consts.TcStatDate] = db.Gte(startDate)
	} else if params.EndDate != nil {
		endDate := time.Time(*params.EndDate).Format(consts.AppDateFormat)
		cond[consts.TcStatDate] = db.Lte(endDate)
	}

	bos, total, err1 := domain.GetIterationStatBoList(page, size, cond)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.BaseDomainError, err1)
	}

	resultList := &[]*vo.IterationStat{}
	err3 := copyer.Copy(bos, resultList)
	if err3 != nil {
		log.Error(err3)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}

	return &vo.IterationStatList{
		Total: total,
		List:  *resultList,
	}, nil

}

func AppendIterationStat(iterationBo bo.IterationBo, date string) errs.SystemErrorInfo {
	return domain.AppendIterationStat(iterationBo, date)
}
