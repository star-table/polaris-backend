package service

import (
	"github.com/star-table/common/core/util/pinyin"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/domain"
)

func ConvertCode(orgId int64, input vo.ConvertCodeReq) (*vo.ConvertCodeResp, errs.SystemErrorInfo) {
	//产品说转换长度改为10
	code := pinyin.ConvertCodeWithMaxLen(input.Name, 10)
	log.Infof("%s convert to %s", input.Name, code)
	rep, err := domain.CheckCodeRepetition(orgId, code)
	if err != nil {
		log.Error(err)
		return nil, errs.BuildSystemErrorInfo(errs.ProjectDomainError, err)
	}
	if rep {
		return nil, errs.BuildSystemErrorInfo(errs.ProjectPreCodeExist)
	}
	return &vo.ConvertCodeResp{
		Code: code,
	}, nil
}
