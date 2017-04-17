package component

import "testing"

func TestAccessToken(t *testing.T) {
	access_token, err := PreAuthCode()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(access_token)
}
