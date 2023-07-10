package bo

type InviteInfoBo struct {
	InviterId      int64
	OrgId          int64
	SourcePlatform string
}

// 多个手机号邀请时的信息存储
type InviteInfoForPhonesBo struct {
	InviterId      int64  // 邀请人，即当前操作人 current userId
	OrgId          int64
	SourcePlatform string
	BeInvitedPhones []BeInvitedPhone
}

type BeInvitedPhone struct {
	Origin string `json:"origin"`
	Number string `json:"number"`
}
