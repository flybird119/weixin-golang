package douban

import (
	"strings"
	"testing"

	"github.com/arstd/log"
)

func TestGetBookInfo(t *testing.T) {
	_, err := GetBookInfo("9787513570343")

	if err != nil && strings.Index(err.Error(), "404") == -1 {
		log.Error(err)

	} else {
		log.Debug("not found ...")
	}

}
