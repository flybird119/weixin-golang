/**
 * Copyright 2015-2016, Wothing Co., Ltd.
 * All rights reserved.
 *
 * Created by Elvizlai on 2016/07/05 23:16
 */

package token

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/elgs/gostrgen"

	ro "goushuyun/admin/role"
)

func sign(c Claims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, c)
	tokenStr, err := token.SignedString(getPrivateKeyFromPEM())
	if err != nil {
		panic(err)
	}
	return tokenStr
}

func Sign(from TokenType, uid, mob string, role int64) (tokenStr string, session string) {
	var err error
	session, err = gostrgen.RandGen(4, gostrgen.Lower|gostrgen.Digit, "", "")
	if err != nil {
		panic(err)
	}

	c := Claims{
		UserId:   uid,
		Mobile:   mob,
		Session:  session,
		Scope:    from,
		Role:     ro.Role(role),
		IssuedAt: time.Now().Unix(),
		Version:  currentVersion,
	}

	tokenStr = sign(c)
	return tokenStr, session
}

func SignWithSession(from TokenType, uid, mob, hosp, sess string, role int64) string {
	c := Claims{
		UserId:   uid,
		Mobile:   mob,
		Session:  sess,
		Scope:    from,
		Role:     ro.Role(role),
		IssuedAt: time.Now().Unix(),
		Version:  currentVersion,
	}

	return sign(c)
}

func Refresh(c *Claims) string {
	c.IssuedAt = time.Now().Unix()
	return sign(*c)
}

// Check checks only if the token is illegal, logical check should use c.VerifyXXX
func Check(tokenStr string) (*Claims, error) {
	c := &Claims{}
	_, err := jwt.ParseWithClaims(tokenStr, c, func(token *jwt.Token) (interface{}, error) {
		return jwt.ParseRSAPublicKeyFromPEM(publicKey)
	})

	if err != nil {
		return nil, err
	}

	return c, nil
}

// Get using gctx to get value using key 'claims' from r
// if not, return nil
func Get(r *http.Request) *Claims {
	if claims, ok := r.Context().Value("claims").(*Claims); ok {
		return claims
	} else {
		return nil
	}
}

func ReSign(c *Claims, role ro.Role) string {
	c.Role = role
	return sign(*c)
}

func ReSignAll(c *Claims, userId string, mobile string, hospId string, role ro.Role) string {
	c.UserId = userId
	c.Mobile = mobile
	c.Role = role
	return sign(*c)
}
