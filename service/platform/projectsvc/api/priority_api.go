package api

//func (PostGreeter) PriorityList(input projectvo.PriorityListReqVo) projectvo.PriorityListRespVo {
//	cond := db.Cond{}
//	if input.Type != nil {
//		cond[consts.TcType] = *input.Type
//	}
//	page := input.Page
//	size := input.Size
//
//	pageA := uint(0)
//	sizeA := uint(0)
//	if page != nil && size != nil && *page > 0 && *size > 0 {
//		pageA = uint(*page)
//		sizeA = uint(*size)
//	}
//
//	list, err := service2.PriorityList(input.OrgId, pageA, sizeA, cond)
//	return projectvo.PriorityListRespVo{Err: vo.NewErr(err), PriorityList: list}
//}
//
//func (PostGreeter) CreatePriority(req projectvo.CreatePriorityReqVo) vo.CommonRespVo {
//	res, err := service2.CreatePriority(req.UserId, req.CreatePriorityReq)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}
//
//func (PostGreeter) UpdatePriority(req projectvo.UpdatePriorityReqVo) vo.CommonRespVo {
//	res, err := service2.UpdatePriority(req.UserId, req.UpdatePriorityReq)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}
//
//func (PostGreeter) DeletePriority(req projectvo.DeletePriorityReqVo) vo.CommonRespVo {
//	res, err := service2.DeletePriority(req.OrgId, req.UserId, req.DeletePriorityReq)
//	return vo.CommonRespVo{Err: vo.NewErr(err), Void: res}
//}
//
//func (PostGreeter) VerifyPriority(req projectvo.VerifyPriorityReqVo) (resp projectvo.VerifyPriorityRespVo) {
//	res, err := service2.VerifyPriority(req.OrgId, req.Typ, req.BeVerifyId)
//	return projectvo.VerifyPriorityRespVo{Successful: res, Err: vo.NewErr(err)}
//}

//func (PostGreeter) InitPriority(req projectvo.InitPriorityReqVo) vo.VoidErr {
//	err := service2.InitPriority(req.OrgId)
//	return vo.VoidErr{Err: vo.NewErr(err)}
//}
//
//func (GetGreeter) GetPriorityById(req projectvo.GetPriorityByIdReqVo) projectvo.GetPriorityByIdRespVo {
//	res, err := service2.GetPriorityById(req.OrgId, req.Id)
//	return projectvo.GetPriorityByIdRespVo{PriorityBo: res, Err: vo.NewErr(err)}
//}
