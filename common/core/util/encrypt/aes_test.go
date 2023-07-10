package encrypt

import (
	"fmt"
	"testing"
	"time"

	"github.com/star-table/common/core/types"
	"gotest.tools/assert"
)

func TestAesDecrypt(t *testing.T) {

	key := "CMIPOLADMIN"
	encrypt := "yA0zPqk9BBm9Aa0yLN2IB8Tz3zdy2tr+zPb7K201cIWfH0GhnIdRxdGSXCuxdKmmgpfZ9y1JWz+OPUDxRd8ed307L0LKYxcnxMVVAYz6NdjFRGAQ1DcU9RrVDh0DrX3moYXxmK4jqj0U2JPvqYcRx18L8x64qJnR9qKuG/k9+1sDFVZ0lzMOt98S3PG8p4wYYRsJKDrK2w7VJZ7q4vMYntRuGQChSm+rR+ujE1AhCK3RceyRHIyfbkem2CPlMKvdHgzamkVZ8W9OyWOgTSZgDdKQ8qlKd+hM1/22ezKfKDEYAgOXESrEKStEhIMQK1hiPN8hQRVdw5OtkhxuxUWqWQh5YFdebm6boJ97dRQakmF4iqQkMRUX08+KQncyZ0vr60G3Z9HG02eBw5mH7igNuD6lJBqUZRw341wJcn5GI5H4bssNz99XoZ5MyKwJyzwmC2hcwvztSbPX0YSRLqtyZsRvFD5+c2LRAoR57WB6ykZC83xh21Rroace/ezUZAO/IhKOW+V/lwh5A4i5Z+Vkd1bFQaxIyonSQxtPPQi1jKdTJRiT+X2ly65B2X0lcv/5bvj1J0bN7AvUpqEW3efPCQtfJwUmGH67JrkpSQadwN8SKrjOkIWEOXJsDF9H4L/Bs1+GAXyv+vGJslqqrb/ZdMIBQCtWRehlNpQIRSvOnsCaGYFfssMuO7jbBS/U7AX1EUofpyACwrYUzrx+sId0Qg=="

	d, e := AesDecrypt(key, encrypt)
	t.Log(e)
	t.Log(d)

	assert.Equal(t, d, "hello world")
}

//func TestAesEncrypt(t *testing.T) {
//	key := "rME0VeyPmaYOLEOz"
//	obj := projectvo.PushFeishuShortcutMsgVo{
//		Token:  "UDbcxblq3u84EQ0xYmkoPgq0ykwiCqw7",
//		Type:   "application_event",
//		Events: []projectvo.FsShortcutMsgEventData{
//			{
//				Header:     projectvo.FsShortcutHeader{
//					UserId:     "ou_22f5529e4a1bc5c3a1271ca609b24ed5",
//					TenantId:   "2ed263bf32cf1651",
//					HappenedAt: 1604373279000,
//				},
//				TriggerKey: "trgZ7qbJb2bx1Y6A",
//				EventId:    "aaa",
//				Event:      projectvo.ActionFeishuIssueInfoVo{
//					Id:                    "5015",
//					Title:                 "aa",
//					Remark:                "",
//					StatusName:            "未完成",
//					PriorityName:          "较低",
//					PlanStartTime:         "",
//					PlanEndTime:           "",
//					OwnerName:             "樊宇",
//					FollowerName:          []string{},
//					Tags:                  []string{},
//					Url:                   "https://apptest1.bjx.cloud/project/1021/task/0?issueId=5015&platform=FEISHU&orgId=1004",
//					ProjectName:           "哈哈哈",
//					ProjectObjectTypeName: "任务",
//					EventType:             "",
//				},
//			},
//		},
//	}
//
//	t.Log(AesEncrypt(json.ToJsonIgnoreError(obj), key))
//}

func TestAesDecrypt2(t *testing.T) {
	a := types.NowTime()
	a.UnmarshalJSON([]byte("2020-11-26 00:00:00"))

	z, _ := time.LoadLocation("Asia/Shanghai")
	fmt.Println(time.Time(a).In(z).Unix())
	fmt.Println(time.Time(a).In(z))
	//time.Time(a).In(time.FixedZone("UTC", ))
	fmt.Println(time.Time(a).Unix())
}
