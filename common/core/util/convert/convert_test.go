package convert

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

type test struct {
	Name string `json:"name"`
}

func TestObjectToMap(t *testing.T) {
	te := test{"testname"}
	data := ObjectToMap(te)
	assert.Equal(t, data["name"], "testname")
}
