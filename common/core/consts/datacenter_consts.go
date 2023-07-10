package consts

const (
	TableType = "schema"

	ConditionEqual   = "equal"
	ConditionUnEqual = "un_equal"
	ConditionIn      = "in"
	ConditionLike    = "like"
	ConditionAnd     = "and"

	//RowKeyId          = "id"
	//RowKeyDelFlag     = "delFlag"
	//RowKeyRowId       = "issueId"
	//RowKeyOrgId       = "orgId"
	//RowKeyIssueStatus = "issueStatus"
	//RowKeyPriorityId  = "priorityId"
	//RowKeyTableId     = "tableId"

	DeleteFlagDel    = 1
	DeleteFlagNotDel = 2

	// Set
	SetTypeNormal = 1
	SetTypeJson   = 2

	SetActionSet              = 1
	SetActionJsonArrayAddItem = 2
	SetActionJsonArrayDelItem = 3
	SetActionJsonMapSet       = 4
)

const (
	DataSource = 2 // pg使用的id
	DataBase   = 1 // pg使用的id
)
