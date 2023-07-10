package domain

import (
	"time"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/slice"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo/ordervo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orderfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
)

//特殊功能配置，暂时用不上了
func GetFunctionConfig(orgIds []int64) ([]bo.PpmOrcFunctionConfig, errs.SystemErrorInfo) {
	pos := &[]po.PpmOrcFunctionConfig{}
	err := mysql.SelectAllByCond(consts.TableOrgFunctionConfig, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcOrgId:    db.In(orgIds),
	}, pos)
	if err != nil {
		log.Error(err)
		return nil, errs.MysqlOperateError
	}

	bos := &[]bo.PpmOrcFunctionConfig{}
	_ = copyer.Copy(pos, bos)
	return *bos, nil
}

func GetOrgPayFunction(orgId int64) ([]ordervo.FunctionLimitObj, errs.SystemErrorInfo) {
	orgConfig, err := GetOrgConfig(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//获取等级对应的功能项
	functionsResp := orderfacade.GetFunctionByLevel(ordervo.FunctionReq{Level: int64(orgConfig.PayLevel)})
	if functionsResp.Failure() {
		log.Error(functionsResp.Error())
		return nil, functionsResp.Error()
	}

	return functionsResp.Data, nil
}

func PayLevelIsExist(level int64) (bool, errs.SystemErrorInfo) {
	baseLevelDatas := make([]orgvo.PayBaseLevelData, 0)
	errJson := json.FromJson(consts.BASE_PAY_LEVEL, &baseLevelDatas)
	if errJson != nil {
		log.Errorf("[PayLevelIsExist] json err:%v, level:%v", errJson, level)
		return false, errs.JSONConvertError
	}
	payLevels := []int64{}
	for _, item := range baseLevelDatas {
		payLevels = append(payLevels, item.Id)
	}

	if ok, err := slice.Contain(payLevels, level); err == nil && ok {
		return true, nil
	}
	return false, nil

	//info := &po.PpmBasPayLevel{}
	//err := mysql.SelectOneByCond(consts.TableBasPayLevel, db.Cond{
	//	consts.TcIsDelete: consts.AppIsNoDelete,
	//	consts.TcId:       level,
	//}, info)
	//if err != nil {
	//	if err == db.ErrNoMoreRows {
	//		return false, nil
	//	}
	//	return false, errs.MysqlOperateError
	//}
	//
	//return true, nil
}

// RedirectPayLevelInfoForPrivate 私有化部署时，将 payLevel 等信息进行重置
func RedirectPayLevelInfoForPrivate(orgConfig *bo.OrgConfigBo) {
	orgConfig.PayLevel = consts.PayLevelPrivateDeploy
	// 私有化部署暂时没有到期时间
	orgConfig.PayEndTime = time.Now().AddDate(10, 0, 0)
}
