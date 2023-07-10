package jsonx

import (
	"testing"
	"time"
)

type TA struct {
	D time.Time
	N int64
	M map[int64]string
}

func TestFromJson(t *testing.T) {

}

func TestToJson(t *testing.T) {
	ta := TA{
		D: time.Now(),
		N: 555,
		M: map[int64]string{12: "wad"},
	}
	str, err := ToJson(ta)
	t.Log(err)
	t.Log(str)

	ta1 := &TA{}
	err = FromJson(str, ta1)
	t.Log(err)
	t.Log(ta.N)

	ta1 = &TA{}
	err = FromJson("{\"D\":\"2020-08-27 17:00:47\",\"N\":\"555\",\"M\":{\"12\":\"wad\"}}", ta1)
	t.Log(err)
	t.Log(ta1.N)
	t.Log(ta1.N)
}
