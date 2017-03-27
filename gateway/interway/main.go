package main

import (
	"fmt"
	"goushuyun/db"
	"goushuyun/gateway/interway/router"
	m "goushuyun/gateway/middleware"
	"runtime"

	"github.com/urfave/negroni"
	"github.com/wothing/log"
)

const (
	port    = 10013
	svcName = "interway"
)

var serviceNames = []string{
	"admin",
	"store",
	"mediastore",
	"books",
	"orders",
	"activity",
	"statistic",
}

func main() {
	defer log.Info("Interway stoped !")
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
