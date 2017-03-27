package misc

import (
	"testing"
	"time"

	"github.com/jinzhu/now"
)

func TestFazzyQuery(t *testing.T) {
	t.Log(FazzyQuery("hell"))
}

func TestMd5(t *testing.T) {
	t.Log("The md5 value is ", Md5String("12345678"))
}

func TestDemo(t *testing.T) {
	t.Log(now.BeginningOfDay().Add(-24 * time.Hour))

	t.Log(t)

}
