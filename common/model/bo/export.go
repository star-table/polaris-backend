package bo

type ExportIssueBo struct {
	Id           int64
	ParentId     int64
	IterationId  int64
	ExportValues map[string]interface{}
	Documents    []*DocumentBo
}
