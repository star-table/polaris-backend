package bo

type OssCallbackBo struct {
	CallbackUrl      string `json:"callbackUrl"`
	CallbackHost     string `json:"callbackHost"`
	CallbackBody     string `json:"callbackBody"`
	CallbackBodyType string `json:"callbackBodyType"`
}

//oss callback body
type OssCallBackBody struct {
	Host      string `json:"host"`      //host
	UserId    int64  `json:"userId"`    //用户id
	OrgId     int64  `json:"orgId"`     //组织id
	Type      int    `json:"type"`      //类型
	Path      string `json:"path"`      //文件path
	Filename  string `json:"filename"`  //文件名(不带后缀)
	ProjectId int64  `json:"projectId"` //项目id
	IssueId   int64  `json:"issueId"`   //任务id
	ColumnId  string `json:"columnId"`  //列id
	TableId   string `json:"tableId"`   //表id
	AppId     string `json:"appId"`     //APPid
	FolderId  int64  `json:"folderId"`  //目录id
	Bucket    string `json:"bucket"`    //存储空间
	Size      int64  `json:"size"`      //大小
	Format    string `json:"format"`    //格式（后缀）
	Object    string `json:"object"`    //对象（文件path）
	RealName  string `json:"realName"`  //文件实际的名称，需要前端传x:filename
}

//oss callback body ext properties
type OssCallBackExtProperties struct {
	Size     int64  `json:"size"`     //大小
	Format   string `json:"format"`   //格式（后缀）
	Object   string `json:"object"`   //对象（文件path）
	Filename string `json:"filename"` //文件实际的名称，需要前端传x:filename
}

//oss的key信息
type OssKeyInfo struct {
	OrgId     int64 `json:"orgId"`
	ProjectId int64 `json:"projectId"`
	IssueId   int64 `json:"issueId"`
}
