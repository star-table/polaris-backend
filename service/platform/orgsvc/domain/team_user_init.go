package domain

import (
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/idvo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
)

func TeamUserInit(orgId int64, teamId int64, userId int64, isRoot bool, tx sqlbuilder.Tx) errs.SystemErrorInfo {
	userTeam := &po.PpmTemUserTeam{}

	err := mysql.SelectOneByCond(userTeam.TableName(), db.Cond{
		consts.TcOrgId:  orgId,
		consts.TcTeamId: teamId,
		consts.TcUserId: userId,
	}, userTeam)

	if err != nil {
		//如果团队和用户关联已存在，判断是否被删除
		if userTeam.IsDelete == consts.AppIsDeleted {
			err := mysql.TransUpdateSmart(tx, userTeam.TableName(), userTeam.Id, mysql.Upd{
				consts.TcIsDelete: consts.AppIsNoDelete,
			})
			if err != nil {
				log.Error(strs.ObjectToString(err))
				return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
			}
		}
		return nil
	} else {
		respVo := idfacade.ApplyPrimaryId(idvo.ApplyPrimaryIdReqVo{Code: userTeam.TableName()})
		if respVo.Failure() {
			log.Error(respVo.Message)
			return respVo.Error()
		}
		id := respVo.Id
		userTeam.Id = id
		userTeam.OrgId = orgId
		userTeam.TeamId = teamId
		userTeam.UserId = userId
		if isRoot {
			userTeam.RelationType = consts.UserTeamRelationTypeLeader
		} else {
			userTeam.RelationType = consts.UserTeamRelationTypeMember
		}
		userTeam.Status = consts.AppStatusEnable
		userTeam.IsDelete = consts.AppIsNoDelete

		err2 := mysql.TransInsert(tx, userTeam)
		if err2 != nil {
			log.Error(strs.ObjectToString(err2))
			return errs.BuildSystemErrorInfo(errs.MysqlOperateError, err2)
		}
		return nil
	}

}
