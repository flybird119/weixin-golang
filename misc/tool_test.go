package misc

import (
	"fmt"
	"testing"
)

func TestGenCheckCode(t *testing.T) {
	code := GenCheckCode(4, KC_RAND_KIND_NUM)
	fmt.Println("====>", code)
}
