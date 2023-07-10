package callvo

import (
	"github.com/star-table/polaris-backend/common/model/bo"
)

// GroupChatHandleInfo 处理飞书机器人群聊指令时，用到的数据
type GroupChatHandleInfo struct {
	SourceChannel  string `json:"sourceChannel"`
	TenantClient   interface{}
	DingClient     interface{}
	OrgInfo        bo.BaseOrgInfoBo  `json:"orgInfo"`        // 群聊所在的组织信息
	BindProjectIds []int64           `json:"bindProjectIds"` // 群聊绑定的项目id列表
	IsIssueChat    bool              `json:"isIssueChat"`    // 是否是任务群聊
	IssueInfo      bo.IssueBo        `json:"issueInfo"`      // 临时用的任务信息
	OperateUser    bo.BaseUserInfoBo `json:"operateUser"`    // 操作人信息
	OperatorOpenId string            `json:"operatorOpenId"`
	IssueInfoUrl   string            `json:"issueInfoUrl"` // “查看详情”按钮链接
	IssuePcUrl     string            `json:"issuePcUrl"`   // "应用内查看"按钮链接
}

type MessageReqData struct {
	AppId     string `json:"app_id"`
	TenantKey string `json:"tenant_key"`
	Type      string `json:"type"`

	RootId           string   `json:"root_id"`
	ParentId         string   `json:"parent_id"`
	OpenChatId       string   `json:"open_chat_id"`
	ChatType         string   `json:"chat_type"`
	MsgType          string   `json:"msg_type"`
	OpenId           string   `json:"open_id"`
	OpenMessageId    string   `json:"open_message_id"`
	IsMention        bool     `json:"is_mention"`
	Text             string   `json:"text"`
	TextWithoutAtBot string   `json:"text_without_at_bot"`
	Title            string   `json:"title"`
	ImageKeys        []string `json:"image_keys"`
}
