package util

import (
	"context"
	"fmt"
	"github.com/star-table/common/core/consts"
	consts2 "github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/core/util/http"
	"github.com/gin-gonic/gin"
)

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value("GinContextKey")
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func GetCtxParameters(ctx context.Context, key string) (string, error) {
	gc, err := GinContextFromContext(ctx)
	if err != nil {
		return "", err
	}
	v := gc.GetString(key)
	return v, nil
}

func GetCtxToken(ctx context.Context) (string, error) {
	return GetCtxParameters(ctx, consts.AppHeaderTokenName)
}

func BuildHeaderOptions(ctx context.Context) ([]http.HeaderOption, error) {
	token, err := GetCtxToken(ctx)
	if err != nil {
		return nil, err
	}
	headerOptions := make([]http.HeaderOption, 0)
	if token != "" {
		headerOptions = append(headerOptions, http.HeaderOption{
			Name:  consts.AppHeaderTokenName,
			Value: token,
		})
	}
	return headerOptions, nil
}

func GetCtxVersion(ctx context.Context) (string, error) {
	return GetCtxParameters(ctx, consts2.AppHeaderVerName)
}
