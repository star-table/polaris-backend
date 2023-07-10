package domain

import (
	"strconv"
	"time"

	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/core/util/date"
	"github.com/star-table/common/core/util/encrypt"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/md5"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/facade/idfacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"upper.io/db.v3"
)

//获取我参与的和我负责的任务
func GetRelatedIssues(currentUserId, orgId int64) ([]bo.IssueRelationBo, error) {
	issueRelationCond := db.Cond{}
	issueRelationCond[consts.TcRelationId] = currentUserId
	issueRelationCond[consts.TcOrgId] = orgId
	issueRelationCond[consts.TcRelationType] = db.In([]int64{consts.IssueRelationTypeParticipant, consts.IssueRelationTypeOwner})
	issueRelationCond[consts.TcIsDelete] = consts.AppIsNoDelete
	issueRelations := &[]po.PpmPriIssueRelation{}
	err := mysql.SelectAllByCond(consts.TableIssueRelation, issueRelationCond, issueRelations)
	if err != nil {
		return nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	issueRelationBos := &[]bo.IssueRelationBo{}
	_ = copyer.Copy(issueRelations, issueRelationBos)
	return *issueRelationBos, nil
}

func GetIssueReport(id string) (bo.ShareBo, error) {
	shareInfo := &po.PpmShaShare{}
	err := mysql.SelectOneByCond(consts.TableShare, db.Cond{
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   1,
		consts.TcId:       id,
	}, shareInfo)
	shareBo := &bo.ShareBo{}
	if err != nil {
		return *shareBo, err
	}
	_ = copyer.Copy(shareInfo, shareBo)

	return *shareBo, nil
}

func InsertIssueReport(orgId, currentUserId, total int64, startTime, endTime string, homeIssueBos []bo.HomeIssueInfoBo) (bo.IssueReportDetailBo, error) {
	resp := &bo.IssueReportDetailBo{}
	homeIssueVos := &[]bo.HomeIssueInfoBo{}
	err := copyer.Copy(homeIssueBos, homeIssueVos)
	if err != nil {
		log.Error(err)
		return *resp, errs.BuildSystemErrorInfo(errs.ObjectCopyError, err)
	}
	userInfo, err := orgfacade.GetBaseUserInfoRelaxed(orgId, currentUserId)
	if err != nil {
		return *resp, errs.BuildSystemErrorInfo(errs.GetUserInfoError, err)
	}
	resp = &bo.IssueReportDetailBo{
		Total:          int64(total),
		ReportUserName: userInfo.Name,
		StartTime:      startTime,
		EndTime:        endTime,
		List:           *homeIssueVos,
	}

	listJson, err := json.ToJson(resp)
	if err != nil {
		return *resp, errs.BuildSystemErrorInfo(errs.JSONConvertError, err)
	}
	contentMd5 := md5.Md5V(listJson)

	shareInfo := &po.PpmShaShare{}
	//一天内防止重复插入
	lastDay := (time.Now()).AddDate(0, 0, -1)
	err = mysql.SelectOneByCond(shareInfo.TableName(), db.Cond{
		consts.TcIsDelete:   consts.AppIsNoDelete,
		consts.TcStatus:     1,
		consts.TcContentMd5: contentMd5,
		consts.TcCreator:    currentUserId,
		consts.TcOrgId:      orgId,
		consts.TcCreateTime: db.Gte(date.Format(lastDay)),
	}, shareInfo)

	if err == nil {
		resp.ShareID, _ = encrypt.AesEncrypt(strconv.FormatInt(shareInfo.Id, 10))
	} else {
		tableShareId, err := idfacade.ApplyPrimaryIdRelaxed(shareInfo.TableName())
		if err != nil {
			return *resp, errs.BuildSystemErrorInfo(errs.ApplyIdError, err)
		}
		shareInsert := &po.PpmShaShare{
			Id:         tableShareId,
			OrgId:      orgId,
			Creator:    currentUserId,
			Content:    listJson,
			ContentMd5: contentMd5,
			CreateTime: time.Now(),
			Updator:    currentUserId,
			UpdateTime: time.Now(),
		}
		err2 := mysql.Insert(shareInsert)
		if err2 != nil {
			return *resp, errs.BuildSystemErrorInfo(errs.ApplyIdError, err2)
		}
		resp.ShareID, _ = encrypt.AesEncrypt(strconv.FormatInt(tableShareId, 10))
	}

	return *resp, nil
}
