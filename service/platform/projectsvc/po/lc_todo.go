package po

//type LcTodo struct {
//	Id                     int64     `db:"id,omitempty" json:"id"`                                            // 待办ID
//	OrgId                  int64     `db:"org_id,omitempty" json:"orgId"`                                     // 待办所属-组织ID
//	AppId                  int64     `db:"app_id,omitempty" json:"appId"`                                     // 待办所属-应用ID
//	TableId                int64     `db:"table_id,omitempty" json:"tableId"`                                 // 操作对象-表ID
//	IssueId                int64     `db:"issue_id,omitempty" json:"issueId"`                                 // 操作对象-数据IssueID
//	WorkflowId             int64     `db:"workflow_id,omitempty" json:"workflowId"`                           // 所属工作流-ID
//	WorkflowName           string    `db:"workflow_name,omitempty" json:"workflowName"`                       // 所属工作流-名字
//	ExecutionId            int64     `db:"execution_id,omitempty" json:"executionId"`                         // 所属工作流-执行ID
//	TriggerUserId          int64     `db:"trigger_user_id,omitempty" json:"triggerUserId"`                    // 所属工作流-发起人ID
//	AllowWithdrawByTrigger int       `db:"allow_withdraw_by_trigger,omitempty" json:"allowWithdrawByTrigger"` // 允许发起人撤回
//	AllowUrgeByTrigger     int       `db:"allow_urge_by_trigger,omitempty" json:"allowUrgeByTrigger"`         // 允许发起人催办
//	Type                   int       `db:"type,omitempty" json:"type"`                                        // 待办类型 0:审批,1-填写
//	Status                 int       `db:"status,omitempty" json:"status"`                                    // 待办状态 0-待处理,1-通过,2-驳回,3-撤回
//	Operators              string    `db:"operators,omitempty" json:"operators"`                              // 待办处理人
//	Parameters             string    `db:"parameters,omitempty" json:"parameters"`                            // 配置参数
//	CreatedAt              time.Time `db:"created_at,omitempty" json:"createdAt"`                             // 创建时间
//	UpdatedAt              time.Time `db:"updated_at,omitempty" json:"updatedAt"`                             // 更新时间
//	Creator                int64     `db:"creator,omitempty" json:"creator"`                                  // 创建人
//	Updater                int64     `db:"updater,omitempty" json:"updater"`                                  // 更新人
//
//}

//
//func (*LcTodo) TableName() string {
//	return "lc_todo"
//}
//
//func (l *LcTodo) ConvertFromBo(b *bo.Todo) {
//	l.Id = b.Id
//	l.OrgId = b.OrgId
//	l.AppId = b.AppId
//	l.TableId = b.TableId
//	l.IssueId = b.IssueId
//	l.WorkflowId = b.WorkflowId
//	l.WorkflowName = b.WorkflowName
//	l.ExecutionId = b.ExecutionId
//	l.TriggerUserId = b.TriggerUserId
//	l.AllowWithdrawByTrigger = b.AllowWithdrawByTrigger
//	l.AllowUrgeByTrigger = b.AllowUrgeByTrigger
//	l.Type = b.Type
//	l.Status = b.Status
//	l.Operators = json.ToJsonIgnoreError(b.Operators)
//	l.Parameters = json.ToJsonIgnoreError(b.Parameters)
//	l.CreatedAt = b.CreatedAt
//	l.UpdatedAt = b.UpdatedAt
//	l.Creator = b.Creator
//	l.Updater = b.Updater
//}
//
//func (l *LcTodo) ConvertToBo() *bo.Todo {
//	b := &bo.Todo{}
//	b.Id = l.Id
//	b.OrgId = l.OrgId
//	b.AppId = l.AppId
//	b.TableId = l.TableId
//	b.IssueId = l.IssueId
//	b.WorkflowId = l.WorkflowId
//	b.WorkflowName = l.WorkflowName
//	b.ExecutionId = l.ExecutionId
//	b.TriggerUserId = l.TriggerUserId
//	b.AllowWithdrawByTrigger = l.AllowWithdrawByTrigger
//	b.AllowUrgeByTrigger = l.AllowUrgeByTrigger
//	b.Type = l.Type
//	b.Status = l.Status
//	b.CreatedAt = l.CreatedAt
//	b.UpdatedAt = l.UpdatedAt
//	b.Creator = l.Creator
//	b.Updater = l.Updater
//	json.Unmarshal(convert.UnsafeStringToBytes(l.Operators), &b.Operators)
//	switch l.Type {
//	case int(automationPb.TodoType_Audit):
//		parameters := &automationPb.ParameterTodoAudit{}
//		json.Unmarshal(convert.UnsafeStringToBytes(l.Parameters), parameters)
//		b.Parameters = parameters
//	case int(automationPb.TodoType_FillIn):
//		parameters := &automationPb.ParameterTodoFillIn{}
//		json.Unmarshal(convert.UnsafeStringToBytes(l.Parameters), parameters)
//		b.Parameters = parameters
//	}
//	return b
//}
