package main

import (
	"fmt"
	"goushuyun/db"
	"goushuyun/gateway/appway/router"
	m "goushuyun/gateway/middleware"
	"runtime"

	"github.com/urfave/negroni"
	"github.com/wothing/log"
)

const (
	port    = 10014
	svcName = "appway"
)

var serviceNames = []string{
	"weixin",
	"users",
	"books",
	"orders",
	"address",
	"admin",
	"payment",
	"activity",
	"mediastore",
}

func main() {
	defer log.Info("Appway stoped !")
	runtime.GOMAXPROCS(runtime.NumCPU())

	micro := db.NewMicro(svcName, port)
	micro.ReferServices(serviceNames...)

	n := negroni.New()
	n.Use(m.RecoveryMiddleware())
	n.Use(m.LogMiddleware())
	n.Use(m.JWTMiddleware())
	n.UseHandler(router.SetRouterV1())

	networkAddr := fmt.Sprintf("0.0.0.0:%d", db.GetPort(port))
	log.Infof("Interway servering in %s", networkAddr)
	n.Run(networkAddr)
}
