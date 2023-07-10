package bo

//type TodoResult struct {
//	Op  int    `json:"op"`
//	Msg string `json:"msg"`
//}

//type Todo struct {
//	Id                     int64                              `json:"id"`                     // 待办ID
//	OrgId                  int64                              `json:"orgId"`                  // 待办所属-组织ID
//	AppId                  int64                              `json:"appId"`                  // 待办所属-应用ID
//	TableId                int64                              `json:"tableId"`                // 操作对象-表ID
//	IssueId                int64                              `json:"issueId"`                // 操作对象-数据ID
//	WorkflowId             int64                              `json:"workflowId"`             // 所属工作流-ID
//	WorkflowName           string                             `json:"workflowName"`           // 所属工作流-名字
//	ExecutionId            int64                              `json:"executionId"`            // 所属工作流-执行ID
//	TriggerUserId          int64                              `json:"triggerUserId"`          // 所属工作流-发起人ID
//	AllowWithdrawByTrigger int                                `json:"allowWithdrawByTrigger"` // 允许发起人撤回
//	AllowUrgeByTrigger     int                                `json:"allowUrgeByTrigger"`     // 允许发起人催办
//	Type                   int                                `json:"type"`                   // 待办类型-0:审批,1-填写
//	Status                 int                                `json:"status"`                 // 待办状态-0:待处理,1-已处理,2-已失效
//	Parameters             proto.Message                      `json:"parameters"`             // 配置参数
//	Operators              map[int64]*automationPb.TodoResult `json:"operator"`               // 待办处理人
//	CreatedAt              time.Time                          `json:"createdAt"`              // 创建时间
//	UpdatedAt              time.Time                          `json:"updatedAt"`              // 更新时间
//	Creator                int64                              `json:"creator"`                // 创建人
//	Updater                int64                              `json:"updater"`                // 更新人
//}
