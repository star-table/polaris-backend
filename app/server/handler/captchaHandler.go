package handler

import (
	"github.com/star-table/polaris-backend/common/core/consts"
	"github.com/star-table/polaris-backend/common/model/vo/orgvo"
	"github.com/star-table/polaris-backend/facade/orgfacade"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"strconv"
)

type CaptchaResp struct {
	Code    int32      `json:"code"`
	Message string     `json:"message"`
	Data    CaptchaResponse `json:"data"`
}

type CaptchaResponse struct {
	CaptchaId string `json:"captchaId"` //验证码Id
	ImageUrl  string `json:"imageUrl"`  //验证码图片url
}

type RedisCache struct {
}

func (s *RedisCache) Set(id string, digits []byte) {
	pwd := ""
	for _, b := range digits{
		pwd += strconv.Itoa(int(b))
	}
	resp := orgfacade.SetPwdLoginCode(orgvo.SetPwdLoginCodeReqVo{CaptchaId:id, CaptchaPassword:pwd})
	if resp.Failure() {
		log.Error(resp.Error())
	}
}

func (s *RedisCache) Get(id string, clear bool) (digits []byte)  {
	resp := orgfacade.GetPwdLoginCode(orgvo.GetPwdLoginCodeReqVo{CaptchaId:id})
	if resp.Failure() {
		log.Error(resp.Error())
	}
	res := []byte{}
	for _, i2 := range resp.CaptchaPassword {
		v, _ := strconv.Atoi(string(i2))
		res = append(res, byte(v))
	}
	return res
}

func CaptchaGetHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//length := captcha.DefaultLen
		length := 4
		captchaId := captcha.NewLen(length)
		var captcha CaptchaResponse
		captcha.CaptchaId = captchaId
		captcha.ImageUrl = "captcha/" + captchaId + ".png"
		c.JSON(200, CaptchaResp{Data:captcha, Code:200, Message:"success"})
	}
}

func CaptchaShowHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		width := captcha.StdWidth
		widthParam := c.Query("w")
		if widthParam != consts.BlankString {
			mid, err := strconv.Atoi(widthParam)
			if err != nil {
				log.Error(err)
			} else {
				width = mid
			}
		}

		height := captcha.StdHeight
		heightParam := c.Query("h")
		if heightParam != consts.BlankString {
			mid, err := strconv.Atoi(heightParam)
			if err != nil {
				log.Error(err)
			} else {
				height = mid
			}
		}
		if width > 2000 || height > 2000 {
			c.JSON(500, "You are a bad boy!")
			return
		}
		captcha.Server(width, height).ServeHTTP(c.Writer, c.Request)
	}
}
