package service

import (
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/domain"
	"upper.io/db.v3"
)

var log = logger.GetDefaultLogger()

func IndustryList() (*vo.IndustryListResp, errs.SystemErrorInfo) {

	cond := db.Cond{
		consts.TcIsShow: consts.AppShowEnable,
	}

	bos, err := domain.GetIndustryBoAllList(cond)

	if err != nil {
		return nil, err
	}

	resultList := &[]*vo.IndustryResp{}

	copyErr := copyer.Copy(bos, resultList)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return &vo.IndustryListResp{
		List: *resultList,
	}, nil

}
