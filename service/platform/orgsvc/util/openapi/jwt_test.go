package openapi

import "testing"

func TestEncodeAccessToken(t *testing.T) {
	token, err := EncodeAccessToken("abc", "efg")
	t.Log(err)
	t.Log(token)
}

