package github

import (
	"encoding/hex"
	"reflect"
	"testing"
)

func TestComputeHmac256(t *testing.T) {

	result := ComputeHmac256([]byte("source string"), []byte("secret"))
	expected, _ := hex.DecodeString("51da5fbfc38fbcc0266b8b6a81bcd01fdc7fdc01dea0bcb983a77a8a004ca66a")

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Computed hash is incorrect")
	}
}
