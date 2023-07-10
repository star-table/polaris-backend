package rand

import (
	"encoding/base64"
	"github.com/star-table/common/core/util/strs"
	"math/rand"
	"strings"
	"time"
)

func RandomVerifyCode(l int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < l; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func RandomInviteCode(str string) string {
	afterBase64 := base64.StdEncoding.EncodeToString([]byte(str))
	//afterUpper := strings.ToUpper(afterBase64)
	afterFilter := strings.ReplaceAll(afterBase64, "=", "")
	if strs.Len(afterFilter) > 8 {
		afterFilter = afterFilter[0:8]
	}
	return afterFilter
}
