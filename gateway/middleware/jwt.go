package middleware

import (
	"context"
	"goushuyun/db"
	"net/http"
	"strings"

	"goushuyun/errs"

	"goushuyun/misc"
	"goushuyun/misc/token"

	"github.com/garyburd/redigo/redis"
	"github.com/urfave/negroni"
	"github.com/wothing/log"
)

var whiteList = map[string]bool{
	"v1/user/refresh_token": true,
}

func JWTMiddleware() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			authHeader = r.URL.Query().Get("token")
		}

		if authHeader != "" {
			authHeaderParts := strings.Split(authHeader, " ")
			switch len(authHeaderParts) {
			case 1:
				next(w, r)
				return
			case 2:
				if strings.ToLower(authHeaderParts[0]) != "bearer" || authHeaderParts[1] == "" {
					misc.RespondMessage(w, r, errs.NewError(errs.ErrTokenFormat, "token format error"))
					return
				}
			default:
				misc.RespondMessage(w, r, errs.NewError(errs.ErrTokenFormat, "token length error"))
				return
			}

			c, err := token.Check(authHeaderParts[1])
			if err != nil {
				log.Warn("authHeader: ", authHeader, "err: ", err)
				misc.RespondMessage(w, r, errs.NewError(errs.ErrTokenFormat, "token illegal"))
				return
			}

			//token version check
			if !c.VerifyVersion() {
				misc.RespondMessage(w, r, errs.NewError(errs.ErrTokenRefreshExpired, "need reload"))
				return
			}

			if c.VerifyIsExpired() {
				if c.VerifyCanRefresh() {
					if !whiteList[r.RequestURI] {
						w.Header().Add("X-JWT-Token", token.Refresh(c))
					}
				} else {
					misc.RespondMessage(w, r, errs.NewError(errs.ErrTokenRefreshExpired, "need relogin"))
					return
				}
			}
			// claims is ptr
			r = r.WithContext(context.WithValue(r.Context(), "claims", c))

			log.Debug(misc.SuperPrint(*c))

		}
		next(w, r)
	}
}

func SessionMiddleware() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if c := token.Get(r); c != nil {
			if c.Session == "H5" {
				// do nothing
			} else {
				rc := db.GetRedisConn()
				s, err := redis.String(rc.Do("get", "s:"+c.UserId))

				rc.Close()
				if err == nil && s != c.Session {
					misc.RespondMessage(w, r, errs.NewError(errs.ErrSessionExpired, "please re-signin"))
					return
				}
			}
		}
		next(w, r)
	}
}
