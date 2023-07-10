package resourcefacade

import (
	"fmt"

	"github.com/star-table/common/core/config"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade"
)


func AddFsResourceBatch(req resourcevo.AddFsResourceBatchReq) resourcevo.AddFsResourceBatchResp {
	respVo := &resourcevo.AddFsResourceBatchResp{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/addFsResourceBatch", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["projectId"] = req.ProjectId
	queryParams["issueId"] = req.IssueId
	queryParams["folderId"] = req.FolderId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func AddResourceRelationWithType(req resourcevo.AddResourceRelationReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/addResourceRelationWithType", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CacheResourceSize(req resourcevo.CacheResourceSizeReq) resourcevo.CacheResourceSizeResp {
	respVo := &resourcevo.CacheResourceSizeResp{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/cacheResourceSize", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CompleteDeleteFolder(req resourcevo.CompleteDeleteFolderReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/completeDeleteFolder", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["projectId"] = req.ProjectId
	queryParams["folderId"] = req.FolderId
	queryParams["recycleId"] = req.RecycleId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CompleteDeleteResource(req resourcevo.CompleteDeleteResourceReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/completeDeleteResource", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["resourceId"] = req.ResourceId
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, nil, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateFolder(req resourcevo.CreateFolderReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/createFolder", config.GetPreUrl("resourcesvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateResource(req resourcevo.CreateResourceReqVo) resourcevo.CreateResourceRespVo {
	respVo := &resourcevo.CreateResourceRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/createResource", config.GetPreUrl("resourcesvc"))
	requestBody := &req.CreateResourceBo
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func CreateResourceRelation(req resourcevo.CreateResourceRelationReqVo) resourcevo.CreateResourceRelationRespVo {
	respVo := &resourcevo.CreateResourceRelationRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/createResourceRelation", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteAttachmentRelation(req resourcevo.DeleteAttachmentRelationReq) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/deleteAttachmentRelation", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteFolder(req resourcevo.DeleteFolderReqVo) resourcevo.DeleteFolderRespVo {
	respVo := &resourcevo.DeleteFolderRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/deleteFolder", config.GetPreUrl("resourcesvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteResource(req resourcevo.DeleteResourceReqVo) resourcevo.UpdateResourceInfoResVo {
	respVo := &resourcevo.UpdateResourceInfoResVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/deleteResource", config.GetPreUrl("resourcesvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DeleteResourceRelation(req resourcevo.DeleteResourceRelationReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/deleteResourceRelation", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DingDocumentList(req resourcevo.DingDocReq) resourcevo.DingDocumentRespVo {
	respVo := &resourcevo.DingDocumentRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/dingDocumentList", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DingFileList(req resourcevo.DingFileListReq) resourcevo.DingFileListRespVo {
	respVo := &resourcevo.DingFileListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/dingFileList", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DingFileListById(req resourcevo.DingSpaceFileReqVo) resourcevo.DingSpaceFileResp {
	respVo := &resourcevo.DingSpaceFileResp{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/dingFileListById", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func DingSpaceList(req resourcevo.DingSpaceReqVo) resourcevo.DingSpaceListResp {
	respVo := &resourcevo.DingSpaceListResp{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/dingSpaceList", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func FsDocumentList(req resourcevo.FsDocumentListReq) resourcevo.FsDocumentListResp {
	respVo := &resourcevo.FsDocumentListResp{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/fsDocumentList", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Params
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetFolder(req resourcevo.GetFolderReqVo) resourcevo.GetFolderVoListRespVo {
	respVo := &resourcevo.GetFolderVoListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getFolder", config.GetPreUrl("resourcesvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetFolderInfoBasic(req resourcevo.GetFolderInfoBasicReqVo) resourcevo.GetFolderInfoBasicRespVo {
	respVo := &resourcevo.GetFolderInfoBasicRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getFolderInfoBasic", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	requestBody := &req.FolderIds
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetIssueIdsByResourceIds(req resourcevo.GetIssueIdsByResourceIdsReqVo) resourcevo.GetIssueIdsByResourceIdsRespVo {
	respVo := &resourcevo.GetIssueIdsByResourceIdsRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getIssueIdsByResourceIds", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOssPostPolicy(req resourcevo.GetOssPostPolicyReqVo) resourcevo.GetOssPostPolicyRespVo {
	respVo := &resourcevo.GetOssPostPolicyRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getOssPostPolicy", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetOssSignURL(req resourcevo.OssApplySignURLReqVo) resourcevo.GetOssSignURLRespVo {
	respVo := &resourcevo.GetOssSignURLRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getOssSignURL", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetResource(req resourcevo.GetResourceReqVo) resourcevo.GetResourceVoListRespVo {
	respVo := &resourcevo.GetResourceVoListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getResource", config.GetPreUrl("resourcesvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetResourceBoList(req resourcevo.GetResourceBoListReqVo) resourcevo.GetResourceBoListRespVo {
	respVo := &resourcevo.GetResourceBoListRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getResourceBoList", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["page"] = req.Page
	queryParams["size"] = req.Size
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetResourceById(req resourcevo.GetResourceByIdReqVo) resourcevo.GetResourceByIdRespVo {
	respVo := &resourcevo.GetResourceByIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getResourceById", config.GetPreUrl("resourcesvc"))
	requestBody := &req.GetResourceByIdReqBody
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetResourceInfo(req resourcevo.GetResourceInfoReqVo) resourcevo.GetResourceVoInfoRespVo {
	respVo := &resourcevo.GetResourceVoInfoRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getResourceInfo", config.GetPreUrl("resourcesvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetResourceRelationList(req resourcevo.GetResourceRelationListReq) resourcevo.GetResourceRelationListResp {
	respVo := &resourcevo.GetResourceRelationListResp{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getResourceRelationList", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func GetResourceRelationsByProjectId(req resourcevo.GetResourceRelationsByProjectIdReqVo) resourcevo.GetResourceRelationsByProjectIdRespVo {
	respVo := &resourcevo.GetResourceRelationsByProjectIdRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/getResourceRelationsByProjectId", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["userId"] = req.UserId
	queryParams["orgId"] = req.OrgId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func RecoverFolder(req resourcevo.RecoverFolderReqVo) resourcevo.RecoverFolderRespVo {
	respVo := &resourcevo.RecoverFolderRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/recoverFolder", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func RecoverResource(req resourcevo.RecoverResourceReqVo) resourcevo.RecoverResourceRespVo {
	respVo := &resourcevo.RecoverResourceRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/recoverResource", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	queryParams["sourceChannel"] = req.SourceChannel
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateFolder(req resourcevo.UpdateFolderReqVo) resourcevo.UpdateFolderRespVo {
	respVo := &resourcevo.UpdateFolderRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/updateFolder", config.GetPreUrl("resourcesvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateResourceFolder(req resourcevo.UpdateResourceFolderReqVo) resourcevo.UpdateResourceInfoResVo {
	respVo := &resourcevo.UpdateResourceInfoResVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/updateResourceFolder", config.GetPreUrl("resourcesvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateResourceInfo(req resourcevo.UpdateResourceInfoReqVo) resourcevo.UpdateResourceInfoResVo {
	respVo := &resourcevo.UpdateResourceInfoResVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/updateResourceInfo", config.GetPreUrl("resourcesvc"))
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, nil, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

func UpdateResourceRelationProjectId(req resourcevo.UpdateResourceRelationProjectIdReqVo) vo.CommonRespVo {
	respVo := &vo.CommonRespVo{}
	reqUrl := fmt.Sprintf("%s/api/resourcesvc/updateResourceRelationProjectId", config.GetPreUrl("resourcesvc"))
	queryParams := map[string]interface{}{}
	queryParams["orgId"] = req.OrgId
	queryParams["userId"] = req.UserId
	requestBody := &req.Input
	err := facade.Request(consts.HttpMethodPost, reqUrl, queryParams, nil, requestBody, respVo)
	if err.Failure() {
		respVo.Err = err
	}
	return *respVo
}

