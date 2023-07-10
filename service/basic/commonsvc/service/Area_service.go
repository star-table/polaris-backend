package service

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/commonvo"
	sconsts "github.com/star-table/polaris-backend/service/basic/commonsvc/consts"
	"github.com/star-table/polaris-backend/service/basic/commonsvc/domain"
	"upper.io/db.v3"
)

func OrgAreaInfo(req commonvo.AreaInfoReqVo) (*commonvo.AreaInfoResp, errs.SystemErrorInfo) {

	//&commonvo.AreaInfoResp{
	//	IndustryName:  industryBo.Cname,
	//	CountryCname:  countryBo.Cname,
	//	ProvinceCname: statesBo.Cname,
	//	CityCname:     cityBo.Cname,
	//}
	resp := &commonvo.AreaInfoResp{}

	err := assemblyOrgCountryName(req, resp)

	if err != nil {
		return nil, err
	}

	err = assemblyOrgStatesName(req, resp)

	if err != nil {
		return nil, err
	}

	err = assemblyOrgCityName(req, resp)

	if err != nil {
		return nil, err
	}

	err = assemblyOrgIndustryName(req, resp)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func assemblyOrgCountryName(req commonvo.AreaInfoReqVo, resp *commonvo.AreaInfoResp) errs.SystemErrorInfo {

	if req.CountryID != sconsts.DefalutCountryId {
		countryBo, err := domain.GetCountriesBo(req.CountryID)

		if err != nil {
			return err
		}

		resp.CountryCname = countryBo.Cname
	}
	return nil
}

func assemblyOrgStatesName(req commonvo.AreaInfoReqVo, resp *commonvo.AreaInfoResp) errs.SystemErrorInfo {

	if req.ProvinceID != sconsts.DefalutStatesId {
		statesBo, err := domain.GetStatesBo(req.ProvinceID)

		if err != nil {
			return err
		}

		resp.ProvinceCname = statesBo.Cname
	}
	return nil
}

func assemblyOrgCityName(req commonvo.AreaInfoReqVo, resp *commonvo.AreaInfoResp) errs.SystemErrorInfo {

	if req.CityID != sconsts.DefalutCityId {
		cityBo, err := domain.GetCitiesBo(req.CityID)

		if err != nil {
			return err
		}

		resp.CityCname = cityBo.Cname
	}
	return nil
}

func assemblyOrgIndustryName(req commonvo.AreaInfoReqVo, resp *commonvo.AreaInfoResp) errs.SystemErrorInfo {

	if req.IndustryID != sconsts.DefalutIndustryId {

		industryBo, err := domain.GetIndustryBo(req.IndustryID)

		if err != nil {
			return err
		}

		resp.IndustryName = industryBo.Cname
	}
	return nil

}

func AreaLinkageList(req vo.AreaLinkageListReq) (*vo.AreaLinkageListResp, errs.SystemErrorInfo) {

	if req.IsRoot != nil {
		//获取大陆板块
		return assemblyContinentList()
	} else if req.ContinentID != nil {
		//获取国家
		return assemblyCountryList(*req.ContinentID)
	} else if req.CountryID != nil {
		//获取省
		return assemblyStatesList(*req.CountryID)
	} else if req.StateID != nil {
		//获取城市
		return assemblyCityList(*req.StateID)
	} else if req.CityID != nil {
		//获取县
		return assemblyRegionList(*req.CityID)
	}

	//参数不正确
	return nil, errs.BuildSystemErrorInfo(errs.ReqParamsValidateError)
}

//组装大陆板块返回参数
func assemblyContinentList() (*vo.AreaLinkageListResp, errs.SystemErrorInfo) {

	bos, err := domain.GetContinentList()

	if err != nil {
		return nil, err
	}

	resp := []*vo.AreaLinkageResp{}

	for _, value := range *bos {

		resp = append(resp, &vo.AreaLinkageResp{
			ID:        value.Id,
			Name:      value.Name,
			Cname:     value.Cname,
			IsDefault: value.IsDefault,
		})
	}

	return &vo.AreaLinkageListResp{
		resp,
	}, nil
}

//组装国家返回参数
func assemblyCountryList(continentID int64) (*vo.AreaLinkageListResp, errs.SystemErrorInfo) {

	bos, err := domain.GetCountriesBoListByCond(db.Cond{
		consts.TcIsShow:      consts.AppShowEnable,
		consts.TcContinentId: continentID,
	})

	if err != nil {
		return nil, err
	}

	resp := []*vo.AreaLinkageResp{}

	for _, value := range *bos {

		resp = append(resp, &vo.AreaLinkageResp{
			ID:        value.Id,
			Name:      value.Name,
			Cname:     value.Cname,
			Code:      value.Code,
			IsDefault: value.IsDefault,
		})
	}

	return &vo.AreaLinkageListResp{
		resp,
	}, nil
}

//组装省/州返回参数
func assemblyStatesList(countryID int64) (*vo.AreaLinkageListResp, errs.SystemErrorInfo) {

	bos, err := domain.GetStatesBoListByCond(db.Cond{
		consts.TcIsShow:    consts.AppShowEnable,
		consts.TcCountryId: countryID,
	})

	if err != nil {
		return nil, err
	}

	resp := []*vo.AreaLinkageResp{}

	for _, value := range *bos {

		resp = append(resp, &vo.AreaLinkageResp{
			ID:        value.Id,
			Name:      value.Name,
			Cname:     value.Cname,
			Code:      value.Code,
			IsDefault: value.IsDefault,
		})
	}

	return &vo.AreaLinkageListResp{
		resp,
	}, nil
}

//组装城市的返回参数
func assemblyCityList(stateID int64) (*vo.AreaLinkageListResp, errs.SystemErrorInfo) {

	bos, err := domain.GetCityBoListByCond(db.Cond{
		consts.TcIsShow:  consts.AppShowEnable,
		consts.TcStateId: stateID,
	})

	if err != nil {
		return nil, err
	}

	resp := []*vo.AreaLinkageResp{}

	for _, value := range *bos {

		resp = append(resp, &vo.AreaLinkageResp{
			ID:        value.Id,
			Name:      value.Name,
			Cname:     value.Cname,
			Code:      value.Code,
			IsDefault: value.IsDefault,
		})
	}

	return &vo.AreaLinkageListResp{
		resp,
	}, nil
}

func assemblyRegionList(cityID int64) (*vo.AreaLinkageListResp, errs.SystemErrorInfo) {

	bos, err := domain.GetRegionBoListByCond(db.Cond{
		consts.TcIsShow: consts.AppShowEnable,
		consts.TcCityId: cityID,
	})

	if err != nil {
		return nil, err
	}

	resp := []*vo.AreaLinkageResp{}

	for _, value := range *bos {

		resp = append(resp, &vo.AreaLinkageResp{
			ID:        value.Id,
			Name:      value.Name,
			Cname:     value.Cname,
			Code:      value.Code,
			IsDefault: value.IsDefault,
		})
	}

	return &vo.AreaLinkageListResp{
		resp,
	}, nil
}
