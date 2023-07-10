package format

import "regexp"

func VerifyEmailFormat(emails ...string) bool {
	reg := regexp.MustCompile(EmailPattern)
	for _, email := range emails{
		suc := reg.MatchString(email)
		if ! suc {
			return false
		}
	}
	return true
}

// VerifyChinaPhoneFormat 校验是否是合法的中文手机号码
func VerifyChinaPhoneFormat(phone string) bool {
	reg := regexp.MustCompile(ChinaPhoneWithAreaFlagPattern)
	suc := reg.MatchString(phone)
	if !suc {
		return false
	}
	return true
}

// VerifyChinaPhoneWithoutAreaFormat 校验是否是合法的中文手机号码，不带国家代码
func VerifyChinaPhoneWithoutAreaFormat(phone string) bool {
	reg := regexp.MustCompile(ChinaPhoneWithoutAreaPattern)
	suc := reg.MatchString(phone)
	if !suc {
		return false
	}
	return true
}

// VerifyPhoneRegionFormat 验证手机归属地 code 格式
func VerifyPhoneRegionFormat(regionCode string) bool {
	reg := regexp.MustCompile(PhoneRegionCodePattern)
	suc := reg.MatchString(regionCode)
	if !suc {
		return false
	}
	return true
}

func VerifyForeignPhoneNumberFormat(phone string) bool {
	reg := regexp.MustCompile(ForeignPhoneWithoutAreaPattern)
	suc := reg.MatchString(phone)
	if !suc {
		return false
	}
	return true
}

// VerifyWebLinkFormat 校验是否是 web 链接
func VerifyWebLinkFormat(str string) bool {
	reg := regexp.MustCompile(WebLinkPattern)
	suc := reg.MatchString(str)
	if !suc {
		return false
	}
	return true
}

