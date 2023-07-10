package dao

//func GetIssueCountByObjectType(orgId, projectId int64) ([]po.IssueStatByObjectType, error) {
//	stat := &[]po.IssueStatByObjectType{}
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		return *stat, err
//	}
//
//	err = conn.Select(db.Raw("i.project_object_type_id, count(i.id) as count, o.lang_code")).From("ppm_pri_issue as i").
//		LeftJoin("ppm_prs_project_object_type as o").On("i.project_object_type_id = o.id").Where(db.Cond{
//		"i.org_id":     orgId,
//		"i.project_id": projectId,
//		"i.is_delete":  consts.AppIsNoDelete,
//	}).GroupBy("i.project_object_type_id").All(stat)
//
//	if err != nil {
//		return *stat, err
//	}
//
//	return *stat, nil
//}

//func StatIssue(orgId, projectId, projectObjectTypeId int64) ([]po.IssueStatByStatus, error) {
//	issueStatByStatus := &[]po.IssueStatByStatus{}
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		return *issueStatByStatus, err
//	}
//
//	err = conn.Select(db.Raw("status, count(*) as count")).From(consts.TableIssue).Where(db.Cond{
//		consts.TcOrgId:               orgId,
//		consts.TcProjectId:           projectId,
//		consts.TcProjectObjectTypeId: projectObjectTypeId,
//	}).GroupBy("status").All(issueStatByStatus)
//
//	if err != nil {
//		return *issueStatByStatus, err
//	}
//
//	return *issueStatByStatus, nil
//}

//func IssueStatusStat(orgId int64, projectId *int64, iterationId *int64) ([]po.IssueStatByProjectIdAndObjectId, error) {
//	stat := &[]po.IssueStatByProjectIdAndObjectId{}
//
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		return *stat, err
//	}
//
//	cond := db.Cond{
//		"issue.org_id":    orgId,
//		"issue.status":    db.Raw("sta.id"),
//		"issue.is_delete": consts.AppIsNoDelete,
//	}
//	if projectId != nil {
//		cond["issue.project_id"] = projectId
//	}
//	if iterationId != nil {
//		cond["issue.iteration_id"] = iterationId
//	}
//
//	err = conn.Select(db.Raw("count(issue.id) AS count,sta.type as status_type")).
//		From("ppm_pri_issue issue", "ppm_prs_process_status sta").
//		Where(cond).GroupBy("sta.type").All(stat)
//
//	if err != nil {
//		return *stat, err
//	}
//
//	return *stat, nil
//}
//
//func IssueStatusStatDetail(orgId int64, projectId *int64, iterationId *int64) ([]bo.IssueStatByProjectIdAndObjectId, error) {
//	stat := &[]bo.IssueStatByProjectIdAndObjectId{}
//	issueStatBos, err1 := issuedomain.GetIssueStatusStat(bo.IssueStatusStatCondBo{
//		OrgId: orgId,
//		ProjectId: projectId,
//		IterationId: iterationId,
//	})
//	if err1 != nil{
//		log.Error(err1)
//		return nil, errs.BuildSystemErrorInfo(errs.IssueDomainError, err1)
//	}
//
//	conn, err := mysql.GetConnect()
//	if err != nil {
//		return *stat, err
//	}
//
//	cond := db.Cond{
//		"issue.org_id":                 orgId,
//		"issue.status":                 db.Raw("sta.id"),
//		"issue.project_object_type_id": db.Raw("typ.id"),
//		"issue.is_delete":              consts.AppIsNoDelete,
//	}
//	if projectId != nil {
//		cond["issue.project_id"] = projectId
//	}
//
//	err = conn.Select(db.Raw("count(issue.id) AS count,typ.name AS name,sta.type as status_type")).
//		From("ppm_pri_issue issue", "ppm_prs_process_status sta", "ppm_prs_project_object_type typ").
//		Where(cond).GroupBy("issue.project_object_type_id", "sta.type").All(stat)
//
//	if err != nil {
//		return *stat, err
//	}
//
//	return *stat, nil
//}

//func JudgeIssueIsExist(orgId, id int64) bool {
//	issue := &po.PpmPriIssue{}
//	err := mysql.SelectOneByCond(consts.TableIssue, db.Cond{
//		consts.TcId:       id,
//		consts.TcOrgId:    orgId,
//		consts.TcIsDelete: consts.AppIsNoDelete,
//	}, issue)
//	if err != nil {
//		log.Error(err)
//		return false
//	}
//	return true
//}
