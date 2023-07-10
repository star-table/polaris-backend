package msgvo

type SendLoginSMSReqVo struct {
	PhoneNumber string `json:"phoneNumber"`
	Code string `json:"code"`
}

type SendMailReqVo struct {
	Input SendMailReqData `json:"input"`
}

type SendMailReqData struct {
	Emails []string `json:"emails"`
	//标题
	Subject string `json:"subject"`
	//内容
	Content string `json:"content"`
}

type SendSMSReqVo struct {
	Input SendSMSReqVoReqData `json:"input"`
}

type SendSMSReqVoReqData struct {
	Mobile string `json:"mobile"`
	Params map[string]string `json:"params"`
	SignName string `json:"signName"`
	TemplateCode string `json:"templateCode"`
}