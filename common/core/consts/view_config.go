package consts

const (
	MyOverDueViewConfig = `{
    "tableId":"0",
    "condition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "fieldType":"multiStatus",
                "option":[
                    {
                        "color":"#f5f6f5",
                        "label":"未开始",
                        "value":-5,
                        "fontColor":"#5f5f5f"
                    },
                    {
                        "color":"#ecf5fc",
                        "label":"进行中",
                        "value":-3,
                        "fontColor":"#377aff"
                    },
                    {
                        "color":"#f5f6f5",
                        "label":"待确认",
                        "value":-1,
                        "fontColor":"#5f5f5f"
                    }
                ]
            },
            {
                "type":"lt",
                "column":"planEndTime",
                "value":"${today}",
                "fieldType":"datepicker",
                "props":{
                    "DatePicker":{
                        "selectType":"${today}"
                    }
                }
            },
            {
                "type":"values_in",
                "column":"ownerId",
                "value":[
                    {
                        "id":"${current_user}",
                        "name":"本人",
                        "type":"U_",
                        "avatar":"https://static-polaris-hd2.startable.cn/front_resources/picture.svg"
                    }
                ],
                "fieldType":"member"
            }
        ]
    },
    "realCondition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "values":[
                    -5,
                    -3,
                    -1
                ]
            },
            {
                "type":"and",
                "column":"planEndTime",
                "conds":[
                    {
                        "type":"lt",
                        "value":"${today}",
                        "column":"planEndTime"
                    },
                    {
                        "type":"gt",
                        "column":"planEndTime",
                        "value":"1970-01-01 00:00:00"
                    }
                ]
            },
            {
                "type":"values_in",
                "column":"ownerId",
                "value":[
                    "U_${current_user}"
                ],
                "values":[
                    "U_${current_user}"
                ]
            }
        ]
    },
    "projectObjectTypeId":0
}`

	MyDailyDueOfTodayViewConfig = `{
    "tableId":"0",
    "condition":{
        "type":"and",
        "conds":[
            {
                "type":"equal",
                "column":"planEndTime",
                "value":"${today}",
                "fieldType":"datepicker",
                "props":{
                    "DatePicker":{
                        "selectType":"${today}"
                    }
                }
            },
            {
                "type":"values_in",
                "column":"ownerId",
                "value":[
                    {
                        "id":"${current_user}",
                        "name":"本人",
                        "type":"U_",
                        "avatar":"https://static-polaris-hd2.startable.cn/front_resources/picture.svg"
                    }
                ],
                "fieldType":"member"
            },
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "fieldType":"multiStatus",
                "option":[
                    {
                        "color":"#f5f6f5",
                        "label":"未开始",
                        "value":-5,
                        "fontColor":"#5f5f5f"
                    },
                    {
                        "color":"#ecf5fc",
                        "label":"进行中",
                        "value":-3,
                        "fontColor":"#377aff"
                    },
                    {
                        "color":"#f5f6f5",
                        "label":"待确认",
                        "value":-1,
                        "fontColor":"#5f5f5f"
                    }
                ]
            }
        ]
    },
    "realCondition":{
        "type":"and",
        "conds":[
            {
                "type":"equal",
                "column":"planEndTime",
                "value":"${today}",
                "values":[
                    "${today}"
                ]
            },
            {
                "type":"values_in",
                "column":"ownerId",
                "value":[
                    "U_${current_user}"
                ],
                "values":[
                    "U_${current_user}"
                ]
            },
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "values":[
                    -5,
                    -3,
                    -1
                ]
            }
        ]
    },
    "projectObjectTypeId":0
}`

	OverDueSoonViewConfig = `{
    "tableId":"0",
    "condition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "fieldType":"multiStatus",
                "option":[
                    {
                        "color":"#f5f6f5",
                        "label":"未开始",
                        "value":-5,
                        "fontColor":"#5f5f5f"
                    },
                    {
                        "color":"#ecf5fc",
                        "label":"进行中",
                        "value":-3,
                        "fontColor":"#377aff"
                    },
                    {
                        "color":"#f5f6f5",
                        "label":"待确认",
                        "value":-1,
                        "fontColor":"#5f5f5f"
                    }
                ]
            },
            {
                "type":"gt",
                "column":"planEndTime",
                "value":"${today}",
                "fieldType":"datepicker",
                "props":{
                    "DatePicker":{
                        "selectType":"${today}"
                    }
                }
            },
            {
                "type":"lt",
                "column":"planEndTime",
                "value":"${afterday:3}",
                "fieldType":"datepicker",
                "props":{
                    "DatePicker":{
                        "selectType":"${afterday:N}"
                    }
                }
            },
            {
                "type":"values_in",
                "column":"ownerId",
                "value":[
                    {
                        "id":"${current_user}",
                        "name":"本人",
                        "type":"U_",
                        "avatar":"https://static-polaris-hd2.startable.cn/front_resources/picture.svg"
                    }
                ],
                "fieldType":"member"
            }
        ]
    },
    "groupInfoList":[

    ],
    "realCondition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "values":[
                    -5,
                    -3,
                    -1
                ]
            },
            {
                "type":"gt",
                "column":"planEndTime",
                "value":"${today}",
                "values":[
                    "${today}"
                ]
            },
            {
                "type":"and",
                "column":"planEndTime",
                "conds":[
                    {
                        "type":"lt",
                        "value":"${afterday:3}",
                        "column":"planEndTime"
                    },
                    {
                        "type":"gt",
                        "column":"planEndTime",
                        "value":"1970-01-01 00:00:00"
                    }
                ]
            },
            {
                "type":"values_in",
                "column":"ownerId",
                "value":[
                    "U_${current_user}"
                ],
                "values":[
                    "U_${current_user}"
                ]
            }
        ]
    },
    "projectObjectTypeId":0
}`

	MyPendingViewConfig = `{
    "tableId":"0",
    "condition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "fieldType":"multiStatus",
                "option":[
                    {
                        "color":"#f5f6f5",
                        "label":"未开始",
                        "value":-5,
                        "fontColor":"#5f5f5f"
                    },
                    {
                        "color":"#ecf5fc",
                        "label":"进行中",
                        "value":-3,
                        "fontColor":"#377aff"
                    },
                    {
                        "color":"#f5f6f5",
                        "label":"待确认",
                        "value":-1,
                        "fontColor":"#5f5f5f"
                    }
                ]
            },
            {
                "type":"values_in",
                "column":"ownerId",
                "value":[
                    {
                        "id":"${current_user}",
                        "name":"本人",
                        "type":"U_",
                        "avatar":"https://static-polaris-hd2.startable.cn/front_resources/picture.svg"
                    }
                ],
                "fieldType":"member"
            }
        ]
    },
    "ColumnWidth":{
        "tableColumnWidthtitle":253
    },
    "realCondition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "values":[
                    -5,
                    -3,
                    -1
                ]
            },
            {
                "type":"values_in",
                "column":"ownerId",
                "value":[
                    "U_${current_user}"
                ],
                "values":[
                    "U_${current_user}"
                ]
            }
        ]
    },
    "projectObjectTypeId":0
}`

	AllOverDueSoonViewConfig = `{
    "tableId":"0",
    "condition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "fieldType":"multiStatus",
                "option":[
                    {
                        "color":"#f5f6f5",
                        "label":"未开始",
                        "value":-5,
                        "fontColor":"#5f5f5f"
                    },
                    {
                        "color":"#ecf5fc",
                        "label":"进行中",
                        "value":-3,
                        "fontColor":"#377aff"
                    },
                    {
                        "color":"#f5f6f5",
                        "label":"待确认",
                        "value":-1,
                        "fontColor":"#5f5f5f"
                    }
                ]
            },
            {
                "type":"gt",
                "column":"planEndTime",
                "value":"${today}",
                "fieldType":"datepicker",
                "props":{
                    "DatePicker":{
                        "selectType":"${today}"
                    }
                }
            },
            {
                "type":"lt",
                "column":"planEndTime",
                "value":"${afterday:3}",
                "fieldType":"datepicker",
                "props":{
                    "DatePicker":{
                        "selectType":"${afterday:N}"
                    }
                }
            }
        ]
    },
    "ColumnWidth":{
        "tableColumnWidthtitle":337,
        "tableColumnWidthownerId":195,
        "tableColumnWidthprojectId":262
    },
    "groupInfoList":[

    ],
    "realCondition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "values":[
                    -5,
                    -3,
                    -1
                ]
            },
            {
                "type":"gt",
                "column":"planEndTime",
                "value":"${today}",
                "values":[
                    "${today}"
                ]
            },
            {
                "type":"and",
                "column":"planEndTime",
                "conds":[
                    {
                        "type":"lt",
                        "value":"${afterday:3}",
                        "column":"planEndTime"
                    },
                    {
                        "type":"gt",
                        "column":"planEndTime",
                        "value":"1970-01-01 00:00:00"
                    }
                ]
            }
        ]
    },
    "projectObjectTypeId":0
}`

	AllOverDueViewConfig = `{
    "tableId":"0",
    "condition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "fieldType":"multiStatus",
                "option":[
                    {
                        "color":"#f5f6f5",
                        "label":"未开始",
                        "value":-5,
                        "fontColor":"#5f5f5f"
                    },
                    {
                        "color":"#ecf5fc",
                        "label":"进行中",
                        "value":-3,
                        "fontColor":"#377aff"
                    },
                    {
                        "color":"#f5f6f5",
                        "label":"待确认",
                        "value":-1,
                        "fontColor":"#5f5f5f"
                    }
                ]
            },
            {
                "type":"lt",
                "column":"planEndTime",
                "value":"${today}",
                "fieldType":"datepicker",
                "props":{
                    "DatePicker":{
                        "selectType":"${today}"
                    }
                }
            }
        ]
    },
    "tableOrder":[
        "title",
        "code",
        "ownerId",
        "issueStatus",
        "projectObjectTypeId",
        "planStartTime",
        "planEndTime",
        "iterationId",
        "remark",
        "followerIds",
        "auditorIds",
        "projectId"
    ],
    "ColumnWidth":{
        "tableColumnWidthprojectId":348,
        "tableColumnWidthauditorIds":225
    },
    "groupInfoList":[

    ],
    "realCondition":{
        "type":"and",
        "conds":[
            {
                "type":"in",
                "column":"issueStatus",
                "value":[
                    -5,
                    -3,
                    -1
                ],
                "values":[
                    -5,
                    -3,
                    -1
                ]
            },
            {
                "type":"and",
                "column":"planEndTime",
                "conds":[
                    {
                        "type":"lt",
                        "value":"${today}",
                        "column":"planEndTime"
                    },
                    {
                        "type":"gt",
                        "column":"planEndTime",
                        "value":"1970-01-01 00:00:00"
                    }
                ]
            }
        ]
    },
    "hiddenColumnIds":[

    ],
    "projectObjectTypeId":0
}`

	OverDueViewConfig = `{
    "tableId": "0",
    "condition": {
        "type": "and",
        "conds": [
            {
                "type": "not_in",
                "column": "issueStatus",
                "value": [
                    -2
                ],
                "fieldType": "multiStatus",
                "option": [
                    {
                        "color": "#edf8ed",
                        "label": "已完成",
                        "value": -2,
                        "fontColor": "#54a944"
                    }
                ]
            },
            {
                "type": "lt",
                "column": "planEndTime",
                "value": "${today}",
                "fieldType": "datepicker",
                "props": {
                    "DatePicker": {
                        "selectType": "${today}"
                    }
                }
            },
            {
                "type": "not_null",
                "column": "planEndTime",
                "fieldType": "datepicker"
            }
        ]
    },
    "realCondition": {
        "type": "and",
        "conds": [
            {
                "type": "not_in",
                "column": "issueStatus",
                "value": [
                    -2
                ],
                "values": [
                    -2
                ]
            },
            {
                "type": "lt",
                "column": "planEndTime",
                "value": "${today}",
                "values": [
                    "${today}"
                ]
            },
            {
                "type": "gt",
                "column": "planEndTime",
                "value": "1971-01-01 00:00:00"
            }
        ]
    },
    "projectObjectTypeId": 0
}`

	AllUnFinishedViewConfig = `{
    "group":{
        "id":0,
        "key":0,
        "field":{
            "type":"",
            "props":{

            }
        },
        "title":"不分组",
        "i18_name":"不分组"
    },
    "tableId":"0",
    "taskRank":0,
    "condition":{
        "column":"issueStatus",
        "type":"in",
        "values":[
            [
                -3,
                -5,
                -1
            ]
        ],
        "fieldType":"multiStatus",
        "option":[
            {
                "color":"#ecf5fc",
                "label":"进行中",
                "value":-3,
                "fontColor":"#377aff"
            },
            {
                "color":"#f5f6f5",
                "label":"未开始",
                "value":-5,
                "fontColor":"#5f5f5f"
            },
            {
                "color":"#f5f6f5",
                "label":"待确认",
                "value":-1,
                "fontColor":"#5f5f5f"
            }
        ]
    },
    "orderParams":[
        {
            "asc":false,
            "column":"sort"
        }
    ],
    "groupInfoList":[

    ],
    "realCondition":{
        "column":"issueStatus",
        "type":"in",
        "values":[
            -3,
            -5,
            -1
        ],
        "value":[
            -3,
            -5,
            -1
        ]
    },
    "projectObjectTypeId":0,
    "selectPresentationForm":0
}`

	OwnerViewConfig = `{
    "condition": {
        "type": "values_in",
        "column": "ownerId",
        "values": [
            [
                {
                    "id": "${current_user}",
                    "name": "本人",
                    "avatar": "https://polaris-hd2.oss-cn-shanghai.aliyuncs.com/front_resources/picture.svg"
                }
            ]
        ],
        "fieldType": "member"
    },
    "realCondition": {
        "type": "values_in",
        "value": [
            "U_${current_user}"
        ],
        "column": "ownerId",
        "values": [
            "U_${current_user}"
        ]
    },
    "projectObjectTypeId": 0,
    "tableId": "0"
}`

	AllIssueViewConfig = `{
    "condition": {},
    "realCondition": {},
    "projectObjectTypeId": 0,
    "tableId": "0"
}`
)
