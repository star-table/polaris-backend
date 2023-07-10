package consts

import (
	"github.com/star-table/polaris-backend/common/core/consts"
)

var (
	//用户配置缓存
	CacheUserConfig = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "config"
	//用户基础信息缓存key
	CacheBaseUserInfo = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "baseinfo"
	//用户外部信息缓存key
	CacheBaseUserOutInfo = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "outinfo"

	//组织基础信息
	CacheBaseOrgInfo = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "baseinfo2022"
	//组织外部信息
	CacheBaseOrgOutInfo = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "outinfo"
	//获取外部组织id关联的内部组织id
	CacheOutOrgIdRelationId = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOutOrg + "org_id_info"
	//组织付费功能信息
	CacheOrgPayFunction = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "pay_function"

	//获取外部用户id关联的内部用户id
	CacheOutUserIdRelationId = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfOutUser + "user_id"
	//部门对应关系
	CacheDeptRelation = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "dept_relation_list"
	//用户token
	CacheUserToken = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "user:token:"
	//用户开通范围
	CacheUserCheckPay = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "check_pay_2022"
	//fs用户refresh_token和user_access_token
	CacheFsUserAccessToken = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "fstoken"
	//fs用户authCode
	CacheFsAuthCodeToken = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + "fsAuthCode:"

	//飞书组织初始化企业
	CacheFsOrgInit = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + "fs_org_init:"

	//用户邀请code, 拼接 inviteCode
	CacheUserInviteCode = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "invite_code:"
	//用户邀请code有效时间为: 1小时
	CacheUserInviteCodeExpire = 60 * 60 * 1

	// 短信验证码相关 + 手机号, 验证失败五次间隔调整为五分钟，验证失败50次冻结一天
	// 短信发送时间间隔: 一分钟，五分钟
	CacheSmsSendLoginCodeFreezeTime = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfAuthType + consts.CacheKeyOfAddressType + consts.CacheKeyOfPhone + "sms_auth_code:freeze_time"
	// 短信验证码: 五分钟
	CacheSmsLoginCode = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfAuthType + consts.CacheKeyOfAddressType + consts.CacheKeyOfPhone + "sms_auth_code"
	// 号码白名单
	CachePhoneNumberWhiteList = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "sms_white_list"
	// 登录短信验证失败次数
	CacheSmsLoginCodeVerifyFailTimes = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfAuthType + consts.CacheKeyOfAddressType + consts.CacheKeyOfPhone + "sms_auth_code:verify_times"
	//登录图形验证码：一分钟
	CacheLoginGraphCode = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfLoginName + "graph_auth_code"
	//换绑账号行为记录：5分钟
	ChangeLoginNameSign = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + consts.CacheKeyOfAddressType + "change_login_name_sign"
	//钉钉扫码登录：5分钟
	LoginByDingCode = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSourceChannel + consts.CacheKeyOfOutUser + "login_by_ding_code"
	//分享url
	CacheShareUrl = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "user:share:"
	//白名单付费组织
	CacheVipOrg = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "cache_vip_org"
	//飞书福袋用户
	CacheFeishLuckyTagTenant = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "feishu_lucky_tag"
	//灰度企业id
	CacheGrayLevelOrg = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + "gray_level_org"

	//组织开发平台信息
	CacheOrgAppTicket = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "app_ticket"
	//标准版到期前三天提醒
	CachePayExpireRemind = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "pay_expire_remind"
	//标准版已过期提醒
	CachePayOverdueRemind = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "pay_overdue_remind"

	// 账号密码验证失败次数
	CachePwdLoginCodeVerifyFailTimes = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfAuthType + consts.CacheKeyOfAddressType + consts.CacheKeyOfLoginName + "pwd_auth_verify_times"
	CachePwdLoginCodeVerifyFreeze    = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfSys + consts.CacheKeyOfAuthType + consts.CacheKeyOfAddressType + consts.CacheKeyOfLoginName + "verify_times_freeze"

	// 实验室开关缓存
	CacheLabConfig = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + "lab_config"

	// 弹窗卡片config缓存
	CacheCardConfig = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "card_config"

	// 用户的上一次浏览的位置缓存
	CacheUserLocationConfig = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "user_last_location_config"

	// 2022-11-11活动总开关 临时的
	CacheActivity20221111Switch = "activity20221111_switch"

	// 是否是飞书管理员的缓存
	CacheFsPlatformAdmin = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "fs_platform_admin"

	// 检查是否是三方平台管理员
	CachePlatformAdmin = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "platform_admin"
)

const (
	OrgCodeLength = 50
)

var (
	//角色列表
	CacheRoleList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + "role_list"
	////用户角色列表
	//CacheUserRoleList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "user_role_list"
	//角色权限列表
	CacheRolePermissionOperationList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfProject + consts.CacheKeyOfRole + "role_permission_list"
	//角色操作列表
	CacheRoleOperationList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfSys + "role_operation_list"

	//补偿的角色列表
	CacheCompensatoryRolePermissionPathList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + "compensatory_role_permission_path_list"
	//角色组信息
	CacheRoleGroupList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + "role_group_list"
	//权限项列表
	CachePermissionList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfSys + "permission_list"
	//权限项操作列表
	CachePermissionOperationList = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfSys + "permission_operation_list"
	//用户角色列表(hash)
	CacheUserRoleListHash = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfUser + "user_role_list_hash"
	//角色列表
	CacheRoleListHash = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + "role_list_hash"

	//部门角色列表（hash）
	CacheDepartmentRoleListHash = consts.CacheKeyPrefix + consts.RolesvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfDepartment + "department_role_list_hash"

	//组织外部信息
	CacheOrgJSTicket = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + consts.CacheKeyOfOrg + consts.CacheKeyOfJsapiTicket + "js_ticket"
)

var (
	CacheThirdUserInfo       = consts.CacheKeyPrefix + consts.OrgsvcApplicationName + "third_user_info:%s"
	CacheOfficialAccessToken = "polaris:org:official_access_token"
)
