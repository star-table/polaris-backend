package mvc

import (
	"context"
	"fmt"
	"github.com/star-table/common/core/consts"
	"github.com/star-table/common/core/types"
	"github.com/star-table/common/core/util/json"
	"github.com/star-table/polaris-backend/common/test"
	"github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
	"upper.io/db.v3"
)

type TestInjectStruct struct {
	A int
	B int8
	C int16
	D int32
	F int64
	Ua uint
	Ub uint8
	Uc uint16
	Ud uint32
	Uf uint64

	G float32
	H float64

	I bool
	Ui bool
	J time.Time
	K types.Time

	Z TestBody


	Ap *int
	Bp *int8
	Cp *int16
	Dp *int32
	Fp *int64
	Uap *uint
	Ubp *uint8
	Ucp *uint16
	Udp *uint32
	Ufp *uint64

	Gp *float32
	Hp *float64

	Ip *bool
	Uip *bool
	Jp *time.Time
	Kp *types.Time

	Zp *TestBody
}

type TestBody struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

func Test_Convert(t *testing.T){

	convey.Convey("Test_Convert", t, test.StartUp(func(ctx context.Context) {


		testBody := TestBody{Name:"nico", Value:"abc",}

		reqInfo := RequestInfo{
			Parameters: map[string]string{
				"a": "1",
				"b": "1",
				"c": "1",
				"d": "1",
				"f": "1",

				"ua": "1",
				"ub": "1",
				"uc": "1",
				"ud": "1",
				"uf": "1",

				"g":"12.3",
				"h":"34.2",

				"i":"true",
				"ui":"false",

				"j":"2008-09-01 11:11:11",
				"k":"2008-09-01 11:11:11",

				"bp": "1",
				"cp": "1",
				"dp": "1",
				"fp": "1",

				"uap": "1",
				"ubp": "1",
				"ucp": "1",
				"udp": "1",
				"ufp": "1",

				"gp":"12.3",
				"hp":"34.2",

				"ip":"true",
				"uip":"false",

				"jp":"2008-09-01 11:11:11",
				"kp":"2008-09-01 11:11:11",
			},
			Body: json.ToJsonIgnoreError(testBody),
		}

		v, e := InjectFunc(testFunc, reqInfo)
		t.Log(v, e)

		req := v[0].Interface().(TestInjectStruct)

		fmt.Println(json.ToJsonIgnoreError(req))


		assert.Equal(t, req.A, int(1))
		assert.Equal(t, req.B, int8(1))
		assert.Equal(t, req.C, int16(1))
		assert.Equal(t, req.D, int32(1))
		assert.Equal(t, req.F, int64(1))

		assert.Equal(t, req.Ua, uint(1))
		assert.Equal(t, req.Ub, uint8(1))
		assert.Equal(t, req.Uc, uint16(1))
		assert.Equal(t, req.Ud, uint32(1))
		assert.Equal(t, req.Uf, uint64(1))

		assert.Equal(t, req.G, float32(12.3))
		assert.Equal(t, req.H, float64(34.2))
		assert.Equal(t, req.I, true)
		assert.Equal(t, req.Ui, false)

		assert.Equal(t, req.Ap, (*int)(nil))

		time, _ := time.Parse(consts.AppTimeFormat, "2008-09-01 11:11:11")
		assert.Equal(t, req.J, time)

		myTime := types.Time(time)
		assert.Equal(t, req.K, myTime)

		assert.Equal(t, req.Z, testBody)

	}))

}

func Test_ConvertNil(t *testing.T){
	testBody := TestBody{Name:"nico", Value:"abc",}

	reqInfo := RequestInfo{
		Parameters: map[string]string{
			"b": "1",
		},
		Body: json.ToJsonIgnoreError(testBody),
	}

	v, e := InjectFunc(testFunc, reqInfo)
	t.Log(v, e)

	assert.NotEqual(t, e, nil)
}

func testFunc(req TestInjectStruct) TestInjectStruct{

	return req
}


func Test_Cond(t *testing.T){
	cond := db.Cond{}
	ty := reflect.TypeOf(cond)
	t.Log(ty.Kind())
}