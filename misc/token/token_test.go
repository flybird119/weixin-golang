/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/07/05 23:17
 */

package token

import (
	"fmt"
	"log"
	"reflect"
	"testing"

	"github.com/goushuyun/weixin-golang/seller/role"
)

var tokenStr string
var session string

func TestZeroValue(t *testing.T) {
	x := struct {
		Foo string
		Bar int
	}{"foo", 2}

	v := reflect.ValueOf(x)
	t.Log(">>>>>>>>>>>>>>>>>>>>>")

	t.Logf("%+v\n", v)
	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	fmt.Printf("%+v\n", values)

	// val := ""
	// t.Log(IsZeroOfUnderlyingType(val))
}

func IsZeroOfUnderlyingType(x interface{}) bool {
	log.Println("------", reflect.TypeOf(x), "----------")

	return x == reflect.Zero(reflect.TypeOf(x)).Interface()
}

func TestSignAppToken(t *testing.T) {
	tokenStr = SignUserToken(AppToken, "0000000001", "170424000006")
	t.Log(tokenStr)

	claims, err := Check(tokenStr)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", claims)
}

func TestSign(t *testing.T) {
	tokenStr = SignSellerToken(InterToken, "00000004", "18817953402", "170411000002", role.InterAdmin)
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
	tokenStr := `eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJtb2IiOiIxODgxODAwMDMwNSIsInNlcyI6ImljZXoiLCJzY3AiOjIsInJvbCI6MSwiaWF0IjoxNDkxMzg0MDkyLCJzZWxsZXJfaWQiOiIwMDAwMDAwMSIsImlzcyI6IjEuMCJ9.Nbi8J-emf6ZLx_6y-bJ830Y4YphXO4wJaGU4ssmTK3bsLDMV2ULDcUzl7offO68asUFWFrAF35aazO1YJPqJdktawAUKZtJQXLwf6l0_re3_AOjWa092E0xdRFyFatrAoT53GUjo7UNrOiOOa-KqhfgD9sowP3W4DwC3ehHZn87sIuqJITnuuYHRTkPiSM1bOOwGl57qzjwi_7bhemlgoomy8NdinDS-cJZkS_K5br_WyBbgk4ndD7NzK9TRbXt7nSzm3yXtLeQ-fZ8xUT2kd-Rgv_K1GUy71u77bZLLSYnRZh_rCjOekktHpxSHeg6Fkc7naxc1EmUVeEtkmBRZSg`

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
