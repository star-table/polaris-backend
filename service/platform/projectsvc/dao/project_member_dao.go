package dao

import (
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"upper.io/db.v3"
)

func GetProjectMemberCount(orgId, projectId int64) int64 {
	count, err := mysql.SelectCountByCond(consts.TableProjectRelation, db.Cond{
		"org_id":        orgId,
		"project_id":    projectId,
		"is_delete":     consts.AppIsNoDelete,
		"relation_type": []int64{consts.IssueRelationTypeOwner, consts.IssueRelationTypeParticipant},
	})

	if err != nil {
		return 0
	}

	return int64(count)
}
