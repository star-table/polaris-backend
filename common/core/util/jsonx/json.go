package jsonx

import (
	"strconv"
	"time"
	"unsafe"

	"github.com/star-table/polaris-backend/common/core/consts"
	jsoniter "github.com/json-iterator/go"
	"github.com/modern-go/reflect2"
)

type int64codec struct {
}

type TimeDecoder struct {
}
type TimeEncoder struct {
	precision time.Duration
}

var loc, _ = time.LoadLocation("Local")

func (td *TimeDecoder) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	str := iter.ReadString()

	mayBlank, _ := time.Parse(consts.AppTimeFormat, str)
	now, err := time.ParseInLocation(consts.AppTimeFormat, str, loc)

	if err != nil {
		*((*time.Time)(ptr)) = time.Unix(0, 0)
	} else if mayBlank.IsZero() {
		*((*time.Time)(ptr)) = mayBlank
	} else {
		*((*time.Time)(ptr)) = now
	}
}

func (codec *TimeEncoder) IsEmpty(ptr unsafe.Pointer) bool {
	ts := *((*time.Time)(ptr))
	return ts.IsZero()
}

func (codec *TimeEncoder) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	ts := *((*time.Time)(ptr))
	if !ts.IsZero() {
		timestamp := ts.Unix()
		tm := time.Unix(timestamp, 0)
		format := tm.Format(consts.AppTimeFormat)
		stream.WriteString(format)
	} else {
		mayBlank, _ := time.Parse(consts.AppTimeFormat, consts.BlankString)
		stream.WriteString(mayBlank.Format(consts.AppTimeFormat))
	}

}

func (codec *int64codec) IsEmpty(ptr unsafe.Pointer) bool {
	return *((*int64)(ptr)) == 0
}

func (codec *int64codec) Encode(ptr unsafe.Pointer, stream *jsoniter.Stream) {
	stream.WriteString(strconv.FormatInt(*((*int64)(ptr)), 10))
}

func (codec *int64codec) Decode(ptr unsafe.Pointer, iter *jsoniter.Iterator) {
	data := iter.Read()
	if in, ok := data.(string); ok {
		v, err := strconv.ParseInt(in, 10, 64)
		if err != nil {
			*((*int64)(ptr)) = 0
		} else {
			*((*int64)(ptr)) = v
		}
	} else if in, ok := data.(int64); ok {
		*((*int64)(ptr)) = in
	} else {
		v, err := strconv.ParseInt(iter.ReadString(), 10, 64)
		if err != nil {
			*((*int64)(ptr)) = 0
		} else {
			*((*int64)(ptr)) = v
		}
	}
}

type int64Extension struct {
}

func (extension *int64Extension) CreateEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	if typ == reflect2.TypeOfPtr((*int64)(nil)).Elem() {
		return &int64codec{}
	}
	return nil
}

func (extension *int64Extension) CreateDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	if typ == reflect2.TypeOfPtr((*int64)(nil)).Elem() {
		return &int64codec{}
	}
	return nil
}

// UpdateStructDescriptor No-op
func (extension *int64Extension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
}

// CreateMapKeyDecoder No-op
func (extension *int64Extension) CreateMapKeyDecoder(typ reflect2.Type) jsoniter.ValDecoder {
	if typ == reflect2.TypeOfPtr((*int64)(nil)).Elem() {
		return &int64codec{}
	}
	return nil
}

// CreateMapKeyEncoder No-op
func (extension *int64Extension) CreateMapKeyEncoder(typ reflect2.Type) jsoniter.ValEncoder {
	if typ == reflect2.TypeOfPtr((*int64)(nil)).Elem() {
		return &int64codec{}
	}
	return nil
}

// DecorateDecoder No-op
func (extension *int64Extension) DecorateDecoder(typ reflect2.Type, decoder jsoniter.ValDecoder) jsoniter.ValDecoder {
	return decoder
}

// DecorateEncoder No-op
func (extension *int64Extension) DecorateEncoder(typ reflect2.Type, encoder jsoniter.ValEncoder) jsoniter.ValEncoder {
	return encoder
}

var json = jsoniter.Config{
	EscapeHTML:                    false,
	MarshalFloatWith6Digits:       true, // will lose precession
	ObjectFieldMustBeSimpleString: true, // do not unescape object field
}.Froze()

func JSON() jsoniter.API {
	return json
}

func init() {
	jsoniter.RegisterTypeEncoder("time.Time", &TimeEncoder{})
	jsoniter.RegisterTypeDecoder("time.Time", &TimeDecoder{})
	json.RegisterExtension(&int64Extension{})
}

func ToJson(obj interface{}) (string, error) {
	bs, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func ToJsonIgnoreError(obj interface{}) string {
	if obj == nil {
		return ""
	}
	bs, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return string(bs)
}

func FromJson(jsonStr string, obj interface{}) error {
	return json.Unmarshal([]byte(jsonStr), obj)
}

func Copy(src, dst interface{}) error {
	bs, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bs, dst)
	if err != nil {
		return err
	}
	return nil
}
