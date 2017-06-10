package wanxiang

import "testing"

func TestWanxiang(t *testing.T) {
	book, err := GetBookInfo("9787541695796")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(book)
}
