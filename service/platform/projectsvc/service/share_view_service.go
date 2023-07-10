package service

import (
	"fmt"
	"math/rand"

	"github.com/star-table/common/core/util/md5"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/vo/projectvo"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/dao"
	"github.com/star-table/polaris-backend/service/platform/projectsvc/po"
	"github.com/spf13/cast"
	"upper.io/db.v3"
)

const (
	randomKey     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randKeyLength = 8
)

func CreateShareView(req *projectvo.CreateShareViewReq) (*projectvo.ShareViewInfo, errs.SystemErrorInfo) {
	info, err := GetShareViewInfo(&projectvo.GetShareViewInfoReq{
		OrgId:  req.OrgId,
		UserId: req.UserId,
		Input:  &projectvo.ShareViewIdData{ViewId: req.Input.ViewId},
	})
	if err != nil && err.Code() != errs.ShareViewNotExist.Code() {
		return nil, err
	}

	// 不存在，则创建
	if info == nil {
		shareKey := randKey()
		daoErr := dao.GetShareView().Create(&po.PpmShareView{
			ShareKey:  shareKey,
			OrgId:     req.OrgId,
			UserId:    req.UserId,
			ProjectId: req.Input.ProjectId,
			AppId:     req.Input.AppId,
			ViewId:    req.Input.ViewId,
			TableId:   req.Input.TableId,
			Config:    "{}",
		})
		if daoErr != nil {
			return nil, errs.MysqlOperateError
		}
		info = &projectvo.ShareViewInfo{ShareKey: shareKey, Config: "{}"}
	}

	return info, nil
}

func DeleteShareView(req *projectvo.DeleteShareViewReq) errs.SystemErrorInfo {
	err := dao.GetShareView().Delete(req.UserId, req.Input.ViewId)
	if err != nil {
		return errs.MysqlOperateError
	}

	return nil
}

func randKey() string {
	shareKey := ""
	all := len(randomKey)
	for i := 0; i < randKeyLength; i++ {
		shareKey += fmt.Sprintf("%c", randomKey[rand.Intn(all)])
	}

	return shareKey
}

func ResetShareKey(req *projectvo.ResetShareKeyReq) (*projectvo.ShareViewInfo, errs.SystemErrorInfo) {
	shareKey := randKey()
	err := dao.GetShareView().ResetShareKeyAndPassword(req.OrgId, req.UserId, req.Input.ViewId, shareKey)
	if err != nil {
		return nil, errs.MysqlOperateError
	}
	return &projectvo.ShareViewInfo{
		ShareKey: shareKey,
	}, nil
}

func UpdateSharePassword(req *projectvo.UpdateSharePasswordReq) errs.SystemErrorInfo {
	if req.Input.Password != "" {
		req.Input.Password = md5.Md5V(req.Input.Password)
	}
	err := dao.GetShareView().UpdatePassword(req.OrgId, req.UserId, req.Input.ViewId, req.Input.Password)
	if err != nil {
		return errs.MysqlOperateError
	}
	return nil
}

func UpdateShareConfig(req *projectvo.UpdateShareConfigReq) errs.SystemErrorInfo {
	err := dao.GetShareView().UpdateShareConfig(req.OrgId, req.UserId, req.Input.ViewId, req.Input.Config)
	if err != nil {
		return errs.MysqlOperateError
	}
	return nil
}

func GetShareViewInfo(req *projectvo.GetShareViewInfoReq) (*projectvo.ShareViewInfo, errs.SystemErrorInfo) {
	info, err := dao.GetShareView().GetByViewId(req.OrgId, req.UserId, req.Input.ViewId)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.ShareViewNotExist
		}
		return nil, errs.MysqlOperateError
	}

	return &projectvo.ShareViewInfo{
		ShareKey:      info.ShareKey,
		Config:        info.Config,
		IsSetPassword: info.SharePassword != "",
	}, nil
}

func GetShareViewInfoByKey(req *projectvo.GetShareViewInfoByKeyReq) (*projectvo.ShareViewInfo, errs.SystemErrorInfo) {
	info, err := dao.GetShareView().GetByShareKey(req.Input.ShareKey)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.ShareViewNotExist
		}
		return nil, errs.MysqlOperateError
	}

	return &projectvo.ShareViewInfo{
		ShareKey:      info.ShareKey,
		Config:        info.Config,
		IsSetPassword: info.SharePassword != "",
		AppId:         cast.ToString(info.AppId),
		ViewId:        cast.ToString(info.ViewId),
		TableId:       cast.ToString(info.TableId),
		ProjectId:     info.ProjectId,
	}, nil
}

func CheckShareViewPassword(req *projectvo.CheckPasswordData) (*projectvo.CheckPassword, errs.SystemErrorInfo) {
	info, err := dao.GetShareView().GetByShareKey(req.ShareKey)
	if err != nil {
		if err == db.ErrNoMoreRows {
			return nil, errs.ShareViewNotExist
		}
		return nil, errs.MysqlOperateError
	}

	if info.SharePassword == "" || info.SharePassword == md5.Md5V(req.Password) {
		return &projectvo.CheckPassword{IsCorrect: true, OrgId: info.OrgId, UserId: info.UserId}, nil
	}

	return nil, errs.ShareViewPasswordError
}
