package wanxiang

import "testing"

func TestWanxiang(t *testing.T) {
	_, err := GetBookInfo("9787122087935")
	if err != nil {
		t.Error(err)
	}
}
