package openapi

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func DecodeAccessToken(accessToken string, secret string) (jwt.MapClaims, error){
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error){
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v ", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil{
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	} else {
		return nil, err
	}
}

func EncodeAccessToken(key, secret string) (string, error){
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	claims["iat"] = time.Now().Unix()
	claims["_key"] = key
	token.Claims = claims
	return token.SignedString([]byte(secret))
}

