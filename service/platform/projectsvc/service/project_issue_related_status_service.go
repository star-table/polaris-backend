package service

import (
	"strconv"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	lang2 "github.com/star-table/polaris-backend/common/core/util/lang"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

func ProjectIssueRelatedStatus(orgId int64, input vo.ProjectIssueRelatedStatusReq) ([]*vo.HomeIssueStatusInfo, errs.SystemErrorInfo) {
	projectId := input.ProjectID
	// 表头改造后，参数 ProjectObjectTypeID 变为 tableId
	tableIdStr := input.TableID
	tableId, oriErr := strconv.ParseInt(tableIdStr, 10, 64)
	if oriErr != nil {
		log.Errorf("[ProjectIssueRelatedStatus] parse tableIdStr err: %v", oriErr)
		return nil, errs.BuildSystemErrorInfo(errs.TypeConvertError, oriErr)
	}

	homeStatusInfoBos, err := domain.GetProjectRelatedStatus(orgId, projectId, tableId)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	results := &[]*vo.HomeIssueStatusInfo{}
	err1 := copyer.Copy(homeStatusInfoBos, results)
	if err1 != nil {
		log.Error(err1)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError)
	}
	lang := lang2.GetLang()
	isOtherLang := lang2.IsEnglish()
	otherLanguageMap := make(map[string]string, 0)
	if tmpMap, ok1 := consts.LANG_STATUS_MAP[lang]; ok1 {
		otherLanguageMap = tmpMap
	}
	if isOtherLang {
		for index, item := range *results {
			if tmpVal, ok2 := otherLanguageMap[item.Name]; ok2 {
				(*results)[index].Name = tmpVal
			}
		}
	}
	return *results, nil
}
