/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/07/05 23:17
 */

package token

import (
	"testing"

	"github.com/goushuyun/weixin-golang/seller/role"
)

var tokenStr string
var session string

func TestSign(t *testing.T) {
	tokenStr = SignSellerToken(InterToken, "2487583", "18817953402", "", role.InterAdmin)
	t.Log(">>>>>>>>>>>>>>>token>>>>>>>>>>>>>>")
	t.Logf("%s\n", tokenStr)
	t.Log(">>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>")

	claims, err := Check(tokenStr)

	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", claims)
}

func TestCheck(t *testing.T) {
	tokenStr := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJtb2IiOiIxODgxNzk1MzQwMiIsInNlcyI6Im1wNDAiLCJzY3AiOjIsInJvbCI6MSwiaWF0IjoxNDkxMzE0ODM1LCJzZWxsZXJfaWQiOiIwMDAwMDAwNCIsImlzcyI6IjEuMCJ9.uEbyqA8MpMtUWtm2qVdN8cWehhNyRvYlPQLB4eDaXkH5MJlsBdVDqyEZmy4B38tTiEbLpSeoF0U6zdt-7algwU_wfiNI5x47IqTb3XW2d8Z3WuNh1l7E_80vWi-Q_QkNDjbApY4_w7eRL9Se1ZgfOvCKnTHT4V2THEVlOAWT4Sh_p1gZyclQqmtR2Pq6OLMKM9SkROZeQNgANi9UGeeyV6XIDajeoRh4F8lDwMIiy1ZgO2VfOi8IU-9rduSIGrhOBYX_dK3iNxNjiSN_eAylbEBPCASag0CDGttknSOSPndMB-nHkA2nqtcP857mfM-AXNm2xESN6BOKyVLu89gJUw`

	c, err := Check(tokenStr)

	t.Logf("%+v\n", c)

	if err != nil {
		t.Error("check failed")
	}

	if c.VerifyIsExpired() != false {
		t.Error("should not expire")
	}

	if c.VerifyCanRefresh() != true {
		t.Error("can be refresh")
	}

	temp := tokenStr + "xxxx"
	_, err = Check(temp)
	if err == nil {
		t.Error("should be illegal error")
	}
}

func TestRefresh(t *testing.T) {
	//before := tokenStr
	//<-time.After(time.Second)
	//err := Refresh(&tokenStr)
	//if err != nil {
	//	t.Error("refresh failed")
	//}
	//if before == tokenStr {
	//	t.Error("refresh failed")
	//}
}
