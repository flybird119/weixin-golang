package main

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/urfave/negroni"
	"github.com/wothing/log"

	"github.com/goushuyun/weixin-golang/db"
	"github.com/goushuyun/weixin-golang/gateway/adminway/router"
	m "github.com/goushuyun/weixin-golang/gateway/middleware"
)

const (
	svcName = "bc_adminway"
	port    = 8870
)

var serviceNames = []string{
	"bc_master",
	"bc_store",
}

func main() {
	defer log.Infof("%s stopped, bye bye !", svcName)
	runtime.GOMAXPROCS(runtime.NumCPU())

	micro := db.NewMicro(svcName, port)
	micro.ReferServices(serviceNames...)

	n := negroni.New()
	n.Use(m.RecoveryMiddleware())
	n.Use(m.LogMiddleware())
	n.Use(m.JWTMiddleware())
	n.UseHandler(router.SetRouterV1())
	networkAddr := fmt.Sprintf("0.0.0.0:%d", db.GetPort(port))
	log.Infof("%s servering on %s", svcName, networkAddr)
	log.Fatal(http.ListenAndServe(networkAddr, n))
}
