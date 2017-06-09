package wanxiang

import (
	"path/filepath"
	"testing"
)

func TestWanxiang(t *testing.T) {
	book, err := GetBookInfo("9787541695796")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(book)
}

func TestPath(t *testing.T) {

	t.Log(filepath.Ext("http://api.jisuapi.com/isbn/upload/20170104/212742_44734.jpg"))
}
