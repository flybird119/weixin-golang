package service

import "testing"

func TestQrcode(t *testing.T) {
	t.Log("==========================")
	t.Log("==========================")

	url, err := GenQrcode("Wang", 188, 188)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(url)
}
