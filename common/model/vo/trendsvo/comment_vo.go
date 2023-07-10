package trendsvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type CreateCommentReqVo struct {
	CommentBo bo.CommentBo `json:"commentBo"`
}

type CreateCommentRespVo struct {
	CommentId int64 `json:"data"`
	
	vo.Err
}