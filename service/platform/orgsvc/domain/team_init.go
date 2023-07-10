package domain

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/vo/idvo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

const PpmTemTeamSql = consts.TemplateDirPrefix + "ppm_tem_team.template"

func TeamInit(orgId int64, tx sqlbuilder.Tx) (int64, errs.SystemErrorInfo) {
	team := &po.PpmTemTeam{}

	err := mysql.SelectById(team.TableName(), db.Cond{
		consts.TcOrgId:     orgId,
		consts.TcIsDefault: consts.APPIsDefault,
		consts.TcIsDelete:  consts.AppIsNoDelete,
	}, team)

	//表示默认团队不存在
	if err != nil {
		//当team数量等于0时，初始化默认team
		respVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: team.TableName()})
		if respVo.Failure() {
			log.Error(respVo.Message)
			return 0, respVo.Error()
		}

		contextMap := map[string]interface{}{}
		contextMap["TeamId"] = respVo.Id
		contextMap["OrgId"] = orgId
		contextMap["Name"] = consts.InitDefaultTeamName
		contextMap["NickName"] = consts.InitDefaultTeamNickname
		contextMap["IsDefault"] = consts.APPIsDefault
		contextMap["Status"] = consts.AppStatusEnable
		insertErr := util.ReadAndWrite(PpmTemTeamSql, contextMap, tx)
		if insertErr != nil {
			return 0, errs.BuildSystemErrorInfo(errs.BaseDomainError, insertErr)
		}

		//team.Id = respVo.Id
		//team.OrgId = orgId
		//team.Name = consts.InitDefaultTeamName
		//team.NickName = consts.InitDefaultTeamNickname
		//team.IsDefault = consts.APPIsDefault
		//team.Status = consts.AppStatusEnable
		//team.IsDelete = consts.AppIsNoDelete
		//
		//err2 := mysql.TransInsert(tx, team)
		//if err2 != nil {
		//	log.Error("初始化team过程中，添加team异常", err2)
		//	return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
		//}
		return respVo.Id, nil
	}
	return team.Id, nil
}

func TeamOwnerInit(teamId int64, owner int64, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	team := &po.PpmTemTeam{}
	err := mysql.TransUpdateSmart(tx, team.TableName(), teamId, mysql.Upd{
		consts.TcOwner: owner,
	})
	if err != nil {
		return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}
	return nil
}
