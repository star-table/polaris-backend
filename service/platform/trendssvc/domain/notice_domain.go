package domain

import (
	"github.com/star-table/common/core/util/copyer"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/trendsvo"
	"github.com/star-table/polaris-backend/service/platform/trendssvc/po"
	"upper.io/db.v3"
)

func UnreadNoticeCount(orgId, userId int64) (uint64, errs.SystemErrorInfo) {
	count, err := mysql.SelectCountByCond(consts.TableNotice, db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcNoticer:  userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
		consts.TcStatus:   consts.NoticeUnReadStatus,
		//@人的notice独立在外
		consts.TcRelationType: db.NotIn([]string{consts.NoticeTypeIssueCommentAtSomebody, consts.NoticeTypeIssueRemarkAtSomebody}),
	})
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return count, nil
}

func CallMeCount(req trendsvo.CallMeCountReqVo) (uint64, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcOrgId:        req.OrgId,
		consts.TcNoticer:      req.UserId,
		consts.TcIsDelete:     consts.AppIsNoDelete,
		consts.TcRelationType: db.In([]string{consts.NoticeTypeIssueCommentAtSomebody, consts.NoticeTypeIssueRemarkAtSomebody}),
	}
	if req.ProjectId != nil && *req.ProjectId != 0 {
		cond[consts.TcProjectId] = *req.ProjectId
	}
	if req.IssueId != nil && *req.IssueId != 0 {
		cond[consts.TcIssueId] = *req.IssueId
	}
	count, err := mysql.SelectCountByCond(consts.TableNotice, cond)
	if err != nil {
		return 0, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	return count, nil
}

func GetNoticeList(orgId, userId int64, page, size int, input *vo.NoticeListReq) (uint64, *[]bo.NoticeBo, errs.SystemErrorInfo) {
	cond := db.Cond{
		consts.TcOrgId:    orgId,
		consts.TcNoticer:  userId,
		consts.TcIsDelete: consts.AppIsNoDelete,
	}
	if input != nil {
		if input.Type != nil {
			cond[consts.TcType] = input.Type
		}
		if input.IsCallMe != nil && *input.IsCallMe == 2 {
			cond[consts.TcRelationType] = db.In([]string{consts.NoticeTypeIssueCommentAtSomebody, consts.NoticeTypeIssueRemarkAtSomebody})
		} else {
			//@人的notice独立在外
			cond[consts.TcRelationType] = db.NotIn([]string{consts.NoticeTypeIssueCommentAtSomebody, consts.NoticeTypeIssueRemarkAtSomebody})
		}

		if input.LastID != nil && *input.LastID != 0 {
			//按时间倒序
			cond[consts.TcId] = db.Lt(*input.LastID)
		}
		if input.ProjectID != nil && *input.ProjectID != 0 {
			cond[consts.TcProjectId] = *input.ProjectID
		}
		if input.IssueID != nil && *input.IssueID != 0 {
			cond[consts.TcIssueId] = *input.IssueID
		}
	}

	count, err := mysql.SelectCountByCond(consts.TableNotice, cond)
	if err != nil {
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	noticePo := &[]po.PpmTreNotice{}
	err = mysql.SelectAllByCondWithNumAndOrder(consts.TableNotice, cond, nil, page, size, "create_time desc", noticePo)
	if err != nil {
		return 0, nil, errs.BuildSystemErrorInfo(errs.MysqlOperateError, err)
	}

	noticeBo := &[]bo.NoticeBo{}
	copyErr := copyer.Copy(noticePo, noticeBo)
	if copyErr != nil {
		log.Errorf("对象copy异常: %v", copyErr)
		return 0, nil, errs.BuildSystemErrorInfo(errs.ObjectCopyError, copyErr)
	}

	return count, noticeBo, nil
}
