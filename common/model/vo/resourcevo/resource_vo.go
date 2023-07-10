package resourcevo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
)

type CreateResourceReqVo struct {
	CreateResourceBo bo.CreateResourceBo `json:"createResourceBo"`
}

type CreateResourceRespVo struct {
	ResourceId int64 `json:"data"`

	vo.Err
}

type UpdateResourceInfoReqVo struct {
	Input bo.UpdateResourceInfoBo `json:"updateResourceBo"`
}
type UpdateResourceInfoResVo struct {
	*UpdateResourceData
	vo.Err
}
type UpdateResourceData struct {
	OldBo             []bo.ResourceBo
	NewBo             []bo.ResourceBo
	CurrentFolderName *string
	TargetFolderName  *string
}

type UpdateResourceFolderReqVo struct {
	Input bo.UpdateResourceFolderBo `json:"updateResourceBo"`
}

type DeleteResourceReqVo struct {
	Input bo.DeleteResourceBo `json:"deleteResourceBo"`
}

type GetResourceReqVo struct {
	Input bo.GetResourceBo `json:"getResourceBo"`
}

type GetResourceInfoReqVo struct {
	Input bo.GetResourceInfoBo `json:"getResourceInfoBo"`
}
type InsertResourceReqVo struct {
	Input InsertResourceReqData `json:"input"`
}

type InsertResourceReqData struct {
	ResourcePath  string `json:"resourcePath"`
	OrgId         int64  `json:"orgId"`
	CurrentUserId int64  `json:"currentUserId"`
	ResourceType  int    `json:"resourceType"`
	FileName      string `json:"fileName"`
}

type InsertResourceRespVo struct {
	ResourceId int64 `json:"data"`

	vo.Err
}

type GetResourceByIdReqVo struct {
	GetResourceByIdReqBody GetResourceByIdReqBody `json:"getResourceByIdReqBody"`
}

type GetResourceByIdReqBody struct {
	ResourceIds []int64 `json:"resourceIds"`
}

type GetResourceByIdRespVo struct {
	ResourceBos []bo.ResourceBo `json:"data"`

	vo.Err
}

type GetIdByPathReqVo struct {
	OrgId        int64  `json:"orgId"`
	ResourcePath string `json:"resourcePath"`
	ResourceType int    `json:"resourceType"`
}

type GetIdByPathRespVo struct {
	ResourceId int64 `json:"data"`

	vo.Err
}

type GetResourceBoListReqVo struct {
	Page  uint                  `json:"page"`
	Size  uint                  `json:"size"`
	Input GetResourceBoListCond `json:"cond"`
}

type GetResourceBoListCond struct {
	OrgId       int64    `json:"orgId"`
	ResourceIds *[]int64 `json:"resourceIds"`
	IsDelete    *int     `json:"isDelete"`
}

type GetResourceBoListRespVo struct {
	GetResourceBoListRespData `json:"data"`

	vo.Err
}

type GetResourceBoListRespData struct {
	ResourceBos *[]bo.ResourceBo `json:"resourceBos"`
	Total       int64            `json:"total"`
}

type GetResourceVoListRespVo struct {
	*vo.ResourceList `json:"data"`

	vo.Err
}

type GetResourceVoInfoRespVo struct {
	*vo.Resource `json:"data"`

	vo.Err
}

type GetIssueIdsByResourceIdsReqVo struct {
	OrgId  int64                       `json:"orgId"`
	UserId int64                       `json:"userId"`
	Input  GetIssueIdsByResourceIdsReq `json:"input"`
}

type GetIssueIdsByResourceIdsReq struct {
	ProjectId   int64   `json:"projectId"`
	ResourceIds []int64 `json:"resourceIds"`
}

type GetIssueIdsByResourceIdsRespVo struct {
	vo.Err
	Data []GetIssueIdsByResourceIdsResp `json:"data"`
}

type GetIssueIdsByResourceIdsResp struct {
	ResourceId int64   `json:"resourceId"`
	IssueIds   []int64 `json:"issueIds"`
}

type CreateIssueResourceReqVo struct {
	Input  vo.CreateProjectResourceReq `json:"input"`
	UserId int64                       `json:"userId"`
	OrgId  int64                       `json:"orgId"`
}

type CreateResourceRelationReqVo struct {
	Input  CreateResourceRelationData `json:"input"`
	UserId int64                      `json:"userId"`
	OrgId  int64                      `json:"orgId"`
}

type CreateResourceRelationData struct {
	ProjectId   int64   `json:"projectId"`
	IssueId     int64   `json:"issueId"`
	ResourceIds []int64 `json:"resourceIds"`
}

type CreateResourceRelationRespVo struct {
	ResourceIds []int64 `json:"data"`
	vo.Err
}

type GetResourceRelationsByProjectIdReqVo struct {
	Input  GetResourceRelationsByProjectIdData `json:"input"`
	UserId int64                               `json:"userId"`
	OrgId  int64                               `json:"orgId"`
}

type GetResourceRelationsByProjectIdData struct {
	ProjectId   int64   `json:"projectId"`
	SourceTypes []int32 `json:"sourceTypes"`
}

type GetResourceRelationsByProjectIdRespVo struct {
	ResourceRelations []ResourceRelationVo `json:"data"`
	vo.Err
}

type ResourceRelationVo struct {
	Id           int64  `json:"id"`
	OrgId        int64  `json:"orgId"`
	ProjectId    int64  `json:"projectId"`
	IssueId      int64  `json:"issueId"`
	ResourceId   int64  `json:"resourceId"`
	SourceType   int    `json:"sourceType"`
	ColumnId     string `json:"columnId"`
	ResourceType int    `json:"resourceType"`
}

type RecoverResourceReqVo struct {
	OrgId         int64               `json:"orgId"`
	UserId        int64               `json:"userId"`
	SourceChannel string              `json:"sourceChannel"`
	Input         RecoverResourceData `json:"input"`
}

type RecoverResourceRespVo struct {
	vo.Err
	Data *bo.ResourceBo `json:"data"`
}

type RecoverResourceData struct {
	ProjectId        int64   `json:"projectId"`
	IssueIds         []int64 `json:"issueIds"`
	ResourceId       int64   `json:"resourceId"`
	RecycleVersionId int64   `json:"recycleVersionId"`
	AppId            int64   `json:"appId"`
}

type DeleteResourceRelationReqVo struct {
	OrgId  int64                      `json:"orgId"`
	UserId int64                      `json:"userId"`
	Input  DeleteResourceRelationData `json:"input"`
}

type DeleteResourceRelationData struct {
	ProjectId        int64   `json:"projectId"`
	IssueIds         []int64 `json:"issueIds"`
	SourceTypes      []int   `json:"sourceTypes"`
	RecycleVersionId int64   `json:"recycleVersionId"`
}

type UpdateResourceRelationProjectIdReqVo struct {
	OrgId  int64                               `json:"orgId"`
	UserId int64                               `json:"userId"`
	Input  UpdateResourceRelationProjectIdData `json:"input"`
}

type UpdateResourceRelationProjectIdData struct {
	ProjectId int64   `json:"projectId"`
	IssueIds  []int64 `json:"issueIds"`
}

type CacheResourceSizeReq struct {
	OrgId int64 `json:"orgId"`
}

type CacheResourceSizeResp struct {
	vo.Err
	Size int64 `json:"size"`
}

type CompleteDeleteResourceReq struct {
	OrgId      int64 `json:"orgId"`
	UserId     int64 `json:"userId"`
	ResourceId int64 `json:"resourceId"`
}

type FsDocumentListReq struct {
	OrgId  int64                `json:"orgId"`
	UserId int64                `json:"userId"`
	Page   int                  `json:"page"`
	Size   int                  `json:"size"`
	Params FsDocumentLisReqData `json:"params"`
	//SearchKey string `json:"searchKey"`
}

type FsDocumentLisReqData struct {
	SearchKey string `json:"searchKey"`
}

type FsDocumentListResp struct {
	vo.Err
	Data *vo.FsDocumentListResp `json:"data"`
}

type AddFsResourceBatchReq struct {
	OrgId     int64                          `json:"orgId"`
	UserId    int64                          `json:"userId"`
	ProjectId int64                          `json:"projectId"`
	IssueId   int64                          `json:"issueId"`
	FolderId  int64                          `json:"folderId"`
	Input     []*vo.AddIssueAttachmentFsData `json:"input"`
}

type AddFsResourceBatchResp struct {
	vo.Err
	Ids []int64 `json:"ids"`
}

type DeleteAttachmentRelationReq struct {
	OrgId  int64                         `json:"orgId"`
	UserId int64                         `json:"userId"`
	Input  *DeleteAttachmentRelationData `json:"input"`
}

type DeleteAttachmentRelationData struct {
	ProjectId        int64   `json:"projectId"`
	IssueIds         []int64 `json:"issueIds"`
	RecycleVersionId int64   `json:"recycleVersionId"`
	ColumnId         string  `json:"columnId"`
	ResourceIds      []int64 `json:"resourceIds"`
	IsDeleteResource bool    `json:"isDeleteResource"` // 是否删除资源本身(目前只在删除表、删除附件列的情况下使用)
}

type GetResourceRelationListReq struct {
	OrgId  int64                    `json:"orgId"`
	UserId int64                    `json:"userId"`
	Input  *GetResourceRelationList `json:"input"`
}

type GetResourceRelationList struct {
	ProjectId          int64   `json:"projectId"`
	SourceTypes        []int32 `json:"sourceTypes"` // sourceType是 osspolicy的type类型
	ResourceIds        []int64 `json:"resourceIds"`
	IsDelete           int     `json:"isDelete"`
	VersionId          int     `json:"versionId"`
	IsNeedResourceType bool    `json:"isNeedResourceType"` // 是否需要资源的type类型
}

type GetResourceRelationListResp struct {
	vo.Err
	Data []ResourceRelationVo `json:"data"`
}

type AddResourceRelationReq struct {
	OrgId  int64                    `json:"orgId"`
	UserId int64                    `json:"userId"`
	Input  *AddResourceRelationData `json:"input"`
}

type AddResourceRelationData struct {
	ProjectId   int64   `json:"projectId"`
	IssueId     int64   `json:"issueId"`
	SourceType  int     `json:"sourceType"`
	ResourceIds []int64 `json:"resourceIds"`
	ColumnId    string  `json:"columnId"`
}

type DingFileListReq struct {
	OrgId  int64               `json:"orgId"`
	UserId int64               `json:"userId"`
	Page   int                 `json:"page"`
	Size   int                 `json:"size"`
	Input  DingFileListReqData `json:"input"`
}

type DingFileListReqData struct {
	//Code      string `json:"code"`      // 免登授权码
	//DingUserId string `json:"dingUserId"` // 用户userid
	SpaceType string `json:"spaceType"` // 空间类型  personal：私人空间, org：企业空间
	SpaceName string `json:"spaceName"` // 空间名称
}

type DingFileListRespVo struct {
	vo.Err
	Data *DingFileListResp `json:"data"`
}

type DingFileListResp struct {
	Total int64               `json:"total"`
	List  []*DingFileListData `json:"list"`
}

type DingFileListData struct {
	FileType      string `json:"fileType"`      // 文件类型, 文件、文件夹
	ContentType   string `json:"contentType"`   // 文件内容类型, text、document、image等
	ParentId      string `json:"parentId"`      // 父目录id
	FileId        string `json:"fileId"`        // 文件id
	FileName      string `json:"fileName"`      // 文件名称
	FilePath      string `json:"filePath"`      // 文件路径
	FileExtension string `json:"fileExtension"` // 文件后缀名
	FileSize      int64  `json:"fileSize"`      // 文件大小
	Creator       string `json:"creator"`       // 创建者id
	Modifier      string `json:"modifier"`      // 修改者id
	OwnerName     string `json:"ownerName"`     // 所有者名字
}

type DingDocReq struct {
	OrgId  int64           `json:"orgId"`
	UserId int64           `json:"userId"`
	Page   int             `json:"page"`
	Size   int             `json:"size"`
	Input  *DingDocReqData `json:"input"`
}

type DingDocReqData struct {
}

type DingDocumentData struct {
	DocId   string `json:"docId"`
	DocName string `json:"docName"`
	DocUrl  string `json:"docUrl"`
}

type DingDocumentRespVo struct {
	vo.Err
	Data *DingDocumentResp `json:"data"`
}

type DingDocumentResp struct {
	Total int64               `json:"total"`
	List  []*DingDocumentData `json:"list"`
}

// 钉钉团队空间信息
type SpaceInfo struct {
	SpaceId        string // 空间Id
	SpaceName      string // 空间名称
	Quota          int64  // 空间总额度
	UsedQuota      int64  // 空间已使用额度
	PermissionMode string // 授权模式
}

type DingSpaceListResp struct {
	vo.Err
	Data []*SpaceInfo `json:"data"`
}

type DingSpaceReqVo struct {
	OrgId  int64         `json:"orgId"`
	UserId int64         `json:"userId"`
	Input  *DingSpaceReq `json:"input"`
}

type DingSpaceReq struct {
	SpaceType  string `json:"spaceType"`  // 企业空间 org
	NextToken  string `json:"nextToken"`  // 分页游标
	MaxResults int    `json:"maxResults"` // 分页大小
}

type DingSpaceFileReqVo struct {
	OrgId  int64             `json:"orgId"`
	UserId int64             `json:"userId"`
	Input  *DingSpaceFileReq `json:"input"`
}

type DingSpaceFileReq struct {
	SpaceId string `json:"spaceId"`
	DirId   string `json:"dirId"`
	Size    int    `json:"size"`
}

type DingSpaceFileResp struct {
	vo.Err
	Data []*DingFileListData `json:"data"`
}
