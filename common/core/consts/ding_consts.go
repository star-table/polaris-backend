package consts

// 钉钉订单类型
const (
	DingOrderBy           = "BUY"           // 新购
	DingOrderRenew        = "RENEW"         // 续费
	DingOrderUpgrade      = "UPGRADE"       // 升级
	DingOrderRenewUpgrade = "RENEW_UPGRADE" // 续费升配
	DingOrderRenewDegrade = "RENEW_DEGRADE" // 续费降配

	DingOrderTrail = 15 // 免费试用15天旗舰版功能

	// 针对试用规格
	DingOrderChargeTryout = "TRYOUT" // 试用开通
	DingOrderChargeFree   = "FREE"   // 免费开通

)
