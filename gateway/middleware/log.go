package middleware

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"goushuyun/errs"
	"goushuyun/misc"
	"goushuyun/misc/hack"

	"github.com/pborman/uuid"
	"github.com/urfave/negroni"
	"github.com/wothing/log"
)

func LogMiddleware() negroni.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		start := time.Now()
		tid := uuid.New()
		w.Header().Set("X-Request-ID", tid)

		version := r.Header.Get("X-App-Version")

		body, err := ioutil.ReadAll(r.Body)

		if len(body) > 1000 {
			log.Debugf("%s", body[:1000])
		} else {
			log.Debugf("%s", body)
		}

		r.Body.Close()
		if err != nil {
			log.Terrorf(tid, "error on reading rwquest body, from=%v, method=%v, remote=%v, agent=%v, version=%s", r.RequestURI, r.Method, r.RemoteAddr, r.UserAgent(), version)
			misc.RespondMessage(w, r, misc.NewErrResult(errs.ErrRequestFormat, "error on reading request body"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "tid", tid)
		ctx = context.WithValue(ctx, "body", body)
		r = r.WithContext(ctx)

		if realIp := r.Header.Get("X-Real-IP"); realIp != "" {
			r.RemoteAddr = realIp
		} else {
			if i := strings.LastIndex(r.RemoteAddr, ":"); i > 0 {
				r.RemoteAddr = r.RemoteAddr[:i]
			}
		}

		bodyFormat := replaceHttpReqPassword(hack.String(body))
		log.Tinfof(tid, "started handling request, from=%v, method=%v, remote=%v, agent=%v, version=%s, body=%v", r.RequestURI, r.Method, r.RemoteAddr, r.UserAgent(), version, bodyFormat)
		next(w, r)
		log.Tinfof(tid, "completed handling request, status=%v, took=%v", w.(negroni.ResponseWriter).Status(), time.Since(start))
	}
}

func replaceHttpReqPassword(s string) string {
	if len(s) > 10000 {
		s = s[:10000]
	}

	match := `"password":"`
	if i := strings.Index(s, match); i != -1 {
		if j := strings.Index(s[i+12:], `"`); j != -1 {
			return s[:i+12] + "******" + s[i+12+j:]
		}
	}
	return s
}
