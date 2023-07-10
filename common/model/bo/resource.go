package bo

type IssueCreateResourceReqBo struct {
	// 资源路径
	ResourcePath string `json:"resourcePath"`
	// 资源大小，单位B
	ResourceSize int64 `json:"resourceSize"`
	// 文件名
	FileName string `json:"fileName"`
	// 文件后缀
	FileSuffix string `json:"fileSuffix"`
	// md5
	Md5 *string `json:"md5"`
	// bucketName
	BucketName *string `json:"bucketName"`
	// 操作人
	OperatorId int64
	// 来源类型
	PolicyType int
	// 资源存储类型
	ResourceType int
}

type ProjectCreateResourceReqBo struct {
	ProjectId int64 `json:"projectId"`
	OrgId     int64 `json:"orgId"`
	FolderId  int64 `json:"folderId"`
	// 资源路径
	ResourcePath string `json:"resourcePath"`
	// 资源大小，单位B
	ResourceSize int64 `json:"resourceSize"`
	// 文件名
	FileName string `json:"fileName"`
	// 文件后缀
	FileSuffix string `json:"fileSuffix"`
	// md5
	Md5 *string `json:"md5"`
	// bucketName
	BucketName *string `json:"bucketName"`
	//文件类型
	FileType int
	// 操作人
	OperatorId int64
	// 来源类型
	PolicyType int
	// 资源存储类型
	ResourceType int
}

type CreateResourceBo struct {
	Id         int64 `db:"id,omitempty" json:"id"`
	ProjectId  int64
	OrgId      int64  `db:"org_id,omitempty" json:"orgId"`
	Bucket     string `db:"bucket,omitempty" json:"bucket"`
	Path       string `db:"path,omitempty" json:"path"`
	Name       string `db:"name,omitempty" json:"name"`
	Size       int64
	Suffix     string
	Type       int    `db:"type,omitempty" json:"type"`
	Md5        string `db:"md5,omitempty" json:"md5"`
	OperatorId int64
	//新增folderId用户文件管理创建资源	2019/12/12
	FolderId *int64
	//新增文件类型,用于区分文件和附件 文件来源,0其他,1项目封面,2任务附件,3导入文件,4项目资源,5兼容测试,6头像 2019/12/20
	SourceType *int

	//文件本地路径, 用于图片压缩
	DistPath string `json:"distPath"`
	IssueId  int64  `json:"issueId"`

	//是否跳过创建关联关系
	SkipCreateRelation bool `json:"skipCreateRelation"`
}

type UpdateResourceInfoBo struct {
	UserId int64 `json:"userId"`
	OrgId  int64 `json:"orgId"`
	// 文件id
	ResourceId int64 `json:"resourceId"`
	// 应用id
	AppId int64 `json:"appId"`
	// 项目id
	ProjectId int64 `json:"projectId"`
	// issueId
	IssueId int64 `json:"issueId"`
	// 文件大小
	FileSize int64 `json:"fileSize"`
	// 文件名
	FileName *string `json:"fileName"`
	// 文件后缀
	FileSuffix *string `json:"fileSuffix"`
	// 修改项
	UpdateFields []string `json:"updateFields"`
}

type UpdateResourceFolderBo struct {
	UserId int64 `json:"userId"`
	OrgId  int64 `json:"orgId"`
	// 当前文件夹id
	CurrentFolderId int64 `json:"currentFolderId"`
	// 目标文件夹id
	TargetFolderID int64 `json:"targetFolderId"`
	// 文件id列表
	ResourceIds []int64 `json:"resourceIds"`
	// 项目id
	ProjectId int64 `json:"projectId"`
}

type DeleteResourceBo struct {
	ResourceIds      []int64
	FolderId         *int64
	UserId           int64
	OrgId            int64
	ProjectId        int64
	IssueId          int64
	RecycleVersionId int64
	AppId            int64
}

type GetResourceBo struct {
	FolderId   *int64
	UserId     int64
	OrgId      int64
	ProjectId  int64
	Page       int
	Size       int
	SourceType []int
	FileType   *int
	KeyWord    *string
}

type GetResourceInfoBo struct {
	UserId      int64
	OrgId       int64
	AppId       int64
	ResourceId  int64
	IssueId     int64
	SourceTypes []int
	FileType    *int
	KeyWord     *string
}

type ResourceRelationBo struct {
	OrgId       int64
	ProjectId   int64
	IssueId     int64
	ResourceIds []int64
	SourceType  int
}
