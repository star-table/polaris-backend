package lc_helper

// DefaultGroupSelectJsonForIssueStatus groupSelect-任务状态分组单选配置
var DefaultGroupSelectJsonForIssueStatus = `[
    {
        "children":[
            {
                "fontColor":"#5f5f5f",
                "id":7,
                "value":"未开始",
                "sort":1,
                "parentId":1,
                "color":"#E1E2E4"
            }
        ],
        "value":"未开始",
        "id":1
    },
    {
        "value":"进行中",
        "id":2,
        "children":[
            {
                "color":"#377AFF",
                "fontColor":"#377aff",
                "id":16,
                "value":"进行中",
                "sort":2,
                "parentId":2
            }
        ]
    },
    {
        "children":[
            {
                "sort":3,
                "parentId":3,
                "color":"#45CB7E",
                "fontColor":"#54a944",
                "id":26,
                "value":"已完成"
            }
        ],
        "value":"已完成",
        "id":3
    }
]`

// DefaultSelectJsonForTaskBar select-任务栏默认值 json
var DefaultSelectJsonForTaskBar = `[
{"id": 1, "value": "任务", "color": ""}
]`

// DefaultSelectJsonForIterationId select-迭代默认值 json
var DefaultSelectJsonForIterationId = `[
{"id": 0, "value": "未规划", "color": ""}
]`

// ProcessFieldValue 进度字段默认值 id 900  数字框
var ProcessFieldValue = `[
    {
        "type":2,
        "value":"percentage",
        "fieldName":"字段格式",
        "id":"5"
    },
    {
        "id":"1",
        "type":4,
        "value":"0",
        "fieldName":"小数点位数"
    }
]`

// StoryPointFieldValue Story Points字段默认值，id 901 单选框
var StoryPointFieldValue = `[
    {
        "fieldName":"选项值",
        "id":1,
        "type":1,
        "value":"无"
    },
    {
        "fieldName":"选项值",
        "id":2,
        "type":1,
        "value":"0"
    },
    {
        "value":"0.5",
        "fieldName":"选项值",
        "id":3,
        "type":1
    },
    {
        "value":"1",
        "fieldName":"选项值",
        "id":4,
        "type":1
    },
    {
        "type":1,
        "value":"2",
        "fieldName":"选项值",
        "id":5
    },
    {
        "id":6,
        "type":1,
        "value":"3",
        "fieldName":"选项值"
    },
    {
        "fieldName":"选项值",
        "id":7,
        "type":1,
        "value":"5"
    },
    {
        "value":"8",
        "fieldName":"选项值",
        "id":8,
        "type":1
    },
    {
        "value":"10",
        "fieldName":"选项值",
        "id":9,
        "type":1
    },
    {
        "type":1,
        "value":"20",
        "fieldName":"选项值",
        "id":10
    },
    {
        "type":1,
        "value":"40",
        "fieldName":"选项值",
        "id":11
    },
    {
        "id":12,
        "type":1,
        "value":"100",
        "fieldName":"选项值"
    }
]`

// ScoreFieldValue 评分字段默认值，id 902 单选框
var ScoreFieldValue = `
[
    {
        "id":1,
        "fieldName":"选项值",
        "value":"1分",
        "type":1
    },
    {
        "fieldName":"选项值",
        "value":"2分",
        "type":1,
        "id":2
    },
    {
        "value":"3分",
        "type":1,
        "id":3,
        "fieldName":"选项值"
    },
    {
        "value":"4分",
        "type":1,
        "id":4,
        "fieldName":"选项值"
    },
    {
        "type":1,
        "id":5,
        "fieldName":"选项值",
        "value":"5分"
    },
    {
        "value":"6分",
        "type":1,
        "id":6,
        "fieldName":"选项值"
    },
    {
        "id":7,
        "fieldName":"选项值",
        "value":"7分",
        "type":1
    },
    {
        "type":1,
        "id":8,
        "fieldName":"选项值",
        "value":"8分"
    },
    {
        "value":"9分",
        "type":1,
        "id":9,
        "fieldName":"选项值"
    },
    {
        "fieldName":"选项值",
        "value":"10分",
        "type":1,
        "id":10
    }
]`
