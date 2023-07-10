package service

import (
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/common/library/db/mysql"
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/errs"
	"github.com/star-table/polaris-backend/common/model/bo"
	"github.com/star-table/polaris-backend/common/model/vo"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/domain"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/po"
	"github.com/star-table/polaris-backend/service/platform/orgsvc/util/openapi"
	"github.com/dgrijalva/jwt-go"
	"time"
	"upper.io/db.v3"
)

func OpenAPIAuth(req orgvo.OpenAPIAuthReq) (*orgvo.OpenAPIAuthData, errs.SystemErrorInfo){
	accessToken := req.AccessToken
	if accessToken == ""{
		return nil, errs.OpenAccessTokenIsEmpty
	}
	secrets := make([]po.PpmOrgSecret, 0)
	err := mysql.SelectAllByCond(consts.TableOrgSecret, db.Cond{
		"org_id": req.OrgID,
		"is_delete": consts.AppIsNoDelete,
	}, &secrets)
	if err != nil{
		log.Error(err)
		return nil, errs.SystemError
	}
	if len(secrets) == 0{
		log.Errorf("企业未配置secret, orgId: %d", req.OrgID)
		return nil, errs.OpenAccessTokenInvalid
	}

	var targetClaims jwt.MapClaims
	var targetSecret po.PpmOrgSecret
	for _, secret := range secrets{
		claims, err := openapi.DecodeAccessToken(accessToken, secret.Secret)
		if err == nil{
			targetClaims = claims
			targetSecret = secret
			break
		}
	}
	if targetClaims == nil{
		return nil, errs.OpenAccessTokenInvalid
	}
	authInfoJson := json.ToJsonIgnoreError(targetClaims)
	authInfo := bo.OpenAuthInfo{}
	_ = json.FromJson(authInfoJson, &authInfo)

	if authInfo.Key != targetSecret.Key{
		return nil, errs.OpenAccessTokenInvalid
	}

	if authInfo.Exp > 0 && time.Now().UTC().After(time.Unix(authInfo.Exp, 0)){
		return nil, errs.OpenAccessTokenExpired
	}

	return &orgvo.OpenAPIAuthData{
		OrgID: req.OrgID,
	}, nil
}

func GetAppTicket(input orgvo.GetAppTicketReq) (*vo.GetAppTicketResp, errs.SystemErrorInfo) {
	orgId := input.OrgId
	userId := input.UserId

	//校验当前用户是否具有"团队设置"的权限
	authErr := AuthOrgRole(orgId, userId, consts.RoleOperationPathOrgUser, consts.OperationOrgConfigModify)
	if authErr != nil {
		log.Error(authErr)
		return nil, authErr
	}

	info, err := domain.GetOrgAppTicket(orgId)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	return &vo.GetAppTicketResp{
		AppID:     info.AppId,
		AppSecret: info.AppSecret,
	}, nil
}
