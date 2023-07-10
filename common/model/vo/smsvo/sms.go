package smsvo

import "github.com/star-table/polaris-backend/common/model/vo"

type AddSmsTaskReq struct {
    Input AddSmsTaskReqData `json:"input"`
}

type AddSmsTaskReqData struct {
    TaskName string `json:"taskName"`
    PhoneNumbers []string `json:"phoneNumbers"`
    SignName string `json:"signName"`
    SmsTemplateCode string `json:"smsTemplateCode"`
    SmsInfo string `json:"smsInfo"`
    ProviderName int `json:"providerName"`
    DataSource int `json:"dataSource"`
    SmsTemplateVariate SmsTemplateVariate `json:"smsTemplateVariate"`
}

type SmsTemplateVariate struct {
    Code string `json:"code"`
}

type AddSmsTaskResp struct {
    vo.Err
    Data AddSmsTaskData `json:"data"`
}

type AddSmsTaskData struct {
    Success bool `json:"success"`
    Message string `json:"message"`
}