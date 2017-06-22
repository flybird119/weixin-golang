package wanxiang

import "testing"

func TestWanxiang(t *testing.T) {
	book, err := GetBookInfo("9780596001193")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(book)
}
