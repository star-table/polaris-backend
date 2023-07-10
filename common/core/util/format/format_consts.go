package format

const (
	EmailPattern               = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //用户邮箱
	PasswordPattern            = `^[a-zA-Z]\w{5,17}$`                          //用户密码
	UserNamePattern            = "^[ 0-9A-Za-z-_]{1,40}$"                      //用户名
	OrgNamePattern             = "^[\u4e00-\u9fa5|0-9|a-zA-Z]{1,20}$"          //组织名
	OrgCodePattern             = "^[0-9|a-zA-Z]{1,20}$"                        //网址后缀编号
	OrgAdressPattern           = `^.{0,100}$`                                  //组织地址
	ProjectNamePattern         = `^.{1,30}$`                                   //项目名
	ProjectPreviousCodePattern = `^[a-zA-Z|0-9]{1,50}$`                        //项目前缀编号
	ProjectRemarkPattern       = `^[\s\S]{0,500}$`                             //项目描述(简介)
	//ProjectNoticePattern       = `^.{0,2000}$`                                 //项目公告
	IssueNamePattern = `^(.|\n){0,500}$` //任务名 允许换行符
	//IssueRemarkPattern           = `^.{0,10000}$`                                //任务描述(详情)
	IssueCommenPattern           = `^.{0,10000}$` //任务评论
	ProjectObjectTypeNamePattern = `^.{1,30}$`    //标题栏名
	ProjectTableNamePattern      = `^.{1,200}$`   // 项目表名
	//ResourceNamePattern          = `^[^\\\\/:*?\"<>|]{1,300}(\.[a-zA-Z0-9]+)?$` //资源名
	ResourceNamePattern            = `^.{1,300}$`                                                                 //资源名
	FolderNamePattern              = `^[^\\\\/:*?\"<>|]{1,30}$`                                                   //文件夹名
	RoleNamePattern                = "^[a-zA-Z|0-9]{1,20}$"                                                       //角色名
	NumFloat1                      = `^([1-9]\d{0,9}|0)(\.\d)?$`                                                  //小数点后一位
	SqlFieldPattern                = `(^_([a-zA-Z0-9]_?)*$)|(^[a-zA-Z](_?[a-zA-Z0-9])*_?$)`                       //数字、字母、下划线
	ChinaPhoneWithAreaFlagPattern  = `\+86\ 1[0-9]{10}`                                                           //带区号的中国手机号。如：`+86 15010111001`
	ChinaPhoneWithoutAreaPattern   = `1[0-9]{10}`                                                                 //不带区号的中国手机号。如：`15010111001`
	ForeignPhoneWithoutAreaPattern = `[0-9]{3,15}`                                                                //不带区号的国际手机号。一定是数值，并且一般不超过 15 位
	PhoneRegionCodePattern         = `\+\d{1,4}`                                                                  // 国际的手机号码归属地 code
	WebLinkPattern                 = `^(?:http(s)?:\/\/)?[\w.-]+(?:\.[\w\.-]+)+[\w\-\._~:/?#[\]@!\$&'\*\+,;=.]+$` // web 链接

	IssueCommentMemberPattern = "\\#\\[([\u4E00-\u9FA5A-Za-z0-9_]{0,})\\:\\d+\\]\\&\\$" // @成员 评论内容提取成员名字，如：@#[小牟:29520]&$
)

const (
	ChinesePattern  = "[\u4e00-\u9fa5]+?"
	AllBlankPattern = "^ +$"
)
