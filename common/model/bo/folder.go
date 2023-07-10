package bo

import (
	"time"
)

type FolderBo struct {
	Id         int64     `db:"id,omitempty" json:"id"`
	OrgId      int64     `db:"org_id,omitempty" json:"orgId"`
	ProjectId  int64     `db:"project_id,omitempty" json:"projectId"`
	Name       string    `db:"name,omitempty" json:"name"`
	ParentId   int64     `db:"parent_id,omitempty" json:"parentId"`
	FileType   int64     `db:"file_type,omitempty" json:"fileType"`
	Path       string    `db:"path,omitempty" json:"path"`
	Creator    int64     `db:"creator,omitempty" json:"creator"`
	CreateTime time.Time `db:"create_time,omitempty" json:"createTime"`
	Updator    int64     `db:"updator,omitempty" json:"updator"`
	UpdateTime time.Time `db:"update_time,omitempty" json:"updateTime"`
	Version    int       `db:"version,omitempty" json:"version"`
	IsDelete   int       `db:"is_delete,omitempty" json:"isDelete"`
}

type CreateFolderBo struct {
	Id    int64 `db:"id,omitempty" json:"id"`
	OrgId int64 `json:"orgId"`
	// 项目id
	ProjectId int64 `json:"projectId"`
	// 文件夹名
	Name string `json:"name"`
	// 父级文件夹id
	ParentId int64 `json:"parentId"`
	// 文件夹类型,0其他,1文档,2图片,3视频,4音频
	FileType int64 `json:"fileType"`
	// 操作人
	UserId int64 `json:"userId"`
}

type UpdateFolderBo struct {
	// 文件夹id
	FolderID int64 `json:"folderId"`
	// 项目id
	ProjectID int64 `json:"projectId"`
	// 文件夹名
	Name *string `json:"name"`
	// 父级文件夹id
	ParentID *int64 `json:"parentId"`
	// 文件夹类型,0其他,1文档,2图片,3视频,4音频
	FileType *int64 `json:"fileType"`
	UserId   int64
	OrgId    int64
	// 变动的字段列表
	UpdateFields []string `json:"updateFields"`
}

type DeleteFolderBo struct {
	FolderIds []int64
	UserId    int64
	OrgId     int64
	ProjectId int64
}

type GetFolderBo struct {
	ParentId  *int64
	UserId    int64
	OrgId     int64
	ProjectId int64
	Page      int
	Size      int
}

type FolderTreeBo struct {
	Id int64
	ParentId int64
	Children []*FolderTreeBo
}
