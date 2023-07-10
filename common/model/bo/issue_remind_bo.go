package bo

type RemindConfig struct {
	//提醒類型，1 截止時間，2，开始时间
	RemindType int `json:"remindType"`
	//时间点, 在截止时间/开始时间前/后提醒，单位：分钟
	RemindTime int `json:"remindTime"`
	//配置标识, 对于缓存使用，唯一
	ID string `json:"id"`
}