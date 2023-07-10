package handler

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/star-table/common/core/config"
	"github.com/star-table/common/core/logger"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/core/util/strs"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/core/util"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/common/model/vo/resourcevo"
	"github.com/star-table/polaris-backend/facade/officefacade"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/star-table/polaris-backend/facade/projectfacade"
	"github.com/star-table/polaris-backend/facade/resourcefacade"

	"github.com/gin-gonic/gin"
)

const POS = "."

const DIR = "/"

const PATH = "resource"

var log = logger.GetDefaultLogger()

type Result struct {
	Code    int32      `json:"code"`
	Message string     `json:"message"`
	Data    UploadResp `json:"data"`
}

type UploadResp struct {
	FileList map[string]FileList `json:"fileList"`
}

type FileList struct {
	Url      string `json:"url"`
	Size     string `json:"size"`
	SourceId int64  `json:"sourceId"`
	FileName string `json:"fileName"`
	DistPath string `json:"distPath"`
}

func FileUploadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		runMode := config.GetConfig().Application.RunMode
		if runMode == 1 {
			errorHandle(errs.BuildSystemErrorInfo(errs.RunModeUnsupportUpload), c.Writer)
			return
		}
		r := c.Request
		// 整体限制 100Mb
		//r.Body = http.MaxBytesReader(c.Writer, r.Body, 1024<<20)
		_, err := c.MultipartForm()
		if err != nil {
			log.Error(err)
			errorHandle(errs.BuildSystemErrorInfo(errs.FileTooLarge, err), c.Writer)
			return
		}

		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
		if userErr != nil {
			log.Error(userErr)
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
			return
		}

		projectId, issueId, _, _, folderId, _, policyType := baseParam(r)
		result, err := uploadFile(c, cacheUserInfo, projectId, issueId, folderId, policyType)
		if err != nil {
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), c.Writer)
			return
		} else {
			//jsonStr, _ := json.Marshal(result)
			jsonStr := json.ToJsonIgnoreError(result)
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.Write([]byte(jsonStr))
			return
		}
	}
}

func FileUploadAndUpdateResourceInfoHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var (
			result *Result = nil
		)

		log.Info("[FileHandler] FileUploadAndUpdateResourceInfoHandler")

		runMode := config.GetConfig().Application.RunMode
		if runMode == 1 {
			errorHandle(errs.BuildSystemErrorInfo(errs.RunModeUnsupportUpload), c.Writer)
			return
		}
		r := c.Request
		// 整体限制 100Mb
		r.Body = http.MaxBytesReader(c.Writer, r.Body, 10<<20)
		_, err := c.MultipartForm()
		if err != nil {
			log.Error(err)
			errorHandle(errs.BuildSystemErrorInfo(errs.FileTooLarge, err), c.Writer)
			return
		}

		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
		if userErr != nil {
			log.Error(userErr)
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
			return
		}

		appId, issueId, _, _, _, resourceId, policyType := baseParam(r)
		if resourceId <= 0 {
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), c.Writer)
			return
		}

		result, err = uploadUpdateResourceInfo(c, cacheUserInfo, appId, issueId, resourceId, policyType)
		if err != nil {
			log.Errorf("[FileUploadHandler] uploadUpdateFile Failed: %s", strs.ObjectToString(err))
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), c.Writer)
			return
		}

		if err != nil {
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), c.Writer)
			return
		} else {
			//jsonStr, _ := json.Marshal(result)
			jsonStr := json.ToJsonIgnoreError(result)
			c.Writer.Header().Set("Content-Type", "application/json")
			c.Writer.Write([]byte(jsonStr))
			return
		}
	}
}

func FileReadHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		readFile(c)
	}
}

func baseParam(r *http.Request) (int64, int64, int64, int64, int64, int64, int) {
	var issueId, projectId, tableId, iterationId, folderId, resourceId int64
	var policyType int
	queries := r.URL.Query()
	form := r.Form
	if len(queries["projectId"]) > 0 {
		projectId, _ = strconv.ParseInt(queries["projectId"][0], 10, 64)
	}
	if len(queries["issueId"]) > 0 {
		issueId, _ = strconv.ParseInt(queries["issueId"][0], 10, 64)
	}
	if len(queries["policyType"]) > 0 {
		policyType, _ = strconv.Atoi(queries["policyType"][0])
	}
	if len(queries["iterationId"]) > 0 {
		iterationId, _ = strconv.ParseInt(queries["iterationId"][0], 10, 64)
	}
	if len(form["projectId"]) > 0 {
		projectId, _ = strconv.ParseInt(form["projectId"][0], 10, 64)
	}
	if len(form["issueId"]) > 0 {
		issueId, _ = strconv.ParseInt(form["issueId"][0], 10, 64)
	}
	if len(form["policyType"]) > 0 {
		policyType, _ = strconv.Atoi(form["policyType"][0])
	}
	if len(form["tableId"]) > 0 {
		tableId, _ = strconv.ParseInt(form["tableId"][0], 10, 64)
	}
	if len(form["iterationId"]) > 0 {
		iterationId, _ = strconv.ParseInt(form["iterationId"][0], 10, 64)
	}
	if len(form["folderId"]) > 0 {
		folderId, _ = strconv.ParseInt(form["folderId"][0], 10, 64)
	}
	if len(form["resourceId"]) > 0 {
		resourceId, _ = strconv.ParseInt(form["resourceId"][0], 10, 64)
	}
	return projectId, issueId, tableId, iterationId, folderId, resourceId, policyType
}

// @Security PM-TOEKN
// @Summary 任务导入
// @Description 任务导入
// @Tags 任务
// @accept application/json
// @Produce application/json
// @Success 200 {object} int64
// @Failure 400
// @Router /api/task/importData [post]
func ImportDataHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		r := c.Request
		_, err := c.MultipartForm()
		if err != nil {
			log.Error(err)
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), c.Writer)
			return
		}

		cacheUserInfo, userErr := orgfacade.GetCurrentUserRelaxed(r.Context())
		if userErr != nil {
			log.Error(userErr)
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, userErr), c.Writer)
			return
		}

		projectId, issueId, tableId, iterationId, folderId, _, policyType := baseParam(r)
		result, err := uploadFile(c, cacheUserInfo, projectId, issueId, folderId, policyType)
		if err != nil {
			log.Error(err)
			errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), c.Writer)
			return
		} else {
			log.Infof("上传文件成功，准备导入数据，result： %s", json.ToJsonIgnoreError(result))

			if len(result.Data.FileList) == 0 {
				errorHandle(errs.ImportFileNotExist, c.Writer)
				return
			}
			url := ""
			urlType := consts.UrlTypeDistPath
			for _, v := range result.Data.FileList {
				//url = v.Url
				//改为本地路径
				url = v.DistPath
				break
			}

			importInput := projectvo.ImportIssuesReqVo{
				UserId: cacheUserInfo.UserId,
				OrgId:  cacheUserInfo.OrgId,
				Input: vo.ImportIssuesReq{
					URL:         url,
					ProjectID:   projectId,
					URLType:     urlType,
					TableID:     fmt.Sprintf("%v", tableId),
					IterationID: &iterationId,
				},
			}

			log.Infof("批量导入任务请求结构体 %s", json.ToJsonIgnoreError(importInput))
			c.Writer.Header().Set("Content-Type", "application/json")
			resp := projectfacade.ImportIssues(importInput)
			if resp.Failure() {
				log.Errorf("[ImportDataHandler] err: %v", resp.Error())
				// var errMsg string
				result := projectvo.ImportIssuesErrorRespVo{}
				result.Code = resp.Error().Code()
				result.Message = resp.Error().Message()
				cellErrArr := make([]string, 0)
				if resp.Error().Code() == errs.FileParseFail.Code() {
					//data := resp.Message[len(errs.FileParseFail.Message()):len(resp.Message)]
					result.Message = errs.FileParseFail.Message()
					dataStr := resp.Error().Message()
					if err := json.FromJson(dataStr, &cellErrArr); err != nil {
						log.Errorf("[ImportDataHandler] try decode err msg err. err: %v", err)
					}
					log.Errorf("[ImportDataHandler] importBigTitleSizeError: %s", dataStr)
					/*
						errMsg = "{\"code\":" + strconv.Itoa(resp.Error().Code()) + ",\"message\":\"" + errs.FileParseFail.Message() + "\", \"data\":" + data + "}"
					*/
					result.Data = cellErrArr
				} else {
					result.Data = cellErrArr
					/*
						errMsg = "{\"code\":" + strconv.Itoa(resp.Error().Code()) + ",\"message\":\"" + json.ToJsonIgnoreError(resp.Message) + "\"}"
					*/
				}
				errorHandle(errors.New(json.ToJsonIgnoreError(result)), c.Writer)
				return
			}

			// result := vo.Err{
			// 	Code:    200,
			// 	Message: fmt.Sprintf("共%d条任务数据上传成功！", resp.NewData.Count),
			// }
			result := projectvo.ImportIssuesRespVo{
				Err: vo.Err{
					Code:    200,
					Message: fmt.Sprintf("共%d条任务数据上传成功！", resp.Data.Count),
				},
				Data: resp.Data,
			}
			//jsonStr, _ := json.Marshal(result)
			jsonStr := json.ToJsonIgnoreError(result)
			c.Writer.Write([]byte(jsonStr))
			return
		}
	}
}

//读取文件
func readFile(c *gin.Context) {
	w := c.Writer
	path := c.Param("path")
	file, err := os.Open(config.GetConfig().OSS.RootPath + path)
	if err != nil {
		log.Error(err)
		errorHandle(errs.BuildSystemErrorInfo(errs.SystemError, err), w)
		return
	}

	defer file.Close()
	buff, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error(err)
		errorHandle(errs.FileReadFail, w)
		return
	}
	w.Write(buff)
}

// 统一错误输出接口
func errorHandle(err error, w http.ResponseWriter) {
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func getFileName(ext string, r *http.Request, projectId, issueId int64, policyType int) (string, string, error) {
	cacheUserInfo, err := orgfacade.GetCurrentUserRelaxed(r.Context())
	if err != nil {
		log.Info(strs.ObjectToString(err))
		return consts.BlankString, consts.BlankString, errs.BuildSystemErrorInfo(errs.TokenAuthError, err)
	}

	reqVo := resourcevo.GetOssPostPolicyReqVo{
		Input: vo.OssPostPolicyReq{
			ProjectID:  &projectId,
			IssueID:    &issueId,
			PolicyType: policyType,
		},
		OrgId: cacheUserInfo.OrgId,
	}

	resp := resourcefacade.GetOssPostPolicy(reqVo)

	if resp.Failure() {
		return consts.BlankString, consts.BlankString, resp.Error()
	}
	commonPath := resp.GetOssPostPolicy.Dir

	rootPath := config.GetOSSConfig().RootPath
	dstPath := rootPath + DIR + commonPath
	suffixSplit := strings.Split(ext, POS)
	suffix := suffixSplit[len(suffixSplit)-1]
	fileName := resp.GetOssPostPolicy.FileName + POS + suffix
	dstFile := dstPath + DIR + fileName

	_, err1 := os.Stat(dstPath)
	res := os.IsNotExist(err1)
	if res == true {
		os.MkdirAll(dstPath, os.ModePerm)
	}

	relatePath := commonPath + DIR + fileName
	return dstFile, relatePath, nil
}

//上传文件
func uploadFile(c *gin.Context, cacheUserInfo *bo.CacheUserInfoBo, projectId, issueId, folderId int64, policyType int) (Result, error) {
	result := Result{
		Code:    200,
		Message: "success",
		Data: UploadResp{
			FileList: map[string]FileList{},
		},
	}
	r := c.Request

	//支持多文件上传
	for k, _ := range r.MultipartForm.File {
		file, handler, err := r.FormFile(k)
		if err != nil {
			log.Error(err)
			return result, err
		}
		defer file.Close()

		dstFile, relatePath, err := getFileName(handler.Filename, r, projectId, issueId, policyType)
		if err != nil {
			log.Error(err)
			return result, err
		}

		fp, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Error(err)
			return result, err
		}
		defer fp.Close()

		size, err := io.Copy(fp, file)
		if err != nil {
			log.Error(err)
			return result, err
		}

		url := config.GetOSSConfig().LocalDomain + "/" + relatePath
		////单机部署
		//if runMode == 2 {
		//	url = fmt.Sprintf("%s/read/%s", config.GetOSSConfig().EndPoint, relatePath)
		//}

		var resourceId int64
		switch policyType {
		case consts.OssPolicyTypeIssueResource, consts.OssPolicyTypeLesscodeResource:
			respVo := projectfacade.CreateIssueResource(projectvo.CreateIssueResourceReqVo{
				UserId: cacheUserInfo.UserId,
				OrgId:  cacheUserInfo.OrgId,
				Input: vo.CreateIssueResourceReq{
					ProjectId: projectId,
					IssueId:   issueId,
					// 资源路径
					ResourcePath: url,
					// 资源大小，单位B
					ResourceSize: handler.Size,
					// 文件名
					FileName: handler.Filename,
					// 文件后缀
					FileSuffix: util.ParseFileSuffix(handler.Filename),
					// 来源类型
					PolicyType: policyType,
					// 存储类型
					ResourceType: consts.LocalResource,
				},
			})
			if respVo.Failure() {
				log.Error(respVo.Message)
				continue
			}
			resourceId = respVo.Void.ID

		case consts.OssPolicyTypeProjectResource:
			respVo := projectfacade.CreateProjectResource(projectvo.CreateProjectResourceReqVo{
				UserId: cacheUserInfo.UserId,
				OrgId:  cacheUserInfo.OrgId,
				Input: vo.CreateProjectResourceReq{
					ProjectID: projectId,
					FolderID:  folderId,
					// 资源路径
					ResourcePath: url,
					// 资源大小，单位B
					ResourceSize: handler.Size,
					// 文件名
					FileName: handler.Filename,
					// 文件后缀
					FileSuffix: util.ParseFileSuffix(handler.Filename),
					// 来源类型
					PolicyType: policyType,
					// 文件存储类型
					ResourceType: consts.LocalResource,
				},
			})
			if respVo.Failure() {
				log.Error(respVo.Message)
				continue
			}
			resourceId = respVo.Void.ID

		default:
			respVo := resourcefacade.CreateResource(resourcevo.CreateResourceReqVo{
				CreateResourceBo: bo.CreateResourceBo{
					Name:       handler.Filename,
					Suffix:     util.ParseFileSuffix(handler.Filename),
					Size:       handler.Size,
					Path:       url,
					ProjectId:  projectId,
					IssueId:    issueId,
					OrgId:      cacheUserInfo.OrgId,
					OperatorId: cacheUserInfo.UserId,
					Type:       consts.LocalResource,
					DistPath:   dstFile,
					SourceType: &policyType,
				},
			})
			if respVo.Failure() {
				log.Error(respVo.Error())
				continue
			}
			resourceId = respVo.ResourceId
		}

		result.Data.FileList[k] = FileList{
			Url:      url,
			Size:     strconv.FormatInt(size/1024, 10) + "KB",
			SourceId: resourceId,
			FileName: handler.Filename,
			DistPath: dstFile,
		}
	}

	return result, nil
}

//上传更新文件
func uploadUpdateResourceInfo(c *gin.Context, cacheUserInfo *bo.CacheUserInfoBo, appId, issueId int64, resourceId int64, policyType int) (*Result, error) {
	log.Infof("[uploadUpdateFile] uploadUpdateFile cacheUserInfo: %v, appId: %d, issueId: %d, resourceId: %d, policyType: %d", cacheUserInfo, appId, issueId, resourceId, policyType)

	result := &Result{
		Code:    200,
		Message: "success",
		Data: UploadResp{
			FileList: map[string]FileList{},
		},
	}
	r := c.Request

	//支持多文件上传
	for k, _ := range r.MultipartForm.File {
		file, handler, err := r.FormFile(k)
		if err != nil {
			//log.Error(err)
			log.Errorf("[uploadUpdateFile] FormFile Failed: %s, k: %v", strs.ObjectToString(err), k)
			return result, err
		}
		defer file.Close()

		dstFile, relatePath, err := getFileNameByResourceId(handler.Filename, r, cacheUserInfo, appId, resourceId, issueId, policyType)
		if err != nil {
			//log.Error(err)
			log.Errorf("[uploadUpdateFile] getFileNameByResourceId Failed: %s, Filename: %s, Request: %v, cacheUserInfo: %v, appId: %d, resourceId: %d, policyType: %d", strs.ObjectToString(err), handler.Filename, r, cacheUserInfo, appId, resourceId, policyType)
			return result, err
		}

		fp, err := os.OpenFile(dstFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			//log.Error(err)
			log.Errorf("[uploadUpdateFile] OpenFile Failed: %s, dstFile: %s", strs.ObjectToString(err), dstFile)
			return nil, err
		}
		defer fp.Close()

		log.Infof("[uploadUpdateFile] OpenFile success dstFile: %s", dstFile)

		size, err := io.Copy(fp, file)
		if err != nil {
			log.Errorf("[uploadUpdateFile] io.Copy Failed: %s, file: %v", strs.ObjectToString(err), file)
			return nil, err
		}

		url := relatePath

		updateFields := make([]string, 1)
		updateFields[0] = "fileSize"

		respVo := resourcefacade.UpdateResourceInfo(resourcevo.UpdateResourceInfoReqVo{
			Input: bo.UpdateResourceInfoBo{
				UserId:       cacheUserInfo.UserId,
				OrgId:        cacheUserInfo.OrgId,
				ResourceId:   resourceId,
				AppId:        appId,
				FileSize:     size,
				UpdateFields: updateFields,
			},
		})

		if respVo.Failure() {
			log.Errorf("[uploadUpdateFile] UpdateResourceInfo Failed: %s, UserId: %d, OrgId: %d, ResourceId: %d, AppId: %d, FileSize: %d, UpdateColumns: %v", strs.ObjectToString(err), cacheUserInfo.UserId, cacheUserInfo.OrgId, resourceId, appId, size, updateFields)
			return nil, respVo.Error()
		}

		log.Infof("[uploadUpdateFile] UpdateResourceInfo success UserId: %d, OrgId: %d, ResourceId: %d, AppId: %d, FileSize: %d, UpdateColumns: %v", cacheUserInfo.UserId, cacheUserInfo.OrgId, resourceId, appId, size, updateFields)

		result.Data.FileList[k] = FileList{
			Url:      url,
			Size:     strconv.FormatInt(size/1024, 10) + "KB",
			SourceId: resourceId,
			FileName: handler.Filename,
			DistPath: dstFile,
		}
	}

	return result, nil
}

func getFileNameByResourceId(
	ext string,
	r *http.Request,
	cacheUserInfo *bo.CacheUserInfoBo,
	appId,
	resourceId,
	issueId int64,
	policyType int) (string, string, error) {
	var (
		dstFile    string = ""
		relatePath string = ""
		rootPath   string = ""
		resp       resourcevo.GetResourceVoInfoRespVo
	)

	// log.Infof("[getFileNameByResourceId] uploadUpdateFile cacheUserInfo: %v, appId: %d, issueId: %d, resourceId: %d, policyType: %d", cacheUserInfo, appId, issueId, resourceId, policyType)

	if issueId == 0 {
		reqVo := resourcevo.GetResourceInfoReqVo{
			Input: bo.GetResourceInfoBo{
				UserId:      cacheUserInfo.UserId,
				OrgId:       cacheUserInfo.OrgId,
				AppId:       appId,
				ResourceId:  resourceId,
				SourceTypes: []int{policyType},
			},
		}

		resp = resourcefacade.GetResourceInfo(reqVo)

		if resp.Failure() {
			return consts.BlankString, consts.BlankString, resp.Error()
		}

	} else {

		reqVo := resourcevo.GetResourceInfoReqVo{
			Input: bo.GetResourceInfoBo{
				UserId:      cacheUserInfo.UserId,
				OrgId:       cacheUserInfo.OrgId,
				AppId:       appId,
				ResourceId:  resourceId,
				IssueId:     issueId,
				SourceTypes: []int{consts.OssPolicyTypeIssueResource, consts.OssPolicyTypeLesscodeResource},
			},
		}

		resp = resourcefacade.GetResourceInfo(reqVo)

		if resp.Failure() {
			return consts.BlankString, consts.BlankString, resp.Error()
		}
	}

	rootPath = config.GetOSSConfig().RootPath
	//dstPath := rootPath + DIR + commonPath

	dstFile = fmt.Sprintf("%s%s", rootPath, resp.Path) //resp.Path
	relatePath = resp.Path

	return dstFile, relatePath, nil
}

func GetOfficeConfigHandler() gin.HandlerFunc {
	return func(c *gin.Context) {

		var OrgId int64
		var UserId int64

		cacheUserInfo, err := GetCacheUserInfo(c)
		if err != nil {
			Fail(c, err)
			return
		}

		OrgId = cacheUserInfo.OrgId
		UserId = cacheUserInfo.UserId
		log.Infof("[GetOfficeConfigHandler] GetCacheUserInfo cacheUserInfo: %+v, OrgId: %d, UserId: %d", cacheUserInfo, OrgId, UserId)

		resp := officefacade.GetOfficeConfig(
			OrgId,
			UserId,
		)

		if resp.Failure() {
			Fail(c, resp.Error())
		} else {
			Success(c, resp.Data)
		}
	}
}
