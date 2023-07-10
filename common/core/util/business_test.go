package util

import (
	"fmt"
	"github.com/star-table/common/core/util/json"
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestRenderCommentContentToMarkDown(t *testing.T) {
	t1 := RenderCommentContentToMarkDown("@#[nico:123]&$fjafjj", false)
	assert.Equal(t, t1, "<at id=123></at>fjafjj")

	t2 := RenderCommentContentToMarkDown("@#[nico:123]&$fjaf@#[nico:123]&$jj", false)
	assert.Equal(t, t2, "<at id=123></at>fjaf<at id=123></at>jj")

	t3 := RenderCommentContentToMarkDown("@#[nico:123]&$fjafjj@#[nico:123]&$", false)
	assert.Equal(t, t3, "<at id=123></at>fjafjj<at id=123></at>")

	t4 := RenderCommentContentToMarkDown("@#[nico:123]&$fja@#[nico:123]&$fjj@#[nico:123]&$", false)
	assert.Equal(t, t4, "<at id=123></at>fja<at id=123></at>fjj<at id=123></at>")

	t5 := RenderCommentContentToMarkDown("fjafjj", false)
	assert.Equal(t, t5, "fjafjj")

	t6 := RenderCommentContentToMarkDown("@#[樊宇:ou_87f1b2210acad10a90cc3690802626d7]&$hello这是一条艾特消息321", false)
	assert.Equal(t, t6, "<at id=ou_87f1b2210acad10a90cc3690802626d7></at>hello这是一条艾特消息321")

	t7 := RenderCommentContentToMarkDown("@#[nico:123]&$fjafjj", true)
	assert.Equal(t, t7, "fjafjj")

	t8 := RenderCommentContentToMarkDown("@#[nico:123]&$fjafjj@#[nico:123]&$abc@#[nico:123]&$ 换行", true)
	assert.Equal(t, t8, "fjafjjabc 换行")

	t9 := RenderCommentContentToOrigin("@#[樊宇:ou_87f1b2210acad10a90cc3690802626d7]&$hello这是一条艾特@#[Nico:ou_87f1b2210acad10a90cc3690802626d7]&$消息321")
	fmt.Println(t9)
	assert.Equal(t, t9, "@樊宇 hello这是一条艾特@Nico 消息321")
}

func TestGetCommentAtUserIds(t *testing.T) {
	t.Log(json.ToJsonIgnoreError(GetCommentAtUserIds("@#[nico:123]&$fjafjj@#[nico:1234]&$abc@#[nico:1255]&$ 换行")))
}

func TestRenderCommentContentToMarkDownWithOpenIdMap(t *testing.T) {
	t1 := RenderCommentContentToMarkDown("@#[nico:123]&$fjafjj", false)
	t.Log(t1)

	openIdMap := map[string]string{}
	openIdMap["123"] = "456"
	t1 = RenderCommentContentToMarkDownWithOpenIdMap("@#[nico:123]&$fjafjj", false, openIdMap)
	t.Log(t1)

	t1 = RenderCommentContentToMarkDownWithOpenIdMap("@#[nico:123]&$fjafjj@#[樊宇:ou_87f1b2210acad10a90cc3690802626d7]&$hello这是一条艾特", false, openIdMap)
	t.Log(t1)
}

func TestJointUrl(t *testing.T) {

	u1 := JointUrl("https://abc.com", "haha")
	assert.Equal(t, u1, "https://abc.com/haha")

	u2 := JointUrl("https://abc.com/", "haha")
	assert.Equal(t, u2, "https://abc.com/haha")

	u3 := JointUrl("https://abc.com", "/haha")
	assert.Equal(t, u3, "https://abc.com/haha")

	u4 := JointUrl("https://abc.com/", "/haha")
	assert.Equal(t, u4, "https://abc.com/haha")

	t.Log(u1, u2, u3, u4)
}

func TestModifyFileName(t *testing.T) {
	f1 := ModifyFileName("abc.jpg", "_80")
	t.Log(f1)
	f1 = ModifyFileName("https://12312/abc.jpg", "_80")
	t.Log(f1)
	f1 = ModifyFileName("https://a/1.1/abc.jpg", "_80")
	t.Log(f1)
}

func TestGetOssKeyInfo(t *testing.T) {

	t.Log(json.ToJsonIgnoreError(GetOssKeyInfo("org_1325/project_1525/issue_7083/resource/2019/12/12/376b9ff64f184b958041c4f6b7eb31021576141131899.png")))
	t.Log(json.ToJsonIgnoreError(GetOssKeyInfo("org_/project_1525/issue_7083/resource/2019/12/12/376b9ff64f184b958041c4f6b7eb31021576141131899.png")))
	t.Log(json.ToJsonIgnoreError(GetOssKeyInfo("org_1325/project_1525/resource/2019/12/12/376b9ff64f184b958041c4f6b7eb31021576141131899.png")))
	t.Log(json.ToJsonIgnoreError(GetOssKeyInfo("org_1325/issue_7083/resource/2019/12/12/376b9ff64f184b958041c4f6b7eb31021576141131899.png")))
	t.Log(json.ToJsonIgnoreError(GetOssKeyInfo("org_1325/issue_erf/resource/2019/12/12/376b9ff64f184b958041c4f6b7eb31021576141131899.png")))

}

func TestUnicodeEmojiCodeFilter(t *testing.T) {
	t.Log("testtest😄")
	t.Log(UnicodeEmojiCodeFilter("testtest😄"))
}

func TestParseFileSuffix(t *testing.T) {

	suffix := ParseFileSuffix("abc")
	assert.Equal(t, suffix, "")
	suffix = ParseFileSuffix("abc.txt")
	assert.Equal(t, suffix, "txt")
	suffix = ParseFileSuffix("abc.zip")
	assert.Equal(t, suffix, "zip")

}

func TestRoleOperationCodesMatch(t *testing.T) {

	t.Log(RoleOperationCodesMatch("Modify", "(View)|(Modify)|(Bind)|(Unbind)"))
	t.Log(RoleOperationCodesMatch("Modify", "(View)|(ModifyStatus)|(Bind)|(Unbind)"))
	t.Log(RoleOperationCodesMatch("Modify", "(View)|(Bind)|(Unbind)"))
	t.Log(RoleOperationCodesMatch("Modify", "Modify"))
	t.Log(RoleOperationCodesMatch("Modify", "1"))

}
