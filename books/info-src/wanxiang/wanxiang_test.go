package wanxiang

import (
	"path/filepath"
	"testing"
)

func TestWanxiang(t *testing.T) {
	_, err := GetBookInfo("9787122087935")
	if err != nil {
		t.Error(err)
	}
}

func TestPath(t *testing.T) {

	t.Log(filepath.Ext("http://api.jisuapi.com/isbn/upload/20170104/212742_44734.jpg"))
}
