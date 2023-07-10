package str

import (
	"fmt"
	"testing"
)

func TestUrlParse(t *testing.T) {
	fmt.Println(UrlParse("http://www.baidu.com/a.jpg"))
}

func TestParseOssKey(t *testing.T) {
	t.Log(ParseOssKey("http://www.baidu.com/a.jpg"))
	t.Log(ParseOssKey("/a.jpg"))
	t.Log(ParseOssKey("a.jpg"))
}

func TestCountStrByGBK(t *testing.T) {
	t.Log(CountStrByGBK("aa啊a12啊a"))
}

func TestTrimHtml(t *testing.T) {
	//str := "<p style=\"aa\">abcd <a href=\"http://www.baiud.com\" target=\"_blank\">http://www.baiud.com</a></p><div class=\"media-wrap image-wrap\"><img src='https://attachments.bjx.cloud/org_1004/feedback_resource/2020/7/22/608c011b2d80402ab74d0ecd3c3f4df01595402680773.jpg' style=\"size:100px\"/></div><p><br/>cd <strong>ee <em>ff  <u>dd <span style=\"color:#f39c12\">ffaa</span></u></em></strong></p>"
	str := "<p>abcd <a href=\"http://www.baiud.com\" target=\"_blank\">http://www.baiud.com</a></p><div class=\"media-wrap image-wrap\"><img class=\"media-wrap image-wrap\" src=\"https://attachments.bjx.cloud/org_1004/feedback_resource/2020/7/22/608c011b2d80402ab74d0ecd3c3f4df01595402680773.jpg\"/></div><p><br/>cd <strong>ee <em>ff  <u>dd <span style=\"color:#f39c12\">ffaa</span></u></em></strong></p><p><strong><em><u><span style=\"color:#f39c12\">bc</span></u></em></strong></p><p><strong><em><u><span style=\"color:#f39c12\">aa</span></u></em></strong></p>"
	t.Log(TrimHtml(str))
}
