package bo

import (
	"github.com/star-table/common/core/types"
)

type TrendsBo struct {
	Id              int64      `json:"id"`
	OrgId           int64      `json:"orgId"`
	Uuid            string     `db:"uuid,omitempty" json:"uuid"`
	Module1         string     `json:"module1"`
	Module2Id       int64      `json:"module2Id"`
	Module2         string     `json:"module2"`
	Module3Id       int64      `json:"module3Id"`
	Module3         string     `json:"module3"`
	OperCode        string     `json:"operCode"`
	OperObjId       int64      `json:"operObjId"`
	OperObjType     string     `json:"operObjType"`
	OperObjProperty string     `json:"operObjProperty"`
	RelationObjId   int64      `json:"relationObjId"`
	RelationObjType string     `json:"relationObjType"`
	RelationType    string     `json:"relationType"`
	NewValue        *string    `json:"newValue"`
	OldValue        *string    `json:"oldValue"`
	Ext             string     `json:"ext"`
	Creator         int64      `json:"creator"`
	CreateTime      types.Time `json:"createTime"`
}

/**
动态分页对象
*/
type TrendsPageBo struct {
	Total int64       `json:"total"`
	Page  int64       `json:"page"`
	Size  int64       `json:"size"`
	List  *[]TrendsBo `json:"list"`
}

/**
查询动态条件对象
*/
type TrendsQueryCondBo struct {
	// 上次分页的最后一条动态id
	LastTrendID *int64 `json:"lastTrendId"`
	// 组织id
	OrgId int64 `json:"orgId"`
	// 对象id
	ObjId *int64 `json:"objId"`
	// 对象类型
	ObjType *string `json:"objType"`
	// 操作人id
	OperId *int64 `json:"operId"`

	// 开始时间
	StartTime *types.Time `json:"startTime"`

	// 结束时间
	EndTime *types.Time `json:"endTime"`

	// 分类（1动态2评论）
	Type *int `json:"type"`

	// 页码
	Page *int64 `json:"page"`

	//排序(1时间正序2时间倒序)
	OrderType *int `json:"orderType"`

	// 分页数量
	Size *int64 `json:"size"`

	CurrentUserId int64 `json:"currentUserId"`
}
